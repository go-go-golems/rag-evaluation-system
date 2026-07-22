package ragoperators

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestEmbeddingBatchPlanAndExecutionAreDeterministic(t *testing.T) {
	representations := []Representation{
		fixtureRepresentation("representation:b", "summary", "chunk:b", "unit:b", "two"),
		fixtureRepresentation("representation:a", "summary", "chunk:a", "unit:a", "one"),
		fixtureRepresentation("representation:c", "summary", "chunk:c", "unit:c", "three"),
	}
	plan, err := PlanEmbeddingBatches(representations, ragcontract.Node{Config: json.RawMessage(`{"model":"m","dimensions":2,"normalize":"l2","batchSize":2}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Batches) != 2 || plan.Batches[0].Representations[0].Record.ID != "representation:a" {
		t.Fatalf("plan=%#v", plan)
	}
	result, err := ExecuteEmbeddingBatch(context.Background(), plan, plan.Batches[0], &Environment{Manifests: fixtureResolver(), Embedder: fakeEmbedder{}})
	if err != nil {
		t.Fatal(err)
	}
	if !result.ProviderCall || len(result.Embeddings) != 2 || result.Embeddings[0].Record.RepresentationID != "representation:a" {
		t.Fatalf("result=%#v", result)
	}
}

func TestEmbeddingBatchRejectsDuplicateRepresentationIdentity(t *testing.T) {
	representation := fixtureRepresentation("representation:a", "summary", "chunk:a", "unit:a", "one")
	_, err := PlanEmbeddingBatches([]Representation{representation, representation}, ragcontract.Node{Config: json.RawMessage(`{"model":"m","batchSize":2}`)})
	if err == nil {
		t.Fatal("duplicate representation identity accepted")
	}
}
