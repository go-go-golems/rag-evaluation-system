import type { Meta, StoryObj } from '@storybook/react-vite';
import { FormRow } from './FormRow';

const meta = { title: 'Design System/Layout/FormRow', component: FormRow } satisfies Meta<typeof FormRow>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = { args: { label: 'Limit', control: <input className="input" defaultValue="10" />, hint: 'Maximum number of results.' } };
