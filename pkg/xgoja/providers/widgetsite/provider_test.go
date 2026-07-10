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

func TestRegisterExposesSplitAndWidgetV3Modules(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	cases := []struct {
		name   string
		helper string
	}{
		{widgetdsl.UIModuleName, "panel"},
		{widgetdsl.DataModuleName, "dataTable"},
		{widgetdsl.DataV2ModuleName, "collection"},
		{widgetdsl.WidgetV3ModuleName, "raw"},
		{widgetdsl.ContextWindowModuleName, "contextDiagramPanel"},
		{widgetdsl.CourseModuleName, "courseStudioShell"},
	}
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
	for _, oldName := range []string{"rag.dsl"} {
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
	if len(entries) != 3 {
		t.Fatalf("expected three help entries, got %v", entries)
	}
	for _, want := range []string{"widget-dsl-getting-started", "widget-dsl-js-api-reference", "widget-dsl-spa-bundling"} {
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

func TestGeneratedRuntimeCanRequireSplitWidgetDSLModules(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	runtimePlan := &app.RuntimePlan{
		Schema: app.RuntimePlanSchema,
		Name:   "widgetsite-provider-test",
		Runtime: app.RuntimeSection{Modules: []app.RuntimeModulePlan{
			{Provider: PackageID, Name: widgetdsl.UIModuleName, As: widgetdsl.UIModuleName},
			{Provider: PackageID, Name: widgetdsl.DataModuleName, As: widgetdsl.DataModuleName},
			{Provider: PackageID, Name: widgetdsl.ContextWindowModuleName, As: widgetdsl.ContextWindowModuleName},
			{Provider: PackageID, Name: widgetdsl.CourseModuleName, As: widgetdsl.CourseModuleName},
		}},
	}
	host := app.NewHost(registry, runtimePlan)
	rt, err := host.Factory.NewRuntime(context.Background())
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	ret, err := rt.Owner.Call(context.Background(), "widgetsite-provider.require-split-dsl", func(_ context.Context, vm *goja.Runtime) (any, error) {
		value, runErr := vm.RunString(`
			const ui = require("ui.dsl");
			const data = require("data.dsl");
			const contextWindow = require("context_window.dsl");
			const course = require("course.dsl");
			const node = ui.panel({ title: "Demo" },
				data.dataTable({ rows: [], getRowKey: "id", columns: [{ id: "name", header: "Name", cell: data.cell.field("name") }] }),
				contextWindow.contextDiagramPanel({ snapshot: { id: "ctx", title: "Window", limit: 0, parts: [] } }),
				course.courseStudioShell({ sections: [], title: "Course" })
			);
			JSON.stringify(node);
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
	for _, want := range []string{`"kind":"component"`, `"type":"Panel"`, `"title":"Demo"`, `"type":"DataTable"`, `"type":"ContextDiagramPanel"`, `"type":"CourseStudioShell"`} {
		if !strings.Contains(json, want) {
			t.Fatalf("result missing %s: %s", want, json)
		}
	}
}
