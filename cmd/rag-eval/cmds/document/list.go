package document

import (
	"context"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	cmds2 "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	documentservice "github.com/go-go-golems/rag-evaluation-system/internal/services/document"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ListCommand)(nil)

type ListSettings struct {
	DB     string `glazed:"db"`
	Limit  int    `glazed:"limit"`
	Offset int    `glazed:"offset"`
}

func newListCommand() *cobra.Command {
	glazedCmd, err := newListGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newListGlazeCommand() (*ListCommand, error) {
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

	return &ListCommand{
		CommandDescription: cmds.NewCommandDescription(
			"list",
			cmds.WithShort("List all documents"),
			cmds.WithLong(`List documents in the database with pagination.

Examples:
  rag-eval document list
  rag-eval document list --limit 10 --offset 20
  rag-eval document list --output json
`),
			cmds.WithFlags(
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New(
					"limit",
					fields.TypeInteger,
					fields.WithDefault(50),
					fields.WithHelp("Maximum number of results"),
				),
				fields.New(
					"offset",
					fields.TypeInteger,
					fields.WithDefault(0),
					fields.WithHelp("Offset for pagination"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *ListCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &ListSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

	service := documentservice.NewService(queries)
	docs, err := service.List(ctx, documentservice.ListRequest{Limit: s.Limit, Offset: s.Offset})
	if err != nil {
		return err
	}

	for _, d := range docs {
		row := types.NewRow(
			types.MRP("id", d.ID),
			types.MRP("source_id", d.SourceID),
			types.MRP("title", d.Title),
			types.MRP("content_type", d.ContentType),
			types.MRP("word_count", d.WordCount),
			types.MRP("status", d.Status),
			types.MRP("created_at", d.CreatedAt),
		)
		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
	}

	return nil
}
