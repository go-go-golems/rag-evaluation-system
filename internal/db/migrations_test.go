package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrateUpgradesLegacyChunksWithoutStrategyID(t *testing.T) {
	database, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "legacy.db")+"?_foreign_keys=on")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer func() { _ = database.Close() }()

	legacySchema := []string{
		migrationV1Sources,
		migrationV1Documents,
		`CREATE TABLE chunks (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
			chunk_index INTEGER NOT NULL,
			text TEXT NOT NULL,
			token_count INTEGER NOT NULL DEFAULT 0,
			start_offset INTEGER DEFAULT 0,
			end_offset INTEGER DEFAULT 0,
			boundaries_json TEXT DEFAULT '{}',
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			UNIQUE(document_id, chunk_index)
		)`,
		`INSERT INTO sources (id, name, type, config_json) VALUES ('s1', 'Source', 'filesystem', '{}')`,
		`INSERT INTO documents (id, source_id, title, content_type, content_text, word_count, language, status) VALUES ('d1', 's1', 'Doc', 'text', 'hello world', 2, 'en', 'chunked')`,
		`INSERT INTO chunks (id, document_id, chunk_index, text, token_count, start_offset, end_offset, boundaries_json) VALUES ('c1', 'd1', 0, 'hello', 1, 0, 5, '{"strategy_id":"fixed-old"}')`,
	}
	for _, stmt := range legacySchema {
		if _, err := database.Exec(stmt); err != nil {
			t.Fatalf("legacy setup failed for %q: %v", stmt, err)
		}
	}

	if err := Migrate(database); err != nil {
		t.Fatalf("migrate legacy db: %v", err)
	}

	hasStrategyID, err := tableHasColumn(database, "chunks", "strategy_id")
	if err != nil {
		t.Fatalf("check chunks strategy_id column: %v", err)
	}
	if !hasStrategyID {
		t.Fatal("expected chunks.strategy_id after migration")
	}

	var strategyID string
	if err := database.QueryRow(`SELECT strategy_id FROM chunks WHERE id = 'c1'`).Scan(&strategyID); err != nil {
		t.Fatalf("read migrated chunk: %v", err)
	}
	if strategyID != "fixed-old" {
		t.Fatalf("expected migrated strategy fixed-old, got %q", strategyID)
	}
}
