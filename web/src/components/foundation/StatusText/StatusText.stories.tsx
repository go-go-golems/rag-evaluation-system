import type { Meta, StoryObj } from '@storybook/react-vite';
import { StatusText, type RagStatus } from './StatusText';

const meta = {
  title: 'Design System/Foundation/StatusText',
  component: StatusText,
  args: {
    status: 'running',
    children: 'running',
  },
} satisfies Meta<typeof StatusText>;

export default meta;
type Story = StoryObj<typeof meta>;

const statuses: RagStatus[] = ['pending', 'ready', 'running', 'succeeded', 'partial', 'warning', 'failed', 'canceled'];

export const WorkflowStatuses: Story = {
  render: () => (
    <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap' }}>
      {statuses.map((status) => <StatusText key={status} status={status} icon />)}
    </div>
  ),
};
