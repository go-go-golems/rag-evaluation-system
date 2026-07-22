package ragcompiler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestV2Goldens(t *testing.T) {
	pipeline, err := Normalize(matrixPipeline([]string{"raw", "summary", "question"}, "unit"), nil)
	if err != nil {
		t.Fatal(err)
	}
	product, err := CompileProduct(ragcontract.ProductPlan{SchemaVersion: ragcontract.ProductSchemaVersion, Pipeline: matrixPipeline([]string{"raw"}, "unit"), Bindings: testBindings()[:1], Citations: ragcontract.CitationPolicy{Mode: "required"}, Runtime: ragcontract.RuntimePolicy{MaxResults: 10}, Display: ragcontract.DisplayMetadata{Name: "qualification"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	variants := []ragcontract.Variant{{ID: "raw", Pipeline: matrixPipeline([]string{"raw"}, "unit")}, {ID: "summary", Pipeline: matrixPipeline([]string{"summary"}, "unit")}, {ID: "raw-summary", Pipeline: matrixPipeline([]string{"raw", "summary"}, "unit")}, {ID: "raw-question", Pipeline: matrixPipeline([]string{"raw", "question"}, "unit")}, {ID: "all", Pipeline: matrixPipeline([]string{"raw", "summary", "question"}, "unit")}}
	study := ragcontract.Study{SchemaVersion: ragcontract.StudySchemaVersion, Variants: variants, Factors: []ragcontract.Factor{{ID: "collapse", Values: []ragcontract.FactorValue{{ID: "chunk", Value: json.RawMessage(`"chunk"`), Overrides: []ragcontract.NodeConfigOverride{{NodeID: "final", Config: json.RawMessage(`{"scope":"chunk"}`)}}}, {ID: "unit", Value: json.RawMessage(`"unit"`), Overrides: []ragcontract.NodeConfigOverride{{NodeID: "final", Config: json.RawMessage(`{"scope":"unit"}`)}}}}}}, Bindings: testBindings(), Dataset: ragcontract.DatasetBinding{ManifestDigest: digest("dataset"), Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", ValueKind: "number", Required: true, Config: json.RawMessage(`{}`)}}, Replicates: 2, Display: ragcontract.DisplayMetadata{Name: "sol2 matrix"}}
	cells, err := ExpandStudy(study, nil)
	if err != nil {
		t.Fatal(err)
	}
	trace := ragcontract.QueryTrace{SchemaVersion: ragcontract.TraceSchemaVersion, Query: ragcontract.QueryInputTrace{ID: "q-1", TextDigest: digest("query"), DatasetSplit: "smoke"}, Operators: []ragcontract.OperatorTrace{}, Channels: []ragcontract.ChannelTrace{}, Collapses: []ragcontract.CollapseTrace{}, Results: []ragcontract.ResultTrace{}, Timing: ragcontract.TimingTrace{ByOperator: map[string]int64{}}, Usage: ragcontract.UsageTrace{}, Failures: []ragcontract.FailureTrace{}}
	pipelineID, _ := ragcontract.Digest(pipeline)
	productID, _ := ProductSemanticIdentity(product)
	studyID, _ := StudySemanticIdentity(study)
	identities := map[string]any{"pipeline": pipelineID, "product": productID, "study": studyID, "cells": cellIDs(cells)}
	goldens := map[string]any{"pipeline.golden.json": pipeline, "product.golden.json": product, "study.golden.json": study, "matrix-cells.golden.json": cells, "trace.golden.json": trace, "identities.golden.json": identities}
	for name, value := range goldens {
		assertGolden(t, name, value)
	}
}
func cellIDs(cells []ragcontract.ExpandedCell) []string {
	r := make([]string, len(cells))
	for i := range cells {
		r[i] = cells[i].ID
	}
	return r
}
func assertGolden(t *testing.T, name string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var stored any
	if err := json.Unmarshal(want, &stored); err != nil {
		t.Fatalf("decode golden %s: %v", path, err)
	}
	actualCanonical, err := ragcontract.CanonicalJSON(value)
	if err != nil {
		t.Fatal(err)
	}
	storedCanonical, err := ragcontract.CanonicalJSON(stored)
	if err != nil {
		t.Fatal(err)
	}
	if string(storedCanonical) != string(actualCanonical) {
		t.Fatalf("golden mismatch %s; run UPDATE_GOLDEN=1 go test ./pkg/ragcompiler", path)
	}
}
