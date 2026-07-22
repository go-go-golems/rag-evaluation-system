package main

import (
	"bufio"
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
type operationEvidence struct {
	JSONLPath    string                                     `json:"jsonlPath"`
	ManifestPath string                                     `json:"manifestPath"`
	Manifest     workflowv3.ExternalOperationExportManifest `json:"manifest"`
}

type cellEvidence struct {
	Cell                  workflowv3ttc.SweepCell `json:"cell"`
	RunID                 string                  `json:"runId"`
	PlanDigest            string                  `json:"planDigest"`
	Requests              int                     `json:"requests"`
	Chunks                int                     `json:"chunks"`
	MakespanMicros        int64                   `json:"makespanMicros"`
	ProviderMicros        []int64                 `json:"providerMicros"`
	Batches               []batchEvidence         `json:"batches"`
	Attempts              []attemptEvidence       `json:"attempts"`
	EmbeddingAttempts     []attemptEvidence       `json:"embeddingAttempts"`
	EmbeddingRequests     int                     `json:"embeddingRequests"`
	AttemptOverlapMicros  int64                   `json:"attemptOverlapMicros"`
	ProviderOverlapMicros int64                   `json:"providerOverlapMicros"`
	AttemptPeakActive     int                     `json:"attemptPeakActive"`
	ProviderPeakActive    int                     `json:"providerPeakActive"`
	ChunksPerSecond       float64                 `json:"chunksPerSecond"`
	RequestsPerSecond     float64                 `json:"requestsPerSecond"`
	Usage                 map[string]int64        `json:"usage"`
	Operations            *operationEvidence      `json:"operations,omitempty"`
}
type evidence struct {
	SchemaVersion         string                    `json:"schemaVersion"`
	Profile               string                    `json:"profile"`
	ProviderProfileDigest string                    `json:"providerProfileDigest"`
	GenerationModelDigest string                    `json:"generationModelDigest"`
	WorkflowPlanDigest    string                    `json:"workflowPlanDigest"`
	RegistryGeneration    string                    `json:"registryGeneration"`
	BundleDigest          string                    `json:"bundleDigest"`
	Plan                  workflowv3ttc.SweepPlan   `json:"plan"`
	Cells                 []cellEvidence            `json:"cells"`
	GenerationAuthority   *generationAuthorityState `json:"generationAuthority,omitempty"`
}

type cellCheckpoint struct {
	SchemaVersion         string       `json:"schemaVersion"`
	Profile               string       `json:"profile"`
	ProviderProfileDigest string       `json:"providerProfileDigest"`
	GenerationModelDigest string       `json:"generationModelDigest"`
	Cell                  cellEvidence `json:"cell"`
}

type failedCellCheckpoint struct {
	SchemaVersion string                      `json:"schemaVersion"`
	Cell          workflowv3ttc.SweepCell     `json:"cell"`
	RunID         string                      `json:"runId"`
	RunStatus     string                      `json:"runStatus"`
	Reason        string                      `json:"reason"`
	Attempts      map[string]map[string]int   `json:"attempts"`
	FailureCodes  map[string]int              `json:"failureCodes"`
	Budget        []workflowv3.BudgetProgress `json:"budget"`
	Operations    *operationEvidence          `json:"operations,omitempty"`
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
	var maximumEmbeddingRequests, maximumGenerationRequests, priorGenerationRequests, maximumGenerationRetries int
	var cellTimeout time.Duration
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
	flag.IntVar(&maximumGenerationRequests, "maximum-generation-requests", 0, "authorized cumulative real-run generation-request ceiling")
	flag.IntVar(&priorGenerationRequests, "prior-generation-requests", 0, "conservatively consumed generation requests before this invocation")
	flag.IntVar(&maximumGenerationRetries, "maximum-generation-retries", 0, "authorized generation retry headroom for this matrix")
	flag.DurationVar(&cellTimeout, "cell-timeout", 0, "per-cell terminal deadline (default 30s fixtures, 30m real)")
	flag.Parse()
	concurrency, err := parseLevels(concurrencyText)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	config := runConfig{profile: profile, providerConfig: providerConfig, specification: specification, artifactRoot: artifactRoot, executeReal: executeReal, maximumCost: maximumCost, maximumInputTokens: maximumInputTokens, maximumOutputTokens: maximumOutputTokens, maximumEmbeddingTokens: maximumEmbeddingTokens, maximumEmbeddingRequests: maximumEmbeddingRequests, maximumGenerationRequests: maximumGenerationRequests, priorGenerationRequests: priorGenerationRequests, maximumGenerationRetries: maximumGenerationRetries, cellTimeout: cellTimeout}
	if err := run(context.Background(), output, chunkCount, maximum, concurrency, config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type runConfig struct {
	profile, providerConfig, specification, artifactRoot                         string
	executeReal                                                                  bool
	maximumCost, maximumInputTokens, maximumOutputTokens, maximumEmbeddingTokens int64
	maximumEmbeddingRequests, maximumGenerationRequests, priorGenerationRequests int
	maximumGenerationRetries                                                     int
	cellTimeout                                                                  time.Duration
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
	chunks := fixtureChunks(chunkCount)
	if config.cellTimeout <= 0 {
		config.cellTimeout = 30 * time.Second
		if config.profile == "real" {
			config.cellTimeout = 30 * time.Minute
		}
	}
	if config.profile == "real" {
		chunks, err = loadRealChunks(ctx, config.specification, config.artifactRoot, chunkCount)
		if err != nil {
			return err
		}
		expectedEmbeddingRequests := len(plan.Cells) * chunkCount
		requiredEmbeddingTokens := int64(plan.PlannedRequests) * 65536
		if plan.PlannedRequests != maximum {
			return fmt.Errorf("real request authority must equal the exact planned request count")
		}
		cumulativeGenerationRequests := config.priorGenerationRequests + plan.PlannedRequests + config.maximumGenerationRetries
		if !config.executeReal {
			fmt.Printf("real dry-run profile_digest=%s model_digest=%s frozen_chunks=%d planned_generation_requests=%d prior_generation_requests=%d maximum_generation_retries=%d required_cumulative_generation_requests=%d planned_embedding_requests=%d required_maximum_cost_microunits=%d required_input_tokens=%d required_output_tokens=%d required_embedding_tokens=%d\n", authority.profileDigest, authority.modelDigest, len(chunks), plan.PlannedRequests, config.priorGenerationRequests, config.maximumGenerationRetries, cumulativeGenerationRequests, expectedEmbeddingRequests, int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationCostMicrounitsPerRequest, int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationInputTokensPerRequest, int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationOutputTokensPerRequest, requiredEmbeddingTokens)
			return nil
		}
		if config.priorGenerationRequests < 0 || config.maximumGenerationRetries < 0 || config.maximumGenerationRequests != cumulativeGenerationRequests || config.maximumCost < int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationCostMicrounitsPerRequest || config.maximumInputTokens < int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationInputTokensPerRequest || config.maximumOutputTokens < int64(cumulativeGenerationRequests)*workflowv3ttc.SweepGenerationOutputTokensPerRequest || config.maximumEmbeddingTokens < requiredEmbeddingTokens || config.maximumEmbeddingRequests != expectedEmbeddingRequests {
			return fmt.Errorf("real numeric authority is below the task reservation maximum")
		}
	}
	if config.profile == "real" {
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("RAG_SWEEP_REAL_OUTPUT_EXISTS")
		} else if !os.IsNotExist(err) {
			return err
		}
	} else if err := os.RemoveAll(output); err != nil {
		return err
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return err
	}
	runtimeRoot := filepath.Join(output, "runtime")
	defer func() { _ = os.RemoveAll(runtimeRoot) }()
	var generationAuthority *generationAdmission
	if config.profile == "real" {
		admission, err := newGenerationAdmission(filepath.Join(output, "generation-authority.json"), config.maximumGenerationRequests, config.priorGenerationRequests)
		if err != nil {
			return err
		}
		generationAuthority = admission
		authority.admitGeneration = admission.Admit
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
	bundle, err := workflowv3ttc.Bundle()
	if err != nil {
		return err
	}
	result := evidence{SchemaVersion: "rag-ttc-v3-sweep-evidence/v2", Profile: config.profile, ProviderProfileDigest: authority.profileDigest, GenerationModelDigest: authority.modelDigest, WorkflowPlanDigest: authored.Plan.Digest, RegistryGeneration: registry.Generation(), BundleDigest: bundle.Digest(), Plan: plan}
	for index, cell := range plan.Cells {
		cellName := fmt.Sprintf("cell-%02d-b%d-c%d", index, cell.ChunksPerRequest, cell.Concurrency)
		cellRoot := filepath.Join(runtimeRoot, cellName)
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
		runID := workflowv3.RunID(fmt.Sprintf("%s-b%d-c%d-r%d", config.profile, cell.ChunksPerRequest, cell.Concurrency, cell.Replicate))
		if err := engine.Submit(ctx, runID, authored.Plan, map[string]workflowv3.ArtifactRef{"batches": manifest}); err != nil {
			_ = store.Close()
			return err
		}
		dispatchCtx, cancel := context.WithCancel(ctx)
		done := make(chan error, 1)
		dispatcher := &workflowv3runtime.Dispatcher{Engine: engine, Capacities: map[string]int{workflowv3ttc.ResourceGeneration: cell.Concurrency, workflowv3ttc.ResourceEmbedding: 4}, PollInterval: time.Millisecond}
		go func() { done <- dispatcher.Run(dispatchCtx) }()
		var snapshot workflowv3.RunSnapshot
		deadline := time.Now().Add(config.cellTimeout)
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
		dispatchErr := <-done
		if snapshot.Status != "succeeded" && snapshot.Status != "failed" {
			if current, snapshotErr := engine.Snapshot(ctx, runID); snapshotErr == nil {
				snapshot = current
			}
			budget, budgetErr := store.BudgetSnapshot(ctx, &runID)
			if budgetErr != nil {
				_ = store.Close()
				return fmt.Errorf("snapshot budget for timed-out cell %+v: %w", cell, budgetErr)
			}
			operations, exportErr := exportCellOperations(ctx, store, output, cellName, runID)
			if exportErr != nil {
				_ = store.Close()
				return fmt.Errorf("export operations for timed-out cell %+v: %w", cell, exportErr)
			}
			if checkpointErr := writeFailedCellCheckpoint(output, cellName, snapshot, cell, "timeout", budget, operations); checkpointErr != nil {
				_ = store.Close()
				return fmt.Errorf("write failed-cell checkpoint for timed-out cell %+v: %w", cell, checkpointErr)
			}
			_ = store.Close()
			return fmt.Errorf("cell %+v timed out after %s with status %s (dispatcher: %v)", cell, config.cellTimeout, snapshot.Status, dispatchErr)
		}
		if snapshot.Status != "succeeded" {
			budget, budgetErr := store.BudgetSnapshot(ctx, &runID)
			if budgetErr != nil {
				_ = store.Close()
				return fmt.Errorf("snapshot budget for failed cell %+v: %w", cell, budgetErr)
			}
			operations, exportErr := exportCellOperations(ctx, store, output, cellName, runID)
			if exportErr != nil {
				_ = store.Close()
				return fmt.Errorf("export operations for failed cell %+v: %w", cell, exportErr)
			}
			if checkpointErr := writeFailedCellCheckpoint(output, cellName, snapshot, cell, "terminal", budget, operations); checkpointErr != nil {
				_ = store.Close()
				return fmt.Errorf("write failed-cell checkpoint for failed cell %+v: %w", cell, checkpointErr)
			}
			_ = store.Close()
			return fmt.Errorf("cell %+v status %s", cell, snapshot.Status)
		}
		budget, err := store.BudgetSnapshot(ctx, &runID)
		if err != nil {
			_ = store.Close()
			return err
		}
		operations, err := exportCellOperations(ctx, store, output, cellName, runID)
		if err != nil {
			_ = store.Close()
			return err
		}
		cellResult, err := readCell(ctx, artifacts, snapshot, cell, chunkCount, budget)
		if err != nil {
			_ = store.Close()
			return err
		}
		cellResult.Operations = operations
		checkpoint := cellCheckpoint{SchemaVersion: "rag-ttc-v3-cell-evidence/v1", Profile: config.profile, ProviderProfileDigest: authority.profileDigest, GenerationModelDigest: authority.modelDigest, Cell: cellResult}
		checkpointBody, err := workflowv3.CanonicalJSON(checkpoint)
		if err != nil {
			_ = store.Close()
			return err
		}
		if err := writeFileAtomically(filepath.Join(output, "cells", cellName+".json"), append(checkpointBody, '\n'), 0o644); err != nil {
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
	if generationAuthority != nil {
		state := generationAuthority.State()
		result.GenerationAuthority = &state
	}
	body, err := workflowv3.CanonicalJSON(result)
	if err != nil {
		return err
	}
	if err = writeFileAtomically(filepath.Join(output, "evidence.json"), append(body, '\n'), 0o644); err != nil {
		return err
	}
	if err = writeCSV(filepath.Join(output, "cells.csv"), result.Cells); err != nil {
		return err
	}
	if err = writeJSONL(filepath.Join(output, "measurements.jsonl"), result.Cells); err != nil {
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
	generationProvider, embeddingProvider, err := providerIntervals(batches)
	if err != nil {
		return cellEvidence{}, err
	}
	return cellEvidence{Cell: c, RunID: string(s.RunID), PlanDigest: s.PlanDigest, Requests: len(attempts), Chunks: chunks, MakespanMicros: makespan.Microseconds(), ProviderMicros: durations, Batches: batches, Attempts: evidenceRows(attempts), EmbeddingAttempts: evidenceRows(embeddingAttempts), EmbeddingRequests: embeddingRequests, AttemptOverlapMicros: overlapIntervals(attemptIntervals(attempts), attemptIntervals(embeddingAttempts)), ProviderOverlapMicros: overlapIntervals(generationProvider, embeddingProvider), AttemptPeakActive: peakIntervals(attemptIntervals(attempts)), ProviderPeakActive: peakIntervals(generationProvider), ChunksPerSecond: float64(chunks) / makespan.Seconds(), RequestsPerSecond: float64(len(attempts)) / makespan.Seconds(), Usage: usage}, nil
}

func exportCellOperations(ctx context.Context, store *workflowv3sqlite.Store, output, cellName string, runID workflowv3.RunID) (*operationEvidence, error) {
	jsonlPath := filepath.Join(output, "operations", cellName+".jsonl")
	manifestPath := filepath.Join(output, "operations", cellName+".manifest.json")
	manifest, err := store.ExportExternalOperations(ctx, runID, jsonlPath, manifestPath)
	if err != nil {
		return nil, err
	}
	return &operationEvidence{JSONLPath: filepath.ToSlash(filepath.Join("operations", cellName+".jsonl")), ManifestPath: filepath.ToSlash(filepath.Join("operations", cellName+".manifest.json")), Manifest: manifest}, nil
}

func writeFailedCellCheckpoint(output, cellName string, snapshot workflowv3.RunSnapshot, cell workflowv3ttc.SweepCell, reason string, budget []workflowv3.BudgetProgress, operations *operationEvidence) error {
	attempts := map[string]map[string]int{}
	failureCodes := map[string]int{}
	for _, attempt := range snapshot.Attempts {
		byStatus := attempts[attempt.ResourceClass]
		if byStatus == nil {
			byStatus = map[string]int{}
			attempts[attempt.ResourceClass] = byStatus
		}
		byStatus[attempt.Status]++
		if attempt.Failure != nil && attempt.Failure.Code != "" {
			failureCodes[attempt.Failure.Code]++
		}
	}
	checkpoint := failedCellCheckpoint{SchemaVersion: "rag-ttc-v3-failed-cell-evidence/v1", Cell: cell, RunID: string(snapshot.RunID), RunStatus: snapshot.Status, Reason: reason, Attempts: attempts, FailureCodes: failureCodes, Budget: budget, Operations: operations}
	body, err := workflowv3.CanonicalJSON(checkpoint)
	if err != nil {
		return err
	}
	return writeFileAtomically(filepath.Join(output, "failures", cellName+".json"), append(body, '\n'), 0o644)
}

func evidenceRows(attempts []workflowv3.Attempt) []attemptEvidence {
	rows := make([]attemptEvidence, len(attempts))
	for i, attempt := range attempts {
		rows[i] = attemptEvidence{NodeKey: string(attempt.NodeKey), Number: attempt.Number, StartedAt: attempt.StartedAt.UTC().Format(time.RFC3339Nano), FinishedAt: attempt.FinishedAt.UTC().Format(time.RFC3339Nano)}
	}
	return rows
}

type measuredInterval struct{ start, end time.Time }

func attemptIntervals(attempts []workflowv3.Attempt) []measuredInterval {
	intervals := make([]measuredInterval, len(attempts))
	for index, attempt := range attempts {
		intervals[index] = measuredInterval{start: attempt.StartedAt, end: attempt.FinishedAt}
	}
	return intervals
}

func providerIntervals(batches []batchEvidence) ([]measuredInterval, []measuredInterval, error) {
	generation := make([]measuredInterval, 0, len(batches))
	embedding := make([]measuredInterval, 0, len(batches))
	appendMeasurement := func(target *[]measuredInterval, startedAt string, elapsedMicros int64) error {
		start, err := time.Parse(time.RFC3339Nano, startedAt)
		if err != nil || elapsedMicros < 0 {
			return fmt.Errorf("invalid provider measurement")
		}
		*target = append(*target, measuredInterval{start: start, end: start.Add(time.Duration(elapsedMicros) * time.Microsecond)})
		return nil
	}
	for _, batch := range batches {
		if err := appendMeasurement(&generation, batch.Generation.ProviderStartedAt, batch.Generation.ProviderElapsedMicros); err != nil {
			return nil, nil, err
		}
		if err := appendMeasurement(&embedding, batch.Embedding.ProviderStartedAt, batch.Embedding.ProviderElapsedMicros); err != nil {
			return nil, nil, err
		}
	}
	return generation, embedding, nil
}

func overlapIntervals(left, right []measuredInterval) int64 {
	type event struct {
		at          time.Time
		side, delta int
	}
	events := []event{}
	for _, x := range left {
		events = append(events, event{x.start, 0, 1}, event{x.end, 0, -1})
	}
	for _, x := range right {
		events = append(events, event{x.start, 1, 1}, event{x.end, 1, -1})
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
func peakIntervals(intervals []measuredInterval) int {
	type event struct {
		at    time.Time
		delta int
	}
	events := make([]event, 0, 2*len(intervals))
	for _, interval := range intervals {
		events = append(events, event{at: interval.start, delta: 1}, event{at: interval.end, delta: -1})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].at.Equal(events[j].at) {
			return events[i].delta < events[j].delta
		}
		return events[i].at.Before(events[j].at)
	})
	active, peak := 0, 0
	for _, item := range events {
		active += item.delta
		if active > peak {
			peak = active
		}
	}
	return peak
}
func writeJSONL(path string, cells []cellEvidence) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	write := func(value any) error {
		body, encodeErr := workflowv3.CanonicalJSON(value)
		if encodeErr != nil {
			return encodeErr
		}
		if _, writeErr := writer.Write(append(body, '\n')); writeErr != nil {
			return writeErr
		}
		return nil
	}
	for _, cell := range cells {
		identity := map[string]any{"chunksPerRequest": cell.Cell.ChunksPerRequest, "concurrency": cell.Cell.Concurrency, "replicate": cell.Cell.Replicate, "runId": cell.RunID}
		for _, attempt := range cell.Attempts {
			if err := write(map[string]any{"schemaVersion": "rag-ttc-v3-attempt-measurement/v1", "phase": "generation", "cell": identity, "attempt": attempt}); err != nil {
				_ = file.Close()
				return err
			}
		}
		for _, attempt := range cell.EmbeddingAttempts {
			if err := write(map[string]any{"schemaVersion": "rag-ttc-v3-attempt-measurement/v1", "phase": "embedding", "cell": identity, "attempt": attempt}); err != nil {
				_ = file.Close()
				return err
			}
		}
		for _, batch := range cell.Batches {
			if err := write(map[string]any{"schemaVersion": "rag-ttc-v3-provider-measurement/v1", "cell": identity, "batch": batch}); err != nil {
				_ = file.Close()
				return err
			}
		}
	}
	if err := writer.Flush(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func writeCSV(path string, cells []cellEvidence) error {
	f, e := os.Create(path)
	if e != nil {
		return e
	}
	w := csv.NewWriter(f)
	_ = w.Write([]string{"batch_size", "concurrency", "replicate", "generation_requests", "embedding_requests", "chunks", "makespan_us", "attempt_overlap_us", "provider_overlap_us", "attempt_peak_active", "provider_peak_active", "chunks_per_second", "requests_per_second", "input_tokens", "output_tokens", "embedding_tokens", "cost_microunits"})
	for _, c := range cells {
		_ = w.Write([]string{strconv.Itoa(c.Cell.ChunksPerRequest), strconv.Itoa(c.Cell.Concurrency), strconv.Itoa(c.Cell.Replicate), strconv.Itoa(c.Requests), strconv.Itoa(c.EmbeddingRequests), strconv.Itoa(c.Chunks), strconv.FormatInt(c.MakespanMicros, 10), strconv.FormatInt(c.AttemptOverlapMicros, 10), strconv.FormatInt(c.ProviderOverlapMicros, 10), strconv.Itoa(c.AttemptPeakActive), strconv.Itoa(c.ProviderPeakActive), strconv.FormatFloat(c.ChunksPerSecond, 'f', 6, 64), strconv.FormatFloat(c.RequestsPerSecond, 'f', 6, 64), strconv.FormatInt(c.Usage["input_tokens"], 10), strconv.FormatInt(c.Usage["output_tokens"], 10), strconv.FormatInt(c.Usage["embedding.embedding_tokens"], 10), strconv.FormatInt(c.Usage["cost_microunits"], 10)})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
