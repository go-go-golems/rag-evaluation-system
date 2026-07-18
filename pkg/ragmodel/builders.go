package ragmodel

import (
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type Pipeline struct {
	Name                                                 string
	Corpus                                               CorpusInput
	Unitizer, Chunker, Representations, Embedding, Index *Descriptor
	IndexName                                            string
	Notes                                                []string
	Tags                                                 map[string]string
}
type PipelineBuilder struct{ value Pipeline }

func NewPipelineBuilder(name string) *PipelineBuilder {
	return &PipelineBuilder{value: Pipeline{Name: name, Tags: map[string]string{}}}
}

type Fragment struct {
	Name  string
	patch Pipeline
}

func NewPipeline(name string, configure func(*PipelineBuilder)) *Pipeline {
	b := NewPipelineBuilder(name)
	if configure != nil {
		configure(b)
	}
	v := b.value
	return &v
}
func NewFragment(name string, configure func(*PipelineBuilder)) *Fragment {
	p := NewPipeline(name, configure)
	return FragmentFromPipeline(name, p)
}
func FragmentFromPipeline(name string, pipeline *Pipeline) *Fragment {
	return &Fragment{Name: name, patch: *pipeline}
}
func (b *PipelineBuilder) CorpusInput(v CorpusInput) *PipelineBuilder { b.value.Corpus = v; return b }
func (b *PipelineBuilder) Units(v *Descriptor) *PipelineBuilder       { b.value.Unitizer = v; return b }
func (b *PipelineBuilder) Chunks(v *Descriptor) *PipelineBuilder      { b.value.Chunker = v; return b }
func (b *PipelineBuilder) Represent(v *Descriptor) *PipelineBuilder {
	b.value.Representations = v
	return b
}
func (b *PipelineBuilder) EmbeddingModel(v *Descriptor) *PipelineBuilder {
	b.value.Embedding = v
	return b
}
func (b *PipelineBuilder) IndexNamed(name string, v *Descriptor) *PipelineBuilder {
	b.value.IndexName = name
	b.value.Index = v
	return b
}
func (b *PipelineBuilder) Use(v *Fragment) *PipelineBuilder {
	if v == nil {
		return b
	}
	p := v.patch
	if p.Corpus.Role != "" {
		b.value.Corpus = p.Corpus
	}
	if p.Unitizer != nil {
		b.value.Unitizer = p.Unitizer
	}
	if p.Chunker != nil {
		b.value.Chunker = p.Chunker
	}
	if p.Representations != nil {
		b.value.Representations = p.Representations
	}
	if p.Embedding != nil {
		b.value.Embedding = p.Embedding
	}
	if p.Index != nil {
		b.value.Index = p.Index
		b.value.IndexName = p.IndexName
	}
	b.value.Notes = append(b.value.Notes, p.Notes...)
	for k, x := range p.Tags {
		b.value.Tags[k] = x
	}
	return b
}
func (b *PipelineBuilder) Note(v string) *PipelineBuilder {
	b.value.Notes = append(b.value.Notes, v)
	return b
}
func (b *PipelineBuilder) Tag(k, v string) *PipelineBuilder { b.value.Tags[k] = v; return b }
func (b *PipelineBuilder) Build() *Pipeline                 { v := b.value; return &v }

type QueryPlan struct {
	Name                                              string
	Channels                                          []*Descriptor
	ChannelCollapse, Fusion, FinalCollapse, Hydration *Descriptor
	Results                                           int
}
type QueryBuilder struct{ value QueryPlan }

func NewQueryBuilder(name string) *QueryBuilder {
	return &QueryBuilder{value: QueryPlan{Name: name, Results: 5}}
}

func NewQueryPlan(name string, configure func(*QueryBuilder)) *QueryPlan {
	b := NewQueryBuilder(name)
	if configure != nil {
		configure(b)
	}
	v := b.value
	return &v
}
func (b *QueryBuilder) Channels(values ...*Descriptor) *QueryBuilder {
	b.value.Channels = append([]*Descriptor(nil), values...)
	return b
}
func (b *QueryBuilder) CollapseChannels(v *Descriptor) *QueryBuilder {
	b.value.ChannelCollapse = v
	return b
}
func (b *QueryBuilder) Fuse(v *Descriptor) *QueryBuilder { b.value.Fusion = v; return b }
func (b *QueryBuilder) CollapseFinal(v *Descriptor) *QueryBuilder {
	b.value.FinalCollapse = v
	return b
}
func (b *QueryBuilder) Hydrate(v *Descriptor) *QueryBuilder { b.value.Hydration = v; return b }
func (b *QueryBuilder) ResultCount(v int) *QueryBuilder     { b.value.Results = v; return b }
func (b *QueryBuilder) Build() *QueryPlan                   { v := b.value; return &v }

type FieldContract struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Required  bool   `json:"required"`
	MaxLength int    `json:"maxLength,omitempty"`
}
type RequestContract struct {
	Fields []FieldContract `json:"fields"`
}
type RequestBuilder struct{ value RequestContract }

