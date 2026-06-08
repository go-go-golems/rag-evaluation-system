import { createElement, type CSSProperties, type ReactNode } from 'react';
import { AnnotationBadge, Button, ContextKindSwatch, ErrorCallout, SelectInput, TextInput, TranscriptRoleBadge } from '../components/atoms';
import { Caption, CodeText, Divider, StatusText, Text } from '../components/foundation';
import { AppShell, DashboardGrid, FormRow, Inline, Panel, ScrollRegion, SectionBlock, SidebarShell, SlideShell, SplitPane, Stack, TabList } from '../components/layout';
import { AnnotationNoteCard, AnchoredCommentCard, AppNav, CheckList, ContextBudgetBar, ContextLegend, ContextStackDiagram, ContextStripDiagram, ContextTreemap, CourseStepNav, DataTable, DocumentListPanel, DocumentPreviewToolbar, FigureBlock, KeyPointList, KeyValueStrip, MarkdownArticle, MetadataGrid, PersonSummary, SidebarNav, StepList, TranscriptMessageCard, TranscriptSessionHeader } from '../components/molecules';
import { AnnotationRailPanel, AnchoredCommentRail, ContextDiagramPanel, CourseLessonPanel, CourseSlidePanel, CourseStudioShell, HandoutDocumentShell, TranscriptReaderPanel, TranscriptWorkspacePanel } from '../components/organisms';
import type {
  ActionSpec,
  AnchoredCommentCardWidgetProps,
  AnchoredCommentRailWidgetProps,
  AnnotationNoteCardWidgetProps,
  AnnotationRailPanelWidgetProps,
  AppNavWidgetProps,
  AppShellWidgetProps,
  AnnotationBadgeWidgetProps,
  ButtonWidgetProps,
  CaptionWidgetProps,
  CheckListWidgetProps,
  CodeTextWidgetProps,
  ComponentNode,
  ContextBudgetBarWidgetProps,
  ContextDiagramPanelWidgetProps,
  ContextKindSwatchWidgetProps,
  ContextLegendWidgetProps,
  ContextStackDiagramWidgetProps,
  ContextStripDiagramWidgetProps,
  ContextTreemapWidgetProps,
  CourseLessonPanelWidgetProps,
  CourseSlidePanelWidgetProps,
  CourseStepNavWidgetProps,
  CourseStudioShellWidgetProps,
  DashboardGridWidgetProps,
  DividerWidgetProps,
  DataTableWidgetProps,
  DocumentListPanelWidgetProps,
  DocumentPreviewToolbarWidgetProps,
  ElementNode,
  FigureBlockWidgetProps,
  FormRowWidgetProps,
  HandoutDocumentShellWidgetProps,
  InlineWidgetProps,
  JsonObject,
  KeyPointListWidgetProps,
  KeyValueStripWidgetProps,
  MarkdownArticleWidgetProps,
  MetadataGridWidgetProps,
  PanelWidgetProps,
  PersonSummaryWidgetProps,
  ScrollRegionWidgetProps,
  SectionBlockWidgetProps,
  SelectInputWidgetProps,
  SidebarNavWidgetProps,
  SidebarShellWidgetProps,
  SlideShellWidgetProps,
  SplitPaneWidgetProps,
  StepListWidgetProps,
  StackWidgetProps,
  StatusTextWidgetProps,
  TabListWidgetProps,
  TextInputWidgetProps,
  TextWidgetProps,
  TranscriptMessageCardWidgetProps,
  TranscriptReaderPanelWidgetProps,
  TranscriptRoleBadgeWidgetProps,
  TranscriptSessionHeaderWidgetProps,
  TranscriptWorkspacePanelWidgetProps,
  WidgetNode,
} from './ir';
import { bindAction, type WidgetActionContext, type WidgetActionHandler } from './actions';
import { renderCell, renderRenderable, rowKey } from './cellRenderers';

export interface WidgetRendererProps {
  node: WidgetNode;
  onAction?: WidgetActionHandler;
}

