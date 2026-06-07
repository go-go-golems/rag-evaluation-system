import type { Meta, StoryObj } from '@storybook/react-vite';
import { CheckList } from './CheckList';
const meta = { title: 'Component Library/Molecules/CheckList', component: CheckList, args: { items: ['Read any context window as a budget of competing claims on space', 'Diagram a window four ways — strip, stack, budget, and treemap', 'Decide what to summarize, what to evict, and what to pin'] } } satisfies Meta<typeof CheckList>;
export default meta; type Story = StoryObj<typeof meta>;
export const Outcomes: Story = {};
