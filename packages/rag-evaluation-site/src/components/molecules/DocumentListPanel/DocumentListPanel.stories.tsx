import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { contextHandoutFixture } from "../../../context";
import { DocumentListPanel } from "./DocumentListPanel";

const items = contextHandoutFixture.docs.map((document) => ({
	id: document.id,
	title: document.title,
	format: document.format,
	size: document.size,
	description: document.description,
}));

const meta = {
	title: "Component Library/Molecules/DocumentListPanel",
	component: DocumentListPanel,
	args: {
		title: "Handout",
		description: contextHandoutFixture.intro,
		items,
		selectedItemId: items[0]?.id,
	},
} satisfies Meta<typeof DocumentListPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Documents: Story = {};

export const Interactive: Story = {
	render: () => {
		const [selected, setSelected] = useState(items[0]?.id);
		return (
			<DocumentListPanel
				title="Handout"
				description={contextHandoutFixture.intro}
				items={items}
				selectedItemId={selected}
				onItemSelect={setSelected}
				onDownloadAll={() => {}}
			/>
		);
	},
};

export const Empty: Story = {
	args: { items: [], selectedItemId: undefined },
};
