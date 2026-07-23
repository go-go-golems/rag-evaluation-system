package ragworkflow

import (
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

type LoweringPhase string

const (
	PhasePreparation LoweringPhase = "preparation"
	PhaseQuery       LoweringPhase = "query"
)

type OperatorLowering struct {
	Operator         ragcontract.OperatorRef `json:"operator"`
	Phase            LoweringPhase           `json:"phase"`
	Task             workflowv3.TaskKey      `json:"task"`
	ProviderRequired bool                    `json:"providerRequired,omitempty"`
}

type OperatorRegistry struct{ entries map[string]OperatorLowering }

func NewOperatorRegistry() *OperatorRegistry {
	values := []OperatorLowering{
		entry("units.identity", PhasePreparation, TaskUnits), entry("units.individual-turns", PhasePreparation, TaskUnits), entry("transcript.units.agents-view-runs", PhasePreparation, TaskUnits),
		entry("chunks.identity", PhasePreparation, TaskChunks), entry("chunks.recursive", PhasePreparation, TaskChunks),
		entry("representations.raw", PhasePreparation, TaskRepresent), entry("representations.merge", PhasePreparation, TaskRepresent),
		entry("embed.model", PhasePreparation, TaskEmbed), entry("index.bleve-multi", PhasePreparation, TaskIndex), entry("index.memory-smoke", PhasePreparation, TaskIndex),
		entry("retrieve.bm25", PhaseQuery, TaskQuery), entry("retrieve.vector", PhaseQuery, TaskQuery), entry("collapse.parent", PhaseQuery, TaskQuery), entry("fusion.weighted-rrf", PhaseQuery, TaskQuery), entry("collapse.final", PhaseQuery, TaskQuery), entry("hydrate.source-evidence", PhaseQuery, TaskQuery), entry("evaluate.relevance", PhaseQuery, TaskQuery),
		providerEntry("representations.structured-summary", PhasePreparation, TaskRepresent), providerEntry("representations.synthetic-questions", PhasePreparation, TaskRepresent), providerEntry("representations.combined-summary-questions", PhasePreparation, TaskRepresent), providerEntry("rerank.cross-encoder", PhaseQuery, TaskQuery), providerEntry("generate.answer", PhaseQuery, TaskQuery),
	}
	registry := &OperatorRegistry{entries: make(map[string]OperatorLowering, len(values))}
	for _, value := range values {
		registry.entries[value.Operator.ID()] = value
	}
	return registry
}

func (r *OperatorRegistry) Definition(ref ragcontract.OperatorRef) (OperatorLowering, bool) {
	if r == nil {
		return OperatorLowering{}, false
	}
	value, found := r.entries[ref.ID()]
	return value, found
}

func (r *OperatorRegistry) Definitions() []OperatorLowering {
	if r == nil {
		return nil
	}
	values := make([]OperatorLowering, 0, len(r.entries))
	for _, value := range r.entries {
		values = append(values, value)
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Operator.ID() < values[j].Operator.ID() })
	return values
}

func entry(kind string, phase LoweringPhase, task workflowv3.TaskKey) OperatorLowering {
	return OperatorLowering{Operator: ragcontract.OperatorRef{Kind: kind, Version: "v1"}, Phase: phase, Task: task}
}
func providerEntry(kind string, phase LoweringPhase, task workflowv3.TaskKey) OperatorLowering {
	value := entry(kind, phase, task)
	value.ProviderRequired = true
	return value
}
