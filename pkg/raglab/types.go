// Package raglab defines typed, reproducible RAG laboratory experiment plans.
// It intentionally contains no Goja, database, or provider dependencies.
package raglab

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/experimentspec"
	"github.com/pkg/errors"
)

type ArtifactKind string

const (
	CorpusSnapshotArtifact    ArtifactKind = "corpusSnapshot"
	ChunkSetArtifact          ArtifactKind = "chunkSet"
	EmbeddingSetArtifact      ArtifactKind = "embeddingSet"
	BM25IndexArtifact         ArtifactKind = "bm25Index"
	EvaluationDatasetArtifact ArtifactKind = "evaluationDataset"
	RepresentationSetArtifact ArtifactKind = "representationSet"
)

type ArtifactRef struct {
	Kind ArtifactKind `json:"kind"`
	ID   string       `json:"id"`
}

func Artifact(kind ArtifactKind, id string) ArtifactRef { return ArtifactRef{Kind: kind, ID: id} }
func CorpusSnapshot(id string) ArtifactRef              { return Artifact(CorpusSnapshotArtifact, id) }
func ChunkSet(id string) ArtifactRef                    { return Artifact(ChunkSetArtifact, id) }
func EmbeddingSet(id string) ArtifactRef                { return Artifact(EmbeddingSetArtifact, id) }
func BM25Index(id string) ArtifactRef                   { return Artifact(BM25IndexArtifact, id) }
func EvaluationDataset(id string) ArtifactRef           { return Artifact(EvaluationDatasetArtifact, id) }
func RepresentationSet(id string) ArtifactRef           { return Artifact(RepresentationSetArtifact, id) }

type RelevanceGrade struct {
	Name    string `json:"name"`
	Ordinal int    `json:"ordinal"`
}

var grades = map[string]RelevanceGrade{
	"0_FAIL":          {Name: "0_FAIL", Ordinal: 0},
	"1_PARTIAL":       {Name: "1_PARTIAL", Ordinal: 1},
	"2_SUBSTANTIAL":   {Name: "2_SUBSTANTIAL", Ordinal: 2},
	"3_AUTHORITATIVE": {Name: "3_AUTHORITATIVE", Ordinal: 3},
}

func Grade(name string) (RelevanceGrade, error) {
	grade, ok := grades[name]
	if !ok {
		return RelevanceGrade{}, errors.Errorf("RAG_UNKNOWN_GRADE: %q", name)
	}
	return grade, nil
}

type ValidationSeverity string

const (
	ValidationErrorSeverity   ValidationSeverity = "error"
	ValidationWarningSeverity ValidationSeverity = "warning"
)

type ValidationIssue struct {
	Code     string             `json:"code"`
	Path     string             `json:"path"`
	Message  string             `json:"message"`
	Severity ValidationSeverity `json:"severity"`
}

type ValidationReport struct {
	Issues []ValidationIssue `json:"issues"`
}

func (r ValidationReport) OK() bool {
	for _, issue := range r.Issues {
		if issue.Severity == ValidationErrorSeverity {
			return false
		}
	}
	return true
}

func (r *ValidationReport) add(code, path, message string) {
	r.Issues = append(r.Issues, ValidationIssue{Code: code, Path: path, Message: message, Severity: ValidationErrorSeverity})
}

func (r *ValidationReport) warn(code, path, message string) {
	r.Issues = append(r.Issues, ValidationIssue{Code: code, Path: path, Message: message, Severity: ValidationWarningSeverity})
}

func (r *ValidationReport) Normalize() {
	sort.Slice(r.Issues, func(i, j int) bool {
		if r.Issues[i].Path == r.Issues[j].Path {
			return r.Issues[i].Code < r.Issues[j].Code
		}
		return r.Issues[i].Path < r.Issues[j].Path
	})
}

type ValidationError struct{ Report ValidationReport }

func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.Report.Issues))
	for _, issue := range e.Report.Issues {
		if issue.Severity == ValidationErrorSeverity {
			parts = append(parts, fmt.Sprintf("%s at %s: %s", issue.Code, issue.Path, issue.Message))
		}
	}
	return strings.Join(parts, "; ")
}

type CollapseScope string

const (
	CollapseNone        CollapseScope = "none"
	CollapseParentChunk CollapseScope = "parentChunk"
	CollapseDocument    CollapseScope = "document"
)

type RetrievalBackend string

const (
	BM25Backend   RetrievalBackend = "bm25"
	VectorBackend RetrievalBackend = "vector"
)

