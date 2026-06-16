import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	contextPaletteOptions,
	contextStyleSetForPalette,
	contextThreeLabelStyleSets,
	contextWindowSnapshots,
	type ContextPaletteName,
	type ContextStyleSet,
	type ContextVisualStyle,
	type ContextWindowSnapshot,
} from "../context";
import { WidgetRenderer, type WidgetRendererProps } from "./WidgetRenderer";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type WidgetNode } from "./ir";

type StoryArgs = WidgetRendererProps & { palette: ContextPaletteName };

const meta = {
	title: "Widget IR/Renderer/Context Diagrams",
	args: { registry: defaultWidgetRegistry, palette: "Dusty Magenta / Blue" },
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		node: { control: false },
		registry: { control: false },
		onAction: { control: false },
	},
} satisfies Meta<StoryArgs>;
export default meta;
type Story = StoryObj<StoryArgs>;

const snapshot = contextWindowSnapshots[0]!;
const contentSnapshot: ContextWindowSnapshot = {
	id: "widget-ir-content-blocks",
	title: "Widget IR content blocks",
	subtitle: "Server-provided parts correspond to transcript/tool blocks.",
	limit: 16_000,
	selectedPartId: "tool-result",
	parts: [
		{
			id: "system",
			label: "system",
			styleKey: "system",
			tokens: 720,
			note: "Instructions and tool policy.",
			contentPreview: "You are an expert coding assistant...",
		},
		{
			id: "user-turn",
			label: "T4 user",
			styleKey: "conversation",
			tokens: 180,
			note: "User asks for content-level context visualization.",
			contentPreview:
				"I want the context view visualization to be more about the actual content...",
		},
		{
			id: "tool-call",
			label: "T5 bash call",
			styleKey: "tool",
			tokens: 90,
			note: "Search command.",
			contentPreview: "$ rg -n ContextDiagramPanel packages pkg -S",
			metadata: { turn: 5, toolName: "bash" },
		},
		{
			id: "tool-result",
			label: "T5 search output",
			styleKey: "result",
			tokens: 1180,
			note: "Search results showing the files to change.",
			contentPreview:
				"packages/rag-evaluation-site/src/widgets/ir.ts:229...\npackages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel...",
			metadata: { turn: 5, fullBytes: 4720 },
		},
		{
			id: "answer",
			label: "T6 assistant",
			styleKey: "active",
			tokens: 340,
			note: "Current answer draft.",
			contentPreview: "Yes — the current panel hard-codes a generic legend...",
		},
		{ id: "free", label: "free", styleKey: "empty", tokens: 13_490 },
	],
};
const overBudgetSnapshot: ContextWindowSnapshot = {
	...contextWindowSnapshots[1]!,
	id: "over-budget-widget-ir",
	title: "Over budget context window",
	limit: 12_000,
};
const groupedSessionSnapshot: ContextWindowSnapshot = {
	id: "widget-ir-grouped-session",
	title: "Uploaded session context by turn",
	subtitle: "Balanced fixture with visible turn boundaries.",
	limit: 12_000,
	selectedPartId: "tool-result",
	parts: [
		{
			id: "system",
			label: "system",
			styleKey: "system",
			tokens: 700,
			note: "Stable system/developer instructions.",
			metadata: { turnIndex: "global", blockType: "system" },
		},
		{
			id: "t4-user",
			label: "T4 user",
			styleKey: "conversation",
			tokens: 1800,
			note: "User request enters the window.",
			metadata: { turnIndex: 4 },
		},
		{
			id: "t4-context",
			label: "T4 context",
			styleKey: "context",
			tokens: 1100,
			note: "Existing working context before tool use.",
			metadata: { turnIndex: 4 },
		},
		{
			id: "tool-call",
			label: "T5 bash call",
			styleKey: "tool",
			tokens: 900,
			note: "Search command.",
			metadata: { turnIndex: 5, toolName: "bash" },
		},
		{
			id: "tool-result",
			label: "T5 search output",
			styleKey: "result",
			tokens: 2200,
			note: "Search results showing the files to change.",
			metadata: { turnIndex: 5, fullBytes: 4720 },
		},
		{
			id: "t5-file",
			label: "T5 file read",
			styleKey: "retrieval",
			tokens: 950,
			note: "Source file content used for implementation.",
			metadata: { turnIndex: 5 },
		},
		{
			id: "answer",
			label: "T6 assistant",
			styleKey: "active",
			tokens: 1700,
			note: "Current answer draft.",
			metadata: { turnIndex: 6 },
		},
		{
			id: "t6-summary",
			label: "T6 summary",
			styleKey: "summary",
			tokens: 1050,
			note: "Condensed implementation summary.",
			metadata: { turnIndex: 6 },
		},
		{
			id: "free",
			label: "headroom",
			styleKey: "empty",
			tokens: 1600,
			metadata: { turnIndex: "headroom", blockType: "headroom" },
		},
	],
};
const turnPagerSnapshots: ContextWindowSnapshot[] = [4, 5, 6].map((turn) => ({
	...groupedSessionSnapshot,
	id: `widget-ir-turn-${turn}`,
	title: `Turn ${turn} context window`,
	subtitle:
		turn === 4
			? "User request enters the window."
			: turn === 5
				? "Tool call/result dominate the working set."
				: "Assistant answer becomes active.",
	selectedPartId: turn === 4 ? "user-turn" : turn === 5 ? "tool-result" : "answer",
	metadata: { turnIndex: turn },
	parts: groupedSessionSnapshot.parts.map((part) => ({
		...part,
		tokens:
			part.id === "free"
				? Math.max(2000, part.tokens - turn * 320)
				: Math.max(1, part.tokens + (part.metadata?.turnIndex === turn ? 260 : 0)),
	})),
}));
const threeLabelStyleSet = contextThreeLabelStyleSets[0]!;
const threeLabelSnapshot: ContextWindowSnapshot = {
	id: "widget-ir-three-label",
	title: "Widget IR three-label context",
	limit: 32_000,
	selectedPartId: "retrieved-docs",
	parts: [
		{ id: "prompt", label: "Prompt", styleKey: "prompt", tokens: 1400 },
		{ id: "retrieved-docs", label: "Evidence", styleKey: "evidence", tokens: 9200 },
		{ id: "answer-draft", label: "Draft", styleKey: "answer", tokens: 1800 },
		{ id: "free", label: "Free", styleKey: "free", tokens: 19600 },
	],
};

