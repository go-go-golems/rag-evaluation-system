// Package rag exposes pure Go-backed RAG v2 authoring through require("rag").
package rag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/modules"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
)

const ModuleName = "rag"

type module struct{}

var _ modules.NativeModule = (*module)(nil)
var _ modules.TypeScriptDeclarer = (*module)(nil)

func NewLoader() require.ModuleLoader { return (&module{}).Loader }
func NewRegistrar() engine.RuntimeModuleRegistrar {
	return engine.NativeModuleRegistrar{ModuleID: "rag", ModuleName: ModuleName, Loader: NewLoader()}
}
func (m *module) Name() string                   { return ModuleName }
func (m *module) Doc() string                    { return "Pure Go-backed composable RAG v2 authoring and compilation." }
func (m *module) TypeScriptModule() *spec.Module { return TypeScriptModule() }
func init()                                      { modules.Register(&module{}) }

type runtime struct {
	vm                                                                                                                                                                   *goja.Runtime
	descriptor, corpus, dataset, pipeline, pipelineBuilder, fragment, query, queryBuilder, product, productBuilder, study, studyBuilder, variant, variantBuilder, factor *goja.Symbol
}

func newRuntime(vm *goja.Runtime) *runtime {
	return &runtime{
		vm:         vm,
		descriptor: goja.NewSymbol("rag.descriptor"), corpus: goja.NewSymbol("rag.corpus"), dataset: goja.NewSymbol("rag.dataset"),
		pipeline: goja.NewSymbol("rag.pipeline"), pipelineBuilder: goja.NewSymbol("rag.pipelineBuilder"), fragment: goja.NewSymbol("rag.fragment"),
		query: goja.NewSymbol("rag.query"), queryBuilder: goja.NewSymbol("rag.queryBuilder"),
		product: goja.NewSymbol("rag.product"), productBuilder: goja.NewSymbol("rag.productBuilder"),
		study: goja.NewSymbol("rag.study"), studyBuilder: goja.NewSymbol("rag.studyBuilder"),
		variant: goja.NewSymbol("rag.variant"), variantBuilder: goja.NewSymbol("rag.variantBuilder"), factor: goja.NewSymbol("rag.factor"),
	}
}
func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
	exports := moduleObj.Get("exports").(*goja.Object)
	r := newRuntime(vm)
	set := func(name string, value any) { modules.SetExport(exports, ModuleName, name, value) }
	set("version", "v2")
	set("pipeline", r.pipelineFactory)
	set("fragment", r.fragmentFactory)
	set("queryPlan", r.queryFactory)
	set("product", r.productFactory)
	set("study", r.studyFactory)
	set("variant", r.variantFactory)
	set("validate", r.validate)
	set("explain", r.explain)
	set("compileProduct", r.compileProduct)
	set("compileStudy", r.compileStudy)
	set("preview", r.preview)
	_ = exports.Set("inputs", r.inputsObject())
	_ = exports.Set("units", r.unitsObject())
	_ = exports.Set("transcript", r.transcriptObject())
	_ = exports.Set("chunks", r.chunksObject())
	_ = exports.Set("representations", r.representationsObject())
	_ = exports.Set("embeddings", r.embeddingsObject())
	_ = exports.Set("indexes", r.indexesObject())
	_ = exports.Set("retrieve", r.retrieveObject())
	_ = exports.Set("collapse", r.collapseObject())
	_ = exports.Set("fusion", r.fusionObject())
	_ = exports.Set("hydration", r.hydrationObject())
	_ = exports.Set("rerank", r.rerankObject())
	_ = exports.Set("generation", r.generationObject())
	_ = exports.Set("metrics", r.metricsObject())
	_ = exports.Set("datasets", r.datasetsObject())
	_ = exports.Set("recipes", r.recipesObject())
}
func (r *runtime) obj(methods map[string]any) *goja.Object {
	o := r.vm.NewObject()
	for name, fn := range methods {
		modules.SetExport(o, ModuleName, name, fn)
	}
	return o
}
func (r *runtime) setHidden(o *goja.Object, s *goja.Symbol, v any) *goja.Object {
	if err := o.SetSymbol(s, v); err != nil {
		r.throw(err)
	}
	return o
}
func (r *runtime) hidden(v goja.Value, s *goja.Symbol, label string) any {
	o := v.ToObject(r.vm)
	if o == nil {
		r.throwType("RAG_V2_TYPE", label+" is required")
	}
	x := o.GetSymbol(s)
	if x == nil || goja.IsUndefined(x) {
		r.throwType("RAG_V2_TYPE", "expected "+label)
	}
	return x.Export()
}
func (r *runtime) descriptorValue(v goja.Value) *ragmodel.Descriptor {
	x, ok := r.hidden(v, r.descriptor, "descriptor").(*ragmodel.Descriptor)
	if !ok {
		r.throwType("RAG_V2_TYPE", "expected descriptor")
	}
	return x
}
func (r *runtime) descriptorObject(v *ragmodel.Descriptor) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.descriptor, v)
}
func (r *runtime) corpusObject(v ragmodel.CorpusInput) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.corpus, &v)
}
func (r *runtime) datasetObject(v ragmodel.DatasetRef) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.dataset, &v)
}
func (r *runtime) factorObject(v ragmodel.FactorRef) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.factor, &v)
}

