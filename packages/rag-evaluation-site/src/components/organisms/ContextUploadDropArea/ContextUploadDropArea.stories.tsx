import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { ContextUploadDropArea } from "./ContextUploadDropArea";

const meta = {
	title: "Component Library/Organisms/ContextUploadDropArea",
	component: ContextUploadDropArea,
	decorators: [
		(Story) => (
			<div style={{ maxWidth: 440 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof ContextUploadDropArea>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const DragActive: Story = {
	args: { active: true },
};

export const Disabled: Story = {
	args: { disabled: true },
};

export const SelectionState: Story = {
	render: (args) => {
		const [selected, setSelected] = useState("No context-window file selected.");
		return (
			<Stack gap="sm">
				<ContextUploadDropArea
					{...args}
					onFilesSelected={(files) => setSelected(files.map((file) => file.name).join(", "))}
				/>
				<Caption>{selected}</Caption>
			</Stack>
		);
	},
};
