import { useEffect, useState } from 'react';
import { Button } from '../../atoms';
import { Caption, StatusText } from '../../foundation';
import { Panel, ScrollRegion, TabList } from '../../layout';
import { MetadataGrid, type MetadataGridItem } from '../../molecules';
import type { CorpusDocumentDetail, RetrievalResult } from '../../../services/api';
import styles from './ResultInspectorPanel.module.css';

export type ResultInspectorTab = 'detail' | 'chunk' | 'document';

export interface ResultInspectorPanelProps {
  result: RetrievalResult;
  docDetail: CorpusDocumentDetail | null | undefined;
  onClose: () => void;
  onOpenInCorpus: () => void;
  defaultTab?: ResultInspectorTab;
}

const tabs = [
  { id: 'detail', label: 'detail' },
  { id: 'chunk', label: 'chunk' },
  { id: 'document', label: 'document' },
] satisfies { id: ResultInspectorTab; label: string }[];

export function ResultInspectorPanel({ result, docDetail, onClose, onOpenInCorpus, defaultTab = 'detail' }: ResultInspectorPanelProps) {
  const [activeTab, setActiveTab] = useState<ResultInspectorTab>(defaultTab);

  useEffect(() => {
    setActiveTab(defaultTab);
  }, [defaultTab, result.chunk_id]);

  const currentChunk = docDetail?.chunks?.find(c => c.id === result.chunk_id);
  const chunkIndex = docDetail?.chunks?.findIndex(c => c.id === result.chunk_id) ?? -1;
  const prevChunk = chunkIndex > 0 ? docDetail!.chunks[chunkIndex - 1] : null;
  const nextChunk = chunkIndex >= 0 && chunkIndex < (docDetail?.chunks?.length ?? 0) - 1
    ? docDetail!.chunks[chunkIndex + 1]
    : null;

  const scoreItems: MetadataGridItem[] = [
    { key: 'Retriever', value: <StatusText status={result.retriever === 'bm25' ? 'warning' : result.retriever === 'vector' ? 'running' : 'done'}>{result.retriever.toUpperCase()}</StatusText> },
    { key: 'Score', value: result.score.toFixed(6) },
    ...(result.components?.bm25 ? [{ key: 'BM25', value: <Caption tone="warning">rank #{result.components.bm25.rank}, score {result.components.bm25.score.toFixed(4)}</Caption> }] : []),
    ...(result.components?.vector ? [{ key: 'Vector', value: <Caption tone="accent">rank #{result.components.vector.rank}, score {result.components.vector.score.toFixed(4)}</Caption> }] : []),
  ];

  const identityItems: MetadataGridItem[] = [
    { key: 'Title', value: result.title },
    { key: 'Chunk ID', value: result.chunk_id, copyValue: result.chunk_id },
    { key: 'Document ID', value: result.document_id, copyValue: result.document_id },
    { key: 'Source ID', value: result.source_id, copyValue: result.source_id },
    ...(result.url ? [{ key: 'URL', value: <Caption tone="accent">{result.url}</Caption> }] : []),
    { key: 'Chunk Index', value: result.chunk_index },
    { key: 'Strategy', value: result.strategy_id },
  ];

  const documentItems: MetadataGridItem[] = docDetail ? [
    { key: 'ID', value: docDetail.document.id, copyValue: docDetail.document.id },
    { key: 'Title', value: docDetail.document.title },
    { key: 'Source', value: docDetail.document.source_id },
    ...(docDetail.document.url ? [{ key: 'URL', value: <Caption tone="accent">{docDetail.document.url}</Caption> }] : []),
    { key: 'Words', value: docDetail.document.word_count },
    { key: 'Chunks', value: docDetail.chunks.length },
    { key: 'Status', value: docDetail.document.status },
  ] : [];

  const metadataItems: MetadataGridItem[] = docDetail?.document.metadata
    ? Object.entries(docDetail.document.metadata).slice(0, 12).map(([key, value]) => ({ key, value: String(value) }))
    : [];

  return (
    <Panel
      fill
      className={styles.root}
      title={`Inspector — Rank #${result.rank}`}
      actions={(
        <div className={styles.actions}>
          <Button size="compact" onClick={onOpenInCorpus}>Open in Corpus</Button>
          <Button size="compact" onClick={onClose}>✕</Button>
        </div>
      )}
      data-rag-component="ResultInspectorPanel"
    >
      <TabList items={tabs} activeId={activeTab} onChange={setActiveTab} ariaLabel="Result inspector tabs" />
      <ScrollRegion axis="y" className={styles.body}>
        {activeTab === 'detail' && (
          <>
            <section className={styles.section}>
              <div className={styles.sectionTitle}>Scores</div>
              <MetadataGrid items={scoreItems} />
            </section>
            <section className={styles.section}>
              <div className={styles.sectionTitle}>Identity</div>
              <MetadataGrid items={identityItems} />
            </section>
            <section className={styles.section}>
              <div className={styles.sectionTitle}>Preview</div>
              <div className={styles.textBlock}>{result.preview}</div>
            </section>
          </>
        )}

        {activeTab === 'chunk' && (
          <>
            <section className={styles.section}>
              <div className={styles.sectionTitle}>Chunk Text</div>
              {currentChunk ? (
                <>
                  <MetadataGrid items={[
                    { key: 'Offsets', value: `${currentChunk.start_offset}–${currentChunk.end_offset}` },
                    { key: 'Tokens', value: currentChunk.token_count },
                    { key: 'Embedding', value: <StatusText status={currentChunk.embedding?.present ? 'done' : 'error'}>{currentChunk.embedding?.present ? '● Present' : '○ Missing'}</StatusText> },
                  ]} />
                  <div className={styles.textBlock}>{currentChunk.text}</div>
                </>
              ) : (
                <div className={styles.empty}>Loading chunk text…</div>
              )}
            </section>
            {(prevChunk || nextChunk) && (
              <section className={styles.section}>
                <div className={styles.sectionTitle}>Neighbors</div>
                {prevChunk && (
                  <div>
                    <Caption>Previous chunk #{prevChunk.chunk_index}</Caption>
                    <div className={styles.neighborBlock}>{prevChunk.text.slice(0, 300)}{prevChunk.text.length > 300 ? '…' : ''}</div>
                  </div>
                )}
                {nextChunk && (
                  <div>
                    <Caption>Next chunk #{nextChunk.chunk_index}</Caption>
                    <div className={styles.neighborBlock}>{nextChunk.text.slice(0, 300)}{nextChunk.text.length > 300 ? '…' : ''}</div>
                  </div>
                )}
              </section>
            )}
          </>
        )}

        {activeTab === 'document' && (
          <>
            {docDetail ? (
              <>
                <section className={styles.section}>
                  <div className={styles.sectionTitle}>Document</div>
                  <MetadataGrid items={documentItems} />
                </section>
                {docDetail.chunks.length > 0 && (
                  <section className={styles.section}>
                    <div className={styles.sectionTitle}>Chunk Coverage</div>
                    <div className={styles.coverageStrip}>
                      {docDetail.chunks.map((ch) => {
                        const isCurrentChunk = ch.id === result.chunk_id;
                        const hasEmbedding = ch.embedding?.present;
                        return (
                          <span
                            key={ch.id}
                            title={`Chunk ${ch.chunk_index}${isCurrentChunk ? ' (selected)' : ''}`}
                            className={[
                              styles.coverageDot,
                              isCurrentChunk ? styles.currentChunk : hasEmbedding ? styles.embeddedChunk : styles.missingChunk,
                            ].filter(Boolean).join(' ')}
                          />
                        );
                      })}
                    </div>
                    <Caption>{docDetail.chunks.filter(c => c.embedding?.present).length}/{docDetail.chunks.length} embedded · <span className="accent">■</span> = current</Caption>
                  </section>
                )}
                {metadataItems.length > 0 && (
                  <section className={styles.section}>
                    <div className={styles.sectionTitle}>Metadata</div>
                    <MetadataGrid items={metadataItems} />
                  </section>
                )}
              </>
            ) : (
              <div className={styles.empty}>Loading document…</div>
            )}
          </>
        )}
      </ScrollRegion>
    </Panel>
  );
}
