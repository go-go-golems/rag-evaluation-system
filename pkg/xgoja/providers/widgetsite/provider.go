package widgetsite

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
	"github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc"
)

const PackageID = "rag-widget-site"

// Register exposes the RAG WidgetRenderer authoring modules to generated xgoja
// binaries. Selecting widget.dsl or rag.dsl in xgoja.yaml makes the same
// JSON-compatible Widget IR DSL available to JavaScript via require().
func Register(registry *providerapi.ProviderRegistry) error {
	loader := func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
		return widgetdsl.NewLoader(), nil
	}
	return registry.Package(PackageID,
		providerapi.Module{
			Name:             "widget.dsl",
			DefaultAs:        "widget.dsl",
			Description:      "RAG WidgetRenderer authoring DSL for JSON-compatible Widget IR pages.",
			NewModuleFactory: loader,
		},
		providerapi.Module{
			Name:             "rag.dsl",
			DefaultAs:        "rag.dsl",
			Description:      "Alias for widget.dsl in RAG-oriented xgoja scripts.",
			NewModuleFactory: loader,
		},
		providerapi.HelpSource{
			Name:        "widget-dsl",
			Description: "Getting started and JavaScript API reference for widget.dsl and rag.dsl.",
			FS:          doc.FS,
			Root:        ".",
		},
	)
}
