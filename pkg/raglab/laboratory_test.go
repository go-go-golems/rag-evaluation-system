package raglab

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/experimentspec"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
)

type testStore struct{ createSpecifications, createRuns, appendEvents int }

func (s *testStore) CreateSpecification(_ context.Context, input experimentspec.Input) (*experimentrun.Specification, bool, error) {
	s.createSpecifications++
	return &experimentrun.Specification{ID: "spec-" + input.CorpusSnapshotID}, false, nil
}
func (s *testStore) CreateRun(_ context.Context, specificationID string) (*experimentrun.Run, error) {
	s.createRuns++
	return &experimentrun.Run{ID: "run-" + specificationID, ExperimentSpecID: specificationID}, nil
}
func (s *testStore) AppendEvent(_ context.Context, _, _ string, _ json.RawMessage) (*experimentrun.Event, error) {
	s.appendEvents++
	return &experimentrun.Event{Sequence: s.appendEvents}, nil
}

func validCatalog() testCatalog {
	return testCatalog{
		CorpusSnapshot("snapshot"): {Ref: CorpusSnapshot("snapshot")},
		ChunkSet("chunks"):         {Ref: ChunkSet("chunks"), CorpusSnapshotID: "snapshot"},
		BM25Index("bm25"):          {Ref: BM25Index("bm25"), ChunkSetID: "chunks"},
		EmbeddingSet("embeddings"): {Ref: EmbeddingSet("embeddings"), ChunkSetID: "chunks", Dimensions: 768},
		EvaluationDataset("eval"):  {Ref: EvaluationDataset("eval"), CorpusSnapshotID: "snapshot", Status: "candidate"},
	}
}

func TestLaboratoryReadOnlyRejectsPersistAndStart(t *testing.T) {
	spec, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	lab := NewLaboratory(validCatalog(), &testStore{}, false)
	if _, err := lab.Persist(context.Background(), spec); err == nil {
		t.Fatal("expected read-only persist rejection")
	}
	if _, err := lab.Start(context.Background(), spec); err == nil {
		t.Fatal("expected read-only start rejection")
	}
}

func TestLaboratoryPersistsThenCreatesDistinctRun(t *testing.T) {
	spec, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	store := &testStore{}
	lab := NewLaboratory(validCatalog(), store, true)
	persisted, err := lab.Persist(context.Background(), spec)
	if err != nil || persisted.Specification.ID != "spec-snapshot" || store.createSpecifications != 1 {
		t.Fatalf("persist = %#v store=%#v err=%v", persisted, store, err)
	}
	run, err := lab.Start(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	if run.ID != "run-spec-snapshot" || store.createSpecifications != 2 || store.createRuns != 1 || store.appendEvents != 1 {
		t.Fatalf("run=%#v store=%#v", run, store)
	}
}
