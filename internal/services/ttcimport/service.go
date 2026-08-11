// Package ttcimport imports a deterministic TTC baseline selection from the
// rich WordPress/WooCommerce SQLite export into the RAG Evaluation database.
package ttcimport

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	// ManifestSchemaVersion identifies the document selection contract produced
	// before immutable document revisions are introduced.
	ManifestSchemaVersion = "ttc-snapshot-manifest/v1"

	// DefaultSnapshotName identifies the first bounded TTC retrieval baseline.
	DefaultSnapshotName = "ttc-baseline-v1"

	// DefaultSourceID identifies imported TTC documents in the current mutable
	// operational database. Task 3ydv replaces this operational representation
	// with immutable document revisions and snapshots.
	DefaultSourceID = "ttc-wordpress-rag"
)

var kindOrder = []string{"ttc_guide", "faq", "post", "product", "page"}

// DefaultKindQuotas define the 200-document source-balanced TTC baseline.
var DefaultKindQuotas = map[string]int{
	"ttc_guide": 19,
	"faq":       35,
	"post":      48,
	"product":   80,
	"page":      18,
}

// DefaultEvaluationSeedDocumentIDs are every source document currently named
// by the 20 source-validated candidate evaluation cards. Keeping the list in
// the importer ensures a deterministic filler selection cannot exclude a
// document that the future evaluation dataset must resolve.
var DefaultEvaluationSeedDocumentIDs = []string{
	"wp:3699", "wp:3701", "wp:549614", "wp:3709", "wp:552438",
	"wp:15947", "wp:3703", "wp:7347", "wp:26028", "wp:3717",
	"wp:10069", "wp:812290", "wp:4131", "wp:627148", "wp:4133",
	"wp:4134", "wp:9892", "wp:28084", "wp:751617", "wp:398454",
	"wp:405431", "wp:405509", "wp:405437", "wp:8017", "wp:15288",
	"wp:418694", "wp:19387", "wp:9688", "wp:4355", "wp:224522",
	"wp:4237", "wp:4116", "wp:76495", "wp:76497", "wp:456943",
	"wp:558351", "wp:398600", "wp:398593", "wp:270766", "wp:4140",
	"wp:398551",
}

// SourceDocument is the subset of a source-export document required by the
// importer and manifest.
type SourceDocument struct {
	ID              string
	Kind            string
	Title           string
	URL             string
	SearchText      string
	SearchMarkdown  string
	ContentMarkdown string
}

// ManifestDocument records the exact source document selected for a baseline.
// Hashes allow a later immutable-revision importer to prove which source bytes
// were available when the manifest was created.
type ManifestDocument struct {
	DocID                 string `json:"doc_id"`
	Kind                  string `json:"kind"`
	Title                 string `json:"title"`
	URL                   string `json:"url"`
	SearchTextSHA256      string `json:"search_text_sha256"`
	SearchMarkdownSHA256  string `json:"search_markdown_sha256"`
	ContentMarkdownSHA256 string `json:"content_markdown_sha256"`
}

// Manifest is deterministic JSON input for the later immutable snapshot work.
// It intentionally has no generated timestamp or source filesystem path: those
// values would make byte-identical source selections appear different.
type Manifest struct {
	SchemaVersion      string             `json:"schema_version"`
	SnapshotName       string             `json:"snapshot_name"`
	SourceExportSHA256 string             `json:"source_export_sha256"`
	SelectionAlgorithm string             `json:"selection_algorithm"`
	KindQuotas         map[string]int     `json:"kind_quotas"`
	SeedDocumentIDs    []string           `json:"seed_document_ids"`
	Documents          []ManifestDocument `json:"documents"`
}

// Plan holds source documents and the manifest that selects them.
type Plan struct {
	Manifest  Manifest
	Documents []SourceDocument
}

// BuildRequest configures deterministic source selection.
type BuildRequest struct {
	SourceDBPath                  string
	SnapshotName                  string
	KindQuotas                    map[string]int
	IncludeDefaultEvaluationSeeds bool
	AdditionalSeedDocumentIDs     []string
}