export function WidgetRenderer({ node, onAction }: WidgetRendererProps) {
  return <>{renderWidgetNode(node, onAction)}</>;
}

function renderWidgetNode(node: WidgetNode, onAction?: WidgetActionHandler): ReactNode {
  if (node.kind === 'text') return node.text;
  if (node.kind === 'element') return renderElementNode(node, onAction);
  return renderComponentNode(node, onAction);
}

function renderChildren(children: WidgetNode[] | undefined, onAction?: WidgetActionHandler): ReactNode[] {
  return (children ?? []).map((child, index) => <WidgetRenderer key={index} node={child} onAction={onAction} />);
}

function renderRenderableValue(value: unknown, onAction?: WidgetActionHandler): ReactNode {
  return renderRenderable(value as never, (node) => renderWidgetNode(node, onAction));
}

function renderNodeProp(node: WidgetNode | undefined, onAction?: WidgetActionHandler): ReactNode | undefined {
  return node ? renderWidgetNode(node, onAction) : undefined;
}

function renderElementNode(node: ElementNode, onAction?: WidgetActionHandler): ReactNode {
  const attrs = { ...(node.attrs ?? {}) } as Record<string, unknown>;
  if (attrs.style && typeof attrs.style === 'object') attrs.style = attrs.style as CSSProperties;
  return createElement(node.tag, attrs, ...renderChildren(node.children, onAction));
}

function renderComponentNode(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  switch (node.type) {
    case 'AppShell':
      return renderAppShell(node, onAction);
    case 'AppNav':
      return renderAppNav(node, onAction);
    case 'Button':
      return renderButton(node, onAction);
    case 'Text':
      return renderText(node, onAction);
    case 'CodeText':
      return renderCodeText(node, onAction);
    case 'Divider':
      return renderDivider(node);
    case 'ContextKindSwatch':
      return renderContextKindSwatch(node);
    case 'AnnotationBadge':
      return renderAnnotationBadge(node);
    case 'TranscriptRoleBadge':
      return renderTranscriptRoleBadge(node);
    case 'ContextLegend':
      return renderContextLegend(node);
    case 'ContextBudgetBar':
      return renderContextBudgetBar(node);
    case 'ContextStripDiagram':
      return renderContextStripDiagram(node);
    case 'ContextStackDiagram':
      return renderContextStackDiagram(node);
    case 'ContextTreemap':
      return renderContextTreemap(node);
    case 'ContextDiagramPanel':
      return renderContextDiagramPanel(node);
    case 'TranscriptSessionHeader':
      return renderTranscriptSessionHeader(node, onAction);
    case 'TranscriptMessageCard':
      return renderTranscriptMessageCard(node, onAction);
    case 'AnnotationNoteCard':
      return renderAnnotationNoteCard(node);
    case 'AnnotationRailPanel':
      return renderAnnotationRailPanel(node, onAction);
    case 'TranscriptReaderPanel':
      return renderTranscriptReaderPanel(node, onAction);
    case 'TranscriptWorkspacePanel':
      return renderTranscriptWorkspacePanel(node, onAction);
    case 'AnchoredCommentCard':
      return renderAnchoredCommentCard(node, onAction);
    case 'AnchoredCommentRail':
      return renderAnchoredCommentRail(node, onAction);
    case 'KeyValueStrip':
      return renderKeyValueStrip(node, onAction);
    case 'CheckList':
      return renderCheckList(node, onAction);
    case 'StepList':
      return renderStepList(node, onAction);
    case 'PersonSummary':
      return renderPersonSummary(node, onAction);
    case 'FigureBlock':
      return renderFigureBlock(node, onAction);
    case 'KeyPointList':
      return renderKeyPointList(node, onAction);
    case 'SidebarNav':
      return renderSidebarNav(node, onAction);
    case 'CourseStepNav':
      return renderCourseStepNav(node, onAction);
    case 'MarkdownArticle':
      return renderMarkdownArticle(node);
    case 'DocumentListPanel':
      return renderDocumentListPanel(node, onAction);
    case 'DocumentPreviewToolbar':
      return renderDocumentPreviewToolbar(node, onAction);
    case 'CourseLessonPanel':
      return renderCourseLessonPanel(node, onAction);
    case 'CourseSlidePanel':
      return renderCourseSlidePanel(node, onAction);
    case 'CourseStudioShell':
      return renderCourseStudioShell(node, onAction);
    case 'HandoutDocumentShell':
      return renderHandoutDocumentShell(node, onAction);
    case 'Caption':
      return renderCaption(node, onAction);
    case 'DashboardGrid':
      return renderDashboardGrid(node, onAction);
    case 'DataTable':
      return renderDataTable(node, onAction);
    case 'FormRow':
      return renderFormRow(node, onAction);
    case 'Inline':
      return renderInline(node, onAction);
    case 'MetadataGrid':
      return renderMetadataGrid(node, onAction);
    case 'Panel':
      return renderPanel(node, onAction);
    case 'ScrollRegion':
      return renderScrollRegion(node, onAction);
    case 'SectionBlock':
      return renderSectionBlock(node, onAction);
    case 'SplitPane':
      return renderSplitPane(node, onAction);
    case 'SidebarShell':
      return renderSidebarShell(node, onAction);
    case 'SlideShell':
      return renderSlideShell(node, onAction);
    case 'SelectInput':
      return renderSelectInput(node, onAction);
    case 'Stack':
      return renderStack(node, onAction);
    case 'StatusText':
      return renderStatusText(node, onAction);
    case 'TabList':
      return renderTabList(node, onAction);
    case 'TextInput':
      return renderTextInput(node);
    default:
      return <ErrorCallout>Unknown widget: {node.type}</ErrorCallout>;
  }
}

function renderAppShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AppShellWidgetProps;
  return (
    <AppShell
      className={props.className}
      header={props.header ? renderWidgetNode(props.header, onAction) : undefined}
      sidebar={props.sidebar ? renderWidgetNode(props.sidebar, onAction) : undefined}
    >
      {renderChildren(node.children, onAction)}
    </AppShell>
  );
}

function renderAppNav(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AppNavWidgetProps;
  return (
    <AppNav
      brand={renderRenderableValue(props.brand, onAction)}
      items={props.items.map((item) => ({ id: item.id, label: renderRenderableValue(item.label, onAction) }))}
      activeItemId={props.activeItemId}
      onItemSelect={(itemId) => {
        const item = props.items.find((candidate) => candidate.id === itemId);
        if (item?.action) bindAndRun(item.action, { value: itemId, componentType: 'AppNav' }, onAction);
      }}
    />
  );
}

function renderButton(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as ButtonWidgetProps;
  return (
    <Button
      className={props.className}
      variant={props.variant}
      size={props.size}
      selected={props.selected}
      disabled={props.disabled}
      type={props.type ?? 'button'}
      onClick={bindAction(props.action, { componentType: 'Button' }, onAction)}
    >
      {renderChildren(node.children, onAction)}
    </Button>
  );
}

function renderText(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TextWidgetProps;
  return <Text className={props.className} as={props.as} size={props.size} tone={props.tone} weight={props.weight} align={props.align} truncate={props.truncate}>{renderChildren(node.children, onAction)}</Text>;
}

function renderCodeText(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CodeTextWidgetProps;
  return <CodeText className={props.className} as={props.as} tone={props.tone} display={props.display} copyable={props.copyable}>{renderChildren(node.children, onAction)}</CodeText>;
}

function renderDivider(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as DividerWidgetProps;
  return <Divider className={props.className} orientation={props.orientation} />;
}

function renderContextKindSwatch(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextKindSwatchWidgetProps;
  return <ContextKindSwatch className={props.className} kind={props.kind} mode={props.mode} size={props.size} selected={props.selected} />;
}

