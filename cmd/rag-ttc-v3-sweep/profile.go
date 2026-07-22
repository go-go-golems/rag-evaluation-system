package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/internal/workflowv3ttc"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

type providerAuthority struct {
	profile       string
	execution     ragcontract.PipelineExecution
	options       ragengine.Options
	profileDigest string
	modelDigest   string
	close         func() error
}

type specificationEnvelope struct {
	CanonicalIdentity struct {
		Inputs []struct {
			Role, URI, Digest string
			SizeBytes         int64
		} `json:"inputs"`
		DomainConfig json.RawMessage `json:"domainConfig"`
	} `json:"canonicalIdentity"`
}

func loadProviderAuthority(ctx context.Context, profile, providerConfig, specificationPath string) (*providerAuthority, error) {
	if profile == "fixtures" {
		return &providerAuthority{profile: profile, profileDigest: "sha256:" + strings.Repeat("a", 64), modelDigest: "sha256:" + strings.Repeat("b", 64), close: func() error { return nil }}, nil
	}
	if profile != "real" || providerConfig == "" || specificationPath == "" {
		return nil, fmt.Errorf("real profile requires --provider-config and --specification")
	}
	specification, err := decodeSpecification(specificationPath)
	if err != nil {
		return nil, err
	}
	execution, err := ragcontract.DecodeExecution(strings.NewReader(string(specification.CanonicalIdentity.DomainConfig)))
	if err != nil {
		return nil, err
	}
	set, err := ragproviders.Load(ctx, providerConfig)
	if err != nil {
		return nil, err
	}
	descriptor := set.CapabilityDescriptor()
	profileDigest, err := workflowv3.Digest(descriptor)
	if err != nil {
		_ = set.Close()
		return nil, err
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(execution.Pipeline)
	if err != nil {
		_ = set.Close()
		return nil, err
	}
	var generation struct {
		Model string `json:"model"`
	}
	if err = json.Unmarshal(mapping.CombinedNode.Config, &generation); err != nil {
		_ = set.Close()
		return nil, err
	}
	model, err := set.Manifests.Model(generation.Model)
	if err != nil {
		_ = set.Close()
		return nil, err
	}
	return &providerAuthority{profile: profile, execution: execution, options: set.EngineOptions(), profileDigest: profileDigest, modelDigest: model.ModelDigest, close: set.Close}, nil
}

func (a *providerAuthority) provider(batch int) (*workflowv3ttc.OperatorProvider, error) {
	if a.profile == "fixtures" {
		return fixtureProvider(batch)
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(a.execution.Pipeline)
	if err != nil {
		return nil, err
	}
	var cfg map[string]any
	if err = json.Unmarshal(mapping.CombinedNode.Config, &cfg); err != nil {
		return nil, err
	}
	cfg["batchSize"] = batch
	baseRunes, ok := cfg["maxBatchRunes"].(float64)
	if !ok || baseRunes < 1 {
		return nil, fmt.Errorf("generation maxBatchRunes is required")
	}
	cfg["maxBatchRunes"] = int(baseRunes) * batch
	generationConfig, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	generationNode := mapping.CombinedNode
	generationNode.Config = generationConfig
	cache := ragoperators.NewMemoryCache()
	return workflowv3ttc.NewOperatorProvider(workflowv3ttc.OperatorProviderConfig{GenerationNode: generationNode, EmbeddingNode: mapping.EmbeddingNode, RawRepresentationName: mapping.RawRepresentationName, MaxRepresentationsPerChunk: mapping.MaxRepresentationsPerChunk, ProviderProfileDigest: a.profileDigest, GenerationModelDigest: a.modelDigest, EmbeddingProfileDigest: "sha256:" + strings.Repeat("0", 64), ResolveEnvironment: func(context.Context) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: a.options.Manifests, Schemas: a.options.Schemas, Generator: a.options.Generator, Embedder: a.options.Embedder, Reranker: a.options.Reranker, Cache: cache, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: 1, GenerationSettingsFingerprint: a.options.GenerationSettingsFingerprint}, nil
	}})
}

func loadRealChunks(ctx context.Context, specificationPath, artifactRoot string, limit int) ([]workflowv3ttc.Chunk, error) {
	specification, err := decodeSpecification(specificationPath)
	if err != nil {
		return nil, err
	}
	execution, err := ragcontract.DecodeExecution(strings.NewReader(string(specification.CanonicalIdentity.DomainConfig)))
	if err != nil {
		return nil, err
	}
	var corpusInput *struct {
		Role, URI, Digest string
		SizeBytes         int64
	}
	for i := range specification.CanonicalIdentity.Inputs {
		input := specification.CanonicalIdentity.Inputs[i]
		if input.Role == "corpus" {
			corpusInput = &struct {
				Role, URI, Digest string
				SizeBytes         int64
			}{input.Role, input.URI, input.Digest, input.SizeBytes}
			break
		}
	}
	if corpusInput == nil {
		return nil, fmt.Errorf("specification corpus input is missing")
	}
	path := filepath.Join(artifactRoot, filepath.FromSlash(corpusInput.URI))
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if int64(len(body)) != corpusInput.SizeBytes {
		return nil, fmt.Errorf("corpus artifact size mismatch")
	}
	sum := sha256.Sum256(body)
	digest := fmt.Sprintf("sha256:%x", sum)
	if digest != corpusInput.Digest {
		return nil, fmt.Errorf("corpus artifact digest mismatch")
	}
	var artifact ragoperators.CorpusArtifact
	if err = workflowv3.StrictDecode(body, &artifact); err != nil {
		return nil, err
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(execution.Pipeline)
	if err != nil {
		return nil, err
	}
	inputs, err := ragengine.New(nil).StaticInputs(ctx, execution.Pipeline, artifact.Corpus, ragengine.Options{}, mapping.CombinedNode.ID)
	if err != nil {
		return nil, err
	}
	domainChunks, ok := inputs["chunks"].([]ragoperators.Chunk)
	if !ok {
		return nil, fmt.Errorf("pipeline did not materialize chunks")
	}
	sort.Slice(domainChunks, func(i, j int) bool { return domainChunks[i].Record.ID < domainChunks[j].Record.ID })
	if limit > len(domainChunks) {
		return nil, fmt.Errorf("requested %d chunks but only %d exist", limit, len(domainChunks))
	}
	domainChunks = domainChunks[:limit]
	chunks := make([]workflowv3ttc.Chunk, len(domainChunks))
	for i, chunk := range domainChunks {
		chunks[i] = workflowv3ttc.Chunk{Key: chunk.Record.ID, Chunk: chunk, CitationIDs: []string{chunk.Record.ID}, SourceDigest: artifact.Manifest.Digest}
	}
	return chunks, nil
}

func decodeSpecification(path string) (specificationEnvelope, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return specificationEnvelope{}, err
	}
	var specification specificationEnvelope
	if err = json.Unmarshal(body, &specification); err != nil {
		return specificationEnvelope{}, err
	}
	return specification, nil
}
