import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { anchoredCommentFixtures } from "../../../context";
import { AnchoredCommentRail } from "./AnchoredCommentRail";
const meta = {
	title: "Component Library/Organisms/AnchoredCommentRail",
	component: AnchoredCommentRail,
	args: { comments: anchoredCommentFixtures },
} satisfies Meta<typeof AnchoredCommentRail>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Rail: Story = {
	render: () => {
		const [selected, setSelected] = useState(anchoredCommentFixtures[0]?.id);
		return (
			<AnchoredCommentRail
				comments={anchoredCommentFixtures}
				selectedCommentId={selected}
				onCommentSelect={setSelected}
			/>
		);
	},
};
