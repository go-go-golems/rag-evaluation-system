import React, { useState, useMemo, useCallback } from 'react';
import {
  useListWorkflowsQuery,
  useGetWorkflowQuery,
  useGetWorkflowOpsQuery,
  useSubmitIntakeWorkflowMutation,
  useRetryOpMutation,
  useCancelWorkflowMutation,
  useListQueuesQuery,
  useListSourcesQuery,
  WorkflowListItem,
  WorkflowOpGroup,
  QueueStatus,
  SubmitIntakeRequest,
} from '../../services/api';

// ─── Status helpers ──────────────────────────────────────────────────────────

const STATUS_ICON: Record<string, string> = {
  pending: '◌', ready: '◌', running: '●', succeeded: '✔', failed: '✘', canceled: '⊘',
};
const STATUS_CLASS: Record<string, string> = {
  pending: 'status-pending', ready: 'status-pending', running: 'status-running',
  succeeded: 'status-done', failed: 'status-error', canceled: 'status-canceled',
};

function statusIcon(s: string) { return STATUS_ICON[s] ?? '?'; }
function statusClass(s: string) { return STATUS_CLASS[s] ?? ''; }

function timeAgo(iso: string): string {
  const ms = Date.now() - new Date(iso).getTime();
  if (ms < 60000) return `${Math.floor(ms / 1000)}s`;
  if (ms < 3600000) return `${Math.floor(ms / 60000)}m`;
  if (ms < 86400000) return `${Math.floor(ms / 3600000)}h`;
  return `${Math.floor(ms / 86400000)}d`;
}

// ─── Queue Health Widget ─────────────────────────────────────────────────────

