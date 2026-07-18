import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { KeyboardShortcutHelp } from "./KeyboardShortcutHelp";

const meta = {
	title: "Component Library/Molecules/KeyboardShortcutHelp",
	component: KeyboardShortcutHelp,
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof KeyboardShortcutHelp>;

export default meta;
type Story = StoryObj<typeof meta>;

export const TriageCommands: Story = {
	args: {
		items: [
			{ id: "accept", label: "Yes", chord: "y" },
			{ id: "reject", label: "No", chord: "n" },
			{ id: "skip", label: "Skip", chord: "s" },
		],
		enabled: true,
		onEnabledChange: () => undefined,
	},
	render: (args) => {
		const [enabled, setEnabled] = useState(args.enabled);
		return <KeyboardShortcutHelp {...args} enabled={enabled} onEnabledChange={setEnabled} />;
	},
};
