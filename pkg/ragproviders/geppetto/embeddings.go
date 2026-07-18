package geppetto

import (
	"context"
	"fmt"
	"math"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
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
func (e *Embedder) Embed(ctx context.Context, model string, texts []string) ([][]float64, ragoperators.Usage, error) {
	if e == nil || e.provider == nil {
		return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBEDDER_UNAVAILABLE")
	}
	if model != e.model {
		return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBEDDER_MODEL_MISMATCH")
	}
	vectors, err := e.provider.GenerateBatchEmbeddings(ctx, texts)
	if err != nil {
		return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBED: %w", err)
	}
	if len(vectors) != len(texts) {
		return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBED_COUNT")
	}
	result := make([][]float64, len(vectors))
	for i, vector := range vectors {
		if len(vector) != e.dimensions {
			return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBED_DIMENSIONS")
		}
		result[i] = make([]float64, len(vector))
		for j, value := range vector {
			if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
				return nil, ragoperators.Usage{}, fmt.Errorf("RAG_GEPPETTO_EMBED_NONFINITE")
			}
			result[i][j] = float64(value)
		}
	}
	return result, ragoperators.Usage{}, nil
}
