package widgetprovider

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
)

const PackageID = "rag-widget"

func Register(registry *providerapi.ProviderRegistry) error {
	return registry.Package(PackageID,
		providerapi.Module{
			Name:        "widget.dsl",
			DefaultAs:   "widget.dsl",
			Description: "RAG Widget IR authoring DSL for xgoja experiments",
			NewModuleFactory: func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
				return widgetdsl.NewLoader(widgetdsl.WidgetV3ModuleName), nil
			},
		},
		providerapi.Module{
			Name:        "rag.dsl",
			DefaultAs:   "rag.dsl",
			Description: "Alias for widget.dsl",
			NewModuleFactory: func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
				return widgetdsl.NewLoader(widgetdsl.WidgetV3ModuleName), nil
			},
		},
	)
}
