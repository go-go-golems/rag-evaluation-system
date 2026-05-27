package embedding

import (
	"context"
	"strings"
	"testing"
)

func TestResolveInferenceSettingsDirectOllama(t *testing.T) {
	settings, closeFn, effective, err := ResolveInferenceSettings(context.Background(), ProviderConfig{
		Type:       "ollama",
		Engine:     "nomic-embed-text",
		Dimensions: 768,
		BaseURL:    "http://localhost:11434",
	})
	if err != nil {
		t.Fatalf("resolve inference settings: %v", err)
	}
	if closeFn != nil {
		t.Fatal("direct settings should not return a close function")
	}
	if effective != "direct" {
		t.Fatalf("expected direct effective profile, got %q", effective)
	}
	if settings.Embeddings == nil || settings.Embeddings.Type != "ollama" || settings.Embeddings.Engine != "nomic-embed-text" || settings.Embeddings.Dimensions != 768 {
		t.Fatalf("unexpected embeddings settings: %#v", settings.Embeddings)
	}
	if got := settings.API.BaseUrls["ollama-base-url"]; got != "http://localhost:11434" {
		t.Fatalf("expected ollama base URL, got %q", got)
	}
}

func TestResolveProviderValidatesMissingOpenAIKey(t *testing.T) {
	_, err := ResolveProvider(context.Background(), ProviderConfig{
		Type:       "openai",
		Engine:     "text-embedding-3-small",
		Dimensions: 1536,
	})
	if err == nil {
		t.Fatal("expected missing OpenAI key to fail")
	}
	if !strings.Contains(err.Error(), "openai-api-key") {
		t.Fatalf("expected openai-api-key error, got %v", err)
	}
}

func TestResolveProviderDirectOllamaDoesNotNeedNetworkAtConstruction(t *testing.T) {
	resolved, err := ResolveProvider(context.Background(), ProviderConfig{
		Type:       "ollama",
		Engine:     "nomic-embed-text",
		Dimensions: 768,
	})
	if err != nil {
		t.Fatalf("resolve provider: %v", err)
	}
	if resolved.Provider == nil {
		t.Fatal("expected provider")
	}
	if resolved.ProviderType != "ollama" {
		t.Fatalf("expected ollama provider type, got %q", resolved.ProviderType)
	}
	if resolved.Model.Name != "nomic-embed-text" || resolved.Model.Dimensions != 768 {
		t.Fatalf("unexpected model: %#v", resolved.Model)
	}
}

func TestDefaultPinocchioProfilesPath(t *testing.T) {
	path := DefaultPinocchioProfilesPath()
	if !strings.HasSuffix(path, ".config/pinocchio/profiles.yaml") {
		t.Fatalf("unexpected profile path %q", path)
	}
}