const QueueHealthWidget: React.FC = () => {
  const { data: queues = [], isLoading } = useListQueuesQuery(undefined, {
    pollingInterval: 5000,
  });

  return (
    <div className="panel">
      <div className="panel-header"><span>Queue Health</span></div>
      <div className="panel-body-condensed">
        {isLoading ? (
          <span className="text-dim text-mono">Loading...</span>
        ) : queues.length === 0 ? (
          <span className="text-dim text-mono">No queues active</span>
        ) : (
          <table className="data-table">
            <thead>
              <tr><th>Queue</th><th>Ready</th><th>Running</th><th>Failed</th></tr>
            </thead>
            <tbody>
              {queues.map((q: QueueStatus) => (
                <tr key={q.queue}>
                  <td className="mono">{q.queue.replace('rag-eval:', '')}</td>
                  <td className="num">{q.ready}</td>
                  <td className="num">{q.running}</td>
                  <td className="num" style={{ color: q.failed > 0 ? 'var(--status-error)' : undefined }}>{q.failed}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

// ─── Submit Intake Modal ──────────────────────────────────────────────────────

interface SubmitIntakeModalProps {
  onClose: () => void;
  onSubmitted: (workflowId: string) => void;
}

const SubmitIntakeModal: React.FC<SubmitIntakeModalProps> = ({ onClose, onSubmitted }) => {
  const { data: _sources = [] } = useListSourcesQuery();
  const [submitIntake, { isLoading: submitting }] = useSubmitIntakeWorkflowMutation();

  const [form, setForm] = useState<SubmitIntakeRequest>({
    strategy: 'fixed',
    chunk_size: 1200,
    overlap: 150,
    embeddings_type: 'ollama',
    embeddings_engine: 'nomic-embed-text',
    embeddings_dimensions: 768,
    batch_size: 16,
    skip_embeddings: false,
    skip_bm25: false,
  });

  const [sourceInput, setSourceInput] = useState('');

  const strategyId = useMemo(() =>
    `${form.strategy}-${form.chunk_size}-${form.overlap}`, [form.strategy, form.chunk_size, form.overlap]);

  const handleSubmit = useCallback(async () => {
    try {
      const sourceIds = sourceInput.split(',').map(s => s.trim()).filter(Boolean);
      const result = await submitIntake({ ...form, source_ids: sourceIds.length > 0 ? sourceIds : undefined }).unwrap();
      onSubmitted(result.workflow_id);
    } catch (e) {
      console.error('submit failed', e);
    }
  }, [form, sourceInput, submitIntake, onSubmitted]);

  const set = (key: keyof SubmitIntakeRequest, val: unknown) =>
    setForm(f => ({ ...f, [key]: val }));

  return (
    <div style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, background: 'rgba(0,0,0,0.5)', zIndex: 100, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="panel" style={{ width: 560, maxHeight: '90vh', overflowY: 'auto' }}>
        <div className="panel-header">
          <span>Submit Intake Workflow</span>
          <button className="copy-btn" onClick={onClose}>✕</button>
        </div>
        <div className="panel-body-condensed" style={{ display: 'flex', flexDirection: 'column', gap: 10, padding: 12 }}>
          <fieldset>
            <legend>Document Selection</legend>
            <label>Source IDs <input style={{ width: '100%' }} value={sourceInput} onChange={e => setSourceInput(e.target.value)} placeholder="ttc-guides, ttc-articles" /></label>
            <label>Strategy ID preview: <span className="mono">{strategyId}</span></label>
          </fieldset>
          <fieldset>
            <legend>Chunking</legend>
            <label>Strategy
              <select value={form.strategy} onChange={e => set('strategy', e.target.value)}>
                <option value="fixed">fixed</option>
                <option value="recursive">recursive</option>
              </select>
            </label>
            <label>Chunk Size <input type="number" value={form.chunk_size} onChange={e => set('chunk_size', +e.target.value)} /></label>
            <label>Overlap <input type="number" value={form.overlap} onChange={e => set('overlap', +e.target.value)} /></label>
          </fieldset>
          <fieldset>
            <legend>Embedding</legend>
            <label><input type="checkbox" checked={!form.skip_embeddings} onChange={e => set('skip_embeddings', !e.target.checked)} /> Compute Embeddings</label>
            {!form.skip_embeddings && (
              <>
                <label>Provider
                  <select value={form.embeddings_type} onChange={e => set('embeddings_type', e.target.value)}>
                    <option value="ollama">ollama</option>
                    <option value="openai">openai</option>
                  </select>
                </label>
                <label>Engine <input value={form.embeddings_engine} onChange={e => set('embeddings_engine', e.target.value)} /></label>
                <label>Dimensions <input type="number" value={form.embeddings_dimensions} onChange={e => set('embeddings_dimensions', +e.target.value)} /></label>
              </>
            )}
          </fieldset>
          <fieldset>
            <legend>BM25 Index</legend>
            <label><input type="checkbox" checked={!form.skip_bm25} onChange={e => set('skip_bm25', !e.target.checked)} /> Build BM25 Index</label>
          </fieldset>
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
            <button onClick={onClose}>Cancel</button>
            <button onClick={handleSubmit} disabled={submitting || !sourceInput}>
              {submitting ? 'Submitting…' : 'Submit Workflow'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

// ─── Op Group Row ────────────────────────────────────────────────────────────

const OpGroupRow: React.FC<{ group: WorkflowOpGroup; onInspect: (sample: WorkflowOpGroup['sample']) => void }> = ({ group, onInspect }) => {
  const opName = group.operation.replace(/_/g, ' ');
  return (
    <tr style={{ cursor: 'pointer' }} onClick={() => onInspect(group.sample)}>
      <td className={statusClass(group.status)}>{statusIcon(group.status)} {group.status}</td>
      <td>{opName}</td>
      <td className="mono" style={{ fontSize: 11 }}>{group.queue.replace('rag-eval:', '')}</td>
      <td className="num">{group.count}</td>
    </tr>
  );
};

// ─── Workflow Detail ──────────────────────────────────────────────────────────

interface WorkflowDetailProps {
  workflowId: string;
  onBack: () => void;
  onNavigateToCorpus: (sourceId: string, strategyId: string) => void;
}

const WorkflowDetail: React.FC<WorkflowDetailProps> = ({ workflowId, onBack, onNavigateToCorpus }) => {
  const { data: summary, isLoading: summaryLoading } = useGetWorkflowQuery(workflowId, {
    pollingInterval: 2000,
  });
  const { data: opsData, isLoading: opsLoading } = useGetWorkflowOpsQuery(workflowId, {
    pollingInterval: 2000,
  });
  const [retryOp] = useRetryOpMutation();
  const [cancelWorkflow] = useCancelWorkflowMutation();
  const [inspectedSample, setInspectedSample] = useState<WorkflowOpGroup['sample'] | null>(null);

  if (summaryLoading || opsLoading) return <span className="text-dim text-mono">Loading workflow…</span>;

  const wf = summary?.workflow;
  const groups = opsData?.groups ?? [];
  const totalOps = opsData?.total ?? 0;

  if (!wf) return <span className="text-dim text-mono">Workflow not found.</span>;

  const isRunning = wf.Status === 'running' || wf.Status === 'pending';
  const input = wf.Input as Record<string, unknown>;
  const strategyId = (input?.strategy_id as string) ?? wf.Metadata?.strategy ?? '?';
  const sourceIds = (input?.source_ids as string[]) ?? [];
  const docIds = (input?.document_ids as string[]) ?? [];

  // Compute counts per status from groups
  const succeededCount = groups.filter(g => g.status === 'succeeded').reduce((s, g) => s + g.count, 0);
  const failedCount = groups.filter(g => g.status === 'failed').reduce((s, g) => s + g.count, 0);
  const runningCount = groups.filter(g => g.status === 'running').reduce((s, g) => s + g.count, 0);
  const pendingCount = groups.filter(g => g.status === 'pending' || g.status === 'ready').reduce((s, g) => s + g.count, 0);

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      {/* Header */}
      <div className="panel">
        <div className="panel-header">
          <button className="copy-btn" onClick={onBack}>← Back</button>
          <span>{wf.ID}</span>
          <span className={statusClass(wf.Status)}>{statusIcon(wf.Status)} {wf.Status}</span>
        </div>
        <div className="panel-body-condensed" style={{ display: 'flex', gap: 16, flexWrap: 'wrap', fontSize: 12 }}>
          <span>Strategy: <span className="mono">{strategyId}</span></span>
          <span>Docs: <span className="num">{docIds.length}</span></span>
          <span>Ops: <span className="num">{succeededCount}/{totalOps}</span></span>
          {runningCount > 0 && <span>Running: <span className="num">{runningCount}</span></span>}
          {pendingCount > 0 && <span>Pending: <span className="num">{pendingCount}</span></span>}
          {failedCount > 0 && <span style={{ color: 'var(--status-error)' }}>Failed: <span className="num">{failedCount}</span></span>}
          <span>Age: <span className="mono">{timeAgo(wf.CreatedAt)}</span></span>
          {isRunning && <button onClick={() => cancelWorkflow(workflowId)}>Cancel</button>}
          {wf.Status === 'succeeded' && sourceIds[0] && (
            <button onClick={() => onNavigateToCorpus(sourceIds[0]!, strategyId)}>View in Corpus →</button>
          )}
        </div>
      </div>

      {/* Op Graph (compact) */}
      <div className="panel">
        <div className="panel-header"><span>Op Graph</span></div>
        <div className="panel-body-condensed" style={{ padding: 12, display: 'flex', flexDirection: 'column', gap: 6, alignItems: 'center' }}>
          {groups.filter(g => g.operation === 'chunk_document').map(g => (
            <span key={g.operation + g.status} className={`op-node ${statusClass(g.status)}`}
              style={{ padding: '3px 8px', borderRadius: 3, border: '1px solid var(--panel-border)', fontSize: 12 }}>
              {statusIcon(g.status)} chunk_document × {g.count} ({g.status})
            </span>
          ))}
          {groups.filter(g => g.operation === 'preprocess_document').map(g => (
            <span key={g.operation + g.status} className={`op-node ${statusClass(g.status)}`}
              style={{ padding: '3px 8px', borderRadius: 3, border: '1px solid var(--panel-border)', fontSize: 12 }}>
              {statusIcon(g.status)} preprocess × {g.count} ({g.status})
            </span>
          ))}
          <span style={{ fontSize: 10, color: 'var(--text-dim)' }}>↓ {succeededCount}/{totalOps} completed</span>
          <div style={{ display: 'flex', gap: 8 }}>
            {groups.filter(g => g.operation === 'compute_embeddings').map(g => (
              <span key={g.operation + g.status} className={`op-node ${statusClass(g.status)}`}
                style={{ padding: '4px 10px', borderRadius: 4, border: '1px solid var(--panel-border)', fontSize: 12 }}>
                {statusIcon(g.status)} embed
              </span>
            ))}
            {groups.filter(g => g.operation === 'build_bm25').map(g => (
              <span key={g.operation + g.status} className={`op-node ${statusClass(g.status)}`}
                style={{ padding: '4px 10px', borderRadius: 4, border: '1px solid var(--panel-border)', fontSize: 12 }}>
                {statusIcon(g.status)} bm25
              </span>
            ))}
          </div>
        </div>
      </div>

      {/* Op Groups Table */}
      <div className="panel">
        <div className="panel-header"><span>Ops by Group ({totalOps} total)</span></div>
        <div className="panel-body-condensed">
          <table className="data-table">
            <thead>
              <tr><th>Status</th><th>Operation</th><th>Queue</th><th>Count</th></tr>
            </thead>
            <tbody>
              {groups.map(g => (
                <OpGroupRow key={g.operation + '|' + g.status} group={g} onInspect={setInspectedSample} />
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Inspector for a sample op */}
      {inspectedSample && (
        <div className="panel" style={{ borderLeft: '3px solid var(--status-running)' }}>
          <div className="panel-header">
            <span>Sample Op: {inspectedSample.op.ID.slice(-40)}</span>
            <button className="copy-btn" onClick={() => setInspectedSample(null)}>✕</button>
          </div>
          <div className="panel-body-condensed" style={{ fontSize: 12, display: 'flex', flexDirection: 'column', gap: 6 }}>
            <div>Status: <span className={statusClass(inspectedSample.status)}>{statusIcon(inspectedSample.status)} {inspectedSample.status}</span></div>
            <div>Queue: <span className="mono">{inspectedSample.op.Queue}</span></div>
            <fieldset>
              <legend>Input</legend>
              {Object.entries(inspectedSample.op.Input as Record<string, unknown>).map(([k, v]) => (
                <div key={k}><span className="mono" style={{ color: 'var(--text-dim)' }}>{k}:</span> {typeof v === 'string' ? v : JSON.stringify(v)}</div>
              ))}
            </fieldset>
            {inspectedSample.status === 'failed' && inspectedSample.op.RetryState.LastError && (
              <fieldset style={{ borderColor: 'var(--status-error)' }}>
                <legend>Error</legend>
                <div style={{ color: 'var(--status-error)' }}>{inspectedSample.op.RetryState.LastError}</div>
                <div>Attempt: {inspectedSample.op.RetryState.Attempt}/{inspectedSample.op.Retry.MaxAttempts}</div>
                <button onClick={() => retryOp({ workflowId, opId: inspectedSample.op.ID })}>Retry Now</button>
              </fieldset>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

// ─── Workflows List ──────────────────────────────────────────────────────────

const WorkflowsList: React.FC<{ onSelect: (id: string) => void }> = ({ onSelect }) => {
  const [statusFilter, setStatusFilter] = useState('');
  const { data: result, isLoading } = useListWorkflowsQuery(
    { status: statusFilter || undefined },
    { pollingInterval: 2000 },
  );
  const workflows = result?.workflows ?? [];

  return (
    <div className="panel" style={{ flex: 1 }}>
      <div className="panel-header">
        <span>Workflows ({result?.total ?? 0})</span>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)} style={{ fontSize: 12 }}>
          <option value="">all</option>
          <option value="pending">pending</option>
          <option value="running">running</option>
          <option value="succeeded">succeeded</option>
          <option value="failed">failed</option>
          <option value="canceled">canceled</option>
        </select>
      </div>
      <div className="panel-body-condensed" style={{ overflowY: 'auto', maxHeight: 400 }}>
        {isLoading ? (
          <span className="text-dim text-mono">Loading...</span>
        ) : workflows.length === 0 ? (
          <span className="text-dim">No workflows yet. Submit one!</span>
        ) : (
          <table className="data-table">
            <thead>
              <tr><th>Status</th><th>Workflow ID</th><th>Strategy</th><th>Ops</th><th>Age</th></tr>
            </thead>
            <tbody>
              {workflows.map((w: WorkflowListItem) => (
                <tr key={w.workflow.ID} onClick={() => onSelect(w.workflow.ID)} style={{ cursor: 'pointer' }}>
                  <td className={statusClass(w.workflow.Status)}>{statusIcon(w.workflow.Status)} {w.workflow.Status}</td>
                  <td className="mono" style={{ fontSize: 11, maxWidth: 250, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{w.workflow.ID}</td>
                  <td className="mono" style={{ fontSize: 11 }}>{w.workflow.Metadata?.strategy ?? '—'}</td>
                  <td className="num">{w.opDone}/{w.opTotal}</td>
                  <td className="mono">{timeAgo(w.workflow.CreatedAt)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

// ─── Main View ───────────────────────────────────────────────────────────────

export const WorkflowsView: React.FC = () => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [showSubmit, setShowSubmit] = useState(false);

  const handleNavigateToCorpus = useCallback((sourceId: string, strategyId: string) => {
    const event = new CustomEvent('rag:navigate-to-chunk', {
      detail: { sourceId, strategyId },
    });
    window.dispatchEvent(event);
  }, []);

  if (selectedId) {
    return (
      <WorkflowDetail
        workflowId={selectedId}
        onBack={() => setSelectedId(null)}
        onNavigateToCorpus={handleNavigateToCorpus}
      />
    );
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      <div style={{ display: 'flex', gap: 8 }}>
        <WorkflowsList onSelect={setSelectedId} />
        <QueueHealthWidget />
      </div>
      <div>
        <button onClick={() => setShowSubmit(true)}>+ New Intake Workflow</button>
      </div>
      {showSubmit && (
        <SubmitIntakeModal
          onClose={() => setShowSubmit(false)}
          onSubmitted={(id) => { setShowSubmit(false); setSelectedId(id); }}
        />
      )}
    </div>
  );
};
