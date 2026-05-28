import React, { useEffect, useMemo, useState } from 'react';
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
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      {/* Controls */}
      <div className="panel">
        <div className="panel-header"><span>Embedding Inspector — Controls</span></div>
        <div className="panel-body" style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 8 }}>
          <fieldset className="fieldset">
            <legend>Strategy</legend>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 50 }}>Strategy</span>
                <select className="select" value={strategyId} onChange={(e) => setStrategyId(e.target.value)} style={{ flex: 1 }}>
                  {strategiesLoading ? <option>Loading...</option> : null}
                  {strategies.map((s) => <option key={s.id} value={s.id}>{s.id}</option>)}
                </select>
              </label>
              <span className="text-dim text-small">
                {selectedStrategy ? `${selectedStrategy.type}: ${selectedStrategy.description || 'no description'}` : 'Create chunks before computing embeddings.'}
              </span>
            </div>
          </fieldset>

          <fieldset className="fieldset">
            <legend>Provider Identity</legend>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 50 }}>Provider</span>
                <input className="input" value={providerType} onChange={(e) => setProviderType(e.target.value)} style={{ flex: 1 }} />
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 50 }}>Model</span>
                <input className="input" value={model} onChange={(e) => setModel(e.target.value)} style={{ flex: 1 }} />
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 50 }}>Dims</span>
                <input className="input" type="number" min={1} value={dimensions} onChange={(e) => setDimensions(Number(e.target.value))} style={{ width: 60 }} />
              </label>
            </div>
          </fieldset>

          <fieldset className="fieldset">
            <legend>Snapshot</legend>
            <div className="stat-grid">
              <span className="stat-label">Documents</span>
              <span className="stat-value">{documentsLoading ? '...' : documents.length}</span>
              <span className="stat-label">Strategies</span>
              <span className="stat-value">{strategiesLoading ? '...' : strategies.length}</span>
              <span className="stat-label">Chunks</span>
              <span className="stat-value">{chunksLoading ? '...' : strategyChunks.length}</span>
              <span className="stat-label">Last compute</span>
              <span className="stat-value">{computeState.data ? `${computeState.data.computed} computed, ${computeState.data.skipped_fresh} fresh` : '—'}</span>
            </div>
          </fieldset>
        </div>
      </div>

      {/* Compute + Similarity */}
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
        <div className="panel">
          <div className="panel-header"><span>Compute Embeddings</span></div>
          <div className="panel-body">
            <p className="text-dim text-small">Keep the limit small for first tests.</p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 60 }}>Batch</span>
                <input className="input" type="number" min={1} value={batchSize} onChange={(e) => setBatchSize(Number(e.target.value))} style={{ width: 60 }} />
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 60 }}>Limit</span>
                <input className="input" type="number" min={0} value={computeLimit} onChange={(e) => setComputeLimit(Number(e.target.value))} style={{ width: 60 }} />
              </label>
              <label className="checkbox-row">
                <input type="checkbox" checked={force} onChange={(e) => setForce(e.target.checked)} />
                <span>Force recompute</span>
              </label>
            </div>
            <div style={{ marginTop: 6 }}>
              <button className="btn btn-primary" onClick={handleCompute} disabled={!canCompute || computeState.isLoading}>
                {computeState.isLoading ? 'Computing...' : 'Compute Embeddings'}
              </button>
            </div>
            {computeState.data ? (
              <table className="data-table" style={{ marginTop: 6 }}>
                <tbody>
                  <tr><th>Considered</th><td>{computeState.data.considered}</td></tr>
                  <tr><th>Computed</th><td>{computeState.data.computed}</td></tr>
                  <tr><th>Skipped fresh</th><td>{computeState.data.skipped_fresh}</td></tr>
                  <tr><th>Model</th><td className="mono">{computeState.data.provider_type}/{computeState.data.model} ({computeState.data.dimensions})</td></tr>
                </tbody>
              </table>
            ) : null}
            {computeState.error ? <pre className="error-box">{formatApiError(computeState.error)}</pre> : null}
          </div>
        </div>

        <div className="panel">
          <div className="panel-header"><span>Pairwise / Bounded Similarity</span></div>
          <div className="panel-body">
            <p className="text-dim text-small">Similarity reads stored vectors only. Select a document with chunks, then compare.</p>
            <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span className="text-mono text-dim" style={{ width: 60 }}>Document</span>
              <select className="select" value={documentId} onChange={(e) => setDocumentId(e.target.value)} style={{ flex: 1 }}>
                {documentsLoading ? <option>Loading...</option> : null}
                {documents.map((d) => <option key={d.id} value={d.id}>{d.title || d.id}</option>)}
              </select>
            </label>
            <div className="text-dim text-small" style={{ margin: '4px 0' }}>
              {selectedDocument ? `${selectedDocument.status} · ${selectedDocument.word_count} words · ${selectedDocument.id}` : 'No document selected.'}
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 60 }}>Chunk A</span>
                <select className="select" value={chunkIDA} onChange={(e) => setChunkIDA(e.target.value)} style={{ flex: 1 }}>
                  {strategyChunks.map((c) => <option key={c.id} value={c.id}>#{c.chunk_index} {c.id}</option>)}
                </select>
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 60 }}>Chunk B</span>
                <select className="select" value={chunkIDB} onChange={(e) => setChunkIDB(e.target.value)} style={{ flex: 1 }}>
                  <option value="">Top candidates for A</option>
                  {strategyChunks.map((c) => <option key={c.id} value={c.id}>#{c.chunk_index} {c.id}</option>)}
                </select>
              </label>
              <label style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
                <span className="text-mono text-dim" style={{ width: 60 }}>Limit</span>
                <input className="input" type="number" min={1} value={matchLimit} onChange={(e) => setMatchLimit(Number(e.target.value))} style={{ width: 60 }} />
              </label>
            </div>
            <div style={{ marginTop: 6 }}>
              <button className="btn btn-primary" onClick={handleCompare} disabled={!canCompare || similarityState.isLoading}>
                {similarityState.isLoading ? 'Comparing...' : 'Compare Similarity'}
              </button>
            </div>
            {similarityState.error ? <pre className="error-box">{formatApiError(similarityState.error)}</pre> : null}
          </div>
        </div>
      </div>

      {/* Similarity Results */}
      {similarityState.data ? (
        <div className="panel">
          <div className="panel-header"><span>Similarity Results</span></div>
          <div className="panel-body-condensed">
            <div className="text-dim text-small">
              Source: {similarityState.data.source.chunk_id} · {similarityState.data.considered} candidates
            </div>
            <table className="data-table" style={{ marginTop: 4 }}>
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
                    <td className="mono">{match.score.toFixed(6)}</td>
                    <td className="mono">{match.chunk_id}</td>
                    <td className="mono">{match.document_id}</td>
                    <td className="num">{match.chunk_index}</td>
                    <td className="text-small">{match.text_preview}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
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
