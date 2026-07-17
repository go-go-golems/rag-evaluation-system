// Package ragcontract defines the dependency-free canonical RAG v2 wire contracts.
// It contains no authoring runtime, provider, storage, filesystem, or researchctl code.
package ragcontract

import "encoding/json"

const (
	PipelineSchemaVersion  = "rag-pipeline-ir/v2"
	ProductSchemaVersion   = "rag-product-plan/v2"
	StudySchemaVersion     = "rag-study/v2"
	ExecutionSchemaVersion = "rag-pipeline-execution/v2"
	Domain                 = "rag-pipeline"
	DomainSchemaVersion    = "rag-pipeline/v2"
)

type PortKind string

const (
	PortCorpus          PortKind = "corpus"
	PortUnits           PortKind = "units"
	PortChunks          PortKind = "chunks"
	PortRepresentations PortKind = "representations"
	PortEmbeddings      PortKind = "embeddings"
	PortIndex           PortKind = "index"
	PortQuery           PortKind = "query"
	PortRankedRecords   PortKind = "ranked-records"
	PortRankedParents   PortKind = "ranked-parents"
	PortEvidence        PortKind = "evidence"
	PortAnswer          PortKind = "answer"
	PortEvaluation      PortKind = "evaluation"
)

type OperatorRef struct {
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

func (r OperatorRef) ID() string { return r.Kind + "/" + r.Version }

type InputSlot struct {
	ID             string   `json:"id"`
	Kind           PortKind `json:"kind"`
	BindingMode    string   `json:"bindingMode"`
	ArtifactRole   string   `json:"artifactRole,omitempty"`
	ManifestSchema string   `json:"manifestSchema,omitempty"`
	Digest         string   `json:"digest,omitempty"`
}

type PortRef struct {
	NodeID string `json:"nodeId"`
	Port   string `json:"port"`
}

type InputBinding struct {
	Port string  `json:"port"`
	From PortRef `json:"from"`
}

type Node struct {
	ID       string          `json:"id"`
	Operator OperatorRef     `json:"operator"`
	Inputs   []InputBinding  `json:"inputs"`
	Config   json.RawMessage `json:"config"`
	Order    int             `json:"order,omitempty"`
}

type OutputRef struct {
	Name string   `json:"name"`
	Kind PortKind `json:"kind"`
	From PortRef  `json:"from"`
}

type PipelineIR struct {
	SchemaVersion string      `json:"schemaVersion"`
	Inputs        []InputSlot `json:"inputs"`
	Nodes         []Node      `json:"nodes"`
	Outputs       []OutputRef `json:"outputs"`
	SeedPolicy    *SeedPolicy `json:"seedPolicy,omitempty"`
}

type SeedPolicy struct {
	Mode string `json:"mode"`
	Seed *int64 `json:"seed,omitempty"`
}

type DisplayMetadata struct {
	Name  string            `json:"name,omitempty"`
	Notes []string          `json:"notes,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
}

type ArtifactBinding struct {
	SlotID        string `json:"slotId"`
	Role          string `json:"role"`
	Kind          string `json:"kind"`
	ID            string `json:"id,omitempty"`
	URI           string `json:"uri,omitempty"`
	Digest        string `json:"digest"`
	SizeBytes     *int64 `json:"sizeBytes,omitempty"`
	SchemaVersion string `json:"schemaVersion"`
}

type ModelBinding struct {
	Reference string `json:"reference"`
	Manifest  string `json:"manifest"`
	Digest    string `json:"digest"`
}

type DatasetBinding struct {
	ManifestDigest  string `json:"manifestDigest"`
	Split           string `json:"split"`
	Status          string `json:"status"`
	RelevanceTarget string `json:"relevanceTarget"`
}

type Measure struct {
	Name      string          `json:"name"`
	Version   string          `json:"version"`
	ValueKind string          `json:"valueKind"`
	Unit      string          `json:"unit,omitempty"`
	Required  bool            `json:"required"`
	Config    json.RawMessage `json:"config"`
}

type CitationPolicy struct {
	Mode              string `json:"mode"`
	RequireSourceText bool   `json:"requireSourceText"`
}

type RuntimePolicy struct {
	TimeoutMilliseconds int64  `json:"timeoutMilliseconds,omitempty"`
	MaxResults          int    `json:"maxResults"`
	TracePolicy         string `json:"tracePolicy"`
	FailurePolicy       string `json:"failurePolicy"`
}

type ProductPlan struct {
	SchemaVersion string            `json:"schemaVersion"`
	Pipeline      PipelineIR        `json:"pipeline"`
	Bindings      []ArtifactBinding `json:"bindings"`
	Models        []ModelBinding    `json:"models,omitempty"`
	Citations     CitationPolicy    `json:"citations"`
	Runtime       RuntimePolicy     `json:"runtime"`
	Display       DisplayMetadata   `json:"display,omitempty"`
}

type FactorValue struct {
	ID        string               `json:"id"`
	Value     json.RawMessage      `json:"value"`
	Overrides []NodeConfigOverride `json:"overrides,omitempty"`
}

type Factor struct {
	ID     string        `json:"id"`
	Values []FactorValue `json:"values"`
}

type NodeConfigOverride struct {
	NodeID string          `json:"nodeId"`
	Config json.RawMessage `json:"config"`
}

type Variant struct {
	ID       string          `json:"id"`
	Pipeline PipelineIR      `json:"pipeline"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

type Study struct {
	SchemaVersion string            `json:"schemaVersion"`
	Variants      []Variant         `json:"variants"`
	Factors       []Factor          `json:"factors,omitempty"`
	Bindings      []ArtifactBinding `json:"bindings"`
	Dataset       DatasetBinding    `json:"dataset"`
	Measures      []Measure         `json:"measures"`
	Replicates    int               `json:"replicates"`
	Acceptance    json.RawMessage   `json:"acceptance,omitempty"`
	Display       DisplayMetadata   `json:"display,omitempty"`
}

type FactorSelection struct {
	FactorID string          `json:"factorId"`
	ValueID  string          `json:"valueId"`
	Value    json.RawMessage `json:"value"`
}

type PipelineExecution struct {
	SchemaVersion string            `json:"schemaVersion"`
	Pipeline      PipelineIR        `json:"pipeline"`
	Bindings      []ArtifactBinding `json:"bindings"`
	Dataset       DatasetBinding    `json:"dataset"`
	Measures      []Measure         `json:"measures"`
	VariantID     string            `json:"variantId"`
	Factors       []FactorSelection `json:"factors"`
	CellID        string            `json:"cellId"`
}

type ExpandedCell struct {
	ID         string            `json:"id"`
	VariantID  string            `json:"variantId"`
	Factors    []FactorSelection `json:"factors"`
	Replicates int               `json:"replicates"`
	Execution  PipelineExecution `json:"execution"`
}
