import type { Meta, StoryObj } from "@storybook/react-vite";
import { transcriptFixture } from "../../../context";
import { TranscriptSessionHeader } from "./TranscriptSessionHeader";

const tokenTotal = transcriptFixture.messages.reduce(
	(total, message) => total + (message.tokens ?? 0),
	0,
);

const meta = {
	title: "Component Library/Molecules/TranscriptSessionHeader",
	component: TranscriptSessionHeader,
	args: {
		title: transcriptFixture.title,
		subtitle: transcriptFixture.subtitle,
		messageCount: transcriptFixture.messages.length,
		annotationCount: transcriptFixture.annotations.length,
		tokenTotal,
	},
} satisfies Meta<typeof TranscriptSessionHeader>;

export default meta;
type Story = StoryObj<typeof meta>;

export const SessionSummary: Story = {};
