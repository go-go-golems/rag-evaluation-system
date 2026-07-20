// Package preparationworkflow adapts immutable RAG combined-preparation plans
// to scraper's durable workflow runtime. It deliberately leaves prompt/model
// resolution and validation in ragoperators.
package preparationworkflow

import (
	"context"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

const (
	PackageName      = "rag-preparation/v1"
	CombinedStepKind = "rag-preparation/combined-batch/v1"
	FinalizeStepKind = "rag-preparation/finalize/v1"
	GenerationQueue  = model.QueueKey("rag:generator")
	LocalQueue       = model.QueueKey("rag:local")
)

type Identity struct {
	SchemaVersion  string `json:"schemaVersion"`
	PreparedDigest string `json:"preparedDigest"`
}

type Input struct {
	Identity Identity                             `json:"identity"`
	Plan     ragoperators.CombinedPreparationPlan `json:"plan"`
}

type batchInput struct {
	Identity Identity                              `json:"identity"`
	Plan     ragoperators.CombinedPreparationPlan  `json:"plan"`
	Batch    ragoperators.CombinedPreparationBatch `json:"batch"`
}

type batchOutput struct {
	Representations []ragoperators.Representation `json:"representations"`
	CacheHit        bool                          `json:"cacheHit"`
	ProviderCall    bool                          `json:"providerCall"`
}

type finalizerInput struct {
	ExpectedBatches int `json:"expectedBatches"`
}

// EnvironmentResolver resolves process-local provider configuration without
// placing credentials in durable scraper operation inputs.
type EnvironmentResolver func(context.Context, Identity) (*ragoperators.Environment, error)

// Register installs the preparation package and its two domain executors into a
// scraper runtime. The caller owns runtime lifetime and scheduling.
func Register(runtime *scraperworkflow.Runtime, resolve EnvironmentResolver) error {
	if runtime == nil || resolve == nil {
		return fmt.Errorf("preparation workflow runtime and environment resolver are required")
	}
	if err := runtime.RegisterExecutor(scraperworkflow.NewTypedExecutor(CombinedStepKind, func(ctx context.Context, step *scraperworkflow.StepContext, in batchInput) error {
		env, err := resolve(ctx, in.Identity)
		if err != nil {
			return err
		}
		result, err := ragoperators.ExecuteCombinedPreparationBatch(ctx, in.Plan, in.Batch, env)
		if err != nil {
			return err
		}
		return step.Result(batchOutput{Representations: result.Representations, CacheHit: result.CacheHit, ProviderCall: result.ProviderCall})
	})); err != nil {
		return err
	}
	if err := runtime.RegisterExecutor(scraperworkflow.NewTypedExecutor(FinalizeStepKind, func(_ context.Context, step *scraperworkflow.StepContext, in finalizerInput) error {
		all := make([]ragoperators.Representation, 0)
		for _, dep := range step.Step().DependsOn {
			var output batchOutput
			if err := step.DependencyData(dep.OpID, &output); err != nil {
				return err
			}
			all = append(all, output.Representations...)
		}
		if len(step.Step().DependsOn) != in.ExpectedBatches {
			return fmt.Errorf("RAG_PREPARATION_FINALIZE_DEPENDENCIES: got %d want %d", len(step.Step().DependsOn), in.ExpectedBatches)
		}
		return step.Result(map[string]any{"schemaVersion": "rag-preparation-finalize/v1", "representationCount": len(all)})
	})); err != nil {
		return err
	}
	return runtime.RegisterPackage(scraperworkflow.NewPackage(PackageName).DisplayName("RAG preparation").Entrypoint(scraperworkflow.EntrypointFunc[Input](build)).Build())
}

func build(_ context.Context, run *scraperworkflow.RunBuilder, input Input) error {
	if input.Identity.SchemaVersion != "rag-preparation-workflow/v1" || input.Identity.PreparedDigest == "" {
		return fmt.Errorf("RAG_PREPARATION_IDENTITY")
	}
	if len(input.Plan.Batches) == 0 {
		return fmt.Errorf("RAG_PREPARATION_EMPTY_PLAN")
	}
	steps := make([]scraperworkflow.StepHandle, 0, len(input.Plan.Batches))
	for _, batch := range input.Plan.Batches {
		handle, err := run.Step(fmt.Sprintf("combined-%04d", batch.Index), batchInput{Identity: input.Identity, Plan: input.Plan, Batch: batch}, scraperworkflow.StepOpts{Kind: CombinedStepKind, Queue: GenerationQueue, DedupKey: batchID(input.Identity, batch), Retry: model.RetryPolicy{MaxAttempts: 1}})
		if err != nil {
			return err
		}
		steps = append(steps, handle)
	}
	_, err := run.Step("finalize", finalizerInput{ExpectedBatches: len(steps)}, scraperworkflow.StepOpts{Kind: FinalizeStepKind, Queue: LocalQueue, DependsOn: scraperworkflow.Require(steps...)})
	return err
}

func batchID(identity Identity, batch ragoperators.CombinedPreparationBatch) string {
	return fmt.Sprintf("%s:combined:%04d", identity.PreparedDigest, batch.Index)
}
