package serve

import (
	"context"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/spf13/cobra"
)

type ServeCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ServeCommand)(nil)

type ServeSettings struct {
	Address           string `glazed:"address"`
	DB                string `glazed:"db"`
	WorkflowDB        string `glazed:"workflow-db"`
	WorkflowArtifacts string `glazed:"workflow-artifact-root"`
	IndexRoot         string `glazed:"index-root"`
	LogLevel          string `glazed:"log-level"`
}

func NewCommand() *cobra.Command {
	glazedCmd, err := newServeGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	// Override Run to use the actual server (GlazeCommand is just for flag parsing)
	cobraCmd.RunE = runServer
	return cobraCmd
}

func newServeGlazeCommand() (*ServeCommand, error) {
	return &ServeCommand{
		CommandDescription: cmds.NewCommandDescription(
			"serve",
			cmds.WithShort("Start the RAG evaluation HTTP server"),
			cmds.WithLong(`Start the HTTP server with API endpoints and embedded SPA.

The server provides:
- REST API at /api/v1/*
- Web UI at /
- Health check at /api/v1/health

Examples:
  rag-eval serve
  rag-eval serve --address 0.0.0.0:8772
  rag-eval serve --db /path/to/db --workflow-db /path/to/intake-v3.db --log-level debug
`),
			cmds.WithFlags(
				fields.New(
					"address",
					fields.TypeString,
					fields.WithDefault("127.0.0.1:8772"),
					fields.WithHelp("Address to listen on"),
				),
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Path to the Workflow V3 intake SQLite database")),
				fields.New("workflow-artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow V3 intake artifact root")),
				fields.New("index-root", fields.TypeString, fields.WithDefault("data/indexes"), fields.WithHelp("Host BM25 index root")),
				fields.New(
					"log-level",
					fields.TypeString,
					fields.WithDefault("info"),
					fields.WithHelp("Log level (debug, info, warn, error)"),
				),
			),
		),
	}, nil
}

func (c *ServeCommand) RunIntoGlazeProcessor(
	_ context.Context,
	_ *values.Values,
	_ middlewares.Processor,
) error {
	// Never called; we override RunE in the cobra command.
	return nil
}

func runServer(cmd *cobra.Command, args []string) error {
	address, _ := cmd.Flags().GetString("address")
	dbPath, _ := cmd.Flags().GetString("db")
	workflowDB, _ := cmd.Flags().GetString("workflow-db")
	workflowArtifacts, _ := cmd.Flags().GetString("workflow-artifact-root")
	indexRoot, _ := cmd.Flags().GetString("index-root")
	logLevel, _ := cmd.Flags().GetString("log-level")

	return runHTTPServer(cmd, address, dbPath, workflowDB, workflowArtifacts, indexRoot, logLevel)
}
