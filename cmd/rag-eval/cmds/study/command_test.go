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
	studyPath := filepath.Join("..", "..", "..", "..", "examples", "rag-v2", "02-five-variant-study.js")
	if _, err := LoadStudy(studyPath); err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{Records: []ragoperators.SourceRecord{{ID: "s", Text: "rank fusion"}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{Queries: []ragoperators.Query{{ID: "q", Text: "rank"}}}, "fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest)
	corpusPath := writeFixture(t, root, "corpus.json", corpus)
	datasetPath := writeFixture(t, root, "dataset.json", dataset)
	inputsPath := writeFixture(t, root, "inputs.json", map[string]any{"inputs": map[string]any{"corpus": map[string]any{"uri": corpusPath}, "evaluation-dataset": map[string]any{"uri": datasetPath}}})
	for _, subcommand := range []string{"validate", "explain", "compile"} {
		command := NewCommand()
		var output bytes.Buffer
		command.SetOut(&output)
		command.SetErr(&output)
		command.SetArgs([]string{subcommand, studyPath, "--inputs", inputsPath, "--artifact-root", filepath.Join(root, "artifacts")})
		if err := command.Execute(); err != nil {
			t.Fatalf("%s: %v\n%s", subcommand, err, output.String())
		}
		if !strings.Contains(output.String(), "rag-study") && !strings.Contains(output.String(), `"valid": true`) {
			t.Fatalf("%s output=%s", subcommand, output.String())
		}
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
