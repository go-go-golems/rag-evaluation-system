package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	storecontract "github.com/go-go-golems/scraper/pkg/engine/store"
	sqlitestore "github.com/go-go-golems/scraper/pkg/engine/store/sqlite"
)

type SubmitIntakeRequest struct {
	EngineDB string
	DBPath   string

	WorkflowID string
	Name       string

	DocumentIDs   []string
	SourceIDs     []string
	DocumentLimit int

	Strategy  string
	ChunkSize int
	Overlap   int

	ProfileRegistries []string
	Profile           string
	BaseProfile       string
	EmbeddingType     string
	EmbeddingEngine   string
	Dimensions        int
	APIKey            string
	BaseURL           string
	CacheType         string
	CacheDirectory    string
	BatchSize         int
	EmbeddingLimit    int
	ForceEmbeddings   bool
	SkipEmbeddings    bool

	IndexRoot  string
	IndexID    string
	IndexLimit int
	ForceIndex bool
	SkipBM25   bool
}

type SubmitIntakeResult struct {
	WorkflowID  string   `json:"workflow_id"`
	EngineDB    string   `json:"engine_db"`
	DBPath      string   `json:"db_path"`
	DocumentIDs []string `json:"document_ids"`
	StrategyID  string   `json:"strategy_id"`
	OpIDs       []string `json:"op_ids"`
}

