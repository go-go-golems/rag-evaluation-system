package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/engine/runner"
)

type ProviderResolver func(ctx context.Context, input IntakeOpInput) (*embeddingservice.ResolvedProvider, error)

// IntakeRunner dispatches durable scraper ops into rag-eval intake services.
type IntakeRunner struct {
	ResolveProvider ProviderResolver
	IndexRoot       string
}

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
	case OperationComputeEmbeddings:
		return r.runComputeEmbeddings(ctx, runCtx, input)
	case OperationBuildBM25:
		return r.runBuildBM25(ctx, runCtx, input)
	default:
		return opErrorResult(runCtx.Op.ID, "unknown_operation", fmt.Sprintf("unknown intake operation %q", input.Operation), false, map[string]any{"operation": input.Operation}), nil
	}
}

func (r *IntakeRunner) runChunkDocument(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	queries, opErr := openOpQueries(runCtx.Op.ID, input.DBPath)
	if opErr != nil {
		return opErr, nil
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

func (r *IntakeRunner) runComputeEmbeddings(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	queries, opErr := openOpQueries(runCtx.Op.ID, input.DBPath)
	if opErr != nil {
		return opErr, nil
	}
	defer queries.Close()

	resolved, err := r.resolveProvider(ctx, input)
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "resolve_embedding_provider_failed", err.Error(), false, nil), nil
	}
	if resolved.Close != nil {
		defer func() { _ = resolved.Close() }()
	}

	result, err := embeddingservice.NewService(queries).Compute(ctx, embeddingservice.ComputeRequest{
		StrategyID:   input.StrategyID,
		SourceIDs:    input.SourceIDs,
		DocumentIDs:  input.DocumentIDs,
		Provider:     resolved.Provider,
		ProviderType: resolved.ProviderType,
		BatchSize:    input.BatchSize,
		Limit:        input.Limit,
		Force:        input.Force,
	})
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "compute_embeddings_failed", err.Error(), true, map[string]any{"strategy_id": input.StrategyID}), nil
	}

	data, err := json.Marshal(ComputeEmbeddingsOutput{
		StrategyID:   result.StrategyID,
		SourceIDs:    result.SourceIDs,
		DocumentIDs:  result.DocumentIDs,
		ProviderType: result.ProviderType,
		Model:        result.Model,
		Dimensions:   result.Dimensions,
		Considered:   result.Considered,
		Computed:     result.Computed,
		SkippedFresh: result.SkippedFresh,
	})
	if err != nil {
		return nil, err
	}
	return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
}

func (r *IntakeRunner) runBuildBM25(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	queries, opErr := openOpQueries(runCtx.Op.ID, input.DBPath)
	if opErr != nil {
		return opErr, nil
	}
	defer queries.Close()

	indexRoot := input.IndexRoot
	if indexRoot == "" {
		indexRoot = r.IndexRoot
	}
	result, err := searchservice.NewService(queries, indexRoot).BuildBM25(ctx, searchservice.BuildIndexRequest{
		IndexID:     input.IndexID,
		StrategyID:  input.StrategyID,
		SourceIDs:   input.SourceIDs,
		DocumentIDs: input.DocumentIDs,
		Force:       input.Force,
		Limit:       input.Limit,
	})
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "build_bm25_failed", err.Error(), false, map[string]any{"strategy_id": input.StrategyID, "index_id": input.IndexID}), nil
	}

	data, err := json.Marshal(BuildBM25Output{
		IndexID:       result.IndexID,
		StrategyID:    result.StrategyID,
		SourceIDs:     result.SourceIDs,
		DocumentIDs:   result.DocumentIDs,
		IndexPath:     result.IndexPath,
		ChunkCount:    result.ChunkCount,
		DocumentCount: result.DocumentCount,
		Rebuilt:       result.Rebuilt,
	})
	if err != nil {
		return nil, err
	}
	return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
}

func (r *IntakeRunner) resolveProvider(ctx context.Context, input IntakeOpInput) (*embeddingservice.ResolvedProvider, error) {
	if r.ResolveProvider != nil {
		return r.ResolveProvider(ctx, input)
	}
	return embeddingservice.ResolveProvider(ctx, embeddingservice.ProviderConfig{
		ProfileRegistries: input.ProfileRegistries,
		Profile:           input.Profile,
		BaseProfile:       input.BaseProfile,
		Type:              input.EmbeddingType,
		Engine:            input.EmbeddingEngine,
		Dimensions:        input.Dimensions,
		APIKey:            input.APIKey,
		BaseURL:           input.BaseURL,
		CacheType:         input.CacheType,
		CacheDirectory:    input.CacheDirectory,
	})
}

func openOpQueries(opID model.OpID, path string) (*db.Queries, *model.OpResult) {
	if path == "" {
		return nil, opErrorResult(opID, "missing_db_path", "db_path is required", false, nil)
	}
	queries, err := openQueries(path)
	if err != nil {
		return nil, opErrorResult(opID, "open_db_failed", err.Error(), true, map[string]any{"db_path": path})
	}
	return queries, nil
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
