package ttcimport

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func TestBuildPlanIsDeterministicAndIncludesSeeds(t *testing.T) {
	sourcePath := createSourceFixture(t, []fixtureDocument{
		{ID: "wp:3699", Kind: "product", Title: "Seed product"},
		{ID: "wp:product-a", Kind: "product", Title: "Product A"},
		{ID: "wp:product-b", Kind: "product", Title: "Product B"},
		{ID: "wp:812290", Kind: "ttc_guide", Title: "Seed guide"},
		{ID: "wp:guide-a", Kind: "ttc_guide", Title: "Guide A"},
		{ID: "wp:4131", Kind: "faq", Title: "Seed FAQ"},
		{ID: "wp:faq-a", Kind: "faq", Title: "FAQ A"},
		{ID: "wp:627148", Kind: "post", Title: "Seed post"},
		{ID: "wp:post-a", Kind: "post", Title: "Post A"},
		{ID: "wp:4237", Kind: "page", Title: "Seed page"},
		{ID: "wp:page-a", Kind: "page", Title: "Page A"},
	})

	request := BuildRequest{
		SourceDBPath: sourcePath,
		SnapshotName: "fixture-v1",
		KindQuotas: map[string]int{
			"ttc_guide": 2,
			"faq":       2,
			"post":      2,
			"product":   2,
			"page":      2,
		},
		AdditionalSeedDocumentIDs: []string{"wp:product-a"},
	}
	planOne, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("BuildPlan first run: %v", err)
	}
	planTwo, err := BuildPlan(context.Background(), request)
	if err != nil {
		t.Fatalf("BuildPlan second run: %v", err)
	}

	if len(planOne.Documents) != 10 {
		t.Fatalf("selected documents = %d, want 10", len(planOne.Documents))
	}
	oneJSON, err := json.Marshal(planOne.Manifest)
	if err != nil {
		t.Fatal(err)
	}
	twoJSON, err := json.Marshal(planTwo.Manifest)
	if err != nil {
		t.Fatal(err)
	}
	if string(oneJSON) != string(twoJSON) {
		t.Fatalf("manifest differs across identical runs\nfirst: %s\nsecond: %s", oneJSON, twoJSON)
	}
	assertSelected(t, planOne, "wp:3699")
	assertSelected(t, planOne, "wp:812290")
	assertSelected(t, planOne, "wp:4131")
	assertSelected(t, planOne, "wp:627148")
	assertSelected(t, planOne, "wp:4237")
	assertSelected(t, planOne, "wp:product-a")
}

func TestPersistWritesManifestAndOperationalDocuments(t *testing.T) {
	sourcePath := createSourceFixture(t, []fixtureDocument{
		{ID: "wp:3699", Kind: "product", Title: "Seed product"},
		{ID: "wp:812290", Kind: "ttc_guide", Title: "Seed guide"},
		{ID: "wp:4131", Kind: "faq", Title: "Seed FAQ"},
		{ID: "wp:627148", Kind: "post", Title: "Seed post"},
		{ID: "wp:4237", Kind: "page", Title: "Seed page"},
	})
	plan, err := BuildPlan(context.Background(), BuildRequest{
		SourceDBPath: sourcePath,
		SnapshotName: "fixture-v1",
		KindQuotas: map[string]int{
			"ttc_guide": 1,
			"faq":       1,
			"post":      1,
			"product":   1,
			"page":      1,
		},
	})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}

	targetPath := filepath.Join(t.TempDir(), "target.db")
	target, err := db.OpenDB(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = target.Close() })
	if err := db.Migrate(target); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	result, err := Persist(context.Background(), db.NewQueries(target), plan, PersistRequest{
		ManifestPath: manifestPath,
	})
	if err != nil {
		t.Fatalf("Persist: %v", err)
	}
	if result.DocumentCount != 5 || result.DryRun {
		t.Fatalf("unexpected result: %#v", result)
	}
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest not written: %v", err)
	}
	var count int
	if err := target.QueryRow(`SELECT COUNT(*) FROM documents WHERE source_id = ?`, DefaultSourceID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Fatalf("imported documents = %d, want 5", count)
	}
	var externalID, metadata string
	if err := target.QueryRow(`SELECT external_id, metadata_json FROM documents WHERE id = ?`, "ttc:wp:3699").Scan(&externalID, &metadata); err != nil {
		t.Fatal(err)
	}
	if externalID != "wp:3699" {
		t.Fatalf("external ID = %q, want wp:3699", externalID)
	}
	if !strings.Contains(metadata, "source_search_text_sha256") {
		t.Fatalf("metadata does not preserve source hashes: %s", metadata)
	}
}

func TestBuildPlanRejectsMissingSeed(t *testing.T) {
	sourcePath := createSourceFixture(t, []fixtureDocument{
		{ID: "wp:3699", Kind: "product", Title: "Seed product"},
		{ID: "wp:812290", Kind: "ttc_guide", Title: "Seed guide"},
		{ID: "wp:4131", Kind: "faq", Title: "Seed FAQ"},
		{ID: "wp:627148", Kind: "post", Title: "Seed post"},
		{ID: "wp:4237", Kind: "page", Title: "Seed page"},
	})
	_, err := BuildPlan(context.Background(), BuildRequest{
		SourceDBPath: sourcePath,
		KindQuotas: map[string]int{
			"ttc_guide": 1,
			"faq":       1,
			"post":      1,
			"product":   1,
			"page":      1,
		},
		AdditionalSeedDocumentIDs: []string{"wp:not-present"},
	})
	if err == nil {
		t.Fatal("BuildPlan succeeded with missing seed")
	}
}

type fixtureDocument struct {
	ID    string
	Kind  string
	Title string
}

func createSourceFixture(t *testing.T, documents []fixtureDocument) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "ttc-source.db")
	source, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = source.Close() })
	if _, err := source.Exec(`
		CREATE TABLE documents (
			doc_id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			search_text TEXT NOT NULL,
			search_markdown TEXT NOT NULL,
			content_markdown TEXT NOT NULL
		)`); err != nil {
		t.Fatal(err)
	}
	for _, document := range documents {
		text := document.Title + " searchable content"
		if _, err := source.Exec(
			`INSERT INTO documents (doc_id, kind, title, url, search_text, search_markdown, content_markdown) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			document.ID, document.Kind, document.Title, "https://example.test/"+document.ID, text, "# "+document.Title, "# "+document.Title,
		); err != nil {
			t.Fatal(err)
		}
	}
	return path
}

func assertSelected(t *testing.T, plan *Plan, documentID string) {
	t.Helper()
	for _, document := range plan.Documents {
		if document.ID == documentID {
			return
		}
	}
	t.Fatalf("expected selected document %q", documentID)
}
