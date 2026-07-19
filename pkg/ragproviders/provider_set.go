package ragproviders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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

	"github.com/go-go-golems/geppetto/pkg/engineprofiles"
)

type EffectiveProviderIdentity struct {
	Role                string `json:"role"`
	ProfileSlug         string `json:"profileSlug,omitempty"`
	ProfileSource       string `json:"profileSource,omitempty"`
	ModelManifestDigest string `json:"modelManifestDigest"`
	ModelID             string `json:"modelId"`
	SettingsFingerprint string `json:"settingsFingerprint"`
	ConcurrencyLimit    int    `json:"concurrencyLimit"`
}

type CapabilityDescriptor struct {
	SchemaVersion               string                      `json:"schemaVersion"`
	ProfileID                   string                      `json:"profileId"`
	FixtureProviders            bool                        `json:"fixtureProviders"`
	Capabilities                []string                    `json:"capabilities"`
	ModelManifestDigests        []string                    `json:"modelManifestDigests"`
	PromptManifestDigests       []string                    `json:"promptManifestDigests"`
	EffectiveProviderIdentities []EffectiveProviderIdentity `json:"effectiveProviderIdentities,omitempty"`
}

type ProviderSet struct {
	ProfileID                   string
	Manifests                   *FileManifestRegistry
	Schemas                     *FileSchemaRegistry
	Generator                   ragoperators.TextGenerator
	Embedder                    ragoperators.Embedder
	Reranker                    ragoperators.Reranker
	Cache                       ragoperators.Cache
	GenerationConcurrency       int
	QueryCheckpoints            ragengine.QueryCheckpointStore
	PreparedStore               ragengine.PreparedCorpusStore
	EffectiveProviderIdentities map[string]EffectiveProviderIdentity

	profileChains map[string]*engineprofiles.ChainedRegistry
	closers       []io.Closer
	closeOnce     sync.Once
	closeErr      error
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
	set := &ProviderSet{ProfileID: cfg.ProfileID, Manifests: manifests, Schemas: schemas, Cache: cache, EffectiveProviderIdentities: map[string]EffectiveProviderIdentity{}, profileChains: map[string]*engineprofiles.ChainedRegistry{}, closers: []io.Closer{cache}}
	if cfg.Checkpoints.Kind != "" {
		set.QueryCheckpoints, err = ragengine.NewFileQueryCheckpointStore(filepath.Join(cfg.Checkpoints.Directory, "queries"))
		if err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_CHECKPOINTS: %w", err)
		}
		set.PreparedStore, err = ragengine.NewFilePreparedCorpusStore(filepath.Join(cfg.Checkpoints.Directory, "prepared"))
		if err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_CHECKPOINTS: %w", err)
		}
	}
	if spec, ok := cfg.Providers["embedding-primary"]; ok {
		model, err := manifests.Model(spec.ModelManifest)
		if err != nil {
			return nil, err
		}
		var provider embeddings.Provider
		if spec.Profile != "" {
			provider, err = newEmbeddingProviderFromProfile(ctx, spec, model, set)
			if err != nil {
				return nil, err
			}
		} else {
			endpoint, err := resolveEnvRef(spec.EndpointRef, true)
			if err != nil {
				return nil, err
			}
			if err := validateEndpoint(endpoint, spec); err != nil {
				return nil, fmt.Errorf("RAG_PROVIDER_EMBEDDING_ENDPOINT: %w", err)
			}
			provider, err = newEmbeddingProvider(endpoint, model, spec)
			if err != nil {
				return nil, err
			}
		}
		embedder, err := geppettoadapter.NewEmbedder(provider, model.ModelID, model.Dimensions)
		if err != nil {
			return nil, err
		}
		set.Embedder, err = newLimitedEmbedder(embedder, providerConcurrencyLimit(spec))
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
		var provider rankprovider.Provider
		if spec.Profile != "" {
			provider, err = newRerankProviderFromProfile(ctx, spec, model, set)
			if err != nil {
				return nil, err
			}
		} else {
			endpoint, err := resolveEnvRef(spec.EndpointRef, true)
			if err != nil {
				return nil, err
			}
			if err := validateEndpoint(endpoint, spec); err != nil {
				return nil, fmt.Errorf("RAG_PROVIDER_RERANK_ENDPOINT: %w", err)
			}
			provider, err = newRerankProvider(endpoint, model, spec)
			if err != nil {
				return nil, err
			}
		}
		reranker, err := geppettoadapter.NewReranker(provider)
		if err != nil {
			return nil, err
		}
		set.Reranker, err = newLimitedReranker(reranker, providerConcurrencyLimit(spec))
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
		var generatorSettings *settings.InferenceSettings
		if spec.Profile != "" {
			generatorSettings, err = newGeneratorSettingsFromProfile(ctx, spec, model, set)
			if err != nil {
				return nil, err
			}
		} else {
			endpoint, err := resolveEnvRef(spec.EndpointRef, true)
			if err != nil {
				return nil, err
			}
			if err := validateEndpoint(endpoint, spec); err != nil {
				return nil, fmt.Errorf("RAG_PROVIDER_GENERATOR_ENDPOINT: %w", err)
			}
			generatorSettings, err = newGeneratorSettings(endpoint, model, spec)
			if err != nil {
				return nil, err
			}
		}
		generator, err := geppettoadapter.NewGenerator(generatorSettings, manifests, schemas)
		if err != nil {
			return nil, err
		}
		set.GenerationConcurrency = providerConcurrencyLimit(spec)
		set.Generator, err = newLimitedGenerator(generator, set.GenerationConcurrency)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("RAG_PROVIDER_GENERATOR_REQUIRED")
	}
	_ = ctx
	return set, nil
}
func (p *ProviderSet) CapabilityDescriptor() CapabilityDescriptor {
	if p == nil {
		return CapabilityDescriptor{SchemaVersion: "rag-provider-capabilities/v1"}
	}
	descriptor := CapabilityDescriptor{SchemaVersion: "rag-provider-capabilities/v1", ProfileID: p.ProfileID, Capabilities: []string{"schema-validator", "persistent-cache"}}
	if p.Generator != nil {
		descriptor.Capabilities = append(descriptor.Capabilities, "generator")
	}
	if p.Embedder != nil {
		descriptor.Capabilities = append(descriptor.Capabilities, "embedder")
	}
	if p.Reranker != nil {
		descriptor.Capabilities = append(descriptor.Capabilities, "reranker")
	}
	seenModels, seenPrompts := map[string]bool{}, map[string]bool{}
	if p.Manifests != nil {
		for _, manifest := range p.Manifests.models {
			if !seenModels[manifest.Digest] {
				seenModels[manifest.Digest] = true
				descriptor.ModelManifestDigests = append(descriptor.ModelManifestDigests, manifest.Digest)
			}
		}
		for _, manifest := range p.Manifests.prompts {
			if !seenPrompts[manifest.Digest] {
				seenPrompts[manifest.Digest] = true
				descriptor.PromptManifestDigests = append(descriptor.PromptManifestDigests, manifest.Digest)
			}
		}
	}
	for _, identity := range p.EffectiveProviderIdentities {
		descriptor.EffectiveProviderIdentities = append(descriptor.EffectiveProviderIdentities, identity)
	}
	sort.Slice(descriptor.EffectiveProviderIdentities, func(i, j int) bool {
		return descriptor.EffectiveProviderIdentities[i].Role < descriptor.EffectiveProviderIdentities[j].Role
	})
	sort.Strings(descriptor.Capabilities)
	sort.Strings(descriptor.ModelManifestDigests)
	sort.Strings(descriptor.PromptManifestDigests)
	return descriptor
}

