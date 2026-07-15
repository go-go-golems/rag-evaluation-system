// Package immutableretrieval provides retrieval over content-addressed artifacts.
package immutableretrieval

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
	"github.com/pkg/errors"
)

const bm25Schema = "rag-eval-bm25-artifact/v1"

type ChunkHit struct {
	Rank               int     `json:"rank"`
	ChunkID            string  `json:"chunk_id"`
	DocumentRevisionID string  `json:"document_revision_id"`
	Title              string  `json:"title"`
	URL                string  `json:"url"`
	Kind               string  `json:"kind"`
	ChunkIndex         int     `json:"chunk_index"`
	StartRunes         int     `json:"start_runes"`
	EndRunes           int     `json:"end_runes"`
	Score              float64 `json:"score"`
	Text               string  `json:"text"`
	Channel            string  `json:"channel"`
}
type BM25BuildRequest struct{ ChunkSetID, ArtifactRoot string }
type BM25BuildResult struct {
	ArtifactID, Path string
	ChunkCount       int
	Reused           bool
}
type bm25Config struct {
	SchemaVersion  string `json:"schema_version"`
	Analyzer       string `json:"analyzer"`
	Implementation string `json:"implementation"`
}
type bm25Manifest struct {
	SchemaVersion string     `json:"schema_version"`
	ChunkSetID    string     `json:"chunk_set_id"`
	Config        bm25Config `json:"config"`
	ChunkIDs      []string   `json:"chunk_ids"`
}
type indexed struct {
	ChunkID            string `json:"chunk_id"`
	DocumentRevisionID string `json:"document_revision_id"`
	Title              string `json:"title"`
	URL                string `json:"url"`
	Kind               string `json:"kind"`
	ChunkIndex         int    `json:"chunk_index"`
	StartRunes         int    `json:"start_runes"`
	EndRunes           int    `json:"end_runes"`
	Text               string `json:"text"`
}

