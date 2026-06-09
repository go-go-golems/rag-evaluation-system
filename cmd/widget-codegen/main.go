package main

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/go-go-golems/rag-evaluation-system/internal/widgetmanifest"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "widget-codegen",
		Short:   "Discover, validate, and generate Widget IR registry/codegen artifacts",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return logging.InitLoggerFromCobra(cmd)
		},
	}

	if err := logging.AddLoggingSectionToRootCommand(rootCmd, "widget-codegen"); err != nil {
		cobra.CheckErr(err)
	}

	listCmd, err := newListCobraCommand()
	cobra.CheckErr(err)
	checkCmd, err := newCheckCobraCommand()
	cobra.CheckErr(err)

	rootCmd.AddCommand(listCmd, checkCmd)
	cobra.CheckErr(rootCmd.Execute())
}

type listCommand struct{ *cmds.CommandDescription }

type listSettings struct {
	Root string `glazed:"root"`
	Kind string `glazed:"kind"`
}

func newListCobraCommand() (*cobra.Command, error) {
	glazedCmd, err := newListCommand()
	if err != nil {
		return nil, err
	}
	return cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "widget-codegen"}))
}

func newListCommand() (*listCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &listCommand{CommandDescription: cmds.NewCommandDescription(
		"list",
		cmds.WithShort("List discovered Widget IR manifests"),
		cmds.WithLong(`List discovered colocated widget and recipe manifests.

Examples:
  widget-codegen list .
  widget-codegen list ./... --kind widget --output json
`),
		cmds.WithArguments(
			fields.New("root", fields.TypeString, fields.WithDefault("."), fields.WithHelp("Repository root to scan; './...' is accepted as a shorthand for '.'"), fields.WithIsArgument(true)),
		),
		cmds.WithFlags(
			fields.New("kind", fields.TypeString, fields.WithDefault("all"), fields.WithHelp("Manifest kind to list: all, widget, recipe")),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *listCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &listSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	catalog, err := widgetmanifest.Discover(settings.Root)
	if err != nil {
		return err
	}
	kind := settings.Kind
	if kind == "" {
		kind = "all"
	}
	if kind != "all" && kind != "widget" && kind != "recipe" {
		return fmt.Errorf("unknown kind %q (expected all, widget, or recipe)", kind)
	}
	sort.Slice(catalog.Widgets, func(i, j int) bool {
		left, right := catalog.Widgets[i].Manifest, catalog.Widgets[j].Manifest
		if left.Module != right.Module {
			return left.Module < right.Module
		}
		return left.Type < right.Type
	})
	sort.Slice(catalog.Recipes, func(i, j int) bool {
		left, right := catalog.Recipes[i].Manifest, catalog.Recipes[j].Manifest
		if left.Module != right.Module {
			return left.Module < right.Module
		}
		return left.Name < right.Name
	})
	if kind == "all" || kind == "widget" {
		for _, widget := range catalog.Widgets {
			m := widget.Manifest
			if err := gp.AddRow(ctx, types.NewRow(
				types.MRP("kind", "widget"),
				types.MRP("name", m.Type),
				types.MRP("module", m.Module),
				types.MRP("helper", m.Helper),
				types.MRP("props", m.Props),
				types.MRP("adapter_path", m.Adapter.Path),
				types.MRP("adapter_export", m.Adapter.Export),
				types.MRP("react_component", m.ReactComponent),
				types.MRP("status", m.Status),
				types.MRP("path", rel(catalog.Root, widget.Path)),
			)); err != nil {
				return err
			}
		}
	}
	if kind == "all" || kind == "recipe" {
		for _, recipe := range catalog.Recipes {
			m := recipe.Manifest
			if err := gp.AddRow(ctx, types.NewRow(
				types.MRP("kind", "recipe"),
				types.MRP("name", m.Name),
				types.MRP("module", m.Module),
				types.MRP("input", m.Input),
				types.MRP("output", m.Output),
				types.MRP("implementation_path", m.Implementation.Path),
				types.MRP("implementation_symbol", m.Implementation.Symbol),
				types.MRP("status", m.Status),
				types.MRP("path", rel(catalog.Root, recipe.Path)),
			)); err != nil {
				return err
			}
		}
	}
	return nil
}

type checkCommand struct{ *cmds.CommandDescription }

type checkSettings struct {
	Root string `glazed:"root"`
}

func newCheckCobraCommand() (*cobra.Command, error) {
	glazedCmd, err := newCheckCommand()
	if err != nil {
		return nil, err
	}
	return cli.BuildCobraCommand(glazedCmd, cli.WithParserConfig(cli.CobraParserConfig{AppName: "widget-codegen"}))
}

func newCheckCommand() (*checkCommand, error) {
	glazedSection, err := settings.NewGlazedSchema(settings.WithOutputSectionOptions(schema.WithDefaults(map[string]interface{}{"output": "table"})))
	if err != nil {
		return nil, err
	}
	return &checkCommand{CommandDescription: cmds.NewCommandDescription(
		"check",
		cmds.WithShort("Validate discovered Widget IR manifests"),
		cmds.WithLong(`Validate colocated widget and recipe manifests.

The Phase 1 checker treats missing adapter files as warnings because Phase 2 creates
colocated adapters. Duplicate types/helpers, unknown modules, and malformed manifests
are errors.

Examples:
  widget-codegen check .
  widget-codegen check ./... --output json
`),
		cmds.WithArguments(
			fields.New("root", fields.TypeString, fields.WithDefault("."), fields.WithHelp("Repository root to scan; './...' is accepted as a shorthand for '.'"), fields.WithIsArgument(true)),
		),
		cmds.WithSections(glazedSection),
	)}, nil
}

func (c *checkCommand) RunIntoGlazeProcessor(ctx context.Context, vals *values.Values, gp middlewares.Processor) error {
	settings := &checkSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	catalog, err := widgetmanifest.Discover(settings.Root)
	if err != nil {
		return err
	}
	findings := widgetmanifest.Validate(catalog)
	errors := 0
	for _, finding := range findings {
		if finding.IsError() {
			errors++
		}
		if err := gp.AddRow(ctx, types.NewRow(
			types.MRP("severity", finding.Severity),
			types.MRP("check", finding.Check),
			types.MRP("path", finding.Path),
			types.MRP("subject", finding.Subject),
			types.MRP("message", finding.Message),
		)); err != nil {
			return err
		}
	}
	if errors > 0 {
		return fmt.Errorf("widget manifest check failed with %d error(s)", errors)
	}
	return nil
}

func rel(root, path string) string {
	if root == "" {
		return filepath.ToSlash(path)
	}
	if r, err := filepath.Rel(root, path); err == nil {
		return filepath.ToSlash(r)
	}
	return filepath.ToSlash(path)
}
