package raglab

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/pkg/errors"
)

// QueryEmbedder is the narrow capability needed for vector channels. It is
// deliberately injected: selecting an immutable embedding set does not grant
// a plan permission to create query embeddings with an implicit provider.
type QueryEmbedder interface {
	GenerateEmbedding(context.Context, string) ([]float32, error)
}

// ChannelRetriever isolates execution planning from SQLite and permits a
// deterministic fake backend in tests.
type ChannelRetriever interface {
	BM25(context.Context, string, string, int) ([]immutableretrieval.ChunkHit, error)
	Vector(context.Context, string, []float32, int) ([]immutableretrieval.ChunkHit, error)
}

// RunRecorder is the append-only subset of experimentrun.Service used by a
// single execution. The executor never writes SQL directly.
type RunRecorder interface {
	AppendEvent(context.Context, string, string, json.RawMessage) (*experimentrun.Event, error)
	RecordQueryTrace(context.Context, string, experimentrun.QueryTraceInput) error
	CompleteRun(context.Context, string, experimentrun.SummaryInput) (*experimentrun.Summary, error)
}

var _ RunRecorder = (*experimentrun.Service)(nil)

type EvaluationCard struct {
	ID                          string   `json:"id"`
	Query                       string   `json:"query"`
	RelevantDocumentRevisionIDs []string `json:"relevantDocumentRevisionIds,omitempty"`
}

type ExecutionOptions struct {
	Embedder QueryEmbedder
}

type ExecutionResult struct {
	RunID       string         `json:"runId"`
	QueryCount  int            `json:"queryCount"`
	Metrics     map[string]any `json:"metrics"`
	Timing      map[string]any `json:"timing"`
	CompletedAt time.Time      `json:"completedAt"`
}

// Executor executes raw immutable chunk retrieval plans and stores every
// observation through the append-only experiment-run service. Summary/question
// representations and parent-chunk collapse are rejected until their
// materialisation and parent mapping schemas exist.
type Executor struct {
	backend  ChannelRetriever
	recorder RunRecorder
}

func NewExecutor(backend ChannelRetriever, recorder RunRecorder) *Executor {
	return &Executor{backend: backend, recorder: recorder}
}

func (e *Executor) Execute(ctx context.Context, runID string, specification ExperimentSpecification, cards []EvaluationCard, options ExecutionOptions) (_ *ExecutionResult, retErr error) {
	if e == nil || e.backend == nil || e.recorder == nil {
		return nil, errors.New("RAG_EXECUTOR_REQUIRED: retrieval backend and run recorder are required")
	}
	if runID == "" || len(cards) == 0 {
		return nil, errors.New("RAG_EXECUTION_INPUT_REQUIRED: run ID and evaluation cards are required")
	}
	if err := executable(specification, options); err != nil {
		return nil, err
	}
	if _, err := e.recorder.AppendEvent(ctx, runID, "execution_started", json.RawMessage(`{"executor":"raglab/raw-v1"}`)); err != nil {
		return nil, errors.Wrap(err, "append execution start event")
	}
	started := time.Now()
	var totalMilliseconds int64
	var reciprocalRank, recallAtResults float64
	answerable := 0
	for _, card := range cards {
		trace, err := e.executeCard(ctx, specification, card, options)
		if err != nil {
			return e.fail(ctx, runID, err)
		}
		totalMilliseconds += trace.Timing.TotalMilliseconds
		if len(card.RelevantDocumentRevisionIDs) > 0 {
			answerable++
			rank, recall := relevance(trace.Results, card.RelevantDocumentRevisionIDs, specification.Retrieval.Results)
			if rank > 0 {
				reciprocalRank += 1 / float64(rank)
			}
			recallAtResults += recall
		}
		traceJSON, err := json.Marshal(trace)
		if err != nil {
			return e.fail(ctx, runID, errors.Wrap(err, "encode query trace"))
		}
		metricsJSON, err := json.Marshal(map[string]any{"firstRelevantRank": trace.FirstRelevantRank, "recallAtResults": trace.RecallAtResults})
		if err != nil {
			return e.fail(ctx, runID, errors.Wrap(err, "encode query metrics"))
		}
		timingJSON, err := json.Marshal(trace.Timing)
		if err != nil {
			return e.fail(ctx, runID, errors.Wrap(err, "encode query timing"))
		}
		if err := e.recorder.RecordQueryTrace(ctx, runID, experimentrun.QueryTraceInput{QueryCardID: card.ID, Trace: traceJSON, Metrics: metricsJSON, Timing: timingJSON, Cost: json.RawMessage(`{}`), Storage: json.RawMessage(`{}`)}); err != nil {
			return e.fail(ctx, runID, errors.Wrap(err, "record query trace"))
		}
	}
	metrics := map[string]any{"queries": len(cards), "answerableQueries": answerable}
	if answerable > 0 {
		metrics["meanReciprocalRank"] = reciprocalRank / float64(answerable)
		metrics["meanRelevantRecallAtResults"] = recallAtResults / float64(answerable)
	}
	timing := map[string]any{"totalMilliseconds": totalMilliseconds, "wallClockMilliseconds": time.Since(started).Milliseconds()}
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return e.fail(ctx, runID, errors.Wrap(err, "encode run metrics"))
	}
	timingJSON, err := json.Marshal(timing)
	if err != nil {
		return e.fail(ctx, runID, errors.Wrap(err, "encode run timing"))
	}
	if _, err := e.recorder.AppendEvent(ctx, runID, "execution_completed", json.RawMessage(`{"status":"succeeded"}`)); err != nil {
		return e.fail(ctx, runID, errors.Wrap(err, "append execution completion event"))
	}
	if _, err := e.recorder.CompleteRun(ctx, runID, experimentrun.SummaryInput{Status: "succeeded", Metrics: metricsJSON, Cost: json.RawMessage(`{}`), Storage: timingJSON, Error: json.RawMessage(`{}`)}); err != nil {
		return nil, errors.Wrap(err, "complete experiment run")
	}
	return &ExecutionResult{RunID: runID, QueryCount: len(cards), Metrics: metrics, Timing: timing, CompletedAt: time.Now().UTC()}, nil
}

