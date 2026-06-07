import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots } from '../../../context';
import { Panel, Stack } from '../../layout';
import { ContextBudgetBar } from './ContextBudgetBar';

const [underBudget, selectedBudget, nearLimit, overBudget] = contextWindowSnapshots;

const meta = {
  title: 'Component Library/Molecules/ContextBudgetBar',
  component: ContextBudgetBar,
  args: {
    snapshot: selectedBudget!,
  },
} satisfies Meta<typeof ContextBudgetBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const BudgetStates: Story = {
  render: () => (
    <Stack gap="md">
      <Panel title="under budget"><ContextBudgetBar snapshot={underBudget!} /></Panel>
      <Panel title="near limit"><ContextBudgetBar snapshot={nearLimit!} /></Panel>
      <Panel title="over budget"><ContextBudgetBar snapshot={overBudget!} /></Panel>
    </Stack>
  ),
};

export const SelectedSegment: Story = {
  render: () => (
    <Panel title="selected context tenant">
      <ContextBudgetBar snapshot={selectedBudget!} selectedPartId="t14-file-reads" />
    </Panel>
  ),
};
