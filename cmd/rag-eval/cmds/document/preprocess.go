package document

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
	documentprocessing "github.com/go-go-golems/rag-evaluation-system/internal/services/documentprocessing"
	"github.com/spf13/cobra"
)

type PreprocessCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*PreprocessCommand)(nil)

type PreprocessSettings struct {
	DB                string   `glazed:"db"`
	DocumentID        string   `glazed:"document-id"`
	ArtifactType      string   `glazed:"artifact-type"`
	PromptVersion     string   `glazed:"prompt-version"`
	Provider          string   `glazed:"provider"`
	Model             string   `glazed:"model"`
	Profile           string   `glazed:"profile"`
	ProfileRegistries []string `glazed:"profile-registries"`
	Force             bool     `glazed:"force"`
}

func newPreprocessCommand() *cobra.Command {
	command, err := NewPreprocessCommand()
	cobra.CheckErr(err)
	cobraCommand, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return cobraCommand
}
func NewPreprocessCommand() (*PreprocessCommand, error) {
	return &PreprocessCommand{CommandDescription: cmds.NewCommandDescription("preprocess", cmds.WithShort("Create a non-destructive document processing artifact"), cmds.WithLong("Create a document-level preprocessing artifact without modifying documents.content_text."), cmds.WithFlags(
		fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("Path to the SQLite database")), fields.New("document-id", fields.TypeString, fields.WithHelp("Document ID to preprocess"), fields.WithRequired(true)), fields.New("artifact-type", fields.TypeString, fields.WithDefault("clean_text"), fields.WithHelp("Artifact type to produce")), fields.New("prompt-version", fields.TypeString, fields.WithDefault("v1"), fields.WithHelp("Prompt version identity")), fields.New("provider", fields.TypeString, fields.WithDefault("fake"), fields.WithHelp("Document processing provider: fake or openai-responses")), fields.New("model", fields.TypeString, fields.WithDefault("fake-document-processor"), fields.WithHelp("Document processing model identity")), fields.New("profile", fields.TypeString, fields.WithHelp("Pinocchio profile slug for live document processing")), fields.New("profile-registries", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Profile registry sources for live document processing")), fields.New("force", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Recompute even if an artifact is fresh"))))}, nil
}
func (c *PreprocessCommand) RunIntoWriter(ctx context.Context, vals *values.Values, writer io.Writer) error {
	s := &PreprocessSettings{}
	if err := vals.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	queries, err := cmdhelpers.OpenDBAtPath(s.DB)
	if err != nil {
		return err
	}
	defer func() { _ = queries.Close() }()
	var provider documentprocessing.Provider
	switch s.Provider {
	case "fake":
		provider = documentprocessing.FakeProvider{ProviderName: s.Provider, ModelName: s.Model}
	case "openai-responses":
		profile := s.Profile
		if profile == "" {
			profile = s.Model
		}
		provider, err = documentprocessing.NewOpenAIResponsesProvider(ctx, profile, s.ProfileRegistries)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported document processing provider %q", s.Provider)
	}
	result, err := documentprocessing.NewService(queries).Process(ctx, documentprocessing.ProcessRequest{DocumentID: s.DocumentID, ArtifactType: s.ArtifactType, PromptVersion: s.PromptVersion, Provider: provider, Force: s.Force})
	if err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(result)
}
