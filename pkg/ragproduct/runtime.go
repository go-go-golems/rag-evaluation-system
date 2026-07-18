// Package ragproduct executes canonical product plans with host-owned lifecycle.
// It deliberately imports neither researchctl nor the RAG research adapter.
package ragproduct

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type ArtifactSink interface {
	Put(context.Context, string, string, []byte) error
}

type Bindings struct {
	Corpus    ragoperators.CorpusArtifact
	Manifests ragoperators.ManifestResolver
	Schemas   ragoperators.OutputSchemaValidator
	Generator ragoperators.TextGenerator
	Embedder  ragoperators.Embedder
	Reranker  ragoperators.Reranker
	Cache     ragoperators.Cache
	Traces    ArtifactSink
}

type Request struct {
	ID     string         `json:"id,omitempty"`
	Values map[string]any `json:"values"`
}

type Citation struct {
	ChunkID string                  `json:"chunkId"`
	Source  ragcontract.CitationRef `json:"source"`
}

type Result struct {
	Rank      int                          `json:"rank"`
	Collapse  ragcontract.CollapseIdentity `json:"collapse"`
	ChunkID   string                       `json:"chunkId"`
	Citations []Citation                   `json:"citations"`
}

type Response struct {
	TraceID   string                    `json:"traceId,omitempty"`
	Answer    string                    `json:"answer,omitempty"`
	Results   []Result                  `json:"results"`
	Citations []Citation                `json:"citations"`
	Abstained bool                      `json:"abstained"`
	Failure   *ragcontract.FailureTrace `json:"failure,omitempty"`
	Trace     *ragcontract.QueryTrace   `json:"trace,omitempty"`
}

type Runtime struct {
	plan      ragcontract.ProductPlan
	planID    string
	engine    *ragengine.Engine
	options   ragengine.Options
	prepared  *ragengine.Prepared
	semaphore chan struct{}
	request   requestContract
	response  responseContract
	closeMu   sync.Mutex
	closed    bool
	traceSink ArtifactSink
}

type requestField struct {
	Name, Type string
	Required   bool
	MaxLength  int
}
type requestContract struct {
	Fields []requestField `json:"fields"`
}
type responseContract struct {
	AnswerFormat, CitationMode string
	IncludeTraceID             bool `json:"includeTraceId"`
}

func Load(reader io.Reader) (ragcontract.ProductPlan, error) {
	plan, err := ragcontract.DecodeProduct(reader)
	if err != nil {
		return plan, err
	}
	return ragcompiler.CompileProduct(plan, nil)
}

func LoadBytes(data []byte) (ragcontract.ProductPlan, error) { return Load(bytes.NewReader(data)) }

func New(ctx context.Context, plan ragcontract.ProductPlan, bindings Bindings) (*Runtime, error) {
	compiled, err := ragcompiler.CompileProduct(plan, nil)
	if err != nil {
		return nil, err
	}
	planID, err := ragcompiler.ProductSemanticIdentity(compiled)
	if err != nil {
		return nil, err
	}
	if err := validateCorpusBinding(compiled, bindings.Corpus); err != nil {
		return nil, err
	}
	if err := validateResolvedModels(compiled.Models, bindings.Manifests); err != nil {
		return nil, err
	}
	var req requestContract
	if err := ragcontract.DecodeStrict(compiled.Request, &req); err != nil {
		return nil, err
	}
	var res responseContract
	if err := ragcontract.DecodeStrict(compiled.Response, &res); err != nil {
		return nil, err
	}
	options := ragengine.Options{Manifests: bindings.Manifests, Schemas: bindings.Schemas, Generator: bindings.Generator, Embedder: bindings.Embedder, Reranker: bindings.Reranker, Cache: bindings.Cache}
	engine := ragengine.New(nil)
	prepared, err := engine.Prepare(ctx, compiled.Pipeline, bindings.Corpus.Corpus, options)
	if err != nil {
		return nil, fmt.Errorf("RAG_PRODUCT_PREPARE: %w", err)
	}
	return &Runtime{plan: compiled, planID: planID, engine: engine, options: options, prepared: prepared, semaphore: make(chan struct{}, compiled.Runtime.MaxConcurrent), request: req, response: res, traceSink: bindings.Traces}, nil
}

func (r *Runtime) PlanID() string { return r.planID }

