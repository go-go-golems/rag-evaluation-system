import type { RagStatus } from "@go-go-golems/rag-evaluation-site";

const STATUS_ICON: Record<string, string> = {
	pending: "◌",
	ready: "◌",
	running: "●",
	succeeded: "✔",
	failed: "✘",
	canceled: "⊘",
};

const STATUS_CLASS: Record<string, string> = {
	pending: "status-pending",
	ready: "status-pending",
	running: "status-running",
	succeeded: "status-done",
	failed: "status-error",
	canceled: "status-canceled",
};

export function statusIcon(status: string) {
	return STATUS_ICON[status] ?? "?";
}
export function statusClass(status: string) {
	return STATUS_CLASS[status] ?? "";
}

export function workflowStatusToRagStatus(status: string): RagStatus {
	if (status === "succeeded") return "succeeded";
	if (status === "failed") return "failed";
	if (status === "canceled") return "canceled";
	if (status === "running") return "running";
	if (status === "ready") return "ready";
	return "pending";
}

export function friendlyOpName(op: string): string {
	return op.replace(/_/g, " ");
}

export function timeAgo(iso: string): string {
	const ms = Date.now() - new Date(iso).getTime();
	if (ms < 0) return "now";
	if (ms < 60000) return `${Math.floor(ms / 1000)}s`;
	if (ms < 3600000) return `${Math.floor(ms / 60000)}m`;
	if (ms < 86400000) return `${Math.floor(ms / 3600000)}h`;
	return `${Math.floor(ms / 86400000)}d`;
}
