import type { Meta, StoryObj } from '@storybook/react-vite';
import { WorkflowSummaryPanel } from './WorkflowSummaryPanel';
import type { WorkflowListItem } from '../../../services/api';

const meta = {
  title: 'Component Library/Workflows/WorkflowSummaryPanel',
  component: WorkflowSummaryPanel,
} satisfies Meta<typeof WorkflowSummaryPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const workflow: WorkflowListItem['workflow'] = {
  ID: 'wf_intake_ttc_guides_001',
  Site: 'local',
  Name: 'TTC guide intake',
  Status: 'running',
  Input: {},
  Metadata: { strategy: 'fixed-1200-150' },
  CreatedAt: new Date(Date.now() - 24 * 60_000).toISOString(),
  UpdatedAt: new Date().toISOString(),
};

export const Running: Story = {
  args: {
    workflow,
    strategyId: 'fixed-1200-150',
    docCount: 128,
    sourceId: 'ttc-guides',
    totalOps: 120,
    succeededCount: 42,
    runningCount: 4,
    pendingCount: 70,
    failedCount: 4,
    isRunning: true,
    onBack: () => undefined,
    onCancel: () => undefined,
    onNavigateToCorpus: () => undefined,
  },
};

export const Succeeded: Story = {
  args: {
    ...Running.args,
    workflow: { ...workflow, Status: 'succeeded' },
    succeededCount: 120,
    runningCount: 0,
    pendingCount: 0,
    failedCount: 0,
    isRunning: false,
  },
};
