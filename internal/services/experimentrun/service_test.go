package experimentrun

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutablechunk"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

func TestAppendOnlySpecificationRunEventsTraceAndSummary(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	queries := db.NewQueries(database)
	plan := &ttcimport.Plan{Manifest: ttcimport.Manifest{
		SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "fixture", SourceExportSHA256: "fixture", SelectionAlgorithm: "fixture", KindQuotas: map[string]int{"faq": 1},
	}, Documents: []ttcimport.SourceDocument{{ID: "wp:1", Kind: "faq", Title: "Fixture", URL: "https://example.test/1", SearchText: "Blue cypress fixture text."}}}
	snapshot, err := corpussnapshot.Persist(context.Background(), queries, plan, corpussnapshot.PersistRequest{SourceByteSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	chunks, err := immutablechunk.Build(context.Background(), queries, immutablechunk.Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "fixed", ChunkSize: 200})
	if err != nil {
		t.Fatal(err)
	}
	bm25, err := immutableretrieval.BuildBM25(context.Background(), queries, immutableretrieval.BM25BuildRequest{ChunkSetID: chunks.ChunkSetID, ArtifactRoot: filepath.Join(t.TempDir(), "bm25")})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(queries)
	input := SpecificationInput{CorpusSnapshotID: snapshot.SnapshotID, ChunkSetID: chunks.ChunkSetID, BM25ArtifactID: bm25.ArtifactID, EvaluationDatasetID: "candidate:ttc-baseline-v1", Config: map[string]any{"rrf_k": 60, "limit": 10}}
	spec, reused, err := service.CreateSpecification(context.Background(), input)
	if err != nil || reused {
		t.Fatalf("first specification = %#v reused=%v err=%v", spec, reused, err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(spec.Manifest, &manifest); err != nil {
		t.Fatalf("decode exported manifest: %v", err)
	}
	if manifest["schema_version"] != "rag-eval-experiment-spec/v1" {
		t.Fatalf("exported manifest schema_version = %#v", manifest["schema_version"])
	}
	second, reused, err := service.CreateSpecification(context.Background(), input)
	if err != nil || !reused || second.ID != spec.ID {
		t.Fatalf("second specification = %#v reused=%v err=%v", second, reused, err)
	}
	run, err := service.CreateRun(context.Background(), spec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.Events) != 1 || run.Events[0].Type != "created" {
		t.Fatalf("created run = %#v", run)
	}
	if event, err := service.AppendEvent(context.Background(), run.ID, "retrieval_started", []byte(`{"cards":1}`)); err != nil || event.Sequence != 2 {
		t.Fatalf("event = %#v err=%v", event, err)
	}
	if err := service.RecordQueryTrace(context.Background(), run.ID, QueryTraceInput{QueryCardID: "ttc-eval-001", Trace: []byte(`{"query":"blue cypress"}`), Metrics: []byte(`{"recall_at_10":1}`), Timing: []byte(`{"total_ms":3}`), Cost: []byte(`{"billed_usd":0}`), Storage: []byte(`{"trace_bytes":10}`)}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteRun(context.Background(), run.ID, SummaryInput{Status: "succeeded", Metrics: []byte(`{"mrr":1}`), Cost: []byte(`{"billed_usd":0}`), Storage: []byte(`{"trace_bytes":10}`), Error: []byte(`{}`)}); err != nil {
		t.Fatal(err)
	}
	completed, err := service.GetRun(context.Background(), run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != "succeeded" || completed.Summary == nil || len(completed.Events) != 3 {
		t.Fatalf("completed run = %#v", completed)
	}
	traces, err := service.ListQueryTraces(context.Background(), run.ID)
	if err != nil || len(traces) != 1 || traces[0].QueryCardID != "ttc-eval-001" {
		t.Fatalf("traces = %#v err=%v", traces, err)
	}
	if _, err := service.AppendEvent(context.Background(), run.ID, "should_fail", []byte(`{}`)); err == nil {
		t.Fatal("expected append after terminal summary to fail")
	}
	if err := service.RecordQueryTrace(context.Background(), run.ID, QueryTraceInput{QueryCardID: "ttc-eval-002", Trace: []byte(`{}`), Metrics: []byte(`{}`), Timing: []byte(`{}`), Cost: []byte(`{}`), Storage: []byte(`{}`)}); err == nil {
		t.Fatal("expected trace after terminal summary to fail")
	}
	if _, err := database.Exec(`UPDATE experiment_runs SET experiment_spec_id = ? WHERE id = ?`, spec.ID, run.ID); err == nil {
		t.Fatal("expected immutable-run update trigger to reject update")
	}
	// The database deliberately has one connection. List methods must close
	// their ID cursor before hydrating each item through a second query.
	specifications, err := service.ListSpecifications(context.Background())
	if err != nil || len(specifications) != 1 || specifications[0].ID != spec.ID {
		t.Fatalf("specifications = %#v err=%v", specifications, err)
	}
	runs, err := service.ListRuns(context.Background(), spec.ID)
	if err != nil || len(runs) != 1 || runs[0].ID != run.ID || runs[0].Status != "succeeded" {
		t.Fatalf("runs = %#v err=%v", runs, err)
	}
}
