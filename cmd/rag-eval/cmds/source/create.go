package source

import (
	"context"
	"encoding/json"

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

type CreateCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*CreateCommand)(nil)

type CreateSettings struct {
	DB     string `glazed:"db"`
	ID     string `glazed:"id"`
	Name   string `glazed:"name"`
	Type   string `glazed:"type"`
	Config string `glazed:"config"`
}

func newCreateCommand() *cobra.Command {
	glazedCmd, err := newCreateGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newCreateGlazeCommand() (*CreateCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(
		settings.WithOutputSectionOptions(
			schema.WithDefaults(map[string]interface{}{
				"output": "json",
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	return &CreateCommand{
		CommandDescription: cmds.NewCommandDescription(
			"create",
			cmds.WithShort("Create a new source"),
			cmds.WithLong(`Create a new document source for ingesting files.

Examples:
  rag-eval source create --id my-docs --name "My Documents" --type filesystem
  rag-eval source create --id web-data --name "Web Data" --type url --config '{"base_url":"https://example.com"}'
`),
			cmds.WithFlags(
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New(
					"id",
					fields.TypeString,
					fields.WithHelp("Unique identifier for the source"),
				),
				fields.New(
					"name",
					fields.TypeString,
					fields.WithHelp("Human-readable name for the source"),
				),
				fields.New(
					"type",
					fields.TypeString,
					fields.WithDefault("filesystem"),
					fields.WithHelp("Source type (filesystem, url, api)"),
				),
				fields.New(
					"config",
					fields.TypeString,
					fields.WithDefault("{}"),
					fields.WithHelp("JSON configuration for the source"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *CreateCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &CreateSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	// Validate config is valid JSON
	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(s.Config), &configMap); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()

	service := sourceservice.NewService(queries)
	result, err := service.Create(ctx, sourceservice.CreateRequest{
		ID:     s.ID,
		Name:   s.Name,
		Type:   s.Type,
		Config: configMap,
	})
	if err != nil {
		return err
	}

	row := types.NewRow(
		types.MRP("id", result.ID),
		types.MRP("name", result.Name),
		types.MRP("type", result.Type),
		types.MRP("status", result.Status),
	)
	return gp.AddRow(ctx, row)
}
