package ragcontract

import "encoding/json"

const TraceSchemaVersion = "rag-query-trace/v2"

type RepresentationIdentity struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	ParentChunkID string `json:"parentChunkId"`
	ParentUnitID  string `json:"parentUnitId"`
	ContentDigest string `json:"contentDigest"`
	EvidenceRole  string `json:"evidenceRole"`
}
type CollapseIdentity struct {
	Scope string `json:"scope"`
	ID    string `json:"id"`
}
type EvidenceIdentity struct {
	ChunkID  string      `json:"chunkId"`
	Digest   string      `json:"digest"`
	Citation CitationRef `json:"citation"`
}
type CitationRef struct {
	SourceID     string `json:"sourceId"`
	ByteStart    int64  `json:"byteStart,omitempty"`
	ByteEnd      int64  `json:"byteEnd,omitempty"`
	OrdinalStart int64  `json:"ordinalStart,omitempty"`
	OrdinalEnd   int64  `json:"ordinalEnd,omitempty"`
}

type QueryTrace struct {
	SchemaVersion string           `json:"schemaVersion"`
	Query         QueryInputTrace  `json:"query"`
	Operators     []OperatorTrace  `json:"operators"`
	Channels      []ChannelTrace   `json:"channels"`
	Collapses     []CollapseTrace  `json:"collapses"`
	Fusion        *FusionTrace     `json:"fusion,omitempty"`
	Hydration     *HydrationTrace  `json:"hydration,omitempty"`
	Reranking     *RerankingTrace  `json:"reranking,omitempty"`
	Generation    *GenerationTrace `json:"generation,omitempty"`
	Results       []ResultTrace    `json:"results"`
	Relevance     *RelevanceTrace  `json:"relevance,omitempty"`
	Timing        TimingTrace      `json:"timing"`
	Usage         UsageTrace       `json:"usage"`
	Failures      []FailureTrace   `json:"failures"`
}
type QueryInputTrace struct {
	ID           string          `json:"id"`
	TextDigest   string          `json:"textDigest"`
	DatasetSplit string          `json:"datasetSplit"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}
type OperatorTrace struct {
	NodeID               string      `json:"nodeId"`
	Operator             OperatorRef `json:"operator"`
	Status               string      `json:"status"`
	InputCount           int         `json:"inputCount,omitempty"`
	OutputCount          int         `json:"outputCount,omitempty"`
	DurationMilliseconds int64       `json:"durationMilliseconds"`
}
type ChannelHit struct {
	Rank            int                    `json:"rank"`
	Representation  RepresentationIdentity `json:"representation"`
	RawScore        float64                `json:"rawScore"`
	NormalizedScore *float64               `json:"normalizedScore,omitempty"`
	Filter          json.RawMessage        `json:"filter,omitempty"`
}
type ChannelTrace struct {
	Name     string       `json:"name"`
	Operator OperatorRef  `json:"operator"`
	Hits     []ChannelHit `json:"hits"`
}
type CollapseMember struct {
	RepresentationID string  `json:"representationId"`
	Rank             int     `json:"rank"`
	Score            float64 `json:"score"`
}
type CollapseGroup struct {
	Key                      CollapseIdentity `json:"key"`
	Members                  []CollapseMember `json:"members"`
	SelectedRepresentationID string           `json:"selectedRepresentationId"`
	Rank                     int              `json:"rank"`
}
type CollapseTrace struct {
	Stage    string          `json:"stage"`
	Channel  string          `json:"channel,omitempty"`
	Operator OperatorRef     `json:"operator"`
	Groups   []CollapseGroup `json:"groups"`
}
type FusionContribution struct {
	Channel string  `json:"channel"`
	Rank    int     `json:"rank"`
	Weight  float64 `json:"weight"`
	Value   float64 `json:"value"`
}
type FusionResult struct {
	Rank          int                  `json:"rank"`
	Identity      CollapseIdentity     `json:"identity"`
	Contributions []FusionContribution `json:"contributions"`
	Score         float64              `json:"score"`
}
type FusionTrace struct {
	Operator             OperatorRef    `json:"operator"`
	Results              []FusionResult `json:"results"`
	MissingChannelPolicy string         `json:"missingChannelPolicy"`
	TieBreak             string         `json:"tieBreak"`
}
type HydrationCandidate struct {
	Collapse     CollapseIdentity `json:"collapse"`
	Evidence     EvidenceIdentity `json:"evidence"`
	Contribution float64          `json:"contribution"`
}
type HydrationTrace struct {
	Operator   OperatorRef          `json:"operator"`
	Candidates []HydrationCandidate `json:"candidates"`
	Selected   []EvidenceIdentity   `json:"selected"`
}
type RerankingEntry struct {
	Evidence       EvidenceIdentity `json:"evidence"`
	BeforeRank     int              `json:"beforeRank"`
	AfterRank      int              `json:"afterRank"`
	RetrievalScore float64          `json:"retrievalScore"`
	RerankerScore  float64          `json:"rerankerScore"`
}
type RerankingTrace struct {
	Operator             OperatorRef      `json:"operator"`
	ModelManifestDigest  string           `json:"modelManifestDigest"`
	InputPolicy          string           `json:"inputPolicy"`
	InputTemplate        string           `json:"inputTemplate"`
	Truncation           string           `json:"truncation"`
	Tokenization         string           `json:"tokenization"`
	CandidateCount       int              `json:"candidateCount"`
	ResultsLimit         int              `json:"resultsLimit"`
	TimeoutMilliseconds  int64            `json:"timeoutMilliseconds"`
	Entries              []RerankingEntry `json:"entries"`
	DurationMilliseconds int64            `json:"durationMilliseconds"`
}
type GenerationTrace struct {
	Operator             OperatorRef        `json:"operator"`
	ModelManifestDigest  string             `json:"modelManifestDigest"`
	PromptManifestDigest string             `json:"promptManifestDigest"`
	Evidence             []EvidenceIdentity `json:"evidence"`
	InputArtifactDigest  string             `json:"inputArtifactDigest,omitempty"`
	OutputArtifactDigest string             `json:"outputArtifactDigest,omitempty"`
	FinishReason         string             `json:"finishReason"`
	CitationsValid       bool               `json:"citationsValid"`
	DurationMilliseconds int64              `json:"durationMilliseconds"`
}
type ResultScore struct {
	Fusion   *float64 `json:"fusion,omitempty"`
	Reranker *float64 `json:"reranker,omitempty"`
}
type MatchedRepresentation struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Channel string `json:"channel"`
	Rank    int    `json:"rank"`
}
type ResultTrace struct {
	Rank                   int                     `json:"rank"`
	Collapse               CollapseIdentity        `json:"collapse"`
	MatchedRepresentations []MatchedRepresentation `json:"matchedRepresentations"`
	Evidence               EvidenceIdentity        `json:"evidence"`
	Scores                 ResultScore             `json:"scores"`
}
type RelevanceTrace struct {
	Target      string                     `json:"target"`
	ExpectedIDs []string                   `json:"expectedIds"`
	Grades      map[string]float64         `json:"grades,omitempty"`
	Measures    map[string]json.RawMessage `json:"measures"`
}
type TimingTrace struct {
	TotalMilliseconds int64            `json:"totalMilliseconds"`
	ByOperator        map[string]int64 `json:"byOperator"`
}
type UsageTrace struct {
	InputTokens     int64              `json:"inputTokens"`
	OutputTokens    int64              `json:"outputTokens"`
	EmbeddingTokens int64              `json:"embeddingTokens"`
	ProviderCost    map[string]float64 `json:"providerCost,omitempty"`
}
type FailureTrace struct {
	Code       string          `json:"code"`
	Path       string          `json:"path"`
	Message    string          `json:"message"`
	OperatorID string          `json:"operatorId,omitempty"`
	Retryable  bool            `json:"retryable"`
	Details    json.RawMessage `json:"details,omitempty"`
}
