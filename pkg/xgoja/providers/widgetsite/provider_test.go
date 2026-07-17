package widgetsite

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/app"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	"github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl"
)

func TestRegisterExposesOnlyWidgetDSL(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	cases := []struct{ name, helper string }{{widgetdsl.WidgetV3ModuleName, "page"}}
	for _, tc := range cases {
		mod, ok := registry.ResolveModule(PackageID, tc.name)
		if !ok {
			t.Fatalf("expected module %q", tc.name)
		}
		if mod.TypeScript == nil {
			t.Fatalf("expected module %q to carry TypeScript descriptor", tc.name)
		}
		loader, err := mod.NewModuleFactory(providerapi.ModuleSetupContext{Context: context.Background(), Name: tc.name, As: tc.name})
		if err != nil {
			t.Fatalf("new loader for %s: %v", tc.name, err)
		}
		vm := goja.New()
		moduleObj := vm.NewObject()
		exports := vm.NewObject()
		if err := moduleObj.Set("exports", exports); err != nil {
			t.Fatalf("set exports: %v", err)
		}
		loader(vm, moduleObj)
		if got := exports.Get(tc.helper); got == nil || got.String() == "undefined" {
			t.Fatalf("module %s did not expose %s(): %#v", tc.name, tc.helper, got)
		}
	}
	for _, oldName := range []string{"rag.dsl", widgetdsl.UIModuleName, widgetdsl.DataModuleName, widgetdsl.DataV2ModuleName, widgetdsl.ContextWindowModuleName, widgetdsl.CourseModuleName, widgetdsl.CmsModuleName} {
		if _, ok := registry.ResolveModule(PackageID, oldName); ok {
			t.Fatalf("old bucket module %q should not be exposed", oldName)
		}
	}
}

func TestRegisterExposesWidgetDSLHelpSource(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	source, ok := registry.ResolveHelpSource(PackageID, "widget-dsl")
	if !ok {
		t.Fatalf("expected widget-dsl help source")
	}
	entries, err := fs.Glob(source.FS, "*.md")
	if err != nil {
		t.Fatalf("glob help entries: %v", err)
	}
	if len(entries) != 5 {
		t.Fatalf("expected five help entries, got %v", entries)
	}
	for _, want := range []string{
		"widget-dsl-getting-started",
		"widget-dsl-js-api-reference",
		"widget-dsl-spa-bundling",
		"widget-dsl-v3-examples",
		"widget-dsl-v3-api-reference",
	} {
		found := false
		for _, entry := range entries {
			data, err := fs.ReadFile(source.FS, entry)
			if err != nil {
				t.Fatalf("read %s: %v", entry, err)
			}
			if strings.Contains(string(data), "Slug: "+want) {
				found = true
			}
		}
		if !found {
			t.Fatalf("missing help slug %s in %v", want, entries)
		}
	}
}

func TestGeneratedRuntimeCanRequireWidgetDSL(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	runtimePlan := &app.RuntimePlan{
		Schema:  app.RuntimePlanSchema,
		Name:    "widgetsite-provider-test",
		Runtime: app.RuntimeSection{Modules: []app.RuntimeModulePlan{{Provider: PackageID, Name: widgetdsl.WidgetV3ModuleName, As: widgetdsl.WidgetV3ModuleName}}},
	}
	host := app.NewHost(registry, runtimePlan)
	rt, err := host.Factory.NewRuntime(context.Background())
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	ret, err := rt.Owner.Call(context.Background(), "widgetsite-provider.require-widget-dsl", func(_ context.Context, vm *goja.Runtime) (any, error) {
		value, runErr := vm.RunString(`
			const widget = require("widget.dsl");
			const page = widget.page("Demo", p => p.section("Content", s => s.view(widget.ui.text("Ready"))));
			JSON.stringify(page.toPage());
		`)
		if runErr != nil {
			return nil, runErr
		}
		return value.String(), nil
	})
	if err != nil {
		t.Fatalf("run script: %v", err)
	}
	json := ret.(string)
	for _, want := range []string{`"kind":"component"`, `"type":"SectionBlock"`, `"title":"Demo"`, `"type":"Text"`} {
		if !strings.Contains(json, want) {
			t.Fatalf("result missing %s: %s", want, json)
		}
	}
}
