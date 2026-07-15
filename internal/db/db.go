package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// OpenDB opens a SQLite database in read-write mode with WAL enabled
func OpenDB(dbPath string) (*sql.DB, error) {
	// Ensure parent directory exists before SQLite creates the DB and WAL files.
	// The CLI intentionally accepts a caller-selected SQLite database path.
	dir := filepath.Dir(dbPath)
	// codeql[go/path-injection] dbPath is an explicit CLI/config value for the local SQLite database location.
	// lgtm[go/path-injection]
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", dir, err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection pool settings for SQLite
	db.SetMaxOpenConns(1) // SQLite only allows one writer
	db.SetMaxIdleConns(1)

	return db, nil
}

// Migrate runs all schema migrations
func Migrate(db *sql.DB) error {
	migrations := []string{
		migrationV1Sources,
		migrationV1Documents,
		migrationV1Chunks,
		migrationV1ChunkingStrategies,
		migrationV1ChunkEmbeddings,
		migrationV1ChunkEnrichments,
		migrationV1DocumentProcessingArtifacts,
		migrationV1SearchIndexes,
		migrationV1EvalQueries,
		migrationV1EvalRuns,
		migrationV1EvalResults,
		migrationV2SourceArtifacts,
		migrationV2DocumentRevisions,
		migrationV2CorpusSnapshots,
		migrationV2ChunkPlans,
		migrationV2ChunkSets,
		migrationV2EmbeddingPlans,
		migrationV2EmbeddingSets,
		migrationV2RetrievalArtifacts,
		migrationV3ExperimentRuns,
	}

	for i, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	if err := ensureChunksStrategyID(db); err != nil {
		return fmt.Errorf("ensure chunks strategy_id: %w", err)
	}

	return nil
}

const migrationV1Sources = `
CREATE TABLE IF NOT EXISTS sources (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config_json TEXT NOT NULL DEFAULT '{}',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1Documents = `
CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    source_id TEXT NOT NULL REFERENCES sources(id),
    external_id TEXT,
    title TEXT NOT NULL,
    author TEXT DEFAULT '',
    url TEXT DEFAULT '',
    content_type TEXT DEFAULT 'text',
    raw_content TEXT,
    content_text TEXT,
    content_html TEXT,
    word_count INTEGER DEFAULT 0,
    language TEXT DEFAULT 'en',
    metadata_json TEXT DEFAULT '{}',
    status TEXT DEFAULT 'pending',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1Chunks = `
CREATE TABLE IF NOT EXISTS chunks (
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
);
`

const migrationV1ChunkingStrategies = `
CREATE TABLE IF NOT EXISTS chunking_strategies (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config_json TEXT NOT NULL DEFAULT '{}',
    description TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1ChunkEmbeddings = `
CREATE TABLE IF NOT EXISTS chunk_embeddings (
    chunk_id TEXT NOT NULL REFERENCES chunks(id) ON DELETE CASCADE,
    strategy_id TEXT NOT NULL REFERENCES chunking_strategies(id),
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    dimensions INTEGER NOT NULL,
    text_hash TEXT NOT NULL,
    embedding BLOB NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (chunk_id, strategy_id, provider, model, dimensions)
);
`

const migrationV1ChunkEnrichments = `
CREATE TABLE IF NOT EXISTS chunk_enrichments (
    chunk_id TEXT NOT NULL REFERENCES chunks(id) ON DELETE CASCADE,
    strategy_id TEXT NOT NULL REFERENCES chunking_strategies(id),
    prompt_version TEXT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    short_summary TEXT,
    long_summary TEXT,
    key_topics_json TEXT DEFAULT '[]',
    entities_json TEXT DEFAULT '[]',
    hypothetical_questions_json TEXT DEFAULT '[]',
    quality_score REAL,
    text_hash TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (chunk_id, strategy_id, prompt_version)
);
`

const migrationV1DocumentProcessingArtifacts = `
CREATE TABLE IF NOT EXISTS document_processing_artifacts (
    document_id TEXT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    artifact_type TEXT NOT NULL,
    prompt_version TEXT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    input_hash TEXT NOT NULL,
    output_text TEXT,
    output_json TEXT DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'succeeded',
    error_code TEXT DEFAULT '',
    error_message TEXT DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (document_id, artifact_type, prompt_version, provider, model)
);

CREATE INDEX IF NOT EXISTS idx_document_processing_artifacts_type
    ON document_processing_artifacts(artifact_type, prompt_version, provider, model, status);
