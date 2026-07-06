import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	type ContextPaletteName,
	type ContextStyleSet,
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextStyleSetForPalette,
} from "../context";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type WidgetNode } from "./ir";
import { WidgetRenderer, type WidgetRendererProps } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/Foundation and Atoms",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;
type PaletteControlsArgs = WidgetRendererProps & { palette: ContextPaletteName };

function panel(title: string, children: WidgetNode[]): WidgetNode {
	return component("Panel", { title, density: "condensed" }, children);
}

function semanticBadgesAndSwatchesNode(styleSet: ContextStyleSet): WidgetNode {
	return component("Stack", { gap: "md" }, [
		panel("Context style swatches", [
			component(
				"Inline",
				{ gap: "md", wrap: true },
				Object.entries(styleSet.styles)
					.slice(0, 9)
					.map(([styleKey, visualStyle], index) =>
						component("Inline", { gap: "xs", wrap: false }, [
							component("ContextStyleSwatch", { visualStyle, selected: index === 2 }),
							component("Text", { size: "compact" }, [text(styleKey)]),
						]),
					),
			),
		]),
		panel("Annotation and role badges", [
			component("Inline", { gap: "sm", wrap: true }, [
				component("AnnotationBadge", {
					visualStyle: styleSet.styles.context,
					label: "task framing",
				}),
				component("AnnotationBadge", {
					visualStyle: styleSet.styles.result,
					label: "tool result",
					selected: true,
				}),
				component("AnnotationBadge", {
					visualStyle: styleSet.styles.generated,
					label: "scratchpad",
				}),
				component("TranscriptRoleBadge", { role: "user" }),
				component("TranscriptRoleBadge", { role: "assistant" }),
				component("TranscriptRoleBadge", { role: "tool", name: "read_file" }),
			]),
		]),
	]);
}

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		registry: defaultWidgetRegistry,
		palette: "Dusty Magenta / Blue",
		node: semanticBadgesAndSwatchesNode(contextDefaultStyleSet),
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		node: { control: false },
		registry: { control: false },
		onAction: { control: false },
	},
	render: ({ palette, registry, onAction }) => (
		<WidgetRenderer
			registry={registry}
			onAction={onAction}
			node={semanticBadgesAndSwatchesNode(contextStyleSetForPalette(palette))}
		/>
	),
};

export const TypographyTokenSampler: Story = {
	args: {
		node: component("Stack", { gap: "md" }, [
			panel("Text roles rendered from Widget IR", [
				component("Stack", { gap: "sm" }, [
					component("Text", { size: "body" }, [
						text(
							"Body prose: the renderer maps a JSON component node to the package Text primitive.",
						),
					]),
					component("Text", { size: "compact", tone: "muted" }, [
						text("Compact muted metadata for dense interfaces."),
					]),
					component("Text", { size: "metadata", tone: "accent" }, [text("Metadata accent text")]),
					component("Text", { size: "label", weight: "bold" }, [text("UPPERCASE LABEL ROLE")]),
					component("Text", { size: "metric", tone: "success", weight: "bold" }, [
						text("42,918 tokens"),
					]),
				]),
			]),
			component("Divider", { orientation: "horizontal" }),
			panel("Code and captions", [
				component("Stack", { gap: "sm" }, [
					component("CodeText", { tone: "primary" }, [text("ctx.window.parts[3].sourceId")]),
					component("CodeText", { tone: "muted", display: "block" }, [
						text("/api/widget/pages/context-studio"),
					]),
					component("Caption", { tone: "warning" }, [
						text("These nodes are rendered through WidgetRenderer, not direct JSX."),
					]),
				]),
			]),
		]),
	},
};

export const SemanticBadgesAndSwatches: Story = {
	args: {
		node: semanticBadgesAndSwatchesNode(contextDefaultStyleSet),
	},
};

export const InlineStatusHeader: Story = {
	args: {
		node: panel("Inline status header", [
			component("Inline", { gap: "sm", justify: "between", wrap: true }, [
				component("Inline", { gap: "sm", wrap: true }, [
					component("StatusText", { status: "running", icon: true }, [text("running")]),
					component("AnnotationBadge", {
						visualStyle: contextDefaultStyleSet.styles.active,
						label: "active context",
						selected: true,
					}),
					component("Caption", { tone: "muted" }, [text("3 notes · 8 messages")]),
				]),
				component("Inline", { gap: "sm", wrap: true }, [
					component("Button", { size: "compact" }, [text("Copy summary")]),
					component("Button", { size: "compact", variant: "primary" }, [text("Open")]),
				]),
			]),
		]),
	},
};