func (r *runtime) inputsObject() *goja.Object {
	return r.obj(map[string]any{"corpus": func(call goja.FunctionCall) goja.Value {
		return r.corpusObject(ragmodel.Corpus(call.Argument(0).String()))
	}})
}
func (r *runtime) unitsObject() *goja.Object {
	return r.obj(map[string]any{"identity": func(goja.FunctionCall) goja.Value { return r.descriptorObject(ragmodel.UnitsIdentity()) }, "individualTurns": func(goja.FunctionCall) goja.Value { return r.descriptorObject(ragmodel.IndividualTurns()) }})
}
func (r *runtime) transcriptObject() *goja.Object {
	o := r.vm.NewObject()
	_ = o.Set("units", r.obj(map[string]any{"agentsViewRuns": func(goja.FunctionCall) goja.Value { return r.descriptorObject(ragmodel.AgentsViewRuns()) }}))
	return o
}
func (r *runtime) chunksObject() *goja.Object {
	return r.obj(map[string]any{"recursive": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.RecursiveChunkConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.RecursiveChunks(c))
	}})
}
func (r *runtime) generationObject() *goja.Object {
	return r.obj(map[string]any{"structured": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.StructuredGenerationConfig
		r.decode(call.Argument(1), &c)
		if c.Model == "" {
			c.Model = call.Argument(0).String()
		}
		return r.descriptorObject(ragmodel.StructuredGenerator(call.Argument(0).String(), c))
	}, "answer": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.AnswerConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.Answer(c))
	}})
}
func (r *runtime) representationsObject() *goja.Object {
	return r.obj(map[string]any{"raw": func(call goja.FunctionCall) goja.Value {
		return r.descriptorObject(ragmodel.RawRepresentation(call.Argument(0).String()))
	}, "structuredSummary": func(call goja.FunctionCall) goja.Value {
		options := call.Argument(1).ToObject(r.vm)
		generator := r.descriptorValue(options.Get("generator"))
		var c ragmodel.StructuredGenerationConfig
		r.decodeRaw(generator.Config, &c)
		return r.descriptorObject(ragmodel.StructuredSummary(call.Argument(0).String(), ragmodel.StructuredSummaryConfig{Generator: c}))
	}, "syntheticQuestions": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.SyntheticQuestionsConfig
		r.decode(call.Argument(1), &c)
		return r.descriptorObject(ragmodel.SyntheticQuestions(call.Argument(0).String(), c))
	}, "compose": func(call goja.FunctionCall) goja.Value {
		values := make([]*ragmodel.Descriptor, len(call.Arguments))
		for i, v := range call.Arguments {
			values[i] = r.descriptorValue(v)
		}
		return r.descriptorObject(ragmodel.ComposeRepresentations(values...))
	}})
}
func (r *runtime) embeddingsObject() *goja.Object {
	return r.obj(map[string]any{"model": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.EmbeddingConfig
		if len(call.Arguments) > 1 {
			r.decode(call.Argument(1), &c)
		}
		return r.descriptorObject(ragmodel.EmbeddingModel(call.Argument(0).String(), c))
	}})
}
func (r *runtime) indexesObject() *goja.Object {
	return r.obj(map[string]any{"bleveMulti": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.BleveMultiConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.BleveMulti(c))
	}})
}
func (r *runtime) retrieveObject() *goja.Object {
	factory := func(vector bool) func(goja.FunctionCall) goja.Value {
		return func(call goja.FunctionCall) goja.Value {
			var c ragmodel.RetrieveConfig
			r.decode(call.Argument(1), &c)
			if vector {
				return r.descriptorObject(ragmodel.Vector(call.Argument(0).String(), c))
			}
			return r.descriptorObject(ragmodel.BM25(call.Argument(0).String(), c))
		}
	}
	return r.obj(map[string]any{"bm25": factory(false), "vector": factory(true)})
}
func (r *runtime) collapseObject() *goja.Object {
	return r.obj(map[string]any{"parent": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.CollapseConfig
		r.decodeCollapse(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.ParentCollapse(c))
	}})
}
func (r *runtime) fusionObject() *goja.Object {
	return r.obj(map[string]any{"weightedRRF": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.WeightedRRFConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.WeightedRRF(c))
	}})
}
func (r *runtime) hydrationObject() *goja.Object {
	return r.obj(map[string]any{"sourceEvidence": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.HydrationConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.SourceEvidence(c))
	}})
}
func (r *runtime) rerankObject() *goja.Object {
	return r.obj(map[string]any{"crossEncoder": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.CrossEncoderConfig
		r.decode(call.Argument(0), &c)
		return r.descriptorObject(ragmodel.CrossEncoder(c))
	}})
}
func (r *runtime) metricsObject() *goja.Object {
	return r.obj(map[string]any{"mrr": func(goja.FunctionCall) goja.Value { return r.vm.ToValue("rag.mrr/v1") }})
}
func (r *runtime) datasetsObject() *goja.Object {
	return r.obj(map[string]any{"artifact": func(call goja.FunctionCall) goja.Value {
		var c struct {
			Split           string `json:"split"`
			Status          string `json:"status"`
			RelevanceTarget string `json:"relevanceTarget"`
		}
		r.decode(call.Argument(1), &c)
		return r.datasetObject(ragmodel.DatasetArtifact(call.Argument(0).String(), c.Split, c.Status, c.RelevanceTarget))
	}})
}
func (r *runtime) recipesObject() *goja.Object {
	return r.obj(map[string]any{"transcriptPreparation": func(call goja.FunctionCall) goja.Value {
		var c ragmodel.RecursiveChunkConfig
		r.decode(call.Argument(0), &c)
		return r.fragmentObject(ragmodel.NewFragment("transcript-preparation", func(p *ragmodel.PipelineBuilder) {
			p.Units(ragmodel.AgentsViewRuns()).Chunks(ragmodel.RecursiveChunks(c))
		}))
	}})
}

