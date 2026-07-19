package ragproviders

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sync"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	rankprovider "github.com/go-go-golems/geppetto/pkg/rerank"
	rankcore "github.com/go-go-golems/geppetto/pkg/rerank/config"
	geppettorerank "github.com/go-go-golems/geppetto/pkg/rerank/factory"
	"github.com/go-go-golems/geppetto/pkg/security"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/types"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	geppettoadapter "github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders/geppetto"
)

type ProviderSet struct {
	ProfileID string
	Manifests *FileManifestRegistry
	Schemas   *FileSchemaRegistry
	Generator ragoperators.TextGenerator
	Embedder  ragoperators.Embedder
	Reranker  ragoperators.Reranker
	Cache     ragoperators.Cache

	closers   []io.Closer
	closeOnce sync.Once
	closeErr  error
}

func Load(ctx context.Context, path string) (*ProviderSet, error) {
	cfg, _, err := loadConfig(path)
	if err != nil {
		return nil, err
	}
	manifests, err := LoadManifestRegistry(cfg.Manifests.ModelsDir, cfg.Manifests.PromptsDir)
	if err != nil {
		return nil, err
	}
	schemas, err := LoadSchemaRegistry(cfg.Schemas.Directory)
	if err != nil {
		return nil, err
	}
	if err := validateProviderConfiguration(cfg.Providers, manifests, schemas); err != nil {
		return nil, err
	}
	cache, err := NewDiskCache(cfg.Cache)
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_CACHE: %w", err)
	}
	set := &ProviderSet{ProfileID: cfg.ProfileID, Manifests: manifests, Schemas: schemas, Cache: cache, closers: []io.Closer{cache}}
	if spec, ok := cfg.Providers["embedding-primary"]; ok {
		model, err := manifests.Model(spec.ModelManifest)
		if err != nil {
			return nil, err
		}
		endpoint, err := resolveEnvRef(spec.EndpointRef, true)
		if err != nil {
			return nil, err
		}
		if err := validateEndpoint(endpoint, spec); err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_EMBEDDING_ENDPOINT: %w", err)
		}
		provider, err := newEmbeddingProvider(endpoint, model, spec)
		if err != nil {
			return nil, err
		}
		set.Embedder, err = geppettoadapter.NewEmbedder(provider, model.ModelID, model.Dimensions)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("RAG_PROVIDER_EMBEDDING_REQUIRED")
	}
	if spec, ok := cfg.Providers["reranker-primary"]; ok {
		model, err := manifests.Model(spec.ModelManifest)
		if err != nil {
			return nil, err
		}
		endpoint, err := resolveEnvRef(spec.EndpointRef, true)
		if err != nil {
			return nil, err
		}
		if err := validateEndpoint(endpoint, spec); err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_RERANK_ENDPOINT: %w", err)
		}
		provider, err := newRerankProvider(endpoint, model, spec)
		if err != nil {
			return nil, err
		}
		set.Reranker, err = geppettoadapter.NewReranker(provider)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("RAG_PROVIDER_RERANKER_REQUIRED")
	}
	if spec, ok := cfg.Providers["generator-primary"]; ok {
		model, err := manifests.Model(spec.ModelManifest)
		if err != nil {
			return nil, err
		}
		endpoint, err := resolveEnvRef(spec.EndpointRef, true)
		if err != nil {
			return nil, err
		}
		if err := validateEndpoint(endpoint, spec); err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_GENERATOR_ENDPOINT: %w", err)
		}
		generatorSettings, err := newGeneratorSettings(endpoint, model, spec)
		if err != nil {
			return nil, err
		}
		set.Generator, err = geppettoadapter.NewGenerator(generatorSettings, manifests, schemas)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("RAG_PROVIDER_GENERATOR_REQUIRED")
	}
	_ = ctx
	return set, nil
}
func (p *ProviderSet) EngineOptions() ragengine.Options {
	return ragengine.Options{Manifests: p.Manifests, Schemas: p.Schemas, Generator: p.Generator, Embedder: p.Embedder, Reranker: p.Reranker, Cache: p.Cache}
}
func (p *ProviderSet) Close() error {
	if p == nil {
		return nil
	}
	p.closeOnce.Do(func() {
		var closeErrors []error
		for i := len(p.closers) - 1; i >= 0; i-- {
			if err := p.closers[i].Close(); err != nil {
				closeErrors = append(closeErrors, err)
			}
		}
		p.closeErr = errors.Join(closeErrors...)
	})
	return p.closeErr
}

