import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { DocumentBrowser } from './DocumentBrowser';
import type { CorpusDocumentRow } from '../../../services/api';

const meta = {
  title: 'Component Library/Corpus/DocumentBrowser',
  component: DocumentBrowser,
} satisfies Meta<typeof DocumentBrowser>;

export default meta;
type Story = StoryObj;

const documents: CorpusDocumentRow[] = [
  { id: 'doc_privacy_trees', source_id: 'ttc-guides', title: 'Fast Growing Trees for Privacy Screens', url: 'https://example.test/privacy', word_count: 1824, status: 'chunked', chunk_count: 12, embedded_count: 12, missing_embedding_count: 0 },
  { id: 'doc_hydrangea_pruning', source_id: 'ttc-guides', title: 'Hydrangea Pruning Guide', url: 'https://example.test/hydrangea', word_count: 1320, status: 'chunked', chunk_count: 9, embedded_count: 5, missing_embedding_count: 4 },
  { id: 'doc_arborvitae_spacing', source_id: 'ttc-articles', title: 'Emerald Green Arborvitae Spacing', url: 'https://example.test/arborvitae', word_count: 980, status: 'pending', chunk_count: 0, embedded_count: 0, missing_embedding_count: 0 },
];

export const Interactive: Story = {
  render: () => {
    const [selectedDocId, setSelectedDocId] = useState('doc_privacy_trees');
    return (
      <div style={{ height: 360, display: 'flex' }}>
        <DocumentBrowser
          documents={documents}
          isLoading={false}
          sourceName="TTC Guides"
          totalDocs={42}
          selectedDocId={selectedDocId}
          onSelectDocument={setSelectedDocId}
          onLoadMore={() => undefined}
          hasMore
        />
      </div>
    );
  },
};

export const NoSourceSelected: Story = {
  render: () => <DocumentBrowser documents={[]} isLoading={false} sourceName={null} totalDocs={0} selectedDocId="" onSelectDocument={() => undefined} onLoadMore={() => undefined} hasMore={false} />,
};

export const Loading: Story = {
  render: () => <DocumentBrowser documents={[]} isLoading sourceName="TTC Guides" totalDocs={0} selectedDocId="" onSelectDocument={() => undefined} onLoadMore={() => undefined} hasMore={false} />,
};
