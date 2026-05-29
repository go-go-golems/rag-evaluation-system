package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/services/engineview"
	"github.com/spf13/cobra"
)

type opsOptions struct {
	engineDB   string
	workflowID string
}

func newOpsCommand() *cobra.Command {
	opts := &opsOptions{}
	cmd := &cobra.Command{
		Use:   "ops [workflow-id]",
		Short: "List operations for a workflow",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 && opts.workflowID == "" {
				return fmt.Errorf("provide workflow-id argument or --workflow-id")
			}
			if len(args) > 0 {
				opts.workflowID = args[0]
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			service := engineview.NewService(opts.engineDB)
			ops, err := service.WorkflowOps(cmd.Context(), model.WorkflowID(opts.workflowID))
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(ops, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	addEngineDBFlag(cmd, &opts.engineDB)
	cmd.Flags().StringVar(&opts.workflowID, "workflow-id", "", "Workflow ID to inspect")
	return cmd
}
