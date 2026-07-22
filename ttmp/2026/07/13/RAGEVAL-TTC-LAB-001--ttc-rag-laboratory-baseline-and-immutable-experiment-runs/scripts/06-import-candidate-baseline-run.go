//go:build ignore

// 06-import-candidate-baseline-run.go imports the already-measured 20-card
// retrieval artifact as one append-only experiment run. It does not claim the
// candidate card set is human-frozen evaluation truth.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	corpusSnapshotID = "sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409"
	chunkSetID       = "sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392"
	bm25ArtifactID   = "sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691"
	embeddingSetID   = "sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0"
)

type trace struct {
	ID                            string `json:"id"`
	BM25DurationMilliseconds      int64  `json:"bm25_duration_ms"`
	EmbeddingDurationMilliseconds int64  `json:"embedding_duration_ms"`
	VectorDurationMilliseconds    int64  `json:"vector_duration_ms"`
	FusionDurationMilliseconds    int64  `json:"fusion_duration_ms"`
	TotalDurationMilliseconds     int64  `json:"total_duration_ms"`
}

type traceFile struct {
	SchemaVersion string            `json:"schema_version"`
	BM25Artifact  string            `json:"bm25_artifact"`
	EmbeddingSet  string            `json:"embedding_set"`
	Provider      string            `json:"provider"`
	TraceLimit    int               `json:"trace_limit"`
	Traces        []json.RawMessage `json:"traces"`
}

type metricFile struct {
	DatasetStatus           string          `json:"dataset_status"`
	Methods                 json.RawMessage `json:"methods"`
	LatencyMilliseconds     json.RawMessage `json:"latency_ms"`
	EmbeddingGenerationCost json.RawMessage `json:"embedding_generation_cost"`
	StorageBytes            json.RawMessage `json:"storage_bytes"`
}

func main() {
	dbPath := flag.String("db", "data/rag-eval.db", "SQLite database")
	tracesPath := flag.String("traces", "data/artifacts/traces/ttc-baseline-v1.json", "retrieval trace JSON")
	metricsPath := flag.String("metrics", "data/artifacts/metrics/ttc-baseline-v1-candidate-retrieval.json", "candidate metric JSON")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	l, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(l)
	queries, err := cmdhelpers.OpenDBAtPath(*dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer queries.Close()
	traces := readJSON[traceFile](*tracesPath)
	metrics := readJSON[metricFile](*metricsPath)
	if traces.BM25Artifact != bm25ArtifactID || traces.EmbeddingSet != embeddingSetID {
		log.Fatal().Str("bm25", traces.BM25Artifact).Str("embedding", traces.EmbeddingSet).Msg("trace artifact identity mismatch")
	}
	service := experimentrun.NewService(queries)
	specification, reused, err := service.CreateSpecification(context.Background(), experimentrun.SpecificationInput{
		CorpusSnapshotID: corpusSnapshotID, ChunkSetID: chunkSetID, BM25ArtifactID: bm25ArtifactID, EmbeddingSetID: embeddingSetID,
		EvaluationDatasetID: "candidate:ttc-baseline-v1",
		Config:              map[string]any{"retrieval_schema": traces.SchemaVersion, "channels": []string{"bm25", "vector"}, "fusion": map[string]any{"algorithm": "rrf", "rank_constant": 60, "document_collapse": true}, "trace_limit": traces.TraceLimit, "provider": traces.Provider},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("create specification")
	}
	run, err := service.CreateRun(context.Background(), specification.ID)
	if err != nil {
		log.Fatal().Err(err).Msg("create run")
	}
	if _, err := service.AppendEvent(context.Background(), run.ID, "candidate_trace_import_started", json.RawMessage(`{"source":"ticket-local measured trace artifact"}`)); err != nil {
		log.Fatal().Err(err).Msg("append import event")
	}
	for _, raw := range traces.Traces {
		var item trace
		if err := json.Unmarshal(raw, &item); err != nil {
			log.Fatal().Err(err).Msg("decode trace")
		}
		timing, err := json.Marshal(map[string]int64{"bm25_ms": item.BM25DurationMilliseconds, "embedding_ms": item.EmbeddingDurationMilliseconds, "vector_ms": item.VectorDurationMilliseconds, "fusion_ms": item.FusionDurationMilliseconds, "total_ms": item.TotalDurationMilliseconds})
		if err != nil {
			log.Fatal().Err(err).Msg("encode timing")
		}
		if err := service.RecordQueryTrace(context.Background(), run.ID, experimentrun.QueryTraceInput{QueryCardID: item.ID, Trace: raw, Metrics: json.RawMessage(`{"dataset_status":"candidate-source-validated-not-human-frozen"}`), Timing: timing, Cost: metrics.EmbeddingGenerationCost, Storage: json.RawMessage(`{}`)}); err != nil {
			log.Fatal().Err(err).Str("query_card", item.ID).Msg("record trace")
		}
	}
	if _, err := service.AppendEvent(context.Background(), run.ID, "candidate_trace_imported", json.RawMessage(`{"query_count":20}`)); err != nil {
		log.Fatal().Err(err).Msg("append completion event")
	}
	if _, err := service.CompleteRun(context.Background(), run.ID, experimentrun.SummaryInput{Status: "succeeded", Metrics: json.RawMessage(`{"methods":` + string(metrics.Methods) + `,"latency_ms":` + string(metrics.LatencyMilliseconds) + `}`), Cost: metrics.EmbeddingGenerationCost, Storage: metrics.StorageBytes, Error: json.RawMessage(`{}`)}); err != nil {
		log.Fatal().Err(err).Msg("complete run")
	}
	log.Info().Str("specification_id", specification.ID).Bool("specification_reused", reused).Str("run_id", run.ID).Int("traces", len(traces.Traces)).Msg("imported append-only candidate baseline run")
}

func readJSON[T any](path string) T {
	var result T
	f, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Str("path", path).Msg("open JSON")
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&result); err != nil {
		log.Fatal().Err(err).Str("path", path).Msg("decode JSON")
	}
	return result
}
