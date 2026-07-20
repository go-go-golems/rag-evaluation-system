package ragoperators

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type Operator interface {
	Ref() ragcontract.OperatorRef
	Execute(context.Context, ragcontract.Node, map[string]any, *Environment) (map[string]any, error)
}
type Registry struct{ operators map[string]Operator }

func NewRegistry() *Registry { return &Registry{operators: map[string]Operator{}} }
func (r *Registry) Register(value Operator) error {
	if value == nil {
		return fmt.Errorf("RAG_OPERATOR_NIL")
	}
	id := value.Ref().ID()
	if _, ok := r.operators[id]; ok {
		return fmt.Errorf("RAG_OPERATOR_DUPLICATE: %s", id)
	}
	r.operators[id] = value
	return nil
}
func (r *Registry) Lookup(ref ragcontract.OperatorRef) (Operator, bool) {
	v, ok := r.operators[ref.ID()]
	return v, ok
}
func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.operators))
	for id := range r.operators {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
func NativeRegistry() *Registry {
	r := NewRegistry()
	for _, op := range []Operator{unitOperator{"units.identity"}, unitOperator{"units.individual-turns"}, unitOperator{"transcript.units.agents-view-runs"}, chunkOperator{}, chunkOperator{"chunks.identity"}, representationOperator{"representations.raw"}, representationOperator{"representations.structured-summary"}, representationOperator{"representations.synthetic-questions"}, combinedPreparationOperator{}, mergeOperator{}, embeddingOperator{}, indexOperator{}, indexOperator{"index.memory-smoke"}, retrieveOperator{"retrieve.bm25"}, retrieveOperator{"retrieve.vector"}, collapseOperator{"collapse.parent"}, fusionOperator{}, collapseOperator{"collapse.final"}, hydrateOperator{}, rerankOperator{}, answerOperator{}, evaluateOperator{}} {
		if err := r.Register(op); err != nil {
			panic(err)
		}
	}
	return r
}
