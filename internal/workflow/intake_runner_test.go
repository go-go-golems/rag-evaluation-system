package workflow

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	geppettoembeddings "github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
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
	defer func() { _ = engineStore.Close() }()

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

	result := successfulResult(t, ctx, engineStore, workflow.ID, op.ID)
	var output ChunkDocumentOutput
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.DocumentID != "doc-1" || output.StrategyID != "fixed-20-5" || output.ChunkCount == 0 {
		t.Fatalf("unexpected output: %+v", output)
	}

	queries := openAppQueries(t, appDB)
	defer func() { _ = queries.Close() }()
	chunks, err := queries.ListChunks("doc-1")
	if err != nil {
		t.Fatalf("list chunks: %v", err)
	}
	if len(chunks) != output.ChunkCount {
		t.Fatalf("expected %d chunks, got %d", output.ChunkCount, len(chunks))
	}
	assertWorkflowStatus(t, ctx, engineStore, workflow.ID, model.WorkflowStatusSucceeded)
}

func TestIntakeRunnerPreprocessDocumentWorkflow(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	engineStore, sched := newWorkflowScheduler(t, &IntakeRunner{})
	workflow := model.WorkflowRun{ID: "wf-preprocess", Site: "rag-eval", Name: "Preprocess document", Status: model.WorkflowStatusPending, Input: json.RawMessage(`{}`)}
	op := model.OpSpec{
		ID:         "wf-preprocess:preprocess:doc-1",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueLLM,
		DedupKey:   "preprocess:doc-1:clean_text:v1:fake",
		Input: mustJSON(t, IntakeOpInput{
			Operation:                  OperationPreprocessDocument,
			DBPath:                     appDB,
			DocumentID:                 "doc-1",
			ArtifactType:               "clean_text",
			PromptVersion:              "v1",
			DocumentProcessingProvider: "fake",
			DocumentProcessingModel:    "fake-document-processor",
		}),
	}
	if err := sched.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: []model.OpSpec{op}}); err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run preprocess cycle: %v", err)
	}
	result := successfulResult(t, ctx, engineStore, workflow.ID, op.ID)
	var output PreprocessDocumentOutput
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.DocumentID != "doc-1" || output.ArtifactType != "clean_text" || output.Provider != "fake" || output.InputHash == "" {
		t.Fatalf("unexpected output: %+v", output)
	}
	queries := openAppQueries(t, appDB)
	defer func() { _ = queries.Close() }()
	artifact, ok, err := queries.GetDocumentProcessingArtifact("doc-1", "clean_text", "v1", "fake", "fake-document-processor")
	if err != nil {
		t.Fatalf("get artifact: %v", err)
	}
	if !ok || artifact.OutputText == "" || artifact.InputHash != output.InputHash {
		t.Fatalf("unexpected artifact ok=%v artifact=%+v", ok, artifact)
	}
	content, err := queries.GetDocumentContent("doc-1")
	if err != nil {
		t.Fatalf("get document content: %v", err)
	}
	if content == artifact.OutputText {
		t.Fatalf("preprocessing artifact overwrote canonical document content")
	}
	assertWorkflowStatus(t, ctx, engineStore, workflow.ID, model.WorkflowStatusSucceeded)
}

func TestIntakeRunnerEnrichChunkWorkflow(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	queries := openAppQueries(t, appDB)
	if _, err := chunkservice.NewService(queries).Apply(ctx, chunkservice.ApplyRequest{DocumentID: "doc-1", Strategy: "fixed", ChunkSize: 20, Overlap: 5}); err != nil {
		_ = queries.Close()
		t.Fatalf("seed chunks: %v", err)
	}
	chunks, err := queries.ListChunks("doc-1")
	_ = queries.Close()
	if err != nil {
		t.Fatalf("list chunks: %v", err)
	}
	if len(chunks) == 0 {
		t.Fatalf("expected seeded chunks")
	}

	engineStore, sched := newWorkflowScheduler(t, &IntakeRunner{})
	workflow := model.WorkflowRun{ID: "wf-enrich", Site: "rag-eval", Name: "Enrich chunk", Status: model.WorkflowStatusPending, Input: json.RawMessage(`{}`)}
	op := model.OpSpec{
		ID:         "wf-enrich:chunk-0",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueLLM,
		DedupKey:   "enrich:" + chunks[0].ID + ":v1:fake",
		Input: mustJSON(t, IntakeOpInput{
			Operation:               OperationEnrichChunk,
			DBPath:                  appDB,
			ChunkID:                 chunks[0].ID,
			StrategyID:              "fixed-20-5",
			PromptVersion:           "v1",
			ChunkEnrichmentProvider: "fake",
			ChunkEnrichmentModel:    "fake-chunk-enricher",
		}),
	}
	if err := sched.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: []model.OpSpec{op}}); err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run enrich cycle: %v", err)
	}
	result := successfulResult(t, ctx, engineStore, workflow.ID, op.ID)
	var output EnrichChunkOutput
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.ChunkID != chunks[0].ID || output.Provider != "fake" || output.TextHash == "" {
		t.Fatalf("unexpected output: %+v", output)
	}
	queries = openAppQueries(t, appDB)
	defer func() { _ = queries.Close() }()
	enrichment, ok, err := queries.GetChunkEnrichment(chunks[0].ID, "fixed-20-5", "v1")
	if err != nil {
		t.Fatalf("get enrichment: %v", err)
	}
	if !ok || enrichment.ShortSummary == "" || enrichment.TextHash != output.TextHash {
		t.Fatalf("unexpected enrichment ok=%v enrichment=%+v", ok, enrichment)
	}
	assertWorkflowStatus(t, ctx, engineStore, workflow.ID, model.WorkflowStatusSucceeded)
}

