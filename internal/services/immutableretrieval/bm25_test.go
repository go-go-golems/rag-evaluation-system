package immutableretrieval

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutablechunk"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

func TestBuildBM25ReusesArtifactAndHydratesRankedEvidence(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	queries := db.NewQueries(database)
	plan := &ttcimport.Plan{Manifest: ttcimport.Manifest{
		SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "fixture", SourceExportSHA256: "fixture", SelectionAlgorithm: "fixture", KindQuotas: map[string]int{"faq": 2},
	}, Documents: []ttcimport.SourceDocument{
		{ID: "wp:alpha", Kind: "faq", Title: "Blue Cypress", URL: "https://example.test/alpha", SearchText: "Blue Cypress is very drought resistant."},
		{ID: "wp:beta", Kind: "faq", Title: "Green Plant", URL: "https://example.test/beta", SearchText: "Green plants require regular watering."},
	}}
	snapshot, err := corpussnapshot.Persist(context.Background(), queries, plan, corpussnapshot.PersistRequest{SourceByteSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	chunks, err := immutablechunk.Build(context.Background(), queries, immutablechunk.Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "fixed", ChunkSize: 200})
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(t.TempDir(), "bm25")
	first, err := BuildBM25(context.Background(), queries, BM25BuildRequest{ChunkSetID: chunks.ChunkSetID, ArtifactRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	if first.Reused || first.ChunkCount != 2 {
		t.Fatalf("first build = %#v", first)
	}
	second, err := BuildBM25(context.Background(), queries, BM25BuildRequest{ChunkSetID: chunks.ChunkSetID, ArtifactRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	if !second.Reused || second.ArtifactID != first.ArtifactID || second.Path != first.Path {
		t.Fatalf("second build = %#v, first = %#v", second, first)
	}
	hits, err := QueryBM25(context.Background(), queries, first.ArtifactID, "blue cypress drought", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) == 0 {
		t.Fatal("expected lexical hit")
	}
	got := hits[0]
	if got.DocumentRevisionID == "" || got.Title != "Blue Cypress" || got.URL != "https://example.test/alpha" || got.Text != "Blue Cypress is very drought resistant." {
		t.Fatalf("hydrated hit = %#v", got)
	}
	if got.Rank != 1 || got.Channel != "bm25" || got.StartRunes != 0 || got.EndRunes != len([]rune(got.Text)) {
		t.Fatalf("hit rank/citation fields = %#v", got)
	}
}
