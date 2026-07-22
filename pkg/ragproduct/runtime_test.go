package ragproduct

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

func TestLoadRejectsUnknownAndInvalidPolicies(t *testing.T) {
	plan, artifact := fixturePlan(t, "authoritative", "fail", 2)
	data, _ := json.Marshal(plan)
	var value map[string]any
	_ = json.Unmarshal(data, &value)
	value["unknown"] = true
	bad, _ := json.Marshal(value)
	if _, err := LoadBytes(bad); err == nil {
		t.Fatal("unknown product field accepted")
	}
	plan.Runtime.TracePolicy = "sample-silently"
	data, _ = json.Marshal(plan)
	if _, err := LoadBytes(data); err == nil {
		t.Fatal("unknown trace policy accepted")
	}
	_ = artifact
}

func TestRuntimeRejectsMismatchedResolvedModelBinding(t *testing.T) {
	plan, artifact := fixturePlan(t, "authoritative", "fail", 1)
	fixtures := ragoperators.NewFixtureProviders()
	manifest, _ := fixtures.Resolver.Model(ragoperators.FixtureEmbeddingModel)
	plan.Models = []ragcontract.ModelBinding{{Reference: ragoperators.FixtureEmbeddingModel, Manifest: ragcontract.ModelManifestSchema, Digest: "sha256:" + string(bytes.Repeat([]byte("f"), 64))}}
	if _, err := New(context.Background(), plan, Bindings{Corpus: artifact, Manifests: fixtures.Resolver}); err == nil {
		t.Fatal("mismatched model binding accepted")
	}
	plan.Models[0].Digest = manifest.Digest
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Manifests: fixtures.Resolver})
	if err != nil {
		t.Fatal(err)
	}
	_ = runtime.Close()
}

func TestRuntimePreparedConcurrentCitationsAndRequestValidation(t *testing.T) {
	plan, artifact := fixturePlan(t, "authoritative", "fail", 4)
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close() }()
	if _, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"unknown": "x"}}); err == nil {
		t.Fatal("unknown request field accepted")
	}
	var wg sync.WaitGroup
	errorsFound := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			response, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"query": "reciprocal rank fusion"}})
			if err != nil {
				errorsFound <- err
				return
			}
			if len(response.Results) == 0 || len(response.Citations) == 0 || response.Trace == nil || response.TraceID == "" {
				errorsFound <- errors.New("incomplete response")
			}
		}()
	}
	wg.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Error(err)
	}
}

func TestArtifactBackedTraceAndCancellation(t *testing.T) {
	plan, artifact := fixturePlan(t, "artifact-backed", "fail", 1)
	sink := &memorySink{}
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache(), Traces: sink})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close() }()
	response, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"query": "fusion SENSITIVE_QUERY_CANARY"}})
	if err != nil {
		t.Fatal(err)
	}
	if response.Trace != nil || sink.values[response.TraceID] == nil {
		t.Fatalf("trace policy not applied: %#v", response)
	}
	if bytes.Contains(sink.values[response.TraceID], []byte("SENSITIVE_QUERY_CANARY")) {
		t.Fatal("query text leaked into trace artifact")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := runtime.Execute(ctx, Request{Values: map[string]any{"query": "fusion"}}); err == nil {
		t.Fatal("cancelled request succeeded")
	}
}

func TestDeclaredProviderFailurePolicies(t *testing.T) {
	fixtures := ragoperators.NewFixtureProviders()
	for _, policy := range []string{"fail", "abstain", "retrieval-only"} {
		t.Run(policy, func(t *testing.T) {
			plan, artifact := fixtureGeneratedPlan(t, policy)
			runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Manifests: fixtures.Resolver, Generator: failingGenerator{}, Cache: ragoperators.NewMemoryCache()})
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = runtime.Close() }()
			response, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"query": "fusion"}})
			switch policy {
			case "fail":
				if err == nil {
					t.Fatal("fail policy returned success")
				}
			case "abstain":
				if err != nil || !response.Abstained || len(response.Results) != 0 {
					t.Fatalf("response=%#v err=%v", response, err)
				}
			case "retrieval-only":
				if err != nil || len(response.Results) == 0 || response.Failure == nil {
					t.Fatalf("response=%#v err=%v", response, err)
				}
			}
			encoded, _ := json.Marshal(response)
			if bytes.Contains(encoded, []byte("SECRET_CANARY")) {
				t.Fatal("provider secret leaked into response")
			}
		})
	}
}

