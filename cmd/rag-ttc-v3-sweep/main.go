package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/workflowv3ttc"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
)

type batchEvidence struct {
	Key        string                             `json:"key"`
	Generation workflowv3ttc.ProviderMeasurement  `json:"generation"`
	Embedding  workflowv3ttc.EmbeddingMeasurement `json:"embedding"`
}
type attemptEvidence struct {
	NodeKey    string `json:"nodeKey"`
	Number     int    `json:"number"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
}
type cellEvidence struct {
	Cell              workflowv3ttc.SweepCell `json:"cell"`
	RunID             string                  `json:"runId"`
	PlanDigest        string                  `json:"planDigest"`
	Requests          int                     `json:"requests"`
	Chunks            int                     `json:"chunks"`
	MakespanMicros    int64                   `json:"makespanMicros"`
	ProviderMicros    []int64                 `json:"providerMicros"`
	Batches           []batchEvidence         `json:"batches"`
	Attempts          []attemptEvidence       `json:"attempts"`
	EmbeddingAttempts []attemptEvidence       `json:"embeddingAttempts"`
	EmbeddingRequests int                     `json:"embeddingRequests"`
	OverlapMicros     int64                   `json:"overlapMicros"`
	PeakActive        int                     `json:"peakActive"`
	ChunksPerSecond   float64                 `json:"chunksPerSecond"`
	RequestsPerSecond float64                 `json:"requestsPerSecond"`
	Usage             map[string]int64        `json:"usage"`
}
type evidence struct {
	SchemaVersion string                  `json:"schemaVersion"`
	Profile       string                  `json:"profile"`
	Plan          workflowv3ttc.SweepPlan `json:"plan"`
	Cells         []cellEvidence          `json:"cells"`
}

type delayedGenerator struct{ ragoperators.FixtureProviders }

func (d delayedGenerator) Generate(ctx context.Context, req ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	var payload struct {
		Items []json.RawMessage `json:"items"`
	}
	_ = json.Unmarshal([]byte(req.Text), &payload)
	delay := time.Duration(4+2*len(payload.Items)) * time.Millisecond
	select {
	case <-ctx.Done():
		return ragoperators.GenerationResult{}, ctx.Err()
	case <-time.After(delay):
	}
	return d.FixtureProviders.Generate(ctx, req)
}

func main() {
	var output, concurrencyText, profile, providerConfig, specification, artifactRoot string
	var chunkCount, maximum int
	var executeReal bool
	var maximumCost, maximumInputTokens, maximumOutputTokens, maximumEmbeddingTokens int64
	var maximumEmbeddingRequests int
	flag.StringVar(&output, "output", "sweep-evidence", "output directory")
	flag.IntVar(&chunkCount, "chunks", 16, "fixed chunk count")
	flag.IntVar(&maximum, "maximum-requests", 90, "hard request ceiling")
	flag.StringVar(&concurrencyText, "concurrency", "1,2,4", "comma-separated concurrency levels, each at most four")
	flag.StringVar(&profile, "profile", "fixtures", "provider profile: fixtures or real")
	flag.StringVar(&providerConfig, "provider-config", "", "host-only real provider YAML")
	flag.StringVar(&specification, "specification", "", "canonical researchctl execution specification JSON")
	flag.StringVar(&artifactRoot, "artifact-root", "", "root containing specification input URIs")
	flag.BoolVar(&executeReal, "execute-real", false, "authorize submission after all exact budget checks")
	flag.Int64Var(&maximumCost, "maximum-cost-microunits", 0, "authorized real-run cost ceiling")
	flag.Int64Var(&maximumInputTokens, "maximum-input-tokens", 0, "authorized real-run input-token ceiling")
	flag.Int64Var(&maximumOutputTokens, "maximum-output-tokens", 0, "authorized real-run output-token ceiling")
	flag.Int64Var(&maximumEmbeddingTokens, "maximum-embedding-tokens", 0, "authorized real-run embedding-token ceiling")
	flag.IntVar(&maximumEmbeddingRequests, "maximum-embedding-requests", 0, "authorized real-run embedding-request ceiling")
	flag.Parse()
	concurrency, err := parseLevels(concurrencyText)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	config := runConfig{profile: profile, providerConfig: providerConfig, specification: specification, artifactRoot: artifactRoot, executeReal: executeReal, maximumCost: maximumCost, maximumInputTokens: maximumInputTokens, maximumOutputTokens: maximumOutputTokens, maximumEmbeddingTokens: maximumEmbeddingTokens, maximumEmbeddingRequests: maximumEmbeddingRequests}
	if err := run(context.Background(), output, chunkCount, maximum, concurrency, config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type runConfig struct {
	profile, providerConfig, specification, artifactRoot                         string
	executeReal                                                                  bool
	maximumCost, maximumInputTokens, maximumOutputTokens, maximumEmbeddingTokens int64
	maximumEmbeddingRequests                                                     int
}

func run(ctx context.Context, output string, chunkCount, maximum int, concurrency []int, config runConfig) error {
	absoluteOutput, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	output = absoluteOutput
	plan, err := workflowv3ttc.PlanSweep(workflowv3ttc.SweepSpec{ChunkCount: chunkCount, BatchSizes: []int{1, 2, 4, 8}, Concurrency: concurrency, Replicates: 1, MaximumRequests: maximum})
	if err != nil {
		return err
	}
	authority, err := loadProviderAuthority(ctx, config.profile, config.providerConfig, config.specification)
	if err != nil {
		return err
	}
	defer func() { _ = authority.close() }()
	if config.profile == "real" {
		expectedEmbeddingRequests := len(plan.Cells) * chunkCount
		requiredEmbeddingTokens := int64(plan.PlannedRequests) * 65536
		if plan.PlannedRequests != maximum {
			return fmt.Errorf("real request authority must equal the exact planned request count")
		}
		if !config.executeReal {
			fmt.Printf("real dry-run profile_digest=%s model_digest=%s planned_generation_requests=%d planned_embedding_requests=%d required_maximum_cost_microunits=%d required_input_tokens=%d required_output_tokens=%d required_embedding_tokens=%d\n", authority.profileDigest, authority.modelDigest, plan.PlannedRequests, expectedEmbeddingRequests, int64(plan.PlannedRequests)*160000, int64(plan.PlannedRequests)*16384, int64(plan.PlannedRequests)*16384, requiredEmbeddingTokens)
			return nil
		}
		if config.maximumCost < int64(plan.PlannedRequests)*160000 || config.maximumInputTokens < int64(plan.PlannedRequests)*16384 || config.maximumOutputTokens < int64(plan.PlannedRequests)*16384 || config.maximumEmbeddingTokens < requiredEmbeddingTokens || config.maximumEmbeddingRequests != expectedEmbeddingRequests {
			return fmt.Errorf("real numeric authority is below the task reservation maximum")
		}
	}
	if err := os.RemoveAll(output); err != nil {
		return err
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return err
	}
	registry, err := workflowv3ttc.Registry()
	if err != nil {
		return err
	}
	catalog, err := registry.Catalog()
	if err != nil {
		return err
	}
	authored, err := workflowmodule.Author(ctx, workflowv3ttc.SweepWorkflowSource(), catalog, workflowv3ttc.DescriptorModule())
	if err != nil {
		return err
	}
	result := evidence{SchemaVersion: "rag-ttc-v3-sweep-evidence/v1", Profile: config.profile, Plan: plan}
	chunks := fixtureChunks(chunkCount)
	if config.profile == "real" {
		chunks, err = loadRealChunks(ctx, config.specification, config.artifactRoot, chunkCount)
		if err != nil {
			return err
		}
	}
	for index, cell := range plan.Cells {
		cellRoot := filepath.Join(output, fmt.Sprintf("cell-%02d-b%d-c%d", index, cell.ChunksPerRequest, cell.Concurrency))
		artifacts, err := workflowv3.NewFileArtifactStore(filepath.Join(cellRoot, "artifacts"), 1<<30)
		if err != nil {
			return err
		}
		manifest, err := workflowv3ttc.MaterializeBatches(ctx, artifacts, chunks, cell.ChunksPerRequest)
		if err != nil {
			return err
		}
		provider, err := authority.provider(cell.ChunksPerRequest)
		if err != nil {
			return err
		}
		modules, err := workflowv3runtime.NewTaskModuleRegistry(workflowv3ttc.Module(provider))
		if err != nil {
			return err
		}
		store, err := workflowv3sqlite.Open(ctx, filepath.Join(cellRoot, "workflow.sqlite"))
		if err != nil {
			return err
		}
		engine := &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: time.Minute}
		runID := workflowv3.RunID(fmt.Sprintf("fixture-b%d-c%d-r%d", cell.ChunksPerRequest, cell.Concurrency, cell.Replicate))
		if err := engine.Submit(ctx, runID, authored.Plan, map[string]workflowv3.ArtifactRef{"batches": manifest}); err != nil {
			_ = store.Close()
			return err
		}
		dispatchCtx, cancel := context.WithCancel(ctx)
		done := make(chan error, 1)
		dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{workflowv3ttc.ResourceGeneration: cell.Concurrency, workflowv3ttc.ResourceEmbedding: 4}, PollInterval: time.Millisecond}
		go func() { done <- dispatcher.Run(dispatchCtx) }()
		var snapshot workflowv3.RunSnapshot
		deadline := time.Now().Add(30 * time.Second)
		for time.Now().Before(deadline) {
			snapshot, err = engine.Snapshot(ctx, runID)
			if err == nil && (snapshot.Status == "succeeded" || snapshot.Status == "failed") {
				break
			}
			select {
			case dispatchErr := <-done:
				cancel()
				_ = store.Close()
				return fmt.Errorf("dispatcher: %w", dispatchErr)
			case <-time.After(time.Millisecond):
			}
		}
		cancel()
		<-done
		if snapshot.Status != "succeeded" {
			_ = store.Close()
			return fmt.Errorf("cell %+v status %s", cell, snapshot.Status)
		}
		budget, err := store.BudgetSnapshot(ctx, &runID)
		if err != nil {
			_ = store.Close()
			return err
		}
		cellResult, err := readCell(ctx, artifacts, snapshot, cell, chunkCount, budget)
		if err != nil {
			_ = store.Close()
			return err
		}
		result.Cells = append(result.Cells, cellResult)
		if err := store.Close(); err != nil {
			return err
		}
		if err := os.RemoveAll(cellRoot); err != nil {
			return err
		}
	}
	body, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(filepath.Join(output, "evidence.json"), append(body, '\n'), 0o644); err != nil {
		return err
	}
	if err = writeCSV(filepath.Join(output, "cells.csv"), result.Cells); err != nil {
		return err
	}
	fmt.Printf("profile=%s cells=%d planned_requests=%d evidence=%s\n", config.profile, len(result.Cells), plan.PlannedRequests, filepath.Join(output, "evidence.json"))
	return nil
}

func parseLevels(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	levels := make([]int, len(parts))
	for i, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid concurrency %q", part)
		}
		levels[i] = n
	}
	return levels, nil
}

func fixtureChunks(count int) []workflowv3ttc.Chunk {
	chunks := make([]workflowv3ttc.Chunk, count)
	for i := range chunks {
		key := fmt.Sprintf("chunk-%04d", i)
		text := fmt.Sprintf("Fixture source %04d records deterministic batching and concurrency behavior for the Workflow V3 measurement control.", i)
		td, _ := workflowv3.Digest(text)
		chunks[i] = workflowv3ttc.Chunk{Key: key, Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: key, ParentUnitID: fmt.Sprintf("unit-%04d", i), TextDigest: td, Citation: ragcontract.CitationRef{SourceID: fmt.Sprintf("source-%04d", i)}}, Text: text}, CitationIDs: []string{fmt.Sprintf("source-%04d", i)}, SourceDigest: "sha256:" + strings.Repeat("d", 64)}
	}
	return chunks
}
func fixtureProvider(batch int) (*workflowv3ttc.OperatorProvider, error) {
	fixtures := ragoperators.NewFixtureProviders()
	cfg, _ := json.Marshal(map[string]any{"model": ragoperators.FixtureSummaryModel, "prompt": ragoperators.FixtureSummaryPrompt, "outputSchema": ragoperators.FixtureSummarySchema, "batchSize": batch, "questionsPerChunk": 2, "maxBatchRunes": 100000})
	return workflowv3ttc.NewOperatorProvider(workflowv3ttc.OperatorProviderConfig{GenerationNode: ragcontract.Node{Config: cfg}, EmbeddingNode: ragcontract.Node{Config: json.RawMessage(`{"model":"fixture-hash-32-v1","dimensions":32,"normalize":"none","batchSize":8}`)}, RawRepresentationName: "raw", MaxRepresentationsPerChunk: 8, ProviderProfileDigest: "sha256:" + strings.Repeat("a", 64), GenerationModelDigest: "sha256:" + strings.Repeat("b", 64), EmbeddingProfileDigest: "sha256:" + strings.Repeat("c", 64), ResolveEnvironment: func(context.Context) (*ragoperators.Environment, error) {
		return &ragoperators.Environment{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: delayedGenerator{fixtures}, Embedder: fixtures, Cache: ragoperators.NewMemoryCache(), Usage: ragoperators.Usage{Cost: map[string]float64{}}}, nil
	}})
}
func readCell(ctx context.Context, a workflowv3.ArtifactStore, s workflowv3.RunSnapshot, c workflowv3ttc.SweepCell, chunks int, budget []workflowv3.BudgetProgress) (cellEvidence, error) {
	body, err := workflowv3.ReadArtifact(ctx, a, s.Outputs["measured"])
	if err != nil {
		return cellEvidence{}, err
	}
	m, err := workflowv3.DecodeItemManifest(body)
	if err != nil {
		return cellEvidence{}, err
	}
	durations := make([]int64, 0, len(m.Items))
	batches := make([]batchEvidence, 0, len(m.Items))
	embeddingRequests := 0
	for _, item := range m.Items {
		b, err := workflowv3.ReadArtifact(ctx, a, item.Value)
		if err != nil {
			return cellEvidence{}, err
		}
		var measured workflowv3ttc.MeasuredBatch
		if err = workflowv3.StrictDecode(b, &measured); err != nil {
			return cellEvidence{}, err
		}
		durations = append(durations, measured.Generation.ProviderElapsedMicros)
		batches = append(batches, batchEvidence{Key: measured.Key, Generation: measured.Generation, Embedding: measured.Embedding})
		embeddingRequests += measured.Embedding.ProviderRequests
	}
	attempts, embeddingAttempts := []workflowv3.Attempt{}, []workflowv3.Attempt{}
	for _, x := range s.Attempts {
		if x.Status != "succeeded" {
			continue
		}
		if x.ResourceClass == workflowv3ttc.ResourceGeneration {
			attempts = append(attempts, x)
		}
		if x.ResourceClass == workflowv3ttc.ResourceEmbedding {
			embeddingAttempts = append(embeddingAttempts, x)
		}
	}
	sort.Slice(attempts, func(i, j int) bool { return attempts[i].StartedAt.Before(attempts[j].StartedAt) })
	start, end := attempts[0].StartedAt, attempts[0].FinishedAt
	for _, x := range append(append([]workflowv3.Attempt{}, attempts...), embeddingAttempts...) {
		if x.FinishedAt.After(end) {
			end = x.FinishedAt
		}
	}
	makespan := end.Sub(start)
	usage := map[string]int64{}
	for _, amount := range budget {
		usage[amount.Account+"."+amount.Dimension] = amount.Used
		if amount.Account == "generation" {
			usage[amount.Dimension] = amount.Used
		}
	}
	return cellEvidence{Cell: c, RunID: string(s.RunID), PlanDigest: s.PlanDigest, Requests: len(attempts), Chunks: chunks, MakespanMicros: makespan.Microseconds(), ProviderMicros: durations, Batches: batches, Attempts: evidenceRows(attempts), EmbeddingAttempts: evidenceRows(embeddingAttempts), EmbeddingRequests: embeddingRequests, OverlapMicros: overlapMicros(attempts, embeddingAttempts), PeakActive: peakActive(attempts), ChunksPerSecond: float64(chunks) / makespan.Seconds(), RequestsPerSecond: float64(len(attempts)) / makespan.Seconds(), Usage: usage}, nil
}

func evidenceRows(attempts []workflowv3.Attempt) []attemptEvidence {
	rows := make([]attemptEvidence, len(attempts))
	for i, attempt := range attempts {
		rows[i] = attemptEvidence{NodeKey: string(attempt.NodeKey), Number: attempt.Number, StartedAt: attempt.StartedAt.UTC().Format(time.RFC3339Nano), FinishedAt: attempt.FinishedAt.UTC().Format(time.RFC3339Nano)}
	}
	return rows
}

func overlapMicros(left, right []workflowv3.Attempt) int64 {
	type event struct {
		at          time.Time
		side, delta int
	}
	events := []event{}
	for _, x := range left {
		events = append(events, event{x.StartedAt, 0, 1}, event{x.FinishedAt, 0, -1})
	}
	for _, x := range right {
		events = append(events, event{x.StartedAt, 1, 1}, event{x.FinishedAt, 1, -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].at.Equal(events[j].at) {
			return events[i].delta < events[j].delta
		}
		return events[i].at.Before(events[j].at)
	})
	active := [2]int{}
	var total time.Duration
	var previous time.Time
	for _, e := range events {
		if !previous.IsZero() && active[0] > 0 && active[1] > 0 {
			total += e.at.Sub(previous)
		}
		active[e.side] += e.delta
		previous = e.at
	}
	return total.Microseconds()
}
func peakActive(a []workflowv3.Attempt) int {
	peak := 0
	for _, x := range a {
		n := 0
		for _, y := range a {
			if !y.StartedAt.After(x.StartedAt) && y.FinishedAt.After(x.StartedAt) {
				n++
			}
		}
		if n > peak {
			peak = n
		}
	}
	return peak
}
func writeCSV(path string, cells []cellEvidence) error {
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	w := csv.NewWriter(f)
	_ = w.Write([]string{"batch_size", "concurrency", "replicate", "requests", "chunks", "makespan_us", "peak_active", "chunks_per_second", "requests_per_second"})
	for _, c := range cells {
		_ = w.Write([]string{strconv.Itoa(c.Cell.ChunksPerRequest), strconv.Itoa(c.Cell.Concurrency), strconv.Itoa(c.Cell.Replicate), strconv.Itoa(c.Requests), strconv.Itoa(c.Chunks), strconv.FormatInt(c.MakespanMicros, 10), strconv.Itoa(c.PeakActive), strconv.FormatFloat(c.ChunksPerSecond, 'f', 6, 64), strconv.FormatFloat(c.RequestsPerSecond, 'f', 6, 64)})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
