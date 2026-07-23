package ragworkflow

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
	"github.com/stretchr/testify/require"
)

func TestWorkflowExecutionMatchesRAGEngineAndSurvivesRestart(t *testing.T) {
	ctx := context.Background()
	fixture, err := NewProviderFreeFixture(true)
	require.NoError(t, err)
	execution, corpus, dataset := fixture.Execution, fixture.Corpus, fixture.Dataset
	lowered, err := NewLowerer().Lower(ctx, execution)
	require.NoError(t, err)
	root := t.TempDir()
	artifacts, err := workflowv3.NewFileArtifactStore(filepath.Join(root, "artifacts"), 16<<20)
	require.NoError(t, err)
	bundle, err := Bundle()
	require.NoError(t, err)
	builder := workflowv3.NewRegistryBuilder()
	require.NoError(t, builder.AdvertiseModules(ModuleAlias))
	require.NoError(t, builder.AddBundle(bundle))
	registry, err := builder.Seal()
	require.NoError(t, err)
	modules, err := workflowv3runtime.NewTaskModuleRegistry(TaskModuleFactory())
	require.NoError(t, err)
	storePath := filepath.Join(root, "workflow.db")
	store, err := workflowv3sqlite.Open(ctx, storePath)
	require.NoError(t, err)
	engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: 200 * time.Millisecond}
	inputs := stageWorkflowInputs(t, ctx, artifacts, execution, corpus, dataset)
	require.NoError(t, engine.Submit(ctx, "rag-parity", lowered.Plan, inputs))
	preparationTaskCount := len(staticNodes(execution.Pipeline)) + 1
	for index := 0; index < preparationTaskCount; index++ {
		ran, runErr := engine.RunOne(ctx)
		require.NoError(t, runErr)
		require.True(t, ran)
	}
	beforeRestart, err := engine.Snapshot(ctx, "rag-parity")
	require.NoError(t, err)
	require.Equal(t, "running", beforeRestart.Status)
	require.Len(t, beforeRestart.Attempts, preparationTaskCount)
	for _, attempt := range beforeRestart.Attempts {
		require.True(t, attempt.NodeKey == "prepare-start" || strings.HasPrefix(string(attempt.NodeKey), "prepare-"), attempt.NodeKey)
		require.Equal(t, "succeeded", attempt.Status)
	}
	require.NoError(t, store.Close())

	store, err = workflowv3sqlite.Open(ctx, storePath)
	require.NoError(t, err)
	defer func() { _ = store.Close() }()
	engine = &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: 200 * time.Millisecond}
	dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{"cpu.rag.prepare": 1, "cpu.rag.query": 2, "cpu.rag.reduce": 1}, PollInterval: time.Millisecond}
	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() { done <- dispatcher.Run(workerCtx) }()
	var snapshot workflowv3.RunSnapshot
	deadline := time.Now().Add(20 * time.Second)
	for {
		snapshot, err = engine.Snapshot(ctx, "rag-parity")
		require.NoError(t, err)
		if snapshot.Status != "running" {
			break
		}
		require.True(t, time.Now().Before(deadline), "workflow did not terminate")
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	require.Equal(t, "succeeded", snapshot.Status)
	resultBody, err := workflowv3.ReadArtifact(ctx, artifacts, snapshot.Outputs["result"])
	require.NoError(t, err)
	workflowResult, err := DecodeResult(resultBody)
	require.NoError(t, err)
	require.Equal(t, ResultSchema, workflowResult.SchemaVersion)
	require.Len(t, workflowResult.Results, 2)
	require.Equal(t, lowered.PreparationFingerprint, workflowResult.PreparationFingerprint)
	projection, err := (DomainProjector{}).Project(ctx, researchrunner.DomainProjectionInput{Outputs: map[string]researchrunner.DomainOutput{"result": {Name: "result", SchemaVersion: ResultSchema, Data: resultBody}}})
	require.NoError(t, err)
	require.Len(t, projection.Metrics, 2)
	require.Len(t, projection.Traces, 2)
	require.Equal(t, "rag.mrr", projection.Metrics[0].Name)
	require.Contains(t, projection.Metrics[0].Scope, "rag.query.")

	canary := "SECRET-PROJECTION-CANARY"
	privateResult := workflowResult
	privateResult.Results = append([]QueryResult(nil), workflowResult.Results...)
	privateResult.Results[0].Trace.Query.Metadata = json.RawMessage(`{"secret":"` + canary + `"}`)
	privateResult.Results[0].Trace.Failures = []ragcontract.FailureTrace{{Code: "PRIVATE", Path: "query", Message: canary, Details: json.RawMessage(`{"secret":"` + canary + `"}`)}}
	privateResult.Results[0].Metrics = append([]Metric(nil), workflowResult.Results[0].Metrics...)
	privateResult.Results[0].Metrics[0].Metadata = json.RawMessage(`{"secret":"` + canary + `"}`)
	privateResult.Digest, err = resultDigest(privateResult)
	require.NoError(t, err)
	privateBody, err := ragcontract.CanonicalJSON(privateResult)
	require.NoError(t, err)
	privateProjection, err := (DomainProjector{}).Project(ctx, researchrunner.DomainProjectionInput{Outputs: map[string]researchrunner.DomainOutput{"result": {Name: "result", SchemaVersion: ResultSchema, Data: privateBody}}})
	require.NoError(t, err)
	projectedBody, err := json.Marshal(privateProjection)
	require.NoError(t, err)
	require.NotContains(t, string(projectedBody), canary)
	require.Contains(t, string(projectedBody), "measureMetadataDigest")

	environment, err := providerFreeEnvironment(execution)
	require.NoError(t, err)
	for index := range workflowResult.Results {
		baseline, executeErr := ragengine.New(nil).Execute(ctx, execution, corpus, ragoperators.EvaluationDataset{SchemaVersion: dataset.SchemaVersion, Queries: []ragoperators.Query{dataset.Queries[index]}}, nil, ragengine.Options{Manifests: environment.Manifests, Embedder: environment.Embedder, EmbeddingFingerprint: "fixture-embedding/v1"})
		require.NoError(t, executeErr)
		require.Len(t, baseline.Traces, 1)
		require.Equal(t, baseline.Traces[0].Results, workflowResult.Results[index].Trace.Results)
		require.Equal(t, semanticMetrics(workflowMetrics(baseline.Metrics)), semanticMetrics(workflowResult.Results[index].Metrics))
	}
	firstDigest := workflowResult.Digest
	decoded, err := DecodeResult(resultBody)
	require.NoError(t, err)
	require.Equal(t, firstDigest, decoded.Digest)
	tampered := decoded
	tampered.VariantID = "changed"
	require.ErrorContains(t, ValidateResult(tampered), "RESULT_INVALID")
	unknown := append([]byte(nil), resultBody[:len(resultBody)-1]...)
	unknown = append(unknown, []byte(`,"unknown":true}`)...)
	_, err = DecodeResult(unknown)
	require.ErrorContains(t, err, "unknown field")
}

