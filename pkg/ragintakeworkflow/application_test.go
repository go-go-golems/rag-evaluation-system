package ragintakeworkflow

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	geppettoembeddings "github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/stretchr/testify/require"
)

type fakeEmbeddingProvider struct{ dimensions int }

func (f *fakeEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	values, err := f.GenerateBatchEmbeddings(ctx, []string{text})
	return values[0], err
}
func (f *fakeEmbeddingProvider) GenerateBatchEmbeddings(_ context.Context, texts []string) ([][]float32, error) {
	ret := make([][]float32, len(texts))
	for i, text := range texts {
		ret[i] = make([]float32, f.dimensions)
		for j := range ret[i] {
			ret[i][j] = float32(len(text) + j)
		}
	}
	return ret, nil
}
func (f *fakeEmbeddingProvider) GetModel() geppettoembeddings.EmbeddingModel {
	return geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: f.dimensions}
}

func TestIntakeWorkflowV3RestartAndDatabaseEffects(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	databasePath := seedDocument(t, root)
	config := DefaultConfig(databasePath)
	config.WorkflowDatabase = filepath.Join(root, "workflow.db")
	config.ArtifactRoot = filepath.Join(root, "artifacts")
	config.Runtime.IndexRoot = filepath.Join(root, "indexes")
	config.PollInterval = time.Millisecond
	config.Runtime.ResolveProvider = func(context.Context, Request) (*embeddingservice.ResolvedProvider, error) {
		return &embeddingservice.ResolvedProvider{Provider: &fakeEmbeddingProvider{dimensions: 4}, EffectiveProfile: "fake", ProviderType: "fake", Model: geppettoembeddings.EmbeddingModel{Name: "fake-embedding", Dimensions: 4}}, nil
	}
	request, err := PrepareRequest(ctx, config.Runtime, Request{IndexID: "bm25-test", EmbeddingType: "fake", ForceBM25: true}, Selection{DocumentIDs: []string{"doc-1"}, ChunksPerDocumentToEnrich: 1})
	require.NoError(t, err)
	first, err := Open(ctx, config)
	require.NoError(t, err)
	submission, err := first.SubmitRequest(ctx, request, workflowv3.RunID("intake-restart"))
	require.NoError(t, err)
	require.NoError(t, first.Close())
	second, err := Open(ctx, config)
	require.NoError(t, err)
	defer func() { require.NoError(t, second.Close()) }()
	view, err := second.RunUntilTerminal(ctx, submission.RunID)
	require.NoError(t, err)
	for _, attempt := range view.Snapshot.Attempts {
		if attempt.Failure != nil {
			t.Logf("restart failure node=%s value=%+v", attempt.NodeKey, *attempt.Failure)
		}
	}
	require.Equal(t, "succeeded", view.Snapshot.Status, "attempts=%#v operations=%#v", view.Snapshot.Attempts, view.Operations)
	observations, err := second.Observations(ctx, submission.RunID)
	require.NoError(t, err)
	require.Equal(t, "succeeded", observations.RunStatus)
	database, err := db.OpenDB(databasePath)
	require.NoError(t, err)
	defer func() { require.NoError(t, database.Close()) }()
	var chunks, embeddings, indexes int
	require.NoError(t, database.QueryRow(`select count(*) from chunks where strategy_id='fixed-1200-150'`).Scan(&chunks))
	require.NoError(t, database.QueryRow(`select count(*) from chunk_embeddings`).Scan(&embeddings))
	require.NoError(t, database.QueryRow(`select count(*) from search_indexes where id='bm25-test'`).Scan(&indexes))
	require.Positive(t, chunks)
	require.Positive(t, embeddings)
	require.Equal(t, 1, indexes)
}

type flakyEmbeddingProvider struct {
	calls      atomic.Int64
	dimensions int
}

func (f *flakyEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	values, err := f.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return values[0], nil
}
func (f *flakyEmbeddingProvider) GenerateBatchEmbeddings(_ context.Context, texts []string) ([][]float32, error) {
	if f.calls.Add(1) == 1 {
		return nil, errors.New("secret-canary transport unavailable")
	}
	ret := make([][]float32, len(texts))
	for i := range ret {
		ret[i] = make([]float32, f.dimensions)
	}
	return ret, nil
}
func (f *flakyEmbeddingProvider) GetModel() geppettoembeddings.EmbeddingModel {
	return geppettoembeddings.EmbeddingModel{Name: "flaky", Dimensions: f.dimensions}
}

