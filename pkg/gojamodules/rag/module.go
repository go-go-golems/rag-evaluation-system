// Package rag exposes the typed RAG laboratory builder through require("rag").
package rag

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
	"github.com/go-go-golems/rag-evaluation-system/pkg/raglab"
)

const ModuleName = "rag"

const hiddenExperiment = "__ragExperiment"
const hiddenFragment = "__ragFragment"
const hiddenLaboratory = "__ragLaboratory"

type LaboratoryFactory func(raglab.OpenOptions) (*raglab.Laboratory, error)

type module struct{ factory LaboratoryFactory }

var _ modules.NativeModule = (*module)(nil)
var _ modules.TypeScriptDeclarer = (*module)(nil)

func NewLoader(factory ...LaboratoryFactory) require.ModuleLoader {
	selected := raglab.OpenSQLite
	if len(factory) > 0 && factory[0] != nil {
		selected = factory[0]
	}
	return (&module{factory: selected}).Loader
}

func NewRegistrar() engine.RuntimeModuleRegistrar {
	return engine.NativeModuleRegistrar{ModuleID: "raglab", ModuleName: ModuleName, Loader: NewLoader()}
}

func (m *module) Name() string { return ModuleName }
func (m *module) Doc() string {
	return "Typed fluent RAG laboratory builders that compile to immutable experiment specifications."
}
func (m *module) TypeScriptModule() *spec.Module { return TypeScriptModule() }

func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	exports := moduleObj.Get("exports").(*goja.Object)
	runtime := &runtime{vm: vm, factory: m.factory}
	modules.SetExport(exports, ModuleName, "open", runtime.open)
	modules.SetExport(exports, ModuleName, "experiment", runtime.experiment)
	modules.SetExport(exports, ModuleName, "fragment", runtime.fragment)
	modules.SetExport(exports, ModuleName, "artifact", runtime.artifact)
	modules.SetExport(exports, ModuleName, "grade", runtime.grade)
	modules.SetExport(exports, ModuleName, "version", "v1")
}

func init() { modules.Register(&module{factory: raglab.OpenSQLite}) }

type runtime struct {
	vm      *goja.Runtime
	factory LaboratoryFactory
}

type experimentHandle struct{ builder *raglab.ExperimentBuilder }
type fragmentHandle struct{ fragment raglab.Fragment }
type laboratoryHandle struct{ laboratory *raglab.Laboratory }

type gojaQueryEmbedder struct {
	runtime  *runtime
	callback goja.Callable
}

var _ raglab.QueryEmbedder = (*gojaQueryEmbedder)(nil)

func (r *runtime) open(call goja.FunctionCall) goja.Value {
	options := r.objectArgument(call.Argument(0), "open options")
	database := r.stringProperty(options, "database")
	if database == "" {
		r.throwType("RAG_DATABASE_REQUIRED", "database path is required")
	}
	execution := r.stringProperty(options, "execution")
	if execution != "" && execution != "readOnly" && execution != "allowRuns" {
		r.throwType("RAG_INVALID_EXECUTION", "execution must be readOnly or allowRuns")
	}
	openOptions := raglab.OpenOptions{Database: database, AllowRuns: execution == "allowRuns"}
	if callbackValue := options.Get("queryEmbed"); callbackValue != nil && !goja.IsUndefined(callbackValue) && !goja.IsNull(callbackValue) {
		callback, ok := goja.AssertFunction(callbackValue)
		if !ok {
			r.throwType("RAG_INVALID_QUERY_EMBED", "queryEmbed must be a synchronous function returning number[]")
		}
		openOptions.QueryEmbedder = &gojaQueryEmbedder{runtime: r, callback: callback}
	}
	if rerankerValue := options.Get("reranker"); rerankerValue != nil && !goja.IsUndefined(rerankerValue) && !goja.IsNull(rerankerValue) {
		rerankerOptions := r.objectArgument(rerankerValue, "reranker options")
		if r.stringProperty(rerankerOptions, "kind") != "llama.cpp" {
			r.throwType("RAG_INVALID_RERANKER", "reranker.kind must be llama.cpp")
		}
		maxRequestBytes := 0
		if value := rerankerOptions.Get("maxRequestBytes"); value != nil && !goja.IsUndefined(value) && !goja.IsNull(value) {
			maxRequestBytes = int(value.ToInteger())
		}
		reranker, err := raglab.NewLlamaCPPReranker(raglab.LlamaCPPRerankerOptions{
			BaseURL:         r.stringProperty(rerankerOptions, "baseURL"),
			Model:           r.stringProperty(rerankerOptions, "model"),
			MaxRequestBytes: maxRequestBytes,
		})
		if err != nil {
			r.throw(err)
		}
		openOptions.Reranker = reranker
	}
	laboratory, err := r.factory(openOptions)
	if err != nil {
		r.throw(err)
	}
	return r.laboratoryObject(&laboratoryHandle{laboratory: laboratory})
}

