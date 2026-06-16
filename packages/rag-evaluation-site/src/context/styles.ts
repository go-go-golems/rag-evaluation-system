import type { CSSProperties } from "react";
import type {
	ContextLegendItemSpec,
	ContextPatternName,
	ContextStyleSet,
	ContextVisualStyle,
	ContextWindowPart,
} from "./types";

export interface PaletteColors {
	paper: string;
	ink: string;
	grid: string;
	shadow: string;
	accent_a: string;
	accent_b: string;
	accent_c: string;
}

export interface PaletteDefinition {
	name: string;
	mood: string;
	colors: PaletteColors;
	use: string;
}

export const dustyMagentaBlue: PaletteDefinition = {
	name: "Dusty Magenta / Blue",
	mood: "subtle synthwave, not neon",
	colors: {
		paper: "#F2EEF2",
		ink: "#141214",
		grid: "#C0B5C1",
		shadow: "#776C7A",
		accent_a: "#9C527E",
		accent_b: "#4F74A8",
		accent_c: "#D6B8CB",
	},
	use: "Use magenta for agent/persona tags; blue for machine state.",
};

export const signalOrangeCyan: PaletteDefinition = {
	name: "Signal Orange / Cyan",
	mood: "warm terminal + cool selection",
	colors: {
		paper: "#F5F1E8",
		ink: "#141414",
		grid: "#BDB7AA",
		shadow: "#7E7A72",
		accent_a: "#D86F2A",
		accent_b: "#2EA6A6",
		accent_c: "#F0B45A",
	},
	use: "Orange for alerts/active states; cyan for links, cursor, selection.",
};

export const slateCoral: PaletteDefinition = {
	name: "Slate / Coral",
	mood: "modern mac-like restraint",
	colors: {
		paper: "#F2F3EF",
		ink: "#111314",
		grid: "#B8BDBB",
		shadow: "#707878",
		accent_a: "#5F7F89",
		accent_b: "#C46A55",
		accent_c: "#D8E0DF",
	},
	use: "Slate as primary accent; coral only for interruption.",
};

export const cobaltSand: PaletteDefinition = {
	name: "Cobalt / Sand",
	mood: "technical atlas",
	colors: {
		paper: "#F3EEDC",
		ink: "#111318",
		grid: "#BFB59A",
		shadow: "#766F5E",
		accent_a: "#315D91",
		accent_b: "#D2A84A",
		accent_c: "#A8B8CC",
	},
	use: "Cobalt for active navigation; sand for passive emphasis.",
};

export const preferredContextPalettes = [
	signalOrangeCyan,
	slateCoral,
	cobaltSand,
	dustyMagentaBlue,
] as const;

function mix(color: string, pct: number, paper: string) {
	return `color-mix(in srgb, ${color} ${pct}%, ${paper})`;
}

export function contextVisualStyleToCssVars(style: ContextVisualStyle): CSSProperties {
	return {
		"--ctx-fill": style.fill,
		"--ctx-line": style.line ?? style.stroke ?? "var(--mac-border)",
		"--ctx-stroke": style.stroke ?? "var(--mac-border)",
		"--ctx-label": style.labelColor ?? "var(--mac-text)",
		"--ctx-stroke-width": `${style.strokeWidth ?? 1}px`,
		...(style.vars ?? {}),
	} as CSSProperties;
}

export function getPartStyleKey(part: ContextWindowPart): string {
	return part.styleKey;
}

export function resolveContextVisualStyle(
	styleKey: string,
	styleSet: ContextStyleSet,
): ContextVisualStyle {
	const style = styleSet.styles[styleKey];
	if (style) return style;
	if (styleSet.fallbackStyle) return styleSet.fallbackStyle;
	// eslint-disable-next-line no-console
	console.error(`Unknown context diagram styleKey "${styleKey}"`, { styleKey, styleSet });
	return {
		pattern: "overflow",
		fill: "var(--mac-surface)",
		line: "var(--mac-accent-2)",
		stroke: "var(--mac-accent-2)",
		labelColor: "var(--mac-text)",
	};
}

export function legendItemsForStyleSet(styleSet: ContextStyleSet): ContextLegendItemSpec[] {
	return [...styleSet.legend]
		.filter((item) => !item.hidden)
		.sort((a, b) => (a.order ?? 0) - (b.order ?? 0));
}

export interface StyleSetEntrySpec {
	id: string;
	label: string;
	description?: string;
	pattern?: ContextPatternName;
	accent?: "a" | "b" | "c" | "grid" | "shadow" | "ink";
	fillPct?: number;
	linePct?: number;
	solid?: boolean;
	hidden?: boolean;
}

function accentColor(palette: PaletteColors, accent: StyleSetEntrySpec["accent"] = "a") {
	switch (accent) {
		case "b":
			return palette.accent_b;
		case "c":
			return palette.accent_c;
		case "grid":
			return palette.grid;
		case "shadow":
			return palette.shadow;
		case "ink":
			return palette.ink;
		case "a":
		default:
			return palette.accent_a;
	}
}

