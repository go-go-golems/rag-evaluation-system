package providers

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
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/spf13/cobra"
)

type ValidateCommand struct{ *cmds.CommandDescription }

type ValidateSettings struct {
	ProviderConfig string `glazed:"provider-config"`
}

func NewValidateCommand() (*cobra.Command, error) {
	glazedSection, err := settings.NewGlazedSchema()
	if err != nil {
		return nil, err
	}
	commandSection, err := cli.NewCommandSettingsSection()
	if err != nil {
		return nil, err
	}
	description := cmds.NewCommandDescription(
		"validate",
		cmds.WithShort("Validate and construct the real-provider host"),
		cmds.WithFlags(fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only real-provider configuration YAML"), fields.WithRequired(true))),
		cmds.WithSections(glazedSection, commandSection),
	)
	command := &ValidateCommand{CommandDescription: description}
	return cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{ShortHelpSections: []string{schema.DefaultSlug}, MiddlewaresFunc: cli.CobraCommandDefaultMiddlewares}))
}

func (c *ValidateCommand) RunIntoGlazeProcessor(ctx context.Context, values *values.Values, processor middlewares.Processor) error {
	settings := &ValidateSettings{}
	if err := values.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	providerSet, err := ragproviders.Load(ctx, settings.ProviderConfig)
	if err != nil {
		return fmt.Errorf("load real provider host: %w", err)
	}
	defer func() { _ = providerSet.Close() }()
	return processor.AddRow(ctx, types.NewRow(types.MRP("profile_id", providerSet.ProfileID), types.MRP("embedding", providerSet.Embedder != nil), types.MRP("reranker", providerSet.Reranker != nil), types.MRP("generation", providerSet.Generator != nil)))
}
