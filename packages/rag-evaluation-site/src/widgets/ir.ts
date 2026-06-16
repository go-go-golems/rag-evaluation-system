import type {
	CaptionTone,
	CaptionTransform,
	RagStatus,
	TextAlign,
	TextAs,
	TextSize,
	TextTone,
	TextWeight,
} from "../components/foundation";
import type { ButtonSize, ButtonVariant, ContextStudioNavIconId } from "../components/atoms";
import type {
	DashboardGridRecipe,
	InlineGap,
	InlineJustify,
	StackAlign,
	StackGap,
} from "../components/layout";
import type {
	AnchoredComment,
	ArticleBlock,
	ContextCourse,
	ContextCourseAgendaItem,
	ContextDiagramView,
	ContextHandoutDocument,
	ContextLegendItemSpec,
	ContextSlide,
	ContextStyleSet,
	ContextVisualStyle,
	ContextWindowSnapshot,
	TranscriptAnnotation,
	TranscriptMessage,
	TranscriptRole,
} from "../context";

export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | JsonValue[] | { [key: string]: JsonValue };
export type JsonObject = { [key: string]: JsonValue };

export type WidgetNode = TextNode | ElementNode | ComponentNode;

export interface TextNode {
	kind: "text";
	text: string;
}

export interface ElementNode {
	kind: "element";
	tag: string;
	attrs?: JsonObject;
	children?: WidgetNode[];
}

export type RagWidgetType =
	| "AppShell"
	| "AppNav"
	| "Button"
	| "Caption"
	| "CodeText"
	| "ContextStyleSwatch"
	| "ContextStudioNavIcon"
	| "AnnotationBadge"
	| "ContextLegend"
	| "ContextBudgetBar"
	| "ContextStripDiagram"
	| "ContextGroupedStripDiagram"
	| "ContextStackDiagram"
	| "ContextTreemap"
	| "ContextDiagramPanel"
	| "DashboardGrid"
	| "DataTable"
	| "Divider"
	| "FormPanel"
	| "FormRow"
	| "Inline"
	| "MetadataGrid"
	| "Panel"
	| "ScrollRegion"
	| "SectionBlock"
	| "SelectInput"
	| "SidebarShell"
	| "SlideShell"
	| "SplitPane"
	| "Stack"
	| "StatusText"
	| "TabList"
	| "Text"
	| "TextInput"
	| "TextareaInput"
	| "TranscriptRoleBadge"
	| "TranscriptSessionHeader"
	| "TranscriptMessageCard"
	| "AnnotationNoteCard"
	| "AnnotationRailPanel"
	| "TranscriptReaderPanel"
	| "TranscriptWorkspacePanel"
	| "AnchoredCommentCard"
	| "AnchoredCommentRail"
	| "KeyValueStrip"
	| "CheckList"
	| "StepList"
	| "PersonSummary"
	| "FigureBlock"
	| "KeyPointList"
	| "SidebarNav"
	| "CourseStepNav"
	| "MarkdownArticle"
	| "RichArticle"
	| "DocumentListPanel"
	| "DocumentPreviewToolbar"
	| "CourseLessonPanel"
	| "CourseSlidePanel"
	| "CourseStudioShell"
	| "HandoutDocumentShell"
	| "ContextUploadDropArea";

export interface ComponentNode {
	kind: "component";
	type: RagWidgetType | string;
	props?: WidgetProps;
	children?: WidgetNode[];
}

export type RenderableValue = WidgetNode | string | number | boolean | null;

export type ActionSpec =
	| NavigateActionSpec
	| DownloadActionSpec
	| ServerActionSpec
	| EventActionSpec
	| CopyActionSpec;

export interface NavigateActionSpec {
	kind: "navigate";
	to: string;
	params?: JsonObject;
}

export interface DownloadActionSpec {
	kind: "download";
	to: string;
	params?: JsonObject;
}

export interface ServerActionSpec {
	kind: "server";
	name: string;
	payload?: JsonObject;
}

export interface EventActionSpec {
	kind: "event";
	event: string;
	detail?: JsonObject;
}

export interface CopyActionSpec {
	kind: "copy";
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
	type?: "button" | "submit" | "reset";
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
	as?: "code" | "span" | "div";
	tone?: "primary" | "muted" | "accent";
	display?: "inline" | "block";
	copyable?: boolean;
}

