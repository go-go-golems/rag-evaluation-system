import type { Meta, StoryObj } from "@storybook/react-vite";
import { KeyPointList } from "./KeyPointList";

const meta = {
	title: "Component Library/Molecules/KeyPointList",
	component: KeyPointList,
	args: {
		items: [
			"Laid flat, the window reads left to right as a strip of segments.",
			"Separators mark where one tenant ends and the next begins.",
			"This is the fastest way to sanity-check what you packed.",
		],
	},
} satisfies Meta<typeof KeyPointList>;
export default meta;
type Story = StoryObj<typeof meta>;
export const SlidePoints: Story = {};
export const WithTitles: Story = {
	args: {
		items: [
			{
				title: "Read the budget",
				text: "Start with the largest tenant and ask whether it still pays rent.",
			},
			{
				title: "Mark the churn",
				text: "Scratchpad and tool output should age out quickly.",
				meta: "Rule 02",
			},
		],
	},
};
