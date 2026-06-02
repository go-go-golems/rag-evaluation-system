import type { Meta, StoryObj } from '@storybook/react-vite';
import { RetrievalResultsPanel } from './RetrievalResultsPanel';
import type { RetrievalResult, SearchResult } from '../../../services/api';

const meta = {
  title: 'Component Library/Organisms/RetrievalResultsPanel',
  component: RetrievalResultsPanel,
} satisfies Meta<typeof RetrievalResultsPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const hybridItems: RetrievalResult[] = [
  {
    rank: 1,
    chunk_id: 'chunk_privacy_001',
    document_id: 'doc_privacy_trees',
    source_id: 'ttc-guides',
    title: 'Fast Growing Trees for Privacy Screens',
    url: 'https://example.test/privacy-trees',
    strategy_id: 'fixed-1200-150',
    chunk_index: 4,
    score: 0.0387,
    retriever: 'hybrid',
    preview: 'Thuja Green Giant and Leyland Cypress are commonly used where a fast evergreen screen is needed...',
    components: {
      bm25: { rank: 1, score: 12.44 },
      vector: { rank: 3, score: 0.812 },
    },
  },
  {
    rank: 2,
    chunk_id: 'chunk_arborvitae_014',
    document_id: 'doc_arborvitae_spacing',
    source_id: 'ttc-articles',
    title: 'Emerald Green Arborvitae Spacing',
    strategy_id: 'fixed-1200-150',
    chunk_index: 14,
    score: 0.0321,
    retriever: 'hybrid',
    preview: 'For a dense privacy hedge, spacing depends on mature width and how quickly the screen should close...',
    components: {
      bm25: { rank: 4, score: 8.1 },
      vector: { rank: 1, score: 0.889 },
    },
  },
];

const hybridResult: SearchResult = {
  query: 'fast growing trees for privacy screen',
  retriever: 'hybrid',
  items: hybridItems,
};

export const Empty: Story = {
  args: {
    searchResult: null,
    items: [],
    retriever: 'bm25',
    selectedResult: null,
    searchError: null,
    isLoading: false,
    onSelectResult: () => undefined,
  },
};

export const Loading: Story = {
  args: {
    ...Empty.args,
    isLoading: true,
  },
};

export const HybridResults: Story = {
  args: {
    searchResult: hybridResult,
    items: hybridItems,
    retriever: 'hybrid',
    selectedResult: hybridItems[0]!,
    searchError: null,
    isLoading: false,
    onSelectResult: () => undefined,
  },
};

export const ErrorState: Story = {
  args: {
    ...Empty.args,
    searchError: 'vector_search_failed: provider profile openai-embedding-small is not configured',
  },
};
