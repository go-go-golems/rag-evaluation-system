import type { ContextStyleSet } from "../../context";
import type { FieldOption, FieldType, FieldValue } from "../../crm/types";
import type { ActionSpec } from "./actions";
import type { CellSpec, RowKeySpec } from "./cells";
import type { BaseWidgetProps, JsonObject, RenderableValue, WidgetNode } from "./core";

// ─── MatrixGrid (generic grid engine) ────────────────────────────────────

/**
 * Defunctionalized color function: value -> styleKey -> ContextVisualStyle.
 * The IR/DSL sibling of ActionSpec (event handler) and CellSpec (renderer).
 * Add once and MatrixGrid / MonthGrid / SegmentedBar / context diagrams all
 * become recolorable purely from serialized data.
 */
export interface StyleBySpec {
	/** Field to key on; defaults to the resolved cell value. */
	field?: string;
	styleSet: ContextStyleSet;
	/** Optional value -> styleKey remap before lookup. */
	map?: Record<string, string>;
	fallbackStyleKey?: string;
}

/** Renders a MatrixGrid cell as an n-state CycleCell (availability, RSVP, ...). */
export interface CycleCellSpec {
	kind: "cycle";
	states: string[];
	glyphs?: Record<string, RenderableValue>;
	/** Palette override; otherwise the grid's `styleSet` is used. */
	styleSet?: ContextStyleSet;
}

/** Renders the resolved (row,col) value as text; pair with `colorBy` for heatmaps. */
export interface ValueCellSpec {
	kind: "value";
}

export type MatrixCellSpec = CellSpec | CycleCellSpec | ValueCellSpec;

export interface MatrixColumnWidgetSpec {
	id: string;
	header: RenderableValue;
	meta?: JsonObject;
}

/** Accessor for the value at (row, column). */
export type MatrixValueSpec = { mapField: string } | { template: string };

export interface MatrixGridWidgetProps extends BaseWidgetProps {
	rows: JsonObject[];
	columns: MatrixColumnWidgetSpec[];
	/** Value at (row,col). Defaults to `row[col.id]`. */
	valueAt?: MatrixValueSpec;
	/** Mode A — one cell renderer applied per (row,col). */
	cell?: MatrixCellSpec;
	/** Mode B — an explicit matrix of prebuilt nodes. */
	cells?: WidgetNode[][];
	rowHeader?: CellSpec;
	styleSet?: ContextStyleSet;
	colorBy?: StyleBySpec;
	/** Footer row; the CellSpec is evaluated against each column's `meta`. */
	footer?: { header?: RenderableValue; cell: CellSpec };
	getRowKey?: RowKeySpec;
	editableRowKey?: string;
	selectedCell?: { rowKey: string; colId: string };
	cornerCell?: RenderableValue;
	stickyHeader?: boolean;
	ariaLabel?: string;
	onCellAction?: ActionSpec;
}

// ─── SegmentedBar / MonthGrid / TimeGrid (engine props) ──────────────────

export interface SegmentedBarSegmentSpec {
	value: number;
	styleKey: string;
	label?: RenderableValue;
}

export interface SegmentedBarMarkerSpec {
	at: number;
	styleKey?: string;
	label?: RenderableValue;
}

export interface SegmentedBarWidgetProps extends BaseWidgetProps {
	segments: SegmentedBarSegmentSpec[];
	styleSet: ContextStyleSet;
	total?: number;
	showCounts?: boolean;
	markers?: SegmentedBarMarkerSpec[];
	size?: "sm" | "md" | "lg";
	onSegmentAction?: ActionSpec;
}

export interface MonthGridMarkerSpec {
	count?: number;
	styleKey?: string;
	label?: RenderableValue;
}

export interface MonthGridWidgetProps extends BaseWidgetProps {
	monthISO: string;
	markers?: Record<string, MonthGridMarkerSpec>;
	styleSet?: ContextStyleSet;
	selectedDateISO?: string;
	todayISO?: string;
	minDateISO?: string;
	maxDateISO?: string;
	weekStartsOn?: 0 | 1;
	showHeader?: boolean;
	onDaySelectAction?: ActionSpec;
	onMonthChangeAction?: ActionSpec;
}

export interface TimeGridBlockSpec {
	id: string;
	dayISO: string;
	startISO: string;
	endISO: string;
	styleKey: string;
	label: RenderableValue;
	meta?: JsonObject;
}

export interface TimeGridColumnWidgetSpec {
	dayISO: string;
	header?: RenderableValue;
}

