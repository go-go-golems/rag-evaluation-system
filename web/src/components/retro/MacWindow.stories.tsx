import type { Meta, StoryObj } from "@storybook/react-vite";
import { MacButton } from "./MacButton";
import { MacWindow } from "./MacWindow";

const meta = {
	title: "Legacy Retro/MacWindow",
	component: MacWindow,
	args: {
		title: "RAG Evaluation",
		children: "Window content",
	},
} satisfies Meta<typeof MacWindow>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithCloseAndActions: Story = {
	args: {
		title: "Pipeline Control",
		onClose: () => {},
		children: (
			<div style={{ display: "flex", gap: 8, alignItems: "center" }}>
				<span>Ready to process documents.</span>
				<MacButton label="Start" primary onClick={() => {}} />
			</div>
		),
	},
};
