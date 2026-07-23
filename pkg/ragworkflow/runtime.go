package ragworkflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/dop251/goja"
	gggengine "github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

type environmentFactory func(workflowv3runtime.TaskModuleContext, ragcontract.PipelineExecution) (*ragoperators.Environment, error)

func TaskModuleFactory() workflowv3runtime.TaskModuleFactory {
	return newTaskModuleFactory(ModuleAlias, nil, "fixture-embedding/v1", func(_ workflowv3runtime.TaskModuleContext, execution ragcontract.PipelineExecution) (*ragoperators.Environment, error) {
		return providerFreeEnvironment(execution)
	})
}

func newTaskModuleFactory(alias string, operations []workflowv3.ExternalOperationDescriptor, preparationIdentity string, factory environmentFactory) workflowv3runtime.TaskModuleFactory {
	return workflowv3runtime.TaskModuleFactory{
		Alias:      alias,
		Operations: operations,
		Build: func(moduleContext workflowv3runtime.TaskModuleContext) (gggengine.RuntimeModuleRegistrar, error) {
			if factory == nil || preparationIdentity == "" {
				return nil, fmt.Errorf("RAG_WORKFLOW_ENVIRONMENT_FACTORY")
			}
			runtime := &taskRuntime{context: moduleContext, environmentFactory: factory, preparationIdentity: preparationIdentity}
			loader := func(vm *goja.Runtime, moduleObject *goja.Object) {
				exports := moduleObject.Get("exports").ToObject(vm)
				for name, operation := range map[string]func() (any, error){"prepare": runtime.prepare, "query": runtime.query, "merge": runtime.merge, "publish": runtime.publish} {
					operation := operation
					if err := exports.Set(name, func(goja.FunctionCall) goja.Value {
						value, err := operation()
						if err != nil {
							panic(vm.NewGoError(err))
						}
						return vm.ToValue(value)
					}); err != nil {
						panic(vm.NewGoError(err))
					}
				}
			}
			return gggengine.NativeModuleRegistrar{ModuleID: "rag-workflow-runtime-v1", ModuleName: ModuleAlias, Loader: loader}, nil
		},
	}
}

type taskRuntime struct {
	context             workflowv3runtime.TaskModuleContext
	environmentFactory  environmentFactory
	preparationIdentity string
}

func (r *taskRuntime) prepare() (any, error) {
	execution, corpus, err := r.executionAndCorpus()
	if err != nil {
		return nil, err
	}
	executionDigest, _ := ragcontract.Digest(execution)
	pipelineDigest, _ := ragcontract.Digest(execution.Pipeline)
	corpusDigest, _ := ragcontract.Digest(corpus)
	fingerprint, err := preparationFingerprint(execution, r.preparationIdentity)
	if err != nil {
		return nil, err
	}
	values := map[string]any{"corpus/out": corpus}
	bundle := PreparedBundle{SchemaVersion: PreparedSchema, ExecutionDigest: executionDigest, PipelineDigest: pipelineDigest, CorpusDigest: corpusDigest, PreparationFingerprint: fingerprint}
	if previous, ok := r.context.Request.Inputs["prepared"]; ok {
		if err := r.decodeRef(previous, &bundle); err != nil {
			return nil, err
		}
		if err := validatePrepared(bundle, executionDigest, pipelineDigest, corpusDigest, fingerprint); err != nil {
			return nil, err
		}
		decoded, err := ragengine.DeserializePreparedValues(bundle.Values)
		if err != nil {
			return nil, err
		}
		for key, value := range decoded {
			values[key] = value
		}
	}
	if r.context.Request.NodeKey == "prepare-start" {
		return finalizePrepared(bundle, values)
	}
	ordinal, err := preparationOrdinal(r.context.Request.NodeKey)
	if err != nil {
		return nil, err
	}
	static := staticNodes(execution.Pipeline)
	if ordinal >= len(static) {
		return nil, fmt.Errorf("RAG_WORKFLOW_PREPARE_NODE")
	}
	node := static[ordinal]
	operator, ok := ragoperators.NativeRegistry().Lookup(node.Operator)
	if !ok {
		return nil, fmt.Errorf("RAG_WORKFLOW_OPERATOR_UNAVAILABLE")
	}
	inputs := map[string]any{}
	for _, binding := range node.Inputs {
		value, found := values[binding.From.NodeID+"/"+binding.From.Port]
		if !found {
			return nil, fmt.Errorf("RAG_WORKFLOW_PREPARED_INPUT: %s.%s", binding.From.NodeID, binding.From.Port)
		}
		inputs[binding.Port] = value
	}
	environment, err := r.environmentFactory(r.context, execution)
	if err != nil {
		return nil, err
	}
	outputs, err := operator.Execute(r.context.Context, node, inputs, environment)
	if err != nil {
		return nil, err
	}
	for port, value := range outputs {
		if port == "artifact" || port == "manifest" {
			continue
		}
		values[node.ID+"/"+port] = value
		if index, ok := value.(*ragoperators.MultiIndex); ok {
			recordDigest := sha256.Sum256(index.Artifact)
			bundle.Indexes = append(bundle.Indexes, IndexEvidence{NodeID: node.ID, Manifest: index.Manifest, RecordsDigest: "sha256:" + hex.EncodeToString(recordDigest[:]), RecordsSize: int64(len(index.Artifact))})
			_ = index.Close()
		}
	}
	return finalizePrepared(bundle, values)
}

