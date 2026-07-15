// Package experimentspec defines the immutable experiment specification
// manifest shared by authoring surfaces and the persistence service.
package experimentspec

import (
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
)

// SchemaVersion identifies the persisted experiment-specification contract.
// An incompatible meaning change requires a new schema version rather than a
// reinterpretation of existing immutable rows.
const SchemaVersion = "rag-eval-experiment-spec/v1"

// Input identifies immutable inputs and carries the retrieval/evaluation plan
// as JSON-compatible configuration data.
type Input struct {
	CorpusSnapshotID    string         `json:"corpus_snapshot_id"`
	ChunkSetID          string         `json:"chunk_set_id"`
	BM25ArtifactID      string         `json:"bm25_artifact_id,omitempty"`
	EmbeddingSetID      string         `json:"embedding_set_id,omitempty"`
	EvaluationDatasetID string         `json:"evaluation_dataset_id,omitempty"`
	Config              map[string]any `json:"config"`
}

// Manifest is the complete content-addressed payload for a persisted
// experiment specification. It deliberately has no creation timestamp or run
// metadata, which would make repeated semantic plans hash differently.
type Manifest struct {
	SchemaVersion       string         `json:"schema_version"`
	CorpusSnapshotID    string         `json:"corpus_snapshot_id"`
	ChunkSetID          string         `json:"chunk_set_id"`
	BM25ArtifactID      string         `json:"bm25_artifact_id,omitempty"`
	EmbeddingSetID      string         `json:"embedding_set_id,omitempty"`
	EvaluationDatasetID string         `json:"evaluation_dataset_id,omitempty"`
	Config              map[string]any `json:"config"`
}

// Normalize turns omitted configuration into the same empty-object value used
// by persisted specifications. It does not mutate a caller's map.
func Normalize(input Input) Input {
	if input.Config == nil {
		input.Config = map[string]any{}
	}
	return input
}

// NewManifest produces the stable persisted payload for input.
func NewManifest(input Input) Manifest {
	input = Normalize(input)
	return Manifest{
		SchemaVersion:       SchemaVersion,
		CorpusSnapshotID:    input.CorpusSnapshotID,
		ChunkSetID:          input.ChunkSetID,
		BM25ArtifactID:      input.BM25ArtifactID,
		EmbeddingSetID:      input.EmbeddingSetID,
		EvaluationDatasetID: input.EvaluationDatasetID,
		Config:              input.Config,
	}
}

// Fingerprint returns the schema-namespaced identity of input.
func Fingerprint(input Input) (string, error) {
	return experiments.Fingerprint(SchemaVersion, NewManifest(input))
}
