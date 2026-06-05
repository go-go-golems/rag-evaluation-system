package db

import (
	"database/sql"
	"fmt"
)

// ensureChunksStrategyID upgrades older development databases whose chunks table
// predated first-class strategy_id support. New databases already use the
// correct schema from migrationV1Chunks.
func ensureChunksStrategyID(db *sql.DB) error {
	hasStrategyID, err := tableHasColumn(db, "chunks", "strategy_id")
	if err != nil {
		return err
	}
	if hasStrategyID {
		return nil
	}

	// Existing rows were created before strategy_id existed. Preserve them under
	// a legacy strategy so old dev DBs remain inspectable, then rebuild the table
	// to get the correct UNIQUE(document_id, strategy_id, chunk_index) key.
	if _, err := db.Exec(`
		INSERT OR IGNORE INTO chunking_strategies (id, name, type, config_json, description)
		VALUES ('legacy', 'legacy', 'legacy', '{}', 'Migrated chunks from schema before strategy_id existed')
	`); err != nil {
		return fmt.Errorf("create legacy chunking strategy: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		return err
	}
	defer func() { _, _ = db.Exec(`PRAGMA foreign_keys = ON`) }()

	stmts := []string{
		`CREATE TABLE chunks_new (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
			strategy_id TEXT NOT NULL REFERENCES chunking_strategies(id),
			chunk_index INTEGER NOT NULL,
			text TEXT NOT NULL,
			token_count INTEGER NOT NULL DEFAULT 0,
			start_offset INTEGER DEFAULT 0,
			end_offset INTEGER DEFAULT 0,
			boundaries_json TEXT DEFAULT '{}',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE(document_id, strategy_id, chunk_index)
		)`,
		`INSERT INTO chunks_new (id, document_id, strategy_id, chunk_index, text, token_count, start_offset, end_offset, boundaries_json, created_at, updated_at)
		 SELECT id, document_id,
		        COALESCE(NULLIF(json_extract(boundaries_json, '$.strategy_id'), ''), 'legacy'),
		        chunk_index, text, token_count, start_offset, end_offset, boundaries_json, created_at, updated_at
		 FROM chunks`,
		`DROP TABLE chunks`,
		`ALTER TABLE chunks_new RENAME TO chunks`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func tableHasColumn(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf(`PRAGMA table_info(%s)`, tableName))
	if err != nil {
		return false, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}
	return false, rows.Err()
}
