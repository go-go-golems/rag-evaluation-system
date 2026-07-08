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
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.UIModuleName),
			NewModuleFactory: loader(widgetdsl.UIModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.DataModuleName,
			DefaultAs:        widgetdsl.DataModuleName,
			Description:      "Widget IR data-display helpers, table cell helpers, and data recipes.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.DataModuleName),
			NewModuleFactory: loader(widgetdsl.DataModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.DataV2ModuleName,
			DefaultAs:        widgetdsl.DataV2ModuleName,
			Description:      "Experimental hard-cutover typed/fluent data DSL v2 builders.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.DataV2ModuleName),
			NewModuleFactory: loader(widgetdsl.DataV2ModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.WidgetV3ModuleName,
			DefaultAs:        widgetdsl.WidgetV3ModuleName,
			Description:      "Parallel clean Widget DSL v3 root namespace.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.WidgetV3ModuleName),
			NewModuleFactory: loader(widgetdsl.WidgetV3ModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.ContextWindowModuleName,
			DefaultAs:        widgetdsl.ContextWindowModuleName,
			Description:      "Context-window, transcript, annotation, and anchored-comment Widget IR helpers.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.ContextWindowModuleName),
			NewModuleFactory: loader(widgetdsl.ContextWindowModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.CourseModuleName,
			DefaultAs:        widgetdsl.CourseModuleName,
			Description:      "Course, lesson, slide, handout, and course-studio Widget IR helpers.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.CourseModuleName),
			NewModuleFactory: loader(widgetdsl.CourseModuleName),
		},
		providerapi.Module{
			Name:             widgetdsl.CmsModuleName,
			DefaultAs:        widgetdsl.CmsModuleName,
			Description:      "Media, asset, and article-management Widget IR helpers.",
			TypeScript:       widgetdsl.TypeScriptModule(widgetdsl.CmsModuleName),
			NewModuleFactory: loader(widgetdsl.CmsModuleName),
		},
		providerapi.HelpSource{
			Name:        "widget-dsl",
			Description: "Getting started and JavaScript API reference for split Widget IR DSL modules.",
			FS:          doc.FS,
			Root:        ".",
		},
	)
}
