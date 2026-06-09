import type { Meta, StoryObj } from '@storybook/react-vite';
import { WidgetRenderer } from './WidgetRenderer';
import { defaultWidgetRegistry } from './defaultRegistry';
import { component, text, type WidgetNode } from './ir';

const meta = { title: 'Widget IR/Renderer/Foundation and Atoms', component: WidgetRenderer, args: { registry: defaultWidgetRegistry } } satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

function panel(title: string, children: WidgetNode[]): WidgetNode {
  return component('Panel', { title, density: 'condensed' }, children);
}

export const TypographyTokenSampler: Story = {
  args: {
    node: component('Stack', { gap: 'md' }, [
      panel('Text roles rendered from Widget IR', [
        component('Stack', { gap: 'sm' }, [
          component('Text', { size: 'body' }, [text('Body prose: the renderer maps a JSON component node to the package Text primitive.')]),
          component('Text', { size: 'compact', tone: 'muted' }, [text('Compact muted metadata for dense interfaces.')]),
          component('Text', { size: 'metadata', tone: 'accent' }, [text('Metadata accent text')]),
          component('Text', { size: 'label', weight: 'bold' }, [text('UPPERCASE LABEL ROLE')]),
          component('Text', { size: 'metric', tone: 'success', weight: 'bold' }, [text('42,918 tokens')]),
        ]),
      ]),
      component('Divider', { orientation: 'horizontal' }),
      panel('Code and captions', [
        component('Stack', { gap: 'sm' }, [
          component('CodeText', { tone: 'primary' }, [text('ctx.window.parts[3].sourceId')]),
          component('CodeText', { tone: 'muted', display: 'block' }, [text('/api/widget/pages/context-studio')]),
          component('Caption', { tone: 'warning' }, [text('These nodes are rendered through WidgetRenderer, not direct JSX.')]),
        ]),
      ]),
    ]),
  },
};

export const SemanticBadgesAndSwatches: Story = {
  args: {
    node: component('Stack', { gap: 'md' }, [
      panel('Context kind swatches', [
        component('Inline', { gap: 'md', wrap: true }, ['system', 'instruction', 'conversation', 'retrieval', 'tool', 'result', 'generated', 'active', 'evicted'].map((kind, index) =>
          component('Inline', { gap: 'xs', wrap: false }, [
            component('ContextKindSwatch', { kind, selected: index === 2 }),
            component('Text', { size: 'compact' }, [text(kind)]),
          ]),
        )),
      ]),
      panel('Annotation and role badges', [
        component('Inline', { gap: 'sm', wrap: true }, [
          component('AnnotationBadge', { kind: 'context', label: 'task framing' }),
          component('AnnotationBadge', { kind: 'result', label: 'tool result', selected: true }),
          component('AnnotationBadge', { kind: 'generated', label: 'scratchpad' }),
          component('TranscriptRoleBadge', { role: 'user' }),
          component('TranscriptRoleBadge', { role: 'assistant' }),
          component('TranscriptRoleBadge', { role: 'tool', name: 'read_file' }),
        ]),
      ]),
    ]),
  },
};

export const InlineStatusHeader: Story = {
  args: {
    node: panel('Inline status header', [
      component('Inline', { gap: 'sm', justify: 'between', wrap: true }, [
        component('Inline', { gap: 'sm', wrap: true }, [
          component('StatusText', { status: 'running', icon: true }, [text('running')]),
          component('AnnotationBadge', { kind: 'active', label: 'active context', selected: true }),
          component('Caption', { tone: 'muted' }, [text('3 notes · 8 messages')]),
        ]),
        component('Inline', { gap: 'sm', wrap: true }, [
          component('Button', { size: 'compact' }, [text('Copy summary')]),
          component('Button', { size: 'compact', variant: 'primary' }, [text('Open')]),
        ]),
      ]),
    ]),
  },
};
