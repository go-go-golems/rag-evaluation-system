package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	databasePath := flag.String("db", "data/rag-eval.db", "SQLite database")
	chunkSetID := flag.String("chunk-set-id", "sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392", "immutable chunk set ID")
	artifactRoot := flag.String("artifact-root", "data/artifacts/bm25", "artifact root")
	logLevel := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	level, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(level)
	q, err := cmdhelpers.OpenDBAtPath(*databasePath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer q.Close()
	r, err := immutableretrieval.BuildBM25(context.Background(), q, immutableretrieval.BM25BuildRequest{ChunkSetID: *chunkSetID, ArtifactRoot: *artifactRoot})
	if err != nil {
		log.Fatal().Err(err).Msg("build immutable BM25")
	}
	if err := json.NewEncoder(os.Stdout).Encode(r); err != nil {
		log.Fatal().Err(err).Msg("encode result")
	}
}
