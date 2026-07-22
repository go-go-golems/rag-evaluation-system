package ragengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

const preparedCorpusSchemaVersion = "rag-prepared-corpus-manifest/v1"

type PreparedCorpusIdentity struct {
	SchemaVersion                 string `json:"schemaVersion"`
	CorpusDigest                  string `json:"corpusDigest"`
	PipelineDigest                string `json:"pipelineDigest"`
	GenerationSettingsFingerprint string `json:"generationSettingsFingerprint"`
	EmbeddingFingerprint          string `json:"embeddingFingerprint"`
}

type PreparedValue struct {
	Key  string          `json:"key"`
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

type PreparedCorpusBundle struct {
	SchemaVersion string                 `json:"schemaVersion"`
	Identity      PreparedCorpusIdentity `json:"identity"`
	Digest        string                 `json:"digest"`
	Values        []PreparedValue        `json:"values"`
}

type PreparedCorpusStore interface {
	Open(context.Context, *Engine, ragcontract.PipelineIR, ragoperators.Corpus, Options, PreparedCorpusIdentity) (*Prepared, bool, error)
	Put(context.Context, *Prepared, PreparedCorpusIdentity) (string, error)
}

// PreparedCorpusPublication supplies only validated static values to the existing
// atomic prepared-corpus store. Provider clients remain process-local in Options.
type PreparedCorpusPublication struct {
	Store    PreparedCorpusStore
	Engine   *Engine
	Pipeline ragcontract.PipelineIR
	Corpus   ragoperators.Corpus
	Options  Options
	Identity PreparedCorpusIdentity
	Values   map[string]any
}

// PublishPreparedCorpus atomically persists validated static values, then reopens
// the bundle to verify that required live indexes rebuild without provider calls.
func PublishPreparedCorpus(ctx context.Context, publication PreparedCorpusPublication) (string, error) {
	if publication.Store == nil || publication.Engine == nil {
		return "", fmt.Errorf("RAG_PREPARED_PUBLICATION_STORE")
	}
	pipelineDigest, err := ragcontract.Digest(publication.Pipeline)
	if err != nil {
		return "", err
	}
	if publication.Identity.SchemaVersion != preparedCorpusSchemaVersion || publication.Identity.PipelineDigest != pipelineDigest || publication.Identity.CorpusDigest == "" {
		return "", fmt.Errorf("RAG_PREPARED_PUBLICATION_IDENTITY")
	}
	prepared, err := NewPreparedFromStaticValues(publication.Pipeline, publication.Values)
	if err != nil {
		return "", err
	}
	defer func() { _ = prepared.Close() }()
	digest, err := publication.Store.Put(ctx, prepared, publication.Identity)
	if err != nil {
		return "", err
	}
	reopened, found, err := publication.Store.Open(ctx, publication.Engine, publication.Pipeline, publication.Corpus, publication.Options, publication.Identity)
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("RAG_PREPARED_PUBLICATION_REOPEN")
	}
	return digest, reopened.Close()
}

type FilePreparedCorpusStore struct{ directory string }

func NewFilePreparedCorpusStore(directory string) (*FilePreparedCorpusStore, error) {
	if directory == "" {
		return nil, os.ErrInvalid
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, err
	}
	return &FilePreparedCorpusStore{directory: directory}, nil
}

func (s *FilePreparedCorpusStore) path(identity PreparedCorpusIdentity) (string, error) {
	canonical, err := ragcontract.CanonicalJSON(identity)
	if err != nil {
		return "", err
	}
	digest, err := ragcontract.Digest(canonical)
	if err != nil {
		return "", err
	}
	return filepath.Join(s.directory, digest[7:]+".json"), nil
}

func (s *FilePreparedCorpusStore) Put(ctx context.Context, prepared *Prepared, identity PreparedCorpusIdentity) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if prepared == nil {
		return "", fmt.Errorf("RAG_PREPARED_STORE_NIL")
	}
	values, err := serializePreparedValues(prepared.values)
	if err != nil {
		return "", err
	}
	bundle := PreparedCorpusBundle{SchemaVersion: preparedCorpusSchemaVersion, Identity: identity, Values: values}
	bundle.Digest, err = preparedBundleDigest(bundle)
	if err != nil {
		return "", err
	}
	data, err := ragcontract.CanonicalJSON(bundle)
	if err != nil {
		return "", err
	}
	path, err := s.path(identity)
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp(s.directory, ".prepared-*")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return "", err
	}
	return bundle.Digest, nil
}

