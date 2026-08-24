package ragintakeworkflow

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dop251/goja"
	geppettoembeddings "github.com/go-go-golems/geppetto/pkg/embeddings"
	gggengine "github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/documentprocessing"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

type ProviderResolver func(context.Context, Request) (*embeddingservice.ResolvedProvider, error)
type DocumentProcessorResolver func(context.Context, Request) (documentprocessing.Provider, error)
type ChunkEnricherResolver func(context.Context, Request) (chunkenrichment.Provider, error)

type RuntimeConfig struct {
	DatabasePath             string
	IndexRoot                string
	APIKey                   string
	BaseURL                  string
	CacheDirectory           string
	ProviderAuthorityDigest  string
	MaxProviderOperations    int
	ProviderFinishTimeout    time.Duration
	ResolveProvider          ProviderResolver
	ResolveDocumentProcessor DocumentProcessorResolver
	ResolveChunkEnricher     ChunkEnricherResolver
}

func (c RuntimeConfig) Validate() error {
	if strings.TrimSpace(c.DatabasePath) == "" {
		return fmt.Errorf("RAG_INTAKE_DATABASE")
	}
	if c.ProviderAuthorityDigest != "" {
		_, err := ragworkflowops.NewDescriptors(c.ProviderAuthorityDigest, maxProviderOperations(c))
		return err
	}
	return nil
}
func maxProviderOperations(c RuntimeConfig) int {
	if c.MaxProviderOperations > 0 {
		return c.MaxProviderOperations
	}
	return 10_000
}
func providerFinishTimeout(c RuntimeConfig) time.Duration {
	if c.ProviderFinishTimeout > 0 {
		return c.ProviderFinishTimeout
	}
	return 5 * time.Second
}

func TaskModuleFactory(config RuntimeConfig) workflowv3runtime.TaskModuleFactory {
	descriptors := []workflowv3.ExternalOperationDescriptor{}
	if config.ProviderAuthorityDigest != "" {
		descriptors, _ = ragworkflowops.NewDescriptors(config.ProviderAuthorityDigest, maxProviderOperations(config))
	}
	return workflowv3runtime.TaskModuleFactory{Alias: ModuleAlias, Validate: config.Validate, Operations: descriptors, Build: func(moduleContext workflowv3runtime.TaskModuleContext) (gggengine.RuntimeModuleRegistrar, error) {
		runtime := &taskRuntime{context: moduleContext, config: config}
		loader := func(vm *goja.Runtime, moduleObject *goja.Object) {
			exports := moduleObject.Get("exports").ToObject(vm)
			if err := exports.Set("run", func(goja.FunctionCall) goja.Value {
				value, schema, err := runtime.run()
				result := map[string]any{"ok": err == nil, "value": value, "schema": schema}
				if runtime.providerUsage != nil {
					result["providerUsage"] = runtime.providerUsage
				}
				if err != nil {
					result["failure"] = taskFailure(err)
				}
				return vm.ToValue(result)
			}); err != nil {
				panic(vm.NewGoError(err))
			}
		}
		return gggengine.NativeModuleRegistrar{ModuleID: "rag-intake-runtime-v1", ModuleName: ModuleAlias, Loader: loader}, nil
	}}
}

type intakeError struct {
	code      string
	retryable bool
	cause     error
}

func (e *intakeError) Error() string { return e.code + ": " + e.cause.Error() }
func (e *intakeError) Unwrap() error { return e.cause }
func fail(code string, retryable bool, err error) error {
	return &intakeError{code: code, retryable: retryable, cause: err}
}
func taskFailure(err error) map[string]any {
	failure := map[string]any{"class": "execution", "code": "RAG_INTAKE_TASK_FAILED", "retryable": false, "message": "RAG intake task failed"}
	if typed, ok := err.(*intakeError); ok {
		failure["code"] = typed.code
		failure["retryable"] = typed.retryable
	}
	return failure
}

type taskRuntime struct {
	context       workflowv3runtime.TaskModuleContext
	config        RuntimeConfig
	providerUsage map[string]int64
}

