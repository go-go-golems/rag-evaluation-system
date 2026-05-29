package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/services/engineview"
	"github.com/spf13/cobra"
)

type statusOptions struct {
	engineDB   string
	workflowID string
	list       bool
	status     string
	limit      int
	offset     int
}

func newStatusCommand() *cobra.Command {
	opts := &statusOptions{}
	cmd := &cobra.Command{
		Use:   "status [workflow-id]",
		Short: "Show workflow status or list workflows",
		Args: func(cmd *cobra.Command, args []string) error {
			if opts.list {
				return nil
			}
			if len(args) != 1 && opts.workflowID == "" {
				return fmt.Errorf("provide workflow-id argument, --workflow-id, or --list")
			}
			if len(args) > 0 {
				opts.workflowID = args[0]
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			service := engineview.NewService(opts.engineDB)
			var value any
			var err error
			if opts.list {
				value, err = service.ListWorkflows(cmd.Context(), engineview.ListWorkflowsOptions{Site: "rag-eval", Status: model.WorkflowStatus(opts.status), Limit: opts.limit, Offset: opts.offset})
			} else {
				value, err = service.Workflow(cmd.Context(), model.WorkflowID(opts.workflowID))
			}
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	addEngineDBFlag(cmd, &opts.engineDB)
	cmd.Flags().StringVar(&opts.workflowID, "workflow-id", "", "Workflow ID to inspect")
	cmd.Flags().BoolVar(&opts.list, "list", false, "List workflows instead of showing one workflow")
	cmd.Flags().StringVar(&opts.status, "status", "", "Optional workflow status filter for --list")
	cmd.Flags().IntVar(&opts.limit, "limit", 50, "Maximum workflows to list")
	cmd.Flags().IntVar(&opts.offset, "offset", 0, "Workflow list offset")
	return cmd
}
