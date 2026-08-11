//go:build ignore

// 05-score-candidate-retrieval-traces.go scores the reviewed-candidate card
// judgments against immutable retrieval traces. It intentionally labels the
// output provisional: a TTC policy owner must still adjudicate the cards
// before the dataset is frozen as evaluation-dataset/v1.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type trace struct {
	ID                            string                        `json:"id"`
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

type methodMetric struct {
	Queries                int     `json:"queries"`
	AnswerableQueries      int     `json:"answerable_queries"`
	RecallAt1              float64 `json:"recall_at_1"`
	RecallAt3              float64 `json:"recall_at_3"`
	RecallAt10             float64 `json:"recall_at_10"`
	MeanReciprocalRank     float64 `json:"mean_reciprocal_rank"`
	MeanRelevantRecallAt10 float64 `json:"mean_relevant_recall_at_10"`
}

type report struct {
	SchemaVersion             string                    `json:"schema_version"`
	DatasetStatus             string                    `json:"dataset_status"`
	TraceSchemaVersion        string                    `json:"trace_schema_version"`
	BM25Artifact              string                    `json:"bm25_artifact"`
	EmbeddingSet              string                    `json:"embedding_set"`
	Provider                  string                    `json:"provider"`
	Methods                   map[string]methodMetric   `json:"methods"`
	LatencyMilliseconds       map[string]summary        `json:"latency_ms"`
	EmbeddingGenerationCost   cost                      `json:"embedding_generation_cost"`
	StorageBytes              storage                   `json:"storage_bytes"`
	ScoredCardIDs             []string                  `json:"scored_card_ids"`
	Judgments                 map[string]map[string]int `json:"named_relevance_judgments"`
	HumanAdjudicationRequired string                    `json:"human_adjudication_required"`
}

type summary struct{ Total, Mean, P50, P95 int64 }
type cost struct {
	Currency string  `json:"currency"`
	Billed   float64 `json:"billed"`
	Note     string  `json:"note"`
}
type storage struct{ BM25Artifact, TraceFile, SQLiteEmbeddingBlobs, SQLiteDatabase int64 }

