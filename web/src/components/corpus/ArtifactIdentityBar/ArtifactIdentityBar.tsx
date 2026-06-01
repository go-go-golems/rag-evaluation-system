import React from 'react';
import {
  DocumentProcessingIdentity,
  ChunkEnrichmentIdentity,
  DocumentProcessingCoverageArgs,
  ChunkEnrichmentCoverageArgs,
} from '../../../services/api';

// ─── Document Processing Identity Selector ────────────────────────────────

interface DocProcessingIdentityBarProps {
  identities: DocumentProcessingIdentity[];
  selected: DocumentProcessingCoverageArgs;
  onChange: (args: DocumentProcessingCoverageArgs) => void;
}

export const DocProcessingIdentityBar: React.FC<DocProcessingIdentityBarProps> = ({
  identities,
  selected,
  onChange,
}) => {
  if (identities.length === 0) {
    return (
      <div className="text-dim text-mono" style={{ fontSize: 10, padding: '2px 4px' }}>
        No preprocessing identities found
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', gap: 4, alignItems: 'center', flexWrap: 'wrap', padding: '2px 4px' }}>
      <span className="text-mono" style={{ fontSize: 10, fontWeight: 'bold' }}>Preprocessing:</span>
      {identities.map((id) => {
        const key = `${id.artifact_type}|${id.prompt_version}|${id.provider}|${id.model}`;
        const isSelected =
          selected.artifact_type === id.artifact_type &&
          selected.prompt_version === id.prompt_version &&
          selected.provider === id.provider &&
          selected.model === id.model;
        return (
          <button
            key={key}
            className={`btn ${isSelected ? 'btn-primary' : ''}`}
            style={{ fontSize: 10, padding: '1px 6px' }}
            onClick={() => onChange({
              artifact_type: id.artifact_type,
              prompt_version: id.prompt_version,
              provider: id.provider,
              model: id.model,
            })}
            title={`${id.provider}/${id.model} — ${id.artifact_count} artifacts`}
          >
            {id.artifact_type}/{id.prompt_version} ({id.provider}/{id.model}) [{id.artifact_count}]
          </button>
        );
      })}
    </div>
  );
};

// ─── Chunk Enrichment Identity Selector ────────────────────────────────────

interface ChunkEnrichmentIdentityBarProps {
  identities: ChunkEnrichmentIdentity[];
  selected: ChunkEnrichmentCoverageArgs;
  onChange: (args: ChunkEnrichmentCoverageArgs) => void;
}

export const ChunkEnrichmentIdentityBar: React.FC<ChunkEnrichmentIdentityBarProps> = ({
  identities,
  selected,
  onChange,
}) => {
  if (identities.length === 0) {
    return null;
  }

  return (
    <div style={{ display: 'flex', gap: 4, alignItems: 'center', flexWrap: 'wrap', padding: '2px 4px' }}>
      <span className="text-mono" style={{ fontSize: 10, fontWeight: 'bold' }}>Enrichment:</span>
      {identities.map((id) => {
        const key = `${id.strategy_id}|${id.prompt_version}`;
        const isSelected =
          selected.strategy_id === id.strategy_id &&
          selected.prompt_version === id.prompt_version;
        return (
          <button
            key={key}
            className={`btn ${isSelected ? 'btn-primary' : ''}`}
            style={{ fontSize: 10, padding: '1px 6px' }}
            onClick={() => onChange({
              strategy_id: id.strategy_id,
              prompt_version: id.prompt_version,
            })}
            title={`${id.provider}/${id.model} — ${id.enriched_count} enriched`}
          >
            {id.strategy_id}/{id.prompt_version} ({id.provider}/{id.model}) [{id.enriched_count}]
          </button>
        );
      })}
    </div>
  );
};
