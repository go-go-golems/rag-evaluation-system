import type { Meta, StoryObj } from '@storybook/react-vite';
import type { ContextDiagramStyle } from '../../../context';
import { Caption } from '../../foundation';
import { Panel, Stack } from '../../layout';
import { ContextLegend } from './ContextLegend';

const meta = {
  title: 'Component Library/Molecules/ContextLegend',
  component: ContextLegend,
  args: {
    mode: 'pattern',
    compact: false,
  },
} satisfies Meta<typeof ContextLegend>;

export default meta;
type Story = StoryObj<typeof meta>;

const modes: ContextDiagramStyle[] = ['pattern', 'tone', 'outline'];

export const Modes: Story = {
  render: () => (
    <Stack gap="md">
      {modes.map((mode) => (
        <Panel key={mode} title={`${mode} mode`}>
          <ContextLegend mode={mode} selectedKind={mode === 'pattern' ? 'active' : undefined} />
        </Panel>
      ))}
    </Stack>
  ),
};

export const CompactCoreKinds: Story = {
  render: () => (
    <Stack gap="sm">
      <Caption>Core context-window tenants</Caption>
      <ContextLegend
        compact
        kinds={['system', 'context', 'summary', 'result', 'generated', 'evicted', 'active', 'empty']}
        selectedKind="result"
      />
    </Stack>
  ),
};

export const AnnotationAndCourseKinds: Story = {
  render: () => (
    <Panel title="annotation/course taxonomy">
      <ContextLegend kinds={['annotation', 'course', 'conversation', 'retrieval', 'tool', 'instruction']} mode="tone" />
    </Panel>
  ),
};
