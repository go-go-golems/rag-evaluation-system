package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPrototypeExperimentLifecycleRoutesAreRemoved(t *testing.T) {
	handler := newTestAPIHandler(t)
	for _, request := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/lab/specifications"},
		{http.MethodPost, "/api/v1/lab/specifications"},
		{http.MethodGet, "/api/v1/lab/runs"},
		{http.MethodPost, "/api/v1/lab/runs/legacy/events"},
		{http.MethodPost, "/api/v1/lab/runs/legacy/complete"},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(request.method, request.path, nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("%s %s = %d, want removed route 404", request.method, request.path, recorder.Code)
		}
	}
}

func TestDomainArtifactCatalogRouteRemainsReadOnly(t *testing.T) {
	handler := newTestAPIHandler(t)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/lab/catalog", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("catalog = %d body=%s", recorder.Code, recorder.Body.String())
	}
}
