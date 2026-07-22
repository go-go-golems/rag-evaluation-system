package workflowv3ttc

import (
	_ "embed"

	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

//go:embed tasks.cjs
var taskSource []byte

//go:embed workflow.js
var workflowSource string

//go:embed production_workflow.js
var productionWorkflowSource string

//go:embed sweep_workflow.js
var sweepWorkflowSource string

var (
	GenerateKey            = workflowv3.TaskKey{Kind: "rag.ttc.generate-representations", Version: "v1"}
	GenerateBatchKey       = workflowv3.TaskKey{Kind: "rag.ttc.generate-representation-batch", Version: "v1"}
	EmbedBatchKey          = workflowv3.TaskKey{Kind: "rag.ttc.measure-embedding-batch", Version: "v1"}
	EmbedKey               = workflowv3.TaskKey{Kind: "rag.ttc.embed-representations", Version: "v1"}
	MergeKey               = workflowv3.TaskKey{Kind: "rag.ttc.merge-prepared-shard", Version: "v1"}
	ValidatePublicationKey = workflowv3.TaskKey{Kind: "rag.ttc.validate-publication", Version: "v1"}
	PublishKey             = workflowv3.TaskKey{Kind: "rag.ttc.publish-prepared", Version: "v1"}
	EvaluateKey            = workflowv3.TaskKey{Kind: "rag.ttc.evaluate-query", Version: "v1"}
	MergeEvidenceKey       = workflowv3.TaskKey{Kind: "rag.ttc.merge-study-evidence", Version: "v1"}
)

func Bundle() (*workflowv3.Bundle, error) {
	return workflowv3.NewBundle(workflowv3.BundleManifest{
		Name: "rag-ttc-v3", Version: "1.0.0", ABI: workflowv3.TaskABI,
		Tasks: []workflowv3.BundleTask{
			{TaskKey: GenerateBatchKey, Entrypoint: "tasks.cjs#generateBatch", Inputs: map[string]string{"batch": ChunkBatchSchema}, Outputs: map[string]string{"generated": GeneratedBatchSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceGeneration, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "generation", Reserve: []workflowv3.BudgetAmount{{Dimension: "cost_microunits", Units: 160_000}, {Dimension: "input_tokens", Units: 16_384}, {Dimension: "output_tokens", Units: 16_384}, {Dimension: "requests", Units: 1}}, OnExhausted: "fail-run"}},
			{TaskKey: EmbedBatchKey, Entrypoint: "tasks.cjs#embedBatch", Inputs: map[string]string{"generated": GeneratedBatchSchema}, Outputs: map[string]string{"measured": MeasuredBatchSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceEmbedding, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "embedding", Reserve: []workflowv3.BudgetAmount{{Dimension: "embedding_tokens", Units: 65536}, {Dimension: "requests", Units: 8}}, OnExhausted: "fail-run"}},
			{TaskKey: GenerateKey, Entrypoint: "tasks.cjs#generate", Inputs: map[string]string{"chunk": ChunkSchema}, Outputs: map[string]string{"generated": GeneratedSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceGeneration, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "generation", Reserve: []workflowv3.BudgetAmount{{Dimension: "cost_microunits", Units: 20_000}, {Dimension: "input_tokens", Units: 2_048}, {Dimension: "output_tokens", Units: 2_048}, {Dimension: "requests", Units: 1}}, OnExhausted: "fail-run"}},
			{TaskKey: EmbedKey, Entrypoint: "tasks.cjs#embed", Inputs: map[string]string{"generated": GeneratedSchema}, Outputs: map[string]string{"embedded": ShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceEmbedding, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "embedding", Reserve: []workflowv3.BudgetAmount{{Dimension: "embedding_tokens", Units: 4_096}}, OnExhausted: "fail-run"}},
			{TaskKey: MergeKey, Entrypoint: "tasks.cjs#merge", Inputs: map[string]string{"partition": workflowv3.ReductionPartitionSchemaV1}, Outputs: map[string]string{"shard": ShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceLocal},
			{TaskKey: ValidatePublicationKey, Entrypoint: "tasks.cjs#validatePublication", Inputs: map[string]string{"shard": ShardSchema}, Outputs: map[string]string{"receipt": ValidationReceiptSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceLocal},
			{TaskKey: PublishKey, Entrypoint: "tasks.cjs#publish", Inputs: map[string]string{"shard": ShardSchema, "decision": PublicationDecisionSchema}, Outputs: map[string]string{"publication": PublicationReceiptSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceLocal},
			{TaskKey: EvaluateKey, Entrypoint: "tasks.cjs#evaluate", Inputs: map[string]string{"publication": PublicationReceiptSchema, "query": QuerySchema}, Outputs: map[string]string{"evidence": StudyEvidenceShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceEvaluation, Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 100}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "evaluation", Reserve: []workflowv3.BudgetAmount{{Dimension: "cost_microunits", Units: 50_000}, {Dimension: "embedding_tokens", Units: 1_000}, {Dimension: "input_tokens", Units: 10_000}, {Dimension: "output_tokens", Units: 2_000}, {Dimension: "requests", Units: 3}}, OnExhausted: "fail-run"}},
			{TaskKey: MergeEvidenceKey, Entrypoint: "tasks.cjs#mergeEvidence", Inputs: map[string]string{"partition": workflowv3.ReductionPartitionSchemaV1}, Outputs: map[string]string{"evidence": StudyEvidenceShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceLocal},
		},
	}, map[string][]byte{"tasks.cjs": taskSource})
}

func Registry() (*workflowv3.SealedRegistry, error) {
	bundle, err := Bundle()
	if err != nil {
		return nil, err
	}
	builder := workflowv3.NewRegistryBuilder()
	if err := builder.AdvertiseModules(ModuleAlias); err != nil {
		return nil, err
	}
	if err := builder.AddBundle(bundle); err != nil {
		return nil, err
	}
	return builder.Seal()
}

func DescriptorModule() workflowmodule.DescriptorModule {
	return workflowmodule.DescriptorModule{Name: "rag-ttc-v3-tasks", Factories: map[string]workflowv3.TaskKey{"generateBatch": GenerateBatchKey, "embedBatch": EmbedBatchKey, "generate": GenerateKey, "embed": EmbedKey, "merge": MergeKey, "validatePublication": ValidatePublicationKey, "publish": PublishKey, "evaluate": EvaluateKey, "mergeEvidence": MergeEvidenceKey}}
}

func WorkflowSource() string           { return workflowSource }
func ProductionWorkflowSource() string { return productionWorkflowSource }
func SweepWorkflowSource() string      { return sweepWorkflowSource }
