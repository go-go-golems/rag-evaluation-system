package embedding

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	geppettoembeddings "github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	sourceservice "github.com/go-go-golems/rag-evaluation-system/internal/services/source"
)

type fakeProvider struct {
	dimensions int
	calls      int
}

func (f *fakeProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	vectors, err := f.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vectors[0], nil
}

func (f *fakeProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	f.calls++
	vectors := make([][]float32, len(texts))
	for i, text := range texts {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		vectors[i] = make([]float32, f.dimensions)
		for j := range vectors[i] {
			vectors[i][j] = float32(len(text) + j)
		}
	}
	return vectors, nil
}

func (f *fakeProvider) GetModel() geppettoembeddings.EmbeddingModel {
	return geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: f.dimensions}
}

type badDimensionProvider struct{ fakeProvider }

func (f *badDimensionProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i := range texts {
		vectors[i] = []float32{1}
	}
	return vectors, nil
}

func openEmbeddingTestQueries(t *testing.T) *db.Queries {
	t.Helper()
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return db.NewQueries(database)
}

func seedChunkedStrategy(t *testing.T, queries *db.Queries) string {
	t.Helper()
	ctx := context.Background()
	sourceSvc := sourceservice.NewService(queries)
	if _, err := sourceSvc.Create(ctx, sourceservice.CreateRequest{ID: "docs", Name: "Docs", Type: "filesystem"}); err != nil {
		t.Fatalf("create source: %v", err)
	}
	dir := t.TempDir()
	for i := 0; i < 2; i++ {
		path := filepath.Join(dir, fmt.Sprintf("doc-%d.md", i))
		if err := os.WriteFile(path, []byte("# Doc\nThis document has enough text to produce chunks for embeddings."), 0644); err != nil {
			t.Fatalf("write doc: %v", err)
		}
	}
	if _, err := sourceSvc.Scan(ctx, sourceservice.ScanRequest{SourceID: "docs", Dir: dir}); err != nil {
		t.Fatalf("scan source: %v", err)
	}
	docs, err := queries.ListDocuments(10, 0)
	if err != nil {
		t.Fatalf("list docs: %v", err)
	}
	chunkSvc := chunkservice.NewService(queries)
	var strategyID string
	for _, doc := range docs {
		result, err := chunkSvc.Apply(ctx, chunkservice.ApplyRequest{DocumentID: doc.ID, Strategy: "fixed", ChunkSize: 30, Overlap: 5})
		if err != nil {
			t.Fatalf("apply chunks: %v", err)
		}
		strategyID = result.StrategyID
	}
	return strategyID
}

