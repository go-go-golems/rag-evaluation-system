import type { Meta, StoryObj } from "@storybook/react-vite";
import { PipelineOverview } from "./PipelineOverview";
import type { Document, Source } from "../../services/api";

const sources: Source[] = [
	{
		id: "src-tree-guides",
		name: "Tree Center Guides",
		type: "web_crawl",
		created_at: "2026-05-21T09:12:00Z",
		updated_at: "2026-05-21T09:12:00Z",
	},
	{
		id: "src-product-pages",
		name: "Product Pages",
		type: "sitemap",
		created_at: "2026-05-22T11:30:00Z",
		updated_at: "2026-05-22T11:30:00Z",
	},
];

const documents: Document[] = [
	{
		id: "doc-privacy-trees",
		source_id: "src-tree-guides",
		title: "Fast Growing Trees for Privacy Screens",
		author: "Tree Center Editorial",
		content_type: "text/html",
		word_count: 1824,
		language: "en",
		status: "indexed",
		created_at: "2026-05-21T09:20:00Z",
		updated_at: "2026-05-21T09:22:00Z",
	},
	{
		id: "doc-hydrangea-pruning",
		source_id: "src-tree-guides",
		title: "Hydrangea Pruning by Bloom Type",
		author: "Tree Center Editorial",
		content_type: "text/html",
		word_count: 1420,
		language: "en",
		status: "chunked",
		created_at: "2026-05-21T10:15:00Z",
		updated_at: "2026-05-21T10:17:00Z",
	},
];

const meta = {
	title: "Component Library/Pipeline/PipelineOverview",
	component: PipelineOverview,
	args: {
		sources,
		documents,
	},
} satisfies Meta<typeof PipelineOverview>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Populated: Story = {};

export const Loading: Story = {
	args: {
		sources: [],
		documents: [],
		sourcesLoading: true,
		documentsLoading: true,
	},
};

export const Empty: Story = {
	args: {
		sources: [],
		documents: [],
	},
};
