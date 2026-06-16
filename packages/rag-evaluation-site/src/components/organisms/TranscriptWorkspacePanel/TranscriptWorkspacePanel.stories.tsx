import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { transcriptFixture } from "../../../context";
import { TranscriptWorkspacePanel } from "./TranscriptWorkspacePanel";

const meta = {
	title: "Component Library/Organisms/TranscriptWorkspacePanel",
	component: TranscriptWorkspacePanel,
	args: {
		title: transcriptFixture.title,
		subtitle: transcriptFixture.subtitle,
		messages: transcriptFixture.messages,
	},
} satisfies Meta<typeof TranscriptWorkspacePanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const WithoutNotes: Story = {
	args: {
		annotations: [],
		showNotes: false,
	},
};

export const WithNotes: Story = {
	render: () => {
		const [selected, setSelected] = useState(transcriptFixture.selectedAnnotationId);
		return (
			<TranscriptWorkspacePanel
				title={transcriptFixture.title}
				subtitle={transcriptFixture.subtitle}
				messages={transcriptFixture.messages}
				annotations={transcriptFixture.annotations}
				selectedAnnotationId={selected}
				onAnnotationSelect={setSelected}
			/>
		);
	},
};

export const ToolHeavyWithoutNotes: Story = {
	args: {
		title: "Tool-heavy excerpt",
		messages: transcriptFixture.messages.slice(2, 8),
		annotations: [],
		showNotes: false,
	},
};
