package ragcompiler

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestNormalizeIsIdempotentAndExpandsDefaults(t *testing.T) {
	pipeline := matrixPipeline([]string{"raw"}, "unit")
	first, err := Normalize(pipeline, nil)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Normalize(first, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		a, _ := json.Marshal(first)
		b, _ := json.Marshal(second)
		t.Fatalf("not idempotent:\n%s\n%s", a, b)
	}
	found := false
	for _, node := range first.Nodes {
		if node.Operator.Kind == "fusion.weighted-rrf" {
			found = true
			if !strings.Contains(string(node.Config), `"rankConstant":60`) {
				t.Fatalf("defaults missing: %s", node.Config)
			}
		}
	}
	if !found {
		t.Fatal("fusion missing")
	}
}

func TestNormalizeRejectsCycleAndPortMismatch(t *testing.T) {
	p := ragcontract.PipelineIR{SchemaVersion: ragcontract.PipelineSchemaVersion, Nodes: []ragcontract.Node{
		{ID: "a", Operator: ragcontract.OperatorRef{Kind: "collapse.final", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "parents", From: ragcontract.PortRef{NodeID: "b", Port: "parents"}}}, Config: json.RawMessage(`{}`)},
		{ID: "b", Operator: ragcontract.OperatorRef{Kind: "collapse.final", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "parents", From: ragcontract.PortRef{NodeID: "a", Port: "parents"}}}, Config: json.RawMessage(`{}`)},
	}}
	_, err := Normalize(p, nil)
	if err == nil || !strings.Contains(err.Error(), "RAG_V2_GRAPH_CYCLE") {
		t.Fatalf("cycle error=%v", err)
	}
	p = matrixPipeline([]string{"raw"}, "unit")
	p.Nodes[len(p.Nodes)-1].Inputs[0].From.Port = "answer"
	_, err = Normalize(p, nil)
	if err == nil || !strings.Contains(err.Error(), "RAG_V2_OUTPUT_PORT") {
		t.Fatalf("port error=%v", err)
	}
}

func TestNormalizeRejectsUnknownConfigAndUnsafeQuestionMultiplicity(t *testing.T) {
	pipeline := matrixPipeline([]string{"raw"}, "unit")
	pipeline.Nodes[1].Config = json.RawMessage(`{"size":800,"mystery":true}`)
	_, err := Normalize(pipeline, nil)
	if err == nil || !strings.Contains(err.Error(), "unknown config field") {
		t.Fatalf("unknown config error = %v", err)
	}

	pipeline = matrixPipeline([]string{"question"}, "unit")
	for i, node := range pipeline.Nodes {
		if node.Operator.Kind == "collapse.parent" {
			pipeline.Nodes[i].Operator = ragcontract.OperatorRef{Kind: "fusion.weighted-rrf", Version: "v1"}
			pipeline.Nodes[i].Inputs[0].Port = "channel.unsafe"
			pipeline.Nodes[i].Config = json.RawMessage(`{}`)
			break
		}
	}
	_, err = Normalize(pipeline, nil)
	if err == nil || !strings.Contains(err.Error(), "RAG_V2_COLLAPSE_REQUIRED") {
		t.Fatalf("unsafe question error = %v", err)
	}
}

