package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragmodel"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragproduct"
)

func TestHTTPReferenceHost(t *testing.T) {
	corpus := ragoperators.Corpus{SchemaVersion: "rag-corpus-data/v1", Records: []ragoperators.SourceRecord{{ID: "s1", SessionID: "s", Ordinal: 1, Role: "user", Text: "reciprocal rank fusion"}}}
	artifact := ragoperators.NewCorpusArtifact(corpus, "http-test")
	pipeline := ragmodel.NewPipeline("p", func(p *ragmodel.PipelineBuilder) {
		p.CorpusInput(ragmodel.Corpus("corpus")).Units(ragmodel.UnitsIdentity()).Chunks(ragmodel.RecursiveChunks(ragmodel.RecursiveChunkConfig{MaxRunes: 100})).Represent(ragmodel.RawRepresentation("raw")).IndexNamed("idx", ragmodel.BleveMulti(ragmodel.BleveMultiConfig{Lexical: true}))
	})
	query := ragmodel.NewQueryPlan("q", func(q *ragmodel.QueryBuilder) {
		q.Channels(ragmodel.BM25("lexical", ragmodel.RetrieveConfig{Index: "idx", Representation: "raw", TopK: 5})).CollapseChannels(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "scoreThenRepresentationId"})).Fuse(ragmodel.WeightedRRF(ragmodel.WeightedRRFConfig{RankConstant: 60})).CollapseFinal(ragmodel.ParentCollapse(ragmodel.CollapseConfig{Scope: "unit", Representative: "bestFusionContributionThenId"})).Hydrate(ragmodel.SourceEvidence(ragmodel.HydrationConfig{Selection: "bestContributionThenId"}))
	})
	product := ragmodel.NewProduct("http", func(p *ragmodel.ProductBuilder) {
		p.PipelineValue(pipeline).QueryPlan(query).RequestContract(func(r *ragmodel.RequestBuilder) { r.Field("query", "string", true, 128) }).ResponseContract(func(r *ragmodel.ResponseBuilder) { r.Citations("required").TraceID(true) }).RuntimePolicy(func(r *ragmodel.RuntimeBuilder) { r.Concurrent(2).Trace("metadata-only").ProviderFailure("fail") })
	})
	plan, err := ragmodel.CompileProduct(product, ragmodel.CompileOptions{Inputs: map[string]ragcontract.ArtifactBinding{"corpus": {SlotID: "corpus", Role: "corpus", Kind: "json", Digest: artifact.Manifest.Digest, SchemaVersion: ragcontract.CorpusManifestSchema}}})
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := ragproduct.New(context.Background(), plan, ragproduct.Bindings{Corpus: artifact, Cache: ragoperators.NewMemoryCache()})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = runtime.Close() }()
	server := httptest.NewServer(newHandler(runtime))
	defer server.Close()
	requestBody := []byte(`{"values":{"query":"fusion"}}`)
	request, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/query", bytes.NewReader(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", response.StatusCode)
	}
	var value ragproduct.Response
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}
	if len(value.Citations) == 0 || value.TraceID == "" {
		t.Fatalf("response=%#v", value)
	}
	bad, err := http.Post(server.URL+"/v1/query", "text/plain", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	defer bad.Body.Close()
	if bad.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("bad status=%d", bad.StatusCode)
	}
}
