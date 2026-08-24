package ragworkflow

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/stretchr/testify/require"
)

func TestTaskFailureIsClosedPermanentAndRedacted(t *testing.T) {
	ctx := context.Background()
	fixture, err := NewProviderFreeFixture(false)
	require.NoError(t, err)
	lowered, err := NewLowerer().Lower(ctx, fixture.Execution)
	require.NoError(t, err)
	artifacts, err := workflowv3.NewFileArtifactStore(t.TempDir(), 1<<20)
	require.NoError(t, err)
	executionBody, err := ragcontract.CanonicalJSON(fixture.Execution)
	require.NoError(t, err)
	executionRef, err := artifacts.Put(ctx, ragcontract.ExecutionSchemaVersion, "application/json", executionBody)
	require.NoError(t, err)
	secret := "SECRET-CORPUS-CANARY"
	badCorpus := []byte(`{"schemaVersion":"bad","records":[{"text":"` + secret + `"}]}`)
	corpusRef, err := artifacts.Put(ctx, CorpusSchema, "application/json", badCorpus)
	require.NoError(t, err)
	bundle, err := Bundle()
	require.NoError(t, err)
	builder := workflowv3.NewRegistryBuilder()
	require.NoError(t, builder.AdvertiseModules(ModuleAlias))
	require.NoError(t, builder.AddBundle(bundle))
	registry, err := builder.Seal()
	require.NoError(t, err)
	var start workflowv3.PlanNode
	for _, node := range lowered.Plan.Nodes {
		if node.Key == "prepare-start" {
			start = node
			break
		}
	}
	task, err := registry.ResolveNode(start)
	require.NoError(t, err)
	modules, err := workflowv3runtime.NewTaskModuleRegistry(TaskModuleFactory())
	require.NoError(t, err)
	_, err = workflowv3runtime.RunTask(ctx, workflowv3runtime.TaskRequest{RunID: "failure", NodeKey: start.Key, Attempt: 1, Task: task, Inputs: map[string]workflowv3.ArtifactRef{"execution": executionRef, "corpus": corpusRef}, Artifacts: artifacts, Modules: modules})
	require.Error(t, err)
	var failure *workflowv3runtime.TaskFailureError
	require.True(t, errors.As(err, &failure))
	require.Equal(t, "RAG_WORKFLOW_TASK_FAILED", failure.Failure.Code)
	require.False(t, failure.Failure.Retryable)
	require.NotContains(t, strings.ToUpper(err.Error()), secret)
}
