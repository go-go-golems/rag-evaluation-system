import type { Meta, StoryObj } from "@storybook/react-vite";
import { WidgetRenderer } from "./WidgetRenderer";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component } from "./ir";

const meta = {
	title: "Widget IR/Renderer/Upload",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

export const ContextUploadDropArea: Story = {
	args: {
		node: component("ContextUploadDropArea", {
			title: "Drop a .json file here",
			description: "or paste below · max 200k tokens",
			accept: "application/json,.json",
			onFilesSelectedAction: { kind: "event", event: "widget-ir:context-upload-files-selected" },
		}),
	},
};

export const ActiveContextUploadDropArea: Story = {
	args: {
		node: component("ContextUploadDropArea", {
			active: true,
			title: "Release to import context JSON",
			description: "the renderer keeps file selection in the browser action context",
			onFilesSelectedAction: { kind: "event", event: "widget-ir:context-upload-files-selected" },
		}),
	},
};
