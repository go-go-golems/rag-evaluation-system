import type { KeyboardEvent } from 'react';
import { Button } from '../../atoms';
import { Caption } from '../../foundation';
import { FormRow, Panel, ScrollRegion, Stack } from '../../layout';
import styles from './SearchControlsPanel.module.css';

export type RetrieverType = 'bm25' | 'vector' | 'hybrid';

export interface SearchControlsPanelProps {
  query: string;
  setQuery: (q: string) => void;
  retriever: RetrieverType;
  setRetriever: (r: RetrieverType) => void;
  indexId: string;
  setIndexId: (v: string) => void;
  strategyId: string;
  setStrategyId: (v: string) => void;
  profile: string;
  setProfile: (v: string) => void;
  limit: number;
  setLimit: (v: number) => void;
  candidateLimit: number;
  setCandidateLimit: (v: number) => void;
  previewRunes: number;
  setPreviewRunes: (v: number) => void;
  sources: { source_id: string; source_name: string }[];
  selectedSourceIds: string[];
  toggleSource: (id: string) => void;
  onSearch: () => void;
  isLoading: boolean;
  onKeyDown: (e: KeyboardEvent) => void;
}

export function SearchControlsPanel({
  query,
  setQuery,
  retriever,
  setRetriever,
  indexId,
  setIndexId,
  strategyId,
  setStrategyId,
  profile,
  setProfile,
  limit,
  setLimit,
  candidateLimit,
  setCandidateLimit,
  previewRunes,
  setPreviewRunes,
  sources,
  selectedSourceIds,
  toggleSource,
  onSearch,
  isLoading,
  onKeyDown,
}: SearchControlsPanelProps) {
  return (
    <Panel title="Search" density="condensed" data-rag-component="SearchControlsPanel">
      <Stack gap="sm">
        <div className={styles.queryRow}>
          <input
            className={`input ${styles.queryInput}`}
            placeholder="Enter query…"
            value={query}
            onChange={e => setQuery(e.target.value)}
            onKeyDown={onKeyDown}
            autoFocus
          />
          <Button variant="primary" onClick={onSearch} disabled={isLoading || !query.trim()}>
            {isLoading ? '…' : '▶'}
          </Button>
        </div>

        <Caption transform="uppercase">Retriever</Caption>
        <div className={styles.retrieverGroup}>
          {(['bm25', 'vector', 'hybrid'] as RetrieverType[]).map(r => (
            <Button
              key={r}
              className={styles.retrieverButton}
              selected={retriever === r}
              aria-pressed={retriever === r}
              onClick={() => setRetriever(r)}
            >
              {r}
            </Button>
          ))}
        </div>

        {(retriever === 'bm25' || retriever === 'hybrid') && (
          <FormRow
            label="BM25 Index"
            control={<input className={`input ${styles.fullInput}`} value={indexId} onChange={e => setIndexId(e.target.value)} placeholder="Index ID" />}
          />
        )}

        {(retriever === 'vector' || retriever === 'hybrid') && (
          <Stack gap="xs">
            <Caption transform="uppercase">Vector</Caption>
            <FormRow label="Strategy" control={<input className={`input ${styles.fullInput}`} value={strategyId} onChange={e => setStrategyId(e.target.value)} />} />
            <FormRow label="Profile" control={<input className={`input ${styles.fullInput}`} value={profile} onChange={e => setProfile(e.target.value)} />} />
          </Stack>
        )}

        <Caption transform="uppercase">Source Filter</Caption>
        <ScrollRegion axis="y" className={styles.sourceList}>
          {sources.map(s => (
            <label key={s.source_id} className={`checkbox-row ${styles.sourceLabel}`}>
              <input
                type="checkbox"
                checked={selectedSourceIds.includes(s.source_id)}
                onChange={() => toggleSource(s.source_id)}
              />
              <span className="truncate">{s.source_name}</span>
            </label>
          ))}
          {sources.length === 0 && <Caption>Loading sources…</Caption>}
        </ScrollRegion>

        <Caption transform="uppercase">Limits</Caption>
        <FormRow label="Limit" control={<input className={`input ${styles.narrowInput}`} type="number" value={limit} onChange={e => setLimit(Number(e.target.value))} min={1} max={100} />} />
        {(retriever === 'vector' || retriever === 'hybrid') && (
          <FormRow label="Candidates" control={<input className={`input ${styles.narrowInput}`} type="number" value={candidateLimit} onChange={e => setCandidateLimit(Number(e.target.value))} min={1} />} />
        )}
        <FormRow label="Preview" control={<input className={`input ${styles.narrowInput}`} type="number" value={previewRunes} onChange={e => setPreviewRunes(Number(e.target.value))} min={0} />} />
      </Stack>
    </Panel>
  );
}
