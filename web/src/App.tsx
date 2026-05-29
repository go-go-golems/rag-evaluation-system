import React, { useState, useEffect, useCallback } from 'react';
import { PipelineView } from './components/pipeline/PipelineView';
import { EmbeddingsView } from './components/embeddings/EmbeddingsView';
import { SearchView } from './components/search/SearchView';
import { EvaluationView } from './components/evaluation/EvaluationView';
import { CorpusExplorerView } from './components/corpus/CorpusExplorerView';

const views = [
  { id: 'search', label: 'Search' },
  { id: 'corpus', label: 'Corpus' },
  { id: 'pipeline', label: 'Pipeline' },
  { id: 'embeddings', label: 'Embeddings' },
  { id: 'evaluation', label: 'Evaluation' },
];

export interface ChunkNavigationTarget {
  documentId: string;
  chunkId: string;
  sourceId: string;
}

export const App: React.FC = () => {
  const [activeView, setActiveView] = useState('search');
  const [chunkTarget, setChunkTarget] = useState<ChunkNavigationTarget | null>(null);

  // Listen for navigation events from SearchView
  const handleNavigateToChunk = useCallback((e: Event) => {
    const detail = (e as CustomEvent<ChunkNavigationTarget>).detail;
    setChunkTarget(detail);
    setActiveView('corpus');
  }, []);

  useEffect(() => {
    window.addEventListener('rag:navigate-to-chunk', handleNavigateToChunk);
    return () => window.removeEventListener('rag:navigate-to-chunk', handleNavigateToChunk);
  }, [handleNavigateToChunk]);

  // Clear target once consumed
  const handleChunkTargetConsumed = useCallback(() => {
    setChunkTarget(null);
  }, []);

  const renderView = () => {
    switch (activeView) {
      case 'corpus':
        return <CorpusExplorerView initialTarget={chunkTarget} onTargetConsumed={handleChunkTargetConsumed} />;
      case 'pipeline':
        return <PipelineView />;
      case 'embeddings':
        return <EmbeddingsView />;
      case 'search':
        return <SearchView />;
      case 'evaluation':
        return <EvaluationView />;
      default:
        return <SearchView />;
    }
  };

  return (
    <div style={{ minHeight: '100vh', background: 'var(--mac-bg)' }}>
      <nav className="nav-strip">
        <span className="nav-brand">◉ RAG Eval</span>
        {views.map((view) => (
          <span
            key={view.id}
            className={`nav-item ${activeView === view.id ? 'active' : ''}`}
            onClick={() => setActiveView(view.id)}
          >
            {view.label}
          </span>
        ))}
      </nav>
      <div style={{ padding: '8px' }}>
        {renderView()}
      </div>
    </div>
  );
};
