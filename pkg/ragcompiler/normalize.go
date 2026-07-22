package ragcompiler

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type CompileError struct{ Report ragcontract.Report }

func (e *CompileError) Error() string {
	return (&ragcontract.ValidationError{Report: e.Report}).Error()
}

func Normalize(input ragcontract.PipelineIR, registry *Registry) (ragcontract.PipelineIR, error) {
	if registry == nil {
		registry = BuiltinRegistry()
	}
	if input.SchemaVersion == "" {
		input.SchemaVersion = ragcontract.PipelineSchemaVersion
	}
	expanded, err := expandRecipes(input, registry)
	if err != nil {
		return ragcontract.PipelineIR{}, err
	}
	report := ragcontract.ValidatePipeline(expanded)
	aliases := map[string]ragcontract.PortKind{}
	for _, in := range expanded.Inputs {
		aliases[in.ID] = in.Kind
	}
	nodes := map[string]ragcontract.Node{}
	defs := map[string]OperatorDefinition{}
	for i, n := range expanded.Nodes {
		d, ok := registry.Definition(n.Operator)
		if !ok {
			report.Add("RAG_V2_OPERATOR_UNKNOWN", fmt.Sprintf("$.nodes[%d].operator", i), "operator "+n.Operator.ID()+" is not registered")
			continue
		}
		config, e := normalizeConfig(n.Config, d)
		if e != nil {
			report.Add("RAG_V2_CONFIG", fmt.Sprintf("$.nodes[%d].config", i), e.Error())
		} else {
			n.Config = config
		}
		nodes[n.ID] = n
		defs[n.ID] = d
	}
	validateSemantics(nodes, &report)
	if !report.OK() {
		report.Normalize()
		return ragcontract.PipelineIR{}, &CompileError{Report: report}
	}
	order, cycle := topological(expanded.Nodes)
	if len(cycle) > 0 {
		report.Add("RAG_V2_GRAPH_CYCLE", "$.nodes", "cycle: "+strings.Join(cycle, " -> "))
		report.Normalize()
		return ragcontract.PipelineIR{}, &CompileError{Report: report}
	}
	resolvedIDs := map[string]string{}
	outputKinds := map[string]map[string]ragcontract.PortKind{}
	for _, in := range expanded.Inputs {
		resolvedIDs[in.ID] = in.ID
		outputKinds[in.ID] = map[string]ragcontract.PortKind{"out": in.Kind}
	}
	result := expanded
	result.Nodes = nil
	semanticNodeIDs := map[string]string{}
	for _, oldID := range order {
		n := nodes[oldID]
		d := defs[oldID]
		sort.Slice(n.Inputs, func(i, j int) bool { return n.Inputs[i].Port < n.Inputs[j].Port })
		for i := range n.Inputs {
			binding := &n.Inputs[i]
			sourceAlias := binding.From.NodeID
			sourceID, ok := resolvedIDs[sourceAlias]
			if !ok {
				report.Add("RAG_V2_UNKNOWN_REFERENCE", "$.nodes["+oldID+"].inputs", "unknown source "+sourceAlias)
				continue
			}
			binding.From.NodeID = sourceID
			actual, ok := outputKinds[sourceAlias][binding.From.Port]
			if !ok {
				report.Add("RAG_V2_OUTPUT_PORT", "$.nodes["+oldID+"].inputs", "unknown source port "+binding.From.Port)
				continue
			}
			expected, ok := inputKind(d, binding.Port)
			if !ok {
				report.Add("RAG_V2_INPUT_PORT", "$.nodes["+oldID+"].inputs", "unknown input port "+binding.Port)
				continue
			}
			if actual != expected {
				report.Add("RAG_V2_PORT_MISMATCH", "$.nodes["+oldID+"].inputs", fmt.Sprintf("%s is %s, expected %s", binding.Port, actual, expected))
			}
		}
		identity := struct {
			Operator ragcontract.OperatorRef    `json:"operator"`
			Inputs   []ragcontract.InputBinding `json:"inputs"`
			Config   any                        `json:"config"`
			Order    int                        `json:"order,omitempty"`
		}{n.Operator, n.Inputs, n.Config, n.Order}
		digest, _ := ragcontract.Digest(identity)
		n.ID = "n-" + strings.TrimPrefix(digest, "sha256:")[:20]
		if previous, exists := semanticNodeIDs[n.ID]; exists {
			report.Add("RAG_V2_NODE_ID_COLLISION", "$.nodes["+oldID+"]", "semantically duplicate node "+previous)
		}
		semanticNodeIDs[n.ID] = oldID
		resolvedIDs[oldID] = n.ID
		outputKinds[oldID] = map[string]ragcontract.PortKind{}
		for _, p := range d.Outputs {
			outputKinds[oldID][p.Name] = p.Kind
		}
		result.Nodes = append(result.Nodes, n)
	}
	canonicalOrder, _ := topological(result.Nodes)
	canonicalNodes := make(map[string]ragcontract.Node, len(result.Nodes))
	for _, node := range result.Nodes {
		canonicalNodes[node.ID] = node
	}
	result.Nodes = result.Nodes[:0]
	for _, id := range canonicalOrder {
		result.Nodes = append(result.Nodes, canonicalNodes[id])
	}
	for i := range result.Outputs {
		old := result.Outputs[i].From.NodeID
		if id, ok := resolvedIDs[old]; ok {
			result.Outputs[i].From.NodeID = id
		}
		actual := outputKinds[old][result.Outputs[i].From.Port]
		if actual != "" && actual != result.Outputs[i].Kind {
			report.Add("RAG_V2_OUTPUT_MISMATCH", fmt.Sprintf("$.outputs[%d]", i), fmt.Sprintf("declared %s, actual %s", result.Outputs[i].Kind, actual))
		}
	}
	sort.Slice(result.Inputs, func(i, j int) bool { return result.Inputs[i].ID < result.Inputs[j].ID })
	sort.Slice(result.Outputs, func(i, j int) bool { return result.Outputs[i].Name < result.Outputs[j].Name })
	if !report.OK() {
		report.Normalize()
		return ragcontract.PipelineIR{}, &CompileError{Report: report}
	}
	return result, nil
}

