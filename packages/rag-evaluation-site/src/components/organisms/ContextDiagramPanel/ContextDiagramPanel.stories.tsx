import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import {
	type ContextStyleSet,
	type ContextWindowSnapshot,
	contextCobaltSandStyleSet,
	contextDefaultStyleSet,
	contextSignalOrangeStyleSet,
	contextSlateCoralStyleSet,
	contextThreeLabelStyleSets,
	contextWindowSnapshots,
} from "../../../context";
import { Button } from "../../atoms";
import { Inline, Stack } from "../../layout";
import { ContextDiagramPanel, type ContextDiagramPanelProps } from "./ContextDiagramPanel";

const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;

const threeLabelSnapshot: ContextWindowSnapshot = {
	id: "three-label-rag",
	title: "Three-label RAG context",
	subtitle: "Caller-defined labels, palette size, and styles.",
	limit: 32_000,
	selectedPartId: "retrieved-docs",
	parts: [
		{
			id: "prompt",
			label: "Prompt",
			styleKey: "prompt",
			tokens: 1400,
			note: "System and task scaffolding.",
		},
		{
			id: "retrieved-docs",
			label: "Evidence",
			styleKey: "evidence",
			tokens: 9200,
			note: "Retrieved chunks and citations.",
		},
		{
			id: "answer-draft",
			label: "Draft",
			styleKey: "answer",
			tokens: 1800,
			note: "Current generated answer.",
		},
		{
			id: "free",
			label: "Free",
			styleKey: "free",
			tokens: 19600,
			note: "Remaining model context budget.",
		},
	],
};

const contentBlocks: ContextWindowSnapshot = {
	id: "content-blocks",
	title: "Turn 8 — actual context blocks",
	subtitle: "Each segment is a transcript or tool block, not an aggregate bucket.",
	limit: 32_000,
	selectedPartId: "turn-7-tool-result-bash",
	parts: [
		{
			id: "system",
			label: "system prompt",
			styleKey: "system",
			tokens: 950,
			note: "Runtime instructions and tool-use policy.",
			contentPreview:
				"You are an expert coding assistant. Prefer reading files before editing. Keep responses concise.",
			metadata: { source: "session", blockType: "system" },
		},
		{
			id: "turn-5-user",
			label: "T5 user",
			styleKey: "conversation",
			tokens: 210,
			note: "The user asks for a real context visualization instead of aggregate metrics.",
			contentPreview:
				"I want the context view visualization to be more about the actual content, seeing the different turns...",
			metadata: { turn: 5, role: "user" },
		},
		{
			id: "turn-6-assistant-plan",
			label: "T6 assistant plan",
			styleKey: "generated",
			tokens: 480,
			note: "Assistant planning text that may or may not be preserved by the upstream agent runtime.",
			contentPreview:
				"I’ll inspect the current ContextDiagram IR/component contract and our minitrace context model...",
			metadata: { turn: 6, role: "assistant", blockType: "plan" },
		},
		{
			id: "turn-7-tool-call-rg",
			label: "T7 tool call rg",
			styleKey: "tool",
			tokens: 120,
			note: "The command invocation itself.",
			contentPreview: '$ rg -n "ContextDiagramPanel|ContextWindowPart" packages pkg -S',
			metadata: { turn: 7, role: "tool", toolName: "bash", operation: "execute" },
		},
		{
			id: "turn-7-tool-result-bash",
			label: "T7 rg output",
			styleKey: "result",
			tokens: 1320,
			note: "Search output that explains which IR files and components need to change.",
			contentPreview:
				"packages/rag-evaluation-site/src/widgets/ir.ts:229: export interface ContextDiagramPanelWidgetProps...\npackages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/ContextDiagramPanel.tsx...",
			metadata: { turn: 7, role: "tool", toolName: "bash", fullBytes: 5280 },
		},
		{
			id: "turn-8-current",
			label: "T8 current answer",
			styleKey: "active",
			tokens: 390,
			note: "The currently generated response.",
			contentPreview:
				"Yes — your read is right. Current ContextDiagramPanel hard-codes the legend...",
			metadata: { turn: 8, role: "assistant" },
		},
		{
			id: "free",
			label: "free space",
			styleKey: "empty",
			tokens: 28_530,
			note: "Remaining model context budget.",
		},
	],
};

type PaletteName =
	| "Dusty Magenta / Blue"
	| "Signal Orange / Cyan"
	| "Slate / Coral"
	| "Cobalt / Sand";
type PaletteControlsArgs = Omit<ContextDiagramPanelProps, "styleSet"> & { palette: PaletteName };

const paletteControlStyleSets: Record<PaletteName, ContextStyleSet> = {
	"Dusty Magenta / Blue": contextDefaultStyleSet,
	"Signal Orange / Cyan": contextSignalOrangeStyleSet,
	"Slate / Coral": contextSlateCoralStyleSet,
	"Cobalt / Sand": contextCobaltSandStyleSet,
};

const meta = {
	title: "Component Library/Organisms/ContextDiagramPanel",
	args: {
		snapshot: deepBug!,
		palette: "Dusty Magenta / Blue",
		initialView: "strip",
		views: ["strip", "stack", "budget", "treemap"],
		showPartDetails: true,
	},
	argTypes: {
		palette: { control: "select", options: Object.keys(paletteControlStyleSets) },
		initialView: { control: "select", options: ["strip", "stack", "budget", "treemap"] },
		showLegend: { control: "boolean" },
		showPartDetails: { control: "boolean" },
		snapshot: { control: false },
		selectedPartId: { control: "text" },
		views: { control: false },
	},
} satisfies Meta<PaletteControlsArgs>;

