package ragoperators

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type fakeGenerator struct {
	calls     int
	fail      bool
	invalid   bool
	questions []string
	abstained bool
	request   GenerationRequest
}

func (g *fakeGenerator) Generate(_ context.Context, r GenerationRequest) (GenerationResult, error) {
	g.calls++
	g.request = r
	if g.fail {
		return GenerationResult{}, errors.New("generated failure")
	}
	if r.Kind == "representations.structured-summary" {
		if g.invalid {
			return GenerationResult{Text: `not-json`}, nil
		}
		return GenerationResult{Text: `{"summary":"ok"}`, InputTokens: 2, OutputTokens: 1, FinishReason: "stop"}, nil
	}
	if r.Kind == "representations.synthetic-questions" {
		q := append([]string(nil), g.questions...)
		if q == nil {
			q = make([]string, r.Count)
			for i := range q {
				q[i] = "question " + itoa(i+1)
			}
		}
		return GenerationResult{Questions: q, InputTokens: 2, OutputTokens: int64(len(q))}, nil
	}
	ids := []string{}
	for _, e := range r.Evidence {
		ids = append(ids, e.Chunk.Record.ID)
	}
	if g.abstained {
		return GenerationResult{Text: "", Abstained: true, FinishReason: "stop"}, nil
	}
	return GenerationResult{Text: "answer", CitationChunkIDs: ids, FinishReason: "stop"}, nil
}

type fakeSchemaValidator struct{}

func (fakeSchemaValidator) Validate(schema string, document json.RawMessage) error {
	if schema != "summary/v1" || !json.Valid(document) {
		return errors.New("invalid fixture schema")
	}
	return nil
}

type fakeEmbedder struct{ badCount, badDimension bool }

func (e fakeEmbedder) Embed(_ context.Context, _ string, texts []string) ([][]float64, Usage, error) {
	if e.badCount {
		return [][]float64{}, Usage{}, nil
	}
	out := make([][]float64, len(texts))
	for i, text := range texts {
		out[i] = []float64{float64(len(text)), 1}
		if e.badDimension && i == len(texts)-1 {
			out[i] = []float64{1}
		}
	}
	return out, Usage{EmbeddingTokens: int64(len(texts))}, nil
}

type fakeReranker struct{ incomplete bool }

func (r fakeReranker) Rerank(_ context.Context, request RerankRequest) ([]RerankScore, error) {
	out := []RerankScore{}
	for i, e := range request.Candidates {
		if r.incomplete && i == 0 {
			continue
		}
		out = append(out, RerankScore{ChunkID: e.Chunk.Record.ID, Score: float64(len(request.Candidates) - i)})
	}
	return out, nil
}

func TestGenerationUsageKeepsUnknownCostAbsent(t *testing.T) {
	env := &Environment{Usage: Usage{Cost: map[string]float64{}}}
	addGenerationUsage(env, "generator-primary", GenerationResult{InputTokens: 2, OutputTokens: 3})
	if _, found := env.Usage.Cost["generator-primary"]; found {
		t.Fatalf("unknown provider cost was recorded as %#v", env.Usage.Cost)
	}
	cost := 0.25
	addGenerationUsage(env, "generator-primary", GenerationResult{Cost: &cost})
	if got := env.Usage.Cost["generator-primary"]; got != cost {
		t.Fatalf("provider cost = %v, want %v", got, cost)
	}
	if env.Usage.InputTokens != 2 || env.Usage.OutputTokens != 3 {
		t.Fatalf("usage = %#v", env.Usage)
	}
}

