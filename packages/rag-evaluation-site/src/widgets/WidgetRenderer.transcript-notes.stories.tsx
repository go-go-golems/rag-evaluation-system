import type { Meta, StoryObj } from '@storybook/react-vite';
import { anchoredCommentFixtures, transcriptFixture } from '../context';
import { WidgetRenderer } from './WidgetRenderer';
import { component, text, type WidgetNode } from './ir';

const meta = { title: 'Widget IR/Renderer/Transcript and Notes', component: WidgetRenderer } satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

function panel(title: string, children: WidgetNode[]): WidgetNode {
  return component('Panel', { title, density: 'condensed' }, children);
}

export const TranscriptWithoutNotes: Story = {
  args: {
    node: component('TranscriptWorkspacePanel', {
      title: transcriptFixture.title,
      subtitle: transcriptFixture.subtitle,
      messages: transcriptFixture.messages,
      annotations: [],
      showNotes: false,
    }),
  },
};

export const AnnotatedTranscriptWithNotesRail: Story = {
  args: {
    node: component('TranscriptWorkspacePanel', {
      title: transcriptFixture.title,
      subtitle: transcriptFixture.subtitle,
      messages: transcriptFixture.messages,
      annotations: transcriptFixture.annotations,
      selectedAnnotationId: transcriptFixture.selectedAnnotationId,
      onAnnotationSelectAction: { kind: 'event', event: 'widget-ir:annotation-select' },
    }),
  },
};

export const TranscriptReaderPlusCustomRail: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'rightNarrow',
      divider: true,
      left: component('TranscriptReaderPanel', {
        title: 'Manual transcript composition',
        subtitle: transcriptFixture.subtitle,
        messages: transcriptFixture.messages,
        annotations: transcriptFixture.annotations,
        selectedAnnotationId: transcriptFixture.selectedAnnotationId,
        onAnnotationSelectAction: { kind: 'event', event: 'widget-ir:annotation-select' },
      }),
      right: component('AnnotationRailPanel', {
        annotations: transcriptFixture.annotations,
        selectedAnnotationId: transcriptFixture.selectedAnnotationId,
        onAnnotationSelectAction: { kind: 'event', event: 'widget-ir:annotation-select' },
      }),
    }),
  },
};

export const MessageCardStates: Story = {
  args: {
    node: component('Stack', { gap: 'sm' }, transcriptFixture.messages.slice(0, 5).map((message) =>
      component('TranscriptMessageCard', {
        message,
        annotations: transcriptFixture.annotations,
        selectedAnnotationId: transcriptFixture.selectedAnnotationId,
        onAnnotationSelectAction: { kind: 'event', event: 'widget-ir:annotation-select' },
      }),
    )),
  },
};

export const AnchoredCommentsOverTranscript: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'rightNarrow',
      divider: true,
      left: component('TranscriptReaderPanel', {
        title: 'Transcript with anchored comments nearby',
        messages: transcriptFixture.messages.slice(0, 5),
        annotations: transcriptFixture.annotations,
        selectedAnnotationId: transcriptFixture.selectedAnnotationId,
      }),
      right: component('AnchoredCommentRail', {
        comments: anchoredCommentFixtures,
        selectedCommentId: anchoredCommentFixtures[0]?.id,
        onCommentSelectAction: { kind: 'event', event: 'widget-ir:comment-select' },
      }),
    }),
  },
};

export const TranscriptActionLogger: Story = {
  args: { node: component('Caption', {}, [text('action logger')]) },
  render: () => (
    <WidgetRenderer
      node={component('Stack', { gap: 'md' }, [
        panel('Action logger harness', [
          component('Caption', { tone: 'muted' }, [text('Open the browser console and click note chips or comments to inspect emitted Widget IR actions.')]),
        ]),
        component('TranscriptWorkspacePanel', {
          title: transcriptFixture.title,
          messages: transcriptFixture.messages,
          annotations: transcriptFixture.annotations,
          selectedAnnotationId: transcriptFixture.selectedAnnotationId,
          onAnnotationSelectAction: { kind: 'event', event: 'widget-ir:annotation-select' },
        }),
      ])}
      onAction={(action, context) => {
        console.info('Widget IR action', { action, context });
      }}
    />
  ),
};