func (r *runtime) experiment(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	builder := raglab.NewExperiment(name)
	if !goja.IsUndefined(call.Argument(1)) {
		r.configure(call.Argument(1), r.experimentObject(&experimentHandle{builder: builder}))
	}
	return r.experimentObject(&experimentHandle{builder: builder})
}

func (r *runtime) fragment(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	configure := call.Argument(1)
	if _, ok := goja.AssertFunction(configure); !ok {
		r.throwType("RAG_INVALID_FRAGMENT", "fragment configurator must be a function")
	}
	fragment := raglab.NewFragment(name, func(builder *raglab.ExperimentBuilder) {
		r.configure(configure, r.experimentObject(&experimentHandle{builder: builder}))
	})
	object := r.vm.NewObject()
	_ = object.Set(hiddenFragment, &fragmentHandle{fragment: fragment})
	return object
}

func (r *runtime) artifact(call goja.FunctionCall) goja.Value {
	ref := raglab.Artifact(raglab.ArtifactKind(call.Argument(0).String()), call.Argument(1).String())
	return r.vm.ToValue(map[string]any{"kind": string(ref.Kind), "id": ref.ID})
}

func (r *runtime) grade(call goja.FunctionCall) goja.Value {
	grade, err := raglab.Grade(call.Argument(0).String())
	if err != nil {
		r.throw(err)
	}
	return r.vm.ToValue(map[string]any{"name": grade.Name, "ordinal": grade.Ordinal})
}

func (r *runtime) experimentObject(handle *experimentHandle) *goja.Object {
	object := r.vm.NewObject()
	_ = object.Set(hiddenExperiment, handle)
	set := func(name string, fn func(goja.FunctionCall) goja.Value) {
		modules.SetExport(object, ModuleName, name, fn)
	}
	set("use", func(call goja.FunctionCall) goja.Value {
		handle.builder.Use(r.fragmentArgument(call.Argument(0)))
		return object
	})
	set("corpus", func(call goja.FunctionCall) goja.Value {
		handle.builder.Corpus(r.artifactArgument(call.Argument(0), raglab.CorpusSnapshotArtifact))
		return object
	})
	set("chunks", func(call goja.FunctionCall) goja.Value {
		handle.builder.Chunks(r.artifactArgument(call.Argument(0), raglab.ChunkSetArtifact))
		return object
	})
	set("bm25", func(call goja.FunctionCall) goja.Value {
		handle.builder.BM25(r.artifactArgument(call.Argument(0), raglab.BM25IndexArtifact))
		return object
	})
	set("embeddings", func(call goja.FunctionCall) goja.Value {
		handle.builder.Embeddings(r.artifactArgument(call.Argument(0), raglab.EmbeddingSetArtifact))
		return object
	})
	set("evaluation", func(call goja.FunctionCall) goja.Value {
		handle.builder.Evaluation(r.artifactArgument(call.Argument(0), raglab.EvaluationDatasetArtifact))
		return object
	})
	set("note", func(call goja.FunctionCall) goja.Value { handle.builder.Note(call.Argument(0).String()); return object })
	set("tag", func(call goja.FunctionCall) goja.Value {
		handle.builder.Tag(call.Argument(0).String(), call.Argument(1).String())
		return object
	})
	set("representations", func(call goja.FunctionCall) goja.Value {
		handle.builder.Representations(func(builder *raglab.RepresentationBuilder) {
			r.configure(call.Argument(0), r.representationsObject(builder))
		})
		return object
	})
	set("retrieval", func(call goja.FunctionCall) goja.Value {
		handle.builder.Retrieval(func(builder *raglab.RetrievalBuilder) { r.configure(call.Argument(0), r.retrievalObject(builder)) })
		return object
	})
	set("metrics", func(call goja.FunctionCall) goja.Value {
		handle.builder.Metrics(func(builder *raglab.MetricsBuilder) { r.configure(call.Argument(0), r.metricsObject(builder)) })
		return object
	})
	set("validate", func(call goja.FunctionCall) goja.Value { return r.validateExperiment(handle, call.Argument(0)) })
	set("toSpec", func(goja.FunctionCall) goja.Value { return r.specValue(r.build(handle)) })
	set("toJSON", func(goja.FunctionCall) goja.Value { return r.specValue(r.build(handle)) })
	return object
}

