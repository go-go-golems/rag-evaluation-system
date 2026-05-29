package workflow

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/engine/runner"
	"github.com/go-go-golems/scraper/pkg/engine/scheduler"
	storecontract "github.com/go-go-golems/scraper/pkg/engine/store"
	sqlitestore "github.com/go-go-golems/scraper/pkg/engine/store/sqlite"
)

func TestIntakeRunnerChunkDocumentWorkflow(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	engineStore, err := sqlitestore.Open(ctx, filepath.Join(t.TempDir(), "engine.db"))
	if err != nil {
		t.Fatalf("open engine store: %v", err)
	}
	defer engineStore.Close()

	registry := runner.NewRegistry()
	if err := registry.Register(&IntakeRunner{}); err != nil {
		t.Fatalf("register intake runner: %v", err)
	}
	sched, err := scheduler.New(engineStore, registry, scheduler.Config{
		MaxWorkers:           1,
		PollInterval:         time.Millisecond,
		DefaultLeaseDuration: time.Minute,
	}, "test-worker", nil)
	if err != nil {
		t.Fatalf("new scheduler: %v", err)
	}

	workflow := model.WorkflowRun{ID: "wf-chunk", Site: "rag-eval", Name: "Chunk document", Status: model.WorkflowStatusPending, Input: json.RawMessage(`{}`)}
	op := model.OpSpec{
		ID:         "wf-chunk:chunk",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueCPU,
		DedupKey:   "chunk:doc-1:fixed-20-5",
		Input: mustJSON(t, IntakeOpInput{
			Operation:   OperationChunkDocument,
			DBPath:      appDB,
			DocumentID:  "doc-1",
			Strategy:    "fixed",
			ChunkSize:   20,
			Overlap:     5,
			Description: "workflow test",
		}),
	}
	if err := sched.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: []model.OpSpec{op}}); err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run once: %v", err)
	}

	result, err := engineStore.GetResult(ctx, workflow.ID, op.ID)
	if err != nil {
		t.Fatalf("get op result: %v", err)
	}
	if result == nil {
		t.Fatalf("expected successful result")
	}
	if result.Error != nil && result.Error.Code != "" {
		t.Fatalf("expected successful result, got error %+v", result.Error)
	}
	var output ChunkDocumentOutput
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.DocumentID != "doc-1" || output.StrategyID != "fixed-20-5" || output.ChunkCount == 0 {
		t.Fatalf("unexpected output: %+v", output)
	}

	queries := openAppQueries(t, appDB)
	defer queries.Close()
	chunks, err := queries.ListChunks("doc-1")
	if err != nil {
		t.Fatalf("list chunks: %v", err)
	}
	if len(chunks) != output.ChunkCount {
		t.Fatalf("expected %d chunks, got %d", output.ChunkCount, len(chunks))
	}
	storedWorkflow, err := engineStore.GetWorkflow(ctx, workflow.ID)
	if err != nil {
		t.Fatalf("get workflow: %v", err)
	}
	if storedWorkflow.Status != model.WorkflowStatusSucceeded {
		t.Fatalf("expected succeeded workflow, got %s", storedWorkflow.Status)
	}
}

func TestIntakeRunnerMissingDBPathFailsNonRetryably(t *testing.T) {
	ctx := context.Background()
	r := &IntakeRunner{}
	result, err := r.Run(ctx, runner.RunContext{
		Workflow: model.WorkflowRun{ID: "wf", Site: "rag-eval"},
		Op: model.OpSpec{
			ID:    "op",
			Input: mustJSON(t, IntakeOpInput{Operation: OperationChunkDocument, DocumentID: "doc-1"}),
		},
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result == nil || result.Error == nil {
		t.Fatalf("expected op error")
	}
	if result.Error.Code != "missing_db_path" || result.Error.Retryable {
		t.Fatalf("unexpected error: %+v", result.Error)
	}
}

func seedWorkflowTestDocument(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rag-eval.db")
	queries := openAppQueries(t, path)
	defer queries.Close()
	if err := queries.InsertSource("test-source", "Test Source", "test", "{}"); err != nil {
		t.Fatalf("insert source: %v", err)
	}
	content := "This is a workflow chunking test document. It has enough text to produce several overlapping fixed-size chunks."
	if err := queries.InsertDocument("doc-1", "test-source", "doc-1", "Workflow Test", "", "", "text", content, content, "", len(content), "en", "extracted"); err != nil {
		t.Fatalf("insert document: %v", err)
	}
	return path
}

func openAppQueries(t *testing.T, path string) *db.Queries {
	t.Helper()
	database, err := db.OpenDB(path)
	if err != nil {
		t.Fatalf("open app db: %v", err)
	}
	if err := db.Migrate(database); err != nil {
		_ = database.Close()
		t.Fatalf("migrate app db: %v", err)
	}
	return db.NewQueries(database)
}
