package ragengine

import (
	"context"
	"os"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestFilePreparedCorpusStoreReopensStaticValuesAndRebuildsIndex(t *testing.T) {
	execution := rawExecution(t)
	corpus, dataset := fixtureData()
	engine := New(nil)
	identity := PreparedCorpusIdentity{SchemaVersion: preparedCorpusSchemaVersion, CorpusDigest: "sha256:corpus", PipelineDigest: mustDigest(execution.Pipeline), GenerationSettingsFingerprint: "sha256:g", EmbeddingFingerprint: "sha256:e"}
	store, err := NewFilePreparedCorpusStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := engine.Prepare(context.Background(), execution.Pipeline, corpus, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Put(context.Background(), prepared, identity); err != nil {
		t.Fatal(err)
	}
	if err := prepared.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, found, err := store.Open(context.Background(), engine, execution.Pipeline, corpus, Options{}, identity)
	if err != nil || !found {
		t.Fatalf("Open() found=%t err=%v", found, err)
	}
	defer func() { _ = reopened.Close() }()
	result, err := engine.Execute(context.Background(), execution, corpus, dataset, &collector{}, Options{Prepared: reopened})
	if err != nil || len(result.Traces) != len(dataset.Queries) {
		t.Fatalf("Execute() result=%#v err=%v", result, err)
	}
	if err := reopened.Close(); err != nil {
		t.Fatal(err)
	}
	path, err := store.path(identity)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"corrupt":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, found, err = store.Open(context.Background(), engine, execution.Pipeline, corpus, Options{}, identity)
	if err != nil || found {
		t.Fatalf("corrupt Open() found=%t err=%v", found, err)
	}
}

func TestNewPreparedFromStaticValuesPublishesAndReopens(t *testing.T) {
	execution := rawExecution(t)
	corpus, _ := fixtureData()
	engine := New(nil)
	original, err := engine.Prepare(context.Background(), execution.Pipeline, corpus, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = original.Close() }()
	values, closed := original.snapshot()
	if closed {
		t.Fatal("prepared unexpectedly closed")
	}
	store, err := NewFilePreparedCorpusStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	identity := PreparedCorpusIdentity{SchemaVersion: preparedCorpusSchemaVersion, CorpusDigest: "sha256:corpus", PipelineDigest: mustDigest(execution.Pipeline)}
	if _, err := PublishPreparedCorpus(context.Background(), PreparedCorpusPublication{Store: store, Engine: engine, Pipeline: execution.Pipeline, Corpus: corpus, Identity: identity, Values: values}); err != nil {
		t.Fatal(err)
	}
	reopened, found, err := store.Open(context.Background(), engine, execution.Pipeline, corpus, Options{}, identity)
	if err != nil || !found {
		t.Fatalf("Open() found=%t err=%v", found, err)
	}
	defer func() { _ = reopened.Close() }()
}

func TestFileQueryCheckpointStoreRoundTripAndRejectsCorruption(t *testing.T) {
	store, err := NewFileQueryCheckpointStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	identity := QueryCheckpointIdentity{
		SchemaVersion:   queryCheckpointSchemaVersion,
		ExecutionDigest: "sha256:execution", QueryID: "q1", QueryTextDigest: "sha256:query",
		PreparedCorpusDigest: "sha256:corpus", GeneratorFingerprint: "sha256:generator",
		RerankerFingerprint: "sha256:reranker", EvaluationPolicyDigest: "sha256:policy",
	}
	checkpoint := QueryCheckpoint{SchemaVersion: queryCheckpointSchemaVersion, Identity: identity, Trace: ragcontract.QueryTrace{SchemaVersion: ragcontract.TraceSchemaVersion}}
	if err := store.Put(context.Background(), checkpoint); err != nil {
		t.Fatal(err)
	}
	got, found, err := store.Get(context.Background(), identity)
	if err != nil || !found || got.Identity.QueryID != "q1" {
		t.Fatalf("Get() = %#v, %t, %v", got, found, err)
	}
	path, err := store.path(identity)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"invalid":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, found, err = store.Get(context.Background(), identity)
	if err != nil || found {
		t.Fatalf("corrupt Get() found=%t err=%v", found, err)
	}
}
