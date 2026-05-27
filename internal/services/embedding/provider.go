package embedding

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	embeddingconfig "github.com/go-go-golems/geppetto/pkg/embeddings/config"
	profiles "github.com/go-go-golems/geppetto/pkg/engineprofiles"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
)

// ProviderConfig describes how to resolve an embedding provider. It supports
// direct embedding settings and profile-backed settings resolved from Pinocchio
// profile registries via Geppetto engineprofiles.
type ProviderConfig struct {
	ProfileRegistries []string
	Profile           string
	BaseProfile       string

	Type       string
	Engine     string
	Dimensions int
	APIKey     string
	BaseURL    string

	CacheType       string
	CacheDirectory  string
	CacheMaxEntries int
	CacheMaxSize    int64
}

type ResolvedProvider struct {
	Provider         embeddings.Provider
	EffectiveProfile string
	ProviderType     string
	Model            embeddings.EmbeddingModel
	Close            func() error
}

func ResolveProvider(ctx context.Context, cfg ProviderConfig) (*ResolvedProvider, error) {
	inferenceSettings, closeFn, effectiveProfile, err := ResolveInferenceSettings(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := embeddings.ValidateInferenceSettingsForEmbeddings(inferenceSettings); err != nil {
		if closeFn != nil {
			_ = closeFn()
		}
		return nil, err
	}

	provider, err := embeddings.NewSettingsFactoryFromInferenceSettings(inferenceSettings).NewProvider()
	if err != nil {
		if closeFn != nil {
			_ = closeFn()
		}
		return nil, err
	}
	model := provider.GetModel()
	return &ResolvedProvider{
		Provider:         provider,
		EffectiveProfile: effectiveProfile,
		ProviderType:     inferenceSettings.Embeddings.Type,
		Model:            model,
		Close:            closeFn,
	}, nil
}

func ResolveInferenceSettings(ctx context.Context, cfg ProviderConfig) (*settings.InferenceSettings, func() error, string, error) {
	if strings.TrimSpace(cfg.Profile) != "" || strings.TrimSpace(cfg.BaseProfile) != "" {
		return resolveProfileBackedSettings(ctx, cfg)
	}
	return directInferenceSettings(cfg), nil, "direct", nil
}

func resolveProfileBackedSettings(ctx context.Context, cfg ProviderConfig) (*settings.InferenceSettings, func() error, string, error) {
	registries := cfg.ProfileRegistries
	if len(registries) == 0 {
		registries = []string{DefaultPinocchioProfilesPath()}
	}

	specs, err := profiles.ParseRegistrySourceSpecs(registries)
	if err != nil {
		return nil, nil, "", err
	}
	chain, err := profiles.NewChainedRegistryFromSourceSpecs(ctx, specs)
	if err != nil {
		return nil, nil, "", err
	}

	if strings.TrimSpace(cfg.Profile) != "" {
		profileSlug, err := profiles.ParseEngineProfileSlug(cfg.Profile)
		if err != nil {
			_ = chain.Close()
			return nil, nil, "", err
		}
		resolved, err := chain.ResolveEngineProfile(ctx, profiles.ResolveInput{EngineProfileSlug: profileSlug})
		if err != nil {
			_ = chain.Close()
			return nil, nil, "", err
		}
		return resolved.InferenceSettings, chain.Close, cfg.Profile, nil
	}

	baseSlug, err := profiles.ParseEngineProfileSlug(cfg.BaseProfile)
	if err != nil {
		_ = chain.Close()
		return nil, nil, "", err
	}
	baseResolved, err := chain.ResolveEngineProfile(ctx, profiles.ResolveInput{EngineProfileSlug: baseSlug})
	if err != nil {
		_ = chain.Close()
		return nil, nil, "", err
	}

	merged, err := profiles.MergeInferenceSettings(baseResolved.InferenceSettings, directInferenceSettings(cfg))
	if err != nil {
		_ = chain.Close()
		return nil, nil, "", err
	}
	return merged, chain.Close, fmt.Sprintf("%s + embeddings(%s/%s)", cfg.BaseProfile, cfg.Type, cfg.Engine), nil
}

func directInferenceSettings(cfg ProviderConfig) *settings.InferenceSettings {
	api := settings.NewAPISettings()
	providerType := strings.TrimSpace(cfg.Type)
	if providerType == "" {
		providerType = "ollama"
	}
	if strings.TrimSpace(cfg.APIKey) != "" {
		api.APIKeys[providerType+"-api-key"] = cfg.APIKey
		if providerType == "openai" {
			api.APIKeys["openai-api-key"] = cfg.APIKey
		}
	}
	if strings.TrimSpace(cfg.BaseURL) != "" {
		api.BaseUrls[providerType+"-base-url"] = cfg.BaseURL
	}

	return &settings.InferenceSettings{
		API: api,
		Embeddings: &embeddingconfig.EmbeddingsConfig{
			Type:            providerType,
			Engine:          cfg.Engine,
			Dimensions:      cfg.Dimensions,
			CacheType:       cfg.CacheType,
			CacheDirectory:  cfg.CacheDirectory,
			CacheMaxEntries: cfg.CacheMaxEntries,
			CacheMaxSize:    cfg.CacheMaxSize,
			APIKeys:         api.APIKeys,
			BaseURLs:        api.BaseUrls,
		},
	}
}

func DefaultPinocchioProfilesPath() string {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return filepath.Join(".config", "pinocchio", "profiles.yaml")
	}
	return filepath.Join(home, ".config", "pinocchio", "profiles.yaml")
}
