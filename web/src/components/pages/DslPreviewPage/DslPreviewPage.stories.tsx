import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../../storybook/MockApiProvider";
import { DslPreviewPage } from "./DslPreviewPage";

const meta = {
	title: "Pages/DslPreviewPage",
	component: DslPreviewPage,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
} satisfies Meta<typeof DslPreviewPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Demo: Story = {
	args: { pageId: "demo" },
};
