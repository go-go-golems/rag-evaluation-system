package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

func TestDslDemoPageEndpointReturnsRenderableWidgetIR(t *testing.T) {
	handler := newTestAPIHandler(t)
	body := assertStatus(t, handler, "/api/v1/dsl/pages/demo", http.StatusOK)

	if body["id"] != "demo" {
		t.Fatalf("id = %#v, want demo", body["id"])
	}
	if body["title"] != "Widget IR Demo" {
		t.Fatalf("title = %#v", body["title"])
	}

	root, ok := body["root"].(map[string]any)
	if !ok {
		t.Fatalf("root = %#v, want object", body["root"])
	}
	if root["kind"] != "component" || root["type"] != "AppShell" {
		t.Fatalf("root = %#v, want AppShell component", root)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("response should be JSON serializable: %v", err)
	}
	if !json.Valid(encoded) {
		t.Fatalf("response marshaled to invalid JSON: %s", string(encoded))
	}
}

func TestDslDemoPageEndpointRejectsUnknownPage(t *testing.T) {
	handler := newTestAPIHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dsl/pages/unknown", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET unknown DSL page expected 404, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func newTestAPIHandler(t *testing.T) http.Handler {
	t.Helper()
	path := filepath.Join(t.TempDir(), "app.db")
	database, err := db.OpenDB(path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	if err := db.Migrate(database); err != nil {
		t.Fatalf("migrate db: %v", err)
	}

	mux := http.NewServeMux()
	RegisterHandlersWithOptions(mux, database, Options{DatabasePath: path, WorkflowDB: filepath.Join(t.TempDir(), "workflow.db"), WorkflowArtifactRoot: filepath.Join(t.TempDir(), "artifacts")})
	return mux
}

func TestDslDemoPageHandlerDirectPathValue(t *testing.T) {
	mux := http.NewServeMux()
	h := &handler{}
	mux.HandleFunc("GET /api/v1/dsl/pages/{id}", h.handleDslDemoPage)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/dsl/pages/demo", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("direct handler expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}
