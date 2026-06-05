package corpus

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

// EmbeddingIdentity selects which provider/model/dimensions tuple to use
// for computing embedding coverage.
type EmbeddingIdentity struct {
	StrategyID   string `json:"strategy_id"`
	ProviderType string `json:"provider_type"`
	Model        string `json:"model"`
	Dimensions   int    `json:"dimensions"`
}

// SourceSummary is a source-level overview with document, word, chunk, and
// embedding counts.
type SourceSummary struct {
	SourceID              string `json:"source_id"`
	SourceName            string `json:"source_name"`
	SourceType            string `json:"source_type"`
	DocumentCount         int    `json:"document_count"`
	WordCount             int    `json:"word_count"`
	ChunkCount            int    `json:"chunk_count"`
	EmbeddedCount         int    `json:"embedded_count"`
	MissingEmbeddingCount int    `json:"missing_embedding_count"`
}

// CorpusDocumentRow is a document row returned by the corpus document browser.
type CorpusDocumentRow struct {
	ID                    string `json:"id"`
	SourceID              string `json:"source_id"`
	Title                 string `json:"title"`
	URL                   string `json:"url"`
	WordCount             int    `json:"word_count"`
	Status                string `json:"status"`
	ChunkCount            int    `json:"chunk_count"`
	EmbeddedCount         int    `json:"embedded_count"`
	MissingEmbeddingCount int    `json:"missing_embedding_count"`
}

