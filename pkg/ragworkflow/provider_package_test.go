package ragworkflow

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
	"github.com/stretchr/testify/require"
)

type fixtureReranker struct{}

func (fixtureReranker) Rerank(_ context.Context, request ragoperators.RerankRequest) (ragoperators.RerankResult, error) {
	scores := make([]ragoperators.RerankScore, len(request.Candidates))
	for index, candidate := range request.Candidates {
		scores[index] = ragoperators.RerankScore{ChunkID: candidate.Chunk.Record.ID, Score: float64(len(scores) - index)}
	}
	return ragoperators.RerankResult{Scores: scores}, nil
}

func fixtureProviderPackage(t *testing.T) (*ProviderPackage, ragoperators.FixtureProviders) {
	t.Helper()
	providers := ragoperators.NewFixtureProviders()
	authority := ProviderAuthority{SchemaVersion: ProviderAuthoritySchema, ProfileID: "fixture-geppetto-v1", Capabilities: []string{"embedder", "generator", "reranker", "schema-validator"}, ModelManifestDigests: []string{"sha256:" + strings.Repeat("1", 64), "sha256:" + strings.Repeat("2", 64), "sha256:" + strings.Repeat("3", 64)}, PromptManifestDigests: []string{"sha256:" + strings.Repeat("4", 64), "sha256:" + strings.Repeat("5", 64)}, Providers: []WorkflowProviderIdentity{{Role: "embedding-primary", ProfileSlug: "fixture-embedding", ModelManifestDigest: "sha256:" + strings.Repeat("3", 64), ModelID: ragoperators.FixtureEmbeddingModel, SettingsFingerprint: "sha256:" + strings.Repeat("a", 64), ConcurrencyLimit: 1}, {Role: "generator-primary", ProfileSlug: "fixture-generation", ModelManifestDigest: "sha256:" + strings.Repeat("1", 64), ModelID: ragoperators.FixtureSummaryModel, SettingsFingerprint: "sha256:" + strings.Repeat("b", 64), ConcurrencyLimit: 1}, {Role: "reranker-primary", ProfileSlug: "fixture-rerank", ModelManifestDigest: "sha256:" + strings.Repeat("3", 64), ModelID: "fixture-rerank-v1", SettingsFingerprint: "sha256:" + strings.Repeat("c", 64), ConcurrencyLimit: 1}}}
	canonicalizeAuthority(&authority)
	digest, err := providerAuthorityDigest(authority)
	require.NoError(t, err)
	authority.Digest = digest
	providerPackage, err := NewProviderPackage(ProviderServices{Authority: authority, Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Reranker: fixtureReranker{}, GenerationConcurrency: 1}, ragworkflowops.Policy{MaxPerAttempt: 100, FinishTimeout: time.Second})
	require.NoError(t, err)
	return providerPackage, providers
}

func providerExecution(t *testing.T) Fixture {
	t.Helper()
	fixture, err := NewProviderFreeFixture(true)
	require.NoError(t, err)
	for index, node := range fixture.Execution.Pipeline.Nodes {
		switch node.Operator.Kind {
		case "representations.raw":
			fixture.Execution.Pipeline.Nodes[index].Operator = ragcontract.OperatorRef{Kind: "representations.structured-summary", Version: "v1"}
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"name":"summary","model":"fixture-summary-v1","prompt":"fixture-transcript-summary-v1","outputSchema":"transcript-rag-summary/v1"}`)
		case "embed.model":
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"model":"fixture-hash-32-v1","dimensions":32,"distance":"cosine","normalize":"l2","batchSize":8}`)
		}
	}
	fixture.Execution.Pipeline, err = ragcompiler.Normalize(fixture.Execution.Pipeline, nil)
	require.NoError(t, err)
	fixture.Execution.CellID = ""
	fixture.Execution.CellID, err = ragcontract.Digest(fixture.Execution)
	require.NoError(t, err)
	return fixture
}

func TestProviderPackageBindsAuthorityAndExecutesDurableOperations(t *testing.T) {
	ctx := context.Background()
	providerPackage, providers := fixtureProviderPackage(t)
	fixture := providerExecution(t)
	_, err := NewLowerer().Lower(ctx, fixture.Execution)
	require.ErrorContains(t, err, "RAG_WORKFLOW_PROVIDER_REQUIRED")
	lowerer, err := NewProviderLowerer(providerPackage)
	require.NoError(t, err)
	lowered, err := lowerer.Lower(ctx, fixture.Execution)
	require.NoError(t, err)
	expectedFingerprint, err := preparationFingerprint(fixture.Execution, providerPackage.authority.Digest)
	require.NoError(t, err)
	require.Equal(t, expectedFingerprint, lowered.PreparationFingerprint)
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
	require.Equal(t, "succeeded", snapshot.Status)
	operations, err := store.ExternalOperations(ctx, "provider-fixture")
	require.NoError(t, err)
	require.NotEmpty(t, operations)
	kinds := map[string]int{}
	for _, operation := range operations {
		require.NotNil(t, operation.Completion)
		require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, operation.Completion.Outcome)
		require.Equal(t, providerPackage.authority.Digest, operation.AuthorityDigest)
		kinds[operation.Kind.Name]++
	}
	require.Positive(t, kinds[ragworkflowops.GenerateOperation])
	require.Positive(t, kinds[ragworkflowops.EmbedOperation])
	require.Zero(t, kinds[ragworkflowops.RerankOperation])
	resultBody, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	workflowResult, err := DecodeResult(resultBody)
	require.NoError(t, err)
	environment := &ragoperators.Environment{Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Reranker: fixtureReranker{}, GenerationConcurrency: 1, Usage: ragoperators.Usage{Cost: map[string]float64{}}}
	baseline, err := ragengine.New(nil).Execute(ctx, fixture.Execution, fixture.Corpus, fixture.Dataset, nil, ragengine.Options{Manifests: environment.Manifests, Schemas: environment.Schemas, Generator: environment.Generator, Embedder: environment.Embedder, Reranker: environment.Reranker, GenerationConcurrency: 1, GenerationSettingsFingerprint: providerPackage.authority.Digest, EmbeddingFingerprint: providerPackage.authority.Digest})
	require.NoError(t, err)
	require.Equal(t, baseline.Traces[0].Results, workflowResult.Results[0].Trace.Results)
}
