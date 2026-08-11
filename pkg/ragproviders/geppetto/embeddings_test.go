package geppetto

import (
	"context"
	"errors"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
)

type fakeEmbedProvider struct {
	vectors [][]float32
	err     error
	block   bool
}

func (p *fakeEmbedProvider) GenerateEmbedding(context.Context, string) ([]float32, error) {
	if len(p.vectors) == 0 {
		return nil, p.err
	}
	return p.vectors[0], p.err
}
func (p *fakeEmbedProvider) GenerateBatchEmbeddings(ctx context.Context, _ []string) ([][]float32, error) {
	if p.block {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return p.vectors, p.err
}
func (p *fakeEmbedProvider) GetModel() embeddings.EmbeddingModel {
	return embeddings.EmbeddingModel{Name: "test", Dimensions: 2}
}

func TestEmbedderPreservesOrderAndRejectsInvalidVectors(t *testing.T) {
	adapter, err := NewEmbedder(&fakeEmbedProvider{vectors: [][]float32{{1, 2}, {3, 4}}}, "embedding-primary", 2)
	if err != nil {
		t.Fatal(err)
	}
	vectors, _, err := adapter.Embed(context.Background(), "embedding-primary", []string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	if vectors[0][0] != 1 || vectors[1][0] != 3 {
		t.Fatalf("vectors=%#v", vectors)
	}

	cases := []struct {
		name    string
		vectors [][]float32
		want    string
	}{
		{"count", [][]float32{{1, 2}}, "RAG_GEPPETTO_EMBED_COUNT"},
		{"dimensions", [][]float32{{1}, {2, 3}}, "RAG_GEPPETTO_EMBED_DIMENSIONS"},
		{"nonfinite", [][]float32{{float32(math.NaN()), 2}, {3, 4}}, "RAG_GEPPETTO_EMBED_NONFINITE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, err := NewEmbedder(&fakeEmbedProvider{vectors: tc.vectors}, "embedding-primary", 2)
			if err != nil {
				t.Fatal(err)
			}
			_, _, err = a.Embed(context.Background(), "embedding-primary", []string{"a", "b"})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%v want=%s", err, tc.want)
			}
		})
	}
}

func TestEmbedderRejectsModelMismatchAndPreservesCancellation(t *testing.T) {
	adapter, err := NewEmbedder(&fakeEmbedProvider{vectors: [][]float32{{1, 2}}}, "embedding-primary", 2)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = adapter.Embed(context.Background(), "other", []string{"a"})
	if err == nil || !strings.Contains(err.Error(), "RAG_GEPPETTO_EMBEDDER_MODEL_MISMATCH") {
		t.Fatalf("error=%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	adapter, err = NewEmbedder(&fakeEmbedProvider{block: true}, "embedding-primary", 2)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = adapter.Embed(ctx, "embedding-primary", []string{"a"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v want cancellation", err)
	}
}

func newLiveEmbeddingSettings(baseURL, model string) (*settings.InferenceSettings, error) {
	in, err := settings.NewInferenceSettings()
	if err != nil {
		return nil, err
	}
	in.Embeddings.Type = "ollama"
	in.Embeddings.Engine = model
	in.Embeddings.Dimensions = 768
	in.Embeddings.BaseURLs = map[string]string{"ollama-base-url": baseURL}
	in.API.BaseUrls["ollama-base-url"] = baseURL
	in.API.AllowHTTP["ollama"] = true
	in.API.AllowLocalNetworks["ollama"] = true
	return in, nil
}

func TestEmbedderLiveOllama(t *testing.T) {
	if os.Getenv("RAG_EMBEDDING_LIVE_TEST") != "1" {
		t.Skip("set RAG_EMBEDDING_LIVE_TEST=1 to run")
	}
	baseURL := os.Getenv("RAG_EMBEDDING_LIVE_BASE_URL")
	if baseURL == "" {
		t.Fatal("RAG_EMBEDDING_LIVE_BASE_URL is required")
	}
	model := os.Getenv("RAG_EMBEDDING_LIVE_MODEL")
	if model == "" {
		model = "nomic-embed-text:latest"
	}
	in, err := newLiveEmbeddingSettings(baseURL, model)
	if err != nil {
		t.Fatal(err)
	}
	provider, err := embeddings.NewSettingsFactoryFromInferenceSettings(in).NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	adapter, err := NewEmbedder(provider, model, 768)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	vectors, _, err := adapter.Embed(ctx, model, []string{"A payroll adjustment corrects wages."})
	if err != nil {
		t.Fatal(err)
	}
	if len(vectors) != 1 || len(vectors[0]) != 768 {
		t.Fatalf("vectors=%d dims=%d", len(vectors), len(vectors[0]))
	}
}