func (p *ProviderSet) EngineOptions() ragengine.Options {
	generatorFingerprint := ""
	if p != nil {
		generatorFingerprint = p.EffectiveProviderIdentities["generator-primary"].SettingsFingerprint
	}
	return ragengine.Options{Manifests: p.Manifests, Schemas: p.Schemas, Generator: p.Generator, Embedder: p.Embedder, Reranker: p.Reranker, Cache: p.Cache, GenerationConcurrency: p.GenerationConcurrency, GenerationSettingsFingerprint: generatorFingerprint, GeneratorFingerprint: generatorFingerprint, RerankerFingerprint: p.EffectiveProviderIdentities["reranker-primary"].SettingsFingerprint, QueryCheckpoints: p.QueryCheckpoints, PreparedStore: p.PreparedStore, EmbeddingFingerprint: p.EffectiveProviderIdentities["embedding-primary"].SettingsFingerprint}
}

// CheckExecution verifies that this host can satisfy every provider-backed
// operator in a canonical execution before the engine performs provider work.
func (p *ProviderSet) CheckExecution(execution ragcontract.PipelineExecution) error {
	if p == nil || p.Manifests == nil || p.Schemas == nil {
		return fmt.Errorf("RAG_PROVIDER_SET_UNAVAILABLE")
	}
	for _, node := range execution.Pipeline.Nodes {
		var config struct {
			Model, Prompt, OutputSchema, Truncation, Tokenization string
			Dimensions                                            int
		}
		if err := json.Unmarshal(node.Config, &config); err != nil {
			return fmt.Errorf("RAG_PROVIDER_EXECUTION_CONFIG")
		}
		switch node.Operator.Kind {
		case "representations.structured-summary", "representations.synthetic-questions":
			if p.Generator == nil || p.Cache == nil {
				return fmt.Errorf("RAG_PROVIDER_GENERATOR_REQUIRED")
			}
			if _, err := p.Manifests.Model(config.Model); err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_MODEL")
			}
			prompt, err := p.Manifests.Prompt(config.Prompt)
			if err != nil || (config.OutputSchema != "" && prompt.OutputSchema != config.OutputSchema) {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_PROMPT")
			}
			if _, err := p.Schemas.Raw(prompt.OutputSchema); err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_SCHEMA")
			}
		case "embed.model":
			if p.Embedder == nil {
				return fmt.Errorf("RAG_PROVIDER_EMBEDDING_REQUIRED")
			}
			model, err := p.Manifests.Model(config.Model)
			if err != nil || (config.Dimensions > 0 && model.Dimensions != config.Dimensions) {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_MODEL")
			}
		case "rerank.cross-encoder":
			if p.Reranker == nil {
				return fmt.Errorf("RAG_PROVIDER_RERANKER_REQUIRED")
			}
			model, err := p.Manifests.Model(config.Model)
			if err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_MODEL")
			}
			if config.Truncation == "" {
				config.Truncation = model.Truncation
			}
			if config.Tokenization == "" {
				config.Tokenization = model.Tokenization
			}
			if model.Truncation != config.Truncation || model.Tokenization != config.Tokenization {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_MODEL")
			}
		case "generate.answer":
			if p.Generator == nil {
				return fmt.Errorf("RAG_PROVIDER_GENERATOR_REQUIRED")
			}
			if _, err := p.Manifests.Model(config.Model); err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_MODEL")
			}
			prompt, err := p.Manifests.Prompt(config.Prompt)
			if err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_PROMPT")
			}
			if _, err := p.Schemas.Raw(prompt.OutputSchema); err != nil {
				return fmt.Errorf("RAG_PROVIDER_EXECUTION_SCHEMA")
			}
		}
	}
	return nil
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

