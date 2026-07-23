package ragworkflow

import (
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3product"
)

func BuildRunnerExecution(lowered LoweredExecution, bindings map[string]researchrunner.InputBinding) (researchrunner.WorkflowExecution, error) {
	if lowered.SchemaVersion != LoweredExecutionSchema || lowered.Plan.Digest == "" {
		return researchrunner.WorkflowExecution{}, fmt.Errorf("RAG_WORKFLOW_LOWERED_IDENTITY")
	}
	packages, err := workflowv3product.BuildPackageSet([]string{PackageName}, NewPackage())
	if err != nil {
		return researchrunner.WorkflowExecution{}, err
	}
	return researchrunner.BuildExecution(lowered.Plan, packages, bindings, researchrunner.ObservationPolicy{ExportOutputs: true, ExportExternalOperations: true, ExportCanonicalObservations: true})
}

func BuildProviderRunnerExecution(lowered LoweredExecution, bindings map[string]researchrunner.InputBinding, providerPackage *ProviderPackage) (researchrunner.WorkflowExecution, error) {
	if lowered.SchemaVersion != LoweredExecutionSchema || lowered.Plan.Digest == "" || providerPackage == nil {
		return researchrunner.WorkflowExecution{}, fmt.Errorf("RAG_WORKFLOW_PROVIDER_LOWERED_IDENTITY")
	}
	packages, err := workflowv3product.BuildPackageSet([]string{ProviderPackageName}, providerPackage)
	if err != nil {
		return researchrunner.WorkflowExecution{}, err
	}
	return researchrunner.BuildExecution(lowered.Plan, packages, bindings, researchrunner.ObservationPolicy{ExportOutputs: true, ExportExternalOperations: true, ExportCanonicalObservations: true})
}

func BuildQueryArchive(execution ragcontract.PipelineExecution, dataset ragoperators.EvaluationDataset) (researchrunner.SetInputArchive, error) {
	if dataset.SchemaVersion != "rag-evaluation-data/v1" || len(dataset.Queries) == 0 || len(dataset.Queries) > 10_000 {
		return researchrunner.SetInputArchive{}, fmt.Errorf("RAG_WORKFLOW_DATASET")
	}
	datasetDigest, err := ragcontract.Digest(dataset)
	if err != nil {
		return researchrunner.SetInputArchive{}, err
	}
	if datasetDigest != execution.Dataset.ManifestDigest {
		return researchrunner.SetInputArchive{}, fmt.Errorf("RAG_WORKFLOW_DATASET_DIGEST: got %s want %s", datasetDigest, execution.Dataset.ManifestDigest)
	}
	archive := researchrunner.SetInputArchive{SchemaVersion: researchrunner.SetInputArchiveSchema, ItemSchema: QuerySchema, ManifestSchema: workflowv3.ItemManifestSchemaV1}
	previous := ""
	for _, query := range dataset.Queries {
		if query.ID == "" || query.Text == "" || query.ID <= previous {
			return researchrunner.SetInputArchive{}, fmt.Errorf("RAG_WORKFLOW_QUERY_ORDER")
		}
		body, err := ragcontract.CanonicalJSON(QueryItem{SchemaVersion: QuerySchema, DatasetManifestDigest: datasetDigest, Query: query})
		if err != nil {
			return researchrunner.SetInputArchive{}, err
		}
		archive.Items = append(archive.Items, researchrunner.SetInputArchiveItem{Key: query.ID, MediaType: "application/json", Data: body})
		previous = query.ID
	}
	return archive, nil
}
