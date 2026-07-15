package raglab

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutablechunk"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
)

func TestSQLiteCatalogResolvesImmutableLineage(t *testing.T) {
	database, err := db.OpenDB(filepath.Join(t.TempDir(), "rag-eval.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = database.Close() }()
	if err := db.Migrate(database); err != nil {
		t.Fatal(err)
	}
	queries := db.NewQueries(database)
	ctx := context.Background()
	plan := &ttcimport.Plan{Manifest: ttcimport.Manifest{
		SchemaVersion: ttcimport.ManifestSchemaVersion, SnapshotName: "fixture", SourceExportSHA256: "fixture", SelectionAlgorithm: "fixture", KindQuotas: map[string]int{"faq": 1},
	}, Documents: []ttcimport.SourceDocument{{ID: "wp:1", Kind: "faq", Title: "Fixture", URL: "https://example.test/1", SearchText: "Blue cypress fixture text."}}}
	snapshot, err := corpussnapshot.Persist(ctx, queries, plan, corpussnapshot.PersistRequest{SourceByteSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	chunks, err := immutablechunk.Build(ctx, queries, immutablechunk.Request{CorpusSnapshotID: snapshot.SnapshotID, Strategy: "fixed", ChunkSize: 200})
	if err != nil {
		t.Fatal(err)
	}
	bm25, err := immutableretrieval.BuildBM25(ctx, queries, immutableretrieval.BM25BuildRequest{ChunkSetID: chunks.ChunkSetID, ArtifactRoot: filepath.Join(t.TempDir(), "bm25")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO embedding_plans (id,schema_version,provider_type,model,dimensions,normalization,implementation_version) VALUES (?,?,?,?,?,?,?)`, "embedding-plan", "fixture/v1", "fixture", "fixture", 768, "unit", "unit"); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO embedding_sets (id,chunk_set_id,embedding_plan_id,manifest_json,embedding_count) VALUES (?,?,?,?,?)`, "embeddings", chunks.ChunkSetID, "embedding-plan", `{}`, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := database.Exec(`INSERT INTO evaluation_datasets (id,schema_version,corpus_snapshot_id,status,manifest_json,query_count) VALUES (?,?,?,?,?,?)`, "evaluation", "fixture/v1", snapshot.SnapshotID, "candidate", `{}`, 1); err != nil {
		t.Fatal(err)
	}
	catalog := NewSQLiteCatalog(queries)
	cases := []struct {
		ref    ArtifactRef
		assert func(ArtifactMetadata) bool
	}{
		{CorpusSnapshot(snapshot.SnapshotID), func(metadata ArtifactMetadata) bool { return metadata.Ref.ID == snapshot.SnapshotID }},
		{ChunkSet(chunks.ChunkSetID), func(metadata ArtifactMetadata) bool { return metadata.CorpusSnapshotID == snapshot.SnapshotID }},
		{BM25Index(bm25.ArtifactID), func(metadata ArtifactMetadata) bool { return metadata.ChunkSetID == chunks.ChunkSetID }},
		{EmbeddingSet("embeddings"), func(metadata ArtifactMetadata) bool {
			return metadata.ChunkSetID == chunks.ChunkSetID && metadata.Dimensions == 768
		}},
		{EvaluationDataset("evaluation"), func(metadata ArtifactMetadata) bool {
			return metadata.CorpusSnapshotID == snapshot.SnapshotID && metadata.Status == "candidate"
		}},
	}
	for _, tc := range cases {
		metadata, err := catalog.LookupArtifact(ctx, tc.ref)
		if err != nil {
			t.Fatalf("lookup %s:%s: %v", tc.ref.Kind, tc.ref.ID, err)
		}
		if !tc.assert(metadata) {
			t.Fatalf("metadata for %s:%s = %#v", tc.ref.Kind, tc.ref.ID, metadata)
		}
	}
}