// PersistRequest configures import of a deterministic plan into the current
// operational database.
type PersistRequest struct {
	SourceID     string
	SourceName   string
	ManifestPath string
}

// ImportResult summarizes a completed or dry-run import.
type ImportResult struct {
	SnapshotName       string
	SourceExportSHA256 string
	DocumentCount      int
	KindCounts         map[string]int
	ManifestPath       string
	DryRun             bool
}

// BuildPlan opens the rich TTC export read-only, validates source documents,
// and produces a deterministic baseline membership plan.
func BuildPlan(ctx context.Context, request BuildRequest) (*Plan, error) {
	if request.SourceDBPath == "" {
		return nil, errors.New("source SQLite database path is required")
	}
	if request.SnapshotName == "" {
		request.SnapshotName = DefaultSnapshotName
	}
	quotas, err := normalizeQuotas(request.KindQuotas)
	if err != nil {
		return nil, err
	}

	sourceHash, err := fileSHA256(request.SourceDBPath)
	if err != nil {
		return nil, errors.Wrap(err, "hash TTC source export")
	}

	source, err := sql.Open("sqlite3", "file:"+request.SourceDBPath+"?mode=ro")
	if err != nil {
		return nil, errors.Wrap(err, "open TTC source export read-only")
	}
	defer func() { _ = source.Close() }()
	if err := source.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "ping TTC source export")
	}

	documents, err := loadSourceDocuments(ctx, source)
	if err != nil {
		return nil, err
	}
	seedDocumentIDs := append([]string{}, request.AdditionalSeedDocumentIDs...)
	if request.IncludeDefaultEvaluationSeeds {
		seedDocumentIDs = append(seedDocumentIDs, DefaultEvaluationSeedDocumentIDs...)
	}
	selected, seedIDs, err := selectDocuments(documents, quotas, seedDocumentIDs)
	if err != nil {
		return nil, err
	}

	manifestDocuments := make([]ManifestDocument, 0, len(selected))
	for _, document := range selected {
		manifestDocuments = append(manifestDocuments, ManifestDocument{
			DocID:                 document.ID,
			Kind:                  document.Kind,
			Title:                 document.Title,
			URL:                   document.URL,
			SearchTextSHA256:      stringSHA256(document.SearchText),
			SearchMarkdownSHA256:  stringSHA256(document.SearchMarkdown),
			ContentMarkdownSHA256: stringSHA256(document.ContentMarkdown),
		})
	}

	return &Plan{
		Manifest: Manifest{
			SchemaVersion:      ManifestSchemaVersion,
			SnapshotName:       request.SnapshotName,
			SourceExportSHA256: sourceHash,
			SelectionAlgorithm: "seed documents first; remaining candidates ordered by sha256(doc_id) within fixed kind quotas",
			KindQuotas:         quotas,
			SeedDocumentIDs:    seedIDs,
			Documents:          manifestDocuments,
		},
		Documents: selected,
	}, nil
}

