import type { Meta, StoryObj } from "@storybook/react-vite";
import { anchoredCommentFixtures, transcriptFixture } from "../context";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type WidgetNode } from "./ir";
import { WidgetRenderer, type WidgetRendererProps } from "./WidgetRenderer";

type StoryArgs = WidgetRendererProps;

const meta = {
	title: "Widget IR/Renderer/Transcript and Notes",
	args: { registry: defaultWidgetRegistry },
	argTypes: {
		node: { control: false },
		registry: { control: false },
		onAction: { control: false },
	},
} satisfies Meta<StoryArgs>;
export default meta;
type Story = StoryObj<StoryArgs>;

function panel(title: string, children: WidgetNode[]): WidgetNode {
	return component("Panel", { title, density: "condensed" }, children);
}

function renderNode(node: WidgetNode, registry = defaultWidgetRegistry) {
	return <WidgetRenderer node={node} registry={registry} />;
}

function transcriptWorkspaceNode(showNotes = true): WidgetNode {
	return component("TranscriptWorkspacePanel", {
		title: transcriptFixture.title,
		subtitle: transcriptFixture.subtitle,
		messages: transcriptFixture.messages,
		annotations: showNotes ? transcriptFixture.annotations : [],
		selectedAnnotationId: showNotes ? transcriptFixture.selectedAnnotationId : undefined,
		showNotes,
		onAnnotationSelectAction: { kind: "event", event: "widget-ir:annotation-select" },
	});
}

function messageCardStatesNode(): WidgetNode {
	return component(
		"Stack",
		{ gap: "sm" },
		transcriptFixture.messages.slice(0, 5).map((message) =>
			component("TranscriptMessageCard", {
				message,
				annotations: transcriptFixture.annotations,
				selectedAnnotationId: transcriptFixture.selectedAnnotationId,
				onAnnotationSelectAction: { kind: "event", event: "widget-ir:annotation-select" },
			}),
		),
	);
}

export const TranscriptWithoutNotes: Story = {
	render: ({ registry }) => renderNode(transcriptWorkspaceNode(false), registry),
};

export const AnnotatedTranscriptWithNotesRail: Story = {
	render: ({ registry }) => renderNode(transcriptWorkspaceNode(true), registry),
};

export const TranscriptReaderPlusCustomRail: Story = {
	render: ({ registry }) =>
		renderNode(
			component("SplitPane", {
				ratio: "rightNarrow",
				divider: true,
				left: component("TranscriptReaderPanel", {
					title: "Manual transcript composition",
					subtitle: transcriptFixture.subtitle,
					messages: transcriptFixture.messages,
					annotations: transcriptFixture.annotations,
					selectedAnnotationId: transcriptFixture.selectedAnnotationId,
					onAnnotationSelectAction: { kind: "event", event: "widget-ir:annotation-select" },
				}),
				right: component("AnnotationRailPanel", {
					annotations: transcriptFixture.annotations,
					selectedAnnotationId: transcriptFixture.selectedAnnotationId,
					onAnnotationSelectAction: { kind: "event", event: "widget-ir:annotation-select" },
				}),
			}),
			registry,
		),
};

export const MessageCardStates: Story = {
	render: ({ registry }) => renderNode(messageCardStatesNode(), registry),
};

export const AnchoredCommentsOverTranscript: Story = {
	render: ({ registry }) =>
		renderNode(
			component("SplitPane", {
				ratio: "rightNarrow",
				divider: true,
				left: component("TranscriptReaderPanel", {
					title: "Transcript with anchored comments nearby",
					messages: transcriptFixture.messages.slice(0, 5),
					annotations: transcriptFixture.annotations,
					selectedAnnotationId: transcriptFixture.selectedAnnotationId,
				}),
				right: component("AnchoredCommentRail", {
					comments: anchoredCommentFixtures,
					selectedCommentId: anchoredCommentFixtures[0]?.id,
					onCommentSelectAction: { kind: "event", event: "widget-ir:comment-select" },
				}),
			}),
			registry,
		),
};

export const TranscriptActionLogger: Story = {
	args: { node: component("Caption", {}, [text("action logger")]) },
	render: () => (
		<WidgetRenderer
			node={component("Stack", { gap: "md" }, [
				panel("Action logger harness", [
					component("Caption", { tone: "muted" }, [
						text(
							"Open the browser console and click note chips or comments to inspect emitted Widget IR actions.",
						),
					]),
				]),
				transcriptWorkspaceNode(true),
			])}
			registry={defaultWidgetRegistry}
			onAction={(action, context) => {
				console.info("Widget IR action", { action, context });
			}}
		/>
	),
};
