import React from 'react';

export const EvaluationView: React.FC = () => {
  return (
    <div className="panel">
      <div className="panel-header"><span>Evaluation Dashboard</span></div>
      <div className="panel-body">
        <p className="text-dim">Evaluation dashboard coming soon. This view will let you:</p>
        <ul style={{ fontSize: 12, paddingLeft: 16 }}>
          <li>Create and manage evaluation query sets</li>
          <li>Run evaluations across different search configurations</li>
          <li>Compare Recall@K, MRR, nDCG@K metrics</li>
          <li>See per-query pass/fail matrices</li>
        </ul>
      </div>
    </div>
  );
};
