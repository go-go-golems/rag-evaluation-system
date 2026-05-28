package embedding

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type ComputeRequest struct {
	StrategyID   string
	SourceIDs    []string
	Provider     embeddings.Provider
	ProviderType string
	BatchSize    int
	Limit        int
	Force        bool
}

type ComputeResult struct {
	StrategyID   string   `json:"strategy_id"`
	SourceIDs    []string `json:"source_ids,omitempty"`
	ProviderType string   `json:"provider_type"`
	Model        string   `json:"model"`
	Dimensions   int      `json:"dimensions"`
	Considered   int      `json:"considered"`
	Computed     int      `json:"computed"`
	SkippedFresh int      `json:"skipped_fresh"`
}

type CoverageRequest struct {
	StrategyID   string
	ProviderType string
	Model        string
	Dimensions   int
}

type CoverageResult struct {
	StrategyID   string                   `json:"strategy_id"`
	ProviderType string                   `json:"provider_type"`
	Model        string                   `json:"model"`
	Dimensions   int                      `json:"dimensions"`
	Items        []db.EmbeddingCoverage   `json:"items"`
	Totals       EmbeddingCoverageSummary `json:"totals"`
}

type EmbeddingCoverageSummary struct {
	ChunkCount    int `json:"chunk_count"`
	EmbeddedCount int `json:"embedded_count"`
	MissingCount  int `json:"missing_count"`
}

func (s *Service) Compute(ctx context.Context, req ComputeRequest) (*ComputeResult, error) {
	if req.StrategyID == "" {
		return nil, fmt.Errorf("strategy id is required")
	}
	if req.Provider == nil {
		return nil, fmt.Errorf("embedding provider is required")
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 16
	}
	if req.ProviderType == "" {
		req.ProviderType = "unknown"
	}

	model := req.Provider.GetModel()
	if model.Name == "" || model.Dimensions <= 0 {
		return nil, fmt.Errorf("embedding provider returned invalid model metadata: %#v", model)
	}

	chunks, err := s.queries.ListChunksForStrategySources(req.StrategyID, req.SourceIDs, req.Limit)
	if err != nil {
		return nil, err
	}
	result := &ComputeResult{
		StrategyID:   req.StrategyID,
		SourceIDs:    req.SourceIDs,
		ProviderType: req.ProviderType,
		Model:        model.Name,
		Dimensions:   model.Dimensions,
		Considered:   len(chunks),
	}

	type pendingChunk struct {
		chunk db.Chunk
		hash  string
	}
	pending := make([]pendingChunk, 0, len(chunks))
	for _, ch := range chunks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		hash := TextHash(ch.Text)
		storedHash, ok, err := s.queries.GetChunkEmbeddingTextHash(ch.ID, req.StrategyID, req.ProviderType, model.Name, model.Dimensions)
		if err != nil {
			return nil, err
		}
		if ok && storedHash == hash && !req.Force {
			result.SkippedFresh++
			continue
		}
		pending = append(pending, pendingChunk{chunk: ch, hash: hash})
	}

	for start := 0; start < len(pending); start += req.BatchSize {
		end := start + req.BatchSize
		if end > len(pending) {
			end = len(pending)
		}
		batch := pending[start:end]
		texts := make([]string, len(batch))
		for i, item := range batch {
			texts[i] = item.chunk.Text
		}

		vectors, err := req.Provider.GenerateBatchEmbeddings(ctx, texts)
		if err != nil {
			return nil, err
		}
		if len(vectors) != len(batch) {
			return nil, fmt.Errorf("embedding provider returned %d vectors for %d texts", len(vectors), len(batch))
		}

		for i, vector := range vectors {
			if len(vector) != model.Dimensions {
				return nil, fmt.Errorf("embedding dimension mismatch for chunk %s: model=%d vector=%d", batch[i].chunk.ID, model.Dimensions, len(vector))
			}
			if err := s.queries.UpsertChunkEmbedding(
				batch[i].chunk.ID,
				req.StrategyID,
				req.ProviderType,
				model.Name,
				model.Dimensions,
				batch[i].hash,
				EncodeFloat32Vector(vector),
			); err != nil {
				return nil, err
			}
			result.Computed++
		}
	}

	return result, nil
}

func (s *Service) Coverage(ctx context.Context, req CoverageRequest) (*CoverageResult, error) {
	if req.StrategyID == "" {
		return nil, fmt.Errorf("strategy id is required")
	}
	if req.ProviderType == "" {
		return nil, fmt.Errorf("provider type is required")
	}
	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if req.Dimensions <= 0 {
		return nil, fmt.Errorf("dimensions must be positive")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	items, err := s.queries.ListEmbeddingCoverageBySource(req.StrategyID, req.ProviderType, req.Model, req.Dimensions)
	if err != nil {
		return nil, err
	}
	result := &CoverageResult{
		StrategyID:   req.StrategyID,
		ProviderType: req.ProviderType,
		Model:        req.Model,
		Dimensions:   req.Dimensions,
		Items:        items,
	}
	for _, item := range items {
		result.Totals.ChunkCount += item.ChunkCount
		result.Totals.EmbeddedCount += item.EmbeddedCount
		result.Totals.MissingCount += item.MissingCount
	}
	return result, nil
}

func TextHash(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}
