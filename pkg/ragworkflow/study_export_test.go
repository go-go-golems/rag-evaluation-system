package ragworkflow

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/stretchr/testify/require"
)

func TestWriteStudyBundleIsStableAndWorkflowNative(t *testing.T) {
	fixture, err := NewProviderFreeFixture(true)
	require.NoError(t, err)
	root := t.TempDir()
	cases := []StudyWorkflowCase{
		{ID: "case-a", Execution: fixture.Execution, Corpus: fixture.Corpus, Dataset: fixture.Dataset, Replicates: 2},
		{ID: "case-b", Execution: fixture.Execution, Corpus: fixture.Corpus, Dataset: fixture.Dataset, Replicates: 1},
	}
	first, err := WriteStudyBundle(context.Background(), root, filepath.Join(root, "inputs", "study"), "fixture study", "EXP-FIXTURE", cases, nil)
	require.NoError(t, err)
	plan, err := os.ReadFile(filepath.Join(root, first.Plan.Path))
	require.NoError(t, err)
	require.Contains(t, string(plan), `domain: "scraper-workflow"`)
	require.NotContains(t, string(plan), "rag-worker")
	require.NotContains(t, string(plan), "execute-spec")
	require.Len(t, first.Cases, 2)
	require.Equal(t, 2, first.Cases[0].Replicates)

	secondRoot := t.TempDir()
	second, err := WriteStudyBundle(context.Background(), secondRoot, filepath.Join(secondRoot, "inputs", "study"), "fixture study", "EXP-FIXTURE", cases, nil)
	require.NoError(t, err)
	secondPlan, err := os.ReadFile(filepath.Join(secondRoot, second.Plan.Path))
	require.NoError(t, err)
	require.Equal(t, plan, secondPlan)
	require.Equal(t, first.Plan.Digest, second.Plan.Digest)
	_, err = WriteStudyBundle(context.Background(), root, filepath.Join(root, "inputs", "study"), "fixture study", "EXP-FIXTURE", cases, nil)
	require.NoError(t, err, "identical compilation must be idempotent")
	changed := append([]StudyWorkflowCase(nil), cases...)
	changed[0].Corpus.Records = append([]ragoperators.SourceRecord(nil), changed[0].Corpus.Records...)
	changed[0].Corpus.Records[0].Text += " changed"
	_, err = WriteStudyBundle(context.Background(), root, filepath.Join(root, "inputs", "study"), "fixture study", "EXP-FIXTURE", changed, nil)
	require.ErrorContains(t, err, "RAG_WORKFLOW_STUDY_CONFLICT")
	for index := range first.Cases {
		require.Equal(t, first.Cases[index].DomainConfig.Digest, second.Cases[index].DomainConfig.Digest)
	}
}

func TestWriteImmutableStudyFileUsesOwnerOnlyPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "artifact.json")
	require.NoError(t, writeImmutableStudyFile(path, []byte("{}\n")))
	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}

func TestWriteStudyBundleBindsProviderAuthority(t *testing.T) {
	fixture, err := NewDeterministicProviderFixture()
	require.NoError(t, err)
	services, err := NewDeterministicProviderServices()
	require.NoError(t, err)
	providerPackage, err := NewProviderPackage(services, defaultProviderPolicy())
	require.NoError(t, err)
	root := t.TempDir()
	bundle, err := WriteStudyBundle(context.Background(), root, filepath.Join(root, "provider"), "provider study", "EXP-PROVIDER", []StudyWorkflowCase{{ID: "provider", Execution: fixture.Execution, Corpus: fixture.Corpus, Dataset: fixture.Dataset, Replicates: 2}}, providerPackage)
	require.NoError(t, err)
	domainBody, err := os.ReadFile(filepath.Join(root, bundle.Cases[0].DomainConfig.Path))
	require.NoError(t, err)
	require.Contains(t, string(domainBody), ProviderPackageName)
	require.Contains(t, string(domainBody), providerPackage.authority.Digest)
	planBody, err := os.ReadFile(filepath.Join(root, bundle.Plan.Path))
	require.NoError(t, err)
	require.NotContains(t, string(planBody), "fixture answer evidence")
	require.NotContains(t, string(planBody), "provider-config")
}

func TestWriteStudyBundleRejectsOutputEscapeAndProviderMismatch(t *testing.T) {
	fixture, err := NewProviderFreeFixture(true)
	require.NoError(t, err)
	root := t.TempDir()
	_, err = WriteStudyBundle(context.Background(), root, filepath.Join(filepath.Dir(root), "escape"), "fixture", "EXP", []StudyWorkflowCase{{ID: "case", Execution: fixture.Execution, Corpus: fixture.Corpus, Dataset: fixture.Dataset, Replicates: 1}}, nil)
	require.ErrorContains(t, err, "RAG_WORKFLOW_STUDY_OUTPUT_BOUNDARY")

	providerFixture, err := NewDeterministicProviderFixture()
	require.NoError(t, err)
	_, err = WriteStudyBundle(context.Background(), root, filepath.Join(root, "provider"), "provider", "EXP", []StudyWorkflowCase{{ID: "provider", Execution: providerFixture.Execution, Corpus: providerFixture.Corpus, Dataset: providerFixture.Dataset, Replicates: 1}}, nil)
	require.ErrorContains(t, err, "RAG_WORKFLOW_PROVIDER_REQUIRED")
}
