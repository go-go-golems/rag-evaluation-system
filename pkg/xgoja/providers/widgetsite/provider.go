package widgetsite

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
	"github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc"
)

const PackageID = "rag-widget-site"

// Register exposes the single hard-cutover widget.dsl module to generated xgoja binaries.
func Register(registry *providerapi.ProviderRegistry) error {
	loader := func(moduleName string) func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
		return func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
			return widgetdsl.NewLoader(moduleName), nil
		}
	}
	return registry.Package(PackageID,
		providerapi.Module{
			Name:             widgetdsl.WidgetV3ModuleName,
			DefaultAs:        widgetdsl.WidgetV3ModuleName,
			Description:      "Preferred Widget DSL v3 root namespace for new Widget IR pages.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.WidgetV3ModuleName),
			NewModuleFactory: loader(widgetdsl.WidgetV3ModuleName),
		},
		providerapi.HelpSource{
			Name:        "widget-dsl",
			Description: "Getting started, v3 examples, API reference, and migration help for Widget IR DSL modules.",
			FS:          doc.FS,
			Root:        ".",
		},
	)
}
