package researchctladapter

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/stretchr/testify/require"
)

func TestResolveAndLoadDomainArtifactsPreservesCanonicalInputs(t *testing.T) {
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(
		ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{
			{ID: "source", SessionID: "s", Ordinal: 1, Role: "user", Text: "weighted reciprocal rank fusion"},
		}},
		"fixture",
	)
	dataset := ragoperators.NewEvaluationArtifact(
		ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{
			{ID: "q", Text: "rank fusion", RelevantIDs: []string{"source"}},
		}},
		"fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest,
	)
	corpusPath := writeInput(t, root, "corpus-source.json", corpus)
	datasetPath := writeInput(t, root, "dataset-source.json", dataset)
	document := InputDocument{Inputs: map[string]InputReference{
		"corpus":             {Role: "corpus", URI: corpusPath},
		"evaluation-dataset": {Role: "evaluation-dataset", URI: datasetPath},
	}}
	artifactRoot := filepath.Join(root, "artifacts")
	resolved, err := ResolveInputs(context.Background(), document, root, artifactRoot, nil)
	require.NoError(t, err)
	loadedCorpus, loadedDataset, err := LoadDomainArtifacts(artifactRoot, resolved)
	require.NoError(t, err)
	require.Equal(t, corpus.Manifest.Digest, loadedCorpus.Manifest.Digest)
	require.Equal(t, dataset.Manifest.Digest, loadedDataset.Manifest.Digest)
	require.Equal(t, corpus.Corpus, loadedCorpus.Corpus)
	require.Equal(t, dataset.Dataset, loadedDataset.Dataset)
}

func TestLoadDomainArtifactsRejectsTamperedStagedInput(t *testing.T) {
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(
		ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "source", Text: "safe"}}},
		"fixture",
	)
	resolved, err := StageEnvelope(InputReference{Role: "corpus"}, corpus, root)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(root, filepath.FromSlash(resolved.Reference.URI)), []byte(`{"tampered":true}`), 0o644))
	_, _, err = LoadDomainArtifacts(root, ResolvedInputs{ByRole: map[string]ResolvedInput{"corpus": resolved, "evaluation-dataset": resolved}})
	require.ErrorContains(t, err, "RAG_INPUT_DIGEST")
}

func writeInput(t *testing.T, root, name string, value any) string {
	t.Helper()
	body, err := ragcontract.CanonicalJSON(value)
	require.NoError(t, err)
	path := filepath.Join(root, name)
	require.NoError(t, os.WriteFile(path, body, 0o644))
	return path
}
