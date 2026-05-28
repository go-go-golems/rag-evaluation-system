import React, { useState, useMemo, useCallback } from 'react';
import {
  useListCorpusSourcesQuery,
  useListCorpusDocumentsQuery,
  useGetCorpusDocumentQuery,
  useListChunkingStrategiesQuery,
  CorpusSourceSummary,
  CorpusIdentityArgs,
  CorpusChunk,
} from '../../services/api';

// Default identity matches the current OpenAI smoke data
const DEFAULT_IDENTITY: CorpusIdentityArgs = {
  strategy_id: 'fixed-1200-150',
  provider_type: 'openai',
  model: 'text-embedding-3-small',
  dimensions: 1536,
};

export const CorpusExplorerView: React.FC = () => {
  const [identity, setIdentity] = useState<CorpusIdentityArgs>(DEFAULT_IDENTITY);
  const [sourceId, setSourceId] = useState('');
  const [documentId, setDocumentId] = useState('');

  const { data: strategies = [] } = useListChunkingStrategiesQuery();
  const { data: sources = [], isLoading: sourcesLoading } = useListCorpusSourcesQuery(identity);
  const { data: documents = [], isLoading: docsLoading } = useListCorpusDocumentsQuery(
    { ...identity, source_id: sourceId },
    { skip: !sourceId },
  );
  const { data: detail, isLoading: detailLoading } = useGetCorpusDocumentQuery(
    { ...identity, document_id: documentId, include_text: true },
    { skip: !documentId },
  );

  const selectedSource = useMemo(
    () => sources.find((s) => s.source_id === sourceId),
    [sources, sourceId],
  );

  const selectedDoc = useMemo(
    () => documents.find((d) => d.id === documentId),
    [documents, documentId],
  );

  const handleSelectSource = useCallback((id: string) => {
    setSourceId(id);
    setDocumentId('');
  }, []);

  const handleSelectDocument = useCallback((id: string) => {
    setDocumentId(id);
  }, []);

  const totalDocs = useMemo(() => sources.reduce((s, x) => s + x.document_count, 0), [sources]);
  const totalChunks = useMemo(() => sources.reduce((s, x) => s + x.chunk_count, 0), [sources]);
  const totalEmbedded = useMemo(() => sources.reduce((s, x) => s + x.embedded_count, 0), [sources]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      {/* Identity Bar */}
      <div className="panel">
        <div className="panel-header">
          <span>Embedding Identity</span>
        </div>
        <div className="panel-body-condensed" style={{ display: 'flex', gap: 12, alignItems: 'center', flexWrap: 'wrap' }}>
          <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <span className="text-mono text-dim">Strategy</span>
            <select
              className="select"
              value={identity.strategy_id}
              onChange={(e) => setIdentity((i) => ({ ...i, strategy_id: e.target.value }))}
            >
              {strategies.map((s) => (
                <option key={s.id} value={s.id}>{s.id}</option>
              ))}
            </select>
          </label>
          <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <span className="text-mono text-dim">Provider</span>
            <input
              className="input"
              value={identity.provider_type}
              onChange={(e) => setIdentity((i) => ({ ...i, provider_type: e.target.value }))}
              style={{ width: 80 }}
            />
          </label>
          <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <span className="text-mono text-dim">Model</span>
            <input
              className="input"
              value={identity.model}
              onChange={(e) => setIdentity((i) => ({ ...i, model: e.target.value }))}
              style={{ width: 160 }}
            />
          </label>
          <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <span className="text-mono text-dim">Dims</span>
            <input
              className="input"
              type="number"
              value={identity.dimensions}
              onChange={(e) => setIdentity((i) => ({ ...i, dimensions: Number(e.target.value) }))}
              style={{ width: 60 }}
            />
          </label>
          <span className="text-mono text-dim" style={{ marginLeft: 'auto' }}>
            {totalDocs} docs · {totalChunks} chunks · {totalEmbedded} embedded
          </span>
        </div>
      </div>

      {/* Main layout: Source list | Documents | Detail */}
      <div style={{ display: 'flex', gap: 8, minHeight: 500 }}>
        {/* Source List */}
        <div className="panel" style={{ width: 220, flexShrink: 0 }}>
          <div className="panel-header">
            <span>Sources</span>
            <span className="text-mono" style={{ fontSize: 10 }}>{sources.length}</span>
          </div>
          <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 600 }}>
            {sourcesLoading ? (
              <span className="text-dim text-mono">Loading...</span>
            ) : (
              sources.map((src) => (
                <SourceItem
                  key={src.source_id}
                  source={src}
                  selected={src.source_id === sourceId}
                  onClick={() => handleSelectSource(src.source_id)}
                />
              ))
            )}
          </div>
        </div>

        {/* Document Browser */}
        <div className="panel" style={{ flex: 1, minWidth: 0 }}>
          <div className="panel-header">
            <span>{selectedSource ? `${selectedSource.source_name} — Documents` : 'Documents'}</span>
            {selectedSource && (
              <span className="text-mono" style={{ fontSize: 10 }}>
                {selectedSource.document_count} docs
              </span>
            )}
          </div>
          <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 600 }}>
            {!sourceId ? (
              <span className="text-dim text-mono">Select a source to browse documents.</span>
            ) : docsLoading ? (
              <span className="text-dim text-mono">Loading...</span>
            ) : documents.length === 0 ? (
              <span className="text-dim text-mono">No documents found.</span>
            ) : (
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Title</th>
                    <th>Words</th>
                    <th>Chunks</th>
                    <th>Embed</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {documents.map((doc) => (
                    <tr
                      key={doc.id}
                      className={`selectable ${doc.id === documentId ? 'selected' : ''}`}
                      onClick={() => handleSelectDocument(doc.id)}
                    >
                      <td className="truncate" style={{ maxWidth: 300 }}>{doc.title}</td>
                      <td className="num">{doc.word_count.toLocaleString()}</td>
                      <td className="num">{doc.chunk_count}</td>
                      <td className="num">
                        {doc.chunk_count > 0 ? (
                          <span className={doc.embedded_count === doc.chunk_count ? 'accent-green' : doc.embedded_count > 0 ? 'accent-amber' : 'accent-dim'}>
                            {doc.embedded_count}/{doc.chunk_count}
                          </span>
                        ) : '—'}
                      </td>
                      <td>
                        <span className={`status-${doc.status === 'chunked' ? 'done' : doc.status}`}>{doc.status}</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>

        {/* Document Inspector */}
        <div className="panel" style={{ flex: 1, minWidth: 0 }}>
          <div className="panel-header">
            <span>{selectedDoc ? selectedDoc.title : 'Document Inspector'}</span>
            {selectedDoc && (
              <button
                className="copy-btn"
                title="Copy document ID"
                onClick={() => navigator.clipboard.writeText(selectedDoc.id)}
              >
                #{selectedDoc.id}
              </button>
            )}
          </div>
          <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 600 }}>
            {!documentId ? (
              <span className="text-dim text-mono">Select a document to inspect.</span>
            ) : detailLoading ? (
              <span className="text-dim text-mono">Loading...</span>
            ) : detail ? (
              <DocumentInspector detail={detail} chunks={detail.chunks} />
            ) : (
              <span className="text-dim text-mono">Document not found.</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

// --- Source Item ---

const SourceItem: React.FC<{
  source: CorpusSourceSummary;
  selected: boolean;
  onClick: () => void;
}> = ({ source, selected, onClick }) => {
  const pct = source.chunk_count > 0
    ? Math.round((source.embedded_count / source.chunk_count) * 100)
    : 0;

  return (
    <div
      onClick={onClick}
      style={{
        padding: '4px 6px',
        cursor: 'pointer',
        background: selected ? 'var(--mac-bg-dark)' : 'transparent',
        color: selected ? 'var(--mac-text-inv)' : 'inherit',
        borderBottom: '1px dotted #CCCCCC',
      }}
    >
      <div className="text-bold truncate" style={{ fontSize: 11 }}>{source.source_name}</div>
      <div className="text-mono" style={{ fontSize: 10, color: selected ? '#AAAAAA' : undefined }}>
        {source.document_count} docs · {source.word_count.toLocaleString()} words
      </div>
      {source.chunk_count > 0 && (
        <div className="text-mono" style={{ fontSize: 10, color: selected ? '#AAAAAA' : undefined }}>
          {source.embedded_count}/{source.chunk_count} embedded{' '}
          <span style={{ color: selected ? undefined : pct === 100 ? 'var(--mac-green)' : pct > 0 ? 'var(--mac-amber)' : 'var(--mac-text-dim)' }}>
            ({pct}%)
          </span>
        </div>
      )}
    </div>
  );
};

// --- Document Inspector ---

const DocumentInspector: React.FC<{
  detail: NonNullable<ReturnType<typeof useGetCorpusDocumentQuery>['data']>;
  chunks: CorpusChunk[];
}> = ({ detail, chunks }) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'text' | 'chunks' | 'coverage'>('overview');
  const [selectedChunkIdx, setSelectedChunkIdx] = useState<number | null>(null);

  const doc = detail.document;
  const metaKeys = Object.keys(doc.metadata || {});

  const embeddedCount = chunks.filter((c) => c.embedding?.present).length;
  const missingCount = chunks.length - embeddedCount;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 0 }}>
      {/* Tabs */}
      <div className="tab-bar">
        {(['overview', 'text', 'chunks', 'coverage'] as const).map((tab) => (
          <span
            key={tab}
            className={`tab-item ${activeTab === tab ? 'active' : ''}`}
            onClick={() => setActiveTab(tab)}
          >
            {tab}
          </span>
        ))}
      </div>

      <div style={{ padding: '6px 0' }}>
        {activeTab === 'overview' && (
          <>
            <div className="stat-grid">
              <span className="stat-label">ID</span>
              <span className="stat-value">
                {doc.id}
                <button className="copy-btn" onClick={() => navigator.clipboard.writeText(doc.id)} title="Copy"> [copy]</button>
              </span>
              <span className="stat-label">Source</span>
              <span className="stat-value">{doc.source_id}</span>
              <span className="stat-label">URL</span>
              <span className="stat-value">
                {doc.url ? (
                  <a href={doc.url} target="_blank" rel="noreferrer" className="accent" style={{ textDecoration: 'none' }}>
                    {doc.url}
                  </a>
                ) : '—'}
              </span>
              <span className="stat-label">Words</span>
              <span className="stat-value">{doc.word_count.toLocaleString()}</span>
              <span className="stat-label">Chunks</span>
              <span className="stat-value">{chunks.length}</span>
              <span className="stat-label">Embedded</span>
              <span className="stat-value">
                <span className={embeddedCount === chunks.length && chunks.length > 0 ? 'accent-green' : embeddedCount > 0 ? 'accent-amber' : ''}>
                  {embeddedCount}/{chunks.length}
                </span>
              </span>
              <span className="stat-label">Status</span>
              <span className="stat-value">
                <span className={`status-${doc.status === 'chunked' ? 'done' : doc.status}`}>{doc.status}</span>
              </span>
            </div>

            {metaKeys.length > 0 && (
              <>
                <div className="section-title" style={{ marginTop: 8 }}>Metadata</div>
                <div className="meta-grid">
                  {metaKeys.map((key) => (
                    <React.Fragment key={key}>
                      <span className="meta-key">{key}</span>
                      <span className="meta-value">{String(doc.metadata[key] ?? '')}</span>
                    </React.Fragment>
                  ))}
                </div>
              </>
            )}
          </>
        )}

        {activeTab === 'text' && (
          <>
            <div className="section-title">Extracted Text ({doc.word_count.toLocaleString()} words)</div>
            <div className="text-content">
              {doc.content_text || <span className="text-dim">No content text available.</span>}
            </div>
          </>
        )}

        {activeTab === 'chunks' && (
          <>
            <div className="section-title">
              Chunks ({chunks.length}) — {embeddedCount} embedded, {missingCount} missing
            </div>

            {/* Mini timeline bar */}
            <ChunkTimelineBar
              chunks={chunks}
              selectedIdx={selectedChunkIdx}
              onSelect={setSelectedChunkIdx}
            />

            <table className="data-table" style={{ marginTop: 4 }}>
              <thead>
                <tr>
                  <th>#</th>
                  <th>Range</th>
                  <th>Tokens</th>
                  <th>Embed</th>
                  <th>ID</th>
                </tr>
              </thead>
              <tbody>
                {chunks.map((chunk, idx) => (
                  <tr
                    key={chunk.id}
                    className={`selectable ${selectedChunkIdx === idx ? 'selected' : ''}`}
                    onClick={() => setSelectedChunkIdx(idx === selectedChunkIdx ? null : idx)}
                  >
                    <td className="mono">{chunk.chunk_index}</td>
                    <td className="mono">{chunk.start_offset}–{chunk.end_offset}</td>
                    <td className="num">{chunk.token_count}</td>
                    <td>
                      <span className={chunk.embedding?.present ? 'accent-green' : 'accent-dim'}>
                        {chunk.embedding?.present ? '●' : '○'}
                      </span>
                    </td>
                    <td className="mono">
                      <button
                        className="copy-btn"
                        onClick={(e) => {
                          e.stopPropagation();
                          navigator.clipboard.writeText(chunk.id);
                        }}
                        title="Copy chunk ID"
                      >
                        {chunk.id.slice(0, 12)}…
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Selected chunk text */}
            {selectedChunkIdx !== null && chunks[selectedChunkIdx] && (
              <div style={{ marginTop: 6 }}>
                <div className="section-title">
                  Chunk #{chunks[selectedChunkIdx].chunk_index} — {chunks[selectedChunkIdx].token_count} tokens
                </div>
                <div className="text-content">
                  {chunks[selectedChunkIdx].text}
                </div>
              </div>
            )}
          </>
        )}

        {activeTab === 'coverage' && (
          <>
            <div className="section-title">
              Embedding Coverage — {identity.label}
            </div>
            <div className="stat-grid" style={{ marginBottom: 8 }}>
              <span className="stat-label">Total</span>
              <span className="stat-value">{chunks.length}</span>
              <span className="stat-label">Embedded</span>
              <span className="stat-value accent-green">{embeddedCount}</span>
              <span className="stat-label">Missing</span>
              <span className="stat-value accent-amber">{missingCount}</span>
              <span className="stat-label">Coverage</span>
              <span className="stat-value">
                {chunks.length > 0 ? Math.round((embeddedCount / chunks.length) * 100) : 0}%
              </span>
            </div>

            {/* Coverage strip */}
            <div className="section-title">Per-chunk status</div>
            <div className="coverage-strip">
              {chunks.map((chunk) => (
                <span
                  key={chunk.id}
                  className={`coverage-dot ${chunk.embedding?.present ? 'present' : 'missing'}`}
                  title={`#${chunk.chunk_index} ${chunk.embedding?.present ? 'embedded' : 'missing'}`}
                />
              ))}
            </div>
            <div className="text-mono text-dim" style={{ marginTop: 4, fontSize: 10 }}>
              ● embedded · ○ missing
            </div>

            {/* Missing chunks list */}
            {missingCount > 0 && (
              <>
                <div className="section-title" style={{ marginTop: 8 }}>Missing chunks</div>
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>#</th>
                      <th>Range</th>
                      <th>Tokens</th>
                      <th>ID</th>
                    </tr>
                  </thead>
                  <tbody>
                    {chunks.filter((c) => !c.embedding?.present).map((chunk) => (
                      <tr key={chunk.id}>
                        <td className="mono">{chunk.chunk_index}</td>
                        <td className="mono">{chunk.start_offset}–{chunk.end_offset}</td>
                        <td className="num">{chunk.token_count}</td>
                        <td className="mono">{chunk.id.slice(0, 16)}…</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </>
            )}
          </>
        )}
      </div>
    </div>
  );
};

// --- Chunk Timeline Bar ---

const ChunkTimelineBar: React.FC<{
  chunks: CorpusChunk[];
  selectedIdx: number | null;
  onSelect: (idx: number | null) => void;
}> = ({ chunks, selectedIdx, onSelect }) => {
  if (chunks.length === 0) return null;

  const maxEnd = Math.max(...chunks.map((c) => c.end_offset));
  if (maxEnd <= 0) return null;

  return (
    <div className="chunk-bar-container">
      {chunks.map((chunk, idx) => {
        const left = (chunk.start_offset / maxEnd) * 100;
        const width = Math.max(((chunk.end_offset - chunk.start_offset) / maxEnd) * 100, 0.5);
        return (
          <div
            key={chunk.id}
            className={`chunk-bar ${chunk.embedding?.present ? 'embedded' : 'not-embedded'} ${selectedIdx === idx ? 'selected' : ''}`}
            style={{
              left: `${left}%`,
              width: `${width}%`,
            }}
            title={`#${chunk.chunk_index} ${chunk.start_offset}–${chunk.end_offset} ${chunk.embedding?.present ? '●' : '○'}`}
            onClick={() => onSelect(idx === selectedIdx ? null : idx)}
          />
        );
      })}
    </div>
  );
};

// Identity label for coverage tab
const identity = { label: 'openai/text-embedding-3-small @ 1536' };
