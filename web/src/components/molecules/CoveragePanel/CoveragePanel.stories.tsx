import type { Meta, StoryObj } from '@storybook/react-vite';
import { CoveragePanel } from './CoveragePanel';

const meta = {
  title: 'Component Library/Molecules/CoveragePanel',
  component: CoveragePanel,
} satisfies Meta<typeof CoveragePanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Complete: Story = {
  args: {
    coveragePct: 100,
    coverage: {
      strategy_id: 'fixed-1200-150',
      provider_type: 'openai',
      model: 'text-embedding-3-small',
      dimensions: 1536,
      total_chunks: 204,
      embedded_chunks: 204,
      sources: [
        { source_id: 'ttc-guides', source_name: 'TTC Guides', total_chunks: 90, embedded_chunks: 90, coverage_pct: 100, missing_chunks: 0 },
        { source_id: 'ttc-articles', source_name: 'TTC Articles', total_chunks: 114, embedded_chunks: 114, coverage_pct: 100, missing_chunks: 0 },
      ],
    },
  },
};

export const Sparse: Story = {
  args: {
    coveragePct: 32,
    coverage: {
      strategy_id: 'fixed-1200-150',
      provider_type: 'openai',
      model: 'text-embedding-3-small',
      dimensions: 1536,
      total_chunks: 204,
      embedded_chunks: 65,
      sources: [
        { source_id: 'ttc-guides', source_name: 'TTC Guides', total_chunks: 90, embedded_chunks: 40, coverage_pct: 44, missing_chunks: 50 },
        { source_id: 'ttc-articles', source_name: 'TTC Articles', total_chunks: 114, embedded_chunks: 25, coverage_pct: 22, missing_chunks: 89 },
      ],
    },
  },
};
