package ragengine

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type collector struct {
	events    []ragoperators.Event
	traces    []ragcontract.QueryTrace
	metrics   []ragoperators.Metric
	artifacts []ragoperators.Artifact
	cancel    context.CancelFunc
}

func (c *collector) Event(_ context.Context, v ragoperators.Event) error {
	c.events = append(c.events, v)
	if c.cancel != nil && len(c.events) == 1 {
		c.cancel()
	}
	return nil
}
func (c *collector) Trace(_ context.Context, v ragcontract.QueryTrace) error {
	c.traces = append(c.traces, v)
	return nil
}
func (c *collector) Metric(_ context.Context, v ragoperators.Metric) error {
	c.metrics = append(c.metrics, v)
	return nil
}
func (c *collector) Artifact(_ context.Context, v ragoperators.Artifact) error {
	c.artifacts = append(c.artifacts, v)
	return nil
}
func TestEngineExecutesTopologicalRawBM25AndEmitsEvidence(t *testing.T) {
	execution := rawExecution(t)
	corpus, dataset := fixtureData()
	observer := &collector{}
	started := time.Now()
	result, err := New(nil).Execute(context.Background(), execution, corpus, dataset, observer, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed > 5*time.Second {
		t.Fatalf("raw BM25 smoke budget exceeded: %s", elapsed)
	}
	artifactBytes := 0
	for _, artifact := range result.Artifacts {
		artifactBytes += len(artifact.Data)
	}
	if artifactBytes > 1<<20 {
		t.Fatalf("raw BM25 smoke storage budget exceeded: %d", artifactBytes)
	}
	materialized := map[string]bool{}
	for _, artifact := range observer.artifacts {
		if len(artifact.Metadata) == 0 {
			continue
		}
		var base ragcontract.ManifestBase
		if err := json.Unmarshal(artifact.Metadata, &base); err != nil {
			t.Fatal(err)
		}
		if err := ragcontract.ValidateManifestBase(base, artifact.SchemaVersion, true); err != nil {
			t.Fatalf("artifact %s: %v", artifact.Role, err)
		}
		materialized[artifact.Role] = true
	}
	for _, role := range []string{"unit-set", "chunk-set", "representation-set"} {
		if !materialized[role] {
			t.Fatalf("missing materialized %s", role)
		}
	}
	if len(result.Traces) != 1 || len(observer.traces) != 1 || len(observer.metrics) != 1 {
		t.Fatalf("result=%#v observer=%#v", result, observer)
	}
	trace := result.Traces[0]
	if len(trace.Results) == 0 || trace.Results[0].Evidence.ChunkID == "" || len(trace.Results[0].MatchedRepresentations) == 0 {
		t.Fatalf("trace=%#v", trace)
	}
	if len(observer.artifacts) < 3 {
		t.Fatalf("artifacts=%d", len(observer.artifacts))
	}
	if string(observer.metrics[0].Value) != "1" {
		t.Fatalf("mrr=%s", observer.metrics[0].Value)
	}
}
func TestEngineCancellationPreservesPartialEvents(t *testing.T) {
	execution := rawExecution(t)
	corpus, dataset := fixtureData()
	ctx, cancel := context.WithCancel(context.Background())
	observer := &collector{cancel: cancel}
	result, err := New(nil).Execute(ctx, execution, corpus, dataset, observer, Options{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
	if result == nil || len(observer.events) == 0 {
		t.Fatalf("partial evidence lost: %#v", observer)
	}
}

type acceptingSchema struct{}

func (acceptingSchema) Validate(string, json.RawMessage) error { return nil }

type failingGenerator struct{}

func (failingGenerator) Generate(context.Context, ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	return ragoperators.GenerationResult{}, errors.New("fixture generation failed")
}

func TestEnginePreparesStaticPipelineOnceForMultipleQueries(t *testing.T) {
	execution := rawExecution(t)
	corpus, dataset := fixtureData()
	dataset.Queries = append(dataset.Queries, ragoperators.Query{ID: "q2", Text: "fusion"})
	observer := &collector{}
	result, err := New(nil).Execute(context.Background(), execution, corpus, dataset, observer, Options{})
	if err != nil {
		t.Fatal(err)
	}
	unitArtifacts := 0
	for _, artifact := range observer.artifacts {
		if artifact.Role == "unit-set" {
			unitArtifacts++
		}
	}
	if unitArtifacts != 1 {
		t.Fatalf("unit artifacts=%d", unitArtifacts)
	}
	if len(result.Traces) != 2 || len(result.Traces[1].Operators) >= len(result.Traces[0].Operators) {
		t.Fatalf("traces=%#v", result.Traces)
	}
}

func TestEngineRecordsGeneratedFailureAsPartialEvidence(t *testing.T) {
	execution := executionWithRepresentation(t, ragmodel.StructuredSummary("summary", ragmodel.StructuredSummaryConfig{Generator: ragmodel.StructuredGenerationConfig{Model: "m", Prompt: "p", OutputSchema: "summary/v1"}}), "summary")
	modelDigest := "sha256:" + strings.Repeat("a", 64)
	promptDigest := "sha256:" + strings.Repeat("b", 64)
	resolver := ragoperators.StaticManifestResolver{Models: map[string]ragcontract.ModelManifest{"m": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: modelDigest}, ModelID: "m"}}, Prompts: map[string]ragcontract.PromptManifest{"p": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: promptDigest}, PromptID: "p", OutputSchema: "summary/v1"}}}
	corpus, dataset := fixtureData()
	observer := &collector{}
	result, err := New(nil).Execute(context.Background(), execution, corpus, dataset, observer, Options{Manifests: resolver, Schemas: acceptingSchema{}, Generator: failingGenerator{}, Cache: ragoperators.NewMemoryCache()})
	if err == nil || !strings.Contains(err.Error(), "RAG_GENERATION_FAILED") {
		t.Fatalf("err=%v", err)
	}
	if result == nil || len(result.Failures) != 1 || len(observer.events) == 0 || len(observer.traces) == 0 {
		t.Fatalf("partial result=%#v observer=%#v", result, observer)
	}
}

