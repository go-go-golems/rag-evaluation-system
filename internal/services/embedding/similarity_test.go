package embedding

import (
	"context"
	"math"
	"testing"
)

func seedStoredSimilarityEmbeddings(t *testing.T) (strategyID string, chunkIDs []string, service *Service) {
	t.Helper()
	queries := openEmbeddingTestQueries(t)
	strategyID = seedChunkedStrategy(t, queries)
	chunks, err := queries.ListChunksForStrategy(strategyID, 3)
	if err != nil {
		t.Fatalf("list chunks: %v", err)
	}
	if len(chunks) < 3 {
		t.Fatalf("expected at least 3 chunks, got %d", len(chunks))
	}

	vectors := [][]float32{
		{1, 0},
		{0, 1},
		{1, 1},
	}
	for i, ch := range chunks[:3] {
		if err := queries.UpsertChunkEmbedding(ch.ID, strategyID, "fake", "unit", 2, TextHash(ch.Text), EncodeFloat32Vector(vectors[i])); err != nil {
			t.Fatalf("upsert embedding: %v", err)
		}
		chunkIDs = append(chunkIDs, ch.ID)
	}
	return strategyID, chunkIDs, NewService(queries)
}

func TestCosineSimilarity(t *testing.T) {
	score, err := CosineSimilarity([]float32{1, 0}, []float32{0, 1})
	if err != nil {
		t.Fatalf("cosine: %v", err)
	}
	if math.Abs(score) > 1e-9 {
		t.Fatalf("expected orthogonal vectors to score 0, got %v", score)
	}

	score, err = CosineSimilarity([]float32{1, 0}, []float32{1, 1})
	if err != nil {
		t.Fatalf("cosine: %v", err)
	}
	want := 1 / math.Sqrt2
	if math.Abs(score-want) > 1e-6 {
		t.Fatalf("expected %v, got %v", want, score)
	}
}

func TestCosineSimilarityRejectsInvalidVectors(t *testing.T) {
	if _, err := CosineSimilarity([]float32{1}, []float32{1, 2}); err == nil {
		t.Fatal("expected dimension mismatch to fail")
	}
	if _, err := CosineSimilarity([]float32{0, 0}, []float32{1, 2}); err == nil {
		t.Fatal("expected zero vector to fail")
	}
}

func TestSimilarityPairwise(t *testing.T) {
	ctx := context.Background()
	strategyID, chunkIDs, service := seedStoredSimilarityEmbeddings(t)

	result, err := service.Similarity(ctx, SimilarityRequest{
		StrategyID:   strategyID,
		ProviderType: "fake",
		Model:        "unit",
		Dimensions:   2,
		ChunkIDA:     chunkIDs[0],
		ChunkIDB:     chunkIDs[2],
		PreviewRunes: 12,
	})
	if err != nil {
		t.Fatalf("similarity: %v", err)
	}
	if result.Considered != 1 || len(result.Matches) != 1 {
		t.Fatalf("unexpected result shape: %#v", result)
	}
	want := 1 / math.Sqrt2
	if math.Abs(result.Matches[0].Score-want) > 1e-6 {
		t.Fatalf("expected %v, got %v", want, result.Matches[0].Score)
	}
	if result.Source.TextPreview == "" || result.Matches[0].TextPreview == "" {
		t.Fatalf("expected text previews in result: %#v", result)
	}
}

func TestSimilarityAgainstStrategyReturnsBoundedSortedMatches(t *testing.T) {
	ctx := context.Background()
	strategyID, chunkIDs, service := seedStoredSimilarityEmbeddings(t)

	result, err := service.Similarity(ctx, SimilarityRequest{
		StrategyID:     strategyID,
		ProviderType:   "fake",
		Model:          "unit",
		Dimensions:     2,
		ChunkIDA:       chunkIDs[0],
		Limit:          1,
		CandidateLimit: 10,
	})
	if err != nil {
		t.Fatalf("similarity: %v", err)
	}
	if result.Considered != 2 {
		t.Fatalf("expected to consider two non-source candidates, got %d", result.Considered)
	}
	if len(result.Matches) != 1 {
		t.Fatalf("expected one bounded match, got %d", len(result.Matches))
	}
	if result.Matches[0].ChunkID != chunkIDs[2] {
		t.Fatalf("expected diagonal chunk to rank first, got %s", result.Matches[0].ChunkID)
	}
}

func TestSimilarityReportsMissingEmbedding(t *testing.T) {
	ctx := context.Background()
	strategyID, _, service := seedStoredSimilarityEmbeddings(t)

	_, err := service.Similarity(ctx, SimilarityRequest{
		StrategyID:   strategyID,
		ProviderType: "fake",
		Model:        "unit",
		Dimensions:   2,
		ChunkIDA:     "missing",
	})
	if err == nil {
		t.Fatal("expected missing embedding to fail")
	}
}
