package study

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type compiled struct {
	Study          ragcontract.Study          `json:"study"`
	Cells          []ragcontract.ExpandedCell `json:"cells"`
	Specifications []any                      `json:"specifications"`
}
type studyCommand struct {
	*cmds.CommandDescription
	action string
	writer io.Writer
}

var _ cmds.WriterCommand = (*studyCommand)(nil)

type settings struct {
	StudyPath     string   `glazed:"study"`
	Inputs        string   `glazed:"inputs"`
	ArtifactRoot  string   `glazed:"artifact-root"`
	TTCDatabase   string   `glazed:"ttc-database"`
	SpecOutputDir string   `glazed:"spec-output-dir"`
	Project       string   `glazed:"project"`
	Database      string   `glazed:"database"`
	Experiment    string   `glazed:"experiment-id"`
	Researchctl   string   `glazed:"researchctl-command"`
	Worker        string   `glazed:"worker-command"`
	WorkerArgs    []string `glazed:"worker-arg"`
	Secrets       []string `glazed:"secret-env"`
	MaxAttempts   int      `glazed:"max-attempts"`
	Timeout       string   `glazed:"timeout"`
}

func NewCommand() *cobra.Command {
	r := &cobra.Command{Use: "study", Short: "Validate, explain, compile, and run RAG v2 studies"}
	for _, a := range []string{"validate", "explain", "compile", "run"} {
		c, e := newCommand(a)
		cobra.CheckErr(e)
		cc, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
		cobra.CheckErr(e)
		cc.PreRunE = func(cmd *cobra.Command, _ []string) error { c.writer = cmd.OutOrStdout(); return nil }
		r.AddCommand(cc)
	}
	return r
}
func newCommand(action string) (*studyCommand, error) {
	f := []*fields.Definition{fields.New("inputs", fields.TypeString, fields.WithRequired(true), fields.WithHelp("RAG input bindings/catalog aliases JSON")), fields.New("artifact-root", fields.TypeString, fields.WithHelp("Researchctl artifact root")), fields.New("ttc-database", fields.TypeString, fields.WithHelp("Read-only TTC catalog SQLite database"))}
	if action == "compile" {
		f = append(f, fields.New("spec-output-dir", fields.TypeString, fields.WithHelp("Write canonical specifications here")))
	}
	if action == "run" {
		f = append(f, fields.New("project", fields.TypeString, fields.WithDefault("project.yaml"), fields.WithHelp("Researchctl project file")), fields.New("database", fields.TypeString, fields.WithHelp("Researchctl laboratory database")), fields.New("experiment-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Researchctl experiment receiving runs")), fields.New("researchctl-command", fields.TypeString, fields.WithDefault("researchctl"), fields.WithHelp("Researchctl executable")), fields.New("worker-command", fields.TypeString, fields.WithDefault("rag-worker"), fields.WithHelp("RAG worker executable")), fields.New("worker-arg", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("RAG worker argument")), fields.New("secret-env", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Secret environment variable")), fields.New("max-attempts", fields.TypeInteger, fields.WithDefault(1), fields.WithHelp("Maximum attempts")), fields.New("timeout", fields.TypeString, fields.WithDefault("0s"), fields.WithHelp("Attempt timeout")), fields.New("spec-output-dir", fields.TypeString, fields.WithHelp("Canonical spec directory")))
	}
	return &studyCommand{CommandDescription: cmds.NewCommandDescription(action, cmds.WithShort(action+" RAG v2 study"), cmds.WithFlags(f...), cmds.WithArguments(fields.New("study", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Study JavaScript file")))), action: action}, nil
}
func resolve(ctx context.Context, path string, s *settings) (ragcontract.Study, researchctladapter.ResolvedInputs, []ragcontract.ExpandedCell, func(), error) {
	study, e := LoadStudy(path)
	if e != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, func() {}, e
	}
	doc, base, e := researchctladapter.LoadInputs(s.Inputs)
	if e != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, func() {}, e
	}
	root := s.ArtifactRoot
	clean := func() {}
	if root == "" {
		root, e = os.MkdirTemp("", "rag-study-inputs-")
		if e != nil {
			return study, researchctladapter.ResolvedInputs{}, nil, clean, e
		}
		clean = func() { _ = os.RemoveAll(root) }
	}
	var catalog researchctladapter.CatalogResolver
	if s.TTCDatabase != "" {
		catalog = researchctladapter.NewTTCCatalog(s.TTCDatabase)
	}
	resolved, e := researchctladapter.ResolveInputs(ctx, doc, base, root, catalog)
	if e != nil {
		clean()
		return study, resolved, nil, func() {}, e
	}
	study, cells, e := researchctladapter.Expand(study, resolved)
	return study, resolved, cells, clean, e
}
func (c *studyCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &settings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	if c.action == "run" && s.ArtifactRoot == "" {
		a, e := filepath.Abs(s.Project)
		if e != nil {
			return e
		}
		s.ArtifactRoot = filepath.Join(filepath.Dir(a), ".researchctl", "artifacts")
	}
	study, resolved, cells, clean, e := resolve(ctx, s.StudyPath, s)
	defer clean()
	if e != nil {
		return e
	}
	var out any
	switch c.action {
	case "validate":
		out = map[string]any{"valid": true, "schemaVersion": study.SchemaVersion, "variants": len(study.Variants), "cells": len(cells)}
	case "explain":
		ops := map[string]bool{}
		for _, x := range study.Variants {
			for _, n := range x.Pipeline.Nodes {
				ops[n.Operator.ID()] = true
			}
		}
		ids := []string{}
		for x := range ops {
			ids = append(ids, x)
		}
		sort.Strings(ids)
		out = map[string]any{"schemaVersion": "rag-study-explanation/v2", "name": study.Display.Name, "variants": study.Variants, "factors": study.Factors, "cellCount": len(cells), "operators": ids}
	default:
		specs := []any{}
		results := []researchctladapter.RunResult{}
		for _, cell := range cells {
			spec, e := researchctladapter.WrapExecution(cell.Execution, resolved, study.Display.Name+" / "+cell.VariantID)
			if e != nil {
				return e
			}
			specs = append(specs, spec)
			if s.SpecOutputDir != "" {
				if e = researchctladapter.WriteSpecification(filepath.Join(s.SpecOutputDir, spec.ID+".json"), spec); e != nil {
					return e
				}
			}
			if c.action == "run" {
				d, e := time.ParseDuration(s.Timeout)
				if e != nil {
					return e
				}
				for i := 0; i < max(1, cell.Replicates); i++ {
					r, e := researchctladapter.ExecuteSpecification(ctx, spec, researchctladapter.RunOptions{ResearchctlCommand: s.Researchctl, Project: s.Project, Database: s.Database, ExperimentID: s.Experiment, Worker: researchctladapter.WorkerCommand{Executable: s.Worker, Args: s.WorkerArgs}, MaxAttempts: s.MaxAttempts, Timeout: d, SecretEnvironment: s.Secrets, OutputDirectory: s.SpecOutputDir})
					if e != nil {
						return e
					}
					results = append(results, r)
				}
			}
		}
		if c.action == "run" {
			out = map[string]any{"study": study.Display.Name, "cellCount": len(cells), "runCount": len(results), "results": results}
		} else {
			out = compiled{study, cells, specs}
		}
	}
	if c.writer != nil {
		w = c.writer
	}
	return json.NewEncoder(w).Encode(out)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var _ = fmt.Sprintf
