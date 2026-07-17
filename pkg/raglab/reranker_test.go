package raglab

import (
	"reflect"
	"testing"
)

func TestRerankingPolicyIsNormalizedAndChangesAuthoringValue(t *testing.T) {
	build := func(model string) ExperimentSpecification {
		t.Helper()
		spec, err := validExperiment(t).Retrieval(func(r *RetrievalBuilder) {
			r.RerankCrossEncoder(model, 50, 20)
		}).Build()
		if err != nil {
			t.Fatal(err)
		}
		return spec
	}
	first := build("bge-reranker-v2-m3-q4_k_m")
	second := build("bge-reranker-v2-m3-q4_k_m")
	third := build("qwen3-reranker-4b-q4_k_m")
	if first.Retrieval.Reranking == nil || first.Retrieval.Reranking.CandidateCount != 50 {
		t.Fatalf("reranking = %#v", first.Retrieval.Reranking)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("same reranking policy changed normalized value: %#v != %#v", first, second)
	}
	if reflect.DeepEqual(first, third) {
		t.Fatal("reranker model must participate in exported authoring semantics")
	}
}

func TestRerankingPolicyRejectsInvalidBoundsAndDuplicates(t *testing.T) {
	builder := validExperiment(t).Retrieval(func(r *RetrievalBuilder) {
		r.RerankCrossEncoder("", 10, 5).RerankCrossEncoder("other", 20, 10)
	})
	report := builder.Validate()
	issues := map[string]bool{}
	for _, issue := range report.Issues {
		issues[issue.Code] = true
	}
	if !issues["RAG_INVALID_RERANKING"] || !issues["RAG_CONFLICTING_RERANKING"] {
		t.Fatalf("issues = %#v", report.Issues)
	}
}

func TestRerankingResultsMustCoverFinalResults(t *testing.T) {
	builder := validExperiment(t).Retrieval(func(r *RetrievalBuilder) {
		r.RerankCrossEncoder("bge", 20, 5)
	})
	if builder.Validate().OK() {
		t.Fatal("expected reranking result count validation failure")
	}
}
