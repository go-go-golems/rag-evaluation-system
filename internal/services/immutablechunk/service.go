// Package immutablechunk builds chunk artifacts from immutable corpus snapshots.
package immutablechunk

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	chunkcore "github.com/go-go-golems/rag-evaluation-system/internal/chunking"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
	"github.com/pkg/errors"
)

const (
	planSchemaVersion     = "rag-eval-chunk-plan/v1"
	setSchemaVersion      = "rag-eval-chunk-set/v1"
	implementationVersion = "go-chunkers-source-runes/v1"
)

type Request struct {
	CorpusSnapshotID string
	Strategy         string
	InputVariant     string
	ChunkSize        int
	Overlap          int
}

type Result struct {
	ChunkPlanID string
	ChunkSetID  string
	ChunkCount  int
	Reused      bool
}

type planContract struct {
	SchemaVersion         string `json:"schema_version"`
	Strategy              string `json:"strategy"`
	InputVariant          string `json:"input_variant"`
	ChunkSize             int    `json:"chunk_size"`
	Overlap               int    `json:"overlap"`
	ImplementationVersion string `json:"implementation_version"`
}
type chunkSetContract struct {
	SchemaVersion    string   `json:"schema_version"`
	CorpusSnapshotID string   `json:"corpus_snapshot_id"`
	ChunkPlanID      string   `json:"chunk_plan_id"`
	ChunkIDs         []string `json:"chunk_ids"`
}
type builtChunk struct {
	ID, RevisionID, Text          string
	Index, TokenCount, Start, End int
}

// Build creates or verifies an immutable chunk plan and complete chunk set.
func Build(ctx context.Context, queries *db.Queries, request Request) (*Result, error) {
	if queries == nil {
		return nil, errors.New("database queries are required")
	}
	if request.CorpusSnapshotID == "" {
		return nil, errors.New("corpus snapshot ID is required")
	}
	if request.Strategy == "" {
		request.Strategy = "fixed"
	}
	if request.InputVariant == "" {
		request.InputVariant = "search_text"
	}
	if request.ChunkSize <= 0 {
		request.ChunkSize = 1200
	}
	if request.Strategy == "markdown-heading" {
		request.InputVariant = "search_markdown"
	}
	contract := planContract{planSchemaVersion, request.Strategy, request.InputVariant, request.ChunkSize, request.Overlap, implementationVersion}
	planID, err := experiments.Fingerprint(planSchemaVersion, contract)
	if err != nil {
		return nil, err
	}
	configJSON, err := experiments.CanonicalJSON(contract)
	if err != nil {
		return nil, err
	}
	documents, err := loadSnapshotDocuments(ctx, queries.DB(), request.CorpusSnapshotID, request.InputVariant)
	if err != nil {
		return nil, err
	}
	built := make([]builtChunk, 0)
	for _, document := range documents {
		chunker, err := newChunker(request.Strategy, request.ChunkSize, request.Overlap, planID)
		if err != nil {
			return nil, err
		}
		chunks, err := chunker.Chunk(document.RevisionID, document.Text)
		if err != nil {
			return nil, errors.Wrapf(err, "chunk revision %s", document.RevisionID)
		}
		for _, chunk := range chunks {
			if chunk.StartOffset < 0 || chunk.EndOffset < chunk.StartOffset || chunk.EndOffset > len([]rune(document.Text)) || string([]rune(document.Text)[chunk.StartOffset:chunk.EndOffset]) != chunk.Text {
				return nil, errors.Errorf("chunk source-range invariant failed for revision %s index %d", document.RevisionID, chunk.ChunkIndex)
			}
			id, err := experiments.Fingerprint("rag-eval-chunk/v1", struct {
				ChunkSetInput string `json:"chunk_set_input"`
				RevisionID    string `json:"revision_id"`
				Index         int    `json:"index"`
				Text          string `json:"text"`
				Start         int    `json:"start"`
				End           int    `json:"end"`
			}{request.CorpusSnapshotID + ":" + planID, document.RevisionID, chunk.ChunkIndex, chunk.Text, chunk.StartOffset, chunk.EndOffset})
			if err != nil {
				return nil, err
			}
			built = append(built, builtChunk{id, document.RevisionID, chunk.Text, chunk.ChunkIndex, chunk.TokenCount, chunk.StartOffset, chunk.EndOffset})
		}
	}
	ids := make([]string, 0, len(built))
	for _, chunk := range built {
		ids = append(ids, chunk.ID)
	}
	setContract := chunkSetContract{setSchemaVersion, request.CorpusSnapshotID, planID, ids}
	setID, err := experiments.Fingerprint(setSchemaVersion, setContract)
	if err != nil {
		return nil, err
	}
	manifest, err := experiments.CanonicalJSON(setContract)
	if err != nil {
		return nil, err
	}
	tx, err := queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin chunk set transaction")
	}
	defer func() { _ = tx.Rollback() }()
	if err := ensurePlan(ctx, tx, planID, contract, string(configJSON)); err != nil {
		return nil, err
	}
	reused, err := ensureSet(ctx, tx, setID, request.CorpusSnapshotID, planID, string(manifest), built)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit chunk set")
	}
	return &Result{planID, setID, len(built), reused}, nil
}

