import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type ActionSpec, type WidgetNode } from "./ir";
import { WidgetRenderer, type WidgetRendererProps } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/V3 Regression Fixtures",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof WidgetRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

type WidgetActionContext = Parameters<NonNullable<WidgetRendererProps["onAction"]>>[1];

const handoutDocs = [
	{
		id: "overview",
		title: "Overview",
		format: "Markdown",
		body: "# Overview\n\nRead this before starting the lab.",
	},
	{
		id: "lab",
		title: "Lab guide",
		format: "Markdown",
		body: "# Lab\n\nFollow the worksheet and capture your findings.",
	},
];

const speakerSnapshot = {
	id: "speaker-context",
	title: "Speaker context",
	limit: 2048,
	parts: [
		{
			id: "system",
			label: "System prompt",
			styleKey: "system",
			tokens: 320,
			contentPreview: "Course facilitation guardrails and output format.",
		},
		{
			id: "retrieval",
			label: "Retrieved notes",
			styleKey: "retrieval",
			tokens: 960,
			contentPreview: "Context budget examples and reclaim-policy notes.",
		},
		{
			id: "speaker",
			label: "Speaker notes",
			styleKey: "assistant",
			tokens: 420,
			contentPreview: "Explain the budget and show the diagram.",
		},
	],
};

function pageShell(title: string, child: WidgetNode): WidgetNode {
	return component("Stack", { gap: "lg", style: { padding: "16px" } }, [
		component("Text", { as: "h1", size: "title" }, [text(title)]),
		child,
	]);
}

function courseShellNode(activeItemId: string): WidgetNode {
	return component(
		"CourseStudioShell",
		{
			sections: [
				{
					id: "nav",
					label: "Navigation",
					items: [
						{ id: "overview", label: "Overview" },
						{ id: "slides", label: "Slides" },
					],
				},
			],
			activeItemId,
			title: "Course Shell",
			onNavigateAction: { kind: "event", event: "storybook:course-shell-navigate" },
		},
		[activeItemId === "slides" ? slideDeckNode() : overviewNode()],
	);
}

function overviewNode(): WidgetNode {
	return component("Panel", { title: "Overview" }, [
		component("Caption", {}, [text("Course body inside shell. Use Slides to switch this panel.")]),
	]);
}

function slideDeckNode(): WidgetNode {
	return component("Stack", { gap: "md" }, [
		component("Panel", { title: "Slides" }, [
			component("Caption", {}, [
				text("Review speaker notes, deck exports, and presentation state."),
			]),
		]),
		component("CourseSlidePanel", {
			slide: {
				id: "s1",
				title: "Course shell slide",
				view: "stack",
				notes: ["Navigation changes the shell body"],
			},
			snapshot: speakerSnapshot,
			index: 0,
			total: 1,
			mode: "speaker",
		}),
	]);
}

function handoutNode(selectedDocumentId: string): WidgetNode {
	return component("HandoutDocumentShell", {
		intro: "Downloadable workshop material",
		documents: handoutDocs,
		selectedDocumentId,
		onDocumentSelectAction: { kind: "event", event: "storybook:handout-select" },
		onDownloadAction: { kind: "download", to: "/handouts/${document.id}" },
	});
}

function transcriptWorkspaceNode(): WidgetNode {
	return component("TranscriptWorkspacePanel", {
		title: "Debug session",
		subtitle: "Retrieval failure triage",
		messages: [
			{ id: "m1", role: "user", text: "Why did retrieval fail?", tokens: 42 },
			{
				id: "m2",
				role: "assistant",
				text: "The corpus filter was too narrow, so the retriever excluded the answer document.",
				tokens: 96,
				annotationIds: ["a1"],
			},
		],
		annotations: [
			{
				id: "a1",
				targetMessageId: "m2",
				styleKey: "assistant",
				label: "Root cause",
				text: "The filter scoped retrieval to one corpus instead of the full evaluation fixture set.",
				confidence: 0.91,
			},
		],
		selectedAnnotationId: "a1",
		onAnnotationSelectAction: { kind: "event", event: "storybook:annotation-select" },
	});
}

