package ragworkflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
)

const StudyBundleSchema = "rag-workflow-study-bundle/v1"

var safeStudyCaseID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

type StudyWorkflowCase struct {
	ID         string
	Execution  ragcontract.PipelineExecution
	Corpus     ragoperators.Corpus
	Dataset    ragoperators.EvaluationDataset
	Replicates int
}

type StudyBundleCase struct {
	ID           string            `json:"id"`
	Replicates   int               `json:"replicates"`
	Factors      map[string]string `json:"factors,omitempty"`
	Execution    FixtureFile       `json:"execution"`
	Corpus       FixtureFile       `json:"corpus"`
	Queries      FixtureFile       `json:"queries"`
	DomainConfig FixtureFile       `json:"domainConfig"`
}

type StudyBundle struct {
	SchemaVersion string            `json:"schemaVersion"`
	Name          string            `json:"name"`
	ExperimentID  string            `json:"experimentId"`
	Plan          FixtureFile       `json:"plan"`
	Cases         []StudyBundleCase `json:"cases"`
}

type planMeasure struct {
	Name      string          `json:"name"`
	ValueKind string          `json:"valueKind"`
	Unit      string          `json:"unit,omitempty"`
	Required  bool            `json:"required"`
	Config    json.RawMessage `json:"config,omitempty"`
}

type planCase struct {
	ID           string            `json:"id"`
	Replicates   int               `json:"replicates"`
	Factors      map[string]string `json:"factors"`
	DomainConfig json.RawMessage   `json:"domainConfig"`
	Inputs       []map[string]any  `json:"inputs"`
	Measures     []planMeasure     `json:"measures"`
}

