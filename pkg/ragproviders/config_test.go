package ragproviders

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigRejectsEscapingAndMultipleDocuments(t *testing.T) {
	tests := []struct{ name, body, want string }{
		{"escape", "schemaVersion: rag-provider-host-config/v1\nprofileId: test\nmanifests: {modelsDir: ../escape}\nproviders: {x: {kind: test, modelManifest: test}}\n", "RAG_PROVIDER_CONFIG_PATH_ESCAPE"},
		{"absolute", "schemaVersion: rag-provider-host-config/v1\nprofileId: test\nmanifests: {modelsDir: /tmp/models}\nproviders: {x: {kind: test, modelManifest: test}}\n", "RAG_PROVIDER_CONFIG_PATH_ABSOLUTE"},
		{"multiple", "schemaVersion: rag-provider-host-config/v1\nprofileId: test\nproviders: {x: {kind: test, modelManifest: test}}\n---\n{}\n", "RAG_PROVIDER_CONFIG_MULTIPLE_DOCUMENTS"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "providers.yaml")
			if err := os.WriteFile(path, []byte(test.body), 0o644); err != nil {
				t.Fatal(err)
			}
			_, _, err := loadConfig(path)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %s", err, test.want)
			}
		})
	}
}