func TestRuntimeEnforcesDeclaredTimeout(t *testing.T) {
	plan, artifact := fixtureGeneratedPlan(t, "fail")
	plan.Runtime.TimeoutMilliseconds = 2
	fixtures := ragoperators.NewFixtureProviders()
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Manifests: fixtures.Resolver, Generator: blockingGenerator{}, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close() }()
	if _, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"query": "fusion"}}); err == nil || !errors.Is(errors.Unwrap(err), context.DeadlineExceeded) && !bytes.Contains([]byte(err.Error()), []byte("deadline exceeded")) {
		t.Fatalf("timeout err=%v", err)
	}
}

type blockingGenerator struct{}

func (blockingGenerator) Generate(ctx context.Context, _ ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	<-ctx.Done()
	return ragoperators.GenerationResult{}, ctx.Err()
}

type failingGenerator struct{}

func (failingGenerator) Generate(context.Context, ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	return ragoperators.GenerationResult{}, errors.New("provider unavailable SECRET_CANARY")
}

func TestQualificationResolvesExactModelAndPromptBindings(t *testing.T) {
	plan, _ := fixtureGeneratedPlan(t, "fail")
	fixtures := ragoperators.NewFixtureProviders()
	dataset := ragcontract.DatasetBinding{ManifestDigest: "sha256:" + string(bytes.Repeat([]byte("b"), 64)), Split: "qualification", Status: "candidate", RelevanceTarget: "unit"}
	if _, err := Qualify(plan, dataset, nil); err == nil {
		t.Fatal("unresolved provider references qualified")
	}
	qualification, err := QualifyResolved(plan, dataset, nil, fixtures.Resolver)
	if err != nil {
		t.Fatal(err)
	}
	if len(qualification.Models) != 1 || len(qualification.Prompts) != 1 || qualification.Models[0].Digest == "" || qualification.Prompts[0].Digest == "" {
		t.Fatalf("qualification=%#v", qualification)
	}
}

