import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../storybook/MockApiProvider";
import { WorkflowsView } from "./WorkflowsView";

const meta = {
	title: "Views/Workflows/WorkflowsView",
	component: WorkflowsView,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof WorkflowsView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
