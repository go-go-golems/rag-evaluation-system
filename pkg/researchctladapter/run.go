package researchctladapter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/researchctl/pkg/lab"
)

const RunnerName = "rag-worker"
const RunnerVersion = "rag-worker/v2"
const ProtocolVersion = "researchctl-runner-stdio/v1"

type WorkerCommand struct {
	Executable string
	Args       []string
}
type RunOptions struct {
	ResearchctlCommand string
	Project            string
	Database           string
	ExperimentID       string
	Worker             WorkerCommand
	MaxAttempts        int
	Timeout            time.Duration
	SecretEnvironment  []string
	OutputDirectory    string
}
type RunResult struct {
	SpecificationID string          `json:"specificationId"`
	Output          json.RawMessage `json:"output"`
}

func CheckWorker(ctx context.Context, worker WorkerCommand) error {
	if worker.Executable == "" {
		return fmt.Errorf("RAG_WORKER_COMMAND_REQUIRED")
	}
	command := exec.CommandContext(ctx, worker.Executable, worker.Args...)
	request := map[string]any{"protocolVersion": ProtocolVersion, "attempt": map[string]any{"specification": map[string]any{"canonicalIdentity": map[string]any{"domain": ragcontract.Domain, "domainSchemaVersion": ragcontract.DomainSchemaVersion, "domainConfig": map[string]any{}}}}}
	payload, _ := json.Marshal(request)
	command.Stdin = bytes.NewReader(append(payload, '\n'))
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		return fmt.Errorf("RAG_WORKER_LAUNCH: %w", err)
	}
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		_ = command.Wait()
		return fmt.Errorf("RAG_WORKER_HELLO_MISSING: %s", stderr.String())
	}
	var frame struct {
		Type  string `json:"type"`
		Hello struct {
			ProtocolVersion string              `json:"protocolVersion"`
			Runner          lab.RunnerRecord    `json:"runner"`
			Domains         []lab.DomainVersion `json:"domains"`
		} `json:"hello"`
	}
	if err := json.Unmarshal(scanner.Bytes(), &frame); err != nil {
		_ = command.Process.Kill()
		return err
	}
	_ = command.Process.Kill()
	_ = command.Wait()
	if frame.Type != "hello" || frame.Hello.ProtocolVersion != ProtocolVersion || frame.Hello.Runner.Name != RunnerName || frame.Hello.Runner.ResolvedVersion != RunnerVersion {
		return fmt.Errorf("RAG_WORKER_CAPABILITY: %#v", frame)
	}
	if len(frame.Hello.Domains) != 1 || frame.Hello.Domains[0].Domain != ragcontract.Domain || frame.Hello.Domains[0].SchemaVersion != ragcontract.DomainSchemaVersion {
		return fmt.Errorf("RAG_WORKER_DOMAIN_CAPABILITY")
	}
	if workerUsesRealProfile(worker.Args) {
		if err := checkWorkerCapabilities(ctx, worker); err != nil {
			return err
		}
	}
	return nil
}

func checkWorkerCapabilities(ctx context.Context, worker WorkerCommand) error {
	command := exec.CommandContext(ctx, worker.Executable, append(append([]string{}, worker.Args...), "--capabilities")...)
	output, err := command.Output()
	if err != nil {
		return fmt.Errorf("RAG_WORKER_PROVIDER_CAPABILITY")
	}
	var capabilities struct {
		SchemaVersion    string   `json:"schemaVersion"`
		FixtureProviders bool     `json:"fixtureProviders"`
		Capabilities     []string `json:"capabilities"`
	}
	if err := json.Unmarshal(output, &capabilities); err != nil || capabilities.SchemaVersion != "rag-provider-capabilities/v1" || capabilities.FixtureProviders || !hasCapabilities(capabilities.Capabilities, "generator", "embedder", "reranker", "schema-validator", "persistent-cache") {
		return fmt.Errorf("RAG_WORKER_PROVIDER_CAPABILITY")
	}
	return nil
}

