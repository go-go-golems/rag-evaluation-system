import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots } from '../../../context';
import { Panel, Stack } from '../../layout';
import { ContextLegend } from '../ContextLegend';
import { ContextStripDiagram } from './ContextStripDiagram';

const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;

const meta = {
  title: 'Component Library/Molecules/ContextStripDiagram',
  component: ContextStripDiagram,
  args: { snapshot: deepBug! },
} satisfies Meta<typeof ContextStripDiagram>;

export default meta;
type Story = StoryObj<typeof meta>;

export const DenseSegments: Story = {
  render: () => (
    <Panel title="turn 14 context strip">
      <Stack gap="sm">
        <ContextStripDiagram snapshot={deepBug!} />
        <ContextLegend compact kinds={['system', 'context', 'summary', 'result', 'generated', 'active', 'empty']} />
      </Stack>
    </Panel>
  ),
};

export const SelectedSegment: Story = {
  render: () => (
    <Panel title="selected file reads">
      <ContextStripDiagram snapshot={deepBug!} selectedPartId="t14-file-reads" />
    </Panel>
  ),
};

export const LimitComparison: Story = {
  render: () => (
    <Stack gap="md">
      <Panel title="at limit"><ContextStripDiagram snapshot={atLimit!} /></Panel>
      <Panel title="before reclaim"><ContextStripDiagram snapshot={overBudget!} /></Panel>
    </Stack>
  ),
};
