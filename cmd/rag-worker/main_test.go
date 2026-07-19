package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

func TestExternalWorkerSpeaksGenericProtocol(t *testing.T) {
	binary := buildWorker(t)
	corpusArtifact := ragoperators.NewCorpusArtifact(ragoperators.Corpus{Records: []ragoperators.SourceRecord{{ID: "s1", SessionID: "s", Ordinal: 1, Role: "user", Text: "reciprocal rank fusion"}}}, "fixture")
	evaluationArtifact := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q1", Text: "rank fusion"}}}, "fixture", "smoke", "candidate", "unit", corpusArtifact.Manifest.Digest)
	execution := workerExecution(t, corpusArtifact.Manifest.Digest, evaluationArtifact.Manifest.Digest)
	root := t.TempDir()
	corpusPath := writeJSON(t, root, "corpus.json", corpusArtifact)
	datasetPath := writeJSON(t, root, "dataset.json", evaluationArtifact)
	request := map[string]any{"protocolVersion": protocolVersion, "attempt": map[string]any{"specification": map[string]any{"canonicalIdentity": map[string]any{"domain": ragcontract.Domain, "domainSchemaVersion": ragcontract.DomainSchemaVersion, "domainConfig": execution}}}, "inputs": []map[string]any{{"reference": map[string]any{"role": "corpus", "schemaVersion": ragcontract.CorpusManifestSchema}, "path": corpusPath}, {"reference": map[string]any{"role": "evaluation-dataset", "schemaVersion": ragcontract.EvaluationManifestSchema}, "path": datasetPath}}}
	frames, stderr, err := runWorker(binary, request)
	if err != nil {
		t.Fatalf("worker: %v stderr=%s", err, stderr)
	}
	if len(frames) < 5 || frames[0]["type"] != "hello" || frames[len(frames)-1]["type"] != "complete" {
		t.Fatalf("frames=%#v", frames)
	}
	seen := map[string]bool{}
	for _, frame := range frames {
		seen[frame["type"].(string)] = true
		if frame["type"] == "trace" && frame["trace"].(map[string]any)["kind"] != ragcontract.TraceSchemaVersion {
			t.Fatalf("trace=%#v", frame)
		}
	}
	for _, kind := range []string{"hello", "event", "trace", "metric", "artifact", "complete"} {
		if !seen[kind] {
			t.Fatalf("missing %s: %#v", kind, frames)
		}
	}
}
func TestExternalWorkerAdvertisesBeforeDomainError(t *testing.T) {
	binary := buildWorker(t)
	request := map[string]any{"protocolVersion": protocolVersion, "attempt": map[string]any{"specification": map[string]any{"canonicalIdentity": map[string]any{"domain": "other", "domainSchemaVersion": "other/v1", "domainConfig": map[string]any{}}}}}
	frames, _, err := runWorker(binary, request)
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 2 || frames[0]["type"] != "hello" || frames[1]["type"] != "error" {
		t.Fatalf("frames=%#v", frames)
	}
}
func buildWorker(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rag-worker")
	command := exec.Command("go", "build", "-o", path, ".")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	return path
}
func runWorker(binary string, request any) ([]map[string]any, string, error) {
	payload, _ := json.Marshal(request)
	command := exec.Command(binary, "--provider-profile", "fixtures")
	command.Stdin = bytes.NewReader(append(payload, '\n'))
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	frames := []map[string]any{}
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		var frame map[string]any
		if decodeErr := json.Unmarshal(scanner.Bytes(), &frame); decodeErr != nil {
			return frames, stderr.String(), decodeErr
		}
		frames = append(frames, frame)
	}
	if err == nil {
		err = scanner.Err()
	}
	return frames, stderr.String(), err
}
func writeJSON(t *testing.T, root, name string, value any) string {
	t.Helper()
	path := filepath.Join(root, name)
	data, _ := json.Marshal(value)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
func workerExecution(t *testing.T, corpusDigest, evaluationDigest string) ragcontract.PipelineExecution {
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
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "json", Digest: corpusDigest, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: plan.Pipeline, Bindings: plan.Bindings, Dataset: ragcontract.DatasetBinding{ManifestDigest: evaluationDigest, Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", Version: "v1", ValueKind: "number", Unit: "ratio", Required: true, Config: json.RawMessage(`{}`)}}, VariantID: "raw", Factors: []ragcontract.FactorSelection{}}
	execution.CellID, _ = ragcontract.Digest(execution)
	return execution
}
