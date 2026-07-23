package ragworkflow

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestAuthorizedRealProviderWorkflowAcceptance(t *testing.T) {
	if os.Getenv("RAG_WORKFLOW_REAL_PROVIDER_ACCEPT") != "1" {
		t.Skip("explicit authorization required")
	}
	configPath := os.Getenv("RAG_PROVIDER_CONFIG")
	if configPath == "" {
		t.Fatal("RAG_PROVIDER_CONFIG is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	providers, err := ragproviders.Load(ctx, configPath)
	require.NoError(t, err)
	defer func() { require.NoError(t, providers.Close()) }()
	services, err := ProviderServicesFromSet(providers)
	require.NoError(t, err)
	providerPackage, err := NewProviderPackage(services, ragworkflowops.Policy{MaxPerAttempt: 100, FinishTimeout: 5 * time.Second})
	require.NoError(t, err)
	fixture, err := NewRealProviderSmokeFixture()
	require.NoError(t, err)
	fixture = bindTTCAcceptanceWorkload(t, fixture, os.Getenv("RAG_TTC_DATABASE"))
	lowerer, err := NewProviderLowerer(providerPackage)
	require.NoError(t, err)
	lowered, err := lowerer.Lower(ctx, fixture.Execution)
	require.NoError(t, err)
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
	engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: 2 * time.Second}
	inputs := stageWorkflowInputs(t, ctx, artifacts, fixture.Execution, fixture.Corpus, fixture.Dataset)
	preflight := &taskRuntime{context: workflowv3runtime.TaskModuleContext{Context: ctx, Request: workflowv3runtime.TaskRequest{NodeKey: "prepare-start", Inputs: inputs, Artifacts: artifacts}}, preparationIdentity: providerPackage.authority.Digest}
	_, err = preflight.prepare()
	require.NoError(t, err)
	require.NoError(t, engine.Submit(ctx, "real-provider-acceptance", lowered.Plan, inputs))
	dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{"cpu.rag.prepare": 1, "cpu.rag.query": 1, "cpu.rag.reduce": 1}, PollInterval: 5 * time.Millisecond}
	done := make(chan error, 1)
	go func() { done <- dispatcher.Run(ctx) }()
	var snapshot workflowv3.RunSnapshot
	for {
		snapshot, err = engine.Snapshot(ctx, "real-provider-acceptance")
		require.NoError(t, err)
		if snapshot.Status != "running" {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	operations, err := store.ExternalOperations(context.Background(), "real-provider-acceptance")
	require.NoError(t, err)
	require.Equal(t, "succeeded", snapshot.Status, "attempts=%#v operations=%#v", snapshot.Attempts, operations)
	require.Len(t, operations, 10)
	for _, attempt := range snapshot.Attempts {
		require.Equal(t, 1, attempt.Number)
	}
	require.False(t, providerPackage.authority.FixtureProviders)
	kinds := map[string]int{}
	for _, operation := range operations {
		require.NotNil(t, operation.Completion)
		require.Equal(t, workflowv3.ExternalOperationOutcomeSucceeded, operation.Completion.Outcome)
		require.Equal(t, providerPackage.authority.Digest, operation.AuthorityDigest)
		counters := map[string]int64{}
		for _, counter := range operation.Completion.Counters {
			counters[counter.Name] = counter.Units
		}
		require.Equal(t, int64(1), counters["requests"])
		switch operation.Kind.Name {
		case ragworkflowops.GenerateOperation:
			require.Positive(t, counters["input_tokens"])
			require.Positive(t, counters["output_tokens"])
			require.Positive(t, counters["cost_microunits"])
		case ragworkflowops.EmbedOperation:
			require.Positive(t, counters["output_items"])
		case ragworkflowops.RerankOperation:
			require.Positive(t, counters["input_tokens"])
			require.Positive(t, counters["output_items"])
		}
		kinds[operation.Kind.Name]++
	}
	require.Equal(t, 7, kinds[ragworkflowops.GenerateOperation])
	require.Equal(t, 2, kinds[ragworkflowops.EmbedOperation])
	require.Equal(t, 1, kinds[ragworkflowops.RerankOperation])
	body, err := workflowv3.ReadArtifact(context.Background(), artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	result, err := DecodeResult(body)
	require.NoError(t, err)
	require.Len(t, result.Results, 1)
	for _, item := range result.Results {
		require.NotEmpty(t, item.Answers)
		require.NotEmpty(t, item.Answers[0].CitationChunkIDs)
	}
}

func bindTTCAcceptanceWorkload(t *testing.T, fixture Fixture, databasePath string) Fixture {
	t.Helper()
	if databasePath == "" {
		t.Fatal("RAG_TTC_DATABASE is required")
	}
	database, err := sql.Open("sqlite3", databasePath)
	require.NoError(t, err)
	defer func() { require.NoError(t, database.Close()) }()
	var title, content string
	require.NoError(t, database.QueryRow(`SELECT title,content_text FROM documents WHERE content_text <> '' ORDER BY doc_id LIMIT 1`).Scan(&title, &content))
	runes := []rune(content)
	if len(runes) > 1200 {
		runes = runes[:1200]
	}
	require.NotEmpty(t, runes)
	fixture.Corpus = ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "ttc-source-1", SessionID: "ttc-wordpress", Ordinal: 1, Role: "user", Text: string(runes)}}}
	fixture.Dataset = ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{{ID: "ttc-q1", Text: "What does " + title + " describe?", RelevantIDs: []string{"ttc-source-1"}, Grades: map[string]float64{"ttc-source-1": 1}}}}
	corpusBody, err := ragcontract.CanonicalJSON(fixture.Corpus)
	require.NoError(t, err)
	corpusDigest, err := ragcontract.Digest(fixture.Corpus)
	require.NoError(t, err)
	corpusSize := int64(len(corpusBody) + 1)
	fixture.Execution.Bindings[0].Digest = corpusDigest
	fixture.Execution.Bindings[0].SizeBytes = &corpusSize
	fixture.Execution.Dataset.ManifestDigest, err = ragcontract.Digest(fixture.Dataset)
	require.NoError(t, err)
	fixture.Execution.CellID = ""
	fixture.Execution.CellID, err = ragcontract.Digest(fixture.Execution)
	require.NoError(t, err)
	_, err = ragengine.SerializePreparedValues(map[string]any{"corpus/out": fixture.Corpus})
	require.NoError(t, err)
	executionBody, err := ragcontract.CanonicalJSON(fixture.Execution)
	require.NoError(t, err)
	decoded, err := ragcontract.DecodeExecution(bytes.NewReader(executionBody))
	require.NoError(t, err)
	require.NoError(t, validateExecution(decoded))
	return fixture
}
