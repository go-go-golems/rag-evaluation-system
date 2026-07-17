package ragcontract

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestStrictDecodersRejectUnknownAndTrailingValues(t *testing.T) {
	for _, raw := range []string{`{"schemaVersion":"rag-pipeline-ir/v2","inputs":[],"nodes":[],"outputs":[],"unknown":true}`, `{"schemaVersion":"rag-pipeline-ir/v2","inputs":[],"nodes":[],"outputs":[]} {}`} {
		if _, err := DecodePipeline(strings.NewReader(raw)); err == nil {
			t.Fatalf("accepted %s", raw)
		}
	}
}
func TestCanonicalJSONAndDigestIgnoreMapOrder(t *testing.T) {
	a := map[string]any{"b": 2, "a": 1}
	b := map[string]any{"a": 1, "b": 2}
	ca, _ := CanonicalJSON(a)
	cb, _ := CanonicalJSON(b)
	if !bytes.Equal(ca, cb) {
		t.Fatalf("%s != %s", ca, cb)
	}
	da, _ := Digest(a)
	db, _ := Digest(b)
	if da != db {
		t.Fatalf("%s != %s", da, db)
	}
}
func TestSchemasAreValidJSON(t *testing.T) {
	entries, err := os.ReadDir("schema")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 5 {
		t.Fatalf("schemas = %d", len(entries))
	}
	for _, entry := range entries {
		data, err := os.ReadFile("schema/" + entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatalf("%s: %v", entry.Name(), err)
		}
		if schema["$id"] == "" {
			t.Fatalf("%s has no $id", entry.Name())
		}
	}
}

func TestDistinctRepresentationCollapseAndEvidenceIdentity(t *testing.T) {
	result := ResultTrace{Rank: 1, Collapse: CollapseIdentity{Scope: "unit", ID: "unit:1"}, MatchedRepresentations: []MatchedRepresentation{{ID: "representation:question:1", Kind: "question", Channel: "question.vector", Rank: 1}}, Evidence: EvidenceIdentity{ChunkID: "chunk:source:1", Digest: "sha256:source", Citation: CitationRef{SourceID: "turn:1"}}}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, want := range []string{"representation:question:1", "unit:1", "chunk:source:1", "sha256:source"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %s in %s", want, text)
		}
	}
}
func TestManifestLineageRequiresParentAndProduction(t *testing.T) {
	base := ManifestBase{SchemaVersion: ChunkSetManifestSchema, Digest: "sha256:" + strings.Repeat("a", 64)}
	if err := ValidateManifestBase(base, ChunkSetManifestSchema, true); err == nil {
		t.Fatal("accepted child manifest without lineage")
	}
	base.Parents = []ParentDigest{{Role: "units", Digest: "sha256:" + strings.Repeat("b", 64), SchemaVersion: UnitSetManifestSchema}}
	base.Production = &Production{Operator: OperatorRef{Kind: "chunks.recursive", Version: "v1"}, Config: json.RawMessage(`{}`)}
	if err := ValidateManifestBase(base, ChunkSetManifestSchema, true); err != nil {
		t.Fatal(err)
	}
}

func FuzzDecodePipelineNeverPanics(f *testing.F) {
	f.Add([]byte(`{"schemaVersion":"rag-pipeline-ir/v2","inputs":[],"nodes":[],"outputs":[]}`))
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = DecodePipeline(bytes.NewReader(data)) })
}
