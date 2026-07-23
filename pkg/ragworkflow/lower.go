package ragworkflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const LoweredExecutionSchema = "rag-workflow-lowered-execution/v1"

type embeddingConfig struct {
	Model         string          `json:"model"`
	Dimensions    int             `json:"dimensions"`
	Distance      string          `json:"distance"`
	Normalize     string          `json:"normalize"`
	BatchSize     int             `json:"batchSize"`
	Preprocessing json.RawMessage `json:"preprocessing,omitempty"`
}

type Lowerer struct {
	definitions         *ragcompiler.Registry
	lowerings           *OperatorRegistry
	bundle              func() (*workflowv3.Bundle, error)
	allowProviders      bool
	preparationIdentity string
}

func NewLowerer() *Lowerer {
	return &Lowerer{definitions: ragcompiler.BuiltinRegistry(), lowerings: NewOperatorRegistry(), bundle: Bundle, preparationIdentity: "fixture-embedding/v1"}
}

func NewProviderLowerer(providerPackage *ProviderPackage) (*Lowerer, error) {
	if providerPackage == nil {
		return nil, fmt.Errorf("RAG_WORKFLOW_PROVIDER_PACKAGE")
	}
	if _, err := providerPackage.Bundle(); err != nil {
		return nil, err
	}
	return &Lowerer{definitions: ragcompiler.BuiltinRegistry(), lowerings: NewOperatorRegistry(), bundle: providerPackage.Bundle, allowProviders: true, preparationIdentity: providerPackage.authority.Digest}, nil
}

