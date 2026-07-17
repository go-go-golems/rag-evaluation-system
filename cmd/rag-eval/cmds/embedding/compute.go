package embedding

import (
	"context"
	"strings"

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

type ComputeCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ComputeCommand)(nil)

type ComputeSettings struct {
	DB                string   `glazed:"db"`
	StrategyID        string   `glazed:"strategy-id"`
	SourceIDs         []string `glazed:"source-ids"`
	ProfileRegistries []string `glazed:"profile-registries"`
	Profile           string   `glazed:"profile"`
	BaseProfile       string   `glazed:"base-profile"`
	EmbeddingType     string   `glazed:"embeddings-type"`
	EmbeddingEngine   string   `glazed:"embeddings-engine"`
	Dimensions        int      `glazed:"embeddings-dimensions"`
	APIKey            string   `glazed:"api-key"`
	BaseURL           string   `glazed:"base-url"`
	CacheType         string   `glazed:"cache-type"`
	CacheDirectory    string   `glazed:"cache-directory"`
	BatchSize         int      `glazed:"batch-size"`
	Limit             int      `glazed:"limit"`
	Force             bool     `glazed:"force"`
}

func newComputeCommand() *cobra.Command {
	glazedCmd, err := newComputeGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newComputeGlazeCommand() (*ComputeCommand, error) {
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

	return &ComputeCommand{
		CommandDescription: cmds.NewCommandDescription(
			"compute",
			cmds.WithShort("Compute embeddings for chunks in a strategy"),
			cmds.WithLong(`Compute embeddings for all chunks belonging to one chunking strategy.

The command uses Geppetto embedding providers. It can resolve a provider from
Pinocchio/Geppetto profiles or from direct embedding flags.

Examples:
  rag-eval embedding compute --strategy-id fixed-300-50 --embeddings-type ollama --embeddings-engine nomic-embed-text --embeddings-dimensions 768 --limit 10
  rag-eval embedding compute --strategy-id fixed-300-50 --profile my-embedding-profile --profile-registries ~/.config/pinocchio/profiles.yaml
`),
			cmds.WithFlags(
				fields.New("db", fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New("strategy-id", fields.TypeString,
					fields.WithHelp("Chunking strategy ID whose chunks should be embedded"),
				),
				fields.New("source-ids", fields.TypeStringList,
					fields.WithHelp("Optional source IDs to restrict chunks, for example ttc-dump-guides,ttc-dump-products"),
				),
				fields.New("profile-registries", fields.TypeStringList,
					fields.WithHelp("Profile registry sources (yaml/sqlite/sqlite-dsn); defaults to ~/.config/pinocchio/profiles.yaml when using profiles"),
				),
				fields.New("profile", fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Embedding-capable profile to resolve"),
				),
				fields.New("base-profile", fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Base profile to overlay direct embedding flags onto"),
				),
				fields.New("embeddings-type", fields.TypeString,
					fields.WithDefault("ollama"),
					fields.WithHelp("Embedding provider type: ollama or openai"),
				),
				fields.New("embeddings-engine", fields.TypeString,
					fields.WithDefault("nomic-embed-text"),
					fields.WithHelp("Embedding model/engine"),
				),
				fields.New("embeddings-dimensions", fields.TypeInteger,
					fields.WithDefault(768),
					fields.WithHelp("Embedding vector dimensions"),
				),
				fields.New("api-key", fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Provider API key (for OpenAI, maps to openai-api-key)"),
				),
				fields.New("base-url", fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Provider base URL (for Ollama, maps to ollama-base-url)"),
				),
				fields.New("cache-type", fields.TypeString,
					fields.WithDefault("none"),
					fields.WithHelp("Geppetto embedding cache type: none, memory, or file"),
				),
				fields.New("cache-directory", fields.TypeString,
					fields.WithDefault("state/embedding-cache"),
					fields.WithHelp("Directory for file-based embedding cache"),
				),
				fields.New("batch-size", fields.TypeInteger,
					fields.WithDefault(16),
					fields.WithHelp("Maximum texts per provider batch"),
				),
				fields.New("limit", fields.TypeInteger,
					fields.WithDefault(0),
					fields.WithHelp("Maximum chunks to consider; 0 means all chunks"),
				),
				fields.New("force", fields.TypeBool,
					fields.WithDefault(false),
					fields.WithHelp("Recompute embeddings even when text_hash is fresh"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *ComputeCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &ComputeSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()

	resolved, err := embeddingservice.ResolveProvider(ctx, embeddingservice.ProviderConfig{
		ProfileRegistries: s.ProfileRegistries,
		Profile:           s.Profile,
		BaseProfile:       s.BaseProfile,
		Type:              s.EmbeddingType,
		Engine:            s.EmbeddingEngine,
		Dimensions:        s.Dimensions,
		APIKey:            s.APIKey,
		BaseURL:           s.BaseURL,
		CacheType:         s.CacheType,
		CacheDirectory:    s.CacheDirectory,
	})
	if err != nil {
		return err
	}
	if resolved.Close != nil {
		defer func() { _ = resolved.Close() }()
	}

	service := embeddingservice.NewService(queries)
	result, err := service.Compute(ctx, embeddingservice.ComputeRequest{
		StrategyID:   s.StrategyID,
		SourceIDs:    s.SourceIDs,
		Provider:     resolved.Provider,
		ProviderType: resolved.ProviderType,
		BatchSize:    s.BatchSize,
		Limit:        s.Limit,
		Force:        s.Force,
	})
	if err != nil {
		return err
	}

	row := types.NewRow(
		types.MRP("strategy_id", result.StrategyID),
		types.MRP("source_ids", strings.Join(result.SourceIDs, ",")),
		types.MRP("provider_type", result.ProviderType),
		types.MRP("model", result.Model),
		types.MRP("dimensions", result.Dimensions),
		types.MRP("effective_profile", resolved.EffectiveProfile),
		types.MRP("profile_registries", strings.Join(s.ProfileRegistries, ",")),
		types.MRP("considered", result.Considered),
		types.MRP("computed", result.Computed),
		types.MRP("skipped_fresh", result.SkippedFresh),
	)
	return gp.AddRow(ctx, row)
}
