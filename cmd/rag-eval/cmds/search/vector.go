package search

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
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/spf13/cobra"
)

type VectorCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*VectorCommand)(nil)

type VectorSettings struct {
	DB                string   `glazed:"db"`
	StrategyID        string   `glazed:"strategy-id"`
	SourceIDs         []string `glazed:"source-ids"`
	Query             string   `glazed:"query"`
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
	Limit             int      `glazed:"limit"`
	CandidateLimit    int      `glazed:"candidate-limit"`
	PreviewRunes      int      `glazed:"preview-runes"`
}

func newVectorCommand() *cobra.Command {
	glazedCmd, err := newVectorGlazeCommand()
	cobra.CheckErr(err)
	cobraCmd, err := cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCmd
}

func newVectorGlazeCommand() (*VectorCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &VectorCommand{CommandDescription: cmds.NewCommandDescription(
		"vector",
		cmds.WithShort("Run query-vector search over stored chunk embeddings"),
		cmds.WithLong(`Embed a user query with Geppetto/Pinocchio and compare it to stored chunk embeddings.

Examples:
  rag-eval search vector --query "which trees make a good privacy screen" --strategy-id fixed-1200-150 --profile openai-embedding-small --profile-registries ~/.config/pinocchio/profiles.yaml --limit 10 --candidate-limit 200
`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("strategy-id", fields.TypeString, fields.WithHelp("Chunking strategy ID whose stored embeddings should be searched")),
			fields.New("source-ids", fields.TypeStringList, fields.WithHelp("Optional source IDs to restrict stored embedding candidates")),
			fields.New("query", fields.TypeString, fields.WithHelp("Search query text to embed")),
			fields.New("profile-registries", fields.TypeStringList, fields.WithHelp("Profile registry sources; defaults to ~/.config/pinocchio/profiles.yaml when using profiles")),
			fields.New("profile", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Embedding-capable profile to resolve")),
			fields.New("base-profile", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Base profile to overlay direct embedding flags onto")),
			fields.New("embeddings-type", fields.TypeString, fields.WithDefault("ollama"), fields.WithHelp("Embedding provider type: ollama or openai")),
			fields.New("embeddings-engine", fields.TypeString, fields.WithDefault("nomic-embed-text"), fields.WithHelp("Embedding model/engine")),
			fields.New("embeddings-dimensions", fields.TypeInteger, fields.WithDefault(768), fields.WithHelp("Embedding vector dimensions")),
			fields.New("api-key", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Provider API key")),
			fields.New("base-url", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Provider base URL")),
			fields.New("cache-type", fields.TypeString, fields.WithDefault("none"), fields.WithHelp("Geppetto embedding cache type: none, memory, or file")),
			fields.New("cache-directory", fields.TypeString, fields.WithDefault("state/embedding-cache"), fields.WithHelp("Directory for file-based embedding cache")),
			fields.New("limit", fields.TypeInteger, fields.WithDefault(searchservice.DefaultLimit), fields.WithHelp("Maximum results to emit")),
			fields.New("candidate-limit", fields.TypeInteger, fields.WithDefault(searchservice.DefaultCandidateLimit), fields.WithHelp("Maximum stored embeddings to compare")),
			fields.New("preview-runes", fields.TypeInteger, fields.WithDefault(searchservice.DefaultPreviewRunes), fields.WithHelp("Maximum runes of preview text")),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *VectorCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	s := &VectorSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

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

	service := searchservice.NewService(queries, searchservice.DefaultIndexRoot)
	result, err := service.QueryVector(ctx, searchservice.VectorQueryRequest{
		Query:          s.Query,
		StrategyID:     s.StrategyID,
		SourceIDs:      s.SourceIDs,
		Provider:       resolved.Provider,
		ProviderType:   resolved.ProviderType,
		Limit:          s.Limit,
		CandidateLimit: s.CandidateLimit,
		PreviewRunes:   s.PreviewRunes,
	})
	if err != nil {
		return err
	}
	for _, item := range result.Items {
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("rank", item.Rank),
			types.MRP("retriever", item.Retriever),
			types.MRP("query", result.Query),
			types.MRP("score", item.Score),
			types.MRP("chunk_id", item.ChunkID),
			types.MRP("document_id", item.DocumentID),
			types.MRP("source_id", item.SourceID),
			types.MRP("title", item.Title),
			types.MRP("chunk_index", item.ChunkIndex),
			types.MRP("provider_type", resolved.ProviderType),
			types.MRP("model", resolved.Provider.GetModel().Name),
			types.MRP("dimensions", resolved.Provider.GetModel().Dimensions),
			types.MRP("effective_profile", resolved.EffectiveProfile),
			types.MRP("source_ids", strings.Join(s.SourceIDs, ",")),
			types.MRP("preview", item.Preview),
		)); err != nil {
			return err
		}
	}
	return nil
}
