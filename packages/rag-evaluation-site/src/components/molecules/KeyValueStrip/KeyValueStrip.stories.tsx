import type { Meta, StoryObj } from '@storybook/react-vite';
import { KeyValueStrip } from './KeyValueStrip';
const meta = { title: 'Component Library/Molecules/KeyValueStrip', component: KeyValueStrip, args: { items: [{ key: 'When', value: 'Thursday, June 18, 2026 · 6:30 – 8:30 PM' }, { key: 'Where', value: 'South Park Commons · San Francisco' }, { key: 'Format', value: 'In person · 24 seats · laptops required' }] } } satisfies Meta<typeof KeyValueStrip>;
export default meta; type Story = StoryObj<typeof meta>;
export const MetaStrip: Story = {};
