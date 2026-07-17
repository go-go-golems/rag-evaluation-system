package ragcompiler

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func CompileProduct(input ragcontract.ProductPlan, registry *Registry) (ragcontract.ProductPlan, error) {
	if input.SchemaVersion == "" {
		input.SchemaVersion = ragcontract.ProductSchemaVersion
	}
	if input.SchemaVersion != ragcontract.ProductSchemaVersion {
		return input, fmt.Errorf("RAG_V2_PRODUCT_SCHEMA: expected %s", ragcontract.ProductSchemaVersion)
	}
	pipeline, err := Normalize(input.Pipeline, registry)
	if err != nil {
		return input, err
	}
	input.Pipeline = pipeline
	sort.Slice(input.Bindings, func(i, j int) bool { return input.Bindings[i].SlotID < input.Bindings[j].SlotID })
	sort.Slice(input.Models, func(i, j int) bool { return input.Models[i].Reference < input.Models[j].Reference })
	if input.Runtime.MaxResults <= 0 {
		input.Runtime.MaxResults = 10
	}
	if input.Runtime.TracePolicy == "" {
		input.Runtime.TracePolicy = "authoritative"
	}
	if input.Runtime.FailurePolicy == "" {
		input.Runtime.FailurePolicy = "fail-closed"
	}
	if input.Citations.Mode == "" {
		input.Citations.Mode = "required"
	}
	if input.Citations.Mode == "required" {
		input.Citations.RequireSourceText = true
	}
	if err := validateBindings(pipeline, input.Bindings, false); err != nil {
		return input, err
	}
	return input, nil
}

func ExpandStudy(input ragcontract.Study, registry *Registry) ([]ragcontract.ExpandedCell, error) {
	if input.SchemaVersion == "" {
		input.SchemaVersion = ragcontract.StudySchemaVersion
	}
	if input.SchemaVersion != ragcontract.StudySchemaVersion {
		return nil, fmt.Errorf("RAG_V2_STUDY_SCHEMA: expected %s", ragcontract.StudySchemaVersion)
	}
	if input.Replicates <= 0 {
		input.Replicates = 1
	}
	if len(input.Variants) == 0 {
		return nil, fmt.Errorf("RAG_V2_STUDY_VARIANTS: at least one variant is required")
	}
	measures, err := DeriveRequestedMeasures(input.Measures)
	if err != nil {
		return nil, err
	}
	variants := append([]ragcontract.Variant(nil), input.Variants...)
	sort.Slice(variants, func(i, j int) bool { return variants[i].ID < variants[j].ID })
	for i := 1; i < len(variants); i++ {
		if variants[i].ID == variants[i-1].ID {
			return nil, fmt.Errorf("RAG_V2_VARIANT_DUPLICATE: %s", variants[i].ID)
		}
	}
	factors := append([]ragcontract.Factor(nil), input.Factors...)
	sort.Slice(factors, func(i, j int) bool { return factors[i].ID < factors[j].ID })
	for i := range factors {
		if i > 0 && factors[i].ID == factors[i-1].ID {
			return nil, fmt.Errorf("RAG_V2_FACTOR_DUPLICATE: %s", factors[i].ID)
		}
		sort.Slice(factors[i].Values, func(a, b int) bool { return factors[i].Values[a].ID < factors[i].Values[b].ID })
		if len(factors[i].Values) == 0 {
			return nil, fmt.Errorf("RAG_V2_FACTOR_EMPTY: %s", factors[i].ID)
		}
		for valueIndex := 1; valueIndex < len(factors[i].Values); valueIndex++ {
			if factors[i].Values[valueIndex].ID == factors[i].Values[valueIndex-1].ID {
				return nil, fmt.Errorf("RAG_V2_FACTOR_VALUE_DUPLICATE: %s/%s", factors[i].ID, factors[i].Values[valueIndex].ID)
			}
		}
	}
	combinations := factorCombinations(factors)
	cells := []ragcontract.ExpandedCell{}
	for _, variant := range variants {
		for _, selection := range combinations {
			pipeline := clonePipeline(variant.Pipeline)
			var err error
			pipeline, err = applyOverrides(pipeline, factors, selection)
			if err != nil {
				return nil, err
			}
			pipeline, err = Normalize(pipeline, registry)
			if err != nil {
				return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
			}
			if err = validateBindings(pipeline, input.Bindings, true); err != nil {
				return nil, fmt.Errorf("variant %s: %w", variant.ID, err)
			}
			execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: pipeline, Bindings: selectedBindings(pipeline, input.Bindings), Dataset: input.Dataset, Measures: measures, VariantID: variant.ID, Factors: selection}
			identity := execution
			identity.CellID = ""
			cellID, err := ragcontract.Digest(identity)
			if err != nil {
				return nil, err
			}
			execution.CellID = cellID
			cells = append(cells, ragcontract.ExpandedCell{ID: cellID, VariantID: variant.ID, Factors: selection, Replicates: input.Replicates, Execution: execution})
		}
	}
	return cells, nil
}

