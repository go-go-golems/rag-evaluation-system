import type { Meta, StoryObj } from '@storybook/react-vite';
import { contextCourseFixture, contextHandoutFixture, contextSlides, contextWindowSnapshots } from '../context';
import { WidgetRenderer } from './WidgetRenderer';
import { defaultWidgetRegistry } from './defaultRegistry';
import { component, text, type WidgetNode } from './ir';

const meta = { title: 'Widget IR/Renderer/Course Studio', component: WidgetRenderer, args: { registry: defaultWidgetRegistry } } satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

const snapshot = contextWindowSnapshots[0]!;
const slide = contextSlides[0]!;
const selectedDoc = contextHandoutFixture.docs[0]!;

const irCourseStudioNavSections = [
  { id: 'present', label: 'Present', items: [{ id: 'course', label: 'Course', icon: component('ContextStudioNavIcon', { id: 'course' }) }, { id: 'slides', label: 'Slides', icon: component('ContextStudioNavIcon', { id: 'slides' }) }] },
  { id: 'analyze', label: 'Analyze', items: [{ id: 'visualize', label: 'Visualize', icon: component('ContextStudioNavIcon', { id: 'visualize' }) }, { id: 'upload', label: 'Upload', icon: component('ContextStudioNavIcon', { id: 'upload' }) }] },
  { id: 'review', label: 'Review', items: [{ id: 'transcript', label: 'Transcript', icon: component('ContextStudioNavIcon', { id: 'transcript' }) }, { id: 'comments', label: 'Comments', icon: component('ContextStudioNavIcon', { id: 'comments' }) }] },
  { id: 'take-home', label: 'Take-home', items: [{ id: 'handout', label: 'Handout', icon: component('ContextStudioNavIcon', { id: 'handout' }) }] },
];

function studio(main: WidgetNode, activeItemId = 'slides'): WidgetNode {
  return component('CourseStudioShell', {
    sections: irCourseStudioNavSections,
    activeItemId,
    title: 'Context Window Engineering',
    subtitle: 'Widget IR rendered course studio',
    sidebarFooter: component('Caption', { tone: 'muted' }, [text('IR story · no backend')]),
    onNavigateAction: { kind: 'event', event: 'widget-ir:navigate' },
  }, [main]);
}

export const CourseLessonLanding: Story = {
  args: {
    node: studio(component('CourseLessonPanel', {
      course: contextCourseFixture,
      activeAgendaItemId: contextCourseFixture.agenda[1]?.id,
      onAgendaItemSelectAction: { kind: 'event', event: 'widget-ir:agenda-select' },
    }), 'course'),
  },
};

export const CourseSlideWithContextVisual: Story = {
  args: {
    node: studio(component('CourseSlidePanel', {
      slide,
      snapshot,
      index: 0,
      total: contextSlides.length,
      visualSide: 'right',
      onPreviousAction: { kind: 'event', event: 'widget-ir:slide-prev' },
      onNextAction: { kind: 'event', event: 'widget-ir:slide-next' },
    }), 'slides'),
  },
};

export const TeachingSlideComposition: Story = {
  args: {
    node: component('SlideShell', {
      eyebrow: 'custom IR composition',
      counter: '02 / 06',
      title: 'A slide can be assembled without CourseSlidePanel',
      primarySide: 'left',
      primary: component('FigureBlock', {
        frame: 'bordered',
        caption: 'Context budget rendered inside a FigureBlock slot',
        legend: component('ContextLegend', { compact: true }),
      }, [component('ContextBudgetBar', { snapshot })]),
      secondary: component('KeyPointList', { items: [
        { id: 'one', title: 'Direct nodes', text: 'Use low-level nodes when the page needs custom structure.' },
        { id: 'two', title: 'Recipes later', text: 'Use recipes when the product composition is stable.' },
        { id: 'three', title: 'React renderer', text: 'All visuals still come from package components.' },
      ] }),
      footer: component('Inline', { justify: 'between' }, [
        component('Button', { size: 'compact' }, [text('◂ Prev')]),
        component('Button', { size: 'compact', variant: 'primary' }, [text('Next ▸')]),
      ]),
    }),
  },
};

export const HandoutDocumentShell: Story = {
  args: {
    node: component('HandoutDocumentShell', {
      intro: contextHandoutFixture.intro,
      documents: contextHandoutFixture.docs,
      selectedDocumentId: selectedDoc.id,
      title: 'Workshop handout',
      onDocumentSelectAction: { kind: 'event', event: 'widget-ir:document-select' },
      onDownloadAction: { kind: 'event', event: 'widget-ir:document-download' },
      onDownloadAllAction: { kind: 'event', event: 'widget-ir:download-all' },
    }),
  },
};

export const CoursePlusHandoutSplitView: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'balanced',
      divider: true,
      left: component('CourseSlidePanel', { slide, snapshot, index: 0, total: contextSlides.length }),
      right: component('HandoutDocumentShell', {
        intro: contextHandoutFixture.intro,
        documents: contextHandoutFixture.docs,
        selectedDocumentId: selectedDoc.id,
        title: 'Reference docs',
      }),
    }),
  },
};

export const DocumentListAndPreviewToolbar: Story = {
  args: {
    node: component('SplitPane', {
      ratio: 'sidebar',
      divider: true,
      gutter: 'lg',
      left: component('DocumentListPanel', {
        title: 'Documents',
        description: 'A custom handout layout assembled from molecules.',
        items: contextHandoutFixture.docs.map((doc) => ({ id: doc.id, title: doc.title, format: doc.format, size: doc.size, description: doc.description })),
        selectedItemId: selectedDoc.id,
        showItemDescriptions: true,
        onItemSelectAction: { kind: 'event', event: 'widget-ir:document-select' },
      }),
      right: component('Stack', { gap: 'sm' }, [
        component('DocumentPreviewToolbar', {
          file: selectedDoc.file,
          format: selectedDoc.format,
          size: selectedDoc.size,
          onDownloadAction: { kind: 'event', event: 'widget-ir:document-download' },
          rightSlot: component('AnnotationBadge', { kind: 'course', label: 'selected' }),
        }),
        component('MarkdownArticle', { source: selectedDoc.body }),
      ]),
    }),
  },
};
