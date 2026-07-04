import type { Meta, StoryObj } from "@storybook/react-vite";
import { transcriptFixture } from "../../../context";
import { Stack } from "../../layout";
import { AnnotationNoteCard } from "./AnnotationNoteCard";

const meta = {
	title: "Component Library/Molecules/AnnotationNoteCard",
	component: AnnotationNoteCard,
	args: { annotation: transcriptFixture.annotations[0]!, index: 1 },
} satisfies Meta<typeof AnnotationNoteCard>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Notes: Story = {
	render: () => (
		<Stack gap="sm">
			{transcriptFixture.annotations.map((annotation, index) => (
				<AnnotationNoteCard
					key={annotation.id}
					annotation={annotation}
					index={index + 1}
					selected={annotation.id === transcriptFixture.selectedAnnotationId}
				/>
			))}
		</Stack>
	),
};
