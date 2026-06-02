import type { Meta, StoryObj } from '@storybook/react-vite';
import { DataTable } from './DataTable';

const meta = { title: 'Component Library/Molecules/DataTable', component: DataTable } satisfies Meta<typeof DataTable>;
export default meta;
type Story = StoryObj;

interface Row { id: string; rank: number; title: string; score: number; }
const rows: Row[] = [
  { id: 'a', rank: 1, title: 'Fast Growing Trees', score: 0.0387 },
  { id: 'b', rank: 2, title: 'Arborvitae Spacing', score: 0.0321 },
];

export const RetrievalRows: Story = {
  render: () => <DataTable rows={rows} getRowKey={(r) => r.id} selectedKey="a" columns={[{ id: 'rank', header: '#', align: 'end', cell: (r) => r.rank }, { id: 'title', header: 'Title', cell: (r) => r.title }, { id: 'score', header: 'Score', align: 'end', cell: (r) => r.score.toFixed(4) }]} />,
};
