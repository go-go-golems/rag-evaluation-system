package main

import (
	"testing"

	"github.com/go-go-golems/scraper/pkg/engine/model"
)

func TestPreparationWorkerCount(t *testing.T) {
	for input, want := range map[int]int{-1: 1, 0: 1, 1: 1, 3: 3} {
		if got := preparationWorkerCount(input); got != want {
			t.Fatalf("preparationWorkerCount(%d) = %d, want %d", input, got, want)
		}
	}
}

func TestPreparationRuntimeConfigAppliesProviderConcurrencyToWorkQueues(t *testing.T) {
	config := preparationRuntimeConfig("test.sqlite", 3)
	if config.MaxWorkers != 3 {
		t.Fatalf("MaxWorkers = %d, want 3", config.MaxWorkers)
	}
	for queue, want := range map[string]int{"rag:generator": 3, "rag:embedding": 3, "rag:local": 1} {
		if got := config.Queues[model.QueueKey(queue)].MaxWorkers; got != want {
			t.Fatalf("queue %q MaxWorkers = %d, want %d", queue, got, want)
		}
	}
}
