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

export const ragApi = createApi({
  reducerPath: 'ragApi',
  baseQuery: fetchBaseQuery({ baseUrl: '/api/v1' }),
  tagTypes: ['Sources', 'Documents', 'Chunks', 'Strategies', 'Embeddings'],
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
} = ragApi;
