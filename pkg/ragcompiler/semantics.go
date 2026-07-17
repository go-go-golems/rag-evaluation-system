package ragcompiler

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func validateSemantics(nodes map[string]ragcontract.Node, report *ragcontract.Report) {
	consumers := map[string][]ragcontract.Node{}
	for _, node := range nodes {
		for _, input := range node.Inputs {
			consumers[input.From.NodeID] = append(consumers[input.From.NodeID], node)
		}
	}
	for id, node := range nodes {
		config := map[string]any{}
		_ = json.Unmarshal(node.Config, &config)
		switch node.Operator.Kind {
		case "retrieve.bm25", "retrieve.vector":
			if value, ok := config["topK"].(float64); !ok || value <= 0 || value != float64(int(value)) {
				report.Add("RAG_V2_TOP_K", "$.nodes["+id+"].config.topK", "topK must be a positive integer")
			}
			if config["representation"] == "question" {
				safe := false
				for _, consumer := range consumers[id] {
					if consumer.Operator.Kind == "collapse.parent" {
						safe = true
					}
				}
				if !safe {
					report.Add("RAG_V2_COLLAPSE_REQUIRED", "$.nodes["+id+"]", "question representations require channel-local parent collapse")
				}
			}
		case "collapse.parent", "collapse.final":
			scope := fmt.Sprint(config["scope"])
			if scope != "chunk" && scope != "unit" {
				report.Add("RAG_V2_COLLAPSE_SCOPE", "$.nodes["+id+"].config.scope", "scope must be chunk or unit")
			}
		case "fusion.weighted-rrf":
			value, ok := config["rankConstant"].(float64)
			if !ok || value <= 0 || value != float64(int(value)) {
				report.Add("RAG_V2_RRF_RANK_CONSTANT", "$.nodes["+id+"].config.rankConstant", "rankConstant must be a positive integer")
			}
			for _, input := range node.Inputs {
				source, ok := nodes[input.From.NodeID]
				if ok && source.Operator.Kind != "collapse.parent" {
					report.Add("RAG_V2_CHANNEL_COLLAPSE", "$.nodes["+id+"].inputs", "fusion inputs must come from channel-local collapse.parent nodes")
				}
			}
		}
	}
}
