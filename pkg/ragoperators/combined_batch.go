package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

// CombinedPreparationPlan is the immutable, deterministic partition of chunks
// into provider requests for representations.combined-summary-questions/v1.
type CombinedPreparationPlan struct {
	Node    ragcontract.Node
	Config  CombinedPreparationConfig
	Batches []CombinedPreparationBatch
}

// CombinedPreparationConfig is the provider-independent batching contract.
type CombinedPreparationConfig struct {
	Model, Prompt, OutputSchema  string
	BatchSize, QuestionsPerChunk int
	MaxBatchRunes                int
}

// CombinedPreparationBatch is one atomic provider request. Chunks are sorted
// before planning, so its membership and request payload are stable.
type CombinedPreparationBatch struct {
	Index       int
	Chunks      []Chunk
	RequestText string
}

// CombinedBatchResult is a validated, materialized batch result. A cache hit
// still has to pass validation before it is returned.
type CombinedBatchResult struct {
	Representations []Representation
	CacheHit        bool
	ProviderCall    bool
}

// PlanCombinedPreparation validates the node configuration and deterministically
// partitions chunks. It performs no provider or cache I/O.
func PlanCombinedPreparation(chunks []Chunk, node ragcontract.Node) (CombinedPreparationPlan, error) {
	var cfg CombinedPreparationConfig
	if err := decodeConfig(node.Config, &cfg); err != nil {
		return CombinedPreparationPlan{}, err
	}
	if cfg.BatchSize < 1 || cfg.QuestionsPerChunk < 1 || cfg.MaxBatchRunes < 1 {
		return CombinedPreparationPlan{}, fmt.Errorf("RAG_COMBINED_CONFIG")
	}
	batches, err := makeCombinedBatches(chunks, combinedPreparationConfig(cfg))
	if err != nil {
		return CombinedPreparationPlan{}, err
	}
	plan := CombinedPreparationPlan{Node: node, Config: cfg, Batches: make([]CombinedPreparationBatch, len(batches))}
	for i, batch := range batches {
		plan.Batches[i] = CombinedPreparationBatch{Index: i, Chunks: append([]Chunk(nil), batch.chunks...), RequestText: batch.text}
	}
	return plan, nil
}

// ExecuteCombinedPreparationBatch performs one batch request. Provider output is
// strictly validated before it may be written to the cache or returned to a
// durable workflow executor.
func ExecuteCombinedPreparationBatch(ctx context.Context, plan CombinedPreparationPlan, batch CombinedPreparationBatch, env *Environment) (CombinedBatchResult, error) {
	if env == nil || env.Generator == nil {
		return CombinedBatchResult{}, fmt.Errorf("RAG_GENERATOR_UNAVAILABLE: combined preparation")
	}
	model, err := resolveModel(env, plan.Config.Model)
	if err != nil {
		return CombinedBatchResult{}, err
	}
	prompt, err := resolvePrompt(env, plan.Config.Prompt)
	if err != nil {
		return CombinedBatchResult{}, err
	}
	outputSchema := plan.Config.OutputSchema
	if outputSchema == "" {
		outputSchema = prompt.OutputSchema
	}
	if outputSchema == "" {
		return CombinedBatchResult{}, fmt.Errorf("RAG_COMBINED_SCHEMA")
	}
	parentDigest, err := ragcontract.Digest(struct {
		Chunks []string
		Config json.RawMessage
	}{chunkDigests(batch.Chunks), plan.Node.Config})
	if err != nil {
		return CombinedBatchResult{}, err
	}
	identity := GenerationCacheIdentityV2{SchemaVersion: "rag-generation-cache-identity/v2", Operator: combinedPreparationOperator{}.Ref(), CanonicalOperatorConfig: plan.Node.Config, ParentDigest: parentDigest, ModelManifestDigest: model.Digest, PromptManifestDigest: prompt.Digest, OutputSchemaFingerprint: outputSchemaFingerprint(env, outputSchema), EffectiveSettingsFingerprint: env.GenerationSettingsFingerprint}
	key, err := ragcontract.Digest(identity)
	if err != nil {
		return CombinedBatchResult{}, err
	}
	if env.Cache != nil {
		var envelope generationCacheEnvelopeV2
		if raw, found := env.Cache.Get(key); found && json.Unmarshal(raw, &envelope) == nil && envelope.SchemaVersion == "rag-generation-cache-envelope/v2" && cacheIdentityEqual(envelope.Identity, identity) {
			values, validationErr := validateCombinedBatch(batch.Chunks, envelope.Value.CombinedItems, plan.Config.QuestionsPerChunk, model, prompt)
			if validationErr == nil {
				return CombinedBatchResult{Representations: values, CacheHit: true}, nil
			}
		}
	}
	result, err := env.Generator.Generate(ctx, GenerationRequest{Kind: "representations.combined-summary-questions", Model: model.ModelID, Prompt: prompt.PromptID, OutputSchema: outputSchema, ParentID: parentDigest, Text: batch.RequestText, Count: plan.Config.QuestionsPerChunk})
	if err != nil {
		return CombinedBatchResult{}, err
	}
	values, err := validateCombinedBatch(batch.Chunks, result.CombinedItems, plan.Config.QuestionsPerChunk, model, prompt)
	if err != nil {
		return CombinedBatchResult{}, err
	}
	if env.Cache != nil {
		data, marshalErr := json.Marshal(generationCacheEnvelopeV2{SchemaVersion: "rag-generation-cache-envelope/v2", Identity: identity, Value: result})
		if marshalErr != nil {
			return CombinedBatchResult{}, marshalErr
		}
		env.Cache.Put(key, data)
	}
	return CombinedBatchResult{Representations: values, ProviderCall: true}, nil
}
