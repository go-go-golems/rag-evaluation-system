package study

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

func TestStudyValidateExplainCompileCommands(t *testing.T) {
	studyPath := filepath.Join("..", "..", "..", "..", "examples", "rag-v2", "06-raw-study.js")
	if _, err := LoadStudy(studyPath); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "s", Text: "rank fusion"}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{{ID: "q", Text: "rank", RelevantIDs: []string{"s"}, Grades: map[string]float64{"s": 1}}}}, "fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest)
	corpusPath := writeFixture(t, root, "corpus.json", corpus)
	datasetPath := writeFixture(t, root, "dataset.json", dataset)
	inputsPath := writeFixture(t, root, "inputs.json", map[string]any{"inputs": map[string]any{"corpus": map[string]any{"uri": corpusPath}, "evaluation-dataset": map[string]any{"uri": datasetPath}}})
	artifactRoot := filepath.Join(root, "artifacts")
	for _, subcommand := range []string{"validate", "explain", "compile"} {
		command := NewCommand()
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetErr(&output)
		args := []string{subcommand, studyPath, "--inputs", inputsPath, "--artifact-root", artifactRoot}
		if subcommand == "compile" {
			args = append(args, "--output-dir", filepath.Join(artifactRoot, "inputs", "study"), "--experiment-id", "EXP-STUDY")
		}
		command.SetArgs(args)
		if err := command.Execute(); err != nil {
			t.Fatalf("%s: %v\n%s", subcommand, err, output.String())
		}
		if !strings.Contains(output.String(), "rag-study") && !strings.Contains(output.String(), "rag-workflow-study-bundle") && !strings.Contains(output.String(), `"valid": true`) {
			t.Fatalf("%s output=%s", subcommand, output.String())
		}
	}
	plan, err := os.ReadFile(filepath.Join(artifactRoot, "inputs", "study", "researchctl-plan.js"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(plan), "rag-worker") || strings.Contains(string(plan), "execute-spec") {
		t.Fatalf("compiled plan retained direct execution: %s", plan)
	}
	command := NewCommand()
	command.SetArgs([]string{"run"})
	if err := command.Execute(); err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("study run remains available: %v", err)
	}
}
func writeFixture(t *testing.T, root, name string, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}
