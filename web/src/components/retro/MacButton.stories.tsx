import type { Meta, StoryObj } from "@storybook/react-vite";
import { MacButton } from "./MacButton";

const meta = {
	title: "Legacy Retro/MacButton",
	component: MacButton,
	args: {
		label: "Run pipeline",
		onClick: () => {},
	},
} satisfies Meta<typeof MacButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Primary: Story = {
	args: { primary: true, label: "Search" },
};

export const Disabled: Story = {
	args: { disabled: true, label: "Disabled" },
};