export interface DividerWidgetProps extends BaseWidgetProps {
	orientation?: "horizontal" | "vertical";
}

export interface ContextStudioNavIconWidgetProps extends BaseWidgetProps {
	id: ContextStudioNavIconId;
	title?: string;
}

export interface ContextStyleSwatchWidgetProps extends BaseWidgetProps {
	visualStyle: ContextVisualStyle;
	size?: "xs" | "sm" | "md" | "lg";
	selected?: boolean;
}

export interface AnnotationBadgeWidgetProps extends BaseWidgetProps {
	visualStyle: ContextVisualStyle;
	label: string;
	selected?: boolean;
}

export interface TranscriptRoleBadgeWidgetProps extends BaseWidgetProps {
	role: TranscriptRole;
	name?: string;
}

export interface ContextLegendWidgetProps extends BaseWidgetProps {
	items: ContextLegendItemSpec[];
	styles: Record<string, ContextVisualStyle>;
	size?: "xs" | "sm" | "md" | "lg";
	selectedId?: string;
}

export interface ContextBudgetBarWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	showLegend?: boolean;
	selectedPartId?: string;
}

export interface ContextStripDiagramWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
	showLabels?: boolean;
	showSelection?: boolean;
}

export interface ContextGroupedStripDiagramWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
	groupBy?: "turn" | "styleKey" | "sourceId";
	showGroupLabels?: boolean;
	showPartLabels?: boolean;
	showSelection?: boolean;
}

export interface ContextStackDiagramWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
}

export interface ContextTreemapWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	selectedPartId?: string;
}

export interface ContextDiagramPanelWidgetProps extends BaseWidgetProps {
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
	initialView?: ContextDiagramView;
	selectedPartId?: string;
	views?: ContextDiagramView[];
	showLegend?: boolean;
	showPartDetails?: boolean;
	chrome?: "panel" | "inline";
}

export interface ContextTurnPagerPanelWidgetProps extends BaseWidgetProps {
	snapshots: ContextWindowSnapshot[];
	styleSet: ContextStyleSet;
	initialSnapshotId?: string;
	selectedPartId?: string;
	diagram?: "grouped-strip" | "strip";
	groupBy?: "turn" | "styleKey" | "sourceId";
	mode?: "turn-only" | "snapshot";
	includeGlobalParts?: boolean;
	showLegend?: boolean;
	title?: string;
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
	styleSet?: ContextStyleSet;
	onAnnotationSelectAction?: ActionSpec;
}

export interface AnnotationNoteCardWidgetProps extends BaseWidgetProps {
	annotation: TranscriptAnnotation;
	selected?: boolean;
	index?: number;
	styleSet?: ContextStyleSet;
}

export interface AnnotationRailPanelWidgetProps extends BaseWidgetProps {
	title?: string;
	description?: string;
	annotations: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	styleSet?: ContextStyleSet;
	onAnnotationSelectAction?: ActionSpec;
}

export interface TranscriptReaderPanelWidgetProps extends BaseWidgetProps {
	title?: string;
	subtitle?: string;
	messages: TranscriptMessage[];
	annotations?: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	showAnnotationChips?: boolean;
	styleSet?: ContextStyleSet;
	onAnnotationSelectAction?: ActionSpec;
}

