package search

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
)

const DefaultCandidateLimit = 500

type VectorQueryRequest struct {
	Query          string
	StrategyID     string
	SourceIDs      []string
	Provider       embeddings.Provider
	ProviderType   string
	Limit          int
	CandidateLimit int
	PreviewRunes   int
}

// QueryVector embeds the user query and compares it to stored chunk embeddings.
func (s *Service) QueryVector(ctx context.Context, req VectorQueryRequest) (*QueryResult, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	if req.StrategyID == "" {
		return nil, fmt.Errorf("strategy id is required")
	}
	if req.Provider == nil {
		return nil, fmt.Errorf("embedding provider is required")
	}
	if req.ProviderType == "" {
		req.ProviderType = "unknown"
	}
	if req.CandidateLimit <= 0 {
		req.CandidateLimit = DefaultCandidateLimit
	}
	if req.PreviewRunes == 0 {
		req.PreviewRunes = DefaultPreviewRunes
	}
	limit := normalizeLimit(req.Limit)
	sourceIDs := normalizeSourceIDs(req.SourceIDs)

	model := req.Provider.GetModel()
	if model.Name == "" || model.Dimensions <= 0 {
		return nil, fmt.Errorf("embedding provider returned invalid model metadata: %#v", model)
	}
	queryVector, err := req.Provider.GenerateEmbedding(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("generate query embedding: %w", err)
	}
	if len(queryVector) != model.Dimensions {
		return nil, fmt.Errorf("query embedding dimension mismatch: model=%d vector=%d", model.Dimensions, len(queryVector))
	}

	candidates, err := s.queries.ListChunkEmbeddingsForStrategySourcesWithContext(
		req.StrategyID, sourceIDs, req.ProviderType, model.Name, model.Dimensions, req.CandidateLimit,
	)
	if err != nil {
		return nil, err
	}

	items := make([]RetrievalResult, 0, len(candidates))
	for _, candidate := range candidates {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		vector, err := embeddingservice.DecodeFloat32Vector(candidate.Embedding)
		if err != nil {
			return nil, fmt.Errorf("decode embedding for chunk %s: %w", candidate.ChunkID, err)
		}
		if len(vector) != model.Dimensions {
			return nil, fmt.Errorf("stored embedding dimension mismatch for chunk %s: expected %d got %d", candidate.ChunkID, model.Dimensions, len(vector))
		}
		score, err := embeddingservice.CosineSimilarity(queryVector, vector)
		if err != nil {
			return nil, fmt.Errorf("compute vector similarity for chunk %s: %w", candidate.ChunkID, err)
		}
		items = append(items, RetrievalResult{
			ChunkID:    candidate.ChunkID,
			DocumentID: candidate.DocumentID,
			SourceID:   candidate.SourceID,
			Title:      candidate.Title,
			URL:        candidate.URL,
			StrategyID: candidate.StrategyID,
			ChunkIndex: candidate.ChunkIndex,
			Score:      score,
			Retriever:  "vector",
			Preview:    preview(candidate.Text, req.PreviewRunes),
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})
	if len(items) > limit {
		items = items[:limit]
	}
	for i := range items {
		items[i].Rank = i + 1
	}
	return &QueryResult{Query: req.Query, Retriever: "vector", Items: items}, nil
}
