import type { Meta, StoryObj } from '@storybook/react-vite';
import { Inline } from '../../layout';
import { Button } from './Button';

const meta = {
  title: 'Design System/Atoms/Button',
  component: Button,
  args: {
    children: 'Button',
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Variants: Story = {
  render: () => (
    <Inline>
      <Button>Default</Button>
      <Button variant="primary">Primary</Button>
      <Button selected aria-pressed>Selected</Button>
      <Button disabled>Disabled</Button>
    </Inline>
  ),
};

export const Compact: Story = {
  render: () => (
    <Inline gap="xs">
      <Button size="compact">summary/v1 [128]</Button>
      <Button size="compact" selected aria-pressed>selected/v2 [42]</Button>
    </Inline>
  ),
};
