import type { Meta, StoryObj } from '@storybook/react-vite';
import { QueueHealthPanel } from './QueueHealthPanel';
import type { QueueStatus } from '../../../services/api';

const meta = {
  title: 'Component Library/Workflows/QueueHealthPanel',
  component: QueueHealthPanel,
} satisfies Meta<typeof QueueHealthPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const queues: QueueStatus[] = [
  { site: 'local', queue: 'rag-eval:chunk_document', pending: 0, ready: 8, running: 2, succeeded: 120, failed: 0, inFlight: 2, maxInFlight: 4 },
  { site: 'local', queue: 'rag-eval:compute_embeddings', pending: 0, ready: 18, running: 4, succeeded: 86, failed: 3, inFlight: 4, maxInFlight: 4 },
];

export const Active: Story = { args: { queues } };
export const Empty: Story = { args: { queues: [] } };
export const Loading: Story = { args: { queues: [], isLoading: true } };
