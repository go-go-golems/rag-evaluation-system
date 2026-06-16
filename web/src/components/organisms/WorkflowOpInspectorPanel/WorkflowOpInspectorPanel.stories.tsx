import type { Meta, StoryObj } from "@storybook/react-vite";
import { WorkflowOpInspectorPanel } from "./WorkflowOpInspectorPanel";
import type { WorkflowOp } from "../../../services/api";

const meta = {
	title: "Component Library/Workflows/WorkflowOpInspectorPanel",
	component: WorkflowOpInspectorPanel,
} satisfies Meta<typeof WorkflowOpInspectorPanel>;

export default meta;
type Story = StoryObj<typeof meta>;

const sample: WorkflowOp = {
	op: {
		ID: "op_compute_embeddings_001",
		WorkflowID: "wf_intake_001",
		Kind: "compute_embeddings",
		Queue: "rag-eval:compute_embeddings",
		DedupKey: "embed:doc_privacy_trees:fixed-1200-150",
		Input: { document_id: "doc_privacy_trees", strategy_id: "fixed-1200-150", batch_size: 16 },
		DependsOn: [{ OpID: "op_chunk_document_001", Required: true }],
		Retry: { MaxAttempts: 3, BackoffKind: "exponential", InitialBackoff: 30 },
		RetryState: { Attempt: 1, NextAttemptAt: null, LastError: "" },
		Metadata: { source_id: "ttc-guides" },
	},
	status: "succeeded",
	createdAt: new Date(Date.now() - 10 * 60_000).toISOString(),
	updatedAt: new Date().toISOString(),
};

export const Succeeded: Story = {
	args: {
		sample,
		onClose: () => undefined,
	},
};

export const FailedRetryable: Story = {
	args: {
		sample: {
			...sample,
			status: "failed",
			op: {
				...sample.op,
				RetryState: {
					Attempt: 2,
					NextAttemptAt: null,
					LastError: "embedding provider returned 429 rate limit",
				},
			},
		},
		onClose: () => undefined,
		onRetry: () => undefined,
	},
};
