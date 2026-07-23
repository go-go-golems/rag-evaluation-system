package ragworkflow

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type Fixture struct {
	Execution ragcontract.PipelineExecution  `json:"execution"`
	Corpus    ragoperators.Corpus            `json:"corpus"`
	Dataset   ragoperators.EvaluationDataset `json:"dataset"`
}

func NewProviderFreeFixture(vector bool) (Fixture, error) {
	pipeline := ragmodel.NewPipeline("pipeline", func(p *ragmodel.PipelineBuilder) {
		builder := p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 200})).Represent(ragmodel.RawRepresentation("raw"))
		if vector {
			builder.EmbeddingModel(ragmodel.EmbeddingModel("fixture-embedding-v1", ragmodel.EmbeddingConfig{Dimensions: 8, Distance: "cosine", Normalize: "l2", BatchSize: 4}))
		}
		index := ragmodel.BleveMultiConfig{Lexical: true}
		if vector {
			index.Vector = &ragmodel.VectorIndexConfig{Distance: "cosine"}
		}
		builder.IndexNamed("representations", ragmodel.BleveMulti(index))
	})
	query := ragmodel.NewQueryPlan("query", func(q *ragmodel.QueryBuilder) {
		if vector {
			q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10}), ragmodel.Vector("raw.vector", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10}))
		} else {
			q.Channels(ragmodel.BM25("raw.lexical", ragmodel.RetrieveConfig{Index: "representations", Representation: "raw", TopK: 10}))
		}
		q.CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"})).ResultCount(5)
	})
	product := ragmodel.NewProduct("fixture", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Citations("source") })
	})
	size := int64(1)
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {Role: "corpus", Kind: "json", Digest: "sha256:" + strings.Repeat("a", 64), SizeBytes: &size, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		return Fixture{}, err
	}
	execution := ragcontract.PipelineExecution{SchemaVersion: ragcontract.ExecutionSchemaVersion, Pipeline: plan.Pipeline, Bindings: plan.Bindings, Dataset: ragcontract.DatasetBinding{ManifestDigest: "sha256:" + strings.Repeat("b", 64), Split: "smoke", Status: "candidate", RelevanceTarget: "unit"}, Measures: []ragcontract.Measure{{Name: "rag.mrr", Version: "v1", ValueKind: "number", Unit: "ratio", Required: true, Config: json.RawMessage(`{}`)}}, VariantID: "fixture", Factors: []ragcontract.FactorSelection{}}
	corpus := ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "s1", SessionID: "session", Ordinal: 1, Role: "user", Text: "weighted reciprocal rank fusion decision"}, {ID: "s2", SessionID: "session", Ordinal: 2, Role: "assistant", Text: "vector retrieval uses deterministic embeddings"}, {ID: "s3", SessionID: "session", Ordinal: 3, Role: "assistant", Text: "unrelated material"}}}
	dataset := ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-data/v1", Queries: []ragoperators.Query{{ID: "q1", Text: "reciprocal rank fusion", RelevantIDs: []string{fixtureUnitID(corpus.Records[0])}}, {ID: "q2", Text: "deterministic embeddings", RelevantIDs: []string{fixtureUnitID(corpus.Records[1])}}}}
	corpusDigest, _ := ragcontract.Digest(corpus)
	corpusBody, _ := ragcontract.CanonicalJSON(corpus)
	corpusSize := int64(len(corpusBody) + 1)
	datasetDigest, _ := ragcontract.Digest(dataset)
	execution.Bindings[0].Digest = corpusDigest
	execution.Bindings[0].SizeBytes = &corpusSize
	execution.Dataset.ManifestDigest = datasetDigest
	execution.CellID, err = ragcontract.Digest(execution)
	if err != nil {
		return Fixture{}, err
	}
	return Fixture{Execution: execution, Corpus: corpus, Dataset: dataset}, nil
}

func fixtureUnitID(record ragoperators.SourceRecord) string {
	textDigest, _ := ragcontract.Digest(record.Text)
	digest, _ := ragcontract.Digest(struct {
		Kind   string
		IDs    []string
		Digest string
	}{"units.identity", []string{record.ID}, textDigest})
	if len(digest) < 23 {
		panic(fmt.Sprintf("invalid fixture digest %q", digest))
	}
	return "unit:" + digest[7:23]
}
