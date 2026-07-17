package raglab

import (
	"context"
	"errors"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

var errUnexpected = errors.New("unexpected fake retrieval request")

type fakeRetrievalBackend struct{}

func (fakeRetrievalBackend) BM25(_ context.Context, artifactID, _ string, limit int) ([]immutableretrieval.ChunkHit, error) {
	if artifactID != "bm25" || limit != 10 {
		return nil, errUnexpected
	}
	return []immutableretrieval.ChunkHit{
		{Rank: 1, ChunkID: "a-1", DocumentRevisionID: "a", Title: "A", URL: "https://example.test/a", Score: 5},
		{Rank: 2, ChunkID: "b-1", DocumentRevisionID: "b", Title: "B", URL: "https://example.test/b", Score: 4},
	}, nil
}

func (fakeRetrievalBackend) Vector(_ context.Context, artifactID string, vector []float32, limit int) ([]immutableretrieval.ChunkHit, error) {
	if artifactID != "embeddings" || len(vector) != 2 || limit != 10 {
		return nil, errUnexpected
	}
	return []immutableretrieval.ChunkHit{
		{Rank: 1, ChunkID: "b-1", DocumentRevisionID: "b", Title: "B", URL: "https://example.test/b", Score: .9},
		{Rank: 2, ChunkID: "a-1", DocumentRevisionID: "a", Title: "A", URL: "https://example.test/a", Score: .8},
	}, nil
}

type fakeEmbedder struct{ calls int }

func (f *fakeEmbedder) GenerateEmbedding(_ context.Context, _ string) ([]float32, error) {
	f.calls++
	return []float32{1, 2}, nil
}

type collectingObserver struct {
	events    []DomainEvent
	traces    []ragcontract.QueryTrace
	metrics   []DomainMetric
	artifacts []DomainArtifact
}

func (o *collectingObserver) Event(_ context.Context, value DomainEvent) error {
	o.events = append(o.events, value)
	return nil
}
func (o *collectingObserver) QueryTrace(_ context.Context, value ragcontract.QueryTrace) error {
	o.traces = append(o.traces, value)
	return nil
}
func (o *collectingObserver) Metric(_ context.Context, value DomainMetric) error {
	o.metrics = append(o.metrics, value)
	return nil
}
func (o *collectingObserver) Artifact(_ context.Context, value DomainArtifact) error {
	o.artifacts = append(o.artifacts, value)
	return nil
}

func TestObservationExecutorEmitsPublicEvidenceWithoutRunLifecycle(t *testing.T) {
	grade, _ := Grade("2_SUBSTANTIAL")
	specification, err := NewExperiment("observer-hybrid").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Embeddings(EmbeddingSet("embeddings")).Evaluation(EvaluationDataset("evaluation")).
		Representations(func(builder *RepresentationBuilder) { builder.RawChunks("raw") }).
		Retrieval(func(builder *RetrievalBuilder) {
			builder.Channel("lexical", func(channel *ChannelBuilder) { channel.BM25().Representation("raw").TopK(10) })
			builder.Channel("semantic", func(channel *ChannelBuilder) { channel.Vector().Representation("raw").TopK(10) })
			builder.FuseRRF(60).Weight("semantic", 2).Collapse(CollapseDocument).Results(2)
		}).Metrics(func(metrics *MetricsBuilder) {
		metrics.RelevanceAt(grade).RecallAt(2).PrecisionAt(1, 2).NDCGAt(2).MRR()
	}).Build()
	if err != nil {
		t.Fatal(err)
	}
	observer := &collectingObserver{}
	err = NewObservationExecutor(fakeRetrievalBackend{}).Execute(context.Background(), ObservationExecutionRequest{
		Specification: specification, DatasetSplit: "development",
		Cards:   []EvaluationCard{{ID: "q-1", Query: "test", RelevantDocumentRevisionIDs: []string{"b"}}},
		Options: ExecutionOptions{Embedder: &fakeEmbedder{}},
	}, observer)
	if err != nil {
		t.Fatal(err)
	}
	if len(observer.events) != 4 || observer.events[0].Type != "rag.execution.started" || observer.events[3].Type != "rag.execution.completed" {
		t.Fatalf("events = %#v", observer.events)
	}
	if len(observer.traces) != 1 || observer.traces[0].SchemaVersion != ragcontract.TraceSchemaVersion {
		t.Fatalf("traces = %#v", observer.traces)
	}
	trace := observer.traces[0]
	if len(trace.Channels) != 2 || trace.Channels[0].Name != "lexical" || trace.Channels[0].Backend != "bm25" || trace.Channels[1].Name != "semantic" {
		t.Fatalf("channel order/metadata = %#v", trace.Channels)
	}
	if trace.Fusion == nil || trace.Fusion.RankConstant != 60 || trace.Results[0].DocumentRevisionID != "b" {
		t.Fatalf("trace = %#v", trace)
	}
	metrics := map[string]float64{}
	for _, metric := range observer.metrics {
		metrics[metric.Name] = metric.Value
	}
	for _, name := range []string{"rag.query_count", "rag.answerable_query_count", "rag.mean_reciprocal_rank", "rag.mean_relevant_document_recall_at_2", "rag.mean_precision_at_1", "rag.mean_precision_at_2", "rag.mean_ndcg_at_2", "rag.wall_clock_duration_ms"} {
		if _, ok := metrics[name]; !ok {
			t.Fatalf("missing metric %s in %#v", name, observer.metrics)
		}
	}
	if metrics["rag.mean_reciprocal_rank"] != 1 || metrics["rag.mean_precision_at_1"] != 1 || metrics["rag.mean_relevant_document_recall_at_2"] != 1 {
		t.Fatalf("unexpected relevance metrics: %#v", metrics)
	}
}

type failSecondBackend struct{ calls int }

func (b *failSecondBackend) BM25(_ context.Context, _, _ string, _ int) ([]immutableretrieval.ChunkHit, error) {
	b.calls++
	if b.calls == 2 {
		return nil, errors.New("second query failed")
	}
	return []immutableretrieval.ChunkHit{{Rank: 1, ChunkID: "a-1", DocumentRevisionID: "a", Score: 1}}, nil
}
func (*failSecondBackend) Vector(context.Context, string, []float32, int) ([]immutableretrieval.ChunkHit, error) {
	return nil, errors.New("unexpected vector call")
}

func TestObservationExecutorPreservesPartialEvidenceOnRunnerError(t *testing.T) {
	specification := ExperimentSpecification{
		Inputs:    InputSpec{BM25Index: pointerArtifact(BM25Index("bm25")), Representations: []RepresentationSpec{{Name: "raw", Kind: RawChunksRepresentation}}},
		Retrieval: RetrievalPlan{Channels: []ChannelSpec{{Name: "lexical", Backend: BM25Backend, TopK: 1}}, Collapse: CollapseNone, Results: 1},
		Metrics:   MetricsPlan{RelevanceAt: &RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}, MRR: true},
	}
	observer := &collectingObserver{}
	err := NewObservationExecutor(&failSecondBackend{}).Execute(context.Background(), ObservationExecutionRequest{
		Specification: specification, DatasetSplit: "development",
		Cards: []EvaluationCard{{ID: "q-1", Query: "first"}, {ID: "q-2", Query: "second"}},
	}, observer)
	if err == nil || err.Error() != "retrieve channel \"lexical\": second query failed" {
		t.Fatalf("expected second-query failure, got %v", err)
	}
	if len(observer.traces) != 1 || observer.traces[0].QueryCardID != "q-1" {
		t.Fatalf("partial traces = %#v", observer.traces)
	}
	if len(observer.metrics) != 0 || observer.events[len(observer.events)-1].Type == "rag.execution.completed" {
		t.Fatalf("failure must not emit aggregate metrics/completion: metrics=%#v events=%#v", observer.metrics, observer.events)
	}
}

