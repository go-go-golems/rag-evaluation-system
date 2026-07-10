import type { ContextStyleSet } from "../context";
import type { AvailabilityState } from "./types";

/**
 * Canonical scheduling palettes, typed as the shared `ContextStyleSet` contract.
 * These used to be duplicated inside individual stories; they live here once so
 * a poll cell and a calendar event stay visually consistent everywhere.
 */

export const availabilityStyleSet: ContextStyleSet = {
	id: "availability",
	styles: {
		yes: { fill: "var(--mac-green)", labelColor: "var(--mac-text-inv)" },
		ifneedbe: { fill: "var(--mac-amber)", labelColor: "var(--mac-text)" },
		no: { fill: "var(--mac-accent-2)", labelColor: "var(--mac-text-inv)" },
		unknown: { fill: "var(--mac-surface)", labelColor: "var(--mac-text-dim)" },
	},
	legend: [
		{ id: "yes", label: "Yes", styleKey: "yes" },
		{ id: "ifneedbe", label: "If need be", styleKey: "ifneedbe" },
		{ id: "no", label: "No", styleKey: "no" },
	],
};

export const AVAILABILITY_STATES: AvailabilityState[] = ["yes", "ifneedbe", "no", "unknown"];

export const AVAILABILITY_GLYPHS: Record<AvailabilityState, string> = {
	yes: "✓",
	ifneedbe: "~",
	no: "✕",
	unknown: "·",
};

export const eventStyleSet: ContextStyleSet = {
	id: "events",
	styles: {
		meeting: {
			fill: "color-mix(in srgb, var(--mac-accent) 22%, var(--mac-surface))",
			line: "var(--mac-accent)",
			labelColor: "var(--mac-text)",
		},
		focus: {
			fill: "color-mix(in srgb, var(--mac-green) 22%, var(--mac-surface))",
			line: "var(--mac-green)",
			labelColor: "var(--mac-text)",
		},
		personal: {
			fill: "color-mix(in srgb, var(--mac-amber) 26%, var(--mac-surface))",
			line: "var(--mac-amber)",
			labelColor: "var(--mac-text)",
		},
	},
	legend: [
		{ id: "meeting", label: "Meeting", styleKey: "meeting" },
		{ id: "focus", label: "Focus", styleKey: "focus" },
		{ id: "personal", label: "Personal", styleKey: "personal" },
	],
};