func (e *Executor) fail(ctx context.Context, runID string, cause error) (*ExecutionResult, error) {
	errorJSON, err := json.Marshal(map[string]string{"message": cause.Error()})
	if err == nil {
		_, _ = e.recorder.CompleteRun(ctx, runID, experimentrun.SummaryInput{Status: "failed", Metrics: json.RawMessage(`{}`), Cost: json.RawMessage(`{}`), Storage: json.RawMessage(`{}`), Error: errorJSON})
	}
	return nil, cause
}

func executable(specification ExperimentSpecification, options ExecutionOptions) error {
	for _, representation := range specification.Inputs.Representations {
		if representation.Kind != RawChunksRepresentation {
			return errors.New("RAG_EXECUTION_UNSUPPORTED: materialized representations need a representation executor")
		}
	}
	if specification.Retrieval.Collapse == CollapseParentChunk {
		return errors.New("RAG_EXECUTION_UNSUPPORTED: parentChunk collapse needs materialized parent mappings")
	}
	for _, channel := range specification.Retrieval.Channels {
		if channel.Backend == VectorBackend && options.Embedder == nil {
			return errors.New("RAG_EMBEDDER_REQUIRED: vector retrieval needs an explicit query embedder")
		}
	}
	return nil
}

type executionTrace struct {
	QueryID           string                                   `json:"queryId"`
	Query             string                                   `json:"query"`
	Channels          map[string][]immutableretrieval.ChunkHit `json:"channels"`
	Results           []immutableretrieval.ChunkHit            `json:"results"`
	Fusion            []immutableretrieval.FusedHit            `json:"fusion,omitempty"`
	Timing            executionTiming                          `json:"timing"`
	FirstRelevantRank int                                      `json:"firstRelevantRank,omitempty"`
	RecallAtResults   float64                                  `json:"recallAtResults,omitempty"`
}

type executionTiming struct {
	EmbeddingMilliseconds int64 `json:"embeddingMilliseconds"`
	RetrievalMilliseconds int64 `json:"retrievalMilliseconds"`
	FusionMilliseconds    int64 `json:"fusionMilliseconds"`
	TotalMilliseconds     int64 `json:"totalMilliseconds"`
}

