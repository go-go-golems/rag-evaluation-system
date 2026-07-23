package db

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestCleanBaselineExcludesRetiredExperimentLifecycle(t *testing.T) {
	database, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "clean.db")+"?_foreign_keys=on")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := Migrate(database); err != nil {
		t.Fatal(err)
	}
	experiment := "experiment_"
	for _, retired := range []string{experiment + "specs", experiment + "runs", experiment + "run_events", experiment + "run_summaries", experiment + "query_traces", "idx_" + experiment + "runs_spec_created", "idx_" + experiment + "events_run_sequence", experiment + "specs_no_update", experiment + "specs_no_delete", experiment + "runs_no_update", experiment + "runs_no_delete", experiment + "events_no_update", experiment + "events_no_delete", experiment + "summaries_no_update", experiment + "summaries_no_delete", experiment + "query_traces_no_update", experiment + "query_traces_no_delete"} {
		var count int
		if err := database.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = ?`, retired).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("retired database object remains: %s", retired)
		}
	}
	for _, current := range []string{"corpus_snapshots", "chunk_sets", "embedding_sets", "retrieval_artifacts", "evaluation_datasets", "representation_sets"} {
		var count int
		if err := database.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, current).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("current baseline table missing: %s", current)
		}
	}
}
