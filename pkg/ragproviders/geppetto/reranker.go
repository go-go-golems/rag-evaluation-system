// Package geppetto adapts host-owned Geppetto services to the narrow RAG
// operator interfaces. It owns no endpoints, credentials, or provider
// transport policy; those remain in the host configuration used to construct
// the supplied Geppetto providers.
package geppetto

import (
	"context"
	"fmt"
	"math"

	geppettorerank "github.com/go-go-golems/geppetto/pkg/rerank"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
)

// Reranker adapts a configured Geppetto rerank provider to the RAG operator
// contract. RAG keeps durable chunk identity, source-evidence selection, and
// complete-score semantics; Geppetto owns provider transport and index mapping.
type Reranker struct {
	provider geppettorerank.Provider
}

var _ ragoperators.Reranker = (*Reranker)(nil)

// NewReranker constructs a RAG adapter around a host-configured Geppetto
// provider. Callers must not construct provider transport from RAG execution
// configuration or JavaScript authoring code.
func NewReranker(provider geppettorerank.Provider) (*Reranker, error) {
	if provider == nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANKER_PROVIDER_REQUIRED")
	}
	return &Reranker{provider: provider}, nil
}

// Rerank submits one document per hydrated source-evidence chunk. It always
// requests every candidate score: req.Results is the RAG display/final limit
// and is deliberately applied later by the native rerank operator.
type preparedRerank struct {
	provider geppettorerank.Provider
	request  geppettorerank.Request
	expected map[string]struct{}
}

func (r *Reranker) Rerank(ctx context.Context, req ragoperators.RerankRequest) (ragoperators.RerankResult, error) {
	prepared, err := r.PrepareRerank(req)
	if err != nil {
		return ragoperators.RerankResult{}, err
	}
	return prepared.Execute(ctx)
}

func (r *Reranker) PrepareRerank(req ragoperators.RerankRequest) (ragworkflowops.PreparedRerank, error) {
	if r == nil || r.provider == nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANKER_UNAVAILABLE")
	}
	if req.Model == "" {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_MODEL_REQUIRED")
	}
	if len(req.Candidates) == 0 {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_CANDIDATES_REQUIRED")
	}
	documents := make([]geppettorerank.Document, len(req.Candidates))
	expected := make(map[string]struct{}, len(req.Candidates))
	for i, candidate := range req.Candidates {
		id := candidate.Chunk.Record.ID
		if id == "" {
			return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_CHUNK_ID_REQUIRED")
		}
		if _, duplicate := expected[id]; duplicate {
			return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_DUPLICATE_CHUNK_ID")
		}
		expected[id] = struct{}{}
		documents[i] = geppettorerank.Document{ID: id, Text: candidate.Chunk.Text}
	}
	return &preparedRerank{provider: r.provider, request: geppettorerank.Request{Model: req.Model, Query: req.Query, Documents: documents, TopN: len(documents)}, expected: expected}, nil
}

func (p *preparedRerank) Execute(ctx context.Context) (ragoperators.RerankResult, error) {
	response, err := p.provider.Rerank(ctx, p.request)
	if err != nil {
		return ragoperators.RerankResult{}, classifyProviderError(err)
	}
	invalid := func(code string) (ragoperators.RerankResult, error) {
		result := ragoperators.RerankResult{Cost: response.Cost}
		if response.Usage != nil {
			result.InputTokens = int64(response.Usage.InputTokens)
		}
		return result, ragworkflowops.ProviderSucceededWithInvalidResult(code)
	}
	if response.Model != p.request.Model {
		return invalid("RAG_GEPPETTO_RERANK_MODEL_MISMATCH")
	}
	if len(response.Results) != len(p.request.Documents) {
		return invalid("RAG_GEPPETTO_RERANK_INCOMPLETE")
	}
	scores := make([]ragoperators.RerankScore, 0, len(response.Results))
	seen := make(map[string]struct{}, len(response.Results))
	for _, result := range response.Results {
		if _, ok := p.expected[result.DocumentID]; !ok {
			return invalid("RAG_GEPPETTO_RERANK_UNKNOWN_CHUNK_ID")
		}
		if _, duplicate := seen[result.DocumentID]; duplicate {
			return invalid("RAG_GEPPETTO_RERANK_DUPLICATE_CHUNK_ID")
		}
		if math.IsNaN(result.Score) || math.IsInf(result.Score, 0) {
			return invalid("RAG_GEPPETTO_RERANK_NONFINITE_SCORE")
		}
		seen[result.DocumentID] = struct{}{}
		scores = append(scores, ragoperators.RerankScore{ChunkID: result.DocumentID, Score: result.Score})
	}
	result := ragoperators.RerankResult{Scores: scores, Cost: response.Cost}
	if response.Usage != nil {
		result.InputTokens = int64(response.Usage.InputTokens)
	}
	return result, nil
}
