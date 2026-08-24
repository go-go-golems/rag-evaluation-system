// Package researchctladapter maps canonical RAG v2 values to researchctl's
// public, domain-neutral laboratory contracts.
package researchctladapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/researchctl/pkg/lab"
)

type InputDocument struct {
	Inputs map[string]InputReference `json:"inputs"`
}
type InputReference struct {
	Role          string            `json:"role"`
	Kind          string            `json:"kind"`
	ID            string            `json:"id,omitempty"`
	URI           string            `json:"uri,omitempty"`
	Digest        string            `json:"digest,omitempty"`
	SizeBytes     *int64            `json:"sizeBytes,omitempty"`
	SchemaVersion string            `json:"schemaVersion,omitempty"`
	MediaType     string            `json:"mediaType,omitempty"`
	Catalog       *lab.CatalogRef   `json:"catalog,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}
type ResolvedInput struct {
	Reference lab.ArtifactRef
	Binding   ragcontract.ArtifactBinding
}
type ResolvedInputs struct{ ByRole map[string]ResolvedInput }
type CatalogResolver interface {
	Resolve(context.Context, InputReference, string) (ResolvedInput, error)
}

func LoadInputs(path string) (InputDocument, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return InputDocument{}, "", err
	}
	defer func() { _ = file.Close() }()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var document InputDocument
	if err := decoder.Decode(&document); err != nil {
		return InputDocument{}, "", fmt.Errorf("RAG_INPUTS_DECODE: %w", err)
	}
	if len(document.Inputs) == 0 {
		return InputDocument{}, "", fmt.Errorf("RAG_INPUTS_EMPTY")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return InputDocument{}, "", err
	}
	return document, filepath.Dir(absolute), nil
}

func ResolveInputs(ctx context.Context, document InputDocument, baseDir, artifactRoot string, catalog CatalogResolver) (ResolvedInputs, error) {
	result := ResolvedInputs{ByRole: map[string]ResolvedInput{}}
	roles := make([]string, 0, len(document.Inputs))
	for role := range document.Inputs {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	for _, role := range roles {
		if err := ctx.Err(); err != nil {
			return ResolvedInputs{}, err
		}
		reference := document.Inputs[role]
		if reference.Role == "" {
			reference.Role = role
		}
		if reference.Role != role {
			return ResolvedInputs{}, fmt.Errorf("RAG_INPUT_ROLE: key %s value %s", role, reference.Role)
		}
		var resolved ResolvedInput
		var err error
		if reference.Catalog != nil {
			if catalog == nil {
				return ResolvedInputs{}, fmt.Errorf("RAG_CATALOG_RESOLVER_REQUIRED: %s", role)
			}
			resolved, err = catalog.Resolve(ctx, reference, artifactRoot)
		} else {
			resolved, err = resolveFile(reference, baseDir, artifactRoot)
		}
		if err != nil {
			return ResolvedInputs{}, err
		}
		result.ByRole[role] = resolved
	}
	return result, nil
}

func resolveFile(reference InputReference, baseDir, artifactRoot string) (ResolvedInput, error) {
	if reference.URI == "" {
		return ResolvedInput{}, fmt.Errorf("RAG_INPUT_URI_REQUIRED: %s", reference.Role)
	}
	source := reference.URI
	if !filepath.IsAbs(source) {
		source = filepath.Join(baseDir, filepath.FromSlash(source))
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return ResolvedInput{}, fmt.Errorf("RAG_INPUT_READ %s: %w", reference.Role, err)
	}
	if reference.Digest != "" && lab.DigestBytes(data) != reference.Digest {
		return ResolvedInput{}, fmt.Errorf("RAG_INPUT_DIGEST: %s", reference.Role)
	}
	return stageEnvelope(reference, data, artifactRoot)
}

func StageEnvelope(reference InputReference, value any, artifactRoot string) (ResolvedInput, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return ResolvedInput{}, err
	}
	return stageEnvelope(reference, data, artifactRoot)
}
func stageEnvelope(reference InputReference, data []byte, artifactRoot string) (ResolvedInput, error) {
	manifestDigest, schema, err := manifestIdentity(reference.Role, data)
	if err != nil {
		return ResolvedInput{}, err
	}
	if reference.SchemaVersion != "" && reference.SchemaVersion != schema {
		return ResolvedInput{}, fmt.Errorf("RAG_INPUT_SCHEMA: %s", reference.Role)
	}
	fileDigest := lab.DigestBytes(data)
	name := strings.TrimPrefix(fileDigest, "sha256:") + ".json"
	uri := filepath.ToSlash(filepath.Join("inputs", reference.Role, name))
	destination, err := lab.PrepareArtifactPath(artifactRoot, uri)
	if err != nil {
		return ResolvedInput{}, err
	}
	// #nosec G703 -- PrepareArtifactPath rejects traversal and symlink escapes under artifactRoot.
	if existing, readErr := os.ReadFile(destination); readErr == nil {
		if !bytes.Equal(existing, data) {
			return ResolvedInput{}, fmt.Errorf("RAG_INPUT_CONFLICT: %s", uri)
		}
	} else if !os.IsNotExist(readErr) {
		return ResolvedInput{}, readErr
	} else if err := os.WriteFile(destination, data, 0o644); err != nil { // #nosec G703 -- destination was validated by PrepareArtifactPath.

		return ResolvedInput{}, err
	}
	size := int64(len(data))
	kind := reference.Kind
	if kind == "" {
		kind = "manifest-envelope"
	}
	media := reference.MediaType
	if media == "" {
		media = "application/json"
	}
	artifact := lab.ArtifactRef{Role: reference.Role, Kind: kind, ID: reference.ID, Digest: fileDigest, SizeBytes: &size, SchemaVersion: schema, URI: uri, MediaType: media, Metadata: reference.Metadata}
	binding := ragcontract.ArtifactBinding{SlotID: reference.Role, Role: reference.Role, Kind: kind, ID: reference.ID, URI: uri, Digest: manifestDigest, SizeBytes: &size, SchemaVersion: schema}
	return ResolvedInput{Reference: artifact, Binding: binding}, nil
}

func manifestIdentity(role string, data []byte) (string, string, error) {
	switch role {
	case "corpus":
		var value ragoperators.CorpusArtifact
		if err := strictJSON(data, &value); err != nil {
			return "", "", err
		}
		if err := ragcontract.ValidateManifestBase(value.Manifest.ManifestBase, ragcontract.CorpusManifestSchema, false); err != nil {
			return "", "", err
		}
		return value.Manifest.Digest, ragcontract.CorpusManifestSchema, nil
	case "evaluation-dataset", "judgments":
		var value ragoperators.EvaluationArtifact
		if err := strictJSON(data, &value); err != nil {
			return "", "", err
		}
		if err := ragcontract.ValidateManifestBase(value.Manifest.ManifestBase, ragcontract.EvaluationManifestSchema, true); err != nil {
			return "", "", err
		}
		return value.Manifest.Digest, ragcontract.EvaluationManifestSchema, nil
	default:
		return "", "", fmt.Errorf("RAG_INPUT_ROLE_UNSUPPORTED: %s", role)
	}
}
func strictJSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("trailing JSON value")
		}
		return err
	}
	return nil
}

func ApplyInputs(study ragcontract.Study, resolved ResolvedInputs) (ragcontract.Study, error) {
	corpus, ok := resolved.ByRole["corpus"]
	if !ok {
		return study, fmt.Errorf("RAG_INPUT_CORPUS_REQUIRED")
	}
	evaluation, ok := resolved.ByRole["evaluation-dataset"]
	if !ok {
		evaluation, ok = resolved.ByRole["judgments"]
	}
	if !ok {
		return study, fmt.Errorf("RAG_INPUT_EVALUATION_REQUIRED")
	}
	for index := range study.Bindings {
		if study.Bindings[index].Role == "corpus" {
			binding := corpus.Binding
			binding.SlotID = study.Bindings[index].SlotID
			study.Bindings[index] = binding
		}
	}
	study.Dataset.ManifestDigest = evaluation.Binding.Digest
	return study, nil
}

func LoadDomainArtifacts(artifactRoot string, resolved ResolvedInputs) (ragoperators.CorpusArtifact, ragoperators.EvaluationArtifact, error) {
	corpusInput, ok := resolved.ByRole["corpus"]
	if !ok {
		return ragoperators.CorpusArtifact{}, ragoperators.EvaluationArtifact{}, fmt.Errorf("RAG_INPUT_CORPUS_REQUIRED")
	}
	evaluationInput, ok := resolved.ByRole["evaluation-dataset"]
	if !ok {
		evaluationInput, ok = resolved.ByRole["judgments"]
	}
	if !ok {
		return ragoperators.CorpusArtifact{}, ragoperators.EvaluationArtifact{}, fmt.Errorf("RAG_INPUT_EVALUATION_REQUIRED")
	}
	var corpus ragoperators.CorpusArtifact
	if err := decodeStagedArtifact(artifactRoot, corpusInput.Reference, &corpus); err != nil {
		return corpus, ragoperators.EvaluationArtifact{}, err
	}
	var evaluation ragoperators.EvaluationArtifact
	if err := decodeStagedArtifact(artifactRoot, evaluationInput.Reference, &evaluation); err != nil {
		return corpus, evaluation, err
	}
	return corpus, evaluation, nil
}

func decodeStagedArtifact(root string, reference lab.ArtifactRef, target any) error {
	path, err := lab.PrepareArtifactPath(root, reference.URI)
	if err != nil {
		return err
	}
	body, err := os.ReadFile(path) // #nosec G703 -- PrepareArtifactPath confines the URI to root.
	if err != nil {
		return err
	}
	if lab.DigestBytes(body) != reference.Digest {
		return fmt.Errorf("RAG_INPUT_DIGEST: %s", reference.Role)
	}
	return strictJSON(body, target)
}

func Expand(study ragcontract.Study, resolved ResolvedInputs) (ragcontract.Study, []ragcontract.ExpandedCell, error) {
	updated, err := ApplyInputs(study, resolved)
	if err != nil {
		return study, nil, err
	}
	cells, err := ragcompiler.ExpandStudy(updated, nil)
	return updated, cells, err
}