export interface TimeGridWidgetProps extends BaseWidgetProps {
	days: Array<string | TimeGridColumnWidgetSpec>;
	blocks: TimeGridBlockSpec[];
	styleSet: ContextStyleSet;
	hourStart?: number;
	hourEnd?: number;
	hourHeight?: number;
	nowISO?: string;
	selectedBlockId?: string;
	onBlockSelectAction?: ActionSpec;
	onSlotCreateAction?: ActionSpec;
}

// ─── CRM field system (FieldRenderer / RecordFieldList) ──────────────────

/**
 * Defunctionalized field rendering: a field's appearance is described by data
 * (a FieldSpec), and one FieldRenderer interpreter turns it into the right
 * control for its mode. The CRM analogue of CellSpec — add once and contacts,
 * companies, deals, and any custom object all render for free.
 */
export interface FieldSpec {
	/** which key in record.fields */
	key: string;
	type: FieldType;
	label?: RenderableValue;
	/** for select / multiselect */
	options?: FieldOption[];
	/** for relation / user */
	relatedObject?: string;
	readOnly?: boolean;
	/** e.g. "USD" for currency */
	unit?: string;
	/** colors for select/tag values (reuse the palette contract). */
	styleSet?: ContextStyleSet;
}

/** How a relation/user id resolves to a display label + avatar in read mode. */
export interface FieldRefSpec {
	label: string;
	avatarUrl?: string;
	href?: string;
}

export interface FieldRendererWidgetProps extends BaseWidgetProps {
	spec: FieldSpec;
	value: FieldValue;
	mode?: "read" | "edit";
	invalid?: boolean;
	/** id -> display, for relation/user read mode. */
	refs?: Record<string, FieldRefSpec>;
	/** edit mode reports value changes (e.g. a field.update server action). */
	onChangeAction?: ActionSpec;
}

export interface RecordFieldListSectionSpec {
	label?: RenderableValue;
	fields: FieldSpec[];
}

export interface RecordFieldListWidgetProps extends BaseWidgetProps {
	/** the record's field values, keyed by FieldSpec.key. */
	values: JsonObject;
	/** grouped sections; or use the flat `fields` shortcut. */
	sections?: RecordFieldListSectionSpec[];
	fields?: FieldSpec[];
	mode?: "read" | "edit";
	rowLayout?: "inline" | "stacked";
	invalidKeys?: string[];
	refs?: Record<string, FieldRefSpec>;
	onFieldChangeAction?: ActionSpec;
}

// ─── BoardEngine (kanban pipeline — the signature new CRM engine) ─────────

export interface BoardColumnWidgetSpec {
	id: string;
	header: RenderableValue;
	footer?: RenderableValue;
	/** palette styleKey for the column accent. */
	accent?: string;
}

/** How one card renders. v1: a DataTable-style CellSpec stack; extensible later. */
export interface BoardCardSpec {
	title: CellSpec;
	subtitle?: CellSpec;
	meta?: CellSpec;
	/** palette styleKey field to accent the card (e.g. status). */
	accentField?: string;
}

export interface BoardEngineWidgetProps extends BaseWidgetProps {
	columns: BoardColumnWidgetSpec[];
	cards: JsonObject[];
	/** which column a card is in — a field on the card. */
	columnField: string;
	getCardId?: RowKeySpec;
	card: BoardCardSpec;
	styleSet?: ContextStyleSet;
	selectedCardId?: string;
	ariaLabel?: string;
	/** dragging a card between columns. */
	onMoveAction?: ActionSpec;
	onCardSelectAction?: ActionSpec;
}

// ─── ActivityFeed (record timeline) ──────────────────────────────────────

export interface ActivityFeedItemSpec {
	id: string;
	kind: string;
	title: RenderableValue;
	body?: RenderableValue;
	atISO: string;
	actor?: { id?: string; name: string; avatarUrl?: string };
	meta?: JsonObject;
}

export interface ActivityFeedWidgetProps extends BaseWidgetProps {
	activities: ActivityFeedItemSpec[];
	styleSet?: ContextStyleSet;
	/** kind -> glyph for the spine icon. */
	glyphs?: Record<string, RenderableValue>;
	groupByDay?: boolean;
	onOpenAction?: ActionSpec;
	onLoadMoreAction?: ActionSpec;
}

// ─── StatTile (dashboard number) ─────────────────────────────────────────
// Several tiles are laid out with the existing DashboardGrid layout primitive —
// no bespoke row engine (follow the layout/molecule split already in the kit).

export interface StatTileWidgetProps extends BaseWidgetProps {
	label: RenderableValue;
	value: RenderableValue;
	delta?: number;
	deltaLabel?: RenderableValue;
	trend?: "up" | "down" | "flat";
	/** optional inline proportion 0..1 for a MeterBar track. */
	progress?: number;
	tone?: "accent" | "success" | "danger";
	onAction?: ActionSpec;
}
