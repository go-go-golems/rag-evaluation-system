package corpus

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	tmp, err := os.CreateTemp("", "corpus-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	path := tmp.Name()
	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	database, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() })

	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}

	return database
}

func seedTestData(t *testing.T, database *sql.DB) {
	t.Helper()
	queries := db.NewQueries(database)
	ctx := context.Background()

	// Insert sources
	_ = queries.InsertSource("src-a", "Source A", "filesystem", "{}")
	_ = queries.InsertSource("src-b", "Source B", "sqlite-corpus", "{}")

	// Insert documents
	_ = queries.InsertDocument("doc-1", "src-a", "", "Document One", "", "", "text", "", "Hello world from doc one", "", 5, "en", "chunked")
	_ = queries.InsertDocument("doc-2", "src-a", "", "Document Two", "", "", "text", "", "Hello world from doc two", "", 5, "en", "chunked")
	_ = queries.InsertDocument("doc-3", "src-b", "", "Document Three", "", "https://example.com", "text", "", "Hello world from doc three", "", 4, "en", "extracted")

	// Insert chunking strategy
	_ = queries.InsertChunkingStrategy("fixed-100-20", "Fixed 100/20", "fixed", `{"chunk_size":100,"overlap":20}`, "Test strategy")

	// Insert chunks for doc-1
	_ = queries.InsertChunk("chk-001", "doc-1", "fixed-100-20", 0, "Hello world from doc one chunk zero", 25, 0, 100, "{}")
	_ = queries.InsertChunk("chk-002", "doc-1", "fixed-100-20", 1, "chunk one text here", 20, 80, 180, "{}")
	_ = queries.InsertChunk("chk-003", "doc-1", "fixed-100-20", 2, "chunk two text here", 20, 160, 260, "{}")

	// Insert chunks for doc-2
	_ = queries.InsertChunk("chk-004", "doc-2", "fixed-100-20", 0, "Hello world from doc two chunk zero", 25, 0, 100, "{}")

	// Insert embedding for chk-001 only
	if _, err := database.ExecContext(ctx, `
		INSERT INTO chunk_embeddings (chunk_id, strategy_id, provider, model, dimensions, text_hash, embedding)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"chk-001", "fixed-100-20", "openai", "text-embedding-3-small", 1536, "hash001", []byte{0x01, 0x02, 0x03}); err != nil {
		t.Fatal(err)
	}
}

func TestSourceSummaries_NoIdentity(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	results, err := svc.SourceSummaries(context.Background(), EmbeddingIdentity{})
	if err != nil {
		t.Fatalf("SourceSummaries: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 sources, got %d", len(results))
	}

	// Source A: 2 documents, 10 words
	if results[0].SourceID != "src-a" {
		t.Errorf("expected src-a, got %s", results[0].SourceID)
	}
	if results[0].DocumentCount != 2 {
		t.Errorf("expected 2 documents, got %d", results[0].DocumentCount)
	}
	if results[0].WordCount != 10 {
		t.Errorf("expected 10 words, got %d", results[0].WordCount)
	}

	// Source B: 1 document, 4 words
	if results[1].SourceID != "src-b" {
		t.Errorf("expected src-b, got %s", results[1].SourceID)
	}
	if results[1].DocumentCount != 1 {
		t.Errorf("expected 1 document, got %d", results[1].DocumentCount)
	}
}

func TestSourceSummaries_WithIdentity(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	results, err := svc.SourceSummaries(context.Background(), EmbeddingIdentity{
		StrategyID:   "fixed-100-20",
		ProviderType: "openai",
		Model:        "text-embedding-3-small",
		Dimensions:   1536,
	})
	if err != nil {
		t.Fatalf("SourceSummaries: %v", err)
	}

	// Source A: 4 chunks, 1 embedded
	srcA := results[0]
	if srcA.SourceID != "src-a" {
		t.Errorf("expected src-a, got %s", srcA.SourceID)
	}
	if srcA.ChunkCount != 4 {
		t.Errorf("expected 4 chunks, got %d", srcA.ChunkCount)
	}
	if srcA.EmbeddedCount != 1 {
		t.Errorf("expected 1 embedded, got %d", srcA.EmbeddedCount)
	}
	if srcA.MissingEmbeddingCount != 3 {
		t.Errorf("expected 3 missing, got %d", srcA.MissingEmbeddingCount)
	}

	// Source B: 0 chunks, 0 embedded
	srcB := results[1]
	if srcB.ChunkCount != 0 {
		t.Errorf("expected 0 chunks for src-b, got %d", srcB.ChunkCount)
	}
}

func TestDocumentBrowser(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	results, err := svc.DocumentBrowser(context.Background(), "src-a", EmbeddingIdentity{
		StrategyID:   "fixed-100-20",
		ProviderType: "openai",
		Model:        "text-embedding-3-small",
		Dimensions:   1536,
	}, 100, 0)
	if err != nil {
		t.Fatalf("DocumentBrowser: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(results))
	}

	// doc-1 has 3 chunks, 1 embedded
	d1 := results[0]
	if d1.ID != "doc-1" {
		t.Errorf("expected doc-1, got %s", d1.ID)
	}
	if d1.ChunkCount != 3 {
		t.Errorf("expected 3 chunks, got %d", d1.ChunkCount)
	}
	if d1.EmbeddedCount != 1 {
		t.Errorf("expected 1 embedded, got %d", d1.EmbeddedCount)
	}
}

func TestDocumentBrowser_Pagination(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	results, err := svc.DocumentBrowser(context.Background(), "src-a", EmbeddingIdentity{}, 1, 0)
	if err != nil {
		t.Fatalf("DocumentBrowser: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 document (page limit), got %d", len(results))
	}
}

func TestDocumentDetail(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	result, err := svc.DocumentDetail(context.Background(), "doc-1", EmbeddingIdentity{
		StrategyID:   "fixed-100-20",
		ProviderType: "openai",
		Model:        "text-embedding-3-small",
		Dimensions:   1536,
	}, true)
	if err != nil {
		t.Fatalf("DocumentDetail: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	doc := result.Document
	if doc.ID != "doc-1" {
		t.Errorf("expected doc-1, got %s", doc.ID)
	}
	if doc.Title != "Document One" {
		t.Errorf("expected 'Document One', got %s", doc.Title)
	}
	if doc.WordCount != 5 {
		t.Errorf("expected 5 words, got %d", doc.WordCount)
	}
	if doc.ContentText != "Hello world from doc one" {
		t.Errorf("unexpected content_text: %s", doc.ContentText)
	}

	if len(result.Chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(result.Chunks))
	}

	// Chunk 0 should have embedding
	c0 := result.Chunks[0]
	if c0.ChunkIndex != 0 {
		t.Errorf("expected chunk_index 0, got %d", c0.ChunkIndex)
	}
	if c0.Embedding == nil || !c0.Embedding.Present {
		t.Error("expected chunk 0 to have embedding")
	}

	// Chunk 1 should not have embedding
	c1 := result.Chunks[1]
	if c1.Embedding != nil && c1.Embedding.Present {
		t.Error("expected chunk 1 to not have embedding")
	}
}

func TestDocumentDetail_NotFound(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	result, err := svc.DocumentDetail(context.Background(), "nonexistent", EmbeddingIdentity{}, false)
	if err != nil {
		t.Fatalf("DocumentDetail: %v", err)
	}
	if result != nil {
		t.Error("expected nil result for nonexistent document")
	}
}

func TestDocumentDetail_NoText(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	result, err := svc.DocumentDetail(context.Background(), "doc-1", EmbeddingIdentity{
		StrategyID: "fixed-100-20",
	}, false)
	if err != nil {
		t.Fatalf("DocumentDetail: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// ContentText should be empty when includeText=false
	if result.Document.ContentText != "" {
		t.Errorf("expected empty content_text when includeText=false, got: %s", result.Document.ContentText)
	}
}

func TestDocumentDetail_NoStrategy(t *testing.T) {
	database := setupTestDB(t)
	seedTestData(t, database)

	svc := NewService(database)
	result, err := svc.DocumentDetail(context.Background(), "doc-1", EmbeddingIdentity{}, true)
	if err != nil {
		t.Fatalf("DocumentDetail: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// No strategy means no chunks fetched
	if len(result.Chunks) != 0 {
		t.Errorf("expected 0 chunks without strategy, got %d", len(result.Chunks))
	}
}