func (r *taskRuntime) run() (any, string, error) {
	r.providerUsage = nil
	request, err := r.request()
	if err != nil {
		return nil, "", err
	}
	kind := r.context.Request.Task.Spec.Identity.Kind
	switch kind {
	case TaskPreprocess.Kind:
		return r.preprocess(request)
	case TaskChunk.Kind:
		return r.chunk(request)
	case TaskEnrich.Kind:
		return r.enrich(request)
	case TaskEmbed.Kind:
		return r.embed(request)
	case TaskBM25.Kind:
		return r.bm25(request)
	case TaskPublish.Kind:
		return Summary{SchemaVersion: SummarySchema, TargetDigest: request.TargetDigest, DocumentCount: len(request.DocumentIDs), ChunkTargetCount: len(request.ChunkIDs)}, SummarySchema, nil
	default:
		return nil, "", fail("RAG_INTAKE_TASK_KIND", false, fmt.Errorf("unsupported task"))
	}
}
func (r *taskRuntime) request() (Request, error) {
	ref, ok := r.context.Request.Inputs["request"]
	if !ok {
		return Request{}, fail("RAG_INTAKE_REQUEST_MISSING", false, fmt.Errorf("request input required"))
	}
	reader, err := r.context.Request.Artifacts.Open(r.context.Context, ref)
	if err != nil {
		return Request{}, fail("RAG_INTAKE_REQUEST_OPEN", true, err)
	}
	defer func() { _ = reader.Close() }()
	request, err := DecodeRequest(reader, TargetDigest(r.config.DatabasePath))
	if err == nil && request.ProviderAuthorityDigest != r.config.ProviderAuthorityDigest {
		err = fmt.Errorf("RAG_INTAKE_PROVIDER_AUTHORITY")
	}
	if err != nil {
		return Request{}, fail("RAG_INTAKE_REQUEST_INVALID", false, err)
	}
	return request, nil
}
func nodeOrdinal(key workflowv3.NodeKey) (int, error) {
	value := string(key)
	index := strings.LastIndexByte(value, '-')
	if index < 0 {
		return 0, fmt.Errorf("node ordinal")
	}
	ordinal, err := strconv.Atoi(value[index+1:])
	if err != nil {
		return 0, fmt.Errorf("node ordinal")
	}
	return ordinal, nil
}
func (r *taskRuntime) queries() (*db.Queries, error) {
	database, err := db.OpenDB(r.config.DatabasePath)
	if err != nil {
		return nil, err
	}
	if err := db.Migrate(database); err != nil {
		_ = database.Close()
		return nil, err
	}
	return db.NewQueries(database), nil
}
func (r *taskRuntime) preprocess(request Request) (any, string, error) {
	ordinal, err := nodeOrdinal(r.context.Request.NodeKey)
	if err != nil || ordinal >= len(request.DocumentIDs) {
		return nil, "", fail("RAG_INTAKE_NODE", false, fmt.Errorf("preprocess ordinal"))
	}
	queries, err := r.queries()
	if err != nil {
		return nil, "", fail("RAG_INTAKE_DATABASE_OPEN", true, err)
	}
	defer func() { _ = queries.Close() }()
	provider, err := r.documentProcessor(request)
	if err != nil {
		return nil, "", fail("RAG_INTAKE_PROVIDER", false, err)
	}
	result, err := documentprocessing.NewService(queries).Process(r.context.Context, documentprocessing.ProcessRequest{DocumentID: request.DocumentIDs[ordinal], ArtifactType: request.PreprocessArtifactType, PromptVersion: request.PreprocessPromptVersion, Provider: provider, Force: request.ForcePreprocessing})
	if err != nil {
		return nil, "", fail("RAG_INTAKE_PREPROCESS", true, err)
	}
	return Result{SchemaVersion: ResultSchema, Operation: "preprocess_document", DocumentID: result.DocumentID, Identity: map[string]string{"provider": result.Provider, "model": result.Model, "inputHash": result.InputHash}, SkippedFresh: result.SkippedFresh}, ResultSchema, nil
}
func (r *taskRuntime) chunk(request Request) (any, string, error) {
	ordinal, err := nodeOrdinal(r.context.Request.NodeKey)
	if err != nil || ordinal >= len(request.DocumentIDs) {
		return nil, "", fail("RAG_INTAKE_NODE", false, fmt.Errorf("chunk ordinal"))
	}
	queries, err := r.queries()
	if err != nil {
		return nil, "", fail("RAG_INTAKE_DATABASE_OPEN", true, err)
	}
	defer func() { _ = queries.Close() }()
	result, err := chunkservice.NewService(queries).Apply(r.context.Context, chunkservice.ApplyRequest{DocumentID: request.DocumentIDs[ordinal], Strategy: request.Strategy, ChunkSize: request.ChunkSize, Overlap: request.Overlap})
	if err != nil {
		return nil, "", fail("RAG_INTAKE_CHUNK", false, err)
	}
	return Result{SchemaVersion: ResultSchema, Operation: "chunk_document", DocumentID: result.DocumentID, StrategyID: result.StrategyID, Counts: map[string]int{"chunks": result.ChunkCount}}, ResultSchema, nil
}
func (r *taskRuntime) enrich(request Request) (any, string, error) {
	ordinal, err := nodeOrdinal(r.context.Request.NodeKey)
	if err != nil || ordinal >= len(request.ChunkIDs) {
		return nil, "", fail("RAG_INTAKE_NODE", false, fmt.Errorf("enrich ordinal"))
	}
	queries, err := r.queries()
	if err != nil {
		return nil, "", fail("RAG_INTAKE_DATABASE_OPEN", true, err)
	}
	defer func() { _ = queries.Close() }()
	provider, err := r.chunkEnricher(request)
	if err != nil {
		return nil, "", fail("RAG_INTAKE_PROVIDER", false, err)
	}
	strategyID := fmt.Sprintf("%s-%d-%d", request.Strategy, request.ChunkSize, request.Overlap)
	result, err := chunkenrichment.NewService(queries).Enrich(r.context.Context, chunkenrichment.EnrichRequest{ChunkID: request.ChunkIDs[ordinal], StrategyID: strategyID, PromptVersion: request.ChunkEnrichmentPrompt, Provider: provider, Force: request.ForceChunkEnrichment})
	if err != nil {
		return nil, "", fail("RAG_INTAKE_ENRICH", true, err)
	}
	return Result{SchemaVersion: ResultSchema, Operation: "enrich_chunk", ChunkID: result.ChunkID, DocumentID: result.DocumentID, StrategyID: result.StrategyID, Identity: map[string]string{"provider": result.Provider, "model": result.Model, "textHash": result.TextHash}, SkippedFresh: result.SkippedFresh}, ResultSchema, nil
}
func (r *taskRuntime) embed(request Request) (any, string, error) {
	queries, err := r.queries()
	if err != nil {
		return nil, "", fail("RAG_INTAKE_DATABASE_OPEN", true, err)
	}
	defer func() { _ = queries.Close() }()
	provider, err := r.embeddingProvider(request)
	if err != nil {
		return nil, "", fail("RAG_INTAKE_PROVIDER", false, err)
	}
	if provider.ProviderType != "fake" {
		r.providerUsage = map[string]int64{"requests": 1}
		if r.config.ProviderAuthorityDigest == "" {
			return nil, "", fail("RAG_INTAKE_PROVIDER_AUTHORITY", false, fmt.Errorf("provider authority required"))
		}
		decorator, decorateErr := ragworkflowops.NewDecorator(r.context.ExternalOperations, ragworkflowops.Policy{AuthorityDigest: r.config.ProviderAuthorityDigest, MaxPerAttempt: maxProviderOperations(r.config), FinishTimeout: providerFinishTimeout(r.config), Reservations: map[string][]workflowv3.ExternalOperationCounter{ragworkflowops.EmbedOperation: {{Name: "requests", Units: 1}}}}, nil)
		if decorateErr != nil {
			return nil, "", fail("RAG_INTAKE_PROVIDER_AUTHORITY", false, decorateErr)
		}
		embedder, decorateErr := decorator.Embedder(providerEmbedder{provider: provider.Provider})
		if decorateErr != nil {
			return nil, "", fail("RAG_INTAKE_PROVIDER_AUTHORITY", false, decorateErr)
		}
		provider.Provider = custodiedEmbeddingProvider{embedder: embedder, model: provider.Provider.GetModel()}
	}
	if provider.Close != nil {
		defer func() { _ = provider.Close() }()
	}
	strategyID := fmt.Sprintf("%s-%d-%d", request.Strategy, request.ChunkSize, request.Overlap)
	result, err := embeddingservice.NewService(queries).Compute(r.context.Context, embeddingservice.ComputeRequest{StrategyID: strategyID, SourceIDs: request.SourceIDs, DocumentIDs: request.DocumentIDs, Provider: provider.Provider, ProviderType: provider.ProviderType, BatchSize: request.BatchSize, Limit: request.EmbeddingLimit, Force: request.ForceEmbeddings})
	if err != nil {
		return nil, "", fail("RAG_INTAKE_EMBED", true, err)
	}
	return Result{SchemaVersion: ResultSchema, Operation: "compute_embeddings", StrategyID: result.StrategyID, Counts: map[string]int{"considered": result.Considered, "computed": result.Computed, "skippedFresh": result.SkippedFresh}, Identity: map[string]string{"providerType": result.ProviderType, "model": result.Model, "dimensions": strconv.Itoa(result.Dimensions)}}, ResultSchema, nil
}
func (r *taskRuntime) bm25(request Request) (any, string, error) {
	queries, err := r.queries()
	if err != nil {
		return nil, "", fail("RAG_INTAKE_DATABASE_OPEN", true, err)
	}
	defer func() { _ = queries.Close() }()
	strategyID := fmt.Sprintf("%s-%d-%d", request.Strategy, request.ChunkSize, request.Overlap)
	result, err := searchservice.NewService(queries, r.config.IndexRoot).BuildBM25(r.context.Context, searchservice.BuildIndexRequest{IndexID: request.IndexID, StrategyID: strategyID, SourceIDs: request.SourceIDs, DocumentIDs: request.DocumentIDs, Force: request.ForceBM25, Limit: request.IndexLimit})
	if err != nil {
		return nil, "", fail("RAG_INTAKE_BM25", false, err)
	}
	return Result{SchemaVersion: ResultSchema, Operation: "build_bm25", IndexID: result.IndexID, StrategyID: result.StrategyID, Counts: map[string]int{"chunks": result.ChunkCount, "documents": result.DocumentCount}, Identity: map[string]string{"indexPath": result.IndexPath}}, ResultSchema, nil
}

