// Package rag exposes the typed RAG laboratory builder through require("rag").
package rag

import (
	"encoding/json"
	"errors"
	"fmt"

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

type module struct{}

var _ modules.NativeModule = (*module)(nil)
var _ modules.TypeScriptDeclarer = (*module)(nil)

func NewLoader() require.ModuleLoader { return (&module{}).Loader }

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
	runtime := &runtime{vm: vm}
	modules.SetExport(exports, ModuleName, "experiment", runtime.experiment)
	modules.SetExport(exports, ModuleName, "fragment", runtime.fragment)
	modules.SetExport(exports, ModuleName, "artifact", runtime.artifact)
	modules.SetExport(exports, ModuleName, "grade", runtime.grade)
	modules.SetExport(exports, ModuleName, "version", "v1")
}

func init() { modules.Register(&module{}) }

type runtime struct{ vm *goja.Runtime }

type experimentHandle struct{ builder *raglab.ExperimentBuilder }
type fragmentHandle struct{ fragment raglab.Fragment }

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
	set("validate", func(goja.FunctionCall) goja.Value { return r.validateExperiment(handle) })
	set("exportSpecification", func(call goja.FunctionCall) goja.Value {
		options := r.objectArgument(call.Argument(0), "export options")
		exported, err := raglab.ExportSpecificationV1(r.build(handle), raglab.ExportOptions{DatasetSplit: r.stringProperty(options, "datasetSplit")})
		if err != nil {
			r.throw(err)
		}
		return r.jsonValue(exported)
	})
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

func (r *runtime) validateExperiment(handle *experimentHandle) goja.Value {
	return r.reportValue(handle.builder.Validate())
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

func (r *runtime) fragmentArgument(value goja.Value) raglab.Fragment {
	object := r.objectArgument(value, "fragment")
	handle, ok := object.Get(hiddenFragment).Export().(*fragmentHandle)
	if !ok || handle == nil {
		r.throwType("RAG_INVALID_FRAGMENT", "expected a rag.fragment value")
	}
	return handle.fragment
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
func (r *runtime) jsonValue(value any) goja.Value {
	encoded, err := json.Marshal(value)
	if err != nil {
		r.throw(err)
	}
	var plain any
	if err := json.Unmarshal(encoded, &plain); err != nil {
		r.throw(err)
	}
	return r.vm.ToValue(plain)
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
