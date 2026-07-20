// rag-preparation-smoke runs a deliberately tiny, real-provider durable
// preparation workflow. It is an operator smoke command, not a benchmark.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

func main() {
	providerConfig := flag.String("provider-config", "", "real-provider host configuration YAML")
	stateDB := flag.String("state-db", "", "durable scraper SQLite workflow database")
	text := flag.String("text", "A young tree benefits from consistent watering after planting.", "one small source chunk")
	withEmbedding := flag.Bool("with-embedding", false, "also execute one durable embedding batch")
	flag.Parse()
	if *providerConfig == "" || *stateDB == "" {
		fmt.Fprintln(os.Stderr, "--provider-config and --state-db are required")
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	providers, err := ragproviders.Load(ctx, *providerConfig)
	if err != nil {
		fail(err)
	}
	defer providers.Close()
	digest, err := ragcontract.Digest(*text)
	if err != nil {
		fail(err)
	}
	plan, err := ragoperators.PlanCombinedPreparation([]ragoperators.Chunk{{Record: ragcontract.ChunkRecord{ID: "smoke-chunk-0000", ParentUnitID: "smoke-unit-0000", TextDigest: digest, Citation: ragcontract.CitationRef{SourceID: "smoke-source"}}, Text: *text, ManifestDigest: digest}}, ragcontract.Node{Config: []byte(`{"model":"generator-umans-flash","prompt":"ttc-combined-preparation-v2","outputSchema":"rag-combined-preparation/v2","batchSize":1,"questionsPerChunk":4,"maxBatchRunes":6000}`)})
	if err != nil {
		fail(err)
	}
	identityDigest, err := ragcontract.Digest(struct {
		Plan          ragoperators.CombinedPreparationPlan
		Provider      string
		WithEmbedding bool
	}{plan, providers.ProfileID, *withEmbedding})
	if err != nil {
		fail(err)
	}
	identity := preparationworkflow.Identity{SchemaVersion: "rag-preparation-workflow/v1", PreparedDigest: identityDigest}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(*stateDB), WorkerID: "rag-preparation-smoke", MaxWorkers: 1, LeaseDuration: time.Minute})
	if err != nil {
		fail(err)
	}
	defer runtime.Close()
	if err := preparationworkflow.Register(runtime, func(context.Context, preparationworkflow.Identity) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: providers.Manifests, Schemas: providers.Schemas, Generator: providers.Generator, Embedder: providers.Embedder, Cache: providers.Cache, GenerationConcurrency: 1, GenerationSettingsFingerprint: providers.EngineOptions().GenerationSettingsFingerprint}, nil
	}); err != nil {
		fail(err)
	}
	input := preparationworkflow.Input{Identity: identity, Plan: plan}
	if *withEmbedding {
		input.Embedding = &preparationworkflow.EmbeddingSpec{RawRepresentationName: "raw", MaxRepresentationsPerChunk: 3, Node: ragcontract.Node{Config: []byte(`{"model":"embedding-primary","dimensions":768,"normalize":"l2","batchSize":16}`)}}
	}
	handle, err := runtime.EnsureRun(ctx, preparationworkflow.PackageName, input, scraperworkflow.WithRunID("rag-preparation-smoke-"+identityDigest[len("sha256:"):len("sha256:")+16]), scraperworkflow.WithRunIdentity(identity))
	if err != nil {
		fail(err)
	}
	for {
		snapshot, err := runtime.Snapshot(ctx, handle.ID)
		if err != nil {
			fail(err)
		}
		if snapshot.Stats.Succeeded == snapshot.Stats.Total {
			fmt.Printf("workflow=%s created=%t total=%d succeeded=%d provider=real\n", handle.ID, handle.Created, snapshot.Stats.Total, snapshot.Stats.Succeeded)
			return
		}
		if snapshot.Stats.Failed > 0 || snapshot.Stats.Canceled > 0 || snapshot.Stats.Blocked > 0 {
			fail(fmt.Errorf("workflow=%s total=%d succeeded=%d failed=%d blocked=%d canceled=%d", handle.ID, snapshot.Stats.Total, snapshot.Stats.Succeeded, snapshot.Stats.Failed, snapshot.Stats.Blocked, snapshot.Stats.Canceled))
		}
		if _, err := runtime.RunOnce(ctx); err != nil {
			fail(err)
		}
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "rag-preparation-smoke:", err)
	os.Exit(1)
}
