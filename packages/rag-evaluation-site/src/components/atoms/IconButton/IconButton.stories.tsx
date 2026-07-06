import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inline } from "../../layout";
import { IconButton } from "./IconButton";

const meta = {
	title: "Design System/Atoms/IconButton",
	component: IconButton,
} satisfies Meta<typeof IconButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Actions: Story = {
	render: () => (
		<Inline>
			<IconButton label="Close">✕</IconButton>
			<IconButton label="Copy chunk identifier">⧉</IconButton>
			<IconButton label="Back">← Back</IconButton>
			<IconButton label="Disabled" disabled>
				✕
			</IconButton>
		</Inline>
	),
};

export const Sizes: Story = {
	render: () => (
		<Inline gap="md">
			<IconButton label="Compact edit">✎</IconButton>
			<IconButton size="normal" label="Normal edit">
				✎
			</IconButton>
			<IconButton size="large" label="Large edit">
				✎
			</IconButton>
			<IconButton size="large" label="Large archive">
				▣
			</IconButton>
			<IconButton size="large" label="Large delete">
				×
			</IconButton>
		</Inline>
	),
};

export const Boxed: Story = {
	render: () => (
		<Inline gap="xs">
			<IconButton size="large" variant="boxed" label="Edit">
				✎
			</IconButton>
			<IconButton size="large" variant="boxed" label="Publish">
				●
			</IconButton>
			<IconButton size="large" variant="boxed" label="Archive">
				▣
			</IconButton>
			<IconButton size="large" variant="boxed" label="Delete">
				×
			</IconButton>
			<IconButton size="large" variant="boxed" label="Disabled" disabled>
				✎
			</IconButton>
		</Inline>
	),
};
