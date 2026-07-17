package ragcontract

const TraceSchemaVersion = "rag-query-trace/v1"

type QueryTrace struct {
	SchemaVersion string          `json:"schemaVersion"`
	QueryCardID   string          `json:"queryCardId"`
	Query         string          `json:"query"`
	DatasetSplit  string          `json:"datasetSplit"`
	Channels      []ChannelTrace  `json:"channels"`
	Fusion        *FusionTrace    `json:"fusion,omitempty"`
	Reranking     *RerankingTrace `json:"reranking,omitempty"`
	Results       []Hit           `json:"results"`
	Relevance     RelevanceTrace  `json:"relevance"`
	Timing        TimingTrace     `json:"timing"`
}

type ChannelTrace struct {
	Name    string `json:"name"`
	Backend string `json:"backend"`
	Hits    []Hit  `json:"hits"`
}

type Hit struct {
	Rank               int     `json:"rank"`
	ChunkID            string  `json:"chunkId"`
	DocumentRevisionID string  `json:"documentRevisionId"`
	Score              float64 `json:"score"`
	Title              string  `json:"title,omitempty"`
	URL                string  `json:"url,omitempty"`
	Channel            string  `json:"channel,omitempty"`
}

type FusionTrace struct {
	Kind         string `json:"kind"`
	RankConstant int    `json:"rankConstant"`
	Hits         []Hit  `json:"hits"`
}

type RerankingTrace struct {
	Kind       string               `json:"kind"`
	Model      string               `json:"model"`
	Candidates []RerankingCandidate `json:"candidates"`
	Results    []RerankingResult    `json:"results"`
}

type RerankingCandidate struct {
	CandidateID    string  `json:"candidateId"`
	PreRerankRank  int     `json:"preRerankRank"`
	RetrievalScore float64 `json:"retrievalScore"`
}

type RerankingResult struct {
	CandidateID string  `json:"candidateId"`
	Rank        int     `json:"rank"`
	Score       float64 `json:"score"`
}

type RelevanceTrace struct {
	ExpectedDocumentRevisionIDs []string `json:"expectedDocumentRevisionIds"`
	FirstRelevantRank           int      `json:"firstRelevantRank"`
	RelevantDocumentRecall      float64  `json:"relevantDocumentRecallAtResults"`
}

type TimingTrace struct {
	EmbeddingMilliseconds int64 `json:"embeddingMilliseconds"`
	RetrievalMilliseconds int64 `json:"retrievalMilliseconds"`
	FusionMilliseconds    int64 `json:"fusionMilliseconds"`
	RerankingMilliseconds int64 `json:"rerankingMilliseconds"`
	TotalMilliseconds     int64 `json:"totalMilliseconds"`
}
