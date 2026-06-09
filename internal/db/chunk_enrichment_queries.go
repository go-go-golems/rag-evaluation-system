package db

import (
	"database/sql"
	"fmt"
)

type ChunkEnrichment struct {
	ChunkID                   string  `json:"chunk_id"`
	DocumentID                string  `json:"document_id,omitempty"`
	StrategyID                string  `json:"strategy_id"`
	PromptVersion             string  `json:"prompt_version"`
	Provider                  string  `json:"provider"`
	Model                     string  `json:"model"`
	ShortSummary              string  `json:"short_summary,omitempty"`
	LongSummary               string  `json:"long_summary,omitempty"`
	KeyTopicsJSON             string  `json:"key_topics_json,omitempty"`
	EntitiesJSON              string  `json:"entities_json,omitempty"`
	HypotheticalQuestionsJSON string  `json:"hypothetical_questions_json,omitempty"`
	QualityScore              float64 `json:"quality_score,omitempty"`
	TextHash                  string  `json:"text_hash"`
	CreatedAt                 string  `json:"created_at"`
	UpdatedAt                 string  `json:"updated_at"`
}

type ChunkEnrichmentCoverage struct {
	SourceID      string `json:"source_id"`
	ChunkCount    int    `json:"chunk_count"`
	EnrichedCount int    `json:"enriched_count"`
	FreshCount    int    `json:"fresh_count"`
	MissingCount  int    `json:"missing_count"`
}

func (q *Queries) GetChunk(chunkID, strategyID string) (*Chunk, bool, error) {
	var c Chunk
	err := q.db.QueryRow(`
		SELECT id, document_id, strategy_id, chunk_index, text, token_count,
		       COALESCE(start_offset, 0), COALESCE(end_offset, 0), created_at
		FROM chunks
		WHERE id = ? AND strategy_id = ?
	`, chunkID, strategyID).Scan(&c.ID, &c.DocumentID, &c.StrategyID, &c.ChunkIndex, &c.Text, &c.TokenCount, &c.StartOffset, &c.EndOffset, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("get chunk: %w", err)
	}
	return &c, true, nil
}

func (q *Queries) ListChunksForDocuments(strategyID string, documentIDs []string, perDocumentLimit int) ([]Chunk, error) {
	if len(documentIDs) == 0 {
		return nil, nil
	}
	// #nosec G202 -- placeholder count is derived from documentIDs length; IDs are bound in args below.
	query := `
		SELECT id, document_id, strategy_id, chunk_index, text, token_count,
		       COALESCE(start_offset, 0), COALESCE(end_offset, 0), created_at
		FROM (
			SELECT c.id, c.document_id, c.strategy_id, c.chunk_index, c.text, c.token_count,
			       c.start_offset, c.end_offset, c.created_at,
			       ROW_NUMBER() OVER (PARTITION BY c.document_id ORDER BY c.chunk_index) AS rn
			FROM chunks c
			WHERE c.strategy_id = ? AND c.document_id IN (` + placeholders(len(documentIDs)) + `)
		)
		WHERE rn <= ?
		ORDER BY document_id, chunk_index
	`
	args := []any{strategyID}
	for _, documentID := range documentIDs {
		args = append(args, documentID)
	}
	if perDocumentLimit <= 0 {
		perDocumentLimit = 1
	}
	args = append(args, perDocumentLimit)
	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list chunks for documents: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []Chunk{}
	for rows.Next() {
		var c Chunk
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.StrategyID, &c.ChunkIndex, &c.Text, &c.TokenCount, &c.StartOffset, &c.EndOffset, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan chunk for documents: %w", err)
		}
		ret = append(ret, c)
	}
	return ret, rows.Err()
}

func (q *Queries) UpsertChunkEnrichment(e ChunkEnrichment) error {
	if e.KeyTopicsJSON == "" {
		e.KeyTopicsJSON = "[]"
	}
	if e.EntitiesJSON == "" {
		e.EntitiesJSON = "[]"
	}
	if e.HypotheticalQuestionsJSON == "" {
		e.HypotheticalQuestionsJSON = "[]"
	}
	_, err := q.db.Exec(`
		INSERT INTO chunk_enrichments (
			chunk_id, strategy_id, prompt_version, provider, model,
			short_summary, long_summary, key_topics_json, entities_json,
			hypothetical_questions_json, quality_score, text_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(chunk_id, strategy_id, prompt_version) DO UPDATE SET
			provider = excluded.provider,
			model = excluded.model,
			short_summary = excluded.short_summary,
			long_summary = excluded.long_summary,
			key_topics_json = excluded.key_topics_json,
			entities_json = excluded.entities_json,
			hypothetical_questions_json = excluded.hypothetical_questions_json,
			quality_score = excluded.quality_score,
			text_hash = excluded.text_hash,
			updated_at = datetime('now')
	`, e.ChunkID, e.StrategyID, e.PromptVersion, e.Provider, e.Model,
		e.ShortSummary, e.LongSummary, e.KeyTopicsJSON, e.EntitiesJSON,
		e.HypotheticalQuestionsJSON, e.QualityScore, e.TextHash)
	if err != nil {
		return fmt.Errorf("upsert chunk enrichment: %w", err)
	}
	return nil
}

