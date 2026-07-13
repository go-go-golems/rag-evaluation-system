import type { ActionSpec } from "./actions";
import type { WidgetProps } from "./props";

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
	| "MatrixGrid"
	| "SegmentedBar"
	| "MonthGrid"
	| "TimeGrid"
	| "FieldRenderer"
	| "RecordFieldList"
	| "BoardEngine"
	| "ActivityFeed"
	| "StatTile"
	| "Divider"
	| "Disclosure"
	| "FormPanel"
	| "FormRow"
	| "Inline"
	| "MetadataGrid"
	| "Panel"
	| "ScrollRegion"
	| "SectionBlock"
	| "FieldGrid"
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
	| "ShareLink"
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
	| "ContextUploadDropArea"
	| "MediaThumb"
	| "Tag"
	| "ContentStatusBadge"
	| "MeterBar"
	| "TileGrid"
	| "AssetTile"
	| "Breadcrumbs"
	| "Pagination"
	| "SearchField"
	| "EmptyState"
	| "MarkdownEditor"
	| "MediaLibraryPanel"
	| "ArticleListPanel"
	| "CmsShell";

export interface ComponentNode {
	kind: "component";
	type: RagWidgetType | string;
	props?: WidgetProps;
	children?: WidgetNode[];
}

export type RenderableValue = WidgetNode | string | number | boolean | null;

export interface BaseWidgetProps {
	className?: string;
	style?: JsonObject;
	id?: string;
	action?: ActionSpec;
	[key: string]: unknown;
}

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
