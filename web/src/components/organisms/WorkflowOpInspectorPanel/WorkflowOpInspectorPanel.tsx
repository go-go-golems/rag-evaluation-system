import type { ReactNode } from 'react';
import { Button } from '../../atoms';
import { Caption, StatusText } from '../../foundation';
import { Panel } from '../../layout';
import { MetadataGrid } from '../../molecules';
import type { WorkflowOp } from '../../../services/api';
import { friendlyOpName, workflowStatusToRagStatus } from '../../workflows/workflowFormat';
import styles from './WorkflowOpInspectorPanel.module.css';

export interface WorkflowOpInspectorPanelProps {
  sample: WorkflowOp;
  retrying?: boolean;
  onClose: () => void;
  onRetry?: () => void;
  children?: ReactNode;
}

export function WorkflowOpInspectorPanel({ sample, retrying = false, onClose, onRetry, children }: WorkflowOpInspectorPanelProps) {
  const inputItems = Object.entries(sample.op.Input as Record<string, unknown>).map(([key, value]) => ({
    key,
    value: typeof value === 'string' ? value : JSON.stringify(value),
  }));
  const canRetry = sample.status === 'failed' && sample.op.RetryState.LastError && onRetry;

  return (
    <Panel
      className={styles.root}
      title={`Sample Op: ${friendlyOpName(sample.op.Kind)}`}
      actions={<button className="copy-btn" onClick={onClose}>✕</button>}
      density="condensed"
      data-rag-component="WorkflowOpInspectorPanel"
    >
      <MetadataGrid items={[
        { key: 'ID', value: sample.op.ID, copyValue: sample.op.ID },
        { key: 'Status', value: <StatusText status={workflowStatusToRagStatus(sample.status)} icon>{sample.status}</StatusText> },
        { key: 'Queue', value: sample.op.Queue },
        { key: 'Dedup', value: sample.op.DedupKey },
      ]} />

      {inputItems.length > 0 && (
        <section className={styles.section}>
          <div className={styles.sectionTitle}>Input</div>
          <MetadataGrid items={inputItems} />
        </section>
      )}

      {children}

      {canRetry && (
        <div className="error-box" style={{ marginTop: 6 }}>
          {sample.op.RetryState.LastError}
          <Caption>Attempt {sample.op.RetryState.Attempt}/{sample.op.Retry.MaxAttempts}</Caption>
          <div className={styles.actions}>
            <Button className={styles.retryButton} onClick={onRetry} disabled={retrying}>
              {retrying ? 'Retrying…' : 'Retry Now'}
            </Button>
          </div>
        </div>
      )}
    </Panel>
  );
}
