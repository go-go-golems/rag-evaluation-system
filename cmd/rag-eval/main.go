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
	"github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/glazedcobra"
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

	// Legacy command behavior is hosted by Glazed command descriptions at the
	// boundary. The adapter is temporary while each implementation is migrated.
	for _, legacy := range []*cobra.Command{
		source.NewCommand(), corpus.NewCommand(), chunk.NewCommand(), document.NewCommand(),
		embedding.NewCommand(), search.NewCommand(), workflowcmd.NewCommand(), serve.NewCommand(),
		study.NewCommand(), preview.NewCommand(),
	} {
		wrapped, err := glazedcobra.WrapTree(legacy)
		cobra.CheckErr(err)
		rootCmd.AddCommand(wrapped)
	}
	providerCommands, err := providers.NewCommand()
	cobra.CheckErr(err)
	rootCmd.AddCommand(providerCommands)

	cobra.CheckErr(rootCmd.Execute())
}
