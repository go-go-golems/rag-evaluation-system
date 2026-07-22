package ragoperators

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type parityFixture struct {
	SchemaVersion string `json:"schemaVersion"`
	Oracle        struct {
		UnitCount, ChunkCount, RepresentationCount int
		RepresentationKinds                        map[string]int     `json:"representationKinds"`
		MetricFixture                              map[string]float64 `json:"metricFixture"`
	} `json:"oracle"`
	IntentionalCorrections []string `json:"intentionalCorrections"`
	BlockingDiscrepancies  []string `json:"blockingDiscrepancies"`
}

func TestRagSol2DeterministicParityFixture(t *testing.T) {
	raw, err := os.ReadFile("testdata/rag-sol2-parity-v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture parityFixture
	if err := json.Unmarshal(raw, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.SchemaVersion != "rag-sol2-parity-fixture/v1" || len(fixture.BlockingDiscrepancies) != 0 {
		t.Fatalf("fixture=%#v", fixture)
	}
	corpus := Corpus{SchemaVersion: "rag-source-record-set/v2", Records: []SourceRecord{{ID: "u1", SessionID: "unicode/session", Ordinal: 1, Role: "user", Text: "Explain café ☕"}, {ID: "a2", SessionID: "unicode/session", Ordinal: 2, Role: "assistant", Text: "The build failed in api.js.\n\nWe decided to update the vector mapping."}, {ID: "tool3", SessionID: "unicode/session", Ordinal: 3, Role: "tool", Text: "tool payload"}, {ID: "a4", SessionID: "unicode/session", Ordinal: 4, Role: "assistant", Text: "Then run the verification 🧪."}}}
	unitOut, err := (unitOperator{"transcript.units.agents-view-runs"}).Execute(context.Background(), ragcontract.Node{ID: "units", Operator: ragcontract.OperatorRef{Kind: "transcript.units.agents-view-runs", Version: "v1"}, Config: json.RawMessage(`{}`)}, map[string]any{"corpus": corpus}, nil)
	if err != nil {
		t.Fatal(err)
	}
	units := unitOut["units"].([]Unit)
	if len(units) != fixture.Oracle.UnitCount || len(units[1].Records) != 2 {
		t.Fatalf("units=%#v", units)
	}
	chunkOut, err := (chunkOperator{}).Execute(context.Background(), ragcontract.Node{ID: "chunks", Operator: (chunkOperator{}).Ref(), Config: json.RawMessage(`{"size":36,"overlap":0}`)}, map[string]any{"units": units}, nil)
	if err != nil {
		t.Fatal(err)
	}
	chunks := chunkOut["chunks"].([]Chunk)
	if len(chunks) != fixture.Oracle.ChunkCount {
		t.Fatalf("chunks=%d want=%d", len(chunks), fixture.Oracle.ChunkCount)
	}
	if err := validateUTF8Ranges(chunks); err != nil {
		t.Fatal(err)
	}
	providers := NewFixtureProviders()
	environment := &Environment{Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Cache: NewMemoryCache(), Usage: Usage{Cost: map[string]float64{}}}
	rawOut, err := (representationOperator{"representations.raw"}).Execute(context.Background(), ragcontract.Node{ID: "raw", Operator: ragcontract.OperatorRef{Kind: "representations.raw", Version: "v1"}, Config: json.RawMessage(`{"name":"raw"}`)}, map[string]any{"chunks": chunks}, environment)
	if err != nil {
		t.Fatal(err)
	}
	summariesOut, err := (representationOperator{"representations.structured-summary"}).Execute(context.Background(), ragcontract.Node{ID: "summary", Operator: ragcontract.OperatorRef{Kind: "representations.structured-summary", Version: "v1"}, Config: json.RawMessage(`{"name":"summary","model":"fixture-summary-v1","prompt":"fixture-transcript-summary-v1","outputSchema":"transcript-rag-summary/v1"}`)}, map[string]any{"chunks": chunks}, environment)
	if err != nil {
		t.Fatal(err)
	}
	summaries := summariesOut["representations"].([]Representation)
	questionsOut, err := (representationOperator{"representations.synthetic-questions"}).Execute(context.Background(), ragcontract.Node{ID: "question", Operator: ragcontract.OperatorRef{Kind: "representations.synthetic-questions", Version: "v1"}, Config: json.RawMessage(`{"name":"question","from":"summary","count":2,"model":"fixture-question-v1","prompt":"fixture-transcript-questions-v1"}`)}, map[string]any{"chunks": chunks, "source": summaries}, environment)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{"raw": len(rawOut["representations"].([]Representation)), "summary": len(summaries), "question": len(questionsOut["representations"].([]Representation))}
	total := 0
	for kind, want := range fixture.Oracle.RepresentationKinds {
		if counts[kind] != want {
			t.Fatalf("%s=%d want=%d", kind, counts[kind], want)
		}
		total += counts[kind]
	}
	if total != fixture.Oracle.RepresentationCount {
		t.Fatalf("representations=%d", total)
	}
	questionRecords := questionsOut["representations"].([]Representation)
	hits := []RankedRecord{}
	for index, record := range questionRecords {
		hits = append(hits, RankedRecord{Rank: index + 1, Representation: record, Score: float64(len(questionRecords) - index), Channel: "question.vector"})
	}
	collapsed, err := (collapseOperator{"collapse.parent"}).Execute(context.Background(), ragcontract.Node{ID: "collapse", Operator: ragcontract.OperatorRef{Kind: "collapse.parent", Version: "v1"}, Config: json.RawMessage(`{"scope":"unit","representative":"scoreThenRepresentationId"}`)}, map[string]any{"hits": hits}, environment)
	if err != nil {
		t.Fatal(err)
	}
	parents := collapsed["parents"].([]RankedParent)
	if len(parents) != fixture.Oracle.UnitCount {
		t.Fatalf("collapsed parents=%d", len(parents))
	}
	for _, parent := range parents {
		if len(parent.Members) < 2 {
			t.Fatalf("multiplicity was lost before collapse: %#v", parent)
		}
	}
	metrics := Evaluate(Query{ID: "graded", RelevantIDs: []string{"a", "c"}, Grades: map[string]float64{"a": 3, "c": 1}}, []Evidence{{Collapse: ragcontract.CollapseIdentity{Scope: "unit", ID: "a"}}, {Collapse: ragcontract.CollapseIdentity{Scope: "unit", ID: "b"}}, {Collapse: ragcontract.CollapseIdentity{Scope: "unit", ID: "c"}}}, nil, []ragcontract.Measure{{Name: "rag.precision", Config: json.RawMessage(`{"cutoffs":[3]}`)}, {Name: "rag.recall", Config: json.RawMessage(`{"cutoffs":[3]}`)}, {Name: "rag.hit-rate", Config: json.RawMessage(`{"cutoffs":[3]}`)}, {Name: "rag.mrr", Config: json.RawMessage(`{}`)}, {Name: "rag.ndcg", Config: json.RawMessage(`{"cutoffs":[3]}`)}}, nil, Usage{}, nil, 0)
	metricValues := map[string]json.RawMessage{}
	for _, metric := range metrics {
		metricValues[metric.Name] = metric.Value
	}
	assertMetric := func(name, key string, want float64) {
		var value any
		if err := json.Unmarshal(metricValues[name], &value); err != nil {
			t.Fatal(err)
		}
		got := 0.0
		if key == "" {
			got = value.(float64)
		} else {
			got = value.(map[string]any)[key].(float64)
		}
		if difference := got - want; difference > 1e-12 || difference < -1e-12 {
			t.Fatalf("%s=%v want=%v", name, got, want)
		}
	}
	assertMetric("rag.precision", "3", fixture.Oracle.MetricFixture["precisionAt3"])
	assertMetric("rag.recall", "3", fixture.Oracle.MetricFixture["recallAt3"])
	assertMetric("rag.hit-rate", "3", fixture.Oracle.MetricFixture["hitRateAt3"])
	assertMetric("rag.mrr", "", fixture.Oracle.MetricFixture["mrr"])
	assertMetric("rag.ndcg", "3", fixture.Oracle.MetricFixture["ndcgAt3"])
	chunkTarget := Evaluate(Query{RelevantIDs: []string{"unit-relevant"}}, []Evidence{{Collapse: ragcontract.CollapseIdentity{Scope: "chunk", ID: "chunk-hit"}, Chunk: Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-hit", ParentUnitID: "unit-relevant"}}}}, nil, []ragcontract.Measure{{Name: "rag.mrr", Config: json.RawMessage(`{}`)}}, nil, Usage{}, nil, 0)
	if len(chunkTarget) != 1 || chunkTarget[0].Numeric == nil || *chunkTarget[0].Numeric != 1 {
		t.Fatalf("chunk collapse lost unit relevance: %#v", chunkTarget)
	}
	deduplicated := EvaluateForTarget(Query{RelevantIDs: []string{"unit-relevant"}}, []Evidence{{Chunk: Chunk{Record: ragcontract.ChunkRecord{ID: "c1", ParentUnitID: "unit-relevant"}}}, {Chunk: Chunk{Record: ragcontract.ChunkRecord{ID: "c2", ParentUnitID: "unit-relevant"}}}, {Chunk: Chunk{Record: ragcontract.ChunkRecord{ID: "c3", ParentUnitID: "unit-other"}}}}, nil, []ragcontract.Measure{{Name: "rag.precision", Config: json.RawMessage(`{"cutoffs":[3]}`)}}, nil, Usage{}, nil, 0, "unit")
	var precision map[string]float64
	if err := json.Unmarshal(deduplicated[0].Value, &precision); err != nil {
		t.Fatal(err)
	}
	if precision["3"] != 0.5 {
		t.Fatalf("unit target multiplicity bias: %#v", precision)
	}
}
