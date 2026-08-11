//go:build ignore

// Runs the first executor-backed TTC observation without requiring an
// embedding server. It is intentionally raw-BM25 only; the vector/hybrid
// companion must receive an explicit live query embedder.
package main

import (
	"context"
	"flag"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/go-go-golems/rag-evaluation-system/pkg/raglab"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	snapshotID = "sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409"
	chunkSetID = "sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392"
	bm25ID     = "sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691"
	datasetID  = "candidate:ttc-baseline-v1"
)

func main() {
	databasePath := flag.String("db", "data/rag-eval.db", "rag-eval SQLite database")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	parsed, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(parsed)
	ctx := context.Background()
	database, err := db.OpenDB(*databasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("migrate database")
	}
	queries := db.NewQueries(database)
	specification, err := raglab.NewExperiment("ttc-raw-bm25-executor-v1").
		Corpus(raglab.CorpusSnapshot(snapshotID)).Chunks(raglab.ChunkSet(chunkSetID)).BM25(raglab.BM25Index(bm25ID)).Evaluation(raglab.EvaluationDataset(datasetID)).
		Representations(func(builder *raglab.RepresentationBuilder) { builder.RawChunks("raw") }).
		Retrieval(func(builder *raglab.RetrievalBuilder) {
			builder.Channel("lexical", func(channel *raglab.ChannelBuilder) { channel.BM25().Representation("raw").TopK(50) }).Collapse(raglab.CollapseDocument).Results(10)
		}).
		Metrics(func(builder *raglab.MetricsBuilder) {
			builder.RelevanceAt(raglab.RelevanceGrade{Name: "2_SUBSTANTIAL", Ordinal: 2}).RecallAt(10).MRR()
		}).
		Build()
	if err != nil {
		log.Fatal().Err(err).Msg("build experiment")
	}
	lab := raglab.NewLaboratory(raglab.NewSQLiteCatalog(queries), experimentrun.NewService(queries), true)
	run, err := lab.Start(ctx, specification)
	if err != nil {
		log.Fatal().Err(err).Msg("persist and start experiment")
	}
	cards, err := raglab.LoadEvaluationCards(ctx, queries, datasetID)
	if err != nil {
		log.Fatal().Err(err).Msg("load immutable evaluation cards")
	}
	result, err := raglab.NewExecutor(raglab.NewSQLiteChannelRetriever(queries), experimentrun.NewService(queries)).Execute(ctx, run.ID, specification, cards, raglab.ExecutionOptions{})
	if err != nil {
		log.Fatal().Err(err).Str("run_id", run.ID).Msg("execute raw BM25 experiment")
	}
	log.Info().Str("run_id", result.RunID).Interface("metrics", result.Metrics).Interface("timing", result.Timing).Msg("completed raw BM25 TTC experiment")
}
