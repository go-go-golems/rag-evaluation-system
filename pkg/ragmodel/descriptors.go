// Package ragmodel provides pure Go authoring builders for canonical RAG v2 plans.
// It performs no filesystem, database, network, provider, index, or lifecycle work.
package ragmodel

import (
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type DescriptorKind string

const (
	KindCorpus          DescriptorKind = "corpus"
	KindUnitizer        DescriptorKind = "unitizer"
	KindChunker         DescriptorKind = "chunker"
	KindRepresentations DescriptorKind = "representations"
	KindEmbedding       DescriptorKind = "embedding"
	KindIndex           DescriptorKind = "index"
	KindRetriever       DescriptorKind = "retriever"
	KindCollapse        DescriptorKind = "collapse"
	KindFusion          DescriptorKind = "fusion"
	KindHydration       DescriptorKind = "hydration"
	KindReranker        DescriptorKind = "reranker"
	KindGeneration      DescriptorKind = "generation"
)

type Descriptor struct {
	Kind     DescriptorKind
	Name     string
	Operator ragcontract.OperatorRef
	Config   json.RawMessage
	Children []*Descriptor
}
type CorpusInput struct {
	Role           string
	ManifestSchema string
}
type DatasetRef struct{ Role, Split, Status, RelevanceTarget string }
type FactorRef struct {
	ID string `json:"$factor"`
}
type RecursiveChunkConfig struct {
	MaxRunes     int      `json:"maxRunes"`
	OverlapSpans int      `json:"overlapSpans"`
	Levels       []string `json:"levels,omitempty"`
	Atomic       []string `json:"atomic,omitempty"`
}
type StructuredGenerationConfig struct {
	Model        string                  `json:"model"`
	Prompt       string                  `json:"prompt"`
	OutputSchema string                  `json:"outputSchema"`
	Decoding     json.RawMessage         `json:"decoding,omitempty"`
	SeedPolicy   *ragcontract.SeedPolicy `json:"seedPolicy,omitempty"`
}
type StructuredSummaryConfig struct {
	Generator StructuredGenerationConfig `json:"generator"`
}
type SyntheticQuestionsConfig struct {
	From   string `json:"from"`
	Count  int    `json:"count"`
	Model  string `json:"model,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}
type CombinedPreparationConfig struct {
	Model             string `json:"model"`
	Prompt            string `json:"prompt"`
	OutputSchema      string `json:"outputSchema"`
	BatchSize         int    `json:"batchSize"`
	QuestionsPerChunk int    `json:"questionsPerChunk"`
	MaxBatchRunes     int    `json:"maxBatchRunes"`
}
type EmbeddingConfig struct {
	Dimensions int    `json:"dimensions,omitempty"`
	Distance   string `json:"distance,omitempty"`
	Normalize  string `json:"normalize,omitempty"`
	BatchSize  int    `json:"batchSize,omitempty"`
}
type BleveMultiConfig struct {
	Lexical bool               `json:"lexical"`
	Vector  *VectorIndexConfig `json:"vector,omitempty"`
}
type VectorIndexConfig struct {
	Distance    string `json:"distance"`
	OptimizeFor string `json:"optimizeFor,omitempty"`
}
type RetrieveConfig struct {
	Index          string       `json:"index"`
	Representation string       `json:"representation"`
	TopK           int          `json:"topK"`
	Filter         FilterConfig `json:"filter,omitempty"`
}
type FilterConfig struct {
	SourceIDs      []string          `json:"sourceIds,omitempty"`
	DocumentIDs    []string          `json:"documentIds,omitempty"`
	ContentTypes   []string          `json:"contentTypes,omitempty"`
	MetadataEquals map[string]string `json:"metadataEquals,omitempty"`
}
type CollapseConfig struct {
	Scope          any    `json:"scope"`
	Representative string `json:"representative"`
}
type WeightedRRFConfig struct {
	RankConstant         int                `json:"rankConstant"`
	Weights              map[string]float64 `json:"weights,omitempty"`
	MissingChannelPolicy string             `json:"missingChannelPolicy,omitempty"`
	TieBreak             string             `json:"tieBreak,omitempty"`
}
type HydrationConfig struct {
	Selection           string `json:"selection"`
	AllSupportingChunks bool   `json:"allSupportingChunks,omitempty"`
}
type CrossEncoderConfig struct {
	Model               string `json:"model"`
	Candidates          int    `json:"candidates"`
	Results             int    `json:"results"`
	Truncation          string `json:"truncation,omitempty"`
	Tokenization        string `json:"tokenization,omitempty"`
	InputTemplate       string `json:"inputTemplate,omitempty"`
	TimeoutMilliseconds int64  `json:"timeoutMilliseconds,omitempty"`
}
type AnswerConfig struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	Citations string `json:"citations"`
	// CitationFailurePolicy is explicit: "error" (default) rejects an
	// ungrounded answer, while "abstain" discards it and returns a safe abstention.
	CitationFailurePolicy string                  `json:"citationFailurePolicy,omitempty"`
	ContextBudgetTokens   int                     `json:"contextBudgetTokens"`
	Decoding              json.RawMessage         `json:"decoding,omitempty"`
	SeedPolicy            *ragcontract.SeedPolicy `json:"seedPolicy,omitempty"`
}

func desc(kind DescriptorKind, name, operator string, config any) *Descriptor {
	raw, _ := json.Marshal(config)
	parts := splitOperator(operator)
	return &Descriptor{Kind: kind, Name: name, Operator: ragcontract.OperatorRef{Kind: parts[0], Version: parts[1]}, Config: raw}
}
func splitOperator(value string) [2]string {
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == '/' {
			return [2]string{value[:i], value[i+1:]}
		}
	}
	return [2]string{value, "v1"}
}
func Corpus(role string) CorpusInput {
	return CorpusInput{Role: role, ManifestSchema: ragcontract.CorpusManifestSchema}
}
func UnitsIdentity() *Descriptor {
	return desc(KindUnitizer, "", "units.identity/v1", map[string]any{})
}
func IndividualTurns() *Descriptor {
	return desc(KindUnitizer, "", "units.individual-turns/v1", map[string]any{})
}
func AgentsViewRuns() *Descriptor {
	return desc(KindUnitizer, "", "transcript.units.agents-view-runs/v1", map[string]any{})
}
func RecursiveChunks(config RecursiveChunkConfig) *Descriptor {
	return desc(KindChunker, "", "chunks.recursive/v1", map[string]any{"size": config.MaxRunes, "overlap": config.OverlapSpans, "levels": config.Levels, "atomic": config.Atomic})
}
func RawRepresentation(name string) *Descriptor {
	return desc(KindRepresentations, name, "representations.raw/v1", map[string]any{"name": name})
}
func StructuredSummary(name string, config StructuredSummaryConfig) *Descriptor {
	return desc(KindRepresentations, name, "representations.structured-summary/v1", map[string]any{"name": name, "model": config.Generator.Model, "prompt": config.Generator.Prompt, "outputSchema": config.Generator.OutputSchema, "decoding": config.Generator.Decoding, "seedPolicy": config.Generator.SeedPolicy})
}
func SyntheticQuestions(name string, config SyntheticQuestionsConfig) *Descriptor {
	return desc(KindRepresentations, name, "representations.synthetic-questions/v1", map[string]any{"name": name, "from": config.From, "count": config.Count, "model": config.Model, "prompt": config.Prompt})
}
func CombinedPreparation(config CombinedPreparationConfig) *Descriptor {
	return desc(KindRepresentations, "", "representations.combined-summary-questions/v1", config)
}
func ComposeRepresentations(values ...*Descriptor) *Descriptor {
	return &Descriptor{Kind: KindRepresentations, Operator: ragcontract.OperatorRef{Kind: "representations.compose", Version: "v1"}, Children: append([]*Descriptor(nil), values...)}
}
func EmbeddingModel(name string, config EmbeddingConfig) *Descriptor {
	return desc(KindEmbedding, name, "embed.model/v1", map[string]any{"model": name, "dimensions": config.Dimensions, "distance": config.Distance, "normalize": config.Normalize, "batchSize": config.BatchSize})
}
func BleveMulti(config BleveMultiConfig) *Descriptor {
	return desc(KindIndex, "", "index.bleve-multi/v1", map[string]any{"lexical": config.Lexical, "vector": config.Vector})
}
func BM25(name string, config RetrieveConfig) *Descriptor {
	return desc(KindRetriever, name, "retrieve.bm25/v1", config)
}
func Vector(name string, config RetrieveConfig) *Descriptor {
	return desc(KindRetriever, name, "retrieve.vector/v1", config)
}
func ParentCollapse(config CollapseConfig) *Descriptor {
	return desc(KindCollapse, "", "collapse.parent/v1", config)
}
func WeightedRRF(config WeightedRRFConfig) *Descriptor {
	return desc(KindFusion, "", "fusion.weighted-rrf/v1", config)
}
func SourceEvidence(config HydrationConfig) *Descriptor {
	return desc(KindHydration, "", "hydrate.source-evidence/v1", map[string]any{"policy": config.Selection, "allSupportingChunks": config.AllSupportingChunks})
}
func CrossEncoder(config CrossEncoderConfig) *Descriptor {
	return desc(KindReranker, "", "rerank.cross-encoder/v1", config)
}
func StructuredGenerator(name string, config StructuredGenerationConfig) *Descriptor {
	return desc(KindGeneration, name, "generation.structured/v1", config)
}
func Answer(config AnswerConfig) *Descriptor {
	return desc(KindGeneration, "", "generate.answer/v1", config)
}
func DatasetArtifact(role string, split, status, target string) DatasetRef {
	return DatasetRef{Role: role, Split: split, Status: status, RelevanceTarget: target}
}