func TestProviderOperationRetryCustodyAndPrivacy(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	databasePath := seedDocument(t, root)
	authority, err := ragcontract.Digest(map[string]string{"provider": "authorized-test"})
	require.NoError(t, err)
	provider := &flakyEmbeddingProvider{dimensions: 3}
	config := DefaultConfig(databasePath)
	config.WorkflowDatabase = filepath.Join(root, "workflow.db")
	config.ArtifactRoot = filepath.Join(root, "artifacts")
	config.Runtime.ProviderAuthorityDigest = authority
	config.Runtime.ResolveProvider = func(context.Context, Request) (*embeddingservice.ResolvedProvider, error) {
		return &embeddingservice.ResolvedProvider{Provider: provider, EffectiveProfile: "authorized-test", ProviderType: "test-provider", Model: provider.GetModel()}, nil
	}
	config.PollInterval = time.Millisecond
	request, err := PrepareRequest(ctx, config.Runtime, Request{IndexID: "unused", SkipPreprocessing: true, SkipChunkEnrichment: true, SkipBM25: true, EmbeddingType: "test-provider"}, Selection{DocumentIDs: []string{"doc-1"}})
	require.NoError(t, err)
	app, err := Open(ctx, config)
	require.NoError(t, err)
	defer func() { require.NoError(t, app.Close()) }()
	_, err = app.SubmitRequest(ctx, request, "provider-retry")
	require.NoError(t, err)
	view, err := app.RunUntilTerminal(ctx, "provider-retry")
	require.NoError(t, err)
	for _, attempt := range view.Snapshot.Attempts {
		if attempt.Failure != nil {
			t.Logf("failure node=%s attempt=%d value=%+v", attempt.NodeKey, attempt.Number, *attempt.Failure)
		}
	}
	t.Logf("provider calls=%d external=%+v", provider.calls.Load(), view.Operations.ExternalOperations)
	require.Equal(t, "succeeded", view.Snapshot.Status, "attempts=%#v operations=%#v", view.Snapshot.Attempts, view.Operations)
	require.GreaterOrEqual(t, view.Operations.RetryAttempts, 1)
	require.NotNil(t, view.Operations.ExternalOperations)
	require.Equal(t, 2, view.Operations.ExternalOperations.Admitted)
	require.Equal(t, 2, view.Operations.ExternalOperations.Completed)
	err = filepath.Walk(config.ArtifactRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		require.NotContains(t, string(body), "secret-canary")
		return nil
	})
	require.NoError(t, err)
}

func TestCancelBeforeWorkerProducesCanonicalCanceledObservation(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	databasePath := seedDocument(t, root)
	config := DefaultConfig(databasePath)
	config.WorkflowDatabase = filepath.Join(root, "workflow.db")
	config.ArtifactRoot = filepath.Join(root, "artifacts")
	request, err := PrepareRequest(ctx, config.Runtime, Request{IndexID: "unused", SkipPreprocessing: true, SkipChunkEnrichment: true, SkipEmbeddings: true, SkipBM25: true}, Selection{DocumentIDs: []string{"doc-1"}})
	require.NoError(t, err)
	app, err := Open(ctx, config)
	require.NoError(t, err)
	defer func() { require.NoError(t, app.Close()) }()
	_, err = app.SubmitRequest(ctx, request, "cancel-before-worker")
	require.NoError(t, err)
	view, err := app.Cancel(ctx, "cancel-before-worker")
	require.NoError(t, err)
	require.Equal(t, "canceled", view.Snapshot.Status)
	observations, err := app.Observations(ctx, "cancel-before-worker")
	require.NoError(t, err)
	require.Equal(t, "canceled", observations.RunStatus)
}

func TestHostSecretsNeverEnterRequest(t *testing.T) {
	request, err := NormalizeRequest(Request{DocumentIDs: []string{"doc"}, IndexID: "index", SkipEmbeddings: true}, TargetDigest("db"))
	require.NoError(t, err)
	body, err := json.Marshal(request)
	require.NoError(t, err)
	require.False(t, strings.Contains(string(body), "apiKey"))
	require.False(t, strings.Contains(string(body), "baseUrl"))
}

func TestRequestStrictIdentityAndPlanStability(t *testing.T) {
	config := RuntimeConfig{DatabasePath: "/tmp/intake.db"}
	request, err := NormalizeRequest(Request{DocumentIDs: []string{"b", "a", "a"}, IndexID: "index", SkipEmbeddings: true, SkipBM25: true, SkipChunkEnrichment: true}, TargetDigest(config.DatabasePath))
	require.NoError(t, err)
	first, err := CompilePlan(context.Background(), request, config)
	require.NoError(t, err)
	second, err := CompilePlan(context.Background(), request, config)
	require.NoError(t, err)
	require.Equal(t, first.Digest, second.Digest)
	require.Equal(t, []string{"a", "b"}, request.DocumentIDs)
}

func seedDocument(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "rag.db")
	database, err := db.OpenDB(path)
	require.NoError(t, err)
	require.NoError(t, db.Migrate(database))
	queries := db.NewQueries(database)
	require.NoError(t, queries.InsertSource("source-1", "Source", "test", "{}"))
	content := "This document has enough content to build chunks embeddings and a BM25 index."
	require.NoError(t, queries.InsertDocument("doc-1", "source-1", "doc-1", "Title", "", "", "text", content, content, "", len(content), "en", "extracted"))
	require.NoError(t, database.Close())
	return path
}
