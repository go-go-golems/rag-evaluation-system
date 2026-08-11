// Package corpussnapshot persists content-addressed document revisions and
// immutable corpus membership for the RAG laboratory.
package corpussnapshot

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
	"github.com/pkg/errors"
)

const (
	sourceArtifactSchemaVersion = "rag-eval-source-artifact/v1"
	documentRevisionSchema      = "rag-eval-document-revision/v1"
	corpusSnapshotSchema        = "rag-eval-corpus-snapshot/v1"
	sourceArtifactKind          = "ttc-wordpress-sqlite"
)

// PersistRequest supplies physical source facts that are intentionally not
// copied into the semantic snapshot identity.
type PersistRequest struct {
	SourceByteSize int64
}

// Result identifies the immutable snapshot constructed or reused.
type Result struct {
	SourceArtifactID string
	SnapshotID       string
	DocumentCount    int
	Reused           bool
}

type sourceArtifactManifest struct {
	SchemaVersion  string `json:"schema_version"`
	Kind           string `json:"kind"`
	ChecksumSHA256 string `json:"checksum_sha256"`
	ByteSize       int64  `json:"byte_size"`
}

type kindQuota struct {
	Kind  string `json:"kind"`
	Count int    `json:"count"`
}

type selectionContract struct {
	SchemaVersion      string      `json:"schema_version"`
	SnapshotName       string      `json:"snapshot_name"`
	SourceExportSHA256 string      `json:"source_export_sha256"`
	SelectionAlgorithm string      `json:"selection_algorithm"`
	KindQuotas         []kindQuota `json:"kind_quotas"`
	SeedDocumentIDs    []string    `json:"seed_document_ids"`
	SourceDocumentIDs  []string    `json:"source_document_ids"`
}

type snapshotManifest struct {
	SchemaVersion       string            `json:"schema_version"`
	SourceArtifactID    string            `json:"source_artifact_id"`
	Selection           selectionContract `json:"selection"`
	DocumentRevisionIDs []string          `json:"document_revision_ids"`
}

// Persist converts one deterministic TTC selection into immutable revisions
// and an ordered corpus snapshot. An existing ID is accepted only after every
// persisted semantic field and every membership ordinal is verified equal.
func Persist(ctx context.Context, queries *db.Queries, plan *ttcimport.Plan, request PersistRequest) (*Result, error) {
	if queries == nil {
		return nil, errors.New("database queries are required")
	}
	if plan == nil || len(plan.Documents) == 0 {
		return nil, errors.New("a non-empty TTC plan is required")
	}
	if request.SourceByteSize < 0 {
		return nil, errors.New("source byte size must not be negative")
	}
	if err := validatePlan(plan); err != nil {
		return nil, err
	}

	sourceID := fingerprint("source-artifact", sourceArtifactSchemaVersion, sourceArtifactKind, plan.Manifest.SourceExportSHA256)
	sourceManifest, err := marshal(sourceArtifactManifest{
		SchemaVersion:  sourceArtifactSchemaVersion,
		Kind:           sourceArtifactKind,
		ChecksumSHA256: plan.Manifest.SourceExportSHA256,
		ByteSize:       request.SourceByteSize,
	})
	if err != nil {
		return nil, err
	}

	revisions := make([]revision, 0, len(plan.Documents))
	for _, document := range plan.Documents {
		value, err := newRevision(sourceID, plan.Manifest.SourceExportSHA256, document)
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, value)
	}
	selection := newSelection(plan.Manifest, plan.Documents)
	revisionIDs := make([]string, 0, len(revisions))
	for _, value := range revisions {
		revisionIDs = append(revisionIDs, value.ID)
	}
	snapshotID := fingerprint("corpus-snapshot", corpusSnapshotSchema, sourceID, mustJSON(selection), strings.Join(revisionIDs, ","))
	selectionJSON, err := marshal(selection)
	if err != nil {
		return nil, err
	}
	manifestJSON, err := marshal(snapshotManifest{
		SchemaVersion: corpusSnapshotSchema, SourceArtifactID: sourceID, Selection: selection, DocumentRevisionIDs: revisionIDs,
	})
	if err != nil {
		return nil, err
	}

	tx, err := queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin immutable corpus transaction")
	}
	defer func() { _ = tx.Rollback() }()
	if err := ensureSourceArtifact(ctx, tx, sourceID, plan.Manifest.SourceExportSHA256, request.SourceByteSize, sourceManifest); err != nil {
		return nil, err
	}
	for _, value := range revisions {
		if err := ensureRevision(ctx, tx, value); err != nil {
			return nil, err
		}
	}
	reused, err := ensureSnapshot(ctx, tx, snapshotID, sourceID, selectionJSON, manifestJSON, revisionIDs)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit immutable corpus snapshot")
	}
	return &Result{SourceArtifactID: sourceID, SnapshotID: snapshotID, DocumentCount: len(revisions), Reused: reused}, nil
}

