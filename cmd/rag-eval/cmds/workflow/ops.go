package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/services/engineview"
	"github.com/spf13/cobra"
	"io"
)

type OpsCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*OpsCommand)(nil)

type OpsSettings struct {
	EngineDB           string `glazed:"engine-db"`
	WorkflowID         string `glazed:"workflow-id"`
	ArgumentWorkflowID string `glazed:"workflow-id-argument"`
}

func newOpsCommand() *cobra.Command {
	c, e := NewOpsCommand()
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(e)
	return r
}
func NewOpsCommand() (*OpsCommand, error) {
	return &OpsCommand{CommandDescription: cmds.NewCommandDescription("ops", cmds.WithShort("List operations for a workflow"), cmds.WithFlags(fields.New("engine-db", fields.TypeString, fields.WithDefault("state/rag-eval-workflows.db"), fields.WithHelp("Path to the scraper workflow engine SQLite database")), fields.New("workflow-id", fields.TypeString, fields.WithHelp("Workflow ID to inspect"))), cmds.WithArguments(fields.New("workflow-id-argument", fields.TypeString, fields.WithIsArgument(true), fields.WithHelp("Workflow ID to inspect"))))}, nil
}
func (c *OpsCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &OpsSettings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	if s.WorkflowID == "" {
		s.WorkflowID = s.ArgumentWorkflowID
	}
	if s.WorkflowID == "" {
		return fmt.Errorf("provide workflow-id argument or --workflow-id")
	}
	value, e := engineview.NewService(s.EngineDB).WorkflowOps(ctx, model.WorkflowID(s.WorkflowID))
	if e != nil {
		return e
	}
	return json.NewEncoder(w).Encode(value)
}
