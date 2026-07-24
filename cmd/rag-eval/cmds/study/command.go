package study

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
	"github.com/spf13/cobra"
)

type studyCommand struct {
	*cmds.CommandDescription
	action string
	writer io.Writer
}

var _ cmds.WriterCommand = (*studyCommand)(nil)

type settings struct {
	StudyPath       string `glazed:"study"`
	Inputs          string `glazed:"inputs"`
	ArtifactRoot    string `glazed:"artifact-root"`
	TTCDatabase     string `glazed:"ttc-database"`
	OutputDir       string `glazed:"output-dir"`
	Experiment      string `glazed:"experiment-id"`
	ProviderConfig  string `glazed:"provider-config"`
	ProviderFixture bool   `glazed:"provider-fixture"`
}

func NewCommand() *cobra.Command {
	root := &cobra.Command{Use: "study", Short: "Validate, explain, and compile RAG v2 studies for Researchctl and Workflow V3"}
	for _, action := range []string{"validate", "explain", "compile"} {
		command, err := newCommand(action)
		cobra.CheckErr(err)
		cobraCommand, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
		cobra.CheckErr(err)
		cobraCommand.PreRunE = func(cmd *cobra.Command, _ []string) error {
			command.writer = cmd.OutOrStdout()
			return nil
		}
		root.AddCommand(cobraCommand)
	}
	return root
}

