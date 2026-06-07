import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots } from '../../../context';
import { Stack } from '../../layout';
import { ContextDiagramPanel } from './ContextDiagramPanel';
const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;
const meta = { title: 'Component Library/Organisms/ContextDiagramPanel', component: ContextDiagramPanel, args: { snapshot: deepBug! } } satisfies Meta<typeof ContextDiagramPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const InteractiveViews: Story = { render: () => <ContextDiagramPanel snapshot={deepBug!} /> };
export const StartingViews: Story = { render: () => <Stack gap="md"><ContextDiagramPanel snapshot={deepBug!} initialView="strip" /><ContextDiagramPanel snapshot={atLimit!} initialView="treemap" /><ContextDiagramPanel snapshot={overBudget!} initialView="budget" /></Stack> };
