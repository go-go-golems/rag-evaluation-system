import { Button, IconButton } from '../../atoms';
import { Caption, StatusText } from '../../foundation';
import { Panel } from '../../layout';
import { MetadataGrid } from '../../molecules';
import type { WorkflowListItem } from '../../../services/api';
import { timeAgo, workflowStatusToRagStatus } from '../../workflows/workflowFormat';
import styles from './WorkflowSummaryPanel.module.css';

export interface WorkflowSummaryPanelProps {
  workflow: WorkflowListItem['workflow'];
  strategyId: string;
  docCount: number;
  sourceId?: string;
  totalOps: number;
  succeededCount: number;
  runningCount: number;
  pendingCount: number;
  failedCount: number;
  isRunning: boolean;
  onBack: () => void;
  onCancel: () => void;
  onNavigateToCorpus: (sourceId: string, strategyId: string) => void;
}

function ProgressBar({ done, total, running = 0, failed = 0 }: { done: number; total: number; running?: number; failed?: number }) {
  const pct = total > 0 ? (done / total) * 100 : 0;
  const runPct = total > 0 ? (running / total) * 100 : 0;
  const failPct = total > 0 ? (failed / total) * 100 : 0;

  return (
    <div className={styles.progressBar} data-rag-component="WorkflowProgressBar">
      <div className={`${styles.progressFill} ${styles.done}`} style={{ width: `${pct}%` }} />
      {running > 0 && <div className={`${styles.progressFill} ${styles.running}`} style={{ width: `${runPct}%`, left: `${pct}%` }} />}
      {failed > 0 && <div className={`${styles.progressFill} ${styles.failed}`} style={{ width: `${failPct}%`, left: `${pct + runPct}%` }} />}
      <span className={styles.progressLabel} style={{ color: pct > 50 ? 'var(--mac-text-inv)' : 'var(--mac-text)' }}>
        {done}/{total} ({Math.round(pct)}%)
      </span>
    </div>
  );
}

export function WorkflowSummaryPanel({
  workflow,
  strategyId,
  docCount,
  sourceId,
  totalOps,
  succeededCount,
  runningCount,
  pendingCount,
  failedCount,
  isRunning,
  onBack,
  onCancel,
  onNavigateToCorpus,
}: WorkflowSummaryPanelProps) {
  return (
    <Panel
      title={workflow.ID}
      actions={(
        <div className={styles.actions}>
          <StatusText status={workflowStatusToRagStatus(workflow.Status)} icon>{workflow.Status}</StatusText>
          <IconButton label="Back to workflow list" onClick={onBack}>← Back</IconButton>
        </div>
      )}
      density="condensed"
      data-rag-component="WorkflowSummaryPanel"
    >
      <ProgressBar done={succeededCount} total={totalOps} running={runningCount} failed={failedCount} />
      <div className={styles.metaRow}>
        <MetadataGrid density="compact" items={[
          { key: 'Strategy', value: <Caption>{strategyId}</Caption> },
          { key: 'Docs', value: docCount },
          { key: 'Age', value: <Caption>{timeAgo(workflow.CreatedAt)}</Caption> },
          ...(runningCount > 0 ? [{ key: 'Running', value: <StatusText status="running">{runningCount}</StatusText> }] : []),
          ...(pendingCount > 0 ? [{ key: 'Pending', value: <StatusText status="pending">{pendingCount}</StatusText> }] : []),
          ...(failedCount > 0 ? [{ key: 'Failed', value: <StatusText status="failed">{failedCount}</StatusText> }] : []),
        ]} />
        <span className={styles.spacer} />
        {isRunning && <Button onClick={onCancel}>Cancel</Button>}
        {workflow.Status === 'succeeded' && sourceId && <Button onClick={() => onNavigateToCorpus(sourceId, strategyId)}>View in Corpus →</Button>}
      </div>
    </Panel>
  );
}
