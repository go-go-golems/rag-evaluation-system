package raglab

import (
	"context"
	"testing"
)

type testCatalog map[ArtifactRef]ArtifactMetadata

func (c testCatalog) LookupArtifact(_ context.Context, ref ArtifactRef) (ArtifactMetadata, error) {
	metadata, ok := c[ref]
	if !ok {
		return ArtifactMetadata{}, ErrArtifactNotFound
	}
	return metadata, nil
}

func TestValidateCompatibilityAcceptsSingleArtifactLineage(t *testing.T) {
	spec, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	catalog := testCatalog{
		CorpusSnapshot("snapshot"): {Ref: CorpusSnapshot("snapshot")},
		ChunkSet("chunks"):         {Ref: ChunkSet("chunks"), CorpusSnapshotID: "snapshot"},
		BM25Index("bm25"):          {Ref: BM25Index("bm25"), ChunkSetID: "chunks"},
		EmbeddingSet("embeddings"): {Ref: EmbeddingSet("embeddings"), ChunkSetID: "chunks", Dimensions: 768},
		EvaluationDataset("eval"):  {Ref: EvaluationDataset("eval"), CorpusSnapshotID: "snapshot", Status: "candidate"},
	}
	if report := spec.ValidateCompatibility(context.Background(), catalog); !report.OK() {
		t.Fatalf("compatibility report = %#v", report)
	}
}

func TestValidateCompatibilityReportsAllLineageProblems(t *testing.T) {
	spec, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	catalog := testCatalog{
		CorpusSnapshot("snapshot"): {Ref: CorpusSnapshot("snapshot")},
		ChunkSet("chunks"):         {Ref: ChunkSet("chunks"), CorpusSnapshotID: "other-snapshot"},
		BM25Index("bm25"):          {Ref: BM25Index("bm25"), ChunkSetID: "other-chunks"},
		EmbeddingSet("embeddings"): {Ref: EmbeddingSet("embeddings"), ChunkSetID: "chunks", Dimensions: 0},
		EvaluationDataset("eval"):  {Ref: EvaluationDataset("eval"), CorpusSnapshotID: "other-snapshot"},
	}
	report := spec.ValidateCompatibility(context.Background(), catalog)
	if report.OK() {
		t.Fatalf("expected incompatibilities: %#v", report)
	}
	issues := map[string]int{}
	for _, issue := range report.Issues {
		issues[issue.Code]++
	}
	if issues["RAG_INCOMPATIBLE_ARTIFACT"] != 4 {
		t.Fatalf("incompatibility count = %d, report=%#v", issues["RAG_INCOMPATIBLE_ARTIFACT"], report.Issues)
	}
}

func TestValidateCompatibilityReportsMissingArtifact(t *testing.T) {
	spec, err := validExperiment(t).Build()
	if err != nil {
		t.Fatal(err)
	}
	report := spec.ValidateCompatibility(context.Background(), testCatalog{})
	if report.OK() || len(report.Issues) != 5 {
		t.Fatalf("report = %#v", report)
	}
	for _, issue := range report.Issues {
		if issue.Code != "RAG_UNKNOWN_ARTIFACT" {
			t.Fatalf("issue = %#v", issue)
		}
	}
}
