package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

func executeDurablePreparation(ctx context.Context, stateDB string, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, options ragengine.Options, identity ragengine.PreparedCorpusIdentity) error {
	if stateDB == "" || options.PreparedStore == nil {
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
	input := preparationworkflow.Input{Identity: workflowIdentity, Plan: plan, Embedding: &preparationworkflow.EmbeddingSpec{Node: mapping.EmbeddingNode, RawRepresentationName: mapping.RawRepresentationName, MaxRepresentationsPerChunk: mapping.MaxRepresentationsPerChunk}, Publication: &preparationworkflow.PublicationSpec{Identity: identity, RawOutputKey: mapping.RawOutputKey, DerivedOutputKey: mapping.DerivedOutputKey, MergedOutputKey: mapping.MergedOutputKey, EmbeddingOutputKey: mapping.EmbeddingOutputKey}}
	handle, err := runtime.EnsureRun(ctx, preparationworkflow.PackageName, input, scraperworkflow.WithRunID("rag-preparation-"+identityDigest[len("sha256:"):len("sha256:")+16]), scraperworkflow.WithRunIdentity(workflowIdentity))
	if err != nil {
		return err
	}
	for {
		snapshot, err := runtime.Snapshot(ctx, handle.ID)
		if err != nil {
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
