// Package preparationworkflow adapts immutable RAG combined-preparation plans
// to scraper's durable workflow runtime. It deliberately leaves prompt/model
// resolution and validation in ragoperators.
package preparationworkflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	scraperworkflow "github.com/go-go-golems/scraper/pkg/workflow"
)

const (
	PackageName       = "rag-preparation/v1"
	CombinedStepKind  = "rag-preparation/combined-batch/v1"
	EmbeddingStepKind = "rag-preparation/embedding-batch/v1"
	FinalizeStepKind  = "rag-preparation/finalize/v1"
	GenerationQueue   = model.QueueKey("rag:generator")
	EmbeddingQueue    = model.QueueKey("rag:embedding")
	LocalQueue        = model.QueueKey("rag:local")
)

type Identity struct {
	SchemaVersion  string `json:"schemaVersion"`
	PreparedDigest string `json:"preparedDigest"`
}

type EmbeddingSpec struct {
	Node                  ragcontract.Node `json:"node"`
	RawRepresentationName string           `json:"rawRepresentationName"`
}

type Input struct {
	Identity  Identity                             `json:"identity"`
	Plan      ragoperators.CombinedPreparationPlan `json:"plan"`
	Embedding *EmbeddingSpec                       `json:"embedding,omitempty"`
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

type embeddingInput struct {
	Identity       Identity                              `json:"identity"`
	Plan           ragoperators.CombinedPreparationPlan  `json:"plan"`
	Batch          ragoperators.CombinedPreparationBatch `json:"batch"`
	CombinedStepID model.OpID                            `json:"combinedStepId"`
	Spec           EmbeddingSpec                         `json:"spec"`
}

type embeddingOutput struct {
	Representations []ragoperators.Representation `json:"representations"`
	Embeddings      []ragoperators.Embedding      `json:"embeddings"`
	ProviderCall    bool                          `json:"providerCall"`
}

type finalizerInput struct {
	ExpectedBatches int  `json:"expectedBatches"`
	Embeddings      bool `json:"embeddings"`
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
		output := batchOutput{Representations: result.Representations, CacheHit: result.CacheHit, ProviderCall: result.ProviderCall}
		body, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("marshal combined batch artifact: %w", err)
		}
		if _, err := step.Artifact("combined-batch-result.json", "application/json", body, scraperworkflow.ArtifactKind("rag-preparation-batch-result")); err != nil {
			return err
		}
		return step.Result(output)
	})); err != nil {
		return err
	}
	if err := runtime.RegisterExecutor(scraperworkflow.NewTypedExecutor(EmbeddingStepKind, func(ctx context.Context, step *scraperworkflow.StepContext, in embeddingInput) error {
		var combined batchOutput
		if err := step.DependencyData(in.CombinedStepID, &combined); err != nil {
			return err
		}
		raw, err := ragoperators.RawRepresentations(in.Spec.RawRepresentationName, in.Batch.Chunks)
		if err != nil {
			return err
		}
		representations := append(raw, combined.Representations...)
		sort.Slice(representations, func(i, j int) bool { return representations[i].Record.ID < representations[j].Record.ID })
		for i := 1; i < len(representations); i++ {
			if representations[i].Record.ID == representations[i-1].Record.ID {
				return fmt.Errorf("RAG_PREPARATION_REPRESENTATION_DUPLICATE: %s", representations[i].Record.ID)
			}
		}
		plan, err := ragoperators.PlanEmbeddingBatches(representations, in.Spec.Node)
		if err != nil {
			return err
		}
		if len(plan.Batches) != 1 {
			return fmt.Errorf("RAG_PREPARATION_EMBED_BATCH_CONFIG: combined batch %d requires %d embedding requests", in.Batch.Index, len(plan.Batches))
		}
		env, err := resolve(ctx, in.Identity)
		if err != nil {
			return err
		}
		result, err := ragoperators.ExecuteEmbeddingBatch(ctx, plan, plan.Batches[0], env)
		if err != nil {
			return err
		}
		output := embeddingOutput{Representations: representations, Embeddings: result.Embeddings, ProviderCall: result.ProviderCall}
		body, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("marshal embedding batch artifact: %w", err)
		}
		if _, err := step.Artifact("embedding-batch-result.json", "application/json", body, scraperworkflow.ArtifactKind("rag-preparation-embedding-batch-result")); err != nil {
			return err
		}
		return step.Result(output)
	})); err != nil {
		return err
	}
	if err := runtime.RegisterExecutor(scraperworkflow.NewTypedExecutor(FinalizeStepKind, func(_ context.Context, step *scraperworkflow.StepContext, in finalizerInput) error {
		if len(step.Step().DependsOn) != in.ExpectedBatches {
			return fmt.Errorf("RAG_PREPARATION_FINALIZE_DEPENDENCIES: got %d want %d", len(step.Step().DependsOn), in.ExpectedBatches)
		}
		representationCount, embeddingCount := 0, 0
		for _, dep := range step.Step().DependsOn {
			if in.Embeddings {
				var output embeddingOutput
				if err := step.DependencyData(dep.OpID, &output); err != nil {
					return err
				}
				representationCount += len(output.Representations)
				embeddingCount += len(output.Embeddings)
				continue
			}
			var output batchOutput
			if err := step.DependencyData(dep.OpID, &output); err != nil {
				return err
			}
			representationCount += len(output.Representations)
		}
		return step.Result(map[string]any{"schemaVersion": "rag-preparation-finalize/v1", "representationCount": representationCount, "embeddingCount": embeddingCount})
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
	combinedSteps := make([]scraperworkflow.StepHandle, 0, len(input.Plan.Batches))
	for _, batch := range input.Plan.Batches {
		handle, err := run.Step(fmt.Sprintf("combined-%04d", batch.Index), batchInput{Identity: input.Identity, Plan: input.Plan, Batch: batch}, scraperworkflow.StepOpts{Kind: CombinedStepKind, Queue: GenerationQueue, DedupKey: batchID(input.Identity, batch), Retry: model.RetryPolicy{MaxAttempts: 1}})
		if err != nil {
			return err
		}
		combinedSteps = append(combinedSteps, handle)
	}
	finalDependencies := combinedSteps
	if input.Embedding != nil {
		if input.Embedding.RawRepresentationName == "" {
			return fmt.Errorf("RAG_PREPARATION_RAW_REPRESENTATION_NAME")
		}
		embeddingSteps := make([]scraperworkflow.StepHandle, 0, len(input.Plan.Batches))
		for index, batch := range input.Plan.Batches {
			handle, err := run.Step(fmt.Sprintf("embedding-%04d", batch.Index), embeddingInput{Identity: input.Identity, Plan: input.Plan, Batch: batch, CombinedStepID: combinedSteps[index].ID, Spec: *input.Embedding}, scraperworkflow.StepOpts{Kind: EmbeddingStepKind, Queue: EmbeddingQueue, DedupKey: fmt.Sprintf("%s:embedding:%04d", input.Identity.PreparedDigest, batch.Index), DependsOn: scraperworkflow.Require(combinedSteps[index]), Retry: model.RetryPolicy{MaxAttempts: 1}})
			if err != nil {
				return err
			}
			embeddingSteps = append(embeddingSteps, handle)
		}
		finalDependencies = embeddingSteps
	}
	_, err := run.Step("finalize", finalizerInput{ExpectedBatches: len(finalDependencies), Embeddings: input.Embedding != nil}, scraperworkflow.StepOpts{Kind: FinalizeStepKind, Queue: LocalQueue, DependsOn: scraperworkflow.Require(finalDependencies...)})
	return err
}

func batchID(identity Identity, batch ragoperators.CombinedPreparationBatch) string {
	return fmt.Sprintf("%s:combined:%04d", identity.PreparedDigest, batch.Index)
}
