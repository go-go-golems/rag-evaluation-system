//go:build ignore

// Runs fresh executor-backed vector-only and weighted-RRF TTC observations.
//
// The script deliberately injects a direct Geppetto/Ollama provider for query
// embeddings. It does not infer a server or model from the immutable embedding
// set selected by either plan. Run it from the repository root after starting
// the documented mimimi SSH tunnel:
//
// GOWORK=off go run ./ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/03-run-ttc-vector-and-weighted-rrf-experiments.go --base-url http://127.0.0.1:11435
package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/go-go-golems/rag-evaluation-system/pkg/raglab"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	snapshotID     = "sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409"
	chunkSetID     = "sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392"
	bm25ID         = "sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691"
	embeddingSetID = "sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0"
	datasetID      = "candidate:ttc-baseline-v1"
)

type observation struct {
	Name        string         `json:"name"`
	Fingerprint string         `json:"fingerprint"`
	RunID       string         `json:"run_id"`
	Metrics     map[string]any `json:"metrics"`
	Timing      map[string]any `json:"timing"`
}

type report struct {
	Provider map[string]any `json:"provider"`
	Cost     map[string]any `json:"cost"`
	Storage  map[string]any `json:"storage"`
	Runs     []observation  `json:"runs"`
}

func main() {
	databasePath := flag.String("db", "data/rag-eval.db", "rag-eval SQLite database")
	baseURL := flag.String("base-url", "http://127.0.0.1:11435", "explicit Ollama base URL, normally the mimimi SSH tunnel")
	engine := flag.String("engine", "nomic-embed-text", "Ollama embedding model")
	vectorWeight := flag.Float64("vector-weight", 2, "positive semantic-channel RRF weight")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()

	parsed, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(parsed)
	if *vectorWeight <= 0 {
		log.Fatal().Float64("vector_weight", *vectorWeight).Msg("--vector-weight must be positive")
	}

	ctx := context.Background()
	provider, err := embedding.ResolveProvider(ctx, embedding.ProviderConfig{
		Type: "ollama", Engine: *engine, Dimensions: 768, BaseURL: *baseURL, CacheType: "none",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("resolve explicit Geppetto Ollama query embedder")
	}
	if provider.Close != nil {
		defer func() { _ = provider.Close() }()
	}
	if provider.Model.Dimensions != 768 {
		log.Fatal().Int("dimensions", provider.Model.Dimensions).Msg("configured embedding model does not match immutable 768D set")
	}

	database, err := db.OpenDB(*databasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("migrate database")
	}
	queries := db.NewQueries(database)
	cards, err := raglab.LoadEvaluationCards(ctx, queries, datasetID)
	if err != nil {
		log.Fatal().Err(err).Msg("load immutable evaluation cards")
	}

	service := experimentrun.NewService(queries)
	lab := raglab.NewLaboratory(raglab.NewSQLiteCatalog(queries), service, true)
	executor := raglab.NewExecutor(raglab.NewSQLiteChannelRetriever(queries), service)
	plans := []raglab.ExperimentSpecification{
		mustBuild(vectorPlan()),
		mustBuild(weightedRRFPlan(*vectorWeight)),
	}
	completed := make([]observation, 0, len(plans))
	for _, specification := range plans {
		run, err := lab.Start(ctx, specification)
		if err != nil {
			log.Fatal().Err(err).Str("experiment", specification.Name).Msg("persist and start immutable experiment")
		}
		result, err := executor.Execute(ctx, run.ID, specification, cards, raglab.ExecutionOptions{Embedder: provider.Provider})
		if err != nil {
			log.Fatal().Err(err).Str("run_id", run.ID).Str("experiment", specification.Name).Msg("execute TTC retrieval observation")
		}
		completed = append(completed, observation{
			Name: specification.Name, Fingerprint: specification.Fingerprint, RunID: result.RunID, Metrics: result.Metrics, Timing: result.Timing,
		})
		log.Info().Str("run_id", result.RunID).Str("experiment", specification.Name).Interface("metrics", result.Metrics).Interface("timing", result.Timing).Msg("completed fresh TTC observation")
	}

	result := report{
		Provider: map[string]any{
			"type": provider.ProviderType, "model": provider.Model.Name, "dimensions": provider.Model.Dimensions, "base_url": *baseURL,
		},
		Cost: map[string]any{
			"provider_billing_usd": 0,
			"scope":                "No billed API-provider cost: embeddings used the user-owned mimimi Mac Ollama service. Hardware amortization and energy are not estimated.",
		},
		Storage: artifactStorage(ctx, queries),
		Runs:    completed,
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		log.Fatal().Err(err).Msg("write experiment report")
	}
}

func vectorPlan() (*raglab.ExperimentBuilder, error) {
	return raglab.NewExperiment("ttc-vector-executor-v1").
		Corpus(raglab.CorpusSnapshot(snapshotID)).Chunks(raglab.ChunkSet(chunkSetID)).Embeddings(raglab.EmbeddingSet(embeddingSetID)).Evaluation(raglab.EvaluationDataset(datasetID)).
		Note("Fresh vector-only TTC observation using an explicit mimimi-tunnel Geppetto query embedder.").
		Tag("provider", "ollama/nomic-embed-text/768").
		Retrieval(func(builder *raglab.RetrievalBuilder) {
			builder.Channel("semantic", func(channel *raglab.ChannelBuilder) { channel.Vector().Representation("raw").TopK(50) }).Collapse(raglab.CollapseDocument).Results(10)
		}).
		Metrics(func(builder *raglab.MetricsBuilder) {
			builder.RelevanceAt(raglab.RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).RecallAt(10).MRR()
		}), nil
}

func weightedRRFPlan(vectorWeight float64) (*raglab.ExperimentBuilder, error) {
	return raglab.NewExperiment("ttc-weighted-rrf-executor-v1").
		Corpus(raglab.CorpusSnapshot(snapshotID)).Chunks(raglab.ChunkSet(chunkSetID)).BM25(raglab.BM25Index(bm25ID)).Embeddings(raglab.EmbeddingSet(embeddingSetID)).Evaluation(raglab.EvaluationDataset(datasetID)).
		Note("Fresh lexical plus semantic TTC observation; explicit weighted RRF gives semantic evidence weight 2 and lexical evidence default weight 1.").
		Tag("provider", "ollama/nomic-embed-text/768").Tag("fusion", "weighted-rrf-60").
		Retrieval(func(builder *raglab.RetrievalBuilder) {
			builder.Channel("lexical", func(channel *raglab.ChannelBuilder) { channel.BM25().Representation("raw").TopK(50) })
			builder.Channel("semantic", func(channel *raglab.ChannelBuilder) { channel.Vector().Representation("raw").TopK(50) })
			builder.FuseRRF(60).Weight("semantic", vectorWeight).Collapse(raglab.CollapseDocument).Results(10)
		}).
		Metrics(func(builder *raglab.MetricsBuilder) {
			builder.RelevanceAt(raglab.RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).RecallAt(10).MRR()
		}), nil
}

func mustBuild(builder *raglab.ExperimentBuilder, err error) raglab.ExperimentSpecification {
	if err != nil {
		log.Fatal().Err(err).Msg("configure experiment builder")
	}
	specification, err := builder.Build()
	if err != nil {
		log.Fatal().Err(err).Msg("build immutable experiment specification")
	}
	return specification
}

func artifactStorage(ctx context.Context, queries *db.Queries) map[string]any {
	storage := map[string]any{"scope": "Shared immutable artifact storage; each append-only run stores traces but does not duplicate vectors or the BM25 index."}
	var vectorCount, vectorBytes int64
	if err := queries.DB().QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(LENGTH(vector)), 0) FROM immutable_embeddings WHERE embedding_set_id=?`, embeddingSetID).Scan(&vectorCount, &vectorBytes); err == nil {
		storage["embedding_vectors"] = vectorCount
		storage["embedding_bytes"] = vectorBytes
	}
	var bm25Path string
	if err := queries.DB().QueryRowContext(ctx, `SELECT artifact_path FROM retrieval_artifacts WHERE id=?`, bm25ID).Scan(&bm25Path); err == nil {
		if bytes, err := directoryBytes(bm25Path); err == nil {
			storage["bm25_bytes"] = bytes
		}
	}
	return storage
}

func directoryBytes(path string) (int64, error) {
	var total int64
	err := filepath.WalkDir(path, func(_ string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		return nil
	})
	return total, err
}