func (r *runtime) pipelineFactory(call goja.FunctionCall) goja.Value {
	builder := ragmodel.NewPipelineBuilder(call.Argument(0).String())
	o := r.pipelineBuilderObject(builder)
	if len(call.Arguments) > 1 && !goja.IsUndefined(call.Argument(1)) {
		r.configure(call.Argument(1), o)
	}
	return r.pipelineObject(builder.Build())
}

func (r *runtime) pipelineBuilderObject(b *ragmodel.PipelineBuilder) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.pipelineBuilder, b)
	chain := func(name string, fn func(goja.FunctionCall)) {
		modules.SetExport(o, ModuleName, name, func(call goja.FunctionCall) goja.Value { fn(call); return o })
	}
	chain("corpus", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.corpus, "corpus input").(*ragmodel.CorpusInput)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected corpus input")
		}
		b.CorpusInput(*v)
	})
	chain("units", func(c goja.FunctionCall) { b.Units(r.descriptorValue(c.Argument(0))) })
	chain("chunks", func(c goja.FunctionCall) { b.Chunks(r.descriptorValue(c.Argument(0))) })
	chain("representations", func(c goja.FunctionCall) { b.Represent(r.descriptorValue(c.Argument(0))) })
	chain("embedding", func(c goja.FunctionCall) { b.EmbeddingModel(r.descriptorValue(c.Argument(0))) })
	chain("index", func(c goja.FunctionCall) { b.IndexNamed(c.Argument(0).String(), r.descriptorValue(c.Argument(1))) })
	chain("use", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.fragment, "fragment").(*ragmodel.Fragment)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected fragment")
		}
		b.Use(v)
	})
	chain("note", func(c goja.FunctionCall) { b.Note(c.Argument(0).String()) })
	chain("tag", func(c goja.FunctionCall) { b.Tag(c.Argument(0).String(), c.Argument(1).String()) })
	return o
}
func (r *runtime) pipelineObject(v *ragmodel.Pipeline) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.pipeline, v)
	modules.SetExport(o, ModuleName, "validate", func(goja.FunctionCall) goja.Value { return r.validationValue(ragmodel.ValidatePipeline(v)) })
	modules.SetExport(o, ModuleName, "explain", func(goja.FunctionCall) goja.Value {
		x, err := ragmodel.Explain(v)
		if err != nil {
			r.throw(err)
		}
		return r.jsonValue(x)
	})
	return o
}
func (r *runtime) fragmentFactory(c goja.FunctionCall) goja.Value {
	b := ragmodel.NewPipelineBuilder(c.Argument(0).String())
	r.configure(c.Argument(1), r.pipelineBuilderObject(b))
	return r.fragmentObject(ragmodel.FragmentFromPipeline(c.Argument(0).String(), b.Build()))
}
func (r *runtime) fragmentObject(v *ragmodel.Fragment) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.fragment, v)
}
func (r *runtime) queryFactory(c goja.FunctionCall) goja.Value {
	b := ragmodel.NewQueryBuilder(c.Argument(0).String())
	r.configure(c.Argument(1), r.queryBuilderObject(b))
	return r.queryObject(b.Build())
}
func (r *runtime) queryBuilderObject(b *ragmodel.QueryBuilder) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.queryBuilder, b)
	chain := func(name string, fn func(goja.FunctionCall)) {
		modules.SetExport(o, ModuleName, name, func(c goja.FunctionCall) goja.Value { fn(c); return o })
	}
	chain("channels", func(c goja.FunctionCall) { b.Channels(r.descriptorArray(c.Argument(0))...) })
	chain("collapseChannels", func(c goja.FunctionCall) { b.CollapseChannels(r.descriptorValue(c.Argument(0))) })
	chain("fuse", func(c goja.FunctionCall) { b.Fuse(r.descriptorValue(c.Argument(0))) })
	chain("collapseFinal", func(c goja.FunctionCall) { b.CollapseFinal(r.descriptorValue(c.Argument(0))) })
	chain("hydrate", func(c goja.FunctionCall) { b.Hydrate(r.descriptorValue(c.Argument(0))) })
	chain("results", func(c goja.FunctionCall) { b.ResultCount(int(c.Argument(0).ToInteger())) })
	return o
}
func (r *runtime) queryObject(v *ragmodel.QueryPlan) *goja.Object {
	return r.setHidden(r.vm.NewObject(), r.query, v)
}
func (r *runtime) variantFactory(c goja.FunctionCall) goja.Value {
	value := ragmodel.NewVariant(c.Argument(0).String(), func(b *ragmodel.VariantBuilder) {
		r.configure(c.Argument(1), r.variantBuilderObject(b))
	})
	return r.setHidden(r.vm.NewObject(), r.variant, value)
}
func (r *runtime) variantBuilderObject(b *ragmodel.VariantBuilder) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.variantBuilder, b)
	chain := func(name string, fn func(goja.FunctionCall)) {
		modules.SetExport(o, ModuleName, name, func(c goja.FunctionCall) goja.Value { fn(c); return o })
	}
	chain("selectRepresentations", func(c goja.FunctionCall) { b.SelectRepresentations(r.stringArray(c.Argument(0))...) })
	chain("query", func(c goja.FunctionCall) {
		value := c.Argument(0)
		if fn, ok := goja.AssertFunction(value); ok {
			ctx := r.obj(map[string]any{"factor": func(x goja.FunctionCall) goja.Value {
				return r.factorObject(ragmodel.FactorRef{ID: x.Argument(0).String()})
			}})
			returned, err := fn(goja.Undefined(), ctx)
			if err != nil {
				panic(err)
			}
			value = returned
		}
		q, ok := r.hidden(value, r.query, "query plan").(*ragmodel.QueryPlan)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected query plan")
		}
		b.QueryPlan(q)
	})
	chain("rerank", func(c goja.FunctionCall) { b.Rerank(r.descriptorValue(c.Argument(0))) })
	chain("generate", func(c goja.FunctionCall) { b.Generate(r.descriptorValue(c.Argument(0))) })
	return o
}

