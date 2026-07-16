package rag

import "github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"

// TypeScriptModule mirrors the public JavaScript authoring surface. The Go
// builder remains authoritative at runtime.
func TypeScriptModule() *spec.Module {
	return &spec.Module{Name: ModuleName, Description: "Typed fluent immutable RAG laboratory experiments.", RawDTS: []string{
		"export type ArtifactKind = 'corpusSnapshot' | 'chunkSet' | 'embeddingSet' | 'bm25Index' | 'evaluationDataset' | 'representationSet';",
		"export type RelevanceGradeName = '0_FAIL' | '1_PARTIAL' | '2_SUBSTANTIAL' | '3_AUTHORITATIVE';",
		"export interface ArtifactRef { kind: ArtifactKind; id: string; }",
		"export interface ValidationIssue { code: string; path: string; message: string; severity: 'error' | 'warning'; }",
		"export interface ValidationReport { ok: boolean; issues: ValidationIssue[]; }",
		"export interface ChannelBuilder { bm25(): this; vector(): this; representation(name: string): this; topK(count: number): this; filter(configure: (filter: FilterBuilder) => void): this; }",
		"export interface FilterBuilder { sourceIds(ids: string[]): this; documentIds(ids: string[]): this; contentTypes(types: string[]): this; metadataEquals(key: string, value: string): this; }",
		"export interface FusionBuilder { rrf(): this; rankConstant(value: number): this; weight(channel: string, value: number): this; }",
		"export interface RerankingBuilder { crossEncoder(model: string): this; candidates(count: number): this; results(count: number): this; }",
		"export interface RetrievalBuilder { channel(name: string, configure: (channel: ChannelBuilder) => void): this; filter(configure: (filter: FilterBuilder) => void): this; fuse(configure: (fusion: FusionBuilder) => void): this; rerank(configure: (reranker: RerankingBuilder) => void): this; collapse(scope: 'none' | 'parentChunk' | 'document'): this; results(count: number): this; }",
		"export interface MetricsBuilder { relevanceAt(grade: RelevanceGrade): this; precisionAt(cutoffs: number[]): this; recallAt(cutoffs: number[]): this; hitRateAt(cutoffs: number[]): this; ndcgAt(cutoff: number): this; mrr(): this; meanRelevantRecallAt(cutoff: number): this; abstention(): this; }",
		"export interface RepresentationConfig { artifact(ref: string | ArtifactRef): this; parent(value: 'sourceChunk'): this; }",
		"export interface RepresentationsBuilder { rawChunks(name?: string): this; summaries(name: string, configure: (representation: RepresentationConfig) => void): this; questions(name: string, configure: (representation: RepresentationConfig) => void): this; }",
		"export interface Experiment { use(fragment: Fragment): this; corpus(ref: string | ArtifactRef): this; chunks(ref: string | ArtifactRef): this; bm25(ref: string | ArtifactRef): this; embeddings(ref: string | ArtifactRef): this; evaluation(ref: string | ArtifactRef): this; note(text: string): this; tag(name: string, value?: string): this; representations(configure: (builder: RepresentationsBuilder) => void): this; retrieval(configure: (builder: RetrievalBuilder) => void): this; metrics(configure: (builder: MetricsBuilder) => void): this; validate(lab?: Laboratory): ValidationReport; toSpec(): ExperimentSpecification; toJSON(): ExperimentSpecification; }",
		"export interface Fragment { readonly __ragFragment?: never; }",
		"export type QueryEmbed = (query: string) => number[];",
		"export interface LlamaCPPRerankerOptions { kind: 'llama.cpp'; baseURL: string; model: string; maxRequestBytes?: number; }",
		"export interface OpenOptions { database: string; execution?: 'readOnly' | 'allowRuns'; queryEmbed?: QueryEmbed; reranker?: LlamaCPPRerankerOptions; }",
		"export interface ExecutionResult { runId: string; queryCount: number; metrics: Record<string, any>; timing: Record<string, any>; completedAt: string; }",
		"export interface Laboratory { validate(experiment: Experiment): ValidationReport; persist(experiment: Experiment): { id: string; reused: boolean; schemaVersion: string }; start(experiment: Experiment): { id: string; experimentSpecId: string; status: string }; execute(experiment: Experiment): ExecutionResult; close(): void; }",
		"export interface RelevanceGrade { name: RelevanceGradeName; ordinal: number; }",
		"export interface ExperimentSpecification { schemaVersion: 'rag-eval-experiment-spec/v1'; fingerprint: string; name: string; provenance: Record<string, any>; inputs: Record<string, any>; retrieval: Record<string, any>; metrics: Record<string, any>; }",
		"export function open(options: OpenOptions): Laboratory;",
		"export function fragment(name: string, configure: (experiment: Experiment) => void): Fragment;",
		"export function experiment(name: string, configure?: (experiment: Experiment) => void): Experiment;",
		"export function artifact(kind: ArtifactKind, id: string): ArtifactRef;",
		"export function grade(name: RelevanceGradeName): RelevanceGrade;",
		"export const version: 'v1';",
	}}
}