// newEmbeddingProviderFromProfile resolves embedding settings from a Geppetto
// engine profile registry, mirroring the generator profile path.
func newEmbeddingProviderFromProfile(ctx context.Context, spec ProviderSpec, model ragcontract.ModelManifest, set *ProviderSet) (embeddings.Provider, error) {
	ss, err := resolveProfileSettings(ctx, spec, set)
	if err != nil {
		return nil, err
	}
	if ss.Embeddings == nil || ss.Embeddings.Type == "" {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_NO_EMBEDDINGS")
	}
	if err := validateProfileModelIdentity("embedding", model, ss); err != nil {
		return nil, err
	}
	if err := set.recordProfileIdentity("embedding-primary", spec, model, ss); err != nil {
		return nil, err
	}
	return embeddings.NewSettingsFactoryFromInferenceSettings(ss).NewProvider()
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

// newRerankProviderFromProfile resolves reranker settings from a Geppetto
// engine profile registry, mirroring the generator profile path.
func newRerankProviderFromProfile(ctx context.Context, spec ProviderSpec, model ragcontract.ModelManifest, set *ProviderSet) (rankprovider.Provider, error) {
	ss, err := resolveProfileSettings(ctx, spec, set)
	if err != nil {
		return nil, err
	}
	if ss.Rerank == nil || ss.Rerank.Type == "" {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_NO_RERANK")
	}
	if err := validateProfileModelIdentity("reranker", model, ss); err != nil {
		return nil, err
	}
	if err := set.recordProfileIdentity("reranker-primary", spec, model, ss); err != nil {
		return nil, err
	}
	factory, err := geppettorerank.NewSettingsFactoryFromInferenceSettings(ss)
	if err != nil {
		return nil, err
	}
	return factory.NewProvider()
}

// resolveProfileSettings is the shared profile resolution logic used by the
// generator, embedding, and reranker profile paths. It loads the Geppetto
// engine profile registry (defaulting to ~/.config/pinocchio/profiles.yaml),
// resolves the named profile, and merges it onto a base InferenceSettings so
// that default Client, API, and other infrastructure fields are initialized.
func resolveProfileSettings(ctx context.Context, spec ProviderSpec, set *ProviderSet) (*settings.InferenceSettings, error) {
	chain, err := set.profileRegistry(ctx, spec.ProfileRegistries)
	if err != nil {
		return nil, err
	}
	resolved, err := chain.ResolveEngineProfile(ctx, engineprofiles.ResolveInput{
		EngineProfileSlug: engineprofiles.MustEngineProfileSlug(spec.Profile),
	})
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_RESOLVE: %w", err)
	}
	if resolved.InferenceSettings == nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_NO_SETTINGS")
	}
	base, err := settings.NewInferenceSettings()
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_BASE: %w", err)
	}
	merged, err := engineprofiles.MergeInferenceSettings(base, resolved.InferenceSettings)
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_MERGE: %w", err)
	}
	return merged, nil
}

