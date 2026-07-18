package ragproviders

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestLoadProviderSetFromStrictHostConfig(t *testing.T) {
	root := t.TempDir()
	models := filepath.Join(root, "models")
	prompts := filepath.Join(root, "prompts")
	schemas := filepath.Join(root, "schemas")
	for _, dir := range []string{models, prompts, schemas} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeJSONFile(t, filepath.Join(models, "embedding.json"), ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + repeat("a", 64)}, ProviderAdapterVersion: "geppetto-embedding/v1", ModelID: "nomic-embed-text", ModelDigest: "sha256:" + repeat("b", 64), Dimensions: 768, Tokenization: "exact", Truncation: "none", Normalization: "l2", ImplementationVersion: "test", RequestParameters: json.RawMessage(`{}`)})
	writeJSONFile(t, filepath.Join(models, "generator.json"), ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + repeat("c", 64)}, ProviderAdapterVersion: "geppetto-generation/v1", ModelID: "qwen3:8b", ModelDigest: "sha256:" + repeat("d", 64), Tokenization: "exact", Truncation: "none", Normalization: "none", ImplementationVersion: "test", RequestParameters: json.RawMessage(`{"temperature":0}`)})
	writeJSONFile(t, filepath.Join(models, "reranker.json"), ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + repeat("e", 64)}, ProviderAdapterVersion: "geppetto-reranker/v1", ModelID: "bge", ModelDigest: "sha256:" + repeat("f", 64), Tokenization: "exact", Truncation: "none", Normalization: "none", ImplementationVersion: "test", RequestParameters: json.RawMessage(`{}`)})
	for _, p := range []struct{ id, schema string }{{"summary", "summary/v1"}, {"questions", "questions/v1"}, {"answer", "answer/v1"}} {
		writeJSONFile(t, filepath.Join(prompts, p.id+".json"), ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + repeat("1", 64)}, PromptID: p.id, TemplateDigest: "sha256:" + repeat("2", 64), InputSchema: "text", OutputSchema: p.schema})
		if err := os.WriteFile(filepath.Join(prompts, p.id+".txt"), []byte("Return the requested JSON."), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeJSONFile(t, filepath.Join(schemas, "summary.json"), map[string]any{"type": "object"})
	writeJSONFile(t, filepath.Join(schemas, "questions.json"), map[string]any{"type": "object"})
	writeJSONFile(t, filepath.Join(schemas, "answer.json"), map[string]any{"type": "object"})
	config := filepath.Join(root, "providers.yaml")
	if err := os.WriteFile(config, []byte(`schemaVersion: rag-provider-host-config/v1
profileId: test-profile
manifests:
  modelsDir: models
  promptsDir: prompts
schemas:
  directory: schemas
cache:
  kind: filesystem-content-addressed/v1
  directory: cache
providers:
  embedding-primary:
    kind: geppetto-embedding/v1
    modelManifest: embedding
    endpointRef: env:TEST_RAG_EMBED_URL
    allowHttp: true
    allowLocalNetworks: true
  generator-primary:
    kind: geppetto-generation/v1
    modelManifest: generator
    endpointRef: env:TEST_RAG_GENERATE_URL
    allowHttp: true
    allowLocalNetworks: true
  reranker-primary:
    kind: geppetto-reranker/v1
    modelManifest: reranker
    endpointRef: env:TEST_RAG_RERANK_URL
    allowHttp: true
    allowLocalNetworks: true
`), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, value := range map[string]string{"TEST_RAG_EMBED_URL": "http://127.0.0.1:11434", "TEST_RAG_GENERATE_URL": "http://127.0.0.1:11434/v1", "TEST_RAG_RERANK_URL": "http://127.0.0.1:18012"} {
		t.Setenv(name, value)
	}
	set, err := Load(context.Background(), config)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	defer func() { _ = set.Close() }()
	if set.ProfileID != "test-profile" || set.Generator == nil || set.Embedder == nil || set.Reranker == nil || set.Cache == nil {
		t.Fatalf("incomplete provider set: %#v", set)
	}
	if got, err := set.Manifests.PromptText("summary"); err != nil || got == "" {
		t.Fatalf("PromptText() = %q, %v", got, err)
	}
	if err := set.Schemas.Validate("summary", json.RawMessage(`{}`)); err != nil {
		t.Fatalf("schema Validate() = %v", err)
	}
}

func TestLoadProviderSetRejectsMissingSecretReference(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "providers.yaml"), []byte("schemaVersion: rag-provider-host-config/v1\nprofileId: x\nproviders: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(context.Background(), filepath.Join(root, "providers.yaml"))
	if err == nil {
		t.Fatal("Load() error = nil")
	}
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
func repeat(value string, count int) string {
	out := ""
	for i := 0; i < count; i++ {
		out += value
	}
	return out
}
