package widgetserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetrunner"
)

func TestPageEndpointReturnsRunnerPage(t *testing.T) {
	handler := newTestServer(t, `
		const rag = require("widget.dsl");
		exports.pages = {
			demo(ctx) {
				return { id: ctx.pageId, title: "Demo " + ctx.query.q, root: rag.panel({ title: "Demo" }, "ok") };
			}
		};
	`).Handler()

	body := getJSON(t, handler, "/api/widget/pages/demo?q=World", http.StatusOK)
	if body["id"] != "demo" || body["title"] != "Demo World" {
		t.Fatalf("page body = %#v", body)
	}
	root := body["root"].(map[string]any)
	if root["kind"] != "component" || root["type"] != "Panel" {
		t.Fatalf("root = %#v", root)
	}
}

func TestPageEndpointMapsMissingPageTo404(t *testing.T) {
	handler := newTestServer(t, `exports.pages = {};`).Handler()
	body := getJSON(t, handler, "/api/widget/pages/missing", http.StatusNotFound)
	errorBody := body["error"].(map[string]any)
	if errorBody["code"] != "page_not_found" {
		t.Fatalf("error = %#v", errorBody)
	}
}

func TestSchemaAndHealthEndpoints(t *testing.T) {
	handler := newTestServer(t, `exports.pages = {};`).Handler()
	health := getJSON(t, handler, "/api/widget/health", http.StatusOK)
	if health["status"] != "ok" {
		t.Fatalf("health = %#v", health)
	}
	schema := getJSON(t, handler, "/api/widget/schema", http.StatusOK)
	if schema["schemaVersion"] != "0.1.0" {
		t.Fatalf("schema = %#v", schema)
	}
	components := schema["components"].([]any)
	if len(components) == 0 || len(schema["cellKinds"].([]any)) == 0 {
		t.Fatalf("schema lists are empty: %#v", schema)
	}
	for _, component := range []string{"Text", "ContextKindSwatch", "SectionBlock", "SplitPane", "SidebarShell", "SlideShell", "ContextDiagramPanel", "ContextBudgetBar", "TranscriptWorkspacePanel", "AnnotationRailPanel", "AnchoredCommentRail", "CourseStudioShell", "CourseSlidePanel", "HandoutDocumentShell"} {
		if !containsString(components, component) {
			t.Fatalf("schema components missing %s: %#v", component, components)
		}
	}
	jsonSchema, ok := schema["jsonSchema"].(map[string]any)
	if !ok || jsonSchema["$schema"] == "" {
		t.Fatalf("jsonSchema missing: %#v", schema)
	}
}

func TestActionEndpointInvokesRunnerAction(t *testing.T) {
	handler := newTestServer(t, `
		exports.actions = {
			save(ctx, payload) {
				return { ok: true, refresh: true, toast: "saved " + payload.id, data: { action: ctx.action, rowKey: ctx.rowKey } };
			}
		};
	`).Handler()
	body := postJSONBody(t, handler, "/api/widget/actions/save", `{"payload":{"id":"42"},"context":{"rowKey":"row-42"}}`, http.StatusOK)
	if body["ok"] != true || body["refresh"] != true || body["toast"] != "saved 42" {
		t.Fatalf("action body = %#v", body)
	}
	data := body["data"].(map[string]any)
	if data["action"] != "save" || data["rowKey"] != "row-42" {
		t.Fatalf("data = %#v", data)
	}
}

func TestActionEndpointMapsMissingActionTo404(t *testing.T) {
	handler := newTestServer(t, `exports.actions = {};`).Handler()
	body := postJSON(t, handler, "/api/widget/actions/save", http.StatusNotFound)
	errorBody := body["error"].(map[string]any)
	if errorBody["code"] != "action_not_found" {
		t.Fatalf("error = %#v", errorBody)
	}
}

func TestAPIOnlyDoesNotServeFrontend(t *testing.T) {
	handler := newTestServer(t, `exports.pages = {};`).Handler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET / = %d, want 404", rec.Code)
	}
}

func TestEmbeddedFrontendServesDefaultSPA(t *testing.T) {
	runner := newTestRunner(t, `exports.pages = {};`)
	server, err := New(Config{Runner: runner, FrontendMode: FrontendEmbedded})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	body := getText(t, server.Handler(), "/pages/demo", http.StatusOK)
	if !strings.Contains(body, `id="root"`) {
		t.Fatalf("embedded fallback did not return app index: %s", body)
	}
}

func TestEmbeddedFrontendAllowsCustomHandlerOverride(t *testing.T) {
	runner := newTestRunner(t, `exports.pages = {};`)
	server, err := New(Config{
		Runner:       runner,
		FrontendMode: FrontendEmbedded,
		FrontendHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("custom"))
		}),
	})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	body := getText(t, server.Handler(), "/", http.StatusOK)
	if body != "custom" {
		t.Fatalf("custom handler body = %q", body)
	}
}

func TestDirFrontendServesStaticAndFallsBackToIndex(t *testing.T) {
	runner := newTestRunner(t, `exports.pages = {};`)
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), `<html><body>app shell</body></html>`)
	writeFile(t, filepath.Join(dir, "asset.txt"), `asset`)
	server, err := New(Config{Runner: runner, FrontendMode: FrontendDir, FrontendDir: dir})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	text := getText(t, server.Handler(), "/asset.txt", http.StatusOK)
	if text != "asset" {
		t.Fatalf("asset text = %q", text)
	}
	fallback := getText(t, server.Handler(), "/unknown/route", http.StatusOK)
	if fallback != "<html><body>app shell</body></html>" {
		t.Fatalf("fallback = %q", fallback)
	}
}

func newTestServer(t *testing.T, script string) *Server {
	t.Helper()
	runner := newTestRunner(t, script)
	server, err := New(Config{Runner: runner, FrontendMode: FrontendAPIOnly, Dev: true})
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	return server
}

func newTestRunner(t *testing.T, script string) *widgetrunner.Runner {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "page.js"), script)
	runner, err := widgetrunner.New(context.Background(), widgetrunner.Config{ScriptDirs: []string{dir}})
	if err != nil {
		t.Fatalf("new runner: %v", err)
	}
	t.Cleanup(func() { _ = runner.Close(context.Background()) })
	return runner
}

func containsString(values []any, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func getJSON(t *testing.T, handler http.Handler, path string, want int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != want {
		t.Fatalf("GET %s = %d, want %d body=%s", path, rec.Code, want, rec.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v body=%s", err, rec.Body.String())
	}
	return body
}

func postJSON(t *testing.T, handler http.Handler, path string, want int) map[string]any {
	t.Helper()
	return postJSONBody(t, handler, path, "", want)
}

func postJSONBody(t *testing.T, handler http.Handler, path, body string, want int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != want {
		t.Fatalf("POST %s = %d, want %d body=%s", path, rec.Code, want, rec.Body.String())
	}
	var decoded map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode body: %v body=%s", err, rec.Body.String())
	}
	return decoded
}

func getText(t *testing.T, handler http.Handler, path string, want int) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != want {
		t.Fatalf("GET %s = %d, want %d body=%s", path, rec.Code, want, rec.Body.String())
	}
	return rec.Body.String()
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
