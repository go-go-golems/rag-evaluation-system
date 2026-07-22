package workflowv3ttc

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	gggengine "github.com/go-go-golems/go-go-goja/pkg/engine"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
)

func Module(provider Provider) workflowv3runtime.TaskModuleFactory {
	return ModuleWithAuthorities(provider, nil, nil)
}

func ModuleWithPublication(provider Provider, publication PublicationService) workflowv3runtime.TaskModuleFactory {
	return ModuleWithAuthorities(provider, publication, nil)
}

func ModuleWithAuthorities(provider Provider, publication PublicationService, evaluation EvaluationService) workflowv3runtime.TaskModuleFactory {
	return workflowv3runtime.TaskModuleFactory{
		Alias: ModuleAlias,
		Validate: func() error {
			if provider == nil {
				return fmt.Errorf("TTC provider is required")
			}
			return nil
		},
		Build: func(moduleContext workflowv3runtime.TaskModuleContext) (gggengine.RuntimeModuleRegistrar, error) {
			loader := func(vm *goja.Runtime, moduleObject *goja.Object) {
				exports := moduleObject.Get("exports").ToObject(vm)
				mustSet := func(name string, function func(goja.FunctionCall) goja.Value) {
					if err := exports.Set(name, function); err != nil {
						panic(vm.NewGoError(err))
					}
				}
				mustSet("generate", func(goja.FunctionCall) goja.Value {
					var chunk Chunk
					if err := readInput(moduleContext, "chunk", &chunk); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_CHUNK_INPUT", false))
					}
					if err := validateChunk(chunk); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_CHUNK_INVALID", false))
					}
					result, err := provider.Generate(moduleContext.Context, chunk)
					if err != nil {
						return vm.ToValue(providerFailure(err))
					}
					if err := validateGenerated(chunk, result.Value); err != nil {
						return vm.ToValue(closedFailureWithUsage("malformed-output", "RAG_TTC_GENERATED_INVALID", true, result.Usage))
					}
					return vm.ToValue(successResponse(result.Value, result.Usage))
				})
				mustSet("embed", func(goja.FunctionCall) goja.Value {
					var generated Generated
					if err := readInput(moduleContext, "generated", &generated); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_GENERATED_INPUT", false))
					}
					result, err := provider.Embed(moduleContext.Context, generated)
					if err != nil {
						return vm.ToValue(providerFailure(err))
					}
					if err := validateEmbedded(generated, result.Value); err != nil {
						return vm.ToValue(closedFailureWithUsage("malformed-output", "RAG_TTC_EMBEDDED_INVALID", true, result.Usage))
					}
					shard, err := mergeShards([]Embedded{result.Value})
					if err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_SHARD_INVALID", false))
					}
					return vm.ToValue(successResponse(shard, result.Usage))
				})
				mustSet("merge", func(goja.FunctionCall) goja.Value {
					items, err := readPartition(moduleContext)
					if err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_PARTITION_INVALID", false))
					}
					shard, err := mergeShards(items)
					if err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_SHARD_INVALID", false))
					}
					return vm.ToValue(successResponse(shard, nil))
				})
				mustSet("validatePublication", func(goja.FunctionCall) goja.Value {
					if publication == nil {
						return vm.ToValue(closedFailure("configuration", "RAG_TTC_PUBLICATION_UNAVAILABLE", false))
					}
					var shard PreparedShard
					if err := readInput(moduleContext, "shard", &shard); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_SHARD_INPUT", false))
					}
					receipt, err := publication.Validate(shard)
					if err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_PUBLICATION_INVALID", false))
					}
					return vm.ToValue(successResponse(receipt, nil))
				})
				mustSet("publish", func(goja.FunctionCall) goja.Value {
					if publication == nil {
						return vm.ToValue(closedFailure("configuration", "RAG_TTC_PUBLICATION_UNAVAILABLE", false))
					}
					var shard PreparedShard
					var decision PublicationDecision
					if err := readInput(moduleContext, "shard", &shard); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_SHARD_INPUT", false))
					}
					if err := readInput(moduleContext, "decision", &decision); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_PUBLICATION_DECISION", false))
					}
					receipt, err := publication.Publish(moduleContext.Context, shard, decision)
					if err != nil {
						return vm.ToValue(closedFailure("publication", "RAG_TTC_PUBLICATION_FAILED", false))
					}
					return vm.ToValue(successResponse(receipt, nil))
				})
				mustSet("evaluate", func(goja.FunctionCall) goja.Value {
					if evaluation == nil {
						return vm.ToValue(closedFailure("configuration", "RAG_TTC_EVALUATION_UNAVAILABLE", false))
					}
					var publicationReceipt PublicationReceipt
					var query QueryEnvelope
					if err := readInput(moduleContext, "publication", &publicationReceipt); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_PUBLICATION_RECEIPT", false))
					}
					if err := readInput(moduleContext, "query", &query); err != nil {
						return vm.ToValue(closedFailure("validation", "RAG_TTC_QUERY_INPUT", false))
					}
					evidence, err := evaluation.Evaluate(moduleContext.Context, publicationReceipt, query)
					if err != nil {
						return vm.ToValue(closedFailure("evaluation", "RAG_TTC_QUERY_EVALUATION", true))
					}
					return vm.ToValue(successResponse(evidence, evidence.Usage))
				})
			}
			return gggengine.NativeModuleRegistrar{ModuleID: "rag-evaluation-system:workflowv3-ttc", ModuleName: ModuleAlias, Loader: loader}, nil
		},
	}
}

