package rag

import "github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"

func TypeScriptModule() *spec.Module {
	return &spec.Module{Name: ModuleName, Description: "Pure Go-backed composable RAG v2 authoring and compilation.", RawDTS: []string{
		`export type JSONPrimitive = string | number | boolean | null;
export type JSONValue = JSONPrimitive | JSONObject | JSONValue[];
export interface JSONObject { [key: string]: JSONValue; }
export type RepresentationKind = "raw" | "summary" | "question";
export type CollapseScope = "chunk" | "unit";
export type DatasetStatus = "smoke" | "candidate" | "frozen" | "adjudicated";
export type Distance = "cosine" | "dot" | "euclidean";
export interface ArtifactBinding { slotId?: string; role?: string; kind: string; id?: string; uri?: string; digest: string; sizeBytes?: number; schemaVersion: string; }
export interface CompileOptions { inputs: Record<string, ArtifactBinding>; models?: ModelBinding[]; }
export interface ModelBinding { reference: string; manifest: string; digest: string; }
export interface ValidationIssue { code: string; path: string; message: string; severity: "error" | "warning"; operator?: string; hint?: string; }
export interface ValidationReport { ok: boolean; issues: ValidationIssue[]; }
export interface Descriptor<I extends string, O extends string> { readonly __ragDescriptor?: [I, O]; }
export interface Fragment<T> { readonly __ragFragment?: T; }
export interface CorpusInput { readonly __ragCorpus?: true; }
export interface DatasetRef { readonly __ragDataset?: true; }
export interface FactorRef { readonly __ragFactor?: true; }
export type UnitOperator = Descriptor<"corpus", "units">;
export type ChunkOperator = Descriptor<"units", "chunks">;
export type RepresentationOperator = Descriptor<"chunks", "representations">;
export type EmbeddingOperator = Descriptor<"representations", "embeddings">;
export type IndexOperator = Descriptor<"representations", "index">;
export type RetrieverOperator = Descriptor<"index", "ranked-records">;
export type CollapseOperator = Descriptor<"ranked-records", "ranked-parents">;
export type FusionOperator = Descriptor<"ranked-parents", "ranked-parents">;
export type HydrationOperator = Descriptor<"ranked-parents", "evidence">;
export type RerankerOperator = Descriptor<"evidence", "evidence">;
export type GenerationOperator = Descriptor<"evidence", "answer">;
export interface RecursiveChunkConfig { maxRunes: number; overlapSpans?: number; levels?: string[]; atomic?: string[]; }
export interface StructuredGenerationConfig { model?: string; prompt: string; outputSchema: string; decoding?: JSONObject; seedPolicy?: { mode: string; seed?: number }; }
export interface StructuredSummaryConfig { generator: GenerationOperator; }
export interface SyntheticQuestionsConfig { from: string; count: number; model?: string; prompt?: string; }
export interface EmbeddingConfig { dimensions?: number; distance?: Distance; normalize?: "l2" | "none"; batchSize?: number; }
export interface BleveMultiConfig { lexical: boolean; vector?: { distance: Distance; optimizeFor?: "recall" | "latency" }; }
export interface FilterConfig { sourceIds?: string[]; documentIds?: string[]; contentTypes?: string[]; metadataEquals?: Record<string,string>; }
export interface RetrieveConfig { index: string; representation: RepresentationKind; topK: number; filter?: FilterConfig; }
export interface CollapseConfig { scope: CollapseScope | FactorRef; representative: "scoreThenRepresentationId" | "bestFusionContributionThenId"; }
export interface WeightedRRFConfig { rankConstant?: number; weights?: Record<string,number>; missingChannelPolicy?: "zero" | "reject"; tieBreak?: string; }
export interface HydrationConfig { selection: "bestContributionThenId"; allSupportingChunks?: boolean; }
export interface CrossEncoderConfig { model: string; candidates: number; results: number; truncation?: string; tokenization?: string; inputTemplate?: string; timeoutMilliseconds?: number; }
export interface AnswerConfig { model: string; prompt: string; citations: "required" | "optional"; citationFailurePolicy?: "error" | "abstain"; contextBudgetTokens: number; decoding?: JSONObject; seedPolicy?: { mode: string; seed?: number }; }
export interface DatasetConfig { split: string; status: DatasetStatus; relevanceTarget: "document" | "unit" | "chunk" | "answer"; }
export interface PipelineBuilder { corpus(input: CorpusInput): this; units(op: UnitOperator): this; chunks(op: ChunkOperator): this; representations(op: RepresentationOperator): this; embedding(op: EmbeddingOperator): this; index(name: string, op: IndexOperator): this; use(fragment: Fragment<PipelineBuilder>): this; note(text: string): this; tag(name: string, value: string): this; }
export interface Pipeline { validate(): ValidationReport; explain(): Explanation; readonly __ragPipeline?: true; }
export interface QueryBuilder { channels(values: RetrieverOperator[]): this; collapseChannels(op: CollapseOperator): this; fuse(op: FusionOperator): this; collapseFinal(op: CollapseOperator): this; hydrate(op: HydrationOperator): this; results(count: number): this; }
export interface QueryPlan { readonly __ragQueryPlan?: true; }
export interface FactorContext { factor(name: string): FactorRef; }
export interface VariantBuilder { selectRepresentations(kinds: RepresentationKind[]): this; query(plan: QueryPlan | ((context: FactorContext) => QueryPlan)): this; }
export interface Variant { readonly __ragVariant?: true; }
export interface VariantsBuilder { add(name: string, configure: (variant: VariantBuilder) => void): this; }
export interface FactorsBuilder { enum(name: string, values: string[]): this; }
export interface MetricsBuilder { precisionAt(cutoffs: number[]): this; recallAt(cutoffs: number[]): this; hitRateAt(cutoffs: number[]): this; mrr(): this; ndcgAt(cutoffs: number[]): this; latency(stages: string[]): this; tokenUsage(): this; providerCost(): this; storageBytes(): this; failureRates(): this; }
export interface InvariantsBuilder { require(identifier: string): this; }
export interface RequestBuilder { field(name: string, type: "string" | "string[]", options: { required: boolean; maxLength?: number }): this; }
export interface ResponseBuilder { answer(format: "markdown" | "text"): this; citations(mode: "source" | "required" | "optional"): this; includeTraceId(value: boolean): this; }
export interface RuntimeBuilder { timeoutMs(value: number): this; maxConcurrent(value: number): this; onProviderFailure(value: "fail" | "abstain"): this; trace(value: "authoritative" | "sampled"): this; }
export interface ProductBuilder { pipeline(value: Pipeline): this; query(value: QueryPlan): this; rerank(op: RerankerOperator): this; generate(op: GenerationOperator): this; request(configure: (request: RequestBuilder) => void): this; response(configure: (response: ResponseBuilder) => void): this; runtime(configure: (runtime: RuntimeBuilder) => void): this; tag(name: string,value: string): this; }
export interface Product { validate(): ValidationReport; explain(): Explanation; compileProduct(options: CompileOptions): RAGProductPlanV2; readonly __ragProduct?: true; }
export interface StudyBuilder { pipeline(value: Pipeline): this; dataset(value: DatasetRef): this; variants(configure: (variants: VariantsBuilder) => void): this; factors(configure: (factors: FactorsBuilder) => void): this; replicates(count: number): this; metrics(configure: (metrics: MetricsBuilder) => void): this; invariants(configure: (invariants: InvariantsBuilder) => void): this; tag(name: string,value: string): this; }
export interface Study { validate(): ValidationReport; explain(): Explanation; compileStudy(options: CompileOptions): RAGStudyV2; readonly __ragStudy?: true; }
export interface Explanation { schemaVersion: "rag-explanation/v1"; kind: "pipeline" | "product" | "study"; name: string; nodeCount?: number; variantCount?: number; cellCount?: number; operators?: string[]; factors?: string[]; warnings?: string[]; }
export interface RAGProductPlanV2 { schemaVersion: "rag-product-plan/v2"; pipeline: RAGPipelineIRV2; bindings: ArtifactBinding[]; citations: { mode: string; requireSourceText: boolean }; request: JSONObject; response: JSONObject; runtime: JSONObject; }
export interface RAGStudyV2 { schemaVersion: "rag-study/v2"; variants: Array<{ id: string; pipeline: RAGPipelineIRV2; metadata?: JSONObject }>; factors?: Array<{ id: string; values: Array<{ id: string; value: JSONValue }> }>; bindings: ArtifactBinding[]; dataset: JSONObject; measures: JSONObject[]; replicates: number; }
export interface RAGPipelineIRV2 { schemaVersion: "rag-pipeline-ir/v2"; inputs: JSONObject[]; nodes: Array<{ id: string; operator: { kind: string; version: string }; inputs: JSONObject[]; config: JSONObject; order?: number }>; outputs: JSONObject[]; }
export interface PreviewOptions extends CompileOptions { variant: string; factors: Record<string,string>; query: string; trace: "full" | "summary"; }
export interface PreviewRequest { schemaVersion: "rag-preview-request/v1"; cell: JSONObject; query: string; trace: string; }
export function pipeline(name: string, configure: (pipeline: PipelineBuilder) => void): Pipeline;
export function fragment(name: string, configure: (pipeline: PipelineBuilder) => void): Fragment<PipelineBuilder>;
export function queryPlan(name: string, configure: (query: QueryBuilder) => void): QueryPlan;
export function product(name: string, configure: (product: ProductBuilder) => void): Product;
export function study(name: string, configure: (study: StudyBuilder) => void): Study;
export function variant(name: string, configure: (variant: VariantBuilder) => void): Variant;
export function validate(value: Pipeline | Product | Study): ValidationReport;
export function explain(value: Pipeline | Product | Study): Explanation;
export function compileProduct(value: Product, options: CompileOptions): RAGProductPlanV2;
export function compileStudy(value: Study, options: CompileOptions): RAGStudyV2;
export function preview(value: Study, options: PreviewOptions): PreviewRequest;
export const inputs: { corpus(role: string): CorpusInput };
export const units: { identity(): UnitOperator; individualTurns(): UnitOperator };
export const transcript: { units: { agentsViewRuns(): UnitOperator } };
export const chunks: { recursive(config: RecursiveChunkConfig): ChunkOperator };
export const representations: { raw(name: string): RepresentationOperator; structuredSummary(name: string, config: StructuredSummaryConfig): RepresentationOperator; syntheticQuestions(name: string, config: SyntheticQuestionsConfig): RepresentationOperator; compose(...values: RepresentationOperator[]): RepresentationOperator };
export const embeddings: { model(name: string, config?: EmbeddingConfig): EmbeddingOperator };
export const indexes: { bleveMulti(config: BleveMultiConfig): IndexOperator };
export const retrieve: { bm25(name: string, config: RetrieveConfig): RetrieverOperator; vector(name: string, config: RetrieveConfig): RetrieverOperator };
export const collapse: { parent(config: CollapseConfig): CollapseOperator };
export const fusion: { weightedRRF(config: WeightedRRFConfig): FusionOperator };
export const hydration: { sourceEvidence(config: HydrationConfig): HydrationOperator };
export const rerank: { crossEncoder(config: CrossEncoderConfig): RerankerOperator };
export const generation: { structured(name: string, config: StructuredGenerationConfig): GenerationOperator; answer(config: AnswerConfig): GenerationOperator };
export const metrics: { mrr(): "rag.mrr/v1" };
export const datasets: { artifact(role: string, config: DatasetConfig): DatasetRef };
export const recipes: { transcriptPreparation(config: RecursiveChunkConfig): Fragment<PipelineBuilder> };
export const version: "v2";`}}
}
