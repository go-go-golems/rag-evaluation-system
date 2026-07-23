// Builds a researchctl import export from an already completed TTC sweep.
// It reads only compact evidence and operation JSONL/manifest paths; it never
// opens the source corpus, provider configuration, runtime SQLite state, or
// provider payloads.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/go-go-golems/researchctl/pkg/lab"
)

type evidence struct {
	Cells []cell `json:"cells"`
}

type cell struct {
	Cell struct {
		ChunksPerRequest int `json:"chunksPerRequest"`
		Concurrency      int `json:"concurrency"`
	} `json:"cell"`
	Requests          int   `json:"requests"`
	EmbeddingRequests int   `json:"embeddingRequests"`
	MakespanMicros    int64 `json:"makespanMicros"`
	Operations        struct {
		JSONLPath    string `json:"jsonlPath"`
		ManifestPath string `json:"manifestPath"`
	} `json:"operations"`
}

func main() {
	var root, specificationPath, runID, attemptID, externalRunID, recordedAt, output string
	flag.StringVar(&root, "root", "", "completed sweep output root")
	flag.StringVar(&specificationPath, "specification", "", "canonical researchctl specification")
	flag.StringVar(&runID, "run-id", "", "explicit researchctl run ID")
	flag.StringVar(&attemptID, "attempt-id", "", "explicit researchctl attempt ID")
	flag.StringVar(&externalRunID, "external-run-id", "", "stable external sweep identity")
	flag.StringVar(&recordedAt, "recorded-at", "", "RFC3339 completion timestamp")
	flag.StringVar(&output, "output", "", "output export JSON path")
	flag.Parse()
	if root == "" || specificationPath == "" || runID == "" || attemptID == "" || externalRunID == "" || recordedAt == "" || output == "" {
		panic("all flags are required")
	}
	stamp, err := time.Parse(time.RFC3339Nano, recordedAt)
	if err != nil {
		panic(err)
	}
	body, err := os.ReadFile(filepath.Join(root, "evidence.json"))
	if err != nil {
		panic(err)
	}
	var measured evidence
	if err := json.Unmarshal(body, &measured); err != nil {
		panic(err)
	}
	if len(measured.Cells) == 0 {
		panic("sweep has no cells")
	}
	specification, err := lab.ReadSpecificationRecord(specificationPath)
	if err != nil {
		panic(err)
	}
	artifacts := []researchctladapter.OperationCustodyArtifact{{Role: "sweep-evidence", Kind: "rag-ttc-sweep-evidence", URI: "evidence.json", Source: filepath.Join(root, "evidence.json"), SchemaVersion: "rag-ttc-v3-sweep-evidence/v2", MediaType: "application/json"}}
	var generationRequests, embeddingRequests, makespanMicros int64
	for index, value := range measured.Cells {
		cellName := fmt.Sprintf("cell-%02d-b%d-c%d", index, value.Cell.ChunksPerRequest, value.Cell.Concurrency)
		if value.Operations.JSONLPath == "" || value.Operations.ManifestPath == "" {
			panic(fmt.Sprintf("cell %d lacks operation custody", index))
		}
		artifacts = append(artifacts,
			researchctladapter.OperationCustodyArtifact{Role: "cell-evidence", Kind: "rag-ttc-cell-evidence", URI: filepath.ToSlash(filepath.Join("cells", cellName+".json")), Source: filepath.Join(root, "cells", cellName+".json"), SchemaVersion: "rag-ttc-v3-cell-evidence/v1", MediaType: "application/json"},
			researchctladapter.OperationCustodyArtifact{Role: "operation-ledger", Kind: "workflow-operation-jsonl", URI: value.Operations.JSONLPath, Source: filepath.Join(root, filepath.FromSlash(value.Operations.JSONLPath)), SchemaVersion: "workflow-v3-external-operation-export/v1", MediaType: "application/x-ndjson"},
			researchctladapter.OperationCustodyArtifact{Role: "operation-manifest", Kind: "workflow-operation-manifest", URI: value.Operations.ManifestPath, Source: filepath.Join(root, filepath.FromSlash(value.Operations.ManifestPath)), SchemaVersion: "workflow-v3-external-operation-export/v1", MediaType: "application/json"},
		)
		generationRequests += int64(value.Requests)
		embeddingRequests += int64(value.EmbeddingRequests)
		makespanMicros += value.MakespanMicros
	}
	export, err := researchctladapter.BuildOperationCustodyRunExport(researchctladapter.OperationCustodyExportInput{
		Specification: *specification,
		Source:        lab.ExportSource{Namespace: "rag-ttc-v3-sweep", ExternalRunID: externalRunID},
		RunID:         runID,
		AttemptID:     attemptID,
		RecordedAt:    stamp,
		Status:        "succeeded",
		Artifacts:     artifacts,
		Metrics: []researchctladapter.OperationCustodyMetric{
			{Name: "operation.cell_count", Units: int64(len(measured.Cells)), Unit: "cells"},
			{Name: "operation.generation_requests", Units: generationRequests, Unit: "requests"},
			{Name: "operation.embedding_requests", Units: embeddingRequests, Unit: "requests"},
			{Name: "operation.makespan_micros", Units: makespanMicros, Unit: "micros"},
		},
	})
	if err != nil {
		panic(err)
	}
	canonical, err := lab.CanonicalJSON(export)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(output, append(canonical, '\n'), 0o644); err != nil {
		panic(err)
	}
}
