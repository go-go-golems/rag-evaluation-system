package ragoperators

import (
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type StaticManifestResolver struct {
	Models  map[string]ragcontract.ModelManifest
	Prompts map[string]ragcontract.PromptManifest
}

func (r StaticManifestResolver) Model(reference string) (ragcontract.ModelManifest, error) {
	value, ok := r.Models[reference]
	if !ok {
		for _, candidate := range r.Models {
			if candidate.Digest == reference {
				value = candidate
				ok = true
				break
			}
		}
	}
	if !ok {
		return ragcontract.ModelManifest{}, fmt.Errorf("RAG_MODEL_MANIFEST_MISSING: %s", reference)
	}
	if err := ragcontract.ValidateManifestBase(value.ManifestBase, ragcontract.ModelManifestSchema, false); err != nil {
		return value, err
	}
	return value, nil
}
func (r StaticManifestResolver) Prompt(reference string) (ragcontract.PromptManifest, error) {
	value, ok := r.Prompts[reference]
	if !ok {
		for _, candidate := range r.Prompts {
			if candidate.Digest == reference {
				value = candidate
				ok = true
				break
			}
		}
	}
	if !ok {
		return ragcontract.PromptManifest{}, fmt.Errorf("RAG_PROMPT_MANIFEST_MISSING: %s", reference)
	}
	if err := ragcontract.ValidateManifestBase(value.ManifestBase, ragcontract.PromptManifestSchema, false); err != nil {
		return value, err
	}
	return value, nil
}
func resolveModel(env *Environment, reference string) (ragcontract.ModelManifest, error) {
	if env == nil || env.Manifests == nil {
		return ragcontract.ModelManifest{}, fmt.Errorf("RAG_MODEL_MANIFEST_RESOLVER_UNAVAILABLE")
	}
	return env.Manifests.Model(reference)
}
func resolvePrompt(env *Environment, reference string) (ragcontract.PromptManifest, error) {
	if env == nil || env.Manifests == nil {
		return ragcontract.PromptManifest{}, fmt.Errorf("RAG_PROMPT_MANIFEST_RESOLVER_UNAVAILABLE")
	}
	return env.Manifests.Prompt(reference)
}
