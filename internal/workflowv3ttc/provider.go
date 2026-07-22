package workflowv3ttc

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type EnvironmentResolver func(context.Context) (*ragoperators.Environment, error)

type OperatorProviderConfig struct {
	GenerationNode             ragcontract.Node
	EmbeddingNode              ragcontract.Node
	RawRepresentationName      string
	MaxRepresentationsPerChunk int
	ProviderProfileDigest      string
	GenerationModelDigest      string
	EmbeddingProfileDigest     string
	ResolveEnvironment         EnvironmentResolver
}

type OperatorProvider struct{ config OperatorProviderConfig }

func NewOperatorProvider(config OperatorProviderConfig) (*OperatorProvider, error) {
	if config.ResolveEnvironment == nil || config.ProviderProfileDigest == "" || config.GenerationModelDigest == "" ||
		config.EmbeddingProfileDigest == "" || config.RawRepresentationName == "" || config.MaxRepresentationsPerChunk < 1 {
		return nil, fmt.Errorf("complete immutable TTC provider configuration is required")
	}
	// Planning one bounded probe validates both operator configurations without
	// constructing providers or touching credentials.
	probe := ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "probe", ParentUnitID: "probe", TextDigest: "sha256:probe", Citation: ragcontract.CitationRef{SourceID: "probe"}}, Text: "probe"}
	if _, err := ragoperators.PlanCombinedPreparation([]ragoperators.Chunk{probe}, config.GenerationNode); err != nil {
		return nil, fmt.Errorf("generation configuration: %w", err)
	}
	if _, err := ragoperators.PlanEmbeddingBatches([]ragoperators.Representation{{Record: ragcontract.RepresentationRecord{ID: "probe", ParentChunkID: "probe"}, Text: "probe"}}, config.EmbeddingNode); err != nil {
		return nil, fmt.Errorf("embedding configuration: %w", err)
	}
	return &OperatorProvider{config: config}, nil
}

func (p *OperatorProvider) Generate(ctx context.Context, chunk Chunk) (Result[Generated], error) {
	plan, err := ragoperators.PlanCombinedPreparation([]ragoperators.Chunk{chunk.Chunk}, p.config.GenerationNode)
	if err != nil || len(plan.Batches) != 1 {
		return Result[Generated]{}, &Failure{Class: "configuration", Code: "RAG_TTC_GENERATION_PLAN", Retryable: false}
	}
	env, err := p.config.ResolveEnvironment(ctx)
	if err != nil || env == nil {
		return Result[Generated]{}, &Failure{Class: "configuration", Code: "RAG_TTC_PROVIDER_RESOLUTION", Retryable: false}
	}
	before := env.Usage
	result, err := ragoperators.ExecuteCombinedPreparationBatch(ctx, plan, plan.Batches[0], env)
	if err != nil {
		return Result[Generated]{}, classifyOperatorError(err, true)
	}
	sort.Slice(result.Representations, func(i, j int) bool { return result.Representations[i].Record.ID < result.Representations[j].Record.ID })
	generated := Generated{Key: chunk.Key, Chunk: chunk.Chunk, Representations: result.Representations, CitationIDs: append([]string(nil), chunk.CitationIDs...), ProviderProfileDigest: p.config.ProviderProfileDigest, ModelDigest: p.config.GenerationModelDigest}
	usage, err := usageDelta(before, env.Usage)
	if err != nil {
		return Result[Generated]{}, err
	}
	usage, err = completeGenerationUsage(usage)
	if err != nil {
		return Result[Generated]{}, err
	}
	return Result[Generated]{Value: generated, Usage: usage}, nil
}