export default meta;
type Story = StoryObj<PaletteControlsArgs>;

export const PaletteControls: Story = {
	args: {
		snapshot: deepBug!,
		palette: "Dusty Magenta / Blue",
		initialView: "strip",
		views: ["strip", "stack", "budget", "treemap"],
		showPartDetails: true,
	},
	render: ({ palette, ...args }) => (
		<ContextDiagramPanel {...args} styleSet={paletteControlStyleSets[palette]} />
	),
};

export const InteractiveViews: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel snapshot={deepBug!} styleSet={paletteControlStyleSets[palette]} />
	),
};

export const StartingViews: Story = {
	render: ({ palette }) => {
		const styleSet = paletteControlStyleSets[palette];
		return (
			<Stack gap="md">
				<ContextDiagramPanel snapshot={deepBug!} styleSet={styleSet} initialView="strip" />
				<ContextDiagramPanel snapshot={atLimit!} styleSet={styleSet} initialView="treemap" />
				<ContextDiagramPanel snapshot={overBudget!} styleSet={styleSet} initialView="budget" />
			</Stack>
		);
	},
};

export const ContentBlocksWithPartDetails: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={contentBlocks}
			styleSet={paletteControlStyleSets[palette]}
			initialView="stack"
			views={["stack", "strip", "budget"]}
			showPartDetails
		/>
	),
};

export const LegendDerivedFromSnapshotParts: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={contentBlocks}
			styleSet={paletteControlStyleSets[palette]}
			initialView="strip"
			showPartDetails
		/>
	),
};

export const CustomThreeLabelStyleSets: Story = {
	render: () => (
		<Stack gap="md">
			{contextThreeLabelStyleSets.map((styleSet) => (
				<ContextDiagramPanel
					key={styleSet.id}
					snapshot={threeLabelSnapshot}
					styleSet={styleSet}
					initialView="strip"
					views={["strip", "budget", "treemap"]}
					showPartDetails
				/>
			))}
		</Stack>
	),
};

const switchableStyleSets: ContextStyleSet[] = [
	contextDefaultStyleSet,
	contextSignalOrangeStyleSet,
	contextSlateCoralStyleSet,
	contextCobaltSandStyleSet,
	...contextThreeLabelStyleSets,
];

export const InteractiveStyleSwitcher: Story = {
	render: ({ palette }) => {
		const [styleSet, setStyleSet] = useState<ContextStyleSet>(paletteControlStyleSets[palette]);
		const usesThreeLabel = styleSet.id?.startsWith("three-label");
		return (
			<Stack gap="md">
				<Inline gap="xs" wrap>
					{switchableStyleSets.map((candidate) => (
						<Button
							key={candidate.id ?? candidate.name}
							size="compact"
							selected={candidate.id === styleSet.id}
							onClick={() => setStyleSet(candidate)}
						>
							{candidate.name ?? candidate.id}
						</Button>
					))}
				</Inline>
				<ContextDiagramPanel
					snapshot={usesThreeLabel ? threeLabelSnapshot : deepBug!}
					styleSet={styleSet}
					initialView="strip"
					views={["strip", "stack", "budget", "treemap"]}
					showPartDetails
				/>
			</Stack>
		);
	},
};

export const InlineEmbeddedHandoutNoTreemap: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={contentBlocks}
			styleSet={paletteControlStyleSets[palette]}
			initialView="budget"
			views={["budget", "strip", "stack"]}
			chrome="inline"
			showPartDetails={false}
		/>
	),
};

export const InlineEmbeddedHandoutWithDetails: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={contentBlocks}
			styleSet={paletteControlStyleSets[palette]}
			initialView="stack"
			views={["stack", "strip", "budget"]}
			chrome="inline"
			showPartDetails
		/>
	),
};

export const FilteredLegendBudgetOmitsHeadroom: Story = {
	render: () => (
		<ContextDiagramPanel
			snapshot={threeLabelSnapshot}
			styleSet={contextThreeLabelStyleSets[0]!}
			initialView="budget"
			views={["budget"]}
			chrome="inline"
		/>
	),
};

export const FilteredLegendOnlyVisibleStyleKeys: Story = {
	render: ({ palette }) => (
		<Stack gap="md">
			<ContextDiagramPanel
				snapshot={contextWindowSnapshots[0]!}
				styleSet={paletteControlStyleSets[palette]}
				initialView="budget"
				views={["budget"]}
				chrome="inline"
			/>
			<ContextDiagramPanel
				snapshot={deepBug!}
				styleSet={paletteControlStyleSets[palette]}
				initialView="strip"
				views={["strip"]}
				chrome="inline"
			/>
			<ContextDiagramPanel
				snapshot={overBudget!}
				styleSet={paletteControlStyleSets[palette]}
				initialView="stack"
				views={["stack"]}
				chrome="inline"
			/>
		</Stack>
	),
};

export const PanelChromeStillAvailable: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={deepBug!}
			styleSet={paletteControlStyleSets[palette]}
			initialView="strip"
			views={["strip", "stack", "budget", "treemap"]}
			chrome="panel"
			showPartDetails
		/>
	),
};

export const InlineWithoutLegendForDensePrint: Story = {
	render: ({ palette }) => (
		<ContextDiagramPanel
			snapshot={atLimit!}
			styleSet={paletteControlStyleSets[palette]}
			initialView="budget"
			views={["budget", "strip", "stack"]}
			chrome="inline"
			showLegend={false}
			showPartDetails={false}
		/>
	),
};
