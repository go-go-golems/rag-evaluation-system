import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import {
	contextDefaultStyleSet,
	contextHandoutFixture,
	contextSignalOrangeStyleSet,
	contextWindowSnapshots,
} from "../../../context";
import { HandoutDocumentShell } from "./HandoutDocumentShell";

const meta = {
	title: "Component Library/Organisms/HandoutDocumentShell",
	component: HandoutDocumentShell,
	args: {
		intro: contextHandoutFixture.intro,
		documents: contextHandoutFixture.docs,
		selectedDocumentId: contextHandoutFixture.docs[0]?.id,
		onDownload: () => {},
		onDownloadAll: () => {},
	},
} satisfies Meta<typeof HandoutDocumentShell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const FieldGuideSelected: Story = {};

export const ChecklistSelected: Story = {
	args: { selectedDocumentId: "checklist" },
};

const illustratedDocument = {
	id: "illustrated",
	title: "Illustrated Context Checklist",
	file: "illustrated-context-checklist.md",
	format: "Markdown",
	size: "5 KB",
	description: "A handout with prose, an image, and an embedded context-window diagram.",
	body: "# Illustrated Context Checklist\n\n![Budget sketch](/course-assets/context-window-token-budget.svg)\n",
	blocks: [
		{
			kind: "markdown" as const,
			id: "intro",
			source:
				"# Illustrated Context Checklist\n\nImages render as bounded figures.\n\n![Budget sketch](/course-assets/context-window-token-budget.svg)",
		},
		{
			kind: "context-window" as const,
			id: "diagram",
			snapshot: contextWindowSnapshots[0]!,
			view: "budget" as const,
			caption: "Token budget",
		},
		{
			kind: "markdown" as const,
			id: "outro",
			source: "## Review prompt\n\nWhat should survive into the next model call?",
		},
	],
};

const illustratedSecondDocument = {
	id: "illustrated-stack",
	title: "Illustrated Stack Walkthrough",
	file: "illustrated-stack-walkthrough.md",
	format: "Markdown",
	size: "7 KB",
	description: "A second illustrated handout that starts in stack view.",
	body: "# Illustrated Stack Walkthrough\n",
	blocks: [
		{
			kind: "markdown" as const,
			id: "stack-intro",
			source:
				"# Illustrated Stack Walkthrough\n\nThis handout starts with prose and then shows the same session as a stack.",
		},
		{
			kind: "context-window" as const,
			id: "stack-diagram",
			snapshot: contextWindowSnapshots[1]!,
			view: "stack" as const,
			caption: "Deep bug context",
		},
		{
			kind: "markdown" as const,
			id: "stack-outro",
			source:
				"- Notice that the legend is limited to visible styles.\n- Notice there is no heavy panel chrome inside the article.",
		},
	],
};

export const IllustratedHandout: Story = {
	args: {
		selectedDocumentId: "illustrated",
		documents: [illustratedDocument],
		styleSet: contextDefaultStyleSet,
	},
};

export const IllustratedSignalOrange: Story = {
	args: {
		selectedDocumentId: "illustrated",
		documents: [illustratedDocument],
		styleSet: contextSignalOrangeStyleSet,
	},
};

export const MixedDocumentBundle: Story = {
	args: {
		selectedDocumentId: "illustrated-stack",
		documents: [illustratedDocument, illustratedSecondDocument, ...contextHandoutFixture.docs],
		styleSet: contextDefaultStyleSet,
	},
};

export const Interactive: Story = {
	render: () => {
		const [selectedDocumentId, setSelectedDocumentId] = useState(contextHandoutFixture.docs[0]?.id);
		return (
			<HandoutDocumentShell
				intro={contextHandoutFixture.intro}
				documents={contextHandoutFixture.docs}
				selectedDocumentId={selectedDocumentId}
				onDocumentSelect={setSelectedDocumentId}
				onDownload={() => {}}
				onDownloadAll={() => {}}
			/>
		);
	},
};
