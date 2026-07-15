import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import type { WidgetNode } from "@go-go-golems/rag-evaluation-site/ir";

export interface DslPageResponse {
	id: string;
	title: string;
	root: WidgetNode;
}

export interface Source {
	id: string;
	name: string;
	type: string;
	config_json?: string;
	created_at: string;
	updated_at: string;
}

export interface Document {
	id: string;
	source_id: string;
	external_id?: string;
	title: string;
	author: string;
	url?: string;
	content_type: string;
	word_count: number;
	language: string;
	status: string;
	created_at: string;
	updated_at: string;
}

export interface Chunk {
	id: string;
	document_id: string;
	strategy_id: string;
	chunk_index: number;
	text: string;
	token_count: number;
	start_offset: number;
	end_offset: number;
	created_at: string;
}

export interface ChunkingStrategy {
	id: string;
	name: string;
	type: string;
	description: string;
	created_at: string;
}

export interface ComputeEmbeddingsRequest {
	strategy_id: string;
	profile_registries?: string[];
	profile?: string;
	base_profile?: string;
	embeddings_type: string;
	embeddings_engine: string;
	embeddings_dimensions: number;
	api_key?: string;
	base_url?: string;
	cache_type?: string;
	cache_directory?: string;
	batch_size?: number;
	limit?: number;
	force?: boolean;
}

export interface ComputeEmbeddingsResponse {
	strategy_id: string;
	provider_type: string;
	model: string;
	dimensions: number;
	effective_profile?: string;
	considered: number;
	computed: number;
	skipped_fresh: number;
}

export interface SimilarityChunk {
	chunk_id: string;
	document_id: string;
	strategy_id: string;
	chunk_index: number;
	text_preview?: string;
}

export interface SimilarityMatch extends SimilarityChunk {
	score: number;
}

export interface SimilarityResponse {
	strategy_id: string;
	provider_type: string;
	model: string;
	dimensions: number;
	source: SimilarityChunk;
	matches: SimilarityMatch[];
	considered: number;
	candidate_limit: number;
}

export interface SimilarityRequest {
	strategy_id: string;
	provider_type: string;
	model: string;
	dimensions: number;
	chunk_id_a: string;
	chunk_id_b?: string;
	limit?: number;
	candidate_limit?: number;
	preview_runes?: number;
}

// --- Corpus Explorer Types ---

export interface CorpusIdentityArgs {
	strategy_id?: string;
	provider_type?: string;
	model?: string;
	dimensions?: number;
}

export interface CorpusSourceSummary {
	source_id: string;
	source_name: string;
	source_type: string;
	document_count: number;
	word_count: number;
	chunk_count: number;
	embedded_count: number;
	missing_embedding_count: number;
}

export interface CorpusDocumentRow {
	id: string;
	source_id: string;
	title: string;
	url: string;
	word_count: number;
	status: string;
	chunk_count: number;
	embedded_count: number;
	missing_embedding_count: number;
}

export interface CorpusChunk {
	id: string;
	strategy_id: string;
	chunk_index: number;
	start_offset: number;
	end_offset: number;
	token_count: number;
	text: string;
	embedding?: {
		present: boolean;
		provider_type: string;
		model: string;
		dimensions: number;
		text_hash?: string;
		updated_at?: string;
	};
	enrichment?: {
		present: boolean;
		prompt_version?: string;
		short_summary?: string;
		quality_score?: number;
		updated_at?: string;
	};
}

export interface CorpusDocumentDetail {
	document: {
		id: string;
		source_id: string;
		external_id: string;
		title: string;
		url: string;
		word_count: number;
		status: string;
		metadata: Record<string, unknown>;
		content_text?: string;
		content_html?: string;
		content_type: string;
		author: string;
		language: string;
		created_at: string;
		updated_at: string;
	};
	chunks: CorpusChunk[];
}

export interface CorpusDocumentArgs extends CorpusIdentityArgs {
	source_id: string;
	limit?: number;
	offset?: number;
}

export interface CorpusDocumentDetailArgs extends CorpusIdentityArgs {
	document_id: string;
	include_text?: boolean;
}

// --- Search Types ---

export interface RetrievalComponent {
	rank: number;
	score: number;
}

export interface RetrievalResult {
	rank: number;
	chunk_id: string;
	document_id: string;
	source_id: string;
	title: string;
	url?: string;
	strategy_id: string;
	chunk_index: number;
	score: number;
	retriever: string;
	preview: string;
	components?: Record<string, RetrievalComponent>;
}

