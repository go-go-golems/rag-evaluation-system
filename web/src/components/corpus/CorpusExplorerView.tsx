import React, { useState, useMemo, useCallback, useEffect } from 'react';
import {
  useListCorpusSourcesQuery,
  useListCorpusDocumentsQuery,
  useGetCorpusDocumentQuery,
  useListChunkingStrategiesQuery,
  useGetDocumentProcessingCoverageQuery,
  CorpusIdentityArgs,
  DocumentProcessingCoverageItem,
} from '../../services/api';
import { IdentityBar } from './IdentityBar';
import { SourcePanel } from './SourcePanel';
import { DocumentBrowser } from './DocumentBrowser';
import { DocumentInspector } from './DocumentInspector';

export interface ChunkNavigationTarget {
  documentId: string;
  chunkId: string;
  sourceId: string;
}

interface CorpusExplorerViewProps {
  initialTarget?: ChunkNavigationTarget | null;
  onTargetConsumed?: () => void;
}

const DEFAULT_IDENTITY: CorpusIdentityArgs = {
  strategy_id: 'fixed-1200-150',
  provider_type: 'openai',
  model: 'text-embedding-3-small',
  dimensions: 1536,
};

const PAGE_SIZE = 100;

export const CorpusExplorerView: React.FC<CorpusExplorerViewProps> = ({ initialTarget, onTargetConsumed }) => {
  const [identity, setIdentity] = useState<CorpusIdentityArgs>(DEFAULT_IDENTITY);
  const [sourceId, setSourceId] = useState('');
  const [documentId, setDocumentId] = useState('');
  const [highlightChunkId, setHighlightChunkId] = useState<string | null>(null);
  const [pageLimit, setPageLimit] = useState(PAGE_SIZE);

  // Handle navigation from Search workbench
  useEffect(() => {
    if (initialTarget) {
      setSourceId(initialTarget.sourceId);
      setDocumentId(initialTarget.documentId);
      setHighlightChunkId(initialTarget.chunkId);
      onTargetConsumed?.();
    }
  }, [initialTarget, onTargetConsumed]);

  const { data: strategies = [] } = useListChunkingStrategiesQuery();
  const { data: sources = [], isLoading: sourcesLoading } = useListCorpusSourcesQuery(identity);
  const { data: documents = [], isLoading: docsLoading } = useListCorpusDocumentsQuery(
    { ...identity, source_id: sourceId, limit: pageLimit },
    { skip: !sourceId },
  );
  const { data: detail, isLoading: detailLoading } = useGetCorpusDocumentQuery(
    { ...identity, document_id: documentId, include_text: true },
    { skip: !documentId },
  );

  // Preprocessing coverage — uses default identity (fake provider from workflow)
  const { data: preprocessingCoverageData } = useGetDocumentProcessingCoverageQuery({
    artifact_type: 'clean_text',
    prompt_version: 'v1',
    provider: 'fake',
    model: 'fake-document-processor',
  });
  const preprocessingCoverage = useMemo(() => {
    const map: Record<string, DocumentProcessingCoverageItem> = {};
    for (const item of preprocessingCoverageData?.items ?? []) {
      map[item.source_id] = item;
    }
    return map;
  }, [preprocessingCoverageData]);

  const selectedSource = useMemo(
    () => sources.find((s) => s.source_id === sourceId),
    [sources, sourceId],
  );

  const selectedDoc = useMemo(
    () => documents.find((d) => d.id === documentId),
    [documents, documentId],
  );

  const handleSelectSource = useCallback((id: string) => {
    setSourceId(id);
    setDocumentId('');
    setHighlightChunkId(null);
    setPageLimit(PAGE_SIZE);
  }, []);

  const handleSelectDocument = useCallback((id: string) => {
    setDocumentId(id);
    setHighlightChunkId(null);
  }, []);

  const handleLoadMore = useCallback(() => {
    setPageLimit((prev) => prev + PAGE_SIZE);
  }, []);

  const totalDocs = useMemo(() => sources.reduce((s, x) => s + x.document_count, 0), [sources]);
  const totalChunks = useMemo(() => sources.reduce((s, x) => s + x.chunk_count, 0), [sources]);
  const totalEmbedded = useMemo(() => sources.reduce((s, x) => s + x.embedded_count, 0), [sources]);

  const hasMore = selectedSource ? documents.length < selectedSource.document_count : false;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      <IdentityBar
        identity={identity}
        onChange={setIdentity}
        strategies={strategies}
        totalDocs={totalDocs}
        totalChunks={totalChunks}
        totalEmbedded={totalEmbedded}
      />

      <div style={{ display: 'flex', gap: 8, minHeight: 500 }}>
        <SourcePanel
          sources={sources}
          isLoading={sourcesLoading}
          selectedId={sourceId}
          onSelect={handleSelectSource}
          preprocessingCoverage={preprocessingCoverage}
        />

        <DocumentBrowser
          documents={documents}
          isLoading={docsLoading}
          sourceName={selectedSource?.source_name ?? null}
          totalDocs={selectedSource?.document_count ?? 0}
          selectedDocId={documentId}
          onSelectDocument={handleSelectDocument}
          onLoadMore={handleLoadMore}
          hasMore={hasMore}
        />

        <div className="panel" style={{ flex: 1, minWidth: 0 }}>
          <div className="panel-header">
            <span>{selectedDoc ? selectedDoc.title : 'Document Inspector'}</span>
            {selectedDoc && (
              <button
                className="copy-btn"
                title="Copy document ID"
                onClick={() => navigator.clipboard.writeText(selectedDoc.id)}
              >
                #{selectedDoc.id}
              </button>
            )}
          </div>
          <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 600 }}>
            {!documentId ? (
              <span className="text-dim text-mono">Select a document to inspect.</span>
            ) : detailLoading ? (
              <span className="text-dim text-mono">Loading...</span>
            ) : detail ? (
              <DocumentInspector detail={detail} chunks={detail.chunks} identity={identity} highlightChunkId={highlightChunkId} />
            ) : (
              <span className="text-dim text-mono">Document not found.</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