func validateProviderConfiguration(providers map[string]ProviderSpec, manifests *FileManifestRegistry, schemas *FileSchemaRegistry) error {
	expected := map[string]string{
		"embedding-primary": "geppetto-embedding/v1",
		"generator-primary": "geppetto-generation/v1",
		"reranker-primary":  "geppetto-reranker/v1",
	}
	for name, spec := range providers {
		kind, known := expected[name]
		if !known || spec.Kind != kind {
			return fmt.Errorf("RAG_PROVIDER_KIND_UNSUPPORTED")
		}
		model, err := manifests.Model(spec.ModelManifest)
		if err != nil {
			return err
		}
		if model.ProviderAdapterVersion != kind {
			return fmt.Errorf("RAG_PROVIDER_MODEL_INCOMPATIBLE")
		}
	}
	for name := range expected {
		if _, ok := providers[name]; !ok {
			return fmt.Errorf("RAG_PROVIDER_REQUIRED")
		}
	}
	seen := map[string]bool{}
	for _, prompt := range manifests.prompts {
		if seen[prompt.PromptID] {
			continue
		}
		seen[prompt.PromptID] = true
		if _, err := schemas.Raw(prompt.OutputSchema); err != nil {
			return fmt.Errorf("RAG_PROVIDER_PROMPT_SCHEMA_MISSING")
		}
	}
	return nil
}
func newEmbeddingProvider(endpoint string, model ragcontract.ModelManifest, spec ProviderSpec) (embeddings.Provider, error) {
	in, err := settings.NewInferenceSettings()
	if err != nil {
		return nil, err
	}
	in.Embeddings.Type = "ollama"
	in.Embeddings.Engine = model.ModelID
	in.Embeddings.Dimensions = model.Dimensions
	in.Embeddings.BaseURLs = map[string]string{"ollama-base-url": endpoint}
	in.API.BaseUrls["ollama-base-url"] = endpoint
	in.API.AllowHTTP["ollama"] = spec.AllowHTTP
	in.API.AllowLocalNetworks["ollama"] = spec.AllowLocalNetworks
	return embeddings.NewSettingsFactoryFromInferenceSettings(in).NewProvider()
}
func newRerankProvider(endpoint string, model ragcontract.ModelManifest, spec ProviderSpec) (rankprovider.Provider, error) {
	in, err := settings.NewInferenceSettings()
	if err != nil {
		return nil, err
	}
	in.Rerank = &rankcore.RerankConfig{Type: "llamacpp", Engine: model.ModelID}
	in.API.BaseUrls["rerank-base-url"] = endpoint
	in.API.AllowHTTP["rerank"] = spec.AllowHTTP
	in.API.AllowLocalNetworks["rerank"] = spec.AllowLocalNetworks
	factory, err := geppettorerank.NewSettingsFactoryFromInferenceSettings(in)
	if err != nil {
		return nil, err
	}
	return factory.NewProvider()
}
func newGeneratorSettings(endpoint string, model ragcontract.ModelManifest, spec ProviderSpec) (*settings.InferenceSettings, error) {
	in, err := settings.NewInferenceSettings()
	if err != nil {
		return nil, err
	}
	apiType := types.ApiTypeOpenAI
	in.Chat.ApiType = &apiType
	in.Chat.Engine = stringPtr(model.ModelID)
	in.Chat.Stream = true
	in.API.BaseUrls["openai-base-url"] = endpoint
	in.API.AllowHTTP["openai"] = spec.AllowHTTP
	in.API.AllowLocalNetworks["openai"] = spec.AllowLocalNetworks
	key := "ollama"
	if spec.CredentialRef != "" {
		key, err = resolveEnvRef(spec.CredentialRef, true)
		if err != nil {
			return nil, err
		}
	}
	in.API.APIKeys["openai-api-key"] = key
	return in, nil
}
func stringPtr(value string) *string { return &value }
func validateEndpoint(raw string, spec ProviderSpec) error {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.User != nil {
		return fmt.Errorf("RAG_PROVIDER_ENDPOINT_POLICY")
	}
	if err := security.ValidateOutboundURL(raw, security.OutboundURLOptions{AllowHTTP: spec.AllowHTTP, AllowLocalNetworks: spec.AllowLocalNetworks}); err != nil {
		return fmt.Errorf("RAG_PROVIDER_ENDPOINT_POLICY")
	}
	return nil
}