export interface TranscriptWorkspacePanelWidgetProps extends BaseWidgetProps {
	title?: string;
	subtitle?: string;
	messages: TranscriptMessage[];
	annotations?: TranscriptAnnotation[];
	selectedAnnotationId?: string;
	showNotes?: boolean;
	styleSet?: ContextStyleSet;
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

export interface KeyValueStripWidgetProps extends BaseWidgetProps {
	items: Array<{ key: RenderableValue; value: RenderableValue }>;
}
export interface CheckListWidgetProps extends BaseWidgetProps {
	items: Array<{ id?: string; text: RenderableValue } | RenderableValue>;
	marker?: RenderableValue;
}
export interface StepListWidgetProps extends BaseWidgetProps {
	items: Array<{
		id: string;
		index: RenderableValue;
		title: RenderableValue;
		description?: RenderableValue;
		meta?: RenderableValue;
	}>;
	activeItemId?: string;
	onItemSelectAction?: ActionSpec;
}
export interface PersonSummaryWidgetProps extends BaseWidgetProps {
	name: RenderableValue;
	subtitle?: RenderableValue;
	bio?: RenderableValue;
	avatar?: RenderableValue;
}
export interface FigureBlockWidgetProps extends BaseWidgetProps {
	as?: "figure" | "div";
	caption?: RenderableValue;
	legend?: WidgetNode;
	frame?: "none" | "bordered" | "inset";
}
export interface KeyPointListWidgetProps extends BaseWidgetProps {
	items: Array<
		| {
				id?: string;
				index?: RenderableValue;
				title?: RenderableValue;
				text: RenderableValue;
				meta?: RenderableValue;
		  }
		| RenderableValue
	>;
	markerTone?: "accent" | "muted";
}

export interface SidebarNavItemWidgetSpec {
	id: string;
	label: RenderableValue;
	icon?: RenderableValue;
	badge?: RenderableValue;
	disabled?: boolean;
}
export interface SidebarNavSectionWidgetSpec {
	id: string;
	label: RenderableValue;
	items: SidebarNavItemWidgetSpec[];
}
export interface SidebarNavWidgetProps extends BaseWidgetProps {
	sections: SidebarNavSectionWidgetSpec[];
	activeItemId?: string;
	onItemSelectAction?: ActionSpec;
	ariaLabel?: string;
}

export interface CourseStepNavWidgetProps extends BaseWidgetProps {
	items: ContextCourseAgendaItem[];
	activeItemId?: string;
	onItemSelectAction?: ActionSpec;
}
export interface MarkdownArticleWidgetProps extends BaseWidgetProps {
	source: string;
}
export interface RichArticleWidgetProps extends BaseWidgetProps {
	blocks: ArticleBlock[];
	styleSet?: ContextStyleSet;
}
export interface DocumentListItemWidgetSpec {
	id: string;
	title: RenderableValue;
	format?: RenderableValue;
	size?: RenderableValue;
	description?: RenderableValue;
	icon?: RenderableValue;
	disabled?: boolean;
}
export interface DocumentListPanelWidgetProps extends BaseWidgetProps {
	title: RenderableValue;
	description?: RenderableValue;
	items: DocumentListItemWidgetSpec[];
	selectedItemId?: string;
	onItemSelectAction?: ActionSpec;
	onDownloadAllAction?: ActionSpec;
	downloadAllLabel?: RenderableValue;
	emptyMessage?: RenderableValue;
	showItemDescriptions?: boolean;
}
export interface DocumentPreviewToolbarWidgetProps extends BaseWidgetProps {
	file: RenderableValue;
	format?: RenderableValue;
	size?: RenderableValue;
	onDownloadAction?: ActionSpec;
	downloadLabel?: RenderableValue;
	onPrintAction?: ActionSpec;
	printLabel?: RenderableValue;
	rightSlot?: WidgetNode;
}
export interface CourseLessonPanelWidgetProps extends BaseWidgetProps {
	course: ContextCourse;
	activeAgendaItemId?: string;
	onAgendaItemSelectAction?: ActionSpec;
	onPrimaryCtaAction?: ActionSpec;
	onSecondaryCtaAction?: ActionSpec;
}
export interface CourseSlidePanelWidgetProps extends BaseWidgetProps {
	slide: ContextSlide;
	snapshot: ContextWindowSnapshot;
	styleSet?: ContextStyleSet;
	index?: number;
	total?: number;
	visualSide?: "left" | "right";
	mode?: "default" | "present";
	onPreviousAction?: ActionSpec;
	onNextAction?: ActionSpec;
	onPresentAction?: ActionSpec;
	onFullscreenAction?: ActionSpec;
}
export interface CourseStudioShellWidgetProps extends BaseWidgetProps {
	sections: SidebarNavSectionWidgetSpec[];
	activeItemId?: string;
	onNavigateAction?: ActionSpec;
	title?: RenderableValue;
	subtitle?: RenderableValue;
	sidebarFooter?: WidgetNode;
	contentPadding?: "default" | "none";
}
export interface HandoutDocumentShellWidgetProps extends BaseWidgetProps {
	intro: RenderableValue;
	documents: ContextHandoutDocument[];
	selectedDocumentId?: string;
	onDocumentSelectAction?: ActionSpec;
	onDownloadAction?: ActionSpec;
	onDownloadAllAction?: ActionSpec;
	onPrintAction?: ActionSpec;
	styleSet?: ContextStyleSet;
	title?: RenderableValue;
	emptyMessage?: RenderableValue;
}
export interface ContextUploadDropAreaWidgetProps extends BaseWidgetProps {
	title?: RenderableValue;
	description?: RenderableValue;
	accept?: string;
	disabled?: boolean;
	active?: boolean;
	onFilesSelectedAction?: ActionSpec;
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
	align?: "start" | "end" | "center";
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
	| LinkButtonCellSpec
	| ActionButtonCellSpec
	| ConstantCellSpec;

export interface FieldCellSpec {
	kind: "field";
	field: string;
	fallback?: string;
}

export interface NumberCellSpec {
	kind: "number";
	field: string;
	format?: "integer" | "fixed";
	digits?: number;
	fallback?: string;
}

export interface StatusCellSpec {
	kind: "status";
	field: string;
	icon?: boolean;
	fallback?: RagStatus | string;
}

export interface CaptionCellSpec {
	kind: "caption";
	field: string;
	tone?: CaptionTone;
	fallback?: string;
}

export interface TemplateCellSpec {
	kind: "template";
	template: string;
}

export interface LinkCellSpec {
	kind: "link";
	hrefField: string;
	labelField: string;
	target?: "_blank" | "_self" | "_parent" | "_top";
	fallbackLabel?: string;
}

export interface LinkButtonCellSpec {
	kind: "linkButton";
	hrefField: string;
	labelField: string;
	variant?: ButtonVariant;
	size?: ButtonSize;
	fallbackLabel?: string;
}

export interface ActionButtonCellSpec {
	kind: "actionButton";
	label: RenderableValue;
	action: ActionSpec;
	variant?: ButtonVariant;
	size?: ButtonSize;
	disabled?: boolean;
}

export interface ConstantCellSpec {
	kind: "constant";
	value: RenderableValue;
}

export interface FormPanelWidgetProps extends BaseWidgetProps {
	title: RenderableValue;
	subtitle?: RenderableValue;
	actions?: RenderableValue;
	formAction?: string;
	method?: "get" | "post";
	status?: "idle" | "saving" | "success" | "error";
	statusMessage?: RenderableValue;
	submitLabel?: RenderableValue;
	resetLabel?: RenderableValue;
	footer?: RenderableValue;
	disabled?: boolean;
}

export interface FormRowWidgetProps extends BaseWidgetProps {
	label: RenderableValue;
	control: WidgetNode;
	description?: RenderableValue;
	hint?: RenderableValue;
	error?: RenderableValue;
	success?: RenderableValue;
	counter?: RenderableValue;
	required?: boolean;
	orientation?: "inline" | "stacked";
}

export interface InlineWidgetProps extends BaseWidgetProps {
	gap?: InlineGap;
	justify?: InlineJustify;
	wrap?: boolean;
}

export interface MetadataGridWidgetProps extends BaseWidgetProps {
	items: MetadataGridItemSpec[];
	density?: "normal" | "compact";
}

export interface MetadataGridItemSpec {
	key: RenderableValue;
	value: RenderableValue;
	copyValue?: string;
}

export interface PanelWidgetProps extends BaseWidgetProps {
	title?: RenderableValue;
	actions?: RenderableValue;
	density?: "normal" | "condensed";
	fill?: boolean;
}

export interface ScrollRegionWidgetProps extends BaseWidgetProps {
	axis?: "y" | "x" | "both";
}

export interface SectionBlockWidgetProps extends BaseWidgetProps {
	as?: "section" | "article" | "div";
	label?: RenderableValue;
	caption?: RenderableValue;
	density?: "normal" | "spacious";
	divider?: "none" | "top" | "bottom" | "both";
}

export interface SplitPaneWidgetProps extends BaseWidgetProps {
	left: WidgetNode;
	right: WidgetNode;
	ratio?: "balanced" | "leftNarrow" | "rightNarrow" | "course" | "sidebar";
	divider?: boolean;
	gutter?: "none" | "md" | "lg";
}

export interface SidebarShellWidgetProps extends BaseWidgetProps {
	sidebar?: WidgetNode;
	sidebarWidth?: number;
	contentPadding?: "none" | "md" | "lg";
	header?: WidgetNode;
	footer?: WidgetNode;
}

export interface SlideShellWidgetProps extends BaseWidgetProps {
	as?: "article" | "section" | "div";
	eyebrow?: RenderableValue;
	counter?: RenderableValue;
	title: RenderableValue;
	subtitle?: RenderableValue;
	primary: WidgetNode;
	secondary?: WidgetNode;
	primarySide?: "left" | "right";
	ratio?: "balanced" | "primaryWide" | "secondaryWide";
	divider?: boolean;
	footer?: WidgetNode;
}

export interface SelectInputWidgetProps extends BaseWidgetProps {
	name?: string;
	value?: string | number;
	defaultValue?: string | number;
	disabled?: boolean;
	required?: boolean;
	ariaInvalid?: boolean;
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
	defaultValue?: string | number;
	placeholder?: string;
	type?: string;
	disabled?: boolean;
	readOnly?: boolean;
	required?: boolean;
	min?: number;
	max?: number;
	minLength?: number;
	maxLength?: number;
	autoComplete?: string;
	ariaInvalid?: boolean;
}

export interface TextareaInputWidgetProps extends BaseWidgetProps {
	name?: string;
	value?: string | number;
	defaultValue?: string | number;
	placeholder?: string;
	disabled?: boolean;
	readOnly?: boolean;
	required?: boolean;
	minLength?: number;
	maxLength?: number;
	rows?: number;
	resize?: "vertical" | "none";
	ariaInvalid?: boolean;
}

export type WidgetProps =
	| AppShellWidgetProps
	| AppNavWidgetProps
	| ButtonWidgetProps
	| TextWidgetProps
	| CodeTextWidgetProps
	| DividerWidgetProps
	| ContextStyleSwatchWidgetProps
	| ContextStudioNavIconWidgetProps
	| AnnotationBadgeWidgetProps
	| TranscriptRoleBadgeWidgetProps
	| ContextLegendWidgetProps
	| ContextBudgetBarWidgetProps
	| ContextStripDiagramWidgetProps
	| ContextGroupedStripDiagramWidgetProps
	| ContextStackDiagramWidgetProps
	| ContextTreemapWidgetProps
	| ContextDiagramPanelWidgetProps
	| ContextTurnPagerPanelWidgetProps
	| TranscriptSessionHeaderWidgetProps
	| TranscriptMessageCardWidgetProps
	| AnnotationNoteCardWidgetProps
	| AnnotationRailPanelWidgetProps
	| TranscriptReaderPanelWidgetProps
	| TranscriptWorkspacePanelWidgetProps
	| AnchoredCommentCardWidgetProps
	| AnchoredCommentRailWidgetProps
	| KeyValueStripWidgetProps
	| CheckListWidgetProps
	| StepListWidgetProps
	| PersonSummaryWidgetProps
	| FigureBlockWidgetProps
	| KeyPointListWidgetProps
	| SidebarNavWidgetProps
	| CourseStepNavWidgetProps
	| MarkdownArticleWidgetProps
	| RichArticleWidgetProps
	| DocumentListPanelWidgetProps
	| DocumentPreviewToolbarWidgetProps
	| CourseLessonPanelWidgetProps
	| CourseSlidePanelWidgetProps
	| CourseStudioShellWidgetProps
	| HandoutDocumentShellWidgetProps
	| ContextUploadDropAreaWidgetProps
	| CaptionWidgetProps
	| DashboardGridWidgetProps
	| DataTableWidgetProps
	| FormPanelWidgetProps
	| FormRowWidgetProps
	| InlineWidgetProps
	| MetadataGridWidgetProps
	| PanelWidgetProps
	| ScrollRegionWidgetProps
	| SectionBlockWidgetProps
	| SplitPaneWidgetProps
	| SidebarShellWidgetProps
	| SlideShellWidgetProps
	| SelectInputWidgetProps
	| StackWidgetProps
	| StatusTextWidgetProps
	| TabListWidgetProps
	| TextInputWidgetProps
	| TextareaInputWidgetProps
	| BaseWidgetProps;

export function text(value: string | number | boolean): TextNode {
	return { kind: "text", text: String(value) };
}

export function element(tag: string, attrs?: JsonObject, children: WidgetNode[] = []): ElementNode {
	return { kind: "element", tag, attrs, children };
}

export function component(
	type: RagWidgetType | string,
	props?: WidgetProps,
	children: WidgetNode[] = [],
): ComponentNode {
	return { kind: "component", type, props, children };
}

export function isWidgetNode(value: unknown): value is WidgetNode {
	if (!value || typeof value !== "object") return false;
	const kind = (value as { kind?: unknown }).kind;
	return kind === "text" || kind === "element" || kind === "component";
}
