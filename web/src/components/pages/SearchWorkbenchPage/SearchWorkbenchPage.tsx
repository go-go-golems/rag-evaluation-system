import React, { useState, useCallback, useEffect } from 'react';
import { CoveragePanel, QueryPresetList } from '../../molecules';
import { DashboardGrid, Stack } from '../../layout';
import { ResultInspectorPanel, RetrievalResultsPanel, SearchControlsPanel, type RetrieverType } from '../../organisms';
import styles from './SearchWorkbenchPage.module.css';
import {
  useSearchBM25Mutation,
  useSearchVectorMutation,
  useSearchHybridMutation,
  useEmbeddingCoverageMutation,
  useListCorpusSourcesQuery,
  useGetCorpusDocumentQuery,
  SearchResult,
  RetrievalResult,
  EmbeddingCoverageResult,
  CorpusIdentityArgs,
} from '../../../services/api';

// ─── Defaults ───

const DEFAULT_INDEX_ID = 'bm25-ttc-guides-articles-fixed-1200-150';
const DEFAULT_STRATEGY_ID = 'fixed-1200-150';
const DEFAULT_PROFILE = 'openai-embedding-small';
const DEFAULT_PROVIDER_TYPE = 'openai';
const DEFAULT_MODEL = 'text-embedding-3-small';
const DEFAULT_DIMENSIONS = 1536;

const TEST_QUERIES = [
  'crape myrtle varieties',
  'how to plant arborvitae',
  'hydrangea pruning',
  'emerald green arborvitae spacing',
  'fast growing trees for privacy screen',
  'which trees make a good privacy screen',
  'zone 5 flowering trees',
  'compact evergreen shrubs',
];

// ─── Main component ───

