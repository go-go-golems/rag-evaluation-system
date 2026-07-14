// Probe the Geppetto embedding factory against the local Ollama service.
// Run from the repository root:
// GOWORK=off go run ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableembedding"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logLevel := flag.String("log-level", "info", "zerolog level")
	engine := flag.String("engine", "nomic-embed-text", "Ollama embedding model")
	batchSize := flag.Int("batch-size", 1, "number of texts for GenerateBatchEmbeddings; 1 uses GenerateEmbedding")
	databasePath := flag.String("db", "", "RAG SQLite database; requires --chunk-set-id")
	chunkSetID := flag.String("chunk-set-id", "", "immutable chunk set ID; requires --db")
	flag.Parse()
	level, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(level)
	resolved, err := embedding.ResolveProvider(context.Background(), embedding.ProviderConfig{Type: "ollama", Engine: *engine, Dimensions: 768, CacheType: "none"})
	if err != nil {
		log.Fatal().Err(err).Msg("resolve Geppetto Ollama embedding provider")
	}
	if resolved.Close != nil {
		defer func() { _ = resolved.Close() }()
	}
	if *batchSize <= 0 {
		log.Fatal().Msg("batch size must be positive")
	}
	texts := make([]string, *batchSize)
	for i := range texts {
		texts[i] = "TTC Geppetto embedding probe"
	}
	var vectors [][]float32
	if *batchSize == 1 {
		vector, err := resolved.Provider.GenerateEmbedding(context.Background(), texts[0])
		if err != nil {
			log.Fatal().Err(err).Msg("generate embedding")
		}
		vectors = [][]float32{vector}
	} else {
		vectors, err = resolved.Provider.GenerateBatchEmbeddings(context.Background(), texts)
		if err != nil {
			log.Fatal().Err(err).Msg("generate embedding batch")
		}
	}
	if len(vectors) != *batchSize || len(vectors[0]) != resolved.Model.Dimensions {
		log.Fatal().Int("vectors", len(vectors)).Int("first_dimensions", len(vectors[0])).Msg("provider returned unexpected batch shape")
	}
	result := map[string]any{"provider_type": resolved.ProviderType, "model": resolved.Model.Name, "model_dimensions": resolved.Model.Dimensions, "batch_size": *batchSize, "vector_dimensions": len(vectors[0])}
	if *databasePath != "" || *chunkSetID != "" {
		if *databasePath == "" || *chunkSetID == "" {
			log.Fatal().Msg("--db and --chunk-set-id must be supplied together")
		}
		queries, err := cmdhelpers.OpenDBAtPath(*databasePath)
		if err != nil {
			log.Fatal().Err(err).Msg("open RAG database")
		}
		defer func() { _ = queries.Close() }()
		set, err := immutableembedding.Build(context.Background(), queries, immutableembedding.Request{ChunkSetID: *chunkSetID, ProviderType: resolved.ProviderType, Provider: resolved.Provider, BatchSize: 16})
		if err != nil {
			log.Fatal().Err(err).Msg("build immutable embedding set")
		}
		result["embedding_set_id"] = set.EmbeddingSetID
		result["embedding_count"] = set.EmbeddingCount
		result["reused"] = set.Reused
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		log.Fatal().Err(err).Msg("write result")
	}
}
