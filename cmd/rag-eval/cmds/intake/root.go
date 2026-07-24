package intake

import "github.com/spf13/cobra"

// NewCommand creates the document-intake command group. It is deliberately
// separate from canonical RAG v2 study execution.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intake",
		Short: "Operate the document-intake product pipeline",
		Long:  `Submit document chunk/embed/BM25 intake work and inspect its product-local status. Scientific RAG v2 studies compile with "rag-eval study compile" and execute through Researchctl plus Workflow V3.`,
	}
	cmd.AddCommand(newSubmitIntakeCommand())
	cmd.AddCommand(newRunOnceCommand())
	cmd.AddCommand(newRunWorkerCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newOpsCommand())
	cmd.AddCommand(newCancelCommand())
	return cmd
}
