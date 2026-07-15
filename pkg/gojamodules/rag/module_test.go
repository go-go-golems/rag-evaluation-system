package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/rag-evaluation-system/internal/experimentspec"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/go-go-golems/rag-evaluation-system/pkg/raglab"
)

type moduleCatalog map[raglab.ArtifactRef]raglab.ArtifactMetadata

func (c moduleCatalog) LookupArtifact(_ context.Context, ref raglab.ArtifactRef) (raglab.ArtifactMetadata, error) {
	metadata, ok := c[ref]
	if !ok {
		return raglab.ArtifactMetadata{}, raglab.ErrArtifactNotFound
	}
	return metadata, nil
}

type moduleStore struct{ specifications, runs, events int }

func (s *moduleStore) CreateSpecification(_ context.Context, input experimentspec.Input) (*experimentrun.Specification, bool, error) {
	s.specifications++
	return &experimentrun.Specification{ID: "spec-" + input.CorpusSnapshotID, SchemaVersion: experimentspec.SchemaVersion}, false, nil
}
func (s *moduleStore) CreateRun(_ context.Context, id string) (*experimentrun.Run, error) {
	s.runs++
	return &experimentrun.Run{ID: "run-" + id, ExperimentSpecID: id, Status: "running"}, nil
}
func (s *moduleStore) AppendEvent(_ context.Context, _ string, _ string, _ json.RawMessage) (*experimentrun.Event, error) {
	s.events++
	return &experimentrun.Event{Sequence: s.events}, nil
}

func testFactory(store *moduleStore) LaboratoryFactory {
	return func(options raglab.OpenOptions) (*raglab.Laboratory, error) {
		if options.Database != "fixture.db" {
			return nil, fmt.Errorf("unexpected database %q", options.Database)
		}
		catalog := moduleCatalog{
			raglab.CorpusSnapshot("snapshot"): {Ref: raglab.CorpusSnapshot("snapshot")},
			raglab.ChunkSet("chunks"):         {Ref: raglab.ChunkSet("chunks"), CorpusSnapshotID: "snapshot"},
			raglab.BM25Index("bm25"):          {Ref: raglab.BM25Index("bm25"), ChunkSetID: "chunks"},
			raglab.EmbeddingSet("embeddings"): {Ref: raglab.EmbeddingSet("embeddings"), ChunkSetID: "chunks", Dimensions: 768},
			raglab.EvaluationDataset("eval"):  {Ref: raglab.EvaluationDataset("eval"), CorpusSnapshotID: "snapshot", Status: "candidate"},
		}
		return raglab.NewLaboratory(catalog, store, options.AllowRuns), nil
	}
}

func newVM(factory LaboratoryFactory) *goja.Runtime {
	vm := goja.New()
	registry := require.NewRegistry()
	registry.RegisterNativeModule(ModuleName, NewLoader(factory))
	registry.Enable(vm)
	return vm
}

func TestModuleBuildsValidSpecAndStartsExplicitRun(t *testing.T) {
	store := &moduleStore{}
	vm := newVM(testFactory(store))
	value, err := vm.RunString(`
		const rag = require("rag");
		const baseline = rag.fragment("ttc-inputs", e => e
			.corpus("snapshot").chunks("chunks").bm25("bm25").embeddings("embeddings").evaluation("eval"));
		const lab = rag.open({ database: "fixture.db", execution: "allowRuns" });
		const experiment = rag.experiment("hybrid", e => e
			.use(baseline)
			.retrieval(r => r
				.channel("lexical", c => c.bm25().topK(50))
				.channel("semantic", c => c.vector().topK(50))
				.fuse(f => f.rrf().rankConstant(60).weight("semantic", 1.25))
				.collapse("document").results(10))
			.metrics(m => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([10, 1, 3]).mrr()));
		const report = experiment.validate(lab);
		const spec = experiment.toSpec();
		const persisted = lab.persist(experiment);
		const run = lab.start(experiment);
		({ report, spec, persisted, run, version: rag.version });
	`)
	if err != nil {
		t.Fatalf("run rag module: %v", err)
	}
	got := value.Export().(map[string]any)
	report := got["report"].(map[string]any)
	if report["ok"] != true {
		t.Fatalf("report = %#v", report)
	}
	specification := got["spec"].(map[string]any)
	if specification["schemaVersion"] != experimentspec.SchemaVersion || specification["fingerprint"] == "" {
		t.Fatalf("spec = %#v", specification)
	}
	inputs := specification["inputs"].(map[string]any)
	if _, ok := inputs["corpusSnapshot"]; !ok {
		t.Fatalf("JavaScript spec must use lower-camel keys: %#v", inputs)
	}
	if _, leakedGoName := inputs["CorpusSnapshot"]; leakedGoName {
		t.Fatalf("JavaScript spec leaked a Go field name: %#v", inputs)
	}
	retrieval := specification["retrieval"].(map[string]any)
	channels := retrieval["channels"].([]map[string]any)
	if channels[0]["topK"] != 50 {
		t.Fatalf("JavaScript retrieval projection = %#v", retrieval)
	}
	if got["version"] != "v1" || got["persisted"].(map[string]any)["id"] != "spec-snapshot" || got["run"].(map[string]any)["id"] != "run-spec-snapshot" {
		t.Fatalf("result = %#v", got)
	}
	if store.specifications != 2 || store.runs != 1 || store.events != 1 {
		t.Fatalf("store = %#v", store)
	}
}

