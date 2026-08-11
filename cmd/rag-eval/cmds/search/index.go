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
	searchservice "github.com/go-go-golems/rag-evaluation-system/internal/services/search"
	"github.com/spf13/cobra"
)

type IndexCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*IndexCommand)(nil)

type IndexSettings struct {
	DB         string   `glazed:"db"`
	IndexRoot  string   `glazed:"index-root"`
	IndexID    string   `glazed:"index-id"`
	StrategyID string   `glazed:"strategy-id"`
	SourceIDs  []string `glazed:"source-ids"`
	Force      bool     `glazed:"force"`
	Limit      int      `glazed:"limit"`
}

func newIndexCommand() *cobra.Command {
	glazedCmd, err := newIndexGlazeCommand()
	cobra.CheckErr(err)
	cobraCmd, err := cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCmd
}

func newIndexGlazeCommand() (*IndexCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &IndexCommand{CommandDescription: cmds.NewCommandDescription(
		"index",
		cmds.WithShort("Build or rebuild a BM25 index over chunks"),
		cmds.WithLong(`Build a disposable Bleve/BM25 index from canonical SQLite chunks.

Examples:
  rag-eval search index --strategy-id fixed-1200-150 --source-ids ttc-dump-guides,ttc-dump-articles --index-id bm25-ttc-guides-articles-fixed-1200-150 --force
`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("index-root", fields.TypeString, fields.WithDefault(searchservice.DefaultIndexRoot), fields.WithHelp("Root directory for derived search indexes")),
			fields.New("index-id", fields.TypeString, fields.WithDefault(""), fields.WithHelp("Stable index ID; derived from strategy/source filters when omitted")),
			fields.New("strategy-id", fields.TypeString, fields.WithHelp("Chunking strategy ID to index")),
			fields.New("source-ids", fields.TypeStringList, fields.WithHelp("Optional source IDs to restrict indexed chunks")),
			fields.New("force", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Rebuild an existing index")),
			fields.New("limit", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Maximum chunks to index for smoke builds; 0 means all matching chunks")),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *IndexCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	s := &IndexSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()

	service := searchservice.NewService(queries, s.IndexRoot)
	result, err := service.BuildBM25(ctx, searchservice.BuildIndexRequest{
		IndexID:    s.IndexID,
		StrategyID: s.StrategyID,
		SourceIDs:  s.SourceIDs,
		Force:      s.Force,
		Limit:      s.Limit,
	})
	if err != nil {
		return err
	}
	return gp.AddRow(ctx, types.NewRow(
		types.MRP("index_id", result.IndexID),
		types.MRP("strategy_id", result.StrategyID),
		types.MRP("source_ids", strings.Join(result.SourceIDs, ",")),
		types.MRP("index_path", result.IndexPath),
		types.MRP("document_count", result.DocumentCount),
		types.MRP("chunk_count", result.ChunkCount),
		types.MRP("rebuilt", result.Rebuilt),
	))
}
