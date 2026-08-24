package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragintakeworkflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

func TestWorkflowAndArtifactVisibilityEndpoints(t *testing.T) {
	ctx := context.Background()
	appDBPath := filepath.Join(t.TempDir(), "app.db")
	workflowDBPath := filepath.Join(t.TempDir(), "workflow.db")
	artifactRoot := filepath.Join(t.TempDir(), "artifacts")
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

	config := ragintakeworkflow.DefaultConfig(appDBPath)
	config.WorkflowDatabase, config.ArtifactRoot = workflowDBPath, artifactRoot
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		t.Fatalf("open intake: %v", err)
	}
	request, err := ragintakeworkflow.PrepareRequest(ctx, config.Runtime, ragintakeworkflow.Request{IndexID: "unused", SkipPreprocessing: true, SkipEmbeddings: true, SkipBM25: true, SkipChunkEnrichment: true}, ragintakeworkflow.Selection{DocumentIDs: []string{"doc-1"}})
	if err != nil {
		t.Fatalf("prepare intake: %v", err)
	}
	if _, err = app.SubmitRequest(ctx, request, workflowv3.RunID("wf-api-visibility")); err != nil {
		t.Fatalf("submit intake: %v", err)
	}
	if _, err = app.RunUntilTerminal(ctx, workflowv3.RunID("wf-api-visibility")); err != nil {
		t.Fatalf("run intake: %v", err)
	}
	if err = app.Close(); err != nil {
		t.Fatalf("close intake: %v", err)
	}

	mux := http.NewServeMux()
	RegisterHandlersWithOptions(mux, database, Options{DatabasePath: appDBPath, WorkflowDB: workflowDBPath, WorkflowArtifactRoot: artifactRoot})

	assertStatus(t, mux, "/api/v1/intake/runs", http.StatusOK)
	assertStatus(t, mux, "/api/v1/intake/runs/wf-api-visibility", http.StatusOK)
	assertStatus(t, mux, "/api/v1/intake/runs/wf-api-visibility/observations", http.StatusOK)
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

func TestApplyIntakeEmbeddingDefaultsRestoresProviderSettings(t *testing.T) {
	// Omitted fields fall back to the conventional defaults when embeddings run.
	input := &intakeSubmitRequest{SkipEmbeddings: false}
	applyIntakeEmbeddingDefaults(input)
	if input.EmbeddingType != "ollama" || input.EmbeddingEngine != "nomic-embed-text" || input.Dimensions != 768 {
		t.Fatalf("defaults not restored: type=%q engine=%q dimensions=%d", input.EmbeddingType, input.EmbeddingEngine, input.Dimensions)
	}

	// Explicit values are preserved.
	explicit := &intakeSubmitRequest{SkipEmbeddings: false, EmbeddingType: "openai", EmbeddingEngine: "text-embedding-3-small", Dimensions: 1536}
	applyIntakeEmbeddingDefaults(explicit)
	if explicit.EmbeddingType != "openai" || explicit.EmbeddingEngine != "text-embedding-3-small" || explicit.Dimensions != 1536 {
		t.Fatalf("explicit values overwritten: type=%q engine=%q dimensions=%d", explicit.EmbeddingType, explicit.EmbeddingEngine, explicit.Dimensions)
	}

	// Skipping embeddings leaves the provider fields untouched.
	skipped := &intakeSubmitRequest{SkipEmbeddings: true}
	applyIntakeEmbeddingDefaults(skipped)
	if skipped.EmbeddingType != "" || skipped.EmbeddingEngine != "" || skipped.Dimensions != 0 {
		t.Fatalf("skipped embeddings should not be defaulted: type=%q engine=%q dimensions=%d", skipped.EmbeddingType, skipped.EmbeddingEngine, skipped.Dimensions)
	}
}
