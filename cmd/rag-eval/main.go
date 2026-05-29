package main

import (
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/chunk"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/document"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/embedding"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/search"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/serve"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/source"
	workflowcmd "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/workflow"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "rag-eval",
		Short:   "RAG Evaluation System — workflow-driven document indexing with interactive playground",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
	}

	if err := logging.AddLoggingSectionToRootCommand(rootCmd, "rag-eval"); err != nil {
		cobra.CheckErr(err)
	}

	// Add command groups
	rootCmd.AddCommand(source.NewCommand())
	rootCmd.AddCommand(chunk.NewCommand())
	rootCmd.AddCommand(document.NewCommand())
	rootCmd.AddCommand(embedding.NewCommand())
	rootCmd.AddCommand(search.NewCommand())
	rootCmd.AddCommand(workflowcmd.NewCommand())
	rootCmd.AddCommand(serve.NewCommand())

	cobra.CheckErr(rootCmd.Execute())
}
