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

type CancelCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*CancelCommand)(nil)

type CancelSettings struct {
	DB           string `glazed:"db"`
	WorkflowDB   string `glazed:"workflow-db"`
	ArtifactRoot string `glazed:"artifact-root"`
	RunID        string `glazed:"run-id"`
}

func newCancelCommand() *cobra.Command {
	c, err := NewCancelCommand()
	cobra.CheckErr(err)
	ret, err := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return ret
}
func NewCancelCommand() (*CancelCommand, error) {
	return &CancelCommand{CommandDescription: cmds.NewCommandDescription("cancel", cmds.WithShort("Cancel a Workflow V3 intake run"), cmds.WithFlags(fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("RAG domain database")), fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Workflow V3 database")), fields.New("artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow artifact root"))), cmds.WithArguments(fields.New("run-id", fields.TypeString, fields.WithIsArgument(true), fields.WithRequired(true), fields.WithHelp("Run ID"))))}, nil
}
func (c *CancelCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &CancelSettings{}
	if err := v.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	if s.RunID == "" {
		return fmt.Errorf("run ID required")
	}
	config := ragintakeworkflow.DefaultConfig(s.DB)
	config.WorkflowDatabase = s.WorkflowDB
	config.ArtifactRoot = s.ArtifactRoot
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()
	view, err := app.Cancel(ctx, workflowv3.RunID(s.RunID))
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(view)
}