func main() {
	dbPath := flag.String("db", "data/rag-eval.db", "SQLite database")
	cardsPath := flag.String("cards", "ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md", "candidate cards")
	tracesPath := flag.String("traces", "data/artifacts/traces/ttc-baseline-v1.json", "trace JSON")
	out := flag.String("out", "data/artifacts/metrics/ttc-baseline-v1-candidate-retrieval.json", "output JSON")
	bm25Path := flag.String("bm25-path", "data/artifacts/bm25/sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691", "BM25 artifact directory")
	level := flag.String("log-level", "info", "zerolog level")
	flag.Parse()
	l, err := zerolog.ParseLevel(*level)
	if err != nil {
		log.Fatal().Err(err).Msg("parse log level")
	}
	zerolog.SetGlobalLevel(l)

	judgments, ids, err := parseCards(*cardsPath)
	if err != nil {
		log.Fatal().Err(err).Msg("parse candidate cards")
	}
	relevant := thresholdJudgments(judgments, 2)
	var traces traceFile
	file, err := os.Open(*tracesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("open traces")
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&traces); err != nil {
		log.Fatal().Err(err).Msg("decode traces")
	}
	q, err := cmdhelpers.OpenDBAtPath(*dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer q.Close()

	methods := map[string]methodMetric{}
	latencies := map[string][]int64{"bm25": {}, "embedding": {}, "vector": {}, "fusion": {}, "total": {}}
	for _, t := range traces.Traces {
		relevantDocs, known := relevant[t.ID]
		if !known {
			log.Fatal().Str("id", t.ID).Msg("trace is not a scored candidate card")
		}
		for name, hits := range map[string][]immutableretrieval.ChunkHit{"bm25": t.BM25, "vector": t.Vector, "hybrid": fusedChunks(t.Hybrid)} {
			metric := methods[name]
			metric.Queries++
			if len(relevantDocs) > 0 {
				metric.AnswerableQueries++
				rank, found := firstRelevantRank(hits, relevantDocs, q)
				if found {
					metric.MeanReciprocalRank += 1 / float64(rank)
				}
				if rank > 0 && rank <= 1 {
					metric.RecallAt1++
				}
				if rank > 0 && rank <= 3 {
					metric.RecallAt3++
				}
				if rank > 0 && rank <= 10 {
					metric.RecallAt10++
				}
				metric.MeanRelevantRecallAt10 += relevantRecallAt10(hits, relevantDocs, q)
			}
			methods[name] = metric
		}
		latencies["bm25"] = append(latencies["bm25"], t.BM25DurationMilliseconds)
		latencies["embedding"] = append(latencies["embedding"], t.EmbeddingDurationMilliseconds)
		latencies["vector"] = append(latencies["vector"], t.VectorDurationMilliseconds)
		latencies["fusion"] = append(latencies["fusion"], t.FusionDurationMilliseconds)
		latencies["total"] = append(latencies["total"], t.TotalDurationMilliseconds)
	}
	for name, metric := range methods {
		if metric.AnswerableQueries > 0 {
			n := float64(metric.AnswerableQueries)
			metric.RecallAt1 /= n
			metric.RecallAt3 /= n
			metric.RecallAt10 /= n
			metric.MeanReciprocalRank /= n
			metric.MeanRelevantRecallAt10 /= n
		}
		methods[name] = metric
	}
	latencySummary := map[string]summary{}
	for name, values := range latencies {
		latencySummary[name] = summarise(values)
	}
	artifacts, err := sizePath(*bm25Path)
	if err != nil {
		log.Fatal().Err(err).Msg("measure BM25 artifact")
	}
	traceSize, err := sizePath(*tracesPath)
	if err != nil {
		log.Fatal().Err(err).Msg("measure trace file")
	}
	var vectorBytes int64
	if err := q.DB().QueryRowContext(context.Background(), `SELECT COALESCE(SUM(LENGTH(vector)), 0) FROM immutable_embeddings WHERE embedding_set_id = ?`, traces.EmbeddingSet).Scan(&vectorBytes); err != nil {
		log.Fatal().Err(err).Msg("measure embedding blobs")
	}
	dbSize, err := sizePath(*dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("measure SQLite database")
	}
	report := report{"rag-eval-candidate-retrieval-score/v1", "candidate-source-validated-not-human-frozen", traces.SchemaVersion, traces.BM25Artifact, traces.EmbeddingSet, traces.Provider, methods, latencySummary, cost{"USD", 0, "Ollama nomic-embed-text ran on the user-owned Mac through SSH loopback; no provider-billed embedding cost. Local compute energy and hardware amortization are not measured."}, storage{artifacts, traceSize, vectorBytes, dbSize}, ids, judgments, "A TTC policy owner must review named judgments, source precedence, and answerability rules before evaluation-dataset/v1 is frozen."}
	if err := os.MkdirAll(filepath.Dir(*out), 0o750); err != nil {
		log.Fatal().Err(err).Msg("create output directory")
	}
	o, err := os.Create(*out)
	if err != nil {
		log.Fatal().Err(err).Msg("create report")
	}
	defer o.Close()
	if err := json.NewEncoder(o).Encode(report); err != nil {
		log.Fatal().Err(err).Msg("write report")
	}
	log.Info().Int("queries", len(traces.Traces)).Str("out", *out).Msg("scored provisional candidate retrieval traces")
}

func parseCards(path string) (map[string]map[string]int, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	idRx := regexp.MustCompile("^#### `([^`]+)`")
	judgmentRx := regexp.MustCompile("^- `([0-3])_(?:NOT_RELEVANT|PARTIAL|SUBSTANTIAL|AUTHORITATIVE)` .*`(wp:[0-9]+)`")
	judgments := map[string]map[string]int{}
	var ids []string
	var id string
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if m := idRx.FindStringSubmatch(line); m != nil {
			id = m[1]
			judgments[id] = map[string]int{}
			ids = append(ids, id)
			continue
		}
		if m := judgmentRx.FindStringSubmatch(line); m != nil && id != "" {
			judgments[id][m[2]] = int(m[1][0] - '0')
		}
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	return judgments, ids, nil
}

func thresholdJudgments(judgments map[string]map[string]int, threshold int) map[string]map[string]struct{} {
	result := make(map[string]map[string]struct{}, len(judgments))
	for cardID, byDocument := range judgments {
		result[cardID] = map[string]struct{}{}
		for documentID, grade := range byDocument {
			if grade >= threshold {
				result[cardID][documentID] = struct{}{}
			}
		}
	}
	return result
}

func fusedChunks(hits []immutableretrieval.FusedHit) []immutableretrieval.ChunkHit {
	out := make([]immutableretrieval.ChunkHit, len(hits))
	for i := range hits {
		out[i] = hits[i].ChunkHit
	}
	return out
}

func firstRelevantRank(hits []immutableretrieval.ChunkHit, relevant map[string]struct{}, q interface{ DB() *sql.DB }) (int, bool) {
	for _, hit := range hits {
		if stableDocumentID(hit.DocumentRevisionID, q) != "" {
			if _, ok := relevant[stableDocumentID(hit.DocumentRevisionID, q)]; ok {
				return hit.Rank, true
			}
		}
	}
	return 0, false
}

func relevantRecallAt10(hits []immutableretrieval.ChunkHit, relevant map[string]struct{}, q interface{ DB() *sql.DB }) float64 {
	if len(relevant) == 0 {
		return 0
	}
	found := map[string]struct{}{}
	for _, hit := range hits {
		if hit.Rank > 10 {
			break
		}
		if id := stableDocumentID(hit.DocumentRevisionID, q); id != "" {
			if _, ok := relevant[id]; ok {
				found[id] = struct{}{}
			}
		}
	}
	return float64(len(found)) / float64(len(relevant))
}

func stableDocumentID(revisionID string, q interface{ DB() *sql.DB }) string {
	var stable string
	if err := q.DB().QueryRowContext(context.Background(), `SELECT stable_document_id FROM document_revisions WHERE id = ?`, revisionID).Scan(&stable); err != nil {
		return ""
	}
	return stable
}

func summarise(values []int64) summary {
	if len(values) == 0 {
		return summary{}
	}
	sorted := append([]int64(nil), values...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	var total int64
	for _, value := range sorted {
		total += value
	}
	return summary{Total: total, Mean: total / int64(len(sorted)), P50: sorted[(len(sorted)-1)/2], P95: sorted[(95*len(sorted)+99)/100-1]}
}

func sizePath(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return info.Size(), nil
	}
	var total int64
	err = filepath.Walk(path, func(_ string, item os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !item.IsDir() {
			total += item.Size()
		}
		return nil
	})
	return total, err
}