func (r *runtime) productFactory(c goja.FunctionCall) goja.Value {
	b := ragmodel.NewProductBuilder(c.Argument(0).String())
	r.configure(c.Argument(1), r.productBuilderObject(b))
	return r.productObject(b.Build())
}
func (r *runtime) productBuilderObject(b *ragmodel.ProductBuilder) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.productBuilder, b)
	chain := func(name string, fn func(goja.FunctionCall)) {
		modules.SetExport(o, ModuleName, name, func(c goja.FunctionCall) goja.Value { fn(c); return o })
	}
	chain("pipeline", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.pipeline, "pipeline").(*ragmodel.Pipeline)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected pipeline")
		}
		b.PipelineValue(v)
	})
	chain("query", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.query, "query plan").(*ragmodel.QueryPlan)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected query plan")
		}
		b.QueryPlan(v)
	})
	chain("rerank", func(c goja.FunctionCall) { b.Rerank(r.descriptorValue(c.Argument(0))) })
	chain("generate", func(c goja.FunctionCall) { b.Generate(r.descriptorValue(c.Argument(0))) })
	chain("request", func(c goja.FunctionCall) {
		b.RequestContract(func(x *ragmodel.RequestBuilder) { r.configure(c.Argument(0), r.requestBuilderObject(x)) })
	})
	chain("response", func(c goja.FunctionCall) {
		b.ResponseContract(func(x *ragmodel.ResponseBuilder) { r.configure(c.Argument(0), r.responseBuilderObject(x)) })
	})
	chain("runtime", func(c goja.FunctionCall) {
		b.RuntimePolicy(func(x *ragmodel.RuntimeBuilder) { r.configure(c.Argument(0), r.runtimeBuilderObject(x)) })
	})
	chain("tag", func(c goja.FunctionCall) { b.Tag(c.Argument(0).String(), c.Argument(1).String()) })
	return o
}
func (r *runtime) requestBuilderObject(b *ragmodel.RequestBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "field", func(c goja.FunctionCall) goja.Value {
		var x struct {
			Required  bool `json:"required"`
			MaxLength int  `json:"maxLength"`
		}
		if len(c.Arguments) > 2 {
			r.decode(c.Argument(2), &x)
		}
		b.Field(c.Argument(0).String(), c.Argument(1).String(), x.Required, x.MaxLength)
		return o
	})
	return o
}
func (r *runtime) responseBuilderObject(b *ragmodel.ResponseBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "answer", func(c goja.FunctionCall) goja.Value { b.Answer(c.Argument(0).String()); return o })
	modules.SetExport(o, ModuleName, "citations", func(c goja.FunctionCall) goja.Value { b.Citations(c.Argument(0).String()); return o })
	modules.SetExport(o, ModuleName, "includeTraceId", func(c goja.FunctionCall) goja.Value { b.TraceID(c.Argument(0).ToBoolean()); return o })
	return o
}
func (r *runtime) runtimeBuilderObject(b *ragmodel.RuntimeBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "timeoutMs", func(c goja.FunctionCall) goja.Value { b.Timeout(c.Argument(0).ToInteger()); return o })
	modules.SetExport(o, ModuleName, "maxConcurrent", func(c goja.FunctionCall) goja.Value { b.Concurrent(int(c.Argument(0).ToInteger())); return o })
	modules.SetExport(o, ModuleName, "onProviderFailure", func(c goja.FunctionCall) goja.Value { b.ProviderFailure(c.Argument(0).String()); return o })
	modules.SetExport(o, ModuleName, "trace", func(c goja.FunctionCall) goja.Value { b.Trace(c.Argument(0).String()); return o })
	return o
}
func (r *runtime) productObject(v *ragmodel.Product) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.product, v)
	modules.SetExport(o, ModuleName, "validate", func(goja.FunctionCall) goja.Value { return r.validationValue(validateTarget(v)) })
	modules.SetExport(o, ModuleName, "explain", func(goja.FunctionCall) goja.Value {
		x, err := ragmodel.Explain(v)
		if err != nil {
			r.throw(err)
		}
		return r.jsonValue(x)
	})
	modules.SetExport(o, ModuleName, "compileProduct", func(c goja.FunctionCall) goja.Value { return r.compileProductValue(v, c.Argument(0)) })
	return o
}