func (r *runtime) representationsObject(builder *raglab.RepresentationBuilder) *goja.Object {
	object := r.vm.NewObject()
	set := func(name string, fn func(goja.FunctionCall) goja.Value) {
		modules.SetExport(object, ModuleName, name, fn)
	}
	set("rawChunks", func(call goja.FunctionCall) goja.Value {
		name := ""
		if !goja.IsUndefined(call.Argument(0)) {
			name = call.Argument(0).String()
		}
		builder.RawChunks(name)
		return object
	})
	materialized := func(kind string, call goja.FunctionCall) goja.Value {
		name := call.Argument(0).String()
		config := &representationConfig{runtime: r, kind: kind}
		r.configure(call.Argument(1), config.object())
		if config.parent != "sourceChunk" {
			r.throwType("RAG_INVALID_REPRESENTATION", "materialized representations require parent(\"sourceChunk\")")
		}
		if kind == "summary" {
			builder.Summaries(name, config.artifact)
		} else {
			builder.Questions(name, config.artifact)
		}
		return object
	}
	set("summaries", func(call goja.FunctionCall) goja.Value { return materialized("summary", call) })
	set("questions", func(call goja.FunctionCall) goja.Value { return materialized("question", call) })
	return object
}

type representationConfig struct {
	runtime  *runtime
	kind     string
	artifact raglab.ArtifactRef
	parent   string
}

func (c *representationConfig) object() *goja.Object {
	object := c.runtime.vm.NewObject()
	modules.SetExport(object, ModuleName, "artifact", func(call goja.FunctionCall) goja.Value {
		c.artifact = c.runtime.artifactArgument(call.Argument(0), raglab.RepresentationSetArtifact)
		return object
	})
	modules.SetExport(object, ModuleName, "parent", func(call goja.FunctionCall) goja.Value { c.parent = call.Argument(0).String(); return object })
	return object
}

func (r *runtime) retrievalObject(builder *raglab.RetrievalBuilder) *goja.Object {
	object := r.vm.NewObject()
	set := func(name string, fn func(goja.FunctionCall) goja.Value) {
		modules.SetExport(object, ModuleName, name, fn)
	}
	set("channel", func(call goja.FunctionCall) goja.Value {
		builder.Channel(call.Argument(0).String(), func(channel *raglab.ChannelBuilder) { r.configure(call.Argument(1), r.channelObject(channel)) })
		return object
	})
	set("filter", func(call goja.FunctionCall) goja.Value {
		builder.Filter(func(filter *raglab.FilterBuilder) { r.configure(call.Argument(0), r.filterObject(filter)) })
		return object
	})
	set("fuse", func(call goja.FunctionCall) goja.Value {
		r.configure(call.Argument(0), r.fusionObject(builder))
		return object
	})
	set("rerank", func(call goja.FunctionCall) goja.Value {
		config := &rerankingConfig{runtime: r}
		r.configure(call.Argument(0), config.object())
		builder.RerankCrossEncoder(config.model, config.candidates, config.results)
		return object
	})
	set("collapse", func(call goja.FunctionCall) goja.Value {
		builder.Collapse(raglab.CollapseScope(call.Argument(0).String()))
		return object
	})
	set("results", func(call goja.FunctionCall) goja.Value {
		builder.Results(int(call.Argument(0).ToInteger()))
		return object
	})
	return object
}

