import React, { useEffect, useMemo, useState } from 'react';
import { MacWindow } from '../retro/MacWindow';
import { MacButton } from '../retro/MacButton';
import {
  Chunk,
  useComputeEmbeddingsMutation,
  useEmbeddingSimilarityMutation,
  useListChunkingStrategiesQuery,
  useListChunksQuery,
  useListDocumentsQuery,
} from '../../services/api';

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
    if (!firstChunk) {
      return;
    }
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
    <div className="flex flex-col gap-2">
      <MacWindow title="Embedding Inspector — Controls">
        <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
          <fieldset className="mac-fieldset">
            <legend>Strategy</legend>
            <label className="mac-form-row">
              <span>Chunking strategy</span>
              <select className="mac-select" value={strategyId} onChange={(event) => setStrategyId(event.target.value)}>
                {strategiesLoading ? <option>Loading...</option> : null}
                {strategies.map((strategy) => (
                  <option key={strategy.id} value={strategy.id}>
                    {strategy.id}
                  </option>
                ))}
              </select>
            </label>
            <div className="mac-help-text">
              {selectedStrategy ? `${selectedStrategy.type}: ${selectedStrategy.description || 'no description'}` : 'Create chunks before computing embeddings.'}
            </div>
          </fieldset>

          <fieldset className="mac-fieldset">
            <legend>Provider identity</legend>
            <label className="mac-form-row">
              <span>Provider</span>
              <input className="mac-input" value={providerType} onChange={(event) => setProviderType(event.target.value)} />
            </label>
            <label className="mac-form-row">
              <span>Model</span>
              <input className="mac-input" value={model} onChange={(event) => setModel(event.target.value)} />
            </label>
            <label className="mac-form-row">
              <span>Dimensions</span>
              <input
                className="mac-input"
                type="number"
                min={1}
                value={dimensions}
                onChange={(event) => setDimensions(Number(event.target.value))}
              />
            </label>
          </fieldset>

          <fieldset className="mac-fieldset">
            <legend>Coverage snapshot</legend>
            <dl className="mac-stat-list">
              <dt>Documents loaded</dt>
              <dd>{documentsLoading ? 'Loading...' : documents.length}</dd>
              <dt>Strategies</dt>
              <dd>{strategiesLoading ? 'Loading...' : strategies.length}</dd>
              <dt>Visible chunks for document/strategy</dt>
              <dd>{chunksLoading ? 'Loading...' : strategyChunks.length}</dd>
              <dt>Last compute</dt>
              <dd>{computeState.data ? `${computeState.data.computed} computed, ${computeState.data.skipped_fresh} fresh` : 'Not run in this session'}</dd>
            </dl>
          </fieldset>
        </div>
      </MacWindow>

      <div className="grid grid-cols-1 gap-2 lg:grid-cols-2">
        <MacWindow title="Compute Embeddings">
          <p className="mac-help-text">
            This triggers the backend compute endpoint. Keep the limit small for first tests, especially when using a live Ollama or OpenAI provider.
          </p>
          <div className="mac-form-grid">
            <label className="mac-form-row">
              <span>Batch size</span>
              <input className="mac-input" type="number" min={1} value={batchSize} onChange={(event) => setBatchSize(Number(event.target.value))} />
            </label>
            <label className="mac-form-row">
              <span>Chunk limit</span>
              <input className="mac-input" type="number" min={0} value={computeLimit} onChange={(event) => setComputeLimit(Number(event.target.value))} />
            </label>
            <label className="mac-checkbox-row">
              <input type="checkbox" checked={force} onChange={(event) => setForce(event.target.checked)} />
              <span>Force recompute fresh embeddings</span>
            </label>
          </div>
          <div className="mt-2">
            <MacButton label={computeState.isLoading ? 'Computing...' : 'Compute Embeddings'} onClick={handleCompute} primary disabled={!canCompute || computeState.isLoading} />
          </div>
          {computeState.data ? (
            <table className="mac-table mt-2">
              <tbody>
                <tr><th>Considered</th><td>{computeState.data.considered}</td></tr>
                <tr><th>Computed</th><td>{computeState.data.computed}</td></tr>
                <tr><th>Skipped fresh</th><td>{computeState.data.skipped_fresh}</td></tr>
                <tr><th>Model</th><td>{computeState.data.provider_type}/{computeState.data.model} ({computeState.data.dimensions})</td></tr>
              </tbody>
            </table>
          ) : null}
          {computeState.error ? <pre className="mac-error-box">{formatApiError(computeState.error)}</pre> : null}
        </MacWindow>

        <MacWindow title="Pairwise / Bounded Similarity">
          <p className="mac-help-text">
            Similarity reads stored vectors only. It does not call the embedding provider. Select a document with chunks, then compare two chunk IDs.
          </p>
          <label className="mac-form-row">
            <span>Document</span>
            <select className="mac-select" value={documentId} onChange={(event) => setDocumentId(event.target.value)}>
              {documentsLoading ? <option>Loading...</option> : null}
              {documents.map((document) => (
                <option key={document.id} value={document.id}>
                  {document.title || document.id}
                </option>
              ))}
            </select>
          </label>
          <div className="mac-help-text">
            {selectedDocument ? `${selectedDocument.status} · ${selectedDocument.word_count} words · ${selectedDocument.id}` : 'No document selected.'}
          </div>
          <div className="mac-form-grid mt-2">
            <label className="mac-form-row">
              <span>Chunk A</span>
              <select className="mac-select" value={chunkIDA} onChange={(event) => setChunkIDA(event.target.value)}>
                {strategyChunks.map((chunk) => (
                  <option key={chunk.id} value={chunk.id}>
                    #{chunk.chunk_index} {chunk.id}
                  </option>
                ))}
              </select>
            </label>
            <label className="mac-form-row">
              <span>Chunk B</span>
              <select className="mac-select" value={chunkIDB} onChange={(event) => setChunkIDB(event.target.value)}>
                <option value="">Top candidates for A</option>
                {strategyChunks.map((chunk) => (
                  <option key={chunk.id} value={chunk.id}>
                    #{chunk.chunk_index} {chunk.id}
                  </option>
                ))}
              </select>
            </label>
            <label className="mac-form-row">
              <span>Result limit</span>
              <input className="mac-input" type="number" min={1} value={matchLimit} onChange={(event) => setMatchLimit(Number(event.target.value))} />
            </label>
          </div>
          <div className="mt-2">
            <MacButton label={similarityState.isLoading ? 'Comparing...' : 'Compare Similarity'} onClick={handleCompare} primary disabled={!canCompare || similarityState.isLoading} />
          </div>
          {similarityState.error ? <pre className="mac-error-box">{formatApiError(similarityState.error)}</pre> : null}
        </MacWindow>
      </div>

      {similarityState.data ? (
        <MacWindow title="Similarity Results">
          <div className="mac-help-text">
            Source: {similarityState.data.source.chunk_id} · considered {similarityState.data.considered} candidates · candidate limit {similarityState.data.candidate_limit}
          </div>
          <table className="mac-table mt-2">
            <thead>
              <tr>
                <th>Score</th>
                <th>Chunk</th>
                <th>Document</th>
                <th>Index</th>
                <th>Preview</th>
              </tr>
            </thead>
            <tbody>
              {similarityState.data.matches.map((match) => (
                <tr key={match.chunk_id}>
                  <td>{match.score.toFixed(6)}</td>
                  <td>{match.chunk_id}</td>
                  <td>{match.document_id}</td>
                  <td>{match.chunk_index}</td>
                  <td>{match.text_preview}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </MacWindow>
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
