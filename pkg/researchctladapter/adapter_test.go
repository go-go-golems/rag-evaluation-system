package researchctladapter

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/researchctl/pkg/lab"
	"github.com/go-go-golems/researchctl/pkg/lab/processrunner"
)

type sink struct {
	events    []lab.EventInput
	artifacts []lab.ArtifactInput
	metrics   []lab.MetricInput
	traces    []lab.TraceInput
	complete  *lab.AttemptSummary
	cancel    context.CancelFunc
}

func (s *sink) AppendEvent(_ context.Context, v lab.EventInput) (lab.EventRecord, error) {
	s.events = append(s.events, v)
	return lab.EventRecord{}, nil
}
func (s *sink) RecordArtifact(_ context.Context, v lab.ArtifactInput) (lab.RunArtifact, error) {
	s.artifacts = append(s.artifacts, v)
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	return lab.RunArtifact{}, nil
}
func (s *sink) RecordMetric(_ context.Context, v lab.MetricInput) (lab.MetricRecord, error) {
	s.metrics = append(s.metrics, v)
	return lab.MetricRecord{}, nil
}
func (s *sink) RecordTrace(_ context.Context, v lab.TraceInput) (lab.TraceRecord, error) {
	s.traces = append(s.traces, v)
	return lab.TraceRecord{}, nil
}
func (s *sink) CompleteAttempt(_ context.Context, v lab.AttemptSummary) error {
	s.complete = &v
	return nil
}

