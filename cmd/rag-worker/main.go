package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproviders"
	"github.com/spf13/cobra"
)

const protocolVersion = "researchctl-runner-stdio/v1"

type artifactRef struct{ Role, SchemaVersion string }
type resolvedInput struct {
	Reference artifactRef `json:"reference"`
	Path      string      `json:"path"`
}
type request struct {
	ProtocolVersion string `json:"protocolVersion"`
	Attempt         struct {
		Specification struct {
			CanonicalIdentity struct {
				Domain, DomainSchemaVersion string
				DomainConfig                json.RawMessage
			}
		}
	}
	Inputs []resolvedInput `json:"inputs"`
}
type frame struct {
	Type     string `json:"type"`
	Hello    any    `json:"hello,omitempty"`
	Event    any    `json:"event,omitempty"`
	Trace    any    `json:"trace,omitempty"`
	Metric   any    `json:"metric,omitempty"`
	Artifact any    `json:"artifact,omitempty"`
	Error    any    `json:"error,omitempty"`
	Complete any    `json:"complete,omitempty"`
}
type observer struct{ encoder *json.Encoder }

func (o observer) Event(_ context.Context, v ragoperators.Event) error {
	return o.encoder.Encode(frame{Type: "event", Event: map[string]any{"type": v.Type, "payload": v.Payload}})
}
func (o observer) Trace(_ context.Context, v ragcontract.QueryTrace) error {
	return o.encoder.Encode(frame{Type: "trace", Trace: map[string]any{"kind": ragcontract.TraceSchemaVersion, "value": v}})
}
func (o observer) Metric(_ context.Context, v ragoperators.Metric) error {
	metadata := v.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return o.encoder.Encode(frame{Type: "metric", Metric: map[string]any{"name": v.Name, "value": v.Value, "numericProjection": v.Numeric, "unit": v.Unit, "metadata": metadata}})
}
func (o observer) Artifact(_ context.Context, v ragoperators.Artifact) error {
	metadata := v.Metadata
	if len(metadata) == 0 {
		metadata = json.RawMessage(`{}`)
	}
	return o.encoder.Encode(frame{Type: "artifact", Artifact: map[string]any{"role": v.Role, "kind": v.Kind, "name": v.Name, "schemaVersion": v.SchemaVersion, "mediaType": v.MediaType, "metadata": metadata, "data": v.Data}})
}

type workerCommand struct{ *cmds.CommandDescription }

type workerSettings struct {
	ProviderProfile string `glazed:"provider-profile"`
	ProviderConfig  string `glazed:"provider-config"`
}

var _ cmds.BareCommand = (*workerCommand)(nil)

func newWorkerCommand() (*cobra.Command, error) {
	description := cmds.NewCommandDescription(
		"rag-worker",
		cmds.WithShort("Execute a canonical RAG request over the stdio runner protocol"),
		cmds.WithFlags(
			fields.New("provider-profile", fields.TypeString, fields.WithHelp("Explicit provider profile: fixtures or real"), fields.WithRequired(true)),
			fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only real-provider configuration YAML; required for provider-profile=real")),
		),
	)
	command := &workerCommand{CommandDescription: description}
	return cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-worker", ShortHelpSections: []string{schema.DefaultSlug}}))
}

func (c *workerCommand) Run(ctx context.Context, parsed *values.Values) error {
	settings := &workerSettings{}
	if err := parsed.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return err
	}
	executeWorker(ctx, settings)
	return nil
}

func main() {
	command, err := newWorkerCommand()
	cobra.CheckErr(err)
	cobra.CheckErr(command.Execute())
}

