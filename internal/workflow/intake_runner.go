package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkenrichment "github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	documentprocessing "github.com/go-go-golems/rag-evaluation-system/internal/services/documentprocessing"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/engine/runner"
)

type ProviderResolver func(ctx context.Context, input IntakeOpInput) (*embeddingservice.ResolvedProvider, error)
type DocumentProcessorResolver func(ctx context.Context, input IntakeOpInput) (documentprocessing.Provider, error)
type ChunkEnricherResolver func(ctx context.Context, input IntakeOpInput) (chunkenrichment.Provider, error)

// IntakeRunner dispatches durable scraper ops into rag-eval intake services.
type IntakeRunner struct {
	ResolveProvider          ProviderResolver
	ResolveDocumentProcessor DocumentProcessorResolver
	ResolveChunkEnricher     ChunkEnricherResolver
	IndexRoot                string
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
	case OperationPreprocessDocument:
		return r.runPreprocessDocument(ctx, runCtx, input)
	case OperationEnrichChunk:
		return r.runEnrichChunk(ctx, runCtx, input)
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
	defer func() { _ = queries.Close() }()

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

func (r *IntakeRunner) runPreprocessDocument(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	queries, opErr := openOpQueries(runCtx.Op.ID, input.DBPath)
	if opErr != nil {
		return opErr, nil
	}
	defer func() { _ = queries.Close() }()

	provider, err := r.resolveDocumentProcessor(ctx, input)
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "resolve_document_processor_failed", err.Error(), false, nil), nil
	}
	result, err := documentprocessing.NewService(queries).Process(ctx, documentprocessing.ProcessRequest{
		DocumentID:    input.DocumentID,
		ArtifactType:  input.ArtifactType,
		PromptVersion: input.PromptVersion,
		Provider:      provider,
		Force:         input.Force,
	})
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "preprocess_document_failed", err.Error(), true, map[string]any{"document_id": input.DocumentID, "artifact_type": input.ArtifactType}), nil
	}
	data, err := json.Marshal(PreprocessDocumentOutput{
		DocumentID:    result.DocumentID,
		ArtifactType:  result.ArtifactType,
		PromptVersion: result.PromptVersion,
		Provider:      result.Provider,
		Model:         result.Model,
		InputHash:     result.InputHash,
		Status:        result.Status,
		SkippedFresh:  result.SkippedFresh,
	})
	if err != nil {
		return nil, err
	}
	return &model.OpResult{OpID: runCtx.Op.ID, Data: data}, nil
}

func (r *IntakeRunner) runEnrichChunk(ctx context.Context, runCtx runner.RunContext, input IntakeOpInput) (*model.OpResult, error) {
	queries, opErr := openOpQueries(runCtx.Op.ID, input.DBPath)
	if opErr != nil {
		return opErr, nil
	}
	defer func() { _ = queries.Close() }()

	provider, err := r.resolveChunkEnricher(ctx, input)
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "resolve_chunk_enricher_failed", err.Error(), false, nil), nil
	}
	result, err := chunkenrichment.NewService(queries).Enrich(ctx, chunkenrichment.EnrichRequest{
		ChunkID:       input.ChunkID,
		StrategyID:    input.StrategyID,
		PromptVersion: input.PromptVersion,
		Provider:      provider,
		Force:         input.Force,
	})
	if err != nil {
		return opErrorResult(runCtx.Op.ID, "enrich_chunk_failed", err.Error(), true, map[string]any{"chunk_id": input.ChunkID, "strategy_id": input.StrategyID}), nil
	}
	data, err := json.Marshal(EnrichChunkOutput{
		ChunkID:       result.ChunkID,
		DocumentID:    result.DocumentID,
		StrategyID:    result.StrategyID,
		PromptVersion: result.PromptVersion,
		Provider:      result.Provider,
		Model:         result.Model,
		TextHash:      result.TextHash,
		SkippedFresh:  result.SkippedFresh,
		QualityScore:  result.QualityScore,
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
	defer func() { _ = queries.Close() }()

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
	defer func() { _ = queries.Close() }()

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

func (r *IntakeRunner) resolveDocumentProcessor(ctx context.Context, input IntakeOpInput) (documentprocessing.Provider, error) {
	if r.ResolveDocumentProcessor != nil {
		return r.ResolveDocumentProcessor(ctx, input)
	}
	provider := input.DocumentProcessingProvider
	if provider == "" {
		provider = "fake"
	}
	modelName := input.DocumentProcessingModel
	if modelName == "" {
		modelName = "fake-document-processor"
	}
	if provider == "openai-responses" {
		profile := modelName
		if profile == "" || profile == "fake-document-processor" {
			profile = input.Profile
		}
		if profile == "" {
			return nil, fmt.Errorf("document processing profile is required for openai-responses")
		}
		return documentprocessing.NewOpenAIResponsesProvider(ctx, profile, input.ProfileRegistries)
	}
	if provider != "fake" {
		return nil, fmt.Errorf("unsupported document processing provider %q", provider)
	}
	return documentprocessing.FakeProvider{ProviderName: provider, ModelName: modelName}, nil
}

func (r *IntakeRunner) resolveChunkEnricher(ctx context.Context, input IntakeOpInput) (chunkenrichment.Provider, error) {
	if r.ResolveChunkEnricher != nil {
		return r.ResolveChunkEnricher(ctx, input)
	}
	provider := input.ChunkEnrichmentProvider
	if provider == "" {
		provider = "fake"
	}
	modelName := input.ChunkEnrichmentModel
	if modelName == "" {
		modelName = "fake-chunk-enricher"
	}
	if provider != "fake" {
		return nil, fmt.Errorf("unsupported chunk enrichment provider %q; only fake is available before live provider Phase 5", provider)
	}
	return chunkenrichment.FakeProvider{ProviderName: provider, ModelName: modelName}, nil
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
