// Package experimentrun persists immutable experiment specifications and
// append-only observations of their execution.
package experimentrun

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experiments"
	"github.com/pkg/errors"
)

const specificationSchema = "rag-eval-experiment-spec/v1"

type Service struct{ queries *db.Queries }

func NewService(queries *db.Queries) *Service { return &Service{queries: queries} }

type SpecificationInput struct {
	CorpusSnapshotID    string         `json:"corpus_snapshot_id"`
	ChunkSetID          string         `json:"chunk_set_id"`
	BM25ArtifactID      string         `json:"bm25_artifact_id,omitempty"`
	EmbeddingSetID      string         `json:"embedding_set_id,omitempty"`
	EvaluationDatasetID string         `json:"evaluation_dataset_id,omitempty"`
	Config              map[string]any `json:"config"`
}

type Specification struct {
	ID                  string         `json:"id"`
	SchemaVersion       string         `json:"schema_version"`
	CorpusSnapshotID    string         `json:"corpus_snapshot_id"`
	ChunkSetID          string         `json:"chunk_set_id"`
	BM25ArtifactID      string         `json:"bm25_artifact_id,omitempty"`
	EmbeddingSetID      string         `json:"embedding_set_id,omitempty"`
	EvaluationDatasetID string         `json:"evaluation_dataset_id,omitempty"`
	Config              map[string]any `json:"config"`
	CreatedAt           string         `json:"created_at"`
}

type Event struct {
	Sequence   int             `json:"sequence"`
	Type       string          `json:"type"`
	OccurredAt string          `json:"occurred_at"`
	Payload    json.RawMessage `json:"payload"`
}

type QueryTraceInput struct {
	QueryCardID string          `json:"query_card_id"`
	Trace       json.RawMessage `json:"trace"`
	Metrics     json.RawMessage `json:"metrics"`
	Timing      json.RawMessage `json:"timing"`
	Cost        json.RawMessage `json:"cost"`
	Storage     json.RawMessage `json:"storage"`
}

type QueryTrace struct {
	QueryTraceInput
}

type SummaryInput struct {
	Status  string          `json:"status"`
	Metrics json.RawMessage `json:"metrics"`
	Cost    json.RawMessage `json:"cost"`
	Storage json.RawMessage `json:"storage"`
	Error   json.RawMessage `json:"error"`
}

type Summary struct {
	SummaryInput
	FinishedAt string `json:"finished_at"`
}

type Run struct {
	ID               string   `json:"id"`
	ExperimentSpecID string   `json:"experiment_spec_id"`
	CreatedAt        string   `json:"created_at"`
	Status           string   `json:"status"`
	Events           []Event  `json:"events,omitempty"`
	Summary          *Summary `json:"summary,omitempty"`
}

type specManifest struct {
	SchemaVersion       string         `json:"schema_version"`
	CorpusSnapshotID    string         `json:"corpus_snapshot_id"`
	ChunkSetID          string         `json:"chunk_set_id"`
	BM25ArtifactID      string         `json:"bm25_artifact_id,omitempty"`
	EmbeddingSetID      string         `json:"embedding_set_id,omitempty"`
	EvaluationDatasetID string         `json:"evaluation_dataset_id,omitempty"`
	Config              map[string]any `json:"config"`
}

func (s *Service) CreateSpecification(ctx context.Context, input SpecificationInput) (*Specification, bool, error) {
	if s == nil || s.queries == nil {
		return nil, false, errors.New("experiment-run service requires database queries")
	}
	if err := validateSpecificationInput(ctx, s.queries, input); err != nil {
		return nil, false, err
	}
	if input.Config == nil {
		input.Config = map[string]any{}
	}
	manifest := specManifest{specificationSchema, input.CorpusSnapshotID, input.ChunkSetID, input.BM25ArtifactID, input.EmbeddingSetID, input.EvaluationDatasetID, input.Config}
	id, err := experiments.Fingerprint(specificationSchema, manifest)
	if err != nil {
		return nil, false, err
	}
	configJSON, err := experiments.CanonicalJSON(input.Config)
	if err != nil {
		return nil, false, err
	}
	manifestJSON, err := experiments.CanonicalJSON(manifest)
	if err != nil {
		return nil, false, err
	}
	var existing string
	err = s.queries.DB().QueryRowContext(ctx, `SELECT manifest_json FROM experiment_specs WHERE id = ?`, id).Scan(&existing)
	if err == nil {
		if existing != string(manifestJSON) {
			return nil, false, errors.New("experiment specification fingerprint conflict")
		}
		result, getErr := s.GetSpecification(ctx, id)
		return result, true, getErr
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, false, errors.Wrap(err, "read experiment specification")
	}
	if _, err := s.queries.DB().ExecContext(ctx, `INSERT INTO experiment_specs (id,schema_version,corpus_snapshot_id,chunk_set_id,bm25_artifact_id,embedding_set_id,evaluation_dataset_id,config_json,manifest_json) VALUES (?,?,?,?,?,?,?,?,?)`, id, specificationSchema, input.CorpusSnapshotID, input.ChunkSetID, nullable(input.BM25ArtifactID), nullable(input.EmbeddingSetID), input.EvaluationDatasetID, string(configJSON), string(manifestJSON)); err != nil {
		return nil, false, errors.Wrap(err, "insert experiment specification")
	}
	result, err := s.GetSpecification(ctx, id)
	return result, false, err
}

