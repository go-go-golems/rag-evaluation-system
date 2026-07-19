package ragproviders

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestLoadProviderSetFromStrictHostConfig(t *testing.T) {
	config := writeProviderFixture(t)
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
	if err := set.Schemas.Validate("summary/v1", json.RawMessage(`{}`)); err != nil {
		t.Fatalf("schema Validate() = %v", err)
	}
}

func TestLoadProviderSetRejectsInvalidProviderSpecifications(t *testing.T) {
	tests := []struct {
		name, old, new, want string
	}{
		{"unknown kind", "kind: geppetto-embedding/v1", "kind: unknown/v1", "RAG_PROVIDER_KIND_UNSUPPORTED"},
		{"model adapter mismatch", "modelManifest: embedding", "modelManifest: generator", "RAG_PROVIDER_MODEL_INCOMPATIBLE"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := writeProviderFixture(t)
			replaceConfig(t, config, test.old, test.new)
			_, err := Load(context.Background(), config)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load() error = %v, want %s", err, test.want)
			}
		})
	}
}

func TestLoadProviderSetRejectsMissingDirectories(t *testing.T) {
	tests := []struct {
		name, old, new, want string
	}{
		{"models", "modelsDir: models", "modelsDir: missing-models", "RAG_MODEL_MANIFEST_DIRECTORY"},
		{"schemas", "directory: schemas", "directory: missing-schemas", "RAG_SCHEMA_DIRECTORY"},
		{"cache configuration", "directory: cache", "directory: ", "RAG_PROVIDER_CONFIG_CACHE"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := writeProviderFixture(t)
			replaceConfig(t, config, test.old, test.new)
			_, err := Load(context.Background(), config)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load() error = %v, want %s", err, test.want)
			}
		})
	}
}

func TestLoadProviderSetRedactsEndpointAndMissingCredential(t *testing.T) {
	t.Run("endpoint userinfo", func(t *testing.T) {
		config := writeProviderFixture(t)
		const secret = "do-not-echo-endpoint-secret"
		t.Setenv("TEST_RAG_EMBED_URL", "http://user:"+secret+"@127.0.0.1:11434")
		_, err := Load(context.Background(), config)
		if err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_ENDPOINT_POLICY") || strings.Contains(err.Error(), secret) {
			t.Fatalf("Load() error = %v; endpoint policy failure must not expose credentials", err)
		}
	})
	t.Run("missing credential", func(t *testing.T) {
		config := writeProviderFixture(t)
		replaceConfig(t, config, "allowLocalNetworks: true\n  reranker-primary:", "allowLocalNetworks: true\n    credentialRef: env:TEST_RAG_GENERATOR_SECRET\n  reranker-primary:")
		_, err := Load(context.Background(), config)
		if err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_ENV_MISSING") || strings.Contains(err.Error(), "TEST_RAG_GENERATOR_SECRET") {
			t.Fatalf("Load() error = %v; missing credential failure must not expose its reference", err)
		}
	})
}

func TestLoadProviderSetRejectsPromptSchemaMissing(t *testing.T) {
	config := writeProviderFixture(t)
	replaceConfig(t, filepath.Join(filepath.Dir(config), "prompts", "summary.json"), `"outputSchema":"summary/v1"`, `"outputSchema":"missing/v1"`)
	_, err := Load(context.Background(), config)
	if err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_PROMPT_SCHEMA_MISSING") {
		t.Fatalf("Load() error = %v, want prompt schema failure", err)
	}
}

func TestProviderSetChecksExecutionRequirements(t *testing.T) {
	set, err := Load(context.Background(), writeProviderFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = set.Close() }()
	execution := ragcontract.PipelineExecution{Pipeline: ragcontract.PipelineIR{Nodes: []ragcontract.Node{{Operator: ragcontract.OperatorRef{Kind: "embed.model", Version: "v1"}, Config: json.RawMessage(`{"model":"embedding","dimensions":768}`)}}}}
	if err := set.CheckExecution(execution); err != nil {
		t.Fatalf("CheckExecution() error = %v", err)
	}
	execution.Pipeline.Nodes[0].Config = json.RawMessage(`{"model":"embedding","dimensions":42}`)
	if err := set.CheckExecution(execution); err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_EXECUTION_MODEL") {
		t.Fatalf("CheckExecution() error = %v", err)
	}
}

func TestProviderSetCloseIsIdempotentAndClosesAllResources(t *testing.T) {
	firstErr := errors.New("first close failure")
	first, second := &recordingCloser{err: firstErr}, &recordingCloser{}
	set := &ProviderSet{closers: []io.Closer{first, second}}
	if err := set.Close(); !errors.Is(err, firstErr) {
		t.Fatalf("Close() error = %v, want %v", err, firstErr)
	}
	if first.calls != 1 || second.calls != 1 {
		t.Fatalf("close calls = first:%d second:%d, want one each", first.calls, second.calls)
	}
	if err := set.Close(); !errors.Is(err, firstErr) {
		t.Fatalf("second Close() error = %v, want cached close error", err)
	}
	if first.calls != 1 || second.calls != 1 {
		t.Fatalf("second close repeated resource shutdown: first:%d second:%d", first.calls, second.calls)
	}
}

type recordingCloser struct {
	calls int
	err   error
}

func (c *recordingCloser) Close() error {
	c.calls++
	return c.err
}

func writeProviderFixture(t *testing.T) string {
	t.Helper()
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
	for _, prompt := range []struct{ id, schema string }{{"summary", "summary/v1"}, {"questions", "questions/v1"}, {"answer", "answer/v1"}} {
		writeJSONFile(t, filepath.Join(prompts, prompt.id+".json"), ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + repeat("1", 64)}, PromptID: prompt.id, TemplateDigest: "sha256:" + repeat("2", 64), InputSchema: "text", OutputSchema: prompt.schema})
		if err := os.WriteFile(filepath.Join(prompts, prompt.id+".txt"), []byte("Return the requested JSON."), 0o644); err != nil {
			t.Fatal(err)
		}
		writeJSONFile(t, filepath.Join(schemas, prompt.id+".json"), map[string]any{"$id": prompt.schema, "type": "object"})
	}
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
	return config
}

func replaceConfig(t *testing.T, path, old, replacement string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(data), old, replacement, 1)
	if updated == string(data) {
		t.Fatalf("did not find %q in %s", old, path)
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
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

func repeat(value string, count int) string { return strings.Repeat(value, count) }