function panel(title: string, children: WidgetNode[]): WidgetNode {
	return component("Panel", { title, density: "condensed" }, children);
}
function renderNode(node: WidgetNode, registry = defaultWidgetRegistry) {
	return <WidgetRenderer registry={registry} node={node} />;
}

const fallbackVisualStyle: ContextVisualStyle = {
	pattern: "overflow",
	fill: "var(--mac-surface)",
	line: "var(--mac-border)",
	stroke: "var(--mac-border)",
	labelColor: "var(--mac-text)",
};

function visual(base: ContextStyleSet, key: string): ContextVisualStyle {
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

function compactPanelExamples(base: ContextStyleSet): Array<{
	title: string;
	view: "strip" | "budget" | "stack" | "treemap";
	snapshot: ContextWindowSnapshot;
	styleSet: ContextStyleSet;
}> {
	const ragStyleSet = compactStyleSet(base, "story-rag-answer", "RAG answer triad", [
		{ id: "brief", label: "Brief", styleFrom: "system" },
		{ id: "evidence", label: "Evidence", styleFrom: "result" },
		{ id: "draft", label: "Draft", styleFrom: "active" },
		{ id: "free", label: "Headroom", styleFrom: "empty", hidden: true },
	]);
	const opsStyleSet = compactStyleSet(base, "story-ops-payload", "Ops payload", [
		{ id: "guardrails", label: "Guardrails", styleFrom: "system" },
		{ id: "chat", label: "Chat", styleFrom: "conversation" },
		{ id: "commands", label: "Commands", styleFrom: "tool" },
		{ id: "output", label: "Output", styleFrom: "result" },
		{ id: "free", label: "Free", styleFrom: "empty", hidden: true },
	]);
	const editorStyleSet = compactStyleSet(base, "story-editor-budget", "Editor budget", [
		{ id: "spec", label: "Spec", styleFrom: "context" },
		{ id: "source", label: "Source files", styleFrom: "retrieval" },
		{ id: "tests", label: "Test logs", styleFrom: "result" },
		{ id: "answer", label: "Answer", styleFrom: "active" },
		{ id: "spare", label: "Spare room", styleFrom: "empty" },
	]);
	const compressionStyleSet = compactStyleSet(base, "story-compression-plan", "Compression plan", [
		{ id: "keep", label: "Keep verbatim", styleFrom: "active" },
		{ id: "summarize", label: "Summarize", styleFrom: "summary" },
		{ id: "archive", label: "Archive refs", styleFrom: "generated" },
		{ id: "discard", label: "Drop noise", styleFrom: "evicted" },
		{ id: "new", label: "New turn", styleFrom: "conversation" },
		{ id: "room", label: "Room after trim", styleFrom: "empty" },
	]);
	return [
		{
			title: "3 visible labels: RAG answer",
			view: "strip",
			styleSet: ragStyleSet,
			snapshot: {
				id: "rag-answer-mini-vocab",
				title: "RAG answer budget",
				subtitle: "Only caller-defined Brief / Evidence / Draft entries appear in the legend.",
				limit: 32_000,
				selectedPartId: "rag-evidence",
				parts: [
					{ id: "rag-brief", label: "task brief", styleKey: "brief", tokens: 1_500 },
					{ id: "rag-evidence", label: "retrieved evidence", styleKey: "evidence", tokens: 9_800 },
					{ id: "rag-draft", label: "draft answer", styleKey: "draft", tokens: 2_600 },
					{ id: "rag-free", label: "headroom", styleKey: "free", tokens: 18_100 },
				],
			},
		},
		{
			title: "4 visible labels: Ops payload",
			view: "budget",
			styleSet: opsStyleSet,
			snapshot: {
				id: "ops-payload-mini-vocab",
				title: "Ops payload mix",
				subtitle: "A server can rename system/tool/result into its own operational vocabulary.",
				limit: 50_000,
				selectedPartId: "ops-output",
				parts: [
					{
						id: "ops-guardrails",
						label: "policy + guardrails",
						styleKey: "guardrails",
						tokens: 4_200,
					},
					{ id: "ops-chat", label: "operator chat", styleKey: "chat", tokens: 7_400 },
					{ id: "ops-commands", label: "commands", styleKey: "commands", tokens: 3_600 },
					{ id: "ops-output", label: "stdout/stderr", styleKey: "output", tokens: 24_000 },
					{ id: "ops-free", label: "free", styleKey: "free", tokens: 10_800 },
				],
			},
		},
		{
			title: "5 labels: Editor budget",
			view: "stack",
			styleSet: editorStyleSet,
			snapshot: {
				id: "editor-budget-mini-vocab",
				title: "Patch authoring window",
				subtitle: "Five visible legend entries, including free space as a first-class category.",
				limit: 100_000,
				selectedPartId: "editor-source",
				parts: [
					{ id: "editor-spec", label: "ticket spec", styleKey: "spec", tokens: 5_400 },
					{ id: "editor-source", label: "source files", styleKey: "source", tokens: 41_000 },
					{ id: "editor-tests", label: "test logs", styleKey: "tests", tokens: 18_000 },
					{ id: "editor-answer", label: "current patch answer", styleKey: "answer", tokens: 3_200 },
					{ id: "editor-spare", label: "spare room", styleKey: "spare", tokens: 32_400 },
				],
			},
		},
		{
			title: "6 labels: Compression plan",
			view: "treemap",
			styleSet: compressionStyleSet,
			snapshot: {
				id: "compression-plan-mini-vocab",
				title: "Before/after compression plan",
				subtitle: "The legend expresses product actions, not baked-in part kinds.",
				limit: 80_000,
				selectedPartId: "compress-summarize",
				parts: [
					{ id: "compress-keep", label: "key instructions", styleKey: "keep", tokens: 8_000 },
					{
						id: "compress-summarize",
						label: "summarize prior turns",
						styleKey: "summarize",
						tokens: 22_000,
					},
					{
						id: "compress-archive",
						label: "archive references",
						styleKey: "archive",
						tokens: 9_500,
					},
					{ id: "compress-discard", label: "drop noisy logs", styleKey: "discard", tokens: 14_000 },
					{ id: "compress-new", label: "new user turn", styleKey: "new", tokens: 4_500 },
					{ id: "compress-room", label: "room after trim", styleKey: "room", tokens: 22_000 },
				],
			},
		},
	];
}

function contextDiagramGalleryNode(styleSet: ContextStyleSet): WidgetNode {
	return component("Stack", { gap: "md" }, [
		panel("Same snapshot, four diagram renderers", [
			component("DashboardGrid", { recipe: "twoColumn" }, [
				component("ContextBudgetBar", {
					snapshot,
					styleSet,
					selectedPartId: snapshot.selectedPartId,
					showLegend: true,
				}),
				component("ContextStripDiagram", {
					snapshot,
					styleSet,
					selectedPartId: snapshot.selectedPartId,
				}),
				component("ContextStackDiagram", {
					snapshot,
					styleSet,
					selectedPartId: snapshot.selectedPartId,
				}),
				component("ContextTreemap", {
					snapshot,
					styleSet,
					selectedPartId: snapshot.selectedPartId,
				}),
			]),
		]),
		component("ContextLegend", {
			items: styleSet.legend,
			styles: styleSet.styles,
			selectedId: "conversation",
		}),
	]);
}

export const ContextDiagramGallery: Story = {
	render: ({ palette, registry }) =>
		renderNode(contextDiagramGalleryNode(contextStyleSetForPalette(palette)), registry),
};

export const ContextDiagramPanelViews: Story = {
	render: ({ palette, registry }) => {
		const examples = compactPanelExamples(contextStyleSetForPalette(palette));
		return renderNode(
			component(
				"DashboardGrid",
				{ recipe: "twoColumn" },
				examples.map((example) =>
					component("ContextDiagramPanel", {
						snapshot: example.snapshot,
						styleSet: example.styleSet,
						initialView: example.view,
						selectedPartId: example.snapshot.selectedPartId,
						views: ["strip", "budget", "stack", "treemap"],
						showPartDetails: true,
					}),
				),
			),
			registry,
		);
	},
};

export const ContextDiagramPanelContentDetails: Story = {
	render: ({ palette, registry }) =>
		renderNode(
			component("ContextDiagramPanel", {
				snapshot: contentSnapshot,
				styleSet: contextStyleSetForPalette(palette),
				initialView: "stack",
				views: ["stack", "strip", "budget"],
				showPartDetails: true,
			}),
			registry,
		),
};

export const ContextGroupedStripByTurn: Story = {
	render: ({ palette, registry }) =>
		renderNode(
			component("ContextGroupedStripDiagram", {
				snapshot: groupedSessionSnapshot,
				styleSet: contextStyleSetForPalette(palette),
				groupBy: "turn",
				showGroupLabels: true,
				showPartLabels: true,
			}),
			registry,
		),
};

export const ContextTurnPagerPanelStory: Story = {
	name: "ContextTurnPagerPanel",
	render: ({ palette, registry }) =>
		renderNode(
			component("ContextTurnPagerPanel", {
				snapshots: turnPagerSnapshots,
				styleSet: contextStyleSetForPalette(palette),
				initialSnapshotId: "widget-ir-turn-5",
				title: "Uploaded session turn pager",
				diagram: "grouped-strip",
				groupBy: "turn",
				mode: "turn-only",
				includeGlobalParts: true,
				showLegend: true,
			}),
			registry,
		),
};

export const CustomThreeLabelWidgetIR: Story = {
	parameters: { controls: { include: [] } },
	args: { palette: "Dusty Magenta / Blue" },
	render: ({ registry }) =>
		renderNode(
			component("ContextDiagramPanel", {
				snapshot: threeLabelSnapshot,
				styleSet: threeLabelStyleSet,
				initialView: "strip",
				views: ["strip", "budget", "treemap"],
				showPartDetails: true,
			}),
			registry,
		),
};

export const ContextDiagramWithMetadataSidebar: Story = {
	render: ({ palette, registry }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return renderNode(
			component("SplitPane", {
				ratio: "rightNarrow",
				divider: true,
				left: component("ContextDiagramPanel", {
					snapshot,
					styleSet,
					initialView: "treemap",
					selectedPartId: snapshot.selectedPartId,
				}),
				right: component("Stack", { gap: "md" }, [
					panel("Window metadata", [
						component("MetadataGrid", {
							density: "compact",
							items: [
								{
									key: "Window ID",
									value: component("CodeText", {}, [text(snapshot.id)]),
									copyValue: snapshot.id,
								},
								{ key: "Limit", value: `${snapshot.limit.toLocaleString()} tokens` },
								{ key: "Parts", value: snapshot.parts.length },
								{ key: "Selected", value: snapshot.selectedPartId ?? "none" },
							],
						}),
					]),
					panel("Legend", [
						component("ContextLegend", {
							items: styleSet.legend,
							styles: styleSet.styles,
							size: "sm",
							selectedId: "conversation",
						}),
					]),
				]),
			}),
			registry,
		);
	},
};

export const OverBudgetContextWindow: Story = {
	render: ({ palette, registry }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return renderNode(
			component("Stack", { gap: "md" }, [
				component("ContextBudgetBar", { snapshot: overBudgetSnapshot, styleSet, showLegend: true }),
				component("Inline", { gap: "sm", wrap: true }, [
					component("AnnotationBadge", {
						visualStyle: styleSet.styles.evicted,
						label: "eviction risk",
						selected: true,
					}),
					component("Caption", { tone: "warning" }, [
						text("Budget state is intentionally over the configured limit."),
					]),
				]),
				component("ContextStripDiagram", {
					snapshot: overBudgetSnapshot,
					styleSet,
					showLabels: true,
				}),
			]),
			registry,
		);
	},
};
