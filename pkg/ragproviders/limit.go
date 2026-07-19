package ragproviders

import (
	"context"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

// providerConcurrencyLimit returns a safe default for every real provider.
// Explicit host configuration may raise it only after an operator has qualified
// that provider. Generation workers use the same limit, while wrappers ensure
// every provider call (including answer generation) is bounded.
func providerConcurrencyLimit(spec ProviderSpec) int {
	if spec.Concurrency.MaxInFlight > 0 {
		return spec.Concurrency.MaxInFlight
	}
	return 1
}

type limitedGenerator struct {
	inner ragoperators.TextGenerator
	slots chan struct{}
}

func newLimitedGenerator(inner ragoperators.TextGenerator, limit int) (ragoperators.TextGenerator, error) {
	if inner == nil || limit < 1 {
		return nil, fmt.Errorf("RAG_PROVIDER_CONCURRENCY")
	}
	return &limitedGenerator{inner: inner, slots: make(chan struct{}, limit)}, nil
}

func (g *limitedGenerator) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	select {
	case g.slots <- struct{}{}:
		defer func() { <-g.slots }()
	case <-ctx.Done():
		return ragoperators.GenerationResult{}, ctx.Err()
	}
	return g.inner.Generate(ctx, request)
}

type limitedEmbedder struct {
	inner ragoperators.Embedder
	slots chan struct{}
}

func newLimitedEmbedder(inner ragoperators.Embedder, limit int) (ragoperators.Embedder, error) {
	if inner == nil || limit < 1 {
		return nil, fmt.Errorf("RAG_PROVIDER_CONCURRENCY")
	}
	return &limitedEmbedder{inner: inner, slots: make(chan struct{}, limit)}, nil
}

func (e *limitedEmbedder) Embed(ctx context.Context, model string, texts []string) ([][]float64, ragoperators.Usage, error) {
	select {
	case e.slots <- struct{}{}:
		defer func() { <-e.slots }()
	case <-ctx.Done():
		return nil, ragoperators.Usage{}, ctx.Err()
	}
	return e.inner.Embed(ctx, model, texts)
}

type limitedReranker struct {
	inner ragoperators.Reranker
	slots chan struct{}
}

func newLimitedReranker(inner ragoperators.Reranker, limit int) (ragoperators.Reranker, error) {
	if inner == nil || limit < 1 {
		return nil, fmt.Errorf("RAG_PROVIDER_CONCURRENCY")
	}
	return &limitedReranker{inner: inner, slots: make(chan struct{}, limit)}, nil
}

func (r *limitedReranker) Rerank(ctx context.Context, request ragoperators.RerankRequest) ([]ragoperators.RerankScore, error) {
	select {
	case r.slots <- struct{}{}:
		defer func() { <-r.slots }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return r.inner.Rerank(ctx, request)
}
