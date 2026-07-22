package corpus

import "github.com/spf13/cobra"

// NewCommand creates the `corpus` command group.
func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "corpus",
		Short: "Build and inspect explicit retrieval corpora",
		Long:  "Create deterministic corpus selections that later chunking and experiment runs can reference.",
	}
	command.AddCommand(newImportTTCCommand())
	command.AddCommand(newSnapshotTTCCommand())
	return command
}
