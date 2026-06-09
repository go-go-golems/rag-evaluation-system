package widgetmanifest

import (
	"path/filepath"
	"testing"
)

func TestCurrentManifestCatalogHasAdapters(t *testing.T) {
	root := filepath.Clean("../..")
	catalog, err := Discover(root)
	if err != nil {
		t.Fatalf("discover manifests: %v", err)
	}
	if len(catalog.Widgets) < 5 {
		t.Fatalf("expected at least the five pilot widget manifests, got %d", len(catalog.Widgets))
	}
	findings := Validate(catalog)
	for _, finding := range findings {
		if finding.IsError() {
			t.Fatalf("unexpected manifest error: %+v", finding)
		}
		if finding.Check == "adapter_exists" {
			t.Fatalf("adapter should exist after Phase 2 adapter slice: %+v", finding)
		}
	}
}