func (r *runtime) studyFactory(c goja.FunctionCall) goja.Value {
	b := ragmodel.NewStudyBuilder(c.Argument(0).String())
	r.configure(c.Argument(1), r.studyBuilderObject(b))
	return r.studyObject(b.Build())
}
func (r *runtime) studyBuilderObject(b *ragmodel.StudyBuilder) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.studyBuilder, b)
	chain := func(name string, fn func(goja.FunctionCall)) {
		modules.SetExport(o, ModuleName, name, func(c goja.FunctionCall) goja.Value { fn(c); return o })
	}
	chain("pipeline", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.pipeline, "pipeline").(*ragmodel.Pipeline)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected pipeline")
		}
		b.PipelineValue(v)
	})
	chain("dataset", func(c goja.FunctionCall) {
		v, ok := r.hidden(c.Argument(0), r.dataset, "dataset").(*ragmodel.DatasetRef)
		if !ok {
			r.throwType("RAG_V2_TYPE", "expected dataset")
		}
		b.DatasetRef(*v)
	})
	chain("variants", func(c goja.FunctionCall) {
		b.VariantsList(func(x *ragmodel.VariantsBuilder) { r.configure(c.Argument(0), r.variantsBuilderObject(x)) })
	})
	chain("factors", func(c goja.FunctionCall) {
		b.FactorsList(func(x *ragmodel.FactorsBuilder) { r.configure(c.Argument(0), r.factorsBuilderObject(x)) })
	})
	chain("replicates", func(c goja.FunctionCall) { b.ReplicateCount(int(c.Argument(0).ToInteger())) })
	chain("metrics", func(c goja.FunctionCall) {
		b.MetricsList(func(x *ragmodel.MetricsBuilder) { r.configure(c.Argument(0), r.metricBuilderObject(x)) })
	})
	chain("invariants", func(c goja.FunctionCall) {
		b.InvariantsList(func(x *ragmodel.InvariantsBuilder) { r.configure(c.Argument(0), r.invariantsBuilderObject(x)) })
	})
	chain("tag", func(c goja.FunctionCall) { b.Tag(c.Argument(0).String(), c.Argument(1).String()) })
	return o
}
func (r *runtime) variantsBuilderObject(b *ragmodel.VariantsBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "add", func(c goja.FunctionCall) goja.Value {
		b.Add(c.Argument(0).String(), func(x *ragmodel.VariantBuilder) { r.configure(c.Argument(1), r.variantBuilderObject(x)) })
		return o
	})
	return o
}
func (r *runtime) factorsBuilderObject(b *ragmodel.FactorsBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "enum", func(c goja.FunctionCall) goja.Value {
		b.Enum(c.Argument(0).String(), r.stringArray(c.Argument(1))...)
		return o
	})
	return o
}
func (r *runtime) metricBuilderObject(b *ragmodel.MetricsBuilder) *goja.Object {
	o := r.vm.NewObject()
	setCut := func(name string, fn func([]int)) {
		modules.SetExport(o, ModuleName, name, func(c goja.FunctionCall) goja.Value { fn(r.intArray(c.Argument(0))); return o })
	}
	setCut("precisionAt", func(v []int) { b.PrecisionAt(v) })
	setCut("recallAt", func(v []int) { b.RecallAt(v) })
	setCut("hitRateAt", func(v []int) { b.HitRateAt(v) })
	setCut("ndcgAt", func(v []int) { b.NDCGAt(v) })
	modules.SetExport(o, ModuleName, "mrr", func(goja.FunctionCall) goja.Value { b.MRR(); return o })
	modules.SetExport(o, ModuleName, "latency", func(c goja.FunctionCall) goja.Value { b.Latency(r.stringArray(c.Argument(0))); return o })
	modules.SetExport(o, ModuleName, "tokenUsage", func(goja.FunctionCall) goja.Value { b.TokenUsage(); return o })
	modules.SetExport(o, ModuleName, "providerCost", func(goja.FunctionCall) goja.Value { b.ProviderCost(); return o })
	modules.SetExport(o, ModuleName, "storageBytes", func(goja.FunctionCall) goja.Value { b.StorageBytes(); return o })
	modules.SetExport(o, ModuleName, "failureRates", func(goja.FunctionCall) goja.Value { b.FailureRates(); return o })
	return o
}
func (r *runtime) invariantsBuilderObject(b *ragmodel.InvariantsBuilder) *goja.Object {
	o := r.vm.NewObject()
	modules.SetExport(o, ModuleName, "require", func(c goja.FunctionCall) goja.Value { b.Require(c.Argument(0).String()); return o })
	return o
}
func (r *runtime) studyObject(v *ragmodel.Study) *goja.Object {
	o := r.setHidden(r.vm.NewObject(), r.study, v)
	modules.SetExport(o, ModuleName, "validate", func(goja.FunctionCall) goja.Value { return r.validationValue(validateTarget(v)) })
	modules.SetExport(o, ModuleName, "explain", func(goja.FunctionCall) goja.Value {
		x, err := ragmodel.Explain(v)
		if err != nil {
			r.throw(err)
		}
		return r.jsonValue(x)
	})
	modules.SetExport(o, ModuleName, "compileStudy", func(c goja.FunctionCall) goja.Value { return r.compileStudyValue(v, c.Argument(0)) })
	return o
}

