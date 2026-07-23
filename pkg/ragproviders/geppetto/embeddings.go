package geppetto

import (
	"context"
	"fmt"
	"math"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
)

type Embedder struct {
	provider   embeddings.Provider
	model      string
	dimensions int
}

var _ ragoperators.Embedder = (*Embedder)(nil)

func NewEmbedder(provider embeddings.Provider, model string, dimensions int) (*Embedder, error) {
	if provider == nil || model == "" || dimensions < 1 {
		return nil, fmt.Errorf("RAG_GEPPETTO_EMBEDDER_CONFIG")
	}
	return &Embedder{provider: provider, model: model, dimensions: dimensions}, nil
}

type preparedEmbedding struct {
	embedder *Embedder
	texts    []string
}

func (e *Embedder) Embed(ctx context.Context, model string, texts []string) ([][]float64, ragoperators.Usage, error) {
	prepared, err := e.PrepareEmbedding(model, texts)
	if err != nil {
		return nil, ragoperators.Usage{}, err
	}
	return prepared.Execute(ctx)
}
func (e *Embedder) PrepareEmbedding(model string, texts []string) (ragworkflowops.PreparedEmbedding, error) {
	if e == nil || e.provider == nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_EMBEDDER_UNAVAILABLE")
	}
	if model != e.model {
		return nil, fmt.Errorf("RAG_GEPPETTO_EMBEDDER_MODEL_MISMATCH")
	}
	if len(texts) == 0 {
		return nil, fmt.Errorf("RAG_GEPPETTO_EMBED_INPUT")
	}
	return &preparedEmbedding{embedder: e, texts: append([]string(nil), texts...)}, nil
}
func (p *preparedEmbedding) Execute(ctx context.Context) ([][]float64, ragoperators.Usage, error) {
	vectors, err := p.embedder.provider.GenerateBatchEmbeddings(ctx, p.texts)
	if err != nil {
		return nil, ragoperators.Usage{}, classifyProviderError(err)
	}
	if len(vectors) != len(p.texts) {
		return nil, ragoperators.Usage{}, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GEPPETTO_EMBED_COUNT")
	}
	result := make([][]float64, len(vectors))
	for i, vector := range vectors {
		if len(vector) != p.embedder.dimensions {
			return nil, ragoperators.Usage{}, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GEPPETTO_EMBED_DIMENSIONS")
		}
		result[i] = make([]float64, len(vector))
		for j, value := range vector {
			if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
				return nil, ragoperators.Usage{}, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GEPPETTO_EMBED_NONFINITE")
			}
			result[i][j] = float64(value)
		}
	}
	return result, ragoperators.Usage{}, nil
}
