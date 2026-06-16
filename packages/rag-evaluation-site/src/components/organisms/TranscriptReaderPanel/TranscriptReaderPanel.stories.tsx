import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { transcriptFixture } from "../../../context";
import { TranscriptReaderPanel } from "./TranscriptReaderPanel";

const meta = {
	title: "Component Library/Organisms/TranscriptReaderPanel",
	component: TranscriptReaderPanel,
	args: {
		title: transcriptFixture.title,
		subtitle: transcriptFixture.subtitle,
		messages: transcriptFixture.messages,
		annotations: transcriptFixture.annotations,
	},
} satisfies Meta<typeof TranscriptReaderPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const AnnotatedTranscript: Story = {
	render: () => {
		const [selected, setSelected] = useState(transcriptFixture.selectedAnnotationId);
		return (
			<TranscriptReaderPanel
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

export const ToolHeavyExcerpt: Story = {
	render: () => (
		<TranscriptReaderPanel
			title="Tool-heavy excerpt"
			messages={transcriptFixture.messages.slice(2, 8)}
			annotations={transcriptFixture.annotations}
			selectedAnnotationId="stack-trace-noise"
		/>
	),
};
