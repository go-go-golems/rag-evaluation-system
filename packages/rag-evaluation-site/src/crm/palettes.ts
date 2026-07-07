import type { ContextStyleSet } from "../context";
import type { ActivityKind } from "./types";

/**
 * Canonical CRM palettes, typed as the shared `ContextStyleSet` contract — the
 * same coloring mechanism used by the context diagrams and scheduling widgets.
 * One place to key stage colors, activity-kind glyphs, and tag tints so a deal
 * card, a pipeline column, and a timeline row stay visually consistent.
 */

/** Pipeline stage colors, keyed by `Stage.colorKey`. */
export const stageStyleSet: ContextStyleSet = {
	id: "crm-stages",
	styles: {
		lead: {
			fill: "color-mix(in srgb, var(--mac-text-dim) 20%, var(--mac-surface))",
			line: "var(--mac-text-dim)",
			labelColor: "var(--mac-text)",
		},
		qualified: {
			fill: "color-mix(in srgb, var(--mac-accent) 22%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		proposal: {
			fill: "color-mix(in srgb, var(--mac-amber) 26%, var(--mac-surface))",
			line: "var(--mac-amber)",
			labelColor: "var(--mac-text)",
		},
		negotiation: {
			fill: "color-mix(in srgb, var(--mac-accent-2) 24%, var(--mac-surface))",
			line: "var(--mac-accent-2)",
			labelColor: "var(--mac-text)",
		},
		won: {
			fill: "color-mix(in srgb, var(--mac-green) 26%, var(--mac-surface))",
			line: "var(--mac-green)",
			labelColor: "var(--mac-text)",
		},
		lost: {
			fill: "color-mix(in srgb, var(--mac-accent-2) 30%, var(--mac-surface))",
			line: "var(--mac-accent-2)",
			labelColor: "var(--mac-text)",
		},
	},
	legend: [
		{ id: "lead", label: "Lead", styleKey: "lead" },
		{ id: "qualified", label: "Qualified", styleKey: "qualified" },
		{ id: "proposal", label: "Proposal", styleKey: "proposal" },
		{ id: "negotiation", label: "Negotiation", styleKey: "negotiation" },
		{ id: "won", label: "Won", styleKey: "won" },
	],
};

/** Activity-kind colors, keyed by `Activity.kind`. */
export const activityStyleSet: ContextStyleSet = {
	id: "crm-activities",
	styles: {
		note: {
			fill: "color-mix(in srgb, var(--mac-amber) 22%, var(--mac-surface))",
			line: "var(--mac-amber)",
			labelColor: "var(--mac-text)",
		},
		email: {
			fill: "color-mix(in srgb, var(--mac-accent) 20%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		call: {
			fill: "color-mix(in srgb, var(--mac-green) 22%, var(--mac-surface))",
			line: "var(--mac-green)",
			labelColor: "var(--mac-text)",
		},
		meeting: {
			fill: "color-mix(in srgb, var(--mac-accent) 26%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		task: {
			fill: "color-mix(in srgb, var(--mac-text-dim) 18%, var(--mac-surface))",
			line: "var(--mac-text-dim)",
			labelColor: "var(--mac-text)",
		},
		stage_change: {
			fill: "color-mix(in srgb, var(--mac-accent-2) 20%, var(--mac-surface))",
			line: "var(--mac-accent-2)",
			labelColor: "var(--mac-text)",
		},
		field_change: {
			fill: "var(--mac-surface-2)",
			line: "var(--mac-border)",
			labelColor: "var(--mac-text-dim)",
		},
	},
	legend: [
		{ id: "note", label: "Note", styleKey: "note" },
		{ id: "email", label: "Email", styleKey: "email" },
		{ id: "call", label: "Call", styleKey: "call" },
		{ id: "meeting", label: "Meeting", styleKey: "meeting" },
	],
};

/** Tag tints — a small cycling palette keyed by tag slug is overkill; one soft chip. */
export const tagStyleSet: ContextStyleSet = {
	id: "crm-tags",
	styles: {
		default: {
			fill: "var(--mac-surface-2)",
			line: "var(--mac-border)",
			labelColor: "var(--mac-text)",
		},
		enterprise: {
			fill: "color-mix(in srgb, var(--mac-accent) 18%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		"mid-market": {
			fill: "color-mix(in srgb, var(--mac-green) 18%, var(--mac-surface))",
			line: "var(--mac-green)",
			labelColor: "var(--mac-text)",
		},
		churn_risk: {
			fill: "color-mix(in srgb, var(--mac-accent-2) 22%, var(--mac-surface))",
			line: "var(--mac-accent-2)",
			labelColor: "var(--mac-text)",
		},
	},
	legend: [],
};

/** One glyph per activity kind, for the timeline spine and card icons. */
export const ACTIVITY_GLYPHS: Record<ActivityKind, string> = {
	note: "📝",
	email: "✉",
	call: "☎",
	meeting: "📅",
	task: "☑",
	stage_change: "▧",
	field_change: "✎",
};
