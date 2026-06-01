import React from 'react';
import { CorpusIdentityArgs } from '../../../services/api';

interface IdentityBarProps {
  identity: CorpusIdentityArgs;
  onChange: (identity: CorpusIdentityArgs) => void;
  strategies: { id: string }[];
  totalDocs: number;
  totalChunks: number;
  totalEmbedded: number;
}

export const IdentityBar: React.FC<IdentityBarProps> = ({
  identity, onChange, strategies, totalDocs, totalChunks, totalEmbedded,
}) => {
  return (
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
            onChange={(e) => onChange({ ...identity, strategy_id: e.target.value })}
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
            onChange={(e) => onChange({ ...identity, provider_type: e.target.value })}
            style={{ width: 80 }}
          />
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
          <span className="text-mono text-dim">Model</span>
          <input
            className="input"
            value={identity.model}
            onChange={(e) => onChange({ ...identity, model: e.target.value })}
            style={{ width: 160 }}
          />
        </label>
        <label style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
          <span className="text-mono text-dim">Dims</span>
          <input
            className="input"
            type="number"
            value={identity.dimensions}
            onChange={(e) => onChange({ ...identity, dimensions: Number(e.target.value) })}
            style={{ width: 60 }}
          />
        </label>
        <span className="text-mono text-dim" style={{ marginLeft: 'auto' }}>
          {totalDocs} docs · {totalChunks} chunks · {totalEmbedded} embedded
        </span>
      </div>
    </div>
  );
};