export const SearchWorkbenchPage: React.FC = () => {
  // Search params
  const [query, setQuery] = useState('');
  const [retriever, setRetriever] = useState<RetrieverType>('bm25');
  const [indexId, setIndexId] = useState(DEFAULT_INDEX_ID);
  const [strategyId, setStrategyId] = useState(DEFAULT_STRATEGY_ID);
  const [profile, setProfile] = useState(DEFAULT_PROFILE);
  const [providerType] = useState(DEFAULT_PROVIDER_TYPE);
  const [model] = useState(DEFAULT_MODEL);
  const [dimensions] = useState(DEFAULT_DIMENSIONS);
  const [limit, setLimit] = useState(10);
  const [candidateLimit, setCandidateLimit] = useState(200);
  const [previewRunes, setPreviewRunes] = useState(300);

  // Source filter
  const [selectedSourceIds, setSelectedSourceIds] = useState<string[]>([]);

  // Results
  const [searchResult, setSearchResult] = useState<SearchResult | null>(null);
  const [selectedResult, setSelectedResult] = useState<RetrievalResult | null>(null);
  const [searchError, setSearchError] = useState<string | null>(null);

  // Coverage
  const [coverage, setCoverage] = useState<EmbeddingCoverageResult | null>(null);

  // Mutations
  const [searchBM25, { isLoading: loadingBM25 }] = useSearchBM25Mutation();
  const [searchVector, { isLoading: loadingVector }] = useSearchVectorMutation();
  const [searchHybrid, { isLoading: loadingHybrid }] = useSearchHybridMutation();
  const [fetchCoverage] = useEmbeddingCoverageMutation();

  // Sources for filter
  const identity: CorpusIdentityArgs = {
    strategy_id: strategyId,
    provider_type: providerType,
    model,
    dimensions,
  };
  const { data: sources } = useListCorpusSourcesQuery(identity);

  const isLoading = loadingBM25 || loadingVector || loadingHybrid;

  // Fetch coverage when retriever needs vector
  useEffect(() => {
    if (retriever === 'vector' || retriever === 'hybrid') {
      fetchCoverage({
        strategy_id: strategyId,
        provider_type: providerType,
        model,
        dimensions,
      }).unwrap().then(setCoverage).catch(() => setCoverage(null));
    } else {
      setCoverage(null);
    }
  }, [retriever, strategyId, providerType, model, dimensions, fetchCoverage]);

  const runSearch = useCallback(async () => {
    if (!query.trim()) return;
    setSearchError(null);
    setSelectedResult(null);

    try {
      let result: SearchResult;
      const preview = previewRunes > 0 ? previewRunes : undefined;

      if (retriever === 'bm25') {
        result = await searchBM25({
          index_id: indexId,
          query: query.trim(),
          limit,
          preview_runes: preview,
        }).unwrap();
      } else if (retriever === 'vector') {
        result = await searchVector({
          strategy_id: strategyId,
          source_ids: selectedSourceIds.length > 0 ? selectedSourceIds : undefined,
          query: query.trim(),
          profile,
          embeddings_type: providerType,
          embeddings_engine: model,
          embeddings_dimensions: dimensions,
          limit,
          candidate_limit: candidateLimit,
          preview_runes: preview,
        }).unwrap();
      } else {
        result = await searchHybrid({
          index_id: indexId,
          strategy_id: strategyId,
          source_ids: selectedSourceIds.length > 0 ? selectedSourceIds : undefined,
          query: query.trim(),
          profile,
          embeddings_type: providerType,
          embeddings_engine: model,
          embeddings_dimensions: dimensions,
          limit,
          candidate_limit: candidateLimit,
          preview_runes: preview,
          bm25_limit: 50,
          vector_limit: 50,
          rrf_k: 60,
        }).unwrap();
      }
      setSearchResult(result);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      setSearchError(msg);
      setSearchResult(null);
    }
  }, [query, retriever, indexId, strategyId, profile, providerType, model, dimensions, limit, candidateLimit, previewRunes, selectedSourceIds, searchBM25, searchVector, searchHybrid]);

  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !isLoading) runSearch();
  }, [runSearch, isLoading]);

  const toggleSource = useCallback((sourceId: string) => {
    setSelectedSourceIds(prev =>
      prev.includes(sourceId) ? prev.filter(id => id !== sourceId) : [...prev, sourceId]
    );
  }, []);

  const items = searchResult?.items ?? [];
  const coveragePct = coverage ? (coverage.total_chunks > 0 ? Math.round(coverage.embedded_chunks / coverage.total_chunks * 100) : 0) : null;

  const { data: selectedDocDetail } = useGetCorpusDocumentQuery({
    document_id: selectedResult?.document_id ?? '',
    strategy_id: strategyId,
    provider_type: providerType,
    model,
    dimensions,
    include_text: true,
  }, { skip: !selectedResult?.document_id });

  const openSelectedResultInCorpus = useCallback(() => {
    if (!selectedResult) return;
    window.dispatchEvent(new CustomEvent('rag:navigate-to-chunk', {
      detail: { documentId: selectedResult.document_id, chunkId: selectedResult.chunk_id, sourceId: selectedResult.source_id },
    }));
  }, [selectedResult]);

  return (
    <DashboardGrid recipe="searchWorkbench" className={styles.workbench}>
      {/* Left: Controls */}
      <Stack gap="sm" className={styles.controlsRail}>
        <SearchControlsPanel
          query={query} setQuery={setQuery}
          retriever={retriever} setRetriever={setRetriever}
          indexId={indexId} setIndexId={setIndexId}
          strategyId={strategyId} setStrategyId={setStrategyId}
          profile={profile} setProfile={setProfile}
          limit={limit} setLimit={setLimit}
          candidateLimit={candidateLimit} setCandidateLimit={setCandidateLimit}
          previewRunes={previewRunes} setPreviewRunes={setPreviewRunes}
          sources={sources ?? []}
          selectedSourceIds={selectedSourceIds}
          toggleSource={toggleSource}
          onSearch={runSearch}
          isLoading={isLoading}
          onKeyDown={handleKeyDown}
        />
        {/* Coverage warning */}
        {coverage && (retriever === 'vector' || retriever === 'hybrid') && (
          <CoveragePanel coverage={coverage} coveragePct={coveragePct!} />
        )}
        {/* Quick queries */}
        <QueryPresetList queries={TEST_QUERIES} onSelect={setQuery} />
      </Stack>

      {/* Center: Results */}
      <Stack gap="sm" className={styles.resultsColumn}>
        <RetrievalResultsPanel
          searchResult={searchResult}
          items={items}
          retriever={retriever}
          selectedResult={selectedResult}
          searchError={searchError}
          isLoading={isLoading}
          onSelectResult={setSelectedResult}
        />
      </Stack>

      {/* Right: Inspector */}
      <div className={styles.inspectorRail}>
        {selectedResult ? (
          <ResultInspectorPanel
            result={selectedResult}
            docDetail={selectedDocDetail}
            onClose={() => setSelectedResult(null)}
            onOpenInCorpus={openSelectedResultInCorpus}
          />
        ) : (
          <div className={`panel ${styles.emptyInspector}`}>
            <div className="panel-header"><span>Inspector</span></div>
            <div className={`text-dim ${styles.emptyInspectorBody}`}>
              Click a result to inspect.
            </div>
          </div>
        )}
      </div>
    </DashboardGrid>
  );
};

export default SearchWorkbenchPage;
