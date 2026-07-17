package source

import (
	"context"
	"fmt"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	cmds2 "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	sourceservice "github.com/go-go-golems/rag-evaluation-system/internal/services/source"
	"github.com/spf13/cobra"
)

type ScanCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ScanCommand)(nil)

type ScanSettings struct {
	DB       string `glazed:"db"`
	SourceID string `glazed:"source-id"`
	Dir      string `glazed:"dir"`
}

func newScanCommand() *cobra.Command {
	glazedCmd, err := newScanGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newScanGlazeCommand() (*ScanCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(
		settings.WithOutputSectionOptions(
			schema.WithDefaults(map[string]interface{}{
				"output": "table",
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return &ScanCommand{
		CommandDescription: cmds.NewCommandDescription(
			"scan",
			cmds.WithShort("Scan a directory and ingest documents"),
			cmds.WithLong(`Scan a filesystem directory and ingest all text-readable files as documents.

Walks the directory tree, reads each text file, computes word count and
title (from markdown heading or filename), and inserts into the database.

Hidden directories (starting with .) are skipped. Only known text file
extensions are ingested.

Examples:
  rag-eval source scan --source-id my-docs --dir ./docs
  rag-eval source scan --source-id src-code --dir ./src --output json
`),
			cmds.WithFlags(
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New(
					"source-id",
					fields.TypeString,
					fields.WithHelp("ID of the source to associate documents with"),
				),
				fields.New(
					"dir",
					fields.TypeString,
					fields.WithHelp("Directory to scan for documents"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *ScanCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &ScanSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	if s.SourceID == "" {
		return fmt.Errorf("--source-id is required")
	}
	if s.Dir == "" {
		return fmt.Errorf("--dir is required")
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()

	service := sourceservice.NewService(queries)
	result, err := service.Scan(ctx, sourceservice.ScanRequest{SourceID: s.SourceID, Dir: s.Dir})
	if err != nil {
		return err
	}

	for _, id := range result.Documents {
		row := types.NewRow(
			types.MRP("id", id),
			types.MRP("source_id", s.SourceID),
			types.MRP("status", "ingested"),
		)
		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
	}

	summaryRow := types.NewRow(
		types.MRP("id", "_summary"),
		types.MRP("source_id", s.SourceID),
		types.MRP("status", fmt.Sprintf("scanned %d documents", result.DocumentCount)),
	)
	return gp.AddRow(ctx, summaryRow)
}
