package embedding

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
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/spf13/cobra"
)

type CoverageCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*CoverageCommand)(nil)

type CoverageSettings struct {
	DB           string `glazed:"db"`
	StrategyID   string `glazed:"strategy-id"`
	ProviderType string `glazed:"provider-type"`
	Model        string `glazed:"model"`
	Dimensions   int    `glazed:"dimensions"`
}

func newCoverageCommand() *cobra.Command {
	glazedCmd, err := newCoverageGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newCoverageGlazeCommand() (*CoverageCommand, error) {
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

	return &CoverageCommand{
		CommandDescription: cmds.NewCommandDescription(
			"coverage",
			cmds.WithShort("Show stored embedding coverage by source"),
			cmds.WithLong(`Show chunk embedding coverage grouped by document source for one strategy/provider/model identity.

This command does not call an embedding provider. It only counts chunks and stored
embedding rows in SQLite, so it is safe to run before deciding on a larger live
embedding job.`),
			cmds.WithFlags(
				fields.New("db", fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New("strategy-id", fields.TypeString,
					fields.WithHelp("Chunking strategy ID"),
				),
				fields.New("provider-type", fields.TypeString,
					fields.WithHelp("Embedding provider type, for example openai or ollama"),
				),
				fields.New("model", fields.TypeString,
					fields.WithHelp("Embedding model name"),
				),
				fields.New("dimensions", fields.TypeInteger,
					fields.WithHelp("Embedding dimensions"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *CoverageCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &CoverageSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

	service := embeddingservice.NewService(queries)
	result, err := service.Coverage(ctx, embeddingservice.CoverageRequest{
		StrategyID:   s.StrategyID,
		ProviderType: s.ProviderType,
		Model:        s.Model,
		Dimensions:   s.Dimensions,
	})
	if err != nil {
		return err
	}

	for _, item := range result.Items {
		row := types.NewRow(
			types.MRP("strategy_id", result.StrategyID),
			types.MRP("provider_type", result.ProviderType),
			types.MRP("model", result.Model),
			types.MRP("dimensions", result.Dimensions),
			types.MRP("source_id", item.SourceID),
			types.MRP("chunk_count", item.ChunkCount),
			types.MRP("embedded_count", item.EmbeddedCount),
			types.MRP("missing_count", item.MissingCount),
		)
		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
	}
	return nil
}
