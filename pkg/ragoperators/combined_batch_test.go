package ragoperators

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestCombinedBatchPlanAndExecutionValidateBeforeCaching(t *testing.T) {
	node := ragcontract.Node{Config: json.RawMessage(`{"model":"m","prompt":"p","outputSchema":"summary/v1","batchSize":2,"questionsPerChunk":2,"maxBatchRunes":100}`)}
	chunks := []Chunk{fixtureChunk("c2", "u2", "two"), fixtureChunk("c1", "u1", "one")}
	plan, err := PlanCombinedPreparation(chunks, node)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Batches) != 1 || plan.Batches[0].Chunks[0].Record.ID != "c1" {
		t.Fatalf("plan=%#v", plan)
	}
	cache := NewMemoryCache()
	bad := &fakeGenerator{combinedItems: []CombinedGenerationItem{{ChunkID: "wrong", Summary: "summary", Questions: []string{"one", "two"}}}}
	_, err = ExecuteCombinedPreparationBatch(context.Background(), plan, plan.Batches[0], &Environment{Manifests: fixtureResolver(), Generator: bad, Cache: cache})
	if err == nil || !strings.Contains(err.Error(), "CARDINALITY") {
		t.Fatalf("malformed error=%v", err)
	}
	if len(cache.values) != 0 {
		t.Fatalf("malformed response was cached: %#v", cache.values)
	}
	good := &fakeGenerator{}
	result, err := ExecuteCombinedPreparationBatch(context.Background(), plan, plan.Batches[0], &Environment{Manifests: fixtureResolver(), Generator: good, Cache: cache})
	if err != nil {
		t.Fatal(err)
	}
	if !result.ProviderCall || len(result.Representations) != 6 {
		t.Fatalf("result=%#v", result)
	}
	again, err := ExecuteCombinedPreparationBatch(context.Background(), plan, plan.Batches[0], &Environment{Manifests: fixtureResolver(), Generator: good, Cache: cache})
	if err != nil {
		t.Fatal(err)
	}
	if !again.CacheHit || good.calls != 1 {
		t.Fatalf("cache result=%#v calls=%d", again, good.calls)
	}
}
