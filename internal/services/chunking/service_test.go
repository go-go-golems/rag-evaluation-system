package chunking

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	sourceservice "github.com/go-go-golems/rag-evaluation-system/internal/services/source"
)

func openTestQueries(t *testing.T) *db.Queries {
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

func seedDocument(t *testing.T, queries *db.Queries) string {
	t.Helper()
	ctx := context.Background()
	sourceSvc := sourceservice.NewService(queries)
	if _, err := sourceSvc.Create(ctx, sourceservice.CreateRequest{ID: "docs", Name: "Docs", Type: "filesystem"}); err != nil {
		t.Fatalf("create source: %v", err)
	}
	dir := t.TempDir()
	text := "# Test\n" + "abcdefghijklmnopqrstuvwxyz " + "abcdefghijklmnopqrstuvwxyz " + "abcdefghijklmnopqrstuvwxyz " + "abcdefghijklmnopqrstuvwxyz"
	if err := os.WriteFile(filepath.Join(dir, "doc.md"), []byte(text), 0644); err != nil {
		t.Fatalf("write doc: %v", err)
	}
	if _, err := sourceSvc.Scan(ctx, sourceservice.ScanRequest{SourceID: "docs", Dir: dir}); err != nil {
		t.Fatalf("scan source: %v", err)
	}
	docs, err := queries.ListDocuments(10, 0)
	if err != nil {
		t.Fatalf("list docs: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	return docs[0].ID
}

func TestApplyIsRerunSafeForSameStrategy(t *testing.T) {
	ctx := context.Background()
	queries := openTestQueries(t)
	docID := seedDocument(t, queries)
	service := NewService(queries)

	first, err := service.Apply(ctx, ApplyRequest{DocumentID: docID, Strategy: "fixed", ChunkSize: 40, Overlap: 10})
	if err != nil {
		t.Fatalf("first apply: %v", err)
	}
	second, err := service.Apply(ctx, ApplyRequest{DocumentID: docID, Strategy: "fixed", ChunkSize: 40, Overlap: 10})
	if err != nil {
		t.Fatalf("second apply: %v", err)
	}
	if first.ChunkCount == 0 || first.ChunkCount != second.ChunkCount {
		t.Fatalf("expected stable non-zero chunk count, got %d and %d", first.ChunkCount, second.ChunkCount)
	}

	var stored int
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM chunks WHERE document_id = ? AND strategy_id = ?`, docID, first.StrategyID).Scan(&stored); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if stored != first.ChunkCount {
		t.Fatalf("expected %d stored chunks after rerun, got %d", first.ChunkCount, stored)
	}
}

func TestApplyKeepsStrategiesSeparate(t *testing.T) {
	ctx := context.Background()
	queries := openTestQueries(t)
	docID := seedDocument(t, queries)
	service := NewService(queries)

	fixed, err := service.Apply(ctx, ApplyRequest{DocumentID: docID, Strategy: "fixed", ChunkSize: 40, Overlap: 10})
	if err != nil {
		t.Fatalf("fixed apply: %v", err)
	}
	sentence, err := service.Apply(ctx, ApplyRequest{DocumentID: docID, Strategy: "sentence", ChunkSize: 80, Overlap: 10})
	if err != nil {
		t.Fatalf("sentence apply: %v", err)
	}
	if fixed.StrategyID == sentence.StrategyID {
		t.Fatalf("expected distinct strategy IDs, got %q", fixed.StrategyID)
	}

	rows, err := queries.DB().Query(`SELECT strategy_id, COUNT(*) FROM chunks GROUP BY strategy_id`)
	if err != nil {
		t.Fatalf("group chunks: %v", err)
	}
	defer func() { _ = rows.Close() }()
	counts := map[string]int{}
	for rows.Next() {
		var strategyID string
		var count int
		if err := rows.Scan(&strategyID, &count); err != nil {
			t.Fatalf("scan group: %v", err)
		}
		counts[strategyID] = count
	}
	if counts[fixed.StrategyID] != fixed.ChunkCount {
		t.Fatalf("fixed chunks not isolated: counts=%#v result=%d", counts, fixed.ChunkCount)
	}
	if counts[sentence.StrategyID] != sentence.ChunkCount {
		t.Fatalf("sentence chunks not isolated: counts=%#v result=%d", counts, sentence.ChunkCount)
	}
}
