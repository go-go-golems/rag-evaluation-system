import type { Meta, StoryObj } from '@storybook/react-vite';
import { IconButton } from '../../atoms';
import { MetadataGrid } from '../../molecules';
import { Panel } from './Panel';

const meta = { title: 'Design System/Layout/Panel', component: Panel } satisfies Meta<typeof Panel>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: 'Retrieval Results',
    children: 'Panel body content for a dashboard organism.',
  },
};

export const WithActions: Story = {
  render: () => (
    <Panel title="Inspector" actions={<IconButton label="Close inspector">✕</IconButton>} density="condensed">
      <MetadataGrid items={[
        { key: 'Chunk', value: 'chunk_01HX7RAGDEMO' },
        { key: 'Score', value: '0.8421' },
      ]} />
    </Panel>
  ),
};
