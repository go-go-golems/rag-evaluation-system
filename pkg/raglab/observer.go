package raglab

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/pkg/errors"
)

type DomainEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type DomainMetric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type DomainArtifact struct {
	Role          string `json:"role"`
	Kind          string `json:"kind"`
	SchemaVersion string `json:"schemaVersion,omitempty"`
	MediaType     string `json:"mediaType,omitempty"`
	Name          string `json:"name"`
	Data          []byte `json:"-"`
}

// Observer receives domain observations only. It has no run ID, attempt ID,
// sequence, timestamp, terminal-summary, or database methods; those remain the
// responsibility of the calling laboratory.
type Observer interface {
	Event(context.Context, DomainEvent) error
	QueryTrace(context.Context, ragcontract.QueryTrace) error
	Metric(context.Context, DomainMetric) error
	Artifact(context.Context, DomainArtifact) error
}

type ObservationExecutionRequest struct {
	Specification ExperimentSpecification
	DatasetSplit  string
	Cards         []EvaluationCard
	Options       ExecutionOptions
}

type ObservationExecutor struct {
	backend ChannelRetriever
}

func NewObservationExecutor(backend ChannelRetriever) *ObservationExecutor {
	return &ObservationExecutor{backend: backend}
}

// Execute runs retrieval and emits observations without creating, completing,
// or mutating a laboratory run. Returning nil means domain work finished; the
// caller still owns required-measure checks and terminal attempt completion.
func (e *ObservationExecutor) Execute(ctx context.Context, request ObservationExecutionRequest, observer Observer) error {
	if e == nil || e.backend == nil || observer == nil {
		return errors.New("RAG_OBSERVER_REQUIRED: retrieval backend and observer are required")
	}
	if len(request.Cards) == 0 || request.DatasetSplit == "" {
		return errors.New("RAG_EXECUTION_INPUT_REQUIRED: cards and dataset split are required")
	}
	if err := executable(request.Specification, request.Options); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := observer.Event(ctx, DomainEvent{Type: "rag.execution.started", Payload: mustJSON(map[string]any{"queryCount": len(request.Cards)})}); err != nil {
		return err
	}
	started := time.Now()
	algorithm := &retrievalExecutor{backend: e.backend}
	metricState := newMetricAccumulator(request.Specification.Metrics)
	for _, card := range request.Cards {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := observer.Event(ctx, DomainEvent{Type: "rag.query.started", Payload: mustJSON(map[string]string{"queryCardId": card.ID})}); err != nil {
			return err
		}
		trace, err := algorithm.executeCard(ctx, request.Specification, card, request.Options)
		if err != nil {
			return err
		}
		publicTrace := exportQueryTrace(request.Specification, request.DatasetSplit, card, trace)
		if err := observer.QueryTrace(ctx, publicTrace); err != nil {
			return err
		}
		metricState.add(card, publicTrace.Results)
		if err := observer.Event(ctx, DomainEvent{Type: "rag.query.completed", Payload: mustJSON(map[string]any{"queryCardId": card.ID, "resultCount": len(publicTrace.Results)})}); err != nil {
			return err
		}
	}
	for _, metric := range metricState.metrics(time.Since(started).Milliseconds()) {
		if err := observer.Metric(ctx, metric); err != nil {
			return err
		}
	}
	return observer.Event(ctx, DomainEvent{Type: "rag.execution.completed", Payload: mustJSON(map[string]any{"queryCount": len(request.Cards)})})
}

