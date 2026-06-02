import React from 'react';
import { Caption } from '../../foundation';
import { Panel, ScrollRegion } from '../../layout';
import {
  CorpusSourceSummary,
  DocumentProcessingCoverageItem,
} from '../../../services/api';
import styles from './SourcePanel.module.css';

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
  const prePct = preprocessing && preprocessing.document_count > 0
    ? Math.round((preprocessing.artifact_count / preprocessing.document_count) * 100)
    : 0;

  return (
    <div onClick={onClick} className={[styles.item, selected ? styles.selected : ''].filter(Boolean).join(' ')}>
      <div className={styles.title}>{source.source_name}</div>
      <Caption className={selected ? styles.selectedMeta : undefined}>
        {source.document_count} docs · {source.word_count.toLocaleString()} words
      </Caption>
      {source.chunk_count > 0 && (
        <Caption className={selected ? styles.selectedMeta : undefined} tone={selected ? 'inherit' : pct === 100 ? 'success' : pct > 0 ? 'warning' : 'muted'}>
          {source.embedded_count}/{source.chunk_count} embedded ({pct}%)
        </Caption>
      )}
      {preprocessing && (
        <Caption className={selected ? styles.selectedMeta : undefined} tone={selected ? 'inherit' : preprocessing.missing_count === 0 ? 'success' : preprocessing.artifact_count > 0 ? 'warning' : 'muted'}>
          {preprocessing.artifact_count}/{preprocessing.document_count} preprocessed ({prePct}%)
        </Caption>
      )}
    </div>
  );
};

export const SourcePanel: React.FC<SourcePanelProps> = ({ sources, isLoading, selectedId, onSelect, preprocessingCoverage }) => {
  return (
    <Panel className={styles.root} title="Sources" actions={<Caption>{sources.length}</Caption>} density="condensed">
      <ScrollRegion axis="y" style={{ maxHeight: 600 }}>
        {isLoading ? (
          <Caption>Loading...</Caption>
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
      </ScrollRegion>
    </Panel>
  );
};
