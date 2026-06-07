import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextKindOrder, getContextKindLabel, type ContextDiagramStyle } from '../../../context';
import { Caption } from '../../foundation';
import { Inline, Stack } from '../../layout';
import { ContextKindSwatch } from './ContextKindSwatch';

const meta = {
  title: 'Design System/Atoms/ContextKindSwatch',
  component: ContextKindSwatch,
  args: {
    kind: 'context',
    mode: 'pattern',
    size: 'md',
  },
} satisfies Meta<typeof ContextKindSwatch>;

export default meta;
type Story = StoryObj<typeof meta>;

const modes: ContextDiagramStyle[] = ['pattern', 'tone', 'outline'];

export const AllKinds: Story = {
  render: () => (
    <Stack gap="sm">
      {modes.map((mode) => (
        <Stack key={mode} gap="xs">
          <Caption transform="uppercase">{mode}</Caption>
          <Inline gap="md" wrap>
            {contextKindOrder.map((kind) => (
              <Inline key={`${mode}-${kind}`} gap="xs" style={{ alignItems: 'center' }}>
                <ContextKindSwatch kind={kind} mode={mode} />
                <Caption>{getContextKindLabel(kind)}</Caption>
              </Inline>
            ))}
          </Inline>
        </Stack>
      ))}
    </Stack>
  ),
};

export const SizesAndSelection: Story = {
  render: () => (
    <Inline gap="md" style={{ alignItems: 'center' }}>
      <ContextKindSwatch kind="result" size="sm" />
      <ContextKindSwatch kind="result" size="md" />
      <ContextKindSwatch kind="result" size="lg" />
      <ContextKindSwatch kind="active" size="lg" selected />
      <Caption>selected active context</Caption>
    </Inline>
  ),
};
