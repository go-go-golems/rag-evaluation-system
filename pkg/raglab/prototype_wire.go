package raglab

// Prototype wire values are isolated from canonical v2 contracts and deleted with pkg/raglab.
const (
	PrototypeSchemaVersion      = "rag-retrieval-spec/v1"
	PrototypeTraceSchemaVersion = "rag-query-trace/v1"
)

type PrototypeSpecification struct {
	SchemaVersion   string                    `json:"schemaVersion"`
	Name            string                    `json:"name"`
	Dataset         PrototypeDatasetSelection `json:"dataset"`
	Representations []PrototypeRepresentation `json:"representations"`
	Retrieval       PrototypeRetrievalPlan    `json:"retrieval"`
	Metrics         PrototypeMetricPlan       `json:"metrics"`
	Tags            map[string]string         `json:"tags,omitempty"`
}

type PrototypeDatasetSelection struct {
	Split string `json:"split"`
}

type PrototypeRepresentation struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type PrototypeRetrievalPlan struct {
	Filter    PrototypeFilterSpec     `json:"filter,omitempty"`
	Channels  []PrototypeChannel      `json:"channels"`
	Fusion    *PrototypeFusionPlan    `json:"fusion,omitempty"`
	Reranking *PrototypeRerankingPlan `json:"reranking,omitempty"`
	Collapse  string                  `json:"collapse"`
	Results   int                     `json:"results"`
}

type PrototypeChannel struct {
	Name           string              `json:"name"`
	Backend        string              `json:"backend"`
	Representation string              `json:"representation"`
	TopK           int                 `json:"topK"`
	Filter         PrototypeFilterSpec `json:"filter,omitempty"`
}

type PrototypeFilterSpec struct {
	SourceIDs      []string          `json:"sourceIds,omitempty"`
	DocumentIDs    []string          `json:"documentIds,omitempty"`
	ContentTypes   []string          `json:"contentTypes,omitempty"`
	MetadataEquals map[string]string `json:"metadataEquals,omitempty"`
}

type PrototypeFusionPlan struct {
	Kind         string             `json:"kind"`
	RankConstant int                `json:"rankConstant"`
	Weights      map[string]float64 `json:"weights,omitempty"`
}

type PrototypeRerankingPlan struct {
	Kind           string `json:"kind"`
	Model          string `json:"model"`
	CandidateCount int    `json:"candidateCount"`
	Results        int    `json:"results"`
}

type PrototypeMetricPlan struct {
	RelevanceAt string `json:"relevanceAt"`
	RecallAt    []int  `json:"recallAt,omitempty"`
	PrecisionAt []int  `json:"precisionAt,omitempty"`
	NDCGAt      []int  `json:"ndcgAt,omitempty"`
	MRR         bool   `json:"mrr"`
}

type PrototypeQueryTrace struct {
	SchemaVersion string                   `json:"schemaVersion"`
	QueryCardID   string                   `json:"queryCardId"`
	Query         string                   `json:"query"`
	DatasetSplit  string                   `json:"datasetSplit"`
	Channels      []PrototypeChannelTrace  `json:"channels"`
	Fusion        *PrototypeFusionTrace    `json:"fusion,omitempty"`
	Reranking     *PrototypeRerankingTrace `json:"reranking,omitempty"`
	Results       []PrototypeHit           `json:"results"`
	Relevance     PrototypeRelevanceTrace  `json:"relevance"`
	Timing        PrototypeTimingTrace     `json:"timing"`
}

type PrototypeChannelTrace struct {
	Name    string         `json:"name"`
	Backend string         `json:"backend"`
	Hits    []PrototypeHit `json:"hits"`
}

type PrototypeHit struct {
	Rank               int     `json:"rank"`
	ChunkID            string  `json:"chunkId"`
	DocumentRevisionID string  `json:"documentRevisionId"`
	Score              float64 `json:"score"`
	Title              string  `json:"title,omitempty"`
	URL                string  `json:"url,omitempty"`
	Channel            string  `json:"channel,omitempty"`
}

type PrototypeFusionTrace struct {
	Kind         string         `json:"kind"`
	RankConstant int            `json:"rankConstant"`
	Hits         []PrototypeHit `json:"hits"`
}

type PrototypeRerankingTrace struct {
	Kind       string                        `json:"kind"`
	Model      string                        `json:"model"`
	Candidates []PrototypeRerankingCandidate `json:"candidates"`
	Results    []PrototypeRerankingResult    `json:"results"`
}

type PrototypeRerankingCandidate struct {
	CandidateID    string  `json:"candidateId"`
	PreRerankRank  int     `json:"preRerankRank"`
	RetrievalScore float64 `json:"retrievalScore"`
}

type PrototypeRerankingResult struct {
	CandidateID string  `json:"candidateId"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
}

type PrototypeRelevanceTrace struct {
	ExpectedDocumentRevisionIDs []string `json:"expectedDocumentRevisionIds"`
	FirstRelevantRank           int      `json:"firstRelevantRank"`
	RelevantDocumentRecall      float64  `json:"relevantDocumentRecallAtResults"`
}

type PrototypeTimingTrace struct {
	EmbeddingMilliseconds int64 `json:"embeddingMilliseconds"`
	RetrievalMilliseconds int64 `json:"retrievalMilliseconds"`
	FusionMilliseconds    int64 `json:"fusionMilliseconds"`
	RerankingMilliseconds int64 `json:"rerankingMilliseconds"`
	TotalMilliseconds     int64 `json:"totalMilliseconds"`
}
