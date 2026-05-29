import React from 'react';
import {
  CorpusSourceSummary,
  DocumentProcessingCoverageItem,
} from '../../services/api';

interface SourcePanelProps {
  sources: CorpusSourceSummary[];
  isLoading: boolean;
  selectedId: string;
  onSelect: (id: string) => void;
  /** Preprocessing coverage per source, keyed by source_id */
  preprocessingCoverage?: Record<string, DocumentProcessingCoverageItem>;
}

interface SourceItemProps {
  source: CorpusSourceSummary;
  selected: boolean;
  onClick: () => void;
  preprocessing?: DocumentProcessingCoverageItem;
}

const SourceItem: React.FC<SourceItemProps> = ({ source, selected, onClick, preprocessing }) => {
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
      {preprocessing && (
        <div className="text-mono" style={{ fontSize: 10, color: selected ? '#AAAAAA' : undefined }}>
          {preprocessing.artifact_count}/{preprocessing.document_count} preprocessed{' '}
          <span style={{ color: selected ? undefined : preprocessing.missing_count === 0 ? 'var(--mac-green)' : preprocessing.artifact_count > 0 ? 'var(--mac-amber)' : 'var(--mac-text-dim)' }}>
            ({preprocessing.document_count > 0 ? Math.round((preprocessing.artifact_count / preprocessing.document_count) * 100) : 0}%)
          </span>
        </div>
      )}
    </div>
  );
};

export const SourcePanel: React.FC<SourcePanelProps> = ({ sources, isLoading, selectedId, onSelect, preprocessingCoverage }) => {
  return (
    <div className="panel" style={{ width: 220, flexShrink: 0 }}>
      <div className="panel-header">
        <span>Sources</span>
        <span className="text-mono" style={{ fontSize: 10 }}>{sources.length}</span>
      </div>
      <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 600 }}>
        {isLoading ? (
          <span className="text-dim text-mono">Loading...</span>
        ) : (
          sources.map((src) => (
            <SourceItem
              key={src.source_id}
              source={src}
              selected={src.source_id === selectedId}
              onClick={() => onSelect(src.source_id)}
              preprocessing={preprocessingCoverage?.[src.source_id]}
            />
          ))
        )}
      </div>
    </div>
  );
};