// rerankingConfig collects JavaScript fluent settings before the typed Go
// builder materializes the canonical RerankingSpec.
type rerankingConfig struct {
	runtime    *runtime
	model      string
	candidates int
	results    int
}

func (c *rerankingConfig) object() *goja.Object {
	object := c.runtime.vm.NewObject()
	modules.SetExport(object, ModuleName, "crossEncoder", func(call goja.FunctionCall) goja.Value {
		c.model = call.Argument(0).String()
		return object
	})
	modules.SetExport(object, ModuleName, "candidates", func(call goja.FunctionCall) goja.Value {
		c.candidates = int(call.Argument(0).ToInteger())
		return object
	})
	modules.SetExport(object, ModuleName, "results", func(call goja.FunctionCall) goja.Value {
		c.results = int(call.Argument(0).ToInteger())
		return object
	})
	return object
}

func (r *runtime) channelObject(builder *raglab.ChannelBuilder) *goja.Object {
	object := r.vm.NewObject()
	modules.SetExport(object, ModuleName, "bm25", func(goja.FunctionCall) goja.Value { builder.BM25(); return object })
	modules.SetExport(object, ModuleName, "vector", func(goja.FunctionCall) goja.Value { builder.Vector(); return object })
	modules.SetExport(object, ModuleName, "representation", func(call goja.FunctionCall) goja.Value {
		builder.Representation(call.Argument(0).String())
		return object
	})
	modules.SetExport(object, ModuleName, "topK", func(call goja.FunctionCall) goja.Value {
		builder.TopK(int(call.Argument(0).ToInteger()))
		return object
	})
	modules.SetExport(object, ModuleName, "filter", func(call goja.FunctionCall) goja.Value {
		builder.Filter(func(filter *raglab.FilterBuilder) { r.configure(call.Argument(0), r.filterObject(filter)) })
		return object
	})
	return object
}

func (r *runtime) fusionObject(builder *raglab.RetrievalBuilder) *goja.Object {
	object := r.vm.NewObject()
	modules.SetExport(object, ModuleName, "rrf", func(goja.FunctionCall) goja.Value { builder.FuseRRF(60); return object })
	modules.SetExport(object, ModuleName, "rankConstant", func(call goja.FunctionCall) goja.Value {
		builder.FuseRRF(int(call.Argument(0).ToInteger()))
		return object
	})
	modules.SetExport(object, ModuleName, "weight", func(call goja.FunctionCall) goja.Value {
		builder.Weight(call.Argument(0).String(), call.Argument(1).ToFloat())
		return object
	})
	return object
}

