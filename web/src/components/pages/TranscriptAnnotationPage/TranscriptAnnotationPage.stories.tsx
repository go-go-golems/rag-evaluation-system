import type { Meta, StoryObj } from "@storybook/react-vite";
import { TranscriptAnnotationPage } from "./TranscriptAnnotationPage";
const meta = {
	title: "Pages/TranscriptAnnotationPage",
	component: TranscriptAnnotationPage,
} satisfies Meta<typeof TranscriptAnnotationPage>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Default: Story = {};
