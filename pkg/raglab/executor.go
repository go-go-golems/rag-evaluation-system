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
	Reranker Reranker
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

func (e *Executor) Execute(ctx context.Context, runID string, specification ExperimentSpecification, cards []EvaluationCard, options ExecutionOptions) (*ExecutionResult, error) {
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
	if !emptyFilter(specification.Retrieval.Filter) {
		return errors.New("RAG_EXECUTION_UNSUPPORTED: global filters are not executable until every channel enforces and traces them")
	}
	for _, representation := range specification.Inputs.Representations {
		if representation.Kind != RawChunksRepresentation {
			return errors.New("RAG_EXECUTION_UNSUPPORTED: materialized representations need a representation executor")
		}
	}
	if specification.Retrieval.Collapse == CollapseParentChunk {
		return errors.New("RAG_EXECUTION_UNSUPPORTED: parentChunk collapse needs materialized parent mappings")
	}
	for _, channel := range specification.Retrieval.Channels {
		if !emptyFilter(channel.Filter) {
			return errors.Errorf("RAG_EXECUTION_UNSUPPORTED: channel %q filter is not executable until retrieval enforces and traces it", channel.Name)
		}
		if channel.Backend == VectorBackend && options.Embedder == nil {
			return errors.New("RAG_EMBEDDER_REQUIRED: vector retrieval needs an explicit query embedder")
		}
	}
	if specification.Retrieval.Reranking != nil && options.Reranker == nil {
		return errors.New("RAG_RERANKER_REQUIRED: reranking needs an explicit reranker")
	}
	return nil
}

func emptyFilter(filter FilterSpec) bool {
	return len(filter.SourceIDs) == 0 && len(filter.DocumentIDs) == 0 && len(filter.ContentTypes) == 0 && len(filter.MetadataEquals) == 0
}

type executionTrace struct {
	QueryID           string                                   `json:"queryId"`
	Query             string                                   `json:"query"`
	Channels          map[string][]immutableretrieval.ChunkHit `json:"channels"`
	Results           []immutableretrieval.ChunkHit            `json:"results"`
	Fusion            []immutableretrieval.FusedHit            `json:"fusion,omitempty"`
	Reranking         *rerankingTrace                          `json:"reranking,omitempty"`
	Timing            executionTiming                          `json:"timing"`
	FirstRelevantRank int                                      `json:"firstRelevantRank,omitempty"`
	RecallAtResults   float64                                  `json:"recallAtResults,omitempty"`
}

type rerankingTrace struct {
	Identity   RerankerIdentity          `json:"identity"`
	Candidates []rerankingCandidateTrace `json:"candidates"`
	Results    []RerankResult            `json:"results"`
}

type rerankingCandidateTrace struct {
	CandidateID    string  `json:"candidateId"`
	PreRerankRank  int     `json:"preRerankRank"`
	RetrievalScore float64 `json:"retrievalScore"`
}

