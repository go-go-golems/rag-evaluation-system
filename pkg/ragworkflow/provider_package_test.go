package ragworkflow

import (
	"context"
	"encoding/json"
	"fmt"
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

type completeFixtureProviders struct{ ragoperators.FixtureProviders }

func (p completeFixtureProviders) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	if request.Kind != "generate.answer" {
		return p.FixtureProviders.Generate(ctx, request)
	}
	if len(request.Evidence) == 0 {
		return ragoperators.GenerationResult{}, fmt.Errorf("fixture answer evidence required")
	}
	cost := 0.000004
	return ragoperators.GenerationResult{Text: "fixture grounded answer", CitationChunkIDs: []string{request.Evidence[0].Chunk.Record.ID}, InputTokens: 11, OutputTokens: 4, Cost: &cost, FinishReason: "fixture"}, nil
}
func (completeFixtureProviders) Rerank(_ context.Context, request ragoperators.RerankRequest) (ragoperators.RerankResult, error) {
	scores := make([]ragoperators.RerankScore, len(request.Candidates))
	for index, candidate := range request.Candidates {
		scores[index] = ragoperators.RerankScore{ChunkID: candidate.Chunk.Record.ID, Score: float64(len(scores) - index)}
	}
	cost := 0.000002
	return ragoperators.RerankResult{Scores: scores, InputTokens: 7, Cost: &cost}, nil
}