func (p *OperatorProvider) GenerateBatch(ctx context.Context, batch ChunkBatch) (Result[GeneratedBatch], error) {
	if batch.Key == "" || len(batch.Chunks) == 0 {
		return Result[GeneratedBatch]{}, &Failure{Class: "validation", Code: "RAG_TTC_BATCH_INVALID", Retryable: false}
	}
	chunks := make([]ragoperators.Chunk, len(batch.Chunks))
	inputRunes := 0
	for i, chunk := range batch.Chunks {
		if err := validateChunk(chunk); err != nil {
			return Result[GeneratedBatch]{}, &Failure{Class: "validation", Code: "RAG_TTC_BATCH_INVALID", Retryable: false}
		}
		chunks[i] = chunk.Chunk
		inputRunes += len([]rune(chunk.Chunk.Text))
	}
	plan, err := ragoperators.PlanCombinedPreparation(chunks, p.config.GenerationNode)
	if err != nil || len(plan.Batches) != 1 {
		return Result[GeneratedBatch]{}, &Failure{Class: "configuration", Code: "RAG_TTC_GENERATION_PLAN", Retryable: false}
	}
	env, err := p.config.ResolveEnvironment(ctx)
	if err != nil {
		return Result[GeneratedBatch]{}, &Failure{Class: "configuration", Code: "RAG_TTC_PROVIDER_RESOLUTION", Retryable: false}
	}
	before := env.Usage
	started := time.Now()
	result, err := ragoperators.ExecuteCombinedPreparationBatch(ctx, plan, plan.Batches[0], env)
	elapsed := time.Since(started).Microseconds()
	if err != nil {
		return Result[GeneratedBatch]{}, classifyOperatorError(err, true)
	}
	byChunk := map[string][]ragoperators.Representation{}
	for _, representation := range result.Representations {
		byChunk[representation.Record.ParentChunkID] = append(byChunk[representation.Record.ParentChunkID], representation)
	}
	items := make([]Generated, len(batch.Chunks))
	for i, chunk := range batch.Chunks {
		representations := byChunk[chunk.Key]
		sort.Slice(representations, func(i, j int) bool { return representations[i].Record.ID < representations[j].Record.ID })
		items[i] = Generated{Key: chunk.Key, Chunk: chunk.Chunk, Representations: representations, CitationIDs: append([]string(nil), chunk.CitationIDs...), ProviderProfileDigest: p.config.ProviderProfileDigest, ModelDigest: p.config.GenerationModelDigest}
		if err := validateGenerated(chunk, items[i]); err != nil {
			return Result[GeneratedBatch]{}, &Failure{Class: "malformed-output", Code: "RAG_TTC_GENERATED_INVALID", Retryable: true}
		}
	}
	usage, err := usageDelta(before, env.Usage)
	if err != nil {
		return Result[GeneratedBatch]{}, err
	}
	usage, err = completeGenerationUsage(usage)
	if err != nil {
		return Result[GeneratedBatch]{}, err
	}
	return Result[GeneratedBatch]{Value: GeneratedBatch{Key: batch.Key, Items: items, Measurement: ProviderMeasurement{ProviderStartedAt: started.UTC().Format(time.RFC3339Nano), ProviderElapsedMicros: elapsed, ChunkCount: len(batch.Chunks), InputRunes: inputRunes}}, Usage: usage}, nil
}

func completeGenerationUsage(usage []Usage) ([]Usage, error) {
	values := map[string]int64{"cost_microunits": 0, "input_tokens": 0, "output_tokens": 0, "requests": 1}
	for _, amount := range usage {
		if amount.Dimension == "requests" {
			return nil, &Failure{Class: "internal", Code: "RAG_TTC_USAGE_INVALID", Retryable: false}
		}
		if _, ok := values[amount.Dimension]; !ok || amount.Units < 0 {
			return nil, &Failure{Class: "internal", Code: "RAG_TTC_USAGE_INVALID", Retryable: false}
		}
		values[amount.Dimension] = amount.Units
	}
	return []Usage{{Dimension: "cost_microunits", Units: values["cost_microunits"]}, {Dimension: "input_tokens", Units: values["input_tokens"]}, {Dimension: "output_tokens", Units: values["output_tokens"]}, {Dimension: "requests", Units: 1}}, nil
}

func (p *OperatorProvider) EmbedBatch(ctx context.Context, generated GeneratedBatch) (Result[MeasuredBatch], error) {
	if generated.Key == "" || len(generated.Items) == 0 {
		return Result[MeasuredBatch]{}, &Failure{Class: "validation", Code: "RAG_TTC_GENERATED_BATCH_INVALID", Retryable: false}
	}
	started := time.Now()
	totalRepresentations := 0
	embeddingTokens := int64(0)
	for _, item := range generated.Items {
		result, err := p.Embed(ctx, item)
		if err != nil {
			return Result[MeasuredBatch]{}, err
		}
		totalRepresentations += len(result.Value.Embeddings)
		for _, amount := range result.Usage {
			if amount.Dimension == "embedding_tokens" {
				embeddingTokens += amount.Units
			}
		}
	}
	return Result[MeasuredBatch]{Value: MeasuredBatch{Key: generated.Key, Generation: generated.Measurement, Embedding: EmbeddingMeasurement{ProviderStartedAt: started.UTC().Format(time.RFC3339Nano), ProviderElapsedMicros: time.Since(started).Microseconds(), RepresentationCount: totalRepresentations, ProviderRequests: len(generated.Items)}}, Usage: []Usage{{Dimension: "embedding_tokens", Units: embeddingTokens}, {Dimension: "requests", Units: int64(len(generated.Items))}}}, nil
}

