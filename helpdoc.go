package rageval

import (
	"embed"

	"github.com/go-go-golems/glazed/pkg/help"
	widgetdoc "github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc"
)

// coreHelpFS contains the repository-level Glazed help entries published with rag-eval.
//
//go:embed docs/guides/ttc-data-handbook.md docs/howtos/how-to-write-rag-eval-js-scripts.md docs/howtos/how-to-write-rag-eval-js-scripts-quick-reference.md
var coreHelpFS embed.FS

// AddDocToHelpSystem loads all public rag-eval and Widget DSL help entries.
func AddDocToHelpSystem(helpSystem *help.HelpSystem) error {
	if err := helpSystem.LoadSectionsFromFS(coreHelpFS, "."); err != nil {
		return err
	}
	return helpSystem.LoadSectionsFromFS(widgetdoc.FS, ".")
}