type snapshotDocument struct{ RevisionID, Text string }

func loadSnapshotDocuments(ctx context.Context, database *sql.DB, snapshotID, variant string) ([]snapshotDocument, error) {
	column := map[string]string{"search_text": "search_text", "search_markdown": "search_markdown", "content_text": "content_text", "content_markdown": "content_markdown"}[variant]
	if column == "" {
		return nil, errors.Errorf("unsupported input variant %q", variant)
	}
	rows, err := database.QueryContext(ctx, fmt.Sprintf(`SELECT csd.document_revision_id, dr.%s FROM corpus_snapshot_documents csd JOIN document_revisions dr ON dr.id = csd.document_revision_id WHERE csd.snapshot_id = ? ORDER BY csd.ordinal`, column), snapshotID)
	if err != nil {
		return nil, errors.Wrap(err, "load snapshot documents")
	}
	defer func() { _ = rows.Close() }()
	var result []snapshotDocument
	for rows.Next() {
		var doc snapshotDocument
		if err := rows.Scan(&doc.RevisionID, &doc.Text); err != nil {
			return nil, errors.Wrap(err, "scan snapshot document")
		}
		result = append(result, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate snapshot documents")
	}
	if len(result) == 0 {
		return nil, errors.Errorf("corpus snapshot %s has no documents", snapshotID)
	}
	return result, nil
}
func newChunker(strategy string, size, overlap int, id string) (chunkcore.Chunker, error) {
	switch strategy {
	case "fixed":
		return chunkcore.NewFixedSizeChunker(size, overlap, id), nil
	case "sentence":
		return chunkcore.NewSentenceChunker(size, overlap, id), nil
	case "markdown-heading":
		return chunkcore.NewMarkdownHeadingChunker(size, id), nil
	default:
		return nil, errors.Errorf("unsupported chunk strategy %q", strategy)
	}
}
func ensurePlan(ctx context.Context, tx *sql.Tx, id string, c planContract, config string) error {
	var existing string
	err := tx.QueryRowContext(ctx, `SELECT config_json FROM chunk_plans WHERE id=?`, id).Scan(&existing)
	if err == nil {
		if existing != config {
			return errors.Errorf("chunk plan %s conflicts", id)
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "look up chunk plan")
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO chunk_plans (id,schema_version,strategy,input_variant,config_json,implementation_version) VALUES (?,?,?,?,?,?)`, id, c.SchemaVersion, c.Strategy, c.InputVariant, config, c.ImplementationVersion)
	return errors.Wrap(err, "insert chunk plan")
}
func ensureSet(ctx context.Context, tx *sql.Tx, id, snapshot, plan, manifest string, chunks []builtChunk) (bool, error) {
	var existing string
	var count int
	err := tx.QueryRowContext(ctx, `SELECT manifest_json,chunk_count FROM chunk_sets WHERE id=?`, id).Scan(&existing, &count)
	if err == nil {
		if existing != manifest || count != len(chunks) {
			return false, errors.Errorf("chunk set %s conflicts", id)
		}
		return true, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, errors.Wrap(err, "look up chunk set")
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO chunk_sets (id,corpus_snapshot_id,chunk_plan_id,manifest_json,chunk_count) VALUES (?,?,?,?,?)`, id, snapshot, plan, manifest, len(chunks)); err != nil {
		return false, errors.Wrap(err, "insert chunk set")
	}
	for _, c := range chunks {
		if _, err = tx.ExecContext(ctx, `INSERT INTO immutable_chunks (id,chunk_set_id,document_revision_id,chunk_index,text,token_count,source_start_runes,source_end_runes) VALUES (?,?,?,?,?,?,?,?)`, c.ID, id, c.RevisionID, c.Index, c.Text, c.TokenCount, c.Start, c.End); err != nil {
			return false, errors.Wrap(err, "insert immutable chunk")
		}
	}
	return false, nil
}

var _ = json.Valid