function renderAnnotationBadge(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as AnnotationBadgeWidgetProps;
  return <AnnotationBadge className={props.className} kind={props.kind} label={props.label} selected={props.selected} />;
}

function renderTranscriptRoleBadge(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as TranscriptRoleBadgeWidgetProps;
  return <TranscriptRoleBadge className={props.className} role={props.role} name={props.name} />;
}

function renderContextLegend(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextLegendWidgetProps;
  return <ContextLegend className={props.className} kinds={props.kinds} mode={props.mode} compact={props.compact} selectedKind={props.selectedKind} />;
}

function renderContextBudgetBar(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextBudgetBarWidgetProps;
  return <ContextBudgetBar className={props.className} snapshot={props.snapshot} mode={props.mode} showLegend={props.showLegend} selectedPartId={props.selectedPartId} />;
}

function renderContextStripDiagram(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextStripDiagramWidgetProps;
  return <ContextStripDiagram className={props.className} snapshot={props.snapshot} mode={props.mode} selectedPartId={props.selectedPartId} showLabels={props.showLabels} />;
}

function renderContextStackDiagram(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextStackDiagramWidgetProps;
  return <ContextStackDiagram className={props.className} snapshot={props.snapshot} selectedPartId={props.selectedPartId} />;
}

function renderContextTreemap(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextTreemapWidgetProps;
  return <ContextTreemap className={props.className} snapshot={props.snapshot} selectedPartId={props.selectedPartId} />;
}

function renderContextDiagramPanel(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as ContextDiagramPanelWidgetProps;
  return <ContextDiagramPanel className={props.className} snapshot={props.snapshot} initialView={props.initialView} selectedPartId={props.selectedPartId} />;
}

function renderTranscriptSessionHeader(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TranscriptSessionHeaderWidgetProps;
  return <TranscriptSessionHeader className={props.className} title={renderRenderableValue(props.title, onAction)} subtitle={renderRenderableValue(props.subtitle, onAction)} messageCount={props.messageCount} annotationCount={props.annotationCount} tokenTotal={props.tokenTotal} rightSlot={renderNodeProp(props.rightSlot, onAction)} />;
}

function renderTranscriptMessageCard(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TranscriptMessageCardWidgetProps;
  return <TranscriptMessageCard className={props.className} message={props.message} annotations={props.annotations} selectedAnnotationId={props.selectedAnnotationId} showAnnotationChips={props.showAnnotationChips} onAnnotationSelect={annotationSelectHandler('TranscriptMessageCard', props.onAnnotationSelectAction, onAction)} />;
}

function renderAnnotationNoteCard(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as AnnotationNoteCardWidgetProps;
  return <AnnotationNoteCard className={props.className} annotation={props.annotation} selected={props.selected} index={props.index} />;
}

function renderAnnotationRailPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AnnotationRailPanelWidgetProps;
  return <AnnotationRailPanel className={props.className} title={props.title} description={props.description} annotations={props.annotations} selectedAnnotationId={props.selectedAnnotationId} onAnnotationSelect={annotationSelectHandler('AnnotationRailPanel', props.onAnnotationSelectAction, onAction)} />;
}

function renderTranscriptReaderPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TranscriptReaderPanelWidgetProps;
  return <TranscriptReaderPanel className={props.className} title={props.title} subtitle={props.subtitle} messages={props.messages} annotations={props.annotations} selectedAnnotationId={props.selectedAnnotationId} showAnnotationChips={props.showAnnotationChips} onAnnotationSelect={annotationSelectHandler('TranscriptReaderPanel', props.onAnnotationSelectAction, onAction)} />;
}

function renderTranscriptWorkspacePanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TranscriptWorkspacePanelWidgetProps;
  return <TranscriptWorkspacePanel className={props.className} title={props.title} subtitle={props.subtitle} messages={props.messages} annotations={props.annotations} selectedAnnotationId={props.selectedAnnotationId} showNotes={props.showNotes} onAnnotationSelect={annotationSelectHandler('TranscriptWorkspacePanel', props.onAnnotationSelectAction, onAction)} />;
}

