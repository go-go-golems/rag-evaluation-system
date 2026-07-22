package corpus

import (
	"context"
	"os"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/corpussnapshot"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type snapshotTTCCommand struct{ *cmds.CommandDescription }

var _ cmds.GlazeCommand = (*snapshotTTCCommand)(nil)

type snapshotTTCSettings struct {
	DB                        string   `glazed:"db"`
	SourceDB                  string   `glazed:"source-db"`
	SnapshotName              string   `glazed:"snapshot-name"`
	AdditionalSeedDocumentIDs []string `glazed:"seed-document-ids"`
}

func newSnapshotTTCCommand() *cobra.Command {
	command, err := newSnapshotTTCGlazeCommand()
	cobra.CheckErr(err)
	cobraCommand, err := cli.BuildCobraCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}))
	cobra.CheckErr(err)
	return cobraCommand
}

func newSnapshotTTCGlazeCommand() (*snapshotTTCCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "json"})))
	if err != nil {
		return nil, errors.Wrap(err, "create Glazed output schema")
	}
	return &snapshotTTCCommand{CommandDescription: cmds.NewCommandDescription(
		"snapshot-ttc",
		cmds.WithShort("Create or reuse an immutable TTC corpus snapshot"),
		cmds.WithLong(`Read the rich TTC export, reproduce the deterministic baseline selection,
and persist content-addressed document revisions plus ordered corpus membership.

This command writes only immutable corpus tables. It does not update the legacy
documents table, so later chunking and experiment work can reference a stable
snapshot ID without a hidden mutable import.

Example:
  rag-eval corpus snapshot-ttc --source-db data/ttc-wordpress-rag.sqlite --db data/rag-eval.db
`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the RAG SQLite database")),
			fields.New("source-db", fields.TypeString, fields.WithDefault("data/ttc-wordpress-rag.sqlite"), fields.WithHelp("Path to the rebuilt rich TTC SQLite export")),
			fields.New("snapshot-name", fields.TypeString, fields.WithDefault(ttcimport.DefaultSnapshotName), fields.WithHelp("Human-readable baseline snapshot name")),
			fields.New("seed-document-ids", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Additional comma-separated TTC source document IDs that must be selected")),
		), cmds.WithSections(glazedSection),
	)}, nil
}

func (command *snapshotTTCCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, processor middlewares.Processor) error {
	settings := &snapshotTTCSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "decode snapshot-ttc settings")
	}
	plan, err := ttcimport.BuildPlan(ctx, ttcimport.BuildRequest{SourceDBPath: settings.SourceDB, SnapshotName: settings.SnapshotName, IncludeDefaultEvaluationSeeds: true, AdditionalSeedDocumentIDs: normalizeIDs(settings.AdditionalSeedDocumentIDs)})
	if err != nil {
		return err
	}
	info, err := os.Stat(settings.SourceDB)
	if err != nil {
		return errors.Wrap(err, "stat TTC source export")
	}
	queries, err := cmdhelpers.OpenDBAtPath(settings.DB)
	if err != nil {
		return errors.Wrap(err, "open RAG database")
	}
	defer func() { _ = queries.Close() }()
	result, err := corpussnapshot.Persist(ctx, queries, plan, corpussnapshot.PersistRequest{SourceByteSize: info.Size()})
	if err != nil {
		return err
	}
	return processor.AddRow(ctx, types.NewRow(types.MRP("snapshot_id", result.SnapshotID), types.MRP("source_artifact_id", result.SourceArtifactID), types.MRP("document_count", result.DocumentCount), types.MRP("reused", result.Reused)))
}