// WriteStudyBundle compiles canonical RAG executions into immutable Workflow V3
// inputs and a Researchctl experiment plan. outputDirectory must be contained by
// artifactRoot because plan input URIs are resolved relative to that custody root.
func WriteStudyBundle(ctx context.Context, artifactRoot, outputDirectory, name, experimentID string, cases []StudyWorkflowCase, providerPackage *ProviderPackage) (StudyBundle, error) {
	if artifactRoot == "" || outputDirectory == "" || name == "" || experimentID == "" || len(cases) == 0 {
		return StudyBundle{}, fmt.Errorf("RAG_WORKFLOW_STUDY_ARGUMENT")
	}
	root, err := filepath.Abs(artifactRoot)
	if err != nil {
		return StudyBundle{}, err
	}
	output, err := filepath.Abs(outputDirectory)
	if err != nil {
		return StudyBundle{}, err
	}
	relativeOutput, err := filepath.Rel(root, output)
	if err != nil || relativeOutput == ".." || strings.HasPrefix(relativeOutput, ".."+string(filepath.Separator)) {
		return StudyBundle{}, fmt.Errorf("RAG_WORKFLOW_STUDY_OUTPUT_BOUNDARY")
	}
	if err := os.MkdirAll(output, 0o755); err != nil {
		return StudyBundle{}, err
	}
	bundle := StudyBundle{SchemaVersion: StudyBundleSchema, Name: name, ExperimentID: experimentID}
	planCases := make([]planCase, 0, len(cases))
	seen := map[string]bool{}
	for _, item := range cases {
		if !safeStudyCaseID.MatchString(item.ID) || seen[item.ID] || item.Replicates < 1 {
			return StudyBundle{}, fmt.Errorf("RAG_WORKFLOW_STUDY_CASE: %s", item.ID)
		}
		seen[item.ID] = true
		if err := ctx.Err(); err != nil {
			return StudyBundle{}, err
		}
		lowerer := NewLowerer()
		var lowered LoweredExecution
		if providerPackage == nil {
			lowered, err = lowerer.Lower(ctx, item.Execution)
		} else {
			var providerLowerer *Lowerer
			providerLowerer, err = NewProviderLowerer(providerPackage)
			if err == nil {
				lowered, err = providerLowerer.Lower(ctx, item.Execution)
			}
		}
		if err != nil {
			return StudyBundle{}, fmt.Errorf("RAG_WORKFLOW_STUDY_LOWER %s: %w", item.ID, err)
		}
		bindings := map[string]researchrunner.InputBinding{
			"execution": {Role: "workflow-input", Kind: "rag-execution", ID: "execution-" + item.ID},
			"corpus":    {Role: "workflow-input", Kind: "rag-corpus", ID: "corpus-" + item.ID},
			"queries":   {Role: "workflow-input", Kind: "rag-query-set", ID: "queries-" + item.ID},
		}
		var domain researchrunner.WorkflowExecution
		if providerPackage == nil {
			domain, err = BuildRunnerExecution(lowered, bindings)
		} else {
			domain, err = BuildProviderRunnerExecution(lowered, bindings, providerPackage)
		}
		if err != nil {
			return StudyBundle{}, err
		}
		archive, err := BuildQueryArchive(item.Execution, item.Dataset)
		if err != nil {
			return StudyBundle{}, err
		}
		caseDirectory := filepath.Join(output, item.ID)
		if err := os.MkdirAll(caseDirectory, 0o755); err != nil {
			return StudyBundle{}, err
		}
		current := StudyBundleCase{ID: item.ID, Replicates: item.Replicates, Factors: factorIDs(item.Execution.Factors)}
		values := []struct {
			name   string
			schema string
			value  any
			target *FixtureFile
		}{
			{"execution.json", ragcontract.ExecutionSchemaVersion, item.Execution, &current.Execution},
			{"corpus.json", CorpusSchema, item.Corpus, &current.Corpus},
			{"queries.json", researchrunner.SetInputArchiveSchema, archive, &current.Queries},
			{"domain-config.json", researchrunner.DomainSchemaVersion, domain, &current.DomainConfig},
		}
		inputs := make([]map[string]any, 0, 3)
		for _, value := range values {
			file, body, writeErr := writeStudyFile(root, caseDirectory, value.name, value.schema, value.value)
			if writeErr != nil {
				return StudyBundle{}, writeErr
			}
			*value.target = file
			if value.name != "domain-config.json" {
				role, kind, id := "workflow-input", "rag-execution", "execution-"+item.ID
				switch value.name {
				case "corpus.json":
					kind, id = "rag-corpus", "corpus-"+item.ID
				case "queries.json":
					kind, id = "rag-query-set", "queries-"+item.ID
				}
				inputs = append(inputs, map[string]any{"role": role, "kind": kind, "id": id, "digest": file.Digest, "sizeBytes": file.SizeBytes, "schemaVersion": file.SchemaVersion, "mediaType": "application/json", "uri": file.Path})
			}
			_ = body
		}
		sort.Slice(inputs, func(i, j int) bool { return inputs[i]["id"].(string) < inputs[j]["id"].(string) })
		domainBody, err := os.ReadFile(filepath.Join(caseDirectory, "domain-config.json"))
		if err != nil {
			return StudyBundle{}, err
		}
		planCases = append(planCases, planCase{ID: item.ID, Replicates: item.Replicates, Factors: current.Factors, DomainConfig: json.RawMessage(strings.TrimSpace(string(domainBody))), Inputs: inputs, Measures: planMeasures(item.Execution.Measures)})
		bundle.Cases = append(bundle.Cases, current)
	}
	planBody, err := renderStudyPlan(name, experimentID, planCases)
	if err != nil {
		return StudyBundle{}, err
	}
	planPath := filepath.Join(output, "researchctl-plan.js")
	if err := writeImmutableStudyFile(planPath, planBody); err != nil {
		return StudyBundle{}, err
	}
	bundle.Plan = fileIdentity(filepath.ToSlash(filepath.Join(relativeOutput, "researchctl-plan.js")), "researchctl-experiment-plan-js/v1", planBody)
	manifestBody, err := ragcontract.CanonicalJSON(bundle)
	if err != nil {
		return StudyBundle{}, err
	}
	if err := writeImmutableStudyFile(filepath.Join(output, "manifest.json"), append(manifestBody, '\n')); err != nil {
		return StudyBundle{}, err
	}
	return bundle, nil
}

