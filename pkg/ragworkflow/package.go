package ragworkflow

import (
	_ "embed"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

//go:embed task.cjs
var taskSource []byte

type Package struct{}

func NewPackage() Package                 { return Package{} }
func (Package) Name() string              { return PackageName }
func (Package) Version() string           { return PackageVersion }
func (Package) RequiredModules() []string { return []string{ModuleAlias} }
func (Package) TaskModuleFactories() []workflowv3runtime.TaskModuleFactory {
	return []workflowv3runtime.TaskModuleFactory{TaskModuleFactory()}
}
func (Package) Bundle() (*workflowv3.Bundle, error) { return Bundle() }
func (Package) DescriptorModules() []workflowmodule.DescriptorModule {
	return []workflowmodule.DescriptorModule{DescriptorModule()}
}

func Bundle() (*workflowv3.Bundle, error) {
	stageInputs := map[string]string{"execution": ragcontractExecutionSchema(), "corpus": CorpusSchema, "prepared": PreparedSchema}
	stage := func(key workflowv3.TaskKey) workflowv3.BundleTask {
		return workflowv3.BundleTask{TaskKey: key, Entrypoint: "task.cjs#prepare", Inputs: cloneSchemas(stageInputs), Outputs: map[string]string{"prepared": PreparedSchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.prepare", Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}}
	}
	return workflowv3.NewBundle(workflowv3.BundleManifest{
		Name: PackageName, Version: PackageVersion, ABI: workflowv3.TaskABI,
		Tasks: []workflowv3.BundleTask{
			{TaskKey: TaskCorpusLoad, Entrypoint: "task.cjs#loadCorpus", Inputs: map[string]string{"execution": ragcontractExecutionSchema(), "corpus": CorpusSchema}, Outputs: map[string]string{"prepared": PreparedSchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.prepare", Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}},
			stage(TaskUnits), stage(TaskChunks), stage(TaskRepresent), stage(TaskEmbed), stage(TaskIndex),
			{TaskKey: TaskQuery, Entrypoint: "task.cjs#query", Inputs: map[string]string{"execution": ragcontractExecutionSchema(), "corpus": CorpusSchema, "prepared": PreparedSchema, "query": QuerySchema}, Outputs: map[string]string{"result": ResultPartitionSchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.query", Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}},
			{TaskKey: TaskMerge, Entrypoint: "task.cjs#merge", Inputs: map[string]string{"partition": workflowv3.ReductionPartitionSchemaV1}, Outputs: map[string]string{"result": ResultPartitionSchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.reduce", Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}},
			{TaskKey: TaskPublish, Entrypoint: "task.cjs#publish", Inputs: map[string]string{"execution": ragcontractExecutionSchema(), "results": ResultPartitionSchema}, Outputs: map[string]string{"result": ResultSchema}, Modules: []string{ModuleAlias}, ResourceClass: "cpu.rag.reduce", Retry: workflowv3.RetryPolicy{MaxAttempts: 2, BackoffMillis: 10}},
		},
	}, map[string][]byte{"task.cjs": taskSource})
}

func DescriptorModule() workflowmodule.DescriptorModule {
	return workflowmodule.DescriptorModule{Name: "rag-workflow-tasks", Factories: map[string]workflowv3.TaskKey{
		"loadCorpus": TaskCorpusLoad, "prepareUnits": TaskUnits, "createChunks": TaskChunks,
		"representRaw": TaskRepresent, "embedFixture": TaskEmbed, "buildIndex": TaskIndex,
		"evaluateQuery": TaskQuery, "mergeResults": TaskMerge, "publishResults": TaskPublish,
	}}
}

func cloneSchemas(input map[string]string) map[string]string {
	ret := make(map[string]string, len(input))
	for key, value := range input {
		ret[key] = value
	}
	return ret
}

func ragcontractExecutionSchema() string { return ragcontract.ExecutionSchemaVersion }
