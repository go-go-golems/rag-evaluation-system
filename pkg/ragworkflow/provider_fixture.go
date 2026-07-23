package ragworkflow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type deterministicProviderFixture struct{ ragoperators.FixtureProviders }

func (p deterministicProviderFixture) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	if request.Kind != "generate.answer" {
		return p.FixtureProviders.Generate(ctx, request)
	}
	if len(request.Evidence) == 0 {
		return ragoperators.GenerationResult{}, fmt.Errorf("RAG_PROVIDER_FIXTURE_ANSWER_EVIDENCE")
	}
	cost := 0.000004
	return ragoperators.GenerationResult{Text: "fixture grounded answer", CitationChunkIDs: []string{request.Evidence[0].Chunk.Record.ID}, InputTokens: 11, OutputTokens: 4, Cost: &cost, FinishReason: "fixture"}, nil
}
func (deterministicProviderFixture) Rerank(_ context.Context, request ragoperators.RerankRequest) (ragoperators.RerankResult, error) {
	scores := make([]ragoperators.RerankScore, len(request.Candidates))
	for index, candidate := range request.Candidates {
		scores[index] = ragoperators.RerankScore{ChunkID: candidate.Chunk.Record.ID, Score: float64(len(scores) - index)}
	}
	cost := 0.000002
	return ragoperators.RerankResult{Scores: scores, InputTokens: 7, Cost: &cost}, nil
}

func NewDeterministicProviderFixture() (Fixture, error) {
	fixture, err := NewProviderFreeFixture(true)
	if err != nil {
		return Fixture{}, err
	}
	for index, node := range fixture.Execution.Pipeline.Nodes {
		switch node.Operator.Kind {
		case "representations.raw":
			fixture.Execution.Pipeline.Nodes[index].Operator = ragcontract.OperatorRef{Kind: "representations.structured-summary", Version: "v1"}
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"name":"summary","model":"fixture-summary-v1","prompt":"fixture-transcript-summary-v1","outputSchema":"transcript-rag-summary/v1"}`)
		case "embed.model":
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"model":"fixture-hash-32-v1","dimensions":32,"distance":"cosine","normalize":"l2","batchSize":8}`)
		case "retrieve.bm25", "retrieve.vector":
			fixture.Execution.Pipeline.Nodes[index].Config = json.RawMessage(`{"index":"representations","representation":"summary","topK":10,"filter":{}}`)
		}
	}
	hydrationID := ""
	for _, node := range fixture.Execution.Pipeline.Nodes {
		if node.Operator.Kind == "hydrate.source-evidence" {
			hydrationID = node.ID
		}
	}
	if hydrationID == "" {
		return Fixture{}, fmt.Errorf("RAG_PROVIDER_FIXTURE_HYDRATION")
	}
	fixture.Execution.Pipeline.Nodes = append(fixture.Execution.Pipeline.Nodes,
		ragcontract.Node{ID: "fixture-rerank", Operator: ragcontract.OperatorRef{Kind: "rerank.cross-encoder", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "evidence", From: ragcontract.PortRef{NodeID: hydrationID, Port: "evidence"}}}, Config: json.RawMessage(`{"model":"fixture-rerank-v1","candidateCount":20,"results":5,"inputTemplate":"query-document","truncation":"none","tokenization":"fixture-utf16","timeoutMilliseconds":5000}`)},
		ragcontract.Node{ID: "fixture-answer", Operator: ragcontract.OperatorRef{Kind: "generate.answer", Version: "v1"}, Inputs: []ragcontract.InputBinding{{Port: "evidence", From: ragcontract.PortRef{NodeID: "fixture-rerank", Port: "evidence"}}}, Config: json.RawMessage(`{"model":"fixture-summary-v1","prompt":"fixture-answer-v1","citations":"required","citationFailurePolicy":"abstain","contextBudgetTokens":2048}`)},
	)
	fixture.Execution.Pipeline.Outputs[0].From = ragcontract.PortRef{NodeID: "fixture-rerank", Port: "evidence"}
	fixture.Execution.Pipeline, err = ragcompiler.Normalize(fixture.Execution.Pipeline, nil)
	if err != nil {
		return Fixture{}, err
	}
	fixture.Execution.CellID = ""
	fixture.Execution.CellID, err = ragcontract.Digest(fixture.Execution)
	if err != nil {
		return Fixture{}, err
	}
	return fixture, nil
}

func NewDeterministicProviderServices() (ProviderServices, error) {
	providers := deterministicProviderFixture{FixtureProviders: ragoperators.NewFixtureProviders()}
	providers.Resolver.Models["fixture-rerank-v1"] = ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: "sha256:" + strings.Repeat("6", 64)}, ModelID: "fixture-rerank-v1", ModelDigest: "sha256:" + strings.Repeat("6", 64), Tokenization: "fixture-utf16", Truncation: "none", Normalization: "none", ImplementationVersion: "fixture/v1", RequestParameters: json.RawMessage(`{}`)}
	providers.Resolver.Prompts["fixture-answer-v1"] = ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: "sha256:" + strings.Repeat("7", 64)}, PromptID: "fixture-answer-v1", TemplateDigest: "sha256:" + strings.Repeat("7", 64), InputSchema: "text/plain", OutputSchema: "fixture-answer/v1"}
	authority := ProviderAuthority{SchemaVersion: ProviderAuthoritySchema, ProfileID: "fixture-geppetto-v1", FixtureProviders: true, Capabilities: []string{"embedder", "generator", "reranker", "schema-validator"}, ModelManifestDigests: []string{"sha256:" + strings.Repeat("1", 64), "sha256:" + strings.Repeat("2", 64), "sha256:" + strings.Repeat("3", 64), "sha256:" + strings.Repeat("6", 64)}, PromptManifestDigests: []string{"sha256:" + strings.Repeat("4", 64), "sha256:" + strings.Repeat("5", 64), "sha256:" + strings.Repeat("7", 64)}, Providers: []WorkflowProviderIdentity{{Role: "embedding-primary", ProfileSlug: "fixture-embedding", ModelManifestDigest: "sha256:" + strings.Repeat("3", 64), ModelID: ragoperators.FixtureEmbeddingModel, SettingsFingerprint: "sha256:" + strings.Repeat("a", 64), ConcurrencyLimit: 1}, {Role: "generator-primary", ProfileSlug: "fixture-generation", ModelManifestDigest: "sha256:" + strings.Repeat("1", 64), ModelID: ragoperators.FixtureSummaryModel, SettingsFingerprint: "sha256:" + strings.Repeat("b", 64), ConcurrencyLimit: 1}, {Role: "reranker-primary", ProfileSlug: "fixture-rerank", ModelManifestDigest: "sha256:" + strings.Repeat("6", 64), ModelID: "fixture-rerank-v1", SettingsFingerprint: "sha256:" + strings.Repeat("c", 64), ConcurrencyLimit: 1}}}
	canonicalizeAuthority(&authority)
	digest, err := providerAuthorityDigest(authority)
	if err != nil {
		return ProviderServices{}, err
	}
	authority.Digest = digest
	if err := ValidateProviderAuthority(authority); err != nil {
		return ProviderServices{}, err
	}
	return ProviderServices{Authority: authority, Manifests: providers.Resolver, Schemas: providers, Generator: providers, Embedder: providers, Reranker: providers, GenerationConcurrency: 1}, nil
}
