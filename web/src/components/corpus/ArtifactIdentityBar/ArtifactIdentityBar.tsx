import React from 'react';
import { Button } from '../../atoms';
import { Caption } from '../../foundation';
import { Inline } from '../../layout';
import {
  DocumentProcessingIdentity,
  ChunkEnrichmentIdentity,
  DocumentProcessingCoverageArgs,
  ChunkEnrichmentCoverageArgs,
} from '../../../services/api';
import styles from './ArtifactIdentityBar.module.css';

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
      <Caption className={styles.empty}>
        No preprocessing identities found
      </Caption>
    );
  }

  return (
    <Inline gap="xs" className={styles.root} data-rag-component="DocProcessingIdentityBar">
      <Caption tone="inherit" className={styles.label}>Preprocessing:</Caption>
      {identities.map((identity) => {
        const key = `${identity.artifact_type}|${identity.prompt_version}|${identity.provider}|${identity.model}`;
        const isSelected =
          selected.artifact_type === identity.artifact_type &&
          selected.prompt_version === identity.prompt_version &&
          selected.provider === identity.provider &&
          selected.model === identity.model;
        return (
          <Button
            key={key}
            size="compact"
            selected={isSelected}
            aria-pressed={isSelected}
            onClick={() => onChange({
              artifact_type: identity.artifact_type,
              prompt_version: identity.prompt_version,
              provider: identity.provider,
              model: identity.model,
            })}
            title={`${identity.provider}/${identity.model} — ${identity.artifact_count} artifacts`}
          >
            {identity.artifact_type}/{identity.prompt_version} ({identity.provider}/{identity.model}) [{identity.artifact_count}]
          </Button>
        );
      })}
    </Inline>
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
    <Inline gap="xs" className={styles.root} data-rag-component="ChunkEnrichmentIdentityBar">
      <Caption tone="inherit" className={styles.label}>Enrichment:</Caption>
      {identities.map((identity) => {
        const key = `${identity.strategy_id}|${identity.prompt_version}`;
        const isSelected =
          selected.strategy_id === identity.strategy_id &&
          selected.prompt_version === identity.prompt_version;
        return (
          <Button
            key={key}
            size="compact"
            selected={isSelected}
            aria-pressed={isSelected}
            onClick={() => onChange({
              strategy_id: identity.strategy_id,
              prompt_version: identity.prompt_version,
            })}
            title={`${identity.provider}/${identity.model} — ${identity.enriched_count} enriched`}
          >
            {identity.strategy_id}/{identity.prompt_version} ({identity.provider}/{identity.model}) [{identity.enriched_count}]
          </Button>
        );
      })}
    </Inline>
  );
};