export function createContextStyleSetFromPalette(options: {
	id?: string;
	name?: string;
	palette?: PaletteDefinition | PaletteColors;
	entries: StyleSetEntrySpec[];
	legendSize?: ContextStyleSet["legendSize"];
	swatchSize?: ContextStyleSet["swatchSize"];
}): ContextStyleSet {
	const colors =
		"colors" in (options.palette ?? {})
			? (options.palette as PaletteDefinition).colors
			: ((options.palette as PaletteColors | undefined) ?? dustyMagentaBlue.colors);
	const styles: Record<string, ContextVisualStyle> = {};
	const legend: ContextLegendItemSpec[] = [];
	for (const entry of options.entries) {
		const accent = accentColor(colors, entry.accent);
		const pattern = entry.pattern ?? (entry.solid ? "solid" : "diagonal");
		const solid = entry.solid || pattern === "solid";
		styles[entry.id] = {
			pattern,
			fill: solid ? accent : mix(accent, entry.fillPct ?? 18, colors.paper),
			line: solid ? "transparent" : mix(accent, entry.linePct ?? 60, colors.paper),
			stroke: colors.ink,
			labelColor: solid ? "#ffffff" : colors.ink,
		};
		legend.push({
			id: entry.id,
			label: entry.label,
			description: entry.description,
			hidden: entry.hidden,
		});
	}
	return {
		id: options.id,
		name: options.name,
		legend,
		styles,
		legendSize: options.legendSize ?? "md",
		swatchSize: options.swatchSize ?? "md",
		fallbackStyle: {
			pattern: "overflow",
			fill: colors.paper,
			line: colors.shadow,
			stroke: colors.ink,
			labelColor: colors.ink,
		},
	};
}

export function transcriptStyleSet(
	palette: PaletteDefinition | PaletteColors = dustyMagentaBlue,
): ContextStyleSet {
	return createContextStyleSetFromPalette({
		id: "transcript-tones",
		name: "Transcript Tones",
		palette,
		legendSize: "sm",
		swatchSize: "sm",
		entries: [
			{ id: "system", label: "System", accent: "b", pattern: "checker", fillPct: 16, linePct: 70 },
			{
				id: "developer",
				label: "Developer",
				accent: "b",
				pattern: "diagonal",
				fillPct: 14,
				linePct: 60,
			},
			{ id: "user", label: "User", accent: "a", pattern: "solid", solid: true },
			{
				id: "assistant",
				label: "Assistant",
				accent: "b",
				pattern: "diagonalDense",
				fillPct: 16,
				linePct: 65,
			},
			{ id: "tool", label: "Tool call", accent: "a", pattern: "stipple", fillPct: 16, linePct: 60 },
			{
				id: "result",
				label: "Tool result",
				accent: "a",
				pattern: "cross",
				fillPct: 14,
				linePct: 65,
			},
			{
				id: "context",
				label: "Context note",
				accent: "grid",
				pattern: "diagonal",
				fillPct: 10,
				linePct: 40,
			},
			{
				id: "generated",
				label: "Generated note",
				accent: "c",
				pattern: "diagonalDense",
				fillPct: 18,
				linePct: 70,
			},
			{ id: "active", label: "Active note", accent: "a", pattern: "solid", solid: true },
			{ id: "note", label: "Note", accent: "c", pattern: "checker", fillPct: 18, linePct: 80 },
			{ id: "header", label: "Header", accent: "grid", pattern: "none", hidden: true },
			{
				id: "rail",
				label: "Notes rail",
				accent: "grid",
				pattern: "diagonal",
				fillPct: 10,
				linePct: 35,
				hidden: true,
			},
			{
				id: "code",
				label: "Code output",
				accent: "shadow",
				pattern: "stipple",
				fillPct: 10,
				linePct: 45,
				hidden: true,
			},
			{ id: "other", label: "Other", accent: "grid", pattern: "none" },
		],
	});
}

export function defaultContextStyleSet(
	palette: PaletteDefinition | PaletteColors = dustyMagentaBlue,
): ContextStyleSet {
	return createContextStyleSetFromPalette({
		id: "context-window-default",
		name: "Context Window Default",
		palette,
		entries: [
			{
				id: "system",
				label: "system / instructions",
				accent: "b",
				pattern: "checker",
				fillPct: 20,
				linePct: 100,
			},
			{
				id: "instruction",
				label: "instruction",
				accent: "b",
				pattern: "checker",
				fillPct: 18,
				linePct: 55,
			},
			{ id: "context", label: "context / history", accent: "grid", pattern: "none" },
			{ id: "conversation", label: "conversation", accent: "grid", pattern: "none" },
			{
				id: "summary",
				label: "summary / memory",
				accent: "b",
				pattern: "diagonal",
				fillPct: 18,
				linePct: 55,
			},
			{
				id: "retrieval",
				label: "retrieval",
				accent: "a",
				pattern: "stipple",
				fillPct: 20,
				linePct: 60,
			},
			{ id: "tool", label: "tool call", accent: "a", pattern: "stipple", fillPct: 15, linePct: 50 },
			{
				id: "result",
				label: "tool result",
				accent: "a",
				pattern: "stipple",
				fillPct: 18,
				linePct: 55,
			},
			{
				id: "generated",
				label: "generated / scratchpad",
				accent: "c",
				pattern: "diagonalDense",
				fillPct: 25,
				linePct: 100,
			},
			{
				id: "annotation",
				label: "annotation",
				accent: "c",
				pattern: "diagonal",
				fillPct: 0,
				linePct: 40,
			},
			{ id: "course", label: "course material", accent: "grid", pattern: "none" },
			{
				id: "evicted",
				label: "evicted / compressed",
				accent: "b",
				pattern: "cross",
				fillPct: 15,
				linePct: 70,
			},
			{ id: "active", label: "active (current)", accent: "a", pattern: "solid", solid: true },
			{ id: "empty", label: "free space", accent: "grid", pattern: "none" },
			{ id: "other", label: "other", accent: "grid", pattern: "overflow" },
		],
	});
}
