export type TranscriptRole = "system" | "developer" | "user" | "assistant" | "tool" | "other";

export type ContextJsonPrimitive = string | number | boolean | null;
export type ContextJsonValue =
	| ContextJsonPrimitive
	| ContextJsonValue[]
	| { [key: string]: ContextJsonValue };
export type ContextJsonRecord = Record<string, ContextJsonValue>;

export type ContextPatternName =
	| "none"
	| "solid"
	| "checker"
	| "diagonal"
	| "diagonalWide"
	| "diagonalDense"
	| "stipple"
	| "cross"
	| "overflow";

export type ContextStyleSize = "xs" | "sm" | "md" | "lg";

export interface ContextVisualStyle {
	fill: string;
	line?: string;
	stroke?: string;
	labelColor?: string;
	pattern?: ContextPatternName;
	dashed?: boolean;
	dotted?: boolean;
	strokeWidth?: number;
	vars?: Record<string, string>;
}

export interface ContextLegendItemSpec {
	id: string;
	label: string;
	description?: string;
	styleKey?: string;
	order?: number;
	hidden?: boolean;
}

export interface ContextStyleSet {
	id?: string;
	name?: string;
	styles: Record<string, ContextVisualStyle>;
	legend: ContextLegendItemSpec[];
	fallbackStyle?: ContextVisualStyle;
	legendSize?: ContextStyleSize;
	swatchSize?: ContextStyleSize;
}

export interface ContextWindowPart {
	id: string;
	label: string;
	styleKey: string;
	tokens: number;
	note?: string;
	contentPreview?: string;
	content?: string;
	sourceId?: string;
	startToken?: number;
	endToken?: number;
	metadata?: ContextJsonRecord;
}

export interface ContextWindowSnapshot {
	id: string;
	title: string;
	subtitle?: string;
	limit: number;
	parts: ContextWindowPart[];
	selectedPartId?: string;
	metadata?: ContextJsonRecord;
}

export type ContextDiagramView = "strip" | "stack" | "budget" | "treemap";

export interface ContextDiagramSegment {
	id: string;
	label: string;
	styleKey: string;
	tokens: number;
	note?: string;
	metadata?: ContextJsonRecord;
}

export interface TranscriptAnnotation {
	id: string;
	targetMessageId: string;
	styleKey: string;
	label: string;
	text: string;
	confidence?: number;
	metadata?: ContextJsonRecord;
}

export interface TranscriptMessage {
	id: string;
	role: TranscriptRole;
	text: string;
	tokens?: number;
	name?: string;
	timestamp?: string;
	annotationIds?: string[];
	metadata?: ContextJsonRecord;
}

export interface TranscriptFixture {
	id: string;
	title: string;
	subtitle?: string;
	messages: TranscriptMessage[];
	annotations: TranscriptAnnotation[];
	selectedAnnotationId?: string;
}

export interface ContextCourseInstructor {
	name: string;
	role?: string;
	bio?: string;
}

export interface ContextCourseAgendaItem {
	id: string;
	number: string;
	title: string;
	description: string;
	duration?: string;
}

export interface ContextCourseLabels {
	outcomes?: string;
	agenda?: string;
	instructor?: string;
	when?: string;
	where?: string;
	format?: string;
	primaryCta?: string;
	secondaryCta?: string;
}

export interface ContextCourse {
	id: string;
	kicker?: string;
	title: string;
	tagline: string;
	when?: string;
	where?: string;
	format?: string;
	price?: string;
	instructor?: ContextCourseInstructor;
	blurb?: string;
	labels?: ContextCourseLabels;
	outcomes: string[];
	agenda: ContextCourseAgendaItem[];
}

export interface ContextSlide {
	id: string;
	number: string;
	kicker?: string;
	title: string;
	view: ContextDiagramView;
	snapshotId: string;
	notes: string[];
}

export interface ArticleMarkdownBlock {
	kind: "markdown";
	id: string;
	source: string;
}

export interface ArticleContextWindowBlock {
	kind: "context-window";
	id: string;
	snapshot: ContextWindowSnapshot;
	view?: ContextDiagramView;
	caption?: string;
}

export interface ArticleImageBlock {
	kind: "image";
	id: string;
	src: string;
	alt: string;
	caption?: string;
}

export type ArticleBlock = ArticleMarkdownBlock | ArticleContextWindowBlock | ArticleImageBlock;

export interface ContextHandoutDocument {
	id: string;
	title: string;
	file: string;
	format: string;
	size?: string;
	description: string;
	body: string;
	blocks?: ArticleBlock[];
	downloadHref?: string;
	printHref?: string;
}

export type AnchoredCommentStatus = "open" | "resolved";

export interface AnchoredComment {
	id: string;
	anchorX: number;
	anchorY: number;
	author: string;
	text: string;
	time?: string;
	status?: AnchoredCommentStatus;
	metadata?: ContextJsonRecord;
}

export interface ContextHandoutBundle {
	intro: string;
	docs: ContextHandoutDocument[];
}
