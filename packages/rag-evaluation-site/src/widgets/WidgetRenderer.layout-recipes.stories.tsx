import type { Meta, StoryObj } from '@storybook/react-vite';
import { WidgetRenderer } from './WidgetRenderer';
import { defaultWidgetRegistry } from './defaultRegistry';
import { component, text, type WidgetNode } from './ir';

const meta = { title: 'Widget IR/Renderer/Layout Recipes', component: WidgetRenderer, args: { registry: defaultWidgetRegistry } } satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

function panel(title: string, children: WidgetNode[], actions?: WidgetNode): WidgetNode {
  return component('Panel', { title, density: 'condensed', actions }, children);
}

const sidebar = component('Stack', { gap: 'sm' }, [
  component('Text', { size: 'label', tone: 'muted' }, [text('CONTEXT STUDIO')]),
  component('Button', { selected: true }, [text('Overview')]),
  component('Button', {}, [text('Transcript')]),
  component('Button', {}, [text('Handout')]),
  component('Divider', {}),
  component('Caption', { tone: 'muted' }, [text('Widget IR slot props rendered by SidebarShell.')]),
]);

export const SectionBlocksInAStack: Story = {
  args: {
    node: component('Stack', { gap: 'md' }, [
      component('SectionBlock', { label: 'Overview', caption: 'A plain section rendered from Widget IR.' }, [
        component('Text', { size: 'body' }, [text('Use SectionBlock to make editorial/course content without dashboard panel chrome.')]),
      ]),
      component('SectionBlock', { label: 'Implementation Notes', density: 'spacious', divider: 'top' }, [
        component('Stack', { gap: 'sm' }, [
          component('Text', { size: 'compact' }, [text('Slot props and children are both serialized as WidgetNode trees.')]),
          component('CodeText', { display: 'block' }, [text("component('SectionBlock', { label: '...' }, [...])")]),
        ]),
      ]),
      component('SectionBlock', { label: 'Review checklist', divider: 'both' }, [
        component('Inline', { gap: 'sm' }, [
          component('AnnotationBadge', { kind: 'active', label: 'renderer' }),
          component('AnnotationBadge', { kind: 'course', label: 'storybook' }),
          component('AnnotationBadge', { kind: 'result', label: 'typecheck' }),
        ]),
      ]),
    ]),
  },
};

export const SplitPaneInspector: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'leftNarrow',
      divider: true,
      left: panel('Sources', [
        component('Stack', { gap: 'sm' }, [
          component('Button', { selected: true }, [text('The Tree Center')]),
          component('Button', {}, [text('ArXiv RAG Papers')]),
          component('Button', {}, [text('Broken import batch')]),
        ]),
      ]),
      right: panel('Selected source metadata', [
        component('MetadataGrid', { density: 'compact', items: [
          { key: 'Source ID', value: component('CodeText', {}, [text('src_tree_center')]), copyValue: 'src_tree_center' },
          { key: 'Documents', value: '1,284' },
          { key: 'Chunks', value: '38,492' },
          { key: 'Status', value: component('StatusText', { status: 'done', icon: true }, [text('done')]) },
        ] }),
      ]),
    }),
  },
};

export const SidebarShellWithHeaderAndFooter: Story = {
  args: {
    node: component('SidebarShell', {
      sidebarWidth: 204,
      contentPadding: 'lg',
      sidebar,
      header: component('Inline', { justify: 'between' }, [
        component('Text', { size: 'metric', weight: 'bold' }, [text('IR Studio')]),
        component('Caption', { tone: 'muted' }, [text('schema 0.1.0')]),
      ]),
      footer: component('Caption', { tone: 'muted' }, [text('No backend calls in this story.')]),
    }, [
      component('Stack', { gap: 'md' }, [
        component('SectionBlock', { label: 'Main region', caption: 'SidebarShell renders header/sidebar/footer slots from WidgetNode props.' }, [
          component('Text', {}, [text('This story proves that node-valued props can be recursively rendered by WidgetRenderer.')]),
        ]),
        panel('Embedded panel', [component('Caption', { tone: 'accent' }, [text('Existing panel and new layout primitives compose cleanly.')])]),
      ]),
    ]),
  },
};
