package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	workflowservice "github.com/go-go-golems/rag-evaluation-system/internal/workflow"
	"github.com/spf13/cobra"
)

type workerOptions struct {
	engineDB      string
	workerID      string
	maxWorkers    int
	pollInterval  time.Duration
	leaseDuration time.Duration
	cycles        int
}

func newRunOnceCommand() *cobra.Command {
	opts := &workerOptions{cycles: 1}
	cmd := &cobra.Command{
		Use:   "run-once",
		Short: "Run one local scraper scheduler cycle for rag-eval intake ops",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCycles(cmd, opts)
		},
	}
	addWorkerFlags(cmd, opts)
	return cmd
}

func newRunWorkerCommand() *cobra.Command {
	opts := &workerOptions{}
	cmd := &cobra.Command{
		Use:   "run-worker",
		Short: "Run a local scraper scheduler worker until interrupted",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.cycles > 0 {
				return runCycles(cmd, opts)
			}
			store, sched, err := workflowservice.NewIntakeScheduler(cmd.Context(), workerConfig(opts))
			if err != nil {
				return err
			}
			defer func() { _ = store.Close() }()
			err = sched.Run(cmd.Context())
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		},
	}
	addWorkerFlags(cmd, opts)
	cmd.Flags().IntVar(&opts.cycles, "cycles", 0, "Optional finite number of cycles before exiting; 0 means run until interrupted")
	return cmd
}

func addWorkerFlags(cmd *cobra.Command, opts *workerOptions) {
	addEngineDBFlag(cmd, &opts.engineDB)
	cmd.Flags().StringVar(&opts.workerID, "worker-id", "rag-eval-worker", "Worker ID for scraper leases")
	cmd.Flags().IntVar(&opts.maxWorkers, "max-workers", 1, "Maximum ops to process per cycle")
	cmd.Flags().DurationVar(&opts.pollInterval, "poll-interval", 100*time.Millisecond, "Worker poll interval")
	cmd.Flags().DurationVar(&opts.leaseDuration, "lease-duration", time.Minute, "Op lease duration")
}

func runCycles(cmd *cobra.Command, opts *workerOptions) error {
	cycles := opts.cycles
	if cycles <= 0 {
		cycles = 1
	}
	store, sched, err := workflowservice.NewIntakeScheduler(cmd.Context(), workerConfig(opts))
	if err != nil {
		return err
	}
	defer func() { _ = store.Close() }()
	for i := 1; i <= cycles; i++ {
		result, err := sched.RunOnce(cmd.Context())
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(workflowservice.WorkerCycle{Cycle: i, Result: result}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(b))
	}
	return nil
}

func workerConfig(opts *workerOptions) workflowservice.WorkerConfig {
	return workflowservice.WorkerConfig{
		EngineDB:      opts.engineDB,
		WorkerID:      opts.workerID,
		MaxWorkers:    opts.maxWorkers,
		PollInterval:  opts.pollInterval,
		LeaseDuration: opts.leaseDuration,
	}
}
