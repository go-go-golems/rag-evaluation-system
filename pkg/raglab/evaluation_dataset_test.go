package raglab

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

func TestLoadEvaluationCardsReadsAndOrdersImmutableManifest(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO source_artifacts (id,schema_version,kind,checksum_sha256,byte_size,manifest_json) VALUES ('source','v1','fixture','fixture',1,'{}')`); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO corpus_snapshots (id,schema_version,source_artifact_id,selection_json,manifest_json,document_count) VALUES ('snapshot','v1','source','{}','{}',0)`); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schemaVersion":"rag-eval-evaluation-dataset/v1","cards":[{"id":"b","query":"second"},{"id":"a","query":"first","relevantDocumentRevisionIds":["revision"]}]}`
	if _, err := database.Exec(`INSERT INTO evaluation_datasets (id,schema_version,corpus_snapshot_id,status,manifest_json,query_count) VALUES ('dataset','v1','snapshot','candidate',?,2)`, manifest); err != nil {
		t.Fatal(err)
	}
	cards, err := LoadEvaluationCards(context.Background(), db.NewQueries(database), "dataset")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 2 || cards[0].ID != "a" || cards[0].RelevantDocumentRevisionIDs[0] != "revision" {
		t.Fatalf("cards = %#v", cards)
	}
}
