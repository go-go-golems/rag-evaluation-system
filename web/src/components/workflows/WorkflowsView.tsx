import React, { useState, useMemo, useCallback } from 'react';
import { Button } from '../atoms';
import { Caption } from '../foundation';
import { FormRow, Panel, Stack } from '../layout';
import { QueueHealthPanel, WorkflowListPanel, WorkflowOpGraphPanel, WorkflowOpGroupsPanel, WorkflowOpInspectorPanel, WorkflowOpResultPanel, WorkflowSummaryPanel, workflowGroupKey } from '../organisms';
import {
  useListWorkflowsQuery,
  useGetWorkflowQuery,
  useGetWorkflowOpsQuery,
  useGetOpResultQuery,
  useSubmitIntakeWorkflowMutation,
  useRetryOpMutation,
  useCancelWorkflowMutation,
  useListQueuesQuery,
  useListSourcesQuery,
  SubmitIntakeRequest,
} from '../../services/api';

// ─── Queue Health Widget ─────────────────────────────────────────────────────

const QueueHealthWidget: React.FC = () => {
  const { data: queues = [], isLoading } = useListQueuesQuery(undefined, {
    pollingInterval: 5000,
  });

  return <QueueHealthPanel queues={queues} isLoading={isLoading} />;
};

// ─── Submit Intake Modal ──────────────────────────────────────────────────────

interface SubmitIntakeModalProps {
  onClose: () => void;
  onSubmitted: (workflowId: string) => void;
}

const SubmitIntakeModal: React.FC<SubmitIntakeModalProps> = ({ onClose, onSubmitted }) => {
  const { data: sourcesData } = useListSourcesQuery();
  const [submitIntake, { isLoading: submitting }] = useSubmitIntakeWorkflowMutation();
  const [error, setError] = useState<string | null>(null);

  const sourceNames = (sourcesData as unknown as { sources?: Array<{ source_id: string }> })?.sources
    ?.map(s => s.source_id) ?? [];

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
    setError(null);
    try {
      const sourceIds = sourceInput.split(',').map(s => s.trim()).filter(Boolean);
      if (sourceIds.length === 0) {
        setError('At least one source ID is required');
        return;
      }
      const result = await submitIntake({ ...form, source_ids: sourceIds }).unwrap();
      onSubmitted(result.workflow_id);
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e);
      setError(msg || 'Submission failed');
    }
  }, [form, sourceInput, submitIntake, onSubmitted]);

  const set = (key: keyof SubmitIntakeRequest, val: unknown) =>
    setForm(f => ({ ...f, [key]: val }));

  return (
    <div className="modal-overlay" onClick={e => { if (e.target === e.currentTarget) onClose(); }}>
      <Panel
        className="modal-panel"
        title="Submit Intake Workflow"
        actions={<button className="copy-btn" onClick={onClose}>✕</button>}
        density="condensed"
      >
        <Stack gap="sm" style={{ padding: 10 }}>
          <fieldset className="form-section">
            <legend>Document Selection</legend>
            <Stack gap="xs">
              <FormRow
                label="Source IDs"
                control={<input className="input" style={{ width: '100%' }} value={sourceInput} onChange={e => setSourceInput(e.target.value)} placeholder="ttc-guides, ttc-articles" />}
              />
              {sourceNames.length > 0 && <Caption>Available: {sourceNames.join(', ')}</Caption>}
              <FormRow label="Strategy" control={<Caption tone="accent">{strategyId}</Caption>} />
            </Stack>
          </fieldset>

          <fieldset className="form-section">
            <legend>Chunking</legend>
            <Stack gap="xs">
              <FormRow
                label="Strategy"
                control={(
                  <select className="select" value={form.strategy} onChange={e => set('strategy', e.target.value)}>
                    <option value="fixed">fixed</option>
                    <option value="recursive">recursive</option>
                  </select>
                )}
              />
              <FormRow label="Chunk Size" control={<input className="input" type="number" value={form.chunk_size} onChange={e => set('chunk_size', +e.target.value)} style={{ width: 80 }} />} />
              <FormRow label="Overlap" control={<input className="input" type="number" value={form.overlap} onChange={e => set('overlap', +e.target.value)} style={{ width: 80 }} />} />
            </Stack>
          </fieldset>

          <fieldset className="form-section">
            <legend>Embedding</legend>
            <Stack gap="xs">
              <FormRow label="Compute" control={<input type="checkbox" checked={!form.skip_embeddings} onChange={e => set('skip_embeddings', !e.target.checked)} />} />
              {!form.skip_embeddings && (
                <>
                  <FormRow
                    label="Provider"
                    control={(
                      <select className="select" value={form.embeddings_type} onChange={e => set('embeddings_type', e.target.value)}>
                        <option value="ollama">ollama</option>
                        <option value="openai">openai</option>
                      </select>
                    )}
                  />
                  <FormRow label="Engine" control={<input className="input" value={form.embeddings_engine} onChange={e => set('embeddings_engine', e.target.value)} />} />
                  <FormRow label="Dimensions" control={<input className="input" type="number" value={form.embeddings_dimensions} onChange={e => set('embeddings_dimensions', +e.target.value)} style={{ width: 80 }} />} />
                  <FormRow label="Batch Size" control={<input className="input" type="number" value={form.batch_size} onChange={e => set('batch_size', +e.target.value)} style={{ width: 80 }} />} />
                </>
              )}
            </Stack>
          </fieldset>

          <fieldset className="form-section">
            <legend>BM25 Index</legend>
            <FormRow label="Build" control={<input type="checkbox" checked={!form.skip_bm25} onChange={e => set('skip_bm25', !e.target.checked)} />} />
          </fieldset>

          {error && <div className="error-box">{error}</div>}

          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 4 }}>
            <Button onClick={onClose}>Cancel</Button>
            <Button variant="primary" onClick={handleSubmit} disabled={submitting || !sourceInput}>
              {submitting ? 'Submitting…' : 'Submit Workflow'}
            </Button>
          </div>
        </Stack>
      </Panel>
    </div>
  );
};