func (s *Service) GetSpecification(ctx context.Context, id string) (*Specification, error) {
	if id == "" {
		return nil, errors.New("experiment specification ID is required")
	}
	var result Specification
	var configJSON string
	err := s.queries.DB().QueryRowContext(ctx, `SELECT id,schema_version,corpus_snapshot_id,chunk_set_id,COALESCE(bm25_artifact_id,''),COALESCE(embedding_set_id,''),evaluation_dataset_id,config_json,created_at FROM experiment_specs WHERE id=?`, id).Scan(&result.ID, &result.SchemaVersion, &result.CorpusSnapshotID, &result.ChunkSetID, &result.BM25ArtifactID, &result.EmbeddingSetID, &result.EvaluationDatasetID, &configJSON, &result.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "get experiment specification")
	}
	if err := json.Unmarshal([]byte(configJSON), &result.Config); err != nil {
		return nil, errors.Wrap(err, "decode experiment specification config")
	}
	return &result, nil
}

func (s *Service) ListSpecifications(ctx context.Context) ([]Specification, error) {
	rows, err := s.queries.DB().QueryContext(ctx, `SELECT id FROM experiment_specs ORDER BY created_at DESC,id`)
	if err != nil {
		return nil, errors.Wrap(err, "list experiment specifications")
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	result := make([]Specification, 0, len(ids))
	for _, id := range ids {
		item, err := s.GetSpecification(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, *item)
	}
	return result, nil
}

func (s *Service) CreateRun(ctx context.Context, specificationID string) (*Run, error) {
	if _, err := s.GetSpecification(ctx, specificationID); err != nil {
		return nil, err
	}
	id, err := newRunID()
	if err != nil {
		return nil, err
	}
	tx, err := s.queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin experiment run")
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `INSERT INTO experiment_runs (id,experiment_spec_id) VALUES (?,?)`, id, specificationID); err != nil {
		return nil, errors.Wrap(err, "insert experiment run")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO experiment_run_events (run_id,sequence,event_type,occurred_at,payload_json) VALUES (?,?,?,?,?)`, id, 1, "created", now(), "{}"); err != nil {
		return nil, errors.Wrap(err, "append created event")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit experiment run")
	}
	return s.GetRun(ctx, id)
}

func (s *Service) AppendEvent(ctx context.Context, runID, eventType string, payload json.RawMessage) (*Event, error) {
	if runID == "" || eventType == "" {
		return nil, errors.New("run ID and event type are required")
	}
	payloadJSON, err := canonicalRawJSON(payload)
	if err != nil {
		return nil, errors.Wrap(err, "canonical event payload")
	}
	tx, err := s.queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin append event")
	}
	defer func() { _ = tx.Rollback() }()
	if finished, err := runFinished(ctx, tx, runID); err != nil {
		return nil, err
	} else if finished {
		return nil, errors.New("cannot append event after terminal summary")
	}
	var sequence int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(sequence), 0) + 1 FROM experiment_run_events WHERE run_id=?`, runID).Scan(&sequence); err != nil {
		return nil, errors.Wrap(err, "allocate event sequence")
	}
	occurredAt := now()
	if _, err := tx.ExecContext(ctx, `INSERT INTO experiment_run_events (run_id,sequence,event_type,occurred_at,payload_json) VALUES (?,?,?,?,?)`, runID, sequence, eventType, occurredAt, string(payloadJSON)); err != nil {
		return nil, errors.Wrap(err, "insert experiment event")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit experiment event")
	}
	return &Event{Sequence: sequence, Type: eventType, OccurredAt: occurredAt, Payload: payloadJSON}, nil
}