func (r *runtime) validate(c goja.FunctionCall) goja.Value {
	return r.validationValue(validateTarget(r.target(c.Argument(0))))
}
func validateTarget(v any) error {
	switch x := v.(type) {
	case *ragmodel.Pipeline:
		return ragmodel.ValidatePipeline(x)
	case *ragmodel.Product:
		return ragmodel.ValidateProduct(x)
	case *ragmodel.Study:
		return ragmodel.ValidateStudy(x)
	default:
		return fmt.Errorf("RAG_V2_TYPE")
	}
}
func (r *runtime) explain(c goja.FunctionCall) goja.Value {
	x, err := ragmodel.Explain(r.target(c.Argument(0)))
	if err != nil {
		r.throw(err)
	}
	return r.jsonValue(x)
}
func (r *runtime) compileProduct(c goja.FunctionCall) goja.Value {
	v, ok := r.target(c.Argument(0)).(*ragmodel.Product)
	if !ok {
		r.throwType("RAG_V2_TYPE", "expected product")
	}
	return r.compileProductValue(v, c.Argument(1))
}
func (r *runtime) compileProductValue(v *ragmodel.Product, options goja.Value) goja.Value {
	var o ragmodel.CompileOptions
	r.decode(options, &o)
	value, err := ragmodel.CompileProduct(v, o)
	if err != nil {
		r.throw(err)
	}
	return r.jsonValue(value)
}
func (r *runtime) compileStudy(c goja.FunctionCall) goja.Value {
	v, ok := r.target(c.Argument(0)).(*ragmodel.Study)
	if !ok {
		r.throwType("RAG_V2_TYPE", "expected study")
	}
	return r.compileStudyValue(v, c.Argument(1))
}
func (r *runtime) compileStudyValue(v *ragmodel.Study, options goja.Value) goja.Value {
	var o ragmodel.CompileOptions
	r.decode(options, &o)
	study, _, err := ragmodel.CompileStudy(v, o)
	if err != nil {
		r.throw(err)
	}
	return r.jsonValue(study)
}
func (r *runtime) preview(c goja.FunctionCall) goja.Value {
	v, ok := r.target(c.Argument(0)).(*ragmodel.Study)
	if !ok {
		r.throwType("RAG_V2_TYPE", "expected study")
	}
	var options struct {
		Inputs  map[string]ragcontract.ArtifactBinding `json:"inputs"`
		Variant string                                 `json:"variant"`
		Factors map[string]string                      `json:"factors"`
		Query   string                                 `json:"query"`
		Trace   string                                 `json:"trace"`
	}
	r.decode(c.Argument(1), &options)
	value, err := ragmodel.Preview(v, ragmodel.CompileOptions{Inputs: options.Inputs}, ragmodel.PreviewOptions{Variant: options.Variant, Factors: options.Factors, Query: options.Query, Trace: options.Trace})
	if err != nil {
		r.throw(err)
	}
	return r.jsonValue(value)
}
func (r *runtime) target(v goja.Value) any {
	o := v.ToObject(r.vm)
	for _, s := range []*goja.Symbol{r.pipeline, r.product, r.study} {
		if x := o.GetSymbol(s); x != nil && !goja.IsUndefined(x) {
			return x.Export()
		}
	}
	r.throwType("RAG_V2_TYPE", "expected pipeline, product, or study")
	return nil
}

