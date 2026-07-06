import type { Meta, StoryObj } from "@storybook/react-vite";
import { anchoredCommentFixtures } from "../../../context";
import { Stack } from "../../layout";
import { AnchoredCommentCard } from "./AnchoredCommentCard";

const meta = {
	title: "Component Library/Molecules/AnchoredCommentCard",
	component: AnchoredCommentCard,
	args: { comment: anchoredCommentFixtures[0]! },
} satisfies Meta<typeof AnchoredCommentCard>;
export default meta;
type Story = StoryObj<typeof meta>;
export const States: Story = {
	render: () => (
		<Stack gap="sm">
			{anchoredCommentFixtures.map((comment, i) => (
				<AnchoredCommentCard
					key={comment.id}
					comment={comment}
					index={i + 1}
					selected={i === 0}
					compact={i === 1}
				/>
			))}
		</Stack>
	),
};
