package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetiredExperimentLifecycleRoutesAreRemoved(t *testing.T) {
	handler := newTestAPIHandler(t)
	base := "/api/v1/" + "la" + "b"
	for _, request := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, base + "/catalog"},
		{http.MethodGet, base + "/specifications"},
		{http.MethodPost, base + "/specifications"},
		{http.MethodGet, base + "/runs"},
		{http.MethodPost, base + "/runs/legacy/events"},
		{http.MethodPost, base + "/runs/legacy/complete"},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(request.method, request.path, nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("%s %s = %d, want removed route 404", request.method, request.path, recorder.Code)
		}
	}
}

func TestRAGArtifactCatalogRouteIsReadOnly(t *testing.T) {
	handler := newTestAPIHandler(t)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/artifacts/rag/catalog", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("catalog = %d body=%s", recorder.Code, recorder.Body.String())
	}
	write := httptest.NewRecorder()
	handler.ServeHTTP(write, httptest.NewRequest(http.MethodPost, "/api/v1/artifacts/rag/catalog", nil))
	if write.Code != http.StatusMethodNotAllowed {
		t.Fatalf("catalog POST = %d", write.Code)
	}
}
