package raglab

import (
	"context"

	"github.com/pkg/errors"
)

// ArtifactMetadata is the compatibility-relevant portion of an immutable
// artifact. It intentionally excludes mutable paths and presentation fields.
type ArtifactMetadata struct {
	Ref              ArtifactRef
	CorpusSnapshotID string
	ChunkSetID       string
	Dimensions       int
	Status           string
}

// ArtifactCatalog resolves immutable references without exposing database
// implementation details to the builder.
type ArtifactCatalog interface {
	LookupArtifact(context.Context, ArtifactRef) (ArtifactMetadata, error)
}

// ValidateCompatibility confirms that a structurally valid plan selects
// artifacts from the same immutable lineage. It performs read-only catalog
// lookup; it does not persist or execute a plan.
func (s ExperimentSpecification) ValidateCompatibility(ctx context.Context, catalog ArtifactCatalog) ValidationReport {
	var report ValidationReport
	if catalog == nil {
		report.add("RAG_CATALOG_REQUIRED", "$", "artifact catalog is required for compatibility validation")
		return report
	}
	snapshot, ok := lookup(&report, ctx, catalog, "$.inputs.corpusSnapshot", s.Inputs.CorpusSnapshot)
	chunks, chunksOK := lookup(&report, ctx, catalog, "$.inputs.chunkSet", s.Inputs.ChunkSet)
	if ok && chunksOK && chunks.CorpusSnapshotID != snapshot.Ref.ID {
		report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.chunkSet", "chunk set does not belong to the selected corpus snapshot")
	}
	if s.Inputs.BM25Index != nil {
		bm25, bm25OK := lookup(&report, ctx, catalog, "$.inputs.bm25Index", *s.Inputs.BM25Index)
		if chunksOK && bm25OK && bm25.ChunkSetID != chunks.Ref.ID {
			report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.bm25Index", "BM25 index does not belong to the selected chunk set")
		}
	}
	if s.Inputs.EmbeddingSet != nil {
		embeddings, embeddingsOK := lookup(&report, ctx, catalog, "$.inputs.embeddingSet", *s.Inputs.EmbeddingSet)
		if chunksOK && embeddingsOK && embeddings.ChunkSetID != chunks.Ref.ID {
			report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.embeddingSet", "embedding set does not belong to the selected chunk set")
		}
		if embeddingsOK && embeddings.Dimensions <= 0 {
			report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.embeddingSet", "embedding set has no positive vector dimension")
		}
	}
	dataset, datasetOK := lookup(&report, ctx, catalog, "$.inputs.evaluationDataset", s.Inputs.EvaluationDataset)
	if ok && datasetOK && dataset.CorpusSnapshotID != snapshot.Ref.ID {
		report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.evaluationDataset", "evaluation dataset is bound to a different corpus snapshot")
	}
	for index, representation := range s.Inputs.Representations {
		if representation.Kind == RawChunksRepresentation {
			continue
		}
		metadata, representationOK := lookup(&report, ctx, catalog, "$.inputs.representations["+itoa(index)+"]", RepresentationSet(representation.ArtifactID))
		if chunksOK && representationOK && metadata.ChunkSetID != chunks.Ref.ID {
			report.add("RAG_INCOMPATIBLE_ARTIFACT", "$.inputs.representations["+itoa(index)+"]", "representation set does not belong to the selected chunk set")
		}
	}
	report.Normalize()
	return report
}

func lookup(report *ValidationReport, ctx context.Context, catalog ArtifactCatalog, path string, ref ArtifactRef) (ArtifactMetadata, bool) {
	metadata, err := catalog.LookupArtifact(ctx, ref)
	if err != nil {
		if errors.Is(err, ErrArtifactNotFound) {
			report.add("RAG_UNKNOWN_ARTIFACT", path, "immutable artifact was not found")
		} else {
			report.add("RAG_CATALOG_FAILURE", path, err.Error())
		}
		return ArtifactMetadata{}, false
	}
	if metadata.Ref.Kind != ref.Kind || metadata.Ref.ID != ref.ID {
		report.add("RAG_CATALOG_FAILURE", path, "catalog returned metadata for a different artifact")
		return ArtifactMetadata{}, false
	}
	return metadata, true
}

var ErrArtifactNotFound = errors.New("immutable artifact not found")
