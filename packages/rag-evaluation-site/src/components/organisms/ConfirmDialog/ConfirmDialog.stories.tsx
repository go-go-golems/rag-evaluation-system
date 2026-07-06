import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Button } from "../../atoms";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { ConfirmDialog } from "./ConfirmDialog";

const meta = {
	title: "Component Library/Organisms/ConfirmDialog",
	component: ConfirmDialog,
	args: {
		open: true,
		mode: "inline",
		title: "Archive article",
		message: "Archive “Workshop intro”?",
		onConfirm: () => {},
		onCancel: () => {},
	},
} satisfies Meta<typeof ConfirmDialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Static: Story = {};

export const Destructive: Story = {
	args: {
		title: "Delete asset",
		message: "Delete “missing-figure.png” permanently?",
		detail: "Used in 2 articles — those references will break.",
		confirmLabel: "Delete",
		destructive: true,
	},
};

export const Interactive: Story = {
	render: () => {
		const [open, setOpen] = useState(false);
		const [result, setResult] = useState("(none)");
		return (
			<Stack gap="sm">
				<Button onClick={() => setOpen(true)}>Delete asset…</Button>
				<Caption>last result: {result}</Caption>
				<ConfirmDialog
					open={open}
					title="Delete asset"
					message="Delete “missing-figure.png” permanently?"
					confirmLabel="Delete"
					destructive
					onConfirm={() => {
						setResult("confirmed");
						setOpen(false);
					}}
					onCancel={() => {
						setResult("canceled");
						setOpen(false);
					}}
				/>
			</Stack>
		);
	},
};
