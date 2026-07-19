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
	normalized, err := ragcompiler.Normalize(pipeline, nil)
	if err != nil {
		return nil, fmt.Errorf("RAG_ENGINE_PIPELINE: %w", err)
	}
	got, _ := ragcontract.CanonicalJSON(pipeline)
	want, _ := ragcontract.CanonicalJSON(normalized)
	if string(got) != string(want) {
		return nil, fmt.Errorf("RAG_ENGINE_PIPELINE_NONCANONICAL")
	}
	static := staticNodeIDs(pipeline)
	values := map[string]any{"corpus/out": corpus}
	trace := &ragcontract.QueryTrace{SchemaVersion: ragcontract.TraceSchemaVersion}
	env := &ragoperators.Environment{Manifests: options.Manifests, Schemas: options.Schemas, Generator: options.Generator, Embedder: options.Embedder, Reranker: options.Reranker, Cache: options.Cache, Trace: trace, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: options.GenerationConcurrency, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint}
	failed := true
	defer func() {
		if failed {
			closeIndexes(values)
		}
	}()
	for _, node := range pipeline.Nodes {
		if !static[node.ID] {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		op, ok := e.Registry.Lookup(node.Operator)
		if !ok {
			return nil, fmt.Errorf("RAG_OPERATOR_UNAVAILABLE: %s", node.Operator.ID())
		}
		inputs := map[string]any{}
		for _, binding := range node.Inputs {
			value, exists := values[binding.From.NodeID+"/"+binding.From.Port]
			if !exists {
				return nil, fmt.Errorf("RAG_RUNTIME_INPUT: %s.%s", binding.From.NodeID, binding.From.Port)
			}
			inputs[binding.Port] = value
		}
		outputs, err := op.Execute(ctx, node, inputs, env)
		if err != nil {
			return nil, fmt.Errorf("operator %s: %w", node.Operator.ID(), err)
		}
		definition, ok := e.Definitions.Definition(node.Operator)
		if !ok {
			return nil, fmt.Errorf("RAG_RUNTIME_DEFINITION: %s", node.Operator.ID())
		}
		allowed := map[string]bool{}
		for _, output := range definition.Outputs {
			allowed[output.Name] = true
		}
		for port, value := range outputs {
			if port == "manifest" {
				continue
			}
			if port == "artifact" {
				continue
			}
			if !allowed[port] {
				return nil, fmt.Errorf("RAG_RUNTIME_OUTPUT_PORT: %s.%s", node.ID, port)
			}
			values[node.ID+"/"+port] = value
		}
	}
	prepared := &Prepared{pipelineDigest: mustDigest(pipeline), values: selectPrepared(values, static)}
	failed = false
	return prepared, nil
}

func mustDigest(value any) string {
	digest, _ := ragcontract.Digest(value)
	return digest
}
