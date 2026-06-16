import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../storybook/MockApiProvider";
import { PipelineView } from "./PipelineView";

const meta = {
	title: "Views/Pipeline/PipelineView",
	component: PipelineView,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
} satisfies Meta<typeof PipelineView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
