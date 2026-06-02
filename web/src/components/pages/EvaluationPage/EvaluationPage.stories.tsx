import type { Meta, StoryObj } from '@storybook/react-vite';
import { EvaluationPage } from './EvaluationPage';

const meta = {
  title: 'Pages/EvaluationPage',
  component: EvaluationPage,
} satisfies Meta<typeof EvaluationPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Placeholder: Story = {};

export const CustomCapabilities: Story = {
  args: {
    capabilities: [
      'Review golden query sets by source and topic',
      'Compare retrieval configs across BM25, vector, and hybrid search',
      'Inspect per-query misses and weak evidence chunks',
    ],
  },
};
