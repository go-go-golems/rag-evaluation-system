package chunk

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
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutablechunk"
	"github.com/spf13/cobra"
)

type buildImmutableCommand struct{ *cmds.CommandDescription }

var _ cmds.GlazeCommand = (*buildImmutableCommand)(nil)

type buildImmutableSettings struct {
	DB         string `glazed:"db"`
	SnapshotID string `glazed:"snapshot-id"`
	Strategy   string `glazed:"strategy"`
	Input      string `glazed:"input-variant"`
	Size       int    `glazed:"chunk-size"`
	Overlap    int    `glazed:"overlap"`
}

func newBuildImmutableCommand() *cobra.Command {
	c, e := newBuildImmutableGlazeCommand()
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(e)
	return r
}
func newBuildImmutableGlazeCommand() (*buildImmutableCommand, error) {
	o, e := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "json"})))
	if e != nil {
		return nil, e
	}
	return &buildImmutableCommand{cmds.NewCommandDescription("build-immutable", cmds.WithShort("Build or reuse an immutable chunk set"), cmds.WithFlags(fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db")), fields.New("snapshot-id", fields.TypeString, fields.WithHelp("Immutable corpus snapshot ID")), fields.New("strategy", fields.TypeString, fields.WithDefault("fixed")), fields.New("input-variant", fields.TypeString, fields.WithDefault("search_text")), fields.New("chunk-size", fields.TypeInteger, fields.WithDefault(1200)), fields.New("overlap", fields.TypeInteger, fields.WithDefault(150))), cmds.WithSections(o))}, nil
}
func (c *buildImmutableCommand) RunIntoGlazeProcessor(ctx context.Context, v *values.Values, p middlewares.Processor) error {
	s := &buildImmutableSettings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	q, e := cmdhelpers.OpenDBAtPath(s.DB)
	if e != nil {
		return e
	}
	defer q.Close()
	r, e := immutablechunk.Build(ctx, q, immutablechunk.Request{CorpusSnapshotID: s.SnapshotID, Strategy: s.Strategy, InputVariant: s.Input, ChunkSize: s.Size, Overlap: s.Overlap})
	if e != nil {
		return e
	}
	return p.AddRow(ctx, types.NewRow(types.MRP("chunk_plan_id", r.ChunkPlanID), types.MRP("chunk_set_id", r.ChunkSetID), types.MRP("chunk_count", r.ChunkCount), types.MRP("reused", r.Reused)))
}
