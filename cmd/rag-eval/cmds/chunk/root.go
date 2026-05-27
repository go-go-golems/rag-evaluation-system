package chunk

import (
	"github.com/spf13/cobra"
)

// NewCommand creates the `chunk` command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chunk",
		Short: "Chunk documents for RAG indexing",
		Long:  `Apply chunking strategies to documents, creating text chunks for embedding and search.`,
	}

	cmd.AddCommand(newApplyCommand())
	cmd.AddCommand(newStrategiesCommand())

	return cmd
}
