import { useState } from 'react';
import type { Meta, StoryObj } from '@storybook/react-vite';
import { Stack } from '@go-go-golems/rag-evaluation-site';
import { ChunkEnrichmentIdentityBar, DocProcessingIdentityBar } from './ArtifactIdentityBar';
import type {
  ChunkEnrichmentCoverageArgs,
  ChunkEnrichmentIdentity,
  DocumentProcessingCoverageArgs,
  DocumentProcessingIdentity,
} from '../../../services/api';

const processingIdentities: DocumentProcessingIdentity[] = [
  { artifact_type: 'summary', prompt_version: 'summary-v1', provider: 'openai', model: 'gpt-4.1-mini', artifact_count: 128 },
  { artifact_type: 'keywords', prompt_version: 'keywords-v2', provider: 'openai', model: 'gpt-4.1-mini', artifact_count: 124 },
  { artifact_type: 'qa_pairs', prompt_version: 'qa-v1', provider: 'ollama', model: 'llama3.2', artifact_count: 42 },
];

const enrichmentIdentities: ChunkEnrichmentIdentity[] = [
  { strategy_id: 'fixed-1200-150', prompt_version: 'chunk-summary-v1', provider: 'openai', model: 'gpt-4.1-mini', enriched_count: 920 },
  { strategy_id: 'semantic-sections-v2', prompt_version: 'chunk-summary-v1', provider: 'openai', model: 'gpt-4.1-mini', enriched_count: 320 },
  { strategy_id: 'fixed-1200-150', prompt_version: 'entity-tags-v2', provider: 'ollama', model: 'llama3.2', enriched_count: 180 },
];

const meta = {
  title: 'Component Library/Corpus/ArtifactIdentityBars',
  parameters: { layout: 'padded' },
} satisfies Meta;

export default meta;
type Story = StoryObj<typeof meta>;

export const DocumentProcessing: Story = {
  render: () => {
    const [selected, setSelected] = useState<DocumentProcessingCoverageArgs>({
      artifact_type: 'summary',
      prompt_version: 'summary-v1',
      provider: 'openai',
      model: 'gpt-4.1-mini',
    });
    return <DocProcessingIdentityBar identities={processingIdentities} selected={selected} onChange={setSelected} />;
  },
};

export const ChunkEnrichment: Story = {
  render: () => {
    const [selected, setSelected] = useState<ChunkEnrichmentCoverageArgs>({
      strategy_id: 'fixed-1200-150',
      prompt_version: 'chunk-summary-v1',
    });
    return <ChunkEnrichmentIdentityBar identities={enrichmentIdentities} selected={selected} onChange={setSelected} />;
  },
};

export const Combined: Story = {
  render: () => {
    const [processing, setProcessing] = useState<DocumentProcessingCoverageArgs>({
      artifact_type: 'summary',
      prompt_version: 'summary-v1',
      provider: 'openai',
      model: 'gpt-4.1-mini',
    });
    const [enrichment, setEnrichment] = useState<ChunkEnrichmentCoverageArgs>({
      strategy_id: 'fixed-1200-150',
      prompt_version: 'chunk-summary-v1',
    });
    return (
      <Stack gap="xs">
        <DocProcessingIdentityBar identities={processingIdentities} selected={processing} onChange={setProcessing} />
        <ChunkEnrichmentIdentityBar identities={enrichmentIdentities} selected={enrichment} onChange={setEnrichment} />
      </Stack>
    );
  },
};

export const EmptyProcessing: Story = {
  render: () => (
    <DocProcessingIdentityBar
      identities={[]}
      selected={{ artifact_type: '', prompt_version: '', provider: '', model: '' }}
      onChange={() => {}}
    />
  ),
};

export const EmptyEnrichment: Story = {
  render: () => (
    <Stack gap="xs">
      <ChunkEnrichmentIdentityBar
        identities={[]}
        selected={{ strategy_id: '', prompt_version: '' }}
        onChange={() => {}}
      />
      <span>Empty enrichment identities intentionally render no control.</span>
    </Stack>
  ),
};