func TestIntakeRunnerChunkToEmbeddingWorkflow(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	engineStore, sched := newWorkflowScheduler(t, &IntakeRunner{ResolveProvider: fakeProviderResolver(4)})
	workflow := model.WorkflowRun{ID: "wf-chunk-embed", Site: "rag-eval", Name: "Chunk and embed", Status: model.WorkflowStatusPending, Input: json.RawMessage(`{}`)}
	chunkOp := model.OpSpec{
		ID:         "wf-chunk-embed:chunk",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueCPU,
		DedupKey:   "chunk:doc-1:fixed-20-5",
		Input: mustJSON(t, IntakeOpInput{
			Operation:  OperationChunkDocument,
			DBPath:     appDB,
			DocumentID: "doc-1",
			Strategy:   "fixed",
			ChunkSize:  20,
			Overlap:    5,
		}),
	}
	embedOp := model.OpSpec{
		ID:         "wf-chunk-embed:embed",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueEmbedding,
		DedupKey:   "embed:fixed-20-5:fake",
		DependsOn:  []model.Dependency{{OpID: chunkOp.ID, Required: true}},
		Input: mustJSON(t, IntakeOpInput{
			Operation:  OperationComputeEmbeddings,
			DBPath:     appDB,
			StrategyID: "fixed-20-5",
			BatchSize:  2,
		}),
	}
	freshCheckOp := model.OpSpec{
		ID:         "wf-chunk-embed:embed-fresh-check",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueEmbedding,
		DedupKey:   "embed:fixed-20-5:fake:fresh-check",
		DependsOn:  []model.Dependency{{OpID: embedOp.ID, Required: true}},
		Input: mustJSON(t, IntakeOpInput{
			Operation:  OperationComputeEmbeddings,
			DBPath:     appDB,
			StrategyID: "fixed-20-5",
			BatchSize:  2,
		}),
	}
	if err := sched.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: []model.OpSpec{chunkOp, embedOp, freshCheckOp}}); err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run chunk cycle: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run embed cycle: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run fresh-check cycle: %v", err)
	}

	var output ComputeEmbeddingsOutput
	result := successfulResult(t, ctx, engineStore, workflow.ID, embedOp.ID)
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode embedding output: %v", err)
	}
	if output.StrategyID != "fixed-20-5" || output.ProviderType != "fake" || output.Model != "fake-embedding" || output.Dimensions != 4 {
		t.Fatalf("unexpected embedding output: %+v", output)
	}
	if output.Considered == 0 || output.Computed != output.Considered {
		t.Fatalf("expected embeddings to be computed, got %+v", output)
	}

	var freshOutput ComputeEmbeddingsOutput
	freshResult := successfulResult(t, ctx, engineStore, workflow.ID, freshCheckOp.ID)
	if err := json.Unmarshal(freshResult.Data, &freshOutput); err != nil {
		t.Fatalf("decode fresh embedding output: %v", err)
	}
	if freshOutput.Computed != 0 || freshOutput.SkippedFresh != output.Considered {
		t.Fatalf("expected fresh embeddings to be skipped, got %+v", freshOutput)
	}

	queries := openAppQueries(t, appDB)
	defer func() { _ = queries.Close() }()
	var stored int
	if err := queries.DB().QueryRow(`SELECT COUNT(*) FROM chunk_embeddings WHERE strategy_id = ? AND provider = ?`, "fixed-20-5", "fake").Scan(&stored); err != nil {
		t.Fatalf("count embeddings: %v", err)
	}
	if stored != output.Computed {
		t.Fatalf("expected %d stored embeddings, got %d", output.Computed, stored)
	}
	assertWorkflowStatus(t, ctx, engineStore, workflow.ID, model.WorkflowStatusSucceeded)
}

