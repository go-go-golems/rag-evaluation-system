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

func TestLoadProviderSetUsesOneDefaultProfileRegistryAndValidatesModels(t *testing.T) {
	config := writeProviderFixture(t)
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	profilesDir := filepath.Join(configDir, "pinocchio")
	if err := os.MkdirAll(profilesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	profiles := `slug: default
profiles:
  test-embed:
    inference_settings:
      embeddings:
        type: ollama
        engine: nomic-embed-text
        dimensions: 768
        base_urls:
          ollama-base-url: http://127.0.0.1:11434
      api:
        base_urls:
          ollama-base-url: http://127.0.0.1:11434
        allow_http:
          ollama: true
        allow_local_networks:
          ollama: true
  test-generator:
    inference_settings:
      chat:
        api_type: openai
        engine: qwen3:8b
      api:
        api_keys:
          openai-api-key: test-key
        base_urls:
          openai-base-url: http://127.0.0.1:11434/v1
        allow_http:
          openai: true
        allow_local_networks:
          openai: true
  test-reranker:
    inference_settings:
      rerank:
        type: llamacpp
        engine: bge
      api:
        base_urls:
          rerank-base-url: http://127.0.0.1:18012
        allow_http:
          rerank: true
        allow_local_networks:
          rerank: true
`
	profilePath := filepath.Join(profilesDir, "profiles.yaml")
	if err := os.WriteFile(profilePath, []byte(profiles), 0o600); err != nil {
		t.Fatal(err)
	}
	replaceConfig(t, config, "endpointRef: env:TEST_RAG_EMBED_URL\n    allowHttp: true\n    allowLocalNetworks: true", "profile: test-embed")
	replaceConfig(t, config, "endpointRef: env:TEST_RAG_GENERATE_URL\n    allowHttp: true\n    allowLocalNetworks: true", "profile: test-generator\n    concurrency:\n      maxInFlight: 3\n    generation:\n      maxResponseTokens: 8192\n      pricing:\n        inputMicrounitsPerMillion: 150000\n        outputMicrounitsPerMillion: 1000000\n        cacheReadMicrounitsPerMillion: 50000\n        cacheWriteMicrounitsPerMillion: 150000")
	replaceConfig(t, config, "endpointRef: env:TEST_RAG_RERANK_URL\n    allowHttp: true\n    allowLocalNetworks: true", "profile: test-reranker")
	set, err := Load(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = set.Close() }()
	if got := len(set.profileChains); got != 1 {
		t.Fatalf("profile chains = %d, want 1", got)
	}
	if set.GenerationConcurrency != 3 {
		t.Fatalf("generation concurrency = %d, want 3", set.GenerationConcurrency)
	}
	if set.EngineOptions().GenerationConcurrency != 3 {
		t.Fatalf("engine generation concurrency = %d, want 3", set.EngineOptions().GenerationConcurrency)
	}
	descriptor, err := json.Marshal(set.CapabilityDescriptor())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(descriptor), "test-key") {
		t.Fatalf("capability descriptor leaked profile credential: %s", descriptor)
	}
	if len(set.CapabilityDescriptor().EffectiveProviderIdentities) != 3 {
		t.Fatalf("effective identities = %#v", set.CapabilityDescriptor().EffectiveProviderIdentities)
	}
	for _, identity := range set.CapabilityDescriptor().EffectiveProviderIdentities {
		if identity.SettingsFingerprint == "" || identity.ModelID == "" {
			t.Fatalf("incomplete identity: %#v", identity)
		}
		if identity.Role == "generator-primary" && (!identity.PricingConfigured || identity.MaxResponseTokens != 8192 || identity.InputCostMicrounitsPerMillion != 150000 || identity.OutputCostMicrounitsPerMillion != 1000000 || identity.CacheReadCostMicrounitsPerMillion != 50000 || identity.CacheWriteCostMicrounitsPerMillion != 150000) {
			t.Fatalf("generation policy identity = %#v", identity)
		}
	}

	if err := os.WriteFile(profilePath, []byte(strings.Replace(profiles, "engine: qwen3:8b", "engine: wrong-model", 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = Load(context.Background(), config)
	if err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_PROFILE_MODEL_MISMATCH") {
		t.Fatalf("Load() mismatch error = %v", err)
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

func TestLoadProviderSetRejectsIncompleteGenerationPolicy(t *testing.T) {
	config := writeProviderFixture(t)
	replaceConfig(t, config, "    allowLocalNetworks: true\n  reranker-primary:", "    allowLocalNetworks: true\n    generation:\n      maxResponseTokens: 8192\n  reranker-primary:")
	_, err := Load(context.Background(), config)
	if err == nil || !strings.Contains(err.Error(), "RAG_PROVIDER_CONFIG_GENERATION_POLICY") {
		t.Fatalf("Load() error = %v", err)
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

func TestLoadProviderSetRejectsModelManifestDigestMismatch(t *testing.T) {
	config := writeProviderFixture(t)
	modelPath := filepath.Join(filepath.Dir(config), "models", "generator.json")
	data, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest ragcontract.ModelManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.ModelID = "drifted-model"
	writeJSONFileWithoutNormalizing(t, modelPath, manifest)
	_, err = Load(context.Background(), config)
	if err == nil || !strings.Contains(err.Error(), "RAG_MODEL_MANIFEST_DIGEST") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestLoadProviderSetRejectsPromptSchemaMissing(t *testing.T) {
	config := writeProviderFixture(t)
	manifestPath := filepath.Join(filepath.Dir(config), "prompts", "summary.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest ragcontract.PromptManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	manifest.OutputSchema = "missing/v1"
	manifest.Digest, err = promptManifestContentDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	writeJSONFileWithoutNormalizing(t, manifestPath, manifest)
	_, err = Load(context.Background(), config)
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
		text := "Return the requested JSON."
		if err := os.WriteFile(filepath.Join(prompts, prompt.id+".txt"), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
		templateDigest, err := ragcontract.Digest(text)
		if err != nil {
			t.Fatal(err)
		}
		manifest := ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema}, PromptID: prompt.id, TemplateDigest: templateDigest, InputSchema: "text", OutputSchema: prompt.schema}
		manifest.Digest, err = promptManifestContentDigest(manifest)
		if err != nil {
			t.Fatal(err)
		}
		writeJSONFile(t, filepath.Join(prompts, prompt.id+".json"), manifest)
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

func writeJSONFileWithoutNormalizing(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeJSONFile(t *testing.T, path string, value any) {
	t.Helper()
	if model, ok := value.(ragcontract.ModelManifest); ok {
		digest, err := modelManifestContentDigest(model)
		if err != nil {
			t.Fatal(err)
		}
		model.Digest = digest
		value = model
	}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func repeat(value string, count int) string { return strings.Repeat(value, count) }