func stageWorkflowInputs(t *testing.T, ctx context.Context, artifacts workflowv3.ArtifactStore, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, dataset ragoperators.EvaluationDataset) map[string]workflowv3.ArtifactRef {
	t.Helper()
	executionBody, err := ragcontract.CanonicalJSON(execution)
	require.NoError(t, err)
	executionRef, err := artifacts.Put(ctx, ragcontract.ExecutionSchemaVersion, "application/json", executionBody)
	require.NoError(t, err)
	corpusBody, err := ragcontract.CanonicalJSON(corpus)
	require.NoError(t, err)
	corpusRef, err := artifacts.Put(ctx, CorpusSchema, "application/json", corpusBody)
	require.NoError(t, err)
	items := make([]workflowv3.ManifestItem, len(dataset.Queries))
	for index, query := range dataset.Queries {
		body, marshalErr := ragcontract.CanonicalJSON(QueryItem{SchemaVersion: QuerySchema, DatasetManifestDigest: execution.Dataset.ManifestDigest, Query: query})
		require.NoError(t, marshalErr)
		ref, putErr := artifacts.Put(ctx, QuerySchema, "application/json", body)
		require.NoError(t, putErr)
		items[index] = workflowv3.ManifestItem{Key: query.ID, Value: ref}
	}
	manifest, err := workflowv3.NewItemManifest(QuerySchema, items)
	require.NoError(t, err)
	manifestBody, err := workflowv3.EncodeItemManifest(manifest)
	require.NoError(t, err)
	manifestRef, err := artifacts.Put(ctx, workflowv3.ItemManifestSchemaV1, "application/json", manifestBody)
	require.NoError(t, err)
	return map[string]workflowv3.ArtifactRef{"execution": executionRef, "corpus": corpusRef, "queries": manifestRef}
}

func semanticMetrics(metrics []Metric) map[string]json.RawMessage {
	ret := map[string]json.RawMessage{}
	for _, metric := range metrics {
		ret[metric.Name+"\x00"+string(metric.Metadata)] = metric.Value
	}
	return ret
}

func TestPreparedFingerprintRejectsCorpusAndPipelineDrift(t *testing.T) {
	execution := fixtureExecution(t, false)
	first, err := preparationFingerprint(execution)
	require.NoError(t, err)
	execution.Bindings[0].Digest = "sha256:" + strings.Repeat("f", 64)
	second, err := preparationFingerprint(execution)
	require.NoError(t, err)
	require.NotEqual(t, first, second)
	values := []string{first, second}
	sort.Strings(values)
	require.NotEqual(t, values[0], values[1])
}