export interface SearchResult {
	query: string;
	index_id?: string;
	retriever: string;
	items: RetrievalResult[];
}

export interface BM25SearchRequest {
	index_id: string;
	query: string;
	limit?: number;
	preview_runes?: number;
}

export interface VectorSearchRequest {
	strategy_id: string;
	source_ids?: string[];
	query: string;
	profile_registries?: string[];
	profile?: string;
	base_profile?: string;
	embeddings_type: string;
	embeddings_engine: string;
	embeddings_dimensions: number;
	api_key?: string;
	base_url?: string;
	cache_type?: string;
	cache_directory?: string;
	limit?: number;
	candidate_limit?: number;
	preview_runes?: number;
}

export interface HybridSearchRequest extends VectorSearchRequest {
	index_id: string;
	bm25_limit?: number;
	vector_limit?: number;
	rrf_k?: number;
}

export interface EmbeddingCoverageRequest {
	strategy_id: string;
	provider_type: string;
	model: string;
	dimensions: number;
}

export interface EmbeddingCoverageResult {
	strategy_id: string;
	provider_type: string;
	model: string;
	dimensions: number;
	total_chunks: number;
	embedded_chunks: number;
	sources?: Array<{
		source_id: string;
		source_name: string;
		total_chunks: number;
		embedded_chunks: number;
		coverage_pct: number;
		missing_chunks: number;
	}>;
}

// --- Immutable experiment laboratory ---

export interface LabSnapshot {
	id: string;
	document_count: number;
	created_at: string;
}

export interface LabChunkSet {
	id: string;
	corpus_snapshot_id: string;
	chunk_plan_id: string;
	chunk_count: number;
	created_at: string;
}

export interface LabEmbeddingSet {
	id: string;
	chunk_set_id: string;
	embedding_plan_id: string;
	embedding_count: number;
	created_at: string;
}

export interface LabBM25Artifact {
	id: string;
	chunk_set_id: string;
	chunk_count: number;
	created_at: string;
}

export interface LabCatalog {
	snapshots: LabSnapshot[];
	chunk_sets: LabChunkSet[];
	embedding_sets: LabEmbeddingSet[];
	bm25_artifacts: LabBM25Artifact[];
}

export interface ExperimentSpecificationInput {
	corpus_snapshot_id: string;
	chunk_set_id: string;
	bm25_artifact_id?: string;
	embedding_set_id?: string;
	evaluation_dataset_id?: string;
	config: Record<string, unknown>;
}

export interface ExperimentSpecification extends ExperimentSpecificationInput {
	id: string;
	schema_version: string;
	created_at: string;
}

export interface ExperimentRunEvent {
	sequence: number;
	type: string;
	occurred_at: string;
	payload: Record<string, unknown>;
}

export interface ExperimentRunSummary {
	status: "succeeded" | "failed" | "cancelled";
	finished_at: string;
	metrics: Record<string, unknown>;
	cost: Record<string, unknown>;
	storage: Record<string, unknown>;
	error: Record<string, unknown>;
}

export interface ExperimentRun {
	id: string;
	experiment_spec_id: string;
	created_at: string;
	status: string;
	events: ExperimentRunEvent[];
	summary?: ExperimentRunSummary;
}

export interface ExperimentQueryTrace {
	query_card_id: string;
	trace: Record<string, unknown>;
	metrics: Record<string, unknown>;
	timing: Record<string, unknown>;
	cost: Record<string, unknown>;
	storage: Record<string, unknown>;
}

export interface ExperimentComparison {
	left: ExperimentRun;
	right: ExperimentRun;
	left_traces: ExperimentQueryTrace[];
	right_traces: ExperimentQueryTrace[];
}

function filterIdentityParams(
	args: CorpusIdentityArgs,
): Record<string, string | number | undefined> {
	const params: Record<string, string | number | undefined> = {};
	if (args.strategy_id) params.strategy_id = args.strategy_id;
	if (args.provider_type) params.provider_type = args.provider_type;
	if (args.model) params.model = args.model;
	if (args.dimensions) params.dimensions = args.dimensions;
	return params;
}

