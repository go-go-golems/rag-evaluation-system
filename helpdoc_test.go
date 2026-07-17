package rageval

import (
	"testing"

	"github.com/go-go-golems/glazed/pkg/help"
)

func TestAddDocToHelpSystemLoadsPublicEntries(t *testing.T) {
	t.Parallel()

	helpSystem := help.NewHelpSystem()
	if err := AddDocToHelpSystem(helpSystem); err != nil {
		t.Fatalf("AddDocToHelpSystem() error = %v", err)
	}

	for _, slug := range []string{
		"ttc-data-handbook",
		"how-to-write-rag-eval-js-scripts",
		"how-to-write-rag-eval-js-scripts-quick-reference",
		"widget-dsl-js-api-reference",
		"widget-dsl-v3-examples",
		"widget-dsl-v3-api-reference",
	} {
		if _, err := helpSystem.GetSectionWithSlug(slug); err != nil {
			t.Errorf("GetSectionWithSlug(%q) error = %v", slug, err)
		}
	}
}
