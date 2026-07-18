package geppetto

import (
	"context"
	"strings"
	"testing"

	geppettorerank "github.com/go-go-golems/geppetto/pkg/rerank"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type fakeRerankProvider struct {
	model    geppettorerank.Model
	response geppettorerank.Response
	err      error
	request  geppettorerank.Request
}

func (p *fakeRerankProvider) Rerank(_ context.Context, request geppettorerank.Request) (geppettorerank.Response, error) {
	p.request = request
	return p.response, p.err
}

func (p *fakeRerankProvider) Model() geppettorerank.Model { return p.model }

func TestRerankerForcesCompleteCoverageAndMapsChunkIDs(t *testing.T) {
	provider := &fakeRerankProvider{
		model: geppettorerank.Model{Provider: "llama.cpp", Name: "bge-exact"},
		response: geppettorerank.Response{
			Provider: "llama.cpp",
			Model:    "bge-exact",
			Results: []geppettorerank.Result{
				{DocumentID: "chunk-2", Index: 1, Score: -2.5, Rank: 1},
				{DocumentID: "chunk-1", Index: 0, Score: -3.5, Rank: 2},
			},
		},
	}
	adapter, err := NewReranker(provider)
	if err != nil {
		t.Fatalf("NewReranker() error = %v", err)
	}

	scores, err := adapter.Rerank(context.Background(), ragoperators.RerankRequest{
		Model:   "bge-exact",
		Query:   "How are payroll adjustments handled?",
		Results: 1,
		Candidates: []ragoperators.Evidence{
			{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-1"}, Text: "A payroll adjustment corrects wages."}},
			{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-2"}, Text: "Cypress trees tolerate dry conditions."}},
		},
	})
	if err != nil {
		t.Fatalf("Rerank() error = %v", err)
	}
	if got, want := provider.request.TopN, 2; got != want {
		t.Errorf("provider TopN = %d, want complete candidate count %d", got, want)
	}
	if got, want := provider.request.Documents[0].ID, "chunk-1"; got != want {
		t.Errorf("first provider document ID = %q, want %q", got, want)
	}
	if got, want := len(scores), 2; got != want {
		t.Fatalf("score count = %d, want %d", got, want)
	}
	if scores[0] != (ragoperators.RerankScore{ChunkID: "chunk-2", Score: -2.5}) {
		t.Errorf("first score = %#v, want chunk-2's raw negative score", scores[0])
	}
	if scores[1] != (ragoperators.RerankScore{ChunkID: "chunk-1", Score: -3.5}) {
		t.Errorf("second score = %#v, want chunk-1's raw negative score", scores[1])
	}
}

func TestRerankerRejectsIncompleteOrUnknownProviderResults(t *testing.T) {
	for name, results := range map[string][]geppettorerank.Result{
		"incomplete": {{DocumentID: "chunk-1", Score: 1}},
		"unknown id": {{DocumentID: "not-a-chunk", Score: 1}, {DocumentID: "chunk-2", Score: 0}},
	} {
		t.Run(name, func(t *testing.T) {
			provider := &fakeRerankProvider{response: geppettorerank.Response{Model: "bge-exact", Results: results}}
			adapter, err := NewReranker(provider)
			if err != nil {
				t.Fatal(err)
			}
			_, err = adapter.Rerank(context.Background(), ragoperators.RerankRequest{
				Model: "bge-exact", Query: "payroll",
				Candidates: []ragoperators.Evidence{
					{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-1"}, Text: "payroll"}},
					{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-2"}, Text: "trees"}},
				},
			})
			if err == nil {
				t.Fatal("Rerank() error = nil, want malformed provider response rejection")
			}
			if !strings.Contains(err.Error(), "RAG_GEPPETTO_RERANK") {
				t.Errorf("Rerank() error = %v, want stable RAG rerank error", err)
			}
		})
	}
}
