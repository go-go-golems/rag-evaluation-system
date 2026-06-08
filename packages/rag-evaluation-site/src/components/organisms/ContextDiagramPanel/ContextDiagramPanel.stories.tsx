import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots, type ContextWindowSnapshot } from '../../../context';
import { Stack } from '../../layout';
import { ContextDiagramPanel } from './ContextDiagramPanel';

const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;

const contentBlocks: ContextWindowSnapshot = {
  id: 'content-blocks',
  title: 'Turn 8 — actual context blocks',
  subtitle: 'Each segment is a transcript or tool block, not an aggregate bucket.',
  limit: 32_000,
  selectedPartId: 'turn-7-tool-result-bash',
  parts: [
    {
      id: 'system',
      label: 'system prompt',
      kind: 'system',
      tokens: 950,
      note: 'Runtime instructions and tool-use policy.',
      contentPreview: 'You are an expert coding assistant. Prefer reading files before editing. Keep responses concise.',
      metadata: { source: 'session', blockType: 'system' },
    },
    {
      id: 'turn-5-user',
      label: 'T5 user',
      kind: 'conversation',
      tokens: 210,
      note: 'The user asks for a real context visualization instead of aggregate metrics.',
      contentPreview: 'I want the context view visualization to be more about the actual content, seeing the different turns...',
      metadata: { turn: 5, role: 'user' },
    },
    {
      id: 'turn-6-assistant-plan',
      label: 'T6 assistant plan',
      kind: 'generated',
      tokens: 480,
      note: 'Assistant planning text that may or may not be preserved by the upstream agent runtime.',
      contentPreview: 'I’ll inspect the current ContextDiagram IR/component contract and our minitrace context model...',
      metadata: { turn: 6, role: 'assistant', blockType: 'plan' },
    },
    {
      id: 'turn-7-tool-call-rg',
      label: 'T7 tool call rg',
      kind: 'tool',
      tokens: 120,
      note: 'The command invocation itself.',
      contentPreview: '$ rg -n "ContextDiagramPanel|ContextWindowPart" packages pkg -S',
      metadata: { turn: 7, role: 'tool', toolName: 'bash', operation: 'execute' },
    },
    {
      id: 'turn-7-tool-result-bash',
      label: 'T7 rg output',
      kind: 'result',
      tokens: 1320,
      note: 'Search output that explains which IR files and components need to change.',
      contentPreview: 'packages/rag-evaluation-site/src/widgets/ir.ts:229: export interface ContextDiagramPanelWidgetProps...\npackages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/ContextDiagramPanel.tsx...',
      metadata: { turn: 7, role: 'tool', toolName: 'bash', fullBytes: 5280 },
    },
    {
      id: 'turn-8-current',
      label: 'T8 current answer',
      kind: 'active',
      tokens: 390,
      note: 'The currently generated response.',
      contentPreview: 'Yes — your read is right. Current ContextDiagramPanel hard-codes the legend...',
      metadata: { turn: 8, role: 'assistant' },
    },
    { id: 'free', label: 'free space', kind: 'empty', tokens: 28_530, note: 'Remaining model context budget.' },
  ],
};

const meta = {
  title: 'Component Library/Organisms/ContextDiagramPanel',
  component: ContextDiagramPanel,
  args: { snapshot: deepBug! },
} satisfies Meta<typeof ContextDiagramPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const InteractiveViews: Story = {
  render: () => <ContextDiagramPanel snapshot={deepBug!} />,
};

export const StartingViews: Story = {
  render: () => (
    <Stack gap="md">
      <ContextDiagramPanel snapshot={deepBug!} initialView="strip" />
      <ContextDiagramPanel snapshot={atLimit!} initialView="treemap" />
      <ContextDiagramPanel snapshot={overBudget!} initialView="budget" />
    </Stack>
  ),
};

export const ContentBlocksWithPartDetails: Story = {
  render: () => (
    <ContextDiagramPanel
      snapshot={contentBlocks}
      initialView="stack"
      views={['stack', 'strip', 'budget']}
      showPartDetails
    />
  ),
};

export const LegendDerivedFromSnapshotParts: Story = {
  render: () => (
    <ContextDiagramPanel
      snapshot={contentBlocks}
      initialView="strip"
      legendKinds={['conversation', 'tool', 'result', 'active', 'empty']}
      showPartDetails
    />
  ),
};