func workerUsesRealProfile(args []string) bool {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--provider-profile" && args[i+1] == "real" {
			return true
		}
	}
	return false
}

func hasCapabilities(got []string, required ...string) bool {
	available := map[string]bool{}
	for _, capability := range got {
		available[capability] = true
	}
	for _, capability := range required {
		if !available[capability] {
			return false
		}
	}
	return true
}

func WriteSpecification(path string, specification lab.SpecificationRecord) error {
	data, err := lab.CanonicalJSON(specification)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func ExecuteSpecification(ctx context.Context, specification lab.SpecificationRecord, options RunOptions) (RunResult, error) {
	if err := CheckWorker(ctx, options.Worker); err != nil {
		return RunResult{}, err
	}
	if options.ResearchctlCommand == "" {
		options.ResearchctlCommand = "researchctl"
	}
	if options.MaxAttempts == 0 {
		options.MaxAttempts = 1
	}
	if options.OutputDirectory == "" {
		options.OutputDirectory = os.TempDir()
	}
	if err := os.MkdirAll(options.OutputDirectory, 0o755); err != nil {
		return RunResult{}, err
	}
	data, err := lab.CanonicalJSON(specification)
	if err != nil {
		return RunResult{}, err
	}
	path := filepath.Join(options.OutputDirectory, specification.ID+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return RunResult{}, err
	}
	args := []string{"experiment", "execute-spec", path, "--project", options.Project, "--experiment-id", options.ExperimentID, "--runner-command", options.Worker.Executable, "--runner-name", RunnerName, "--runner-version", RunnerVersion, "--max-attempts", fmt.Sprint(options.MaxAttempts), "--output", "json"}
	for _, arg := range options.Worker.Args {
		args = append(args, "--runner-arg", arg)
	}
	if options.Database != "" {
		args = append(args, "--database", options.Database)
	}
	if options.Timeout > 0 {
		args = append(args, "--timeout", options.Timeout.String())
	}
	for _, name := range options.SecretEnvironment {
		args = append(args, "--secret-env", name)
	}
	command := exec.CommandContext(ctx, options.ResearchctlCommand, args...)
	output, runErr := command.CombinedOutput()
	if runErr != nil {
		return RunResult{SpecificationID: specification.ID, Output: json.RawMessage(output)}, fmt.Errorf("RAG_RESEARCHCTL_EXECUTE: %w: %s", runErr, output)
	}
	canonical, err := lab.CanonicalizeJSON(output)
	if err != nil {
		return RunResult{}, fmt.Errorf("RAG_RESEARCHCTL_OUTPUT: %w", err)
	}
	return RunResult{SpecificationID: specification.ID, Output: canonical}, nil
}

func ReconstructSpecification(export lab.RunExport) (lab.SpecificationRecord, error) {
	execution, err := ragcontract.DecodeExecution(bytes.NewReader(export.Specification.CanonicalIdentity.DomainConfig))
	if err != nil {
		return lab.SpecificationRecord{}, err
	}
	resolved := ResolvedInputs{ByRole: map[string]ResolvedInput{}}
	for _, reference := range export.Specification.CanonicalIdentity.Inputs {
		binding := ragcontract.ArtifactBinding{SlotID: reference.Role, Role: reference.Role, Kind: reference.Kind, ID: reference.ID, URI: reference.URI, SchemaVersion: reference.SchemaVersion}
		for _, candidate := range execution.Bindings {
			if candidate.Role == reference.Role {
				binding.Digest = candidate.Digest
			}
		}
		if reference.Role == "evaluation-dataset" || reference.Role == "judgments" {
			binding.Digest = execution.Dataset.ManifestDigest
		}
		resolved.ByRole[reference.Role] = ResolvedInput{Reference: reference, Binding: binding}
	}
	return WrapExecution(execution, resolved, export.Specification.DisplayName)
}