func writeStudyFile(root, directory, name, schema string, value any) (FixtureFile, []byte, error) {
	body, err := ragcontract.CanonicalJSON(value)
	if err != nil {
		return FixtureFile{}, nil, err
	}
	body = append(body, '\n')
	path := filepath.Join(directory, name)
	if err := writeImmutableStudyFile(path, body); err != nil {
		return FixtureFile{}, nil, err
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return FixtureFile{}, nil, err
	}
	return fileIdentity(filepath.ToSlash(relative), schema, body), body, nil
}

func writeImmutableStudyFile(path string, body []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600) // #nosec G304 -- path is constructed below the validated output boundary.
	if err == nil {
		if _, writeErr := file.Write(body); writeErr != nil {
			_ = file.Close()
			_ = os.Remove(path)
			return writeErr
		}
		if closeErr := file.Close(); closeErr != nil {
			_ = os.Remove(path)
			return closeErr
		}
		return nil
	}
	if !os.IsExist(err) {
		return err
	}
	existing, readErr := os.ReadFile(path) // #nosec G304 -- path is constructed below the validated output boundary.
	if readErr != nil {
		return readErr
	}
	if !bytes.Equal(existing, body) {
		return fmt.Errorf("RAG_WORKFLOW_STUDY_CONFLICT: %s", filepath.Base(path))
	}
	return nil
}

func fileIdentity(path, schema string, body []byte) FixtureFile {
	sum := sha256.Sum256(body)
	return FixtureFile{Path: path, SchemaVersion: schema, Digest: "sha256:" + hex.EncodeToString(sum[:]), SizeBytes: int64(len(body))}
}

func planMeasures(measures []ragcontract.Measure) []planMeasure {
	result := make([]planMeasure, len(measures))
	for index, measure := range measures {
		result[index] = planMeasure{Name: measure.Name, ValueKind: measure.ValueKind, Unit: measure.Unit, Required: measure.Required, Config: measure.Config}
	}
	return result
}

func factorIDs(selections []ragcontract.FactorSelection) map[string]string {
	values := make(map[string]string, len(selections))
	for _, selection := range selections {
		values[selection.FactorID] = selection.ValueID
	}
	return values
}

func renderStudyPlan(name, experimentID string, cases []planCase) ([]byte, error) {
	caseBody, err := ragcontract.CanonicalJSON(cases)
	if err != nil {
		return nil, err
	}
	nameBody, _ := json.Marshal(name)
	experimentBody, _ := json.Marshal(experimentID)
	body := fmt.Sprintf(`const research = require("researchctl");
const cases = %s;
function specification(item) {
  return {
    canonicalIdentity: {
      schemaVersion: "researchctl-execution-spec/v1",
      identityScheme: "researchctl-execution-identity/v1",
      domain: "scraper-workflow",
      domainSchemaVersion: "scraper-workflow-execution/v2",
      inputs: item.inputs,
      domainConfig: item.domainConfig,
      requestedMeasures: item.measures,
      factors: item.factors
    },
    displayName: %s + " / " + item.id,
    provenance: {authoring: "rag-eval study compile"},
    labels: {path: "rag-v2-workflow-v3"}
  };
}
module.exports = research.experimentPlan(%s, plan => {
  let current = plan.experiment(%s);
  for (const item of cases) {
    current = current.case(item.id, value => value.specification(specification(item)).factors(item.factors).replicates(item.replicates));
  }
  return current;
});
`, caseBody, nameBody, nameBody, experimentBody)
	return []byte(body), nil
}