func TestStudyExpandsTenStableUniqueCellsAndExcludesDisplayIdentity(t *testing.T) {
	variants := []ragcontract.Variant{}
	for _, v := range []struct {
		id    string
		kinds []string
	}{{"raw", []string{"raw"}}, {"summary", []string{"summary"}}, {"raw-summary", []string{"raw", "summary"}}, {"raw-question", []string{"raw", "question"}}, {"all", []string{"raw", "summary", "question"}}} {
		variants = append(variants, ragcontract.Variant{ID: v.id, Pipeline: matrixPipeline(v.kinds, "unit")})
	}
	study := ragcontract.Study{SchemaVersion: ragcontract.StudySchemaVersion, Variants: variants, Factors: []ragcontract.Factor{{ID: "collapse", Values: []ragcontract.FactorValue{{ID: "chunk", Value: json.RawMessage(`"chunk"`), Overrides: []ragcontract.NodeConfigOverride{{NodeID: "final", Config: json.RawMessage(`{"scope":"chunk"}`)}}}, {ID: "unit", Value: json.RawMessage(`"unit"`), Overrides: []ragcontract.NodeConfigOverride{{NodeID: "final", Config: json.RawMessage(`{"scope":"unit"}`)}}}}}}, Bindings: testBindings(), Dataset: ragcontract.DatasetBinding{ManifestDigest: digest("dataset"), Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", ValueKind: "number", Required: true}}, Replicates: 2, Display: ragcontract.DisplayMetadata{Name: "matrix"}}
	cells, err := ExpandStudy(study, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(cells) != 10 {
		t.Fatalf("cells=%d", len(cells))
	}
	wantOrder := []string{"all/chunk", "all/unit", "raw/chunk", "raw/unit", "raw-question/chunk", "raw-question/unit", "raw-summary/chunk", "raw-summary/unit", "summary/chunk", "summary/unit"}
	seen := map[string]bool{}
	for index, c := range cells {
		gotOrder := c.VariantID + "/" + c.Factors[0].ValueID
		if gotOrder != wantOrder[index] {
			t.Fatalf("cell[%d] order = %s, want %s", index, gotOrder, wantOrder[index])
		}
		if seen[c.ID] {
			t.Fatalf("duplicate %s", c.ID)
		}
		seen[c.ID] = true
		if c.Replicates != 2 {
			t.Fatalf("replicates=%d", c.Replicates)
		}
	}
	id1, err := StudySemanticIdentity(study)
	if err != nil {
		t.Fatal(err)
	}
	study.Display.Name = "renamed"
	study.Display.Notes = []string{"non semantic"}
	id2, _ := StudySemanticIdentity(study)
	if id1 != id2 {
		t.Fatalf("display changed identity: %s %s", id1, id2)
	}
}

func TestRecipeExpandsToOrdinaryNodes(t *testing.T) {
	p := ragcontract.PipelineIR{SchemaVersion: ragcontract.PipelineSchemaVersion, Inputs: []ragcontract.InputSlot{{ID: "chunks", Kind: ragcontract.PortChunks, BindingMode: "artifact", ArtifactRole: "chunks", ManifestSchema: ragcontract.ChunkSetManifestSchema}}, Nodes: []ragcontract.Node{{ID: "recipe", Operator: ragcontract.OperatorRef{Kind: "representations.compose", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "in", From: ragcontract.PortRef{NodeID: "chunks", Port: "out"}}}, Config: json.RawMessage(`{"operators":[{"kind":"representations.raw","version":"v1","config":{}},{"kind":"representations.structured-summary","version":"v1","config":{"name":"summary"}}]}`)}}, Outputs: []ragcontract.OutputRef{{Name: "representations", Kind: ragcontract.PortRepresentations, From: ragcontract.PortRef{NodeID: "recipe", Port: "out"}}}}
	normalized, err := Normalize(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(normalized.Nodes) != 3 {
		t.Fatalf("expanded nodes = %d, want 3", len(normalized.Nodes))
	}
	for _, node := range normalized.Nodes {
		if node.Operator.Kind == "representations.compose" {
			t.Fatalf("recipe remains: %#v", normalized.Nodes)
		}
	}
}

func matrixPipeline(kinds []string, scope string) ragcontract.PipelineIR {
	inputs := []ragcontract.InputSlot{
		{ID: "corpus", Kind: ragcontract.PortCorpus, BindingMode: "artifact", ArtifactRole: "corpus", ManifestSchema: ragcontract.CorpusManifestSchema, Digest: digest("corpus")},
		{ID: "query", Kind: ragcontract.PortQuery, BindingMode: "dataset"},
	}
	nodes := []ragcontract.Node{
		node("units", "units.identity", "v1", map[string]ragcontract.PortRef{"corpus": {NodeID: "corpus", Port: "out"}}, `{}`),
		node("chunks", "chunks.recursive", "v1", map[string]ragcontract.PortRef{"units": {NodeID: "units", Port: "units"}}, `{"size":800}`),
	}
	indexInputs := map[string]ragcontract.PortRef{}
	for _, kind := range kinds {
		if kind == "raw" {
			nodes = append(nodes, node("raw-representations", "representations.raw", "v1", map[string]ragcontract.PortRef{"chunks": {NodeID: "chunks", Port: "chunks"}}, `{"name":"raw"}`))
			indexInputs["representations.raw"] = ragcontract.PortRef{NodeID: "raw-representations", Port: "representations"}
			continue
		}
		slotID := kind + "-representations"
		inputs = append(inputs, ragcontract.InputSlot{ID: slotID, Kind: ragcontract.PortRepresentations, BindingMode: "artifact", ArtifactRole: slotID, ManifestSchema: ragcontract.RepresentationManifestSchema, Digest: digest(slotID)})
		indexInputs["representations."+kind] = ragcontract.PortRef{NodeID: slotID, Port: "out"}
	}
	nodes = append(nodes, node("index", "index.bleve-multi", "v1", indexInputs, `{"representationKinds":`+mustJSONString(kinds)+`}`))
	fusionInputs := map[string]ragcontract.PortRef{}
	order := 0
	for _, kind := range kinds {
		for _, backend := range []string{"bm25", "vector"} {
			name := kind + "." + backend
			retrieve := name + ".retrieve"
			collapse := name + ".collapse"
			nodes = append(nodes, node(retrieve, "retrieve."+backend, "v1", map[string]ragcontract.PortRef{"index": {NodeID: "index", Port: "index"}, "query": {NodeID: "query", Port: "out"}}, `{"representation":"`+kind+`","topK":30}`))
			nodes = append(nodes, node(collapse, "collapse.parent", "v1", map[string]ragcontract.PortRef{"hits": {NodeID: retrieve, Port: "hits"}}, `{"scope":"unit"}`))
			nodes[len(nodes)-2].Order = order
			nodes[len(nodes)-1].Order = order
			order++
			fusionInputs["channel."+name] = ragcontract.PortRef{NodeID: collapse, Port: "parents"}
		}
	}
	nodes = append(nodes,
		node("fusion", "fusion.weighted-rrf", "v1", fusionInputs, `{"weights":{"raw.vector":2}}`),
		node("final", "collapse.final", "v1", map[string]ragcontract.PortRef{"parents": {NodeID: "fusion", Port: "parents"}}, `{"scope":"`+scope+`"}`),
		node("hydrate", "hydrate.source-evidence", "v1", map[string]ragcontract.PortRef{"parents": {NodeID: "final", Port: "parents"}, "chunks": {NodeID: "chunks", Port: "chunks"}}, `{"policy":"best-contribution-then-id"}`),
	)
	return ragcontract.PipelineIR{SchemaVersion: ragcontract.PipelineSchemaVersion, Inputs: inputs, Nodes: nodes, Outputs: []ragcontract.OutputRef{{Name: "evidence", Kind: ragcontract.PortEvidence, From: ragcontract.PortRef{NodeID: "hydrate", Port: "evidence"}}}}
}
func FuzzNormalizeNeverPanics(f *testing.F) {
	seed, _ := json.Marshal(matrixPipeline([]string{"raw"}, "unit"))
	f.Add(seed)
	f.Fuzz(func(t *testing.T, data []byte) {
		var pipeline ragcontract.PipelineIR
		if json.Unmarshal(data, &pipeline) != nil {
			return
		}
		_, _ = Normalize(pipeline, nil)
	})
}

func node(id, kind, version string, inputs map[string]ragcontract.PortRef, config string) ragcontract.Node {
	bindings := []ragcontract.InputBinding{}
	for name, from := range inputs {
		bindings = append(bindings, ragcontract.InputBinding{Port: name, From: from})
	}
	sort.Slice(bindings, func(i, j int) bool { return bindings[i].Port < bindings[j].Port })
	return ragcontract.Node{ID: id, Operator: ragcontract.OperatorRef{Kind: kind, Version: version}, Inputs: bindings, Config: json.RawMessage(config)}
}
func testBindings() []ragcontract.ArtifactBinding {
	size := int64(10)
	return []ragcontract.ArtifactBinding{
		{SlotID: "corpus", Role: "corpus", Kind: "manifest", Digest: digest("corpus"), SizeBytes: &size, SchemaVersion: ragcontract.CorpusManifestSchema},
		{SlotID: "summary-representations", Role: "summary-representations", Kind: "manifest", Digest: digest("summary-representations"), SizeBytes: &size, SchemaVersion: ragcontract.RepresentationManifestSchema},
		{SlotID: "question-representations", Role: "question-representations", Kind: "manifest", Digest: digest("question-representations"), SizeBytes: &size, SchemaVersion: ragcontract.RepresentationManifestSchema},
	}
}
func mustJSONString(value any) string { data, _ := json.Marshal(value); return string(data) }
func digest(v string) string          { d, _ := ragcontract.Digest(v); return d }
