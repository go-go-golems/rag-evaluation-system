package preview

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	studycmd "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/study"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var inputs, artifactRoot, ttcDatabase, project, database, experiment, researchctlCommand, worker, variant, query, trace, outputDirectory string
	var workerArgs, factors, secrets []string
	var timeout time.Duration
	var maxAttempts int
	command := &cobra.Command{Use: "preview <study.js>", Short: "Run one query from one RAG study cell through a scratch or existing researchctl project", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		factorMap, err := parseFactors(factors)
		if err != nil {
			return err
		}
		scratch := ""
		if project == "" {
			project, scratch, err = createScratch(command.Context(), researchctlCommand)
			if err != nil {
				return err
			}
			experiment = "PREVIEW-001"
		}
		if artifactRoot == "" {
			absolute, _ := filepath.Abs(project)
			artifactRoot = filepath.Join(filepath.Dir(absolute), ".researchctl", "artifacts")
		}
		study, err := studycmd.LoadStudy(args[0])
		if err != nil {
			return err
		}
		document, base, err := researchctladapter.LoadInputs(inputs)
		if err != nil {
			return err
		}
		var catalog researchctladapter.CatalogResolver
		if ttcDatabase != "" {
			catalog = researchctladapter.NewTTCCatalog(ttcDatabase)
		}
		resolved, err := researchctladapter.ResolveInputs(command.Context(), document, base, artifactRoot, catalog)
		if err != nil {
			return err
		}
		corpus, ok := resolved.ByRole["corpus"]
		if !ok {
			return fmt.Errorf("preview requires corpus input")
		}
		dataset := ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-queries/v2", Queries: []ragoperators.Query{{ID: "preview-query", Text: query}}}
		envelope := ragoperators.NewEvaluationArtifact(dataset, "preview-query", "preview", "candidate", study.Dataset.RelevanceTarget, corpus.Binding.Digest)
		evaluation, err := researchctladapter.StageEnvelope(researchctladapter.InputReference{Role: "evaluation-dataset", Kind: "manifest-envelope", SchemaVersion: ragcontract.EvaluationManifestSchema}, envelope, artifactRoot)
		if err != nil {
			return err
		}
		delete(resolved.ByRole, "judgments")
		resolved.ByRole["evaluation-dataset"] = evaluation
		study.Dataset.Split = "preview"
		study.Dataset.Status = "candidate"
		study, cells, err := researchctladapter.Expand(study, resolved)
		if err != nil {
			return err
		}
		cell, err := selectCell(cells, variant, factorMap)
		if err != nil {
			return err
		}
		spec, err := researchctladapter.WrapExecution(cell.Execution, resolved, "preview / "+cell.VariantID)
		if err != nil {
			return err
		}
		result, err := researchctladapter.ExecuteSpecification(command.Context(), spec, researchctladapter.RunOptions{ResearchctlCommand: researchctlCommand, Project: project, Database: database, ExperimentID: experiment, Worker: researchctladapter.WorkerCommand{Executable: worker, Args: workerArgs}, MaxAttempts: maxAttempts, Timeout: timeout, SecretEnvironment: secrets, OutputDirectory: outputDirectory})
		if err != nil {
			return err
		}
		return json.NewEncoder(command.OutOrStdout()).Encode(map[string]any{"schemaVersion": "rag-preview-result/v1", "scratch": scratch, "project": project, "trace": trace, "cell": cell, "result": result})
	}}
	command.Flags().StringVar(&inputs, "inputs", "", "RAG input bindings/catalog aliases JSON")
	command.Flags().StringVar(&artifactRoot, "artifact-root", "", "Researchctl artifact root")
	command.Flags().StringVar(&ttcDatabase, "ttc-database", "", "Read-only TTC catalog SQLite database")
	command.Flags().StringVarP(&project, "project", "p", "", "Existing researchctl project; omitted creates a scratch project")
	command.Flags().StringVar(&database, "database", "", "Researchctl laboratory database")
	command.Flags().StringVar(&experiment, "experiment-id", "PREVIEW-001", "Experiment receiving the preview run")
	command.Flags().StringVar(&researchctlCommand, "researchctl-command", "researchctl", "Researchctl executable")
	command.Flags().StringVar(&worker, "worker-command", "rag-worker", "RAG worker executable")
	command.Flags().StringSliceVar(&workerArgs, "worker-arg", nil, "RAG worker argument")
	command.Flags().StringVar(&variant, "variant", "", "Variant ID (defaults to first cell)")
	command.Flags().StringSliceVar(&factors, "factor", nil, "Factor selection key=value")
	command.Flags().StringVar(&query, "query", "", "Preview query")
	command.Flags().StringVar(&trace, "trace", "full", "Trace presentation policy")
	command.Flags().StringSliceVar(&secrets, "secret-env", nil, "Secret environment variable")
	command.Flags().IntVar(&maxAttempts, "max-attempts", 1, "Maximum attempts")
	command.Flags().DurationVar(&timeout, "timeout", 0, "Attempt timeout")
	command.Flags().StringVar(&outputDirectory, "spec-output-dir", "", "Canonical specification output directory")
	_ = command.MarkFlagRequired("inputs")
	_ = command.MarkFlagRequired("query")
	return command
}
func parseFactors(values []string) (map[string]string, error) {
	result := map[string]string{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid --factor %q", value)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
func selectCell(cells []ragcontract.ExpandedCell, variant string, factors map[string]string) (ragcontract.ExpandedCell, error) {
	for _, cell := range cells {
		if variant != "" && cell.VariantID != variant {
			continue
		}
		matches := true
		for key, value := range factors {
			found := false
			for _, selection := range cell.Factors {
				if selection.FactorID == key && selection.ValueID == value {
					found = true
				}
			}
			if !found {
				matches = false
			}
		}
		if matches {
			return cell, nil
		}
	}
	return ragcontract.ExpandedCell{}, fmt.Errorf("RAG_PREVIEW_CELL_NOT_FOUND")
}
func createScratch(ctx context.Context, researchctlCommand string) (string, string, error) {
	root, err := os.MkdirTemp("", "rag-preview-")
	if err != nil {
		return "", "", err
	}
	project := filepath.Join(root, "project.yaml")
	content := `schemaVersion: 2
kind: ResearchProject
id: PRJ-RAG-PREVIEW
name: RAG preview scratch project
description: Scratch project created by rag-eval preview.
goals:
  - id: GOAL-001
    title: Preview one query
    status: active
    priority: P0
    asks: [Q-001]
questions:
  - id: Q-001
    text: Does the selected RAG cell execute?
    status: active
    priority: P0
    hypotheses: [H-001]
hypotheses:
  - id: H-001
    claim: The selected cell executes through the generic laboratory.
    status: open
    priority: P0
    testedBy: [PREVIEW-001]
    confidence: unknown
experiments:
  - id: PREVIEW-001
    title: One-query RAG preview
    status: planned
    priority: P0
    hypotheses: [H-001]
    expectedArtifacts:
      - type: trace
        required: true
        description: RAG query trace.
    metrics:
      - name: rag.mrr
        unit: ratio
        required: false
    successCriteria:
      - The worker completes without lifecycle violations.
`
	if err := os.WriteFile(project, []byte(content), 0o644); err != nil {
		return "", "", err
	}
	command := exec.CommandContext(ctx, researchctlCommand, "lab", "init", "--project", project, "--output", "json")
	if output, err := command.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("initialize scratch laboratory: %w: %s", err, output)
	}
	return project, root, nil
}