func (s *FilePreparedCorpusStore) Open(ctx context.Context, engine *Engine, pipeline ragcontract.PipelineIR, corpus ragoperators.Corpus, options Options, identity PreparedCorpusIdentity) (*Prepared, bool, error) {
	if err := ctx.Err(); err != nil {
		return nil, false, err
	}
	path, err := s.path(identity)
	if err != nil {
		return nil, false, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var bundle PreparedCorpusBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return nil, false, nil
	}
	want, err := preparedBundleDigest(bundle)
	if err != nil || bundle.SchemaVersion != preparedCorpusSchemaVersion || bundle.Digest != want || !preparedIdentityEqual(bundle.Identity, identity) {
		return nil, false, nil
	}
	values, err := deserializePreparedValues(bundle.Values)
	if err != nil {
		return nil, false, nil
	}
	values["corpus/out"] = corpus
	static := staticNodeIDs(pipeline)
	for _, node := range pipeline.Nodes {
		if !static[node.ID] || node.Operator.Kind != "index.bleve-multi" {
			continue
		}
		if _, present := values[node.ID+"/index"]; present {
			continue
		}
		op, ok := engine.Registry.Lookup(node.Operator)
		if !ok {
			return nil, false, fmt.Errorf("RAG_OPERATOR_UNAVAILABLE: %s", node.Operator.ID())
		}
		inputs := map[string]any{}
		for _, binding := range node.Inputs {
			value, exists := values[binding.From.NodeID+"/"+binding.From.Port]
			if !exists {
				return nil, false, nil
			}
			inputs[binding.Port] = value
		}
		env := &ragoperators.Environment{Manifests: options.Manifests, Schemas: options.Schemas, Generator: options.Generator, Embedder: options.Embedder, Reranker: options.Reranker, Cache: options.Cache, GenerationConcurrency: options.GenerationConcurrency, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint}
		outputs, err := op.Execute(ctx, node, inputs, env)
		if err != nil {
			return nil, false, err
		}
		for port, value := range outputs {
			if port != "artifact" && port != "manifest" {
				values[node.ID+"/"+port] = value
			}
		}
	}
	prepared := &Prepared{pipelineDigest: mustDigest(pipeline), values: selectPrepared(values, static)}
	return prepared, true, nil
}

func serializePreparedValues(values map[string]any) ([]PreparedValue, error) {
	out := make([]PreparedValue, 0, len(values))
	for key, value := range values {
		var kind string
		switch value.(type) {
		case []ragoperators.Unit:
			kind = "units"
		case []ragoperators.Chunk:
			kind = "chunks"
		case []ragoperators.Representation:
			kind = "representations"
		case []ragoperators.Embedding:
			kind = "embeddings"
		default:
			continue // corpus and live indexes are restored/rebuilt separately.
		}
		data, err := ragcontract.CanonicalJSON(value)
		if err != nil {
			return nil, err
		}
		out = append(out, PreparedValue{Key: key, Kind: kind, Data: data})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}

func deserializePreparedValues(values []PreparedValue) (map[string]any, error) {
	out := map[string]any{}
	for _, value := range values {
		if value.Key == "" || value.Kind == "" {
			return nil, fmt.Errorf("RAG_PREPARED_BUNDLE_VALUE")
		}
		switch value.Kind {
		case "units":
			var decoded []ragoperators.Unit
			if err := json.Unmarshal(value.Data, &decoded); err != nil {
				return nil, err
			}
			out[value.Key] = decoded
		case "chunks":
			var decoded []ragoperators.Chunk
			if err := json.Unmarshal(value.Data, &decoded); err != nil {
				return nil, err
			}
			out[value.Key] = decoded
		case "representations":
			var decoded []ragoperators.Representation
			if err := json.Unmarshal(value.Data, &decoded); err != nil {
				return nil, err
			}
			out[value.Key] = decoded
		case "embeddings":
			var decoded []ragoperators.Embedding
			if err := json.Unmarshal(value.Data, &decoded); err != nil {
				return nil, err
			}
			out[value.Key] = decoded
		default:
			return nil, fmt.Errorf("RAG_PREPARED_BUNDLE_KIND")
		}
	}
	return out, nil
}

func preparedBundleDigest(bundle PreparedCorpusBundle) (string, error) {
	bundle.Digest = ""
	return ragcontract.Digest(bundle)
}

func preparedIdentityEqual(left, right PreparedCorpusIdentity) bool {
	leftJSON, leftErr := ragcontract.CanonicalJSON(left)
	rightJSON, rightErr := ragcontract.CanonicalJSON(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}