func pointerArtifact(value ArtifactRef) *ArtifactRef { return &value }

func TestObservationExecutorHonorsCancellationBeforeObservation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	observer := &collectingObserver{}
	specification := ExperimentSpecification{
		Inputs:    InputSpec{Representations: []RepresentationSpec{{Name: "raw", Kind: RawChunksRepresentation}}},
		Retrieval: RetrievalPlan{Channels: []ChannelSpec{{Name: "lexical", Backend: BM25Backend, TopK: 1}}, Collapse: CollapseNone, Results: 1},
	}
	err := NewObservationExecutor(fakeRetrievalBackend{}).Execute(ctx, ObservationExecutionRequest{
		Specification: specification, DatasetSplit: "development", Cards: []EvaluationCard{{ID: "q", Query: "test"}},
	}, observer)
	if !errors.Is(err, context.Canceled) || len(observer.events) != 0 {
		t.Fatalf("err=%v events=%#v", err, observer.events)
	}
}

func TestObservationExecutorRejectsFilterBeforeObservation(t *testing.T) {
	observer := &collectingObserver{}
	specification := ExperimentSpecification{
		Inputs:    InputSpec{Representations: []RepresentationSpec{{Name: "raw", Kind: RawChunksRepresentation}}},
		Retrieval: RetrievalPlan{Channels: []ChannelSpec{{Name: "lexical", Backend: BM25Backend, TopK: 1, Filter: FilterSpec{SourceIDs: []string{"wp:1"}}}}, Collapse: CollapseNone, Results: 1},
	}
	err := NewObservationExecutor(fakeRetrievalBackend{}).Execute(context.Background(), ObservationExecutionRequest{
		Specification: specification, DatasetSplit: "development", Cards: []EvaluationCard{{ID: "q", Query: "test"}},
	}, observer)
	if err == nil || len(observer.events) != 0 {
		t.Fatalf("expected pre-observation capability failure, err=%v events=%#v", err, observer.events)
	}
}
