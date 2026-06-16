import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { AppNav, type AppNavItem } from "./AppNav";

const items: AppNavItem[] = [
	{ id: "search", label: "Search" },
	{ id: "corpus", label: "Corpus" },
	{ id: "workflows", label: "Workflows" },
	{ id: "pipeline", label: "Pipeline" },
	{ id: "embeddings", label: "Embeddings" },
	{ id: "evaluation", label: "Evaluation" },
];

const meta = {
	title: "Component Library/Molecules/AppNav",
	component: AppNav,
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof AppNav>;

export default meta;
type Story = StoryObj<typeof meta>;

export const SearchActive: Story = {
	args: {
		brand: "◉ RAG Eval",
		items,
		activeItemId: "search",
		onItemSelect: () => undefined,
	},
};

export const Interactive: Story = {
	args: {
		brand: "◉ RAG Eval",
		items,
		activeItemId: "workflows",
		onItemSelect: () => undefined,
	},
	render: () => {
		const [activeItemId, setActiveItemId] = useState("workflows");
		return (
			<AppNav
				brand="◉ RAG Eval"
				items={items}
				activeItemId={activeItemId}
				onItemSelect={setActiveItemId}
			/>
		);
	},
};
