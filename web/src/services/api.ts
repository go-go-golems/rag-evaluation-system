import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

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

function filterIdentityParams(args: CorpusIdentityArgs): Record<string, string | number | undefined> {
  const params: Record<string, string | number | undefined> = {};
  if (args.strategy_id) params.strategy_id = args.strategy_id;
  if (args.provider_type) params.provider_type = args.provider_type;
  if (args.model) params.model = args.model;
  if (args.dimensions) params.dimensions = args.dimensions;
  return params;
}

export const ragApi = createApi({
  reducerPath: 'ragApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Sources', 'Documents', 'Chunks', 'Strategies', 'Embeddings', 'Corpus'],
  endpoints: (builder) => ({
    listSources: builder.query<Source[], void>({
      query: () => 'sources',
      transformResponse: (response: { items: Source[] }) => response.items ?? [],
      providesTags: ['Sources'],
    }),
    createSource: builder.mutation<{ id: string; name: string }, Partial<Source> & { config?: Record<string, unknown> }>({
      query: (body) => ({ url: 'sources', method: 'POST', body }),
      invalidatesTags: ['Sources'],
    }),
    listDocuments: builder.query<Document[], void>({
      query: () => 'documents',
      transformResponse: (response: { items: Document[] }) => response.items ?? [],
      providesTags: ['Documents'],
    }),
    getDocument: builder.query<Document, string>({
      query: (id) => `documents/${id}`,
    }),
    listChunks: builder.query<Chunk[], string>({
      query: (docId) => `documents/${docId}/chunks`,
      transformResponse: (response: { items: Chunk[] }) => response.items ?? [],
      providesTags: ['Chunks'],
    }),
    listChunkingStrategies: builder.query<ChunkingStrategy[], void>({
      query: () => 'chunking-strategies',
      transformResponse: (response: { items: ChunkingStrategy[] }) => response.items ?? [],
      providesTags: ['Strategies'],
    }),
    computeEmbeddings: builder.mutation<ComputeEmbeddingsResponse, ComputeEmbeddingsRequest>({
      query: (body) => ({ url: 'embeddings/compute', method: 'POST', body }),
      invalidatesTags: ['Embeddings'],
    }),
    embeddingSimilarity: builder.mutation<SimilarityResponse, SimilarityRequest>({
      query: (body) => ({ url: 'embeddings/similarity', method: 'POST', body }),
    }),

    // --- Corpus Explorer ---
    listCorpusSources: builder.query<CorpusSourceSummary[], CorpusIdentityArgs>({
      query: (args) => ({
        url: 'corpus/sources',
        params: filterIdentityParams(args),
      }),
      transformResponse: (response: { items: CorpusSourceSummary[] }) => response.items ?? [],
      providesTags: ['Corpus'],
    }),
    listCorpusDocuments: builder.query<CorpusDocumentRow[], CorpusDocumentArgs>({
      query: (args) => ({
        url: 'corpus/documents',
        params: {
          source_id: args.source_id,
          limit: args.limit ?? 100,
          offset: args.offset ?? 0,
          ...filterIdentityParams(args),
        },
      }),
      transformResponse: (response: { items: CorpusDocumentRow[] }) => response.items ?? [],
      providesTags: ['Corpus'],
    }),
    getCorpusDocument: builder.query<CorpusDocumentDetail, CorpusDocumentDetailArgs>({
      query: (args) => ({
        url: `corpus/documents/${encodeURIComponent(args.document_id)}`,
        params: {
          include_text: args.include_text ? 'true' : undefined,
          ...filterIdentityParams(args),
        },
      }),
      transformResponse: (response: CorpusDocumentDetail) => ({
        ...response,
        chunks: response.chunks ?? [],
      }),
      providesTags: ['Corpus'],
    }),
  }),
});

export const {
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
} = ragApi;
