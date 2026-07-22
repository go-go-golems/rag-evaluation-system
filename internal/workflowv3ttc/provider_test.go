package workflowv3ttc

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

func TestOperatorProviderUsesValidatedRAGOperators(t *testing.T) {
	fixtures := ragoperators.NewFixtureProviders()
	admitted := 0
	provider, err := NewOperatorProvider(OperatorProviderConfig{
		GenerationNode:        ragcontract.Node{Config: json.RawMessage(`{"model":"fixture-summary-v1","prompt":"fixture-transcript-summary-v1","outputSchema":"transcript-rag-summary/v1","batchSize":1,"questionsPerChunk":2,"maxBatchRunes":1000}`)},
		EmbeddingNode:         ragcontract.Node{Config: json.RawMessage(`{"model":"fixture-hash-32-v1","dimensions":32,"normalize":"none","batchSize":8}`)},
		RawRepresentationName: "raw", MaxRepresentationsPerChunk: 8,
		ProviderProfileDigest: digestOf("a"), GenerationModelDigest: digestOf("b"), EmbeddingProfileDigest: digestOf("c"),
		AdmitGeneration: func() error { admitted++; return nil },
		ResolveEnvironment: func(context.Context) (*ragoperators.Environment, error) {
			return &ragoperators.Environment{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache(), Usage: ragoperators.Usage{Cost: map[string]float64{}}}, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	chunk := Chunk{Key: "chunk-1", Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-1", ParentUnitID: "unit-1", TextDigest: digestOf("d"), Citation: ragcontract.CitationRef{SourceID: "source-1"}}, Text: "A bounded source passage for fixture execution."}, CitationIDs: []string{"source-1"}, SourceDigest: digestOf("e")}
	generated, err := provider.Generate(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}
	if admitted != 1 {
		t.Fatalf("generation admissions = %d, want 1", admitted)
	}
	if len(generated.Value.Representations) != 3 {
		t.Fatalf("representations=%d", len(generated.Value.Representations))
	}
	embedded, err := provider.Embed(context.Background(), generated.Value)
	if err != nil {
		t.Fatal(err)
	}
	if len(embedded.Value.Representations) != 4 || len(embedded.Value.Embeddings) != 4 {
		t.Fatalf("representations=%d embeddings=%d", len(embedded.Value.Representations), len(embedded.Value.Embeddings))
	}
	for index := range embedded.Value.Embeddings {
		if embedded.Value.Embeddings[index].Record.RepresentationID != embedded.Value.Representations[index].Record.ID || len(embedded.Value.Embeddings[index].Vector) != 32 {
			t.Fatalf("embedding %d identity mismatch", index)
		}
	}
}

type recordingOperationRecorder struct {
	specs       []workflowv3.ExternalOperationSpec
	completions []workflowv3.ExternalOperationCompletion
}

func (r *recordingOperationRecorder) BeginExternalOperation(_ context.Context, spec workflowv3.ExternalOperationSpec) (workflowv3.ExternalOperationTicket, error) {
	r.specs = append(r.specs, spec)
	return workflowv3.ExternalOperationTicket{OperationID: "operation", CompletionKey: "test"}, nil
}
func (r *recordingOperationRecorder) FinishExternalOperation(_ context.Context, _ workflowv3.ExternalOperationTicket, completion workflowv3.ExternalOperationCompletion) error {
	r.completions = append(r.completions, completion)
	return nil
}

func TestOperatorProviderRecordsGenerationBeforeReturningSuccess(t *testing.T) {
	fixtures := ragoperators.NewFixtureProviders()
	provider, err := NewOperatorProvider(OperatorProviderConfig{GenerationNode: ragcontract.Node{Config: json.RawMessage(`{"model":"fixture-summary-v1","prompt":"fixture-transcript-summary-v1","outputSchema":"transcript-rag-summary/v1","batchSize":1,"questionsPerChunk":2,"maxBatchRunes":1000}`)}, EmbeddingNode: ragcontract.Node{Config: json.RawMessage(`{"model":"fixture-hash-32-v1","dimensions":32,"normalize":"none","batchSize":8}`)}, RawRepresentationName: "raw", MaxRepresentationsPerChunk: 8, ProviderProfileDigest: digestOf("a"), GenerationModelDigest: digestOf("b"), EmbeddingProfileDigest: digestOf("c"), ResolveEnvironment: func(context.Context) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache(), Usage: ragoperators.Usage{Cost: map[string]float64{}}}, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	chunk := Chunk{Key: "chunk-1", Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-1", ParentUnitID: "unit-1", TextDigest: digestOf("d"), Citation: ragcontract.CitationRef{SourceID: "source-1"}}, Text: "A bounded source passage for fixture execution."}, CitationIDs: []string{"source-1"}, SourceDigest: digestOf("e")}
	recorder := &recordingOperationRecorder{}
	_, err = provider.GenerateBatchWithOperations(context.Background(), recorder, ChunkBatch{Key: "batch-1", Chunks: []Chunk{chunk}})
	if err != nil {
		t.Fatal(err)
	}
	if len(recorder.specs) != 1 || len(recorder.completions) != 1 {
		t.Fatalf("operations=%d completions=%d", len(recorder.specs), len(recorder.completions))
	}
	if recorder.specs[0].Reservation[3].Name != "requests" || recorder.completions[0].Outcome != workflowv3.ExternalOperationOutcomeSucceeded {
		t.Fatalf("unexpected durable provider evidence: %#v %#v", recorder.specs[0], recorder.completions[0])
	}
}

func digestOf(value string) string { return "sha256:" + strings.Repeat(value, 64) }