func (p *OperatorProvider) Embed(ctx context.Context, generated Generated) (Result[Embedded], error) {
	raw, err := ragoperators.RawRepresentations(p.config.RawRepresentationName, []ragoperators.Chunk{generated.Chunk})
	if err != nil {
		return Result[Embedded]{}, &Failure{Class: "configuration", Code: "RAG_TTC_RAW_REPRESENTATION", Retryable: false}
	}
	representations := append(raw, generated.Representations...)
	sort.Slice(representations, func(i, j int) bool { return representations[i].Record.ID < representations[j].Record.ID })
	if len(representations) > p.config.MaxRepresentationsPerChunk {
		return Result[Embedded]{}, &Failure{Class: "validation", Code: "RAG_TTC_REPRESENTATION_CARDINALITY", Retryable: false}
	}
	for index := 1; index < len(representations); index++ {
		if representations[index].Record.ID == representations[index-1].Record.ID {
			return Result[Embedded]{}, &Failure{Class: "validation", Code: "RAG_TTC_REPRESENTATION_DUPLICATE", Retryable: false}
		}
	}
	plan, err := ragoperators.PlanEmbeddingBatches(representations, p.config.EmbeddingNode)
	if err != nil || len(plan.Batches) != 1 {
		return Result[Embedded]{}, &Failure{Class: "configuration", Code: "RAG_TTC_EMBEDDING_PLAN", Retryable: false}
	}
	env, err := p.config.ResolveEnvironment(ctx)
	if err != nil || env == nil {
		return Result[Embedded]{}, &Failure{Class: "configuration", Code: "RAG_TTC_PROVIDER_RESOLUTION", Retryable: false}
	}
	before := env.Usage
	result, err := ragoperators.ExecuteEmbeddingBatch(ctx, plan, plan.Batches[0], env)
	if err != nil {
		return Result[Embedded]{}, classifyOperatorError(err, false)
	}
	embedded := Embedded{Generated: generated, RawRepresentations: raw, Representations: representations, Embeddings: result.Embeddings, EmbeddingProfileDigest: p.config.EmbeddingProfileDigest}
	usage, err := usageDelta(before, env.Usage)
	if err != nil {
		return Result[Embedded]{}, err
	}
	return Result[Embedded]{Value: embedded, Usage: usage}, nil
}

func usageDelta(before, after ragoperators.Usage) ([]Usage, error) {
	if after.InputTokens < before.InputTokens || after.OutputTokens < before.OutputTokens || after.EmbeddingTokens < before.EmbeddingTokens {
		return nil, &Failure{Class: "accounting", Code: "RAG_TTC_USAGE_INVALID", Retryable: false}
	}
	amounts := []Usage{}
	if delta := after.InputTokens - before.InputTokens; delta > 0 {
		amounts = append(amounts, Usage{Dimension: "input_tokens", Units: delta})
	}
	if delta := after.OutputTokens - before.OutputTokens; delta > 0 {
		amounts = append(amounts, Usage{Dimension: "output_tokens", Units: delta})
	}
	if delta := after.EmbeddingTokens - before.EmbeddingTokens; delta > 0 {
		amounts = append(amounts, Usage{Dimension: "embedding_tokens", Units: delta})
	}
	beforeCost, afterCost := 0.0, 0.0
	for _, cost := range before.Cost {
		beforeCost += cost
	}
	for _, cost := range after.Cost {
		afterCost += cost
	}
	if afterCost < beforeCost || math.IsNaN(afterCost) || math.IsInf(afterCost, 0) {
		return nil, &Failure{Class: "accounting", Code: "RAG_TTC_COST_INVALID", Retryable: false}
	}
	if delta := afterCost - beforeCost; delta > 0 {
		microunits := math.Round(delta * 1_000_000)
		if microunits <= 0 || microunits > math.MaxInt64 {
			return nil, &Failure{Class: "accounting", Code: "RAG_TTC_COST_INVALID", Retryable: false}
		}
		amounts = append(amounts, Usage{Dimension: "cost_microunits", Units: int64(microunits)})
	}
	return amounts, nil
}

func classifyOperatorError(err error, generation bool) error {
	text := err.Error()
	if generation && (strings.Contains(text, "RAG_COMBINED_RESPONSE_") || strings.Contains(text, "RAG_GENERATOR_COMBINED_JSON")) {
		return &Failure{Class: "malformed-output", Code: "RAG_TTC_GENERATED_INVALID", Retryable: true}
	}
	if strings.Contains(text, "UNAVAILABLE") || strings.Contains(text, "CONFIG") || strings.Contains(text, "MANIFEST") || strings.Contains(text, "SCHEMA") {
		return &Failure{Class: "configuration", Code: "RAG_TTC_PROVIDER_CONFIGURATION", Retryable: false}
	}
	if generation {
		return &Failure{Class: "provider", Code: "RAG_TTC_GENERATION_PROVIDER", Retryable: true}
	}
	return &Failure{Class: "provider", Code: "RAG_TTC_EMBEDDING_PROVIDER", Retryable: true}
}
