package workflowv3ttc

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"github.com/go-go-golems/scraper/pkg/workflowv3sqlite"
)

type RuntimeConfig struct {
	StateDB                 string
	ArtifactRoot            string
	MaxArtifactBytes        int64
	Execution               ragcontract.PipelineExecution
	Corpus                  ragoperators.Corpus
	Dataset                 ragoperators.EvaluationDataset
	Options                 ragengine.Options
	PreparedStore           ragengine.PreparedCorpusStore
	PreparedIdentity        ragengine.PreparedCorpusIdentity
	CorpusDigest            string
	DatasetDigest           string
	ProviderProfileDigest   string
	GenerationModelDigest   string
	EmbeddingProfileDigest  string
	PublicationPolicyDigest string
}

type Runtime struct {
	Store     *workflowv3sqlite.Store
	Artifacts *workflowv3.FileArtifactStore
	Registry  *workflowv3.SealedRegistry
	Engine    *workflowv3runtime.Engine
	Plan      workflowv3.WorkflowPlan
	Inputs    MaterializedInputs
}

func NewRuntime(ctx context.Context, config RuntimeConfig) (*Runtime, error) {
	if config.StateDB == "" || config.ArtifactRoot == "" || config.MaxArtifactBytes < 1 || config.PreparedStore == nil ||
		!validDigest(config.CorpusDigest) || !validDigest(config.DatasetDigest) || !validDigest(config.PublicationPolicyDigest) {
		return nil, fmt.Errorf("complete Workflow V3 runtime configuration is required")
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(config.Execution.Pipeline)
	if err != nil {
		return nil, err
	}
	provider, err := NewOperatorProvider(OperatorProviderConfig{
		GenerationNode: mapping.CombinedNode, EmbeddingNode: mapping.EmbeddingNode,
		RawRepresentationName: mapping.RawRepresentationName, MaxRepresentationsPerChunk: mapping.MaxRepresentationsPerChunk,
		ProviderProfileDigest: config.ProviderProfileDigest, GenerationModelDigest: config.GenerationModelDigest,
		EmbeddingProfileDigest: config.EmbeddingProfileDigest,
		ResolveEnvironment: func(context.Context) (*ragoperators.Environment, error) {
			return &ragoperators.Environment{Manifests: config.Options.Manifests, Schemas: config.Options.Schemas, Generator: config.Options.Generator, Embedder: config.Options.Embedder, Reranker: config.Options.Reranker, Cache: config.Options.Cache, Usage: ragoperators.Usage{Cost: map[string]float64{}}, GenerationConcurrency: 1, GenerationSettingsFingerprint: config.Options.GenerationSettingsFingerprint}, nil
		},
	})
	if err != nil {
		return nil, err
	}
	domainEngine := ragengine.New(nil)
	publication, err := NewPublicationAuthority(PublicationConfig{Store: config.PreparedStore, Engine: domainEngine, Pipeline: config.Execution.Pipeline, Corpus: config.Corpus, Options: config.Options, Identity: config.PreparedIdentity, ChunksOutputKey: mapping.ChunksOutputKey, RawOutputKey: mapping.RawOutputKey, DerivedOutputKey: mapping.DerivedOutputKey, MergedOutputKey: mapping.MergedOutputKey, EmbeddingOutputKey: mapping.EmbeddingOutputKey, PolicyDigest: config.PublicationPolicyDigest})
	if err != nil {
		return nil, err
	}
	artifacts, err := workflowv3.NewFileArtifactStore(filepath.Clean(config.ArtifactRoot), config.MaxArtifactBytes)
	if err != nil {
		return nil, err
	}
	inputs, err := MaterializeInputs(ctx, artifacts, domainEngine, config.Execution, config.Corpus, config.Dataset, config.Options, config.CorpusDigest, config.DatasetDigest)
	if err != nil {
		return nil, err
	}
	evaluation, err := NewEvaluationAuthority(EvaluationConfig{Store: config.PreparedStore, Engine: domainEngine, Execution: config.Execution, Corpus: config.Corpus, Options: config.Options, Identity: config.PreparedIdentity, DatasetDigest: config.DatasetDigest, ValidCitationIDs: inputs.ValidCitationIDs})
	if err != nil {
		return nil, err
	}
	registry, err := Registry()
	if err != nil {
		return nil, err
	}
	catalog, err := registry.Catalog()
	if err != nil {
		return nil, err
	}
	authored, err := workflowmodule.Author(ctx, ProductionWorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		return nil, err
	}
	modules, err := workflowv3runtime.NewTaskModuleRegistry(ModuleWithAuthorities(provider, publication, evaluation))
	if err != nil {
		return nil, err
	}
	store, err := workflowv3sqlite.Open(ctx, filepath.Clean(config.StateDB))
	if err != nil {
		return nil, err
	}
	runtime := &Runtime{Store: store, Artifacts: artifacts, Registry: registry, Plan: authored.Plan, Inputs: inputs}
	runtime.Engine = &workflowv3runtime.Engine{Store: store, Registry: registry, Artifacts: artifacts, Modules: modules, LeaseDuration: time.Minute}
	return runtime, nil
}

func (r *Runtime) Submit(ctx context.Context, runID workflowv3.RunID) error {
	if r == nil || r.Engine == nil {
		return fmt.Errorf("workflow V3 runtime is required")
	}
	return r.Engine.Submit(ctx, runID, r.Plan, map[string]workflowv3.ArtifactRef{"chunks": r.Inputs.Chunks, "queries": r.Inputs.Queries})
}

func (r *Runtime) Dispatcher(capacities map[string]int) *workflowv3runtime.Dispatcher {
	return &workflowv3runtime.Dispatcher{Engine: r.Engine, Capacities: capacities, PollInterval: 10 * time.Millisecond}
}

func (r *Runtime) Close() error {
	if r == nil || r.Store == nil {
		return nil
	}
	return r.Store.Close()
}