export const ragApi = createApi({
	reducerPath: "ragApi",
	baseQuery: fetchBaseQuery({ baseUrl: "/api/v1" }),
	tagTypes: [
		"Sources",
		"Documents",
		"Chunks",
		"Strategies",
		"Embeddings",
		"Corpus",
		"Workflows",
		"Artifacts",
		"Lab",
	],
	endpoints: (builder) => ({
		getDslPage: builder.query<DslPageResponse, string>({
			query: (id) => `dsl/pages/${encodeURIComponent(id)}`,
		}),
		listSources: builder.query<Source[], void>({
			query: () => "sources",
			transformResponse: (response: { items: Source[] }) => response.items ?? [],
			providesTags: ["Sources"],
		}),
		createSource: builder.mutation<
			{ id: string; name: string },
			Partial<Source> & { config?: Record<string, unknown> }
		>({
			query: (body) => ({ url: "sources", method: "POST", body }),
			invalidatesTags: ["Sources"],
		}),
		listDocuments: builder.query<Document[], void>({
			query: () => "documents",
			transformResponse: (response: { items: Document[] }) => response.items ?? [],
			providesTags: ["Documents"],
		}),
		getDocument: builder.query<Document, string>({
			query: (id) => `documents/${id}`,
		}),
		listChunks: builder.query<Chunk[], string>({
			query: (docId) => `documents/${docId}/chunks`,
			transformResponse: (response: { items: Chunk[] }) => response.items ?? [],
			providesTags: ["Chunks"],
		}),
		listChunkingStrategies: builder.query<ChunkingStrategy[], void>({
			query: () => "chunking-strategies",
			transformResponse: (response: { items: ChunkingStrategy[] }) => response.items ?? [],
			providesTags: ["Strategies"],
		}),
		computeEmbeddings: builder.mutation<ComputeEmbeddingsResponse, ComputeEmbeddingsRequest>({
			query: (body) => ({ url: "embeddings/compute", method: "POST", body }),
			invalidatesTags: ["Embeddings"],
		}),
		embeddingSimilarity: builder.mutation<SimilarityResponse, SimilarityRequest>({
			query: (body) => ({ url: "embeddings/similarity", method: "POST", body }),
		}),

		// --- Search ---
		searchBM25: builder.mutation<SearchResult, BM25SearchRequest>({
			query: (body) => ({ url: "search/query", method: "POST", body }),
		}),
		searchVector: builder.mutation<SearchResult, VectorSearchRequest>({
			query: (body) => ({ url: "search/vector", method: "POST", body }),
		}),
		searchHybrid: builder.mutation<SearchResult, HybridSearchRequest>({
			query: (body) => ({ url: "search/hybrid", method: "POST", body }),
		}),
		embeddingCoverage: builder.mutation<EmbeddingCoverageResult, EmbeddingCoverageRequest>({
			query: (body) => ({ url: "embeddings/coverage", method: "POST", body }),
		}),

		// --- Corpus Explorer ---
		listCorpusSources: builder.query<CorpusSourceSummary[], CorpusIdentityArgs>({
			query: (args) => ({
				url: "corpus/sources",
				params: filterIdentityParams(args),
			}),
			transformResponse: (response: { items: CorpusSourceSummary[] }) => response.items ?? [],
			providesTags: ["Corpus"],
		}),
		listCorpusDocuments: builder.query<CorpusDocumentRow[], CorpusDocumentArgs>({
			query: (args) => ({
				url: "corpus/documents",
				params: {
					source_id: args.source_id,
					limit: args.limit ?? 100,
					offset: args.offset ?? 0,
					...filterIdentityParams(args),
				},
			}),
			transformResponse: (response: { items: CorpusDocumentRow[] }) => response.items ?? [],
			providesTags: ["Corpus"],
		}),
		getCorpusDocument: builder.query<CorpusDocumentDetail, CorpusDocumentDetailArgs>({
			query: (args) => ({
				url: `corpus/documents/${encodeURIComponent(args.document_id)}`,
				params: {
					include_text: args.include_text ? "true" : undefined,
					...filterIdentityParams(args),
				},
			}),
			transformResponse: (response: CorpusDocumentDetail) => ({
				...response,
				chunks: response.chunks ?? [],
			}),
			providesTags: ["Corpus"],
		}),
		// --- Workflow endpoints ---
		listWorkflows: builder.query<
			WorkflowListResponse,
			{ status?: string; limit?: number; offset?: number }
		>({
			query: (params) => ({
				url: "workflows",
				params: { status: params.status, limit: params.limit ?? 50, offset: params.offset ?? 0 },
			}),
			providesTags: ["Workflows"],
		}),
		getWorkflow: builder.query<WorkflowSummary, string>({
			query: (id) => `workflows/${id}`,
			providesTags: ["Workflows"],
		}),
		getWorkflowOps: builder.query<WorkflowOpsResponse, string>({
			query: (id) => `workflows/${id}/ops`,
			providesTags: ["Workflows"],
		}),
		getOpResult: builder.query<OpResult, { workflowId: string; opId: string }>({
			query: ({ workflowId, opId }) =>
				`workflows/${workflowId}/results/${encodeURIComponent(opId)}`,
			providesTags: ["Workflows"],
		}),
		submitIntakeWorkflow: builder.mutation<SubmitIntakeResponse, SubmitIntakeRequest>({
			query: (body) => ({ url: "workflows/intake", method: "POST", body }),
			invalidatesTags: ["Workflows"],
		}),
		retryOp: builder.mutation<void, { workflowId: string; opId: string }>({
			query: ({ workflowId, opId }) => ({
				url: `workflows/${workflowId}/retry/${opId}`,
				method: "POST",
			}),
			invalidatesTags: ["Workflows"],
		}),
		cancelWorkflow: builder.mutation<void, string>({
			query: (id) => ({ url: `workflows/${id}/cancel`, method: "POST" }),
			invalidatesTags: ["Workflows"],
		}),
		listQueues: builder.query<QueueStatus[], void>({
			query: () => "queues",
			transformResponse: (response: { queues: QueueStatus[] }) => response.queues ?? [],
			providesTags: ["Workflows"],
		}),

		// Artifact coverage endpoints (RAGEVAL-006 Phase 6)
		getDocumentProcessingIdentities: builder.query<DocumentProcessingIdentitiesResult, void>({
			query: () => "artifacts/document-processing/identities",
			providesTags: ["Artifacts"],
		}),
		getChunkEnrichmentIdentities: builder.query<ChunkEnrichmentIdentitiesResult, void>({
			query: () => "artifacts/chunk-enrichment/identities",
			providesTags: ["Artifacts"],
		}),

		getDocumentProcessingCoverage: builder.query<
			DocumentProcessingCoverageResult,
			DocumentProcessingCoverageArgs
		>({
			query: (args) =>
				`artifacts/document-processing/coverage?artifact_type=${args.artifact_type}&prompt_version=${args.prompt_version}&provider=${args.provider}&model=${args.model}`,
			providesTags: ["Artifacts"],
		}),

		getChunkEnrichmentCoverage: builder.query<
			ChunkEnrichmentCoverageResult,
			ChunkEnrichmentCoverageArgs
		>({
			query: (args) =>
				`artifacts/chunk-enrichment/coverage?strategy_id=${args.strategy_id}&prompt_version=${args.prompt_version}`,
			providesTags: ["Artifacts"],
		}),

		getDocumentProcessingArtifacts: builder.query<DocumentProcessingArtifactList, string>({
			query: (docId) => `documents/${docId}/processing-artifacts`,
			providesTags: ["Artifacts"],
		}),

		getChunkEnrichments: builder.query<
			ChunkEnrichmentList,
			{ chunkId: string; strategyId?: string; promptVersion?: string }
		>({
			query: (args) =>
				`chunks/${args.chunkId}/enrichments${args.strategyId ? "?strategy_id=" + args.strategyId : ""}${args.promptVersion ? "&prompt_version=" + args.promptVersion : ""}`,
			providesTags: ["Artifacts"],
		}),

		// --- Immutable experiment laboratory ---
		getLabCatalog: builder.query<LabCatalog, void>({
			query: () => "lab/catalog",
			providesTags: ["Lab"],
		}),
		listExperimentSpecifications: builder.query<ExperimentSpecification[], void>({
			query: () => "lab/specifications",
			transformResponse: (response: { items: ExperimentSpecification[] }) => response.items ?? [],
			providesTags: ["Lab"],
		}),
		createExperimentSpecification: builder.mutation<
			{ item: ExperimentSpecification; reused: boolean },
			ExperimentSpecificationInput
		>({
			query: (body) => ({ url: "lab/specifications", method: "POST", body }),
			invalidatesTags: ["Lab"],
		}),
		listExperimentRuns: builder.query<ExperimentRun[], { specificationId?: string } | void>({
			query: (args) => ({
				url: "lab/runs",
				params: args?.specificationId ? { specification_id: args.specificationId } : undefined,
			}),
			transformResponse: (response: { items: ExperimentRun[] }) => response.items ?? [],
			providesTags: ["Lab"],
		}),
		createExperimentRun: builder.mutation<ExperimentRun, string>({
			query: (specificationId) => ({
				url: `lab/specifications/${encodeURIComponent(specificationId)}/runs`,
				method: "POST",
			}),
			invalidatesTags: ["Lab"],
		}),
		getExperimentRun: builder.query<ExperimentRun, string>({
			query: (id) => `lab/runs/${encodeURIComponent(id)}`,
			providesTags: ["Lab"],
		}),
		listExperimentRunTraces: builder.query<ExperimentQueryTrace[], string>({
			query: (id) => `lab/runs/${encodeURIComponent(id)}/traces`,
			transformResponse: (response: { items: ExperimentQueryTrace[] }) => response.items ?? [],
			providesTags: ["Lab"],
		}),
		getExperimentComparison: builder.query<ExperimentComparison, { left: string; right: string }>({
			query: (args) => ({ url: "lab/comparison", params: args }),
			providesTags: ["Lab"],
		}),
	}),
});

