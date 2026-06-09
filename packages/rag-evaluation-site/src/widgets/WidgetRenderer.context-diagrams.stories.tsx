import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextWindowSnapshots } from '../context';
import { WidgetRenderer } from './WidgetRenderer';
import { defaultWidgetRegistry } from './defaultRegistry';
import { component, text, type WidgetNode } from './ir';

const meta = { title: 'Widget IR/Renderer/Context Diagrams', component: WidgetRenderer, args: { registry: defaultWidgetRegistry } } satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

const snapshot = contextWindowSnapshots[0]!;
const contentSnapshot = {
  id: 'widget-ir-content-blocks',
  title: 'Widget IR content blocks',
  subtitle: 'Server-provided parts correspond to transcript/tool blocks.',
  limit: 16_000,
  selectedPartId: 'tool-result',
  parts: [
    { id: 'system', label: 'system', kind: 'system', tokens: 720, note: 'Instructions and tool policy.', contentPreview: 'You are an expert coding assistant...' },
    { id: 'user-turn', label: 'T4 user', kind: 'conversation', tokens: 180, note: 'User asks for content-level context visualization.', contentPreview: 'I want the context view visualization to be more about the actual content...' },
    { id: 'tool-call', label: 'T5 bash call', kind: 'tool', tokens: 90, note: 'Search command.', contentPreview: '$ rg -n ContextDiagramPanel packages pkg -S', metadata: { turn: 5, toolName: 'bash' } },
    { id: 'tool-result', label: 'T5 search output', kind: 'result', tokens: 1180, note: 'Search results showing the files to change.', contentPreview: 'packages/rag-evaluation-site/src/widgets/ir.ts:229...\npackages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel...', metadata: { turn: 5, fullBytes: 4720 } },
    { id: 'answer', label: 'T6 assistant', kind: 'active', tokens: 340, note: 'Current answer draft.', contentPreview: 'Yes — the current panel hard-codes a generic legend...' },
    { id: 'free', label: 'free', kind: 'empty', tokens: 13_490 },
  ],
};
const overBudgetSnapshot = {
  ...contextWindowSnapshots[1]!,
  id: 'over-budget-widget-ir',
  title: 'Over budget context window',
  limit: 12_000,
};

function panel(title: string, children: WidgetNode[]): WidgetNode {
  return component('Panel', { title, density: 'condensed' }, children);
}

export const ContextDiagramGallery: Story = {
  args: {
    node: component('Stack', { gap: 'md' }, [
      panel('Same snapshot, four diagram renderers', [
        component('DashboardGrid', { recipe: 'twoColumn' }, [
          component('ContextBudgetBar', { snapshot, selectedPartId: snapshot.selectedPartId, showLegend: true }),
          component('ContextStripDiagram', { snapshot, selectedPartId: snapshot.selectedPartId }),
          component('ContextStackDiagram', { snapshot, selectedPartId: snapshot.selectedPartId }),
          component('ContextTreemap', { snapshot, selectedPartId: snapshot.selectedPartId }),
        ]),
      ]),
      component('ContextLegend', { compact: false, selectedKind: 'conversation' }),
    ]),
  },
};

export const ContextDiagramPanelModes: Story = {
  args: {
    node: component('DashboardGrid', { recipe: 'twoColumn' }, [
      component('ContextDiagramPanel', { snapshot, initialView: 'strip', selectedPartId: snapshot.selectedPartId }),
      component('ContextDiagramPanel', { snapshot, initialView: 'budget', selectedPartId: snapshot.selectedPartId }),
      component('ContextDiagramPanel', { snapshot, initialView: 'stack', selectedPartId: snapshot.selectedPartId }),
      component('ContextDiagramPanel', { snapshot, initialView: 'treemap', selectedPartId: snapshot.selectedPartId }),
    ]),
  },
};

export const ContextDiagramPanelContentDetails: Story = {
  args: {
    node: component('ContextDiagramPanel', {
      snapshot: contentSnapshot,
      initialView: 'stack',
      views: ['stack', 'strip', 'budget'],
      showPartDetails: true,
      legendKinds: ['conversation', 'tool', 'result', 'active', 'empty'],
    }),
  },
};

export const ContextDiagramWithMetadataSidebar: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'rightNarrow',
      divider: true,
      left: component('ContextDiagramPanel', { snapshot, initialView: 'treemap', selectedPartId: snapshot.selectedPartId }),
      right: component('Stack', { gap: 'md' }, [
        panel('Window metadata', [
          component('MetadataGrid', { density: 'compact', items: [
            { key: 'Window ID', value: component('CodeText', {}, [text(snapshot.id)]), copyValue: snapshot.id },
            { key: 'Limit', value: `${snapshot.limit.toLocaleString()} tokens` },
            { key: 'Parts', value: snapshot.parts.length },
            { key: 'Selected', value: snapshot.selectedPartId ?? 'none' },
          ] }),
        ]),
        panel('Legend', [component('ContextLegend', { compact: true, selectedKind: 'conversation' })]),
      ]),
    }),
  },
};

export const OverBudgetContextWindow: Story = {
  args: {
    node: component('Stack', { gap: 'md' }, [
      component('ContextBudgetBar', { snapshot: overBudgetSnapshot, showLegend: true }),
      component('Inline', { gap: 'sm', wrap: true }, [
        component('AnnotationBadge', { kind: 'evicted', label: 'eviction risk', selected: true }),
        component('Caption', { tone: 'warning' }, [text('Budget state is intentionally over the configured limit.')]),
      ]),
      component('ContextStripDiagram', { snapshot: overBudgetSnapshot, showLabels: true }),
    ]),
  },
};
