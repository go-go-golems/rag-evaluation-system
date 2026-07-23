// Package ragworkflow lowers canonical RAG v2 executions into domain-neutral
// Scraper Workflow V3 plans and supplies the versioned RAG task package.
package ragworkflow

import (
	"encoding/json"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const (
	PackageName    = "rag-v2-provider-free"
	PackageVersion = "1.0.0"
	ModuleAlias    = "rag:workflow"

	CorpusSchema          = "rag-workflow-corpus/v1"
	QuerySchema           = "rag-workflow-query/v1"
	PreparedSchema        = "rag-workflow-prepared/v1"
	ResultPartitionSchema = "rag-workflow-result-partition/v1"
	ResultSchema          = "rag-workflow-result/v1"
)

var (
	TaskCorpusLoad = workflowv3.TaskKey{Kind: "rag.corpus.load", Version: "v1"}
	TaskUnits      = workflowv3.TaskKey{Kind: "rag.units.prepare", Version: "v1"}
	TaskChunks     = workflowv3.TaskKey{Kind: "rag.chunks.create", Version: "v1"}
	TaskRepresent  = workflowv3.TaskKey{Kind: "rag.represent.raw", Version: "v1"}
	TaskEmbed      = workflowv3.TaskKey{Kind: "rag.embed.fixture", Version: "v1"}
	TaskIndex      = workflowv3.TaskKey{Kind: "rag.index.build", Version: "v1"}
	TaskQuery      = workflowv3.TaskKey{Kind: "rag.query.evaluate", Version: "v1"}
	TaskMerge      = workflowv3.TaskKey{Kind: "rag.results.merge", Version: "v1"}
	TaskPublish    = workflowv3.TaskKey{Kind: "rag.results.publish", Version: "v1"}
)

type LoweredExecution struct {
	SchemaVersion          string                  `json:"schemaVersion"`
	ExecutionDigest        string                  `json:"executionDigest"`
	PreparationFingerprint string                  `json:"preparationFingerprint"`
	IR                     workflowv3.WorkflowIR   `json:"ir"`
	Plan                   workflowv3.WorkflowPlan `json:"plan"`
}

type IndexEvidence struct {
	NodeID        string                    `json:"nodeId"`
	Manifest      ragcontract.IndexManifest `json:"manifest"`
	RecordsDigest string                    `json:"recordsDigest"`
	RecordsSize   int64                     `json:"recordsSize"`
}

type PreparedBundle struct {
	SchemaVersion          string                    `json:"schemaVersion"`
	ExecutionDigest        string                    `json:"executionDigest"`
	PipelineDigest         string                    `json:"pipelineDigest"`
	CorpusDigest           string                    `json:"corpusDigest"`
	PreparationFingerprint string                    `json:"preparationFingerprint"`
	Values                 []ragengine.PreparedValue `json:"values"`
	Indexes                []IndexEvidence           `json:"indexes,omitempty"`
	Digest                 string                    `json:"digest"`
}

type Artifact struct {
	Role          string          `json:"role"`
	Kind          string          `json:"kind"`
	Name          string          `json:"name"`
	SchemaVersion string          `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Data          []byte          `json:"data"`
	Metadata      json.RawMessage `json:"metadata,omitempty"`
}

type QueryItem struct {
	SchemaVersion         string             `json:"schemaVersion"`
	DatasetManifestDigest string             `json:"datasetManifestDigest"`
	Query                 ragoperators.Query `json:"query"`
}

type Metric struct {
	Name     string          `json:"name"`
	Unit     string          `json:"unit,omitempty"`
	Value    json.RawMessage `json:"value"`
	Numeric  *float64        `json:"numeric,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

type QueryResult struct {
	QueryID   string                     `json:"queryId"`
	Trace     ragcontract.QueryTrace     `json:"trace"`
	Metrics   []Metric                   `json:"metrics"`
	Artifacts []Artifact                 `json:"artifacts"`
	Answers   []ragoperators.Answer      `json:"answers,omitempty"`
	Failures  []ragcontract.FailureTrace `json:"failures,omitempty"`
}

type ResultPartition struct {
	SchemaVersion          string        `json:"schemaVersion"`
	ExecutionDigest        string        `json:"executionDigest"`
	PreparationFingerprint string        `json:"preparationFingerprint"`
	Results                []QueryResult `json:"results"`
	Digest                 string        `json:"digest"`
}

type Result struct {
	SchemaVersion          string                        `json:"schemaVersion"`
	ExecutionDigest        string                        `json:"executionDigest"`
	CellID                 string                        `json:"cellId"`
	VariantID              string                        `json:"variantId"`
	Factors                []ragcontract.FactorSelection `json:"factors"`
	PreparationFingerprint string                        `json:"preparationFingerprint"`
	Results                []QueryResult                 `json:"results"`
	Digest                 string                        `json:"digest"`
}
