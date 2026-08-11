package preview

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	studycmd "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/study"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Command struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*Command)(nil)

type Settings struct {
	Study           string   `glazed:"study"`
	Inputs          string   `glazed:"inputs"`
	ArtifactRoot    string   `glazed:"artifact-root"`
	TTCDatabase     string   `glazed:"ttc-database"`
	Project         string   `glazed:"project"`
	Database        string   `glazed:"database"`
	Experiment      string   `glazed:"experiment-id"`
	Researchctl     string   `glazed:"researchctl-command"`
	Worker          string   `glazed:"worker-command"`
	WorkerArgs      []string `glazed:"worker-arg"`
	Variant         string   `glazed:"variant"`
	Factors         []string `glazed:"factor"`
	Query           string   `glazed:"query"`
	Trace           string   `glazed:"trace"`
	Secrets         []string `glazed:"secret-env"`
	MaxAttempts     int      `glazed:"max-attempts"`
	Timeout         string   `glazed:"timeout"`
	OutputDirectory string   `glazed:"spec-output-dir"`
}

func NewCommand() *cobra.Command {
	c, e := NewGlazedCommand()
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(e)
	return r
}
func NewGlazedCommand() (*Command, error) {
	f := []*fields.Definition{fields.New("inputs", fields.TypeString, fields.WithRequired(true), fields.WithHelp("RAG input bindings/catalog aliases JSON")), fields.New("artifact-root", fields.TypeString, fields.WithHelp("Researchctl artifact root")), fields.New("ttc-database", fields.TypeString, fields.WithHelp("Read-only TTC catalog SQLite database")), fields.New("project", fields.TypeString, fields.WithHelp("Existing researchctl project")), fields.New("database", fields.TypeString, fields.WithHelp("Researchctl laboratory database")), fields.New("experiment-id", fields.TypeString, fields.WithDefault("PREVIEW-001"), fields.WithHelp("Experiment receiving the preview run")), fields.New("researchctl-command", fields.TypeString, fields.WithDefault("researchctl"), fields.WithHelp("Researchctl executable")), fields.New("worker-command", fields.TypeString, fields.WithDefault("rag-worker"), fields.WithHelp("RAG worker executable")), fields.New("worker-arg", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("RAG worker argument")), fields.New("variant", fields.TypeString, fields.WithHelp("Variant ID")), fields.New("factor", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Factor key=value")), fields.New("query", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Preview query")), fields.New("trace", fields.TypeString, fields.WithDefault("full"), fields.WithHelp("Trace presentation policy")), fields.New("secret-env", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Secret environment variable")), fields.New("max-attempts", fields.TypeInteger, fields.WithDefault(1), fields.WithHelp("Maximum attempts")), fields.New("timeout", fields.TypeString, fields.WithDefault("0s"), fields.WithHelp("Attempt timeout")), fields.New("spec-output-dir", fields.TypeString, fields.WithHelp("Canonical specification output directory"))}
	return &Command{CommandDescription: cmds.NewCommandDescription("preview", cmds.WithShort("Run one query from one RAG study cell"), cmds.WithFlags(f...), cmds.WithArguments(fields.New("study", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Study JavaScript file"))))}, nil
}
func (c *Command) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &Settings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	factors, e := parseFactors(s.Factors)
	if e != nil {
		return e
	}
	scratch := ""
	if s.Project == "" {
		s.Project, scratch, e = createScratch(ctx, s.Researchctl)
		if e != nil {
			return e
		}
		s.Experiment = "PREVIEW-001"
	}
	if s.ArtifactRoot == "" {
		a, _ := filepath.Abs(s.Project)
		s.ArtifactRoot = filepath.Join(filepath.Dir(a), ".researchctl", "artifacts")
	}
	study, e := studycmd.LoadStudy(s.Study)
	if e != nil {
		return e
	}
	document, base, e := researchctladapter.LoadInputs(s.Inputs)
	if e != nil {
		return e
	}
	var catalog researchctladapter.CatalogResolver
	if s.TTCDatabase != "" {
		catalog = researchctladapter.NewTTCCatalog(s.TTCDatabase)
	}
	resolved, e := researchctladapter.ResolveInputs(ctx, document, base, s.ArtifactRoot, catalog)
	if e != nil {
		return e
	}
	corpus, ok := resolved.ByRole["corpus"]
	if !ok {
		return fmt.Errorf("preview requires corpus input")
	}
	dataset := ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-queries/v2", Queries: []ragoperators.Query{{ID: "preview-query", Text: s.Query}}}
	envelope := ragoperators.NewEvaluationArtifact(dataset, "preview-query", "preview", "candidate", study.Dataset.RelevanceTarget, corpus.Binding.Digest)
	evaluation, e := researchctladapter.StageEnvelope(researchctladapter.InputReference{Role: "evaluation-dataset", Kind: "manifest-envelope", SchemaVersion: ragcontract.EvaluationManifestSchema}, envelope, s.ArtifactRoot)
	if e != nil {
		return e
	}
	delete(resolved.ByRole, "judgments")
	resolved.ByRole["evaluation-dataset"] = evaluation
	study.Dataset.Split = "preview"
	study.Dataset.Status = "candidate"
	study, cells, e := researchctladapter.Expand(study, resolved)
	if e != nil {
		return e
	}
	cell, e := selectCell(cells, s.Variant, factors)
	if e != nil {
		return e
	}
	spec, e := researchctladapter.WrapExecution(cell.Execution, resolved, "preview / "+cell.VariantID)
	if e != nil {
		return e
	}
	timeout, e := time.ParseDuration(s.Timeout)
	if e != nil {
		return e
	}
	result, e := researchctladapter.ExecuteSpecification(ctx, spec, researchctladapter.RunOptions{ResearchctlCommand: s.Researchctl, Project: s.Project, Database: s.Database, ExperimentID: s.Experiment, Worker: researchctladapter.WorkerCommand{Executable: s.Worker, Args: s.WorkerArgs}, MaxAttempts: s.MaxAttempts, Timeout: timeout, SecretEnvironment: s.Secrets, OutputDirectory: s.OutputDirectory})
	if e != nil {
		return e
	}
	return json.NewEncoder(w).Encode(map[string]any{"schemaVersion": "rag-preview-result/v1", "scratch": scratch, "project": s.Project, "trace": s.Trace, "cell": cell, "result": result})
}
func parseFactors(values []string) (map[string]string, error) {
	r := map[string]string{}
	for _, v := range values {
		p := strings.SplitN(v, "=", 2)
		if len(p) != 2 || p[0] == "" || p[1] == "" {
			return nil, fmt.Errorf("invalid --factor %q", v)
		}
		r[p[0]] = p[1]
	}
	return r, nil
}
func selectCell(cells []ragcontract.ExpandedCell, variant string, factors map[string]string) (ragcontract.ExpandedCell, error) {
	for _, cell := range cells {
		if variant != "" && variant != cell.VariantID {
			continue
		}
		ok := true
		for k, v := range factors {
			found := false
			for _, x := range cell.Factors {
				if x.FactorID == k && x.ValueID == v {
					found = true
				}
			}
			if !found {
				ok = false
			}
		}
		if ok {
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