func (r *Runtime) Execute(ctx context.Context, request Request) (Response, error) {
	if err := r.validateRequest(request); err != nil {
		return Response{}, err
	}
	r.closeMu.Lock()
	closed := r.closed
	r.closeMu.Unlock()
	if closed {
		return Response{}, fmt.Errorf("RAG_PRODUCT_CLOSED")
	}
	select {
	case r.semaphore <- struct{}{}:
	case <-ctx.Done():
		return Response{}, ctx.Err()
	}
	defer func() { <-r.semaphore }()
	r.closeMu.Lock()
	closed = r.closed
	r.closeMu.Unlock()
	if closed {
		return Response{}, fmt.Errorf("RAG_PRODUCT_CLOSED")
	}
	if timeout := r.plan.Runtime.TimeoutMilliseconds; timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancel()
	}
	queryText, _ := request.Values["query"].(string)
	queryID := request.ID
	if queryID == "" {
		queryID = newTraceID()
	}
	traceID := newTraceID()
	execution := productExecution(r.plan)
	options := r.options
	options.Prepared = r.prepared
	output, err := r.engine.Execute(ctx, execution, r.optionsCorpus(), ragoperators.EvaluationDataset{SchemaVersion: "rag-product-request/v1", Queries: []ragoperators.Query{{ID: queryID, Text: queryText}}}, nil, options)
	response := r.buildResponse(traceID, output, err)
	if err != nil && r.plan.Runtime.FailurePolicy == "fail" {
		return Response{}, fmt.Errorf("RAG_PRODUCT_EXECUTE: %w", err)
	}
	if err != nil && r.plan.Runtime.FailurePolicy == "abstain" {
		response.Abstained = true
		response.Answer = ""
		response.Results = nil
		response.Citations = nil
	}
	if r.plan.Citations.Mode == "required" && !response.Abstained && len(response.Citations) == 0 {
		return Response{}, fmt.Errorf("RAG_PRODUCT_CITATIONS_REQUIRED")
	}
	if err := r.applyTracePolicy(ctx, traceID, &response); err != nil {
		return Response{}, err
	}
	if !r.response.IncludeTraceID {
		response.TraceID = ""
	}
	return response, nil
}

// corpus is retained inside the prepared graph; Execute still requires the same immutable input value.
func (r *Runtime) optionsCorpus() ragoperators.Corpus {
	// Static nodes are skipped, so this value is used only for the reserved corpus input identity.
	return ragoperators.Corpus{}
}

func (r *Runtime) buildResponse(traceID string, output *ragengine.Result, executionErr error) Response {
	response := Response{TraceID: traceID, Results: []Result{}, Citations: []Citation{}}
	if output == nil {
		output = &ragengine.Result{}
	}
	if len(output.Answers) > 0 {
		response.Answer = output.Answers[0].Text
		response.Abstained = output.Answers[0].Abstained
	}
	if len(output.Traces) > 0 {
		trace := output.Traces[0]
		for index := range trace.Failures {
			trace.Failures[index].Message = "product operation failed"
		}
		response.Trace = &trace
		for _, item := range trace.Results {
			citation := Citation{ChunkID: item.Evidence.ChunkID, Source: item.Evidence.Citation}
			response.Results = append(response.Results, Result{Rank: item.Rank, Collapse: item.Collapse, ChunkID: item.Evidence.ChunkID, Citations: []Citation{citation}})
			response.Citations = append(response.Citations, citation)
		}
		if len(response.Results) == 0 && trace.Hydration != nil {
			for index, evidence := range trace.Hydration.Selected {
				citation := Citation{ChunkID: evidence.ChunkID, Source: evidence.Citation}
				response.Results = append(response.Results, Result{Rank: index + 1, ChunkID: evidence.ChunkID, Citations: []Citation{citation}})
				response.Citations = append(response.Citations, citation)
			}
		}
	}
	if executionErr != nil {
		response.Failure = &ragcontract.FailureTrace{Code: "RAG_PRODUCT_EXECUTION_FAILED", Path: "$.request", Message: "product execution failed"}
	}
	return response
}

func (r *Runtime) applyTracePolicy(ctx context.Context, traceID string, response *Response) error {
	switch r.plan.Runtime.TracePolicy {
	case "authoritative":
		return nil
	case "metadata-only":
		if response.Trace != nil {
			response.Trace = metadataTrace(*response.Trace)
		}
	case "artifact-backed":
		if r.traceSink == nil {
			return fmt.Errorf("RAG_PRODUCT_TRACE_SINK_REQUIRED")
		}
		data, _ := ragcontract.CanonicalJSON(response.Trace)
		if err := r.traceSink.Put(ctx, traceID, ragcontract.TraceSchemaVersion, data); err != nil {
			return fmt.Errorf("RAG_PRODUCT_TRACE_STORE: %w", err)
		}
		response.Trace = nil
	case "none":
		response.Trace = nil
	}
	return nil
}

