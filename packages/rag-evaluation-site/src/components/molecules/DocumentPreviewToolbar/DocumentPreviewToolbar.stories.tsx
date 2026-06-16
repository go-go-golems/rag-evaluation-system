import type { Meta, StoryObj } from "@storybook/react-vite";
import { DocumentPreviewToolbar } from "./DocumentPreviewToolbar";

const meta = {
	title: "Component Library/Molecules/DocumentPreviewToolbar",
	component: DocumentPreviewToolbar,
	args: {
		file: "context-window-field-guide.md",
		format: "Markdown",
		size: "11 KB",
		onDownload: () => {},
	},
} satisfies Meta<typeof DocumentPreviewToolbar>;

export default meta;
type Story = StoryObj<typeof meta>;

export const MarkdownDocument: Story = {};

export const LongFilename: Story = {
	args: { file: "very-long-context-window-engineering-workshop-field-guide-with-appendices.md" },
};
