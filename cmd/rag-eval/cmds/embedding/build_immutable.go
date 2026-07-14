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
	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableembedding"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type buildImmutableCommand struct{ *cmds.CommandDescription }

var _ cmds.GlazeCommand = (*buildImmutableCommand)(nil)

type buildImmutableSettings struct {
	DB              string `glazed:"db"`
	ChunkSetID      string `glazed:"chunk-set-id"`
	EmbeddingType   string `glazed:"embeddings-type"`
	EmbeddingEngine string `glazed:"embeddings-engine"`
	BaseURL         string `glazed:"base-url"`
	CacheType       string `glazed:"cache-type"`
	CacheDirectory  string `glazed:"cache-directory"`
	Dimensions      int    `glazed:"embeddings-dimensions"`
	BatchSize       int    `glazed:"batch-size"`
}

func newBuildImmutableCommand() *cobra.Command {
	command, err := newBuildImmutableGlazeCommand()
	cobra.CheckErr(err)
	cobraCommand, err := cli.BuildCobraCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCommand
}
func newBuildImmutableGlazeCommand() (*buildImmutableCommand, error) {
	output, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "json"})))
	if err != nil {
		return nil, err
	}
	return &buildImmutableCommand{CommandDescription: cmds.NewCommandDescription("build-immutable", cmds.WithShort("Build or reuse an immutable embedding set"), cmds.WithLong("Embed every chunk in an immutable chunk set with a resolved Geppetto provider. An identical chunk set and resolved model identity reuse the stored vectors."), cmds.WithFlags(fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db")), fields.New("chunk-set-id", fields.TypeString, fields.WithHelp("Immutable chunk set ID")), fields.New("embeddings-type", fields.TypeString, fields.WithDefault("ollama")), fields.New("embeddings-engine", fields.TypeString, fields.WithDefault("nomic-embed-text")), fields.New("embeddings-dimensions", fields.TypeInteger, fields.WithDefault(768)), fields.New("base-url", fields.TypeString, fields.WithDefault("")), fields.New("cache-type", fields.TypeString, fields.WithDefault("file")), fields.New("cache-directory", fields.TypeString, fields.WithDefault("state/embedding-cache")), fields.New("batch-size", fields.TypeInteger, fields.WithDefault(16))), cmds.WithSections(output))}, nil
}
func (c *buildImmutableCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, processor middlewares.Processor) error {
	s := &buildImmutableSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	q, err := cmdhelpers.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer q.Close()
	resolved, err := embeddingservice.ResolveProvider(ctx, embeddingservice.ProviderConfig{Type: s.EmbeddingType, Engine: s.EmbeddingEngine, Dimensions: s.Dimensions, BaseURL: s.BaseURL, CacheType: s.CacheType, CacheDirectory: s.CacheDirectory})
	if err != nil {
		return errors.Wrap(err, "resolve embedding provider")
	}
	if resolved.Close != nil {
		defer func() { _ = resolved.Close() }()
	}
	r, err := immutableembedding.Build(ctx, q, immutableembedding.Request{ChunkSetID: s.ChunkSetID, ProviderType: resolved.ProviderType, Provider: resolved.Provider, BatchSize: s.BatchSize})
	if err != nil {
		return err
	}
	return processor.AddRow(ctx, types.NewRow(types.MRP("embedding_plan_id", r.EmbeddingPlanID), types.MRP("embedding_set_id", r.EmbeddingSetID), types.MRP("embedding_count", r.EmbeddingCount), types.MRP("reused", r.Reused), types.MRP("model", resolved.Model.Name), types.MRP("dimensions", resolved.Model.Dimensions)))
}
