package raglab

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLlamaCPPRerankerMapsResponseToDurableCandidateIDs(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/v1/rerank" {
			t.Fatalf("unexpected request %s %s", request.Method, request.URL.Path)
		}
		if got := request.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("content type = %q, want application/json", got)
		}
		var body llamaCPPRerankRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "qllama/bge-reranker-v2-m3:q4_k_m" || body.Query != "refund policy" || body.TopN != 2 {
			t.Fatalf("unexpected request body: %#v", body)
		}
		if got, want := strings.Join(body.Documents, "|"), "first|second|third"; got != want {
			t.Fatalf("documents = %q, want %q", got, want)
		}
		_ = json.NewEncoder(w).Encode(llamaCPPRerankResponse{Results: []llamaCPPRerankItem{
			{Index: intPointer(2), RelevanceScore: floatPointer(-5)},
			{Index: intPointer(0), RelevanceScore: floatPointer(-1)},
		}})
	}))
	defer server.Close()

	reranker, err := NewLlamaCPPReranker(LlamaCPPRerankerOptions{BaseURL: server.URL, Model: "qllama/bge-reranker-v2-m3:q4_k_m", Client: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	results, err := reranker.Rerank(context.Background(), validRerankRequest())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("results length = %d, want 2", len(results))
	}
	if results[0].CandidateID != "candidate-0" || results[0].Rank != 1 || results[0].Score != -1 {
		t.Fatalf("first result = %#v", results[0])
	}
	if results[1].CandidateID != "candidate-2" || results[1].Rank != 2 || results[1].Score != -5 {
		t.Fatalf("second result = %#v", results[1])
	}
}

func TestLlamaCPPRerankerRejectsUnsafeNetworkConfiguration(t *testing.T) {
	for _, baseURL := range []string{
		"file:///tmp/socket",
		"http://user:provider-secret@example.test",
		"http://example.test?token=provider-secret",
		"http://example.test/#provider-secret",
	} {
		if _, err := NewLlamaCPPReranker(LlamaCPPRerankerOptions{BaseURL: baseURL, Model: "model"}); err == nil {
			t.Fatalf("unsafe base URL accepted: %s", baseURL)
		}
	}
}

func TestLlamaCPPRerankerRejectsOversizedRequest(t *testing.T) {
	t.Parallel()
	reranker, err := NewLlamaCPPReranker(LlamaCPPRerankerOptions{BaseURL: "http://example.test", Model: "model", MaxRequestBytes: 1})
	if err != nil {
		t.Fatal(err)
	}
	_, err = reranker.Rerank(context.Background(), validRerankRequest())
	if err == nil || !strings.Contains(err.Error(), "RAG_RERANKER_REQUEST_TOO_LARGE") {
		t.Fatalf("error = %v", err)
	}
}

func TestMapLlamaCPPResultsRejectsInvalidResponses(t *testing.T) {
	t.Parallel()
	candidates := validRerankRequest().Candidates
	cases := []struct {
		name  string
		items []llamaCPPRerankItem
		code  string
	}{
		{name: "wrong count", items: []llamaCPPRerankItem{{Index: intPointer(0), RelevanceScore: floatPointer(1)}}, code: "RAG_RERANK_RESPONSE_COUNT_INVALID"},
		{name: "missing index", items: []llamaCPPRerankItem{{RelevanceScore: floatPointer(1)}, {Index: intPointer(1), RelevanceScore: floatPointer(0)}}, code: "RAG_RERANK_RESPONSE_INVALID"},
		{name: "out of range", items: []llamaCPPRerankItem{{Index: intPointer(3), RelevanceScore: floatPointer(1)}, {Index: intPointer(1), RelevanceScore: floatPointer(0)}}, code: "RAG_RERANK_RESPONSE_INDEX_INVALID"},
		{name: "duplicate index", items: []llamaCPPRerankItem{{Index: intPointer(1), RelevanceScore: floatPointer(1)}, {Index: intPointer(1), RelevanceScore: floatPointer(0)}}, code: "RAG_RERANK_RESPONSE_INDEX_DUPLICATE"},
		{name: "non finite score", items: []llamaCPPRerankItem{{Index: intPointer(0), RelevanceScore: floatPointer(math.NaN())}, {Index: intPointer(1), RelevanceScore: floatPointer(0)}}, code: "RAG_RERANK_RESPONSE_SCORE_INVALID"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := mapLlamaCPPResults(candidates, 2, test.items)
			if err == nil || !strings.Contains(err.Error(), test.code) {
				t.Fatalf("error = %v, want %s", err, test.code)
			}
		})
	}
}

func validRerankRequest() RerankRequest {
	return RerankRequest{Query: "refund policy", TopN: 2, Candidates: []RerankCandidate{
		{ID: "candidate-0", Text: "first", OriginalRank: 1},
		{ID: "candidate-1", Text: "second", OriginalRank: 2},
		{ID: "candidate-2", Text: "third", OriginalRank: 3},
	}}
}

func intPointer(value int) *int { return &value }

func floatPointer(value float64) *float64 { return &value }
