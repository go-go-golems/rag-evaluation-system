import { StatusText } from '@go-go-golems/rag-evaluation-site';
import { Panel } from '@go-go-golems/rag-evaluation-site';
import type { WorkflowOpGroup } from '../../../services/api';
import { friendlyOpName, workflowStatusToRagStatus } from '../../workflows/workflowFormat';
import styles from './WorkflowOpGraphPanel.module.css';

export interface WorkflowOpGraphPanelProps {
  groups: WorkflowOpGroup[];
  total: number;
  done: number;
}

export function WorkflowOpGraphPanel({ groups, total, done }: WorkflowOpGraphPanelProps) {
  const fanOutOps = groups.filter(group => group.operation === 'chunk_document' || group.operation === 'preprocess_document');
  const fanInOps = groups.filter(group => group.operation === 'compute_embeddings' || group.operation === 'build_bm25');
  const otherOps = groups.filter(group =>
    !fanOutOps.some(candidate => candidate.operation === group.operation && candidate.status === group.status) &&
    !fanInOps.some(candidate => candidate.operation === group.operation && candidate.status === group.status));

  const renderNode = (group: WorkflowOpGroup, includeCount = true) => (
    <span key={group.operation + '|' + group.status} className={styles.node}>
      <StatusText status={workflowStatusToRagStatus(group.status)} icon>
        {friendlyOpName(group.operation)}{includeCount ? ` × ${group.count}` : ''}
      </StatusText>
    </span>
  );

  return (
    <Panel title="Op Graph" density="condensed" data-rag-component="WorkflowOpGraphPanel">
      <div className={styles.graph}>
        <div className={styles.fanOut}>
          {fanOutOps.map(group => renderNode(group))}
          {otherOps.map(group => renderNode(group))}
        </div>
        <div className={styles.arrow}>↓ {done}/{total} completed</div>
        {fanInOps.length > 0 && <div className={styles.fanOut}>{fanInOps.map(group => renderNode(group, false))}</div>}
      </div>
    </Panel>
  );
}
