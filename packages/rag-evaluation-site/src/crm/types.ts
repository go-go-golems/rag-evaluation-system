/**
 * CRM domain DTOs. Pure data — no React, no Widget IR, no server calls.
 * This is the CRM counterpart to `src/scheduling/types.ts`: engines, adapters,
 * presets, panels, and stories all import their nouns from here so that many
 * widgets share one definition of what a Contact / Deal / Activity is.
 *
 * Two modeling choices drive every widget downstream:
 *   1. Records carry a `fields` bag *in addition to* a few first-class columns.
 *      A CRM's whole selling point is customer-defined fields, so the field set
 *      is data (`FieldDef`), not code.
 *   2. Everything that happens to a record is an `Activity` with a `kind`. One
 *      stream, many kinds — that uniformity is what lets a single timeline
 *      engine render the entire history of a record.
 */

export type Id = string;

/**
 * The value stored under a key in `record.fields`. Deliberately broad: a CRM
 * field can be a scalar, an id (relation), a list of ids (multiselect/relation),
 * or a small structured object (address). `null` means "empty".
 */
export type FieldValue =
	| string
	| number
	| boolean
	| Id
	| Id[]
	| string[]
	| null
	| { [key: string]: string | number | boolean | null };

// ── Field schema (the heart of the CRM — see design-doc Part 4) ──────────────

export type FieldType =
	| "text"
	| "longtext"
	| "email"
	| "phone"
	| "url"
	| "number"
	| "currency"
	| "percent"
	| "date"
	| "datetime"
	| "boolean"
	| "select"
	| "multiselect"
	| "tags"
	| "relation"
	| "user"
	| "address";

export type RelatedObject = "contact" | "company" | "deal" | "user";

export interface FieldOption {
	value: string;
	label: string;
	/** Palette lookup key (never a raw color) for the option's pill. */
	colorKey?: string;
}

/**
 * A field *definition* — the schema a workspace configures. Records carry
 * values keyed by `FieldDef.key`; the definitions live alongside the record and
 * tell the renderer how to display and edit each value.
 */
export interface FieldDef {
	key: string;
	label: string;
	type: FieldType;
	/** for select / multiselect */
	options?: FieldOption[];
	/** for relation / user */
	relatedObject?: RelatedObject;
	required?: boolean;
	readOnly?: boolean;
	/** optional grouping bucket used by RecordFieldList sections. */
	group?: string;
	/** e.g. "USD" for currency; a placeholder hint; etc. */
	unit?: string;
}

// ── Records ──────────────────────────────────────────────────────────────────

export interface Contact {
	id: Id;
	name: string;
	avatarUrl?: string;
	title?: string;
	companyId?: Id;
	/** email, phone, custom fields — keyed by FieldDef.key. */
	fields: Record<string, FieldValue>;
	ownerId?: Id;
	tags?: string[];
	updatedAtISO: string;
}

export interface Company {
	id: Id;
	name: string;
	domain?: string;
	logoUrl?: string;
	fields: Record<string, FieldValue>;
	ownerId?: Id;
	tags?: string[];
	updatedAtISO?: string;
}

export type DealStatus = "open" | "won" | "lost";

export interface Deal {
	id: Id;
	title: string;
	amount?: number;
	currency?: string;
	/** which pipeline column it sits in. */
	stageId: Id;
	pipelineId: Id;
	contactIds?: Id[];
	companyId?: Id;
	ownerId?: Id;
	closeDateISO?: string;
	fields: Record<string, FieldValue>;
	status: DealStatus;
}

export interface Stage {
	id: Id;
	name: string;
	order: number;
	/** Palette lookup key (never a raw color). */
	colorKey: string;
	probability?: number;
}

export interface Pipeline {
	id: Id;
	name: string;
	stages: Stage[];
}

// ── Activity stream (one timeline, many kinds) ───────────────────────────────

export type ActivityKind =
	| "note"
	| "email"
	| "call"
	| "meeting"
	| "task"
	| "stage_change"
	| "field_change";

export interface Activity {
	id: Id;
	kind: ActivityKind;
	actor: { id: Id; name: string; avatarUrl?: string };
	atISO: string;
	/** the record this activity is on. */
	subjectId: Id;
	title: string;
	body?: string;
	/** e.g. { from: "Lead", to: "Qualified" } for stage_change. */
	meta?: Record<string, unknown>;
}

// ── Tasks ────────────────────────────────────────────────────────────────────

export type TaskStatus = "open" | "done";
export type TaskPriority = "low" | "med" | "high";

export interface Task {
	id: Id;
	title: string;
	dueISO?: string;
	status: TaskStatus;
	assigneeId?: Id;
	relatedId?: Id;
	priority?: TaskPriority;
}

// ── User (owner / assignee / actor directory) ────────────────────────────────

export interface CrmUser {
	id: Id;
	name: string;
	avatarUrl?: string;
	email?: string;
}