func TestIntakeRunnerBuildBM25Workflow(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	engineStore, sched := newWorkflowScheduler(t, &IntakeRunner{})
	workflow := model.WorkflowRun{ID: "wf-bm25", Site: "rag-eval", Name: "Build BM25", Status: model.WorkflowStatusPending, Input: json.RawMessage(`{}`)}
	chunkOp := model.OpSpec{
		ID:         "wf-bm25:chunk",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueCPU,
		DedupKey:   "chunk:doc-1:fixed-20-5",
		Input:      mustJSON(t, IntakeOpInput{Operation: OperationChunkDocument, DBPath: appDB, DocumentID: "doc-1", Strategy: "fixed", ChunkSize: 20, Overlap: 5}),
	}
	indexRoot := filepath.Join(t.TempDir(), "indexes")
	bm25Op := model.OpSpec{
		ID:         "wf-bm25:index",
		WorkflowID: workflow.ID,
		Site:       workflow.Site,
		Kind:       IntakeRunnerKind,
		Queue:      QueueIndex,
		DedupKey:   "bm25:fixed-20-5",
		DependsOn:  []model.Dependency{{OpID: chunkOp.ID, Required: true}},
		Input: mustJSON(t, IntakeOpInput{
			Operation:  OperationBuildBM25,
			DBPath:     appDB,
			StrategyID: "fixed-20-5",
			IndexRoot:  indexRoot,
			IndexID:    "bm25-workflow-test",
			Force:      true,
		}),
	}
	if err := sched.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: []model.OpSpec{chunkOp, bm25Op}}); err != nil {
		t.Fatalf("create workflow: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run chunk cycle: %v", err)
	}
	if _, err := sched.RunOnce(ctx); err != nil {
		t.Fatalf("run index cycle: %v", err)
	}

	result := successfulResult(t, ctx, engineStore, workflow.ID, bm25Op.ID)
	var output BuildBM25Output
	if err := json.Unmarshal(result.Data, &output); err != nil {
		t.Fatalf("decode bm25 output: %v", err)
	}
	if output.IndexID != "bm25-workflow-test" || output.ChunkCount == 0 || output.DocumentCount != 1 || output.IndexPath == "" {
		t.Fatalf("unexpected bm25 output: %+v", output)
	}
	queries := openAppQueries(t, appDB)
	defer func() { _ = queries.Close() }()
	idx, ok, err := queries.GetSearchIndex("bm25-workflow-test")
	if err != nil {
		t.Fatalf("get search index: %v", err)
	}
	if !ok || idx.ChunkCount != output.ChunkCount {
		t.Fatalf("expected search index metadata, got ok=%v idx=%+v", ok, idx)
	}
	assertWorkflowStatus(t, ctx, engineStore, workflow.ID, model.WorkflowStatusSucceeded)
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

func newWorkflowScheduler(t *testing.T, intakeRunner *IntakeRunner) (*sqlitestore.Store, *scheduler.Scheduler) {
	t.Helper()
	ctx := context.Background()
	engineStore, err := sqlitestore.Open(ctx, filepath.Join(t.TempDir(), "engine.db"))
	if err != nil {
		t.Fatalf("open engine store: %v", err)
	}
	t.Cleanup(func() { _ = engineStore.Close() })
	registry := runner.NewRegistry()
	if err := registry.Register(intakeRunner); err != nil {
		t.Fatalf("register intake runner: %v", err)
	}
	sched, err := scheduler.New(engineStore, registry, scheduler.Config{MaxWorkers: 2, PollInterval: time.Millisecond, DefaultLeaseDuration: time.Minute}, "test-worker", nil)
	if err != nil {
		t.Fatalf("new scheduler: %v", err)
	}
	return engineStore, sched
}

func successfulResult(t *testing.T, ctx context.Context, store *sqlitestore.Store, workflowID model.WorkflowID, opID model.OpID) *model.OpResult {
	t.Helper()
	result, err := store.GetResult(ctx, workflowID, opID)
	if err != nil {
		t.Fatalf("get op result %s: %v", opID, err)
	}
	if result == nil {
		t.Fatalf("expected result for %s", opID)
	}
	if result.Error != nil && result.Error.Code != "" {
		t.Fatalf("expected successful result for %s, got error %+v", opID, result.Error)
	}
	return result
}

func assertWorkflowStatus(t *testing.T, ctx context.Context, store *sqlitestore.Store, workflowID model.WorkflowID, status model.WorkflowStatus) {
	t.Helper()
	storedWorkflow, err := store.GetWorkflow(ctx, workflowID)
	if err != nil {
		t.Fatalf("get workflow: %v", err)
	}
	if storedWorkflow == nil || storedWorkflow.Status != status {
		t.Fatalf("expected workflow %s, got %+v", status, storedWorkflow)
	}
}

type fakeWorkflowProvider struct{ dimensions int }

func (f *fakeWorkflowProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	vectors, err := f.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vectors[0], nil
}

func (f *fakeWorkflowProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i, text := range texts {
		vectors[i] = make([]float32, f.dimensions)
		for j := range vectors[i] {
			vectors[i][j] = float32(len(text) + j)
		}
	}
	return vectors, nil
}

func (f *fakeWorkflowProvider) GetModel() geppettoembeddings.EmbeddingModel {
	return geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: f.dimensions}
}

func fakeProviderResolver(dimensions int) ProviderResolver {
	return func(ctx context.Context, input IntakeOpInput) (*embeddingservice.ResolvedProvider, error) {
		return &embeddingservice.ResolvedProvider{
			Provider:         &fakeWorkflowProvider{dimensions: dimensions},
			EffectiveProfile: "fake",
			ProviderType:     "fake",
			Model:            geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: dimensions},
		}, nil
	}
}

func seedWorkflowTestDocument(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rag-eval.db")
	queries := openAppQueries(t, path)
	defer func() { _ = queries.Close() }()
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
