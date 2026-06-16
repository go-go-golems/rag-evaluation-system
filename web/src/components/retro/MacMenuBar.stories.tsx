import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { MacMenuBar } from "./MacMenuBar";

const meta = {
	title: "Legacy Retro/MacMenuBar",
	component: MacMenuBar,
	args: {
		activeView: "pipeline",
		onViewChange: () => {},
	},
} satisfies Meta<typeof MacMenuBar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Interactive: Story = {
	render: () => {
		const [activeView, setActiveView] = useState("pipeline");
		return <MacMenuBar activeView={activeView} onViewChange={setActiveView} />;
	},
};
