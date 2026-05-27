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

type SimilarityCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*SimilarityCommand)(nil)

type SimilaritySettings struct {
	DB             string `glazed:"db"`
	StrategyID     string `glazed:"strategy-id"`
	ProviderType   string `glazed:"provider-type"`
	Model          string `glazed:"model"`
	Dimensions     int    `glazed:"dimensions"`
	ChunkIDA       string `glazed:"chunk-id-a"`
	ChunkIDB       string `glazed:"chunk-id-b"`
	Limit          int    `glazed:"limit"`
	CandidateLimit int    `glazed:"candidate-limit"`
	PreviewRunes   int    `glazed:"preview-runes"`
}

func newSimilarityCommand() *cobra.Command {
	glazedCmd, err := newSimilarityGlazeCommand()
	cobra.CheckErr(err)

	cobraCmd, err := cli.BuildCobraCommand(glazedCmd,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)

	return cobraCmd
}

func newSimilarityGlazeCommand() (*SimilarityCommand, error) {
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

	return &SimilarityCommand{
		CommandDescription: cmds.NewCommandDescription(
			"similarity",
			cmds.WithShort("Compare stored chunk embeddings with cosine similarity"),
			cmds.WithLong(`Compare stored chunk embeddings for a shared strategy/provider/model identity.

The command never calls a live embedding provider. It reads vectors already stored in
chunk_embeddings and emits bounded similarity rows. Pass --chunk-id-b for a direct
pairwise comparison, or omit it to compare --chunk-id-a against a bounded candidate
set from the same strategy.

Examples:
  rag-eval embedding similarity --strategy-id fixed-300-50 --provider-type ollama --model nomic-embed-text --dimensions 768 --chunk-id-a chunk-a --chunk-id-b chunk-b
  rag-eval embedding similarity --strategy-id fixed-300-50 --provider-type ollama --model nomic-embed-text --dimensions 768 --chunk-id-a chunk-a --limit 10 --candidate-limit 200
`),
			cmds.WithFlags(
				fields.New("db", fields.TypeString,
					fields.WithDefault("data/rag-eval.db"),
					fields.WithHelp("Path to the SQLite database"),
				),
				fields.New("strategy-id", fields.TypeString,
					fields.WithHelp("Chunking strategy ID for both stored embeddings"),
				),
				fields.New("provider-type", fields.TypeString,
					fields.WithHelp("Stored embedding provider type, for example ollama or openai"),
				),
				fields.New("model", fields.TypeString,
					fields.WithHelp("Stored embedding model name"),
				),
				fields.New("dimensions", fields.TypeInteger,
					fields.WithHelp("Stored embedding dimensions"),
				),
				fields.New("chunk-id-a", fields.TypeString,
					fields.WithHelp("Source chunk ID"),
				),
				fields.New("chunk-id-b", fields.TypeString,
					fields.WithDefault(""),
					fields.WithHelp("Optional target chunk ID for pairwise comparison"),
				),
				fields.New("limit", fields.TypeInteger,
					fields.WithDefault(20),
					fields.WithHelp("Maximum similarity rows to emit when chunk-id-b is omitted"),
				),
				fields.New("candidate-limit", fields.TypeInteger,
					fields.WithDefault(200),
					fields.WithHelp("Maximum candidate embeddings to load when chunk-id-b is omitted"),
				),
				fields.New("preview-runes", fields.TypeInteger,
					fields.WithDefault(120),
					fields.WithHelp("Maximum runes of source/target chunk text to include in rows; 0 disables previews"),
				),
			),
			cmds.WithSections(glazedSection),
		),
	}, nil
}

func (c *SimilarityCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	gp middlewares.Processor,
) error {
	s := &SimilaritySettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}

	queries, err := cmds2.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer queries.Close()

	service := embeddingservice.NewService(queries)
	result, err := service.Similarity(ctx, embeddingservice.SimilarityRequest{
		StrategyID:     s.StrategyID,
		ProviderType:   s.ProviderType,
		Model:          s.Model,
		Dimensions:     s.Dimensions,
		ChunkIDA:       s.ChunkIDA,
		ChunkIDB:       s.ChunkIDB,
		Limit:          s.Limit,
		CandidateLimit: s.CandidateLimit,
		PreviewRunes:   s.PreviewRunes,
	})
	if err != nil {
		return err
	}

	for _, match := range result.Matches {
		row := types.NewRow(
			types.MRP("strategy_id", result.StrategyID),
			types.MRP("provider_type", result.ProviderType),
			types.MRP("model", result.Model),
			types.MRP("dimensions", result.Dimensions),
			types.MRP("source_chunk_id", result.Source.ChunkID),
			types.MRP("source_document_id", result.Source.DocumentID),
			types.MRP("source_chunk_index", result.Source.ChunkIndex),
			types.MRP("target_chunk_id", match.ChunkID),
			types.MRP("target_document_id", match.DocumentID),
			types.MRP("target_chunk_index", match.ChunkIndex),
			types.MRP("score", match.Score),
			types.MRP("considered", result.Considered),
			types.MRP("candidate_limit", result.CandidateLimit),
			types.MRP("source_preview", result.Source.TextPreview),
			types.MRP("target_preview", match.TextPreview),
		)
		if err := gp.AddRow(ctx, row); err != nil {
			return err
		}
	}
	return nil
}
