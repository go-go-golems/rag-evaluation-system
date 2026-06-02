import type { Meta, StoryObj } from '@storybook/react-vite';
import { Provider } from 'react-redux';
import { DocumentInspector } from './DocumentInspector';
import type { CorpusDocumentDetail, CorpusIdentityArgs } from '../../../services/api';
import { store } from '../../../store';

const meta = {
  title: 'Component Library/Corpus/DocumentInspector',
  component: DocumentInspector,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
  ],
} satisfies Meta<typeof DocumentInspector>;

export default meta;
type Story = StoryObj<typeof meta>;

const detail: CorpusDocumentDetail = {
  document: {
    id: 'doc_privacy_trees',
    source_id: 'ttc-guides',
    external_id: 'privacy-trees',
    title: 'Fast Growing Trees for Privacy Screens',
    url: 'https://example.test/privacy-trees',
    word_count: 1824,
    status: 'chunked',
    metadata: {
      category: 'Guides',
      audience: 'homeowners',
      region: 'US',
    },
    content_text: 'Choosing fast growing trees for privacy depends on mature width, climate, soil drainage, and maintenance expectations. Thuja Green Giant, Leyland Cypress, and selected arborvitae cultivars are common options for evergreen screening.',
    content_type: 'text/html',
    author: 'Tree Center Editorial',
    language: 'en',
    created_at: '2026-05-20T12:00:00Z',
    updated_at: '2026-05-25T12:00:00Z',
  },
  chunks: [
    {
      id: 'chunk_privacy_001',
      strategy_id: 'fixed-1200-150',
      chunk_index: 1,
      start_offset: 0,
      end_offset: 1100,
      token_count: 220,
      text: 'Choosing fast growing trees for privacy depends on mature width, climate, soil drainage, and maintenance expectations.',
      embedding: { present: true, provider_type: 'openai', model: 'text-embedding-3-small', dimensions: 1536 },
      enrichment: { present: true, prompt_version: 'summary-v1', short_summary: 'Privacy tree selection factors.', quality_score: 0.9 },
    },
    {
      id: 'chunk_privacy_002',
      strategy_id: 'fixed-1200-150',
      chunk_index: 2,
      start_offset: 950,
      end_offset: 2100,
      token_count: 245,
      text: 'Thuja Green Giant, Leyland Cypress, and selected arborvitae cultivars are common options for evergreen screening.',
      embedding: { present: false, provider_type: 'openai', model: 'text-embedding-3-small', dimensions: 1536 },
      enrichment: { present: false },
    },
  ],
};

const identity: CorpusIdentityArgs = {
  strategy_id: 'fixed-1200-150',
  provider_type: 'openai',
  model: 'text-embedding-3-small',
  dimensions: 1536,
};

export const Overview: Story = {
  args: {
    detail,
    chunks: detail.chunks,
    identity,
  },
};

export const HighlightedChunk: Story = {
  args: {
    detail,
    chunks: detail.chunks,
    identity,
    highlightChunkId: 'chunk_privacy_002',
  },
};

export const NoContentText: Story = {
  args: {
    detail: {
      ...detail,
      document: {
        ...detail.document,
        content_text: '',
        metadata: {},
      },
    },
    chunks: [],
    identity,
  },
};