func (e *Executor) executeCard(ctx context.Context, specification ExperimentSpecification, card EvaluationCard, options ExecutionOptions) (executionTrace, error) {
	if card.ID == "" || card.Query == "" {
		return executionTrace{}, errors.New("RAG_INVALID_EVALUATION_CARD: card ID and query are required")
	}
	started := time.Now()
	channels := map[string][]immutableretrieval.ChunkHit{}
	var queryVector []float32
	var embeddingMilliseconds int64
	for _, channel := range specification.Retrieval.Channels {
		channelStarted := time.Now()
		var hits []immutableretrieval.ChunkHit
		var err error
		switch channel.Backend {
		case BM25Backend:
			hits, err = e.backend.BM25(ctx, specification.Inputs.BM25Index.ID, card.Query, channel.TopK)
		case VectorBackend:
			if queryVector == nil {
				embeddingStarted := time.Now()
				queryVector, err = options.Embedder.GenerateEmbedding(ctx, card.Query)
				embeddingMilliseconds = time.Since(embeddingStarted).Milliseconds()
				if err != nil {
					return executionTrace{}, errors.Wrap(err, "embed evaluation query")
				}
			}
			hits, err = e.backend.Vector(ctx, specification.Inputs.EmbeddingSet.ID, queryVector, channel.TopK)
		default:
			err = errors.New("RAG_EXECUTION_UNSUPPORTED: retrieval backend is not executable")
		}
		if err != nil {
			return executionTrace{}, errors.Wrapf(err, "retrieve channel %q", channel.Name)
		}
		for index := range hits {
			hits[index].Channel = channel.Name
		}
		channels[channel.Name] = hits
		_ = channelStarted
	}
	retrievalMilliseconds := time.Since(started).Milliseconds() - embeddingMilliseconds
	fusionStarted := time.Now()
	trace := executionTrace{QueryID: card.ID, Query: card.Query, Channels: channels}
	if len(specification.Retrieval.Channels) == 1 {
		trace.Results = append([]immutableretrieval.ChunkHit(nil), channels[specification.Retrieval.Channels[0].Name]...)
		if specification.Retrieval.Collapse == CollapseDocument {
			trace.Results = immutableretrieval.CollapseDocuments(trace.Results)
		}
		if len(trace.Results) > specification.Retrieval.Results {
			trace.Results = trace.Results[:specification.Retrieval.Results]
		}
	} else {
		weights := map[string]float64(nil)
		if specification.Retrieval.Fusion != nil {
			weights = specification.Retrieval.Fusion.Weights
		}
		trace.Fusion = immutableretrieval.FuseWeightedRRF(channels, specification.Retrieval.Fusion.RankConstant, weights, specification.Retrieval.Results)
		trace.Results = make([]immutableretrieval.ChunkHit, len(trace.Fusion))
		for i := range trace.Fusion {
			trace.Results[i] = trace.Fusion[i].ChunkHit
		}
	}
	trace.Timing = executionTiming{EmbeddingMilliseconds: embeddingMilliseconds, RetrievalMilliseconds: retrievalMilliseconds, FusionMilliseconds: time.Since(fusionStarted).Milliseconds(), TotalMilliseconds: time.Since(started).Milliseconds()}
	trace.FirstRelevantRank, trace.RecallAtResults = relevance(trace.Results, card.RelevantDocumentRevisionIDs, specification.Retrieval.Results)
	return trace, nil
}

func relevance(hits []immutableretrieval.ChunkHit, relevant []string, cutoff int) (int, float64) {
	if len(relevant) == 0 {
		return 0, 0
	}
	set := make(map[string]struct{}, len(relevant))
	for _, id := range relevant {
		set[id] = struct{}{}
	}
	found := map[string]struct{}{}
	first := 0
	for _, hit := range hits {
		if cutoff > 0 && hit.Rank > cutoff {
			break
		}
		if _, ok := set[hit.DocumentRevisionID]; ok {
			if first == 0 {
				first = hit.Rank
			}
			found[hit.DocumentRevisionID] = struct{}{}
		}
	}
	return first, float64(len(found)) / float64(len(set))
}

// SortCards is provided so callers can make trace insertion ordering explicit.
func SortCards(cards []EvaluationCard) {
	sort.Slice(cards, func(i, j int) bool { return cards[i].ID < cards[j].ID })
}
