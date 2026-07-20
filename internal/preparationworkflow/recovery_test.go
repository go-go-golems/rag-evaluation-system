package preparationworkflow

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

type fixtureEmbedder struct{}

func (fixtureEmbedder) Embed(_ context.Context, _ string, texts []string) ([][]float64, ragoperators.Usage, error) {
	vectors := make([][]float64, len(texts))
	for i := range texts {
		vectors[i] = []float64{3, 4}
	}
	return vectors, ragoperators.Usage{EmbeddingTokens: int64(len(texts))}, nil
}

type batchGenerator struct {
	mu       sync.Mutex
	failOnce map[string]bool
	calls    map[string]int
}

func (g *batchGenerator) Generate(_ context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	var payload struct {
		Items []struct {
			ChunkID string `json:"chunkId"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(request.Text), &payload); err != nil {
		return ragoperators.GenerationResult{}, err
	}
	if len(payload.Items) == 0 {
		return ragoperators.GenerationResult{}, fmt.Errorf("empty fixture request")
	}
	key := payload.Items[0].ChunkID
	g.mu.Lock()
	defer g.mu.Unlock()
	g.calls[key]++
	if g.failOnce[key] {
		delete(g.failOnce, key)
		return ragoperators.GenerationResult{}, fmt.Errorf("fixture transient failure")
	}
	items := make([]ragoperators.CombinedGenerationItem, len(payload.Items))
	for i, item := range payload.Items {
		items[i] = ragoperators.CombinedGenerationItem{ChunkID: item.ChunkID, Summary: "summary " + item.ChunkID, Questions: []string{"question " + item.ChunkID}}
	}
	return ragoperators.GenerationResult{CombinedItems: items}, nil
}

func TestFailedBatchDoesNotStopSiblingAndRetryFinalizesAfterRestart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "workflow.sqlite")
	generator := &batchGenerator{failOnce: map[string]bool{"chunk:a": true}, calls: map[string]int{}}
	input := testInput(t)
	open := func() *scraperworkflow.Runtime {
		t.Helper()
		runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(path), MaxWorkers: 2})
		if err != nil {
			t.Fatal(err)
		}
		resolver := func(context.Context, Identity) (*ragoperators.Environment, error) {
			return &ragoperators.Environment{Manifests: testResolver(), Generator: generator, Cache: ragoperators.NewMemoryCache()}, nil
		}
		if err := Register(runtime, resolver); err != nil {
			_ = runtime.Close()
			t.Fatal(err)
		}
		return runtime
	}
	runtime := open()
	handle, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("recovery-test"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Close(); err != nil {
		t.Fatal(err)
	}

	runtime = open()
	defer runtime.Close()
	snapshot, err := runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats.Failed != 1 || snapshot.Stats.Succeeded != 1 || snapshot.Stats.Blocked != 1 {
		t.Fatalf("after failed batch: %#v", snapshot.Stats)
	}
	if generator.calls["chunk:c"] != 1 {
		t.Fatalf("independent sibling was not run: calls=%#v", generator.calls)
	}
	if err := runtime.RetryStep(ctx, handle.ID, "combined-0000"); err != nil {
		t.Fatal(err)
	}
	for range 3 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err = runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats.Succeeded != 3 || snapshot.Stats.Failed != 0 || snapshot.Stats.Blocked != 0 {
		t.Fatalf("after retry: %#v", snapshot.Stats)
	}
	if generator.calls["chunk:c"] != 1 {
		t.Fatalf("successful sibling was recomputed: calls=%#v", generator.calls)
	}
}

func TestEmbeddingStepsFollowCombinedSteps(t *testing.T) {
	ctx := context.Background()
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(filepath.Join(t.TempDir(), "workflow.sqlite")), MaxWorkers: 2})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	generator := &batchGenerator{calls: map[string]int{}}
	if err := Register(runtime, func(context.Context, Identity) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: testResolver(), Generator: generator, Embedder: fixtureEmbedder{}}, nil
	}); err != nil {
		t.Fatal(err)
	}
	input := testInput(t)
	input.Embedding = &EmbeddingSpec{RawRepresentationName: "raw", MaxRepresentationsPerChunk: 3, Node: ragcontract.Node{Config: json.RawMessage(`{"model":"model","dimensions":2,"normalize":"l2","batchSize":16}`)}}
	handle, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("embedding-test"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	for range 6 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats.Succeeded != 5 || snapshot.Stats.Total != 5 {
		t.Fatalf("snapshot=%#v", snapshot.Stats)
	}
}

func testInput(t *testing.T) Input {
	t.Helper()
	chunks := []ragoperators.Chunk{
		{Record: ragcontract.ChunkRecord{ID: "chunk:a", ParentUnitID: "unit:a", TextDigest: "sha256:a"}, Text: "one"},
		{Record: ragcontract.ChunkRecord{ID: "chunk:b", ParentUnitID: "unit:b", TextDigest: "sha256:b"}, Text: "two"},
		{Record: ragcontract.ChunkRecord{ID: "chunk:c", ParentUnitID: "unit:c", TextDigest: "sha256:c"}, Text: "three"},
	}
	plan, err := ragoperators.PlanCombinedPreparation(chunks, ragcontract.Node{Config: json.RawMessage(`{"model":"model","prompt":"prompt","outputSchema":"schema","batchSize":2,"questionsPerChunk":1,"maxBatchRunes":100}`)})
	if err != nil {
		t.Fatal(err)
	}
	return Input{Identity: Identity{SchemaVersion: "rag-preparation-workflow/v1", PreparedDigest: "sha256:recovery"}, Plan: plan}
}

func testResolver() ragoperators.StaticManifestResolver {
	return ragoperators.StaticManifestResolver{
		Models:  map[string]ragcontract.ModelManifest{"model": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + strings.Repeat("a", 64)}, ModelID: "fixture-model"}},
		Prompts: map[string]ragcontract.PromptManifest{"prompt": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + strings.Repeat("b", 64)}, PromptID: "fixture-prompt", OutputSchema: "schema"}},
	}
}
