package chunk

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
	"github.com/spf13/cobra"
)

type StrategiesCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*StrategiesCommand)(nil)

type StrategiesSettings struct {
	DB string `glazed:"db"`
}

func newStrategiesCommand() *cobra.Command {
	glazedCmd, err := newStrategiesGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newStrategiesGlazeCommand() (*StrategiesCommand, error) {
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

	return &StrategiesCommand{
		CommandDescription: cmds.NewCommandDescription(
			"strategies",
			cmds.WithShort("List chunking strategies"),
			cmds.WithLong(`List all registered chunking strategies.

Examples:
  rag-eval chunk strategies
  rag-eval chunk strategies --output json
`),
			cmds.WithFlags(
				fields.New(
					"db",
					fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *StrategiesCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &StrategiesSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

	rows, err := queries.DB().QueryContext(ctx, `
		SELECT id, name, type, description, created_at
		FROM chunking_strategies ORDER BY created_at DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, name, typ, desc, createdAt string
		if err := rows.Scan(&id, &name, &typ, &desc, &createdAt); err != nil {
			return err
		}
		row := types.NewRow(
			types.MRP("id", id),
			types.MRP("name", name),
			types.MRP("type", typ),
			types.MRP("description", desc),
			types.MRP("created_at", createdAt),
		)
		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
	}

	return rows.Err()
}
