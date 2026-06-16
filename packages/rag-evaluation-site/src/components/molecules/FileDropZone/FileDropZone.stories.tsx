import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { FileDropZone } from "./FileDropZone";

const meta = {
	title: "Component Library/Molecules/FileDropZone",
	component: FileDropZone,
	args: {
		title: "Drop a .json file here",
		description: "or paste below · max 200k tokens",
		accept: "application/json,.json",
	},
	decorators: [
		(Story) => (
			<div style={{ maxWidth: 420 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof FileDropZone>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const DragActive: Story = {
	args: { active: true },
};

export const Disabled: Story = {
	args: { disabled: true },
};

export const WithSelectionLog: Story = {
	render: (args) => {
		const [files, setFiles] = useState<string[]>([]);
		return (
			<Stack gap="sm">
				<FileDropZone
					{...args}
					onFilesSelected={(selected) =>
						setFiles(selected.map((file) => `${file.name} · ${file.size} bytes`))
					}
				/>
				<Caption>{files.length ? files.join(", ") : "No file selected yet."}</Caption>
			</Stack>
		);
	},
};
