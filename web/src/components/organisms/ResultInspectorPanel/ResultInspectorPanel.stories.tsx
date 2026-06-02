import type { Meta, StoryObj } from '@storybook/react-vite';
import { ResultInspectorPanel } from './ResultInspectorPanel';
import type { CorpusDocumentDetail, RetrievalResult } from '../../../services/api';

const meta = {
  title: 'Component Library/Organisms/ResultInspectorPanel',
  component: ResultInspectorPanel,
} satisfies Meta<typeof ResultInspectorPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const selectedResult: RetrievalResult = {
  rank: 1,
  chunk_id: 'chunk_privacy_004',
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
};

const documentDetail: CorpusDocumentDetail = {
  document: {
    id: 'doc_privacy_trees',
    source_id: 'ttc-guides',
    external_id: 'privacy-trees',
    title: 'Fast Growing Trees for Privacy Screens',
    url: 'https://example.test/privacy-trees',
    word_count: 1824,
    status: 'processed',
    metadata: {
      category: 'Guides',
      author: 'Tree Center Editorial',
      audience: 'homeowners',
    },
    content_type: 'text/html',
    author: 'Tree Center Editorial',
    language: 'en',
    created_at: '2026-05-20T12:00:00Z',
    updated_at: '2026-05-25T12:00:00Z',
  },
  chunks: [
    {
      id: 'chunk_privacy_003',
      strategy_id: 'fixed-1200-150',
      chunk_index: 3,
      start_offset: 1800,
      end_offset: 2900,
      token_count: 248,
      text: 'Before choosing a privacy tree, consider the mature width, soil drainage, and available light along the property line.',
      embedding: { present: true, provider_type: 'openai', model: 'text-embedding-3-small', dimensions: 1536 },
    },
    {
      id: 'chunk_privacy_004',
      strategy_id: 'fixed-1200-150',
      chunk_index: 4,
      start_offset: 2901,
      end_offset: 4100,
      token_count: 276,
      text: 'Thuja Green Giant and Leyland Cypress are commonly used where a fast evergreen screen is needed. Spacing should account for mature canopy width and maintenance access.',
      embedding: { present: true, provider_type: 'openai', model: 'text-embedding-3-small', dimensions: 1536 },
    },
    {
      id: 'chunk_privacy_005',
      strategy_id: 'fixed-1200-150',
      chunk_index: 5,
      start_offset: 4101,
      end_offset: 5290,
      token_count: 261,
      text: 'For smaller yards, slower but narrower arborvitae cultivars may provide a better long-term hedge than very fast growers.',
      embedding: { present: false, provider_type: 'openai', model: 'text-embedding-3-small', dimensions: 1536 },
    },
  ],
};

export const Detail: Story = {
  args: {
    result: selectedResult,
    docDetail: documentDetail,
    onClose: () => undefined,
    onOpenInCorpus: () => undefined,
  },
  render: (args) => (
    <div style={{ width: 420, height: 620 }}>
      <ResultInspectorPanel {...args} />
    </div>
  ),
};

export const ChunkTab: Story = {
  args: {
    ...Detail.args,
    defaultTab: 'chunk',
  },
  render: Detail.render,
};

export const DocumentTab: Story = {
  args: {
    ...Detail.args,
    defaultTab: 'document',
  },
  render: Detail.render,
};

export const DocumentLoading: Story = {
  args: {
    ...Detail.args,
    docDetail: null,
    defaultTab: 'document',
  },
  render: Detail.render,
};

export const BM25OnlyResult: Story = {
  args: {
    ...Detail.args,
    result: {
      ...selectedResult,
      retriever: 'bm25',
      score: 13.25,
      components: undefined,
    },
  },
  render: Detail.render,
};
