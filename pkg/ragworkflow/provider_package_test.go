package ragworkflow

import (
	"context"
	"errors"
	"path/filepath"
	"sync/atomic"
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
	runProviderPackageFixtureWithServices(t, variant, expectedGenerationOperations, expectedBudgetUnits, ProviderServices{}, 0)
}

func runProviderPackageFixtureWithServices(t *testing.T, variant string, expectedGenerationOperations int, expectedBudgetUnits int64, services ProviderServices, expectedFailedOperations int) {
	t.Helper()
	ctx := context.Background()
	baselineServices, err := NewDeterministicProviderServices()
	require.NoError(t, err)
	if services.Authority.SchemaVersion == "" {
		services = baselineServices
	}
	providerPackage, err := NewProviderPackage(services, ragworkflowops.Policy{MaxPerAttempt: 100, FinishTimeout: time.Second})
	require.NoError(t, err)
	providers := baselineServices
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
	require.Len(t, lowered.IR.Budgets, 1)
	require.Equal(t, "provider", lowered.IR.Budgets[0].Account)
	require.Equal(t, providerPackage.authority.Digest, lowered.IR.Budgets[0].PolicyDigest)
	accountLimits := budgetAmountsByDimension(lowered.IR.Budgets[0].Limits)
	require.Equal(t, expectedBudgetUnits, accountLimits["requests"])
	require.Equal(t, map[string]int64{"requests": 3}, budgetAmountsByDimension(lowered.IR.Maps[0].Budget.Reserve))
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
	deadline := time.Now().Add(60 * time.Second)
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
	if snapshot.Status != "succeeded" {
		for _, attempt := range snapshot.Attempts {
			if attempt.Failure != nil {
				t.Logf("failure %s/%d: %#v", attempt.NodeKey, attempt.Number, *attempt.Failure)
			}
		}
	}
	require.Equal(t, "succeeded", snapshot.Status, "attempts=%#v operations=%#v", snapshot.Attempts, operations)
	require.NotEmpty(t, operations)
	kinds := map[string]int{}
	failedOperations := 0
	for _, operation := range operations {
		require.NotNil(t, operation.Completion)
		if operation.Completion.Outcome == workflowv3.ExternalOperationOutcomeFailed {
			failedOperations++
		} else {
			require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, operation.Completion.Outcome)
		}
		require.Equal(t, providerPackage.authority.Digest, operation.AuthorityDigest)
		require.Equal(t, int64(1), operationCountersByName(operation.Reservation)["requests"])
		kinds[operation.Kind.Name]++
	}
	require.Equal(t, expectedFailedOperations, failedOperations)
	require.Equal(t, expectedGenerationOperations+expectedFailedOperations, kinds[ragworkflowops.GenerateOperation])
	require.Equal(t, 3, kinds[ragworkflowops.EmbedOperation])
	require.Equal(t, 2, kinds[ragworkflowops.RerankOperation])
	resultBody, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	workflowResult, err := DecodeResult(resultBody)
	require.NoError(t, err)
	require.Equal(t, directResult.Traces[0].Results, workflowResult.Results[0].Trace.Results)
	require.Equal(t, directResult.Answers[0], workflowResult.Results[0].Answers[0])
}

func budgetAmountsByDimension(amounts []workflowv3.BudgetAmount) map[string]int64 {
	ret := map[string]int64{}
	for _, amount := range amounts {
		ret[amount.Dimension] = amount.Units
	}
	return ret
}
func operationCountersByName(counters []workflowv3.ExternalOperationCounter) map[string]int64 {
	ret := map[string]int64{}
	for _, counter := range counters {
		ret[counter.Name] = counter.Units
	}
	return ret
}

type failOnceGenerator struct {
	inner ragoperators.TextGenerator
	calls atomic.Int64
}

func (g *failOnceGenerator) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	if g.calls.Add(1) == 1 {
		return ragoperators.GenerationResult{}, ragworkflowops.NewProviderCallError(errors.New("SECRET_PROVIDER_BODY_CANARY"), "transport", "PROVIDER_TRANSPORT", workflowv3.ExternalOperationOutcomeFailed)
	}
	return g.inner.Generate(ctx, request)
}

func TestProviderWorkflowCacheHitsDoNotCreateProviderOperations(t *testing.T) {
	services, err := NewDeterministicProviderServices()
	require.NoError(t, err)
	services.Cache = ragoperators.NewMemoryCache()
	runProviderPackageFixtureWithServices(t, "structured", 5, 30_200, services, 0)
	// Structured preparation is now cache-resident. Query answers remain live
	// contacts, while cache hits create no durable provider operation.
	runProviderPackageFixtureWithServices(t, "structured", 2, 30_200, services, 0)
}

func TestProviderWorkflowRetriesFailedContactAsDistinctOperation(t *testing.T) {
	services, err := NewDeterministicProviderServices()
	require.NoError(t, err)
	flaky := &failOnceGenerator{inner: services.Generator}
	services.Generator = flaky
	runProviderPackageFixtureWithServices(t, "structured", 5, 30_200, services, 1)
	require.Equal(t, int64(6), flaky.calls.Load())
}