// --- Workflow Types ---

export interface WorkflowListItem {
	workflow: {
		ID: string;
		Site: string;
		Name: string;
		Status: string;
		Input: Record<string, unknown>;
		Metadata: Record<string, string> | null;
		CreatedAt: string;
		UpdatedAt: string;
	};
	opTotal: number;
	opDone: number;
}

export interface WorkflowListResponse {
	workflows: WorkflowListItem[];
	total: number;
}

export interface WorkflowOp {
	op: {
		ID: string;
		WorkflowID: string;
		Kind: string;
		Queue: string;
		DedupKey: string;
		Input: Record<string, unknown>;
		DependsOn: Array<{ OpID: string; Required: boolean }>;
		Retry: { MaxAttempts: number; BackoffKind: string; InitialBackoff: number };
		RetryState: { Attempt: number; NextAttemptAt: string | null; LastError: string };
		Metadata: Record<string, string> | null;
	};
	status: string;
	createdAt: string;
	updatedAt: string;
}

export interface OpResult {
	OpID: string;
	Data: Record<string, unknown> | null;
	Records: Array<{ Table: string; PK: string; Data: Record<string, unknown> }>;
	Artifacts: Array<{ Name: string; Kind: string; ContentType: string; Body: string }>;
	Emitted: Array<{ ID: string; Kind: string; Queue: string }>;
	EmittedIDs: string[];
	Error: { Code: string; Message: string; Retryable: boolean } | null;
	CompletedAt: string;
}

