import type { Meta, StoryObj } from "@storybook/react-vite";
import { useMemo, useState } from "react";
import { cmsArticleFixtures } from "../../../cms";
import { Caption } from "../../foundation";
import { Stack } from "../../layout";
import type { ArticleStatusFilter } from "./ArticleListPanel";
import { ArticleListPanel } from "./ArticleListPanel";

const meta = {
	title: "Component Library/Organisms/ArticleListPanel",
	component: ArticleListPanel,
	args: {
		articles: cmsArticleFixtures,
		onArticleSelect: () => {},
		onRowAction: () => {},
		style: { maxWidth: 960 },
	},
} satisfies Meta<typeof ArticleListPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = {
	args: {
		onCreate: () => {},
		query: "",
		onQueryChange: () => {},
		statusFilter: "all",
		onStatusFilterChange: () => {},
		page: 1,
		pageCount: 2,
		onPageChange: () => {},
	},
};

export const Empty: Story = {
	args: { articles: [] },
};

export const Filtered: Story = {
	args: {
		articles: cmsArticleFixtures.filter((article) => article.status === "draft"),
		statusFilter: "draft",
		onStatusFilterChange: () => {},
	},
};

export const RowSelected: Story = {
	args: { selectedArticleId: cmsArticleFixtures[2]?.id },
};

export const Overflow: Story = {
	args: {
		articles: Array.from({ length: 40 }, (_, index) => ({
			...(cmsArticleFixtures[
				index % cmsArticleFixtures.length
			] as (typeof cmsArticleFixtures)[number]),
			id: `overflow-${index}`,
			title: `Article ${index + 1}`,
			slug: `article-${index + 1}`,
		})),
	},
};

export const Interactive: Story = {
	render: () => {
		const [selectedId, setSelectedId] = useState<string>();
		const [query, setQuery] = useState("");
		const [status, setStatus] = useState<ArticleStatusFilter>("all");
		const [lastAction, setLastAction] = useState("(none)");
		const filtered = useMemo(
			() =>
				cmsArticleFixtures.filter(
					(article) =>
						(status === "all" || article.status === status) &&
						(query === "" || article.title.toLowerCase().includes(query.toLowerCase())),
				),
			[query, status],
		);
		return (
			<Stack gap="sm" style={{ maxWidth: 960 }}>
				<ArticleListPanel
					articles={filtered}
					selectedArticleId={selectedId}
					onArticleSelect={setSelectedId}
					onCreate={() => setLastAction("create")}
					onRowAction={(id, action) => setLastAction(`${action}: ${id}`)}
					query={query}
					onQueryChange={setQuery}
					statusFilter={status}
					onStatusFilterChange={setStatus}
					emptyMessage="No articles match this filter"
				/>
				<Caption>last action: {lastAction}</Caption>
			</Stack>
		);
	},
};
