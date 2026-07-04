import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Button } from "../../atoms";
import { Text } from "../../foundation";
import { Stack } from "../Stack";
import { DialogShell } from "./DialogShell";

const meta = {
	title: "Design System/Layout/DialogShell",
	component: DialogShell,
	args: {
		open: true,
		mode: "inline",
		title: "Dialog",
		onClose: () => {},
	},
} satisfies Meta<typeof DialogShell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Static: Story = {
	args: {
		title: "Choose an asset",
		children: (
			<Text>
				Inline mode renders the same chrome without showModal, so visual diffing and layout review
				work in a plain document flow.
			</Text>
		),
	},
};

export const Sizes: Story = {
	render: () => (
		<Stack gap="md">
			<DialogShell open mode="inline" size="sm" title="Small (420px)" onClose={() => {}}>
				<Text>Confirmations.</Text>
			</DialogShell>
			<DialogShell open mode="inline" size="md" title="Medium (640px)" onClose={() => {}}>
				<Text>Forms and pickers.</Text>
			</DialogShell>
			<DialogShell open mode="inline" size="lg" title="Large (920px)" onClose={() => {}}>
				<Text>Full media library picker.</Text>
			</DialogShell>
		</Stack>
	),
};

export const WithFooterActions: Story = {
	args: {
		title: "Delete asset",
		children: <Text>This action cannot be undone.</Text>,
		footer: (
			<>
				<Button size="compact">Cancel</Button>
				<Button size="compact" variant="primary">
					Delete
				</Button>
			</>
		),
	},
};

export const Interactive: Story = {
	render: () => {
		const [open, setOpen] = useState(false);
		return (
			<Stack gap="md">
				<Button onClick={() => setOpen(true)}>Open modal dialog</Button>
				<DialogShell
					open={open}
					onClose={() => setOpen(false)}
					title="Modal dialog"
					footer={
						<Button size="compact" variant="primary" onClick={() => setOpen(false)}>
							Done
						</Button>
					}
				>
					<Text>Esc, the × button, and the footer all close this dialog.</Text>
				</DialogShell>
			</Stack>
		);
	},
};
