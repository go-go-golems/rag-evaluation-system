package ragproduct

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

const QualificationSchemaVersion = "rag-product-qualification/v1"

type PromptBinding struct {
	Reference     string `json:"reference"`
	SchemaVersion string `json:"schemaVersion"`
	Digest        string `json:"digest"`
}

type Qualification struct {
	SchemaVersion string                     `json:"schemaVersion"`
	ProductID     string                     `json:"productId"`
	Models        []ragcontract.ModelBinding `json:"models"`
	Prompts       []PromptBinding            `json:"prompts"`
	Citations     ragcontract.CitationPolicy `json:"citations"`
	Runtime       ragcontract.RuntimePolicy  `json:"runtime"`
	Study         ragcontract.Study          `json:"study"`
}

// Qualify freezes the exact compiled product bindings and wraps the identical
// normalized pipeline in a research study. It performs no lifecycle operation.
func Qualify(plan ragcontract.ProductPlan, dataset ragcontract.DatasetBinding, measures []ragcontract.Measure) (Qualification, error) {
	return qualify(plan, dataset, measures, nil)
}

// QualifyResolved resolves every model/prompt reference to an exact manifest.
func QualifyResolved(plan ragcontract.ProductPlan, dataset ragcontract.DatasetBinding, measures []ragcontract.Measure, resolver ragoperators.ManifestResolver) (Qualification, error) {
	if resolver == nil {
		return Qualification{}, fmt.Errorf("RAG_PRODUCT_QUALIFICATION_RESOLVER")
	}
	return qualify(plan, dataset, measures, resolver)
}

func qualify(plan ragcontract.ProductPlan, dataset ragcontract.DatasetBinding, measures []ragcontract.Measure, resolver ragoperators.ManifestResolver) (Qualification, error) {
	compiled, err := ragcompiler.CompileProduct(plan, nil)
	if err != nil {
		return Qualification{}, err
	}
	if dataset.ManifestDigest == "" || dataset.Status == "" || dataset.Split == "" || dataset.RelevanceTarget == "" {
		return Qualification{}, fmt.Errorf("RAG_PRODUCT_QUALIFICATION_DATASET")
	}
	productID, err := ragcompiler.ProductSemanticIdentity(compiled)
	if err != nil {
		return Qualification{}, err
	}
	models, prompts, err := qualificationBindings(compiled, resolver)
	if err != nil {
		return Qualification{}, err
	}
	metadata, _ := json.Marshal(map[string]any{"productId": productID, "models": models, "prompts": prompts, "citations": compiled.Citations, "runtime": compiled.Runtime})
	study := ragcontract.Study{SchemaVersion: ragcontract.StudySchemaVersion, Variants: []ragcontract.Variant{{ID: "product-qualification", Pipeline: compiled.Pipeline, Metadata: metadata}}, Bindings: compiled.Bindings, Dataset: dataset, Measures: measures, Replicates: 1, Display: ragcontract.DisplayMetadata{Name: compiled.Display.Name + " qualification", Tags: map[string]string{"qualification.product": productID}}}
	if _, err := ragcompiler.ExpandStudy(study, nil); err != nil {
		return Qualification{}, err
	}
	return Qualification{SchemaVersion: QualificationSchemaVersion, ProductID: productID, Models: models, Prompts: prompts, Citations: compiled.Citations, Runtime: compiled.Runtime, Study: study}, nil
}

func qualificationBindings(plan ragcontract.ProductPlan, resolver ragoperators.ManifestResolver) ([]ragcontract.ModelBinding, []PromptBinding, error) {
	models := append([]ragcontract.ModelBinding(nil), plan.Models...)
	modelSeen, promptSeen := map[string]bool{}, map[string]bool{}
	for _, value := range models {
		modelSeen[value.Reference] = true
	}
	prompts := []PromptBinding{}
	for _, node := range plan.Pipeline.Nodes {
		var config map[string]any
		if json.Unmarshal(node.Config, &config) != nil {
			continue
		}
		if reference, ok := config["model"].(string); ok && reference != "" && !modelSeen[reference] {
			if resolver == nil {
				return nil, nil, fmt.Errorf("RAG_PRODUCT_QUALIFICATION_MODEL_UNRESOLVED: %s", reference)
			}
			manifest, err := resolver.Model(reference)
			if err != nil {
				return nil, nil, err
			}
			models = append(models, ragcontract.ModelBinding{Reference: reference, Manifest: manifest.SchemaVersion, Digest: manifest.Digest})
			modelSeen[reference] = true
		}
		if reference, ok := config["prompt"].(string); ok && reference != "" && !promptSeen[reference] {
			if resolver == nil {
				return nil, nil, fmt.Errorf("RAG_PRODUCT_QUALIFICATION_PROMPT_UNRESOLVED: %s", reference)
			}
			manifest, err := resolver.Prompt(reference)
			if err != nil {
				return nil, nil, err
			}
			prompts = append(prompts, PromptBinding{Reference: reference, SchemaVersion: manifest.SchemaVersion, Digest: manifest.Digest})
			promptSeen[reference] = true
		}
	}
	sort.Slice(models, func(i, j int) bool { return models[i].Reference < models[j].Reference })
	sort.Slice(prompts, func(i, j int) bool { return prompts[i].Reference < prompts[j].Reference })
	return models, prompts, nil
}
