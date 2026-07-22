// Probe the Geppetto embedding factory against the local Ollama service.
// Run from the repository root:
// GOWORK=off go run ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go
//go:build ignore

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
	baseURL := flag.String("base-url", "", "optional Ollama base URL, for example http://127.0.0.1:11435 through an SSH tunnel")
	batchSize := flag.Int("batch-size", 1, "number of texts for GenerateBatchEmbeddings; 1 uses GenerateEmbedding")
	databasePath := flag.String("db", "", "RAG SQLite database; requires --chunk-set-id")
	chunkSetID := flag.String("chunk-set-id", "", "immutable chunk set ID; requires --db")
	realChunks := flag.Bool("real-chunks", false, "use the first --batch-size texts from --chunk-set-id instead of probe text")
	flag.Parse()
	level, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(level)
	resolved, err := embedding.ResolveProvider(context.Background(), embedding.ProviderConfig{Type: "ollama", Engine: *engine, Dimensions: 768, BaseURL: *baseURL, CacheType: "none"})
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
	if *realChunks {
		if *databasePath == "" || *chunkSetID == "" {
			log.Fatal().Msg("--real-chunks requires --db and --chunk-set-id")
		}
		queries, err := cmdhelpers.OpenDBAtPath(*databasePath)
		if err != nil {
			log.Fatal().Err(err).Msg("open RAG database for real chunks")
		}
		rows, err := queries.DB().QueryContext(context.Background(), `SELECT text FROM immutable_chunks WHERE chunk_set_id = ? ORDER BY document_revision_id, chunk_index LIMIT ?`, *chunkSetID, *batchSize)
		if err != nil {
			_ = queries.Close()
			log.Fatal().Err(err).Msg("load real chunks")
		}
		texts = texts[:0]
		for rows.Next() {
			var text string
			if err := rows.Scan(&text); err != nil {
				_ = rows.Close()
				_ = queries.Close()
				log.Fatal().Err(err).Msg("scan real chunk")
			}
			texts = append(texts, text)
		}
		_ = rows.Close()
		_ = queries.Close()
		if len(texts) != *batchSize {
			log.Fatal().Int("chunks", len(texts)).Msg("not enough real chunks")
		}
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
