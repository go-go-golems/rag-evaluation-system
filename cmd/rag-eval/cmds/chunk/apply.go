package chunk

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
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	"github.com/spf13/cobra"
)

type ApplyCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*ApplyCommand)(nil)

type ApplySettings struct {
	DB           string `glazed:"db"`
	DocID        string `glazed:"doc-id"`
	Strategy     string `glazed:"strategy"`
	ChunkSize    int    `glazed:"chunk-size"`
	Overlap      int    `glazed:"overlap"`
	StrategyName string `glazed:"strategy-name"`
	Emit         string `glazed:"emit"`
	PreviewRunes int    `glazed:"preview-runes"`
	Limit        int    `glazed:"limit"`
}

func newApplyCommand() *cobra.Command {
	glazedCmd, err := newApplyGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newApplyGlazeCommand() (*ApplyCommand, error) {
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

	return &ApplyCommand{
		CommandDescription: cmds.NewCommandDescription(
			"apply",
			cmds.WithShort("Apply a chunking strategy to a document"),
			cmds.WithLong(`Apply a chunking strategy to a document and store the resulting chunks.

Available strategies:
  fixed            - Fixed-size character chunks with overlap
  sentence         - Sentence-boundary chunks with overlap
  markdown-heading - Split on markdown headings, sub-chunk large sections

Examples:
  rag-eval chunk apply --doc-id doc-abc --strategy fixed --chunk-size 500 --overlap 50
  rag-eval chunk apply --doc-id doc-abc --strategy sentence --chunk-size 1000
  rag-eval chunk apply --doc-id doc-abc --strategy markdown-heading --emit none --output json
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
					fields.WithHelp("Document ID to chunk"),
				),
				fields.New(
					"strategy",
					fields.TypeString,
					fields.WithDefault("fixed"),
					fields.WithHelp("Chunking strategy (fixed, sentence, markdown-heading)"),
				),
				fields.New(
					"chunk-size",
					fields.TypeInteger,
					fields.WithDefault(500),
					fields.WithHelp("Target chunk size in characters"),
				),
				fields.New(
					"overlap",
					fields.TypeInteger,
					fields.WithDefault(50),
					fields.WithHelp("Overlap between chunks in characters"),
				),
				fields.New(
					"strategy-name",
					fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Human-readable name for this strategy config (auto-generated if empty)"),
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
					fields.WithDefault(80),
					fields.WithHelp("Number of runes to include when --emit preview"),
				),
				fields.New(
					"limit",
					fields.TypeInteger,
					fields.WithDefault(50),
					fields.WithHelp("Maximum number of chunk rows to emit; 0 means no limit"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *ApplyCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &ApplySettings{}
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

	service := chunkservice.NewService(queries)
	result, err := service.Apply(ctx, chunkservice.ApplyRequest{
		DocumentID:   s.DocID,
		Strategy:     s.Strategy,
		ChunkSize:    s.ChunkSize,
		Overlap:      s.Overlap,
		StrategyName: s.StrategyName,
	})
	if err != nil {
		return err
	}

	emitted := 0
	for _, ch := range result.Chunks {
		if s.Limit > 0 && emitted >= s.Limit {
			break
		}

		row := types.NewRow(
			types.MRP("id", ch.ID),
			types.MRP("strategy_id", result.StrategyID),
			types.MRP("chunk_index", ch.ChunkIndex),
			types.MRP("token_count", ch.TokenCount),
			types.MRP("start_offset", ch.StartOffset),
			types.MRP("end_offset", ch.EndOffset),
		)
		switch s.Emit {
		case "full":
			row.Set("text", ch.Text)
		case "preview", "":
			row.Set("text_preview", truncate(ch.Text, s.PreviewRunes))
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

	summaryRow := types.NewRow(
		types.MRP("id", "_summary"),
		types.MRP("strategy_id", result.StrategyID),
		types.MRP("chunk_index", result.ChunkCount),
		types.MRP("token_count", 0),
		types.MRP("start_offset", 0),
		types.MRP("end_offset", 0),
		types.MRP("text_preview", fmt.Sprintf("chunked into %d chunks using %s; emitted %d chunk rows", result.ChunkCount, result.StrategyID, emitted)),
	)
	return gp.AddRow(ctx, summaryRow)
}

func truncate(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}
