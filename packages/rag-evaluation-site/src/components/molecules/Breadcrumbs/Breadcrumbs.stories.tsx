import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import { Breadcrumbs } from "./Breadcrumbs";

const meta = {
	title: "Component Library/Molecules/Breadcrumbs",
	component: Breadcrumbs,
	args: {
		items: [
			{ id: "media", label: "Media" },
			{ id: "course", label: "Course assets" },
			{ id: "diagrams", label: "Diagrams" },
		],
		onNavigate: () => {},
	},
} satisfies Meta<typeof Breadcrumbs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Deep: Story = {
	args: {
		items: [
			{ id: "1", label: "Media" },
			{ id: "2", label: "2026" },
			{ id: "3", label: "Workshops" },
			{ id: "4", label: "Context engineering" },
			{ id: "5", label: "Screenshots with a very long folder name" },
			{ id: "6", label: "Final picks" },
		],
	},
};

export const SingleItem: Story = {
	args: { items: [{ id: "media", label: "Media" }] },
};

export const Interactive: Story = {
	render: () => {
		const all = [
			{ id: "media", label: "Media" },
			{ id: "course", label: "Course assets" },
			{ id: "diagrams", label: "Diagrams" },
		];
		const [depth, setDepth] = useState(3);
		return (
			<Stack gap="sm">
				<Breadcrumbs
					items={all.slice(0, depth)}
					onNavigate={(id) => setDepth(all.findIndex((item) => item.id === id) + 1)}
				/>
				<Caption>Click a crumb to truncate the path.</Caption>
			</Stack>
		);
	},
};
