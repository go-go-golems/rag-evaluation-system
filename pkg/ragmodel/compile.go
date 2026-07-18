package ragmodel

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type CompileOptions struct {
	Inputs map[string]ragcontract.ArtifactBinding `json:"inputs"`
	Models []ragcontract.ModelBinding             `json:"models,omitempty"`
}
type PreviewOptions struct {
	Variant string            `json:"variant"`
	Factors map[string]string `json:"factors"`
	Query   string            `json:"query"`
	Trace   string            `json:"trace"`
}
type PreviewRequest struct {
	SchemaVersion string                        `json:"schemaVersion"`
	Cell          ragcontract.PipelineExecution `json:"cell"`
	Query         string                        `json:"query"`
	Trace         string                        `json:"trace"`
}
type Explanation struct {
	SchemaVersion string   `json:"schemaVersion"`
	Kind          string   `json:"kind"`
	Name          string   `json:"name"`
	NodeCount     int      `json:"nodeCount,omitempty"`
	VariantCount  int      `json:"variantCount,omitempty"`
	CellCount     int      `json:"cellCount,omitempty"`
	Operators     []string `json:"operators,omitempty"`
	Factors       []string `json:"factors,omitempty"`
	Warnings      []string `json:"warnings,omitempty"`
}

func BuildIR(pipeline *Pipeline, query *QueryPlan, reranker, generator *Descriptor) (ragcontract.PipelineIR, error) {
	if err := ValidatePipeline(pipeline); err != nil {
		return ragcontract.PipelineIR{}, err
	}
	ir := ragcontract.PipelineIR{SchemaVersion: ragcontract.PipelineSchemaVersion, Inputs: []ragcontract.InputSlot{{ID: "corpus", Kind: ragcontract.PortCorpus, BindingMode: "artifact", ArtifactRole: pipeline.Corpus.Role, ManifestSchema: pipeline.Corpus.ManifestSchema}, {ID: "query", Kind: ragcontract.PortQuery, BindingMode: "request"}}}
	ir.Nodes = append(ir.Nodes, node("units", pipeline.Unitizer, map[string]ragcontract.PortRef{"corpus": {NodeID: "corpus", Port: "out"}}), node("chunks", pipeline.Chunker, map[string]ragcontract.PortRef{"units": {NodeID: "units", Port: "units"}}))
	representations := pipeline.Representations
	if representations.Operator.Kind == "representations.compose" {
		operators := []map[string]any{}
		for _, child := range representations.Children {
			operators = append(operators, map[string]any{"kind": child.Operator.Kind, "version": child.Operator.Version, "config": json.RawMessage(child.Config)})
		}
		raw, _ := json.Marshal(map[string]any{"operators": operators})
		representations = &Descriptor{Kind: KindRepresentations, Operator: representations.Operator, Config: raw}
	}
	representationInputPort := "chunks"
	representationOutputPort := "representations"
	if representations.Operator.Kind == "representations.compose" {
		representationInputPort = "in"
		representationOutputPort = "out"
	}
	ir.Nodes = append(ir.Nodes, node("representations", representations, map[string]ragcontract.PortRef{representationInputPort: {NodeID: "chunks", Port: "chunks"}}))
	indexInputs := map[string]ragcontract.PortRef{"representations.all": {NodeID: "representations", Port: representationOutputPort}}
	if pipeline.Embedding != nil {
		ir.Nodes = append(ir.Nodes, node("embedding", pipeline.Embedding, map[string]ragcontract.PortRef{"representations": indexInputs["representations.all"]}))
		indexInputs["embeddings"] = ragcontract.PortRef{NodeID: "embedding", Port: "embeddings"}
	}
	ir.Nodes = append(ir.Nodes, node("index", pipeline.Index, indexInputs))
	last := ragcontract.PortRef{NodeID: "index", Port: "index"}
	lastKind := ragcontract.PortIndex
	if query != nil {
		if len(query.Channels) == 0 || query.ChannelCollapse == nil || query.Fusion == nil || query.FinalCollapse == nil || query.Hydration == nil {
			return ir, fmt.Errorf("RAG_V2_QUERY_INCOMPLETE: channels, both collapse stages, fusion, and hydration are required")
		}
		fusionInputs := map[string]ragcontract.PortRef{}
		for i, ch := range query.Channels {
			id := fmt.Sprintf("channel.%03d", i+1)
			n := node(id, ch, map[string]ragcontract.PortRef{"index": {NodeID: "index", Port: "index"}, "query": {NodeID: "query", Port: "out"}})
			n.Order = i
			ir.Nodes = append(ir.Nodes, n)
			cid := id + ".collapse"
			c := node(cid, query.ChannelCollapse, map[string]ragcontract.PortRef{"hits": {NodeID: id, Port: "hits"}})
			c.Order = i
			ir.Nodes = append(ir.Nodes, c)
			fusionInputs["channel."+ch.Name] = ragcontract.PortRef{NodeID: cid, Port: "parents"}
		}
		ir.Nodes = append(ir.Nodes, node("fusion", query.Fusion, fusionInputs))
		final := cloneDescriptor(query.FinalCollapse)
		final.Operator = ragcontract.OperatorRef{Kind: "collapse.final", Version: "v1"}
		hydration := cloneDescriptor(query.Hydration)
		var hydrationConfig map[string]any
		_ = json.Unmarshal(hydration.Config, &hydrationConfig)
		if hydrationConfig == nil {
			hydrationConfig = map[string]any{}
		}
		hydrationConfig["results"] = query.Results
		hydration.Config, _ = json.Marshal(hydrationConfig)
		ir.Nodes = append(ir.Nodes, node("final", final, map[string]ragcontract.PortRef{"parents": {NodeID: "fusion", Port: "parents"}}), node("hydrate", hydration, map[string]ragcontract.PortRef{"parents": {NodeID: "final", Port: "parents"}, "chunks": {NodeID: "chunks", Port: "chunks"}}))
		last = ragcontract.PortRef{NodeID: "hydrate", Port: "evidence"}
		lastKind = ragcontract.PortEvidence
	}
	if reranker != nil {
		ir.Nodes = append(ir.Nodes, node("rerank", reranker, map[string]ragcontract.PortRef{"evidence": last}))
		last = ragcontract.PortRef{NodeID: "rerank", Port: "evidence"}
		lastKind = ragcontract.PortEvidence
	}
	if generator != nil {
		ir.Nodes = append(ir.Nodes, node("generate", generator, map[string]ragcontract.PortRef{"evidence": last}))
		last = ragcontract.PortRef{NodeID: "generate", Port: "answer"}
		lastKind = ragcontract.PortAnswer
	}
	ir.Outputs = []ragcontract.OutputRef{{Name: "result", Kind: lastKind, From: last}}
	return ir, nil
}
func node(id string, d *Descriptor, inputs map[string]ragcontract.PortRef) ragcontract.Node {
	bindings := make([]ragcontract.InputBinding, 0, len(inputs))
	for port, from := range inputs {
		bindings = append(bindings, ragcontract.InputBinding{Port: port, From: from})
	}
	sort.Slice(bindings, func(i, j int) bool { return bindings[i].Port < bindings[j].Port })
	return ragcontract.Node{ID: id, Operator: d.Operator, Inputs: bindings, Config: append(json.RawMessage(nil), d.Config...)}
}
func cloneDescriptor(v *Descriptor) *Descriptor {
	if v == nil {
		return nil
	}
	c := *v
	c.Config = append(json.RawMessage(nil), v.Config...)
	return &c
}

