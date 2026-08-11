package preparationworkflow

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

func TestProviderFailureRetriesWithoutManualIntervention(t *testing.T) {
	ctx := context.Background()
	generator := &batchGenerator{providerFailOnce: map[string]bool{"chunk:a": true}, calls: map[string]int{}}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(filepath.Join(t.TempDir(), "workflow.sqlite")), MaxWorkers: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	if err := Register(runtime, func(context.Context, Identity) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: testResolver(), Generator: generator, Cache: ragoperators.NewMemoryCache()}, nil
	}); err != nil {
		t.Fatal(err)
	}
	input := testInput(t)
	handle, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("provider-retry"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(1100 * time.Millisecond)
	for range 3 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats.Succeeded != 3 || snapshot.Stats.Failed != 0 || generator.calls["chunk:a"] != 2 {
		t.Fatalf("snapshot=%#v calls=%#v", snapshot.Stats, generator.calls)
	}
}

func TestInvalidProviderResponseRetriesWithoutPersistence(t *testing.T) {
	ctx := context.Background()
	generator := &batchGenerator{invalidOnce: map[string]bool{"chunk:a": true}, calls: map[string]int{}}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(filepath.Join(t.TempDir(), "workflow.sqlite")), MaxWorkers: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer runtime.Close()
	if err := Register(runtime, func(context.Context, Identity) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: testResolver(), Generator: generator, Cache: ragoperators.NewMemoryCache()}, nil
	}); err != nil {
		t.Fatal(err)
	}
	input := testInput(t)
	handle, err := runtime.EnsureRun(ctx, PackageName, input, scraperworkflow.WithRunID("validation-retry"), scraperworkflow.WithRunIdentity(input.Identity))
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(1100 * time.Millisecond)
	for range 3 {
		if _, err := runtime.RunOnce(ctx); err != nil {
			t.Fatal(err)
		}
	}
	snapshot, err := runtime.Snapshot(ctx, handle.ID)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Stats.Succeeded != 3 || snapshot.Stats.Failed != 0 || generator.calls["chunk:a"] != 2 {
		t.Fatalf("snapshot=%#v calls=%#v", snapshot.Stats, generator.calls)
	}
}