// Persist imports a plan into the existing operational documents table and
// atomically publishes its manifest. This is intentionally an import boundary,
// not the later immutable document-revision implementation.
func Persist(ctx context.Context, queries *db.Queries, plan *Plan, request PersistRequest) (*ImportResult, error) {
	if queries == nil {
		return nil, errors.New("target database queries are required")
	}
	if plan == nil {
		return nil, errors.New("TTC import plan is required")
	}
	if request.SourceID == "" {
		request.SourceID = DefaultSourceID
	}
	if request.SourceName == "" {
		request.SourceName = "TTC WordPress RAG baseline"
	}
	if request.ManifestPath == "" {
		return nil, errors.New("manifest path is required")
	}

	configJSON, err := json.Marshal(map[string]string{
		"manifest_schema_version": plan.Manifest.SchemaVersion,
		"snapshot_name":           plan.Manifest.SnapshotName,
		"source_export_sha256":    plan.Manifest.SourceExportSHA256,
	})
	if err != nil {
		return nil, errors.Wrap(err, "marshal TTC source configuration")
	}

	tx, err := queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin TTC import transaction")
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO sources (id, name, type, config_json)
		VALUES (?, ?, 'ttc-wordpress-rag', ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			type = excluded.type,
			config_json = excluded.config_json,
			updated_at = datetime('now')
	`, request.SourceID, request.SourceName, string(configJSON)); err != nil {
		return nil, errors.Wrap(err, "upsert TTC source")
	}

	for _, document := range plan.Documents {
		metadataJSON, err := json.Marshal(map[string]string{
			"source_content_markdown_sha256": stringSHA256(document.ContentMarkdown),
			"source_document_id":             document.ID,
			"source_document_kind":           document.Kind,
			"source_search_markdown_sha256":  stringSHA256(document.SearchMarkdown),
			"source_search_text_sha256":      stringSHA256(document.SearchText),
			"source_export_sha256":           plan.Manifest.SourceExportSHA256,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "marshal metadata for %s", document.ID)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO documents (
				id, source_id, external_id, title, url, content_type,
				raw_content, content_text, content_html, word_count,
				language, metadata_json, status
			) VALUES (?, ?, ?, ?, ?, 'text/markdown', ?, ?, '', ?, 'en', ?, 'imported')
			ON CONFLICT(id) DO UPDATE SET
				source_id = excluded.source_id,
				external_id = excluded.external_id,
				title = excluded.title,
				url = excluded.url,
				content_type = excluded.content_type,
				raw_content = excluded.raw_content,
				content_text = excluded.content_text,
				content_html = excluded.content_html,
				word_count = excluded.word_count,
				language = excluded.language,
				metadata_json = excluded.metadata_json,
				status = excluded.status,
				updated_at = datetime('now')
		`, targetDocumentID(document.ID), request.SourceID, document.ID, document.Title, document.URL,
			document.ContentMarkdown, document.SearchText, wordCount(document.SearchText), string(metadataJSON)); err != nil {
			return nil, errors.Wrapf(err, "upsert TTC document %s", document.ID)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit TTC import")
	}
	if err := WriteManifest(request.ManifestPath, plan.Manifest); err != nil {
		return nil, err
	}

	return importResult(plan, request.ManifestPath, false), nil
}

// WriteManifest atomically writes a deterministic, human-readable manifest.
func WriteManifest(path string, manifest Manifest) error {
	if path == "" {
		return errors.New("manifest path is required")
	}
	contents, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return errors.Wrap(err, "marshal TTC manifest")
	}
	contents = append(contents, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return errors.Wrap(err, "create manifest directory")
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".ttc-manifest-*.json")
	if err != nil {
		return errors.Wrap(err, "create temporary manifest")
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	if _, err := temporary.Write(contents); err != nil {
		_ = temporary.Close()
		return errors.Wrap(err, "write temporary manifest")
	}
	if err := temporary.Close(); err != nil {
		return errors.Wrap(err, "close temporary manifest")
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return errors.Wrap(err, "publish TTC manifest")
	}
	return nil
}

// DryRunResult turns a plan into a command result without writing the target
// database or manifest.
func DryRunResult(plan *Plan, manifestPath string) *ImportResult {
	return importResult(plan, manifestPath, true)
}

func importResult(plan *Plan, manifestPath string, dryRun bool) *ImportResult {
	kindCounts := make(map[string]int, len(kindOrder))
	for _, document := range plan.Documents {
		kindCounts[document.Kind]++
	}
	return &ImportResult{
		SnapshotName:       plan.Manifest.SnapshotName,
		SourceExportSHA256: plan.Manifest.SourceExportSHA256,
		DocumentCount:      len(plan.Documents),
		KindCounts:         kindCounts,
		ManifestPath:       manifestPath,
		DryRun:             dryRun,
	}
}

func normalizeQuotas(input map[string]int) (map[string]int, error) {
	quotas := make(map[string]int, len(kindOrder))
	if len(input) == 0 {
		input = DefaultKindQuotas
	}
	for _, kind := range kindOrder {
		quota, ok := input[kind]
		if !ok || quota <= 0 {
			return nil, errors.Errorf("quota for kind %q must be positive", kind)
		}
		quotas[kind] = quota
	}
	for kind := range input {
		if _, ok := quotas[kind]; !ok {
			return nil, errors.Errorf("unsupported TTC kind quota %q", kind)
		}
	}
	return quotas, nil
}

