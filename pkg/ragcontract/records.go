package ragcontract

import "encoding/json"

type SourceRange struct {
	SourceID     string `json:"sourceId"`
	ByteStart    int64  `json:"byteStart,omitempty"`
	ByteEnd      int64  `json:"byteEnd,omitempty"`
	OrdinalStart int64  `json:"ordinalStart,omitempty"`
	OrdinalEnd   int64  `json:"ordinalEnd,omitempty"`
}
type UnitRecord struct {
	ID              string          `json:"id"`
	MemberSourceIDs []string        `json:"memberSourceIds"`
	Ranges          []SourceRange   `json:"ranges"`
	ContentDigest   string          `json:"contentDigest"`
	Metadata        json.RawMessage `json:"metadata,omitempty"`
}
type ChunkRecord struct {
	ID           string      `json:"id"`
	ParentUnitID string      `json:"parentUnitId"`
	ByteStart    int64       `json:"byteStart"`
	ByteEnd      int64       `json:"byteEnd"`
	LogicalStart int64       `json:"logicalStart"`
	LogicalEnd   int64       `json:"logicalEnd"`
	TextDigest   string      `json:"textDigest"`
	Chunker      OperatorRef `json:"chunker"`
	Citation     CitationRef `json:"citation"`
}
type DerivationRef struct {
	Operator                OperatorRef     `json:"operator"`
	ModelManifestDigest     string          `json:"modelManifestDigest,omitempty"`
	PromptManifestDigest    string          `json:"promptManifestDigest,omitempty"`
	Parameters              json.RawMessage `json:"parameters"`
	SeedPolicy              *SeedPolicy     `json:"seedPolicy,omitempty"`
	ParentDigest            string          `json:"parentDigest"`
	ParentRepresentationIDs []string        `json:"parentRepresentationIds,omitempty"`
	SourceRecordIDs         []string        `json:"sourceRecordIds"`
	OutputDigest            string          `json:"outputDigest"`
	InputTokens             int64           `json:"inputTokens,omitempty"`
	OutputTokens            int64           `json:"outputTokens,omitempty"`
	DurationMilliseconds    int64           `json:"durationMilliseconds,omitempty"`
	CacheOutcome            string          `json:"cacheOutcome,omitempty"`
}
type RepresentationRecord struct {
	ID            string         `json:"id"`
	Kind          string         `json:"kind"`
	ParentChunkID string         `json:"parentChunkId"`
	ParentUnitID  string         `json:"parentUnitId"`
	ContentDigest string         `json:"contentDigest"`
	EvidenceRole  string         `json:"evidenceRole"`
	Derivation    *DerivationRef `json:"derivation,omitempty"`
	Citation      CitationRef    `json:"citation"`
}
type EmbeddingRecord struct {
	RepresentationID    string `json:"representationId"`
	ModelManifestDigest string `json:"modelManifestDigest"`
	Dimensions          int    `json:"dimensions"`
	VectorDigest        string `json:"vectorDigest"`
	Position            int64  `json:"position"`
}
