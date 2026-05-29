package search

import "github.com/spf13/cobra"

// NewCommand creates the `search` command group.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Build and query retrieval indexes",
		Long:  `Build disposable search indexes and run bounded retrieval queries over chunked corpus data.`,
	}
	cmd.AddCommand(newIndexCommand())
	cmd.AddCommand(newQueryCommand())
	cmd.AddCommand(newVectorCommand())
	return cmd
}
