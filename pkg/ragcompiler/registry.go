// Package ragcompiler normalizes and compiles ragcontract values without executing providers.
package ragcompiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type PortDefinition struct {
	Name string
	Kind ragcontract.PortKind
}
type OperatorDefinition struct {
	Ref                ragcontract.OperatorRef
	Inputs, Outputs    []PortDefinition
	Defaults           map[string]any
	Capabilities       []string
	Resources          []string
	ObservationSchemas []string
	Recipe             bool
}
type Registry struct{ definitions map[string]OperatorDefinition }

func NewRegistry(definitions ...OperatorDefinition) (*Registry, error) {
	r := &Registry{definitions: map[string]OperatorDefinition{}}
	for _, d := range definitions {
		id := d.Ref.ID()
		if id == "/" {
			return nil, fmt.Errorf("RAG_V2_OPERATOR_REF: empty operator")
		}
		if _, ok := r.definitions[id]; ok {
			return nil, fmt.Errorf("RAG_V2_OPERATOR_DUPLICATE: %s", id)
		}
		r.definitions[id] = d
	}
	return r, nil
}
func (r *Registry) Definition(ref ragcontract.OperatorRef) (OperatorDefinition, bool) {
	if r == nil {
		return OperatorDefinition{}, false
	}
	d, ok := r.definitions[ref.ID()]
	return d, ok
}
func (r *Registry) Definitions() []OperatorDefinition {
	if r == nil {
		return nil
	}
	values := make([]OperatorDefinition, 0, len(r.definitions))
	for _, definition := range r.definitions {
		values = append(values, definition)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Ref.ID() < values[j].Ref.ID() })
	return values
}

