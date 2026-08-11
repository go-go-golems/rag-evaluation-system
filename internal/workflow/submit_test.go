package workflow

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	geppettoembeddings "github.com/go-go-golems/geppetto/pkg/embeddings"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/scraper/pkg/engine/model"
)

func TestSubmitIntakeWorkflowAndRunWorkerWithFakeProvider(t *testing.T) {
	ctx := context.Background()
	appDB := seedWorkflowTestDocument(t)
	engineDB := filepath.Join(t.TempDir(), "engine.db")
	indexRoot := filepath.Join(t.TempDir(), "indexes")

	result, err := SubmitIntakeWorkflow(ctx, SubmitIntakeRequest{
		EngineDB:       engineDB,
		DBPath:         appDB,
		WorkflowID:     "wf-submit-intake",
		DocumentIDs:    []string{"doc-1"},
		Strategy:       "fixed",
		ChunkSize:      20,
		Overlap:        5,
		BatchSize:      2,
		EmbeddingType:  "fake",
		IndexRoot:      indexRoot,
		IndexID:        "bm25-submit-test",
		ForceIndex:     true,
		EmbeddingLimit: 0,
	})
	if err != nil {
		t.Fatalf("submit workflow: %v", err)
	}
	if result.WorkflowID != "wf-submit-intake" || len(result.OpIDs) != 4 || result.StrategyID != "fixed-20-5" {
		t.Fatalf("unexpected submit result: %+v", result)
	}

	store, sched, err := NewIntakeScheduler(ctx, WorkerConfig{
		EngineDB:        engineDB,
		WorkerID:        "test-worker",
		MaxWorkers:      2,
		PollInterval:    time.Millisecond,
		LeaseDuration:   time.Minute,
		ResolveProvider: submitFakeProviderResolver(4),
		IndexRoot:       indexRoot,
	})
	if err != nil {
		t.Fatalf("new scheduler: %v", err)
	}
	defer func() { _ = store.Close() }()
	for i := 0; i < 4; i++ {
		if _, err := sched.RunOnce(ctx); err != nil {
			t.Fatalf("run cycle %d: %v", i+1, err)
		}
	}
	workflow, err := store.GetWorkflow(ctx, model.WorkflowID(result.WorkflowID))
	if err != nil {
		t.Fatalf("get workflow: %v", err)
	}
	if workflow == nil || workflow.Status != model.WorkflowStatusSucceeded {
		t.Fatalf("expected succeeded workflow, got %+v", workflow)
	}
}

type submitFakeProvider struct{ dimensions int }

func (f *submitFakeProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	vectors, err := f.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vectors[0], nil
}

func (f *submitFakeProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i, text := range texts {
		vectors[i] = make([]float32, f.dimensions)
		for j := range vectors[i] {
			vectors[i][j] = float32(len(text) + j)
		}
	}
	return vectors, nil
}

func (f *submitFakeProvider) GetModel() geppettoembeddings.EmbeddingModel {
	return geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: f.dimensions}
}

func submitFakeProviderResolver(dimensions int) ProviderResolver {
	return func(ctx context.Context, input IntakeOpInput) (*embeddingservice.ResolvedProvider, error) {
		return &embeddingservice.ResolvedProvider{
			Provider:         &submitFakeProvider{dimensions: dimensions},
			EffectiveProfile: "fake",
			ProviderType:     "fake",
			Model:            geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: dimensions},
		}, nil
	}
}
