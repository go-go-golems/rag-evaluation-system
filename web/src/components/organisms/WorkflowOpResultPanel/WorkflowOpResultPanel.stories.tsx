import type { Meta, StoryObj } from "@storybook/react-vite";
import { WorkflowOpResultPanel } from "./WorkflowOpResultPanel";
import type { OpResult } from "../../../services/api";

const meta = {
	title: "Component Library/Workflows/WorkflowOpResultPanel",
	component: WorkflowOpResultPanel,
} satisfies Meta<typeof WorkflowOpResultPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const result: OpResult = {
	OpID: "op_compute_embeddings_001",
	Data: { embedded_chunks: 128, provider: "openai", model: "text-embedding-3-small" },
	Records: [
		{ Table: "chunk_embeddings", PK: "chunk_privacy_001", Data: {} },
		{ Table: "chunk_embeddings", PK: "chunk_privacy_002", Data: {} },
	],
	Artifacts: [
		{ Name: "embedding-summary.json", Kind: "json", ContentType: "application/json", Body: "{}" },
	],
	Emitted: [{ ID: "op_build_bm25_001", Kind: "build_bm25", Queue: "rag-eval:build_bm25" }],
	EmittedIDs: ["op_build_bm25_001", "op_quality_check_001"],
	Error: null,
	CompletedAt: "2026-06-01T12:34:56Z",
};

export const Complete: Story = { args: { result } };
export const Loading: Story = { args: { isLoading: true } };
export const NoResult: Story = { args: { error: { status: 404 } } };
export const ErrorResult: Story = {
	args: {
		result: {
			...result,
			Error: { Code: "RATE_LIMIT", Message: "Provider returned 429", Retryable: true },
		},
	},
};
