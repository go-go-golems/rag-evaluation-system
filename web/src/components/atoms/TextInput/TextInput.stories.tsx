import type { Meta, StoryObj } from '@storybook/react-vite';
import { Stack } from '../../layout';
import { TextInput } from './TextInput';

const meta = { title: 'Design System/Atoms/TextInput', component: TextInput } satisfies Meta<typeof TextInput>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Examples: Story = {
  render: () => (
    <Stack gap="xs">
      <TextInput defaultValue="fixed-1200-150" />
      <TextInput type="number" defaultValue={10} />
      <TextInput placeholder="Placeholder" />
    </Stack>
  ),
};
