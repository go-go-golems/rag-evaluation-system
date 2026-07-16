//go:build ignore

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"os"
	"regexp"
	"strings"
	"time"

	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type trace struct {
	ID                            string                        `json:"id"`
	Query                         string                        `json:"query"`
	BM25                          []immutableretrieval.ChunkHit `json:"bm25"`
	Vector                        []immutableretrieval.ChunkHit `json:"vector"`
	Hybrid                        []immutableretrieval.FusedHit `json:"hybrid"`
	BM25DurationMilliseconds      int64                         `json:"bm25_duration_ms"`
	EmbeddingDurationMilliseconds int64                         `json:"embedding_duration_ms"`
	VectorDurationMilliseconds    int64                         `json:"vector_duration_ms"`
	FusionDurationMilliseconds    int64                         `json:"fusion_duration_ms"`
	TotalDurationMilliseconds     int64                         `json:"total_duration_ms"`
}

type traceFile struct {
	SchemaVersion string    `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`
	BM25Artifact  string    `json:"bm25_artifact"`
	EmbeddingSet  string    `json:"embedding_set"`
	Provider      string    `json:"provider"`
	TraceLimit    int       `json:"trace_limit"`
	Traces        []trace   `json:"traces"`
}

func main() {
	dbPath := flag.String("db", "data/rag-eval.db", "SQLite database")
	cards := flag.String("cards", "ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md", "candidate card Markdown paths, comma-separated")
	out := flag.String("out", "data/artifacts/traces/ttc-baseline-v1.json", "output JSON")
	bm25 := flag.String("bm25-artifact", "sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691", "BM25 artifact")
	emb := flag.String("embedding-set", "sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0", "embedding set")
	base := flag.String("base-url", "http://127.0.0.1:11435", "Ollama base URL")
	limit := flag.Int("limit", 50, "candidate chunks per retrieval channel")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	l, e := zerolog.ParseLevel(*level)
	if e != nil {
		log.Fatal().Err(e).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(l)
	q, e := cmdhelpers.OpenDBAtPath(*dbPath)
	if e != nil {
		log.Fatal().Err(e).Msg("open db")
	}
	defer q.Close()
	p, e := embedding.ResolveProvider(context.Background(), embedding.ProviderConfig{Type: "ollama", Engine: "nomic-embed-text", Dimensions: 768, BaseURL: *base, CacheType: "file", CacheDirectory: "state/embedding-cache"})
	if e != nil {
		log.Fatal().Err(e).Msg("resolve provider")
	}
	if p.Close != nil {
		defer p.Close()
	}
	idRx := regexp.MustCompile("^#### `([^`]+)`")
	listIDRx := regexp.MustCompile(`^\s*-\s+\{?id:\s*['\"]?([^,}\s]+)`)
	qRx := regexp.MustCompile(`^query: "(.*)"`)
	var traces []trace
	for _, cardPath := range strings.Split(*cards, ",") {
		cardPath = strings.TrimSpace(cardPath)
		f, e := os.Open(cardPath)
		if e != nil {
			log.Fatal().Err(e).Str("path", cardPath).Msg("open cards")
		}
		var id string
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			if m := idRx.FindStringSubmatch(line); m != nil {
				id = m[1]
				continue
			}
			if m := listIDRx.FindStringSubmatch(line); m != nil {
				id = m[1]
				continue
			}
			if m := qRx.FindStringSubmatch(line); m != nil && id != "" {
				query := m[1]
				started := time.Now()
				channelStarted := time.Now()
				b, e := immutableretrieval.QueryBM25(context.Background(), q, *bm25, query, *limit)
				if e != nil {
					log.Fatal().Err(e).Str("id", id).Msg("bm25")
				}
				bm25Elapsed := time.Since(channelStarted)
				channelStarted = time.Now()
				v, e := p.Provider.GenerateEmbedding(context.Background(), query)
				if e != nil {
					log.Fatal().Err(e).Str("id", id).Msg("query embedding")
				}
				embeddingElapsed := time.Since(channelStarted)
				channelStarted = time.Now()
				vec, e := immutableretrieval.QueryVector(context.Background(), q, *emb, v, *limit)
				if e != nil {
					log.Fatal().Err(e).Str("id", id).Msg("vector")
				}
				vectorElapsed := time.Since(channelStarted)
				channelStarted = time.Now()
				hybrid := immutableretrieval.FuseRRF(map[string][]immutableretrieval.ChunkHit{"bm25": b, "vector": vec}, 60, 10)
				fusionElapsed := time.Since(channelStarted)
				traces = append(traces, trace{
					ID: id, Query: query, BM25: b, Vector: vec, Hybrid: hybrid,
					BM25DurationMilliseconds:      bm25Elapsed.Milliseconds(),
					EmbeddingDurationMilliseconds: embeddingElapsed.Milliseconds(),
					VectorDurationMilliseconds:    vectorElapsed.Milliseconds(),
					FusionDurationMilliseconds:    fusionElapsed.Milliseconds(),
					TotalDurationMilliseconds:     time.Since(started).Milliseconds(),
				})
				log.Info().Str("id", id).Int64("total_duration_ms", time.Since(started).Milliseconds()).Msg("completed retrieval trace")
				id = ""
			}
		}
		if e := s.Err(); e != nil {
			log.Fatal().Err(e).Str("path", cardPath).Msg("scan cards")
		}
		_ = f.Close()
	}
	if e := os.MkdirAll("data/artifacts/traces", 0o750); e != nil {
		log.Fatal().Err(e).Msg("mkdir")
	}
	o, e := os.Create(*out)
	if e != nil {
		log.Fatal().Err(e).Msg("create output")
	}
	defer o.Close()
	file := traceFile{
		SchemaVersion: "rag-eval-immutable-retrieval-traces/v1",
		GeneratedAt:   time.Now().UTC(), BM25Artifact: *bm25, EmbeddingSet: *emb,
		Provider: "ollama:nomic-embed-text@" + *base, TraceLimit: *limit, Traces: traces,
	}
	if e := json.NewEncoder(o).Encode(file); e != nil {
		log.Fatal().Err(e).Msg("write traces")
	}
	log.Info().Int("traces", len(traces)).Str("out", *out).Msg("completed immutable retrieval traces")
}
