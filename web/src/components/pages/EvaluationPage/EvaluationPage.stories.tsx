import type { Meta, StoryObj } from "@storybook/react-vite";
import { EvaluationPage } from "./EvaluationPage";
import { MockApiProvider } from "../../../storybook/MockApiProvider";

const meta = {
	title: "Pages/EvaluationPage",
	component: EvaluationPage,
	decorators: [(Story) => <MockApiProvider><Story /></MockApiProvider>],
} satisfies Meta<typeof EvaluationPage>;

export default meta;
type Story = StoryObj<typeof meta>;

export const BaselineRunInspector: Story = {};