// ─── Op Result Section ────────────────────────────────────────────────────────

const OpResultSection: React.FC<{ workflowId: string; opId: string; opStatus: string }> = ({
  workflowId, opId, opStatus,
}) => {
  // Only fetch result for completed (succeeded/failed) ops
  const { data: result, isLoading, error } = useGetOpResultQuery(
    { workflowId, opId },
    { skip: opStatus !== 'succeeded' && opStatus !== 'failed' },
  );

  if (opStatus !== 'succeeded' && opStatus !== 'failed') return null;

  return <WorkflowOpResultPanel result={result} isLoading={isLoading} error={error} />;
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
  const [retryOp, { isLoading: retrying }] = useRetryOpMutation();
  const [cancelWorkflow] = useCancelWorkflowMutation();
  const [inspectedGroup, setInspectedGroup] = useState<string | null>(null);

  if (summaryLoading || opsLoading) return <div className="text-dim text-mono">Loading workflow…</div>;

  const wf = summary?.workflow;
  const groups = opsData?.groups ?? [];
  const totalOps = opsData?.total ?? 0;

  if (!wf) return <div className="text-dim text-mono">Workflow not found.</div>;

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

  // Find selected group sample
  const selectedGroup = inspectedGroup ? groups.find(g => g.operation + '|' + g.status === inspectedGroup) : null;
  const inspectedSample = selectedGroup?.sample ?? null;

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
      <WorkflowSummaryPanel
        workflow={wf}
        strategyId={strategyId}
        docCount={docIds.length}
        sourceId={sourceIds[0]}
        totalOps={totalOps}
        succeededCount={succeededCount}
        runningCount={runningCount}
        pendingCount={pendingCount}
        failedCount={failedCount}
        isRunning={isRunning}
        onBack={onBack}
        onCancel={() => cancelWorkflow(workflowId)}
        onNavigateToCorpus={onNavigateToCorpus}
      />

      <WorkflowOpGraphPanel groups={groups} total={totalOps} done={succeededCount} />

      <WorkflowOpGroupsPanel
        groups={groups}
        totalOps={totalOps}
        selectedKey={inspectedGroup}
        onInspect={(group) => setInspectedGroup(group.sample ? workflowGroupKey(group) : null)}
      />

      {/* Inspector for a sample op */}
      {inspectedSample && (
        <WorkflowOpInspectorPanel
          sample={inspectedSample}
          onClose={() => setInspectedGroup(null)}
          retrying={retrying}
          onRetry={inspectedSample.status === 'failed' ? async () => {
            const group = groups.find(candidate => workflowGroupKey(candidate) === inspectedGroup);
            if (!group?.sample) return;
            try { await retryOp({ workflowId, opId: group.sample.op.ID }).unwrap(); } catch (error) { console.error('retry failed', error); }
          } : undefined}
        >
          {(inspectedSample.status === 'succeeded' || inspectedSample.status === 'failed') && (
            <OpResultSection workflowId={workflowId} opId={inspectedSample.op.ID} opStatus={inspectedSample.status} />
          )}
        </WorkflowOpInspectorPanel>
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

  return (
    <WorkflowListPanel
      workflows={result?.workflows ?? []}
      total={result?.total ?? 0}
      isLoading={isLoading}
      statusFilter={statusFilter}
      onStatusFilterChange={setStatusFilter}
      onSelect={onSelect}
    />
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
        <Button variant="primary" onClick={() => setShowSubmit(true)}>
          + New Intake Workflow
        </Button>
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
