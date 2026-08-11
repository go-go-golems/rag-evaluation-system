package corpussnapshot

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

func TestPersistCreatesAndReusesImmutableSnapshot(t *testing.T) {
	t.Parallel()
	queries := testQueries(t)
	plan := fixturePlan("alpha")

	first, err := Persist(context.Background(), queries, plan, PersistRequest{SourceByteSize: 123})
	if err != nil {
		t.Fatalf("persist first snapshot: %v", err)
	}
	if first.Reused {
		t.Fatal("first snapshot was unexpectedly reused")
	}
	if first.DocumentCount != 2 {
		t.Fatalf("document count = %d, want 2", first.DocumentCount)
	}
	second, err := Persist(context.Background(), queries, plan, PersistRequest{SourceByteSize: 123})
	if err != nil {
		t.Fatalf("persist same snapshot: %v", err)
	}
	if !second.Reused {
		t.Fatal("second identical snapshot was not reused")
	}
	if second.SnapshotID != first.SnapshotID {
		t.Fatalf("snapshot id changed: %s != %s", second.SnapshotID, first.SnapshotID)
	}

	var artifactCount, revisionCount, snapshotCount, membershipCount int
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM source_artifacts`).Scan(&artifactCount); err != nil {
		t.Fatal(err)
	}
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM document_revisions`).Scan(&revisionCount); err != nil {
		t.Fatal(err)
	}
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM corpus_snapshots`).Scan(&snapshotCount); err != nil {
		t.Fatal(err)
	}
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM corpus_snapshot_documents`).Scan(&membershipCount); err != nil {
		t.Fatal(err)
	}
	if artifactCount != 1 || revisionCount != 2 || snapshotCount != 1 || membershipCount != 2 {
		t.Fatalf("immutable counts = artifacts %d revisions %d snapshots %d memberships %d", artifactCount, revisionCount, snapshotCount, membershipCount)
	}
}

func TestPersistContentChangeCreatesNewRevisionAndSnapshot(t *testing.T) {
	t.Parallel()
	queries := testQueries(t)
	first, err := Persist(context.Background(), queries, fixturePlan("alpha"), PersistRequest{SourceByteSize: 123})
	if err != nil {
		t.Fatalf("persist first snapshot: %v", err)
	}
	second, err := Persist(context.Background(), queries, fixturePlan("changed alpha"), PersistRequest{SourceByteSize: 123})
	if err != nil {
		t.Fatalf("persist changed snapshot: %v", err)
	}
	if first.SnapshotID == second.SnapshotID {
		t.Fatal("content change did not create a new snapshot")
	}
	var revisions int
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM document_revisions`).Scan(&revisions); err != nil {
		t.Fatal(err)
	}
	if revisions != 3 {
		t.Fatalf("revisions = %d, want 3", revisions)
	}
}

func TestPersistRejectsPhysicalArtifactConflict(t *testing.T) {
	t.Parallel()
	queries := testQueries(t)
	plan := fixturePlan("alpha")
	if _, err := Persist(context.Background(), queries, plan, PersistRequest{SourceByteSize: 123}); err != nil {
		t.Fatalf("persist first snapshot: %v", err)
	}
	if _, err := Persist(context.Background(), queries, plan, PersistRequest{SourceByteSize: 456}); err == nil {
		t.Fatal("persist accepted a source artifact whose physical byte size conflicts")
	}
}

func testQueries(t *testing.T) *db.Queries {
	t.Helper()
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate database: %v", err)
	}
	return db.NewQueries(database)
}

func fixturePlan(firstText string) *ttcimport.Plan {
	documents := []ttcimport.SourceDocument{
		{ID: "wp:1", Kind: "faq", Title: "Alpha", URL: "https://example.test/alpha", SearchText: firstText, SearchMarkdown: "# Alpha\n" + firstText, ContentMarkdown: "# Alpha\n" + firstText},
		{ID: "wp:2", Kind: "product", Title: "Beta", URL: "https://example.test/beta", SearchText: "beta", SearchMarkdown: "# Beta\nbeta", ContentMarkdown: "# Beta\nbeta"},
	}
	return &ttcimport.Plan{Manifest: ttcimport.Manifest{SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "fixture", SourceExportSHA256: "abcd", SelectionAlgorithm: "fixture", KindQuotas: map[string]int{"faq": 1, "product": 1}, SeedDocumentIDs: []string{"wp:1"}}, Documents: documents}
}