export interface WorkflowOpGroup {
	operation: string;
	queue: string;
	status: string;
	count: number;
	sample?: WorkflowOp;
}

export interface WorkflowOpsResponse {
	workflow_id: string;
	total: number;
	groups: WorkflowOpGroup[];
}

export interface WorkflowSummary {
	workflow: WorkflowListItem["workflow"];
	stats: {
		Total: number;
		Pending: number;
		Ready: number;
		Running: number;
		Succeeded: number;
		Failed: number;
		Canceled: number;
	};
}

export interface QueueStatus {
	site: string;
	queue: string;
	pending: number;
	ready: number;
	running: number;
	succeeded: number;
	failed: number;
	inFlight: number;
	maxInFlight: number;
	tokens?: number;
	ratePerSecond?: number;
}

export interface SubmitIntakeRequest {
	db_path?: string;
	workflow_id?: string;
	name?: string;
	source_ids?: string[];
	document_ids?: string[];
	document_limit?: number;
	strategy?: string;
	chunk_size?: number;
	overlap?: number;
	skip_embeddings?: boolean;
	skip_bm25?: boolean;
	profile?: string;
	base_profile?: string;
	embeddings_type?: string;
	embeddings_engine?: string;
	embeddings_dimensions?: number;
	batch_size?: number;
	force_embeddings?: boolean;
	index_id?: string;
	index_root?: string;
	force_index?: boolean;
}

export interface SubmitIntakeResponse {
	workflow_id: string;
	engine_db: string;
	db_path: string;
	document_ids: string[];
	strategy_id: string;
	op_ids: string[];
}

// ─── Artifact types (RAGEVAL-006 Phase 6) ─────────────────────────────────

