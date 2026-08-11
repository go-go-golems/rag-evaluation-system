//go:build ignore

// Registers the reviewed TTC candidate cards as one immutable, corpus-bound
// evaluation-dataset manifest. It is intentionally ticket-local because the
// cards are candidate evidence, not a general production import format.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultSnapshot = "sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409"

type card struct {
	ID                          string   `json:"id"`
	Query                       string   `json:"query"`
	RelevantDocumentRevisionIDs []string `json:"relevantDocumentRevisionIds"`
}

type manifest struct {
	SchemaVersion           string `json:"schemaVersion"`
	DatasetStatus           string `json:"datasetStatus"`
	BinaryRelevantAtOrAbove string `json:"binaryRelevantAtOrAbove"`
	Cards                   []card `json:"cards"`
}

func main() {
	databasePath := flag.String("db", "data/rag-eval.db", "rag-eval SQLite database")
	cardsPath := flag.String("cards", "ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md", "candidate card Markdown")
	datasetID := flag.String("dataset-id", "candidate:ttc-baseline-v1", "immutable evaluation dataset ID")
	snapshotID := flag.String("corpus-snapshot-id", defaultSnapshot, "corpus snapshot that owns the cards")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	parsed, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(parsed)

	cards, err := parseCards(*cardsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("parse candidate cards")
	}
	database, err := db.OpenDB(*databasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("migrate database")
	}
	if err := resolveRevisionIDs(context.Background(), database, *snapshotID, cards); err != nil {
		log.Fatal().Err(err).Msg("resolve card evidence in corpus snapshot")
	}
	data := manifest{SchemaVersion: "rag-eval-evaluation-dataset/v1", DatasetStatus: "candidate", BinaryRelevantAtOrAbove: "2_SUBSTANTIAL", Cards: cards}
	manifestJSON, err := experiments.CanonicalJSON(data)
	if err != nil {
		log.Fatal().Err(err).Msg("canonicalize dataset manifest")
	}
	var existing string
	err = database.QueryRowContext(context.Background(), `SELECT manifest_json FROM evaluation_datasets WHERE id=?`, *datasetID).Scan(&existing)
	switch {
	case err == nil:
		if existing != string(manifestJSON) {
			log.Fatal().Str("dataset_id", *datasetID).Msg("immutable dataset ID already has a different manifest")
		}
		log.Info().Str("dataset_id", *datasetID).Int("query_count", len(cards)).Msg("candidate evaluation dataset already registered")
	case err == sql.ErrNoRows:
		if _, err := database.ExecContext(context.Background(), `INSERT INTO evaluation_datasets (id,schema_version,corpus_snapshot_id,status,manifest_json,query_count) VALUES (?,?,?,?,?,?)`, *datasetID, data.SchemaVersion, *snapshotID, data.DatasetStatus, string(manifestJSON), len(cards)); err != nil {
			log.Fatal().Err(err).Msg("insert immutable evaluation dataset")
		}
		log.Info().Str("dataset_id", *datasetID).Int("query_count", len(cards)).Msg("registered candidate evaluation dataset")
	default:
		log.Fatal().Err(err).Msg("read evaluation dataset")
	}
}

func parseCards(path string) ([]card, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	idPattern := regexp.MustCompile("^#### `([^`]+)`")
	queryPattern := regexp.MustCompile(`^query: "(.*)"`)
	judgmentPattern := regexp.MustCompile("^- `([0-3])_(?:NOT_RELEVANT|PARTIAL|SUBSTANTIAL|AUTHORITATIVE)` .*`(wp:[0-9]+)`")
	var result []card
	current := card{}
	flush := func() {
		if current.ID == "" || current.Query == "" {
			return
		}
		sort.Strings(current.RelevantDocumentRevisionIDs)
		result = append(result, current)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match := idPattern.FindStringSubmatch(line); match != nil {
			flush()
			current = card{ID: match[1]}
			continue
		}
		if match := queryPattern.FindStringSubmatch(line); match != nil && current.ID != "" {
			current.Query = match[1]
			continue
		}
		if match := judgmentPattern.FindStringSubmatch(line); match != nil && current.ID != "" && match[1] >= "2" {
			current.RelevantDocumentRevisionIDs = append(current.RelevantDocumentRevisionIDs, match[2])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	flush()
	if len(result) != 20 {
		return nil, fmt.Errorf("candidate card count is not 20: got %d", len(result))
	}
	return result, nil
}

func resolveRevisionIDs(ctx context.Context, database *sql.DB, snapshotID string, cards []card) error {
	for cardIndex := range cards {
		for evidenceIndex, stableID := range cards[cardIndex].RelevantDocumentRevisionIDs {
			var revisionID string
			err := database.QueryRowContext(ctx, `SELECT dr.id FROM corpus_snapshot_documents csd JOIN document_revisions dr ON dr.id=csd.document_revision_id WHERE csd.snapshot_id=? AND dr.stable_document_id=?`, snapshotID, stableID).Scan(&revisionID)
			if err != nil {
				return err
			}
			cards[cardIndex].RelevantDocumentRevisionIDs[evidenceIndex] = revisionID
		}
	}
	return nil
}