func TestQualificationPreservesExactProductPipelineAndBindings(t *testing.T) {
	plan, artifact := fixturePlan(t, "authoritative", "fail", 1)
	dataset := ragcontract.DatasetBinding{ManifestDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", Split: "qualification", Status: "candidate", RelevanceTarget: "unit"}
	qualification, err := Qualify(plan, dataset, []ragcontract.Measure{{Name: "rag.mrr", Version: "v1", ValueKind: "number", Required: true, Config: json.RawMessage(`{}`)}})
	if err != nil {
		t.Fatal(err)
	}
	productPipeline, _ := ragcontract.CanonicalJSON(plan.Pipeline)
	studyPipeline, _ := ragcontract.CanonicalJSON(qualification.Study.Variants[0].Pipeline)
	if !bytes.Equal(productPipeline, studyPipeline) || qualification.ProductID == "" {
		t.Fatal("qualification changed pipeline")
	}
	cells, err := ragcompiler.ExpandStudy(qualification.Study, nil)
	if err != nil {
		t.Fatal(err)
	}
	corpus := artifact.Corpus
	unitID := expectedUnitID(corpus.Records[0])
	studyResult, err := ragengine.New(nil).Execute(context.Background(), cells[0].Execution, corpus, ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q", Text: "reciprocal rank fusion", RelevantIDs: []string{unitID}}}}, nil, ragengine.Options{})
	if err != nil {
		t.Fatal(err)
	}
	productRuntime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = productRuntime.Close() }()
	productResult, err := productRuntime.Execute(context.Background(), Request{ID: "q", Values: map[string]any{"query": "reciprocal rank fusion"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(studyResult.Traces) != 1 || productResult.Trace == nil {
		t.Fatal("missing traces")
	}
	a, _ := ragcontract.CanonicalJSON(studyResult.Traces[0].Results)
	b, _ := ragcontract.CanonicalJSON(productResult.Trace.Results)
	if !bytes.Equal(a, b) {
		t.Fatalf("retrieval traces differ\nstudy=%s\nproduct=%s", a, b)
	}
}

func TestProductLatencyProfile(t *testing.T) {
	plan, artifact := fixturePlan(t, "authoritative", "fail", 1)
	started := time.Now()
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		t.Fatal(err)
	}
	startup := time.Since(started)
	defer func() { _ = runtime.Close() }()
	durations := make([]time.Duration, 200)
	traceBytes := 0
	for index := range durations {
		started = time.Now()
		response, err := runtime.Execute(context.Background(), Request{Values: map[string]any{"query": "reciprocal rank fusion"}})
		if err != nil {
			t.Fatal(err)
		}
		durations[index] = time.Since(started)
		encoded, _ := json.Marshal(response.Trace)
		traceBytes += len(encoded)
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	percentile := func(value float64) time.Duration { return durations[int(float64(len(durations)-1)*value)] }
	t.Logf("prepared_startup=%s p50=%s p95=%s p99=%s mean_trace_bytes=%d", startup, percentile(.50), percentile(.95), percentile(.99), traceBytes/len(durations))
}

func BenchmarkPreparedProductQuery(b *testing.B) {
	plan, artifact := fixturePlan(b, "none", "fail", 1)
	runtime, err := New(context.Background(), plan, Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = runtime.Close() }()
	request := Request{Values: map[string]any{"query": "reciprocal rank fusion"}}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := runtime.Execute(context.Background(), request); err != nil {
			b.Fatal(err)
		}
	}
}

type memorySink struct{ values map[string][]byte }

func (s *memorySink) Put(_ context.Context, id, _ string, value []byte) error {
	if s.values == nil {
		s.values = map[string][]byte{}
	}
	s.values[id] = append([]byte(nil), value...)
	return nil
}

func fixturePlan(t testing.TB, tracePolicy, failurePolicy string, concurrent int) (ragcontract.ProductPlan, ragoperators.CorpusArtifact) {
	t.Helper()
	corpus := ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "s1", SessionID: "session", Ordinal: 1, Role: "user", Text: "weighted reciprocal rank fusion decision"}, {ID: "s2", SessionID: "session", Ordinal: 2, Role: "assistant", Text: "unrelated"}}}
	artifact := ragoperators.NewCorpusArtifact(corpus, "fixture")
	pipeline := ragmodel.NewPipeline("pipeline", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 200})).Represent(ragmodel.RawRepresentation("raw")).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true}))
	})
	query := ragmodel.NewQueryPlan("query", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	product := ragmodel.NewProduct("fixture-product", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).RequestContract(func(r *ragmodel.RequestBuilder) { r.Field("query", "string", true, 256) }).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Answer("text").Citations("required").TraceID(true) }).RuntimePolicy(func(r *ragmodel.RuntimeBuilder) {
			r.Concurrent(concurrent).ProviderFailure(failurePolicy).Trace(tracePolicy).Timeout(5000)
		})
	})
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {SlotID: "corpus", Role: "corpus", Kind: "json", Digest: artifact.Manifest.Digest, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	return plan, artifact
}

func fixtureGeneratedPlan(t testing.TB, failurePolicy string) (ragcontract.ProductPlan, ragoperators.CorpusArtifact) {
	plan, artifact := fixturePlan(t, "metadata-only", failurePolicy, 1)
	// Append answer generation through the authoring model so it remains ordinary canonical IR.
	pipeline := ragmodel.NewPipeline("pipeline", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 200})).Represent(ragmodel.RawRepresentation("raw")).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true}))
	})
	query := ragmodel.NewQueryPlan("query", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	product := ragmodel.NewProduct("generated", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).Generate(ragmodel.Answer(ragmodel.AnswerConfig{Model: ragoperators.FixtureSummaryModel, Prompt: ragoperators.FixtureSummaryPrompt, Citations: "required", ContextBudgetTokens: 512})).RequestContract(func(r *ragmodel.RequestBuilder) { r.Field("query", "string", true, 256) }).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Answer("text").Citations("required").TraceID(true) }).RuntimePolicy(func(r *ragmodel.RuntimeBuilder) {
			r.Concurrent(1).ProviderFailure(failurePolicy).Trace("metadata-only")
		})
	})
	compiled, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {SlotID: "corpus", Role: "corpus", Kind: "json", Digest: artifact.Manifest.Digest, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	_ = plan
	return compiled, artifact
}

func expectedUnitID(record ragoperators.SourceRecord) string {
	digest, _ := ragcontract.Digest(record.Text)
	id, _ := ragcontract.Digest(struct {
		Kind   string
		IDs    []string
		Digest string
	}{"units.identity", []string{record.ID}, digest})
	return "unit:" + id[7:23]
}
