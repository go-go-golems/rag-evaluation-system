import type { Meta, StoryObj } from "@storybook/react-vite";
import { Text } from "./Text";

const meta = {
	title: "Design System/Foundation/Text",
	component: Text,
	args: {
		children: "Retrieval results use compact, information-dense text roles.",
	},
} satisfies Meta<typeof Text>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Body: Story = {};

export const Roles: Story = {
	render: () => (
		<div style={{ display: "grid", gap: 8 }}>
			<Text size="body">Body text for explanatory panel copy.</Text>
			<Text size="compact" tone="muted">
				Compact muted helper text.
			</Text>
			<Text size="metadata">Metadata text: fixed-1200-150 · openai/text-embedding-3-small</Text>
			<Text size="label" tone="muted">
				Section label
			</Text>
			<Text size="metric" weight="bold">
				204 chunks indexed
			</Text>
		</div>
	),
};

export const Tones: Story = {
	render: () => (
		<div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
			<Text as="span">primary</Text>
			<Text as="span" tone="muted">
				muted
			</Text>
			<Text as="span" tone="accent">
				accent
			</Text>
			<Text as="span" tone="success">
				success
			</Text>
			<Text as="span" tone="warning">
				warning
			</Text>
			<Text as="span" tone="danger">
				danger
			</Text>
		</div>
	),
};