func (s *Service) RecordQueryTrace(ctx context.Context, runID string, input QueryTraceInput) error {
	if runID == "" || input.QueryCardID == "" {
		return errors.New("run ID and query card ID are required")
	}
	values, err := canonicalTraceInput(input)
	if err != nil {
		return err
	}
	if finished, err := runFinished(ctx, s.queries.DB(), runID); err != nil {
		return err
	} else if finished {
		return errors.New("cannot record trace after terminal summary")
	}
	_, err = s.queries.DB().ExecContext(ctx, `INSERT INTO experiment_query_traces (run_id,query_card_id,trace_json,metrics_json,timing_json,cost_json,storage_json) VALUES (?,?,?,?,?,?,?)`, runID, input.QueryCardID, values[0], values[1], values[2], values[3], values[4])
	return errors.Wrap(err, "insert immutable query trace")
}

func (s *Service) CompleteRun(ctx context.Context, runID string, input SummaryInput) (*Summary, error) {
	if runID == "" {
		return nil, errors.New("run ID is required")
	}
	if input.Status != "succeeded" && input.Status != "failed" && input.Status != "cancelled" {
		return nil, errors.New("terminal status must be succeeded, failed, or cancelled")
	}
	values, err := canonicalSummaryInput(input)
	if err != nil {
		return nil, err
	}
	tx, err := s.queries.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin complete run")
	}
	defer func() { _ = tx.Rollback() }()
	if finished, err := runFinished(ctx, tx, runID); err != nil {
		return nil, err
	} else if finished {
		return nil, errors.New("experiment run already has a terminal summary")
	}
	finishedAt := now()
	var sequence int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(sequence), 0) + 1 FROM experiment_run_events WHERE run_id=?`, runID).Scan(&sequence); err != nil {
		return nil, errors.Wrap(err, "allocate terminal event sequence")
	}
	payload := fmt.Sprintf(`{"status":%q}`, input.Status)
	if _, err := tx.ExecContext(ctx, `INSERT INTO experiment_run_events (run_id,sequence,event_type,occurred_at,payload_json) VALUES (?,?,?,?,?)`, runID, sequence, "terminal", finishedAt, payload); err != nil {
		return nil, errors.Wrap(err, "append terminal event")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO experiment_run_summaries (run_id,status,finished_at,metrics_json,cost_json,storage_json,error_json) VALUES (?,?,?,?,?,?,?)`, runID, input.Status, finishedAt, values[0], values[1], values[2], values[3]); err != nil {
		return nil, errors.Wrap(err, "insert immutable run summary")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit terminal summary")
	}
	return &Summary{SummaryInput: input, FinishedAt: finishedAt}, nil
}

func (s *Service) GetRun(ctx context.Context, id string) (*Run, error) {
	var result Run
	err := s.queries.DB().QueryRowContext(ctx, `SELECT id,experiment_spec_id,created_at FROM experiment_runs WHERE id=?`, id).Scan(&result.ID, &result.ExperimentSpecID, &result.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "get experiment run")
	}
	result.Status = "running"
	rows, err := s.queries.DB().QueryContext(ctx, `SELECT sequence,event_type,occurred_at,payload_json FROM experiment_run_events WHERE run_id=? ORDER BY sequence`, id)
	if err != nil {
		return nil, errors.Wrap(err, "list experiment events")
	}
	defer rows.Close()
	for rows.Next() {
		var item Event
		var payload string
		if err := rows.Scan(&item.Sequence, &item.Type, &item.OccurredAt, &payload); err != nil {
			return nil, err
		}
		item.Payload = json.RawMessage(payload)
		result.Events = append(result.Events, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var summary Summary
	var metrics, cost, storage, failure string
	err = s.queries.DB().QueryRowContext(ctx, `SELECT status,finished_at,metrics_json,cost_json,storage_json,error_json FROM experiment_run_summaries WHERE run_id=?`, id).Scan(&summary.Status, &summary.FinishedAt, &metrics, &cost, &storage, &failure)
	if err == nil {
		summary.Metrics, summary.Cost, summary.Storage, summary.Error = json.RawMessage(metrics), json.RawMessage(cost), json.RawMessage(storage), json.RawMessage(failure)
		result.Status, result.Summary = summary.Status, &summary
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "get experiment summary")
	}
	return &result, nil
}

