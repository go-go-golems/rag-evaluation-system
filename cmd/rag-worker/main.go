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
	ProviderProfile    string `glazed:"provider-profile"`
	ProviderConfig     string `glazed:"provider-config"`
	Capabilities       bool   `glazed:"capabilities"`
	PreparationStateDB string `glazed:"preparation-state-db"`
}

var _ cmds.BareCommand = (*workerCommand)(nil)

func newWorkerCommand() (*cobra.Command, error) {
	description := cmds.NewCommandDescription(
		"rag-worker",
		cmds.WithShort("Execute a canonical RAG request over the stdio runner protocol"),
		cmds.WithFlags(
			fields.New("provider-profile", fields.TypeString, fields.WithHelp("Explicit provider profile: fixtures or real"), fields.WithRequired(true)),
			fields.New("provider-config", fields.TypeString, fields.WithHelp("Host-only real-provider configuration YAML; required for provider-profile=real")),
			fields.New("capabilities", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Print the non-secret provider capability descriptor and exit")),
			fields.New("preparation-state-db", fields.TypeString, fields.WithHelp("Opt-in scraper SQLite state database for durable canonical combined preparation")),
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
	if settings.Capabilities {
		return printCapabilities(ctx, settings)
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
	ctx, cancel := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	options, capabilities, checkExecution, closeProviders, err := loadWorkerProviders(ctx, settings)
	if err != nil {
		fail(encoder, "RAG_WORKER_PROVIDERS", err)
		return
	}
	defer closeProviders()
	_ = capabilities
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
	if err := checkExecution(execution); err != nil {
		fail(encoder, "RAG_WORKER_PROVIDER_REQUIREMENTS", err)
		return
	}
	engine := ragengine.New(nil)
	options.PreparationEvent = observer{encoder: encoder}.Event
	options.PreparedCorpusDigest = corpusArtifact.Manifest.Digest
	if options.PreparedStore != nil {
		pipelineDigest, err := ragcontract.Digest(execution.Pipeline)
		if err != nil {
			fail(encoder, "RAG_WORKER_PREPARED", err)
			return
		}
		identity := ragengine.PreparedCorpusIdentity{
			SchemaVersion: "rag-prepared-corpus-manifest/v1", CorpusDigest: corpusArtifact.Manifest.Digest,
			PipelineDigest: pipelineDigest, GenerationSettingsFingerprint: options.GenerationSettingsFingerprint,
			EmbeddingFingerprint: options.EmbeddingFingerprint,
		}
		preparedDigest, err := ragcontract.Digest(identity)
		if err != nil {
			fail(encoder, "RAG_WORKER_PREPARED", err)
			return
		}
		if settings.PreparationStateDB != "" {
			if err := executeDurablePreparation(ctx, settings.PreparationStateDB, execution, corpusArtifact.Corpus, options, identity); err != nil {
				fail(encoder, "RAG_WORKER_PREPARATION", err)
				return
			}
		}
		prepared, found, err := options.PreparedStore.Open(ctx, engine, execution.Pipeline, corpusArtifact.Corpus, options, identity)
		if err != nil {
			fail(encoder, "RAG_WORKER_PREPARED", err)
			return
		}
		if !found {
			if settings.PreparationStateDB != "" {
				fail(encoder, "RAG_WORKER_PREPARED", fmt.Errorf("durable preparation completed without a published corpus"))
				return
			}
			prepared, err = engine.Prepare(ctx, execution.Pipeline, corpusArtifact.Corpus, options)
			if err != nil {
				fail(encoder, "RAG_WORKER_PREPARED", err)
				return
			}
			if _, err := options.PreparedStore.Put(ctx, prepared, identity); err != nil {
				_ = prepared.Close()
				fail(encoder, "RAG_WORKER_PREPARED", err)
				return
			}
		}
		defer func() { _ = prepared.Close() }()
		options.Prepared = prepared
		options.PreparedCorpusDigest = preparedDigest
	}
	_, err = engine.Execute(ctx, execution, corpusArtifact.Corpus, evaluationArtifact.Dataset, observer{encoder: encoder}, options)
	if err != nil {
		fail(encoder, "RAG_WORKER_EXECUTE", err)
		return
	}
	_ = encoder.Encode(frame{Type: "complete", Complete: map[string]any{"status": "succeeded", "payload": map[string]any{"domain": ragcontract.Domain}}})
}
func printCapabilities(ctx context.Context, settings *workerSettings) error {
	_, capabilities, _, closeProviders, err := loadWorkerProviders(ctx, settings)
	if err != nil {
		return err
	}
	defer closeProviders()
	return json.NewEncoder(os.Stdout).Encode(capabilities)
}

func loadWorkerProviders(ctx context.Context, settings *workerSettings) (ragengine.Options, any, func(ragcontract.PipelineExecution) error, func(), error) {
	if settings == nil {
		return ragengine.Options{}, nil, nil, nil, fmt.Errorf("worker settings required")
	}
	switch settings.ProviderProfile {
	case "fixtures":
		if settings.ProviderConfig != "" {
			return ragengine.Options{}, nil, nil, nil, fmt.Errorf("fixture profile does not accept provider config")
		}
		fixtures := ragoperators.NewFixtureProviders()
		capabilities := map[string]any{"schemaVersion": "rag-provider-capabilities/v1", "profileId": "fixtures", "fixtureProviders": true, "capabilities": []string{"embedder", "generator", "schema-validator"}}
		return ragengine.Options{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache()}, capabilities, func(ragcontract.PipelineExecution) error { return nil }, func() {}, nil
	case "real":
		if settings.ProviderConfig == "" {
			return ragengine.Options{}, nil, nil, nil, fmt.Errorf("real profile requires provider config")
		}
		set, err := ragproviders.Load(ctx, settings.ProviderConfig)
		if err != nil {
			return ragengine.Options{}, nil, nil, nil, err
		}
		return set.EngineOptions(), set.CapabilityDescriptor(), set.CheckExecution, func() { _ = set.Close() }, nil
	default:
		return ragengine.Options{}, nil, nil, nil, fmt.Errorf("unsupported provider profile")
	}
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
