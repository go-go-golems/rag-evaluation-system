package ragproviders

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const HostConfigSchemaVersion = "rag-provider-host-config/v1"

type HostConfig struct {
	SchemaVersion string                  `yaml:"schemaVersion"`
	ProfileID     string                  `yaml:"profileId"`
	Manifests     ManifestPaths           `yaml:"manifests"`
	Schemas       SchemaPaths             `yaml:"schemas"`
	Cache         CacheConfig             `yaml:"cache"`
	Checkpoints   CheckpointConfig        `yaml:"checkpoints,omitempty"`
	Providers     map[string]ProviderSpec `yaml:"providers"`
}
type ManifestPaths struct {
	ModelsDir  string `yaml:"modelsDir"`
	PromptsDir string `yaml:"promptsDir"`
}
type SchemaPaths struct {
	Directory string `yaml:"directory"`
}
type CacheConfig struct {
	Kind          string `yaml:"kind"`
	Directory     string `yaml:"directory"`
	MaxEntryBytes int64  `yaml:"maxEntryBytes"`
}
type CheckpointConfig struct {
	Kind      string `yaml:"kind,omitempty"`
	Directory string `yaml:"directory,omitempty"`
}
type ConcurrencyConfig struct {
	MaxInFlight int `yaml:"maxInFlight,omitempty"`
}

type ProviderSpec struct {
	Kind               string            `yaml:"kind"`
	ModelManifest      string            `yaml:"modelManifest"`
	EndpointRef        string            `yaml:"endpointRef,omitempty"`
	CredentialRef      string            `yaml:"credentialRef,omitempty"`
	AllowHTTP          bool              `yaml:"allowHttp,omitempty"`
	AllowLocalNetworks bool              `yaml:"allowLocalNetworks,omitempty"`
	Concurrency        ConcurrencyConfig `yaml:"concurrency,omitempty"`
	// Profile, when set, resolves InferenceSettings from a Geppetto engine
	// profile registry (e.g. ~/.config/pinocchio/profiles.yaml) instead of
	// manually constructing endpoint/credential settings. The value is a
	// profile slug like "umans-glm-5.2". ProfileRegistries is an optional
	// comma-separated list of registry source specs (YAML paths, sqlite:...
	// etc.); when empty, the default Pinocchio profiles.yaml is used.
	Profile           string `yaml:"profile,omitempty"`
	ProfileRegistries string `yaml:"profileRegistries,omitempty"`
}

func loadConfig(path string) (HostConfig, string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_PATH: %w", err)
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_READ: %w", err)
	}
	var cfg HostConfig
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_DECODE: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err == nil {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_MULTIPLE_DOCUMENTS")
	} else if err != io.EOF {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_DECODE: %w", err)
	}
	if cfg.SchemaVersion != HostConfigSchemaVersion {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_SCHEMA")
	}
	if strings.TrimSpace(cfg.ProfileID) == "" {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_PROFILE")
	}
	if len(cfg.Providers) == 0 {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_PROVIDERS")
	}
	base := filepath.Dir(absolute)
	resolve := func(value string) (string, error) {
		if value == "" {
			return "", nil
		}
		if filepath.IsAbs(value) {
			return "", fmt.Errorf("RAG_PROVIDER_CONFIG_PATH_ABSOLUTE")
		}
		resolved := filepath.Clean(filepath.Join(base, value))
		relative, err := filepath.Rel(base, resolved)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("RAG_PROVIDER_CONFIG_PATH_ESCAPE")
		}
		return resolved, nil
	}
	for _, target := range []*string{&cfg.Manifests.ModelsDir, &cfg.Manifests.PromptsDir, &cfg.Schemas.Directory, &cfg.Cache.Directory, &cfg.Checkpoints.Directory} {
		resolved, err := resolve(*target)
		if err != nil {
			return HostConfig{}, "", err
		}
		*target = resolved
	}
	if cfg.Manifests.ModelsDir == "" || cfg.Manifests.PromptsDir == "" {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_MANIFESTS_REQUIRED")
	}
	if cfg.Schemas.Directory == "" {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_SCHEMAS_REQUIRED")
	}
	if cfg.Cache.Kind != "filesystem-content-addressed/v1" || cfg.Cache.Directory == "" {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_CACHE")
	}
	if cfg.Checkpoints.Kind != "" && (cfg.Checkpoints.Kind != "filesystem-query-checkpoints/v1" || cfg.Checkpoints.Directory == "") {
		return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_CHECKPOINTS")
	}
	for name, spec := range cfg.Providers {
		if strings.TrimSpace(name) == "" || strings.TrimSpace(spec.Kind) == "" || strings.TrimSpace(spec.ModelManifest) == "" || spec.Concurrency.MaxInFlight < 0 {
			return HostConfig{}, "", fmt.Errorf("RAG_PROVIDER_CONFIG_PROVIDER_INVALID")
		}
		cfg.Providers[name] = spec
	}
	return cfg, absolute, nil
}

func resolveEnvRef(ref string, required bool) (string, error) {
	if ref == "" {
		if required {
			return "", fmt.Errorf("RAG_PROVIDER_ENV_REF_REQUIRED")
		}
		return "", nil
	}
	if !strings.HasPrefix(ref, "env:") || len(ref) <= len("env:") {
		return "", fmt.Errorf("RAG_PROVIDER_ENV_REF_INVALID")
	}
	name := strings.TrimSpace(strings.TrimPrefix(ref, "env:"))
	if name == "" || strings.ContainsAny(name, "= \t\r\n") {
		return "", fmt.Errorf("RAG_PROVIDER_ENV_REF_INVALID")
	}
	value, ok := os.LookupEnv(name)
	if required && (!ok || strings.TrimSpace(value) == "") {
		return "", fmt.Errorf("RAG_PROVIDER_ENV_MISSING")
	}
	return value, nil
}
