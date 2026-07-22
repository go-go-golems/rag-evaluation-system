package workflowv3ttc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
)

type fixtureProvider struct {
	mu    sync.Mutex
	calls map[string]int
}

func (p *fixtureProvider) Generate(_ context.Context, chunk Chunk) (Result[Generated], error) {
	p.mu.Lock()
	p.calls["generate:"+chunk.Key]++
	attempt := p.calls["generate:"+chunk.Key]
	p.mu.Unlock()
	generated := Generated{Key: chunk.Key, TextDigest: chunk.TextDigest, Representation: "representation:" + chunk.Key, CitationIDs: append([]string(nil), chunk.CitationIDs...), ProviderProfileDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ModelDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
	if chunk.Key == "chunk-0007" && attempt == 1 {
		generated.Representation = ""
	}
	return Result[Generated]{Value: generated, Usage: []Usage{{Dimension: "requests", Units: 1}, {Dimension: "input_tokens", Units: 4}, {Dimension: "output_tokens", Units: 2}}}, nil
}

func (p *fixtureProvider) Embed(_ context.Context, generated Generated) (Result[Embedded], error) {
	p.mu.Lock()
	p.calls["embed:"+generated.Key]++
	p.mu.Unlock()
	return Result[Embedded]{Value: Embedded{Generated: generated, Vector: []float64{0.25, 0.5, 0.75}, EmbeddingProfileDigest: "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"}, Usage: []Usage{{Dimension: "embedding_tokens", Units: 3}}}, nil
}

func TestWorkflowV3PreparationFastIntegration(t *testing.T) {
	runPreparationIntegration(t, 65)
}

func TestP1DeterministicWorkflowCompilesAndRunsAll1807Items(t *testing.T) {
	if os.Getenv("RAG_TTC_FULL_PREFLIGHT") != "1" {
		t.Skip("set RAG_TTC_FULL_PREFLIGHT=1 for production-cardinality preflight")
	}
	runPreparationIntegration(t, 1807)
}

func runPreparationIntegration(t *testing.T, itemCount int) {
	t.Helper()
	ctx := context.Background()
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := registry.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	authored, err := workflowmodule.Author(ctx, WorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	artifacts, err := workflowv3.NewFileArtifactStore(filepath.Join(root, "artifacts"), 1<<30)
	if err != nil {
		t.Fatal(err)
	}
	items := make([]workflowv3.ManifestItem, itemCount)
	const sourceCanary = "PRIVATE_TTC_SOURCE_CANARY_8d88c0d5"
	for index := range items {
		text := fmt.Sprintf("%s source %04d", sourceCanary, index)
		digest, _ := workflowv3.Digest(text)
		chunk := Chunk{Key: fmt.Sprintf("chunk-%04d", index), Text: text, TextDigest: digest, CitationIDs: []string{fmt.Sprintf("citation-%04d", index)}, SourceDigest: "sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"}
		body, err := workflowv3.CanonicalJSON(chunk)
		if err != nil {
			t.Fatal(err)
		}
		ref, err := artifacts.Put(ctx, ChunkSchema, "application/json", body)
		if err != nil {
			t.Fatal(err)
		}
		items[index] = workflowv3.ManifestItem{Key: chunk.Key, Value: ref}
	}
	manifest, err := workflowv3.NewItemManifest(ChunkSchema, items)
	if err != nil {
		t.Fatal(err)
	}
	manifestBody, err := workflowv3.EncodeItemManifest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	manifestRef, err := artifacts.Put(ctx, workflowv3.ItemManifestSchemaV1, "application/json", manifestBody)
	if err != nil {
		t.Fatal(err)
	}

	provider := &fixtureProvider{calls: map[string]int{}}
	modules, err := workflowv3runtime.NewTaskModuleRegistry(Module(provider))
	if err != nil {
		t.Fatal(err)
	}
	databasePath := filepath.Join(root, "workflow.sqlite")
	store, err := workflowv3sqlite.Open(ctx, databasePath)
	if err != nil {
		t.Fatal(err)
	}
	engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules}
	if err := engine.Submit(ctx, "ttc-p1-1807", authored.Plan, map[string]workflowv3.ArtifactRef{"chunks": manifestRef}); err != nil {
		t.Fatal(err)
	}
	dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{ResourceGeneration: 12, ResourceEmbedding: 8, ResourceLocal: 4}, PollInterval: time.Millisecond}
	dispatchContext, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- dispatcher.Run(dispatchContext) }()
	timeout := 30 * time.Second
	if itemCount == 1807 {
		timeout = 3 * time.Minute
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case dispatchErr := <-done:
			t.Fatalf("dispatcher stopped before terminal run: %v", dispatchErr)
		default:
		}
		snapshot, err := engine.Snapshot(ctx, "ttc-p1-1807")
		if err == nil && snapshot.Status == "failed" {
			cancel()
			<-done
			for _, attempt := range snapshot.Attempts {
				if attempt.Failure != nil {
					t.Fatalf("workflow failed at %s attempt %d: %+v", attempt.NodeKey, attempt.Number, *attempt.Failure)
				}
			}
			t.Fatal("workflow failed without attempt evidence")
		}
		if err == nil && snapshot.Status == "succeeded" {
			cancel()
			<-done
			body, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["prepared"])
			if err != nil {
				t.Fatal(err)
			}
			var shard PreparedShard
			if err := workflowv3.StrictDecode(body, &shard); err != nil {
				t.Fatal(err)
			}
			lastKey := fmt.Sprintf("chunk-%04d", itemCount-1)
			if len(shard.Items) != itemCount || shard.FirstKey != "chunk-0000" || shard.LastKey != lastKey {
				t.Fatalf("unexpected shard cardinality %d %s %s", len(shard.Items), shard.FirstKey, shard.LastKey)
			}
			provider.mu.Lock()
			malformedCalls := provider.calls["generate:chunk-0007"]
			provider.mu.Unlock()
			if malformedCalls != 2 {
				t.Fatalf("malformed output attempts=%d", malformedCalls)
			}
			malformedEvidence := 0
			for _, attempt := range snapshot.Attempts {
				if attempt.Failure != nil && attempt.Failure.Class == "malformed-output" && attempt.Failure.Code == "RAG_TTC_GENERATED_INVALID" {
					malformedEvidence++
				}
			}
			if malformedEvidence != 1 {
				t.Fatalf("malformed output evidence=%d", malformedEvidence)
			}
			runID := workflowv3.RunID("ttc-p1-1807")
			budgets, err := store.BudgetSnapshot(ctx, &runID)
			if err != nil {
				t.Fatal(err)
			}
			for _, budget := range budgets {
				if budget.Reserved != 0 {
					t.Fatalf("budget reservation leaked: %+v", budget)
				}
			}
			if err := store.Close(); err != nil {
				t.Fatal(err)
			}
			database, err := os.ReadFile(databasePath)
			if err != nil {
				t.Fatal(err)
			}
			if contains(database, []byte(sourceCanary)) {
				t.Fatal("source canary leaked into workflow SQLite")
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	<-done
	snapshot, _ := engine.Snapshot(ctx, "ttc-p1-1807")
	queue, _ := store.QueueSnapshot(ctx, registry, dispatcher.Capacities, time.Now().UTC())
	t.Fatalf("workflow did not succeed: status=%s maps=%+v reductions=%+v blocked=%+v", snapshot.Status, queue.Maps, queue.Reductions, queue.BlockedByReason)
}

func contains(body, fragment []byte) bool {
	if len(fragment) == 0 {
		return true
	}
	for index := 0; index+len(fragment) <= len(body); index++ {
		matched := true
		for offset := range fragment {
			if body[index+offset] != fragment[offset] {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}
