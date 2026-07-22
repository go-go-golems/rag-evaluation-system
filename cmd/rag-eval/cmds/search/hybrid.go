package search

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
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/spf13/cobra"
)

type HybridCommand struct{ *cmds.CommandDescription }

var _ cmds.GlazeCommand = (*HybridCommand)(nil)

type HybridSettings struct {
	DB                string   `glazed:"db"`
	IndexRoot         string   `glazed:"index-root"`
	IndexID           string   `glazed:"index-id"`
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
	BM25Limit         int      `glazed:"bm25-limit"`
	VectorLimit       int      `glazed:"vector-limit"`
	CandidateLimit    int      `glazed:"candidate-limit"`
	PreviewRunes      int      `glazed:"preview-runes"`
	RRFK              int      `glazed:"rrf-k"`
}

func newHybridCommand() *cobra.Command {
	glazedCmd, err := newHybridGlazeCommand()
	cobra.CheckErr(err)
	cobraCmd, err := cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCmd
}

func newHybridGlazeCommand() (*HybridCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &HybridCommand{CommandDescription: cmds.NewCommandDescription(
		"hybrid",
		cmds.WithShort("Run hybrid BM25 + vector retrieval"),
		cmds.WithLong(`Run BM25 and query-vector retrieval independently, then merge results with reciprocal-rank fusion.`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("index-root", fields.TypeString, fields.WithDefault(searchservice.DefaultIndexRoot), fields.WithHelp("Root directory for BM25 indexes")),
			fields.New("index-id", fields.TypeString, fields.WithHelp("BM25 index ID to query")),
			fields.New("strategy-id", fields.TypeString, fields.WithHelp("Chunking strategy ID for vector candidates")),
			fields.New("source-ids", fields.TypeStringList, fields.WithHelp("Optional source IDs for vector candidates")),
			fields.New("query", fields.TypeString, fields.WithHelp("Search query text")),
			fields.New("profile-registries", fields.TypeStringList, fields.WithHelp("Profile registry sources")),
			fields.New("profile", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Embedding-capable profile to resolve")),
			fields.New("base-profile", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Base profile to overlay direct embedding flags onto")),
			fields.New("embeddings-type", fields.TypeString, fields.WithDefault("ollama"), fields.WithHelp("Embedding provider type")),
			fields.New("embeddings-engine", fields.TypeString, fields.WithDefault("nomic-embed-text"), fields.WithHelp("Embedding model/engine")),
			fields.New("embeddings-dimensions", fields.TypeInteger, fields.WithDefault(768), fields.WithHelp("Embedding vector dimensions")),
			fields.New("api-key", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Provider API key")),
			fields.New("base-url", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Provider base URL")),
			fields.New("cache-type", fields.TypeString, fields.WithDefault("none"), fields.WithHelp("Geppetto embedding cache type")),
			fields.New("cache-directory", fields.TypeString, fields.WithDefault("state/embedding-cache"), fields.WithHelp("Directory for file-based embedding cache")),
			fields.New("limit", fields.TypeInteger, fields.WithDefault(searchservice.DefaultLimit), fields.WithHelp("Maximum hybrid results")),
			fields.New("bm25-limit", fields.TypeInteger, fields.WithDefault(50), fields.WithHelp("BM25 candidate count")),
			fields.New("vector-limit", fields.TypeInteger, fields.WithDefault(50), fields.WithHelp("Vector candidate count")),
			fields.New("candidate-limit", fields.TypeInteger, fields.WithDefault(searchservice.DefaultCandidateLimit), fields.WithHelp("Stored embedding candidate scan limit")),
			fields.New("preview-runes", fields.TypeInteger, fields.WithDefault(searchservice.DefaultPreviewRunes), fields.WithHelp("Maximum preview runes")),
			fields.New("rrf-k", fields.TypeInteger, fields.WithDefault(searchservice.DefaultRRFK), fields.WithHelp("Reciprocal-rank fusion k constant")),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *HybridCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	s := &HybridSettings{}
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

	service := searchservice.NewService(queries, s.IndexRoot)
	result, err := service.QueryHybrid(ctx, searchservice.HybridQueryRequest{
		BM25:   searchservice.QueryRequest{IndexID: s.IndexID, Query: s.Query, Limit: s.BM25Limit, PreviewRunes: s.PreviewRunes},
		Vector: searchservice.VectorQueryRequest{Query: s.Query, StrategyID: s.StrategyID, SourceIDs: s.SourceIDs, Provider: resolved.Provider, ProviderType: resolved.ProviderType, Limit: s.VectorLimit, CandidateLimit: s.CandidateLimit, PreviewRunes: s.PreviewRunes},
		Limit:  s.Limit, RRFK: s.RRFK, BM25Limit: s.BM25Limit, VectorLimit: s.VectorLimit,
	})
	if err != nil {
		return err
	}
	for _, item := range result.Items {
		bm25Rank, bm25Score := component(item, "bm25")
		vectorRank, vectorScore := component(item, "vector")
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("rank", item.Rank), types.MRP("retriever", item.Retriever), types.MRP("query", result.Query), types.MRP("score", item.Score),
			types.MRP("bm25_rank", bm25Rank), types.MRP("bm25_score", bm25Score), types.MRP("vector_rank", vectorRank), types.MRP("vector_score", vectorScore),
			types.MRP("chunk_id", item.ChunkID), types.MRP("document_id", item.DocumentID), types.MRP("source_id", item.SourceID), types.MRP("title", item.Title), types.MRP("chunk_index", item.ChunkIndex), types.MRP("preview", item.Preview),
		)); err != nil {
			return err
		}
	}
	return nil
}

func component(item searchservice.RetrievalResult, name string) (int, float64) {
	if item.Components == nil {
		return 0, 0
	}
	c, ok := item.Components[name]
	if !ok {
		return 0, 0
	}
	return c.Rank, c.Score
}
