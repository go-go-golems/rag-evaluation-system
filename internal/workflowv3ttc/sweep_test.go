package workflowv3ttc

import (
	"context"
	"testing"

	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
)

func TestSweepGenerationReservationMatchesAuthoredPlan(t *testing.T) {
	if got := (SweepGenerationInputTokensPerRequest*SweepGenerationInputCostMicrounitsPerMillion+999_999)/1_000_000 + (SweepGenerationOutputTokensPerRequest*SweepGenerationOutputCostMicrounitsPerMillion+999_999)/1_000_000; got != SweepGenerationCostMicrounitsPerRequest {
		t.Fatalf("generation cost reservation = %d, want %d", got, SweepGenerationCostMicrounitsPerRequest)
	}
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := registry.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	authored, err := workflowmodule.Author(context.Background(), SweepWorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		t.Fatal(err)
	}
	var generation map[string]int64
	for _, mapped := range authored.Plan.Maps {
		if mapped.Key != "generate-batches" || mapped.Budget == nil {
			continue
		}
		generation = map[string]int64{}
		for _, amount := range mapped.Budget.Effective {
			generation[amount.Dimension] = amount.Units
		}
	}
	if generation["requests"] != 1 || generation["input_tokens"] != SweepGenerationInputTokensPerRequest || generation["output_tokens"] != SweepGenerationOutputTokensPerRequest || generation["cost_microunits"] != SweepGenerationCostMicrounitsPerRequest {
		t.Fatalf("effective generation reservation = %#v", generation)
	}
}

func TestPlanSweepExactRequestArithmeticAndHardConcurrency(t *testing.T) {
	plan, err := PlanSweep(SweepSpec{ChunkCount: 16, BatchSizes: []int{1, 2, 4, 8}, Concurrency: []int{1, 2, 4}, Replicates: 1, MaximumRequests: 90})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Cells) != 12 || plan.PlannedRequests != 90 || !validDigest(plan.Digest) {
		t.Fatalf("plan=%+v", plan)
	}
	for _, cell := range plan.Cells {
		if cell.Concurrency < 1 || cell.Concurrency > 4 {
			t.Fatalf("cell=%+v", cell)
		}
	}
	if _, err := PlanSweep(SweepSpec{ChunkCount: 16, BatchSizes: []int{1, 2, 4, 8}, Concurrency: []int{1, 2, 4}, Replicates: 1, MaximumRequests: 89}); err == nil {
		t.Fatal("request ceiling accepted")
	}
	if _, err := PlanSweep(SweepSpec{ChunkCount: 16, BatchSizes: []int{1}, Concurrency: []int{5}, Replicates: 1, MaximumRequests: 16}); err == nil {
		t.Fatal("concurrency above four accepted")
	}
}
