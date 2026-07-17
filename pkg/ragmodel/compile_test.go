package ragmodel

import (
	"strings"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestPureGoStudyResolvesFactorReferencesAndProducesTwoCells(t *testing.T) {
	pipeline := testPipeline()
	query := NewQueryPlan("q", func(q *QueryBuilder) {
		q.Channels(BM25("raw.lexical", RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10})).CollapseChannels(ParentCollapse(CollapseConfig{Scope: FactorRef{ID: "collapse"}, Representative: "scoreThenRepresentationId"})).Fuse(WeightedRRF(WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ParentCollapse(CollapseConfig{Scope: FactorRef{ID: "collapse"}, Representative: "bestFusionContributionThenId"})).Hydrate(SourceEvidence(HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	study := NewStudy("s", func(s *StudyBuilder) {
		s.PipelineValue(pipeline).DatasetRef(DatasetArtifact("judgments", "smoke", "candidate", "unit")).VariantsList(func(v *VariantsBuilder) {
			v.Add("raw", func(x *VariantBuilder) { x.SelectRepresentations("raw").QueryPlan(query) })
		}).FactorsList(func(f *FactorsBuilder) { f.Enum("collapse", "chunk", "unit") }).MetricsList(func(m *MetricsBuilder) { m.MRR() })
	})
	bindings := map[string]ragcontract.ArtifactBinding{"corpus": binding("corpus", ragcontract.CorpusManifestSchema, "a"), "judgments": binding("judgments", ragcontract.EvaluationManifestSchema, "b")}
	_, cells, err := CompileStudy(study, CompileOptions{Inputs: bindings})
	if err != nil {
		t.Fatal(err)
	}
	if len(cells) != 2 {
		t.Fatalf("cells=%d", len(cells))
	}
	for _, cell := range cells {
		for _, node := range cell.Execution.Pipeline.Nodes {
			if strings.Contains(string(node.Config), "$factor") {
				t.Fatalf("unresolved factor in %s", node.Config)
			}
		}
	}
}
func TestPureGoCompilationRejectsMissingInputsAndTargets(t *testing.T) {
	if _, err := CompileProduct(&Product{}, CompileOptions{}); err == nil {
		t.Fatal("accepted empty product")
	}
	pipeline := testPipeline()
	query := NewQueryPlan("q", func(q *QueryBuilder) {
		q.Channels(BM25("raw", RetrieveConfig{Index: "i", Representation: "raw", TopK: 1})).CollapseChannels(ParentCollapse(CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(WeightedRRF(WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ParentCollapse(CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(SourceEvidence(HydrationConfig{Selection: "bestContributionThenId"}))
	})
	product := NewProduct("p", func(p *ProductBuilder) { p.PipelineValue(pipeline).QueryPlan(query) })
	if _, err := CompileProduct(product, CompileOptions{}); err == nil || !strings.Contains(err.Error(), "RAG_V2_INPUT_BINDING") {
		t.Fatalf("error=%v", err)
	}
}
func testPipeline() *Pipeline {
	return NewPipeline("p", func(p *PipelineBuilder) {
		p.CorpusInput(Corpus("corpus")).Units(UnitsIdentity()).Chunks(RecursiveChunks(RecursiveChunkConfig{MaxRunes: 100})).Represent(RawRepresentation("raw")).IndexNamed("representations", BleveMulti(BleveMultiConfig{Lexical: true}))
	})
}
func binding(role, schema, letter string) ragcontract.ArtifactBinding {
	return ragcontract.ArtifactBinding{Role: role, Kind: "manifest", Digest: "sha256:" + strings.Repeat(letter, 64), SchemaVersion: schema}
}
