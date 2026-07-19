package main

import (
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/go-go-golems/rag-evaluation-system"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/chunk"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/corpus"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/document"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/embedding"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/preview"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/providers"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/search"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/serve"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/source"
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/study"
	workflowcmd "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/workflow"
	ragdoc "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/doc"
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

	helpSystem := help.NewHelpSystem()
	if err := rageval.AddDocToHelpSystem(helpSystem); err != nil {
		cobra.CheckErr(err)
	}
	if err := ragdoc.AddDocToHelpSystem(helpSystem); err != nil {
		cobra.CheckErr(err)
	}
	help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)

	rootCmd.AddCommand(source.NewCommand())
	rootCmd.AddCommand(corpus.NewCommand())
	rootCmd.AddCommand(chunk.NewCommand())
	rootCmd.AddCommand(document.NewCommand())
	rootCmd.AddCommand(embedding.NewCommand())
	rootCmd.AddCommand(search.NewCommand())
	rootCmd.AddCommand(workflowcmd.NewCommand())
	rootCmd.AddCommand(serve.NewCommand())
	rootCmd.AddCommand(study.NewCommand())
	rootCmd.AddCommand(preview.NewCommand())
	providerCommands, err := providers.NewCommand()
	cobra.CheckErr(err)
	rootCmd.AddCommand(providerCommands)

	cobra.CheckErr(rootCmd.Execute())
}