func (b *RequestBuilder) Field(name, kind string, required bool, maxLength int) *RequestBuilder {
	b.value.Fields = append(b.value.Fields, FieldContract{name, kind, required, maxLength})
	return b
}

type ResponseContract struct {
	AnswerFormat   string `json:"answerFormat"`
	CitationMode   string `json:"citationMode"`
	IncludeTraceID bool   `json:"includeTraceId"`
}
type ResponseBuilder struct{ value ResponseContract }

func (b *ResponseBuilder) Answer(v string) *ResponseBuilder    { b.value.AnswerFormat = v; return b }
func (b *ResponseBuilder) Citations(v string) *ResponseBuilder { b.value.CitationMode = v; return b }
func (b *ResponseBuilder) TraceID(v bool) *ResponseBuilder     { b.value.IncludeTraceID = v; return b }

type RuntimeBuilder struct{ value ragcontract.RuntimePolicy }

func (b *RuntimeBuilder) Timeout(v int64) *RuntimeBuilder  { b.value.TimeoutMilliseconds = v; return b }
func (b *RuntimeBuilder) Concurrent(v int) *RuntimeBuilder { b.value.MaxConcurrent = v; return b }
func (b *RuntimeBuilder) ProviderFailure(v string) *RuntimeBuilder {
	b.value.FailurePolicy = v
	return b
}
func (b *RuntimeBuilder) Trace(v string) *RuntimeBuilder { b.value.TracePolicy = v; return b }

type Product struct {
	Name                string
	Pipeline            *Pipeline
	Query               *QueryPlan
	Reranker, Generator *Descriptor
	Request             RequestContract
	Response            ResponseContract
	Runtime             ragcontract.RuntimePolicy
	Tags                map[string]string
}
type ProductBuilder struct{ value Product }

func NewProductBuilder(name string) *ProductBuilder {
	return &ProductBuilder{value: Product{Name: name, Tags: map[string]string{}}}
}

func NewProduct(name string, configure func(*ProductBuilder)) *Product {
	b := NewProductBuilder(name)
	if configure != nil {
		configure(b)
	}
	v := b.value
	return &v
}
func (b *ProductBuilder) PipelineValue(v *Pipeline) *ProductBuilder { b.value.Pipeline = v; return b }
func (b *ProductBuilder) QueryPlan(v *QueryPlan) *ProductBuilder    { b.value.Query = v; return b }
func (b *ProductBuilder) Rerank(v *Descriptor) *ProductBuilder      { b.value.Reranker = v; return b }
func (b *ProductBuilder) Generate(v *Descriptor) *ProductBuilder    { b.value.Generator = v; return b }
func (b *ProductBuilder) RequestContract(configure func(*RequestBuilder)) *ProductBuilder {
	x := &RequestBuilder{}
	configure(x)
	b.value.Request = x.value
	return b
}
func (b *ProductBuilder) ResponseContract(configure func(*ResponseBuilder)) *ProductBuilder {
	x := &ResponseBuilder{}
	configure(x)
	b.value.Response = x.value
	return b
}
func (b *ProductBuilder) RuntimePolicy(configure func(*RuntimeBuilder)) *ProductBuilder {
	x := &RuntimeBuilder{}
	configure(x)
	b.value.Runtime = x.value
	return b
}
func (b *ProductBuilder) Tag(k, v string) *ProductBuilder { b.value.Tags[k] = v; return b }
func (b *ProductBuilder) Build() *Product                 { v := b.value; return &v }

type Variant struct {
	ID              string
	Representations []string
	Query           *QueryPlan
}
type VariantBuilder struct{ value Variant }

func (b *VariantBuilder) SelectRepresentations(v ...string) *VariantBuilder {
	b.value.Representations = append([]string(nil), v...)
	return b
}
func (b *VariantBuilder) QueryPlan(v *QueryPlan) *VariantBuilder { b.value.Query = v; return b }
func (b *VariantBuilder) Build() *Variant                        { v := b.value; return &v }
func NewVariant(id string, configure func(*VariantBuilder)) *Variant {
	b := &VariantBuilder{value: Variant{ID: id}}
	if configure != nil {
		configure(b)
	}
	return b.Build()
}

type VariantsBuilder struct{ values []Variant }

func (b *VariantsBuilder) Add(id string, configure func(*VariantBuilder)) *VariantsBuilder {
	x := &VariantBuilder{value: Variant{ID: id}}
	configure(x)
	b.values = append(b.values, x.value)
	return b
}

type FactorEnum struct {
	ID     string
	Values []string
}
type FactorsBuilder struct{ values []FactorEnum }

func (b *FactorsBuilder) Enum(id string, values ...string) *FactorsBuilder {
	b.values = append(b.values, FactorEnum{ID: id, Values: append([]string(nil), values...)})
	return b
}

