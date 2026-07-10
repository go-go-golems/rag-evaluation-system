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

const manyItems: AppNavItem[] = [
	"Index",
	"Simple Table",
	"Selectable Table",
	"Master Detail Editor",
	"Row Actions",
	"All Modules Gallery",
	"Admin Course CMS",
	"Handouts And Slide",
	"Markdown Article",
	"Markdown Editor",
	"Course Landing",
	"Course Slide Deck",
	"Course Handouts",
	"Course Shell Layout",
	"Context Budget Diagram",
	"Context Transcript Workspace",
	"Schedule Poll Editable",
	"Schedule Poll Summary",
	"Booking Picker",
	"Time Month",
	"Time Week",
	"CRM Board",
	"Activity Feed",
	"Page Chrome",
	"Matrix Heatmap",
].map((label) => ({ id: label.toLowerCase().replace(/ /g, "-"), label }));

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

export const OverflowManyItems: Story = {
	args: {
		brand: "Widget DSL v3 examples",
		items: manyItems,
		activeItemId: "matrix-heatmap",
		onItemSelect: () => undefined,
	},
	render: () => {
		const [activeItemId, setActiveItemId] = useState("matrix-heatmap");
		return (
			<div style={{ width: "100vw", overflow: "hidden", border: "1px solid var(--mac-border)" }}>
				<AppNav
					brand="Widget DSL v3 examples"
					items={manyItems}
					activeItemId={activeItemId}
					onItemSelect={setActiveItemId}
				/>
			</div>
		);
	},
};
