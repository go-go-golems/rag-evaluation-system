import type { Meta, StoryObj } from '@storybook/react-vite';
import { CodeText } from './CodeText';

const meta = {
  title: 'Design System/Foundation/CodeText',
  component: CodeText,
} satisfies Meta<typeof CodeText>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Identifiers: Story = {
  render: () => (
    <div style={{ display: 'grid', gap: 8 }}>
      <CodeText>fixed-1200-150</CodeText>
      <CodeText tone="muted">doc_ttc_guides_001</CodeText>
      <CodeText tone="accent" copyable>chunk_01HX7RAGDEMO</CodeText>
    </div>
  ),
};
