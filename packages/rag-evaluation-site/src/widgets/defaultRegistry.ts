import { annotationBadgeWidget } from '../components/atoms/AnnotationBadge/AnnotationBadge.widget';
import { buttonWidget } from '../components/atoms/Button/Button.widget';
import { contextKindSwatchWidget } from '../components/atoms/ContextKindSwatch/ContextKindSwatch.widget';
import { selectInputWidget } from '../components/atoms/SelectInput/SelectInput.widget';
import { textInputWidget } from '../components/atoms/TextInput/TextInput.widget';
import { transcriptRoleBadgeWidget } from '../components/atoms/TranscriptRoleBadge/TranscriptRoleBadge.widget';
import { captionWidget } from '../components/foundation/Caption/Caption.widget';
import { codeTextWidget } from '../components/foundation/CodeText/CodeText.widget';
import { dividerWidget } from '../components/foundation/Divider/Divider.widget';
import { statusTextWidget } from '../components/foundation/StatusText/StatusText.widget';
import { textWidget } from '../components/foundation/Text/Text.widget';
import { appShellWidget } from '../components/layout/AppShell/AppShell.widget';
import { dashboardGridWidget } from '../components/layout/DashboardGrid/DashboardGrid.widget';
import { formRowWidget } from '../components/layout/FormRow/FormRow.widget';
import { inlineWidget } from '../components/layout/Inline/Inline.widget';
import { panelWidget } from '../components/layout/Panel/Panel.widget';
import { scrollRegionWidget } from '../components/layout/ScrollRegion/ScrollRegion.widget';
import { sectionBlockWidget } from '../components/layout/SectionBlock/SectionBlock.widget';
import { sidebarShellWidget } from '../components/layout/SidebarShell/SidebarShell.widget';
import { splitPaneWidget } from '../components/layout/SplitPane/SplitPane.widget';
import { stackWidget } from '../components/layout/Stack/Stack.widget';
import { tabListWidget } from '../components/layout/TabList/TabList.widget';
import { anchoredCommentCardWidget } from '../components/molecules/AnchoredCommentCard/AnchoredCommentCard.widget';
import { annotationNoteCardWidget } from '../components/molecules/AnnotationNoteCard/AnnotationNoteCard.widget';
import { appNavWidget } from '../components/molecules/AppNav/AppNav.widget';
import { checkListWidget } from '../components/molecules/CheckList/CheckList.widget';
import { contextBudgetBarWidget } from '../components/molecules/ContextBudgetBar/ContextBudgetBar.widget';
import { contextLegendWidget } from '../components/molecules/ContextLegend/ContextLegend.widget';
import { contextStackDiagramWidget } from '../components/molecules/ContextStackDiagram/ContextStackDiagram.widget';
import { contextStripDiagramWidget } from '../components/molecules/ContextStripDiagram/ContextStripDiagram.widget';
import { contextTreemapWidget } from '../components/molecules/ContextTreemap/ContextTreemap.widget';
import { dataTableWidget } from '../components/molecules/DataTable/DataTable.widget';
import { figureBlockWidget } from '../components/molecules/FigureBlock/FigureBlock.widget';
import { keyPointListWidget } from '../components/molecules/KeyPointList/KeyPointList.widget';
import { keyValueStripWidget } from '../components/molecules/KeyValueStrip/KeyValueStrip.widget';
import { metadataGridWidget } from '../components/molecules/MetadataGrid/MetadataGrid.widget';
import { personSummaryWidget } from '../components/molecules/PersonSummary/PersonSummary.widget';
import { sidebarNavWidget } from '../components/molecules/SidebarNav/SidebarNav.widget';
import { stepListWidget } from '../components/molecules/StepList/StepList.widget';
import { transcriptMessageCardWidget } from '../components/molecules/TranscriptMessageCard/TranscriptMessageCard.widget';
import { transcriptSessionHeaderWidget } from '../components/molecules/TranscriptSessionHeader/TranscriptSessionHeader.widget';
import { anchoredCommentRailWidget } from '../components/organisms/AnchoredCommentRail/AnchoredCommentRail.widget';
import { annotationRailPanelWidget } from '../components/organisms/AnnotationRailPanel/AnnotationRailPanel.widget';
import { contextDiagramPanelWidget } from '../components/organisms/ContextDiagramPanel/ContextDiagramPanel.widget';
import { contextUploadDropAreaWidget } from '../components/organisms/ContextUploadDropArea/ContextUploadDropArea.widget';
import { courseStudioShellWidget } from '../components/organisms/CourseStudioShell/CourseStudioShell.widget';
import { transcriptReaderPanelWidget } from '../components/organisms/TranscriptReaderPanel/TranscriptReaderPanel.widget';
import { transcriptWorkspacePanelWidget } from '../components/organisms/TranscriptWorkspacePanel/TranscriptWorkspacePanel.widget';
import { createWidgetRegistry, mergeWidgetRegistries } from './registry';

export const uiWidgetRegistry = createWidgetRegistry([
  appNavWidget,
  appShellWidget,
  buttonWidget,
  captionWidget,
  checkListWidget,
  codeTextWidget,
  dashboardGridWidget,
  dividerWidget,
  figureBlockWidget,
  formRowWidget,
  inlineWidget,
  keyPointListWidget,
  keyValueStripWidget,
  metadataGridWidget,
  panelWidget,
  personSummaryWidget,
  scrollRegionWidget,
  sectionBlockWidget,
  selectInputWidget,
  sidebarNavWidget,
  sidebarShellWidget,
  splitPaneWidget,
  stackWidget,
  statusTextWidget,
  stepListWidget,
  tabListWidget,
  textInputWidget,
  textWidget,
]);

export const dataWidgetRegistry = createWidgetRegistry([
  dataTableWidget,
]);

export const contextWindowWidgetRegistry = createWidgetRegistry([
  anchoredCommentCardWidget,
  anchoredCommentRailWidget,
  annotationBadgeWidget,
  annotationNoteCardWidget,
  annotationRailPanelWidget,
  contextBudgetBarWidget,
  contextDiagramPanelWidget,
  contextKindSwatchWidget,
  contextLegendWidget,
  contextStackDiagramWidget,
  contextStripDiagramWidget,
  contextTreemapWidget,
  contextUploadDropAreaWidget,
  transcriptMessageCardWidget,
  transcriptReaderPanelWidget,
  transcriptRoleBadgeWidget,
  transcriptSessionHeaderWidget,
  transcriptWorkspacePanelWidget,
]);

export const courseWidgetRegistry = createWidgetRegistry([
  courseStudioShellWidget,
]);

export const defaultWidgetRegistry = mergeWidgetRegistries(
  uiWidgetRegistry,
  dataWidgetRegistry,
  contextWindowWidgetRegistry,
  courseWidgetRegistry,
);