func (r *Runtime) Close() error {
	r.closeMu.Lock()
	if r.closed {
		r.closeMu.Unlock()
		return nil
	}
	r.closed = true
	r.closeMu.Unlock()
	for i := 0; i < cap(r.semaphore); i++ {
		r.semaphore <- struct{}{}
	}
	defer func() {
		for i := 0; i < cap(r.semaphore); i++ {
			<-r.semaphore
		}
	}()
	return r.prepared.Close()
}

func (r *Runtime) validateRequest(request Request) error {
	declared := map[string]requestField{}
	for _, field := range r.request.Fields {
		declared[field.Name] = field
	}
	for name := range request.Values {
		if _, ok := declared[name]; !ok {
			return fmt.Errorf("RAG_PRODUCT_REQUEST_UNKNOWN: %s", name)
		}
	}
	for _, field := range r.request.Fields {
		value, ok := request.Values[field.Name]
		if !ok {
			if field.Required {
				return fmt.Errorf("RAG_PRODUCT_REQUEST_REQUIRED: %s", field.Name)
			}
			continue
		}
		valid := false
		switch field.Type {
		case "string":
			_, valid = value.(string)
		case "number":
			_, valid = value.(float64)
			if !valid {
				_, valid = value.(json.Number)
			}
		case "boolean":
			_, valid = value.(bool)
		case "object":
			_, valid = value.(map[string]any)
		case "array":
			_, valid = value.([]any)
		}
		if !valid {
			return fmt.Errorf("RAG_PRODUCT_REQUEST_TYPE: %s", field.Name)
		}
		if text, ok := value.(string); ok && field.MaxLength > 0 && len([]rune(text)) > field.MaxLength {
			return fmt.Errorf("RAG_PRODUCT_REQUEST_LENGTH: %s", field.Name)
		}
	}
	return nil
}

func validateCorpusBinding(plan ragcontract.ProductPlan, artifact ragoperators.CorpusArtifact) error {
	if err := ragcontract.ValidateManifestBase(artifact.Manifest.ManifestBase, ragcontract.CorpusManifestSchema, false); err != nil {
		return err
	}
	digest, _ := ragcontract.Digest(artifact.Corpus)
	if digest != artifact.Manifest.Digest {
		return fmt.Errorf("RAG_PRODUCT_CORPUS_DIGEST")
	}
	for _, binding := range plan.Bindings {
		if binding.Role == "corpus" {
			if binding.Digest != artifact.Manifest.Digest {
				return fmt.Errorf("RAG_PRODUCT_CORPUS_BINDING")
			}
			return nil
		}
	}
	return fmt.Errorf("RAG_PRODUCT_CORPUS_MISSING")
}

func validateResolvedModels(models []ragcontract.ModelBinding, resolver ragoperators.ManifestResolver) error {
	if len(models) == 0 {
		return nil
	}
	if resolver == nil {
		return fmt.Errorf("RAG_PRODUCT_MANIFEST_RESOLVER_REQUIRED")
	}
	for _, binding := range models {
		manifest, err := resolver.Model(binding.Reference)
		if err != nil {
			return fmt.Errorf("RAG_PRODUCT_MODEL_RESOLVE: %s: %w", binding.Reference, err)
		}
		if manifest.Digest != binding.Digest || manifest.SchemaVersion != binding.Manifest {
			return fmt.Errorf("RAG_PRODUCT_MODEL_BINDING: %s", binding.Reference)
		}
	}
	return nil
}

func productExecution(plan ragcontract.ProductPlan) ragcontract.PipelineExecution {
	execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: plan.Pipeline, Bindings: plan.Bindings, Dataset: ragcontract.DatasetBinding{Split: "product", Status: "online", RelevanceTarget: "none"}, VariantID: "product", Factors: []ragcontract.FactorSelection{}}
	identity := execution
	identity.CellID = ""
	execution.CellID, _ = ragcontract.Digest(identity)
	return execution
}

func metadataTrace(trace ragcontract.QueryTrace) *ragcontract.QueryTrace {
	return &ragcontract.QueryTrace{SchemaVersion: trace.SchemaVersion, Query: trace.Query, Operators: trace.Operators, Timing: trace.Timing, Usage: trace.Usage, Failures: trace.Failures}
}

func newTraceID() string {
	var value [16]byte
	_, _ = rand.Read(value[:])
	return "trc_" + hex.EncodeToString(value[:])
}

func DecodeRequest(reader io.Reader) (Request, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	var request Request
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return request, fmt.Errorf("RAG_PRODUCT_REQUEST_TRAILING")
		}
		return request, err
	}
	if request.Values == nil {
		return request, fmt.Errorf("RAG_PRODUCT_REQUEST_VALUES")
	}
	return request, nil
}