func newCommand(action string) (*studyCommand, error) {
	definitions := []*fields.Definition{
		fields.New("inputs", fields.TypeString, fields.WithRequired(true), fields.WithHelp("RAG input bindings/catalog aliases JSON")),
		fields.New("ttc-database", fields.TypeString, fields.WithHelp("Read-only TTC catalog SQLite database")),
	}
	if action == "compile" {
		definitions = append(definitions,
			fields.New("artifact-root", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Researchctl artifact root receiving immutable Workflow inputs")),
			fields.New("output-dir", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Workflow bundle directory contained by artifact-root")),
			fields.New("experiment-id", fields.TypeString, fields.WithRequired(true), fields.WithHelp("Researchctl experiment identity written into the plan")),
			fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only provider configuration used to bind provider authority")),
			fields.New("provider-fixture", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Bind deterministic fixture provider authority for tests only")),
		)
	} else {
		definitions = append(definitions, fields.New("artifact-root", fields.TypeString, fields.WithHelp("Temporary or explicit input artifact root")))
	}
	return &studyCommand{
		CommandDescription: cmds.NewCommandDescription(action, cmds.WithShort(action+" RAG v2 study"), cmds.WithFlags(definitions...), cmds.WithArguments(fields.New("study", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Study JavaScript file")))),
		action:             action,
	}, nil
}

func resolve(ctx context.Context, path string, settings *settings) (ragcontract.Study, researchctladapter.ResolvedInputs, []ragcontract.ExpandedCell, string, func(), error) {
	study, err := LoadStudy(path)
	if err != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, "", func() {}, err
	}
	document, base, err := researchctladapter.LoadInputs(settings.Inputs)
	if err != nil {
		return study, researchctladapter.ResolvedInputs{}, nil, "", func() {}, err
	}
	root := settings.ArtifactRoot
	cleanup := func() {}
	if root == "" {
		root, err = os.MkdirTemp("", "rag-study-inputs-")
		if err != nil {
			return study, researchctladapter.ResolvedInputs{}, nil, "", cleanup, err
		}
		cleanup = func() { _ = os.RemoveAll(root) }
	}
	var catalog researchctladapter.CatalogResolver
	if settings.TTCDatabase != "" {
		catalog = researchctladapter.NewTTCCatalog(settings.TTCDatabase)
	}
	resolved, err := researchctladapter.ResolveInputs(ctx, document, base, root, catalog)
	if err != nil {
		cleanup()
		return study, resolved, nil, "", func() {}, err
	}
	study, cells, err := researchctladapter.Expand(study, resolved)
	return study, resolved, cells, root, cleanup, err
}

func (command *studyCommand) RunIntoWriter(ctx context.Context, values_ *values.Values, writer io.Writer) error {
	settings := &settings{}
	if err := values_.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	if settings.ProviderConfig != "" && settings.ProviderFixture {
		return fmt.Errorf("RAG_STUDY_PROVIDER_FLAGS: provider-config and provider-fixture are mutually exclusive")
	}
	study, resolved, cells, artifactRoot, cleanup, err := resolve(ctx, settings.StudyPath, settings)
	defer cleanup()
	if err != nil {
		return err
	}
	var output any
	switch command.action {
	case "validate":
		output = map[string]any{"valid": true, "schemaVersion": study.SchemaVersion, "variants": len(study.Variants), "cells": len(cells)}
	case "explain":
		operators := map[string]bool{}
		for _, variant := range study.Variants {
			for _, node := range variant.Pipeline.Nodes {
				operators[node.Operator.ID()] = true
			}
		}
		operatorIDs := make([]string, 0, len(operators))
		for id := range operators {
			operatorIDs = append(operatorIDs, id)
		}
		sort.Strings(operatorIDs)
		output = map[string]any{"schemaVersion": "rag-study-explanation/v2", "name": study.Display.Name, "variants": study.Variants, "factors": study.Factors, "cellCount": len(cells), "operators": operatorIDs}
	case "compile":
		corpus, evaluation, err := researchctladapter.LoadDomainArtifacts(artifactRoot, resolved)
		if err != nil {
			return err
		}
		workflowCases := make([]ragworkflow.StudyWorkflowCase, 0, len(cells))
		for index, cell := range cells {
			if err := ragoperators.ValidateInputArtifacts(cell.Execution, corpus, evaluation); err != nil {
				return err
			}
			caseID := fmt.Sprintf("cell-%03d-%s", index, shortIdentity(cell.Execution.CellID))
			workflowCases = append(workflowCases, ragworkflow.StudyWorkflowCase{ID: caseID, Execution: cell.Execution, Corpus: corpus.Corpus, Dataset: evaluation.Dataset, Replicates: maxInt(1, cell.Replicates)})
		}
		providerPackage, closeProviders, err := compileProviderPackage(ctx, settings)
		if err != nil {
			return err
		}
		defer closeProviders()
		outputDirectory, err := filepath.Abs(settings.OutputDir)
		if err != nil {
			return err
		}
		output, err = ragworkflow.WriteStudyBundle(ctx, artifactRoot, outputDirectory, study.Display.Name, settings.Experiment, workflowCases, providerPackage)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("RAG_STUDY_ACTION: %s", command.action)
	}
	if command.writer != nil {
		writer = command.writer
	}
	return json.NewEncoder(writer).Encode(output)
}

func compileProviderPackage(ctx context.Context, settings *settings) (*ragworkflow.ProviderPackage, func(), error) {
	if settings.ProviderConfig == "" && !settings.ProviderFixture {
		return nil, func() {}, nil
	}
	var services ragworkflow.ProviderServices
	closeProviders := func() {}
	var err error
	if settings.ProviderFixture {
		services, err = ragworkflow.NewDeterministicProviderServices()
	} else {
		providerSet, loadErr := ragproviders.Load(ctx, settings.ProviderConfig)
		if loadErr != nil {
			return nil, closeProviders, loadErr
		}
		closeProviders = func() { _ = providerSet.Close() }
		services, err = ragworkflow.ProviderServicesFromSet(providerSet)
	}
	if err != nil {
		closeProviders()
		return nil, func() {}, err
	}
	providerPackage, err := ragworkflow.NewProviderPackage(services, ragworkflowops.Policy{MaxPerAttempt: 10_000, FinishTimeout: 5 * time.Second})
	if err != nil {
		closeProviders()
		return nil, func() {}, err
	}
	return providerPackage, closeProviders, nil
}

func shortIdentity(identity string) string {
	identity = strings.TrimPrefix(identity, "sha256:")
	if len(identity) > 12 {
		return identity[:12]
	}
	return identity
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