func executeWorker(parent context.Context, settings *workerSettings) {
	encoder := json.NewEncoder(os.Stdout)
	var input request
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fail(encoder, "RAG_WORKER_REQUEST", err)
		return
	}
	hello := map[string]any{"protocolVersion": protocolVersion, "runner": map[string]any{"name": "rag-worker", "resolvedVersion": "rag-worker/v2"}, "domains": []map[string]any{{"domain": ragcontract.Domain, "schemaVersion": ragcontract.DomainSchemaVersion}}}
	if err := encoder.Encode(frame{Type: "hello", Hello: hello}); err != nil {
		return
	}
	if input.ProtocolVersion != protocolVersion {
		fail(encoder, "RAG_WORKER_PROTOCOL", fmt.Errorf("unsupported %q", input.ProtocolVersion))
		return
	}
	if input.Attempt.Specification.CanonicalIdentity.Domain != ragcontract.Domain || input.Attempt.Specification.CanonicalIdentity.DomainSchemaVersion != ragcontract.DomainSchemaVersion {
		fail(encoder, "RAG_WORKER_DOMAIN", fmt.Errorf("unsupported domain"))
		return
	}
	execution, err := ragcontract.DecodeExecution(strings.NewReader(string(input.Attempt.Specification.CanonicalIdentity.DomainConfig)))
	if err != nil {
		fail(encoder, "RAG_WORKER_EXECUTION", err)
		return
	}
	var corpusArtifact ragoperators.CorpusArtifact
	var evaluationArtifact ragoperators.EvaluationArtifact
	for _, resolved := range input.Inputs {
		if resolved.Path == "" {
			continue
		}
		file, err := os.Open(resolved.Path)
		if err != nil {
			fail(encoder, "RAG_WORKER_INPUT", err)
			return
		}
		switch {
		case resolved.Reference.SchemaVersion == ragcontract.CorpusManifestSchema || resolved.Reference.Role == "corpus":
			err = decodeStrict(file, &corpusArtifact)
		case resolved.Reference.SchemaVersion == ragcontract.EvaluationManifestSchema || resolved.Reference.Role == "evaluation-dataset":
			err = decodeStrict(file, &evaluationArtifact)
		}
		_ = file.Close()
		if err != nil {
			fail(encoder, "RAG_WORKER_INPUT", err)
			return
		}
	}
	if err := ragoperators.ValidateInputArtifacts(execution, corpusArtifact, evaluationArtifact); err != nil {
		fail(encoder, "RAG_WORKER_INPUT_LINEAGE", err)
		return
	}
	if len(corpusArtifact.Corpus.Records) == 0 || len(evaluationArtifact.Dataset.Queries) == 0 {
		fail(encoder, "RAG_WORKER_INPUT", fmt.Errorf("corpus and evaluation dataset are required"))
		return
	}
	ctx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	engine := ragengine.New(nil)
	var options ragengine.Options
	switch settings.ProviderProfile {
	case "fixtures":
		if settings.ProviderConfig != "" {
			fail(encoder, "RAG_WORKER_PROVIDER_PROFILE", fmt.Errorf("fixture profile does not accept provider config"))
			return
		}
		fixtures := ragoperators.NewFixtureProviders()
		options = ragengine.Options{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache()}
	case "real":
		if settings.ProviderConfig == "" {
			fail(encoder, "RAG_WORKER_PROVIDER_PROFILE", fmt.Errorf("real profile requires provider config"))
			return
		}
		providerSet, loadErr := ragproviders.Load(ctx, settings.ProviderConfig)
		if loadErr != nil {
			fail(encoder, "RAG_WORKER_PROVIDERS", loadErr)
			return
		}
		defer func() { _ = providerSet.Close() }()
		if err := providerSet.CheckExecution(execution); err != nil {
			fail(encoder, "RAG_WORKER_PROVIDER_REQUIREMENTS", err)
			return
		}
		options = providerSet.EngineOptions()
	default:
		fail(encoder, "RAG_WORKER_PROVIDER_PROFILE", fmt.Errorf("unsupported provider profile"))
		return
	}
	_, err = engine.Execute(ctx, execution, corpusArtifact.Corpus, evaluationArtifact.Dataset, observer{encoder: encoder}, options)
	if err != nil {
		fail(encoder, "RAG_WORKER_EXECUTE", err)
		return
	}
	_ = encoder.Encode(frame{Type: "complete", Complete: map[string]any{"status": "succeeded", "payload": map[string]any{"domain": ragcontract.Domain}}})
}
func decodeStrict(reader io.Reader, target any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("trailing JSON value")
		}
		return err
	}
	return nil
}
func fail(encoder *json.Encoder, code string, err error) {
	_ = encoder.Encode(frame{Type: "error", Error: map[string]any{"code": code, "message": err.Error(), "retryable": false}})
}