func TestWrapExecuteAndReconstructCanonicalSpecification(t *testing.T) {
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{Records: []ragoperators.SourceRecord{{ID: "source", SessionID: "s", Ordinal: 1, Role: "document", Text: "weighted reciprocal rank fusion"}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q", Text: "rank fusion"}}}, "fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest)
	corpusInput, err := StageEnvelope(InputReference{Role: "corpus", Kind: "manifest-envelope"}, corpus, root)
	if err != nil {
		t.Fatal(err)
	}
	datasetInput, err := StageEnvelope(InputReference{Role: "evaluation-dataset", Kind: "manifest-envelope"}, dataset, root)
	if err != nil {
		t.Fatal(err)
	}
	resolved := ResolvedInputs{ByRole: map[string]ResolvedInput{"corpus": corpusInput, "evaluation-dataset": datasetInput}}
	corpusPath := filepath.Join(root, filepath.FromSlash(corpusInput.Reference.URI))
	datasetPath := filepath.Join(root, filepath.FromSlash(datasetInput.Reference.URI))
	corpusBefore, err := os.ReadFile(corpusPath)
	if err != nil {
		t.Fatal(err)
	}
	datasetBefore, err := os.ReadFile(datasetPath)
	if err != nil {
		t.Fatal(err)
	}
	execution := fixtureExecution(t, corpus.Manifest.Digest, dataset.Manifest.Digest)
	specification, err := WrapExecution(execution, resolved, "fixture")
	if err != nil {
		t.Fatal(err)
	}
	binary := buildWorker(t)
	if err := CheckWorker(context.Background(), WorkerCommand{Executable: binary}); err != nil {
		t.Fatal(err)
	}
	runner, err := processrunner.New(processrunner.Config{Command: []string{binary}, Runner: lab.RunnerRecord{Name: RunnerName, ResolvedVersion: RunnerVersion}, Domains: []lab.DomainVersion{{Domain: ragcontract.Domain, SchemaVersion: ragcontract.DomainSchemaVersion}}})
	if err != nil {
		t.Fatal(err)
	}
	collector := &sink{}
	request := lab.AttemptRequest{Specification: specification, Run: lab.RunRecord{ID: "run", CreatedAt: "2026-07-17T00:00:00Z"}, AttemptID: "attempt", AttemptIndex: 1, ArtifactRoot: root}
	if err := runner.Start(context.Background(), request, collector); err != nil {
		t.Fatal(err)
	}
	if collector.complete == nil || len(collector.metrics) != 1 || len(collector.traces) != 1 || len(collector.artifacts) < 4 {
		t.Fatalf("sink=%#v", collector)
	}
	corpusAfter, err := os.ReadFile(corpusPath)
	if err != nil {
		t.Fatal(err)
	}
	datasetAfter, err := os.ReadFile(datasetPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(corpusBefore, corpusAfter) || !bytes.Equal(datasetBefore, datasetAfter) {
		t.Fatal("read-only worker mutated staged input bytes")
	}
	if collector.traces[0].Kind != ragcontract.TraceSchemaVersion {
		t.Fatalf("trace=%#v", collector.traces[0])
	}
	var artifactBytes int64
	for _, artifact := range collector.artifacts {
		info, statErr := os.Stat(filepath.Join(root, filepath.FromSlash(artifact.URI)))
		if statErr != nil {
			t.Fatal(statErr)
		}
		artifactBytes += info.Size()
	}
	var trace ragcontract.QueryTrace
	if err := json.Unmarshal(collector.traces[0].Value, &trace); err != nil {
		t.Fatal(err)
	}
	var metricBytes int
	for _, metric := range collector.metrics {
		metricBytes += len(metric.Value) + len(metric.Metadata)
	}
	t.Logf("study_fixture artifacts=%d artifact_bytes=%d trace_bytes=%d metric_bytes=%d cost=%v", len(collector.artifacts), artifactBytes, len(collector.traces[0].Value), metricBytes, trace.Usage.ProviderCost)
	export := lab.RunExport{Specification: specification}
	reconstructed, err := ReconstructSpecification(export)
	if err != nil {
		t.Fatal(err)
	}
	left, _ := lab.CanonicalJSON(specification)
	right, _ := lab.CanonicalJSON(reconstructed)
	if !bytes.Equal(left, right) {
		t.Fatalf("reconstruction mismatch\n%s\n%s", left, right)
	}
}
func TestExternalWorkerCancellationPreservesEarlierArtifact(t *testing.T) {
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{Records: []ragoperators.SourceRecord{{ID: "s", Text: "x"}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q", Text: "x"}}}, "fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest)
	corpusInput, _ := StageEnvelope(InputReference{Role: "corpus"}, corpus, root)
	datasetInput, _ := StageEnvelope(InputReference{Role: "evaluation-dataset"}, dataset, root)
	resolved := ResolvedInputs{ByRole: map[string]ResolvedInput{"corpus": corpusInput, "evaluation-dataset": datasetInput}}
	specification, _ := WrapExecution(fixtureExecution(t, corpus.Manifest.Digest, dataset.Manifest.Digest), resolved, "cancel")
	binary := buildWorker(t)
	runner, _ := processrunner.New(processrunner.Config{Command: []string{binary}, Runner: lab.RunnerRecord{Name: RunnerName, ResolvedVersion: RunnerVersion}, Domains: []lab.DomainVersion{{Domain: ragcontract.Domain, SchemaVersion: ragcontract.DomainSchemaVersion}}})
	ctx, cancel := context.WithCancel(context.Background())
	collector := &sink{cancel: cancel}
	err := runner.Start(ctx, lab.AttemptRequest{Specification: specification, Run: lab.RunRecord{ID: "r", CreatedAt: "2026-07-17T00:00:00Z"}, AttemptID: "a", AttemptIndex: 1, ArtifactRoot: root}, collector)
	if err == nil || !strings.Contains(err.Error(), "cancel") {
		t.Fatalf("err=%v", err)
	}
	if len(collector.artifacts) == 0 {
		t.Fatal("partial artifact was not preserved")
	}
}
func TestWorkerRejectsManifestLineageBeforeExecution(t *testing.T) {
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{Records: []ragoperators.SourceRecord{{ID: "s", Text: "x"}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q", Text: "x"}}}, "fixture", "smoke", "candidate", "unit", "sha256:"+strings.Repeat("f", 64))
	corpusInput, _ := StageEnvelope(InputReference{Role: "corpus"}, corpus, root)
	datasetInput, _ := StageEnvelope(InputReference{Role: "evaluation-dataset"}, dataset, root)
	resolved := ResolvedInputs{ByRole: map[string]ResolvedInput{"corpus": corpusInput, "evaluation-dataset": datasetInput}}
	specification, err := WrapExecution(fixtureExecution(t, corpus.Manifest.Digest, dataset.Manifest.Digest), resolved, "bad")
	if err != nil {
		t.Fatal(err)
	}
	binary := buildWorker(t)
	runner, _ := processrunner.New(processrunner.Config{Command: []string{binary}, Runner: lab.RunnerRecord{Name: RunnerName, ResolvedVersion: RunnerVersion}, Domains: []lab.DomainVersion{{Domain: ragcontract.Domain, SchemaVersion: ragcontract.DomainSchemaVersion}}})
	err = runner.Start(context.Background(), lab.AttemptRequest{Specification: specification, Run: lab.RunRecord{ID: "r", CreatedAt: "2026-07-17T00:00:00Z"}, AttemptID: "a", AttemptIndex: 1, ArtifactRoot: root}, &sink{})
	if err == nil || !strings.Contains(err.Error(), "RAG_WORKER_INPUT_LINEAGE") {
		t.Fatalf("err=%v", err)
	}
}
func TestCheckWorkerLaunchAndMalformedCapabilityFailures(t *testing.T) {
	if err := CheckWorker(context.Background(), WorkerCommand{Executable: filepath.Join(t.TempDir(), "missing")}); err == nil || !strings.Contains(err.Error(), "RAG_WORKER_LAUNCH") {
		t.Fatalf("err=%v", err)
	}
	script := filepath.Join(t.TempDir(), "worker.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho '{\"type\":\"hello\",\"hello\":{}}'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := CheckWorker(context.Background(), WorkerCommand{Executable: script}); err == nil || !strings.Contains(err.Error(), "RAG_WORKER_CAPABILITY") {
		t.Fatalf("err=%v", err)
	}
}
func buildWorker(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rag-worker")
	command := exec.Command("go", "build", "-o", path, "../../cmd/rag-worker")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	return path
}
func fixtureExecution(t *testing.T, corpusDigest, datasetDigest string) ragcontract.PipelineExecution {
	t.Helper()
	pipeline := ragmodel.NewPipeline("raw", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 100})).Represent(ragmodel.RawRepresentation("raw")).IndexNamed("representations", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true}))
	})
	query := ragmodel.NewQueryPlan("query", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 5})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"}))
	})
	product := ragmodel.NewProduct("raw", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Citations("source") })
	})
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "manifest-envelope", Digest: corpusDigest, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: plan.Pipeline, Bindings: plan.Bindings, Dataset: ragcontract.DatasetBinding{ManifestDigest: datasetDigest, Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", Version: "v1", ValueKind: "number", Unit: "ratio", Required: true, Config: json.RawMessage(`{}`)}}, VariantID: "raw", Factors: []ragcontract.FactorSelection{}}
	execution.CellID, _ = ragcontract.Digest(execution)
	return execution
}
