package widgetsite

import (
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
	"github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc"
)

const PackageID = "rag-widget-site"

// Register exposes the split RAG WidgetRenderer authoring modules to generated
// xgoja binaries. Former bucket-style compatibility modules are intentionally
// not provided; scripts must import the domain module they use.
func Register(registry *providerapi.ProviderRegistry) error {
	loader := func(moduleName string) func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
		return func(providerapi.ModuleSetupContext) (require.ModuleLoader, error) {
			return widgetdsl.NewLoader(moduleName), nil
		}
	}
	return registry.Package(PackageID,
		providerapi.Module{
			Name:             widgetdsl.UIModuleName,
			DefaultAs:        widgetdsl.UIModuleName,
			Description:      "Generic Widget IR page, layout, primitive, and foundation helpers.",
			NewModuleFactory: loader(widgetdsl.UIModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.DataModuleName,
			DefaultAs:        widgetdsl.DataModuleName,
			Description:      "Widget IR data-display helpers, table cell helpers, and data recipes.",
			NewModuleFactory: loader(widgetdsl.DataModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.ContextWindowModuleName,
			DefaultAs:        widgetdsl.ContextWindowModuleName,
			Description:      "Context-window, transcript, annotation, and anchored-comment Widget IR helpers.",
			NewModuleFactory: loader(widgetdsl.ContextWindowModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.CourseModuleName,
			DefaultAs:        widgetdsl.CourseModuleName,
			Description:      "Course, lesson, slide, handout, and course-studio Widget IR helpers.",
			NewModuleFactory: loader(widgetdsl.CourseModuleName),
		},
		providerapi.HelpSource{
			Name:        "widget-dsl",
			Description: "Getting started and JavaScript API reference for split Widget IR DSL modules.",
			FS:          doc.FS,
			Root:        ".",
		},
	)
}
