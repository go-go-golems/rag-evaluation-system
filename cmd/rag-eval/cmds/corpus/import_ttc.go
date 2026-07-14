package corpus

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	cmdhelpers "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/ttcimport"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type importTTCCommand struct {
	*cmds.CommandDescription
}

var _ cmds.GlazeCommand = (*importTTCCommand)(nil)

type importTTCSettings struct {
	DB                        string   `glazed:"db"`
	SourceDB                  string   `glazed:"source-db"`
	Manifest                  string   `glazed:"manifest"`
	SnapshotName              string   `glazed:"snapshot-name"`
	SourceID                  string   `glazed:"source-id"`
	SourceName                string   `glazed:"source-name"`
	AdditionalSeedDocumentIDs []string `glazed:"seed-document-ids"`
	DryRun                    bool     `glazed:"dry-run"`
}

func newImportTTCCommand() *cobra.Command {
	glazedCommand, err := newImportTTCGlazeCommand()
	cobra.CheckErr(err)
	cobraCommand, err := cli.BuildCobraCommand(glazedCommand,
		cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval"}),
	)
	cobra.CheckErr(err)
	return cobraCommand
}

func newImportTTCGlazeCommand() (*importTTCCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(
		settings.WithOutputSectionOptions(
			schema.WithDefaults(map[string]interface{}{"output": "json"}),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create Glazed output schema")
	}
	return &importTTCCommand{CommandDescription: cmds.NewCommandDescription(
		"import-ttc",
		cmds.WithShort("Import the deterministic TTC baseline selection"),
		cmds.WithLong(`Read the rich TTC WordPress SQLite export, select the 200-document
source-balanced baseline, write a deterministic manifest, and import the selected
documents into the current operational database.

The command includes the source documents named by the candidate evaluation cards
before filling remaining quotas in SHA-256 document-ID order. It is intentionally
the pre-immutable import boundary; later work creates document revisions and an
immutable snapshot identity from this manifest.

Examples:
  rag-eval corpus import-ttc --source-db data/ttc-wordpress-rag.sqlite
  rag-eval corpus import-ttc --dry-run --output json
  rag-eval corpus import-ttc --seed-document-ids wp:123,wp:456
`),
		cmds.WithFlags(
			fields.New("db", fields.TypeString,
				fields.WithDefault("data/rag-eval.db"),
				fields.WithHelp("Path to the operational RAG SQLite database"),
			),
			fields.New("source-db", fields.TypeString,
				fields.WithDefault("data/ttc-wordpress-rag.sqlite"),
				fields.WithHelp("Path to the rebuilt rich TTC SQLite export"),
			),
			fields.New("manifest", fields.TypeString,
				fields.WithDefault("data/manifests/ttc-baseline-v1.json"),
				fields.WithHelp("Path for the deterministic TTC baseline manifest"),
			),
			fields.New("snapshot-name", fields.TypeString,
				fields.WithDefault(ttcimport.DefaultSnapshotName),
				fields.WithHelp("Human-readable baseline snapshot name"),
			),
			fields.New("source-id", fields.TypeString,
				fields.WithDefault(ttcimport.DefaultSourceID),
				fields.WithHelp("Operational source ID for imported TTC documents"),
			),
			fields.New("source-name", fields.TypeString,
				fields.WithDefault("TTC WordPress RAG baseline"),
				fields.WithHelp("Human-readable operational source name"),
			),
			fields.New("seed-document-ids", fields.TypeStringList,
				fields.WithDefault([]string{}),
				fields.WithHelp("Additional comma-separated TTC source document IDs that must be selected"),
			),
			fields.New("dry-run", fields.TypeBool,
				fields.WithDefault(false),
				fields.WithHelp("Build and validate the plan without writing the operational database or manifest"),
			),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (command *importTTCCommand) RunIntoGlazeProcessor(
	ctx context.Context,
	vals *values.Values,
	processor middlewares.Processor,
) error {
	settings := &importTTCSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return errors.Wrap(err, "decode import-ttc settings")
	}
	plan, err := ttcimport.BuildPlan(ctx, ttcimport.BuildRequest{
		SourceDBPath:                  settings.SourceDB,
		SnapshotName:                  settings.SnapshotName,
		IncludeDefaultEvaluationSeeds: true,
		AdditionalSeedDocumentIDs:     normalizeIDs(settings.AdditionalSeedDocumentIDs),
	})
	if err != nil {
		return err
	}

	var result *ttcimport.ImportResult
	if settings.DryRun {
		result = ttcimport.DryRunResult(plan, settings.Manifest)
	} else {
		queries, err := cmdhelpers.OpenDBAtPath(settings.DB)
		if err != nil {
			return errors.Wrap(err, "open operational RAG database")
		}
		defer func() { _ = queries.Close() }()
		result, err = ttcimport.Persist(ctx, queries, plan, ttcimport.PersistRequest{
			SourceID:     settings.SourceID,
			SourceName:   settings.SourceName,
			ManifestPath: settings.Manifest,
		})
		if err != nil {
			return err
		}
	}

	kinds := make([]string, 0, len(result.KindCounts))
	for kind := range result.KindCounts {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	counts := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		counts = append(counts, fmt.Sprintf("%s=%d", kind, result.KindCounts[kind]))
	}
	return processor.AddRow(ctx, types.NewRow(
		types.MRP("snapshot_name", result.SnapshotName),
		types.MRP("document_count", result.DocumentCount),
		types.MRP("kind_counts", strings.Join(counts, ",")),
		types.MRP("source_export_sha256", result.SourceExportSHA256),
		types.MRP("manifest", result.ManifestPath),
		types.MRP("dry_run", result.DryRun),
	))
}

func normalizeIDs(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			if normalized := strings.TrimSpace(part); normalized != "" {
				result = append(result, normalized)
			}
		}
	}
	return result
}