type RepresentationKind string

const (
	RawChunksRepresentation RepresentationKind = "rawChunks"
	SummaryRepresentation   RepresentationKind = "summary"
	QuestionRepresentation  RepresentationKind = "question"
)

type RepresentationSpec struct {
	Name       string             `json:"name"`
	Kind       RepresentationKind `json:"kind"`
	ArtifactID string             `json:"artifact_id,omitempty"`
	Parent     string             `json:"parent,omitempty"`
}

type FilterSpec struct {
	SourceIDs      []string          `json:"source_ids,omitempty"`
	DocumentIDs    []string          `json:"document_ids,omitempty"`
	ContentTypes   []string          `json:"content_types,omitempty"`
	MetadataEquals map[string]string `json:"metadata_equals,omitempty"`
}

type ChannelSpec struct {
	Name           string           `json:"name"`
	Backend        RetrievalBackend `json:"backend"`
	Representation string           `json:"representation"`
	TopK           int              `json:"top_k"`
	Filter         FilterSpec       `json:"filter,omitempty"`
}

type FusionSpec struct {
	Kind         string             `json:"kind"`
	RankConstant int                `json:"rank_constant"`
	Weights      map[string]float64 `json:"weights,omitempty"`
}

type RetrievalPlan struct {
	Channels []ChannelSpec `json:"channels"`
	Filter   FilterSpec    `json:"filter,omitempty"`
	Fusion   *FusionSpec   `json:"fusion,omitempty"`
	Collapse CollapseScope `json:"collapse"`
	Results  int           `json:"results"`
}

type MetricsPlan struct {
	RelevanceAt        *RelevanceGrade `json:"relevance_at,omitempty"`
	PrecisionAt        []int           `json:"precision_at,omitempty"`
	RecallAt           []int           `json:"recall_at,omitempty"`
	HitRateAt          []int           `json:"hit_rate_at,omitempty"`
	NDCGAt             int             `json:"ndcg_at,omitempty"`
	MRR                bool            `json:"mrr,omitempty"`
	MeanRelevantRecall int             `json:"mean_relevant_recall_at,omitempty"`
	Abstention         bool            `json:"abstention,omitempty"`
}

func (m MetricsPlan) RequiresRelevance() bool {
	return len(m.PrecisionAt) > 0 || len(m.RecallAt) > 0 || len(m.HitRateAt) > 0 || m.NDCGAt > 0 || m.MRR || m.MeanRelevantRecall > 0
}

type Provenance struct {
	Fragments []string          `json:"fragments,omitempty"`
	Notes     []string          `json:"notes,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
}

type InputSpec struct {
	CorpusSnapshot    ArtifactRef          `json:"corpus_snapshot"`
	ChunkSet          ArtifactRef          `json:"chunk_set"`
	BM25Index         *ArtifactRef         `json:"bm25_index,omitempty"`
	EmbeddingSet      *ArtifactRef         `json:"embedding_set,omitempty"`
	EvaluationDataset ArtifactRef          `json:"evaluation_dataset"`
	Representations   []RepresentationSpec `json:"representations"`
}

type ExperimentSpecification struct {
	SchemaVersion string        `json:"schema_version"`
	Fingerprint   string        `json:"fingerprint"`
	Name          string        `json:"name"`
	Provenance    Provenance    `json:"provenance"`
	Inputs        InputSpec     `json:"inputs"`
	Retrieval     RetrievalPlan `json:"retrieval"`
	Metrics       MetricsPlan   `json:"metrics"`
}

// PersistenceInput is the only conversion from the typed authoring model to
// the pre-existing immutable experiment storage contract.
func (s ExperimentSpecification) PersistenceInput() experimentspec.Input {
	config := map[string]any{
		"name":            s.Name,
		"provenance":      s.Provenance,
		"representations": s.Inputs.Representations,
		"retrieval":       s.Retrieval,
		"metrics":         s.Metrics,
	}
	input := experimentspec.Input{
		CorpusSnapshotID:    s.Inputs.CorpusSnapshot.ID,
		ChunkSetID:          s.Inputs.ChunkSet.ID,
		EvaluationDatasetID: s.Inputs.EvaluationDataset.ID,
		Config:              config,
	}
	if s.Inputs.BM25Index != nil {
		input.BM25ArtifactID = s.Inputs.BM25Index.ID
	}
	if s.Inputs.EmbeddingSet != nil {
		input.EmbeddingSetID = s.Inputs.EmbeddingSet.ID
	}
	return input
}
