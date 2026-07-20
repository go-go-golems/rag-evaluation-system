package preparationworkflow

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

func TestEnsureRunRejectsPublicationIdentityMismatch(t *testing.T) {
	ctx := context.Background()
	plan, err := ragoperators.PlanCombinedPreparation([]ragoperators.Chunk{{Record: ragcontract.ChunkRecord{ID: "chunk:a", TextDigest: "sha256:a"}, Text: "one"}}, ragcontract.Node{Config: json.RawMessage(`{"model":"model","prompt":"prompt","outputSchema":"schema","batchSize":1,"questionsPerChunk":1,"maxBatchRunes":100}`)})
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(filepath.Join(t.TempDir(), "workflow.sqlite"))})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	if err := RegisterWithPublication(runtime, func(context.Context, Identity) (*ragoperators.Environment, error) { return nil, nil }, func(context.Context, Identity, PublicationSpec) (PublicationTarget, error) {
		return PublicationTarget{}, nil
	}); err != nil {
		t.Fatal(err)
	}
	input := Input{Identity: Identity{SchemaVersion: "rag-preparation-workflow/v1", PreparedDigest: "sha256:wrong"}, Plan: plan, Publication: &PublicationSpec{Identity: ragengine.PreparedCorpusIdentity{SchemaVersion: "rag-prepared-corpus-manifest/v1", CorpusDigest: "sha256:corpus", PipelineDigest: "sha256:pipeline"}, RawOutputKey: "raw/out", DerivedOutputKey: "derived/out", MergedOutputKey: "merged/out", EmbeddingOutputKey: "embedding/out"}}
	if _, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("publication-mismatch"), scraperworkflow.WithRunIdentity(input.Identity)); err == nil {
		t.Fatal("publication identity mismatch accepted")
	}
}

func TestEnsureRunBuildsOneOperationPerCombinedBatch(t *testing.T) {
	ctx := context.Background()
	plan, err := ragoperators.PlanCombinedPreparation([]ragoperators.Chunk{
		{Record: ragcontract.ChunkRecord{ID: "chunk:a", TextDigest: "sha256:a"}, Text: "one"},
		{Record: ragcontract.ChunkRecord{ID: "chunk:b", TextDigest: "sha256:b"}, Text: "two"},
		{Record: ragcontract.ChunkRecord{ID: "chunk:c", TextDigest: "sha256:c"}, Text: "three"},
	}, ragcontract.Node{Config: json.RawMessage(`{"model":"model","prompt":"prompt","outputSchema":"schema","batchSize":2,"questionsPerChunk":1,"maxBatchRunes":100}`)})
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(filepath.Join(t.TempDir(), "workflow.sqlite"))})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	if err := Register(runtime, func(context.Context, Identity) (*ragoperators.Environment, error) { return nil, nil }); err != nil {
		t.Fatal(err)
	}
	input := Input{Identity: Identity{SchemaVersion: "rag-preparation-workflow/v1", PreparedDigest: "sha256:prepared"}, Plan: plan}
	handle, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("rag-prep-test"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	if !handle.Created {
		t.Fatal("first ensure did not create workflow")
	}
	snapshot, err := runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats == nil || snapshot.Stats.Total != 3 || snapshot.Stats.Ready != 2 || snapshot.Stats.Pending != 1 {
		t.Fatalf("snapshot=%#v", snapshot)
	}
	attached, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("rag-prep-test"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	if attached.Created {
		t.Fatal("second ensure unexpectedly created workflow")
	}
}
