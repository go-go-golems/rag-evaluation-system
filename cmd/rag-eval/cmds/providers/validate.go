package providers

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
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/spf13/cobra"
)

type ValidateCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*ValidateCommand)(nil)

type ValidateSettings struct {
	ProviderConfig string `glazed:"provider-config"`
}

func NewValidateCommand() (*cobra.Command, error) {
	description := cmds.NewCommandDescription("validate", cmds.WithShort("Validate and construct the real-provider host"), cmds.WithFlags(fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only real-provider configuration YAML"), fields.WithRequired(true))))
	command := &ValidateCommand{CommandDescription: description}
	return cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{ShortHelpSections: []string{schema.DefaultSlug}}))
}

func (c *ValidateCommand) RunIntoWriter(ctx context.Context, values *values.Values, writer io.Writer) error {
	settings := &ValidateSettings{}
	if err := values.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	providerSet, err := ragproviders.Load(ctx, settings.ProviderConfig)
	if err != nil {
		return fmt.Errorf("load real provider host: %w", err)
	}
	defer func() { _ = providerSet.Close() }()
	return json.NewEncoder(writer).Encode(map[string]any{"profile_id": providerSet.ProfileID, "embedding": providerSet.Embedder != nil, "reranker": providerSet.Reranker != nil, "generation": providerSet.Generator != nil})
}
