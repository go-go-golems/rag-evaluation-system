import type { Meta, StoryObj } from "@storybook/react-vite";
import { Caption } from "./Caption";

const meta = {
	title: "Design System/Foundation/Caption",
	component: Caption,
	args: { children: "retriever: hybrid · rank #1", tone: "muted" },
} satisfies Meta<typeof Caption>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
export const Tones: Story = {
	render: () => (
		<div style={{ display: "grid", gap: 6 }}>
			<Caption tone="muted">Muted metadata</Caption>
			<Caption tone="accent">Vector score</Caption>
			<Caption tone="success">Embedded</Caption>
			<Caption tone="warning">Sparse coverage</Caption>
			<Caption tone="danger">Missing embedding</Caption>
		</div>
	),
};
