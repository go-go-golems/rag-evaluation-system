package embedding

import "github.com/spf13/cobra"

// NewCommand creates the `embedding` command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "embedding",
		Short: "Compute and inspect chunk embeddings",
		Long:  `Compute embeddings for chunks and inspect embedding metadata.`,
	}

	cmd.AddCommand(newComputeCommand())
	cmd.AddCommand(newSimilarityCommand())
	return cmd
}
