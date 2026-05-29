package db

import (
	"database/sql"
	"fmt"
)

// ChunkWithDocument joins a chunk with the document/source fields needed for
// search indexing and result rendering.
type ChunkWithDocument struct {
	ChunkID     string `json:"chunk_id"`
	DocumentID  string `json:"document_id"`
	SourceID    string `json:"source_id"`
	Title       string `json:"title"`
	URL         string `json:"url,omitempty"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	Text        string `json:"text,omitempty"`
	TokenCount  int    `json:"token_count"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}

// SearchIndex is metadata for a disposable derived search index.
type SearchIndex struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	StrategyID    string `json:"strategy_id"`
	Provider      string `json:"provider,omitempty"`
	Model         string `json:"model,omitempty"`
	Dimensions    int    `json:"dimensions,omitempty"`
	IndexType     string `json:"index_type"`
	IndexPath     string `json:"index_path"`
	DocumentCount int    `json:"document_count"`
	ChunkCount    int    `json:"chunk_count"`
	LastRebuildAt string `json:"last_rebuild_at,omitempty"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

// ChunkEmbeddingWithContext is a stored embedding plus document/source context
// needed to render retrieval results.
type ChunkEmbeddingWithContext struct {
	ChunkID    string `json:"chunk_id"`
	DocumentID string `json:"document_id"`
	SourceID   string `json:"source_id"`
	Title      string `json:"title"`
	URL        string `json:"url,omitempty"`
	StrategyID string `json:"strategy_id"`
	ChunkIndex int    `json:"chunk_index"`
	Text       string `json:"text,omitempty"`
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions"`
	TextHash   string `json:"text_hash"`
	Embedding  []byte `json:"-"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ListChunksWithDocumentContext returns chunks for a strategy, optionally
// restricted to source IDs, with document metadata needed for search indexing.
func (q *Queries) ListChunksWithDocumentContext(strategyID string, sourceIDs []string, limit int) ([]ChunkWithDocument, error) {
	return q.ListChunksWithDocumentContextFiltered(strategyID, sourceIDs, nil, limit)
}

// ListChunksWithDocumentContextFiltered returns chunks for a strategy, optionally
// restricted to source IDs and document IDs, with document metadata needed for search indexing.
func (q *Queries) ListChunksWithDocumentContextFiltered(strategyID string, sourceIDs []string, documentIDs []string, limit int) ([]ChunkWithDocument, error) {
	query := `
		SELECT c.id, c.document_id, d.source_id, d.title, COALESCE(d.url, ''),
		       c.strategy_id, c.chunk_index, c.text, c.token_count,
		       COALESCE(c.start_offset, 0), COALESCE(c.end_offset, 0)
		FROM chunks c
		JOIN documents d ON d.id = c.document_id
		WHERE c.strategy_id = ?
	`
	args := []interface{}{strategyID}
	if len(sourceIDs) > 0 {
		query += ` AND d.source_id IN (` + placeholders(len(sourceIDs)) + `)`
		for _, sourceID := range sourceIDs {
			args = append(args, sourceID)
		}
	}
	if len(documentIDs) > 0 {
		query += ` AND c.document_id IN (` + placeholders(len(documentIDs)) + `)`
		for _, documentID := range documentIDs {
			args = append(args, documentID)
		}
	}
	query += ` ORDER BY d.source_id, c.document_id, c.chunk_index`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list chunks with document context: %w", err)
	}
	defer rows.Close()

	var chunks []ChunkWithDocument
	for rows.Next() {
		var c ChunkWithDocument
		if err := rows.Scan(&c.ChunkID, &c.DocumentID, &c.SourceID, &c.Title, &c.URL,
			&c.StrategyID, &c.ChunkIndex, &c.Text, &c.TokenCount,
			&c.StartOffset, &c.EndOffset); err != nil {
			return nil, fmt.Errorf("scan chunk with document context: %w", err)
		}
		chunks = append(chunks, c)
	}
	return chunks, rows.Err()
}

// ListChunkEmbeddingsForStrategySourcesWithContext returns stored chunk embeddings
// with document context, optionally restricted to source IDs.
func (q *Queries) ListChunkEmbeddingsForStrategySourcesWithContext(strategyID string, sourceIDs []string, provider, model string, dimensions, limit int) ([]ChunkEmbeddingWithContext, error) {
	query := `
		SELECT ce.chunk_id, c.document_id, d.source_id, d.title, COALESCE(d.url, ''),
		       ce.strategy_id, c.chunk_index, c.text,
		       ce.provider, ce.model, ce.dimensions, ce.text_hash, ce.embedding, ce.created_at, ce.updated_at
		FROM chunk_embeddings ce
		JOIN chunks c ON c.id = ce.chunk_id AND c.strategy_id = ce.strategy_id
		JOIN documents d ON d.id = c.document_id
		WHERE ce.strategy_id = ? AND ce.provider = ? AND ce.model = ? AND ce.dimensions = ?
	`
	args := []interface{}{strategyID, provider, model, dimensions}
	if len(sourceIDs) > 0 {
		query += ` AND d.source_id IN (` + placeholders(len(sourceIDs)) + `)`
		for _, sourceID := range sourceIDs {
			args = append(args, sourceID)
		}
	}
	query += ` ORDER BY d.source_id, c.document_id, c.chunk_index`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}

	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list chunk embeddings with context: %w", err)
	}
	defer rows.Close()

	var embeddings []ChunkEmbeddingWithContext
	for rows.Next() {
		var e ChunkEmbeddingWithContext
		if err := rows.Scan(
			&e.ChunkID, &e.DocumentID, &e.SourceID, &e.Title, &e.URL,
			&e.StrategyID, &e.ChunkIndex, &e.Text,
			&e.Provider, &e.Model, &e.Dimensions, &e.TextHash, &e.Embedding, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan chunk embedding with context: %w", err)
		}
		embeddings = append(embeddings, e)
	}
	return embeddings, rows.Err()
}

// UpsertSearchIndex records metadata for a derived search index.
func (q *Queries) UpsertSearchIndex(index SearchIndex) error {
	_, err := q.db.Exec(`
		INSERT INTO search_indexes (id, name, strategy_id, provider, model, dimensions,
		    index_type, index_path, document_count, chunk_count, last_rebuild_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), ?)
		ON CONFLICT(id) DO UPDATE SET
		    name = excluded.name,
		    strategy_id = excluded.strategy_id,
		    provider = excluded.provider,
		    model = excluded.model,
		    dimensions = excluded.dimensions,
		    index_type = excluded.index_type,
		    index_path = excluded.index_path,
		    document_count = excluded.document_count,
		    chunk_count = excluded.chunk_count,
		    last_rebuild_at = datetime('now'),
		    status = excluded.status,
		    updated_at = datetime('now')
	`, index.ID, index.Name, index.StrategyID, nullString(index.Provider), nullString(index.Model), nullInt(index.Dimensions),
		index.IndexType, index.IndexPath, index.DocumentCount, index.ChunkCount, index.Status)
	if err != nil {
		return fmt.Errorf("upsert search index: %w", err)
	}
	return nil
}

// GetSearchIndex returns search index metadata by ID.
func (q *Queries) GetSearchIndex(id string) (*SearchIndex, bool, error) {
	var idx SearchIndex
	var provider, model, lastRebuild sql.NullString
	var dimensions sql.NullInt64
	err := q.db.QueryRow(`
		SELECT id, name, COALESCE(strategy_id, ''), provider, model, dimensions,
		       index_type, index_path, document_count, chunk_count,
		       last_rebuild_at, status, created_at, updated_at
		FROM search_indexes WHERE id = ?
	`, id).Scan(&idx.ID, &idx.Name, &idx.StrategyID, &provider, &model, &dimensions,
		&idx.IndexType, &idx.IndexPath, &idx.DocumentCount, &idx.ChunkCount,
		&lastRebuild, &idx.Status, &idx.CreatedAt, &idx.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("get search index: %w", err)
	}
	idx.Provider = provider.String
	idx.Model = model.String
	if dimensions.Valid {
		idx.Dimensions = int(dimensions.Int64)
	}
	idx.LastRebuildAt = lastRebuild.String
	return &idx, true, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nullInt(i int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(i), Valid: i > 0}
}
