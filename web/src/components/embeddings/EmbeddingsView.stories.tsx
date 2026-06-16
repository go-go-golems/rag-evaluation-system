import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../storybook/MockApiProvider";
import { EmbeddingsView } from "./EmbeddingsView";

const meta = {
	title: "Views/Embeddings/EmbeddingsView",
	component: EmbeddingsView,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof EmbeddingsView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
