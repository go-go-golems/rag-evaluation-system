package raglab

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const defaultLlamaCPPMaxRequestBytes = 2 << 20

// LlamaCPPRerankerOptions configures operational capability. None of these
// fields are persisted in RerankingSpec: the immutable specification names the
// selected model and candidate policy, while the caller supplies transport.
type LlamaCPPRerankerOptions struct {
	BaseURL         string
	Model           string
	Client          *http.Client
	MaxRequestBytes int
}

type LlamaCPPReranker struct {
	baseURL         string
	model           string
	client          *http.Client
	maxRequestBytes int
}

var _ Reranker = (*LlamaCPPReranker)(nil)

func NewLlamaCPPReranker(options LlamaCPPRerankerOptions) (*LlamaCPPReranker, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(options.BaseURL), "/")
	if baseURL == "" {
		return nil, errors.New("RAG_RERANKER_BASE_URL_REQUIRED: llama.cpp base URL is required")
	}
	if strings.TrimSpace(options.Model) == "" {
		return nil, errors.New("RAG_RERANKER_MODEL_REQUIRED: llama.cpp model is required")
	}
	maxRequestBytes := options.MaxRequestBytes
	if maxRequestBytes == 0 {
		maxRequestBytes = defaultLlamaCPPMaxRequestBytes
	}
	if maxRequestBytes < 1 {
		return nil, errors.New("RAG_RERANKER_REQUEST_LIMIT_INVALID: max request bytes must be positive")
	}
	client := options.Client
	if client == nil {
		client = http.DefaultClient
	}
	return &LlamaCPPReranker{baseURL: baseURL, model: strings.TrimSpace(options.Model), client: client, maxRequestBytes: maxRequestBytes}, nil
}

func (r *LlamaCPPReranker) Identity() RerankerIdentity {
	return RerankerIdentity{Kind: "llama.cpp", Model: r.model}
}

func (r *LlamaCPPReranker) Rerank(ctx context.Context, request RerankRequest) ([]RerankResult, error) {
	if r == nil || r.client == nil {
		return nil, errors.New("RAG_RERANKER_REQUIRED: llama.cpp reranker is required")
	}
	if err := validateRerankRequest(request); err != nil {
		return nil, err
	}
	documents := make([]string, len(request.Candidates))
	for i, candidate := range request.Candidates {
		documents[i] = candidate.Text
	}
	payload, err := json.Marshal(llamaCPPRerankRequest{Model: r.model, Query: request.Query, Documents: documents, TopN: request.TopN})
	if err != nil {
		return nil, errors.Wrap(err, "encode llama.cpp reranking request")
	}
	if len(payload) > r.maxRequestBytes {
		return nil, errors.Errorf("RAG_RERANKER_REQUEST_TOO_LARGE: encoded request is %d bytes, limit is %d", len(payload), r.maxRequestBytes)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/v1/rerank", bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(err, "create llama.cpp reranking request")
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(httpRequest)
	if err != nil {
		return nil, errors.Wrap(err, "call llama.cpp reranking endpoint")
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.Errorf("RAG_RERANKER_HTTP_STATUS: llama.cpp reranking endpoint returned %s", response.Status)
	}
	var decoded llamaCPPRerankResponse
	decoder := json.NewDecoder(http.MaxBytesReader(nil, response.Body, int64(r.maxRequestBytes)))
	if err := decoder.Decode(&decoded); err != nil {
		return nil, errors.Wrap(err, "decode llama.cpp reranking response")
	}
	return mapLlamaCPPResults(request.Candidates, request.TopN, decoded.Results)
}

type llamaCPPRerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n"`
}

type llamaCPPRerankResponse struct {
	Results []llamaCPPRerankItem `json:"results"`
}

type llamaCPPRerankItem struct {
	Index          *int     `json:"index"`
	RelevanceScore *float64 `json:"relevance_score"`
}

func validateRerankRequest(request RerankRequest) error {
	if strings.TrimSpace(request.Query) == "" {
		return errors.New("RAG_RERANK_QUERY_REQUIRED: query is required")
	}
	if len(request.Candidates) == 0 {
		return errors.New("RAG_RERANK_CANDIDATES_REQUIRED: at least one candidate is required")
	}
	if request.TopN < 1 || request.TopN > len(request.Candidates) {
		return errors.Errorf("RAG_RERANK_TOP_N_INVALID: top_n must be between 1 and %d", len(request.Candidates))
	}
	seen := map[string]struct{}{}
	for i, candidate := range request.Candidates {
		if strings.TrimSpace(candidate.ID) == "" || strings.TrimSpace(candidate.Text) == "" || candidate.OriginalRank < 1 {
			return errors.Errorf("RAG_RERANK_CANDIDATE_INVALID: candidate %d requires ID, text, and positive original rank", i)
		}
		if _, exists := seen[candidate.ID]; exists {
			return errors.Errorf("RAG_RERANK_CANDIDATE_DUPLICATE: duplicate candidate ID %q", candidate.ID)
		}
		seen[candidate.ID] = struct{}{}
	}
	return nil
}

func mapLlamaCPPResults(candidates []RerankCandidate, topN int, items []llamaCPPRerankItem) ([]RerankResult, error) {
	if len(items) != topN {
		return nil, errors.Errorf("RAG_RERANK_RESPONSE_COUNT_INVALID: got %d results, want %d", len(items), topN)
	}
	seen := map[int]struct{}{}
	results := make([]RerankResult, 0, len(items))
	for _, item := range items {
		if item.Index == nil || item.RelevanceScore == nil {
			return nil, errors.New("RAG_RERANK_RESPONSE_INVALID: result requires index and relevance_score")
		}
		if *item.Index < 0 || *item.Index >= len(candidates) {
			return nil, errors.Errorf("RAG_RERANK_RESPONSE_INDEX_INVALID: index %d is outside submitted candidates", *item.Index)
		}
		if _, exists := seen[*item.Index]; exists {
			return nil, errors.Errorf("RAG_RERANK_RESPONSE_INDEX_DUPLICATE: index %d appears more than once", *item.Index)
		}
		if math.IsNaN(*item.RelevanceScore) || math.IsInf(*item.RelevanceScore, 0) {
			return nil, errors.Errorf("RAG_RERANK_RESPONSE_SCORE_INVALID: index %d has non-finite relevance score", *item.Index)
		}
		seen[*item.Index] = struct{}{}
		results = append(results, RerankResult{CandidateID: candidates[*item.Index].ID, Index: *item.Index, Score: *item.RelevanceScore})
	}
	// Response order is not treated as a durable ranking guarantee. Sorting is
	// deterministic and preserves original retrieval order for equal scores.
	sortRerankResults(results, candidates)
	return results, nil
}

func sortRerankResults(results []RerankResult, candidates []RerankCandidate) {
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		left, right := candidates[results[i].Index], candidates[results[j].Index]
		if left.OriginalRank != right.OriginalRank {
			return left.OriginalRank < right.OriginalRank
		}
		return left.ID < right.ID
	})
	for i := range results {
		results[i].Rank = i + 1
	}
}