type providerEmbedder struct{ provider geppettoembeddings.Provider }

func (p providerEmbedder) Embed(ctx context.Context, _ string, texts []string) ([][]float64, ragoperators.Usage, error) {
	values, err := p.provider.GenerateBatchEmbeddings(ctx, texts)
	ret := make([][]float64, len(values))
	for i, row := range values {
		ret[i] = make([]float64, len(row))
		for j, value := range row {
			ret[i][j] = float64(value)
		}
	}
	return ret, ragoperators.Usage{}, err
}

type custodiedEmbeddingProvider struct {
	embedder ragoperators.Embedder
	model    geppettoembeddings.EmbeddingModel
}

func (p custodiedEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	values, err := p.GenerateBatchEmbeddings(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return values[0], nil
}
func (p custodiedEmbeddingProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	values, _, err := p.embedder.Embed(ctx, p.model.Name, texts)
	ret := make([][]float32, len(values))
	for i, row := range values {
		ret[i] = make([]float32, len(row))
		for j, value := range row {
			ret[i][j] = float32(value)
		}
	}
	return ret, err
}
func (p custodiedEmbeddingProvider) GetModel() geppettoembeddings.EmbeddingModel { return p.model }

func (r *taskRuntime) embeddingProvider(request Request) (*embeddingservice.ResolvedProvider, error) {
	if r.config.ResolveProvider != nil {
		return r.config.ResolveProvider(r.context.Context, request)
	}
	return embeddingservice.ResolveProvider(r.context.Context, embeddingservice.ProviderConfig{ProfileRegistries: request.ProfileRegistries, Profile: request.Profile, BaseProfile: request.BaseProfile, Type: request.EmbeddingType, Engine: request.EmbeddingEngine, Dimensions: request.Dimensions, APIKey: r.config.APIKey, BaseURL: r.config.BaseURL, CacheType: request.CacheType, CacheDirectory: r.config.CacheDirectory})
}
func (r *taskRuntime) documentProcessor(request Request) (documentprocessing.Provider, error) {
	if r.config.ResolveDocumentProcessor != nil {
		return r.config.ResolveDocumentProcessor(r.context.Context, request)
	}
	if request.PreprocessProvider == "fake" {
		return documentprocessing.FakeProvider{ProviderName: "fake", ModelName: request.PreprocessModel}, nil
	}
	if request.PreprocessProvider == "openai-responses" {
		profile := request.PreprocessModel
		if profile == "" || profile == "fake-document-processor" {
			profile = request.Profile
		}
		if profile == "" {
			return nil, fmt.Errorf("profile required")
		}
		return documentprocessing.NewOpenAIResponsesProvider(r.context.Context, profile, request.ProfileRegistries)
	}
	return nil, fmt.Errorf("unsupported document processor")
}
func (r *taskRuntime) chunkEnricher(request Request) (chunkenrichment.Provider, error) {
	if r.config.ResolveChunkEnricher != nil {
		return r.config.ResolveChunkEnricher(r.context.Context, request)
	}
	if request.ChunkEnrichmentProvider != "fake" {
		return nil, fmt.Errorf("unsupported chunk enricher")
	}
	return chunkenrichment.FakeProvider{ProviderName: "fake", ModelName: request.ChunkEnrichmentModel}, nil
}
