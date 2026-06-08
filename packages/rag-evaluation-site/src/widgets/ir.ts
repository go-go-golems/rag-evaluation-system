import type { CaptionTone, CaptionTransform, RagStatus, TextAlign, TextAs, TextSize, TextTone, TextWeight } from '../components/foundation';
import type { ButtonSize, ButtonVariant, ContextKindSwatchSize } from '../components/atoms';
import type { DashboardGridRecipe, InlineGap, InlineJustify, StackAlign, StackGap } from '../components/layout';
import type { AnchoredComment, ContextDiagramStyle, ContextDiagramView, ContextPartKind, ContextWindowSnapshot, TranscriptAnnotation, TranscriptMessage, TranscriptRole } from '../context';

export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };
export type JsonObject = { [key: string]: JsonValue };

export type WidgetNode = TextNode | ElementNode | ComponentNode;

export interface TextNode {
  kind: 'text';
  text: string;
}

export interface ElementNode {
  kind: 'element';
  tag: string;
  attrs?: JsonObject;
  children?: WidgetNode[];
}

export type RagWidgetType =
  | 'AppShell'
  | 'AppNav'
  | 'Button'
  | 'Caption'
  | 'CodeText'
  | 'ContextKindSwatch'
  | 'AnnotationBadge'
  | 'ContextLegend'
  | 'ContextBudgetBar'
  | 'ContextStripDiagram'
  | 'ContextStackDiagram'
  | 'ContextTreemap'
  | 'ContextDiagramPanel'
  | 'DashboardGrid'
  | 'DataTable'
  | 'Divider'
  | 'FormRow'
  | 'Inline'
  | 'MetadataGrid'
  | 'Panel'
  | 'ScrollRegion'
  | 'SectionBlock'
  | 'SelectInput'
  | 'SidebarShell'
  | 'SplitPane'
  | 'Stack'
  | 'StatusText'
  | 'TabList'
  | 'Text'
  | 'TextInput'
  | 'TranscriptRoleBadge'
  | 'TranscriptSessionHeader'
  | 'TranscriptMessageCard'
  | 'AnnotationNoteCard'
  | 'AnnotationRailPanel'
  | 'TranscriptReaderPanel'
  | 'TranscriptWorkspacePanel'
  | 'AnchoredCommentCard'
  | 'AnchoredCommentRail';

export interface ComponentNode {
  kind: 'component';
  type: RagWidgetType | string;
  props?: WidgetProps;
  children?: WidgetNode[];
}

export type RenderableValue = WidgetNode | string | number | boolean | null;

export type ActionSpec =
  | NavigateActionSpec
  | ServerActionSpec
  | EventActionSpec
  | CopyActionSpec;

export interface NavigateActionSpec {
  kind: 'navigate';
  to: string;
  params?: JsonObject;
}

export interface ServerActionSpec {
  kind: 'server';
  name: string;
  payload?: JsonObject;
}

export interface EventActionSpec {
  kind: 'event';
  event: string;
  detail?: JsonObject;
}

export interface CopyActionSpec {
  kind: 'copy';
  value?: string;
  field?: string;
}

export interface BaseWidgetProps {
  className?: string;
  style?: JsonObject;
  id?: string;
  action?: ActionSpec;
  [key: string]: unknown;
}

export interface AppShellWidgetProps extends BaseWidgetProps {
  header?: WidgetNode;
  sidebar?: WidgetNode;
}

export interface AppNavWidgetProps extends BaseWidgetProps {
  brand: RenderableValue;
  items: AppNavItemSpec[];
  activeItemId: string;
}

export interface AppNavItemSpec {
  id: string;
  label: RenderableValue;
  action?: ActionSpec;
}

export interface ButtonWidgetProps extends BaseWidgetProps {
  variant?: ButtonVariant;
  size?: ButtonSize;
  selected?: boolean;
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
}

export interface TextWidgetProps extends BaseWidgetProps {
  as?: TextAs;
  size?: TextSize;
  tone?: TextTone;
  weight?: TextWeight;
  align?: TextAlign;
  truncate?: boolean;
}

export interface CodeTextWidgetProps extends BaseWidgetProps {
  as?: 'code' | 'span' | 'div';
  tone?: 'primary' | 'muted' | 'accent';
  display?: 'inline' | 'block';
  copyable?: boolean;
}

export interface DividerWidgetProps extends BaseWidgetProps {
  orientation?: 'horizontal' | 'vertical';
}

export interface ContextKindSwatchWidgetProps extends BaseWidgetProps {
  kind: ContextPartKind;
  mode?: ContextDiagramStyle;
  size?: ContextKindSwatchSize;
  selected?: boolean;
}

export interface AnnotationBadgeWidgetProps extends BaseWidgetProps {
  kind: ContextPartKind;
  label: string;
  selected?: boolean;
}

export interface TranscriptRoleBadgeWidgetProps extends BaseWidgetProps {
  role: TranscriptRole;
  name?: string;
}

