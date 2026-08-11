package preparationworkflow

import (
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestDeriveCanonicalMapping(t *testing.T) {
	pipeline := ragcontract.PipelineIR{Nodes: []ragcontract.Node{
		{ID: "chunks", Operator: ragcontract.OperatorRef{Kind: "chunks.identity", Version: "v1"}},
		{ID: "raw", Operator: ragcontract.OperatorRef{Kind: "representations.raw", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "chunks", From: ragcontract.PortRef{NodeID: "chunks", Port: "chunks"}}}, Config: []byte(`{"name":"raw"}`)},
		{ID: "combined", Operator: ragcontract.OperatorRef{Kind: "representations.combined-summary-questions", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "chunks", From: ragcontract.PortRef{NodeID: "chunks", Port: "chunks"}}}, Config: []byte(`{"questionsPerChunk":2}`)},
		{ID: "merge", Operator: ragcontract.OperatorRef{Kind: "representations.merge", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "raw", From: ragcontract.PortRef{NodeID: "raw", Port: "representations"}}, {Port: "derived", From: ragcontract.PortRef{NodeID: "combined", Port: "representations"}}}},
		{ID: "embed", Operator: ragcontract.OperatorRef{Kind: "embed.model", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "representations", From: ragcontract.PortRef{NodeID: "merge", Port: "representations"}}}},
	}}
	mapping, err := DeriveCanonicalMapping(pipeline)
	if err != nil {
		t.Fatal(err)
	}
	if mapping.ChunksOutputKey != "chunks/chunks" || mapping.RawOutputKey != "raw/representations" || mapping.DerivedOutputKey != "combined/representations" || mapping.MergedOutputKey != "merge/representations" || mapping.EmbeddingOutputKey != "embed/embeddings" || mapping.MaxRepresentationsPerChunk != 4 {
		t.Fatalf("mapping=%#v", mapping)
	}
}

func TestDeriveCanonicalMappingRejectsMissingMerge(t *testing.T) {
	_, err := DeriveCanonicalMapping(ragcontract.PipelineIR{Nodes: []ragcontract.Node{{ID: "combined", Operator: ragcontract.OperatorRef{Kind: "representations.combined-summary-questions", Version: "v1"}}}})
	if err == nil {
		t.Fatal("missing merge accepted")
	}
}
