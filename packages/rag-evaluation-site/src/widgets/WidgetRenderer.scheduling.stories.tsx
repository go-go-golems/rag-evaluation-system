import type { Meta, StoryObj } from "@storybook/react-vite";
import type { ContextStyleSet } from "../context";
import { sampleTeamSyncPoll, sampleTeamSyncTallies } from "../scheduling";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type WidgetNode } from "./ir";
import { availabilityMatrix, pollResults } from "./presets/scheduling";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/Scheduling",
	component: WidgetRenderer,
	args: {
		registry: defaultWidgetRegistry,
		// Intercept server actions so cell clicks log instead of hitting the network.
		onAction: (action, context) => {
			// eslint-disable-next-line no-console
			console.log("[widget action]", action, context);
		},
	},
} satisfies Meta<typeof WidgetRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

/**
 * The Doodle poll rendered entirely from a serialized IR node tree produced by
 * the `availabilityMatrix` preset — the same MatrixGrid you see in
 * `Molecules/MatrixGrid`, but driven through the registry rather than React.
 */
export const AvailabilityPoll: Story = {
	args: {
		node: availabilityMatrix(sampleTeamSyncPoll, {
			tallies: sampleTeamSyncTallies,
			editableResponseId: "you",
		}),
	},
};

/** The preset composed inside a Panel — proves cross-widget IR composition. */
export const PollInPanel: Story = {
	args: {
		node: component("Panel", { title: sampleTeamSyncPoll.title, density: "condensed" }, [
			component("Stack", { gap: "sm" }, [
				component("KeyValueStrip", {
					items: [
						{ key: text("Location"), value: text(sampleTeamSyncPoll.location ?? "—") },
						{ key: text("Organizer"), value: text(sampleTeamSyncPoll.organizer.name) },
					],
				}),
				availabilityMatrix(sampleTeamSyncPoll, {
					tallies: sampleTeamSyncTallies,
					editableResponseId: "you",
				}),
			]),
		]) as WidgetNode,
	},
};

/** Organizer results: SegmentedBar per option via the `pollResults` preset. */
export const PollResults: Story = {
	args: {
		node: component("Panel", { title: "Results · Team sync", density: "condensed" }, [
			pollResults(sampleTeamSyncPoll, sampleTeamSyncTallies),
		]) as WidgetNode,
	},
};

/**
 * StyleBySpec in action: a rating heatmap. Value cells rendered from IR, tinted
 * by `colorBy` (value -> styleKey -> ContextVisualStyle) — no per-cell colors.
 */
const heatStyleSet: ContextStyleSet = {
	id: "rating",
	styles: {
		"1": {
			fill: "color-mix(in srgb, var(--mac-accent-2) 55%, var(--mac-surface))",
			labelColor: "var(--mac-text)",
		},
		"2": {
			fill: "color-mix(in srgb, var(--mac-amber) 45%, var(--mac-surface))",
			labelColor: "var(--mac-text)",
		},
		"3": {
			fill: "color-mix(in srgb, var(--mac-amber) 25%, var(--mac-surface))",
			labelColor: "var(--mac-text)",
		},
		"4": {
			fill: "color-mix(in srgb, var(--mac-green) 30%, var(--mac-surface))",
			labelColor: "var(--mac-text)",
		},
		"5": {
			fill: "color-mix(in srgb, var(--mac-green) 55%, var(--mac-surface))",
			labelColor: "var(--mac-text)",
		},
	},
	legend: [],
};

export const ColorByHeatmap: Story = {
	args: {
		node: component("MatrixGrid", {
			ariaLabel: "Skill matrix",
			cornerCell: text("Person"),
			rows: [
				{ id: "a", name: "Alice", go: 5, ts: 4, sql: 2 },
				{ id: "b", name: "Bob", go: 3, ts: 5, sql: 4 },
				{ id: "c", name: "Chen", go: 2, ts: 3, sql: 5 },
			],
			columns: [
				{ id: "go", header: text("Go") },
				{ id: "ts", header: text("TS") },
				{ id: "sql", header: text("SQL") },
			],
			rowHeader: { kind: "field", field: "name" },
			cell: { kind: "value" },
			colorBy: { styleSet: heatStyleSet },
		}),
	},
};

/**
 * Hand-authored IR (no preset) showing the raw node shape a DSL author emits,
 * using Mode B (an explicit `cells` node matrix) for a domain-blind grid.
 */
export const HandAuthoredMatrix: Story = {
	args: {
		node: component("MatrixGrid", {
			ariaLabel: "Plan comparison",
			cornerCell: text("Feature"),
			rows: [
				{ id: "sso", name: "SSO" },
				{ id: "api", name: "API" },
				{ id: "sup", name: "Support" },
			],
			columns: [
				{ id: "free", header: text("Free") },
				{ id: "pro", header: text("Pro") },
				{ id: "team", header: text("Team") },
			],
			rowHeader: { kind: "field", field: "name" },
			cells: [
				[text(""), text("✓"), text("✓")],
				[text(""), text("✓"), text("✓")],
				[text(""), text(""), text("✓")],
			],
		}),
	},
};
