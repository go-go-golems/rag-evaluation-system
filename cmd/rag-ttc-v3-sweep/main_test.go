package main

import (
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/workflowv3ttc"
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

func providerMeasurement(start string, elapsed int64) workflowv3ttc.ProviderMeasurement {
	return workflowv3ttc.ProviderMeasurement{ProviderStartedAt: start, ProviderElapsedMicros: elapsed}
}

func embeddingMeasurement(start string, elapsed int64) workflowv3ttc.EmbeddingMeasurement {
	return workflowv3ttc.EmbeddingMeasurement{ProviderStartedAt: start, ProviderElapsedMicros: elapsed}
}