export interface ContextLegendWidgetProps extends BaseWidgetProps {
  kinds?: ContextPartKind[];
  mode?: ContextDiagramStyle;
  compact?: boolean;
  selectedKind?: ContextPartKind;
}

export interface ContextBudgetBarWidgetProps extends BaseWidgetProps {
  snapshot: ContextWindowSnapshot;
  mode?: ContextDiagramStyle;
  showLegend?: boolean;
  selectedPartId?: string;
}

export interface ContextStripDiagramWidgetProps extends BaseWidgetProps {
  snapshot: ContextWindowSnapshot;
  mode?: ContextDiagramStyle;
  selectedPartId?: string;
  showLabels?: boolean;
}

export interface ContextStackDiagramWidgetProps extends BaseWidgetProps {
  snapshot: ContextWindowSnapshot;
  selectedPartId?: string;
}

export interface ContextTreemapWidgetProps extends BaseWidgetProps {
  snapshot: ContextWindowSnapshot;
  selectedPartId?: string;
}

export interface ContextDiagramPanelWidgetProps extends BaseWidgetProps {
  snapshot: ContextWindowSnapshot;
  initialView?: ContextDiagramView;
  selectedPartId?: string;
}

export interface TranscriptSessionHeaderWidgetProps extends BaseWidgetProps {
  title: RenderableValue;
  subtitle?: RenderableValue;
  messageCount?: number;
  annotationCount?: number;
  tokenTotal?: number;
  rightSlot?: WidgetNode;
}

export interface TranscriptMessageCardWidgetProps extends BaseWidgetProps {
  message: TranscriptMessage;
  annotations?: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  showAnnotationChips?: boolean;
  onAnnotationSelectAction?: ActionSpec;
}

export interface AnnotationNoteCardWidgetProps extends BaseWidgetProps {
  annotation: TranscriptAnnotation;
  selected?: boolean;
  index?: number;
}

export interface AnnotationRailPanelWidgetProps extends BaseWidgetProps {
  title?: string;
  description?: string;
  annotations: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  onAnnotationSelectAction?: ActionSpec;
}

export interface TranscriptReaderPanelWidgetProps extends BaseWidgetProps {
  title?: string;
  subtitle?: string;
  messages: TranscriptMessage[];
  annotations?: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  showAnnotationChips?: boolean;
  onAnnotationSelectAction?: ActionSpec;
}

export interface TranscriptWorkspacePanelWidgetProps extends BaseWidgetProps {
  title?: string;
  subtitle?: string;
  messages: TranscriptMessage[];
  annotations?: TranscriptAnnotation[];
  selectedAnnotationId?: string;
  showNotes?: boolean;
  onAnnotationSelectAction?: ActionSpec;
}

export interface AnchoredCommentCardWidgetProps extends BaseWidgetProps {
  comment: AnchoredComment;
  index?: number;
  selected?: boolean;
  compact?: boolean;
  onDismissAction?: ActionSpec;
}

export interface AnchoredCommentRailWidgetProps extends BaseWidgetProps {
  title?: string;
  comments: AnchoredComment[];
  selectedCommentId?: string;
  onCommentSelectAction?: ActionSpec;
}

export interface CaptionWidgetProps extends BaseWidgetProps {
  tone?: CaptionTone;
  transform?: CaptionTransform;
  truncate?: boolean;
}

export interface DashboardGridWidgetProps extends BaseWidgetProps {
  recipe?: DashboardGridRecipe;
}

export interface DataTableWidgetProps extends BaseWidgetProps {
  columns: DataTableColumnSpec[];
  rows: JsonObject[];
  getRowKey: RowKeySpec;
  selectedKey?: string | null;
  onRowSelect?: ActionSpec;
  emptyMessage?: RenderableValue;
}

export interface DataTableColumnSpec {
  id: string;
  header: RenderableValue;
  cell: CellSpec;
  align?: 'start' | 'end' | 'center';
  maxWidth?: number | string;
}

export type RowKeySpec = string | { field: string } | { template: string };

export type CellSpec =
  | FieldCellSpec
  | NumberCellSpec
  | StatusCellSpec
  | CaptionCellSpec
  | TemplateCellSpec
  | LinkCellSpec
  | ConstantCellSpec;

export interface FieldCellSpec {
  kind: 'field';
  field: string;
  fallback?: string;
}

export interface NumberCellSpec {
  kind: 'number';
  field: string;
  format?: 'integer' | 'fixed';
  digits?: number;
  fallback?: string;
}

export interface StatusCellSpec {
  kind: 'status';
  field: string;
  icon?: boolean;
  fallback?: RagStatus | string;
}

export interface CaptionCellSpec {
  kind: 'caption';
  field: string;
  tone?: CaptionTone;
  fallback?: string;
}

export interface TemplateCellSpec {
  kind: 'template';
  template: string;
}

export interface LinkCellSpec {
  kind: 'link';
  hrefField: string;
  labelField: string;
  target?: '_blank' | '_self' | '_parent' | '_top';
  fallbackLabel?: string;
}

export interface ConstantCellSpec {
  kind: 'constant';
  value: RenderableValue;
}

