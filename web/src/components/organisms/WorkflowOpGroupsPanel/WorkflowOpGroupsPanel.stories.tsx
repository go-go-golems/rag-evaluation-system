import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { WorkflowOpGroupsPanel, workflowGroupKey } from './WorkflowOpGroupsPanel';
import type { WorkflowOpGroup } from '../../../services/api';

const meta = {
  title: 'Component Library/Workflows/WorkflowOpGroupsPanel',
  component: WorkflowOpGroupsPanel,
} satisfies Meta<typeof WorkflowOpGroupsPanel>;

export default meta;
type Story = StoryObj;

const groups: WorkflowOpGroup[] = [
  { operation: 'preprocess_document', queue: 'rag-eval:preprocess_document', status: 'succeeded', count: 32 },
  { operation: 'chunk_document', queue: 'rag-eval:chunk_document', status: 'running', count: 8 },
  { operation: 'compute_embeddings', queue: 'rag-eval:compute_embeddings', status: 'failed', count: 3 },
  { operation: 'build_bm25', queue: 'rag-eval:build_bm25', status: 'ready', count: 1 },
];

export const Interactive: Story = {
  render: () => {
    const [selectedKey, setSelectedKey] = useState<string | null>(workflowGroupKey(groups[1]!));
    return <WorkflowOpGroupsPanel groups={groups} totalOps={44} selectedKey={selectedKey} onInspect={(group) => setSelectedKey(workflowGroupKey(group))} />;
  },
};

export const Empty: Story = {
  render: () => <WorkflowOpGroupsPanel groups={[]} totalOps={0} selectedKey={null} onInspect={() => undefined} />,
};
