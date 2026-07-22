package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/workflowv3ttc"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

func TestPeakIntervalsTreatsTouchingSpansAsNonOverlapping(t *testing.T) {
	start := time.Unix(0, 0)
	intervals := []measuredInterval{
		{start: start, end: start.Add(time.Second)},
		{start: start.Add(time.Second), end: start.Add(2 * time.Second)},
	}
	if got := peakIntervals(intervals); got != 1 {
		t.Fatalf("peak = %d, want 1", got)
	}
}

func TestProviderIntervalsAndOverlapUseProviderWallSpans(t *testing.T) {
	batches := []batchEvidence{
		{
			Generation: providerMeasurement("2026-07-22T00:00:00Z", 20_000),
			Embedding:  embeddingMeasurement("2026-07-22T00:00:00.010Z", 20_000),
		},
		{
			Generation: providerMeasurement("2026-07-22T00:00:00.005Z", 5_000),
			Embedding:  embeddingMeasurement("2026-07-22T00:00:00.030Z", 5_000),
		},
	}
	generation, embedding, err := providerIntervals(batches)
	if err != nil {
		t.Fatal(err)
	}
	if got := peakIntervals(generation); got != 2 {
		t.Fatalf("generation peak = %d, want 2", got)
	}
	if got := overlapIntervals(generation, embedding); got != 10_000 {
		t.Fatalf("provider overlap = %d, want 10000", got)
	}
}

func TestWriteFailedCellCheckpointExcludesFailureMessage(t *testing.T) {
	root := t.TempDir()
	snapshot := workflowv3.RunSnapshot{RunID: "run-1", Status: "failed", Attempts: []workflowv3.Attempt{{ResourceClass: workflowv3ttc.ResourceGeneration, Status: "failed", Failure: &workflowv3.Failure{Class: "provider", Code: "SAFE_CODE", Message: "sensitive provider body"}}}}
	budget := []workflowv3.BudgetProgress{{RunID: "run-1", Account: "generation", Dimension: "requests", Limit: 2, Used: 1, Remaining: 1}}
	if err := writeFailedCellCheckpoint(root, "cell-00", snapshot, workflowv3ttc.SweepCell{ChunksPerRequest: 1, Concurrency: 1, Replicate: 1}, "terminal", budget, nil, failedOperationReduction{}); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(root, "failures", "cell-00.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "sensitive provider body") || !strings.Contains(string(body), "SAFE_CODE") {
		t.Fatalf("unexpected failed checkpoint: %s", body)
	}
}

func TestReduceFailedOperationsUsesOnlyClosedOperationFields(t *testing.T) {
	started := time.Date(2026, 7, 22, 12, 0, 0, 0, time.UTC)
	operations := []workflowv3.ExternalOperation{
		{Kind: workflowv3.ExternalOperationKind{Name: "provider.generate", Version: "v1"}, Completion: &workflowv3.ExternalOperationCompletion{ProviderStartedAt: started, ElapsedMicros: 10, Outcome: workflowv3.ExternalOperationOutcomeSucceeded}},
		{Kind: workflowv3.ExternalOperationKind{Name: "provider.embed", Version: "v1"}, Completion: &workflowv3.ExternalOperationCompletion{ProviderStartedAt: started.Add(5 * time.Microsecond), ElapsedMicros: 10, Outcome: workflowv3.ExternalOperationOutcomeFailed}},
		{Kind: workflowv3.ExternalOperationKind{Name: "provider.embed", Version: "v1"}},
	}
	reduction := reduceFailedOperations(operations)
	if reduction.Admitted != 3 || reduction.Completed != 2 || reduction.Incomplete != 1 || reduction.ProviderElapsedMicros != 20 || reduction.ProviderPeakActive != 2 || reduction.ProviderOverlapMicros != 5 || reduction.GenerationOperationCount != 1 || reduction.EmbeddingOperationCount != 2 || reduction.Outcomes[workflowv3.ExternalOperationOutcomeFailed] != 1 {
		t.Fatalf("unexpected reduction: %#v", reduction)
	}
}

func providerMeasurement(start string, elapsed int64) workflowv3ttc.ProviderMeasurement {
	return workflowv3ttc.ProviderMeasurement{ProviderStartedAt: start, ProviderElapsedMicros: elapsed}
}

func embeddingMeasurement(start string, elapsed int64) workflowv3ttc.EmbeddingMeasurement {
	return workflowv3ttc.EmbeddingMeasurement{ProviderStartedAt: start, ProviderElapsedMicros: elapsed}
}
