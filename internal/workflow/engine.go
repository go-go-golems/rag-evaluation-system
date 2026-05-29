package workflow

import (
	"context"
	"time"

	"github.com/go-go-golems/scraper/pkg/engine/runner"
	"github.com/go-go-golems/scraper/pkg/engine/scheduler"
	sqlitestore "github.com/go-go-golems/scraper/pkg/engine/store/sqlite"
)

type WorkerConfig struct {
	EngineDB        string
	WorkerID        string
	MaxWorkers      int
	PollInterval    time.Duration
	LeaseDuration   time.Duration
	ResolveProvider ProviderResolver
	IndexRoot       string
}

type WorkerCycle struct {
	Cycle  int `json:"cycle"`
	Result any `json:"result"`
}

func NewIntakeScheduler(ctx context.Context, cfg WorkerConfig) (*sqlitestore.Store, *scheduler.Scheduler, error) {
	if cfg.EngineDB == "" {
		cfg.EngineDB = "state/rag-eval-workflows.db"
	}
	if cfg.WorkerID == "" {
		cfg.WorkerID = "rag-eval-worker"
	}
	if cfg.MaxWorkers <= 0 {
		cfg.MaxWorkers = 1
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 100 * time.Millisecond
	}
	if cfg.LeaseDuration <= 0 {
		cfg.LeaseDuration = time.Minute
	}
	store, err := sqlitestore.Open(ctx, cfg.EngineDB)
	if err != nil {
		return nil, nil, err
	}
	registry := runner.NewRegistry()
	if err := registry.Register(&IntakeRunner{ResolveProvider: cfg.ResolveProvider, IndexRoot: cfg.IndexRoot}); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	sched, err := scheduler.New(store, registry, scheduler.Config{
		MaxWorkers:           cfg.MaxWorkers,
		PollInterval:         cfg.PollInterval,
		DefaultLeaseDuration: cfg.LeaseDuration,
	}, cfg.WorkerID, nil)
	if err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	return store, sched, nil
}
