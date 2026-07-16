package raglab

import "context"

// RerankCandidate is a hydrated retrieval candidate supplied to a cross
// encoder. ID is the durable identity; Index is assigned only by the adapter
// when it constructs the remote document array.
type RerankCandidate struct {
	ID             string  `json:"id"`
	Text           string  `json:"text"`
	OriginalRank   int     `json:"original_rank"`
	RetrievalScore float64 `json:"retrieval_score"`
}

type RerankRequest struct {
	Query      string            `json:"query"`
	Candidates []RerankCandidate `json:"candidates"`
	TopN       int               `json:"top_n"`
}

type RerankResult struct {
	CandidateID string  `json:"candidate_id"`
	Index       int     `json:"index"`
	Score       float64 `json:"score"`
	Rank        int     `json:"rank"`
}

type RerankerIdentity struct {
	Kind  string `json:"kind"`
	Model string `json:"model"`
}

// Reranker is intentionally transport-neutral. Implementations must return
// results mapped to durable candidate IDs rather than exposing provider array
// positions as application identity.
type Reranker interface {
	Rerank(ctx context.Context, request RerankRequest) ([]RerankResult, error)
	Identity() RerankerIdentity
}
