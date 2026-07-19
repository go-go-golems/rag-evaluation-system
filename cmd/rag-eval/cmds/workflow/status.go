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

type StatusCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*StatusCommand)(nil)

type StatusSettings struct {
	EngineDB           string `glazed:"engine-db"`
	WorkflowID         string `glazed:"workflow-id"`
	ArgumentWorkflowID string `glazed:"workflow-id-argument"`
	List               bool   `glazed:"list"`
	Status             string `glazed:"status"`
	Limit              int    `glazed:"limit"`
	Offset             int    `glazed:"offset"`
}

func newStatusCommand() *cobra.Command {
	c, e := NewStatusCommand()
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(e)
	return r
}
func NewStatusCommand() (*StatusCommand, error) {
	return &StatusCommand{CommandDescription: cmds.NewCommandDescription("status", cmds.WithShort("Show workflow status or list workflows"), cmds.WithFlags(fields.New("engine-db", fields.TypeString, fields.WithDefault("state/rag-eval-workflows.db"), fields.WithHelp("Path to the scraper workflow engine SQLite database")), fields.New("workflow-id", fields.TypeString, fields.WithHelp("Workflow ID to inspect")), fields.New("list", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("List workflows instead of showing one workflow")), fields.New("status", fields.TypeString, fields.WithHelp("Optional workflow status filter for --list")), fields.New("limit", fields.TypeInteger, fields.WithDefault(50), fields.WithHelp("Maximum workflows to list")), fields.New("offset", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Workflow list offset"))), cmds.WithArguments(fields.New("workflow-id-argument", fields.TypeString, fields.WithIsArgument(true), fields.WithHelp("Workflow ID to inspect"))))}, nil
}
func (c *StatusCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &StatusSettings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	if s.WorkflowID == "" {
		s.WorkflowID = s.ArgumentWorkflowID
	}
	service := engineview.NewService(s.EngineDB)
	var value any
	var e error
	if s.List {
		value, e = service.ListWorkflows(ctx, engineview.ListWorkflowsOptions{Site: "rag-eval", Status: model.WorkflowStatus(s.Status), Limit: s.Limit, Offset: s.Offset})
	} else {
		if s.WorkflowID == "" {
			return fmt.Errorf("provide workflow-id argument, --workflow-id, or --list")
		}
		value, e = service.Workflow(ctx, model.WorkflowID(s.WorkflowID))
	}
	if e != nil {
		return e
	}
	return json.NewEncoder(w).Encode(value)
}