func (r *taskRuntime) query() (any, error) {
	execution, corpus, err := r.executionAndCorpus()
	if err != nil {
		return nil, err
	}
	var item QueryItem
	if err := r.decodeInput("query", &item); err != nil || item.SchemaVersion != QuerySchema || item.DatasetManifestDigest != execution.Dataset.ManifestDigest || item.Query.ID == "" || item.Query.Text == "" {
		return nil, fmt.Errorf("RAG_WORKFLOW_QUERY_INPUT")
	}
	query := item.Query
	var bundle PreparedBundle
	if err := r.decodeInput("prepared", &bundle); err != nil {
		return nil, err
	}
	executionDigest, _ := ragcontract.Digest(execution)
	pipelineDigest, _ := ragcontract.Digest(execution.Pipeline)
	corpusDigest, _ := ragcontract.Digest(corpus)
	fingerprint, _ := preparationFingerprint(execution, r.preparationIdentity)
	if err := validatePrepared(bundle, executionDigest, pipelineDigest, corpusDigest, fingerprint); err != nil {
		return nil, err
	}
	values, err := ragengine.DeserializePreparedValues(bundle.Values)
	if err != nil {
		return nil, err
	}
	values["corpus/out"] = corpus
	environment, err := r.environmentFactory(r.context, execution)
	if err != nil {
		return nil, err
	}
	engine := ragengine.New(nil)
	if err := rebuildIndexes(r.context.Context, engine, execution.Pipeline, values, environment, bundle.Indexes); err != nil {
		return nil, err
	}
	prepared, err := ragengine.NewPreparedFromStaticValues(execution.Pipeline, values)
	if err != nil {
		closeValueIndexes(values)
		return nil, err
	}
	defer func() { _ = prepared.Close() }()
	result, err := engine.Execute(r.context.Context, execution, corpus, ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{query}}, nil, ragengine.Options{Prepared: prepared, Manifests: environment.Manifests, Embedder: environment.Embedder, EmbeddingFingerprint: "fixture-embedding/v1"})
	if err != nil {
		return nil, err
	}
	ragengine.SortMetrics(result.Metrics)
	queryResult := QueryResult{QueryID: query.ID, Metrics: workflowMetrics(result.Metrics), Answers: result.Answers, Failures: result.Failures}
	if len(result.Traces) != 1 {
		return nil, fmt.Errorf("RAG_WORKFLOW_QUERY_TRACE")
	}
	queryResult.Trace = result.Traces[0]
	for _, artifact := range result.Artifacts {
		queryResult.Artifacts = append(queryResult.Artifacts, Artifact{Role: artifact.Role, Kind: artifact.Kind, Name: artifact.Name, SchemaVersion: artifact.SchemaVersion, MediaType: artifact.MediaType, Data: artifact.Data, Metadata: artifact.Metadata})
	}
	partition := ResultPartition{SchemaVersion: ResultPartitionSchema, ExecutionDigest: executionDigest, PreparationFingerprint: fingerprint, Results: []QueryResult{queryResult}}
	partition.Digest, err = resultPartitionDigest(partition)
	if err != nil {
		return nil, err
	}
	return partition, nil
}

