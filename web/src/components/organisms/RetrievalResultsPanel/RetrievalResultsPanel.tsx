import type { ReactNode } from 'react';
import { ErrorCallout } from '@go-go-golems/rag-evaluation-site';
import { Caption } from '@go-go-golems/rag-evaluation-site';
import { Panel, ScrollRegion } from '@go-go-golems/rag-evaluation-site';
import { DataTable, type DataTableColumn } from '@go-go-golems/rag-evaluation-site';
import type { RetrievalResult, SearchResult } from '../../../services/api';
import type { RetrieverType } from '../SearchControlsPanel';
import styles from './RetrievalResultsPanel.module.css';

export interface RetrievalResultsPanelProps {
  searchResult: SearchResult | null;
  items: RetrievalResult[];
  retriever: RetrieverType;
  selectedResult: RetrievalResult | null;
  searchError: string | null;
  isLoading: boolean;
  onSelectResult: (result: RetrievalResult | null) => void;
  emptyHint?: ReactNode;
}

export function RetrievalResultsPanel({
  searchResult,
  items,
  retriever,
  selectedResult,
  searchError,
  isLoading,
  onSelectResult,
  emptyHint,
}: RetrievalResultsPanelProps) {
  const columns: DataTableColumn<RetrievalResult>[] = [
    { id: 'rank', header: '#', align: 'end', cell: (item) => item.rank },
    { id: 'title', header: 'Title', maxWidth: 160, cell: (item) => item.title },
    { id: 'chunkIndex', header: 'Idx', align: 'end', cell: (item) => item.chunk_index },
    { id: 'score', header: 'Score', align: 'end', cell: (item) => item.score.toFixed(4) },
    ...(retriever === 'hybrid'
      ? [
          {
            id: 'bm25',
            header: 'BM25',
            align: 'end' as const,
            cell: (item: RetrievalResult) => item.components?.bm25
              ? <Caption tone="warning">#{item.components.bm25.rank} ({item.components.bm25.score.toFixed(3)})</Caption>
              : <Caption>—</Caption>,
          },
          {
            id: 'vector',
            header: 'Vec',
            align: 'end' as const,
            cell: (item: RetrievalResult) => item.components?.vector
              ? <Caption tone="accent">#{item.components.vector.rank} ({item.components.vector.score.toFixed(3)})</Caption>
              : <Caption>—</Caption>,
          },
        ]
      : []),
    { id: 'preview', header: 'Preview', maxWidth: 300, cell: (item) => item.preview },
  ];

  return (
    <Panel
      fill
      title={searchResult ? `${searchResult.retriever.toUpperCase()} — ${items.length} results` : 'Results'}
      actions={searchResult && <Caption>{searchResult.query}</Caption>}
      data-rag-component="RetrievalResultsPanel"
    >
      <ScrollRegion axis="y" style={{ height: '100%' }}>
        {searchError && <ErrorCallout className={styles.error}>{searchError}</ErrorCallout>}
        {items.length === 0 && !searchError && !isLoading && (
          <Caption className={styles.centeredEmpty}>
            {emptyHint ?? (searchResult ? 'No results found.' : 'Enter a query and press Search.')}
          </Caption>
        )}
        {isLoading && (
          <Caption className={styles.centeredEmpty}>
            Searching…
          </Caption>
        )}
        {items.length > 0 && (
          <DataTable
            rows={items}
            columns={columns}
            getRowKey={(item) => item.chunk_id}
            selectedKey={selectedResult?.chunk_id ?? null}
            onRowSelect={(item) => onSelectResult(selectedResult?.chunk_id === item.chunk_id ? null : item)}
          />
        )}
      </ScrollRegion>
    </Panel>
  );
}
