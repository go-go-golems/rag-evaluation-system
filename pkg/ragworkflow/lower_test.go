package ragworkflow

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/stretchr/testify/require"
)

func fixtureExecution(t testing.TB, vector bool) ragcontract.PipelineExecution {
	t.Helper()
	fixture, err := NewProviderFreeFixture(vector)
	require.NoError(t, err)
	return fixture.Execution
}

func TestLoweringIsDeterministicAndUsesBoundedMapReduction(t *testing.T) {
	execution := fixtureExecution(t, true)
	first, err := NewLowerer().Lower(context.Background(), execution)
	require.NoError(t, err)
	second, err := NewLowerer().Lower(context.Background(), execution)
	require.NoError(t, err)
	require.Equal(t, first, second)
	require.Equal(t, LoweredExecutionSchema, first.SchemaVersion)
	require.NotEmpty(t, first.Plan.Digest)
	require.Len(t, first.IR.Maps, 1)
	require.Equal(t, 10_000, first.IR.Maps[0].Policy.MaxItems)
	require.Equal(t, 8, first.IR.Maps[0].Policy.MaxMaterializedAhead)
	require.Len(t, first.IR.Reductions, 1)
	require.Equal(t, 16, first.IR.Reductions[0].Policy.FanIn)
	require.Equal(t, ResultSchema, first.Plan.Outputs[0].Value.Schema)
	kinds := map[string]bool{}
	for _, node := range first.Plan.Nodes {
		kinds[node.Implementation.Kind] = true
	}
	for _, wanted := range []string{TaskCorpusLoad.Kind, TaskUnits.Kind, TaskChunks.Kind, TaskRepresent.Kind, TaskEmbed.Kind, TaskIndex.Kind, TaskPublish.Kind} {
		require.True(t, kinds[wanted], wanted)
	}
}

func TestQueryArchiveBindsEveryItemToDatasetManifest(t *testing.T) {
	fixture, err := NewProviderFreeFixture(false)
	require.NoError(t, err)
	archive, err := BuildQueryArchive(fixture.Execution, fixture.Dataset)
	require.NoError(t, err)
	require.Len(t, archive.Items, 2)
	for _, archived := range archive.Items {
		var item QueryItem
		require.NoError(t, json.Unmarshal(archived.Data, &item))
		require.Equal(t, fixture.Execution.Dataset.ManifestDigest, item.DatasetManifestDigest)
	}
	fixture.Dataset.Queries[0].Text += " stale"
	_, err = BuildQueryArchive(fixture.Execution, fixture.Dataset)
	require.ErrorContains(t, err, "RAG_WORKFLOW_DATASET_DIGEST")
}

func TestOperatorRegistryIsClosedVersionedAndDeclaresProviderAttachments(t *testing.T) {
	registry := NewOperatorRegistry()
	definitions := registry.Definitions()
	require.NotEmpty(t, definitions)
	for index := 1; index < len(definitions); index++ {
		require.Less(t, definitions[index-1].Operator.ID(), definitions[index].Operator.ID())
	}
	provider, found := registry.Definition(ragcontract.OperatorRef{Kind: "generate.answer", Version: "v1"})
	require.True(t, found)
	require.True(t, provider.ProviderRequired)
	_, found = registry.Definition(ragcontract.OperatorRef{Kind: "generate.answer", Version: "v2"})
	require.False(t, found)
	_, found = registry.Definition(ragcontract.OperatorRef{Kind: "retrieve.future", Version: "v1"})
	require.False(t, found)
}

func TestLoweringRejectsUnknownProviderAndStaleIdentity(t *testing.T) {
	execution := fixtureExecution(t, false)
	execution.Pipeline.Nodes[0].Operator.Kind = "unknown.operator"
	_, err := NewLowerer().Lower(context.Background(), execution)
	require.ErrorContains(t, err, "RAG_V2_OPERATOR_UNKNOWN")
	execution = fixtureExecution(t, false)
	for index := range execution.Pipeline.Nodes {
		if execution.Pipeline.Nodes[index].Operator.Kind == "representations.raw" {
			execution.Pipeline.Nodes[index].Operator.Kind = "representations.structured-summary"
			execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"name":"summary","model":"m","prompt":"p","outputSchema":"s"}`)
		}
	}
	execution.CellID = ""
	execution.CellID, _ = ragcontract.Digest(execution)
	_, err = NewLowerer().Lower(context.Background(), execution)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "PROVIDER_REQUIRED") || strings.Contains(err.Error(), "PIPELINE"), err.Error())
	execution = fixtureExecution(t, false)
	execution.CellID = "sha256:" + strings.Repeat("0", 64)
	_, err = NewLowerer().Lower(context.Background(), execution)
	require.ErrorContains(t, err, "RAG_WORKFLOW_CELL_ID")
}
