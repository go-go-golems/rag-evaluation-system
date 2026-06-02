import type { Meta, StoryObj } from '@storybook/react-vite';
import { SelectInput } from './SelectInput';

const meta = { title: 'Design System/Atoms/SelectInput', component: SelectInput } satisfies Meta<typeof SelectInput>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
  render: () => (
    <SelectInput defaultValue="hybrid">
      <option value="bm25">bm25</option>
      <option value="vector">vector</option>
      <option value="hybrid">hybrid</option>
    </SelectInput>
  ),
};
