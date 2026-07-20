package preparationworkflow

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

// CanonicalMapping is the explicit canonical-pipeline shape supported by the
// durable combined-preparation workflow: raw + combined-derived -> merge ->
// embedding. It prevents heuristic publication-key selection in rag-worker.
type CanonicalMapping struct {
	CombinedNode               ragcontract.Node
	EmbeddingNode              ragcontract.Node
	RawRepresentationName      string
	MaxRepresentationsPerChunk int
	RawOutputKey               string
	DerivedOutputKey           string
	MergedOutputKey            string
	EmbeddingOutputKey         string
}

func DeriveCanonicalMapping(pipeline ragcontract.PipelineIR) (CanonicalMapping, error) {
	combined, found := singleNode(pipeline, "representations.combined-summary-questions")
	if !found {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_COMBINED")
	}
	chunks, ok := inputFrom(combined, "chunks")
	if !ok {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_COMBINED_INPUT")
	}
	raw, found := rawForChunks(pipeline, chunks)
	if !found {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_RAW")
	}
	merge, found := mergeFor(pipeline, ragcontract.PortRef{NodeID: raw.ID, Port: "representations"}, ragcontract.PortRef{NodeID: combined.ID, Port: "representations"})
	if !found {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_MERGE")
	}
	embedding, found := embeddingFor(pipeline, ragcontract.PortRef{NodeID: merge.ID, Port: "representations"})
	if !found {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_EMBEDDING")
	}
	var rawConfig struct {
		Name string `json:"name"`
	}
	var combinedConfig struct {
		QuestionsPerChunk int `json:"questionsPerChunk"`
	}
	if err := json.Unmarshal(raw.Config, &rawConfig); err != nil || rawConfig.Name == "" {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_RAW_CONFIG")
	}
	if err := json.Unmarshal(combined.Config, &combinedConfig); err != nil || combinedConfig.QuestionsPerChunk < 1 {
		return CanonicalMapping{}, fmt.Errorf("RAG_PREPARATION_MAPPING_COMBINED_CONFIG")
	}
	return CanonicalMapping{
		CombinedNode: combined, EmbeddingNode: embedding, RawRepresentationName: rawConfig.Name,
		MaxRepresentationsPerChunk: 2 + combinedConfig.QuestionsPerChunk,
		RawOutputKey:               raw.ID + "/representations", DerivedOutputKey: combined.ID + "/representations",
		MergedOutputKey: merge.ID + "/representations", EmbeddingOutputKey: embedding.ID + "/embeddings",
	}, nil
}

func singleNode(pipeline ragcontract.PipelineIR, kind string) (ragcontract.Node, bool) {
	var matched ragcontract.Node
	for _, node := range pipeline.Nodes {
		if node.Operator.Kind != kind {
			continue
		}
		if matched.ID != "" {
			return ragcontract.Node{}, false
		}
		matched = node
	}
	return matched, matched.ID != ""
}

func inputFrom(node ragcontract.Node, port string) (ragcontract.PortRef, bool) {
	for _, input := range node.Inputs {
		if input.Port == port {
			return input.From, true
		}
	}
	return ragcontract.PortRef{}, false
}

func rawForChunks(pipeline ragcontract.PipelineIR, chunks ragcontract.PortRef) (ragcontract.Node, bool) {
	for _, node := range pipeline.Nodes {
		if node.Operator.Kind == "representations.raw" {
			if from, ok := inputFrom(node, "chunks"); ok && from == chunks {
				return node, true
			}
		}
	}
	return ragcontract.Node{}, false
}

func mergeFor(pipeline ragcontract.PipelineIR, left, right ragcontract.PortRef) (ragcontract.Node, bool) {
	for _, node := range pipeline.Nodes {
		if node.Operator.Kind != "representations.merge" {
			continue
		}
		seenLeft, seenRight := false, false
		for _, input := range node.Inputs {
			seenLeft = seenLeft || input.From == left
			seenRight = seenRight || input.From == right
		}
		if seenLeft && seenRight {
			return node, true
		}
	}
	return ragcontract.Node{}, false
}

func embeddingFor(pipeline ragcontract.PipelineIR, representations ragcontract.PortRef) (ragcontract.Node, bool) {
	for _, node := range pipeline.Nodes {
		if node.Operator.Kind == "embed.model" {
			if from, ok := inputFrom(node, "representations"); ok && from == representations {
				return node, true
			}
		}
	}
	return ragcontract.Node{}, false
}
