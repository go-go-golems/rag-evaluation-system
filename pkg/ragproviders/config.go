package ragproviders

import (
	"fmt"
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
type ProviderSpec struct {
	Kind               string `yaml:"kind"`
	ModelManifest      string `yaml:"modelManifest"`
	EndpointRef        string `yaml:"endpointRef,omitempty"`
	CredentialRef      string `yaml:"credentialRef,omitempty"`
	AllowHTTP          bool   `yaml:"allowHttp,omitempty"`
	AllowLocalNetworks bool   `yaml:"allowLocalNetworks,omitempty"`
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
	resolve := func(value string) string {
		if value == "" || filepath.IsAbs(value) {
			return value
		}
		return filepath.Join(base, value)
	}
	cfg.Manifests.ModelsDir = resolve(cfg.Manifests.ModelsDir)
	cfg.Manifests.PromptsDir = resolve(cfg.Manifests.PromptsDir)
	cfg.Schemas.Directory = resolve(cfg.Schemas.Directory)
	cfg.Cache.Directory = resolve(cfg.Cache.Directory)
	for name, spec := range cfg.Providers {
		if strings.TrimSpace(name) == "" || strings.TrimSpace(spec.Kind) == "" || strings.TrimSpace(spec.ModelManifest) == "" {
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