func (r *runtime) filterObject(builder *raglab.FilterBuilder) *goja.Object {
	object := r.vm.NewObject()
	modules.SetExport(object, ModuleName, "sourceIds", func(call goja.FunctionCall) goja.Value {
		builder.SourceIDs(r.stringsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "documentIds", func(call goja.FunctionCall) goja.Value {
		builder.DocumentIDs(r.stringsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "contentTypes", func(call goja.FunctionCall) goja.Value {
		builder.ContentTypes(r.stringsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "metadataEquals", func(call goja.FunctionCall) goja.Value {
		builder.MetadataEquals(call.Argument(0).String(), call.Argument(1).String())
		return object
	})
	return object
}

func (r *runtime) metricsObject(builder *raglab.MetricsBuilder) *goja.Object {
	object := r.vm.NewObject()
	modules.SetExport(object, ModuleName, "relevanceAt", func(call goja.FunctionCall) goja.Value {
		builder.RelevanceAt(r.gradeArgument(call.Argument(0)))
		return object
	})
	modules.SetExport(object, ModuleName, "precisionAt", func(call goja.FunctionCall) goja.Value {
		builder.PrecisionAt(r.intsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "recallAt", func(call goja.FunctionCall) goja.Value {
		builder.RecallAt(r.intsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "hitRateAt", func(call goja.FunctionCall) goja.Value {
		builder.HitRateAt(r.intsArgument(call.Argument(0))...)
		return object
	})
	modules.SetExport(object, ModuleName, "ndcgAt", func(call goja.FunctionCall) goja.Value {
		builder.NDCGAt(int(call.Argument(0).ToInteger()))
		return object
	})
	modules.SetExport(object, ModuleName, "mrr", func(goja.FunctionCall) goja.Value { builder.MRR(); return object })
	modules.SetExport(object, ModuleName, "meanRelevantRecallAt", func(call goja.FunctionCall) goja.Value {
		builder.MeanRelevantRecallAt(int(call.Argument(0).ToInteger()))
		return object
	})
	modules.SetExport(object, ModuleName, "abstention", func(goja.FunctionCall) goja.Value { builder.Abstention(); return object })
	return object
}

func (r *runtime) laboratoryObject(handle *laboratoryHandle) *goja.Object {
	object := r.vm.NewObject()
	_ = object.Set(hiddenLaboratory, handle)
	modules.SetExport(object, ModuleName, "validate", func(call goja.FunctionCall) goja.Value {
		return r.validateExperiment(r.experimentArgument(call.Argument(0)), object)
	})
	modules.SetExport(object, ModuleName, "persist", func(call goja.FunctionCall) goja.Value {
		persisted, err := handle.laboratory.Persist(context.Background(), r.build(r.experimentArgument(call.Argument(0))))
		if err != nil {
			r.throw(err)
		}
		return r.vm.ToValue(map[string]any{"id": persisted.Specification.ID, "reused": persisted.Reused, "schemaVersion": persisted.Specification.SchemaVersion})
	})
	modules.SetExport(object, ModuleName, "start", func(call goja.FunctionCall) goja.Value {
		run, err := handle.laboratory.Start(context.Background(), r.build(r.experimentArgument(call.Argument(0))))
		if err != nil {
			r.throw(err)
		}
		return r.vm.ToValue(map[string]any{"id": run.ID, "experimentSpecId": run.ExperimentSpecID, "status": run.Status})
	})
	modules.SetExport(object, ModuleName, "execute", func(call goja.FunctionCall) goja.Value {
		result, err := handle.laboratory.Run(context.Background(), r.build(r.experimentArgument(call.Argument(0))))
		if err != nil {
			r.throw(err)
		}
		return r.vm.ToValue(map[string]any{"runId": result.RunID, "queryCount": result.QueryCount, "metrics": result.Metrics, "timing": result.Timing, "completedAt": result.CompletedAt})
	})
	modules.SetExport(object, ModuleName, "close", func(goja.FunctionCall) goja.Value {
		if err := handle.laboratory.Close(); err != nil {
			r.throw(err)
		}
		return goja.Undefined()
	})
	return object
}

func (r *runtime) validateExperiment(handle *experimentHandle, laboratoryValue goja.Value) goja.Value {
	report := handle.builder.Validate()
	if report.OK() && !goja.IsUndefined(laboratoryValue) && !goja.IsNull(laboratoryValue) {
		lab := r.laboratoryArgument(laboratoryValue)
		compatibility := lab.laboratory.Validate(context.Background(), r.build(handle))
		report.Issues = append(report.Issues, compatibility.Issues...)
		report.Normalize()
	}
	return r.reportValue(report)
}

func (r *runtime) build(handle *experimentHandle) raglab.ExperimentSpecification {
	specification, err := handle.builder.Build()
	if err != nil {
		r.throw(err)
	}
	return specification
}

func (r *runtime) configure(value goja.Value, argument goja.Value) {
	callback, ok := goja.AssertFunction(value)
	if !ok {
		r.throwType("RAG_CONFIGURATOR_REQUIRED", "builder configurator must be a function")
	}
	if _, err := callback(goja.Undefined(), argument); err != nil {
		panic(err)
	}
}

func (e *gojaQueryEmbedder) GenerateEmbedding(_ context.Context, text string) ([]float32, error) {
	value, err := e.callback(goja.Undefined(), e.runtime.vm.ToValue(text))
	if err != nil {
		return nil, fmt.Errorf("RAG_QUERY_EMBED_FAILED: callback failed: %w", err)
	}
	if goja.IsUndefined(value) || goja.IsNull(value) {
		return nil, errors.New("RAG_QUERY_EMBED_INVALID: callback must return number[]")
	}
	object := value.ToObject(e.runtime.vm)
	length := int(object.Get("length").ToInteger())
	if length <= 0 {
		return nil, errors.New("RAG_QUERY_EMBED_INVALID: callback must return a non-empty number[]")
	}
	vector := make([]float32, length)
	for i := range vector {
		number := object.Get(fmt.Sprintf("%d", i)).ToFloat()
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, errors.New("RAG_QUERY_EMBED_INVALID: callback vector entries must be finite numbers")
		}
		vector[i] = float32(number)
	}
	return vector, nil
}

func (r *runtime) fragmentArgument(value goja.Value) raglab.Fragment {
	object := r.objectArgument(value, "fragment")
	handle, ok := object.Get(hiddenFragment).Export().(*fragmentHandle)
	if !ok || handle == nil {
		r.throwType("RAG_INVALID_FRAGMENT", "expected a rag.fragment value")
	}
	return handle.fragment
}
func (r *runtime) experimentArgument(value goja.Value) *experimentHandle {
	object := r.objectArgument(value, "experiment")
	handle, ok := object.Get(hiddenExperiment).Export().(*experimentHandle)
	if !ok || handle == nil {
		r.throwType("RAG_INVALID_EXPERIMENT", "expected a rag.experiment value")
	}
	return handle
}
func (r *runtime) laboratoryArgument(value goja.Value) *laboratoryHandle {
	object := r.objectArgument(value, "laboratory")
	handle, ok := object.Get(hiddenLaboratory).Export().(*laboratoryHandle)
	if !ok || handle == nil {
		r.throwType("RAG_LAB_REQUIRED", "expected a rag.open laboratory")
	}
	return handle
}
func (r *runtime) artifactArgument(value goja.Value, expected raglab.ArtifactKind) raglab.ArtifactRef {
	if goja.IsUndefined(value) || goja.IsNull(value) {
		return raglab.Artifact(expected, "")
	}
	if value.ExportType().Kind() == 0 {
		return raglab.Artifact(expected, value.String())
	}
	if value.ToObject(r.vm).ClassName() == "String" {
		return raglab.Artifact(expected, value.String())
	}
	object := value.ToObject(r.vm)
	kind := object.Get("kind").String()
	id := object.Get("id").String()
	if kind != "" && kind != string(expected) {
		r.throwType("RAG_INVALID_ARTIFACT", fmt.Sprintf("expected %s artifact, got %s", expected, kind))
	}
	return raglab.Artifact(expected, id)
}
func (r *runtime) gradeArgument(value goja.Value) raglab.RelevanceGrade {
	name := value.String()
	if object := value.ToObject(r.vm); object != nil && !goja.IsUndefined(object.Get("name")) {
		name = object.Get("name").String()
	}
	grade, err := raglab.Grade(name)
	if err != nil {
		r.throw(err)
	}
	return grade
}
func (r *runtime) objectArgument(value goja.Value, label string) *goja.Object {
	if goja.IsUndefined(value) || goja.IsNull(value) {
		r.throwType("RAG_INVALID_ARGUMENT", label+" is required")
	}
	return value.ToObject(r.vm)
}
func (r *runtime) stringProperty(object *goja.Object, name string) string {
	value := object.Get(name)
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return ""
	}
	return value.String()
}
func (r *runtime) stringsArgument(value goja.Value) []string {
	object := value.ToObject(r.vm)
	length := int(object.Get("length").ToInteger())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, object.Get(fmt.Sprintf("%d", i)).String())
	}
	return result
}
func (r *runtime) intsArgument(value goja.Value) []int {
	object := value.ToObject(r.vm)
	length := int(object.Get("length").ToInteger())
	result := make([]int, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, int(object.Get(fmt.Sprintf("%d", i)).ToInteger()))
	}
	return result
}
func (r *runtime) reportValue(report raglab.ValidationReport) goja.Value {
	issues := make([]map[string]any, 0, len(report.Issues))
	for _, issue := range report.Issues {
		issues = append(issues, map[string]any{"code": issue.Code, "path": issue.Path, "message": issue.Message, "severity": string(issue.Severity)})
	}
	return r.vm.ToValue(map[string]any{"ok": report.OK(), "issues": issues})
}
func (r *runtime) specValue(specification raglab.ExperimentSpecification) goja.Value {
	inputs := map[string]any{
		"corpusSnapshot":    artifactValue(specification.Inputs.CorpusSnapshot),
		"chunkSet":          artifactValue(specification.Inputs.ChunkSet),
		"evaluationDataset": artifactValue(specification.Inputs.EvaluationDataset),
		"representations":   representationValues(specification.Inputs.Representations),
	}
	if specification.Inputs.BM25Index != nil {
		inputs["bm25Index"] = artifactValue(*specification.Inputs.BM25Index)
	}
	if specification.Inputs.EmbeddingSet != nil {
		inputs["embeddingSet"] = artifactValue(*specification.Inputs.EmbeddingSet)
	}
	return r.vm.ToValue(map[string]any{
		"schemaVersion": specification.SchemaVersion,
		"fingerprint":   specification.Fingerprint,
		"name":          specification.Name,
		"provenance": map[string]any{
			"fragments": specification.Provenance.Fragments,
			"notes":     specification.Provenance.Notes,
			"tags":      specification.Provenance.Tags,
		},
		"inputs":    inputs,
		"retrieval": retrievalValue(specification.Retrieval),
		"metrics":   metricsValue(specification.Metrics),
	})
}

// The Go authoring model deliberately uses Go field names internally. The
// JavaScript contract is a separate lower-camel plain-object projection so
// JSON.stringify(), examples, and generated declarations agree exactly.
func artifactValue(ref raglab.ArtifactRef) map[string]any {
	return map[string]any{"kind": string(ref.Kind), "id": ref.ID}
}

func representationValues(items []raglab.RepresentationSpec) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]any{
			"name": item.Name, "kind": string(item.Kind), "artifactId": item.ArtifactID, "parent": item.Parent,
		})
	}
	return result
}

func filterValue(filter raglab.FilterSpec) map[string]any {
	return map[string]any{
		"sourceIds": filter.SourceIDs, "documentIds": filter.DocumentIDs,
		"contentTypes": filter.ContentTypes, "metadataEquals": filter.MetadataEquals,
	}
}

func retrievalValue(plan raglab.RetrievalPlan) map[string]any {
	channels := make([]map[string]any, 0, len(plan.Channels))
	for _, channel := range plan.Channels {
		channels = append(channels, map[string]any{
			"name": channel.Name, "backend": string(channel.Backend), "representation": channel.Representation,
			"topK": channel.TopK, "filter": filterValue(channel.Filter),
		})
	}
	result := map[string]any{
		"channels": channels, "filter": filterValue(plan.Filter),
		"collapse": string(plan.Collapse), "results": plan.Results,
	}
	if plan.Fusion != nil {
		result["fusion"] = map[string]any{
			"kind": plan.Fusion.Kind, "rankConstant": plan.Fusion.RankConstant, "weights": plan.Fusion.Weights,
		}
	}
	if plan.Reranking != nil {
		result["reranking"] = map[string]any{
			"kind": string(plan.Reranking.Kind), "model": plan.Reranking.Model,
			"candidateCount": plan.Reranking.CandidateCount, "results": plan.Reranking.Results,
		}
	}
	return result
}

func metricsValue(plan raglab.MetricsPlan) map[string]any {
	result := map[string]any{
		"precisionAt": plan.PrecisionAt, "recallAt": plan.RecallAt, "hitRateAt": plan.HitRateAt,
		"ndcgAt": plan.NDCGAt, "mrr": plan.MRR, "meanRelevantRecallAt": plan.MeanRelevantRecall,
		"abstention": plan.Abstention,
	}
	if plan.RelevanceAt != nil {
		result["relevanceAt"] = map[string]any{"name": plan.RelevanceAt.Name, "ordinal": plan.RelevanceAt.Ordinal}
	}
	return result
}
func (r *runtime) throw(err error) {
	validation := new(raglab.ValidationError)
	if errors.As(err, &validation) {
		issueCode := "RAG_VALIDATION_FAILED"
		if len(validation.Report.Issues) > 0 {
			issueCode = validation.Report.Issues[0].Code
		}
		r.throwType(issueCode, validation.Error())
	}
	panic(r.vm.NewGoError(err))
}
func (r *runtime) throwType(code, message string) {
	value := r.vm.NewTypeError(message)
	_ = value.Set("code", code)
	panic(value)
}
