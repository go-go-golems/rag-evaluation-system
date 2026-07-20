package ragengine

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

// Prepared owns immutable, query-independent pipeline values, including indexes.
// It is safe for concurrent query execution; Close must follow the last request.
type Prepared struct {
	pipelineDigest string
	values         map[string]any
	mu             sync.Mutex
	closed         bool
}

func (p *Prepared) snapshot() (map[string]any, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.values, p.closed
}

// NewPreparedFromStaticValues constructs a prepared corpus from durable,
// already-validated static pipeline values. It intentionally rebuilds live indexes
// when the value is later opened from a PreparedCorpusStore.
func NewPreparedFromStaticValues(pipeline ragcontract.PipelineIR, values map[string]any) (*Prepared, error) {
	if err := validateCanonicalPipeline(pipeline); err != nil {
		return nil, err
	}
	if values == nil {
		return nil, fmt.Errorf("RAG_PREPARED_VALUES_NIL")
	}
	static := staticNodeIDs(pipeline)
	prepared := selectPrepared(values, static)
	if len(prepared) == 0 {
		return nil, fmt.Errorf("RAG_PREPARED_VALUES_EMPTY")
	}
	return &Prepared{pipelineDigest: mustDigest(pipeline), values: prepared}, nil
}

func (p *Prepared) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	closeIndexes(p.values)
	return nil
}

// Prepare executes only nodes that do not transitively depend on the query input.
func (e *Engine) Prepare(ctx context.Context, pipeline ragcontract.PipelineIR, corpus ragoperators.Corpus, options Options) (*Prepared, error) {
	values, _, err := e.executeStatic(ctx, pipeline, corpus, options, "")
	if err != nil {
		closeIndexes(values)
		return nil, err
	}
	prepared, err := NewPreparedFromStaticValues(pipeline, values)
	if err != nil {
		closeIndexes(values)
		return nil, err
	}
	return prepared, nil
}

// StaticInputs materializes the static predecessors of targetNodeID but never
// executes that target. It is used by durable preparation planners to obtain
// canonical chunks before scheduling provider-backed work.
func (e *Engine) StaticInputs(ctx context.Context, pipeline ragcontract.PipelineIR, corpus ragoperators.Corpus, options Options, targetNodeID string) (map[string]any, error) {
	if targetNodeID == "" {
		return nil, fmt.Errorf("RAG_STATIC_INPUT_TARGET")
	}
	values, inputs, err := e.executeStatic(ctx, pipeline, corpus, options, targetNodeID)
	closeIndexes(values)
	return inputs, err
}

func (e *Engine) executeStatic(ctx context.Context, pipeline ragcontract.PipelineIR, corpus ragoperators.Corpus, options Options, stopBefore string) (map[string]any, map[string]any, error) {
	if err := validateCanonicalPipeline(pipeline); err != nil {
		return nil, nil, err
	}
	static := staticNodeIDs(pipeline)
	values := map[string]any{"corpus/out": corpus}
	trace := &ragcontract.QueryTrace{SchemaVersion: ragcontract.TraceSchemaVersion}
	env := &ragoperators.Environment{Manifests: options.Manifests, Schemas: options.Schemas, Generator: options.Generator, Embedder: options.Embedder, Reranker: options.Reranker, Cache: options.Cache, Trace: trace, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: options.GenerationConcurrency, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint, EmitEvent: options.PreparationEvent}
	for _, node := range pipeline.Nodes {
		if !static[node.ID] {
			continue
		}
		inputs, err := staticNodeInputs(values, node)
		if err != nil {
			return values, nil, err
		}
		if node.ID == stopBefore {
			return values, inputs, nil
		}
		if err := ctx.Err(); err != nil {
			return values, nil, err
		}
		op, ok := e.Registry.Lookup(node.Operator)
		if !ok {
			return values, nil, fmt.Errorf("RAG_OPERATOR_UNAVAILABLE: %s", node.Operator.ID())
		}
		outputs, err := op.Execute(ctx, node, inputs, env)
		if err != nil {
			return values, nil, fmt.Errorf("operator %s: %w", node.Operator.ID(), err)
		}
		if err := storeStaticOutputs(values, e, node, outputs); err != nil {
			return values, nil, err
		}
	}
	if stopBefore != "" {
		return values, nil, fmt.Errorf("RAG_STATIC_INPUT_TARGET: %s", stopBefore)
	}
	return values, nil, nil
}

func staticNodeInputs(values map[string]any, node ragcontract.Node) (map[string]any, error) {
	inputs := map[string]any{}
	for _, binding := range node.Inputs {
		value, exists := values[binding.From.NodeID+"/"+binding.From.Port]
		if !exists {
			return nil, fmt.Errorf("RAG_RUNTIME_INPUT: %s.%s", binding.From.NodeID, binding.From.Port)
		}
		inputs[binding.Port] = value
	}
	return inputs, nil
}

func storeStaticOutputs(values map[string]any, engine *Engine, node ragcontract.Node, outputs map[string]any) error {
	definition, ok := engine.Definitions.Definition(node.Operator)
	if !ok {
		return fmt.Errorf("RAG_RUNTIME_DEFINITION: %s", node.Operator.ID())
	}
	allowed := map[string]bool{}
	for _, output := range definition.Outputs {
		allowed[output.Name] = true
	}
	for port, value := range outputs {
		if port == "manifest" || port == "artifact" {
			continue
		}
		if !allowed[port] {
			return fmt.Errorf("RAG_RUNTIME_OUTPUT_PORT: %s.%s", node.ID, port)
		}
		values[node.ID+"/"+port] = value
	}
	return nil
}

func validateCanonicalPipeline(pipeline ragcontract.PipelineIR) error {
	normalized, err := ragcompiler.Normalize(pipeline, nil)
	if err != nil {
		return fmt.Errorf("RAG_ENGINE_PIPELINE: %w", err)
	}
	got, _ := ragcontract.CanonicalJSON(pipeline)
	want, _ := ragcontract.CanonicalJSON(normalized)
	if string(got) != string(want) {
		return fmt.Errorf("RAG_ENGINE_PIPELINE_NONCANONICAL")
	}
	return nil
}

func mustDigest(value any) string {
	digest, _ := ragcontract.Digest(value)
	return digest
}
