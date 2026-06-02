import type { Meta, StoryObj } from '@storybook/react-vite';
import { CheckboxRow } from './CheckboxRow';

const meta = { title: 'Design System/Atoms/CheckboxRow', component: CheckboxRow } satisfies Meta<typeof CheckboxRow>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Checked: Story = { args: { checked: true, onChange: () => {}, children: 'TTC Guides' } };
export const Unchecked: Story = { args: { checked: false, onChange: () => {}, children: 'TTC Articles' } };
