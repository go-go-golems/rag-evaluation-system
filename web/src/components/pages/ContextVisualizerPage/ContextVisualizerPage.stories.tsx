import type { Meta, StoryObj } from "@storybook/react-vite";
import { ContextVisualizerPage } from "./ContextVisualizerPage";
const meta = {
	title: "Pages/ContextVisualizerPage",
	component: ContextVisualizerPage,
} satisfies Meta<typeof ContextVisualizerPage>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {};