func BuiltinRegistry() *Registry {
	defs := []OperatorDefinition{
		op("units.identity", "v1", []PortDefinition{p("corpus", ragcontract.PortCorpus)}, p("units", ragcontract.PortUnits)),
		op("units.individual-turns", "v1", []PortDefinition{p("corpus", ragcontract.PortCorpus)}, p("units", ragcontract.PortUnits)),
		op("transcript.units.agents-view-runs", "v1", []PortDefinition{p("corpus", ragcontract.PortCorpus)}, p("units", ragcontract.PortUnits)),
		op("chunks.identity", "v1", []PortDefinition{p("units", ragcontract.PortUnits)}, p("chunks", ragcontract.PortChunks)),
		withDefaults(op("chunks.recursive", "v1", []PortDefinition{p("units", ragcontract.PortUnits)}, p("chunks", ragcontract.PortChunks)), map[string]any{"size": 800, "overlap": 120, "unicodePolicy": "utf8-boundary"}),
		op("representations.raw", "v1", []PortDefinition{p("chunks", ragcontract.PortChunks)}, p("representations", ragcontract.PortRepresentations)),
		op("representations.structured-summary", "v1", []PortDefinition{p("chunks", ragcontract.PortChunks)}, p("representations", ragcontract.PortRepresentations)),
		op("representations.synthetic-questions", "v1", []PortDefinition{p("chunks", ragcontract.PortChunks), p("source", ragcontract.PortRepresentations)}, p("representations", ragcontract.PortRepresentations)),
		op("representations.combined-summary-questions", "v1", []PortDefinition{p("chunks", ragcontract.PortChunks)}, p("representations", ragcontract.PortRepresentations)),
		{Ref: ragcontract.OperatorRef{Kind: "representations.compose", Version: "v1"}, Inputs: []PortDefinition{p("in", ragcontract.PortChunks)}, Outputs: []PortDefinition{p("out", ragcontract.PortRepresentations)}, Recipe: true},
		op("representations.merge", "v1", nil, p("representations", ragcontract.PortRepresentations)),
		op("embed.model", "v1", []PortDefinition{p("representations", ragcontract.PortRepresentations)}, p("embeddings", ragcontract.PortEmbeddings)),
		op2("index.bleve-multi", "v1", []PortDefinition{p("embeddings", ragcontract.PortEmbeddings)}, p("index", ragcontract.PortIndex)),
		op("index.memory-smoke", "v1", []PortDefinition{p("representations", ragcontract.PortRepresentations)}, p("index", ragcontract.PortIndex)),
		withDefaults(op2("retrieve.bm25", "v1", []PortDefinition{p("index", ragcontract.PortIndex), p("query", ragcontract.PortQuery)}, p("hits", ragcontract.PortRankedRecords)), map[string]any{"topK": 30, "filter": map[string]any{}}),
		withDefaults(op2("retrieve.vector", "v1", []PortDefinition{p("index", ragcontract.PortIndex), p("query", ragcontract.PortQuery)}, p("hits", ragcontract.PortRankedRecords)), map[string]any{"topK": 30, "filter": map[string]any{}}),
		withDefaults(op("collapse.parent", "v1", []PortDefinition{p("hits", ragcontract.PortRankedRecords)}, p("parents", ragcontract.PortRankedParents)), map[string]any{"scope": "unit", "representative": "scoreThenRepresentationId"}),
		withDefaults(op("fusion.weighted-rrf", "v1", []PortDefinition{p("channels", ragcontract.PortRankedParents)}, p("parents", ragcontract.PortRankedParents)), map[string]any{"rankConstant": 60, "weights": map[string]any{}}),
		withDefaults(op("collapse.final", "v1", []PortDefinition{p("parents", ragcontract.PortRankedParents)}, p("parents", ragcontract.PortRankedParents)), map[string]any{"scope": "unit", "representative": "scoreThenRepresentationId"}),
		op("hydrate.source-evidence", "v1", []PortDefinition{p("parents", ragcontract.PortRankedParents), p("chunks", ragcontract.PortChunks)}, p("evidence", ragcontract.PortEvidence)),
		op("rerank.cross-encoder", "v1", []PortDefinition{p("evidence", ragcontract.PortEvidence)}, p("evidence", ragcontract.PortEvidence)),
		op("generate.answer", "v1", []PortDefinition{p("evidence", ragcontract.PortEvidence)}, p("answer", ragcontract.PortAnswer)),
		op("evaluate.relevance", "v1", []PortDefinition{p("evidence", ragcontract.PortEvidence)}, p("evaluation", ragcontract.PortEvaluation)),
	}
	r, err := NewRegistry(defs...)
	if err != nil {
		panic(err)
	}
	return r
}
func p(name string, kind ragcontract.PortKind) PortDefinition {
	return PortDefinition{Name: name, Kind: kind}
}
func op(kind, version string, inputs []PortDefinition, output PortDefinition) OperatorDefinition {
	capabilities, resources, observations := operatorMetadata(kind)
	return OperatorDefinition{Ref: ragcontract.OperatorRef{Kind: kind, Version: version}, Inputs: inputs, Outputs: []PortDefinition{output}, Capabilities: capabilities, Resources: resources, ObservationSchemas: observations}
}
func operatorMetadata(kind string) ([]string, []string, []string) {
	observations := []string{"rag-operator-event/v2", ragcontract.TraceSchemaVersion}
	switch {
	case strings.HasPrefix(kind, "units.") || strings.HasPrefix(kind, "transcript.units."):
		return []string{"immutable-unitization"}, []string{"memory", "artifact-store"}, append(observations, ragcontract.UnitSetManifestSchema)
	case strings.HasPrefix(kind, "chunks."):
		return []string{"unicode-safe-chunking"}, []string{"memory", "artifact-store"}, append(observations, ragcontract.ChunkSetManifestSchema)
	case kind == "representations.raw" || kind == "representations.merge":
		return []string{"materialize-representations"}, []string{"memory", "artifact-store"}, append(observations, ragcontract.RepresentationManifestSchema)
	case strings.HasPrefix(kind, "representations."):
		return []string{"materialize-generated-representations"}, []string{"model-manifest", "prompt-manifest", "generator-provider", "content-cache", "artifact-store"}, append(observations, ragcontract.RepresentationManifestSchema)
	case kind == "embed.model":
		return []string{"batch-embedding"}, []string{"model-manifest", "embedding-provider"}, append(observations, ragcontract.EmbeddingManifestSchema)
	case strings.HasPrefix(kind, "index."):
		return []string{"immutable-index"}, []string{"memory", "artifact-store"}, append(observations, ragcontract.IndexManifestSchema)
	case strings.HasPrefix(kind, "retrieve."):
		return []string{"filtered-retrieval"}, []string{"index"}, observations
	case kind == "rerank.cross-encoder":
		return []string{"explicit-reranking"}, []string{"model-manifest", "reranker-provider"}, observations
	case kind == "generate.answer":
		return []string{"grounded-generation"}, []string{"model-manifest", "prompt-manifest", "generator-provider"}, observations
	default:
		return []string{"native-execution"}, []string{"memory"}, observations
	}
}
func op2(kind, version string, inputs []PortDefinition, output PortDefinition) OperatorDefinition {
	return op(kind, version, inputs, output)
}
func withDefaults(d OperatorDefinition, v map[string]any) OperatorDefinition {
	d.Defaults = v
	return d
}
func normalizeConfig(raw json.RawMessage, definition OperatorDefinition) (json.RawMessage, error) {
	var value map[string]any
	canonical, err := ragcontract.CanonicalRaw(raw, "{}")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(canonical, &value); err != nil {
		return nil, fmt.Errorf("config must be an object: %w", err)
	}
	if value == nil {
		return nil, fmt.Errorf("config must be an object")
	}
	allowed := configFields(definition.Ref.Kind)
	for key := range value {
		if !allowed[key] {
			return nil, fmt.Errorf("unknown config field %q for %s", key, definition.Ref.ID())
		}
	}
	for k, v := range definition.Defaults {
		if _, ok := value[k]; !ok {
			value[k] = v
		}
	}
	return ragcontract.CanonicalRaw(mustJSON(value), "{}")
}
func configFields(kind string) map[string]bool {
	fields := map[string][]string{
		"chunks.recursive":                           {"size", "overlap", "levels", "atomic", "unicodePolicy", "emptyInputPolicy"},
		"representations.raw":                        {"name"},
		"representations.structured-summary":         {"name", "model", "prompt", "outputSchema", "decoding", "seedPolicy"},
		"representations.synthetic-questions":        {"name", "from", "model", "prompt", "count", "decoding", "seedPolicy"},
		"representations.combined-summary-questions": {"model", "prompt", "outputSchema", "batchSize", "questionsPerChunk", "maxBatchRunes"},
		"representations.merge":                      {},
		"embed.model":                                {"model", "dimensions", "distance", "normalize", "batchSize", "preprocessing"},
		"index.bleve-multi":                          {"engineVersion", "mapping", "representationKinds", "lexical", "vector"},
		"index.memory-smoke":                         {"representationKinds"},
		"retrieve.bm25":                              {"index", "representation", "topK", "filter"},
		"retrieve.vector":                            {"index", "representation", "topK", "filter"},
		"collapse.parent":                            {"scope", "representative"},
		"fusion.weighted-rrf":                        {"rankConstant", "weights", "missingChannelPolicy", "tieBreak"},
		"collapse.final":                             {"scope", "representative"},
		"hydrate.source-evidence":                    {"policy", "allSupportingChunks", "results"},
		"rerank.cross-encoder":                       {"model", "candidateCount", "results", "truncation", "tokenization", "inputTemplate", "timeoutMilliseconds"},
		"generate.answer":                            {"model", "prompt", "citations", "citationFailurePolicy", "contextBudgetTokens", "decoding", "seedPolicy"},
		"evaluate.relevance":                         {"target", "gradeThreshold", "measures"},
	}
	result := map[string]bool{}
	for _, field := range fields[kind] {
		result[field] = true
	}
	return result
}

func mustJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
func strictRaw(raw json.RawMessage, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("trailing JSON content")
	}
	return nil
}
