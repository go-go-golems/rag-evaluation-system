package document

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
	documentservice "github.com/go-go-golems/rag-evaluation-system/internal/services/document"
	"github.com/spf13/cobra"
)

type GetCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*GetCommand)(nil)

type GetSettings struct {
	DB    string `glazed:"db"`
	DocID string `glazed:"doc-id"`
}

func newGetCommand() *cobra.Command {
	glazedCmd, err := newGetGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newGetGlazeCommand() (*GetCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(
		settings.WithOutputSectionOptions(
			schema.WithDefaults(map[string]interface{}{
				"output": "yaml",
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return &GetCommand{
		CommandDescription: cmds.NewCommandDescription(
			"get",
			cmds.WithShort("Get document details"),
			cmds.WithLong(`Get detailed information about a specific document.

Examples:
  rag-eval document get --doc-id doc-abc123
  rag-eval document get --doc-id doc-abc123 --output json
`),
			cmds.WithFlags(
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New(
					"doc-id",
					fields.TypeString,
					fields.WithHelp("Document ID to retrieve"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *GetCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &GetSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	if s.DocID == "" {
		return fmt.Errorf("--doc-id is required")
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

	service := documentservice.NewService(queries)
	doc, err := service.Get(ctx, s.DocID)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("document %s not found", s.DocID)
	}

	row := types.NewRow(
		types.MRP("id", doc.ID),
		types.MRP("source_id", doc.SourceID),
		types.MRP("title", doc.Title),
		types.MRP("author", doc.Author),
		types.MRP("url", doc.URL),
		types.MRP("content_type", doc.ContentType),
		types.MRP("word_count", doc.WordCount),
		types.MRP("language", doc.Language),
		types.MRP("status", doc.Status),
		types.MRP("created_at", doc.CreatedAt),
		types.MRP("updated_at", doc.UpdatedAt),
	)
	return gp.AddRow(ctx, row)
}
