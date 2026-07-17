package rag

import (
	"context"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

func newVM() *goja.Runtime {
	vm := goja.New()
	registry := require.NewRegistry()
	registry.RegisterNativeModule(ModuleName, NewLoader())
	registry.Enable(vm)
	return vm
}

func TestModuleExportsPureResearchctlSpecification(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(`
		const rag = require("rag");
		rag.experiment("hybrid", e => e
			.corpus("snapshot").chunks("chunks").bm25("bm25").evaluation("eval")
			.tag("corpus", "ttc")
			.representations(r => r.rawChunks("raw"))
			.retrieval(r => r.channel("lexical", c => c.bm25().representation("raw").topK(50)).collapse("document").results(10))
			.metrics(m => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([10]).mrr()))
			.exportSpecification({ datasetSplit: "development" });
	`)
	if err != nil {
		t.Fatalf("pure export failed: %v", err)
	}
	exported := value.Export().(map[string]any)
	if exported["schemaVersion"] != "rag-retrieval-spec/v1" || exported["name"] != "hybrid" {
		t.Fatalf("export = %#v", exported)
	}
	if _, leaked := exported["inputs"]; leaked {
		t.Fatalf("domain export must not contain catalog input IDs: %#v", exported)
	}
	dataset := exported["dataset"].(map[string]any)
	if dataset["split"] != "development" {
		t.Fatalf("dataset = %#v", dataset)
	}
}

func TestModuleHasNoLifecycleOrPersistenceAuthority(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(`
		const rag = require("rag");
		const experiment = rag.experiment("pure");
		({
			open: typeof rag.open,
			toSpec: typeof experiment.toSpec,
			toJSON: typeof experiment.toJSON,
			run: typeof experiment.run,
			exportSpecification: typeof experiment.exportSpecification,
		});
	`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	for _, name := range []string{"open", "toSpec", "toJSON", "run"} {
		if got[name] != "undefined" {
			t.Fatalf("%s must be absent: %#v", name, got)
		}
	}
	if got["exportSpecification"] != "function" {
		t.Fatalf("pure export missing: %#v", got)
	}
}

func TestModuleReturnsDiagnosticsAndPreservesConfiguratorException(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(`
		const rag = require("rag");
		rag.experiment("invalid", e => e
			.corpus("snapshot").chunks("chunks").evaluation("eval")
			.retrieval(r => r.channel("semantic", c => c.vector().topK(5)).results(10))
			.metrics(m => m.mrr())).validate();
	`)
	if err != nil {
		t.Fatalf("validate should return diagnostics: %v", err)
	}
	issues := value.Export().(map[string]any)["issues"].([]map[string]any)
	if len(issues) < 3 {
		t.Fatalf("issues = %#v", issues)
	}
	_, err = vm.RunString(`require("rag").experiment("throws", () => { throw new Error("configurator exploded"); });`)
	if err == nil || !strings.Contains(err.Error(), "configurator exploded") {
		t.Fatalf("configurator error = %v", err)
	}
}

func TestEngineRegistrarCanRequirePureRagModule(t *testing.T) {
	factory, err := engine.NewRuntimeFactoryBuilder().WithModules(NewRegistrar()).Build()
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := factory.NewRuntime()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close(context.Background()) }()
	value, err := runtime.VM.RunString(`const rag = require("rag"); ({ open: typeof rag.open, experiment: typeof rag.experiment, version: rag.version });`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	if got["open"] != "undefined" || got["experiment"] != "function" || got["version"] != "v1" {
		t.Fatalf("exports = %#v", got)
	}
}

func TestTypeScriptModuleDeclaresPureRagSurface(t *testing.T) {
	declaration := strings.Join(TypeScriptModule().RawDTS, "\n")
	for _, want := range []string{"export interface Experiment", "export function grade", "export interface RerankingBuilder", "rerank(configure", "exportSpecification(options"} {
		if !strings.Contains(declaration, want) {
			t.Fatalf("TypeScript declaration missing %q: %s", want, declaration)
		}
	}
	for _, forbidden := range []string{"export function open", "interface Laboratory", "execute(experiment", "queryEmbed", "toSpec()", "toJSON()"} {
		if strings.Contains(declaration, forbidden) {
			t.Fatalf("TypeScript declaration retains lifecycle API %q: %s", forbidden, declaration)
		}
	}
}
