import React, { useState } from 'react';
import { PipelineView } from './components/pipeline/PipelineView';
import { EmbeddingsView } from './components/embeddings/EmbeddingsView';
import { SearchView } from './components/search/SearchView';
import { EvaluationView } from './components/evaluation/EvaluationView';
import { CorpusExplorerView } from './components/corpus/CorpusExplorerView';

const views = [
  { id: 'corpus', label: 'Corpus' },
  { id: 'pipeline', label: 'Pipeline' },
  { id: 'embeddings', label: 'Embeddings' },
  { id: 'search', label: 'Search' },
  { id: 'evaluation', label: 'Evaluation' },
];

export const App: React.FC = () => {
  const [activeView, setActiveView] = useState('corpus');

  const renderView = () => {
    switch (activeView) {
      case 'corpus':
        return <CorpusExplorerView />;
      case 'pipeline':
        return <PipelineView />;
      case 'embeddings':
        return <EmbeddingsView />;
      case 'search':
        return <SearchView />;
      case 'evaluation':
        return <EvaluationView />;
      default:
        return <CorpusExplorerView />;
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