// CorpusChunk is a chunk row with embedding status for a given identity.
type CorpusChunk struct {
	ID          string `json:"id"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
	TokenCount  int    `json:"token_count"`
	Text        string `json:"text"`
	Embedding   *struct {
		Present    bool   `json:"present"`
		Provider   string `json:"provider_type"`
		Model      string `json:"model"`
		Dimensions int    `json:"dimensions"`
		TextHash   string `json:"text_hash,omitempty"`
		UpdatedAt  string `json:"updated_at,omitempty"`
	} `json:"embedding,omitempty"`
	Enrichment *struct {
		Present       bool    `json:"present"`
		PromptVersion string  `json:"prompt_version,omitempty"`
		ShortSummary  string  `json:"short_summary,omitempty"`
		QualityScore  float64 `json:"quality_score,omitempty"`
		UpdatedAt     string  `json:"updated_at,omitempty"`
	} `json:"enrichment,omitempty"`
}

// CorpusDocumentDetail is the full document detail payload with metadata,
// content text, and chunk list.
type CorpusDocumentDetail struct {
	Document struct {
		ID          string                 `json:"id"`
		SourceID    string                 `json:"source_id"`
		ExternalID  string                 `json:"external_id"`
		Title       string                 `json:"title"`
		URL         string                 `json:"url"`
		WordCount   int                    `json:"word_count"`
		Status      string                 `json:"status"`
		Metadata    map[string]interface{} `json:"metadata"`
		ContentText string                 `json:"content_text,omitempty"`
		ContentHTML string                 `json:"content_html,omitempty"`
		ContentType string                 `json:"content_type"`
		Author      string                 `json:"author"`
		Language    string                 `json:"language"`
		CreatedAt   string                 `json:"created_at"`
		UpdatedAt   string                 `json:"updated_at"`
	} `json:"document"`
	Chunks []CorpusChunk `json:"chunks"`
}

// Service provides read-only corpus exploration queries.
type Service struct {
	db *sql.DB
}

// NewService creates a corpus service.
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// SourceSummaries returns per-source counts with optional embedding coverage.
func (s *Service) SourceSummaries(ctx context.Context, identity EmbeddingIdentity) ([]SourceSummary, error) {
	query := `
		SELECT
			s.id,
			s.name,
			s.type,
			COUNT(DISTINCT d.id) AS document_count,
			COALESCE(SUM(d.word_count), 0) AS word_count`

	args := []interface{}{}

	if identity.StrategyID != "" {
		query += `,
			COALESCE(SUM(ch.cnt), 0) AS chunk_count`
	}

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		query += `,
			COALESCE(SUM(ce.cnt), 0) AS embedded_count`
	}

	query += `
		FROM sources s
		LEFT JOIN documents d ON d.source_id = s.id`

	if identity.StrategyID != "" {
		query += `
		LEFT JOIN (
			SELECT c.document_id, COUNT(*) AS cnt
			FROM chunks c
			WHERE c.strategy_id = ?
			GROUP BY c.document_id
		) ch ON ch.document_id = d.id`
		args = append(args, identity.StrategyID)
	}

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		query += `
		LEFT JOIN (
			SELECT ce_chunk.document_id, COUNT(*) AS cnt
			FROM chunk_embeddings ce2
			JOIN chunks ce_chunk ON ce_chunk.id = ce2.chunk_id AND ce_chunk.strategy_id = ce2.strategy_id
			WHERE ce2.strategy_id = ? AND ce2.provider = ? AND ce2.model = ? AND ce2.dimensions = ?
			GROUP BY ce_chunk.document_id
		) ce ON ce.document_id = d.id`
		args = append(args, identity.StrategyID, identity.ProviderType, identity.Model, identity.Dimensions)
	}

	query += `
		GROUP BY s.id, s.name, s.type
		ORDER BY s.id`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("corpus source summaries: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []SourceSummary
	for rows.Next() {
		var ss SourceSummary
		var chunkCount, embeddedCount sql.NullInt64

		if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 && identity.StrategyID != "" {
			if err := rows.Scan(&ss.SourceID, &ss.SourceName, &ss.SourceType,
				&ss.DocumentCount, &ss.WordCount, &chunkCount, &embeddedCount); err != nil {
				return nil, fmt.Errorf("corpus source summaries scan: %w", err)
			}
			ss.ChunkCount = int(chunkCount.Int64)
			ss.EmbeddedCount = int(embeddedCount.Int64)
		} else if identity.StrategyID != "" {
			if err := rows.Scan(&ss.SourceID, &ss.SourceName, &ss.SourceType,
				&ss.DocumentCount, &ss.WordCount, &chunkCount); err != nil {
				return nil, fmt.Errorf("corpus source summaries scan: %w", err)
			}
			ss.ChunkCount = int(chunkCount.Int64)
		} else {
			if err := rows.Scan(&ss.SourceID, &ss.SourceName, &ss.SourceType,
				&ss.DocumentCount, &ss.WordCount); err != nil {
				return nil, fmt.Errorf("corpus source summaries scan: %w", err)
			}
		}
		ss.MissingEmbeddingCount = ss.ChunkCount - ss.EmbeddedCount
		results = append(results, ss)
	}

	return results, rows.Err()
}

// DocumentBrowser returns paginated documents for a source, with chunk and
// embedding counts.
func (s *Service) DocumentBrowser(ctx context.Context, sourceID string, identity EmbeddingIdentity, limit, offset int) ([]CorpusDocumentRow, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT
			d.id,
			d.source_id,
			d.title,
			COALESCE(d.url, '') AS url,
			d.word_count,
			d.status`

	args := []interface{}{}

	if identity.StrategyID != "" {
		query += `,
			COALESCE(ch.cnt, 0) AS chunk_count`
	}

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		query += `,
			COALESCE(ce.cnt, 0) AS embedded_count`
	}

	query += `
		FROM documents d`

	if identity.StrategyID != "" {
		query += `
		LEFT JOIN (
			SELECT c.document_id, COUNT(*) AS cnt
			FROM chunks c
			WHERE c.strategy_id = ?
			GROUP BY c.document_id
		) ch ON ch.document_id = d.id`
		args = append(args, identity.StrategyID)
	}

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		query += `
		LEFT JOIN (
			SELECT ce_chunk.document_id, COUNT(*) AS cnt
			FROM chunk_embeddings ce2
			JOIN chunks ce_chunk ON ce_chunk.id = ce2.chunk_id AND ce_chunk.strategy_id = ce2.strategy_id
			WHERE ce2.strategy_id = ? AND ce2.provider = ? AND ce2.model = ? AND ce2.dimensions = ?
			GROUP BY ce_chunk.document_id
		) ce ON ce.document_id = d.id`
		args = append(args, identity.StrategyID, identity.ProviderType, identity.Model, identity.Dimensions)
	}

	query += `
		WHERE d.source_id = ?
		GROUP BY d.id
		ORDER BY d.word_count DESC, d.id
		LIMIT ? OFFSET ?`

	args = append(args, sourceID, limit, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("corpus document browser: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []CorpusDocumentRow
	for rows.Next() {
		var r CorpusDocumentRow
		var chunkCount, embeddedCount sql.NullInt64

		if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 && identity.StrategyID != "" {
			if err := rows.Scan(&r.ID, &r.SourceID, &r.Title, &r.URL, &r.WordCount, &r.Status,
				&chunkCount, &embeddedCount); err != nil {
				return nil, fmt.Errorf("corpus document browser scan: %w", err)
			}
			r.ChunkCount = int(chunkCount.Int64)
			r.EmbeddedCount = int(embeddedCount.Int64)
		} else if identity.StrategyID != "" {
			if err := rows.Scan(&r.ID, &r.SourceID, &r.Title, &r.URL, &r.WordCount, &r.Status,
				&chunkCount); err != nil {
				return nil, fmt.Errorf("corpus document browser scan: %w", err)
			}
			r.ChunkCount = int(chunkCount.Int64)
		} else {
			if err := rows.Scan(&r.ID, &r.SourceID, &r.Title, &r.URL, &r.WordCount, &r.Status); err != nil {
				return nil, fmt.Errorf("corpus document browser scan: %w", err)
			}
		}
		r.MissingEmbeddingCount = r.ChunkCount - r.EmbeddedCount
		results = append(results, r)
	}

	return results, rows.Err()
}

// DocumentDetail returns the full document with chunks and embedding status.
func (s *Service) DocumentDetail(ctx context.Context, documentID string, identity EmbeddingIdentity, includeText bool) (*CorpusDocumentDetail, error) {
	// Fetch document
	var d CorpusDocumentDetail
	var metadataJSON, contentText, contentHTML sql.NullString

	err := s.db.QueryRowContext(ctx, `
		SELECT id, source_id, COALESCE(external_id, ''), title,
		       COALESCE(url, ''), word_count, status,
		       COALESCE(metadata_json, '{}'),
		       COALESCE(content_text, ''),
		       COALESCE(content_html, ''),
		       COALESCE(content_type, 'text'),
		       COALESCE(author, ''),
		       COALESCE(language, 'en'),
		       created_at, updated_at
		FROM documents WHERE id = ?
	`, documentID).Scan(
		&d.Document.ID, &d.Document.SourceID, &d.Document.ExternalID, &d.Document.Title,
		&d.Document.URL, &d.Document.WordCount, &d.Document.Status,
		&metadataJSON, &contentText, &contentHTML,
		&d.Document.ContentType, &d.Document.Author, &d.Document.Language,
		&d.Document.CreatedAt, &d.Document.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("corpus document detail: %w", err)
	}

	// Parse metadata JSON
	d.Document.Metadata = make(map[string]interface{})
	d.Chunks = []CorpusChunk{} // initialize as empty so JSON is [] not null
	if metadataJSON.Valid && metadataJSON.String != "" {
		_ = json.Unmarshal([]byte(metadataJSON.String), &d.Document.Metadata)
	}

	if includeText {
		if contentText.Valid {
			d.Document.ContentText = contentText.String
		}
		if contentHTML.Valid {
			d.Document.ContentHTML = contentHTML.String
		}
	}

	// Fetch chunks with embedding status
	if identity.StrategyID == "" {
		return &d, nil
	}

	chunkQuery := `
		SELECT
			c.id,
			c.strategy_id,
			c.chunk_index,
			COALESCE(c.start_offset, 0),
			COALESCE(c.end_offset, 0),
			c.token_count`

	chunkArgs := []interface{}{}

	if includeText {
		chunkQuery += `,
			c.text`
	} else {
		chunkQuery += `,
			''`
	}

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		chunkQuery += `,
			ce.text_hash,
			ce.updated_at`
	} else {
		chunkQuery += `,
			CAST(NULL AS TEXT) AS text_hash,
			CAST(NULL AS TEXT) AS updated_at`
	}

	// Enrichment status join
	chunkQuery += `,
		cce.short_summary,
		cce.prompt_version,
		cce.quality_score,
		cce.updated_at AS enrichment_updated_at`

	chunkQuery += `
		FROM chunks c`

	if identity.ProviderType != "" && identity.Model != "" && identity.Dimensions > 0 {
		chunkQuery += `
		LEFT JOIN chunk_embeddings ce ON ce.chunk_id = c.id
			AND ce.strategy_id = c.strategy_id
			AND ce.provider = ?
			AND ce.model = ?
			AND ce.dimensions = ?`
		chunkArgs = append(chunkArgs, identity.ProviderType, identity.Model, identity.Dimensions)
	}

	// Left join chunk enrichments for enrichment status
	chunkQuery += `
		LEFT JOIN chunk_enrichments cce ON cce.chunk_id = c.id AND cce.strategy_id = c.strategy_id`

	chunkQuery += `
		WHERE c.document_id = ? AND c.strategy_id = ?
		ORDER BY c.chunk_index`

	chunkArgs = append(chunkArgs, documentID, identity.StrategyID)

	rows, err := s.db.QueryContext(ctx, chunkQuery, chunkArgs...)
	if err != nil {
		return nil, fmt.Errorf("corpus document detail chunks: %w", err)
	}
	defer func() { _ = rows.Close() }()
	if err != nil {
		return nil, fmt.Errorf("corpus document detail chunks: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var c CorpusChunk
		var textHash, updatedAt sql.NullString
		var enrichSummary, enrichVersion, enrichUpdatedAt sql.NullString
		var enrichScore sql.NullFloat64

		if err := rows.Scan(&c.ID, &c.StrategyID, &c.ChunkIndex, &c.StartOffset,
			&c.EndOffset, &c.TokenCount, &c.Text, &textHash, &updatedAt,
			&enrichSummary, &enrichVersion, &enrichScore, &enrichUpdatedAt); err != nil {
			return nil, fmt.Errorf("corpus document detail chunk scan: %w", err)
		}

		if textHash.Valid && textHash.String != "" {
			c.Embedding = &struct {
				Present    bool   `json:"present"`
				Provider   string `json:"provider_type"`
				Model      string `json:"model"`
				Dimensions int    `json:"dimensions"`
				TextHash   string `json:"text_hash,omitempty"`
				UpdatedAt  string `json:"updated_at,omitempty"`
			}{
				Present:    true,
				Provider:   identity.ProviderType,
				Model:      identity.Model,
				Dimensions: identity.Dimensions,
				TextHash:   textHash.String,
				UpdatedAt:  updatedAt.String,
			}
		}

		if enrichSummary.Valid && enrichSummary.String != "" {
			c.Enrichment = &struct {
				Present       bool    `json:"present"`
				PromptVersion string  `json:"prompt_version,omitempty"`
				ShortSummary  string  `json:"short_summary,omitempty"`
				QualityScore  float64 `json:"quality_score,omitempty"`
				UpdatedAt     string  `json:"updated_at,omitempty"`
			}{
				Present:       true,
				PromptVersion: enrichVersion.String,
				ShortSummary:  enrichSummary.String,
				QualityScore:  enrichScore.Float64,
				UpdatedAt:     enrichUpdatedAt.String,
			}
		}

		d.Chunks = append(d.Chunks, c)
	}

	return &d, rows.Err()
}