// profileRegistry returns one shared registry chain per configured source list.
// Empty source lists deliberately fall back to Pinocchio's standard per-user
// profile registry so all provider roles resolve the same host configuration.
func (p *ProviderSet) profileRegistry(ctx context.Context, sources string) (*engineprofiles.ChainedRegistry, error) {
	if sources == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_CONFIG_DIR: %w", err)
		}
		defaultPath := filepath.Join(configDir, "pinocchio", "profiles.yaml")
		if _, err := os.Stat(defaultPath); err != nil {
			return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_REGISTRY_NOT_FOUND: %s", defaultPath)
		}
		sources = defaultPath
	}
	if chain := p.profileChains[sources]; chain != nil {
		return chain, nil
	}
	entries, err := engineprofiles.ParseEngineProfileRegistrySourceEntries(sources)
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_SOURCE: %w", err)
	}
	specs, err := engineprofiles.ParseRegistrySourceSpecs(entries)
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_SPEC: %w", err)
	}
	chain, err := engineprofiles.NewChainedRegistryFromSourceSpecs(ctx, specs)
	if err != nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_CHAIN: %w", err)
	}
	p.profileChains[sources] = chain
	p.closers = append(p.closers, chain)
	return chain, nil
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

// newGeneratorSettingsFromProfile resolves InferenceSettings from a Geppetto
// engine profile registry (e.g. ~/.config/pinocchio/profiles.yaml). This uses
// the same profile stack resolution that Pinocchio and other Geppetto hosts
// use, so credentials, base URLs, model_info, and cost rates all come from the
// profile YAML rather than being manually constructed.
func newGeneratorSettingsFromProfile(ctx context.Context, spec ProviderSpec, model ragcontract.ModelManifest, set *ProviderSet) (*settings.InferenceSettings, error) {
	merged, err := resolveProfileSettings(ctx, spec, set)
	if err != nil {
		return nil, err
	}
	if merged.Chat == nil {
		return nil, fmt.Errorf("RAG_PROVIDER_PROFILE_NO_CHAT")
	}
	if err := validateProfileModelIdentity("generator", model, merged); err != nil {
		return nil, err
	}
	if err := set.recordProfileIdentity("generator-primary", spec, model, merged); err != nil {
		return nil, err
	}
	merged.Chat.Stream = true
	// Reasoning models (GLM-5.2, GPT-5, etc.) burn tokens on reasoning before
	// producing visible output. Set a generous max_response_tokens so the
	// model doesn't hit the token limit mid-reasoning and return empty content.
	if merged.Chat.MaxResponseTokens == nil {
		maxTokens := 8192
		merged.Chat.MaxResponseTokens = &maxTokens
	}
	return merged, nil
}
func validateProfileModelIdentity(role string, model ragcontract.ModelManifest, settings *settings.InferenceSettings) error {
	if settings == nil {
		return fmt.Errorf("RAG_PROVIDER_PROFILE_NO_SETTINGS")
	}
	actual := ""
	switch role {
	case "generator":
		if settings.Chat != nil && settings.Chat.Engine != nil {
			actual = *settings.Chat.Engine
		}
	case "embedding":
		if settings.Embeddings != nil {
			actual = settings.Embeddings.Engine
		}
	case "reranker":
		if settings.Rerank != nil {
			actual = settings.Rerank.Engine
		}
	default:
		return fmt.Errorf("RAG_PROVIDER_PROFILE_ROLE")
	}
	if actual == "" || actual != model.ModelID {
		return fmt.Errorf("RAG_PROVIDER_PROFILE_MODEL_MISMATCH")
	}
	return nil
}