export interface FormRowWidgetProps extends BaseWidgetProps {
  label: RenderableValue;
  control: WidgetNode;
  hint?: RenderableValue;
}

export interface InlineWidgetProps extends BaseWidgetProps {
  gap?: InlineGap;
  justify?: InlineJustify;
  wrap?: boolean;
}

export interface MetadataGridWidgetProps extends BaseWidgetProps {
  items: MetadataGridItemSpec[];
  density?: 'normal' | 'compact';
}

export interface MetadataGridItemSpec {
  key: RenderableValue;
  value: RenderableValue;
  copyValue?: string;
}

export interface PanelWidgetProps extends BaseWidgetProps {
  title?: RenderableValue;
  actions?: RenderableValue;
  density?: 'normal' | 'condensed';
  fill?: boolean;
}

export interface ScrollRegionWidgetProps extends BaseWidgetProps {
  axis?: 'y' | 'x' | 'both';
}

export interface SectionBlockWidgetProps extends BaseWidgetProps {
  as?: 'section' | 'article' | 'div';
  label?: RenderableValue;
  caption?: RenderableValue;
  density?: 'normal' | 'spacious';
  divider?: 'none' | 'top' | 'bottom' | 'both';
}

export interface SplitPaneWidgetProps extends BaseWidgetProps {
  left: WidgetNode;
  right: WidgetNode;
  ratio?: 'balanced' | 'leftNarrow' | 'rightNarrow' | 'course';
  divider?: boolean;
}

export interface SidebarShellWidgetProps extends BaseWidgetProps {
  sidebar?: WidgetNode;
  sidebarWidth?: number;
  header?: WidgetNode;
  footer?: WidgetNode;
}

export interface SelectInputWidgetProps extends BaseWidgetProps {
  name?: string;
  value?: string | number;
  disabled?: boolean;
  options?: SelectOptionSpec[];
}

export interface SelectOptionSpec {
  value: string | number;
  label: RenderableValue;
  selected?: boolean;
  disabled?: boolean;
}

export interface StackWidgetProps extends BaseWidgetProps {
  gap?: StackGap;
  align?: StackAlign;
}

export interface StatusTextWidgetProps extends BaseWidgetProps {
  status: RagStatus | string;
  icon?: boolean;
}

export interface TabListWidgetProps extends BaseWidgetProps {
  items: TabListItemSpec[];
  activeId: string;
  ariaLabel?: string;
  onChange?: ActionSpec;
}

export interface TabListItemSpec {
  id: string;
  label: RenderableValue;
}

export interface TextInputWidgetProps extends BaseWidgetProps {
  name?: string;
  value?: string | number;
  placeholder?: string;
  type?: string;
  disabled?: boolean;
  min?: number;
  max?: number;
}

export type WidgetProps =
  | AppShellWidgetProps
  | AppNavWidgetProps
  | ButtonWidgetProps
  | TextWidgetProps
  | CodeTextWidgetProps
  | DividerWidgetProps
  | ContextKindSwatchWidgetProps
  | AnnotationBadgeWidgetProps
  | TranscriptRoleBadgeWidgetProps
  | ContextLegendWidgetProps
  | ContextBudgetBarWidgetProps
  | ContextStripDiagramWidgetProps
  | ContextStackDiagramWidgetProps
  | ContextTreemapWidgetProps
  | ContextDiagramPanelWidgetProps
  | TranscriptSessionHeaderWidgetProps
  | TranscriptMessageCardWidgetProps
  | AnnotationNoteCardWidgetProps
  | AnnotationRailPanelWidgetProps
  | TranscriptReaderPanelWidgetProps
  | TranscriptWorkspacePanelWidgetProps
  | AnchoredCommentCardWidgetProps
  | AnchoredCommentRailWidgetProps
  | CaptionWidgetProps
  | DashboardGridWidgetProps
  | DataTableWidgetProps
  | FormRowWidgetProps
  | InlineWidgetProps
  | MetadataGridWidgetProps
  | PanelWidgetProps
  | ScrollRegionWidgetProps
  | SectionBlockWidgetProps
  | SplitPaneWidgetProps
  | SidebarShellWidgetProps
  | SelectInputWidgetProps
  | StackWidgetProps
  | StatusTextWidgetProps
  | TabListWidgetProps
  | TextInputWidgetProps
  | BaseWidgetProps;

export function text(value: string | number | boolean): TextNode {
  return { kind: 'text', text: String(value) };
}

export function element(tag: string, attrs?: JsonObject, children: WidgetNode[] = []): ElementNode {
  return { kind: 'element', tag, attrs, children };
}

export function component(type: RagWidgetType | string, props?: WidgetProps, children: WidgetNode[] = []): ComponentNode {
  return { kind: 'component', type, props, children };
}

export function isWidgetNode(value: unknown): value is WidgetNode {
  if (!value || typeof value !== 'object') return false;
  const kind = (value as { kind?: unknown }).kind;
  return kind === 'text' || kind === 'element' || kind === 'component';
}
