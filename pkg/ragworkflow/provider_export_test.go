package ragworkflow

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteDeterministicProviderFixturesIsCanonical(t *testing.T) {
	first, second := filepath.Join(t.TempDir(), "first"), filepath.Join(t.TempDir(), "second")
	left, err := WriteDeterministicProviderFixtures(context.Background(), first)
	require.NoError(t, err)
	right, err := WriteDeterministicProviderFixtures(context.Background(), second)
	require.NoError(t, err)
	require.Equal(t, left, right)
	for _, item := range append([]FixtureFile{{Path: "manifest.json"}}, fixtureFiles(left)...) {
		leftBody, readErr := os.ReadFile(filepath.Join(first, filepath.FromSlash(item.Path)))
		require.NoError(t, readErr)
		rightBody, readErr := os.ReadFile(filepath.Join(second, filepath.FromSlash(item.Path)))
		require.NoError(t, readErr)
		require.Equal(t, leftBody, rightBody, item.Path)
	}
}

func fixtureFiles(manifest FixtureManifest) []FixtureFile {
	files := make([]FixtureFile, 0, len(manifest.Cases)*5)
	for _, item := range manifest.Cases {
		files = append(files, item.Execution, item.Corpus, item.Queries, item.DomainConfig, item.Parity)
	}
	return files
}
