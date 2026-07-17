package rag

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
)

const digestA = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const digestB = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

func newVM() *goja.Runtime {
	vm := goja.New()
	registry := require.NewRegistry()
	registry.RegisterNativeModule(ModuleName, NewLoader())
	registry.Enable(vm)
	return vm
}
func productScript() string {
	return `
const rag=require("rag");
const base=rag.pipeline("base",p=>p
 .corpus(rag.inputs.corpus("corpus"))
 .units(rag.units.identity())
 .chunks(rag.chunks.recursive({maxRunes:800,overlapSpans:120}))
 .representations(rag.representations.raw("raw"))
 .embedding(rag.embeddings.model("embed-v1",{dimensions:3,distance:"cosine",normalize:"l2",batchSize:8}))
 .index("representations",rag.indexes.bleveMulti({lexical:true,vector:{distance:"cosine",optimizeFor:"recall"}})));
const query=rag.queryPlan("raw-query",q=>q
 .channels([rag.retrieve.bm25("raw.lexical",{index:"representations",representation:"raw",topK:30}),rag.retrieve.vector("raw.vector",{index:"representations",representation:"raw",topK:30})])
 .collapseChannels(rag.collapse.parent({scope:"unit",representative:"scoreThenRepresentationId"}))
 .fuse(rag.fusion.weightedRRF({rankConstant:60,weights:{"raw.vector":2}}))
 .collapseFinal(rag.collapse.parent({scope:"unit",representative:"bestFusionContributionThenId"}))
 .hydrate(rag.hydration.sourceEvidence({selection:"bestContributionThenId"})).results(5));
const product=rag.product("assistant",p=>p.pipeline(base).query(query)
 .request(r=>r.field("query","string",{required:true,maxLength:4096}))
 .response(r=>r.answer("markdown").citations("source").includeTraceId(true))
 .runtime(r=>r.timeoutMs(15000).maxConcurrent(16).onProviderFailure("fail").trace("authoritative")));
`
}
func compileOptionsJS() string {
	return `{inputs:{corpus:{role:"corpus",kind:"manifest",digest:"` + digestA + `",sizeBytes:10,schemaVersion:"rag-corpus-snapshot-manifest/v2"}}}`
}