func TestComputeEmbeddingsSkipsFreshRows(t *testing.T) {
	ctx := context.Background()
	queries := openEmbeddingTestQueries(t)
	strategyID := seedChunkedStrategy(t, queries)
	provider := &fakeProvider{dimensions: 4}
	service := NewService(queries)

	first, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: provider, ProviderType: "fake", BatchSize: 2})
	if err != nil {
		t.Fatalf("first compute: %v", err)
	}
	if first.Considered == 0 || first.Computed != first.Considered || first.SkippedFresh != 0 {
		t.Fatalf("unexpected first result: %#v", first)
	}

	second, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: provider, ProviderType: "fake", BatchSize: 2})
	if err != nil {
		t.Fatalf("second compute: %v", err)
	}
	if second.Computed != 0 || second.SkippedFresh != first.Considered {
		t.Fatalf("expected all fresh rows skipped, got %#v", second)
	}

	var stored int
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM chunk_embeddings WHERE strategy_id = ?`, strategyID).Scan(&stored); err != nil {
		t.Fatalf("count embeddings: %v", err)
	}
	if stored != first.Considered {
		t.Fatalf("expected %d embedding rows, got %d", first.Considered, stored)
	}
}

func TestComputeEmbeddingsForceRecomputesFreshRows(t *testing.T) {
	ctx := context.Background()
	queries := openEmbeddingTestQueries(t)
	strategyID := seedChunkedStrategy(t, queries)
	provider := &fakeProvider{dimensions: 4}
	service := NewService(queries)

	first, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: provider, ProviderType: "fake", BatchSize: 2})
	if err != nil {
		t.Fatalf("first compute: %v", err)
	}
	forced, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: provider, ProviderType: "fake", BatchSize: 2, Force: true})
	if err != nil {
		t.Fatalf("forced compute: %v", err)
	}
	if forced.Computed != first.Considered || forced.SkippedFresh != 0 {
		t.Fatalf("expected force to recompute all rows, got %#v", forced)
	}
}

func TestComputeEmbeddingsRejectsDimensionMismatch(t *testing.T) {
	ctx := context.Background()
	queries := openEmbeddingTestQueries(t)
	strategyID := seedChunkedStrategy(t, queries)
	service := NewService(queries)

	_, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: &badDimensionProvider{fakeProvider{dimensions: 4}}, ProviderType: "fake"})
	if err == nil {
		t.Fatal("expected dimension mismatch to fail")
	}
}

func TestComputeEmbeddingsCanFilterBySourceID(t *testing.T) {
	ctx := context.Background()
	queries := openEmbeddingTestQueries(t)
	if err := queries.InsertSource("a", "A", "test", "{}"); err != nil {
		t.Fatalf("insert source a: %v", err)
	}
	if err := queries.InsertSource("b", "B", "test", "{}"); err != nil {
		t.Fatalf("insert source b: %v", err)
	}
	if err := queries.InsertDocument("doc-a", "a", "a.md", "A", "", "", "text", "alpha text", "alpha text", "", 2, "en", "extracted"); err != nil {
		t.Fatalf("insert doc a: %v", err)
	}
	if err := queries.InsertDocument("doc-b", "b", "b.md", "B", "", "", "text", "beta text", "beta text", "", 2, "en", "extracted"); err != nil {
		t.Fatalf("insert doc b: %v", err)
	}
	if err := queries.InsertChunkingStrategy("fixed-test", "fixed-test", "fixed", "{}", "test"); err != nil {
		t.Fatalf("insert strategy: %v", err)
	}
	if err := queries.InsertChunk("chunk-a", "doc-a", "fixed-test", 0, "alpha text", 2, 0, 10, "{}"); err != nil {
		t.Fatalf("insert chunk a: %v", err)
	}
	if err := queries.InsertChunk("chunk-b", "doc-b", "fixed-test", 0, "beta text", 2, 0, 9, "{}"); err != nil {
		t.Fatalf("insert chunk b: %v", err)
	}

	provider := &fakeProvider{dimensions: 4}
	service := NewService(queries)
	result, err := service.Compute(ctx, ComputeRequest{StrategyID: "fixed-test", SourceIDs: []string{"b"}, Provider: provider, ProviderType: "fake"})
	if err != nil {
		t.Fatalf("compute: %v", err)
	}
	if result.Considered != 1 || result.Computed != 1 {
		t.Fatalf("expected only source b to be considered/computed, got %#v", result)
	}
	if _, ok, err := queries.GetChunkEmbeddingTextHash("chunk-a", "fixed-test", "fake", "fake-embedding", 4); err != nil || ok {
		t.Fatalf("expected chunk-a to remain unembedded, ok=%v err=%v", ok, err)
	}
	if _, ok, err := queries.GetChunkEmbeddingTextHash("chunk-b", "fixed-test", "fake", "fake-embedding", 4); err != nil || !ok {
		t.Fatalf("expected chunk-b embedding, ok=%v err=%v", ok, err)
	}
}

func TestEmbeddingCoverageBySource(t *testing.T) {
	ctx := context.Background()
	queries := openEmbeddingTestQueries(t)
	strategyID := seedChunkedStrategy(t, queries)
	provider := &fakeProvider{dimensions: 4}
	service := NewService(queries)
	if _, err := service.Compute(ctx, ComputeRequest{StrategyID: strategyID, Provider: provider, ProviderType: "fake", Limit: 1}); err != nil {
		t.Fatalf("compute one: %v", err)
	}
	coverage, err := service.Coverage(ctx, CoverageRequest{StrategyID: strategyID, ProviderType: "fake", Model: "fake-embedding", Dimensions: 4})
	if err != nil {
		t.Fatalf("coverage: %v", err)
	}
	if coverage.Totals.ChunkCount == 0 || coverage.Totals.EmbeddedCount != 1 || coverage.Totals.MissingCount != coverage.Totals.ChunkCount-1 {
		t.Fatalf("unexpected coverage: %#v", coverage)
	}
}
