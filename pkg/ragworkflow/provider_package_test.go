package ragworkflow

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
	"github.com/stretchr/testify/require"
)

func fixtureProviderPackage(t *testing.T) (*ProviderPackage, ProviderServices) {
	t.Helper()
	services, err := NewDeterministicProviderServices()
	require.NoError(t, err)
	providerPackage, err := NewProviderPackage(services, ragworkflowops.Policy{MaxPerAttempt: 100, FinishTimeout: time.Second})
	require.NoError(t, err)
	return providerPackage, services
}

func TestProviderPackageBindsAuthorityAndExecutesDurableOperations(t *testing.T) {
	for _, test := range []struct {
		variant     string
		generate    int
		budgetUnits int64
	}{
		{variant: "structured", generate: 5, budgetUnits: 30_200},
		{variant: "combined", generate: 4, budgetUnits: 30_200},
		{variant: "synthetic", generate: 8, budgetUnits: 30_300},
	} {
		t.Run(test.variant, func(t *testing.T) { runProviderPackageFixture(t, test.variant, test.generate, test.budgetUnits) })
	}
}

func runProviderPackageFixture(t *testing.T, variant string, expectedGenerationOperations int, expectedBudgetUnits int64) {
	t.Helper()
	ctx := context.Background()
	providerPackage, providers := fixtureProviderPackage(t)
	fixture, err := NewDeterministicProviderFixtureVariant(variant)
	require.NoError(t, err)
	directEnvironment := &ragoperators.Environment{Manifests: providers.Manifests, Schemas: providers.Schemas, Generator: providers.Generator, Embedder: providers.Embedder, Reranker: providers.Reranker, GenerationConcurrency: 1, Usage: ragoperators.Usage{Cost: map[string]float64{}}}
	directResult, directErr := ragengine.New(nil).Execute(ctx, fixture.Execution, fixture.Corpus, fixture.Dataset, nil, ragengine.Options{Manifests: directEnvironment.Manifests, Schemas: directEnvironment.Schemas, Generator: directEnvironment.Generator, Embedder: directEnvironment.Embedder, Reranker: directEnvironment.Reranker, GenerationConcurrency: 1, GenerationSettingsFingerprint: providerPackage.authority.Digest, EmbeddingFingerprint: providerPackage.authority.Digest})
	require.NoError(t, directErr)
	_, err = NewLowerer().Lower(ctx, fixture.Execution)
	require.ErrorContains(t, err, "RAG_WORKFLOW_PROVIDER_REQUIRED")
	lowerer, err := NewProviderLowerer(providerPackage)
	require.NoError(t, err)
	lowered, err := lowerer.Lower(ctx, fixture.Execution)
	require.NoError(t, err)
	expectedFingerprint, err := preparationFingerprint(fixture.Execution, providerPackage.authority.Digest)
	require.NoError(t, err)
	require.Equal(t, expectedFingerprint, lowered.PreparationFingerprint)
	require.Equal(t, []workflowv3.BudgetAccount{{Account: "provider", Limits: []workflowv3.BudgetAmount{{Dimension: "requests", Units: expectedBudgetUnits}}, PolicyDigest: providerPackage.authority.Digest}}, lowered.IR.Budgets)
	require.Equal(t, []workflowv3.BudgetAmount{{Dimension: "requests", Units: 3}}, lowered.IR.Maps[0].Budget.Reserve)
	bundle, err := providerPackage.Bundle()
	require.NoError(t, err)
	builder := workflowv3.NewRegistryBuilder()
	require.NoError(t, builder.AdvertiseModules(ModuleAlias))
	require.NoError(t, builder.AddBundle(bundle))
	registry, err := builder.Seal()
	require.NoError(t, err)
	modules, err := workflowv3runtime.NewTaskModuleRegistry(providerPackage.TaskModuleFactories()...)
	require.NoError(t, err)
	root := t.TempDir()
	artifacts, err := workflowv3.NewFileArtifactStore(filepath.Join(root, "artifacts"), 32<<20)
	require.NoError(t, err)
	store, err := workflowv3sqlite.Open(ctx, filepath.Join(root, "workflow.db"))
	require.NoError(t, err)
	defer func() { _ = store.Close() }()
	engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: time.Second}
	inputs := stageWorkflowInputs(t, ctx, artifacts, fixture.Execution, fixture.Corpus, fixture.Dataset)
	require.NoError(t, engine.Submit(ctx, "provider-fixture", lowered.Plan, inputs))
	dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{"cpu.rag.prepare": 1, "cpu.rag.query": 1, "cpu.rag.reduce": 1}, PollInterval: time.Millisecond}
	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- dispatcher.Run(workerCtx) }()
	var snapshot workflowv3.RunSnapshot
	deadline := time.Now().Add(20 * time.Second)
	for {
		snapshot, err = engine.Snapshot(ctx, "provider-fixture")
		require.NoError(t, err)
		if snapshot.Status != "running" {
			break
		}
		require.True(t, time.Now().Before(deadline))
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	snapshot, err = engine.Snapshot(ctx, "provider-fixture")
	require.NoError(t, err)
	operations, operationErr := store.ExternalOperations(ctx, "provider-fixture")
	require.NoError(t, operationErr)
	require.Equal(t, "succeeded", snapshot.Status, "attempts=%#v operations=%#v", snapshot.Attempts, operations)
	require.NotEmpty(t, operations)
	kinds := map[string]int{}
	for _, operation := range operations {
		require.NotNil(t, operation.Completion)
		require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, operation.Completion.Outcome)
		require.Equal(t, providerPackage.authority.Digest, operation.AuthorityDigest)
		require.Equal(t, []workflowv3.ExternalOperationCounter{{Name: "requests", Units: 1}}, operation.Reservation)
		kinds[operation.Kind.Name]++
	}
	require.Equal(t, expectedGenerationOperations, kinds[ragworkflowops.GenerateOperation])
	require.Equal(t, 3, kinds[ragworkflowops.EmbedOperation])
	require.Equal(t, 2, kinds[ragworkflowops.RerankOperation])
	resultBody, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	workflowResult, err := DecodeResult(resultBody)
	require.NoError(t, err)
	require.Equal(t, directResult.Traces[0].Results, workflowResult.Results[0].Trace.Results)
	require.Equal(t, directResult.Answers[0], workflowResult.Results[0].Answers[0])
}