// recordProfileIdentity captures only non-secret, inference-affecting profile
// values. It deliberately excludes API keys, endpoint URLs, OAuth state, and
// any other credential material from custody artifacts and cache identity.
func (p *ProviderSet) recordProfileIdentity(role string, spec ProviderSpec, model ragcontract.ModelManifest, ss *settings.InferenceSettings) error {
	apiType, chatEngine := "", ""
	var temperature *float64
	var maxResponseTokens *int
	if ss.Chat != nil {
		if ss.Chat.ApiType != nil {
			apiType = string(*ss.Chat.ApiType)
		}
		if ss.Chat.Engine != nil {
			chatEngine = *ss.Chat.Engine
		}
		temperature = ss.Chat.Temperature
		maxResponseTokens = ss.Chat.MaxResponseTokens
	}
	fingerprint, err := ragcontract.Digest(struct {
		SchemaVersion     string
		Role              string
		APIType           string
		ChatEngine        string
		EmbeddingType     string
		EmbeddingEngine   string
		EmbeddingDims     int
		RerankType        string
		RerankEngine      string
		Inference         any
		ModelInfo         any
		Temperature       *float64
		MaxResponseTokens *int
	}{
		SchemaVersion: "rag-effective-provider-settings/v1",
		Role:          role, APIType: apiType, ChatEngine: chatEngine,
		Inference: ss.Inference, ModelInfo: ss.ModelInfo,
		Temperature: temperature, MaxResponseTokens: maxResponseTokens,
	})
	if err != nil {
		return fmt.Errorf("RAG_PROVIDER_PROFILE_FINGERPRINT: %w", err)
	}
	if ss.Embeddings != nil {
		fingerprint, err = ragcontract.Digest(struct {
			Fingerprint string
			Type        string
			Engine      string
			Dimensions  int
		}{fingerprint, ss.Embeddings.Type, ss.Embeddings.Engine, ss.Embeddings.Dimensions})
		if err != nil {
			return fmt.Errorf("RAG_PROVIDER_PROFILE_FINGERPRINT: %w", err)
		}
	}
	if ss.Rerank != nil {
		fingerprint, err = ragcontract.Digest(struct {
			Fingerprint string
			Type        string
			Engine      string
		}{fingerprint, ss.Rerank.Type, ss.Rerank.Engine})
		if err != nil {
			return fmt.Errorf("RAG_PROVIDER_PROFILE_FINGERPRINT: %w", err)
		}
	}
	source := "default"
	if spec.ProfileRegistries != "" {
		source = "configured"
	}
	p.EffectiveProviderIdentities[role] = EffectiveProviderIdentity{
		Role: role, ProfileSlug: spec.Profile, ProfileSource: source,
		ModelManifestDigest: model.Digest, ModelID: model.ModelID,
		SettingsFingerprint: fingerprint, ConcurrencyLimit: providerConcurrencyLimit(spec),
	}
	return nil
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