func (s *Service) ListRuns(ctx context.Context, specificationID string) ([]Run, error) {
	query, args := `SELECT id FROM experiment_runs`, []any{}
	if specificationID != "" {
		query += ` WHERE experiment_spec_id=?`
		args = append(args, specificationID)
	}
	query += ` ORDER BY created_at DESC,id`
	rows, err := s.queries.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "list experiment runs")
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	result := make([]Run, 0, len(ids))
	for _, id := range ids {
		run, err := s.GetRun(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, *run)
	}
	return result, nil
}

func (s *Service) ListQueryTraces(ctx context.Context, runID string) ([]QueryTrace, error) {
	rows, err := s.queries.DB().QueryContext(ctx, `SELECT query_card_id,trace_json,metrics_json,timing_json,cost_json,storage_json FROM experiment_query_traces WHERE run_id=? ORDER BY query_card_id`, runID)
	if err != nil {
		return nil, errors.Wrap(err, "list query traces")
	}
	defer rows.Close()
	var result []QueryTrace
	for rows.Next() {
		var item QueryTrace
		var trace, metrics, timing, cost, storage string
		if err := rows.Scan(&item.QueryCardID, &trace, &metrics, &timing, &cost, &storage); err != nil {
			return nil, err
		}
		item.Trace, item.Metrics, item.Timing, item.Cost, item.Storage = json.RawMessage(trace), json.RawMessage(metrics), json.RawMessage(timing), json.RawMessage(cost), json.RawMessage(storage)
		result = append(result, item)
	}
	return result, rows.Err()
}

func validateSpecificationInput(ctx context.Context, queries *db.Queries, input SpecificationInput) error {
	if input.CorpusSnapshotID == "" || input.ChunkSetID == "" {
		return errors.New("corpus snapshot and chunk set IDs are required")
	}
	var count int
	if err := queries.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM chunk_sets WHERE id=? AND corpus_snapshot_id=?`, input.ChunkSetID, input.CorpusSnapshotID).Scan(&count); err != nil {
		return err
	}
	if count != 1 {
		return errors.New("chunk set does not belong to corpus snapshot")
	}
	for _, reference := range []struct{ ID, Table, Field string }{{input.BM25ArtifactID, "retrieval_artifacts", "id"}, {input.EmbeddingSetID, "embedding_sets", "id"}} {
		if reference.ID == "" {
			continue
		}
		var found int
		if err := queries.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM `+reference.Table+` WHERE `+reference.Field+`=? AND chunk_set_id=?`, reference.ID, input.ChunkSetID).Scan(&found); err != nil {
			return err
		}
		if found != 1 {
			return fmt.Errorf("%s does not belong to chunk set", reference.Table)
		}
	}
	if input.BM25ArtifactID == "" && input.EmbeddingSetID == "" {
		return errors.New("an experiment specification needs BM25 or embedding retrieval input")
	}
	return nil
}

func canonicalTraceInput(input QueryTraceInput) ([5]string, error) {
	values := []json.RawMessage{input.Trace, input.Metrics, input.Timing, input.Cost, input.Storage}
	var result [5]string
	for i, value := range values {
		canonical, err := canonicalRawJSON(value)
		if err != nil {
			return result, err
		}
		result[i] = string(canonical)
	}
	return result, nil
}
func canonicalSummaryInput(input SummaryInput) ([4]string, error) {
	values := []json.RawMessage{input.Metrics, input.Cost, input.Storage, input.Error}
	var result [4]string
	for i, value := range values {
		canonical, err := canonicalRawJSON(value)
		if err != nil {
			return result, err
		}
		result[i] = string(canonical)
	}
	return result, nil
}
func canonicalRawJSON(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		raw = json.RawMessage("{}")
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, errors.Wrap(err, "decode JSON")
	}
	canonical, err := experiments.CanonicalJSON(value)
	return json.RawMessage(canonical), err
}
func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}
func now() string { return time.Now().UTC().Format(time.RFC3339Nano) }
func newRunID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.Wrap(err, "generate run ID")
	}
	return "run_" + hex.EncodeToString(bytes), nil
}
func runFinished(ctx context.Context, queryer interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}, runID string) (bool, error) {
	var count int
	if err := queryer.QueryRowContext(ctx, `SELECT COUNT(*) FROM experiment_run_summaries WHERE run_id=?`, runID).Scan(&count); err != nil {
		return false, errors.Wrap(err, "check terminal summary")
	}
	return count != 0, nil
}