`

const migrationV1SearchIndexes = `
CREATE TABLE IF NOT EXISTS search_indexes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    strategy_id TEXT REFERENCES chunking_strategies(id),
    provider TEXT,
    model TEXT,
    dimensions INTEGER,
    index_type TEXT NOT NULL,
    index_path TEXT NOT NULL,
    document_count INTEGER DEFAULT 0,
    chunk_count INTEGER DEFAULT 0,
    last_rebuild_at TEXT,
    status TEXT DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1EvalQueries = `
CREATE TABLE IF NOT EXISTS eval_queries (
    id TEXT PRIMARY KEY,
    text TEXT NOT NULL,
    relevant_chunk_ids_json TEXT DEFAULT '[]',
    relevant_document_ids_json TEXT DEFAULT '[]',
    notes TEXT DEFAULT '',
    category TEXT DEFAULT 'general',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1EvalRuns = `
CREATE TABLE IF NOT EXISTS eval_runs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    config_json TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    started_at TEXT,
    finished_at TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV1EvalResults = `
CREATE TABLE IF NOT EXISTS eval_results (
    id TEXT PRIMARY KEY,
    run_id TEXT NOT NULL REFERENCES eval_runs(id) ON DELETE CASCADE,
    query_id TEXT NOT NULL REFERENCES eval_queries(id),
    retrieved_chunk_ids_json TEXT NOT NULL,
    scores_json TEXT NOT NULL,
    recall_at_k REAL DEFAULT 0,
    mrr REAL DEFAULT 0,
    ndcg_at_k REAL DEFAULT 0,
    latency_ms INTEGER DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// Immutable corpus objects deliberately live alongside the earlier mutable
// operational tables. New experiment code must use these tables directly;
// they are not a cache or an adapter over documents/sources.
const migrationV2SourceArtifacts = `
CREATE TABLE IF NOT EXISTS source_artifacts (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    kind TEXT NOT NULL,
    checksum_sha256 TEXT NOT NULL,
    byte_size INTEGER NOT NULL,
    manifest_json TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(kind, checksum_sha256)
);
`

const migrationV2DocumentRevisions = `
CREATE TABLE IF NOT EXISTS document_revisions (
    id TEXT PRIMARY KEY,
    stable_document_id TEXT NOT NULL,
    source_artifact_id TEXT NOT NULL REFERENCES source_artifacts(id),
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    content_text TEXT NOT NULL,
    content_markdown TEXT NOT NULL,
    search_text TEXT NOT NULL,
    search_markdown TEXT NOT NULL,
    metadata_json TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_document_revisions_stable_document
    ON document_revisions(stable_document_id);
`

const migrationV2CorpusSnapshots = `
CREATE TABLE IF NOT EXISTS corpus_snapshots (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    source_artifact_id TEXT NOT NULL REFERENCES source_artifacts(id),
    selection_json TEXT NOT NULL,
    manifest_json TEXT NOT NULL,
    document_count INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS corpus_snapshot_documents (
    snapshot_id TEXT NOT NULL REFERENCES corpus_snapshots(id),
    ordinal INTEGER NOT NULL,
    document_revision_id TEXT NOT NULL REFERENCES document_revisions(id),
    PRIMARY KEY (snapshot_id, ordinal),
    UNIQUE (snapshot_id, document_revision_id)
);
`

const migrationV2ChunkPlans = `
CREATE TABLE IF NOT EXISTS chunk_plans (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    strategy TEXT NOT NULL,
    input_variant TEXT NOT NULL,
    config_json TEXT NOT NULL,
    implementation_version TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV2ChunkSets = `
CREATE TABLE IF NOT EXISTS chunk_sets (
    id TEXT PRIMARY KEY,
    corpus_snapshot_id TEXT NOT NULL REFERENCES corpus_snapshots(id),
    chunk_plan_id TEXT NOT NULL REFERENCES chunk_plans(id),
    manifest_json TEXT NOT NULL,
    chunk_count INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(corpus_snapshot_id, chunk_plan_id)
);

CREATE TABLE IF NOT EXISTS immutable_chunks (
    id TEXT PRIMARY KEY,
    chunk_set_id TEXT NOT NULL REFERENCES chunk_sets(id),
    document_revision_id TEXT NOT NULL REFERENCES document_revisions(id),
    chunk_index INTEGER NOT NULL,
    text TEXT NOT NULL,
    token_count INTEGER NOT NULL,
    source_start_runes INTEGER NOT NULL,
    source_end_runes INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(chunk_set_id, document_revision_id, chunk_index)
);

CREATE INDEX IF NOT EXISTS idx_immutable_chunks_set_document
    ON immutable_chunks(chunk_set_id, document_revision_id, chunk_index);
`

const migrationV2EmbeddingPlans = `
CREATE TABLE IF NOT EXISTS embedding_plans (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    provider_type TEXT NOT NULL,
    model TEXT NOT NULL,
    dimensions INTEGER NOT NULL,
    normalization TEXT NOT NULL,
    implementation_version TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const migrationV2EmbeddingSets = `
CREATE TABLE IF NOT EXISTS embedding_sets (
    id TEXT PRIMARY KEY,
    chunk_set_id TEXT NOT NULL REFERENCES chunk_sets(id),
    embedding_plan_id TEXT NOT NULL REFERENCES embedding_plans(id),
    manifest_json TEXT NOT NULL,
    embedding_count INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(chunk_set_id, embedding_plan_id)
);

CREATE TABLE IF NOT EXISTS immutable_embeddings (
    embedding_set_id TEXT NOT NULL REFERENCES embedding_sets(id),
    chunk_id TEXT NOT NULL REFERENCES immutable_chunks(id),
    vector BLOB NOT NULL,
    PRIMARY KEY (embedding_set_id, chunk_id)
);
`

const migrationV2RetrievalArtifacts = `
CREATE TABLE IF NOT EXISTS retrieval_artifacts (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    kind TEXT NOT NULL,
    chunk_set_id TEXT NOT NULL REFERENCES chunk_sets(id),
    embedding_set_id TEXT REFERENCES embedding_sets(id),
    config_json TEXT NOT NULL,
    manifest_json TEXT NOT NULL,
    artifact_path TEXT NOT NULL,
    chunk_count INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// Experiment specifications and run evidence are append-only. A specification
// identifies immutable inputs and retrieval settings. A run is an observation
// of that specification; its lifecycle is represented by ordered events and
// one optional terminal summary rather than mutable status columns.
const migrationV3ExperimentRuns = `
CREATE TABLE IF NOT EXISTS experiment_specs (
    id TEXT PRIMARY KEY,
    schema_version TEXT NOT NULL,
    corpus_snapshot_id TEXT NOT NULL REFERENCES corpus_snapshots(id),
    chunk_set_id TEXT NOT NULL REFERENCES chunk_sets(id),
    bm25_artifact_id TEXT REFERENCES retrieval_artifacts(id),
    embedding_set_id TEXT REFERENCES embedding_sets(id),
    evaluation_dataset_id TEXT NOT NULL DEFAULT '',
    config_json TEXT NOT NULL,
    manifest_json TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS experiment_runs (
    id TEXT PRIMARY KEY,
    experiment_spec_id TEXT NOT NULL REFERENCES experiment_specs(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS experiment_run_events (
    run_id TEXT NOT NULL REFERENCES experiment_runs(id),
    sequence INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    occurred_at TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    PRIMARY KEY (run_id, sequence)
);

CREATE TABLE IF NOT EXISTS experiment_run_summaries (
    run_id TEXT PRIMARY KEY REFERENCES experiment_runs(id),
    status TEXT NOT NULL,
    finished_at TEXT NOT NULL,
    metrics_json TEXT NOT NULL,
    cost_json TEXT NOT NULL,
    storage_json TEXT NOT NULL,
    error_json TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS experiment_query_traces (
    run_id TEXT NOT NULL REFERENCES experiment_runs(id),
    query_card_id TEXT NOT NULL,
    trace_json TEXT NOT NULL,
    metrics_json TEXT NOT NULL,
    timing_json TEXT NOT NULL,
    cost_json TEXT NOT NULL,
    storage_json TEXT NOT NULL,
    PRIMARY KEY (run_id, query_card_id)
);

CREATE INDEX IF NOT EXISTS idx_experiment_runs_spec_created
    ON experiment_runs(experiment_spec_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_experiment_events_run_sequence
    ON experiment_run_events(run_id, sequence);

CREATE TRIGGER IF NOT EXISTS experiment_specs_no_update
BEFORE UPDATE ON experiment_specs BEGIN SELECT RAISE(ABORT, 'experiment specs are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_specs_no_delete
BEFORE DELETE ON experiment_specs BEGIN SELECT RAISE(ABORT, 'experiment specs are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_runs_no_update
BEFORE UPDATE ON experiment_runs BEGIN SELECT RAISE(ABORT, 'experiment runs are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_runs_no_delete
BEFORE DELETE ON experiment_runs BEGIN SELECT RAISE(ABORT, 'experiment runs are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_events_no_update
BEFORE UPDATE ON experiment_run_events BEGIN SELECT RAISE(ABORT, 'experiment events are append-only'); END;
CREATE TRIGGER IF NOT EXISTS experiment_events_no_delete
BEFORE DELETE ON experiment_run_events BEGIN SELECT RAISE(ABORT, 'experiment events are append-only'); END;
CREATE TRIGGER IF NOT EXISTS experiment_summaries_no_update
BEFORE UPDATE ON experiment_run_summaries BEGIN SELECT RAISE(ABORT, 'experiment summaries are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_summaries_no_delete
BEFORE DELETE ON experiment_run_summaries BEGIN SELECT RAISE(ABORT, 'experiment summaries are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_query_traces_no_update
BEFORE UPDATE ON experiment_query_traces BEGIN SELECT RAISE(ABORT, 'query traces are immutable'); END;
CREATE TRIGGER IF NOT EXISTS experiment_query_traces_no_delete
BEFORE DELETE ON experiment_query_traces BEGIN SELECT RAISE(ABORT, 'query traces are immutable'); END;
`