type executionTiming struct {
	EmbeddingMilliseconds int64 `json:"embeddingMilliseconds"`
	RetrievalMilliseconds int64 `json:"retrievalMilliseconds"`
	FusionMilliseconds    int64 `json:"fusionMilliseconds"`
	RerankingMilliseconds int64 `json:"rerankingMilliseconds"`
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
		if err := ctx.Err(); err != nil {
			return executionTrace{}, err
		}
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
		if specification.Retrieval.Reranking == nil && specification.Retrieval.Collapse == CollapseDocument {
			trace.Results = immutableretrieval.CollapseDocuments(trace.Results)
		}
		if specification.Retrieval.Reranking == nil && len(trace.Results) > specification.Retrieval.Results {
			trace.Results = trace.Results[:specification.Retrieval.Results]
		}
	} else {
		weights := map[string]float64(nil)
		if specification.Retrieval.Fusion != nil {
			weights = specification.Retrieval.Fusion.Weights
		}
		candidateLimit := specification.Retrieval.Results
		if specification.Retrieval.Reranking != nil {
			candidateLimit = specification.Retrieval.Reranking.CandidateCount
		}
		trace.Fusion = immutableretrieval.FuseWeightedRRF(channels, specification.Retrieval.Fusion.RankConstant, weights, candidateLimit)
		trace.Results = make([]immutableretrieval.ChunkHit, len(trace.Fusion))
		for i := range trace.Fusion {
			trace.Results[i] = trace.Fusion[i].ChunkHit
		}
	}
	fusionMilliseconds := time.Since(fusionStarted).Milliseconds()
	var rerankingMilliseconds int64
	if specification.Retrieval.Reranking != nil {
		if err := ctx.Err(); err != nil {
			return executionTrace{}, err
		}
		rerankingStarted := time.Now()
		if err := applyReranking(ctx, &trace, card.Query, specification.Retrieval, options.Reranker); err != nil {
			return executionTrace{}, err
		}
		rerankingMilliseconds = time.Since(rerankingStarted).Milliseconds()
	}
	trace.Timing = executionTiming{EmbeddingMilliseconds: embeddingMilliseconds, RetrievalMilliseconds: retrievalMilliseconds, FusionMilliseconds: fusionMilliseconds, RerankingMilliseconds: rerankingMilliseconds, TotalMilliseconds: time.Since(started).Milliseconds()}
	trace.FirstRelevantRank, trace.RecallAtResults = relevance(trace.Results, card.RelevantDocumentRevisionIDs, specification.Retrieval.Results)
	return trace, nil
}

func applyReranking(ctx context.Context, trace *executionTrace, query string, plan RetrievalPlan, reranker Reranker) error {
	policy := plan.Reranking
	if policy == nil || reranker == nil {
		return errors.New("RAG_RERANKER_REQUIRED: reranking policy and capability are required")
	}
	identity := reranker.Identity()
	if identity.Model != policy.Model {
		return errors.Errorf("RAG_RERANKER_IDENTITY_MISMATCH: experiment requires %q, runtime provides %q", policy.Model, identity.Model)
	}
	candidates := append([]immutableretrieval.ChunkHit(nil), trace.Results...)
	if len(candidates) > policy.CandidateCount {
		candidates = candidates[:policy.CandidateCount]
	}
	request := RerankRequest{Query: query, Candidates: make([]RerankCandidate, 0, len(candidates)), TopN: len(candidates)}
	trace.Reranking = &rerankingTrace{Identity: identity, Candidates: make([]rerankingCandidateTrace, 0, len(candidates))}
	byID := make(map[string]immutableretrieval.ChunkHit, len(candidates))
	for i, hit := range candidates {
		if hit.Text == "" {
			return errors.Errorf("RAG_RERANK_CANDIDATE_TEXT_REQUIRED: chunk %q has no text", hit.ChunkID)
		}
		request.Candidates = append(request.Candidates, RerankCandidate{ID: hit.ChunkID, Text: hit.Text, OriginalRank: i + 1, RetrievalScore: hit.Score})
		trace.Reranking.Candidates = append(trace.Reranking.Candidates, rerankingCandidateTrace{CandidateID: hit.ChunkID, PreRerankRank: i + 1, RetrievalScore: hit.Score})
		byID[hit.ChunkID] = hit
	}
	results, err := reranker.Rerank(ctx, request)
	if err != nil {
		return errors.Wrap(err, "rerank retrieved candidates")
	}
	trace.Reranking.Results = append([]RerankResult(nil), results...)
	if len(results) > policy.Results {
		results = results[:policy.Results]
	}
	trace.Results = make([]immutableretrieval.ChunkHit, 0, len(results))
	for _, result := range results {
		hit, ok := byID[result.CandidateID]
		if !ok {
			return errors.Errorf("RAG_RERANK_RESPONSE_CANDIDATE_UNKNOWN: candidate %q was not submitted", result.CandidateID)
		}
		hit.Score = result.Score
		hit.Rank = result.Rank
		trace.Results = append(trace.Results, hit)
	}
	if plan.Collapse == CollapseDocument {
		trace.Results = immutableretrieval.CollapseDocuments(trace.Results)
	}
	if len(trace.Results) > plan.Results {
		trace.Results = trace.Results[:plan.Results]
	}
	for i := range trace.Results {
		trace.Results[i].Rank = i + 1
	}
	return nil
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
