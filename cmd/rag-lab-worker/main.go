package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/pkg/raglab"
	_ "github.com/mattn/go-sqlite3"
)

const protocolVersion = "researchctl-rag-runner-stdio/v1"

type requestEnvelope struct {
	ProtocolVersion string  `json:"protocolVersion"`
	Request         request `json:"request"`
}

type request struct {
	Specification raglab.PrototypeSpecification `json:"specification"`
	Inputs        resolvedInputs                `json:"inputs"`
	DatasetSplit  string                        `json:"datasetSplit"`
}

type resolvedInputs struct {
	ByRole map[string]resolvedInput `json:"byRole"`
}

type resolvedInput struct {
	Role     string            `json:"role"`
	Kind     string            `json:"kind"`
	ID       string            `json:"id,omitempty"`
	Manifest json.RawMessage   `json:"manifest,omitempty"`
	Path     string            `json:"path,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type frame struct {
	Type     string                      `json:"type"`
	Event    *raglab.DomainEvent         `json:"event,omitempty"`
	Trace    *raglab.PrototypeQueryTrace `json:"trace,omitempty"`
	Metric   *raglab.DomainMetric        `json:"metric,omitempty"`
	Artifact *raglab.DomainArtifact      `json:"artifact,omitempty"`
	Error    string                      `json:"error,omitempty"`
}

type protocolObserver struct{ encoder *json.Encoder }

func (o *protocolObserver) Event(_ context.Context, value raglab.DomainEvent) error {
	return o.encoder.Encode(frame{Type: "event", Event: &value})
}
func (o *protocolObserver) QueryTrace(_ context.Context, value raglab.PrototypeQueryTrace) error {
	return o.encoder.Encode(frame{Type: "trace", Trace: &value})
}
func (o *protocolObserver) Metric(_ context.Context, value raglab.DomainMetric) error {
	return o.encoder.Encode(frame{Type: "metric", Metric: &value})
}
func (o *protocolObserver) Artifact(_ context.Context, value raglab.DomainArtifact) error {
	return o.encoder.Encode(frame{Type: "artifact", Artifact: &value})
}

func main() {
	databasePath := flag.String("db", "", "Immutable TTC rag-eval SQLite database")
	artifactBase := flag.String("artifact-base", "", "Base directory for relative catalog artifact paths (defaults to the database parent repository)")
	ollamaBaseURL := flag.String("ollama-base-url", "http://127.0.0.1:11435", "Ollama embedding endpoint")
	cacheDirectory := flag.String("embedding-cache", "state/embedding-cache", "Query embedding cache")
	rerankerBaseURL := flag.String("reranker-base-url", "", "Optional llama.cpp-compatible reranker endpoint")
	flag.Parse()
	if *databasePath == "" {
		fatal("--db is required")
	}
	absoluteDatabase, err := filepath.Abs(*databasePath)
	if err != nil {
		fatal("resolve database: %v", err)
	}
	base := *artifactBase
	if base == "" {
		base = filepath.Dir(filepath.Dir(absoluteDatabase))
	}
	if err := os.Chdir(base); err != nil {
		fatal("change artifact base: %v", err)
	}
	ctx := context.Background()
	var envelope requestEnvelope
	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&envelope); err != nil {
		fatal("decode request: %v", err)
	}
	if envelope.ProtocolVersion != protocolVersion {
		fatal("unsupported protocol %q", envelope.ProtocolVersion)
	}
	database, err := openReadOnly(ctx, absoluteDatabase)
	if err != nil {
		fatal("open database: %v", err)
	}
	defer func() { _ = database.Close() }()
	queries := db.NewQueries(database)
	prototype, datasetID, err := toPrototype(envelope.Request)
	if err != nil {
		fatal("map specification: %v", err)
	}
	cards, err := raglab.LoadEvaluationCards(ctx, queries, datasetID)
	if err != nil {
		fatal("load evaluation cards: %v", err)
	}
	options := raglab.ExecutionOptions{}
	var closeEmbedding func() error
	if hasBackend(prototype, raglab.VectorBackend) {
		provider, resolveErr := embedding.ResolveProvider(ctx, embedding.ProviderConfig{
			Type: "ollama", Engine: "nomic-embed-text", Dimensions: 768,
			BaseURL: *ollamaBaseURL, CacheType: "file", CacheDirectory: *cacheDirectory,
		})
		if resolveErr != nil {
			fatal("resolve embedding provider: %v", resolveErr)
		}
		options.Embedder = provider.Provider
		closeEmbedding = provider.Close
	}
	if closeEmbedding != nil {
		defer func() { _ = closeEmbedding() }()
	}
	if prototype.Retrieval.Reranking != nil {
		if *rerankerBaseURL == "" {
			fatal("--reranker-base-url is required by the specification")
		}
		reranker, rerankErr := raglab.NewLlamaCPPReranker(raglab.LlamaCPPRerankerOptions{BaseURL: *rerankerBaseURL, Model: prototype.Retrieval.Reranking.Model})
		if rerankErr != nil {
			fatal("configure reranker: %v", rerankErr)
		}
		options.Reranker = reranker
	}
	executor := raglab.NewObservationExecutor(raglab.NewSQLiteChannelRetriever(queries))
	observer := &protocolObserver{encoder: json.NewEncoder(os.Stdout)}
	if err := executor.Execute(ctx, raglab.ObservationExecutionRequest{
		Specification: prototype, DatasetSplit: envelope.Request.DatasetSplit,
		Cards: cards, Options: options,
	}, observer); err != nil {
		_ = observer.encoder.Encode(frame{Type: "error", Error: err.Error()})
		os.Exit(1)
	}
}

func openReadOnly(ctx context.Context, path string) (*sql.DB, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("mode", "ro")
	query.Set("_query_only", "1")
	database, err := sql.Open("sqlite3", "file:"+filepath.ToSlash(absolute)+"?"+query.Encode())
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(1)
	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}

func toPrototype(input request) (raglab.ExperimentSpecification, string, error) {
	lookup := func(role string) (resolvedInput, error) {
		value, ok := input.Inputs.ByRole[role]
		if !ok || value.ID == "" {
			return resolvedInput{}, fmt.Errorf("resolved input %q with ID is required", role)
		}
		return value, nil
	}
	corpus, err := lookup("corpus-snapshot")
	if err != nil {
		return raglab.ExperimentSpecification{}, "", err
	}
	chunks, err := lookup("chunk-set")
	if err != nil {
		return raglab.ExperimentSpecification{}, "", err
	}
	evaluation, err := lookup("evaluation-dataset")
	if err != nil {
		return raglab.ExperimentSpecification{}, "", err
	}
	grade, err := raglab.Grade(input.Specification.Metrics.RelevanceAt)
	if err != nil {
		return raglab.ExperimentSpecification{}, "", err
	}
	result := raglab.ExperimentSpecification{
		SchemaVersion: "rag-eval-experiment-spec/v1", Name: input.Specification.Name,
		Inputs: raglab.InputSpec{
			CorpusSnapshot: raglab.CorpusSnapshot(corpus.ID), ChunkSet: raglab.ChunkSet(chunks.ID),
			EvaluationDataset: raglab.EvaluationDataset(evaluation.ID),
		},
		Retrieval: raglab.RetrievalPlan{Collapse: raglab.CollapseScope(input.Specification.Retrieval.Collapse), Results: input.Specification.Retrieval.Results},
		Metrics:   raglab.MetricsPlan{RelevanceAt: &grade, RecallAt: input.Specification.Metrics.RecallAt, PrecisionAt: input.Specification.Metrics.PrecisionAt, MRR: input.Specification.Metrics.MRR},
	}
	if len(input.Specification.Metrics.NDCGAt) > 0 {
		result.Metrics.NDCGAt = input.Specification.Metrics.NDCGAt[0]
	}
	for _, representation := range input.Specification.Representations {
		kind := raglab.RepresentationKind(representation.Kind)
		if representation.Kind == "questions" {
			kind = raglab.QuestionRepresentation
		}
		result.Inputs.Representations = append(result.Inputs.Representations, raglab.RepresentationSpec{Name: representation.Name, Kind: kind})
	}
	if value, ok := input.Inputs.ByRole["bm25-index"]; ok {
		ref := raglab.BM25Index(value.ID)
		result.Inputs.BM25Index = &ref
	}
	if value, ok := input.Inputs.ByRole["embedding-set"]; ok {
		ref := raglab.EmbeddingSet(value.ID)
		result.Inputs.EmbeddingSet = &ref
	}
	for _, channel := range input.Specification.Retrieval.Channels {
		result.Retrieval.Channels = append(result.Retrieval.Channels, raglab.ChannelSpec{Name: channel.Name, Backend: raglab.RetrievalBackend(channel.Backend), Representation: channel.Representation, TopK: channel.TopK, Filter: toFilter(channel.Filter)})
	}
	result.Retrieval.Filter = toFilter(input.Specification.Retrieval.Filter)
	if fusion := input.Specification.Retrieval.Fusion; fusion != nil {
		result.Retrieval.Fusion = &raglab.FusionSpec{Kind: fusion.Kind, RankConstant: fusion.RankConstant, Weights: fusion.Weights}
	}
	if reranking := input.Specification.Retrieval.Reranking; reranking != nil {
		result.Retrieval.Reranking = &raglab.RerankingSpec{Kind: raglab.RerankingKind(reranking.Kind), Model: reranking.Model, CandidateCount: reranking.CandidateCount, Results: reranking.Results}
	}
	return result, evaluation.ID, nil
}

func toFilter(input raglab.PrototypeFilterSpec) raglab.FilterSpec {
	return raglab.FilterSpec(input)
}

func hasBackend(specification raglab.ExperimentSpecification, backend raglab.RetrievalBackend) bool {
	for _, channel := range specification.Retrieval.Channels {
		if channel.Backend == backend {
			return true
		}
	}
	return false
}

func fatal(format string, values ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "rag-lab-worker: "+format+"\n", values...)
	os.Exit(2)
}
