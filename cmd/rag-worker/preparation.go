package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/go-go-golems/researchctl/pkg/lab"
	"github.com/go-go-golems/researchctl/pkg/labprogress"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	storecontract "github.com/go-go-golems/scraper/pkg/engine/store"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

func executeDurablePreparation(ctx context.Context, encoder *json.Encoder, stateDB string, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, options ragengine.Options, identity ragengine.PreparedCorpusIdentity) error {
	if stateDB == "" || options.PreparedStore == nil || encoder == nil {
		return fmt.Errorf("RAG_WORKER_PREPARATION_CONFIG")
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(execution.Pipeline)
	if err != nil {
		return err
	}
	engine := ragengine.New(nil)
	inputs, err := engine.StaticInputs(ctx, execution.Pipeline, corpus, options, mapping.CombinedNode.ID)
	if err != nil {
		return err
	}
	chunks, ok := inputs["chunks"].([]ragoperators.Chunk)
	if !ok || len(chunks) == 0 {
		return fmt.Errorf("RAG_WORKER_PREPARATION_CHUNKS")
	}
	plan, err := ragoperators.PlanCombinedPreparation(chunks, mapping.CombinedNode)
	if err != nil {
		return err
	}
	identityDigest, err := ragcontract.Digest(identity)
	if err != nil {
		return err
	}
	workflowIdentity := preparationworkflow.Identity{SchemaVersion: "rag-preparation-workflow/v1", PreparedDigest: identityDigest}
	runtime, err := scraperworkflow.NewRuntime(ctx, scraperworkflow.Config{Store: scraperworkflow.SQLiteStore(stateDB), WorkerID: "rag-worker", MaxWorkers: 1, LeaseDuration: time.Minute})
	if err != nil {
		return err
	}
	defer runtime.Close()
	resolve := func(context.Context, preparationworkflow.Identity) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: options.Manifests, Schemas: options.Schemas, Generator: options.Generator, Embedder: options.Embedder, Cache: options.Cache, GenerationConcurrency: options.GenerationConcurrency, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint}, nil
	}
	publish := func(context.Context, preparationworkflow.Identity, preparationworkflow.PublicationSpec) (preparationworkflow.PublicationTarget, error) {
		return preparationworkflow.PublicationTarget{Store: options.PreparedStore, Engine: engine, Pipeline: execution.Pipeline, Corpus: corpus, Options: options}, nil
	}
	if err := preparationworkflow.RegisterWithPublication(runtime, resolve, publish); err != nil {
		return err
	}
	input := preparationworkflow.Input{Identity: workflowIdentity, Plan: plan, Embedding: &preparationworkflow.EmbeddingSpec{Node: mapping.EmbeddingNode, RawRepresentationName: mapping.RawRepresentationName, MaxRepresentationsPerChunk: mapping.MaxRepresentationsPerChunk}, Publication: &preparationworkflow.PublicationSpec{Identity: identity, ChunksOutputKey: mapping.ChunksOutputKey, RawOutputKey: mapping.RawOutputKey, DerivedOutputKey: mapping.DerivedOutputKey, MergedOutputKey: mapping.MergedOutputKey, EmbeddingOutputKey: mapping.EmbeddingOutputKey}}
	handle, err := runtime.EnsureRun(ctx, preparationworkflow.PackageName, input, scraperworkflow.WithRunID("rag-preparation-"+identityDigest[len("sha256:"):len("sha256:")+16]), scraperworkflow.WithRunIdentity(workflowIdentity))
	if err != nil {
		return err
	}
	emitter, err := researchctladapter.NewPreparationProgressEmitter(func(event lab.EventInput) error {
		return encoder.Encode(frame{Type: "event", Event: event})
	})
	if err != nil {
		return err
	}
	for {
		snapshot, err := runtime.Snapshot(ctx, handle.ID)
		if err != nil {
			return err
		}
		if err := emitter.Emit(preparationProgress(snapshot, handle.IdentityDigest)); err != nil {
			return err
		}
		switch snapshot.Workflow.Status {
		case model.WorkflowStatusSucceeded:
			return nil
		case model.WorkflowStatusFailed, model.WorkflowStatusCanceled:
			return fmt.Errorf("RAG_WORKER_PREPARATION_TERMINAL: %s", snapshot.Workflow.Status)
		}
		if _, err := runtime.RunOnce(ctx); err != nil {
			return err
		}
	}
}

func preparationProgress(snapshot *storecontract.WorkflowSnapshot, identityDigest string) labprogress.Envelope {
	stats := snapshot.Stats
	return labprogress.Envelope{SchemaVersion: labprogress.SchemaVersionV1, Type: "rag.preparation.workflow.progress/v1", OccurredAt: time.Now().UTC(), WorkflowID: string(snapshot.Workflow.ID), WorkflowIdentityDigest: identityDigest, Phase: string(snapshot.Workflow.Status), Counts: labprogress.Counts{Pending: int64(stats.Pending), Ready: int64(stats.Ready), Running: int64(stats.Running), Succeeded: int64(stats.Succeeded), Failed: int64(stats.Failed), Blocked: int64(stats.Blocked), Canceled: int64(stats.Canceled), Total: int64(stats.Total)}}
}