func exportQueryTrace(specification ExperimentSpecification, split string, card EvaluationCard, input executionTrace) ragcontract.QueryTrace {
	result := ragcontract.QueryTrace{
		SchemaVersion: "rag-query-trace/v1", QueryCardID: card.ID, Query: card.Query, DatasetSplit: split,
		Relevance: ragcontract.RelevanceTrace{
			ExpectedDocumentRevisionIDs: append([]string(nil), card.RelevantDocumentRevisionIDs...),
			FirstRelevantRank:           input.FirstRelevantRank, RelevantDocumentRecall: input.RecallAtResults,
		},
		Timing: ragcontract.TimingTrace{
			EmbeddingMilliseconds: input.Timing.EmbeddingMilliseconds, RetrievalMilliseconds: input.Timing.RetrievalMilliseconds,
			FusionMilliseconds: input.Timing.FusionMilliseconds, RerankingMilliseconds: input.Timing.RerankingMilliseconds,
			TotalMilliseconds: input.Timing.TotalMilliseconds,
		},
	}
	for _, channel := range inputChannelOrder(specification, input) {
		result.Channels = append(result.Channels, ragcontract.ChannelTrace{Name: channel.name, Backend: channel.backend, Hits: exportHits(channel.hits)})
	}
	result.Results = exportHits(input.Results)
	if len(input.Fusion) > 0 {
		result.Fusion = &ragcontract.FusionTrace{Kind: "rrf", RankConstant: specification.Retrieval.Fusion.RankConstant, Hits: exportFusedHits(input.Fusion)}
	}
	if input.Reranking != nil {
		result.Reranking = &ragcontract.RerankingTrace{Kind: input.Reranking.Identity.Kind, Model: input.Reranking.Identity.Model}
		for _, candidate := range input.Reranking.Candidates {
			result.Reranking.Candidates = append(result.Reranking.Candidates, ragcontract.RerankingCandidate{CandidateID: candidate.CandidateID, PreRerankRank: candidate.PreRerankRank, RetrievalScore: candidate.RetrievalScore})
		}
		for _, reranked := range input.Reranking.Results {
			result.Reranking.Results = append(result.Reranking.Results, ragcontract.RerankingResult{CandidateID: reranked.CandidateID, Rank: reranked.Rank, Score: reranked.Score})
		}
	}
	return result
}

type orderedChannel struct {
	name, backend string
	hits          []immutableretrieval.ChunkHit
}

// executionTrace stores channels by name for efficient retrieval assembly. The
// observer contract restores specification order before publishing evidence.
func inputChannelOrder(specification ExperimentSpecification, input executionTrace) []orderedChannel {
	result := make([]orderedChannel, 0, len(specification.Retrieval.Channels))
	for _, channel := range specification.Retrieval.Channels {
		result = append(result, orderedChannel{name: channel.Name, backend: string(channel.Backend), hits: input.Channels[channel.Name]})
	}
	return result
}

func exportHits(input []immutableretrieval.ChunkHit) []ragcontract.Hit {
	result := make([]ragcontract.Hit, 0, len(input))
	for _, hit := range input {
		result = append(result, exportHit(hit))
	}
	return result
}

func exportFusedHits(input []immutableretrieval.FusedHit) []ragcontract.Hit {
	result := make([]ragcontract.Hit, 0, len(input))
	for _, hit := range input {
		result = append(result, exportHit(hit.ChunkHit))
	}
	return result
}

func exportHit(hit immutableretrieval.ChunkHit) ragcontract.Hit {
	return ragcontract.Hit{Rank: hit.Rank, ChunkID: hit.ChunkID, DocumentRevisionID: hit.DocumentRevisionID, Score: hit.Score, Title: hit.Title, URL: hit.URL, Channel: hit.Channel}
}

func mustJSON(value any) json.RawMessage {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return encoded
}

type metricAccumulator struct {
	plan       MetricsPlan
	queries    int
	answerable int
	sums       map[string]float64
}

func newMetricAccumulator(plan MetricsPlan) *metricAccumulator {
	return &metricAccumulator{plan: plan, sums: map[string]float64{}}
}