func (l *Lowerer) Lower(_ context.Context, execution ragcontract.PipelineExecution) (LoweredExecution, error) {
	if err := validateExecution(execution); err != nil {
		return LoweredExecution{}, err
	}
	if l == nil || l.definitions == nil || l.lowerings == nil || l.bundle == nil || l.preparationIdentity == "" {
		return LoweredExecution{}, fmt.Errorf("RAG_WORKFLOW_REGISTRY")
	}
	static := staticNodeIDs(execution.Pipeline)
	for _, node := range execution.Pipeline.Nodes {
		if _, known := l.definitions.Definition(node.Operator); !known {
			return LoweredExecution{}, fmt.Errorf("RAG_WORKFLOW_OPERATOR_UNKNOWN: %s", node.Operator.ID())
		}
		lowering, known := l.lowerings.Definition(node.Operator)
		if !known {
			return LoweredExecution{}, fmt.Errorf("RAG_WORKFLOW_OPERATOR_UNKNOWN: %s", node.Operator.ID())
		}
		if err := validateSupportedNode(node, static[node.ID], lowering, l.allowProviders); err != nil {
			return LoweredExecution{}, err
		}
	}
	bundle, err := l.bundle()
	if err != nil {
		return LoweredExecution{}, err
	}
	builder := workflowv3.NewRegistryBuilder()
	if err := builder.AdvertiseModules(ModuleAlias); err != nil {
		return LoweredExecution{}, err
	}
	if err := builder.AddBundle(bundle); err != nil {
		return LoweredExecution{}, err
	}
	registry, err := builder.Seal()
	if err != nil {
		return LoweredExecution{}, err
	}
	catalog, err := registry.Catalog()
	if err != nil {
		return LoweredExecution{}, err
	}

	executionDigest, _ := ragcontract.Digest(execution)
	pipelineDigest, _ := ragcontract.Digest(execution.Pipeline)
	corpusBinding, err := corpusBinding(execution)
	if err != nil {
		return LoweredExecution{}, err
	}
	preparationFingerprint, err := ragcontract.Digest(struct {
		SchemaVersion           string `json:"schemaVersion"`
		PipelineDigest          string `json:"pipelineDigest"`
		CorpusDigest            string `json:"corpusDigest"`
		EmbeddingImplementation string `json:"embeddingImplementation"`
	}{"rag-workflow-preparation-fingerprint/v1", pipelineDigest, corpusBinding.Digest, l.preparationIdentity})
	if err != nil {
		return LoweredExecution{}, err
	}

	ir := workflowv3.WorkflowIR{
		Schema: workflowv3.IRSchema, Name: "rag-v2-" + executionDigest[7:19],
		Inputs:    []workflowv3.IRInput{{Name: "execution", Schema: ragcontract.ExecutionSchemaVersion}, {Name: "corpus", Schema: CorpusSchema}},
		SetInputs: []workflowv3.IRSetInput{{Name: "queries", ItemSchema: QuerySchema, ManifestSchema: workflowv3.ItemManifestSchemaV1}},
	}
	executionRef := workflowv3.ValueRef{Source: "input", Name: "execution", Schema: ragcontract.ExecutionSchemaVersion}
	corpusRef := workflowv3.ValueRef{Source: "input", Name: "corpus", Schema: CorpusSchema}
	previous := workflowv3.NodeKey("prepare-start")
	ir.Nodes = append(ir.Nodes, workflowv3.IRNode{
		Key: previous, Task: TaskCorpusLoad,
		Bindings: map[string]workflowv3.ValueRef{"execution": executionRef, "corpus": corpusRef},
	})
	staticOrdinal := 0
	for _, node := range execution.Pipeline.Nodes {
		if !static[node.ID] {
			continue
		}
		key := workflowv3.NodeKey(fmt.Sprintf("prepare-%03d", staticOrdinal))
		lowering, _ := l.lowerings.Definition(node.Operator)
		ir.Nodes = append(ir.Nodes, workflowv3.IRNode{
			Key: key, Task: lowering.Task, DependsOn: []workflowv3.NodeKey{previous},
			Bindings: map[string]workflowv3.ValueRef{
				"execution": executionRef, "corpus": corpusRef,
				"prepared": {Source: "node-output", NodeKey: previous, Port: "prepared", Schema: PreparedSchema},
			},
		})
		previous = key
		staticOrdinal++
	}
	mapKey := "evaluate-queries"
	ir.Maps = []workflowv3.IRMap{{
		Key:      mapKey,
		Source:   workflowv3.SetRef{Source: "set-input", Name: "queries", ItemSchema: QuerySchema, ManifestSchema: workflowv3.ItemManifestSchemaV1},
		ItemTask: TaskQuery,
		Bindings: map[string]workflowv3.ValueRef{
			"execution": executionRef, "corpus": corpusRef,
			"prepared": {Source: "node-output", NodeKey: previous, Port: "prepared", Schema: PreparedSchema},
			"query":    {Source: "map-item", MapKey: mapKey, Schema: QuerySchema},
		},
		Policy: workflowv3.MapPolicy{PageSize: 8, MaxItems: 10_000, MaxMaterializedAhead: boundedMaterialization(execution)},
	}}
	reduceKey := "merge-results"
	ir.Reductions = []workflowv3.IRReduce{{
		Key:           reduceKey,
		Source:        workflowv3.SetRef{Source: "map-output", MapKey: mapKey, ItemSchema: ResultPartitionSchema, ManifestSchema: workflowv3.ItemManifestSchemaV1},
		PartitionTask: TaskMerge,
		Bindings:      map[string]workflowv3.ValueRef{"partition": {Source: "reduction-partition", ReduceKey: reduceKey, Schema: workflowv3.ReductionPartitionSchemaV1}},
		Policy:        workflowv3.ReducePolicy{FanIn: 16, MaxLevels: 4},
	}}
	publishKey := workflowv3.NodeKey("publish-results")
	ir.Nodes = append(ir.Nodes, workflowv3.IRNode{
		Key: publishKey, Task: TaskPublish,
		Bindings: map[string]workflowv3.ValueRef{
			"execution": executionRef,
			"results":   {Source: "reduction-output", ReduceKey: reduceKey, Schema: ResultPartitionSchema},
		},
	})
	ir.Outputs = []workflowv3.IROutput{{Name: "result", Value: workflowv3.ValueRef{Source: "node-output", NodeKey: publishKey, Port: "result", Schema: ResultSchema}}}
	plan, err := workflowv3.Compile(ir, catalog)
	if err != nil {
		return LoweredExecution{}, fmt.Errorf("RAG_WORKFLOW_COMPILE: %w", err)
	}
	return LoweredExecution{SchemaVersion: LoweredExecutionSchema, ExecutionDigest: executionDigest, PreparationFingerprint: preparationFingerprint, IR: ir, Plan: plan}, nil
}

