package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflow"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
	"github.com/go-go-golems/scraper/pkg/workflowv3product"
)

func main() {
	config := researchrunner.DefaultConfig()
	config.StateRoot = "state/rag-workflow-runner"
	config.ArtifactRoot = "state/rag-workflow-runner-artifacts"
	config.TaskPackages = []string{ragworkflow.PackageName}
	config.AvailableTaskPackages = []workflowv3product.TaskPackage{ragworkflow.NewPackage()}
	config.DomainProjector = ragworkflow.DomainProjector{}
	config.Capacities = map[string]int{"cpu.rag.prepare": 1, "cpu.rag.query": 4, "cpu.rag.reduce": 1}
	var capacities capacityFlag
	flag.StringVar(&config.StateRoot, "state-root", config.StateRoot, "durable RAG Workflow V3 runner state root")
	flag.StringVar(&config.ArtifactRoot, "artifact-root", config.ArtifactRoot, "RAG Workflow V3 execution artifact root")
	flag.Var(&capacities, "capacity", "resource capacity as name=count (repeatable)")
	flag.DurationVar(&config.LeaseDuration, "lease-duration", config.LeaseDuration, "Workflow V3 lease duration")
	flag.DurationVar(&config.PollInterval, "poll-interval", config.PollInterval, "Workflow V3 dispatch poll interval")
	flag.DurationVar(&config.CancellationTimeout, "cancellation-timeout", config.CancellationTimeout, "bounded workflow cancellation acknowledgement timeout")
	flag.Int64Var(&config.MaxRequestBytes, "max-request-bytes", config.MaxRequestBytes, "maximum protocol request bytes")
	flag.Int64Var(&config.MaxExportBytes, "max-export-bytes", config.MaxExportBytes, "maximum bytes for each exported artifact")
	flag.Parse()
	if flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "runner does not accept positional arguments")
		os.Exit(2)
	}
	if len(capacities) > 0 {
		config.Capacities = capacities
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := researchrunner.Run(ctx, os.Stdin, os.Stdout, config); err != nil {
		fmt.Fprintln(os.Stderr, "RAG workflow runner failed")
		os.Exit(1)
	}
}

type capacityFlag map[string]int

func (f *capacityFlag) String() string { return fmt.Sprint(map[string]int(*f)) }
func (f *capacityFlag) Set(value string) error {
	name, raw, ok := strings.Cut(value, "=")
	if !ok || strings.TrimSpace(name) == "" {
		return fmt.Errorf("capacity must be name=count")
	}
	count, err := strconv.Atoi(raw)
	if err != nil || count < 1 {
		return fmt.Errorf("capacity must be positive")
	}
	if *f == nil {
		*f = capacityFlag{}
	}
	if _, duplicate := (*f)[name]; duplicate {
		return fmt.Errorf("capacity %q is configured more than once", name)
	}
	(*f)[name] = count
	return nil
}
