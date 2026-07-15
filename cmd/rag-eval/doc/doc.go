// Package doc embeds the operator-facing rag-eval help pages.
package doc

import (
	"embed"

	"github.com/go-go-golems/glazed/pkg/help"
)

//go:embed *.md
var docFS embed.FS

// AddDocToHelpSystem registers this command's embedded help pages.
func AddDocToHelpSystem(helpSystem *help.HelpSystem) error {
	return helpSystem.LoadSectionsFromFS(docFS, ".")
}