type MetricsBuilder struct{ values []ragcontract.Measure }

func (b *MetricsBuilder) add(name, kind, unit string, config any) *MetricsBuilder {
	raw, _ := json.Marshal(config)
	b.values = append(b.values, ragcontract.Measure{Name: name, Version: "v1", ValueKind: kind, Unit: unit, Required: true, Config: raw})
	return b
}
func (b *MetricsBuilder) PrecisionAt(v []int) *MetricsBuilder {
	return b.add("rag.precision", "object", "ratio", map[string]any{"cutoffs": v})
}
func (b *MetricsBuilder) RecallAt(v []int) *MetricsBuilder {
	return b.add("rag.recall", "object", "ratio", map[string]any{"cutoffs": v})
}
func (b *MetricsBuilder) HitRateAt(v []int) *MetricsBuilder {
	return b.add("rag.hit-rate", "object", "ratio", map[string]any{"cutoffs": v})
}
func (b *MetricsBuilder) MRR() *MetricsBuilder {
	return b.add("rag.mrr", "number", "ratio", map[string]any{})
}
func (b *MetricsBuilder) NDCGAt(v []int) *MetricsBuilder {
	return b.add("rag.ndcg", "object", "ratio", map[string]any{"cutoffs": v})
}
func (b *MetricsBuilder) Latency(v []string) *MetricsBuilder {
	return b.add("rag.latency", "object", "milliseconds", map[string]any{"stages": v})
}
func (b *MetricsBuilder) TokenUsage() *MetricsBuilder {
	return b.add("rag.token-usage", "object", "tokens", map[string]any{})
}
func (b *MetricsBuilder) ProviderCost() *MetricsBuilder {
	return b.add("rag.provider-cost", "object", "currency", map[string]any{})
}
func (b *MetricsBuilder) StorageBytes() *MetricsBuilder {
	return b.add("rag.storage-bytes", "object", "bytes", map[string]any{})
}
func (b *MetricsBuilder) FailureRates() *MetricsBuilder {
	return b.add("rag.failure-rates", "object", "ratio", map[string]any{})
}

type InvariantsBuilder struct{ values []string }

func (b *InvariantsBuilder) Require(v string) *InvariantsBuilder {
	b.values = append(b.values, v)
	return b
}

type Study struct {
	Name       string
	Pipeline   *Pipeline
	Dataset    DatasetRef
	Variants   []Variant
	Factors    []FactorEnum
	Replicates int
	Measures   []ragcontract.Measure
	Invariants []string
	Tags       map[string]string
}
type StudyBuilder struct{ value Study }

func NewStudyBuilder(name string) *StudyBuilder {
	return &StudyBuilder{value: Study{Name: name, Replicates: 1, Tags: map[string]string{}}}
}

func NewStudy(name string, configure func(*StudyBuilder)) *Study {
	b := NewStudyBuilder(name)
	if configure != nil {
		configure(b)
	}
	v := b.value
	return &v
}
func (b *StudyBuilder) PipelineValue(v *Pipeline) *StudyBuilder { b.value.Pipeline = v; return b }
func (b *StudyBuilder) DatasetRef(v DatasetRef) *StudyBuilder   { b.value.Dataset = v; return b }
func (b *StudyBuilder) VariantsList(configure func(*VariantsBuilder)) *StudyBuilder {
	x := &VariantsBuilder{}
	configure(x)
	b.value.Variants = x.values
	return b
}
func (b *StudyBuilder) FactorsList(configure func(*FactorsBuilder)) *StudyBuilder {
	x := &FactorsBuilder{}
	configure(x)
	b.value.Factors = x.values
	return b
}
func (b *StudyBuilder) ReplicateCount(v int) *StudyBuilder { b.value.Replicates = v; return b }
func (b *StudyBuilder) MetricsList(configure func(*MetricsBuilder)) *StudyBuilder {
	x := &MetricsBuilder{}
	configure(x)
	b.value.Measures = x.values
	return b
}
func (b *StudyBuilder) InvariantsList(configure func(*InvariantsBuilder)) *StudyBuilder {
	x := &InvariantsBuilder{}
	configure(x)
	b.value.Invariants = x.values
	return b
}
func (b *StudyBuilder) Tag(k, v string) *StudyBuilder { b.value.Tags[k] = v; return b }
func (b *StudyBuilder) Build() *Study                 { v := b.value; return &v }

func ValidatePipeline(value *Pipeline) error {
	if value == nil {
		return fmt.Errorf("RAG_V2_PIPELINE_REQUIRED")
	}
	if value.Name == "" || value.Corpus.Role == "" || value.Unitizer == nil || value.Chunker == nil || value.Representations == nil || value.Index == nil {
		return fmt.Errorf("RAG_V2_PIPELINE_INCOMPLETE: name, corpus, units, chunks, representations, and index are required")
	}
	return nil
}
