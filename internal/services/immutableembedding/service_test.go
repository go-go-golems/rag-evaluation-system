package immutableembedding

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/geppetto/pkg/embeddings"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutablechunk"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

type fakeProvider struct{}

func (fakeProvider) GetModel() embeddings.EmbeddingModel {
	return embeddings.EmbeddingModel{Name: "offline-test", Dimensions: 2}
}
func (fakeProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return []float32{float32(len(text)), 1}, nil
}
func (f fakeProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	r := make([][]float32, len(texts))
	for i, text := range texts {
		r[i], _ = f.GenerateEmbedding(ctx, text)
	}
	return r, nil
}

func TestBuildCreatesAndReusesEmbeddingSet(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	q := db.NewQueries(database)
	p := &ttcimport.Plan{Manifest: ttcimport.Manifest{SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "t", SourceExportSHA256: "x", SelectionAlgorithm: "t", KindQuotas: map[string]int{"faq": 1}}, Documents: []ttcimport.SourceDocument{{ID: "wp:1", Kind: "faq", Title: "t", SearchText: "one. two.", SearchMarkdown: "# t\none. two."}}}
	snapshot, err := corpussnapshot.Persist(context.Background(), q, p, corpussnapshot.PersistRequest{SourceByteSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	chunks, err := immutablechunk.Build(context.Background(), q, immutablechunk.Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "fixed", ChunkSize: 4})
	if err != nil {
		t.Fatal(err)
	}
	first, err := Build(context.Background(), q, Request{ChunkSetID: chunks.ChunkSetID, ProviderType: "offline-test", Provider: fakeProvider{}})
	if err != nil {
		t.Fatal(err)
	}
	if first.Reused || first.EmbeddingCount == 0 {
		t.Fatalf("first=%#v", first)
	}
	second, err := Build(context.Background(), q, Request{ChunkSetID: chunks.ChunkSetID, ProviderType: "offline-test", Provider: fakeProvider{}})
	if err != nil {
		t.Fatal(err)
	}
	if !second.Reused || second.EmbeddingSetID != first.EmbeddingSetID {
		t.Fatalf("second=%#v", second)
	}
}