function renderAnchoredCommentCard(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AnchoredCommentCardWidgetProps;
  return <AnchoredCommentCard className={props.className} comment={props.comment} index={props.index} selected={props.selected} compact={props.compact} onDismiss={props.onDismissAction ? () => bindAndRun(props.onDismissAction!, { commentId: props.comment.id, value: props.comment.id, componentType: 'AnchoredCommentCard' }, onAction) : undefined} />;
}

function renderAnchoredCommentRail(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as AnchoredCommentRailWidgetProps;
  return <AnchoredCommentRail className={props.className} title={props.title} comments={props.comments} selectedCommentId={props.selectedCommentId} onCommentSelect={props.onCommentSelectAction ? (commentId) => bindAndRun(props.onCommentSelectAction!, { commentId, value: commentId, componentType: 'AnchoredCommentRail' }, onAction) : undefined} />;
}

function renderKeyValueStrip(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as KeyValueStripWidgetProps;
  return <KeyValueStrip className={props.className} items={props.items.map((item) => ({ key: renderRenderableValue(item.key, onAction), value: renderRenderableValue(item.value, onAction) }))} />;
}

function renderCheckList(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CheckListWidgetProps;
  return <CheckList className={props.className} marker={renderRenderableValue(props.marker, onAction)} items={props.items.map((item) => (isChecklistObject(item) ? { ...item, text: renderRenderableValue(item.text, onAction) } : renderRenderableValue(item, onAction)))} />;
}

function renderStepList(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as StepListWidgetProps;
  return <StepList className={props.className} items={props.items.map((item) => ({ ...item, index: renderRenderableValue(item.index, onAction), title: renderRenderableValue(item.title, onAction), description: renderRenderableValue(item.description, onAction), meta: renderRenderableValue(item.meta, onAction) }))} activeItemId={props.activeItemId} onItemSelect={itemSelectHandler('StepList', props.onItemSelectAction, 'itemId', onAction)} />;
}

function renderPersonSummary(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as PersonSummaryWidgetProps;
  return <PersonSummary className={props.className} name={renderRenderableValue(props.name, onAction)} subtitle={renderRenderableValue(props.subtitle, onAction)} bio={renderRenderableValue(props.bio, onAction)} avatar={renderRenderableValue(props.avatar, onAction)} />;
}

function renderFigureBlock(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as FigureBlockWidgetProps;
  return <FigureBlock className={props.className} as={props.as} caption={renderRenderableValue(props.caption, onAction)} legend={renderNodeProp(props.legend, onAction)} frame={props.frame}>{renderChildren(node.children, onAction)}</FigureBlock>;
}

function renderKeyPointList(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as KeyPointListWidgetProps;
  return <KeyPointList className={props.className} markerTone={props.markerTone} items={props.items.map((item) => (isKeyPointObject(item) ? { ...item, index: renderRenderableValue(item.index, onAction), title: renderRenderableValue(item.title, onAction), text: renderRenderableValue(item.text, onAction), meta: renderRenderableValue(item.meta, onAction) } : renderRenderableValue(item, onAction)))} />;
}

function renderSidebarNav(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SidebarNavWidgetProps;
  return <SidebarNav className={props.className} sections={props.sections.map((section) => ({ ...section, label: renderRenderableValue(section.label, onAction), items: section.items.map((item) => ({ ...item, label: renderRenderableValue(item.label, onAction), icon: renderRenderableValue(item.icon, onAction), badge: renderRenderableValue(item.badge, onAction) })) }))} activeItemId={props.activeItemId} ariaLabel={props.ariaLabel} onItemSelect={itemSelectHandler('SidebarNav', props.onItemSelectAction, 'itemId', onAction)} />;
}

function renderCourseStepNav(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CourseStepNavWidgetProps;
  return <CourseStepNav className={props.className} items={props.items} activeItemId={props.activeItemId} onItemSelect={itemSelectHandler('CourseStepNav', props.onItemSelectAction, 'itemId', onAction)} />;
}

