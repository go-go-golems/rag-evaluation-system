import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inline } from "../../layout";
import { IconButton } from "./IconButton";

const meta = {
	title: "Design System/Atoms/IconButton",
	component: IconButton,
} satisfies Meta<typeof IconButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Actions: Story = {
	render: () => (
		<Inline>
			<IconButton label="Close">✕</IconButton>
			<IconButton label="Copy chunk identifier">⧉</IconButton>
			<IconButton label="Back">← Back</IconButton>
			<IconButton label="Disabled" disabled>
				✕
			</IconButton>
		</Inline>
	),
};