function pageChromeNode(): WidgetNode {
	return component("Stack", { gap: "md" }, [
		component("Breadcrumbs", {
			items: [
				{ id: "home", label: "Home", href: "/pages/index" },
				{ id: "chrome", label: "Chrome" },
			],
		}),
		component(
			"SectionBlock",
			{
				label: "Actions",
				rule: true,
				actions: component("Inline", { gap: "sm", justify: "end" }, [
					component("Button", { action: { kind: "event", event: "storybook:refresh" } }, [
						text("Refresh"),
					]),
				]),
			},
			[component("Caption", {}, [text("Breadcrumbs and section actions.")])],
		),
	]);
}

function metricsNode(): WidgetNode {
	return component("Stack", { gap: "md" }, [
		component("KeyValueStrip", {
			items: [
				{ key: "Documents", label: "Documents", value: "128" },
				{ key: "Chunks", label: "Chunks", value: "42k" },
				{ key: "Recall", label: "Recall", value: "91%" },
			],
		}),
		component("Panel", { tone: "success", title: "Healthy" }, [
			component("Caption", {}, [text("The evaluation corpus is ready.")]),
		]),
	]);
}

function matrixNode(): WidgetNode {
	return component("MatrixGrid", {
		rows: [
			{ id: "q1", label: "Query 1", scores: { precision: 0.8, recall: 0.7 } },
			{ id: "q2", label: "Query 2", scores: { precision: 0.6, recall: 0.9 } },
		],
		columns: [
			{ id: "precision", header: "Precision" },
			{ id: "recall", header: "Recall" },
		],
		valueAt: { mapField: "scores" },
		cell: { kind: "value" },
		rowHeader: { kind: "field", field: "label" },
		ariaLabel: "Retrieval metric heatmap",
	});
}

function speakerSlideNode(): WidgetNode {
	return component("CourseSlidePanel", {
		slide: {
			id: "s1",
			title: "Speaker notes",
			view: "stack",
			notes: ["Explain the budget", "Show diagram"],
		},
		snapshot: speakerSnapshot,
		index: 0,
		total: 1,
		mode: "speaker",
		visualSide: "right",
	});
}

export const CourseShellTabsSwitchMainPanel: Story = {
	args: { node: pageShell("Course shell tab switching", courseShellNode("overview")) },
	render: ({ registry }) => {
		const [activeItemId, setActiveItemId] = useState("overview");
		const onAction = (_action: ActionSpec, context: WidgetActionContext) => {
			if (typeof context.itemId === "string") setActiveItemId(context.itemId);
		};
		return (
			<WidgetRenderer
				registry={registry}
				node={pageShell("Course shell tab switching", courseShellNode(activeItemId))}
				onAction={onAction}
			/>
		);
	},
};

export const HandoutTabsSwitchPreview: Story = {
	args: { node: pageShell("Handout document tab switching", handoutNode("overview")) },
	render: ({ registry }) => {
		const [selectedDocumentId, setSelectedDocumentId] = useState("overview");
		const onAction = (_action: ActionSpec, context: WidgetActionContext) => {
			if (typeof context.documentId === "string") setSelectedDocumentId(context.documentId);
		};
		return (
			<WidgetRenderer
				registry={registry}
				node={pageShell("Handout document tab switching", handoutNode(selectedDocumentId))}
				onAction={onAction}
			/>
		);
	},
};

export const TranscriptWorkspaceShowsMessageAndAnnotationText: Story = {
	args: {
		node: pageShell("Transcript workspace text regression", transcriptWorkspaceNode()),
	},
};

export const PageChromeSectionActionsRenderButtons: Story = {
	args: {
		node: pageShell("Page chrome action rendering", pageChromeNode()),
	},
};

export const DashboardMetricsIncludeLabels: Story = {
	args: {
		node: pageShell("Dashboard metrics labels", metricsNode()),
	},
};

export const MatrixHeatmapShowsHeadersAndValues: Story = {
	args: {
		node: pageShell("Matrix value cells", matrixNode()),
	},
};

export const SpeakerSlideShowsContextDiagram: Story = {
	args: {
		node: pageShell("Speaker slide context diagram", speakerSlideNode()),
	},
};