func TestRuntimeRegistryCoversCompilerDefinitionsAndMetadata(t *testing.T) {
	runtime := NativeRegistry()
	for _, definition := range ragcompiler.BuiltinRegistry().Definitions() {
		if definition.Recipe {
			continue
		}
		if len(definition.Capabilities) == 0 || len(definition.Resources) == 0 || len(definition.ObservationSchemas) == 0 {
			t.Fatalf("metadata=%#v", definition)
		}
		if _, ok := runtime.Lookup(definition.Ref); !ok {
			t.Fatalf("missing runtime %s", definition.Ref.ID())
		}
	}
}
func TestRegistryRejectsDuplicate(t *testing.T) {
	r := NewRegistry()
	op := unitOperator{"units.identity"}
	if err := r.Register(op); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(op); err == nil || !strings.Contains(err.Error(), "RAG_OPERATOR_DUPLICATE") {
		t.Fatalf("%v", err)
	}
}
func TestAgentsViewRunsAndUnicodeChunkRanges(t *testing.T) {
	corpus := Corpus{Records: []SourceRecord{{ID: "a1", SessionID: "s", Ordinal: 2, Role: "assistant", Text: "βeta"}, {ID: "u1", SessionID: "s", Ordinal: 1, Role: "user", Text: "héllo 👋"}, {ID: "a2", SessionID: "s", Ordinal: 3, Role: "assistant", Text: "世界"}}}
	out, err := (unitOperator{"transcript.units.agents-view-runs"}).Execute(context.Background(), ragcontract.Node{}, map[string]any{"corpus": corpus}, nil)
	if err != nil {
		t.Fatal(err)
	}
	units := out["units"].([]Unit)
	if len(units) != 2 || len(units[1].Records) != 2 {
		t.Fatalf("units=%#v", units)
	}
	node := ragcontract.Node{Config: json.RawMessage(`{"size":3,"overlap":1}`)}
	chunked, err := (chunkOperator{}).Execute(context.Background(), node, map[string]any{"units": units}, nil)
	if err != nil {
		t.Fatal(err)
	}
	chunks := chunked["chunks"].([]Chunk)
	if err := validateUTF8Ranges(chunks); err != nil {
		t.Fatal(err)
	}
	bySource := map[string]string{}
	for _, record := range corpus.Records {
		bySource[record.ID] = record.Text
	}
	for _, chunk := range chunks {
		if chunk.Record.LogicalEnd-chunk.Record.LogicalStart > 3 {
			t.Fatalf("range=%#v", chunk.Record)
		}
		rebuilt := ""
		for _, sourceRange := range chunk.Ranges {
			text := bySource[sourceRange.SourceID]
			if sourceRange.ByteStart < 0 || sourceRange.ByteEnd > int64(len(text)) || sourceRange.ByteStart >= sourceRange.ByteEnd {
				t.Fatalf("source range=%#v", sourceRange)
			}
			rebuilt += text[sourceRange.ByteStart:sourceRange.ByteEnd]
		}
		if rebuilt != chunk.Text {
			t.Fatalf("rebuilt=%q chunk=%q ranges=%#v", rebuilt, chunk.Text, chunk.Ranges)
		}
	}
}
func TestRepresentationsDerivationMultiplicityCacheAndFailure(t *testing.T) {
	chunk := fixtureChunk("c", "u", "source")
	generator := &fakeGenerator{}
	env := &Environment{Manifests: fixtureResolver(), Schemas: fakeSchemaValidator{}, Generator: generator, Cache: NewMemoryCache(), Usage: Usage{Cost: map[string]float64{}}}
	summaryNode := ragcontract.Node{Config: json.RawMessage(`{"name":"summary","model":"m","prompt":"p","outputSchema":"summary/v1"}`)}
	summary, err := (representationOperator{"representations.structured-summary"}).Execute(context.Background(), summaryNode, map[string]any{"chunks": []Chunk{chunk}}, env)
	if err != nil {
		t.Fatal(err)
	}
	summaries := summary["representations"].([]Representation)
	if summaries[0].Record.EvidenceRole != "derived" || summaries[0].Record.Derivation == nil {
		t.Fatalf("%#v", summaries[0])
	}
	_, err = (representationOperator{"representations.structured-summary"}).Execute(context.Background(), summaryNode, map[string]any{"chunks": []Chunk{chunk}}, env)
	if err != nil {
		t.Fatal(err)
	}
	if generator.calls != 1 {
		t.Fatalf("cache calls=%d", generator.calls)
	}
	questionNode := ragcontract.Node{Config: json.RawMessage(`{"name":"question","from":"summary","count":3,"model":"m","prompt":"p"}`)}
	questions, err := (representationOperator{"representations.synthetic-questions"}).Execute(context.Background(), questionNode, map[string]any{"chunks": []Chunk{chunk}, "source": summaries}, env)
	if err != nil {
		t.Fatal(err)
	}
	generatedQuestions := questions["representations"].([]Representation)
	if len(generatedQuestions) != 3 || generatedQuestions[0].Record.ID == generatedQuestions[1].Record.ID {
		t.Fatal("question multiplicity and identity")
	}
	merged, err := (mergeOperator{}).Execute(context.Background(), ragcontract.Node{ID: "merge", Operator: (mergeOperator{}).Ref(), Config: json.RawMessage(`{}`)}, map[string]any{"summary": summaries, "questions": generatedQuestions}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var mergedManifest ragcontract.RepresentationSetManifest
	if err := json.Unmarshal(merged["artifact"].(Artifact).Metadata, &mergedManifest); err != nil {
		t.Fatal(err)
	}
	if err := ragcontract.ValidateManifestBase(mergedManifest.ManifestBase, ragcontract.RepresentationManifestSchema, true); err != nil {
		t.Fatal(err)
	}
	derivation := generatedQuestions[0].Record.Derivation
	if derivation == nil || len(derivation.ParentRepresentationIDs) != 1 || derivation.ParentRepresentationIDs[0] != summaries[0].Record.ID || len(derivation.SourceRecordIDs) != 1 || derivation.SourceRecordIDs[0] != "source-c" {
		t.Fatalf("lineage=%#v", derivation)
	}
	other := fixtureChunk("other", "u", "new")
	generator.invalid = true
	_, err = (representationOperator{"representations.structured-summary"}).Execute(context.Background(), summaryNode, map[string]any{"chunks": []Chunk{other}}, env)
	if err == nil || !strings.Contains(err.Error(), "RAG_STRUCTURED_OUTPUT_INVALID") {
		t.Fatalf("%v", err)
	}
	generator.invalid = false
	generator.fail = true
	other = fixtureChunk("other-failure", "u", "new failure")
	_, err = (representationOperator{"representations.structured-summary"}).Execute(context.Background(), summaryNode, map[string]any{"chunks": []Chunk{other}}, env)
	if err == nil || !strings.Contains(err.Error(), "RAG_GENERATION_FAILED") {
		t.Fatalf("%v", err)
	}
}
func TestEmbeddingValidation(t *testing.T) {
	reps := []Representation{fixtureRepresentation("r1", "raw", "c1", "u1", "a"), fixtureRepresentation("r2", "raw", "c2", "u2", "bb")}
	node := ragcontract.Node{Config: json.RawMessage(`{"model":"m","dimensions":2,"normalize":"l2"}`)}
	env := &Environment{Manifests: fixtureResolver(), Embedder: fakeEmbedder{}}
	out, err := (embeddingOperator{}).Execute(context.Background(), node, map[string]any{"representations": reps}, env)
	if err != nil {
		t.Fatal(err)
	}
	embeddings := out["embeddings"].([]Embedding)
	artifact := out["artifact"].(Artifact)
	var manifest ragcontract.EmbeddingSetManifest
	if err := json.Unmarshal(artifact.Metadata, &manifest); err != nil {
		t.Fatal(err)
	}
	if err := ragcontract.ValidateManifestBase(manifest.ManifestBase, ragcontract.EmbeddingManifestSchema, true); err != nil {
		t.Fatal(err)
	}
	if embeddings[0].Record.RepresentationID != "r1" || math.Abs(embeddings[0].Vector[0]*embeddings[0].Vector[0]+embeddings[0].Vector[1]*embeddings[0].Vector[1]-1) > 1e-9 {
		t.Fatalf("%#v", embeddings)
	}
	env.Embedder = fakeEmbedder{badDimension: true}
	if _, err := (embeddingOperator{}).Execute(context.Background(), node, map[string]any{"representations": reps}, env); err == nil {
		t.Fatal("accepted dimensions")
	}
}
func TestIndexRepresentationIsolationRetrievalCollapseFusionHydration(t *testing.T) {
	chunk1 := fixtureChunk("c1", "u1", "alpha source")
	chunk2 := fixtureChunk("c2", "u2", "beta source")
	raw := fixtureRepresentation("raw1", "raw", "c1", "u1", "alpha")
	question1 := fixtureRepresentation("q1", "question", "c2", "u2", "alpha question")
	question2 := fixtureRepresentation("q2", "question", "c2", "u2", "alpha duplicate")
	question1.Record.Citation = chunk2.Record.Citation
	question2.Record.Citation = chunk2.Record.Citation
	indexOut, err := (indexOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{}`)}, map[string]any{"representations.raw": []Representation{raw}, "representations.question": []Representation{question1, question2}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	index := indexOut["index"].(*MultiIndex)
	defer func() { _ = index.Close() }()
	if err := ragcontract.ValidateManifestBase(index.Manifest.ManifestBase, ragcontract.IndexManifestSchema, true); err != nil {
		t.Fatal(err)
	}
	env := &Environment{Trace: &ragcontract.QueryTrace{}}
	retrieve := ragcontract.Node{ID: "question.lexical", Config: json.RawMessage(`{"representation":"question","topK":10,"filter":{}}`)}
	hitsOut, err := (retrieveOperator{"retrieve.bm25"}).Execute(context.Background(), retrieve, map[string]any{"index": index, "query": Query{ID: "q", Text: "alpha"}}, env)
	if err != nil {
		t.Fatal(err)
	}
	hits := hitsOut["hits"].([]RankedRecord)
	if len(hits) != 2 {
		t.Fatalf("hits=%#v", hits)
	}
	for _, hit := range hits {
		if hit.Representation.Record.Kind != "question" {
			t.Fatal("representation isolation")
		}
	}
	collapsed, err := (collapseOperator{"collapse.parent"}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{"scope":"unit","representative":"scoreThenRepresentationId"}`)}, map[string]any{"hits": hits}, env)
	if err != nil {
		t.Fatal(err)
	}
	parents := collapsed["parents"].([]RankedParent)
	if len(parents) != 1 || len(parents[0].Members) != 2 {
		t.Fatalf("parents=%#v", parents)
	}
	rawParents := []RankedParent{{Rank: 1, Identity: ragcontract.CollapseIdentity{Scope: "unit", ID: "u1"}, Score: 1, Representative: raw}}
	fused, err := (fusionOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{"rankConstant":60,"weights":{"raw":2,"question":1}}`)}, map[string]any{"channel.raw": rawParents, "channel.question": parents}, env)
	if err != nil {
		t.Fatal(err)
	}
	fusedParents := fused["parents"].([]RankedParent)
	if len(fusedParents) != 2 || len(fusedParents[0].Contributions) == 0 {
		t.Fatalf("fused=%#v", fusedParents)
	}
	hydrated, err := (hydrateOperator{}).Execute(context.Background(), ragcontract.Node{}, map[string]any{"parents": fusedParents, "chunks": []Chunk{chunk1, chunk2}}, env)
	if err != nil {
		t.Fatal(err)
	}
	evidence := hydrated["evidence"].([]Evidence)
	for _, item := range evidence {
		if item.Chunk.Text == question1.Text || item.Chunk.Text == question2.Text {
			t.Fatal("derived text hydrated as evidence")
		}
	}
}
func TestVectorRetrievalIsStableAndRepresentationIsolated(t *testing.T) {
	raw := fixtureRepresentation("raw", "raw", "c1", "u1", "alpha")
	question := fixtureRepresentation("question", "question", "c2", "u2", "alpha")
	model, err := fixtureResolver().Model("m")
	if err != nil {
		t.Fatal(err)
	}
	embeddings := []Embedding{
		{Record: ragcontract.EmbeddingRecord{RepresentationID: raw.Record.ID, ModelManifestDigest: model.Digest, Dimensions: 2}, Vector: []float64{1, 0}},
		{Record: ragcontract.EmbeddingRecord{RepresentationID: question.Record.ID, ModelManifestDigest: model.Digest, Dimensions: 2}, Vector: []float64{1, 0}},
	}
	indexed, err := (indexOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{}`)}, map[string]any{"representations.all": []Representation{raw, question}, "embeddings": embeddings}, nil)
	if err != nil {
		t.Fatal(err)
	}
	index := indexed["index"].(*MultiIndex)
	defer func() { _ = index.Close() }()
	node := ragcontract.Node{ID: "raw.vector", Config: json.RawMessage(`{"representation":"raw","topK":10,"filter":{}}`)}
	out, err := (retrieveOperator{"retrieve.vector"}).Execute(context.Background(), node, map[string]any{"index": index, "query": Query{ID: "q", Text: "x"}}, &Environment{Manifests: fixtureResolver(), Embedder: fakeEmbedder{}})
	if err != nil {
		t.Fatal(err)
	}
	hits := out["hits"].([]RankedRecord)
	if len(hits) != 1 || hits[0].Representation.Record.Kind != "raw" {
		t.Fatalf("hits=%#v", hits)
	}
}

