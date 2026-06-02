import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { WorkflowListPanel } from './WorkflowListPanel';
import type { WorkflowListItem } from '../../../services/api';

const meta = {
  title: 'Component Library/Workflows/WorkflowListPanel',
  component: WorkflowListPanel,
} satisfies Meta<typeof WorkflowListPanel>;

export default meta;
type Story = StoryObj;

const workflows: WorkflowListItem[] = [
  { workflow: { ID: 'wf_intake_001', Site: 'local', Name: 'TTC guide intake', Status: 'running', Input: {}, Metadata: { strategy: 'fixed-1200-150' }, CreatedAt: new Date(Date.now() - 20 * 60_000).toISOString(), UpdatedAt: new Date().toISOString() }, opDone: 42, opTotal: 120 },
  { workflow: { ID: 'wf_intake_002', Site: 'local', Name: 'Article embeddings', Status: 'succeeded', Input: {}, Metadata: { strategy: 'fixed-1200-150' }, CreatedAt: new Date(Date.now() - 2 * 60 * 60_000).toISOString(), UpdatedAt: new Date().toISOString() }, opDone: 84, opTotal: 84 },
  { workflow: { ID: 'wf_intake_003', Site: 'local', Name: 'Product preprocessing', Status: 'failed', Input: {}, Metadata: { strategy: 'fixed-1200-150' }, CreatedAt: new Date(Date.now() - 6 * 60 * 60_000).toISOString(), UpdatedAt: new Date().toISOString() }, opDone: 12, opTotal: 60 },
];

export const Interactive: Story = {
  render: () => {
    const [statusFilter, setStatusFilter] = useState('');
    return <WorkflowListPanel workflows={workflows} total={workflows.length} statusFilter={statusFilter} onStatusFilterChange={setStatusFilter} onSelect={() => undefined} />;
  },
};

export const Empty: Story = {
  render: () => <WorkflowListPanel workflows={[]} total={0} statusFilter="failed" onStatusFilterChange={() => undefined} onSelect={() => undefined} />,
};

export const Loading: Story = {
  render: () => <WorkflowListPanel workflows={[]} total={0} statusFilter="" isLoading onStatusFilterChange={() => undefined} onSelect={() => undefined} />,
};
