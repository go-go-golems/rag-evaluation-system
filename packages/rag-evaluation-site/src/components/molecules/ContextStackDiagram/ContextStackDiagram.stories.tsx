import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	contextPaletteOptions,
	contextSignalOrangeStyleSet,
	contextStyleSetForPalette,
	type ContextPaletteName,
	type ContextStyleSet,
	type ContextVisualStyle,
	type ContextWindowSnapshot,
} from "../../../context";
import { Panel, Stack } from "../../layout";
import { ContextLegend } from "../ContextLegend";
import { ContextStackDiagram, type ContextStackDiagramProps } from "./ContextStackDiagram";

type StoryArgs = Omit<ContextStackDiagramProps, "styleSet"> & { palette: ContextPaletteName };

const meta = {
	title: "Component Library/Molecules/ContextStackDiagram",
	args: { palette: "Dusty Magenta / Blue" },
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		snapshot: { control: false },
		selectedPartId: { control: "text" },
	},
} satisfies Meta<StoryArgs>;
export default meta;
type Story = StoryObj<StoryArgs>;

const fallbackVisualStyle: ContextVisualStyle = {
	pattern: "overflow",
	fill: "var(--mac-surface)",
	line: "var(--mac-border)",
	stroke: "var(--mac-border)",
	labelColor: "var(--mac-text)",
};
function visual(base: ContextStyleSet, key: string) {
	return base.styles[key] ?? base.fallbackStyle ?? fallbackVisualStyle;
}
function compactStyleSet(
	base: ContextStyleSet,
	id: string,
	name: string,
	entries: Array<{ id: string; label: string; styleFrom: string; hidden?: boolean }>,
): ContextStyleSet {
	return {
		id,
		name,
		legendSize: "sm",
		swatchSize: "sm",
		fallbackStyle: base.fallbackStyle,
		legend: entries.map(({ id, label, hidden }, order) => ({ id, label, hidden, order })),
		styles: Object.fromEntries(entries.map(({ id, styleFrom }) => [id, visual(base, styleFrom)])),
	};
}
function renderStackPanel(
	title: string,
	snapshot: ContextWindowSnapshot,
	styleSet: ContextStyleSet,
	selectedPartId?: string,
) {
	return (
		<Panel title={`${title} · ${styleSet.name}`}>
			<Stack gap="sm">
				<ContextStackDiagram
					snapshot={snapshot}
					styleSet={styleSet}
					selectedPartId={selectedPartId}
				/>
				<ContextLegend
					items={styleSet.legend}
					styles={styleSet.styles}
					size="sm"
					selectedId={
						selectedPartId
							? snapshot.parts.find((part) => part.id === selectedPartId)?.styleKey
							: undefined
					}
				/>
			</Stack>
		</Panel>
	);
}
function triageStyleSet(base: ContextStyleSet) {
	return compactStyleSet(base, "stack-triage", "Triage lanes", [
		{ id: "mustRead", label: "Must read", styleFrom: "active" },
		{ id: "reference", label: "Reference", styleFrom: "retrieval" },
		{ id: "noise", label: "Noise", styleFrom: "evicted" },
		{ id: "room", label: "Room", styleFrom: "empty", hidden: true },
	]);
}
function timelineStyleSet(base: ContextStyleSet) {
	return compactStyleSet(base, "stack-timeline", "Turn timeline", [
		{ id: "setup", label: "Setup", styleFrom: "system" },
		{ id: "question", label: "Question", styleFrom: "conversation" },
		{ id: "work", label: "Work", styleFrom: "tool" },
		{ id: "result", label: "Result", styleFrom: "result" },
		{ id: "reply", label: "Reply", styleFrom: "active" },
		{ id: "room", label: "Unused", styleFrom: "empty" },
	]);
}
const triageSnapshot: ContextWindowSnapshot = {
	id: "stack-triage-snapshot",
	title: "Code-review triage stack",
	subtitle: "Only three visible legend lanes; free room is intentionally hidden.",
	limit: 40_000,
	selectedPartId: "files",
	parts: [
		{ id: "brief", label: "review brief", styleKey: "mustRead", tokens: 2_200 },
		{ id: "files", label: "changed files", styleKey: "mustRead", tokens: 8_600 },
		{ id: "docs", label: "linked docs", styleKey: "reference", tokens: 6_100 },
		{ id: "logs", label: "old logs", styleKey: "noise", tokens: 9_200 },
		{ id: "room", label: "room", styleKey: "room", tokens: 13_900 },
	],
};
const timelineSnapshot: ContextWindowSnapshot = {
	id: "stack-timeline-snapshot",
	title: "Turn timeline stack",
	subtitle: "Six explicit entries named for this product workflow.",
	limit: 60_000,
	selectedPartId: "turn-work",
	parts: [
		{ id: "turn-setup", label: "system setup", styleKey: "setup", tokens: 4_500 },
		{ id: "turn-question", label: "user question", styleKey: "question", tokens: 2_300 },
		{ id: "turn-work", label: "tool work", styleKey: "work", tokens: 18_500 },
		{ id: "turn-result", label: "parsed result", styleKey: "result", tokens: 9_800 },
		{ id: "turn-reply", label: "draft reply", styleKey: "reply", tokens: 3_100 },
		{ id: "turn-room", label: "unused", styleKey: "room", tokens: 21_800 },
	],
};

export const GroupedContextWindow: Story = {
	render: ({ palette }) =>
		renderStackPanel(
			"triage vocabulary",
			triageSnapshot,
			triageStyleSet(contextStyleSetForPalette(palette)),
			"files",
		),
};

export const SelectedLayer: Story = {
	args: { selectedPartId: "turn-work" },
	render: ({ palette, selectedPartId }) =>
		renderStackPanel(
			"timeline vocabulary",
			timelineSnapshot,
			timelineStyleSet(contextStyleSetForPalette(palette)),
			selectedPartId,
		),
};

export const Comparison: Story = {
	render: ({ palette }) => {
		const base = contextStyleSetForPalette(palette);
		return (
			<Stack gap="md">
				{renderStackPanel("3 visible labels", triageSnapshot, triageStyleSet(base), "files")}
				{renderStackPanel("6 labels", timelineSnapshot, timelineStyleSet(base), "turn-work")}
				{renderStackPanel(
					"signal orange reference",
					timelineSnapshot,
					timelineStyleSet(contextSignalOrangeStyleSet),
					"turn-work",
				)}
			</Stack>
		);
	},
};