func (r *runtime) descriptorArray(v goja.Value) []*ragmodel.Descriptor {
	o := v.ToObject(r.vm)
	n := int(o.Get("length").ToInteger())
	out := make([]*ragmodel.Descriptor, n)
	for i := 0; i < n; i++ {
		out[i] = r.descriptorValue(o.Get(fmt.Sprintf("%d", i)))
	}
	return out
}
func (r *runtime) stringArray(v goja.Value) []string {
	o := v.ToObject(r.vm)
	n := int(o.Get("length").ToInteger())
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = o.Get(fmt.Sprintf("%d", i)).String()
	}
	return out
}
func (r *runtime) intArray(v goja.Value) []int {
	o := v.ToObject(r.vm)
	n := int(o.Get("length").ToInteger())
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = int(o.Get(fmt.Sprintf("%d", i)).ToInteger())
	}
	return out
}
func (r *runtime) configure(v goja.Value, arg goja.Value) {
	fn, ok := goja.AssertFunction(v)
	if !ok {
		r.throwType("RAG_V2_CONFIGURATOR", "configurator must be a function")
	}
	if _, err := fn(goja.Undefined(), arg); err != nil {
		panic(err)
	}
}
func (r *runtime) decode(v goja.Value, target any) {
	if goja.IsUndefined(v) || goja.IsNull(v) {
		r.decodeRaw(json.RawMessage(`{}`), target)
		return
	}
	raw, err := json.Marshal(v.Export())
	if err != nil {
		r.throw(err)
	}
	r.decodeRaw(raw, target)
}
func (r *runtime) decodeRaw(raw json.RawMessage, target any) {
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		r.throwType("RAG_V2_OPTIONS", err.Error())
	}
}
func (r *runtime) decodeCollapse(v goja.Value, c *ragmodel.CollapseConfig) {
	o := v.ToObject(r.vm)
	scope := o.Get("scope")
	if f := scope.ToObject(r.vm).GetSymbol(r.factor); f != nil && !goja.IsUndefined(f) {
		c.Scope = f.Export()
	} else {
		c.Scope = scope.String()
	}
	rep := o.Get("representative")
	if rep != nil && !goja.IsUndefined(rep) {
		c.Representative = rep.String()
	}
}
func (r *runtime) jsonValue(v any) goja.Value {
	raw, err := json.Marshal(v)
	if err != nil {
		r.throw(err)
	}
	var plain any
	if err := json.Unmarshal(raw, &plain); err != nil {
		r.throw(err)
	}
	return r.vm.ToValue(plain)
}
func (r *runtime) validationValue(err error) goja.Value {
	if err == nil {
		return r.vm.ToValue(map[string]any{"ok": true, "issues": []any{}})
	}
	return r.vm.ToValue(map[string]any{"ok": false, "issues": []map[string]any{{"code": "RAG_V2_VALIDATION", "path": "$", "message": err.Error(), "severity": "error"}}})
}
func (r *runtime) throw(err error) {
	var validation *ragcontract.ValidationError
	if errors.As(err, &validation) {
		r.throwType("RAG_V2_VALIDATION", validation.Error())
	}
	panic(r.vm.NewGoError(err))
}
func (r *runtime) throwType(code, message string) {
	e := r.vm.NewTypeError(message)
	_ = e.Set("code", code)
	panic(e)
}
