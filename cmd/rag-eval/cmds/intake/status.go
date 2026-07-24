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

type StatusCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*StatusCommand)(nil)

type StatusSettings struct {
	DB                 string `glazed:"db"`
	WorkflowDB         string `glazed:"workflow-db"`
	ArtifactRoot       string `glazed:"artifact-root"`
	WorkflowID         string `glazed:"workflow-id"`
	ArgumentWorkflowID string `glazed:"workflow-id-argument"`
	List               bool   `glazed:"list"`
	Status             string `glazed:"status"`
	Limit              int    `glazed:"limit"`
}

func newStatusCommand() *cobra.Command {
	command, err := NewStatusCommand()
	cobra.CheckErr(err)
	result, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return result
}
func NewStatusCommand() (*StatusCommand, error) {
	return &StatusCommand{CommandDescription: cmds.NewCommandDescription("status", cmds.WithShort("Show or list Workflow V3 intake runs"), cmds.WithFlags(fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("RAG domain database")), fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Workflow V3 database")), fields.New("artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow artifact root")), fields.New("workflow-id", fields.TypeString, fields.WithHelp("Run ID")), fields.New("list", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("List runs")), fields.New("status", fields.TypeString, fields.WithHelp("Status filter")), fields.New("limit", fields.TypeInteger, fields.WithDefault(50), fields.WithHelp("List limit"))), cmds.WithArguments(fields.New("workflow-id-argument", fields.TypeString, fields.WithIsArgument(true), fields.WithHelp("Run ID"))))}, nil
}
func (c *StatusCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &StatusSettings{}
	if err := v.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	if s.WorkflowID == "" {
		s.WorkflowID = s.ArgumentWorkflowID
	}
	config := ragintakeworkflow.DefaultConfig(s.DB)
	config.WorkflowDatabase = s.WorkflowDB
	config.ArtifactRoot = s.ArtifactRoot
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()
	var output any
	if s.List {
		output, err = app.ListRuns(ctx, s.Status, s.Limit)
	} else {
		if s.WorkflowID == "" {
			return fmt.Errorf("provide workflow ID or --list")
		}
		output, err = app.Show(ctx, workflowv3.RunID(s.WorkflowID))
	}
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(output)
}