func inputKind(d OperatorDefinition, name string) (ragcontract.PortKind, bool) {
	for _, p := range d.Inputs {
		if p.Name == name {
			return p.Kind, true
		}
	}
	if d.Ref.Kind == "fusion.weighted-rrf" && strings.HasPrefix(name, "channel.") {
		return ragcontract.PortRankedParents, true
	}
	if d.Ref.Kind == "index.bleve-multi" && strings.HasPrefix(name, "representations.") {
		return ragcontract.PortRepresentations, true
	}
	if d.Ref.Kind == "representations.merge" && strings.HasPrefix(name, "set.") {
		return ragcontract.PortRepresentations, true
	}
	return "", false
}
func topological(nodes []ragcontract.Node) ([]string, []string) {
	by := map[string]ragcontract.Node{}
	indegree := map[string]int{}
	next := map[string][]string{}
	for _, n := range nodes {
		by[n.ID] = n
		indegree[n.ID] = 0
	}
	for _, n := range nodes {
		seen := map[string]bool{}
		for _, b := range n.Inputs {
			if _, ok := by[b.From.NodeID]; ok && !seen[b.From.NodeID] {
				indegree[n.ID]++
				next[b.From.NodeID] = append(next[b.From.NodeID], n.ID)
				seen[b.From.NodeID] = true
			}
		}
	}
	ready := []string{}
	for id, d := range indegree {
		if d == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)
	order := []string{}
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		order = append(order, id)
		sort.Strings(next[id])
		for _, to := range next[id] {
			indegree[to]--
			if indegree[to] == 0 {
				ready = append(ready, to)
				sort.Strings(ready)
			}
		}
	}
	if len(order) == len(nodes) {
		return order, nil
	}
	cycle := []string{}
	for id, d := range indegree {
		if d > 0 {
			cycle = append(cycle, id)
		}
	}
	sort.Strings(cycle)
	return nil, cycle
}

type recipeOperator struct {
	Kind    string          `json:"kind"`
	Version string          `json:"version"`
	Config  json.RawMessage `json:"config"`
}
type composeConfig struct {
	Operators []recipeOperator `json:"operators"`
}

