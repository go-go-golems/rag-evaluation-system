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
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/spf13/cobra"
)

type QueryCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*QueryCommand)(nil)

type QuerySettings struct {
	DB           string `glazed:"db"`
	IndexRoot    string `glazed:"index-root"`
	IndexID      string `glazed:"index-id"`
	Query        string `glazed:"query"`
	Limit        int    `glazed:"limit"`
	PreviewRunes int    `glazed:"preview-runes"`
}

func newQueryCommand() *cobra.Command {
	glazedCmd, err := newQueryGlazeCommand()
	cobra.CheckErr(err)
	cobraCmd, err := cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCmd
}

func newQueryGlazeCommand() (*QueryCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &QueryCommand{CommandDescription: cmds.NewCommandDescription(
		"query",
		cmds.WithShort("Run BM25 lexical search over a built index"),
		cmds.WithLong(`Run bounded BM25 lexical search and emit ranked chunk previews.

Examples:
  rag-eval search query --index-id bm25-ttc-guides-articles-fixed-1200-150 --query "how to plant arborvitae" --limit 10
`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("index-root", fields.TypeString, fields.WithDefault(searchservice.DefaultIndexRoot), fields.WithHelp("Root directory for derived search indexes")),
			fields.New("index-id", fields.TypeString, fields.WithHelp("Index ID to query")),
			fields.New("query", fields.TypeString, fields.WithHelp("Search query text")),
			fields.New("limit", fields.TypeInteger, fields.WithDefault(searchservice.DefaultLimit), fields.WithHelp("Maximum results to emit; capped at 100")),
			fields.New("preview-runes", fields.TypeInteger, fields.WithDefault(searchservice.DefaultPreviewRunes), fields.WithHelp("Maximum runes of chunk text preview; negative disables previews")),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *QueryCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	s := &QuerySettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()

	service := searchservice.NewService(queries, s.IndexRoot)
	result, err := service.QueryBM25(ctx, searchservice.QueryRequest{
		IndexID:      s.IndexID,
		Query:        s.Query,
		Limit:        s.Limit,
		PreviewRunes: s.PreviewRunes,
	})
	if err != nil {
		return err
	}
	for _, item := range result.Items {
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("rank", item.Rank),
			types.MRP("retriever", item.Retriever),
			types.MRP("index_id", result.IndexID),
			types.MRP("query", result.Query),
			types.MRP("score", item.Score),
			types.MRP("chunk_id", item.ChunkID),
			types.MRP("document_id", item.DocumentID),
			types.MRP("source_id", item.SourceID),
			types.MRP("title", item.Title),
			types.MRP("chunk_index", item.ChunkIndex),
			types.MRP("url", item.URL),
			types.MRP("preview", item.Preview),
		)); err != nil {
			return err
		}
	}
	return nil
}
