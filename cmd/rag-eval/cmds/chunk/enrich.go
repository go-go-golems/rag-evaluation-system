package chunk

import (
	"encoding/json"
	"fmt"

	cmds2 "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	chunkenrichment "github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	"github.com/spf13/cobra"
)

type enrichOptions struct {
	db            string
	chunkID       string
	strategyID    string
	promptVersion string
	provider      string
	model         string
	force         bool
}

func newEnrichCommand() *cobra.Command {
	opts := &enrichOptions{}
	cmd := &cobra.Command{
		Use:   "enrich",
		Short: "Create a non-destructive chunk enrichment artifact",
		Long: `Create a chunk-level enrichment artifact without modifying canonical chunk text.

This Phase 4 command intentionally supports the deterministic fake provider first.
Live LLM providers are reserved for the bounded live-provider smoke phase.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			queries, err := cmds2.OpenDBAtPath(opts.db)
			if err != nil {
				return err
			}
			defer func() { _ = queries.Close() }()
			if opts.provider != "fake" {
				return fmt.Errorf("unsupported chunk enrichment provider %q; only fake is available before live provider smoke", opts.provider)
			}
			result, err := chunkenrichment.NewService(queries).Enrich(cmd.Context(), chunkenrichment.EnrichRequest{
				ChunkID:       opts.chunkID,
				StrategyID:    opts.strategyID,
				PromptVersion: opts.promptVersion,
				Provider:      chunkenrichment.FakeProvider{ProviderName: opts.provider, ModelName: opts.model},
				Force:         opts.force,
			})
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.db, "db", "data/rag-eval.db", "Path to the SQLite database")
	cmd.Flags().StringVar(&opts.chunkID, "chunk-id", "", "Chunk ID to enrich")
	cmd.Flags().StringVar(&opts.strategyID, "strategy-id", "", "Chunking strategy ID")
	cmd.Flags().StringVar(&opts.promptVersion, "prompt-version", "v1", "Prompt version identity")
	cmd.Flags().StringVar(&opts.provider, "provider", "fake", "Chunk enrichment provider; currently only fake")
	cmd.Flags().StringVar(&opts.model, "model", "fake-chunk-enricher", "Chunk enrichment model identity")
	cmd.Flags().BoolVar(&opts.force, "force", false, "Recompute even if enrichment is fresh")
	_ = cmd.MarkFlagRequired("chunk-id")
	_ = cmd.MarkFlagRequired("strategy-id")
	return cmd
}
