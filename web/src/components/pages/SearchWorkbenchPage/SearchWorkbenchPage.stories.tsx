import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../../storybook/MockApiProvider";
import { SearchWorkbenchPage } from "./SearchWorkbenchPage";

const meta = {
	title: "Pages/SearchWorkbenchPage",
	component: SearchWorkbenchPage,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof SearchWorkbenchPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
