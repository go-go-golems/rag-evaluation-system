package ragoperators

import (
	"context"
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type SourceRecord struct {
	ID            string          `json:"id"`
	SessionID     string          `json:"sessionId"`
	Ordinal       int64           `json:"ordinal"`
	Role          string          `json:"role"`
	Text          string          `json:"text"`
	ContentDigest string          `json:"contentDigest,omitempty"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}
type Corpus struct {
	SchemaVersion string         `json:"schemaVersion"`
	Records       []SourceRecord `json:"records"`
}
type Unit struct {
	Record         ragcontract.UnitRecord `json:"record"`
	Text           string                 `json:"text"`
	Records        []SourceRecord         `json:"records"`
	ManifestDigest string                 `json:"manifestDigest,omitempty"`
}
type Chunk struct {
	Record         ragcontract.ChunkRecord   `json:"record"`
	Text           string                    `json:"text"`
	Ranges         []ragcontract.SourceRange `json:"ranges"`
	ManifestDigest string                    `json:"manifestDigest,omitempty"`
}
type Representation struct {
	Record         ragcontract.RepresentationRecord `json:"record"`
	Text           string                           `json:"text"`
	ManifestDigest string                           `json:"manifestDigest,omitempty"`
}
type Embedding struct {
	Record         ragcontract.EmbeddingRecord `json:"record"`
	Vector         []float64                   `json:"vector"`
	ManifestDigest string                      `json:"manifestDigest,omitempty"`
}
type Query struct {
	ID          string             `json:"id"`
	Text        string             `json:"text"`
	RelevantIDs []string           `json:"relevantIds,omitempty"`
	Grades      map[string]float64 `json:"grades,omitempty"`
}
type RetrievalFilter struct {
	SourceIDs      []string          `json:"sourceIds,omitempty"`
	DocumentIDs    []string          `json:"documentIds,omitempty"`
	ContentTypes   []string          `json:"contentTypes,omitempty"`
	MetadataEquals map[string]string `json:"metadataEquals,omitempty"`
}
type EvaluationDataset struct {
	SchemaVersion string  `json:"schemaVersion"`
	Queries       []Query `json:"queries"`
}
type RankedRecord struct {
	Rank           int            `json:"rank"`
	Representation Representation `json:"representation"`
	Score          float64        `json:"score"`
	Channel        string         `json:"channel"`
}
type RankedParent struct {
	Rank           int                              `json:"rank"`
	Identity       ragcontract.CollapseIdentity     `json:"identity"`
	Score          float64                          `json:"score"`
	Representative Representation                   `json:"representative"`
	Members        []RankedRecord                   `json:"members"`
	Contributions  []ragcontract.FusionContribution `json:"contributions,omitempty"`
}
type Evidence struct {
	Rank          int                              `json:"rank"`
	Collapse      ragcontract.CollapseIdentity     `json:"collapse"`
	Chunk         Chunk                            `json:"chunk"`
	Score         float64                          `json:"score"`
	Contributions []ragcontract.FusionContribution `json:"contributions,omitempty"`
	Matched       []RankedRecord                   `json:"matchedRepresentations,omitempty"`
	RerankerScore *float64                         `json:"rerankerScore,omitempty"`
}
type Answer struct {
	Text             string   `json:"text"`
	CitationChunkIDs []string `json:"citationChunkIds"`
	FinishReason     string   `json:"finishReason"`
	Abstained        bool     `json:"abstained"`
	InputTokens      int64    `json:"inputTokens"`
	OutputTokens     int64    `json:"outputTokens"`
}
type Usage struct {
	InputTokens, OutputTokens, EmbeddingTokens int64
	Cost                                       map[string]float64
}
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type Artifact struct {
	Role, Kind, Name, SchemaVersion, MediaType string
	Data                                       []byte
	Metadata                                   json.RawMessage
}
type Metric struct {
	Name, Unit string
	Value      json.RawMessage
	Numeric    *float64
	Metadata   json.RawMessage
}

type TextGenerator interface {
	Generate(context.Context, GenerationRequest) (GenerationResult, error)
}
type GenerationRequest struct {
	Kind, Model, Prompt, OutputSchema, ParentID, Text string
	Count                                             int
	Evidence                                          []Evidence
}
type GenerationResult struct {
	Text                      string
	Questions                 []string
	CitationChunkIDs          []string
	InputTokens, OutputTokens int64
	Cost                      float64
	FinishReason              string
	Abstained                 bool
}
type Embedder interface {
	Embed(context.Context, string, []string) ([][]float64, Usage, error)
}
type Reranker interface {
	Rerank(context.Context, RerankRequest) ([]RerankScore, error)
}
type RerankRequest struct {
	Model, InputTemplate, Truncation, Tokenization string
	Query                                          string
	Candidates                                     []Evidence
	Results                                        int
}
type RerankScore struct {
	ChunkID string
	Score   float64
}

type ManifestResolver interface {
	Model(string) (ragcontract.ModelManifest, error)
	Prompt(string) (ragcontract.PromptManifest, error)
}
type OutputSchemaValidator interface {
	Validate(schema string, document json.RawMessage) error
}
type Environment struct {
	Manifests    ManifestResolver
	Schemas      OutputSchemaValidator
	Generator    TextGenerator
	Embedder     Embedder
	Reranker     Reranker
	Cache        Cache
	Trace        *ragcontract.QueryTrace
	CurrentQuery Query
	QueryText    string
	Usage        Usage
}
type Cache interface {
	Get(string) ([]byte, bool)
	Put(string, []byte)
}
