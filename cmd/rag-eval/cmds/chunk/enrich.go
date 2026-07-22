package chunk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	chunkenrichment "github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	"github.com/spf13/cobra"
)

type EnrichCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*EnrichCommand)(nil)

type EnrichSettings struct {
	DB            string `glazed:"db"`
	ChunkID       string `glazed:"chunk-id"`
	StrategyID    string `glazed:"strategy-id"`
	PromptVersion string `glazed:"prompt-version"`
	Provider      string `glazed:"provider"`
	Model         string `glazed:"model"`
	Force         bool   `glazed:"force"`
}

func newEnrichCommand() *cobra.Command {
	command, err := NewEnrichCommand()
	cobra.CheckErr(err)
	cobraCommand, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return cobraCommand
}

func NewEnrichCommand() (*EnrichCommand, error) {
	return &EnrichCommand{CommandDescription: cmds.NewCommandDescription(
		"enrich",
		cmds.WithShort("Create a non-destructive chunk enrichment artifact"),
		cmds.WithLong("Create a chunk-level enrichment artifact without modifying canonical chunk text. Only the deterministic fake provider is currently supported."),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")),
			fields.New("chunk-id", fields.TypeString, fields.WithHelp("Chunk ID to enrich"), fields.WithRequired(true)),
			fields.New("strategy-id", fields.TypeString, fields.WithHelp("Chunking strategy ID"), fields.WithRequired(true)),
			fields.New("prompt-version", fields.TypeString, fields.WithDefault("v1"), fields.WithHelp("Prompt version identity")),
			fields.New("provider", fields.TypeString, fields.WithDefault("fake"), fields.WithHelp("Chunk enrichment provider; currently only fake")),
			fields.New("model", fields.TypeString, fields.WithDefault("fake-chunk-enricher"), fields.WithHelp("Chunk enrichment model identity")),
			fields.New("force", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Recompute even if enrichment is fresh")),
		),
	)}, nil
}

func (c *EnrichCommand) RunIntoWriter(ctx context.Context, vals *values.Values, writer io.Writer) error {
	settings := &EnrichSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	queries, err := cmdhelpers.OpenDBAtPath(settings.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()
	if settings.Provider != "fake" {
		return fmt.Errorf("unsupported chunk enrichment provider %q; only fake is available before live provider smoke", settings.Provider)
	}
	result, err := chunkenrichment.NewService(queries).Enrich(ctx, chunkenrichment.EnrichRequest{ChunkID: settings.ChunkID, StrategyID: settings.StrategyID, PromptVersion: settings.PromptVersion, Provider: chunkenrichment.FakeProvider{ProviderName: settings.Provider, ModelName: settings.Model}, Force: settings.Force})
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(result)
}