func (r *taskRuntime) merge() (any, error) {
	body, err := r.inputBody("partition")
	if err != nil {
		return nil, err
	}
	partition, err := workflowv3.DecodeReductionPartition(body, 16)
	if err != nil || partition.ItemSchema != ResultPartitionSchema {
		return nil, fmt.Errorf("RAG_WORKFLOW_REDUCTION_PARTITION")
	}
	var merged ResultPartition
	merged.SchemaVersion = ResultPartitionSchema
	for _, member := range partition.Members {
		var current ResultPartition
		if err := r.decodeRef(member.Value, &current); err != nil {
			return nil, err
		}
		if err := validateResultPartition(current); err != nil {
			return nil, err
		}
		if merged.ExecutionDigest == "" {
			merged.ExecutionDigest, merged.PreparationFingerprint = current.ExecutionDigest, current.PreparationFingerprint
		}
		if merged.ExecutionDigest != current.ExecutionDigest || merged.PreparationFingerprint != current.PreparationFingerprint {
			return nil, fmt.Errorf("RAG_WORKFLOW_REDUCTION_IDENTITY")
		}
		merged.Results = append(merged.Results, current.Results...)
	}
	sort.Slice(merged.Results, func(i, j int) bool { return merged.Results[i].QueryID < merged.Results[j].QueryID })
	for index := 1; index < len(merged.Results); index++ {
		if merged.Results[index].QueryID == merged.Results[index-1].QueryID {
			return nil, fmt.Errorf("RAG_WORKFLOW_QUERY_DUPLICATE")
		}
	}
	merged.Digest, err = resultPartitionDigest(merged)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func (r *taskRuntime) publish() (any, error) {
	execution, err := r.execution()
	if err != nil {
		return nil, err
	}
	var partition ResultPartition
	if err := r.decodeInput("results", &partition); err != nil {
		return nil, err
	}
	if err := validateResultPartition(partition); err != nil {
		return nil, err
	}
	executionDigest, _ := ragcontract.Digest(execution)
	if partition.ExecutionDigest != executionDigest {
		return nil, fmt.Errorf("RAG_WORKFLOW_RESULT_EXECUTION")
	}
	result := Result{SchemaVersion: ResultSchema, ExecutionDigest: executionDigest, CellID: execution.CellID, VariantID: execution.VariantID, Factors: append([]ragcontract.FactorSelection(nil), execution.Factors...), PreparationFingerprint: partition.PreparationFingerprint, Results: partition.Results}
	result.Digest, err = resultDigest(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *taskRuntime) executionAndCorpus() (ragcontract.PipelineExecution, ragoperators.Corpus, error) {
	execution, err := r.execution()
	if err != nil {
		return execution, ragoperators.Corpus{}, err
	}
	var corpus ragoperators.Corpus
	if err := r.decodeInput("corpus", &corpus); err != nil || corpus.SchemaVersion != "rag-corpus-data/v1" || len(corpus.Records) == 0 {
		return execution, corpus, fmt.Errorf("RAG_WORKFLOW_CORPUS_INPUT")
	}
	return execution, corpus, nil
}

func (r *taskRuntime) execution() (ragcontract.PipelineExecution, error) {
	body, err := r.inputBody("execution")
	if err != nil {
		return ragcontract.PipelineExecution{}, err
	}
	execution, err := ragcontract.DecodeExecution(bytes.NewReader(body))
	if err != nil {
		return ragcontract.PipelineExecution{}, err
	}
	if err := validateExecution(execution); err != nil {
		return ragcontract.PipelineExecution{}, err
	}
	return execution, nil
}

func (r *taskRuntime) inputBody(name string) ([]byte, error) {
	ref, ok := r.context.Request.Inputs[name]
	if !ok {
		return nil, fmt.Errorf("RAG_WORKFLOW_INPUT_MISSING: %s", name)
	}
	return workflowv3.ReadArtifact(r.context.Context, r.context.Request.Artifacts, ref)
}
func (r *taskRuntime) decodeInput(name string, target any) error {
	body, err := r.inputBody(name)
	if err != nil {
		return err
	}
	return strictJSON(body, target)
}
func (r *taskRuntime) decodeRef(ref workflowv3.ArtifactRef, target any) error {
	body, err := workflowv3.ReadArtifact(r.context.Context, r.context.Request.Artifacts, ref)
	if err != nil {
		return err
	}
	return strictJSON(body, target)
}

func preparationOrdinal(key workflowv3.NodeKey) (int, error) {
	value := strings.TrimPrefix(string(key), "prepare-")
	if value == string(key) || len(value) != 3 {
		return 0, fmt.Errorf("RAG_WORKFLOW_PREPARE_KEY")
	}
	ordinal, err := strconv.Atoi(value)
	if err != nil || ordinal < 0 {
		return 0, fmt.Errorf("RAG_WORKFLOW_PREPARE_KEY")
	}
	return ordinal, nil
}
func staticNodes(pipeline ragcontract.PipelineIR) []ragcontract.Node {
	ids := staticNodeIDs(pipeline)
	ret := []ragcontract.Node{}
	for _, node := range pipeline.Nodes {
		if ids[node.ID] {
			ret = append(ret, node)
		}
	}
	return ret
}

func finalizePrepared(bundle PreparedBundle, values map[string]any) (PreparedBundle, error) {
	serialized, err := ragengine.SerializePreparedValues(values)
	if err != nil {
		return PreparedBundle{}, err
	}
	bundle.Values = serialized
	sort.Slice(bundle.Indexes, func(i, j int) bool { return bundle.Indexes[i].NodeID < bundle.Indexes[j].NodeID })
	bundle.Digest, err = preparedDigest(bundle)
	return bundle, err
}
func validatePrepared(bundle PreparedBundle, executionDigest, pipelineDigest, corpusDigest, fingerprint string) error {
	want, err := preparedDigest(bundle)
	if err != nil || bundle.SchemaVersion != PreparedSchema || bundle.Digest != want || bundle.ExecutionDigest != executionDigest || bundle.PipelineDigest != pipelineDigest || bundle.CorpusDigest != corpusDigest || bundle.PreparationFingerprint != fingerprint {
		return fmt.Errorf("RAG_WORKFLOW_PREPARED_IDENTITY")
	}
	return nil
}
func preparedDigest(bundle PreparedBundle) (string, error) {
	bundle.Digest = ""
	return ragcontract.Digest(bundle)
}
func resultPartitionDigest(value ResultPartition) (string, error) {
	value.Digest = ""
	return ragcontract.Digest(value)
}
func validateResultPartition(value ResultPartition) error {
	want, err := resultPartitionDigest(value)
	if err != nil || value.SchemaVersion != ResultPartitionSchema || value.Digest != want || value.ExecutionDigest == "" || value.PreparationFingerprint == "" || len(value.Results) == 0 {
		return fmt.Errorf("RAG_WORKFLOW_RESULT_PARTITION")
	}
	return nil
}
func workflowMetrics(values []ragoperators.Metric) []Metric {
	ret := make([]Metric, len(values))
	for index, value := range values {
		ret[index] = Metric{Name: value.Name, Unit: value.Unit, Value: value.Value, Numeric: value.Numeric, Metadata: value.Metadata}
	}
	return ret
}

func resultDigest(value Result) (string, error) { value.Digest = ""; return ragcontract.Digest(value) }

func preparationFingerprint(execution ragcontract.PipelineExecution, implementationIdentity string) (string, error) {
	pipelineDigest, _ := ragcontract.Digest(execution.Pipeline)
	binding, err := corpusBinding(execution)
	if err != nil {
		return "", err
	}
	return ragcontract.Digest(struct {
		SchemaVersion           string `json:"schemaVersion"`
		PipelineDigest          string `json:"pipelineDigest"`
		CorpusDigest            string `json:"corpusDigest"`
		EmbeddingImplementation string `json:"embeddingImplementation"`
	}{"rag-workflow-preparation-fingerprint/v1", pipelineDigest, binding.Digest, implementationIdentity})
}

func rebuildIndexes(ctx context.Context, engine *ragengine.Engine, pipeline ragcontract.PipelineIR, values map[string]any, env *ragoperators.Environment, evidence []IndexEvidence) error {
	byNode := map[string]IndexEvidence{}
	for _, item := range evidence {
		byNode[item.NodeID] = item
	}
	for _, node := range staticNodes(pipeline) {
		if !strings.HasPrefix(node.Operator.Kind, "index.") {
			continue
		}
		operator, ok := engine.Registry.Lookup(node.Operator)
		if !ok {
			return fmt.Errorf("RAG_WORKFLOW_INDEX_OPERATOR")
		}
		inputs := map[string]any{}
		for _, binding := range node.Inputs {
			value, found := values[binding.From.NodeID+"/"+binding.From.Port]
			if !found {
				return fmt.Errorf("RAG_WORKFLOW_INDEX_INPUT")
			}
			inputs[binding.Port] = value
		}
		outputs, err := operator.Execute(ctx, node, inputs, env)
		if err != nil {
			return err
		}
		index, ok := outputs["index"].(*ragoperators.MultiIndex)
		if !ok {
			return fmt.Errorf("RAG_WORKFLOW_INDEX_OUTPUT")
		}
		expected, ok := byNode[node.ID]
		recordsHash := sha256.Sum256(index.Artifact)
		recordsDigest := "sha256:" + hex.EncodeToString(recordsHash[:])
		if !ok || expected.Manifest.Digest != index.Manifest.Digest || expected.RecordsDigest != recordsDigest || expected.RecordsSize != int64(len(index.Artifact)) {
			_ = index.Close()
			return fmt.Errorf("RAG_WORKFLOW_INDEX_IDENTITY")
		}
		values[node.ID+"/index"] = index
	}
	return nil
}
func closeValueIndexes(values map[string]any) {
	for _, value := range values {
		if index, ok := value.(*ragoperators.MultiIndex); ok {
			_ = index.Close()
		}
	}
}

func providerFreeEnvironment(execution ragcontract.PipelineExecution) (*ragoperators.Environment, error) {
	dimensions := 0
	for _, node := range execution.Pipeline.Nodes {
		if node.Operator.Kind == "embed.model" {
			var config embeddingConfig
			if err := strictJSON(node.Config, &config); err != nil {
				return nil, err
			}
			dimensions = config.Dimensions
		}
	}
	if dimensions == 0 {
		dimensions = 8
	}
	digest, _ := ragcontract.Digest("fixture-embedding-v1")
	manifest := ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: digest}, ProviderAdapterVersion: "provider-free/v1", ModelID: "fixture-embedding-v1", ModelDigest: digest, Dimensions: dimensions, Tokenization: "unicode-codepoint", Truncation: "none", Normalization: "l2", ImplementationVersion: "fixture-embedding/v1", RequestParameters: json.RawMessage(`{}`)}
	return &ragoperators.Environment{Manifests: ragoperators.StaticManifestResolver{Models: map[string]ragcontract.ModelManifest{"fixture-embedding-v1": manifest}}, Embedder: deterministicEmbedder{dimensions: dimensions}, Usage: ragoperators.Usage{Cost: map[string]float64{}}}, nil
}

type deterministicEmbedder struct{ dimensions int }

func (e deterministicEmbedder) Embed(ctx context.Context, _ string, texts []string) ([][]float64, ragoperators.Usage, error) {
	vectors := make([][]float64, len(texts))
	for i, text := range texts {
		if err := ctx.Err(); err != nil {
			return nil, ragoperators.Usage{}, err
		}
		vector := make([]float64, e.dimensions)
		for d := range vector {
			sum := sha256.Sum256([]byte(fmt.Sprintf("fixture-embedding/v1\x00%d\x00%s", d, text)))
			raw := float64(int(sum[0])<<8 | int(sum[1]))
			vector[d] = raw/32767.5 - 1
		}
		norm := 0.0
		for _, v := range vector {
			norm += v * v
		}
		if norm == 0 {
			vector[0] = 1
		} else {
			norm = math.Sqrt(norm)
			for d := range vector {
				vector[d] /= norm
			}
		}
		vectors[i] = vector
	}
	return vectors, ragoperators.Usage{EmbeddingTokens: int64(len(texts))}, nil
}