func TestFusionTieMissingChannelAndHydrationLineageFailures(t *testing.T) {
	rep1 := fixtureRepresentation("r1", "raw", "c1", "u1", "x")
	rep2 := fixtureRepresentation("r2", "raw", "c2", "u2", "x")
	parents := []RankedParent{{Rank: 1, Identity: ragcontract.CollapseIdentity{Scope: "unit", ID: "b"}, Representative: rep1}, {Rank: 1, Identity: ragcontract.CollapseIdentity{Scope: "unit", ID: "a"}, Representative: rep2}}
	out, err := (fusionOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{"rankConstant":60}`)}, map[string]any{"channel.x": parents}, &Environment{})
	if err != nil {
		t.Fatal(err)
	}
	got := out["parents"].([]RankedParent)
	if got[0].Identity.ID != "a" {
		t.Fatalf("tie=%#v", got)
	}
	_, err = (fusionOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{"rankConstant":60,"weights":{"missing":1},"missingChannelPolicy":"reject"}`)}, map[string]any{"channel.x": parents}, &Environment{})
	if err == nil {
		t.Fatal("missing channel accepted")
	}
	_, err = (hydrateOperator{}).Execute(context.Background(), ragcontract.Node{}, map[string]any{"parents": got, "chunks": []Chunk{}}, &Environment{})
	if err == nil || !strings.Contains(err.Error(), "RAG_HYDRATE_LINEAGE") {
		t.Fatalf("%v", err)
	}
}
func TestGenerationContractRejectsQuestionCountMismatchAndPermitsAbstention(t *testing.T) {
	chunk := fixtureChunk("c1", "u1", "source")
	env := &Environment{Manifests: fixtureResolver(), Schemas: fakeSchemaValidator{}, Generator: &fakeGenerator{questions: []string{"only one"}}, Usage: Usage{Cost: map[string]float64{}}}
	questionNode := ragcontract.Node{Config: json.RawMessage(`{"name":"questions","count":2,"model":"m","prompt":"p"}`)}
	if _, err := (representationOperator{"representations.synthetic-questions"}).Execute(context.Background(), questionNode, map[string]any{"chunks": []Chunk{chunk}}, env); err == nil || !strings.Contains(err.Error(), "RAG_QUESTION_COUNT_MISMATCH") {
		t.Fatalf("question count error = %v", err)
	}
	evidence := []Evidence{{Rank: 1, Chunk: chunk, Score: 1}}
	answerNode := ragcontract.Node{Config: json.RawMessage(`{"model":"m","prompt":"p","citations":"source"}`)}
	answerEnv := &Environment{Manifests: fixtureResolver(), Generator: &fakeGenerator{abstained: true}, QueryText: "Unknown question", Usage: Usage{Cost: map[string]float64{}}}
	out, err := (answerOperator{}).Execute(context.Background(), answerNode, map[string]any{"evidence": evidence}, answerEnv)
	if err != nil {
		t.Fatalf("abstained answer error = %v", err)
	}
	if !out["answer"].(Answer).Abstained {
		t.Fatalf("answer = %#v, want abstention", out)
	}
}