func (q *Queries) GetChunkEnrichment(chunkID, strategyID, promptVersion string) (*ChunkEnrichment, bool, error) {
	var e ChunkEnrichment
	var shortSummary, longSummary, keyTopics, entities, questions sql.NullString
	var quality sql.NullFloat64
	err := q.db.QueryRow(`
		SELECT chunk_id, strategy_id, prompt_version, provider, model,
		       short_summary, long_summary, key_topics_json, entities_json,
		       hypothetical_questions_json, quality_score, text_hash, created_at, updated_at
		FROM chunk_enrichments
		WHERE chunk_id = ? AND strategy_id = ? AND prompt_version = ?
	`, chunkID, strategyID, promptVersion).Scan(&e.ChunkID, &e.StrategyID, &e.PromptVersion, &e.Provider, &e.Model,
		&shortSummary, &longSummary, &keyTopics, &entities, &questions, &quality, &e.TextHash, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("get chunk enrichment: %w", err)
	}
	if shortSummary.Valid {
		e.ShortSummary = shortSummary.String
	}
	if longSummary.Valid {
		e.LongSummary = longSummary.String
	}
	if keyTopics.Valid {
		e.KeyTopicsJSON = keyTopics.String
	}
	if entities.Valid {
		e.EntitiesJSON = entities.String
	}
	if questions.Valid {
		e.HypotheticalQuestionsJSON = questions.String
	}
	if quality.Valid {
		e.QualityScore = quality.Float64
	}
	return &e, true, nil
}

func (q *Queries) IsChunkEnrichmentFresh(chunkID, strategyID, promptVersion, textHash string) (bool, error) {
	e, ok, err := q.GetChunkEnrichment(chunkID, strategyID, promptVersion)
	if err != nil || !ok {
		return false, err
	}
	return e.TextHash == textHash, nil
}

func (q *Queries) ListChunkEnrichments(chunkID, strategyID, promptVersion string) ([]ChunkEnrichment, error) {
	query := `
		SELECT e.chunk_id, c.document_id, e.strategy_id, e.prompt_version, e.provider, e.model,
		       COALESCE(e.short_summary, ''), COALESCE(e.long_summary, ''),
		       COALESCE(e.key_topics_json, '[]'), COALESCE(e.entities_json, '[]'),
		       COALESCE(e.hypothetical_questions_json, '[]'), COALESCE(e.quality_score, 0),
		       e.text_hash, e.created_at, e.updated_at
		FROM chunk_enrichments e
		JOIN chunks c ON c.id = e.chunk_id AND c.strategy_id = e.strategy_id
		WHERE e.chunk_id = ?
	`
	args := []any{chunkID}
	if strategyID != "" {
		query += ` AND e.strategy_id = ?`
		args = append(args, strategyID)
	}
	if promptVersion != "" {
		query += ` AND e.prompt_version = ?`
		args = append(args, promptVersion)
	}
	query += ` ORDER BY e.strategy_id, e.prompt_version, e.provider, e.model`
	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list chunk enrichments: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []ChunkEnrichment{}
	for rows.Next() {
		var e ChunkEnrichment
		if err := rows.Scan(&e.ChunkID, &e.DocumentID, &e.StrategyID, &e.PromptVersion, &e.Provider, &e.Model, &e.ShortSummary, &e.LongSummary, &e.KeyTopicsJSON, &e.EntitiesJSON, &e.HypotheticalQuestionsJSON, &e.QualityScore, &e.TextHash, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan chunk enrichment: %w", err)
		}
		ret = append(ret, e)
	}
	return ret, rows.Err()
}

type ChunkEnrichmentIdentity struct {
	StrategyID    string `json:"strategy_id"`
	PromptVersion string `json:"prompt_version"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	EnrichedCount int    `json:"enriched_count"`
}

func (q *Queries) ListChunkEnrichmentIdentities() ([]ChunkEnrichmentIdentity, error) {
	rows, err := q.db.Query(`
		SELECT strategy_id, prompt_version, provider, model, COUNT(*)
		FROM chunk_enrichments
		GROUP BY strategy_id, prompt_version, provider, model
		ORDER BY strategy_id, prompt_version, provider, model
	`)
	if err != nil {
		return nil, fmt.Errorf("list chunk enrichment identities: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []ChunkEnrichmentIdentity{}
	for rows.Next() {
		var id ChunkEnrichmentIdentity
		if err := rows.Scan(&id.StrategyID, &id.PromptVersion, &id.Provider, &id.Model, &id.EnrichedCount); err != nil {
			return nil, fmt.Errorf("scan chunk enrichment identity: %w", err)
		}
		ret = append(ret, id)
	}
	return ret, rows.Err()
}

func (q *Queries) ListChunkEnrichmentCoverage(strategyID, promptVersion string) ([]ChunkEnrichmentCoverage, error) {
	rows, err := q.db.Query(`
		SELECT d.source_id,
		       COUNT(c.id) AS chunk_count,
		       COUNT(e.chunk_id) AS enriched_count,
		       SUM(CASE WHEN e.text_hash IS NOT NULL AND e.text_hash != '' THEN 1 ELSE 0 END) AS fresh_count
		FROM chunks c
		JOIN documents d ON d.id = c.document_id
		LEFT JOIN chunk_enrichments e ON e.chunk_id = c.id AND e.strategy_id = c.strategy_id AND e.prompt_version = ?
		WHERE c.strategy_id = ?
		GROUP BY d.source_id
		ORDER BY d.source_id
	`, promptVersion, strategyID)
	if err != nil {
		return nil, fmt.Errorf("list chunk enrichment coverage: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []ChunkEnrichmentCoverage{}
	for rows.Next() {
		var c ChunkEnrichmentCoverage
		if err := rows.Scan(&c.SourceID, &c.ChunkCount, &c.EnrichedCount, &c.FreshCount); err != nil {
			return nil, fmt.Errorf("scan chunk enrichment coverage: %w", err)
		}
		c.MissingCount = c.ChunkCount - c.FreshCount
		ret = append(ret, c)
	}
	return ret, rows.Err()
}
