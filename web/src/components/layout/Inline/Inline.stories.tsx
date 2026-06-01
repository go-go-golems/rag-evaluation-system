import type { Meta, StoryObj } from '@storybook/react-vite';
import { Button } from '../../atoms';
import { Caption } from '../../foundation';
import { Inline } from './Inline';

const meta = { title: 'Design System/Layout/Inline', component: Inline } satisfies Meta<typeof Inline>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Toolbar: Story = {
  render: () => (
    <Inline gap="sm">
      <Button variant="primary">Search</Button>
      <Button>Reset</Button>
      <Caption>fixed-1200-150</Caption>
    </Inline>
  ),
};