func CompileProduct(value *Product, options CompileOptions) (ragcontract.ProductPlan, error) {
	if value == nil || value.Pipeline == nil || value.Query == nil {
		return ragcontract.ProductPlan{}, fmt.Errorf("RAG_V2_PRODUCT_INCOMPLETE: pipeline and query are required")
	}
	reranker := normalizedReranker(value.Reranker)
	ir, err := BuildIR(value.Pipeline, value.Query, reranker, value.Generator)
	if err != nil {
		return ragcontract.ProductPlan{}, err
	}
	request, _ := json.Marshal(value.Request)
	response, _ := json.Marshal(value.Response)
	bindings, err := bindingsForIR(ir, options.Inputs)
	if err != nil {
		return ragcontract.ProductPlan{}, err
	}
	plan := ragcontract.ProductPlan{SchemaVersion: ragcontract.ProductSchemaVersion, Pipeline: ir, Bindings: bindings, Models: options.Models, Citations: ragcontract.CitationPolicy{Mode: value.Response.CitationMode, RequireSourceText: value.Response.CitationMode == "source" || value.Response.CitationMode == "required"}, Request: request, Response: response, Runtime: value.Runtime, Display: ragcontract.DisplayMetadata{Name: value.Name, Tags: value.Tags}}
	if value.Query.Results > 0 {
		plan.Runtime.MaxResults = value.Query.Results
	}
	return ragcompiler.CompileProduct(plan, nil)
}
func normalizedReranker(v *Descriptor) *Descriptor {
	if v == nil {
		return nil
	}
	c := cloneDescriptor(v)
	var config CrossEncoderConfig
	if json.Unmarshal(c.Config, &config) == nil {
		c.Config, _ = json.Marshal(map[string]any{"model": config.Model, "candidateCount": config.Candidates, "results": config.Results, "truncation": config.Truncation, "tokenization": config.Tokenization, "inputTemplate": config.InputTemplate, "timeoutMilliseconds": config.TimeoutMilliseconds})
	}
	return c
}

