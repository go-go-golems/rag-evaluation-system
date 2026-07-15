package experimentspec

import "testing"

func TestFingerprintNormalizesNilConfigAndMapOrder(t *testing.T) {
	first, err := Fingerprint(Input{
		CorpusSnapshotID: "snapshot",
		ChunkSetID:       "chunks",
		BM25ArtifactID:   "bm25",
		Config:           map[string]any{"z": 1, "a": map[string]any{"b": true, "a": "x"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := Fingerprint(Input{
		CorpusSnapshotID: "snapshot",
		ChunkSetID:       "chunks",
		BM25ArtifactID:   "bm25",
		Config:           map[string]any{"a": map[string]any{"a": "x", "b": true}, "z": 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("semantic map order changed fingerprint: %s != %s", first, second)
	}
	withNil, err := Fingerprint(Input{CorpusSnapshotID: "snapshot", ChunkSetID: "chunks", BM25ArtifactID: "bm25"})
	if err != nil {
		t.Fatal(err)
	}
	withEmpty, err := Fingerprint(Input{CorpusSnapshotID: "snapshot", ChunkSetID: "chunks", BM25ArtifactID: "bm25", Config: map[string]any{}})
	if err != nil {
		t.Fatal(err)
	}
	if withNil != withEmpty {
		t.Fatalf("nil config fingerprint = %s, want empty config fingerprint %s", withNil, withEmpty)
	}
}

func TestManifestCarriesExplicitSchemaVersion(t *testing.T) {
	manifest := NewManifest(Input{CorpusSnapshotID: "snapshot", ChunkSetID: "chunks", BM25ArtifactID: "bm25"})
	if manifest.SchemaVersion != SchemaVersion {
		t.Fatalf("schema version = %q, want %q", manifest.SchemaVersion, SchemaVersion)
	}
}
