package raglab

import (
	"context"
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/experimentspec"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/pkg/errors"
)

type OpenOptions struct {
	Database      string
	AllowRuns     bool
	QueryEmbedder QueryEmbedder
}

type ExperimentStore interface {
	CreateSpecification(context.Context, experimentspec.Input) (*experimentrun.Specification, bool, error)
	CreateRun(context.Context, string) (*experimentrun.Run, error)
	AppendEvent(context.Context, string, string, json.RawMessage) (*experimentrun.Event, error)
}

var _ ExperimentStore = (*experimentrun.Service)(nil)

type PersistedSpecification struct {
	Specification *experimentrun.Specification `json:"specification"`
	Reused        bool                         `json:"reused"`
}

// Laboratory makes side effects explicit. Builder construction never owns one;
// callers must opt into this handle to inspect artifacts or create a run.
type Laboratory struct {
	catalog   ArtifactCatalog
	store     ExperimentStore
	allowRuns bool
	close     func() error
	executor  *Executor
	embedder  QueryEmbedder
	loadCards func(context.Context, string) ([]EvaluationCard, error)
}

func NewLaboratory(catalog ArtifactCatalog, store ExperimentStore, allowRuns bool) *Laboratory {
	return &Laboratory{catalog: catalog, store: store, allowRuns: allowRuns}
}

// OpenSQLite opens an existing rag-eval database. It deliberately does not run
// migrations: opening a laboratory must not create or alter schema state.
func OpenSQLite(options OpenOptions) (*Laboratory, error) {
	if options.Database == "" {
		return nil, errors.New("RAG_DATABASE_REQUIRED: database path is required")
	}
	database, err := db.OpenDB(options.Database)
	if err != nil {
		return nil, errors.Wrap(err, "open rag laboratory database")
	}
	queries := db.NewQueries(database)
	service := experimentrun.NewService(queries)
	lab := NewLaboratory(NewSQLiteCatalog(queries), service, options.AllowRuns)
	lab.executor = NewExecutor(NewSQLiteChannelRetriever(queries), service)
	lab.loadCards = func(ctx context.Context, datasetID string) ([]EvaluationCard, error) {
		return LoadEvaluationCards(ctx, queries, datasetID)
	}
	if options.QueryEmbedder != nil {
		lab.embedder = options.QueryEmbedder
	}
	lab.close = database.Close
	return lab, nil
}

func (l *Laboratory) Close() error {
	if l == nil || l.close == nil {
		return nil
	}
	return l.close()
}

func (l *Laboratory) Validate(ctx context.Context, specification ExperimentSpecification) ValidationReport {
	if l == nil {
		return ValidationReport{Issues: []ValidationIssue{{Code: "RAG_LAB_REQUIRED", Path: "$", Message: "laboratory is required", Severity: ValidationErrorSeverity}}}
	}
	return specification.ValidateCompatibility(ctx, l.catalog)
}

func (l *Laboratory) Persist(ctx context.Context, specification ExperimentSpecification) (*PersistedSpecification, error) {
	if l == nil || l.store == nil {
		return nil, errors.New("RAG_LAB_REQUIRED: laboratory store is required")
	}
	if !l.allowRuns {
		return nil, errors.New("RAG_EXECUTION_DISABLED: laboratory is read-only")
	}
	report := l.Validate(ctx, specification)
	if !report.OK() {
		return nil, &ValidationError{Report: report}
	}
	persisted, reused, err := l.store.CreateSpecification(ctx, specification.PersistenceInput())
	if err != nil {
		return nil, errors.Wrap(err, "persist immutable experiment specification")
	}
	return &PersistedSpecification{Specification: persisted, Reused: reused}, nil
}

// Start persists (or reuses) a specification, then creates a distinct
// append-only run. Retrieval execution is deliberately delegated to the next
// durable executor layer; this method only records durable submission.
func (l *Laboratory) Start(ctx context.Context, specification ExperimentSpecification) (*experimentrun.Run, error) {
	persisted, err := l.Persist(ctx, specification)
	if err != nil {
		return nil, err
	}
	run, err := l.store.CreateRun(ctx, persisted.Specification.ID)
	if err != nil {
		return nil, errors.Wrap(err, "create immutable experiment run")
	}
	if _, err := l.store.AppendEvent(ctx, run.ID, "submitted", []byte(`{"executor":"pending"}`)); err != nil {
		return nil, errors.Wrap(err, "append experiment submission event")
	}
	return run, nil
}

// Execute records the retrieval observations and terminal summary for an
// already-created run. It requires a laboratory opened with an explicit query
// embedder whenever the plan contains vector channels.
func (l *Laboratory) Execute(ctx context.Context, runID string, specification ExperimentSpecification, cards []EvaluationCard, options ExecutionOptions) (*ExecutionResult, error) {
	if l == nil || l.executor == nil {
		return nil, errors.New("RAG_EXECUTOR_REQUIRED: open the laboratory with an explicit query embedder")
	}
	if !l.allowRuns {
		return nil, errors.New("RAG_EXECUTION_DISABLED: laboratory is read-only")
	}
	if options.Embedder == nil {
		options.Embedder = l.embedder
	}
	return l.executor.Execute(ctx, runID, specification, cards, options)
}

// Run loads the frozen evaluation cards selected by the specification, creates
// a fresh append-only run, and executes it. It is the single-call execution
// boundary for user-facing clients such as the JavaScript laboratory module.
func (l *Laboratory) Run(ctx context.Context, specification ExperimentSpecification) (*ExecutionResult, error) {
	if l == nil || l.loadCards == nil {
		return nil, errors.New("RAG_EVALUATION_DATASET_REQUIRED: laboratory must be opened with an immutable evaluation dataset loader")
	}
	cards, err := l.loadCards(ctx, specification.Inputs.EvaluationDataset.ID)
	if err != nil {
		return nil, errors.Wrap(err, "load immutable evaluation cards")
	}
	run, err := l.Start(ctx, specification)
	if err != nil {
		return nil, err
	}
	return l.Execute(ctx, run.ID, specification, cards, ExecutionOptions{})
}