func loadSourceDocuments(ctx context.Context, source *sql.DB) ([]SourceDocument, error) {
	rows, err := source.QueryContext(ctx, `
		SELECT doc_id, kind, title, COALESCE(url, ''),
		       COALESCE(search_text, ''), COALESCE(search_markdown, ''), COALESCE(content_markdown, '')
		FROM documents
		WHERE kind IN ('ttc_guide', 'faq', 'post', 'product', 'page')
	`)
	if err != nil {
		return nil, errors.Wrap(err, "query TTC source documents")
	}
	defer func() { _ = rows.Close() }()

	documents := make([]SourceDocument, 0)
	for rows.Next() {
		var document SourceDocument
		if err := rows.Scan(
			&document.ID,
			&document.Kind,
			&document.Title,
			&document.URL,
			&document.SearchText,
			&document.SearchMarkdown,
			&document.ContentMarkdown,
		); err != nil {
			return nil, errors.Wrap(err, "scan TTC source document")
		}
		if document.ID == "" || document.Title == "" || document.SearchText == "" {
			return nil, errors.Errorf("TTC source document must have id, title, and search text: %q", document.ID)
		}
		documents = append(documents, document)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate TTC source documents")
	}
	return documents, nil
}

func selectDocuments(documents []SourceDocument, quotas map[string]int, seedDocumentIDs []string) ([]SourceDocument, []string, error) {
	byID := make(map[string]SourceDocument, len(documents))
	byKind := make(map[string][]SourceDocument, len(kindOrder))
	for _, document := range documents {
		if _, exists := byID[document.ID]; exists {
			return nil, nil, errors.Errorf("duplicate TTC source document id %q", document.ID)
		}
		byID[document.ID] = document
		byKind[document.Kind] = append(byKind[document.Kind], document)
	}

	seedSet := make(map[string]struct{}, len(seedDocumentIDs))
	for _, documentID := range seedDocumentIDs {
		seedSet[documentID] = struct{}{}
	}
	seedIDs := make([]string, 0, len(seedSet))
	for documentID := range seedSet {
		if _, ok := byID[documentID]; !ok {
			return nil, nil, errors.Errorf("seed TTC source document %q was not found", documentID)
		}
		seedIDs = append(seedIDs, documentID)
	}
	sort.Strings(seedIDs)

	selected := make([]SourceDocument, 0)
	for _, kind := range kindOrder {
		quota := quotas[kind]
		candidates := byKind[kind]
		if len(candidates) < quota {
			return nil, nil, errors.Errorf("TTC source has %d %s documents, need quota %d", len(candidates), kind, quota)
		}
		seeded := make([]SourceDocument, 0)
		fillers := make([]SourceDocument, 0)
		for _, candidate := range candidates {
			if _, ok := seedSet[candidate.ID]; ok {
				seeded = append(seeded, candidate)
			} else {
				fillers = append(fillers, candidate)
			}
		}
		if len(seeded) > quota {
			return nil, nil, errors.Errorf("%d seed documents of kind %q exceed quota %d", len(seeded), kind, quota)
		}
		sort.Slice(seeded, func(i, j int) bool { return seeded[i].ID < seeded[j].ID })
		sort.Slice(fillers, func(i, j int) bool {
			left := stringSHA256(fillers[i].ID)
			right := stringSHA256(fillers[j].ID)
			if left == right {
				return fillers[i].ID < fillers[j].ID
			}
			return left < right
		})
		selected = append(selected, seeded...)
		selected = append(selected, fillers[:quota-len(seeded)]...)
	}
	return selected, seedIDs, nil
}

func targetDocumentID(sourceDocumentID string) string {
	return "ttc:" + sourceDocumentID
}

func stringSHA256(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "open source export")
	}
	defer func() { _ = file.Close() }()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", errors.Wrap(err, "hash source export")
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func wordCount(text string) int {
	return len(strings.Fields(text))
}