func readInput(context workflowv3runtime.TaskModuleContext, port string, target any) error {
	ref, ok := context.Request.Inputs[port]
	if !ok {
		return fmt.Errorf("input is missing")
	}
	body, err := workflowv3.ReadArtifact(context.Context, context.Request.Artifacts, ref)
	if err != nil {
		return err
	}
	return workflowv3.StrictDecode(body, target)
}

func readPartition(context workflowv3runtime.TaskModuleContext) ([]Embedded, error) {
	ref, ok := context.Request.Inputs["partition"]
	if !ok {
		return nil, fmt.Errorf("partition input is missing")
	}
	body, err := workflowv3.ReadArtifact(context.Context, context.Request.Artifacts, ref)
	if err != nil {
		return nil, err
	}
	partition, err := workflowv3.DecodeReductionPartition(body, 32)
	if err != nil {
		return nil, err
	}
	items := []Embedded{}
	for _, member := range partition.Members {
		memberBody, err := workflowv3.ReadArtifact(context.Context, context.Request.Artifacts, member.Value)
		if err != nil {
			return nil, err
		}
		switch member.Value.Schema {
		case ShardSchema:
			var shard PreparedShard
			if err := workflowv3.StrictDecode(memberBody, &shard); err != nil {
				return nil, err
			}
			items = append(items, shard.Items...)
		default:
			return nil, fmt.Errorf("unexpected partition member schema")
		}
	}
	return items, nil
}

func successResponse(value any, usage []Usage) map[string]any {
	body, err := json.Marshal(value)
	if err != nil {
		return closedFailure("internal", "RAG_TTC_RESULT_ENCODE", false)
	}
	var jsonValue any
	if err := json.Unmarshal(body, &jsonValue); err != nil {
		return closedFailure("internal", "RAG_TTC_RESULT_ENCODE", false)
	}
	usageValues := make([]map[string]any, len(usage))
	for index, amount := range usage {
		usageValues[index] = map[string]any{"dimension": amount.Dimension, "units": amount.Units}
	}
	return map[string]any{"ok": true, "value": jsonValue, "usage": usageValues}
}

func closedFailure(class, code string, retryable bool) map[string]any {
	return closedFailureWithUsage(class, code, retryable, nil)
}

func closedFailureWithUsage(class, code string, retryable bool, usage []Usage) map[string]any {
	usageValues := make([]map[string]any, len(usage))
	for index, amount := range usage {
		usageValues[index] = map[string]any{"dimension": amount.Dimension, "units": amount.Units}
	}
	return map[string]any{"ok": false, "usage": usageValues, "failure": map[string]any{
		"class": class, "code": code, "retryable": retryable,
		"message": "RAG task failed closed validation",
	}}
}

func providerFailure(err error) map[string]any {
	if failure, ok := err.(*Failure); ok {
		return closedFailure(failure.Class, failure.Code, failure.Retryable)
	}
	return closedFailure("provider", "RAG_TTC_PROVIDER_FAILURE", true)
}
