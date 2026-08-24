package intake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragintakeworkflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/spf13/cobra"
)

type OpsCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*OpsCommand)(nil)

type OpsSettings struct {
	DB                 string `glazed:"db"`
	WorkflowDB         string `glazed:"workflow-db"`
	ArtifactRoot       string `glazed:"artifact-root"`
	WorkflowID         string `glazed:"workflow-id"`
	ArgumentWorkflowID string `glazed:"workflow-id-argument"`
}

func newOpsCommand() *cobra.Command {
	command, err := NewOpsCommand()
	cobra.CheckErr(err)
	result, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return result
}
func NewOpsCommand() (*OpsCommand, error) {
	return &OpsCommand{CommandDescription: cmds.NewCommandDescription("ops", cmds.WithShort("Show Workflow V3 operational and canonical observation evidence"), cmds.WithFlags(fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("RAG domain database")), fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Workflow V3 database")), fields.New("artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow artifact root")), fields.New("workflow-id", fields.TypeString, fields.WithHelp("Run ID"))), cmds.WithArguments(fields.New("workflow-id-argument", fields.TypeString, fields.WithIsArgument(true), fields.WithHelp("Run ID"))))}, nil
}
func (c *OpsCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &OpsSettings{}
	if err := v.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	if s.WorkflowID == "" {
		s.WorkflowID = s.ArgumentWorkflowID
	}
	if s.WorkflowID == "" {
		return fmt.Errorf("provide workflow ID")
	}
	config := ragintakeworkflow.DefaultConfig(s.DB)
	config.WorkflowDatabase = s.WorkflowDB
	config.ArtifactRoot = s.ArtifactRoot
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()
	runID := workflowv3.RunID(s.WorkflowID)
	view, err := app.Show(ctx, runID)
	if err != nil {
		return err
	}
	observations, err := app.Observations(ctx, runID)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(map[string]any{"run": view, "observations": observations})
}