func TestModuleCompilesPureProductV2(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(productScript() + `product.compileProduct(` + compileOptionsJS() + `);`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	if got["schemaVersion"] != "rag-product-plan/v2" {
		t.Fatalf("%#v", got)
	}
	pipeline := got["pipeline"].(map[string]any)
	if len(pipeline["nodes"].([]any)) < 10 {
		t.Fatalf("nodes=%#v", pipeline["nodes"])
	}
	if got["request"].(map[string]any)["fields"] == nil {
		t.Fatalf("request=%#v", got["request"])
	}
}

func TestModuleCompilesFiveVariantStudyAndPreview(t *testing.T) {
	vm := newVM()
	script := productScript() + `
const variants={raw:["raw"],summary:["summary"],rawSummary:["raw","summary"],rawQuestion:["raw","question"],all:["raw","summary","question"]};
const study=rag.study("matrix",s=>s.pipeline(base)
 .dataset(rag.datasets.artifact("judgments",{split:"smoke",status:"candidate",relevanceTarget:"unit"}))
 .variants(v=>{for(const [name,kinds] of Object.entries(variants)){v.add(name,x=>x.selectRepresentations(kinds).query(ctx=>rag.queryPlan(name,q=>q
  .channels(kinds.flatMap(kind=>[rag.retrieve.bm25(kind+".lexical",{index:"representations",representation:kind,topK:30}),rag.retrieve.vector(kind+".vector",{index:"representations",representation:kind,topK:30})]))
  .collapseChannels(rag.collapse.parent({scope:ctx.factor("collapse"),representative:"scoreThenRepresentationId"}))
  .fuse(rag.fusion.weightedRRF({rankConstant:60})).collapseFinal(rag.collapse.parent({scope:ctx.factor("collapse"),representative:"bestFusionContributionThenId"}))
  .hydrate(rag.hydration.sourceEvidence({selection:"bestContributionThenId"})).results(5))))}})
 .factors(f=>f.enum("collapse",["chunk","unit"])).replicates(3)
 .metrics(m=>m.precisionAt([5]).recallAt([5]).hitRateAt([5]).mrr().ndcgAt([5]).latency(["query"]).tokenUsage().providerCost().storageBytes().failureRates())
 .invariants(i=>i.require("source-hydrated-final-hit/v1")).tag("evaluationStatus","candidate"));
const options={inputs:{corpus:{role:"corpus",kind:"manifest",digest:"` + digestA + `",schemaVersion:"rag-corpus-snapshot-manifest/v2"},judgments:{role:"judgments",kind:"manifest",digest:"` + digestB + `",schemaVersion:"rag-evaluation-dataset-manifest/v2"}}};
({compiled:study.compileStudy(options),explain:rag.explain(study),preview:rag.preview(study,{...options,variant:"all",factors:{collapse:"unit"},query:"why",trace:"full"})});`
	value, err := vm.RunString(script)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	compiled := got["compiled"].(map[string]any)
	if compiled["schemaVersion"] != "rag-study/v2" || len(compiled["variants"].([]any)) != 5 {
		t.Fatalf("compiled=%#v", compiled)
	}
	explain := got["explain"].(map[string]any)
	if explain["cellCount"] != int64(10) && explain["cellCount"] != float64(10) {
		t.Fatalf("explain=%#v", explain)
	}
	preview := got["preview"].(map[string]any)
	if preview["schemaVersion"] != "rag-preview-request/v1" {
		t.Fatalf("preview=%#v", preview)
	}
}

func TestConfiguratorsExecuteImmediatelyAndNoFunctionsReachOutput(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(productScript() + `let calls=0;const f=rag.fragment("f",p=>{calls++;p.note("now")});const p=rag.pipeline("fragmented",p=>p.corpus(rag.inputs.corpus("corpus")).use(f).units(rag.units.identity()).chunks(rag.chunks.recursive({maxRunes:10})).representations(rag.representations.raw("raw")).index("i",rag.indexes.bleveMulti({lexical:true})));({calls,keys:Object.keys(f),json:JSON.stringify(product.compileProduct(` + compileOptionsJS() + `))});`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	if got["calls"] != int64(1) {
		t.Fatalf("calls=%#v", got["calls"])
	}
	if len(got["keys"].([]any)) != 0 {
		t.Fatalf("hidden keys=%#v", got["keys"])
	}
	if strings.Contains(got["json"].(string), "function") {
		t.Fatalf("function escaped: %s", got["json"])
	}
}

func TestRetiredV1LifecycleAndDescriptorNamesAreAbsent(t *testing.T) {
	vm := newVM()
	value, err := vm.RunString(`const rag=require("rag");({version:rag.version,experiment:typeof rag.experiment,artifact:typeof rag.artifact,grade:typeof rag.grade,open:typeof rag.open,rawChunks:typeof rag.representations.rawChunks,summaries:typeof rag.representations.summaries,questions:typeof rag.representations.questions,exportSpecification:typeof rag.exportSpecification});`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	if got["version"] != "v2" {
		t.Fatalf("%#v", got)
	}
	for k, v := range got {
		if k != "version" && v != "undefined" {
			t.Fatalf("retired %s=%v", k, v)
		}
	}
}

func TestConfiguratorExceptionIsPreserved(t *testing.T) {
	_, err := newVM().RunString(`require("rag").pipeline("x",()=>{throw new Error("configurator exploded")})`)
	if err == nil || !strings.Contains(err.Error(), "configurator exploded") {
		t.Fatalf("%v", err)
	}
}

func TestEngineRegistrarCanRequireV2Module(t *testing.T) {
	factory, err := engine.NewRuntimeFactoryBuilder().WithModules(NewRegistrar()).Build()
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := factory.NewRuntime()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close(context.Background()) }()
	value, err := runtime.VM.RunString(`const rag=require("rag");({pipeline:typeof rag.pipeline,study:typeof rag.study,version:rag.version})`)
	if err != nil {
		t.Fatal(err)
	}
	got := value.Export().(map[string]any)
	if got["pipeline"] != "function" || got["study"] != "function" || got["version"] != "v2" {
		t.Fatalf("%#v", got)
	}
}

func TestJavaScriptAndPureGoCompileByteIdentically(t *testing.T) {
	vm := newVM()
	js, err := vm.RunString(productScript() + `JSON.stringify(product.compileProduct(` + compileOptionsJS() + `));`)
	if err != nil {
		t.Fatal(err)
	}
	pipeline := ragmodel.NewPipeline("base", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 800, OverlapSpans: 120})).Represent(ragmodel.RawRepresentation("raw")).EmbeddingModel(ragmodel.EmbeddingModel("embed-v1", ragmodel.EmbeddingConfig{Dimensions: 3, Distance: "cosine", Normalize: "l2", BatchSize: 8})).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true, Vector: &ragmodel.VectorIndexConfig{Distance: "cosine", OptimizeFor: "recall"}}))
	})
	query := ragmodel.NewQueryPlan("raw-query", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 30}), ragmodel.Vector("raw.vector", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 30})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60, Weights: map[string]float64{"raw.vector": 2}})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	product := ragmodel.NewProduct("assistant", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).RequestContract(func(r *ragmodel.RequestBuilder) { r.Field("query", "string", true, 4096) }).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Answer("markdown").Citations("source").TraceID(true) }).RuntimePolicy(func(r *ragmodel.RuntimeBuilder) {
			r.Timeout(15000).Concurrent(16).ProviderFailure("fail").Trace("authoritative")
		})
	})
	size := int64(10)
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "manifest", Digest: digestA, SizeBytes: &size, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	goJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatal(err)
	}
	var jsValue any
	if err := json.Unmarshal([]byte(js.String()), &jsValue); err != nil {
		t.Fatal(err)
	}
	var goValue any
	if err := json.Unmarshal(goJSON, &goValue); err != nil {
		t.Fatal(err)
	}
	jsCanonical, _ := ragcontract.CanonicalJSON(jsValue)
	goCanonical, _ := ragcontract.CanonicalJSON(goValue)
	if string(jsCanonical) != string(goCanonical) {
		t.Fatalf("JS and Go differ\nJS=%s\nGo=%s", jsCanonical, goCanonical)
	}
}

