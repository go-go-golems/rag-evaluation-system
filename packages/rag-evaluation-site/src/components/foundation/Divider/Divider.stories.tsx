import type { Meta, StoryObj } from "@storybook/react-vite";
import { Divider } from "./Divider";

const meta = {
	title: "Design System/Foundation/Divider",
	component: Divider,
} satisfies Meta<typeof Divider>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Horizontal: Story = {
	render: () => (
		<div>
			<p>Before divider</p>
			<Divider />
			<p>After divider</p>
		</div>
	),
};

export const Vertical: Story = {
	render: () => (
		<div style={{ display: "flex", height: 48, alignItems: "stretch" }}>
			<span>Left</span>
			<Divider orientation="vertical" />
			<span>Right</span>
		</div>
	),
};
