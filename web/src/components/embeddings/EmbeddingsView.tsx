import React, { useEffect, useMemo, useState } from 'react';
import { Button, ErrorCallout, SelectInput, TextInput } from '../atoms';
import { Caption, CodeText } from '../foundation';
import { DashboardGrid, FormRow, Panel, Stack } from '../layout';
import { DataTable, MetadataGrid } from '../molecules';
import {
  Chunk,
  useComputeEmbeddingsMutation,
  useEmbeddingSimilarityMutation,
  useListChunkingStrategiesQuery,
  useListChunksQuery,
  useListDocumentsQuery,
} from '../../services/api';
import styles from './EmbeddingsView.module.css';

const DEFAULT_PROVIDER = 'ollama';
const DEFAULT_MODEL = 'nomic-embed-text';
const DEFAULT_DIMENSIONS = 768;

export const EmbeddingsView: React.FC = () => {
  const { data: strategies = [], isLoading: strategiesLoading } = useListChunkingStrategiesQuery();
  const { data: documents = [], isLoading: documentsLoading } = useListDocumentsQuery();

  const [strategyId, setStrategyId] = useState('');
  const [documentId, setDocumentId] = useState('');
  const [providerType, setProviderType] = useState(DEFAULT_PROVIDER);
  const [model, setModel] = useState(DEFAULT_MODEL);
  const [dimensions, setDimensions] = useState(DEFAULT_DIMENSIONS);
  const [batchSize, setBatchSize] = useState(16);
  const [computeLimit, setComputeLimit] = useState(10);
  const [force, setForce] = useState(false);
  const [chunkIDA, setChunkIDA] = useState('');
  const [chunkIDB, setChunkIDB] = useState('');
  const [matchLimit, setMatchLimit] = useState(10);

  const { data: chunks = [], isFetching: chunksLoading } = useListChunksQuery(documentId, {
    skip: !documentId,
  });

  const [computeEmbeddings, computeState] = useComputeEmbeddingsMutation();
  const [embeddingSimilarity, similarityState] = useEmbeddingSimilarityMutation();

  useEffect(() => {
    const firstStrategy = strategies[0];
    if (!strategyId && firstStrategy) {
      setStrategyId(firstStrategy.id);
    }
  }, [strategies, strategyId]);

  useEffect(() => {
    const firstDocument = documents[0];
    if (!documentId && firstDocument) {
      setDocumentId(firstDocument.id);
    }
  }, [documents, documentId]);

  const strategyChunks = useMemo(() => {
    if (!strategyId) return chunks;
    return chunks.filter((chunk: Chunk) => chunk.strategy_id === strategyId);
  }, [chunks, strategyId]);

  useEffect(() => {
    const firstChunk = strategyChunks[0];
    const defaultTargetChunk = strategyChunks[Math.min(1, strategyChunks.length - 1)];
    if (!firstChunk) return;
    if (!chunkIDA || !strategyChunks.some((chunk) => chunk.id === chunkIDA)) {
      setChunkIDA(firstChunk.id);
    }
    if (defaultTargetChunk && (!chunkIDB || !strategyChunks.some((chunk) => chunk.id === chunkIDB))) {
      setChunkIDB(defaultTargetChunk.id);
    }
  }, [strategyChunks, chunkIDA, chunkIDB]);

  const selectedStrategy = strategies.find((strategy) => strategy.id === strategyId);
  const selectedDocument = documents.find((document) => document.id === documentId);

  const canCompute = Boolean(strategyId && providerType && model && dimensions > 0);
  const canCompare = Boolean(strategyId && providerType && model && dimensions > 0 && chunkIDA);

  const handleCompute = async () => {
    if (!canCompute) return;
    await computeEmbeddings({
      strategy_id: strategyId,
      embeddings_type: providerType,
      embeddings_engine: model,
      embeddings_dimensions: dimensions,
      cache_type: 'none',
      batch_size: batchSize,
      limit: computeLimit,
      force,
    });
  };

  const handleCompare = async () => {
    if (!canCompare) return;
    await embeddingSimilarity({
      strategy_id: strategyId,
      provider_type: providerType,
      model,
      dimensions,
      chunk_id_a: chunkIDA,
      chunk_id_b: chunkIDB || undefined,
      limit: matchLimit,
      candidate_limit: 200,
      preview_runes: 160,
    });
  };

  return (
    <div className={styles.root}>
      <Panel title="Embedding Inspector — Controls">
        <div className={styles.controlsGrid}>
          <fieldset className={styles.fieldset}>
            <legend>Strategy</legend>
            <Stack gap="xs">
              <FormRow
                label="Strategy"
                control={(
                  <SelectInput value={strategyId} onChange={(event) => setStrategyId(event.target.value)}>
                    {strategiesLoading ? <option>Loading...</option> : null}
                    {strategies.map((strategy) => <option key={strategy.id} value={strategy.id}>{strategy.id}</option>)}
                  </SelectInput>
                )}
              />
              <Caption>
                {selectedStrategy ? `${selectedStrategy.type}: ${selectedStrategy.description || 'no description'}` : 'Create chunks before computing embeddings.'}
              </Caption>
            </Stack>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Provider Identity</legend>
            <Stack gap="xs">
              <FormRow label="Provider" control={<TextInput value={providerType} onChange={(event) => setProviderType(event.target.value)} />} />
              <FormRow label="Model" control={<TextInput value={model} onChange={(event) => setModel(event.target.value)} />} />
              <FormRow label="Dims" control={<TextInput className={styles.shortInput} type="number" min={1} value={dimensions} onChange={(event) => setDimensions(Number(event.target.value))} />} />
            </Stack>
          </fieldset>

          <fieldset className={styles.fieldset}>
            <legend>Snapshot</legend>
            <MetadataGrid
              density="compact"
              items={[
                { key: 'Documents', value: documentsLoading ? '...' : documents.length },
                { key: 'Strategies', value: strategiesLoading ? '...' : strategies.length },
                { key: 'Chunks', value: chunksLoading ? '...' : strategyChunks.length },
                { key: 'Last compute', value: computeState.data ? `${computeState.data.computed} computed, ${computeState.data.skipped_fresh} fresh` : '—' },
              ]}
            />
          </fieldset>
        </div>
      </Panel>

      <DashboardGrid recipe="twoColumn">
        <Panel title="Compute Embeddings">
          <Stack gap="sm">
            <Caption>Keep the limit small for first tests.</Caption>
            <Stack gap="xs">
              <FormRow label="Batch" control={<TextInput className={styles.shortInput} type="number" min={1} value={batchSize} onChange={(event) => setBatchSize(Number(event.target.value))} />} />
              <FormRow label="Limit" control={<TextInput className={styles.shortInput} type="number" min={0} value={computeLimit} onChange={(event) => setComputeLimit(Number(event.target.value))} />} />
              <label className={styles.checkboxRow}>
                <input type="checkbox" checked={force} onChange={(event) => setForce(event.target.checked)} />
                <span>Force recompute</span>
              </label>
            </Stack>
            <div className={styles.actions}>
              <Button variant="primary" onClick={handleCompute} disabled={!canCompute || computeState.isLoading}>
                {computeState.isLoading ? 'Computing...' : 'Compute Embeddings'}
              </Button>
            </div>
            {computeState.data ? (
              <MetadataGrid
                className={styles.tableBlock}
                density="compact"
                items={[
                  { key: 'Considered', value: computeState.data.considered },
                  { key: 'Computed', value: computeState.data.computed },
                  { key: 'Skipped fresh', value: computeState.data.skipped_fresh },
                  { key: 'Model', value: <CodeText>{computeState.data.provider_type}/{computeState.data.model} ({computeState.data.dimensions})</CodeText> },
                ]}
              />
            ) : null}
            {computeState.error ? <ErrorCallout>{formatApiError(computeState.error)}</ErrorCallout> : null}
          </Stack>
        </Panel>

        <Panel title="Pairwise / Bounded Similarity">
          <Stack gap="sm">
            <Caption>Similarity reads stored vectors only. Select a document with chunks, then compare.</Caption>
            <FormRow
              label="Document"
              control={(
                <SelectInput value={documentId} onChange={(event) => setDocumentId(event.target.value)}>
                  {documentsLoading ? <option>Loading...</option> : null}
                  {documents.map((document) => <option key={document.id} value={document.id}>{document.title || document.id}</option>)}
                </SelectInput>
              )}
              hint={selectedDocument ? `${selectedDocument.status} · ${selectedDocument.word_count} words · ${selectedDocument.id}` : 'No document selected.'}
            />
            <Stack gap="xs">
              <FormRow
                label="Chunk A"
                control={(
                  <SelectInput value={chunkIDA} onChange={(event) => setChunkIDA(event.target.value)}>
                    {strategyChunks.map((chunk) => <option key={chunk.id} value={chunk.id}>#{chunk.chunk_index} {chunk.id}</option>)}
                  </SelectInput>
                )}
              />
              <FormRow
                label="Chunk B"
                control={(
                  <SelectInput value={chunkIDB} onChange={(event) => setChunkIDB(event.target.value)}>
                    <option value="">Top candidates for A</option>
                    {strategyChunks.map((chunk) => <option key={chunk.id} value={chunk.id}>#{chunk.chunk_index} {chunk.id}</option>)}
                  </SelectInput>
                )}
              />
              <FormRow label="Limit" control={<TextInput className={styles.shortInput} type="number" min={1} value={matchLimit} onChange={(event) => setMatchLimit(Number(event.target.value))} />} />
            </Stack>
            <div className={styles.actions}>
              <Button variant="primary" onClick={handleCompare} disabled={!canCompare || similarityState.isLoading}>
                {similarityState.isLoading ? 'Comparing...' : 'Compare Similarity'}
              </Button>
            </div>
            {similarityState.error ? <ErrorCallout>{formatApiError(similarityState.error)}</ErrorCallout> : null}
          </Stack>
        </Panel>
      </DashboardGrid>

      {similarityState.data ? (
        <Panel title="Similarity Results" density="condensed">
          <Caption className={styles.resultIntro}>Source: {similarityState.data.source.chunk_id} · {similarityState.data.considered} candidates</Caption>
          <DataTable
            rows={similarityState.data.matches}
            getRowKey={(match) => match.chunk_id}
            columns={[
              { id: 'score', header: 'Score', cell: (match) => <CodeText>{match.score.toFixed(6)}</CodeText> },
              { id: 'chunk', header: 'Chunk', cell: (match) => <CodeText>{match.chunk_id}</CodeText> },
              { id: 'document', header: 'Document', cell: (match) => <CodeText>{match.document_id}</CodeText> },
              { id: 'index', header: 'Index', align: 'end', cell: (match) => match.chunk_index },
              { id: 'preview', header: 'Preview', cell: (match) => <Caption tone="inherit">{match.text_preview}</Caption> },
            ]}
          />
        </Panel>
      ) : null}
    </div>
  );
};

function formatApiError(error: unknown): string {
  if (typeof error === 'object' && error !== null && 'data' in error) {
    return JSON.stringify((error as { data: unknown }).data, null, 2);
  }
  return JSON.stringify(error, null, 2);
}