func TestEngineRejectsUnknownOperator(t *testing.T) {
	execution := rawExecution(t)
	execution.Pipeline.Nodes[0].Operator = ragcontract.OperatorRef{Kind: "unknown.operator", Version: "v1"}
	corpus, dataset := fixtureData()
	_, err := New(nil).Execute(context.Background(), execution, corpus, dataset, &collector{}, Options{})
	if err == nil {
		t.Fatal("unknown operator accepted")
	}
}
func rawExecution(t testing.TB) ragcontract.PipelineExecution {
	return executionWithRepresentation(t, ragmodel.RawRepresentation("raw"), "raw")
}
func executionWithRepresentation(t testing.TB, representation *ragmodel.Descriptor, representationName string) ragcontract.PipelineExecution {
	t.Helper()
	pipeline := ragmodel.NewPipeline("pipeline", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 200})).Represent(representation).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true}))
	})
	query := ragmodel.NewQueryPlan("query", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25(representationName+".lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: representationName, TopK: 10})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	product := ragmodel.NewProduct("raw", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Citations("source") })
	})
	size := int64(1)
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "json", Digest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", SizeBytes: &size, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: plan.Pipeline, Bindings: plan.Bindings, Dataset: ragcontract.DatasetBinding{Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", Version: "v1", ValueKind: "number", Unit: "ratio", Required: true, Config: json.RawMessage(`{}`)}}, VariantID: representationName, Factors: []ragcontract.FactorSelection{}}
	execution.CellID, _ = ragcontract.Digest(execution)
	return execution
}
func fixtureData() (ragoperators.Corpus, ragoperators.EvaluationDataset) {
	corpus := ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "s1", SessionID: "session", Ordinal: 1, Role: "user", Text: "weighted reciprocal rank fusion decision"}, {ID: "s2", SessionID: "session", Ordinal: 2, Role: "assistant", Text: "unrelated"}}}
	unitID := expectedFirstUnitID(corpus.Records[0])
	dataset := ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{{ID: "q1", Text: "reciprocal rank fusion", RelevantIDs: []string{unitID}}}}
	return corpus, dataset
}
func BenchmarkEngineRawBM25(b *testing.B) {
	execution := rawExecution(b)
	corpus, dataset := fixtureData()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := New(nil).Execute(context.Background(), execution, corpus, dataset, NopObserver{}, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func expectedFirstUnitID(record ragoperators.SourceRecord) string {
	digest, _ := ragcontract.Digest(record.Text)
	id, _ := ragcontract.Digest(struct {
		Kind   string
		IDs    []string
		Digest string
	}{"units.identity", []string{record.ID}, digest})
	return "unit:" + id[7:23]
}
