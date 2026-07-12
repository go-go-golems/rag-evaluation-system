package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
)

const apiReferencePath = "pkg/xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md"

func main() {
	data, err := os.ReadFile(apiReferencePath)
	if err != nil {
		panic(err)
	}
	const frontmatterEnd = "---\n\n"
	parts := strings.SplitN(string(data), frontmatterEnd, 2)
	if len(parts) != 2 {
		panic("API reference has no complete frontmatter")
	}
	const authoredSuffix = "## Using this reference"
	suffixAt := strings.Index(parts[1], authoredSuffix)
	if suffixAt < 0 {
		panic("API reference has no authored suffix")
	}
	generated := strings.TrimPrefix(widgetdsl.WidgetV3APIReferenceMarkdown(), "# widget.dsl API Reference\n\n")
	output := parts[0] + frontmatterEnd + generated + "\n" + parts[1][suffixAt:]
	if err := os.WriteFile(apiReferencePath, []byte(output), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("updated %s\n", apiReferencePath)
}
