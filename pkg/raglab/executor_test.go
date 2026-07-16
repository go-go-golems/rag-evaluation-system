package raglab

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
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

type reverseReranker struct{ requests []RerankRequest }

func (r *reverseReranker) Identity() RerankerIdentity {
	return RerankerIdentity{Kind: "test", Model: "qllama/bge-reranker-v2-m3:q4_k_m"}
}

func (r *reverseReranker) Rerank(_ context.Context, request RerankRequest) ([]RerankResult, error) {
	r.requests = append(r.requests, request)
	results := make([]RerankResult, len(request.Candidates))
	for i := range request.Candidates {
		candidate := request.Candidates[len(request.Candidates)-1-i]
		results[i] = RerankResult{CandidateID: candidate.ID, Index: len(request.Candidates) - 1 - i, Score: float64(len(request.Candidates) - i), Rank: i + 1}
	}
	return results, nil
}

type fakeRunRecorder struct {
	events    []string
	traces    []experimentrun.QueryTraceInput
	completed *experimentrun.SummaryInput
}

func (f *fakeRunRecorder) AppendEvent(_ context.Context, _ string, event string, _ json.RawMessage) (*experimentrun.Event, error) {
	f.events = append(f.events, event)
	return &experimentrun.Event{Sequence: len(f.events), Type: event}, nil
}
func (f *fakeRunRecorder) RecordQueryTrace(_ context.Context, _ string, trace experimentrun.QueryTraceInput) error {
	f.traces = append(f.traces, trace)
	return nil
}
func (f *fakeRunRecorder) CompleteRun(_ context.Context, _ string, summary experimentrun.SummaryInput) (*experimentrun.Summary, error) {
	f.completed = &summary
	return &experimentrun.Summary{SummaryInput: summary}, nil
}

func TestExecutorRunsChannelsFusesCitationsAndRecordsTerminalSummary(t *testing.T) {
	specification, err := NewExperiment("hybrid").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).
		BM25(BM25Index("bm25")).Embeddings(EmbeddingSet("embeddings")).Evaluation(EvaluationDataset("evaluation")).
		Retrieval(func(builder *RetrievalBuilder) {
			builder.Channel("lexical", func(channel *ChannelBuilder) { channel.BM25().TopK(10) })
			builder.Channel("semantic", func(channel *ChannelBuilder) { channel.Vector().TopK(10) })
			builder.FuseRRF(60).Weight("semantic", 2).Collapse(CollapseDocument).Results(2)
		}).
		Metrics(func(metrics *MetricsBuilder) {
			metrics.RelevanceAt(RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).RecallAt(2).MRR()
		}).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	recorder := &fakeRunRecorder{}
	embedder := &fakeEmbedder{}
	result, err := NewExecutor(fakeRetrievalBackend{}, recorder).Execute(context.Background(), "run-1", specification, []EvaluationCard{{ID: "q-1", Query: "test", RelevantDocumentRevisionIDs: []string{"b"}}}, ExecutionOptions{Embedder: embedder})
	if err != nil {
		t.Fatal(err)
	}
	if result.QueryCount != 1 || embedder.calls != 1 || recorder.completed == nil || recorder.completed.Status != "succeeded" {
		t.Fatalf("execution result = %#v recorder=%#v", result, recorder)
	}
	if len(recorder.events) != 2 || recorder.events[0] != "execution_started" || recorder.events[1] != "execution_completed" {
		t.Fatalf("events = %#v", recorder.events)
	}
	if len(recorder.traces) != 1 {
		t.Fatalf("traces = %#v", recorder.traces)
	}
	var trace executionTrace
	if err := json.Unmarshal(recorder.traces[0].Trace, &trace); err != nil {
		t.Fatal(err)
	}
	if len(trace.Results) != 2 || trace.Results[0].DocumentRevisionID != "b" || trace.Results[0].URL == "" || trace.Fusion[0].Components["semantic"].Contribution <= trace.Fusion[0].Components["lexical"].Contribution {
		t.Fatalf("trace = %#v", trace)
	}
}

func TestExecutorRejectsVectorPlanWithoutExplicitEmbedder(t *testing.T) {
	specification, err := NewExperiment("vector").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).Embeddings(EmbeddingSet("embeddings")).Evaluation(EvaluationDataset("evaluation")).
		Retrieval(func(builder *RetrievalBuilder) {
			builder.Channel("semantic", func(channel *ChannelBuilder) { channel.Vector().TopK(10) }).Results(10)
		}).
		Metrics(func(metrics *MetricsBuilder) {
			metrics.RelevanceAt(RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).RecallAt(10)
		}).
		Build()
	if err != nil {
		t.Fatal(err)
	}
	recorder := &fakeRunRecorder{}
	_, err = NewExecutor(fakeRetrievalBackend{}, recorder).Execute(context.Background(), "run-1", specification, []EvaluationCard{{ID: "q", Query: "query"}}, ExecutionOptions{})
	if err == nil || err.Error() != "RAG_EMBEDDER_REQUIRED: vector retrieval needs an explicit query embedder" || len(recorder.events) != 0 {
		t.Fatalf("err=%v events=%#v", err, recorder.events)
	}
}

func TestExecutorReranksBoundedCandidatesAndPersistsBothOrders(t *testing.T) {
	specification, err := NewExperiment("reranked").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Evaluation(EvaluationDataset("evaluation")).
		Retrieval(func(builder *RetrievalBuilder) {
			builder.Channel("lexical", func(channel *ChannelBuilder) { channel.BM25().TopK(10) })
			builder.RerankCrossEncoder("qllama/bge-reranker-v2-m3:q4_k_m", 2, 2).Collapse(CollapseNone).Results(2)
		}).Metrics(func(metrics *MetricsBuilder) {
		metrics.RelevanceAt(RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).MRR()
	}).Build()
	if err != nil {
		t.Fatal(err)
	}
	backend := fakeRetrievalBackendWithText{}
	recorder := &fakeRunRecorder{}
	reranker := &reverseReranker{}
	_, err = NewExecutor(backend, recorder).Execute(context.Background(), "run-rerank", specification, []EvaluationCard{{ID: "q-1", Query: "test"}}, ExecutionOptions{Reranker: reranker})
	if err != nil {
		t.Fatal(err)
	}
	if len(reranker.requests) != 1 || reranker.requests[0].TopN != 2 || reranker.requests[0].Candidates[0].ID != "a-1" {
		t.Fatalf("reranker requests = %#v", reranker.requests)
	}
	var trace executionTrace
	if err := json.Unmarshal(recorder.traces[0].Trace, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Reranking == nil || len(trace.Reranking.Candidates) != 2 || len(trace.Reranking.Results) != 2 {
		t.Fatalf("reranking trace = %#v", trace.Reranking)
	}
	if trace.Results[0].ChunkID != "b-1" || trace.Results[0].Rank != 1 || trace.Results[0].Score != 2 {
		t.Fatalf("reranked results = %#v", trace.Results)
	}
}

type fakeRetrievalBackendWithText struct{ fakeRetrievalBackend }

func (fakeRetrievalBackendWithText) BM25(ctx context.Context, artifactID, query string, limit int) ([]immutableretrieval.ChunkHit, error) {
	hits, err := fakeRetrievalBackend{}.BM25(ctx, artifactID, query, limit)
	for i := range hits {
		hits[i].Text = "candidate text"
	}
	return hits, err
}
