package raglab

import (
	"errors"
	"testing"
)

func substantial(t *testing.T) RelevanceGrade {
	t.Helper()
	grade, err := Grade("2_SUBSTANTIAL")
	if err != nil {
		t.Fatal(err)
	}
	return grade
}

func validExperiment(t *testing.T) *ExperimentBuilder {
	t.Helper()
	return NewExperiment("ttc-hybrid").
		Corpus(CorpusSnapshot("snapshot")).
		Chunks(ChunkSet("chunks")).
		BM25(BM25Index("bm25")).
		Embeddings(EmbeddingSet("embeddings")).
		Evaluation(EvaluationDataset("eval")).
		Retrieval(func(r *RetrievalBuilder) {
			r.Channel("lexical", func(c *ChannelBuilder) { c.BM25().TopK(50) }).
				Channel("semantic", func(c *ChannelBuilder) { c.Vector().TopK(50) }).
				FuseRRF(60).Weight("semantic", 1.25).Collapse(CollapseDocument).Results(10)
		}).
		Metrics(func(m *MetricsBuilder) { m.RelevanceAt(substantial(t)).RecallAt(10, 1, 3, 3).MRR() })
}

func TestBuildProducesStablePersistenceFingerprint(t *testing.T) {
	first, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	second, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	if first.SchemaVersion == "" || first.Fingerprint == "" {
		t.Fatalf("incomplete specification: %#v", first)
	}
	if first.Fingerprint != second.Fingerprint {
		t.Fatalf("fingerprint = %s, want %s", first.Fingerprint, second.Fingerprint)
	}
	if got := first.Metrics.RecallAt; len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 10 {
		t.Fatalf("recall cutoffs = %#v", got)
	}
	input := first.PersistenceInput()
	if input.CorpusSnapshotID != "snapshot" || input.EmbeddingSetID != "embeddings" || input.Config["retrieval"] == nil {
		t.Fatalf("persistence input = %#v", input)
	}
}

func TestUseFragmentIsIdempotentAndDetectsInputConflict(t *testing.T) {
	fragment := NewFragment("ttc-inputs", func(e *ExperimentBuilder) {
		e.Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).Evaluation(EvaluationDataset("eval"))
	})
	builder := NewExperiment("fragment").Use(fragment).Use(fragment).Corpus(CorpusSnapshot("other"))
	report := builder.Validate()
	if report.OK() {
		t.Fatalf("expected conflict report: %#v", report)
	}
	count := 0
	for _, issue := range report.Issues {
		if issue.Code == "RAG_CONFLICTING_FRAGMENT" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("conflicts = %d, report=%#v", count, report)
	}
}

func TestBuildReportsStructuralFailures(t *testing.T) {
	builder := NewExperiment("").
		Corpus(CorpusSnapshot("snapshot")).
		Chunks(ChunkSet("chunks")).
		Evaluation(EvaluationDataset("eval")).
		Retrieval(func(r *RetrievalBuilder) {
			r.Channel("semantic", func(c *ChannelBuilder) { c.Vector().TopK(2) }).Results(10)
		}).
		Metrics(func(m *MetricsBuilder) { m.MRR() })
	_, err := builder.Build()
	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("Build error = %v, want ValidationError", err)
	}
	want := map[string]bool{"RAG_INVALID_NAME": false, "RAG_MISSING_EMBEDDINGS": false, "RAG_INVALID_CHANNEL": false, "RAG_MISSING_RELEVANCE_THRESHOLD": false}
	for _, issue := range validation.Report.Issues {
		if _, ok := want[issue.Code]; ok {
			want[issue.Code] = true
		}
	}
	for code, found := range want {
		if !found {
			t.Errorf("missing %s in %#v", code, validation.Report.Issues)
		}
	}
}

func TestBuildRejectsUnknownRepresentationAndUnknownFusionWeight(t *testing.T) {
	builder := NewExperiment("invalid-representation").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Evaluation(EvaluationDataset("eval")).
		Retrieval(func(r *RetrievalBuilder) {
			r.Channel("lexical", func(c *ChannelBuilder) { c.BM25().Representation("missing").TopK(10) }).Channel("again", func(c *ChannelBuilder) { c.BM25().TopK(10) }).FuseRRF(60).Weight("unknown", 1).Results(10)
		}).
		Metrics(func(m *MetricsBuilder) { m.RelevanceAt(substantial(t)).MRR() })
	report := builder.Validate()
	issues := map[string]bool{}
	for _, issue := range report.Issues {
		issues[issue.Code] = true
	}
	if !issues["RAG_UNKNOWN_REPRESENTATION"] || !issues["RAG_INVALID_FUSION"] {
		t.Fatalf("issues = %#v", report.Issues)
	}
}

func TestBuildCanonicalizesSetLikeFiltersButKeepsFragmentProvenanceOrder(t *testing.T) {
	build := func(ids ...string) ExperimentSpecification {
		t.Helper()
		first := NewFragment("first", func(*ExperimentBuilder) {})
		second := NewFragment("second", func(*ExperimentBuilder) {})
		spec, err := NewExperiment("filtered").
			Use(first).Use(second).
			Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Evaluation(EvaluationDataset("eval")).
			Retrieval(func(r *RetrievalBuilder) {
				r.Filter(func(f *FilterBuilder) {
					f.SourceIDs(ids...).ContentTypes("article", "product", "article").MetadataEquals("locale", "en_US")
				}).Channel("lexical", func(c *ChannelBuilder) { c.BM25().TopK(10) }).Results(10)
			}).
			Metrics(func(m *MetricsBuilder) { m.RelevanceAt(substantial(t)).MRR() }).Build()
		if err != nil {
			t.Fatal(err)
		}
		return spec
	}
	first := build("b", "a", "b")
	second := build("a", "b")
	if first.Fingerprint != second.Fingerprint {
		t.Fatalf("set-like filters changed fingerprint: %s != %s", first.Fingerprint, second.Fingerprint)
	}
	if got := first.Provenance.Fragments; len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Fatalf("fragment provenance order = %#v", got)
	}
}

func TestConflictingMetadataAndCollapseAreStructuralErrors(t *testing.T) {
	builder := NewExperiment("invalid-filter").
		Corpus(CorpusSnapshot("snapshot")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Evaluation(EvaluationDataset("eval")).
		Retrieval(func(r *RetrievalBuilder) {
			r.Filter(func(f *FilterBuilder) {
				f.MetadataEquals("locale", "en_US").MetadataEquals("locale", "fr_FR")
			}).Channel("lexical", func(c *ChannelBuilder) { c.BM25().TopK(10) }).Collapse(CollapseScope("bad")).Results(10)
		}).
		Metrics(func(m *MetricsBuilder) { m.RelevanceAt(substantial(t)).MRR() })
	report := builder.Validate()
	issues := map[string]bool{}
	for _, issue := range report.Issues {
		issues[issue.Code] = true
	}
	if !issues["RAG_CONFLICTING_FILTER"] || !issues["RAG_INVALID_COLLAPSE"] {
		t.Fatalf("issues = %#v", report.Issues)
	}
}
