package ragoperators

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

// EmbeddingBatchConfig defines the deterministic provider request boundary for
// an embedding model. BatchSize is part of the immutable preparation identity.
type EmbeddingBatchConfig struct {
	Model      string
	Dimensions int
	Normalize  string
	BatchSize  int
}

// EmbeddingBatchPlan partitions representations into stable provider requests.
type EmbeddingBatchPlan struct {
	Node    ragcontract.Node
	Config  EmbeddingBatchConfig
	Batches []EmbeddingBatch
}

// EmbeddingBatch is one ordered provider request.
type EmbeddingBatch struct {
	Index           int
	Representations []Representation
}

// EmbeddingBatchResult is a fully validated provider response.
type EmbeddingBatchResult struct {
	Embeddings   []Embedding
	ProviderCall bool
}

// PlanEmbeddingBatches sorts by immutable representation ID and rejects
// duplicate identity before creating a durable batch graph.
func PlanEmbeddingBatches(representations []Representation, node ragcontract.Node) (EmbeddingBatchPlan, error) {
	var cfg EmbeddingBatchConfig
	if err := decodeConfig(node.Config, &cfg); err != nil {
		return EmbeddingBatchPlan{}, err
	}
	if cfg.Model == "" || cfg.BatchSize < 1 {
		return EmbeddingBatchPlan{}, fmt.Errorf("RAG_EMBED_BATCH_CONFIG")
	}
	ordered := append([]Representation(nil), representations...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Record.ID < ordered[j].Record.ID })
	for i, item := range ordered {
		if item.Record.ID == "" || item.Text == "" {
			return EmbeddingBatchPlan{}, fmt.Errorf("RAG_EMBED_BATCH_REPRESENTATION")
		}
		if i > 0 && item.Record.ID == ordered[i-1].Record.ID {
			return EmbeddingBatchPlan{}, fmt.Errorf("RAG_EMBED_BATCH_DUPLICATE: %s", item.Record.ID)
		}
	}
	batches := make([]EmbeddingBatch, 0, (len(ordered)+cfg.BatchSize-1)/cfg.BatchSize)
	for start := 0; start < len(ordered); start += cfg.BatchSize {
		end := start + cfg.BatchSize
		if end > len(ordered) {
			end = len(ordered)
		}
		batches = append(batches, EmbeddingBatch{Index: len(batches), Representations: append([]Representation(nil), ordered[start:end]...)})
	}
	return EmbeddingBatchPlan{Node: node, Config: cfg, Batches: batches}, nil
}

// ExecuteEmbeddingBatch invokes the configured embedder and validates every
// vector before exposing it to a durable artifact or finalizer.
func ExecuteEmbeddingBatch(ctx context.Context, plan EmbeddingBatchPlan, batch EmbeddingBatch, env *Environment) (EmbeddingBatchResult, error) {
	if env == nil || env.Embedder == nil {
		return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBEDDER_UNAVAILABLE")
	}
	if len(batch.Representations) == 0 || len(batch.Representations) > plan.Config.BatchSize {
		return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_BATCH_INPUT")
	}
	modelManifest, err := resolveModel(env, plan.Config.Model)
	if err != nil {
		return EmbeddingBatchResult{}, err
	}
	texts := make([]string, len(batch.Representations))
	for i, representation := range batch.Representations {
		texts[i] = representation.Text
	}
	vectors, usage, err := env.Embedder.Embed(ctx, modelManifest.ModelID, texts)
	if err != nil {
		return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_FAILED: %w", err)
	}
	if len(vectors) != len(batch.Representations) {
		return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_COUNT: got %d want %d", len(vectors), len(batch.Representations))
	}
	result := make([]Embedding, len(batch.Representations))
	for i, vector := range vectors {
		if plan.Config.Dimensions > 0 && len(vector) != plan.Config.Dimensions {
			return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_DIMENSIONS: item %d got %d want %d", i, len(vector), plan.Config.Dimensions)
		}
		norm := 0.0
		for _, value := range vector {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_NONFINITE")
			}
			norm += value * value
		}
		if plan.Config.Normalize == "l2" {
			if norm == 0 {
				return EmbeddingBatchResult{}, fmt.Errorf("RAG_EMBED_ZERO_VECTOR")
			}
			norm = math.Sqrt(norm)
			for j := range vector {
				vector[j] /= norm
			}
		}
		digest, err := ragcontract.Digest(vector)
		if err != nil {
			return EmbeddingBatchResult{}, err
		}
		result[i] = Embedding{Record: ragcontract.EmbeddingRecord{RepresentationID: batch.Representations[i].Record.ID, ModelManifestDigest: modelManifest.Digest, Dimensions: len(vector), VectorDigest: digest, Position: int64(i)}, Vector: vector}
	}
	env.Usage.EmbeddingTokens += usage.EmbeddingTokens
	return EmbeddingBatchResult{Embeddings: result, ProviderCall: true}, nil
}
