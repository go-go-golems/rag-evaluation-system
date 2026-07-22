package geppetto

import (
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

func TestPromptRenderersHaveStableKindSpecificLayouts(t *testing.T) {
	const template = "Return JSON only."
	const source = "source body"
	if got, want := renderSummaryPrompt(template, source), "Return JSON only.\n\nSOURCE TEXT:\nsource body"; got != want {
		t.Fatalf("summary prompt = %q, want %q", got, want)
	}
	if got, want := renderQuestionsPrompt(template, source), "Return JSON only.\n\nSOURCE TEXT:\nsource body"; got != want {
		t.Fatalf("questions prompt = %q, want %q", got, want)
	}
	if got, want := buildUserPrompt(template, ragoperators.GenerationRequest{Kind: "unknown", Text: source}), "Return JSON only.\n\nTEXT:\nsource body"; got != want {
		t.Fatalf("default prompt = %q, want %q", got, want)
	}
}

func TestAnswerPromptRendersQuestionAndDurableEvidenceIDs(t *testing.T) {
	got := renderAnswerPrompt("Answer only.", "What happened?", []ragoperators.Evidence{
		{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-a"}, Text: "first evidence"}},
		{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-b"}, Text: "second evidence"}},
	})
	want := "Answer only.\n\nQUESTION:\nWhat happened?\n\nEVIDENCE:\n[chunk-a]\nfirst evidence\n\n[chunk-b]\nsecond evidence\n\n"
	if got != want {
		t.Fatalf("answer prompt = %q, want %q", got, want)
	}
}
