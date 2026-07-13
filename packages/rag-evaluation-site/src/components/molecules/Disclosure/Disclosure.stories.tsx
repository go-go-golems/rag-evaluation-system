import type { Meta, StoryObj } from "@storybook/react-vite";
import { Disclosure } from "./Disclosure";

const meta = {
	title: "Component Library/Molecules/Disclosure",
	component: Disclosure,
	args: {
		title: "Filter controls",
		open: true,
		children:
			"Grouped controls remain available while the section can be folded when the workspace needs more room.",
	},
} satisfies Meta<typeof Disclosure>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Open: Story = {};
export const Closed: Story = { args: { open: false } };
