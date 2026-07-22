package workflowv3ttc

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const (
	SweepGenerationInputTokensPerRequest  int64 = 16_384
	SweepGenerationOutputTokensPerRequest int64 = 8_192
	// Umans Flash's published rates are USD $0.15/M input tokens and $1.00/M
	// output tokens. Integer microunits and ceiling division keep admission exact.
	SweepGenerationInputCostMicrounitsPerMillion  int64 = 150_000
	SweepGenerationOutputCostMicrounitsPerMillion int64 = 1_000_000
	SweepGenerationCacheReadMicrounitsPerMillion  int64 = 50_000
	SweepGenerationCacheWriteMicrounitsPerMillion int64 = 150_000
	SweepGenerationCostMicrounitsPerRequest       int64 = 10_650
)

type SweepCell struct {
	ChunksPerRequest int `json:"chunksPerRequest"`
	Concurrency      int `json:"concurrency"`
	Replicate        int `json:"replicate"`
}

type SweepSpec struct {
	ChunkCount      int   `json:"chunkCount"`
	BatchSizes      []int `json:"batchSizes"`
	Concurrency     []int `json:"concurrency"`
	Replicates      int   `json:"replicates"`
	MaximumRequests int   `json:"maximumRequests"`
}

type SweepPlan struct {
	Cells           []SweepCell `json:"cells"`
	PlannedRequests int         `json:"plannedRequests"`
	Digest          string      `json:"digest"`
}

func PlanSweep(spec SweepSpec) (SweepPlan, error) {
	if spec.ChunkCount < 1 || spec.Replicates < 1 || spec.MaximumRequests < 1 || len(spec.BatchSizes) == 0 || len(spec.Concurrency) == 0 {
		return SweepPlan{}, fmt.Errorf("RAG_TTC_SWEEP_SPEC")
	}
	allowedBatch := map[int]bool{1: true, 2: true, 4: true, 8: true}
	seenBatch, seenConcurrency := map[int]bool{}, map[int]bool{}
	for _, n := range spec.BatchSizes {
		if !allowedBatch[n] || seenBatch[n] {
			return SweepPlan{}, fmt.Errorf("RAG_TTC_SWEEP_BATCH")
		}
		seenBatch[n] = true
	}
	for _, n := range spec.Concurrency {
		if n < 1 || n > 4 || seenConcurrency[n] {
			return SweepPlan{}, fmt.Errorf("RAG_TTC_SWEEP_CONCURRENCY")
		}
		seenConcurrency[n] = true
	}
	batches := append([]int(nil), spec.BatchSizes...)
	sort.Ints(batches)
	capacities := append([]int(nil), spec.Concurrency...)
	sort.Ints(capacities)
	plan := SweepPlan{}
	for replicate := 1; replicate <= spec.Replicates; replicate++ {
		for ci, concurrency := range capacities {
			order := append([]int(nil), batches...)
			if (replicate+ci)%2 == 0 {
				for left, right := 0, len(order)-1; left < right; left, right = left+1, right-1 {
					order[left], order[right] = order[right], order[left]
				}
			}
			for _, batch := range order {
				plan.Cells = append(plan.Cells, SweepCell{ChunksPerRequest: batch, Concurrency: concurrency, Replicate: replicate})
				plan.PlannedRequests += (spec.ChunkCount + batch - 1) / batch
			}
		}
	}
	if plan.PlannedRequests > spec.MaximumRequests {
		return SweepPlan{}, fmt.Errorf("RAG_TTC_SWEEP_REQUEST_BUDGET")
	}
	digest, err := workflowv3.Digest(struct {
		Spec            SweepSpec   `json:"spec"`
		Cells           []SweepCell `json:"cells"`
		PlannedRequests int         `json:"plannedRequests"`
	}{spec, plan.Cells, plan.PlannedRequests})
	if err != nil {
		return SweepPlan{}, err
	}
	plan.Digest = digest
	return plan, nil
}

func MaterializeBatches(ctx context.Context, artifacts workflowv3.ArtifactStore, chunks []Chunk, size int) (workflowv3.ArtifactRef, error) {
	if artifacts == nil || size < 1 || len(chunks) == 0 {
		return workflowv3.ArtifactRef{}, fmt.Errorf("RAG_TTC_SWEEP_BATCH_INPUT")
	}
	items := make([]workflowv3.ManifestItem, 0, (len(chunks)+size-1)/size)
	for offset := 0; offset < len(chunks); offset += size {
		end := offset + size
		if end > len(chunks) {
			end = len(chunks)
		}
		batchChunks := append([]Chunk(nil), chunks[offset:end]...)
		key := fmt.Sprintf("batch-%04d", len(items))
		body, err := workflowv3.CanonicalJSON(ChunkBatch{Key: key, Chunks: batchChunks})
		if err != nil {
			return workflowv3.ArtifactRef{}, err
		}
		ref, err := artifacts.Put(ctx, ChunkBatchSchema, "application/json", body)
		if err != nil {
			return workflowv3.ArtifactRef{}, err
		}
		items = append(items, workflowv3.ManifestItem{Key: key, Value: ref})
	}
	return putManifest(ctx, artifacts, ChunkBatchSchema, items)
}
