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

var (
	GenerateKey = workflowv3.TaskKey{Kind: "rag.ttc.generate-representations", Version: "v1"}
	EmbedKey    = workflowv3.TaskKey{Kind: "rag.ttc.embed-representations", Version: "v1"}
	MergeKey    = workflowv3.TaskKey{Kind: "rag.ttc.merge-prepared-shard", Version: "v1"}
)

func Bundle() (*workflowv3.Bundle, error) {
	return workflowv3.NewBundle(workflowv3.BundleManifest{
		Name: "rag-ttc-v3", Version: "1.0.0", ABI: workflowv3.TaskABI,
		Tasks: []workflowv3.BundleTask{
			{TaskKey: GenerateKey, Entrypoint: "tasks.cjs#generate", Inputs: map[string]string{"chunk": ChunkSchema}, Outputs: map[string]string{"generated": GeneratedSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceGeneration, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "generation", Reserve: []workflowv3.BudgetAmount{{Dimension: "input_tokens", Units: 4}, {Dimension: "output_tokens", Units: 2}, {Dimension: "requests", Units: 1}}, OnExhausted: "fail-run"}},
			{TaskKey: EmbedKey, Entrypoint: "tasks.cjs#embed", Inputs: map[string]string{"generated": GeneratedSchema}, Outputs: map[string]string{"embedded": ShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceEmbedding, Retry: workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 25}, BudgetMaximum: &workflowv3.BudgetClaim{Account: "embedding", Reserve: []workflowv3.BudgetAmount{{Dimension: "embedding_tokens", Units: 3}}, OnExhausted: "fail-run"}},
			{TaskKey: MergeKey, Entrypoint: "tasks.cjs#merge", Inputs: map[string]string{"partition": workflowv3.ReductionPartitionSchemaV1}, Outputs: map[string]string{"shard": ShardSchema}, Modules: []string{ModuleAlias}, ResourceClass: ResourceLocal},
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
	return workflowmodule.DescriptorModule{Name: "rag-ttc-v3-tasks", Factories: map[string]workflowv3.TaskKey{"generate": GenerateKey, "embed": EmbedKey, "merge": MergeKey}}
}

func WorkflowSource() string { return workflowSource }
