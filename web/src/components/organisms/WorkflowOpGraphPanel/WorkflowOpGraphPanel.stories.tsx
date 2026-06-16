import type { Meta, StoryObj } from "@storybook/react-vite";
import { WorkflowOpGraphPanel } from "./WorkflowOpGraphPanel";
import type { WorkflowOpGroup } from "../../../services/api";

const meta = {
	title: "Component Library/Workflows/WorkflowOpGraphPanel",
	component: WorkflowOpGraphPanel,
} satisfies Meta<typeof WorkflowOpGraphPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const groups: WorkflowOpGroup[] = [
	{
		operation: "preprocess_document",
		queue: "rag-eval:preprocess_document",
		status: "succeeded",
		count: 32,
	},
	{ operation: "chunk_document", queue: "rag-eval:chunk_document", status: "running", count: 8 },
	{
		operation: "compute_embeddings",
		queue: "rag-eval:compute_embeddings",
		status: "pending",
		count: 64,
	},
	{ operation: "build_bm25", queue: "rag-eval:build_bm25", status: "ready", count: 1 },
];

export const ActiveGraph: Story = {
	args: {
		groups,
		total: 105,
		done: 32,
	},
};

export const CompleteGraph: Story = {
	args: {
		groups: groups.map((group) => ({ ...group, status: "succeeded" })),
		total: 105,
		done: 105,
	},
};
