package embedding

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

type SimilarityRequest struct {
	StrategyID     string
	ProviderType   string
	Model          string
	Dimensions     int
	ChunkIDA       string
	ChunkIDB       string
	Limit          int
	CandidateLimit int
	PreviewRunes   int
}

type SimilarityResult struct {
	StrategyID     string            `json:"strategy_id"`
	ProviderType   string            `json:"provider_type"`
	Model          string            `json:"model"`
	Dimensions     int               `json:"dimensions"`
	Source         SimilarityChunk   `json:"source"`
	Matches        []SimilarityMatch `json:"matches"`
	Considered     int               `json:"considered"`
	CandidateLimit int               `json:"candidate_limit"`
}

type SimilarityChunk struct {
	ChunkID     string `json:"chunk_id"`
	DocumentID  string `json:"document_id"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	TextPreview string `json:"text_preview,omitempty"`
}

type SimilarityMatch struct {
	SimilarityChunk
	Score float64 `json:"score"`
}

func (s *Service) Similarity(ctx context.Context, req SimilarityRequest) (*SimilarityResult, error) {
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
	if req.ChunkIDA == "" {
		return nil, fmt.Errorf("chunk-id-a is required")
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.CandidateLimit <= 0 {
		req.CandidateLimit = 200
	}
	if req.PreviewRunes < 0 {
		req.PreviewRunes = 0
	}

	source, sourceVector, err := s.loadEmbeddingVector(req.ChunkIDA, req)
	if err != nil {
		return nil, err
	}

	result := &SimilarityResult{
		StrategyID:     req.StrategyID,
		ProviderType:   req.ProviderType,
		Model:          req.Model,
		Dimensions:     req.Dimensions,
		Source:         embeddingToSimilarityChunk(source, req.PreviewRunes),
		CandidateLimit: req.CandidateLimit,
	}

	if req.ChunkIDB != "" {
		target, targetVector, err := s.loadEmbeddingVector(req.ChunkIDB, req)
		if err != nil {
			return nil, err
		}
		score, err := CosineSimilarity(sourceVector, targetVector)
		if err != nil {
			return nil, fmt.Errorf("compute similarity: %w", err)
		}
		result.Considered = 1
		result.Matches = []SimilarityMatch{{SimilarityChunk: embeddingToSimilarityChunk(target, req.PreviewRunes), Score: score}}
		return result, nil
	}

	candidates, err := s.queries.ListChunkEmbeddingsForStrategy(req.StrategyID, req.ProviderType, req.Model, req.Dimensions, req.CandidateLimit)
	if err != nil {
		return nil, err
	}

	matches := make([]SimilarityMatch, 0, len(candidates))
	for _, candidate := range candidates {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if candidate.ChunkID == source.ChunkID {
			continue
		}
		vector, err := DecodeFloat32Vector(candidate.Embedding)
		if err != nil {
			return nil, fmt.Errorf("decode embedding for chunk %s: %w", candidate.ChunkID, err)
		}
		if len(vector) != req.Dimensions {
			return nil, fmt.Errorf("stored embedding dimension mismatch for chunk %s: expected %d got %d", candidate.ChunkID, req.Dimensions, len(vector))
		}
		score, err := CosineSimilarity(sourceVector, vector)
		if err != nil {
			return nil, fmt.Errorf("compute similarity for chunk %s: %w", candidate.ChunkID, err)
		}
		matches = append(matches, SimilarityMatch{SimilarityChunk: embeddingToSimilarityChunk(&candidate, req.PreviewRunes), Score: score})
	}

	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
	result.Considered = len(matches)
	if len(matches) > req.Limit {
		matches = matches[:req.Limit]
	}
	result.Matches = matches
	return result, nil
}

func (s *Service) loadEmbeddingVector(chunkID string, req SimilarityRequest) (*db.ChunkEmbedding, []float32, error) {
	embedding, ok, err := s.queries.GetChunkEmbedding(chunkID, req.StrategyID, req.ProviderType, req.Model, req.Dimensions)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, fmt.Errorf("embedding not found for chunk %s strategy=%s provider=%s model=%s dimensions=%d", chunkID, req.StrategyID, req.ProviderType, req.Model, req.Dimensions)
	}
	vector, err := DecodeFloat32Vector(embedding.Embedding)
	if err != nil {
		return nil, nil, fmt.Errorf("decode embedding for chunk %s: %w", chunkID, err)
	}
	if len(vector) != req.Dimensions {
		return nil, nil, fmt.Errorf("stored embedding dimension mismatch for chunk %s: expected %d got %d", chunkID, req.Dimensions, len(vector))
	}
	return embedding, vector, nil
}

func embeddingToSimilarityChunk(e *db.ChunkEmbedding, previewRunes int) SimilarityChunk {
	return SimilarityChunk{
		ChunkID:     e.ChunkID,
		DocumentID:  e.DocumentID,
		StrategyID:  e.StrategyID,
		ChunkIndex:  e.ChunkIndex,
		TextPreview: truncateRunes(e.Text, previewRunes),
	}
}

func CosineSimilarity(a, b []float32) (float64, error) {
	if len(a) == 0 || len(b) == 0 {
		return 0, fmt.Errorf("vectors must be non-empty")
	}
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector dimensions differ: %d != %d", len(a), len(b))
	}

	var dot, normA, normB float64
	for i := range a {
		av := float64(a[i])
		bv := float64(b[i])
		dot += av * bv
		normA += av * av
		normB += bv * bv
	}
	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("cosine similarity is undefined for zero vectors")
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}