func validateExecution(execution ragcontract.PipelineExecution) error {
	if execution.SchemaVersion != ragcontract.ExecutionSchemaVersion || execution.Pipeline.SchemaVersion != ragcontract.PipelineSchemaVersion {
		return fmt.Errorf("RAG_WORKFLOW_EXECUTION_SCHEMA")
	}
	normalized, err := ragcompiler.Normalize(execution.Pipeline, nil)
	if err != nil {
		return fmt.Errorf("RAG_WORKFLOW_PIPELINE: %w", err)
	}
	got, _ := ragcontract.CanonicalJSON(execution.Pipeline)
	want, _ := ragcontract.CanonicalJSON(normalized)
	if !bytes.Equal(got, want) {
		return fmt.Errorf("RAG_WORKFLOW_PIPELINE_NONCANONICAL")
	}
	identity := execution
	identity.CellID = ""
	cellID, err := ragcontract.Digest(identity)
	if err != nil || execution.CellID != cellID {
		return fmt.Errorf("RAG_WORKFLOW_CELL_ID")
	}
	return nil
}

func corpusBinding(execution ragcontract.PipelineExecution) (ragcontract.ArtifactBinding, error) {
	for _, binding := range execution.Bindings {
		if binding.SlotID == "corpus" {
			if binding.Digest == "" || binding.SchemaVersion != ragcontract.CorpusManifestSchema {
				return ragcontract.ArtifactBinding{}, fmt.Errorf("RAG_WORKFLOW_CORPUS_BINDING")
			}
			return binding, nil
		}
	}
	return ragcontract.ArtifactBinding{}, fmt.Errorf("RAG_WORKFLOW_CORPUS_BINDING")
}

func validateSupportedNode(node ragcontract.Node, static bool, lowering OperatorLowering, allowProviders bool) error {
	if lowering.ProviderRequired && !allowProviders {
		return fmt.Errorf("RAG_WORKFLOW_PROVIDER_REQUIRED: %s", node.Operator.ID())
	}
	if node.Operator.Kind == "embed.model" {
		var config embeddingConfig
		if err := strictJSON(node.Config, &config); err != nil || config.Dimensions < 1 || config.Dimensions > 65_536 {
			return fmt.Errorf("RAG_WORKFLOW_EMBEDDING: %s", node.ID)
		}
		if !allowProviders && config.Model != "fixture-embedding-v1" {
			return fmt.Errorf("RAG_WORKFLOW_FIXTURE_EMBEDDING: %s", node.ID)
		}
	}
	if static && lowering.Phase != PhasePreparation {
		return fmt.Errorf("RAG_WORKFLOW_STATIC_OPERATOR: %s", node.Operator.ID())
	}
	if !static && lowering.Phase != PhaseQuery {
		return fmt.Errorf("RAG_WORKFLOW_DYNAMIC_OPERATOR: %s", node.Operator.ID())
	}
	return nil
}

func staticNodeIDs(pipeline ragcontract.PipelineIR) map[string]bool {
	dynamic, static := map[string]bool{"query": true}, map[string]bool{}
	for _, node := range pipeline.Nodes {
		isDynamic := false
		for _, input := range node.Inputs {
			if dynamic[input.From.NodeID] {
				isDynamic = true
				break
			}
		}
		if isDynamic {
			dynamic[node.ID] = true
		} else {
			static[node.ID] = true
		}
	}
	return static
}

func boundedMaterialization(_ ragcontract.PipelineExecution) int {
	// Runtime concurrency is not currently part of PipelineExecution. Keep the
	// plan bound deterministic and let worker capacity provide the lower limit.
	return 8
}

func strictJSON(raw json.RawMessage, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("trailing JSON")
	}
	return nil
}
