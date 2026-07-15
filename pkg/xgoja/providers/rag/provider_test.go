package rag

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/app"
	"github.com/go-go-golems/go-go-goja/pkg/xgoja/providerapi"
	ragmodule "github.com/go-go-golems/rag-evaluation-system/pkg/gojamodules/rag"
)

func TestRegisterExposesRagModuleAndTypeScript(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	module, ok := registry.ResolveModule(PackageID, ragmodule.ModuleName)
	if !ok {
		t.Fatalf("expected module %q", ragmodule.ModuleName)
	}
	if module.TypeScript == nil {
		t.Fatal("expected TypeScript descriptor")
	}
}

func TestGeneratedRuntimeCanRequireRagModule(t *testing.T) {
	registry := providerapi.NewProviderRegistry()
	if err := Register(registry); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	host := app.NewHost(registry, &app.RuntimePlan{
		Schema: app.RuntimePlanSchema,
		Name:   "rag-provider-test",
		Runtime: app.RuntimeSection{Modules: []app.RuntimeModulePlan{{
			Provider: PackageID,
			Name:     ragmodule.ModuleName,
			As:       ragmodule.ModuleName,
		}}},
	})
	runtime, err := host.Factory.NewRuntime(context.Background())
	if err != nil {
		t.Fatalf("new runtime: %v", err)
	}
	defer func() { _ = runtime.Close(context.Background()) }()

	result, err := runtime.Owner.Call(context.Background(), "rag-provider-test.require", func(_ context.Context, vm *goja.Runtime) (any, error) {
		value, runErr := vm.RunString(`const rag = require("rag"); ({ experiment: typeof rag.experiment, version: rag.version });`)
		if runErr != nil {
			return nil, runErr
		}
		return value.Export(), nil
	})
	if err != nil {
		t.Fatalf("require rag: %v", err)
	}
	got := result.(map[string]any)
	if got["experiment"] != "function" || got["version"] != "v1" {
		t.Fatalf("unexpected RAG module exports: %#v", got)
	}
}
