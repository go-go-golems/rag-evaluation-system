package researchctladapter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/researchctl/pkg/lab"
)

// OperationCustodyArtifact is one compact, already-published Workflow V3
// operation-evidence file. URI is relative to the import bundle root; Source
// is the local file used only while building the export and is never persisted.
type OperationCustodyArtifact struct {
	Role          string
	Kind          string
	URI           string
	Source        string
	SchemaVersion string
	MediaType     string
}

// OperationCustodyMetric is a scalar derived solely from the bounded exported
// operation model. It intentionally has no arbitrary metadata or text value.
type OperationCustodyMetric struct {
	Name  string
	Units int64
	Unit  string
	Scope string
}

// OperationCustodyExportInput describes one externally-produced TTC sweep
// custody bundle for researchctl import. IDs and RecordedAt are explicit
// operator-owned values so this conversion does not invent non-deterministic
// run identity or timestamps.
type OperationCustodyExportInput struct {
	Specification lab.SpecificationRecord
	Source        lab.ExportSource
	RunID         string
	AttemptID     string
	RecordedAt    time.Time
	Status        string
	Artifacts     []OperationCustodyArtifact
	Metrics       []OperationCustodyMetric
}

// BuildOperationCustodyRunExport maps compact operation JSONL/manifests and
// scalar reductions into researchctl's domain-neutral immutable import
// contract. It calculates file digests and sizes from the local bundle but
// persists only relative URIs and verified identities.
func BuildOperationCustodyRunExport(input OperationCustodyExportInput) (lab.RunExport, error) {
	if input.RecordedAt.IsZero() {
		return lab.RunExport{}, fmt.Errorf("operation custody recorded time is required")
	}
	if input.Status != "succeeded" && input.Status != "failed" {
		return lab.RunExport{}, fmt.Errorf("operation custody status %q is unsupported", input.Status)
	}
	if strings.TrimSpace(input.Source.Namespace) == "" || strings.TrimSpace(input.Source.ExternalRunID) == "" {
		return lab.RunExport{}, fmt.Errorf("operation custody source identity is required")
	}
	artifacts, err := operationCustodyArtifacts(input.Artifacts, input.RecordedAt)
	if err != nil {
		return lab.RunExport{}, err
	}
	metrics, err := operationCustodyMetrics(input.Metrics)
	if err != nil {
		return lab.RunExport{}, err
	}
	stamp := input.RecordedAt.UTC().Format(time.RFC3339Nano)
	attempt := lab.AttemptRecord{ID: input.AttemptID, Index: 1, Runner: lab.RunnerRecord{Name: "rag-ttc-operation-custody", ResolvedVersion: "rag-ttc-operation-custody/v1"}, Environment: json.RawMessage(`{"schemaVersion":"rag-ttc-operation-custody-environment/v1"}`), CreatedAt: stamp, Artifacts: artifacts, Metrics: metrics, TerminalSummary: lab.AttemptSummary{Status: input.Status, ProducerFinishedAt: stamp, RecordedAt: stamp, Payload: json.RawMessage(`{"schemaVersion":"rag-ttc-operation-custody-summary/v1"}`)}}
	result := lab.RunExport{SchemaVersion: lab.RunExportSchemaVersion, Source: input.Source, Specification: input.Specification, Run: lab.RunRecord{ID: input.RunID, ReplicateIndex: 1, CreatedAt: stamp}, Attempts: []lab.AttemptRecord{attempt}, RunSummary: lab.RunSummary{Status: input.Status, RecordedAt: stamp, Payload: json.RawMessage(`{"schemaVersion":"rag-ttc-operation-custody-summary/v1"}`)}}
	if input.Status == "succeeded" {
		result.RunSummary.SelectedAttemptID = &attempt.ID
	}
	if err := lab.ValidateRunExport(result); err != nil {
		return lab.RunExport{}, fmt.Errorf("validate operation custody run export: %w", err)
	}
	return result, nil
}

func operationCustodyArtifacts(inputs []OperationCustodyArtifact, recordedAt time.Time) ([]lab.RunArtifact, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("operation custody artifacts are required")
	}
	ordered := append([]OperationCustodyArtifact(nil), inputs...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].URI < ordered[j].URI })
	result := make([]lab.RunArtifact, 0, len(ordered))
	seen := map[string]bool{}
	for index, input := range ordered {
		if input.Role == "" || input.Kind == "" || input.URI == "" || input.Source == "" || seen[input.URI] {
			return nil, fmt.Errorf("operation custody artifact %d is invalid or duplicate", index)
		}
		seen[input.URI] = true
		if filepath.IsAbs(input.URI) || strings.Contains(filepath.ToSlash(input.URI), "../") {
			return nil, fmt.Errorf("operation custody artifact URI %q is not relative", input.URI)
		}
		info, err := os.Stat(input.Source)
		if err != nil || !info.Mode().IsRegular() {
			if err == nil {
				err = fmt.Errorf("not a regular file")
			}
			return nil, fmt.Errorf("operation custody artifact %q: %w", input.URI, err)
		}
		digest, err := lab.DigestFile(input.Source)
		if err != nil {
			return nil, err
		}
		result = append(result, lab.RunArtifact{Ordinal: int64(index + 1), Role: input.Role, Kind: input.Kind, URI: filepath.ToSlash(input.URI), Digest: digest, SizeBytes: info.Size(), SchemaVersion: input.SchemaVersion, MediaType: input.MediaType, Metadata: json.RawMessage(`{}`), Verification: lab.ArtifactVerification{Status: "verified", VerifiedAt: recordedAt.UTC().Format(time.RFC3339Nano)}})
	}
	return result, nil
}

func operationCustodyMetrics(inputs []OperationCustodyMetric) ([]lab.MetricRecord, error) {
	ordered := append([]OperationCustodyMetric(nil), inputs...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Name < ordered[j].Name })
	result := make([]lab.MetricRecord, 0, len(ordered))
	seen := map[string]bool{}
	for index, input := range ordered {
		if input.Name == "" || input.Units < 0 || seen[input.Name] {
			return nil, fmt.Errorf("operation custody metric %d is invalid or duplicate", index)
		}
		seen[input.Name] = true
		value := json.RawMessage(fmt.Sprintf("%d", input.Units))
		projection := float64(input.Units)
		result = append(result, lab.MetricRecord{Ordinal: int64(index + 1), Name: input.Name, Scope: input.Scope, ValueKind: "number", Value: value, NumericProjection: &projection, Unit: input.Unit, Metadata: json.RawMessage(`{}`)})
	}
	return result, nil
}
