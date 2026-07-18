package study

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/spf13/cobra"
)

type commonFlags struct{ inputs, artifactRoot, ttcDatabase, output string }
type runFlags struct {
	project, database, experiment, researchctl, worker, outputDirectory string
	workerArgs, secrets                                                 []string
	maxAttempts                                                         int
	timeout                                                             time.Duration
}
type compiled struct {
	Study          ragcontract.Study          `json:"study"`
	Cells          []ragcontract.ExpandedCell `json:"cells"`
	Specifications []any                      `json:"specifications"`
}

func NewCommand() *cobra.Command {
	root := &cobra.Command{Use: "study", Short: "Validate, explain, compile, and run RAG v2 studies"}
	root.AddCommand(newValidate(), newExplain(), newCompile(), newRun())
	return root
}
func addCommon(command *cobra.Command, flags *commonFlags) {
	command.Flags().StringVar(&flags.inputs, "inputs", "", "RAG input bindings/catalog aliases JSON")
	command.Flags().StringVar(&flags.artifactRoot, "artifact-root", "", "Researchctl artifact root used to stage verified inputs")
	command.Flags().StringVar(&flags.ttcDatabase, "ttc-database", "", "Read-only TTC catalog SQLite database")
	command.Flags().StringVarP(&flags.output, "output", "o", "json", "Output format (json only)")
	_ = command.MarkFlagRequired("inputs")
}
func addRun(command *cobra.Command, flags *runFlags) {
	command.Flags().StringVarP(&flags.project, "project", "p", "project.yaml", "Researchctl project file")
	command.Flags().StringVar(&flags.database, "database", "", "Researchctl laboratory database")
	command.Flags().StringVar(&flags.experiment, "experiment-id", "", "Researchctl experiment receiving runs")
	command.Flags().StringVar(&flags.researchctl, "researchctl-command", "researchctl", "Researchctl executable")
	command.Flags().StringVar(&flags.worker, "worker-command", "rag-worker", "RAG worker executable")
	command.Flags().StringSliceVar(&flags.workerArgs, "worker-arg", nil, "RAG worker argument (repeatable)")
	command.Flags().StringSliceVar(&flags.secrets, "secret-env", nil, "Secret environment variable inherited by the worker")
	command.Flags().IntVar(&flags.maxAttempts, "max-attempts", 1, "Maximum attempts per replicate")
	command.Flags().DurationVar(&flags.timeout, "timeout", 0, "Attempt timeout")
	command.Flags().StringVar(&flags.outputDirectory, "spec-output-dir", "", "Directory for canonical generic specifications")
	_ = command.MarkFlagRequired("experiment-id")
}
func resolve(ctx context.Context, path string, flags commonFlags) (ragcontract.Study, researchctladapter.ResolvedInputs, []ragcontract.ExpandedCell, func(), error) {
	study, err := LoadStudy(path)
	if err != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, func() {}, err
	}
	document, base, err := researchctladapter.LoadInputs(flags.inputs)
	if err != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, func() {}, err
	}
	root := flags.artifactRoot
	cleanup := func() {}
	if root == "" {
		root, err = os.MkdirTemp("", "rag-study-inputs-")
		if err != nil {
			return study, researchctladapter.ResolvedInputs{}, nil, cleanup, err
		}
		cleanup = func() { _ = os.RemoveAll(root) }
	}
	var catalog researchctladapter.CatalogResolver
	if flags.ttcDatabase != "" {
		catalog = researchctladapter.NewTTCCatalog(flags.ttcDatabase)
	}
	resolved, err := researchctladapter.ResolveInputs(ctx, document, base, root, catalog)
	if err != nil {
		cleanup()
		return study, resolved, nil, func() {}, err
	}
	study, cells, err := researchctladapter.Expand(study, resolved)
	return study, resolved, cells, cleanup, err
}
func newValidate() *cobra.Command {
	var flags commonFlags
	command := &cobra.Command{Use: "validate <study.js>", Short: "Deep-validate and expand one RAG v2 study", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		study, _, cells, cleanup, err := resolve(command.Context(), args[0], flags)
		defer cleanup()
		if err != nil {
			return err
		}
		return write(command, map[string]any{"valid": true, "schemaVersion": study.SchemaVersion, "variants": len(study.Variants), "cells": len(cells)})
	}}
	addCommon(command, &flags)
	return command
}
func newExplain() *cobra.Command {
	var flags commonFlags
	command := &cobra.Command{Use: "explain <study.js>", Short: "Explain expanded variants, factors, operators, and cells", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		study, _, cells, cleanup, err := resolve(command.Context(), args[0], flags)
		defer cleanup()
		if err != nil {
			return err
		}
		operators := map[string]bool{}
		for _, variant := range study.Variants {
			for _, node := range variant.Pipeline.Nodes {
				operators[node.Operator.ID()] = true
			}
		}
		ids := make([]string, 0, len(operators))
		for id := range operators {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		return write(command, map[string]any{"schemaVersion": "rag-study-explanation/v2", "name": study.Display.Name, "variants": study.Variants, "factors": study.Factors, "cellCount": len(cells), "operators": ids})
	}}
	addCommon(command, &flags)
	return command
}
func newCompile() *cobra.Command {
	var flags commonFlags
	var directory string
	command := &cobra.Command{Use: "compile <study.js>", Short: "Compile a study into ordered cells and generic researchctl specifications", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		study, resolved, cells, cleanup, err := resolve(command.Context(), args[0], flags)
		defer cleanup()
		if err != nil {
			return err
		}
		specs := make([]any, 0, len(cells))
		for _, cell := range cells {
			spec, err := researchctladapter.WrapExecution(cell.Execution, resolved, study.Display.Name+" / "+cell.VariantID)
			if err != nil {
				return err
			}
			specs = append(specs, spec)
			if directory != "" {
				if err := researchctladapter.WriteSpecification(filepath.Join(directory, spec.ID+".json"), spec); err != nil {
					return err
				}
			}
		}
		return write(command, compiled{Study: study, Cells: cells, Specifications: specs})
	}}
	addCommon(command, &flags)
	command.Flags().StringVar(&directory, "spec-output-dir", "", "Write each canonical generic specification to this directory")
	return command
}
func newRun() *cobra.Command {
	var common commonFlags
	var run runFlags
	command := &cobra.Command{Use: "run <study.js>", Short: "Execute stable study cells and replicates through researchctl", Args: cobra.ExactArgs(1), RunE: func(command *cobra.Command, args []string) error {
		if common.artifactRoot == "" {
			absolute, err := filepath.Abs(run.project)
			if err != nil {
				return err
			}
			common.artifactRoot = filepath.Join(filepath.Dir(absolute), ".researchctl", "artifacts")
		}
		study, resolved, cells, cleanup, err := resolve(command.Context(), args[0], common)
		defer cleanup()
		if err != nil {
			return err
		}
		results := []researchctladapter.RunResult{}
		for _, cell := range cells {
			spec, err := researchctladapter.WrapExecution(cell.Execution, resolved, study.Display.Name+" / "+cell.VariantID)
			if err != nil {
				return err
			}
			replicates := cell.Replicates
			if replicates < 1 {
				replicates = 1
			}
			for replicate := 0; replicate < replicates; replicate++ {
				result, err := researchctladapter.ExecuteSpecification(command.Context(), spec, runOptions(run))
				if err != nil {
					return err
				}
				results = append(results, result)
			}
		}
		return write(command, map[string]any{"study": study.Display.Name, "cellCount": len(cells), "runCount": len(results), "results": results})
	}}
	addCommon(command, &common)
	addRun(command, &run)
	return command
}
func runOptions(flags runFlags) researchctladapter.RunOptions {
	return researchctladapter.RunOptions{ResearchctlCommand: flags.researchctl, Project: flags.project, Database: flags.database, ExperimentID: flags.experiment, Worker: researchctladapter.WorkerCommand{Executable: flags.worker, Args: flags.workerArgs}, MaxAttempts: flags.maxAttempts, Timeout: flags.timeout, SecretEnvironment: flags.secrets, OutputDirectory: flags.outputDirectory}
}
func write(command *cobra.Command, value any) error {
	if strings.ToLower(command.Flag("output").Value.String()) != "json" {
		return fmt.Errorf("only json output is supported")
	}
	encoder := json.NewEncoder(command.OutOrStdout())
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
