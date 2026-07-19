package ragproviders

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type FileManifestRegistry struct {
	models     map[string]ragcontract.ModelManifest
	prompts    map[string]ragcontract.PromptManifest
	promptText map[string]string
}

var _ ragoperators.ManifestResolver = (*FileManifestRegistry)(nil)

func LoadManifestRegistry(modelsDir, promptsDir string) (*FileManifestRegistry, error) {
	if modelsDir == "" || promptsDir == "" {
		return nil, fmt.Errorf("RAG_MANIFEST_DIRECTORY_REQUIRED")
	}
	models, err := loadModels(modelsDir)
	if err != nil {
		return nil, err
	}
	prompts, text, err := loadPrompts(promptsDir)
	if err != nil {
		return nil, err
	}
	return &FileManifestRegistry{models: models, prompts: prompts, promptText: text}, nil
}
func loadModels(dir string) (map[string]ragcontract.ModelManifest, error) {
	result := map[string]ragcontract.ModelManifest{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("RAG_MODEL_MANIFEST_DIRECTORY: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("RAG_MODEL_MANIFEST_READ: %w", err)
		}
		var value ragcontract.ModelManifest
		if err := ragcontract.DecodeStrict(data, &value); err != nil {
			return nil, fmt.Errorf("RAG_MODEL_MANIFEST_DECODE: %w", err)
		}
		if err := validateModelManifest(value); err != nil {
			return nil, fmt.Errorf("RAG_MODEL_MANIFEST_INVALID: %w", err)
		}
		key := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if _, exists := result[key]; exists || value.ModelID == "" {
			return nil, fmt.Errorf("RAG_MODEL_MANIFEST_DUPLICATE")
		}
		if _, exists := result[value.ModelID]; exists {
			return nil, fmt.Errorf("RAG_MODEL_MANIFEST_DUPLICATE")
		}
		result[key] = value
		result[value.ModelID] = value
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("RAG_MODEL_MANIFEST_EMPTY")
	}
	return result, nil
}
func loadPrompts(dir string) (map[string]ragcontract.PromptManifest, map[string]string, error) {
	result := map[string]ragcontract.PromptManifest{}
	text := map[string]string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_DIRECTORY: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		switch ext {
		case ".json":
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_READ: %w", err)
			}
			var value ragcontract.PromptManifest
			if err := ragcontract.DecodeStrict(data, &value); err != nil {
				return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_DECODE: %w", err)
			}
			if err := validatePromptManifest(value); err != nil {
				return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_INVALID: %w", err)
			}
			key := strings.TrimSuffix(entry.Name(), ext)
			if _, exists := result[key]; exists || value.PromptID == "" {
				return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_DUPLICATE")
			}
			result[key] = value
			result[value.PromptID] = value
		case ".txt", ".md":
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, nil, fmt.Errorf("RAG_PROMPT_TEXT_READ: %w", err)
			}
			key := strings.TrimSuffix(entry.Name(), ext)
			if _, exists := text[key]; exists {
				return nil, nil, fmt.Errorf("RAG_PROMPT_TEXT_DUPLICATE")
			}
			text[key] = string(data)
		}
	}
	if len(result) == 0 {
		return nil, nil, fmt.Errorf("RAG_PROMPT_MANIFEST_EMPTY")
	}
	seen := map[string]bool{}
	for _, manifest := range result {
		if seen[manifest.PromptID] {
			continue
		}
		seen[manifest.PromptID] = true
		if _, ok := text[manifest.PromptID]; !ok {
			return nil, nil, fmt.Errorf("RAG_PROMPT_TEXT_MISSING")
		}
	}
	return result, text, nil
}
func validateModelManifest(v ragcontract.ModelManifest) error {
	if err := ragcontract.ValidateManifestBase(v.ManifestBase, ragcontract.ModelManifestSchema, false); err != nil {
		return err
	}
	if v.ModelID == "" || v.ModelDigest == "" || v.ProviderAdapterVersion == "" || v.ImplementationVersion == "" || v.Tokenization == "" || v.Truncation == "" {
		return fmt.Errorf("required model identity missing")
	}
	return nil
}
func validatePromptManifest(v ragcontract.PromptManifest) error {
	if err := ragcontract.ValidateManifestBase(v.ManifestBase, ragcontract.PromptManifestSchema, false); err != nil {
		return err
	}
	if v.PromptID == "" || v.TemplateDigest == "" || v.OutputSchema == "" {
		return fmt.Errorf("required prompt identity missing")
	}
	return nil
}
func (r *FileManifestRegistry) Model(reference string) (ragcontract.ModelManifest, error) {
	if r == nil {
		return ragcontract.ModelManifest{}, fmt.Errorf("RAG_MODEL_MANIFEST_RESOLVER_UNAVAILABLE")
	}
	if v, ok := r.models[reference]; ok {
		return v, nil
	}
	for _, v := range r.models {
		if v.Digest == reference {
			return v, nil
		}
	}
	return ragcontract.ModelManifest{}, fmt.Errorf("RAG_MODEL_MANIFEST_MISSING")
}
func (r *FileManifestRegistry) Prompt(reference string) (ragcontract.PromptManifest, error) {
	if r == nil {
		return ragcontract.PromptManifest{}, fmt.Errorf("RAG_PROMPT_MANIFEST_RESOLVER_UNAVAILABLE")
	}
	if v, ok := r.prompts[reference]; ok {
		return v, nil
	}
	for _, v := range r.prompts {
		if v.Digest == reference {
			return v, nil
		}
	}
	return ragcontract.PromptManifest{}, fmt.Errorf("RAG_PROMPT_MANIFEST_MISSING")
}
func (r *FileManifestRegistry) PromptText(reference string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("RAG_PROMPT_TEXT_RESOLVER_UNAVAILABLE")
	}
	if value, ok := r.promptText[reference]; ok {
		return value, nil
	}
	manifest, err := r.Prompt(reference)
	if err != nil {
		return "", err
	}
	value, ok := r.promptText[manifest.PromptID]
	if !ok {
		return "", fmt.Errorf("RAG_PROMPT_TEXT_MISSING")
	}
	return value, nil
}
