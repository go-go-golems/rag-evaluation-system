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
func (r *Reranker) Rerank(ctx context.Context, req ragoperators.RerankRequest) ([]ragoperators.RerankScore, error) {
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

	response, err := r.provider.Rerank(ctx, geppettorerank.Request{
		Model:     req.Model,
		Query:     req.Query,
		Documents: documents,
		TopN:      len(documents),
	})
	if err != nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANK: %w", err)
	}
	if response.Model != req.Model {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_MODEL_MISMATCH")
	}
	if len(response.Results) != len(documents) {
		return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_INCOMPLETE")
	}

	scores := make([]ragoperators.RerankScore, 0, len(response.Results))
	seen := make(map[string]struct{}, len(response.Results))
	for _, result := range response.Results {
		if _, ok := expected[result.DocumentID]; !ok {
			return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_UNKNOWN_CHUNK_ID")
		}
		if _, duplicate := seen[result.DocumentID]; duplicate {
			return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_DUPLICATE_CHUNK_ID")
		}
		if math.IsNaN(result.Score) || math.IsInf(result.Score, 0) {
			return nil, fmt.Errorf("RAG_GEPPETTO_RERANK_NONFINITE_SCORE")
		}
		seen[result.DocumentID] = struct{}{}
		scores = append(scores, ragoperators.RerankScore{ChunkID: result.DocumentID, Score: result.Score})
	}
	return scores, nil
}
