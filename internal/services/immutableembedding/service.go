// Package immutableembedding builds vector artifacts from immutable chunk sets.
package immutableembedding

import (
	"context"
	"database/sql"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
	legacyembedding "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/pkg/errors"
)

const planSchema = "rag-eval-embedding-plan/v1"
const setSchema = "rag-eval-embedding-set/v1"

type Request struct {
	ChunkSetID, ProviderType, Normalization string
	Provider                                embeddings.Provider
	BatchSize                               int
}
type Result struct {
	EmbeddingPlanID, EmbeddingSetID string
	EmbeddingCount                  int
	Reused                          bool
}
type plan struct {
	SchemaVersion         string `json:"schema_version"`
	ProviderType          string `json:"provider_type"`
	Model                 string `json:"model"`
	Dimensions            int    `json:"dimensions"`
	Normalization         string `json:"normalization"`
	ImplementationVersion string `json:"implementation_version"`
}
type set struct {
	SchemaVersion   string   `json:"schema_version"`
	ChunkSetID      string   `json:"chunk_set_id"`
	EmbeddingPlanID string   `json:"embedding_plan_id"`
	ChunkIDs        []string `json:"chunk_ids"`
}
type chunk struct{ ID, Text string }

func Build(ctx context.Context, queries *db.Queries, request Request) (*Result, error) {
	if queries == nil || request.Provider == nil {
		return nil, errors.New("database queries and embedding provider are required")
	}
	if request.ChunkSetID == "" {
		return nil, errors.New("chunk set ID is required")
	}
	if request.ProviderType == "" {
		request.ProviderType = "unknown"
	}
	if request.Normalization == "" {
		request.Normalization = "provider-default/v1"
	}
	if request.BatchSize <= 0 {
		request.BatchSize = 16
	}
	model := request.Provider.GetModel()
	if model.Name == "" || model.Dimensions <= 0 {
		return nil, errors.New("embedding provider returned invalid model metadata")
	}
	p := plan{planSchema, request.ProviderType, model.Name, model.Dimensions, request.Normalization, "geppetto-embeddings/v1"}
	planID, err := experiments.Fingerprint(planSchema, p)
	if err != nil {
		return nil, err
	}
	pjson, err := experiments.CanonicalJSON(p)
	if err != nil {
		return nil, err
	}
	chunks, err := load(ctx, queries.DB(), request.ChunkSetID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(chunks))
	for i, c := range chunks {
		ids[i] = c.ID
	}
	s := set{setSchema, request.ChunkSetID, planID, ids}
	setID, err := experiments.Fingerprint(setSchema, s)
	if err != nil {
		return nil, err
	}
	manifest, err := experiments.CanonicalJSON(s)
	if err != nil {
		return nil, err
	}
	vectors := make([][]float32, 0, len(chunks))
	for start := 0; start < len(chunks); start += request.BatchSize {
		end := start + request.BatchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		texts := make([]string, end-start)
		for i, c := range chunks[start:end] {
			texts[i] = c.Text
		}
		batch, err := request.Provider.GenerateBatchEmbeddings(ctx, texts)
		if err != nil {
			return nil, errors.Wrap(err, "generate immutable embeddings")
		}
		if len(batch) != len(texts) {
			return nil, errors.New("embedding provider returned wrong vector count")
		}
		for _, v := range batch {
			if len(v) != model.Dimensions {
				return nil, errors.New("embedding provider returned wrong vector dimensions")
			}
			vectors = append(vectors, v)
		}
	}
	tx, err := queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	if err := ensurePlan(ctx, tx, planID, p, string(pjson)); err != nil {
		return nil, err
	}
	reused, err := ensureSet(ctx, tx, setID, request.ChunkSetID, planID, string(manifest), chunks, vectors)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &Result{planID, setID, len(chunks), reused}, nil
}
func load(ctx context.Context, database *sql.DB, setID string) ([]chunk, error) {
	rows, err := database.QueryContext(ctx, `SELECT id,text FROM immutable_chunks WHERE chunk_set_id=? ORDER BY document_revision_id,chunk_index`, setID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var r []chunk
	for rows.Next() {
		var c chunk
		if err := rows.Scan(&c.ID, &c.Text); err != nil {
			return nil, err
		}
		r = append(r, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, errors.New("chunk set has no chunks")
	}
	return r, nil
}
func ensurePlan(ctx context.Context, tx *sql.Tx, id string, p plan, json string) error {
	var existing plan
	err := tx.QueryRowContext(ctx, `SELECT schema_version, provider_type, model, dimensions, normalization, implementation_version FROM embedding_plans WHERE id=?`, id).Scan(&existing.SchemaVersion, &existing.ProviderType, &existing.Model, &existing.Dimensions, &existing.Normalization, &existing.ImplementationVersion)
	if err == nil {
		existingJSON, marshalErr := experiments.CanonicalJSON(existing)
		if marshalErr != nil {
			return marshalErr
		}
		if string(existingJSON) != json {
			return errors.New("embedding plan conflicts")
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO embedding_plans (id,schema_version,provider_type,model,dimensions,normalization,implementation_version) VALUES (?,?,?,?,?,?,?)`, id, p.SchemaVersion, p.ProviderType, p.Model, p.Dimensions, p.Normalization, p.ImplementationVersion)
	return err
}
func ensureSet(ctx context.Context, tx *sql.Tx, id, chunkSet, planID, manifest string, chunks []chunk, vectors [][]float32) (bool, error) {
	var existing string
	var count int
	err := tx.QueryRowContext(ctx, `SELECT manifest_json,embedding_count FROM embedding_sets WHERE id=?`, id).Scan(&existing, &count)
	if err == nil {
		if existing != manifest || count != len(chunks) {
			return false, errors.New("embedding set conflicts")
		}
		return true, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO embedding_sets (id,chunk_set_id,embedding_plan_id,manifest_json,embedding_count) VALUES (?,?,?,?,?)`, id, chunkSet, planID, manifest, len(chunks)); err != nil {
		return false, err
	}
	for i, c := range chunks {
		if _, err = tx.ExecContext(ctx, `INSERT INTO immutable_embeddings (embedding_set_id,chunk_id,vector) VALUES (?,?,?)`, id, c.ID, legacyembedding.EncodeFloat32Vector(vectors[i])); err != nil {
			return false, err
		}
	}
	return false, nil
}
