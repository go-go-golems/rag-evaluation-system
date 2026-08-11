package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	workflowservice "github.com/go-go-golems/rag-evaluation-system/internal/workflow"
)

func TestWorkflowAndArtifactVisibilityEndpoints(t *testing.T) {
	ctx := context.Background()
	appDBPath := filepath.Join(t.TempDir(), "app.db")
	engineDBPath := filepath.Join(t.TempDir(), "engine.db")
	database, err := db.OpenDB(appDBPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	queries := db.NewQueries(database)
	seedAPIVisibilityData(t, queries)

	if _, err := workflowservice.SubmitIntakeWorkflow(ctx, workflowservice.SubmitIntakeRequest{
		EngineDB:            engineDBPath,
		DBPath:              appDBPath,
		WorkflowID:          "wf-api-visibility",
		DocumentIDs:         []string{"doc-1"},
		Strategy:            "fixed",
		ChunkSize:           20,
		Overlap:             5,
		SkipPreprocessing:   true,
		SkipEmbeddings:      true,
		SkipBM25:            true,
		SkipChunkEnrichment: true,
	}); err != nil {
		t.Fatalf("submit workflow: %v", err)
	}

	mux := http.NewServeMux()
	RegisterHandlersWithOptions(mux, database, Options{EngineDB: engineDBPath})

	assertStatus(t, mux, "/api/v1/workflows", http.StatusOK)
	assertStatus(t, mux, "/api/v1/workflows/wf-api-visibility", http.StatusOK)
	assertStatus(t, mux, "/api/v1/workflows/wf-api-visibility/ops", http.StatusOK)
	assertStatus(t, mux, "/api/v1/artifacts/document-processing/coverage?artifact_type=clean_text&prompt_version=v1&provider=fake&model=fake-document-processor", http.StatusOK)
	assertStatus(t, mux, "/api/v1/documents/doc-1/processing-artifacts", http.StatusOK)
	assertStatus(t, mux, "/api/v1/artifacts/chunk-enrichment/coverage?strategy_id=fixed-20-5&prompt_version=v1", http.StatusOK)
	assertStatus(t, mux, "/api/v1/chunks/chunk-1/enrichments?strategy_id=fixed-20-5", http.StatusOK)
}

func assertStatus(t *testing.T, handler http.Handler, path string, status int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != status {
		t.Fatalf("GET %s expected %d, got %d body=%s", path, status, rec.Code, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response for %s: %v body=%s", path, err, rec.Body.String())
	}
	return body
}

func seedAPIVisibilityData(t *testing.T, queries *db.Queries) {
	t.Helper()
	if err := queries.InsertSource("source-1", "Source", "test", "{}"); err != nil {
		t.Fatalf("insert source: %v", err)
	}
	if err := queries.InsertDocument("doc-1", "source-1", "doc-1", "Title", "", "", "text", "Raw", "Visibility test content.", "", 3, "en", "extracted"); err != nil {
		t.Fatalf("insert document: %v", err)
	}
	if err := queries.InsertChunkingStrategy("fixed-20-5", "fixed", "fixed", "{}", "test"); err != nil {
		t.Fatalf("insert strategy: %v", err)
	}
	if _, err := queries.DB().Exec(`INSERT INTO chunks (id, document_id, strategy_id, chunk_index, text, token_count) VALUES ('chunk-1', 'doc-1', 'fixed-20-5', 0, 'Visibility test content.', 3)`); err != nil {
		t.Fatalf("insert chunk: %v", err)
	}
	if err := queries.UpsertDocumentProcessingArtifact(db.DocumentProcessingArtifact{DocumentID: "doc-1", ArtifactType: "clean_text", PromptVersion: "v1", Provider: "fake", Model: "fake-document-processor", InputHash: "hash", OutputText: "clean", OutputJSON: "{}", Status: "succeeded"}); err != nil {
		t.Fatalf("upsert processing artifact: %v", err)
	}
	if err := queries.UpsertChunkEnrichment(db.ChunkEnrichment{ChunkID: "chunk-1", StrategyID: "fixed-20-5", PromptVersion: "v1", Provider: "fake", Model: "fake-chunk-enricher", ShortSummary: "summary", LongSummary: "long", KeyTopicsJSON: "[]", EntitiesJSON: "[]", HypotheticalQuestionsJSON: "[]", QualityScore: 0.9, TextHash: "hash"}); err != nil {
		t.Fatalf("upsert chunk enrichment: %v", err)
	}
}
