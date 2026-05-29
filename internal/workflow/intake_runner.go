package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/engine/runner"
)

// IntakeRunner dispatches durable scraper ops into rag-eval intake services.
// Phase 1 starts with chunk_document because chunking is already idempotent for
// one document/strategy pair and has no external provider dependency.
type IntakeRunner struct{}

func (r *IntakeRunner) Kind() string { return IntakeRunnerKind }

func (r *IntakeRunner) Run(ctx context.Context, runCtx runner.RunContext) (*model.OpResult, error) {
	var input IntakeOpInput
	if err := json.Unmarshal(runCtx.Op.Input, &input); err != nil {
		return nil, fmt.Errorf("decode intake op input: %w", err)
	}
	if input.Operation == "" {
		return opErrorResult(runCtx.Op.ID, "missing_operation", "operation is required", false, nil), nil
	}

	switch input.Operation {
	case OperationEcho:
		data, err := json.Marshal(EchoOutput{
			Operation:  input.Operation,
			WorkflowID: string(runCtx.Workflow.ID),
			OpID:       string(runCtx.Op.ID),
		})
		if err != nil {
			return nil, err
		}
		return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
	case OperationChunkDocument:
		return r.runChunkDocument(ctx, runCtx, input)
	default:
		return opErrorResult(runCtx.Op.ID, "unknown_operation", fmt.Sprintf("unknown intake operation %q", input.Operation), false, map[string]any{"operation": input.Operation}), nil
	}
}

func (r *IntakeRunner) runChunkDocument(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	if input.DBPath == "" {
		return opErrorResult(runCtx.Op.ID, "missing_db_path", "db_path is required", false, nil), nil
	}
	queries, err := openQueries(input.DBPath)
	if err != nil {
		return nil, err
	}
	defer queries.Close()

	result, err := chunkservice.NewService(queries).Apply(ctx, chunkservice.ApplyRequest{
		DocumentID:   input.DocumentID,
		Strategy:     input.Strategy,
		ChunkSize:    input.ChunkSize,
		Overlap:      input.Overlap,
		StrategyName: input.StrategyName,
		Description:  input.Description,
	})
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "chunk_document_failed", err.Error(), false, map[string]any{"document_id": input.DocumentID}), nil
	}

	data, err := json.Marshal(ChunkDocumentOutput{
		DocumentID: result.DocumentID,
		StrategyID: result.StrategyID,
		ChunkCount: result.ChunkCount,
	})
	if err != nil {
		return nil, err
	}
	return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
}

func openQueries(path string) (*db.Queries, error) {
	database, err := db.OpenDB(path)
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(database); err != nil {
		_ = database.Close()
		return nil, err
	}
	return db.NewQueries(database), nil
}

func opErrorResult(opID model.OpID, code, message string, retryable bool, details map[string]any) *model.OpResult {
	var rawDetails json.RawMessage
	if details != nil {
		if b, err := json.Marshal(details); err == nil {
			rawDetails = b
		}
	}
	return &model.OpResult{
		OpID: opID,
		Error: &model.OpError{
			Code:      code,
			Message:   message,
			Retryable: retryable,
			Details:   rawDetails,
		},
	}
}
