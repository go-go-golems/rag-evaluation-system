package researchctladapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	_ "github.com/mattn/go-sqlite3"
)

const TTCCatalogNamespace = "rag-eval-ttc"

type TTCCatalog struct{ database string }

func NewTTCCatalog(database string) *TTCCatalog { return &TTCCatalog{database: database} }
func (c *TTCCatalog) Resolve(ctx context.Context, reference InputReference, artifactRoot string) (ResolvedInput, error) {
	if reference.Catalog == nil || reference.Catalog.Namespace != TTCCatalogNamespace {
		return ResolvedInput{}, fmt.Errorf("RAG_TTC_CATALOG: unsupported catalog")
	}
	db, err := c.open(ctx)
	if err != nil {
		return ResolvedInput{}, err
	}
	defer func() { _ = db.Close() }()
	switch reference.Role {
	case "corpus":
		return c.resolveCorpus(ctx, db, reference, artifactRoot)
	case "evaluation-dataset", "judgments":
		return c.resolveEvaluation(ctx, db, reference, artifactRoot)
	default:
		return ResolvedInput{}, fmt.Errorf("RAG_TTC_ROLE: %s", reference.Role)
	}
}
func (c *TTCCatalog) open(ctx context.Context) (*sql.DB, error) {
	absolute, err := filepath.Abs(c.database)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("mode", "ro")
	query.Set("_query_only", "1")
	db, err := sql.Open("sqlite3", "file:"+filepath.ToSlash(absolute)+"?"+query.Encode())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
func (c *TTCCatalog) resolveCorpus(ctx context.Context, db *sql.DB, reference InputReference, artifactRoot string) (ResolvedInput, error) {
	corpus, _, err := loadTTCCorpus(ctx, db, reference.Catalog.ID)
	if err != nil {
		return ResolvedInput{}, err
	}
	envelope := ragoperators.NewCorpusArtifact(corpus, TTCCatalogNamespace+":"+reference.Catalog.ID)
	data, _ := json.Marshal(envelope)
	reference.ID = reference.Catalog.ID
	reference.SchemaVersion = ragcontract.CorpusManifestSchema
	return stageEnvelope(reference, data, artifactRoot)
}

type ttcCard struct {
	ID                          string   `json:"id"`
	Query                       string   `json:"query"`
	RelevantDocumentRevisionIDs []string `json:"relevantDocumentRevisionIds"`
}
type ttcEvaluationManifest struct {
	Cards []ttcCard `json:"cards"`
}

func (c *TTCCatalog) resolveEvaluation(ctx context.Context, db *sql.DB, reference InputReference, artifactRoot string) (ResolvedInput, error) {
	var corpusID, status, manifestJSON string
	if err := db.QueryRowContext(ctx, `SELECT corpus_snapshot_id,status,manifest_json FROM evaluation_datasets WHERE id=?`, reference.Catalog.ID).Scan(&corpusID, &status, &manifestJSON); err != nil {
		return ResolvedInput{}, fmt.Errorf("RAG_TTC_EVALUATION_MISSING: %w", err)
	}
	corpus, unitIDs, err := loadTTCCorpus(ctx, db, corpusID)
	if err != nil {
		return ResolvedInput{}, err
	}
	corpusArtifact := ragoperators.NewCorpusArtifact(corpus, TTCCatalogNamespace+":"+corpusID)
	var manifest ttcEvaluationManifest
	if err := json.Unmarshal([]byte(manifestJSON), &manifest); err != nil {
		return ResolvedInput{}, err
	}
	dataset := ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-queries/v2"}
	for _, card := range manifest.Cards {
		query := ragoperators.Query{ID: card.ID, Text: card.Query, Grades: map[string]float64{}}
		for _, revision := range card.RelevantDocumentRevisionIDs {
			if unitID := unitIDs[revision]; unitID != "" {
				query.RelevantIDs = append(query.RelevantIDs, unitID)
				query.Grades[unitID] = 1
			}
		}
		sort.Strings(query.RelevantIDs)
		dataset.Queries = append(dataset.Queries, query)
	}
	split := "candidate"
	if strings.Contains(reference.Catalog.ID, "baseline") {
		split = "smoke"
	}
	envelope := ragoperators.NewEvaluationArtifact(dataset, reference.Catalog.ID, split, status, "unit", corpusArtifact.Manifest.Digest)
	data, _ := json.Marshal(envelope)
	reference.ID = reference.Catalog.ID
	reference.SchemaVersion = ragcontract.EvaluationManifestSchema
	return stageEnvelope(reference, data, artifactRoot)
}
func loadTTCCorpus(ctx context.Context, db *sql.DB, corpusID string) (ragoperators.Corpus, map[string]string, error) {
	rows, err := db.QueryContext(ctx, `SELECT d.id,d.stable_document_id,csd.ordinal,d.kind,d.content_text,d.title,d.url,d.metadata_json FROM corpus_snapshot_documents csd JOIN document_revisions d ON d.id=csd.document_revision_id WHERE csd.snapshot_id=? ORDER BY csd.ordinal`, corpusID)
	if err != nil {
		return ragoperators.Corpus{}, nil, err
	}
	defer func() { _ = rows.Close() }()
	corpus := ragoperators.Corpus{SchemaVersion: "rag-source-record-set/v2"}
	unitIDs := map[string]string{}
	for rows.Next() {
		var record ragoperators.SourceRecord
		var metadataJSON, title, urlValue string
		if err := rows.Scan(&record.ID, &record.SessionID, &record.Ordinal, &record.Role, &record.Text, &title, &urlValue, &metadataJSON); err != nil {
			return ragoperators.Corpus{}, nil, err
		}
		var metadata map[string]any
		if json.Unmarshal([]byte(metadataJSON), &metadata) != nil || metadata == nil {
			metadata = map[string]any{}
		}
		metadata["title"] = title
		metadata["url"] = urlValue
		record.Metadata, _ = json.Marshal(metadata)
		corpus.Records = append(corpus.Records, record)
		digest, _ := ragcontract.Digest(record.Text)
		identity, _ := ragcontract.Digest(struct {
			Kind   string
			IDs    []string
			Digest string
		}{"units.identity", []string{record.ID}, digest})
		unitIDs[record.ID] = "unit:" + identity[7:23]
	}
	if err := rows.Err(); err != nil {
		return ragoperators.Corpus{}, nil, err
	}
	if len(corpus.Records) == 0 {
		return ragoperators.Corpus{}, nil, fmt.Errorf("RAG_TTC_CORPUS_MISSING: %s", corpusID)
	}
	return corpus, unitIDs, nil
}
