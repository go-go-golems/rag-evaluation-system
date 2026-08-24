package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3product"
)

func main() {
	database := flag.String("workflow-db", "", "Workflow V3 database")
	artifacts := flag.String("artifact-root", "", "Workflow artifact root")
	runID := flag.String("run-id", "", "Workflow run ID")
	flag.Parse()
	if *database == "" || *artifacts == "" || *runID == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "--workflow-db, --artifact-root, and --run-id are required")
		os.Exit(2)
	}
	config := workflowv3product.DefaultConfig()
	config.DatabasePath = *database
	config.ArtifactRoot = *artifacts
	config.TaskPackages = []string{ragworkflow.PackageName}
	config.Capacities = map[string]int{"cpu.rag.prepare": 1, "cpu.rag.query": 1, "cpu.rag.reduce": 1}
	app, err := workflowv3product.Open(context.Background(), config, ragworkflow.NewPackage())
	if err != nil {
		fmt.Fprintln(os.Stderr, "open failed")
		os.Exit(1)
	}
	defer func() { _ = app.Close() }()
	observations, err := app.Observations(context.Background(), workflowv3.RunID(*runID))
	if err != nil {
		fmt.Fprintln(os.Stderr, "observation failed")
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(observations); err != nil {
		os.Exit(1)
	}
}
