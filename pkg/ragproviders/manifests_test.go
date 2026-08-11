package ragproviders

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestManifestRegistryRequiresPromptTextAndRejectsDuplicateText(t *testing.T) {
	dir := t.TempDir()
	manifest := ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + strings.Repeat("a", 64)}, PromptID: "answer", TemplateDigest: "sha256:" + strings.Repeat("b", 64), OutputSchema: "answer/v1"}
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "alias.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := loadPrompts(dir); err == nil || !strings.Contains(err.Error(), "RAG_PROMPT_TEXT_MISSING") {
		t.Fatalf("error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "answer.txt"), []byte("prompt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "answer.md"), []byte("prompt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := loadPrompts(dir); err == nil || !strings.Contains(err.Error(), "RAG_PROMPT_TEXT_DUPLICATE") {
		t.Fatalf("error = %v", err)
	}
}
