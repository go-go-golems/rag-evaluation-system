import React, { useState } from 'react';
import {
  CorpusChunk,
  CorpusDocumentDetail,
  CorpusIdentityArgs,
  useGetDocumentProcessingArtifactsQuery,
  DocumentProcessingArtifact,
} from '../../services/api';
import { ChunkTimelineBar } from './ChunkTimelineBar';

function timeAgoShort(iso: string): string {
  const ms = Date.now() - new Date(iso).getTime();
  if (ms < 60000) return `${Math.floor(ms / 1000)}s`;
  if (ms < 3600000) return `${Math.floor(ms / 60000)}m`;
  if (ms < 86400000) return `${Math.floor(ms / 3600000)}h`;
  return `${Math.floor(ms / 86400000)}d`;
}

const ArtifactDetail: React.FC<{ artifacts: DocumentProcessingArtifact[] }> = ({ artifacts }) => {
  const [selectedIdx, setSelectedIdx] = useState<number | null>(null);
  const selected = selectedIdx !== null ? artifacts[selectedIdx] : null;

  return (
    <>
      {selected && (
        <div className="panel" style={{ borderLeft: '3px solid var(--mac-accent)', marginTop: 6 }}>
          <div className="panel-header">
            <span>Artifact: {selected.artifact_type} ({selected.prompt_version})</span>
            <button className="copy-btn" onClick={() => setSelectedIdx(null)}>✕</button>
          </div>
          <div className="panel-body-condensed" style={{ fontSize: 12 }}>
            <div className="meta-grid">
              <span className="meta-key">Provider</span>
              <span className="meta-value">{selected.provider}/{selected.model}</span>
              <span className="meta-key">Status</span>
              <span className={`status-${selected.status === 'fresh' ? 'done' : selected.status === 'failed' ? 'error' : 'pending'}`}>{selected.status}</span>
              <span className="meta-key">Hash</span>
              <span className="meta-value">{selected.input_hash}</span>
            </div>
            {selected.output_text && (
              <fieldset className="form-section" style={{ marginTop: 6 }}>
                <legend>Output Text</legend>
                <div className="text-content" style={{ maxHeight: 200 }}>{selected.output_text}</div>
              </fieldset>
            )}
            {selected.error_message && (
              <div className="error-box" style={{ marginTop: 6 }}>
                {selected.error_message}
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
};

interface DocumentInspectorProps {
  detail: CorpusDocumentDetail;
  chunks: CorpusChunk[];
  identity: CorpusIdentityArgs;
  highlightChunkId?: string | null;
}

export const DocumentInspector: React.FC<DocumentInspectorProps> = ({ detail, chunks, identity, highlightChunkId }) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'text' | 'chunks' | 'coverage' | 'artifacts'>(highlightChunkId ? 'chunks' : 'overview');
  const [selectedChunkIdx, setSelectedChunkIdx] = useState<number | null>(() => {
    if (!highlightChunkId) return null;
    const idx = (chunks ?? []).findIndex(c => c.id === highlightChunkId);
    return idx >= 0 ? idx : null;
  });

  const safeChunks = chunks ?? [];
  const doc = detail.document;
  const metaKeys = Object.keys(doc.metadata || {});

  const embeddedCount = safeChunks.filter((c) => c.embedding?.present).length;
  const missingCount = safeChunks.length - embeddedCount;

  const identityLabel = `${identity.provider_type || '?'}/${identity.model || '?'} @ ${identity.dimensions || '?'}`;

  // Fetch document processing artifacts when artifacts tab is active
  const { data: artifactsData, isLoading: artifactsLoading } = useGetDocumentProcessingArtifactsQuery(
    doc.id,
    { skip: activeTab !== 'artifacts' },
  );
  const artifacts = artifactsData?.items ?? [];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 0 }}>
      {/* Tabs */}
      <div className="tab-bar">
        {(['overview', 'text', 'chunks', 'coverage', 'artifacts'] as const).map((tab) => (
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
              <span className="stat-value">{safeChunks.length}</span>
              <span className="stat-label">Embedded</span>
              <span className="stat-value">
                <span className={embeddedCount === safeChunks.length && safeChunks.length > 0 ? 'accent-green' : embeddedCount > 0 ? 'accent-amber' : ''}>
                  {embeddedCount}/{safeChunks.length}
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
              Chunks ({safeChunks.length}) — {embeddedCount} embedded, {missingCount} missing · {safeChunks.filter(c => c.enrichment?.present).length} enriched
            </div>

            <ChunkTimelineBar
              chunks={safeChunks}
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
                  <th>Enrich</th>
                  <th>ID</th>
                </tr>
              </thead>
              <tbody>
                {safeChunks.map((chunk, idx) => (
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
                    <td>
                      <span className={chunk.enrichment?.present ? 'accent-green' : 'accent-dim'}
                        title={chunk.enrichment?.short_summary ?? ''}>
                        {chunk.enrichment?.present ? '●' : '○'}
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

            {selectedChunkIdx !== null && safeChunks[selectedChunkIdx] && (
              <div style={{ marginTop: 6 }}>
                <div className="section-title">
                  Chunk #{safeChunks[selectedChunkIdx].chunk_index} — {safeChunks[selectedChunkIdx].token_count} tokens
                </div>
                <div className="text-content">
                  {safeChunks[selectedChunkIdx].text}
                </div>
              </div>
            )}
          </>
        )}

        {activeTab === 'coverage' && (
          <>
            <div className="section-title">
              Embedding Coverage — {identityLabel}
            </div>
            <div className="stat-grid" style={{ marginBottom: 8 }}>
              <span className="stat-label">Total</span>
              <span className="stat-value">{safeChunks.length}</span>
              <span className="stat-label">Embedded</span>
              <span className="stat-value accent-green">{embeddedCount}</span>
              <span className="stat-label">Missing</span>
              <span className="stat-value accent-amber">{missingCount}</span>
              <span className="stat-label">Coverage</span>
              <span className="stat-value">
                {safeChunks.length > 0 ? Math.round((embeddedCount / safeChunks.length) * 100) : 0}%
              </span>
            </div>

            <div className="section-title">Per-chunk status</div>
            <div className="coverage-strip">
              {safeChunks.map((chunk) => (
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
                    {safeChunks.filter((c) => !c.embedding?.present).map((chunk) => (
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

        {activeTab === 'artifacts' && (
          <>
            <div className="section-title">
              Document Processing Artifacts ({artifacts.length})
            </div>
            {artifactsLoading ? (
              <span className="text-dim text-mono">Loading...</span>
            ) : artifacts.length === 0 ? (
              <span className="text-dim text-mono">No preprocessing artifacts for this document.</span>
            ) : (
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Type</th>
                    <th>Version</th>
                    <th>Provider</th>
                    <th>Status</th>
                    <th>Age</th>
                  </tr>
                </thead>
                <tbody>
                  {artifacts.map((a: DocumentProcessingArtifact) => (
                    <tr key={`${a.artifact_type}-${a.prompt_version}-${a.provider}-${a.model}`}>
                      <td className="mono">{a.artifact_type}</td>
                      <td className="mono">{a.prompt_version}</td>
                      <td className="mono" style={{ fontSize: 10 }}>{a.provider}/{a.model}</td>
                      <td className={a.status === 'fresh' ? 'status-done' : a.status === 'failed' ? 'status-error' : 'status-pending'}>{a.status}</td>
                      <td className="mono" style={{ fontSize: 10 }}>{timeAgoShort(a.updated_at)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}

            {artifacts.length > 0 && (
              <ArtifactDetail artifacts={artifacts} />
            )}
          </>
        )}
      </div>
    </div>
  );
};
