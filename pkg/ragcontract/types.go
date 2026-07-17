// Package ragcontract defines the pure researchctl RAG retrieval wire payload.
// It has no database, provider, Goja, filesystem, or execution dependencies.
package ragcontract

const SchemaVersion = "rag-retrieval-spec/v1"

type Specification struct {
	SchemaVersion   string            `json:"schemaVersion"`
	Name            string            `json:"name"`
	Dataset         DatasetSelection  `json:"dataset"`
	Representations []Representation  `json:"representations"`
	Retrieval       RetrievalPlan     `json:"retrieval"`
	Metrics         MetricPlan        `json:"metrics"`
	Tags            map[string]string `json:"tags,omitempty"`
}

type DatasetSelection struct {
	Split string `json:"split"`
}

type Representation struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type RetrievalPlan struct {
	Filter    FilterSpec     `json:"filter,omitempty"`
	Channels  []Channel      `json:"channels"`
	Fusion    *FusionPlan    `json:"fusion,omitempty"`
	Reranking *RerankingPlan `json:"reranking,omitempty"`
	Collapse  string         `json:"collapse"`
	Results   int            `json:"results"`
}

type Channel struct {
	Name           string     `json:"name"`
	Backend        string     `json:"backend"`
	Representation string     `json:"representation"`
	TopK           int        `json:"topK"`
	Filter         FilterSpec `json:"filter,omitempty"`
}

type FilterSpec struct {
	SourceIDs      []string          `json:"sourceIds,omitempty"`
	DocumentIDs    []string          `json:"documentIds,omitempty"`
	ContentTypes   []string          `json:"contentTypes,omitempty"`
	MetadataEquals map[string]string `json:"metadataEquals,omitempty"`
}

type FusionPlan struct {
	Kind         string             `json:"kind"`
	RankConstant int                `json:"rankConstant"`
	Weights      map[string]float64 `json:"weights,omitempty"`
}

type RerankingPlan struct {
	Kind           string `json:"kind"`
	Model          string `json:"model"`
	CandidateCount int    `json:"candidateCount"`
	Results        int    `json:"results"`
}

type MetricPlan struct {
	RelevanceAt string `json:"relevanceAt"`
	RecallAt    []int  `json:"recallAt,omitempty"`
	PrecisionAt []int  `json:"precisionAt,omitempty"`
	NDCGAt      []int  `json:"ndcgAt,omitempty"`
	MRR         bool   `json:"mrr"`
}