func TestRuntimeAndDeclarationTopLevelAPIParity(t *testing.T) {
	value, err := newVM().RunString(`Object.keys(require("rag")).sort()`)
	if err != nil {
		t.Fatal(err)
	}
	declaration := strings.Join(TypeScriptModule().RawDTS, "\n")
	for _, exported := range value.Export().([]any) {
		name := exported.(string)
		if name == "version" {
			if !strings.Contains(declaration, "export const version") {
				t.Fatal("version missing from declaration")
			}
			continue
		}
		if !strings.Contains(declaration, "export function "+name) && !strings.Contains(declaration, "export const "+name) {
			t.Fatalf("runtime export %q missing from declaration", name)
		}
	}
}

func TestJavaScriptAndPureGoStudyCompileByteIdentically(t *testing.T) {
	vm := newVM()
	script := productScript() + `
const factorQuery = rag.queryPlan("factor", q => q
 .channels([rag.retrieve.bm25("raw.lexical", {index:"representations", representation:"raw", topK:10})])
 .collapseChannels(rag.collapse.parent({scope:"unit", representative:"scoreThenRepresentationId"}))
 .fuse(rag.fusion.weightedRRF({rankConstant:60}))
 .collapseFinal(rag.collapse.parent({scope:"unit", representative:"bestFusionContributionThenId"}))
 .hydrate(rag.hydration.sourceEvidence({selection:"bestContributionThenId"})));
const study = rag.study("one", s => s.pipeline(base)
 .dataset(rag.datasets.artifact("judgments", {split:"smoke", status:"candidate", relevanceTarget:"unit"}))
 .variants(v => v.add("raw", x => x.selectRepresentations(["raw"]).query(factorQuery)))
 .replicates(2).metrics(m => m.mrr()));
JSON.stringify(study.compileStudy({inputs:{corpus:{role:"corpus",kind:"manifest",digest:"` + digestA + `",schemaVersion:"rag-corpus-snapshot-manifest/v2"},judgments:{role:"judgments",kind:"manifest",digest:"` + digestB + `",schemaVersion:"rag-evaluation-dataset-manifest/v2"}}}));`
	js, err := vm.RunString(script)
	if err != nil {
		t.Fatal(err)
	}
	pipeline := ragmodel.NewPipeline("base", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 800, OverlapSpans: 120})).Represent(ragmodel.RawRepresentation("raw")).EmbeddingModel(ragmodel.EmbeddingModel("embed-v1", ragmodel.EmbeddingConfig{Dimensions: 3, Distance: "cosine", Normalize: "l2", BatchSize: 8})).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true, Vector: &ragmodel.VectorIndexConfig{Distance: "cosine", OptimizeFor: "recall"}}))
	})
	query := ragmodel.NewQueryPlan("factor", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"}))
	})
	study := ragmodel.NewStudy("one", func(s *ragmodel.StudyBuilder) {
		s.PipelineValue(pipeline).DatasetRef(ragmodel.DatasetArtifact("judgments", "smoke", "candidate", "unit")).VariantsList(func(v *ragmodel.VariantsBuilder) {
			v.Add("raw", func(x *ragmodel.VariantBuilder) { x.SelectRepresentations("raw").QueryPlan(query) })
		}).ReplicateCount(2).MetricsList(func(m *ragmodel.MetricsBuilder) { m.MRR() })
	})
	compiled, _, err := ragmodel.CompileStudy(study, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "manifest", Digest: digestA, SchemaVersion: ragcontract.CorpusManifestSchema}, "judgments": {Role: "judgments", Kind: "manifest", Digest: digestB, SchemaVersion: ragcontract.EvaluationManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	goJSON, _ := json.Marshal(compiled)
	var a, b any
	if err := json.Unmarshal([]byte(js.String()), &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(goJSON, &b); err != nil {
		t.Fatal(err)
	}
	ac, _ := ragcontract.CanonicalJSON(a)
	bc, _ := ragcontract.CanonicalJSON(b)
	if string(ac) != string(bc) {
		t.Fatalf("JS and Go studies differ\nJS=%s\nGo=%s", ac, bc)
	}
}

func TestRunnableExamples(t *testing.T) {
	for _, name := range []string{"01-product.js", "02-five-variant-study.js", "03-fragment.js", "04-explain.js", "05-preview.js"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join("..", "..", "..", "examples", "rag-v2", name)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			vm := newVM()
			if _, err := vm.RunString(`const module = { exports: {} };`); err != nil {
				t.Fatal(err)
			}
			if _, err := vm.RunScript(path, string(data)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTypeScriptDeclaresRuntimeParityWithoutSemanticAny(t *testing.T) {
	d := strings.Join(TypeScriptModule().RawDTS, "\n")
	for _, want := range []string{"export function pipeline", "export function product", "export function study", "export function compileProduct", "export function compileStudy", "export function preview", "agentsViewRuns", "interface CrossEncoderConfig", "RAGPipelineIRV2"} {
		if !strings.Contains(d, want) {
			t.Fatalf("missing %q", want)
		}
	}
	for _, bad := range []string{"Record<string, any>", ": any", "export function experiment", "exportSpecification", "rawChunks", "summaries(", "questions("} {
		if strings.Contains(d, bad) {
			t.Fatalf("retired/imprecise %q", bad)
		}
	}
}
