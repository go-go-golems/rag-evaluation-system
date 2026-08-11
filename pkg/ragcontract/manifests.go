package ragcontract

import "encoding/json"

const (
	CorpusManifestSchema         = "rag-corpus-snapshot-manifest/v2"
	UnitSetManifestSchema        = "rag-unit-set-manifest/v2"
	ChunkSetManifestSchema       = "rag-chunk-set-manifest/v2"
	RepresentationManifestSchema = "rag-representation-set-manifest/v2"
	EmbeddingManifestSchema      = "rag-embedding-set-manifest/v2"
	IndexManifestSchema          = "rag-index-manifest/v2"
	EvaluationManifestSchema     = "rag-evaluation-dataset-manifest/v2"
	ModelManifestSchema          = "rag-model-manifest/v1"
	PromptManifestSchema         = "rag-prompt-manifest/v1"
)

type ParentDigest struct {
	Role          string `json:"role"`
	Digest        string `json:"digest"`
	SchemaVersion string `json:"schemaVersion"`
}
type Production struct {
	Operator OperatorRef     `json:"operator"`
	Config   json.RawMessage `json:"config"`
}
type ManifestBase struct {
	SchemaVersion string         `json:"schemaVersion"`
	Digest        string         `json:"digest"`
	Parents       []ParentDigest `json:"parents"`
	Production    *Production    `json:"production,omitempty"`
}
type CorpusManifest struct {
	ManifestBase
	SourceNamespace string          `json:"sourceNamespace"`
	RecordSchema    string          `json:"recordSchema"`
	Ordering        string          `json:"ordering"`
	RecordCount     int64           `json:"recordCount"`
	MetadataSchema  json.RawMessage `json:"metadataSchema,omitempty"`
}
type UnitSetManifest struct {
	ManifestBase
	UnitCount      int64  `json:"unitCount"`
	IdentitySchema string `json:"identitySchema"`
}
type ChunkSetManifest struct {
	ManifestBase
	ChunkCount       int64  `json:"chunkCount"`
	RangeUnit        string `json:"rangeUnit"`
	UnicodePolicy    string `json:"unicodePolicy"`
	EmptyInputPolicy string `json:"emptyInputPolicy"`
}
type RepresentationSetManifest struct {
	ManifestBase
	RepresentationCount int64    `json:"representationCount"`
	Kinds               []string `json:"kinds"`
	EvidenceRoles       []string `json:"evidenceRoles"`
}
type EmbeddingSetManifest struct {
	ManifestBase
	VectorCount         int64  `json:"vectorCount"`
	Dimensions          int    `json:"dimensions"`
	Distance            string `json:"distance"`
	Normalization       string `json:"normalization"`
	ModelManifestDigest string `json:"modelManifestDigest"`
}
type IndexManifest struct {
	ManifestBase
	Engine              string   `json:"engine"`
	EngineVersion       string   `json:"engineVersion"`
	RepresentationKinds []string `json:"representationKinds"`
	VectorDimensions    int      `json:"vectorDimensions,omitempty"`
	Distance            string   `json:"distance,omitempty"`
	DocumentCount       int64    `json:"documentCount"`
	ArtifactTreeDigest  string   `json:"artifactTreeDigest"`
}
type EvaluationDatasetManifest struct {
	ManifestBase
	DatasetID       string          `json:"datasetId"`
	Split           string          `json:"split"`
	Status          string          `json:"status"`
	RelevanceTarget string          `json:"relevanceTarget"`
	QueryCount      int64           `json:"queryCount"`
	GradeSchema     json.RawMessage `json:"gradeSchema"`
}
type ModelManifest struct {
	ManifestBase
	ProviderAdapterVersion string          `json:"providerAdapterVersion"`
	ModelID                string          `json:"modelId"`
	ModelDigest            string          `json:"modelDigest"`
	Dimensions             int             `json:"dimensions,omitempty"`
	Tokenization           string          `json:"tokenization"`
	Truncation             string          `json:"truncation"`
	Normalization          string          `json:"normalization"`
	ImplementationVersion  string          `json:"implementationVersion"`
	RequestParameters      json.RawMessage `json:"requestParameters"`
}
type PromptManifest struct {
	ManifestBase
	PromptID       string `json:"promptId"`
	TemplateDigest string `json:"templateDigest"`
	InputSchema    string `json:"inputSchema"`
	OutputSchema   string `json:"outputSchema"`
}
