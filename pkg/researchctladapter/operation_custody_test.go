package researchctladapter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/researchctl/pkg/lab"
)

func TestBuildOperationCustodyRunExportStagesVerifiedCompactArtifacts(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "operations", "cell-00.manifest.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"schemaVersion":"scraper-workflow-v3-external-operations/v1"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	identity := lab.ExecutionIdentity{SchemaVersion: lab.ExecutionSpecSchemaVersion, IdentityScheme: lab.ExecutionIdentityScheme, Domain: "rag", DomainSchemaVersion: "rag/v2", DomainConfig: []byte(`{"schemaVersion":"rag-ttc-v3-sweep-custody/v1"}`)}
	id, _, err := lab.ExecutionID(identity)
	if err != nil {
		t.Fatal(err)
	}
	specification := lab.SpecificationRecord{ID: id, IdentityScheme: lab.ExecutionIdentityScheme, CanonicalIdentity: identity, DisplayName: "TTC operation custody"}
	recordedAt := time.Date(2026, 7, 23, 12, 0, 0, 0, time.UTC)
	export, err := BuildOperationCustodyRunExport(OperationCustodyExportInput{Specification: specification, Source: lab.ExportSource{Namespace: "rag-ttc", ExternalRunID: "sweep-fixture-1"}, RunID: "run_" + strings.Repeat("0", 26), AttemptID: "attempt_" + strings.Repeat("1", 26), RecordedAt: recordedAt, Status: "succeeded", Artifacts: []OperationCustodyArtifact{{Role: "operation-manifest", Kind: "workflow-operation-manifest", URI: "operations/cell-00.manifest.json", Source: path, SchemaVersion: "scraper-workflow-v3-external-operations/v1", MediaType: "application/json"}}, Metrics: []OperationCustodyMetric{{Name: "operation.completed", Units: 2, Unit: "operations"}, {Name: "operation.elapsed_micros", Units: 17, Unit: "micros"}}})
	if err != nil {
		t.Fatal(err)
	}
	checks, err := lab.VerifyRunExportArtifacts(root, export)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 1 || !checks[0].OK || len(export.Attempts[0].Metrics) != 2 || export.Attempts[0].Metrics[0].Name != "operation.completed" {
		t.Fatalf("unexpected export: %#v %#v", checks, export)
	}
}

func TestBuildOperationCustodyRunExportRejectsUnsafeURI(t *testing.T) {
	_, err := operationCustodyArtifacts([]OperationCustodyArtifact{{Role: "operation", Kind: "manifest", URI: "../secret", Source: "ignored"}}, time.Now())
	if err == nil {
		t.Fatal("expected unsafe URI rejection")
	}
}
