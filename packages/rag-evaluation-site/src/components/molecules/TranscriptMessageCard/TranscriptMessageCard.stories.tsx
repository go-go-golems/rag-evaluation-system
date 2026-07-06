import type { Meta, StoryObj } from "@storybook/react-vite";
import { transcriptFixture } from "../../../context";
import { Stack } from "../../layout";
import { TranscriptMessageCard } from "./TranscriptMessageCard";

const meta = {
	title: "Component Library/Molecules/TranscriptMessageCard",
	component: TranscriptMessageCard,
	args: { message: transcriptFixture.messages[1]!, annotations: transcriptFixture.annotations },
} satisfies Meta<typeof TranscriptMessageCard>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Roles: Story = {
	render: () => (
		<Stack gap="sm">
			{transcriptFixture.messages.slice(0, 6).map((message) => (
				<TranscriptMessageCard
					key={message.id}
					message={message}
					annotations={transcriptFixture.annotations}
					selectedAnnotationId="task-framing"
				/>
			))}
		</Stack>
	),
};
export const ToolOutput: Story = {
	render: () => (
		<TranscriptMessageCard
			message={transcriptFixture.messages[4]!}
			annotations={transcriptFixture.annotations}
			selectedAnnotationId="stack-trace-noise"
		/>
	),
};
