package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	workflowservice "github.com/go-go-golems/rag-evaluation-system/internal/workflow"
	"github.com/spf13/cobra"
	"io"
	"time"
)

type WorkerCommand struct {
	*cmds.CommandDescription
	runOnce bool
}

var _ cmds.WriterCommand = (*WorkerCommand)(nil)

type WorkerSettings struct {
	EngineDB      string `glazed:"engine-db"`
	WorkerID      string `glazed:"worker-id"`
	MaxWorkers    int    `glazed:"max-workers"`
	PollInterval  string `glazed:"poll-interval"`
	LeaseDuration string `glazed:"lease-duration"`
	Cycles        int    `glazed:"cycles"`
}

func workerFields(includeCycles bool) []*fields.Definition {
	f := []*fields.Definition{fields.New("engine-db", fields.TypeString, fields.WithDefault("state/rag-eval-workflows.db"), fields.WithHelp("Path to the scraper workflow engine SQLite database")), fields.New("worker-id", fields.TypeString, fields.WithDefault("rag-eval-worker"), fields.WithHelp("Worker ID for scraper leases")), fields.New("max-workers", fields.TypeInteger, fields.WithDefault(1), fields.WithHelp("Maximum ops to process per cycle")), fields.New("poll-interval", fields.TypeString, fields.WithDefault("100ms"), fields.WithHelp("Worker poll interval")), fields.New("lease-duration", fields.TypeString, fields.WithDefault("1m"), fields.WithHelp("Op lease duration"))}
	if includeCycles {
		f = append(f, fields.New("cycles", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Finite cycles; 0 runs until interrupted")))
	}
	return f
}
func newRunOnceCommand() *cobra.Command   { return buildWorkerCobra("run-once", true) }
func newRunWorkerCommand() *cobra.Command { return buildWorkerCobra("run-worker", false) }
func buildWorkerCobra(name string, once bool) *cobra.Command {
	c, e := NewWorkerCommand(name, once)
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(e)
	return r
}
func NewWorkerCommand(name string, once bool) (*WorkerCommand, error) {
	return &WorkerCommand{CommandDescription: cmds.NewCommandDescription(name, cmds.WithShort("Run local scraper scheduler worker"), cmds.WithFlags(workerFields(!once)...)), runOnce: once}, nil
}
func (c *WorkerCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &WorkerSettings{}
	if e := v.DecodeSectionInto(schema.DefaultSlug, s); e != nil {
		return e
	}
	poll, e := time.ParseDuration(s.PollInterval)
	if e != nil {
		return e
	}
	lease, e := time.ParseDuration(s.LeaseDuration)
	if e != nil {
		return e
	}
	cfg := workflowservice.WorkerConfig{EngineDB: s.EngineDB, WorkerID: s.WorkerID, MaxWorkers: s.MaxWorkers, PollInterval: poll, LeaseDuration: lease}
	store, sched, e := workflowservice.NewIntakeScheduler(ctx, cfg)
	if e != nil {
		return e
	}
	defer func() { _ = store.Close() }()
	cycles := s.Cycles
	if c.runOnce && cycles == 0 {
		cycles = 1
	}
	if cycles > 0 {
		for i := 1; i <= cycles; i++ {
			r, e := sched.RunOnce(ctx)
			if e != nil {
				return e
			}
			if e = json.NewEncoder(w).Encode(workflowservice.WorkerCycle{Cycle: i, Result: r}); e != nil {
				return e
			}
		}
		return nil
	}
	e = sched.Run(ctx)
	if errors.Is(e, context.Canceled) || errors.Is(e, context.DeadlineExceeded) {
		return nil
	}
	return e
}
