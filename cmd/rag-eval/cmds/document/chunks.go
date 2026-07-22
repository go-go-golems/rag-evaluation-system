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

type ChunksCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ChunksCommand)(nil)

type ChunksSettings struct {
	DB           string `glazed:"db"`
	DocID        string `glazed:"doc-id"`
	Emit         string `glazed:"emit"`
	PreviewRunes int    `glazed:"preview-runes"`
	Limit        int    `glazed:"limit"`
}

func newChunksCommand() *cobra.Command {
	glazedCmd, err := newChunksGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newChunksGlazeCommand() (*ChunksCommand, error) {
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

	return &ChunksCommand{
		CommandDescription: cmds.NewCommandDescription(
			"chunks",
			cmds.WithShort("List chunks for a document"),
			cmds.WithLong(`List all chunks for a specific document.

Examples:
  rag-eval document chunks --doc-id doc-abc123
  rag-eval document chunks --doc-id doc-abc123 --emit full --output json
  rag-eval document chunks --doc-id doc-abc123 --limit 10 --preview-runes 160
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
					fields.WithHelp("Document ID to retrieve chunks for"),
				),
				fields.New(
					"emit",
					fields.TypeString,
					fields.WithDefault("preview"),
					fields.WithHelp("Chunk text output mode: preview, full, or none"),
				),
				fields.New(
					"preview-runes",
					fields.TypeInteger,
					fields.WithDefault(120),
					fields.WithHelp("Number of runes to include when --emit preview"),
				),
				fields.New(
					"limit",
					fields.TypeInteger,
					fields.WithDefault(50),
					fields.WithHelp("Maximum number of chunks to emit; 0 means no limit"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *ChunksCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &ChunksSettings{}
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
	defer func() { _ = queries.Close() }()

	service := documentservice.NewService(queries)
	chunks, err := service.Chunks(ctx, documentservice.ChunksRequest{DocumentID: s.DocID})
	if err != nil {
		return err
	}

	emitted := 0
	for _, ch := range chunks {
		if s.Limit > 0 && emitted >= s.Limit {
			break
		}

		row := types.NewRow(
			types.MRP("id", ch.ID),
			types.MRP("strategy_id", ch.StrategyID),
			types.MRP("chunk_index", ch.ChunkIndex),
			types.MRP("token_count", ch.TokenCount),
			types.MRP("start_offset", ch.StartOffset),
			types.MRP("end_offset", ch.EndOffset),
		)
		switch s.Emit {
		case "full":
			row.Set("text", ch.Text)
		case "preview", "":
			row.Set("text_preview", truncateRunes(ch.Text, s.PreviewRunes))
		case "none":
			// omit text
		default:
			return fmt.Errorf("invalid --emit value %q (expected preview, full, or none)", s.Emit)
		}

		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
		emitted++
	}

	if s.Limit > 0 && len(chunks) > emitted {
		return gp.AddRow(ctx, types.NewRow(
			types.MRP("id", "_truncated"),
			types.MRP("strategy_id", ""),
			types.MRP("chunk_index", emitted),
			types.MRP("token_count", 0),
			types.MRP("start_offset", 0),
			types.MRP("end_offset", 0),
			types.MRP("text_preview", fmt.Sprintf("emitted %d of %d chunks; rerun with --limit 0 for all", emitted, len(chunks))),
		))
	}

	return nil
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}
