package ragcontract

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRetiredRuntimeIsAbsent(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	for _, retired := range []string{filepath.Join("pkg", "rag"+"lab"), filepath.Join("cmd", "rag-"+"lab-worker"), filepath.Join("internal", "services", "immutable"+"retrieval")} {
		if _, err := os.Stat(filepath.Join(root, retired)); !os.IsNotExist(err) {
			t.Fatalf("retired runtime path still exists: %s", retired)
		}
	}

	retiredIdentifiers := []string{
		strings.Join([]string{"github.com/go-go-golems/rag-evaluation-system/pkg", "raglab"}, "/"),
		strings.Join([]string{"rag-retrieval-spec", "v1"}, "/"),
		strings.Join([]string{"rag-query-trace", "v1"}, "/"),
		strings.Join([]string{"researchctl-rag-runner-stdio", "v1"}, "/"),
		"rag-" + "lab-worker",
		"run-" + "rag",
		strings.Join([]string{"", "api", "v1", "lab"}, "/"),
		"rag-" + "lab-js",
		"export" + "Specification",
		"raw" + "Chunks",
		"fail-" + "closed",
		"experiment_" + "specs",
		"experiment_" + "runs",
		"experiment_" + "run_events",
		"experiment_" + "run_summaries",
		"experiment_" + "query_traces",
	}
	extensions := map[string]bool{".go": true, ".md": true, ".js": true, ".ts": true, ".tsx": true, ".json": true, ".yaml": true, ".yml": true}
	for _, dir := range []string{"cmd", "docs", "examples", "experiments", "internal", "pkg", "web"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return err
			}
			if strings.Contains(path, string(filepath.Separator)+"dist"+string(filepath.Separator)) || strings.HasSuffix(path, "_test.go") || !extensions[filepath.Ext(path)] {
				return nil
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
