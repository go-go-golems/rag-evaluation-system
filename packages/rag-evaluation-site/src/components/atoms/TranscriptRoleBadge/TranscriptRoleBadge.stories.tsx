import type { Meta, StoryObj } from "@storybook/react-vite";
import type { TranscriptRole } from "../../../context";
import { Inline } from "../../layout";
import { TranscriptRoleBadge } from "./TranscriptRoleBadge";
const meta = {
	title: "Design System/Atoms/TranscriptRoleBadge",
	component: TranscriptRoleBadge,
	args: { role: "assistant" },
} satisfies Meta<typeof TranscriptRoleBadge>;
export default meta;
type Story = StoryObj<typeof meta>;
const roles: TranscriptRole[] = ["system", "developer", "user", "assistant", "tool", "other"];
export const Roles: Story = {
	render: () => (
		<Inline>
			{roles.map((role) => (
				<TranscriptRoleBadge key={role} role={role} />
			))}
			<TranscriptRoleBadge role="tool" name="read_file" />
		</Inline>
	),
};