function renderMarkdownArticle(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as MarkdownArticleWidgetProps;
  return <MarkdownArticle className={props.className} source={props.source} />;
}

function renderDocumentListPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DocumentListPanelWidgetProps;
  return <DocumentListPanel className={props.className} title={renderRenderableValue(props.title, onAction)} description={renderRenderableValue(props.description, onAction)} items={props.items.map((item) => ({ ...item, title: renderRenderableValue(item.title, onAction), format: renderRenderableValue(item.format, onAction), size: renderRenderableValue(item.size, onAction), description: renderRenderableValue(item.description, onAction), icon: renderRenderableValue(item.icon, onAction) }))} selectedItemId={props.selectedItemId} onItemSelect={itemSelectHandler('DocumentListPanel', props.onItemSelectAction, 'itemId', onAction)} onDownloadAll={props.onDownloadAllAction ? () => bindAndRun(props.onDownloadAllAction!, { componentType: 'DocumentListPanel' }, onAction) : undefined} downloadAllLabel={renderRenderableValue(props.downloadAllLabel, onAction)} emptyMessage={renderRenderableValue(props.emptyMessage, onAction)} showItemDescriptions={props.showItemDescriptions} />;
}

function renderDocumentPreviewToolbar(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DocumentPreviewToolbarWidgetProps;
  return <DocumentPreviewToolbar className={props.className} file={renderRenderableValue(props.file, onAction)} format={renderRenderableValue(props.format, onAction)} size={renderRenderableValue(props.size, onAction)} onDownload={props.onDownloadAction ? () => bindAndRun(props.onDownloadAction!, { componentType: 'DocumentPreviewToolbar' }, onAction) : undefined} downloadLabel={renderRenderableValue(props.downloadLabel, onAction)} rightSlot={renderNodeProp(props.rightSlot, onAction)} />;
}

function renderCourseLessonPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CourseLessonPanelWidgetProps;
  return <CourseLessonPanel className={props.className} course={props.course} activeAgendaItemId={props.activeAgendaItemId} onAgendaItemSelect={itemSelectHandler('CourseLessonPanel', props.onAgendaItemSelectAction, 'agendaItemId', onAction)} />;
}

function renderCourseSlidePanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CourseSlidePanelWidgetProps;
  return <CourseSlidePanel className={props.className} slide={props.slide} snapshot={props.snapshot} index={props.index} total={props.total} visualSide={props.visualSide} onPrevious={props.onPreviousAction ? () => bindAndRun(props.onPreviousAction!, { componentType: 'CourseSlidePanel', value: 'previous' }, onAction) : undefined} onNext={props.onNextAction ? () => bindAndRun(props.onNextAction!, { componentType: 'CourseSlidePanel', value: 'next' }, onAction) : undefined} />;
}

function renderCourseStudioShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CourseStudioShellWidgetProps;
  return <CourseStudioShell className={props.className} sections={props.sections.map((section) => ({ ...section, label: renderRenderableValue(section.label, onAction), items: section.items.map((item) => ({ ...item, label: renderRenderableValue(item.label, onAction), icon: renderRenderableValue(item.icon, onAction), badge: renderRenderableValue(item.badge, onAction) })) }))} activeItemId={props.activeItemId} onNavigate={itemSelectHandler('CourseStudioShell', props.onNavigateAction, 'itemId', onAction)} title={renderRenderableValue(props.title, onAction)} subtitle={renderRenderableValue(props.subtitle, onAction)} sidebarFooter={renderNodeProp(props.sidebarFooter, onAction)}>{renderChildren(node.children, onAction)}</CourseStudioShell>;
}

function renderHandoutDocumentShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as HandoutDocumentShellWidgetProps;
  return <HandoutDocumentShell className={props.className} intro={renderRenderableValue(props.intro, onAction)} documents={props.documents} selectedDocumentId={props.selectedDocumentId} onDocumentSelect={itemSelectHandler('HandoutDocumentShell', props.onDocumentSelectAction, 'documentId', onAction)} onDownload={props.onDownloadAction ? (documentId) => bindAndRun(props.onDownloadAction!, { documentId, value: documentId, componentType: 'HandoutDocumentShell' }, onAction) : undefined} onDownloadAll={props.onDownloadAllAction ? () => bindAndRun(props.onDownloadAllAction!, { componentType: 'HandoutDocumentShell' }, onAction) : undefined} title={renderRenderableValue(props.title, onAction)} emptyMessage={renderRenderableValue(props.emptyMessage, onAction)} />;
}

function isChecklistObject(value: unknown): value is { id?: string; text: unknown } {
  return typeof value === 'object' && value !== null && 'text' in value;
}

function isKeyPointObject(value: unknown): value is { id?: string; index?: unknown; title?: unknown; text: unknown; meta?: unknown } {
  return typeof value === 'object' && value !== null && 'text' in value;
}

function itemSelectHandler(componentType: string, action: ActionSpec | undefined, contextKey: string, onAction?: WidgetActionHandler): ((itemId: string) => void) | undefined {
  if (!action) return undefined;
  return (itemId) => bindAndRun(action, { [contextKey]: itemId, value: itemId, componentType }, onAction);
}

function annotationSelectHandler(componentType: string, action: ActionSpec | undefined, onAction?: WidgetActionHandler): ((annotationId: string) => void) | undefined {
  if (!action) return undefined;
  return (annotationId) => bindAndRun(action, { annotationId, value: annotationId, componentType }, onAction);
}

function renderCaption(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as CaptionWidgetProps;
  return <Caption className={props.className} tone={props.tone} transform={props.transform} truncate={props.truncate}>{renderChildren(node.children, onAction)}</Caption>;
}

function renderDashboardGrid(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DashboardGridWidgetProps;
  return <DashboardGrid className={props.className} recipe={props.recipe}>{renderChildren(node.children, onAction)}</DashboardGrid>;
}

function renderDataTable(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as DataTableWidgetProps;
  const rowSelectAction = props.onRowSelect;
  return (
    <DataTable<JsonObject>
      className={props.className}
      rows={props.rows}
      columns={props.columns.map((column) => ({
        id: column.id,
        header: renderRenderableValue(column.header, onAction),
        align: column.align,
        maxWidth: column.maxWidth,
        cell: (row) => renderCell(column.cell, row, (child) => renderWidgetNode(child, onAction)),
      }))}
      getRowKey={(row) => rowKey(row, props.getRowKey)}
      selectedKey={props.selectedKey == null ? props.selectedKey : String(props.selectedKey)}
      emptyMessage={renderRenderableValue(props.emptyMessage, onAction)}
      onRowSelect={rowSelectAction ? (row) => bindAndRun(rowSelectAction, { row, rowKey: rowKey(row, props.getRowKey), componentType: 'DataTable' }, onAction) : undefined}
    />
  );
}

function renderFormRow(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as FormRowWidgetProps;
  return (
    <FormRow
      className={props.className}
      label={renderRenderableValue(props.label, onAction)}
      control={renderWidgetNode(props.control, onAction)}
      hint={renderRenderableValue(props.hint, onAction)}
    />
  );
}

function renderInline(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as InlineWidgetProps;
  return <Inline className={props.className} gap={props.gap} justify={props.justify} wrap={props.wrap}>{renderChildren(node.children, onAction)}</Inline>;
}

function renderMetadataGrid(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as MetadataGridWidgetProps;
  return (
    <MetadataGrid
      className={props.className}
      density={props.density}
      items={props.items.map((item) => ({
        key: renderRenderableValue(item.key, onAction),
        value: renderRenderableValue(item.value, onAction),
        copyValue: item.copyValue,
      }))}
    />
  );
}

