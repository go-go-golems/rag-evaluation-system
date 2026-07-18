package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func Evaluate(query Query, evidence []Evidence, answer *Answer, measures []ragcontract.Measure, timing map[string]int64, usage Usage, failures []ragcontract.FailureTrace, storageBytes int64) []Metric {
	relevant := map[string]float64{}
	for _, id := range query.RelevantIDs {
		relevant[id] = 1
	}
	for id, grade := range query.Grades {
		relevant[id] = grade
	}
	ids := make([]string, len(evidence))
	for i, item := range evidence {
		ids[i] = item.Collapse.ID
		if _, ok := relevant[ids[i]]; !ok {
			if _, exists := relevant[item.Chunk.Record.ID]; exists {
				ids[i] = item.Chunk.Record.ID
			}
		}
	}
	out := []Metric{}
	for _, measure := range measures {
		var config struct {
			Cutoffs []int    `json:"cutoffs"`
			Stages  []string `json:"stages"`
		}
		_ = json.Unmarshal(measure.Config, &config)
		var value any
		switch measure.Name {
		case "rag.precision":
			value = cutoffValues(config.Cutoffs, func(k int) float64 { return precision(ids, relevant, k) })
		case "rag.recall":
			value = cutoffValues(config.Cutoffs, func(k int) float64 { return recall(ids, relevant, k) })
		case "rag.hit-rate":
			value = cutoffValues(config.Cutoffs, func(k int) float64 {
				if recall(ids, relevant, k) > 0 {
					return 1
				}
				return 0
			})
		case "rag.mrr":
			value = mrr(ids, relevant)
		case "rag.ndcg":
			value = cutoffValues(config.Cutoffs, func(k int) float64 { return ndcg(ids, relevant, k) })
		case "rag.latency":
			selected := map[string]int64{}
			for _, stage := range config.Stages {
				selected[stage] = timing[stage]
			}
			value = selected
		case "rag.token-usage":
			value = map[string]int64{"input": usage.InputTokens, "output": usage.OutputTokens, "embedding": usage.EmbeddingTokens}
		case "rag.provider-cost":
			value = usage.Cost
		case "rag.storage-bytes":
			value = map[string]int64{"bytes": storageBytes}
		case "rag.failure-rates":
			value = map[string]any{"failures": len(failures), "failed": len(failures) > 0}
		case "rag.abstention":
			value = answer == nil || answer.Text == ""
		default:
			continue
		}
		raw, _ := json.Marshal(value)
		metric := Metric{Name: measure.Name, Unit: measure.Unit, Value: raw, Metadata: measure.Config}
		if number, ok := value.(float64); ok {
			metric.Numeric = &number
		}
		out = append(out, metric)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func cutoffValues(cutoffs []int, fn func(int) float64) map[string]float64 {
	result := map[string]float64{}
	for _, k := range cutoffs {
		key := itoa(k)
		result[key] = fn(k)
	}
	return result
}
func itoa(v int) string                                         { raw, _ := json.Marshal(v); return string(raw) }
func relevantAt(id string, relevant map[string]float64) float64 { return relevant[id] }
func precision(ids []string, relevant map[string]float64, k int) float64 {
	if k <= 0 {
		return 0
	}
	if k > len(ids) {
		k = len(ids)
	}
	if k == 0 {
		return 0
	}
	hits := 0
	for _, id := range ids[:k] {
		if relevantAt(id, relevant) > 0 {
			hits++
		}
	}
	return float64(hits) / float64(k)
}
func recall(ids []string, relevant map[string]float64, k int) float64 {
	if len(relevant) == 0 {
		return 0
	}
	if k > len(ids) {
		k = len(ids)
	}
	seen := map[string]bool{}
	for _, id := range ids[:max(k, 0)] {
		if relevantAt(id, relevant) > 0 {
			seen[id] = true
		}
	}
	return float64(len(seen)) / float64(len(relevant))
}
func mrr(ids []string, relevant map[string]float64) float64 {
	for i, id := range ids {
		if relevantAt(id, relevant) > 0 {
			return 1 / float64(i+1)
		}
	}
	return 0
}
func ndcg(ids []string, relevant map[string]float64, k int) float64 {
	if k > len(ids) {
		k = len(ids)
	}
	dcg := 0.0
	for i, id := range ids[:max(k, 0)] {
		grade := relevantAt(id, relevant)
		dcg += (math.Pow(2, grade) - 1) / math.Log2(float64(i+2))
	}
	grades := make([]float64, 0, len(relevant))
	for _, grade := range relevant {
		grades = append(grades, grade)
	}
	sort.Slice(grades, func(i, j int) bool { return grades[i] > grades[j] })
	if k > len(grades) {
		k = len(grades)
	}
	ideal := 0.0
	for i, grade := range grades[:max(k, 0)] {
		ideal += (math.Pow(2, grade) - 1) / math.Log2(float64(i+2))
	}
	if ideal == 0 {
		return 0
	}
	return dcg / ideal
}

type evaluateOperator struct{}

func (evaluateOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "evaluate.relevance", Version: "v1"}
}
func (evaluateOperator) Execute(_ context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	evidence, ok := inputs["evidence"].([]Evidence)
	if !ok {
		return nil, fmt.Errorf("RAG_EVALUATE_INPUT")
	}
	var config struct {
		Measures []ragcontract.Measure `json:"measures"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	metrics := Evaluate(env.CurrentQuery, evidence, nil, config.Measures, map[string]int64{}, env.Usage, nil, 0)
	return map[string]any{"evaluation": metrics}, nil
}
