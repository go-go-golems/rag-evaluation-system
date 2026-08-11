package immutablechunk

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

func TestBuildCreatesAndReusesExactChunkSet(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	queries := db.NewQueries(database)
	plan := &ttcimport.Plan{Manifest: ttcimport.Manifest{SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "fixture", SourceExportSHA256: "fixture", SelectionAlgorithm: "fixture", KindQuotas: map[string]int{"faq": 1}}, Documents: []ttcimport.SourceDocument{{ID: "wp:1", Kind: "faq", Title: "Fixture", SearchText: "  First sentence.  Second sentence.  ", SearchMarkdown: "# Fixture\nFirst sentence.\nSecond sentence."}}}
	snapshot, err := corpussnapshot.Persist(context.Background(), queries, plan, corpussnapshot.PersistRequest{SourceByteSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	first, err := Build(context.Background(), queries, Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "sentence", ChunkSize: 20, Overlap: 4})
	if err != nil {
		t.Fatal(err)
	}
	if first.Reused || first.ChunkCount == 0 {
		t.Fatalf("first result = %#v", first)
	}
	second, err := Build(context.Background(), queries, Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "sentence", ChunkSize: 20, Overlap: 4})
	if err != nil {
		t.Fatal(err)
	}
	if !second.Reused || second.ChunkSetID != first.ChunkSetID {
		t.Fatalf("second result = %#v", second)
	}
	var text string
	var start, end int
	if err := database.QueryRow(`SELECT text,source_start_runes,source_end_runes FROM immutable_chunks WHERE chunk_set_id=? ORDER BY chunk_index LIMIT 1`, first.ChunkSetID).Scan(&text, &start, &end); err != nil {
		t.Fatal(err)
	}
	if source := []rune(plan.Documents[0].SearchText); string(source[start:end]) != text {
		t.Fatalf("stored chunk %q != source range %q", text, string(source[start:end]))
	}
}
