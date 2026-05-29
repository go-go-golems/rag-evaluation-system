package workflow

import "github.com/spf13/cobra"

// NewCommand creates the `workflow` command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Submit and operate durable scraper-backed rag-eval workflows",
		Long:  `Submit intake workflows to the scraper durable engine, run local workers, and inspect workflow status/ops.`,
	}
	cmd.AddCommand(newSubmitIntakeCommand())
	cmd.AddCommand(newRunOnceCommand())
	cmd.AddCommand(newRunWorkerCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newOpsCommand())
	return cmd
}
