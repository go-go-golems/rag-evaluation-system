package raglab

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestExportSpecificationV1MatchesResearchctlGolden(t *testing.T) {
	grade, err := Grade("2_SUBSTANTIAL")
	if err != nil {
		t.Fatal(err)
	}
	prototype, err := NewExperiment("ttc-raw-hybrid").
		Corpus(CorpusSnapshot("corpus")).Chunks(ChunkSet("chunks")).BM25(BM25Index("bm25")).Embeddings(EmbeddingSet("embeddings")).Evaluation(EvaluationDataset("evaluation")).
		Tag("corpus", "ttc").Tag("status", "candidate").
		Representations(func(builder *RepresentationBuilder) { builder.RawChunks("raw") }).
		Retrieval(func(builder *RetrievalBuilder) {
			builder.Channel("lexical", func(channel *ChannelBuilder) { channel.BM25().Representation("raw").TopK(50) })
			builder.Channel("semantic", func(channel *ChannelBuilder) { channel.Vector().Representation("raw").TopK(50) })
			builder.FuseRRF(60).Weight("semantic", 1.25).Collapse(CollapseDocument).Results(10)
		}).Metrics(func(metrics *MetricsBuilder) {
		metrics.RelevanceAt(grade).RecallAt(10).PrecisionAt(10, 1, 3).NDCGAt(10).MRR()
	}).Build()
	if err != nil {
		t.Fatal(err)
	}
	exported, err := ExportSpecificationV1(prototype, ExportOptions{DatasetSplit: "development"})
	if err != nil {
		t.Fatal(err)
	}
	goldenData, err := os.ReadFile("testdata/researchctl-rag-domain-specification.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden PrototypeSpecification
	if err := json.Unmarshal(goldenData, &golden); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(exported, golden) {
		actual, _ := json.MarshalIndent(exported, "", "  ")
		t.Fatalf("cross-repository contract drift:\n%s\nwant:\n%s", actual, goldenData)
	}
	if prototype.Inputs.CorpusSnapshot.ID != "corpus" || exported.Name != prototype.Name {
		t.Fatalf("export mutated prototype or lost name: prototype=%+v export=%+v", prototype, exported)
	}
}

func TestExportSpecificationV1RequiresExplicitSplitAndSupportedMetrics(t *testing.T) {
	prototype := ExperimentSpecification{Metrics: MetricsPlan{RelevanceAt: &RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}, MRR: true}}
	if _, err := ExportSpecificationV1(prototype, ExportOptions{}); err == nil || !strings.Contains(err.Error(), "RAG_DATASET_SPLIT_REQUIRED") {
		t.Fatalf("expected explicit split error, got %v", err)
	}
	prototype.Metrics.HitRateAt = []int{10}
	if _, err := ExportSpecificationV1(prototype, ExportOptions{DatasetSplit: "development"}); err == nil || !strings.Contains(err.Error(), "RAG_EXPORT_UNSUPPORTED") {
		t.Fatalf("expected unsupported metric error, got %v", err)
	}
}