func expandRecipes(input ragcontract.PipelineIR, registry *Registry) (ragcontract.PipelineIR, error) {
	result := input
	result.Nodes = nil
	replacements := map[string]ragcontract.PortRef{}
	for _, n := range input.Nodes {
		d, ok := registry.Definition(n.Operator)
		if !ok || !d.Recipe {
			result.Nodes = append(result.Nodes, n)
			continue
		}
		var config composeConfig
		if err := strictRaw(n.Config, &config); err != nil {
			return result, fmt.Errorf("RAG_V2_RECIPE_CONFIG: %w", err)
		}
		if len(config.Operators) == 0 {
			return result, fmt.Errorf("RAG_V2_RECIPE_EMPTY: %s", n.ID)
		}
		if len(n.Inputs) != 1 {
			return result, fmt.Errorf("RAG_V2_RECIPE_INPUT: %s requires one input", n.ID)
		}
		source := n.Inputs[0].From
		mergeInputs := make([]ragcontract.InputBinding, 0, len(config.Operators))
		namedOutputs := map[string]ragcontract.PortRef{}
		for i, entry := range config.Operators {
			ref := ragcontract.OperatorRef{Kind: entry.Kind, Version: entry.Version}
			definition, ok := registry.Definition(ref)
			if !ok || definition.Recipe || len(definition.Inputs) == 0 || definition.Inputs[0].Kind != ragcontract.PortChunks || len(definition.Outputs) != 1 || definition.Outputs[0].Kind != ragcontract.PortRepresentations {
				return result, fmt.Errorf("RAG_V2_RECIPE_OPERATOR: unsupported representation operator %s", ref.ID())
			}
			raw := entry.Config
			if len(raw) == 0 {
				raw = json.RawMessage(`{}`)
			}
			id := fmt.Sprintf("%s.%03d", n.ID, i+1)
			bindings := []ragcontract.InputBinding{{Port: definition.Inputs[0].Name, From: source}}
			var descriptorConfig struct {
				Name string `json:"name"`
				From string `json:"from"`
			}
			_ = json.Unmarshal(raw, &descriptorConfig)
			if descriptorConfig.From != "" {
				upstream, exists := namedOutputs[descriptorConfig.From]
				if !exists {
					return result, fmt.Errorf("RAG_V2_RECIPE_DEPENDENCY: representation %q references unknown or later %q", descriptorConfig.Name, descriptorConfig.From)
				}
				bindings = append(bindings, ragcontract.InputBinding{Port: "source", From: upstream})
			}
			node := ragcontract.Node{ID: id, Operator: ref, Inputs: bindings, Config: raw, Order: n.Order + i}
			result.Nodes = append(result.Nodes, node)
			outputRef := ragcontract.PortRef{NodeID: id, Port: definition.Outputs[0].Name}
			if descriptorConfig.Name != "" {
				namedOutputs[descriptorConfig.Name] = outputRef
			}
			mergeInputs = append(mergeInputs, ragcontract.InputBinding{Port: fmt.Sprintf("set.%03d", i+1), From: outputRef})
		}
		output := mergeInputs[0].From
		if len(mergeInputs) > 1 {
			mergeID := n.ID + ".merge"
			result.Nodes = append(result.Nodes, ragcontract.Node{ID: mergeID, Operator: ragcontract.OperatorRef{Kind: "representations.merge", Version: "v1"}, Inputs: mergeInputs, Config: json.RawMessage(`{}`), Order: n.Order + len(config.Operators)})
			output = ragcontract.PortRef{NodeID: mergeID, Port: "representations"}
		}
		replacements[n.ID] = output
	}
	for i := range result.Nodes {
		for j := range result.Nodes[i].Inputs {
			if replacement, ok := replacements[result.Nodes[i].Inputs[j].From.NodeID]; ok {
				result.Nodes[i].Inputs[j].From = replacement
			}
		}
	}
	for i := range result.Outputs {
		if replacement, ok := replacements[result.Outputs[i].From.NodeID]; ok {
			result.Outputs[i].From = replacement
		}
	}
	return result, nil
}