function renderPanel(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as PanelWidgetProps;
  return (
    <Panel
      className={props.className}
      title={renderRenderableValue(props.title, onAction)}
      actions={renderRenderableValue(props.actions, onAction)}
      density={props.density}
      fill={props.fill}
    >
      {renderChildren(node.children, onAction)}
    </Panel>
  );
}

function renderScrollRegion(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as ScrollRegionWidgetProps;
  return <ScrollRegion className={props.className} axis={props.axis} style={props.style as CSSProperties | undefined}>{renderChildren(node.children, onAction)}</ScrollRegion>;
}

function renderSectionBlock(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SectionBlockWidgetProps;
  return <SectionBlock className={props.className} as={props.as} label={renderRenderableValue(props.label, onAction)} caption={renderRenderableValue(props.caption, onAction)} density={props.density} divider={props.divider}>{renderChildren(node.children, onAction)}</SectionBlock>;
}

function renderSplitPane(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SplitPaneWidgetProps;
  return <SplitPane className={props.className} left={renderWidgetNode(props.left, onAction)} right={renderWidgetNode(props.right, onAction)} ratio={props.ratio} divider={props.divider} />;
}

function renderSidebarShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SidebarShellWidgetProps;
  return <SidebarShell className={props.className} sidebar={renderNodeProp(props.sidebar, onAction)} sidebarWidth={props.sidebarWidth} header={renderNodeProp(props.header, onAction)} footer={renderNodeProp(props.footer, onAction)}>{renderChildren(node.children, onAction)}</SidebarShell>;
}

function renderSlideShell(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SlideShellWidgetProps;
  return <SlideShell className={props.className} as={props.as} eyebrow={renderRenderableValue(props.eyebrow, onAction)} counter={renderRenderableValue(props.counter, onAction)} title={renderRenderableValue(props.title, onAction)} subtitle={renderRenderableValue(props.subtitle, onAction)} primary={renderWidgetNode(props.primary, onAction)} secondary={renderNodeProp(props.secondary, onAction)} primarySide={props.primarySide} ratio={props.ratio} divider={props.divider} footer={renderNodeProp(props.footer, onAction)} />;
}

function renderSelectInput(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as SelectInputWidgetProps;
  return (
    <SelectInput className={props.className} name={props.name} value={props.value} disabled={props.disabled}>
      {(props.options ?? []).map((option) => (
        <option key={String(option.value)} value={option.value} disabled={option.disabled}>
          {renderRenderableValue(option.label, onAction)}
        </option>
      ))}
      {renderChildren(node.children, onAction)}
    </SelectInput>
  );
}

function renderStack(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as StackWidgetProps;
  return <Stack className={props.className} gap={props.gap} align={props.align}>{renderChildren(node.children, onAction)}</Stack>;
}

function renderStatusText(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as StatusTextWidgetProps;
  return <StatusText className={props.className} status={props.status} icon={props.icon}>{renderChildren(node.children, onAction)}</StatusText>;
}

function renderTabList(node: ComponentNode, onAction?: WidgetActionHandler): ReactNode {
  const props = (node.props ?? {}) as TabListWidgetProps;
  return (
    <TabList
      items={props.items.map((item) => ({ id: item.id, label: renderRenderableValue(item.label, onAction) }))}
      activeId={props.activeId}
      ariaLabel={props.ariaLabel}
      onChange={(id) => {
        if (props.onChange) bindAndRun(props.onChange, { value: id, componentType: 'TabList' }, onAction);
      }}
    />
  );
}

function renderTextInput(node: ComponentNode): ReactNode {
  const props = (node.props ?? {}) as TextInputWidgetProps;
  return (
    <TextInput
      className={props.className}
      name={props.name}
      value={props.value}
      placeholder={props.placeholder}
      type={props.type}
      disabled={props.disabled}
      min={props.min}
      max={props.max}
      readOnly
    />
  );
}

function bindAndRun(action: ActionSpec, context: WidgetActionContext, onAction?: WidgetActionHandler): void {
  bindAction(action, context, onAction)?.();
}