func DeriveRequestedMeasures(input []ragcontract.Measure) ([]ragcontract.Measure, error) {
	result := append([]ragcontract.Measure(nil), input...)
	for i := range result {
		if result[i].Name == "" || result[i].ValueKind == "" {
			return nil, fmt.Errorf("RAG_V2_MEASURE_REQUIRED: measure name and value kind are required")
		}
		if result[i].Version == "" {
			result[i].Version = "v1"
		}
		canonical, err := ragcontract.CanonicalRaw(result[i].Config, "{}")
		if err != nil {
			return nil, fmt.Errorf("RAG_V2_MEASURE_CONFIG: %s: %w", result[i].Name, err)
		}
		result[i].Config = canonical
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Name == result[j].Name {
			return result[i].Version < result[j].Version
		}
		return result[i].Name < result[j].Name
	})
	for i := 1; i < len(result); i++ {
		if result[i].Name == result[i-1].Name && result[i].Version == result[i-1].Version {
			return nil, fmt.Errorf("RAG_V2_MEASURE_DUPLICATE: %s/%s", result[i].Name, result[i].Version)
		}
	}
	return result, nil
}
func ProductSemanticIdentity(input ragcontract.ProductPlan) (string, error) {
	input.Display = ragcontract.DisplayMetadata{}
	normalized, err := CompileProduct(input, nil)
	if err != nil {
		return "", err
	}
	normalized.Display = ragcontract.DisplayMetadata{}
	return ragcontract.Digest(normalized)
}
func StudySemanticIdentity(input ragcontract.Study) (string, error) {
	input.Display = ragcontract.DisplayMetadata{}
	for i := range input.Variants {
		input.Variants[i].Metadata = nil
	}
	cells, err := ExpandStudy(input, nil)
	if err != nil {
		return "", err
	}
	identity := struct {
		SchemaVersion string                     `json:"schemaVersion"`
		Cells         []ragcontract.ExpandedCell `json:"cells"`
		Acceptance    json.RawMessage            `json:"acceptance,omitempty"`
	}{SchemaVersion: ragcontract.StudySchemaVersion, Cells: cells, Acceptance: input.Acceptance}
	return ragcontract.Digest(identity)
}