func fixtureProviderPackage(t *testing.T) (*ProviderPackage, completeFixtureProviders) {
	t.Helper()
	providers := completeFixtureProviders{FixtureProviders: ragoperators.NewFixtureProviders()}
	providers.Resolver.Models["fixture-rerank-v1"] = ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + strings.Repeat("6", 64)}, ModelID: "fixture-rerank-v1", ModelDigest: "sha256:" + strings.Repeat("6", 64), Tokenization: "fixture-utf16", Truncation: "none", Normalization: "none", ImplementationVersion: "fixture/v1", RequestParameters: json.RawMessage(`{}`)}
	providers.Resolver.Prompts["fixture-answer-v1"] = ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + strings.Repeat("7", 64)}, PromptID: "fixture-answer-v1", TemplateDigest: "sha256:" + strings.Repeat("7", 64), InputSchema: "text/plain", OutputSchema: "fixture-answer/v1"}
	authority := ProviderAuthority{SchemaVersion: ProviderAuthoritySchema, ProfileID: "fixture-geppetto-v1", Capabilities: []string{"embedder", "generator", "reranker", "schema-validator"}, ModelManifestDigests: []string{"sha256:" + strings.Repeat("1", 64), "sha256:" + strings.Repeat("2", 64), "sha256:" + strings.Repeat("3", 64), "sha256:" + strings.Repeat("6", 64)}, PromptManifestDigests: []string{"sha256:" + strings.Repeat("4", 64), "sha256:" + strings.Repeat("5", 64), "sha256:" + strings.Repeat("7", 64)}, Providers: []WorkflowProviderIdentity{{Role: "embedding-primary", ProfileSlug: "fixture-embedding", ModelManifestDigest: "sha256:" + strings.Repeat("3", 64), ModelID: ragoperators.FixtureEmbeddingModel, SettingsFingerprint: "sha256:" + strings.Repeat("a", 64), ConcurrencyLimit: 1}, {Role: "generator-primary", ProfileSlug: "fixture-generation", ModelManifestDigest: "sha256:" + strings.Repeat("1", 64), ModelID: ragoperators.FixtureSummaryModel, SettingsFingerprint: "sha256:" + strings.Repeat("b", 64), ConcurrencyLimit: 1}, {Role: "reranker-primary", ProfileSlug: "fixture-rerank", ModelManifestDigest: "sha256:" + strings.Repeat("6", 64), ModelID: "fixture-rerank-v1", SettingsFingerprint: "sha256:" + strings.Repeat("c", 64), ConcurrencyLimit: 1}}}
	canonicalizeAuthority(&authority)
	digest, err := providerAuthorityDigest(authority)
	require.NoError(t, err)
	authority.Digest = digest
	providerPackage, err := NewProviderPackage(ProviderServices{Authority: authority, Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Reranker: providers, GenerationConcurrency: 1}, ragworkflowops.Policy{MaxPerAttempt: 100, FinishTimeout: time.Second})
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
		case "retrieve.bm25", "retrieve.vector":
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"index":"representations","representation":"summary","topK":10,"filter":{}}`)
		}
	}
	hydrationID := ""
	for _, node := range fixture.Execution.Pipeline.Nodes {
		if node.Operator.Kind == "hydrate.source-evidence" {
			hydrationID = node.ID
		}
	}
	require.NotEmpty(t, hydrationID)
	fixture.Execution.Pipeline.Nodes = append(fixture.Execution.Pipeline.Nodes,
		ragcontract.Node{ID: "fixture-rerank", Operator: ragcontract.OperatorRef{Kind: "rerank.cross-encoder", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "evidence", From: ragcontract.PortRef{NodeID: hydrationID, Port: "evidence"}}}, Config: json.RawMessage(`{"model":"fixture-rerank-v1","candidateCount":20,"results":5,"inputTemplate":"query-document","truncation":"none","tokenization":"fixture-utf16","timeoutMilliseconds":5000}`)},
		ragcontract.Node{ID: "fixture-answer", Operator: ragcontract.OperatorRef{Kind: "generate.answer", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "evidence", From: ragcontract.PortRef{NodeID: "fixture-rerank", Port: "evidence"}}}, Config: json.RawMessage(`{"model":"fixture-summary-v1","prompt":"fixture-answer-v1","citations":"required","citationFailurePolicy":"abstain","contextBudgetTokens":2048}`)},
	)
	fixture.Execution.Pipeline.Outputs[0].From = ragcontract.PortRef{NodeID: "fixture-rerank", Port: "evidence"}
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
	directEnvironment := &ragoperators.Environment{Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Reranker: providers, GenerationConcurrency: 1, Usage: ragoperators.Usage{Cost: map[string]float64{}}}
	directResult, directErr := ragengine.New(nil).Execute(ctx, fixture.Execution, fixture.Corpus, fixture.Dataset, nil, ragengine.Options{Manifests: directEnvironment.Manifests, Schemas: directEnvironment.Schemas, Generator: directEnvironment.Generator, Embedder: directEnvironment.Embedder, Reranker: directEnvironment.Reranker, GenerationConcurrency: 1, GenerationSettingsFingerprint: providerPackage.authority.Digest, EmbeddingFingerprint: providerPackage.authority.Digest})
	require.NoError(t, directErr)
	_, err := NewLowerer().Lower(ctx, fixture.Execution)
	require.ErrorContains(t, err, "RAG_WORKFLOW_PROVIDER_REQUIRED")
	lowerer, err := NewProviderLowerer(providerPackage)
	require.NoError(t, err)
	lowered, err := lowerer.Lower(ctx, fixture.Execution)
	require.NoError(t, err)
	expectedFingerprint, err := preparationFingerprint(fixture.Execution, providerPackage.authority.Digest)
	require.NoError(t, err)
	require.Equal(t, expectedFingerprint, lowered.PreparationFingerprint)
	require.Equal(t, []workflowv3.BudgetAccount{{Account: "provider", Limits: []workflowv3.BudgetAmount{{Dimension: "requests", Units: 30_200}}, PolicyDigest: providerPackage.authority.Digest}}, lowered.IR.Budgets)
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
	require.Equal(t, 5, kinds[ragworkflowops.GenerateOperation])
	require.Equal(t, 3, kinds[ragworkflowops.EmbedOperation])
	require.Equal(t, 2, kinds[ragworkflowops.RerankOperation])
	resultBody, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	workflowResult, err := DecodeResult(resultBody)
	require.NoError(t, err)
	require.Equal(t, directResult.Traces[0].Results, workflowResult.Results[0].Trace.Results)
	require.Equal(t, directResult.Answers[0], workflowResult.Results[0].Answers[0])
}
