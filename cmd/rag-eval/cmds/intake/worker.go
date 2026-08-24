package intake

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragintakeworkflow"
	"github.com/spf13/cobra"
)

type WorkerCommand struct {
	*cmds.CommandDescription
	runOnce bool
}

var _ cmds.WriterCommand = (*WorkerCommand)(nil)

type WorkerSettings struct {
	DB                      string `glazed:"db"`
	WorkflowDB              string `glazed:"workflow-db"`
	ArtifactRoot            string `glazed:"artifact-root"`
	IndexRoot               string `glazed:"index-root"`
	PollInterval            string `glazed:"poll-interval"`
	LeaseDuration           string `glazed:"lease-duration"`
	Cycles                  int    `glazed:"cycles"`
	BaseURL                 string `glazed:"base-url"`
	APIKey                  string `glazed:"api-key"`
	CacheDirectory          string `glazed:"cache-directory"`
	ProviderAuthorityDigest string `glazed:"provider-authority-digest"`
}

func workerFields(includeCycles bool) []*fields.Definition {
	f := []*fields.Definition{fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("RAG domain database")), fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Workflow V3 database")), fields.New("artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow artifact root")), fields.New("index-root", fields.TypeString, fields.WithDefault("data/indexes"), fields.WithHelp("Host BM25 index root")), fields.New("poll-interval", fields.TypeString, fields.WithDefault("100ms"), fields.WithHelp("Dispatch poll interval")), fields.New("lease-duration", fields.TypeString, fields.WithDefault("30s"), fields.WithHelp("Workflow lease duration")), fields.New("base-url", fields.TypeString, fields.WithHelp("Host-only embedding endpoint")), fields.New("api-key", fields.TypeString, fields.WithHelp("Host-only embedding API key")), fields.New("cache-directory", fields.TypeString, fields.WithDefault("state/embedding-cache"), fields.WithHelp("Host embedding cache directory")), fields.New("provider-authority-digest", fields.TypeString, fields.WithHelp("Exact host provider authority digest"))}
	if includeCycles {
		f = append(f, fields.New("cycles", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Finite dispatch cycles; zero runs until canceled")))
	}
	return f
}
func newRunOnceCommand() *cobra.Command   { return buildWorkerCobra("run-once", true) }
func newRunWorkerCommand() *cobra.Command { return buildWorkerCobra("run-worker", false) }
func buildWorkerCobra(name string, once bool) *cobra.Command {
	command, err := NewWorkerCommand(name, once)
	cobra.CheckErr(err)
	result, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return result
}
func NewWorkerCommand(name string, once bool) (*WorkerCommand, error) {
	return &WorkerCommand{CommandDescription: cmds.NewCommandDescription(name, cmds.WithShort("Run the local Workflow V3 intake dispatcher"), cmds.WithFlags(workerFields(!once)...)), runOnce: once}, nil
}
func (c *WorkerCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &WorkerSettings{}
	if err := v.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	poll, err := time.ParseDuration(s.PollInterval)
	if err != nil {
		return err
	}
	lease, err := time.ParseDuration(s.LeaseDuration)
	if err != nil {
		return err
	}
	config := ragintakeworkflow.DefaultConfig(s.DB)
	config.WorkflowDatabase = s.WorkflowDB
	config.ArtifactRoot = s.ArtifactRoot
	config.Runtime.IndexRoot = s.IndexRoot
	config.Runtime.APIKey = s.APIKey
	config.Runtime.BaseURL = s.BaseURL
	config.Runtime.CacheDirectory = s.CacheDirectory
	config.Runtime.ProviderAuthorityDigest = s.ProviderAuthorityDigest
	config.PollInterval = poll
	config.LeaseDuration = lease
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()
	cycles := s.Cycles
	if c.runOnce && cycles == 0 {
		cycles = 1
	}
	if cycles > 0 {
		encoder := json.NewEncoder(w)
		for i := 1; i <= cycles; i++ {
			lease, dispatchErr := app.Dispatcher.DispatchOnce(ctx)
			if dispatchErr != nil {
				return dispatchErr
			}
			output := map[string]any{"cycle": i, "dispatched": lease != nil}
			if lease != nil {
				output["runId"] = lease.RunID
				output["nodeKey"] = lease.NodeKey
				output["attempt"] = lease.Attempt
				if err := app.Engine.ExecuteLease(ctx, *lease); err != nil {
					return err
				}
			}
			if err := encoder.Encode(output); err != nil {
				return err
			}
		}
		return nil
	}
	err = app.RunWorker(ctx)
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