func (m *metricAccumulator) add(card EvaluationCard, results []ragcontract.Hit) {
	m.queries++
	relevant := map[string]bool{}
	for _, id := range card.RelevantDocumentRevisionIDs {
		relevant[id] = true
	}
	if len(relevant) == 0 {
		return
	}
	m.answerable++
	first := 0
	for i, hit := range results {
		if relevant[hit.DocumentRevisionID] {
			first = i + 1
			break
		}
	}
	if first > 0 {
		m.sums["mrr"] += 1 / float64(first)
	}
	for _, cutoff := range m.plan.RecallAt {
		m.sums[fmt.Sprintf("recall:%d", cutoff)] += recallAt(results, relevant, cutoff)
	}
	for _, cutoff := range m.plan.PrecisionAt {
		m.sums[fmt.Sprintf("precision:%d", cutoff)] += precisionAt(results, relevant, cutoff)
	}
	if m.plan.NDCGAt > 0 {
		m.sums[fmt.Sprintf("ndcg:%d", m.plan.NDCGAt)] += ndcgAt(results, relevant, m.plan.NDCGAt)
	}
}

func (m *metricAccumulator) metrics(wallClockMS int64) []DomainMetric {
	result := []DomainMetric{{Name: "rag.query_count", Value: float64(m.queries), Unit: "queries"}, {Name: "rag.answerable_query_count", Value: float64(m.answerable), Unit: "queries"}}
	mean := func(key string) float64 {
		if m.answerable == 0 {
			return 0
		}
		return m.sums[key] / float64(m.answerable)
	}
	if m.plan.MRR {
		result = append(result, DomainMetric{Name: "rag.mean_reciprocal_rank", Value: mean("mrr"), Unit: "ratio"})
	}
	for _, cutoff := range normalizedExportInts(m.plan.RecallAt) {
		result = append(result, DomainMetric{Name: fmt.Sprintf("rag.mean_relevant_document_recall_at_%d", cutoff), Value: mean(fmt.Sprintf("recall:%d", cutoff)), Unit: "ratio"})
	}
	for _, cutoff := range normalizedExportInts(m.plan.PrecisionAt) {
		result = append(result, DomainMetric{Name: fmt.Sprintf("rag.mean_precision_at_%d", cutoff), Value: mean(fmt.Sprintf("precision:%d", cutoff)), Unit: "ratio"})
	}
	if m.plan.NDCGAt > 0 {
		result = append(result, DomainMetric{Name: fmt.Sprintf("rag.mean_ndcg_at_%d", m.plan.NDCGAt), Value: mean(fmt.Sprintf("ndcg:%d", m.plan.NDCGAt)), Unit: "ratio"})
	}
	return append(result, DomainMetric{Name: "rag.wall_clock_duration_ms", Value: float64(wallClockMS), Unit: "ms"})
}

func recallAt(results []ragcontract.Hit, relevant map[string]bool, cutoff int) float64 {
	seen := map[string]bool{}
	for i, hit := range results {
		if i >= cutoff {
			break
		}
		if relevant[hit.DocumentRevisionID] {
			seen[hit.DocumentRevisionID] = true
		}
	}
	return float64(len(seen)) / float64(len(relevant))
}

func precisionAt(results []ragcontract.Hit, relevant map[string]bool, cutoff int) float64 {
	if cutoff <= 0 {
		return 0
	}
	hits := 0
	seen := map[string]bool{}
	for i, hit := range results {
		if i >= cutoff {
			break
		}
		if relevant[hit.DocumentRevisionID] && !seen[hit.DocumentRevisionID] {
			hits++
			seen[hit.DocumentRevisionID] = true
		}
	}
	return float64(hits) / float64(cutoff)
}

func ndcgAt(results []ragcontract.Hit, relevant map[string]bool, cutoff int) float64 {
	dcg := 0.0
	for i, hit := range results {
		if i >= cutoff {
			break
		}
		if relevant[hit.DocumentRevisionID] {
			dcg += 1 / math.Log2(float64(i+2))
		}
	}
	ideal := 0.0
	limit := len(relevant)
	if cutoff < limit {
		limit = cutoff
	}
	for i := 0; i < limit; i++ {
		ideal += 1 / math.Log2(float64(i+2))
	}
	if ideal == 0 {
		return 0
	}
	return dcg / ideal
}
