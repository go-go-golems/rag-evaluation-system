package ragintakeworkflow

import (
	_ "embed"

	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

const (
	PackageName    = "rag-intake-v1"
	PackageVersion = "1.0.0"
	ModuleAlias    = "rag:intake"
)

var (
	TaskPreprocess = workflowv3.TaskKey{Kind: "rag.intake.preprocess-document", Version: "v1"}
	TaskChunk      = workflowv3.TaskKey{Kind: "rag.intake.chunk-document", Version: "v1"}
	TaskEnrich     = workflowv3.TaskKey{Kind: "rag.intake.enrich-chunk", Version: "v1"}
	TaskEmbed      = workflowv3.TaskKey{Kind: "rag.intake.compute-embeddings", Version: "v1"}
	TaskBM25       = workflowv3.TaskKey{Kind: "rag.intake.build-bm25", Version: "v1"}
	TaskPublish    = workflowv3.TaskKey{Kind: "rag.intake.publish", Version: "v1"}
)

//go:embed task.cjs
var taskSource []byte

type Package struct{ Config RuntimeConfig }

func NewPackage(config RuntimeConfig) Package { return Package{Config: config} }
func (Package) Name() string                  { return PackageName }
func (Package) Version() string               { return PackageVersion }
func (Package) RequiredModules() []string     { return []string{ModuleAlias} }
func (p Package) TaskModuleFactories() []workflowv3runtime.TaskModuleFactory {
	return []workflowv3runtime.TaskModuleFactory{TaskModuleFactory(p.Config)}
}
func (Package) DescriptorModules() []workflowmodule.DescriptorModule {
	return []workflowmodule.DescriptorModule{DescriptorModule()}
}
func (Package) Bundle() (*workflowv3.Bundle, error) { return Bundle() }

func Bundle() (*workflowv3.Bundle, error) {
	task := func(key workflowv3.TaskKey, resource string, retry workflowv3.RetryPolicy) workflowv3.BundleTask {
		return workflowv3.BundleTask{TaskKey: key, Entrypoint: "task.cjs#run", Inputs: map[string]string{"request": RequestSchema}, Outputs: map[string]string{"result": ResultSchema}, Modules: []string{ModuleAlias}, ResourceClass: resource, Retry: retry}
	}
	embed := task(TaskEmbed, "cpu.rag.intake.embedding", workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 100})
	embed.BudgetMaximum = &workflowv3.BudgetClaim{Account: "provider", Reserve: []workflowv3.BudgetAmount{{Dimension: "requests", Units: 10_000}}, OnExhausted: workflowv3.BudgetExhaustFailRun}
	return workflowv3.NewBundle(workflowv3.BundleManifest{Name: PackageName, Version: PackageVersion, ABI: workflowv3.TaskABI, Tasks: []workflowv3.BundleTask{
		task(TaskPreprocess, "cpu.rag.intake.llm", workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 100}),
		task(TaskChunk, "cpu.rag.intake", workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}),
		task(TaskEnrich, "cpu.rag.intake.llm", workflowv3.RetryPolicy{MaxAttempts: 3, BackoffMillis: 100}),
		embed,
		task(TaskBM25, "cpu.rag.intake.index", workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}),
		{TaskKey: TaskPublish, Entrypoint: "task.cjs#run", Inputs: map[string]string{"request": RequestSchema}, Outputs: map[string]string{"result": SummarySchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.intake", Retry: workflowv3.RetryPolicy{MaxAttempts: 1}},
	}}, map[string][]byte{"task.cjs": taskSource})
}

func DescriptorModule() workflowmodule.DescriptorModule {
	return workflowmodule.DescriptorModule{Name: "rag-intake-tasks", Factories: map[string]workflowv3.TaskKey{"preprocessDocument": TaskPreprocess, "chunkDocument": TaskChunk, "enrichChunk": TaskEnrich, "computeEmbeddings": TaskEmbed, "buildBM25": TaskBM25, "publish": TaskPublish}}
}
