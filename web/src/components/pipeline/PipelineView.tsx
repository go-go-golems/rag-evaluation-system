import React from 'react';
import { useListSourcesQuery, useListDocumentsQuery } from '../../services/api';

export const PipelineView: React.FC = () => {
  const { data: sources, isLoading: sourcesLoading } = useListSourcesQuery();
  const { data: documents, isLoading: docsLoading } = useListDocumentsQuery();

  return (
    <div style={{ display: 'flex', gap: 8 }}>
      <div className="panel" style={{ flex: 1 }}>
        <div className="panel-header"><span>Sources</span></div>
        <div className="panel-body-condensed">
          {sourcesLoading ? (
            <span className="text-dim text-mono">Loading...</span>
          ) : sources && sources.length > 0 ? (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Type</th>
                  <th>Created</th>
                </tr>
              </thead>
              <tbody>
                {sources.map((s: { id: string; name: string; type: string; created_at: string }) => (
                  <tr key={s.id}>
                    <td>{s.name}</td>
                    <td className="mono">{s.type}</td>
                    <td className="mono">{s.created_at}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <span className="text-dim">No sources yet.</span>
          )}
        </div>
      </div>

      <div className="panel" style={{ flex: 1 }}>
        <div className="panel-header"><span>Documents</span></div>
        <div className="panel-body-condensed">
          {docsLoading ? (
            <span className="text-dim text-mono">Loading...</span>
          ) : documents && documents.length > 0 ? (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Title</th>
                  <th>Status</th>
                  <th>Words</th>
                </tr>
              </thead>
              <tbody>
                {documents.map((d: { id: string; title: string; status: string; word_count: number }) => (
                  <tr key={d.id}>
                    <td>{d.title}</td>
                    <td className={`status-${d.status === 'indexed' ? 'done' : d.status}`}>{d.status}</td>
                    <td className="num">{d.word_count}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <span className="text-dim">No documents indexed yet.</span>
          )}
        </div>
      </div>
    </div>
  );
};
