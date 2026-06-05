package widgetsite

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/app"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
)

func TestRegisterExposesWidgetModules(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	for _, name := range []string{"widget.dsl", "rag.dsl"} {
		mod, ok := registry.ResolveModule(PackageID, name)
		if !ok {
			t.Fatalf("expected module %q", name)
		}
		loader, err := mod.NewModuleFactory(providerapi.ModuleSetupContext{Context: context.Background(), Name: name, As: name})
		if err != nil {
			t.Fatalf("new loader for %s: %v", name, err)
		}
		vm := goja.New()
		moduleObj := vm.NewObject()
		exports := vm.NewObject()
		if err := moduleObj.Set("exports", exports); err != nil {
			t.Fatalf("set exports: %v", err)
		}
		loader(vm, moduleObj)
		if got := exports.Get("panel"); got == nil || got.String() == "undefined" {
			t.Fatalf("module %s did not expose panel(): %#v", name, got)
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
	if len(entries) != 2 {
		t.Fatalf("expected two help entries, got %v", entries)
	}
	for _, want := range []string{"widget-dsl-getting-started", "widget-dsl-js-api-reference"} {
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
	runtimeSpec := &app.RuntimeSpec{
		Name: "widgetsite-provider-test",
		Modules: []app.ModuleInstanceSpec{{
			Package: PackageID,
			Name:    "widget.dsl",
			As:      "widget.dsl",
		}},
	}
	host := app.NewHost(registry, runtimeSpec)
	rt, err := host.Factory.NewRuntime(context.Background())
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	ret, err := rt.Owner.Call(context.Background(), "widgetsite-provider.require-widget-dsl", func(_ context.Context, vm *goja.Runtime) (any, error) {
		value, runErr := vm.RunString(`
			const rag = require("widget.dsl");
			const node = rag.panel({ title: "Demo" }, "ok");
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
	for _, want := range []string{`"kind":"component"`, `"type":"Panel"`, `"title":"Demo"`} {
		if !strings.Contains(json, want) {
			t.Fatalf("result missing %s: %s", want, json)
		}
	}
}