func validateBindings(p ragcontract.PipelineIR, bindings []ragcontract.ArtifactBinding, allowExtra bool) error {
	by := map[string]ragcontract.ArtifactBinding{}
	for _, b := range bindings {
		if _, ok := by[b.SlotID]; ok {
			return fmt.Errorf("RAG_V2_BINDING_DUPLICATE: %s", b.SlotID)
		}
		by[b.SlotID] = b
		if b.Digest == "" || b.SchemaVersion == "" {
			return fmt.Errorf("RAG_V2_BINDING_IDENTITY: %s requires digest and schema", b.SlotID)
		}
	}
	used := map[string]bool{}
	for _, slot := range p.Inputs {
		if slot.BindingMode != "artifact" {
			continue
		}
		b, ok := by[slot.ID]
		if !ok {
			return fmt.Errorf("RAG_V2_BINDING_MISSING: %s", slot.ID)
		}
		used[slot.ID] = true
		if slot.Digest != "" && slot.Digest != b.Digest {
			return fmt.Errorf("RAG_V2_BINDING_DIGEST: %s", slot.ID)
		}
		if slot.ManifestSchema != b.SchemaVersion {
			return fmt.Errorf("RAG_V2_BINDING_SCHEMA: %s", slot.ID)
		}
	}
	if !allowExtra {
		for slotID := range by {
			if !used[slotID] {
				return fmt.Errorf("RAG_V2_BINDING_UNUSED: %s", slotID)
			}
		}
	}
	return nil
}
func selectedBindings(p ragcontract.PipelineIR, bindings []ragcontract.ArtifactBinding) []ragcontract.ArtifactBinding {
	needed := map[string]bool{}
	for _, slot := range p.Inputs {
		if slot.BindingMode == "artifact" {
			needed[slot.ID] = true
		}
	}
	selected := make([]ragcontract.ArtifactBinding, 0, len(needed))
	for _, binding := range bindings {
		if needed[binding.SlotID] {
			selected = append(selected, binding)
		}
	}
	return sortedBindings(selected)
}
func sortedBindings(v []ragcontract.ArtifactBinding) []ragcontract.ArtifactBinding {
	r := append([]ragcontract.ArtifactBinding(nil), v...)
	sort.Slice(r, func(i, j int) bool { return r[i].SlotID < r[j].SlotID })
	return r
}
func clonePipeline(v ragcontract.PipelineIR) ragcontract.PipelineIR {
	b, _ := json.Marshal(v)
	var r ragcontract.PipelineIR
	_ = json.Unmarshal(b, &r)
	return r
}
func applyOverrides(p ragcontract.PipelineIR, factors []ragcontract.Factor, selections []ragcontract.FactorSelection) (ragcontract.PipelineIR, error) {
	selected := map[string]string{}
	for _, selection := range selections {
		selected[selection.FactorID] = selection.ValueID
	}
	for _, factor := range factors {
		wanted := selected[factor.ID]
		for _, value := range factor.Values {
			if value.ID != wanted {
				continue
			}
			for _, override := range value.Overrides {
				found := false
				for i := range p.Nodes {
					if p.Nodes[i].ID != override.NodeID {
						continue
					}
					base, err := object(p.Nodes[i].Config)
					if err != nil {
						return p, err
					}
					patch, err := object(override.Config)
					if err != nil {
						return p, err
					}
					for key, item := range patch {
						base[key] = item
					}
					p.Nodes[i].Config = mustJSON(base)
					found = true
					break
				}
				if !found {
					return p, fmt.Errorf("RAG_V2_OVERRIDE_NODE: %s", override.NodeID)
				}
			}
		}
	}
	return p, nil
}
func object(raw json.RawMessage) (map[string]any, error) {
	var v map[string]any
	if len(raw) == 0 {
		raw = []byte(`{}`)
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("expected object")
	}
	return v, nil
}
func factorCombinations(factors []ragcontract.Factor) [][]ragcontract.FactorSelection {
	result := [][]ragcontract.FactorSelection{{}}
	for _, f := range factors {
		next := [][]ragcontract.FactorSelection{}
		for _, prefix := range result {
			for _, v := range f.Values {
				item := append([]ragcontract.FactorSelection(nil), prefix...)
				canonical, _ := ragcontract.CanonicalRaw(v.Value, "null")
				item = append(item, ragcontract.FactorSelection{FactorID: f.ID, ValueID: v.ID, Value: canonical})
				next = append(next, item)
			}
		}
		result = next
	}
	return result
}
