package ragproviders

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSchemaRegistryRejectsConcatenatedJSONOutput(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "answer.json"), []byte(`{"type":"object","required":["answer"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	registry, err := LoadSchemaRegistry(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.Validate("answer", json.RawMessage(`{"answer":"ok"}{"answer":"extra"}`)); err == nil || !strings.Contains(err.Error(), "RAG_OUTPUT_SCHEMA_JSON_MULTIPLE_VALUES") {
		t.Fatalf("Validate() = %v", err)
	}
	if err := registry.Validate("answer", json.RawMessage(`{"answer":"ok"}`)); err != nil {
		t.Fatalf("Validate(valid) = %v", err)
	}
}