func TestModuleReturnsDiagnosticsAndPreservesConfiguratorException(t *testing.T) {
	vm := newVM(testFactory(&moduleStore{}))
	value, err := vm.RunString(`
		const rag = require("rag");
		const invalid = rag.experiment("invalid", e => e
			.corpus("snapshot").chunks("chunks").evaluation("eval")
			.retrieval(r => r.channel("semantic", c => c.vector().topK(5)).results(10))
			.metrics(m => m.mrr()));
		invalid.validate();
	`)
	if err != nil {
		t.Fatalf("validate should return diagnostics: %v", err)
	}
	report := value.Export().(map[string]any)
	// goja preserves the concrete Go slice element type supplied by
	// reportValue(), so this exports as []map[string]any rather than []any.
	issues := report["issues"].([]map[string]any)
	if len(issues) < 3 {
		t.Fatalf("issues = %#v", issues)
	}
	_, err = vm.RunString(`require("rag").experiment("throws", () => { throw new Error("configurator exploded"); });`)
	if err == nil || !strings.Contains(err.Error(), "configurator exploded") {
		t.Fatalf("configurator error = %v", err)
	}
}

func TestEngineRegistrarCanRequireRagModule(t *testing.T) {
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
	if got["open"] != "function" || got["experiment"] != "function" || got["version"] != "v1" {
		t.Fatalf("exports = %#v", got)
	}
}

func TestTypeScriptModuleDeclaresRagSurface(t *testing.T) {
	declaration := strings.Join(TypeScriptModule().RawDTS, "\n")
	for _, want := range []string{"declare", "export function open", "export interface Experiment", "export interface Laboratory", "execute(experiment", "queryEmbed?: QueryEmbed", "export function grade"} {
		if !strings.Contains(declaration, want) && want != "declare" {
			t.Fatalf("TypeScript declaration missing %q: %s", want, declaration)
		}
	}
}

func TestGojaQueryEmbedderRequiresFiniteNumberArray(t *testing.T) {
	vm := goja.New()
	value, err := vm.RunString(`query => [query.length, 2.5]`)
	if err != nil {
		t.Fatal(err)
	}
	callback, ok := goja.AssertFunction(value)
	if !ok {
		t.Fatal("expected JavaScript callback")
	}
	embedder := &gojaQueryEmbedder{runtime: &runtime{vm: vm}, callback: callback}
	vector, err := embedder.GenerateEmbedding(context.Background(), "abc")
	if err != nil {
		t.Fatal(err)
	}
	if len(vector) != 2 || vector[0] != 3 || vector[1] != 2.5 {
		t.Fatalf("vector = %#v", vector)
	}
	invalidValue, err := vm.RunString(`() => [1, NaN]`)
	if err != nil {
		t.Fatal(err)
	}
	invalid, _ := goja.AssertFunction(invalidValue)
	embedder.callback = invalid
	if _, err := embedder.GenerateEmbedding(context.Background(), "ignored"); err == nil || !strings.Contains(err.Error(), "finite") {
		t.Fatalf("invalid vector error = %v", err)
	}
}
