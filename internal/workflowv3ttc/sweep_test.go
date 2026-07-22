package workflowv3ttc

import "testing"

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