func BuildBM25(ctx context.Context, q *db.Queries, req BM25BuildRequest) (*BM25BuildResult, error) {
	if q == nil || req.ChunkSetID == "" {
		return nil, errors.New("database queries and chunk set ID are required")
	}
	if req.ArtifactRoot == "" {
		req.ArtifactRoot = "data/artifacts/bm25"
	}
	rows, err := q.DB().QueryContext(ctx, `SELECT ic.id,ic.document_revision_id,dr.title,dr.url,dr.kind,ic.chunk_index,ic.source_start_runes,ic.source_end_runes,ic.text FROM immutable_chunks ic JOIN document_revisions dr ON dr.id=ic.document_revision_id WHERE ic.chunk_set_id=? ORDER BY ic.document_revision_id,ic.chunk_index`, req.ChunkSetID)
	if err != nil {
		return nil, errors.Wrap(err, "load immutable chunks")
	}
	defer rows.Close()
	var docs []indexed
	var ids []string
	for rows.Next() {
		var d indexed
		if err := rows.Scan(&d.ChunkID, &d.DocumentRevisionID, &d.Title, &d.URL, &d.Kind, &d.ChunkIndex, &d.StartRunes, &d.EndRunes, &d.Text); err != nil {
			return nil, err
		}
		docs = append(docs, d)
		ids = append(ids, d.ChunkID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(docs) == 0 {
		return nil, errors.New("chunk set has no chunks")
	}
	cfg := bm25Config{bm25Schema, "standard", "bleve/v2-standard/v1"}
	manifest := bm25Manifest{bm25Schema, req.ChunkSetID, cfg, ids}
	id, err := experiments.Fingerprint(bm25Schema, manifest)
	if err != nil {
		return nil, err
	}
	configJSON, err := experiments.CanonicalJSON(cfg)
	if err != nil {
		return nil, err
	}
	manifestJSON, err := experiments.CanonicalJSON(manifest)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(req.ArtifactRoot, id)
	var existing string
	err = q.DB().QueryRowContext(ctx, `SELECT manifest_json FROM retrieval_artifacts WHERE id=?`, id).Scan(&existing)
	if err == nil {
		if existing != string(manifestJSON) {
			return nil, errors.New("immutable BM25 artifact conflict")
		}
		return &BM25BuildResult{id, path, len(docs), true}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err := os.MkdirAll(req.ArtifactRoot, 0o750); err != nil {
		return nil, err
	}
	tmp := path + ".tmp"
	_ = os.RemoveAll(tmp)
	index, err := bleve.New(tmp, bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}
	batch := index.NewBatch()
	for _, d := range docs {
		if err := batch.Index(d.ChunkID, d); err != nil {
			_ = index.Close()
			return nil, err
		}
	}
	if err := index.Batch(batch); err != nil {
		_ = index.Close()
		return nil, err
	}
	if err := index.Close(); err != nil {
		return nil, err
	}
	if err := os.Rename(tmp, path); err != nil {
		return nil, err
	}
	_, err = q.DB().ExecContext(ctx, `INSERT INTO retrieval_artifacts (id,schema_version,kind,chunk_set_id,config_json,manifest_json,artifact_path,chunk_count) VALUES (?,?,?,?,?,?,?,?)`, id, bm25Schema, "bm25", req.ChunkSetID, string(configJSON), string(manifestJSON), path, len(docs))
	if err != nil {
		return nil, errors.Wrap(err, "record BM25 artifact")
	}
	return &BM25BuildResult{id, path, len(docs), false}, nil
}

func QueryBM25(ctx context.Context, q *db.Queries, artifactID, query string, limit int) ([]ChunkHit, error) {
	if query == "" {
		return nil, errors.New("query is required")
	}
	if limit <= 0 {
		limit = 10
	}
	var path string
	if err := q.DB().QueryRowContext(ctx, `SELECT artifact_path FROM retrieval_artifacts WHERE id=? AND kind='bm25'`, artifactID).Scan(&path); err != nil {
		return nil, errors.Wrap(err, "look up BM25 artifact")
	}
	idx, err := bleve.Open(path)
	if err != nil {
		return nil, err
	}
	defer idx.Close()
	text := bleve.NewMatchQuery(query)
	text.SetField("text")
	title := bleve.NewMatchQuery(query)
	title.SetField("title")
	title.SetBoost(2)
	r, err := idx.SearchInContext(ctx, bleve.NewSearchRequestOptions(bleve.NewDisjunctionQuery(text, title), limit, 0, false))
	if err != nil {
		return nil, err
	}
	hits := make([]ChunkHit, 0, len(r.Hits))
	for i, h := range r.Hits {
		hits = append(hits, ChunkHit{Rank: i + 1, ChunkID: h.ID, Score: h.Score, Channel: "bm25"})
	}
	return hydrate(ctx, q, hits, nil)
}

func hydrate(ctx context.Context, q *db.Queries, hits []ChunkHit, scores map[string]float64) ([]ChunkHit, error) {
	for i := range hits {
		var d indexed
		err := q.DB().QueryRowContext(ctx, `SELECT ic.id,ic.document_revision_id,dr.title,dr.url,dr.kind,ic.chunk_index,ic.source_start_runes,ic.source_end_runes,ic.text FROM immutable_chunks ic JOIN document_revisions dr ON dr.id=ic.document_revision_id WHERE ic.id=?`, hits[i].ChunkID).Scan(&d.ChunkID, &d.DocumentRevisionID, &d.Title, &d.URL, &d.Kind, &d.ChunkIndex, &d.StartRunes, &d.EndRunes, &d.Text)
		if err != nil {
			return nil, fmt.Errorf("hydrate immutable chunk %s: %w", hits[i].ChunkID, err)
		}
		hits[i].DocumentRevisionID = d.DocumentRevisionID
		hits[i].Title = d.Title
		hits[i].URL = d.URL
		hits[i].Kind = d.Kind
		hits[i].ChunkIndex = d.ChunkIndex
		hits[i].StartRunes = d.StartRunes
		hits[i].EndRunes = d.EndRunes
		hits[i].Text = d.Text
		if scores != nil {
			hits[i].Score = scores[hits[i].ChunkID]
		}
	}
	return hits, nil
}
