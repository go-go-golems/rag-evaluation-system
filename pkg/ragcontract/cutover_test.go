package ragcontract

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrototypeRuntimeIsAbsent(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	for _, retired := range []string{"pkg/raglab", "cmd/rag-lab-worker", "internal/services/immutableretrieval"} {
		if _, err := os.Stat(filepath.Join(root, retired)); !os.IsNotExist(err) {
			t.Fatalf("retired runtime path still exists: %s", retired)
		}
	}

	retiredIdentifiers := []string{
		strings.Join([]string{"github.com/go-go-golems/rag-evaluation-system/pkg", "raglab"}, "/"),
		strings.Join([]string{"rag-retrieval-spec", "v1"}, "/"),
		strings.Join([]string{"researchctl-rag-runner-stdio", "v1"}, "/"),
	}
	for _, dir := range []string{"cmd", "internal", "pkg"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
				return err
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			for _, identifier := range retiredIdentifiers {
				if strings.Contains(string(data), identifier) {
					t.Errorf("retired identifier %q remains in %s", identifier, path)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
