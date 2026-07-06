import type { Meta, StoryObj } from "@storybook/react-vite";
import { Button } from "../../atoms";
import { Panel } from "../../layout";
import { EmptyState } from "./EmptyState";

const meta = {
	title: "Component Library/Molecules/EmptyState",
	component: EmptyState,
	args: { title: "No assets yet" },
} satisfies Meta<typeof EmptyState>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: { hint: "Uploaded images and files will appear here." },
};

export const WithAction: Story = {
	args: {
		glyph: "▨",
		hint: "Drop files anywhere in the library, or use the upload button.",
		action: <Button variant="primary">⤓ Upload files</Button>,
	},
};

export const Framed: Story = {
	args: { framed: true, hint: "The dashed frame marks a drop-capable region." },
};

export const InsidePanel: Story = {
	render: () => (
		<Panel title="Media">
			<EmptyState glyph="▨" title="No assets match this filter" hint="Try clearing the search." />
		</Panel>
	),
};
