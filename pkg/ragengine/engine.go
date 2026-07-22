package ragengine

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type Observer interface {
	Event(context.Context, ragoperators.Event) error
	Trace(context.Context, ragcontract.QueryTrace) error
	Metric(context.Context, ragoperators.Metric) error
	Artifact(context.Context, ragoperators.Artifact) error
}
type NopObserver struct{}

func (NopObserver) Event(context.Context, ragoperators.Event) error       { return nil }
func (NopObserver) Trace(context.Context, ragcontract.QueryTrace) error   { return nil }
func (NopObserver) Metric(context.Context, ragoperators.Metric) error     { return nil }
func (NopObserver) Artifact(context.Context, ragoperators.Artifact) error { return nil }

type Options struct {
	Manifests                     ragoperators.ManifestResolver
	Schemas                       ragoperators.OutputSchemaValidator
	Generator                     ragoperators.TextGenerator
	Embedder                      ragoperators.Embedder
	Reranker                      ragoperators.Reranker
	Cache                         ragoperators.Cache
	GenerationConcurrency         int
	GenerationSettingsFingerprint string
	GeneratorFingerprint          string
	RerankerFingerprint           string
	PreparedCorpusDigest          string
	QueryCheckpoints              QueryCheckpointStore
	PreparedStore                 PreparedCorpusStore
	EmbeddingFingerprint          string
	Prepared                      *Prepared
	// PreparationEvent emits runtime-only progress for static preparation.
	PreparationEvent func(context.Context, ragoperators.Event) error
}
type Result struct {
	Traces    []ragcontract.QueryTrace   `json:"traces"`
	Answers   []ragoperators.Answer      `json:"answers,omitempty"`
	Metrics   []ragoperators.Metric      `json:"metrics"`
	Artifacts []ragoperators.Artifact    `json:"artifacts"`
	Failures  []ragcontract.FailureTrace `json:"failures"`
}
type Engine struct {
	Registry    *ragoperators.Registry
	Definitions *ragcompiler.Registry
}