export interface DocumentProcessingIdentity {
	artifact_type: string;
	prompt_version: string;
	provider: string;
	model: string;
	artifact_count: number;
}

export interface DocumentProcessingIdentitiesResult {
	items: DocumentProcessingIdentity[];
}

export interface ChunkEnrichmentIdentity {
	strategy_id: string;
	prompt_version: string;
	provider: string;
	model: string;
	enriched_count: number;
}

export interface ChunkEnrichmentIdentitiesResult {
	items: ChunkEnrichmentIdentity[];
}

export interface DocumentProcessingCoverageArgs {
	artifact_type: string;
	prompt_version: string;
	provider: string;
	model: string;
}

export interface DocumentProcessingCoverageItem {
	source_id: string;
	document_count: number;
	artifact_count: number;
	fresh_count: number;
	failed_count: number;
	missing_count: number;
}

export interface DocumentProcessingCoverageResult {
	artifact_type: string;
	prompt_version: string;
	provider: string;
	model: string;
	items: DocumentProcessingCoverageItem[];
	totals: {
		document_count: number;
		artifact_count: number;
		fresh_count: number;
		failed_count: number;
		missing_count: number;
	};
}

export interface ChunkEnrichmentCoverageArgs {
	strategy_id: string;
	prompt_version: string;
}

export interface ChunkEnrichmentCoverageItem {
	source_id: string;
	chunk_count: number;
	enriched_count: number;
	fresh_count: number;
	missing_count: number;
}

export interface ChunkEnrichmentCoverageResult {
	strategy_id: string;
	prompt_version: string;
	items: ChunkEnrichmentCoverageItem[];
	totals: {
		chunk_count: number;
		enriched_count: number;
		fresh_count: number;
		missing_count: number;
	};
}

export interface DocumentProcessingArtifact {
	document_id: string;
	source_id?: string;
	artifact_type: string;
	prompt_version: string;
	provider: string;
	model: string;
	input_hash: string;
	output_text?: string;
	output_json?: string;
	status: string;
	error_code?: string;
	error_message?: string;
	created_at: string;
	updated_at: string;
}

export interface DocumentProcessingArtifactList {
	document_id: string;
	items: DocumentProcessingArtifact[];
}

export interface ChunkEnrichment {
	chunk_id: string;
	document_id?: string;
	strategy_id: string;
	prompt_version: string;
	provider: string;
	model: string;
	short_summary?: string;
	long_summary?: string;
	key_topics_json?: string;
	entities_json?: string;
	hypothetical_questions_json?: string;
	quality_score?: number;
	text_hash: string;
	created_at: string;
	updated_at: string;
}

export interface ChunkEnrichmentList {
	chunk_id: string;
	items: ChunkEnrichment[];
}

export const {
	useGetDslPageQuery,
	useListSourcesQuery,
	useCreateSourceMutation,
	useListDocumentsQuery,
	useGetDocumentQuery,
	useListChunksQuery,
	useListChunkingStrategiesQuery,
	useComputeEmbeddingsMutation,
	useEmbeddingSimilarityMutation,
	useListCorpusSourcesQuery,
	useListCorpusDocumentsQuery,
	useGetCorpusDocumentQuery,
	useSearchBM25Mutation,
	useSearchVectorMutation,
	useSearchHybridMutation,
	useEmbeddingCoverageMutation,
	// Workflow endpoints
	useListWorkflowsQuery,
	useGetWorkflowQuery,
	useGetWorkflowOpsQuery,
	useSubmitIntakeWorkflowMutation,
	useRetryOpMutation,
	useCancelWorkflowMutation,
	useListQueuesQuery,
	// Artifact endpoints (RAGEVAL-006 Phase 6)
	useGetDocumentProcessingIdentitiesQuery,
	useGetChunkEnrichmentIdentitiesQuery,
	useGetDocumentProcessingCoverageQuery,
	useGetChunkEnrichmentCoverageQuery,
	useGetDocumentProcessingArtifactsQuery,
	useGetChunkEnrichmentsQuery,
	// Immutable experiment laboratory
	useGetLabCatalogQuery,
	useListExperimentSpecificationsQuery,
	useCreateExperimentSpecificationMutation,
	useListExperimentRunsQuery,
	useCreateExperimentRunMutation,
	useGetExperimentRunQuery,
	useListExperimentRunTracesQuery,
	useGetExperimentComparisonQuery,
	// Op result
	useGetOpResultQuery,
} = ragApi;
