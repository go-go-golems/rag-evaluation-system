import type { Meta, StoryObj } from "@storybook/react-vite";
import { EvaluationView } from "./EvaluationView";

const meta = {
	title: "Views/Evaluation/EvaluationView",
	component: EvaluationView,
} satisfies Meta<typeof EvaluationView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