func CompileStudy(value *Study, options CompileOptions) (ragcontract.Study, []ragcontract.ExpandedCell, error) {
	if value == nil || value.Pipeline == nil || value.Dataset.Role == "" || len(value.Variants) == 0 {
		return ragcontract.Study{}, nil, fmt.Errorf("RAG_V2_STUDY_INCOMPLETE: pipeline, dataset, and variants are required")
	}
	variants := make([]ragcontract.Variant, 0, len(value.Variants))
	for _, v := range value.Variants {
		if v.Query == nil {
			return ragcontract.Study{}, nil, fmt.Errorf("RAG_V2_VARIANT_QUERY: %s", v.ID)
		}
		ir, err := BuildIR(value.Pipeline, v.Query, nil, nil)
		if err != nil {
			return ragcontract.Study{}, nil, err
		}
		variants = append(variants, ragcontract.Variant{ID: v.ID, Pipeline: ir, Metadata: mustRaw(map[string]any{"representations": v.Representations})})
	}
	factors := make([]ragcontract.Factor, 0, len(value.Factors))
	for _, factor := range value.Factors {
		values := make([]ragcontract.FactorValue, 0, len(factor.Values))
		for _, v := range factor.Values {
			values = append(values, ragcontract.FactorValue{ID: v, Value: mustRaw(v), Overrides: nil})
		}
		factors = append(factors, ragcontract.Factor{ID: factor.ID, Values: values})
	}
	datasetBinding, ok := options.Inputs[value.Dataset.Role]
	if !ok {
		return ragcontract.Study{}, nil, fmt.Errorf("RAG_V2_DATASET_BINDING: %s", value.Dataset.Role)
	}
	allIRInputs := variants[0].Pipeline
	bindings, err := bindingsForIR(allIRInputs, options.Inputs)
	if err != nil {
		return ragcontract.Study{}, nil, err
	}
	study := ragcontract.Study{SchemaVersion: ragcontract.StudySchemaVersion, Variants: variants, Factors: factors, Bindings: bindings, Dataset: ragcontract.DatasetBinding{ManifestDigest: datasetBinding.Digest, Split: value.Dataset.Split, Status: value.Dataset.Status, RelevanceTarget: value.Dataset.RelevanceTarget}, Measures: value.Measures, Replicates: value.Replicates, Acceptance: mustRaw(map[string]any{"invariants": value.Invariants}), Display: ragcontract.DisplayMetadata{Name: value.Name, Tags: value.Tags}}
	cells, err := ragcompiler.ExpandStudy(study, nil)
	return study, cells, err
}
func bindingsForIR(ir ragcontract.PipelineIR, inputs map[string]ragcontract.ArtifactBinding) ([]ragcontract.ArtifactBinding, error) {
	result := []ragcontract.ArtifactBinding{}
	for _, slot := range ir.Inputs {
		if slot.BindingMode != "artifact" {
			continue
		}
		binding, ok := inputs[slot.ArtifactRole]
		if !ok {
			return nil, fmt.Errorf("RAG_V2_INPUT_BINDING: %s", slot.ArtifactRole)
		}
		binding.SlotID = slot.ID
		if binding.Role == "" {
			binding.Role = slot.ArtifactRole
		}
		if binding.SchemaVersion == "" {
			binding.SchemaVersion = slot.ManifestSchema
		}
		result = append(result, binding)
	}
	return result, nil
}
func mustRaw(v any) json.RawMessage {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return raw
}

func Explain(value any) (Explanation, error) {
	switch v := value.(type) {
	case *Pipeline:
		ir, err := BuildIR(v, nil, nil, nil)
		if err != nil {
			return Explanation{}, err
		}
		return explainIR(v.Name, "pipeline", ir), nil
	case *Product:
		ir, err := BuildIR(v.Pipeline, v.Query, normalizedReranker(v.Reranker), v.Generator)
		if err != nil {
			return Explanation{}, err
		}
		return explainIR(v.Name, "product", ir), nil
	case *Study:
		e := Explanation{SchemaVersion: "rag-explanation/v1", Kind: "study", Name: v.Name, VariantCount: len(v.Variants)}
		for _, f := range v.Factors {
			e.Factors = append(e.Factors, f.ID)
			if e.CellCount == 0 {
				e.CellCount = 1
			}
			e.CellCount *= len(f.Values)
		}
		e.CellCount *= len(v.Variants)
		return e, nil
	default:
		return Explanation{}, fmt.Errorf("RAG_V2_EXPLAIN_TYPE: unsupported value")
	}
}
func explainIR(name, kind string, ir ragcontract.PipelineIR) Explanation {
	e := Explanation{SchemaVersion: "rag-explanation/v1", Kind: kind, Name: name, NodeCount: len(ir.Nodes)}
	for _, n := range ir.Nodes {
		e.Operators = append(e.Operators, n.Operator.ID())
	}
	return e
}
func Preview(value *Study, compile CompileOptions, options PreviewOptions) (PreviewRequest, error) {
	_, cells, err := CompileStudy(value, compile)
	if err != nil {
		return PreviewRequest{}, err
	}
	for _, cell := range cells {
		if cell.VariantID != options.Variant {
			continue
		}
		match := true
		for k, v := range options.Factors {
			found := false
			for _, s := range cell.Factors {
				if s.FactorID == k && s.ValueID == v {
					found = true
				}
			}
			if !found {
				match = false
			}
		}
		if match {
			return PreviewRequest{SchemaVersion: "rag-preview-request/v1", Cell: cell.Execution, Query: options.Query, Trace: options.Trace}, nil
		}
	}
	return PreviewRequest{}, fmt.Errorf("RAG_V2_PREVIEW_CELL: no matching cell")
}