func New(registry *ragoperators.Registry) *Engine {
	if registry == nil {
		registry = ragoperators.NativeRegistry()
	}
	return &Engine{Registry: registry, Definitions: ragcompiler.BuiltinRegistry()}
}
func (e *Engine) Execute(ctx context.Context, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, dataset ragoperators.EvaluationDataset, observer Observer, options Options) (*Result, error) {
	if observer == nil {
		observer = NopObserver{}
	}
	if execution.SchemaVersion != ragcontract.ExecutionSchemaVersion {
		return nil, fmt.Errorf("RAG_ENGINE_SCHEMA: %s", execution.SchemaVersion)
	}
	if execution.Pipeline.SchemaVersion != ragcontract.PipelineSchemaVersion {
		return nil, fmt.Errorf("RAG_ENGINE_PIPELINE_SCHEMA: %s", execution.Pipeline.SchemaVersion)
	}
	if options.Prepared == nil {
		normalized, err := ragcompiler.Normalize(execution.Pipeline, nil)
		if err != nil {
			return nil, fmt.Errorf("RAG_ENGINE_PIPELINE: %w", err)
		}
		gotPipeline, _ := ragcontract.CanonicalJSON(execution.Pipeline)
		wantPipeline, _ := ragcontract.CanonicalJSON(normalized)
		if string(gotPipeline) != string(wantPipeline) {
			return nil, fmt.Errorf("RAG_ENGINE_PIPELINE_NONCANONICAL")
		}
	}
	identity := execution
	identity.CellID = ""
	cellID, err := ragcontract.Digest(identity)
	if err != nil {
		return nil, err
	}
	if execution.CellID != cellID {
		return nil, fmt.Errorf("RAG_ENGINE_CELL_ID: got %s want %s", execution.CellID, cellID)
	}
	result := &Result{}
	irArtifact, _ := json.Marshal(execution.Pipeline)
	pipelineArtifact := ragoperators.Artifact{Role: "pipeline-ir", Kind: "rag-pipeline-ir", Name: "pipeline-ir.json", SchemaVersion: ragcontract.PipelineSchemaVersion, MediaType: "application/json", Data: irArtifact}
	result.Artifacts = append(result.Artifacts, pipelineArtifact)
	if err := observer.Artifact(ctx, pipelineArtifact); err != nil {
		return nil, err
	}
	staticNodes := staticNodeIDs(execution.Pipeline)
	var prepared map[string]any
	ownsPrepared := options.Prepared == nil
	if options.Prepared != nil {
		var closed bool
		prepared, closed = options.Prepared.snapshot()
		if closed {
			return nil, fmt.Errorf("RAG_ENGINE_PREPARED_CLOSED")
		}
		if options.Prepared.pipelineDigest != mustDigest(execution.Pipeline) {
			return nil, fmt.Errorf("RAG_ENGINE_PREPARED_PIPELINE")
		}
	}
	if ownsPrepared {
		defer func() { closeIndexes(prepared) }()
	}
	executionDigest := mustDigest(execution)
	evaluationPolicyDigest := mustDigest(struct {
		Measures        []ragcontract.Measure
		RelevanceTarget string
		DatasetSplit    string
	}{execution.Measures, execution.Dataset.RelevanceTarget, execution.Dataset.Split})
	for queryIndex, query := range dataset.Queries {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		queryDigest, err := ragcontract.Digest(query.Text)
		if err != nil {
			return result, err
		}
		checkpointIdentity := QueryCheckpointIdentity{
			SchemaVersion: queryCheckpointSchemaVersion, ExecutionDigest: executionDigest,
			QueryID: query.ID, QueryTextDigest: queryDigest,
			PreparedCorpusDigest:   options.PreparedCorpusDigest,
			GeneratorFingerprint:   options.GeneratorFingerprint,
			RerankerFingerprint:    options.RerankerFingerprint,
			EvaluationPolicyDigest: evaluationPolicyDigest,
		}
		if options.QueryCheckpoints != nil {
			checkpoint, found, err := options.QueryCheckpoints.Get(ctx, checkpointIdentity)
			if err != nil {
				return result, fmt.Errorf("RAG_QUERY_CHECKPOINT_GET: %w", err)
			}
			if found {
				result.Traces = append(result.Traces, checkpoint.Trace)
				result.Metrics = append(result.Metrics, checkpoint.Metrics...)
				result.Artifacts = append(result.Artifacts, checkpoint.Artifacts...)
				if checkpoint.Answer != nil {
					result.Answers = append(result.Answers, *checkpoint.Answer)
				}
				if err := observer.Event(ctx, event("rag.query.progress/v1", map[string]any{"queryId": query.ID, "queryIndex": queryIndex + 1, "queryCount": len(dataset.Queries), "state": "resumed", "resumed": true})); err != nil {
					return result, err
				}
				if err := observer.Trace(ctx, checkpoint.Trace); err != nil {
					return result, err
				}
				for _, metric := range checkpoint.Metrics {
					if err := observer.Metric(ctx, metric); err != nil {
						return result, err
					}
				}
				for _, artifact := range checkpoint.Artifacts {
					if err := observer.Artifact(ctx, artifact); err != nil {
						return result, err
					}
				}
				continue
			}
		}
		skip := map[string]bool{}
		if queryIndex > 0 || options.Prepared != nil {
			skip = staticNodes
		}
		if err := observer.Event(ctx, event("rag.query.progress/v1", map[string]any{"queryId": query.ID, "queryIndex": queryIndex + 1, "queryCount": len(dataset.Queries), "state": "started"})); err != nil {
			return result, err
		}
		queryStarted := time.Now()
		trace, metrics, artifacts, values, err := e.executeQuery(ctx, execution, corpus, query, options, observer, prepared, skip)
		if queryIndex == 0 && options.Prepared == nil {
			prepared = selectPrepared(values, staticNodes)
		}
		result.Metrics = append(result.Metrics, metrics...)
		result.Artifacts = append(result.Artifacts, artifacts...)
		if answer, ok := answerFromValues(values); ok {
			result.Answers = append(result.Answers, answer)
		}
		if err != nil {
			failure := ragcontract.FailureTrace{Code: "RAG_QUERY_FAILED", Path: "$.queries[" + query.ID + "]", Message: err.Error()}
			result.Failures = append(result.Failures, failure)
			trace.Failures = append(trace.Failures, failure)
			result.Traces = append(result.Traces, trace)
			_ = observer.Trace(context.WithoutCancel(ctx), trace)
			_ = observer.Event(context.WithoutCancel(ctx), event("rag.query.failed", map[string]any{"queryId": query.ID, "error": err.Error()}))
			return result, err
		}
		var answerCheckpoint *ragoperators.Answer
		if answer, ok := answerFromValues(values); ok {
			answerCopy := answer
			answerCheckpoint = &answerCopy
		}
		if options.QueryCheckpoints != nil {
			checkpoint := QueryCheckpoint{SchemaVersion: queryCheckpointSchemaVersion, Identity: checkpointIdentity, Trace: trace, Metrics: metrics, Artifacts: artifacts, Answer: answerCheckpoint}
			if err := options.QueryCheckpoints.Put(ctx, checkpoint); err != nil {
				return result, fmt.Errorf("RAG_QUERY_CHECKPOINT_PUT: %w", err)
			}
		}
		result.Traces = append(result.Traces, trace)
		if err := observer.Event(ctx, event("rag.query.progress/v1", map[string]any{"queryId": query.ID, "queryIndex": queryIndex + 1, "queryCount": len(dataset.Queries), "state": "completed", "elapsedMilliseconds": time.Since(queryStarted).Milliseconds()})); err != nil {
			return result, err
		}
		if err := observer.Trace(ctx, trace); err != nil {
			return result, err
		}
		for _, metric := range metrics {
			if err := observer.Metric(ctx, metric); err != nil {
				return result, err
			}
		}
	}
	traceData, _ := json.Marshal(result.Traces)
	traceArtifact := ragoperators.Artifact{Role: "query-traces", Kind: "rag-query-traces", Name: "query-traces.json", SchemaVersion: "rag-query-trace-bundle/v2", MediaType: "application/json", Data: traceData}
	result.Artifacts = append(result.Artifacts, traceArtifact)
	if err := observer.Artifact(ctx, traceArtifact); err != nil {
		return result, err
	}
	return result, nil
}
func (e *Engine) executeQuery(ctx context.Context, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, query ragoperators.Query, options Options, observer Observer, prepared map[string]any, skip map[string]bool) (ragcontract.QueryTrace, []ragoperators.Metric, []ragoperators.Artifact, map[string]any, error) {
	digest, _ := ragcontract.Digest(query.Text)
	trace := ragcontract.QueryTrace{SchemaVersion: ragcontract.TraceSchemaVersion, Query: ragcontract.QueryInputTrace{ID: query.ID, TextDigest: digest, DatasetSplit: execution.Dataset.Split}, Operators: []ragcontract.OperatorTrace{}, Channels: []ragcontract.ChannelTrace{}, Collapses: []ragcontract.CollapseTrace{}, Results: []ragcontract.ResultTrace{}, Timing: ragcontract.TimingTrace{ByOperator: map[string]int64{}}, Usage: ragcontract.UsageTrace{}, Failures: []ragcontract.FailureTrace{}}
	env := &ragoperators.Environment{Manifests: options.Manifests, Schemas: options.Schemas, Generator: options.Generator, Embedder: options.Embedder, Reranker: options.Reranker, Cache: options.Cache, Trace: &trace, CurrentQuery: query, QueryText: query.Text, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: options.GenerationConcurrency, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint, EmitEvent: observer.Event}
	values := map[string]any{"corpus/out": corpus, "query/out": query}
	for key, value := range prepared {
		values[key] = value
	}
	artifacts := []ragoperators.Artifact{}
	started := time.Now()
	for _, node := range execution.Pipeline.Nodes {
		if skip[node.ID] {
			continue
		}
		if err := ctx.Err(); err != nil {
			return trace, nil, artifacts, values, err
		}
		operator, ok := e.Registry.Lookup(node.Operator)
		if !ok {
			return trace, nil, artifacts, values, fmt.Errorf("RAG_OPERATOR_UNAVAILABLE: %s", node.Operator.ID())
		}
		inputs := map[string]any{}
		for _, binding := range node.Inputs {
			value, exists := values[binding.From.NodeID+"/"+binding.From.Port]
			if !exists {
				return trace, nil, artifacts, values, fmt.Errorf("RAG_RUNTIME_INPUT: %s.%s", binding.From.NodeID, binding.From.Port)
			}
			inputs[binding.Port] = value
		}
		nodeStarted := time.Now()
		outputs, err := operator.Execute(ctx, node, inputs, env)
		duration := time.Since(nodeStarted).Milliseconds()
		status := "succeeded"
		if err != nil {
			status = "failed"
		}
		trace.Operators = append(trace.Operators, ragcontract.OperatorTrace{NodeID: node.ID, Operator: node.Operator, Status: status, InputCount: len(inputs), OutputCount: len(outputs), DurationMilliseconds: duration})
		trace.Timing.ByOperator[node.ID] = duration
		if err != nil {
			return trace, nil, artifacts, values, fmt.Errorf("operator %s: %w", node.Operator.ID(), err)
		}
		if len(outputs) == 0 {
			return trace, nil, artifacts, values, fmt.Errorf("RAG_RUNTIME_OUTPUT_EMPTY: %s", node.ID)
		}
		definition, known := e.Definitions.Definition(node.Operator)
		if !known {
			return trace, nil, artifacts, values, fmt.Errorf("RAG_RUNTIME_DEFINITION: %s", node.Operator.ID())
		}
		allowedOutputs := map[string]bool{}
		for _, output := range definition.Outputs {
			allowedOutputs[output.Name] = true
		}
		for port, value := range outputs {
			if port == "manifest" {
				continue
			}
			if port == "artifact" {
				artifact, ok := value.(ragoperators.Artifact)
				if !ok {
					return trace, nil, artifacts, values, fmt.Errorf("RAG_RUNTIME_ARTIFACT: %s", node.ID)
				}
				ext := filepath.Ext(artifact.Name)
				artifact.Name = artifact.Name[:len(artifact.Name)-len(ext)] + "-" + query.ID + ext
				artifacts = append(artifacts, artifact)
				if err := observer.Artifact(ctx, artifact); err != nil {
					return trace, nil, artifacts, values, err
				}
				continue
			}
			if !allowedOutputs[port] {
				return trace, nil, artifacts, values, fmt.Errorf("RAG_RUNTIME_OUTPUT_PORT: %s.%s", node.ID, port)
			}
			values[node.ID+"/"+port] = value
			if index, ok := value.(*ragoperators.MultiIndex); ok {
				manifest, _ := json.Marshal(index.Manifest)
				artifact := ragoperators.Artifact{Role: "index-manifest", Kind: "rag-index-manifest", Name: "index-manifest-" + query.ID + ".json", SchemaVersion: ragcontract.IndexManifestSchema, MediaType: "application/json", Data: manifest}
				artifacts = append(artifacts, artifact)
				if err := observer.Artifact(ctx, artifact); err != nil {
					return trace, nil, artifacts, values, err
				}
				records := ragoperators.Artifact{Role: "index-records", Kind: "rag-index-records", Name: "index-records-" + query.ID + ".json", SchemaVersion: "rag-index-records/v1", MediaType: "application/json", Data: index.Artifact}
				artifacts = append(artifacts, records)
				if err := observer.Artifact(ctx, records); err != nil {
					return trace, nil, artifacts, values, err
				}
			}
		}
		if err := observer.Event(ctx, event("rag.operator.completed", map[string]any{"queryId": query.ID, "nodeId": node.ID, "operator": node.Operator.ID()})); err != nil {
			return trace, nil, artifacts, values, err
		}
	}
	trace.Timing.TotalMilliseconds = time.Since(started).Milliseconds()
	trace.Usage = ragcontract.UsageTrace{InputTokens: env.Usage.InputTokens, OutputTokens: env.Usage.OutputTokens, EmbeddingTokens: env.Usage.EmbeddingTokens, ProviderCost: env.Usage.Cost}
	var evidence []ragoperators.Evidence
	var answer *ragoperators.Answer
	for i := len(execution.Pipeline.Nodes) - 1; i >= 0 && evidence == nil; i-- {
		if value, exists := values[execution.Pipeline.Nodes[i].ID+"/evidence"]; exists {
			evidence, _ = value.([]ragoperators.Evidence)
		}
	}
	for _, output := range execution.Pipeline.Outputs {
		value := values[output.From.NodeID+"/"+output.From.Port]
		switch typed := value.(type) {
		case []ragoperators.Evidence:
			evidence = typed
		case ragoperators.Answer:
			answer = &typed
		}
	}
	for _, item := range evidence {
		fusion := item.Score
		score := ragcontract.ResultScore{Fusion: &fusion, Reranker: item.RerankerScore}
		matched := make([]ragcontract.MatchedRepresentation, len(item.Matched))
		for index, record := range item.Matched {
			matched[index] = ragcontract.MatchedRepresentation{ID: record.Representation.Record.ID, Kind: record.Representation.Record.Kind, Channel: record.Channel, Rank: record.Rank}
		}
		trace.Results = append(trace.Results, ragcontract.ResultTrace{Rank: item.Rank, Collapse: item.Collapse, MatchedRepresentations: matched, Evidence: ragcontract.EvidenceIdentity{ChunkID: item.Chunk.Record.ID, Digest: item.Chunk.Record.TextDigest, Citation: item.Chunk.Record.Citation}, Scores: score})
	}
	trace.Relevance = &ragcontract.RelevanceTrace{Target: execution.Dataset.RelevanceTarget, ExpectedIDs: append([]string(nil), query.RelevantIDs...), Grades: query.Grades, Measures: map[string]json.RawMessage{}}
	storage := int64(0)
	for _, artifact := range artifacts {
		storage += int64(len(artifact.Data))
	}
	timing := map[string]int64{"query": trace.Timing.TotalMilliseconds}
	for key, value := range trace.Timing.ByOperator {
		timing[key] = value
	}
	metrics := ragoperators.EvaluateForTarget(query, evidence, answer, execution.Measures, timing, env.Usage, trace.Failures, storage, execution.Dataset.RelevanceTarget)
	return trace, metrics, artifacts, values, nil
}
func staticNodeIDs(pipeline ragcontract.PipelineIR) map[string]bool {
	dynamic := map[string]bool{"query": true}
	static := map[string]bool{}
	for _, node := range pipeline.Nodes {
		isDynamic := false
		for _, input := range node.Inputs {
			if dynamic[input.From.NodeID] {
				isDynamic = true
				break
			}
		}
		if isDynamic {
			dynamic[node.ID] = true
		} else {
			static[node.ID] = true
		}
	}
	return static
}
func selectPrepared(values map[string]any, static map[string]bool) map[string]any {
	prepared := map[string]any{}
	for key, value := range values {
		nodeID := key
		for index, character := range key {
			if character == '/' {
				nodeID = key[:index]
				break
			}
		}
		if nodeID == "corpus" || static[nodeID] {
			prepared[key] = value
		}
	}
	return prepared
}
func closeIndexes(values map[string]any) {
	seen := map[*ragoperators.MultiIndex]bool{}
	for _, value := range values {
		if index, ok := value.(*ragoperators.MultiIndex); ok && !seen[index] {
			seen[index] = true
			_ = index.Close()
		}
	}
}
func answerFromValues(values map[string]any) (ragoperators.Answer, bool) {
	for _, value := range values {
		if answer, ok := value.(ragoperators.Answer); ok {
			return answer, true
		}
	}
	return ragoperators.Answer{}, false
}

func event(kind string, value any) ragoperators.Event {
	raw, _ := json.Marshal(value)
	return ragoperators.Event{Type: kind, Payload: raw}
}
func SortMetrics(values []ragoperators.Metric) {
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
}
