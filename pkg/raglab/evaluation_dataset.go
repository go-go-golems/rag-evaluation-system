package raglab

import (
	"context"
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/pkg/errors"
)

// LoadEvaluationCards reads the immutable manifest registered for an
// evaluation dataset. The manifest is the source of truth for this first
// executor; it contains corpus-resolved document revision IDs rather than
// mutable legacy eval_queries rows.
func LoadEvaluationCards(ctx context.Context, queries *db.Queries, datasetID string) ([]EvaluationCard, error) {
	if queries == nil || datasetID == "" {
		return nil, errors.New("RAG_EVALUATION_DATASET_REQUIRED: database queries and dataset ID are required")
	}
	var raw string
	if err := queries.DB().QueryRowContext(ctx, `SELECT manifest_json FROM evaluation_datasets WHERE id=?`, datasetID).Scan(&raw); err != nil {
		return nil, errors.Wrap(err, "load immutable evaluation dataset")
	}
	var manifest struct {
		SchemaVersion string           `json:"schemaVersion"`
		Cards         []EvaluationCard `json:"cards"`
	}
	if err := json.Unmarshal([]byte(raw), &manifest); err != nil {
		return nil, errors.Wrap(err, "decode immutable evaluation dataset manifest")
	}
	if manifest.SchemaVersion != "rag-eval-evaluation-dataset/v1" || len(manifest.Cards) == 0 {
		return nil, errors.New("RAG_INVALID_EVALUATION_DATASET: manifest must use rag-eval-evaluation-dataset/v1 and contain cards")
	}
	for _, card := range manifest.Cards {
		if card.ID == "" || card.Query == "" {
			return nil, errors.New("RAG_INVALID_EVALUATION_DATASET: every card needs ID and query")
		}
	}
	SortCards(manifest.Cards)
	return manifest.Cards, nil
}