func TestRerankAndAnswerNeverFallback(t *testing.T) {
	chunk := fixtureChunk("c1", "u1", "source")
	evidence := []Evidence{{Rank: 1, Collapse: ragcontract.CollapseIdentity{Scope: "unit", ID: "u1"}, Chunk: chunk, Score: 1}}
	rerankNode := ragcontract.Node{Config: json.RawMessage(`{"model":"m","candidateCount":1,"results":1,"truncation":"tail","tokenization":"exact"}`)}
	if _, err := (rerankOperator{}).Execute(context.Background(), rerankNode, map[string]any{"evidence": evidence}, &Environment{}); err == nil || !strings.Contains(err.Error(), "UNAVAILABLE") {
		t.Fatalf("%v", err)
	}
	if _, err := (rerankOperator{}).Execute(context.Background(), rerankNode, map[string]any{"evidence": evidence}, &Environment{Manifests: fixtureResolver(), Reranker: fakeReranker{incomplete: true}}); err == nil || !strings.Contains(err.Error(), "INCOMPLETE") {
		t.Fatalf("%v", err)
	}
	answerNode := ragcontract.Node{Config: json.RawMessage(`{"model":"m","prompt":"p","citations":"required"}`)}
	generator := &fakeGenerator{}
	out, err := (answerOperator{}).Execute(context.Background(), answerNode, map[string]any{"evidence": evidence}, &Environment{Manifests: fixtureResolver(), Generator: generator, QueryText: "What does this source say?"})
	if err != nil {
		t.Fatal(err)
	}
	if out["answer"].(Answer).CitationChunkIDs[0] != "c1" {
		t.Fatalf("%#v", out)
	}
	if generator.request.Kind != "generate.answer" || generator.request.Text != "What does this source say?" || generator.request.OutputSchema != "summary/v1" {
		t.Fatalf("answer request = %#v", generator.request)
	}
}
func TestVersionedEvaluationMetrics(t *testing.T) {
	q := Query{RelevantIDs: []string{"u1"}}
	e := []Evidence{{Rank: 1, Collapse: ragcontract.CollapseIdentity{ID: "u1"}}}
	measures := []ragcontract.Measure{{Name: "rag.precision", Unit: "ratio", Config: json.RawMessage(`{"cutoffs":[1]}`)}, {Name: "rag.recall", Unit: "ratio", Config: json.RawMessage(`{"cutoffs":[1]}`)}, {Name: "rag.hit-rate", Unit: "ratio", Config: json.RawMessage(`{"cutoffs":[1]}`)}, {Name: "rag.mrr", Unit: "ratio", Config: json.RawMessage(`{}`)}, {Name: "rag.ndcg", Unit: "ratio", Config: json.RawMessage(`{"cutoffs":[1]}`)}, {Name: "rag.latency", Config: json.RawMessage(`{"stages":["query"]}`)}, {Name: "rag.token-usage", Config: json.RawMessage(`{}`)}, {Name: "rag.provider-cost", Config: json.RawMessage(`{}`)}, {Name: "rag.storage-bytes", Config: json.RawMessage(`{}`)}, {Name: "rag.failure-rates", Config: json.RawMessage(`{}`)}, {Name: "rag.abstention", Config: json.RawMessage(`{}`)}}
	metrics := Evaluate(q, e, nil, measures, map[string]int64{"query": 2}, Usage{InputTokens: 1, Cost: map[string]float64{"m": 0.1}}, nil, 12)
	if len(metrics) != len(measures) {
		t.Fatalf("metrics=%d", len(metrics))
	}
	for _, metric := range metrics {
		if !json.Valid(metric.Value) {
			t.Fatalf("invalid %s", metric.Name)
		}
	}
}
func fixtureResolver() StaticManifestResolver {
	modelDigest := "sha256:" + strings.Repeat("a", 64)
	promptDigest := "sha256:" + strings.Repeat("b", 64)
	return StaticManifestResolver{
		Models:  map[string]ragcontract.ModelManifest{"m": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: modelDigest}, ModelID: "model-exact", Tokenization: "exact", Truncation: "tail"}},
		Prompts: map[string]ragcontract.PromptManifest{"p": {ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: promptDigest}, PromptID: "prompt-exact", OutputSchema: "summary/v1"}},
	}
}
func fixtureChunk(id, unit, text string) Chunk {
	digest, _ := ragcontract.Digest(text)
	return Chunk{Record: ragcontract.ChunkRecord{ID: id, ParentUnitID: unit, TextDigest: digest, Citation: ragcontract.CitationRef{SourceID: "source-" + id}}, Text: text, ManifestDigest: "sha256:" + strings.Repeat("c", 64)}
}
func fixtureRepresentation(id, kind, chunk, unit, text string) Representation {
	digest, _ := ragcontract.Digest(text)
	return Representation{Record: ragcontract.RepresentationRecord{ID: id, Kind: kind, ParentChunkID: chunk, ParentUnitID: unit, ContentDigest: digest, EvidenceRole: "source", Citation: ragcontract.CitationRef{SourceID: "source-" + chunk}}, Text: text, ManifestDigest: "sha256:" + strings.Repeat("d", 64)}
}
func FuzzChunkRangesNeverSplitUTF8(f *testing.F) {
	f.Add("héllo 世界 👋", 4, 1)
	f.Fuzz(func(t *testing.T, text string, size, overlap int) {
		if !json.Valid([]byte(`"` + strings.ReplaceAll(text, `"`, `\"`) + `"`)) {
			return
		}
		if size < 1 || size > 100 || overlap < 0 || overlap >= size {
			return
		}
		unit := Unit{Record: ragcontract.UnitRecord{ID: "u"}, Text: text, Records: []SourceRecord{{ID: "s", Text: text}}}
		out, err := (chunkOperator{}).Execute(context.Background(), ragcontract.Node{Config: json.RawMessage(`{"size":` + itoa(size) + `,"overlap":` + itoa(overlap) + `}`)}, map[string]any{"units": []Unit{unit}}, nil)
		if err != nil {
			return
		}
		if err := validateUTF8Ranges(out["chunks"].([]Chunk)); err != nil {
			t.Fatal(err)
		}
	})
}
