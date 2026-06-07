import type { Meta, StoryObj } from '@storybook/react-vite';
import { Text } from '../../foundation';
import { SectionBlock } from './SectionBlock';
const meta = { title: 'Component Library/Layout/SectionBlock', component: SectionBlock, args: { label: "What you'll leave with", caption: 'A plain section with no black chrome.', children: <Text>Use this for editorial, course, and marketing surfaces where typography carries hierarchy.</Text> } } satisfies Meta<typeof SectionBlock>;
export default meta; type Story = StoryObj<typeof meta>;
export const Plain: Story = {};
