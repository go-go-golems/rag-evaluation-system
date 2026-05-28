import React from 'react';

export const SearchView: React.FC = () => {
  return (
    <div className="panel">
      <div className="panel-header"><span>Search Sandbox</span></div>
      <div className="panel-body">
        <p className="text-dim">Search sandbox coming soon. This view will let you:</p>
        <ul style={{ fontSize: 12, paddingLeft: 16 }}>
          <li>Run queries against BM25, vector, and hybrid indexes</li>
          <li>Inspect score breakdowns (BM25 rank, vector rank, RRF score)</li>
          <li>Compare results before and after reranking</li>
          <li>Visualize score distributions</li>
        </ul>
      </div>
    </div>
  );
};
