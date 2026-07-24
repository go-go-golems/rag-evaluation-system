package ragintakeworkflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3product"
)

type Config struct {
	WorkflowDatabase string
	ArtifactRoot     string
	Runtime          RuntimeConfig
	LeaseDuration    time.Duration
	PollInterval     time.Duration
	Capacities       map[string]int
}

func DefaultConfig(databasePath string) Config {
	return Config{WorkflowDatabase: "state/rag-eval-intake-v3.db", ArtifactRoot: "state/rag-eval-intake-v3-artifacts", Runtime: RuntimeConfig{DatabasePath: databasePath, IndexRoot: "data/indexes"}, LeaseDuration: 30 * time.Second, PollInterval: 100 * time.Millisecond, Capacities: map[string]int{"cpu.rag.intake": 4, "cpu.rag.intake.llm": 2, "cpu.rag.intake.embedding": 1, "cpu.rag.intake.index": 1}}
}

type Application struct {
	*workflowv3product.Application
	Config Config
}

func Open(ctx context.Context, config Config) (*Application, error) {
	if err := config.Runtime.Validate(); err != nil {
		return nil, err
	}
	product := workflowv3product.DefaultConfig()
	product.DatabasePath = config.WorkflowDatabase
	product.ArtifactRoot = config.ArtifactRoot
	product.TaskPackages = []string{PackageName}
	product.LeaseDuration = config.LeaseDuration
	product.PollInterval = config.PollInterval
	product.Capacities = config.Capacities
	app, err := workflowv3product.Open(ctx, product, NewPackage(config.Runtime))
	if err != nil {
		return nil, err
	}
	return &Application{Application: app, Config: config}, nil
}

type Selection struct {
	DocumentIDs               []string
	SourceIDs                 []string
	DocumentLimit             int
	ChunksPerDocumentToEnrich int
}

func PrepareRequest(ctx context.Context, config RuntimeConfig, request Request, selection Selection) (Request, error) {
	documents := normalizeIDs(selection.DocumentIDs)
	sources := normalizeIDs(selection.SourceIDs)
	if len(documents) == 0 {
		ids, err := selectDocumentIDs(ctx, config.DatabasePath, sources, selection.DocumentLimit)
		if err != nil {
			return Request{}, err
		}
		documents = ids
	}
	if len(documents) == 0 {
		return Request{}, fmt.Errorf("RAG_INTAKE_NO_DOCUMENTS")
	}
	request.DocumentIDs = documents
	request.SourceIDs = sources
	request.ProviderAuthorityDigest = config.ProviderAuthorityDigest
	strategy := request.Strategy
	if strategy == "" {
		strategy = "fixed"
	}
	size := request.ChunkSize
	if size == 0 {
		size = 1200
	}
	overlap := request.Overlap
	if overlap == 0 {
		overlap = 150
	}
	request.Strategy = strategy
	request.ChunkSize = size
	request.Overlap = overlap
	if !request.SkipChunkEnrichment {
		limit := selection.ChunksPerDocumentToEnrich
		if limit <= 0 {
			limit = 1
		}
		chunks, err := selectChunks(ctx, config.DatabasePath, fmt.Sprintf("%s-%d-%d", strategy, size, overlap), documents, limit)
		if err != nil {
			return Request{}, err
		}
		request.ChunkIDs = chunks
	}
	return NormalizeRequest(request, TargetDigest(config.DatabasePath))
}
func (a *Application) SubmitRequest(ctx context.Context, request Request, runID workflowv3.RunID) (workflowv3product.Submission, error) {
	normalized, err := NormalizeRequest(request, TargetDigest(a.Config.Runtime.DatabasePath))
	if err != nil {
		return workflowv3product.Submission{}, err
	}
	plan, err := CompilePlan(ctx, normalized, a.Config.Runtime)
	if err != nil {
		return workflowv3product.Submission{}, err
	}
	body, err := json.Marshal(normalized)
	if err != nil {
		return workflowv3product.Submission{}, err
	}
	ref, err := a.Artifacts.Put(ctx, RequestSchema, "application/json", body)
	if err != nil {
		return workflowv3product.Submission{}, err
	}
	return a.Submit(ctx, plan, map[string]workflowv3product.StagedInput{"request": {Schema: RequestSchema, Reference: &ref}}, "", runID)
}

func selectDocumentIDs(ctx context.Context, path string, sources []string, limit int) ([]string, error) {
	database, err := db.OpenDB(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		return nil, err
	}
	query := "SELECT id FROM documents"
	args := []any{}
	if len(sources) > 0 {
		query += " WHERE source_id IN (" + strings.TrimRight(strings.Repeat("?,", len(sources)), ",") + ")"
		for _, source := range sources {
			args = append(args, source)
		}
	}
	query += " ORDER BY word_count DESC, id"
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	rows, err := database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	ret := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ret, nil
}
func selectChunks(ctx context.Context, path, strategy string, documents []string, limit int) ([]string, error) {
	database, err := db.OpenDB(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		return nil, err
	}
	chunks, err := db.NewQueries(database).ListChunksForDocuments(strategy, documents, limit)
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(chunks))
	for i, chunk := range chunks {
		ret[i] = chunk.ID
	}
	return ret, nil
}
