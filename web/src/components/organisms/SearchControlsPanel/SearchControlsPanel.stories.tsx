import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { SearchControlsPanel, type RetrieverType } from './SearchControlsPanel';

const meta = {
  title: 'Component Library/Organisms/SearchControlsPanel',
  component: SearchControlsPanel,
} satisfies Meta<typeof SearchControlsPanel>;

export default meta;
type Story = StoryObj;

const sources = [
  { source_id: 'ttc-guides', source_name: 'TTC Guides' },
  { source_id: 'ttc-articles', source_name: 'TTC Articles' },
  { source_id: 'ttc-products', source_name: 'TTC Products' },
];

export const Interactive: Story = {
  render: () => {
    const [query, setQuery] = useState('fast growing trees for privacy screen');
    const [retriever, setRetriever] = useState<RetrieverType>('hybrid');
    const [indexId, setIndexId] = useState('bm25-ttc-guides-articles-fixed-1200-150');
    const [strategyId, setStrategyId] = useState('fixed-1200-150');
    const [profile, setProfile] = useState('openai-embedding-small');
    const [limit, setLimit] = useState(10);
    const [candidateLimit, setCandidateLimit] = useState(200);
    const [previewRunes, setPreviewRunes] = useState(300);
    const [selectedSourceIds, setSelectedSourceIds] = useState(['ttc-guides']);
    return (
      <div style={{ width: 280 }}>
        <SearchControlsPanel
          query={query}
          setQuery={setQuery}
          retriever={retriever}
          setRetriever={setRetriever}
          indexId={indexId}
          setIndexId={setIndexId}
          strategyId={strategyId}
          setStrategyId={setStrategyId}
          profile={profile}
          setProfile={setProfile}
          limit={limit}
          setLimit={setLimit}
          candidateLimit={candidateLimit}
          setCandidateLimit={setCandidateLimit}
          previewRunes={previewRunes}
          setPreviewRunes={setPreviewRunes}
          sources={sources}
          selectedSourceIds={selectedSourceIds}
          toggleSource={(id) => setSelectedSourceIds((prev) => prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id])}
          onSearch={() => undefined}
          isLoading={false}
          onKeyDown={() => undefined}
        />
      </div>
    );
  },
};

export const Loading: Story = {
  render: () => (
    <div style={{ width: 280 }}>
      <SearchControlsPanel
        query="hydrangea pruning"
        setQuery={() => undefined}
        retriever="vector"
        setRetriever={() => undefined}
        indexId="bm25-demo"
        setIndexId={() => undefined}
        strategyId="fixed-1200-150"
        setStrategyId={() => undefined}
        profile="openai-embedding-small"
        setProfile={() => undefined}
        limit={10}
        setLimit={() => undefined}
        candidateLimit={200}
        setCandidateLimit={() => undefined}
        previewRunes={300}
        setPreviewRunes={() => undefined}
        sources={[]}
        selectedSourceIds={[]}
        toggleSource={() => undefined}
        onSearch={() => undefined}
        isLoading
        onKeyDown={() => undefined}
      />
    </div>
  ),
};