func SubmitIntakeWorkflow(ctx context.Context, req SubmitIntakeRequest) (*SubmitIntakeResult, error) {
	if req.EngineDB == "" {
		req.EngineDB = "state/rag-eval-workflows.db"
	}
	if req.DBPath == "" {
		req.DBPath = "data/rag-eval.db"
	}
	if req.WorkflowID == "" {
		req.WorkflowID = "intake-" + time.Now().UTC().Format("20060102T150405")
	}
	if req.Name == "" {
		req.Name = "rag-eval intake workflow"
	}
	if req.Strategy == "" {
		req.Strategy = "fixed"
	}
	if req.ChunkSize == 0 {
		req.ChunkSize = 1200
	}
	if req.Overlap == 0 {
		req.Overlap = 150
	}
	if req.BatchSize == 0 {
		req.BatchSize = 16
	}
	if req.IndexID == "" {
		req.IndexID = "bm25-" + req.WorkflowID
	}

	documentIDs := normalizeList(req.DocumentIDs)
	if len(documentIDs) == 0 {
		ids, err := selectDocumentIDs(ctx, req.DBPath, normalizeList(req.SourceIDs), req.DocumentLimit)
		if err != nil {
			return nil, err
		}
		documentIDs = ids
	}
	if len(documentIDs) == 0 {
		return nil, fmt.Errorf("no documents selected; provide --document-ids or --source-ids with matching documents")
	}

	strategyID := fmt.Sprintf("%s-%d-%d", req.Strategy, req.ChunkSize, req.Overlap)
	workflowID := model.WorkflowID(req.WorkflowID)
	workflowInput := mustRawJSON(map[string]any{
		"db_path":         req.DBPath,
		"document_ids":    documentIDs,
		"source_ids":      normalizeList(req.SourceIDs),
		"strategy_id":     strategyID,
		"embedding_limit": req.EmbeddingLimit,
		"index_id":        req.IndexID,
	})
	workflow := model.WorkflowRun{
		ID:     workflowID,
		Site:   "rag-eval",
		Name:   req.Name,
		Status: model.WorkflowStatusPending,
		Input:  workflowInput,
		Metadata: map[string]string{
			"kind":      "rag-eval-intake",
			"strategy":  strategyID,
			"ticket":    "RAGEVAL-006",
			"submitted": time.Now().UTC().Format(time.RFC3339Nano),
		},
	}

	ops := make([]model.OpSpec, 0, len(documentIDs)+2)
	chunkDeps := make([]model.Dependency, 0, len(documentIDs))
	for _, documentID := range documentIDs {
		opID := model.OpID(fmt.Sprintf("%s:chunk:%s", req.WorkflowID, documentID))
		ops = append(ops, model.OpSpec{
			ID:         opID,
			WorkflowID: workflowID,
			Site:       workflow.Site,
			Kind:       IntakeRunnerKind,
			Queue:      QueueCPU,
			DedupKey:   fmt.Sprintf("chunk:%s:%s", documentID, strategyID),
			Input: mustRawJSON(IntakeOpInput{
				Operation:  OperationChunkDocument,
				DBPath:     req.DBPath,
				DocumentID: documentID,
				Strategy:   req.Strategy,
				ChunkSize:  req.ChunkSize,
				Overlap:    req.Overlap,
			}),
		})
		chunkDeps = append(chunkDeps, model.Dependency{OpID: opID, Required: true})
	}

	if !req.SkipEmbeddings {
		embedOpID := model.OpID(req.WorkflowID + ":embed")
		ops = append(ops, model.OpSpec{
			ID:         embedOpID,
			WorkflowID: workflowID,
			Site:       workflow.Site,
			Kind:       IntakeRunnerKind,
			Queue:      QueueEmbedding,
			DedupKey:   fmt.Sprintf("embed:%s:%s", strategyID, strings.Join(normalizeList(req.SourceIDs), ",")),
			DependsOn:  chunkDeps,
			Retry: model.RetryPolicy{
				MaxAttempts:    3,
				BackoffKind:    model.BackoffKindExponential,
				InitialBackoff: time.Second,
				MaxBackoff:     30 * time.Second,
				Multiplier:     2,
			},
			Input: mustRawJSON(IntakeOpInput{
				Operation:         OperationComputeEmbeddings,
				DBPath:            req.DBPath,
				StrategyID:        strategyID,
				SourceIDs:         normalizeList(req.SourceIDs),
				DocumentIDs:       documentIDs,
				ProfileRegistries: req.ProfileRegistries,
				Profile:           req.Profile,
				BaseProfile:       req.BaseProfile,
				EmbeddingType:     req.EmbeddingType,
				EmbeddingEngine:   req.EmbeddingEngine,
				Dimensions:        req.Dimensions,
				APIKey:            req.APIKey,
				BaseURL:           req.BaseURL,
				CacheType:         req.CacheType,
				CacheDirectory:    req.CacheDirectory,
				BatchSize:         req.BatchSize,
				Limit:             req.EmbeddingLimit,
				Force:             req.ForceEmbeddings,
			}),
		})
	}

	if !req.SkipBM25 {
		bm25OpID := model.OpID(req.WorkflowID + ":bm25")
		ops = append(ops, model.OpSpec{
			ID:         bm25OpID,
			WorkflowID: workflowID,
			Site:       workflow.Site,
			Kind:       IntakeRunnerKind,
			Queue:      QueueIndex,
			DedupKey:   fmt.Sprintf("bm25:%s:%s", strategyID, req.IndexID),
			DependsOn:  chunkDeps,
			Input: mustRawJSON(IntakeOpInput{
				Operation:   OperationBuildBM25,
				DBPath:      req.DBPath,
				StrategyID:  strategyID,
				SourceIDs:   normalizeList(req.SourceIDs),
				DocumentIDs: documentIDs,
				IndexRoot:   req.IndexRoot,
				IndexID:     req.IndexID,
				Limit:       req.IndexLimit,
				Force:       req.ForceIndex,
			}),
		})
	}

	store, err := sqlitestore.Open(ctx, req.EngineDB)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	if err := store.CreateWorkflow(ctx, storecontract.CreateWorkflowParams{Workflow: workflow, Initial: ops}); err != nil {
		return nil, err
	}

	opIDs := make([]string, len(ops))
	for i, op := range ops {
		opIDs[i] = string(op.ID)
	}
	return &SubmitIntakeResult{WorkflowID: string(workflowID), EngineDB: req.EngineDB, DBPath: req.DBPath, DocumentIDs: documentIDs, StrategyID: strategyID, OpIDs: opIDs}, nil
}

func selectDocumentIDs(ctx context.Context, dbPath string, sourceIDs []string, limit int) ([]string, error) {
	database, err := db.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}
	defer database.Close()
	if err := db.Migrate(database); err != nil {
		return nil, err
	}
	query := `SELECT id FROM documents`
	args := []any{}
	if len(sourceIDs) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(sourceIDs)), ",")
		query += ` WHERE source_id IN (` + placeholders + `)`
		for _, sourceID := range sourceIDs {
			args = append(args, sourceID)
		}
	}
	query += ` ORDER BY word_count DESC, id`
	if limit > 0 {
		query += ` LIMIT ?`
		args = append(args, limit)
	}
	rows, err := database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ids, nil
}

func normalizeList(values []string) []string {
	seen := map[string]bool{}
	ret := []string{}
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" || seen[trimmed] {
				continue
			}
			seen[trimmed] = true
			ret = append(ret, trimmed)
		}
	}
	return ret
}

func mustRawJSON(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
