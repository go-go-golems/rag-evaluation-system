import type { Meta, StoryObj } from "@storybook/react-vite";
import type { ArticleBlock, ContextWindowSnapshot } from "../../../context";
import {
	contextDefaultStyleSet,
	contextSignalOrangeStyleSet,
	contextWindowSnapshots,
} from "../../../context";
import { Stack } from "../../layout";
import { RichArticle } from "./RichArticle";

const [, deepBug, atLimit] = contextWindowSnapshots;

const compactSnapshot: ContextWindowSnapshot = {
	id: "compact-handout-window",
	title: "Handout example — compact window",
	subtitle: "Only three non-empty style keys should appear in the legend.",
	limit: 100_000,
	selectedPartId: "tool-output",
	parts: [
		{
			id: "durable",
			label: "durable instructions",
			styleKey: "system",
			tokens: 12_000,
			note: "System and project rules.",
		},
		{
			id: "tool-output",
			label: "tool output",
			styleKey: "result",
			tokens: 38_000,
			note: "Large grep/test output.",
		},
		{
			id: "summary",
			label: "summary",
			styleKey: "summary",
			tokens: 8_000,
			note: "Compressed decisions.",
		},
		{ id: "free", label: "free space", styleKey: "empty", tokens: 42_000, note: "Headroom." },
	],
};

const illustratedBlocks: ArticleBlock[] = [
	{
		kind: "markdown",
		id: "intro",
		source: `# Illustrated handout section

A rich handout can mix regular Markdown prose, bounded images, and context-window diagrams without switching authoring systems.

![Context window token budget sketch](/course-assets/context-window-token-budget.svg)

## What to notice

- The image has no decorative border.
- The diagram below uses inline chrome rather than a heavy panel.
- The legend only includes style keys visible in the diagram.`,
	},
	{ kind: "context-window", id: "compact-diagram", snapshot: compactSnapshot, view: "budget" },
	{
		kind: "markdown",
		id: "outro",
		source: `## Review prompt

What should remain exact, what should become a summary, and what can be evicted?`,
	},
];

const multiDiagramBlocks: ArticleBlock[] = [
	{
		kind: "markdown",
		id: "multi-intro",
		source:
			"# Multi-diagram handout\n\nUse multiple diagrams when the prose compares phases of a session.",
	},
	{ kind: "context-window", id: "early", snapshot: deepBug!, view: "stack" },
	{
		kind: "markdown",
		id: "between",
		source:
			"The next diagram starts closer to the limit and should still keep the legend scoped to visible parts.",
	},
	{ kind: "context-window", id: "late", snapshot: atLimit!, view: "budget" },
];

const imageBlockArticle: ArticleBlock[] = [
	{
		kind: "markdown",
		id: "image-intro",
		source: "# Explicit image block\n\nThis story covers the non-Markdown image block path.",
	},
	{
		kind: "image",
		id: "explicit-image",
		src: "/course-assets/context-window-token-budget.svg",
		alt: "Context window budget sketch",
		caption: "A course-authored SVG served through /course-assets.",
	},
	{
		kind: "markdown",
		id: "image-outro",
		source: "The figure should remain borderless and centered.",
	},
];

const galleryBlockArticle: ArticleBlock[] = [
	{
		kind: "markdown",
		id: "gallery-intro",
		source: "# Gallery block\n\nA gallery lays out several bounded figures in a grid.",
	},
	{
		kind: "gallery",
		id: "sketch-gallery",
		columns: 3,
		caption: "Three views of the same context window.",
		images: [
			{ src: "/course-assets/context-window-token-budget.svg", alt: "Budget view" },
			{ src: "/course-assets/context-window-token-budget.svg", alt: "Strip view" },
			{ src: "/course-assets/context-window-token-budget.svg", alt: "Stack view" },
		],
	},
	{
		kind: "gallery",
		id: "two-up-gallery",
		columns: 2,
		images: [
			{
				src: "/course-assets/context-window-token-budget.svg",
				alt: "Before compaction",
				caption: "Turn 3 — warming up.",
			},
			{
				src: "/course-assets/context-window-token-budget.svg",
				alt: "After compaction",
				caption: "Turn 14 — deep in the bug.",
			},
		],
	},
];

const meta = {
	title: "Component Library/Organisms/RichArticle",
	component: RichArticle,
	args: { blocks: illustratedBlocks, styleSet: contextDefaultStyleSet },
} satisfies Meta<typeof RichArticle>;

export default meta;
type Story = StoryObj<typeof meta>;

export const IllustratedHandoutSection: Story = {};

export const SignalOrangePalette: Story = {
	args: { blocks: illustratedBlocks, styleSet: contextSignalOrangeStyleSet },
};

export const MultipleContextDiagrams: Story = {
	args: { blocks: multiDiagramBlocks, styleSet: contextDefaultStyleSet },
};

export const ExplicitImageBlock: Story = {
	args: { blocks: imageBlockArticle, styleSet: contextDefaultStyleSet },
};

export const GalleryBlock: Story = {
	args: { blocks: galleryBlockArticle, styleSet: contextDefaultStyleSet },
};

export const SideBySideDensitySamples: Story = {
	render: () => (
		<Stack gap="lg">
			<RichArticle blocks={illustratedBlocks} styleSet={contextDefaultStyleSet} />
			<RichArticle blocks={multiDiagramBlocks} styleSet={contextSignalOrangeStyleSet} />
		</Stack>
	),
};
