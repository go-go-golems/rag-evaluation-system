package workflow

import "github.com/spf13/cobra"

func addEngineDBFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "engine-db", "state/rag-eval-workflows.db", "Path to the scraper workflow engine SQLite database")
}