type revision struct {
	ID, StableDocumentID, SourceArtifactID, Kind, Title, URL string
	ContentText, ContentMarkdown, SearchText, SearchMarkdown string
	MetadataJSON, ContentHash                                string
}

func newRevision(sourceID, sourceChecksum string, document ttcimport.SourceDocument) (revision, error) {
	metadata, err := marshal(struct {
		SourceDocumentID   string `json:"source_document_id"`
		SourceExportSHA256 string `json:"source_export_sha256"`
	}{document.ID, sourceChecksum})
	if err != nil {
		return revision{}, err
	}
	contentHash := fingerprint("document-content", document.Kind, document.Title, document.URL, document.SearchText, document.SearchMarkdown, document.ContentMarkdown)
	return revision{
		ID:               fingerprint("document-revision", documentRevisionSchema, sourceID, document.ID, contentHash),
		StableDocumentID: document.ID, SourceArtifactID: sourceID, Kind: document.Kind, Title: document.Title, URL: document.URL,
		ContentText: document.SearchText, ContentMarkdown: document.ContentMarkdown, SearchText: document.SearchText, SearchMarkdown: document.SearchMarkdown,
		MetadataJSON: metadata, ContentHash: contentHash,
	}, nil
}

func newSelection(manifest ttcimport.Manifest, documents []ttcimport.SourceDocument) selectionContract {
	keys := make([]string, 0, len(manifest.KindQuotas))
	for key := range manifest.KindQuotas {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	quotas := make([]kindQuota, 0, len(keys))
	for _, key := range keys {
		quotas = append(quotas, kindQuota{Kind: key, Count: manifest.KindQuotas[key]})
	}
	ids := make([]string, 0, len(documents))
	for _, document := range documents {
		ids = append(ids, document.ID)
	}
	return selectionContract{SchemaVersion: manifest.SchemaVersion, SnapshotName: manifest.SnapshotName, SourceExportSHA256: manifest.SourceExportSHA256, SelectionAlgorithm: manifest.SelectionAlgorithm, KindQuotas: quotas, SeedDocumentIDs: append([]string(nil), manifest.SeedDocumentIDs...), SourceDocumentIDs: ids}
}

func ensureSourceArtifact(ctx context.Context, tx *sql.Tx, id, checksum string, byteSize int64, manifest string) error {
	var existingChecksum, existingManifest string
	var existingSize int64
	err := tx.QueryRowContext(ctx, `SELECT checksum_sha256, byte_size, manifest_json FROM source_artifacts WHERE id = ?`, id).Scan(&existingChecksum, &existingSize, &existingManifest)
	if err == nil {
		if existingChecksum != checksum || existingSize != byteSize || existingManifest != manifest {
			return errors.Errorf("source artifact %s conflicts with existing immutable content", id)
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "look up source artifact")
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO source_artifacts (id, schema_version, kind, checksum_sha256, byte_size, manifest_json) VALUES (?, ?, ?, ?, ?, ?)`, id, sourceArtifactSchemaVersion, sourceArtifactKind, checksum, byteSize, manifest)
	return errors.Wrap(err, "insert source artifact")
}

func ensureRevision(ctx context.Context, tx *sql.Tx, value revision) error {
	var existing revision
	err := tx.QueryRowContext(ctx, `SELECT id, stable_document_id, source_artifact_id, kind, title, url, content_text, content_markdown, search_text, search_markdown, metadata_json, content_hash FROM document_revisions WHERE id = ?`, value.ID).Scan(&existing.ID, &existing.StableDocumentID, &existing.SourceArtifactID, &existing.Kind, &existing.Title, &existing.URL, &existing.ContentText, &existing.ContentMarkdown, &existing.SearchText, &existing.SearchMarkdown, &existing.MetadataJSON, &existing.ContentHash)
	if err == nil {
		if existing != value {
			return errors.Errorf("document revision %s conflicts with existing immutable content", value.ID)
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrapf(err, "look up document revision %s", value.ID)
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO document_revisions (id, stable_document_id, source_artifact_id, kind, title, url, content_text, content_markdown, search_text, search_markdown, metadata_json, content_hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, value.ID, value.StableDocumentID, value.SourceArtifactID, value.Kind, value.Title, value.URL, value.ContentText, value.ContentMarkdown, value.SearchText, value.SearchMarkdown, value.MetadataJSON, value.ContentHash)
	return errors.Wrapf(err, "insert document revision %s", value.ID)
}

func ensureSnapshot(ctx context.Context, tx *sql.Tx, id, sourceID, selectionJSON, manifestJSON string, revisionIDs []string) (bool, error) {
	var existingSource, existingSelection, existingManifest string
	var existingCount int
	err := tx.QueryRowContext(ctx, `SELECT source_artifact_id, selection_json, manifest_json, document_count FROM corpus_snapshots WHERE id = ?`, id).Scan(&existingSource, &existingSelection, &existingManifest, &existingCount)
	if err == nil {
		if existingSource != sourceID || existingSelection != selectionJSON || existingManifest != manifestJSON || existingCount != len(revisionIDs) {
			return false, errors.Errorf("corpus snapshot %s conflicts with existing immutable content", id)
		}
		rows, queryErr := tx.QueryContext(ctx, `SELECT document_revision_id FROM corpus_snapshot_documents WHERE snapshot_id = ? ORDER BY ordinal`, id)
		if queryErr != nil {
			return false, errors.Wrap(queryErr, "read existing snapshot membership")
		}
		defer func() { _ = rows.Close() }()
		var actual []string
		for rows.Next() {
			var revisionID string
			if scanErr := rows.Scan(&revisionID); scanErr != nil {
				return false, errors.Wrap(scanErr, "scan existing snapshot membership")
			}
			actual = append(actual, revisionID)
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return false, errors.Wrap(rowsErr, "iterate existing snapshot membership")
		}
		if strings.Join(actual, "\x00") != strings.Join(revisionIDs, "\x00") {
			return false, errors.Errorf("corpus snapshot %s has conflicting membership", id)
		}
		return true, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, errors.Wrap(err, "look up corpus snapshot")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO corpus_snapshots (id, schema_version, source_artifact_id, selection_json, manifest_json, document_count) VALUES (?, ?, ?, ?, ?, ?)`, id, corpusSnapshotSchema, sourceID, selectionJSON, manifestJSON, len(revisionIDs)); err != nil {
		return false, errors.Wrap(err, "insert corpus snapshot")
	}
	for ordinal, revisionID := range revisionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO corpus_snapshot_documents (snapshot_id, ordinal, document_revision_id) VALUES (?, ?, ?)`, id, ordinal, revisionID); err != nil {
			return false, errors.Wrapf(err, "insert snapshot membership ordinal %d", ordinal)
		}
	}
	return false, nil
}

func validatePlan(plan *ttcimport.Plan) error {
	if plan.Manifest.SchemaVersion == "" || plan.Manifest.SourceExportSHA256 == "" {
		return errors.New("TTC plan is missing manifest schema version or source checksum")
	}
	seen := make(map[string]struct{}, len(plan.Documents))
	for _, document := range plan.Documents {
		if document.ID == "" || document.Kind == "" || document.Title == "" || document.SearchText == "" {
			return errors.Errorf("invalid TTC document %q", document.ID)
		}
		if _, ok := seen[document.ID]; ok {
			return errors.Errorf("duplicate TTC document %q", document.ID)
		}
		seen[document.ID] = struct{}{}
	}
	return nil
}

func marshal(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "marshal immutable corpus object")
	}
	return string(bytes), nil
}
func mustJSON(value interface{}) string {
	encoded, err := marshal(value)
	if err != nil {
		panic(err)
	}
	return encoded
}
func fingerprint(domain string, parts ...string) string {
	digest := sha256.New()
	_, _ = digest.Write([]byte(domain))
	_, _ = digest.Write([]byte{0})
	for _, part := range parts {
		_, _ = digest.Write([]byte(strconv.Itoa(len(part))))
		_, _ = digest.Write([]byte{':'})
		_, _ = digest.Write([]byte(part))
		_, _ = digest.Write([]byte{0})
	}
	return "sha256:" + hex.EncodeToString(digest.Sum(nil))
}
