import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { ShareLink } from "./ShareLink";

const meta = {
	title: "Component Library/Molecules/ShareLink",
	component: ShareLink,
	args: {
		label: "Share link",
		href: "/pages/poll?poll=1",
	},
} satisfies Meta<typeof ShareLink>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithDescription: Story = {
	args: {
		description: "Send this URL to participants so they can submit their availability.",
	},
};

export const LongLink: Story = {
	args: {
		href: "/pages/poll?poll=123456789&day=2026-07-11&slot=slot-987654321",
		description: "Long URLs stay on one line and truncate inside the value field.",
	},
};

export const CopyState: Story = {
	render: (args) => {
		const [copied, setCopied] = useState(false);
		return (
			<ShareLink
				{...args}
				copied={copied}
				onCopy={() => {
					setCopied(true);
				}}
			/>
		);
	},
};
