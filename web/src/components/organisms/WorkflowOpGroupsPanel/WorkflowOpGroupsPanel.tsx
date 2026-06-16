import { Caption, StatusText } from "@go-go-golems/rag-evaluation-site";
import { Panel } from "@go-go-golems/rag-evaluation-site";
import { DataTable, type DataTableColumn } from "@go-go-golems/rag-evaluation-site";
import type { WorkflowOpGroup } from "../../../services/api";
import { friendlyOpName, workflowStatusToRagStatus } from "../../workflows/workflowFormat";

export interface WorkflowOpGroupsPanelProps {
	groups: WorkflowOpGroup[];
	totalOps: number;
	selectedKey: string | null;
	onInspect: (group: WorkflowOpGroup) => void;
}

export function workflowGroupKey(group: WorkflowOpGroup) {
	return `${group.operation}|${group.status}`;
}

export function WorkflowOpGroupsPanel({
	groups,
	totalOps,
	selectedKey,
	onInspect,
}: WorkflowOpGroupsPanelProps) {
	const columns: DataTableColumn<WorkflowOpGroup>[] = [
		{
			id: "status",
			header: "Status",
			cell: (group) => (
				<StatusText status={workflowStatusToRagStatus(group.status)} icon>
					{group.status}
				</StatusText>
			),
		},
		{ id: "operation", header: "Operation", cell: (group) => friendlyOpName(group.operation) },
		{ id: "queue", header: "Queue", cell: (group) => group.queue.replace("rag-eval:", "") },
		{ id: "count", header: "Count", align: "end", cell: (group) => group.count },
	];

	return (
		<Panel
			title={`Ops by Group (${totalOps} total)`}
			density="condensed"
			data-rag-component="WorkflowOpGroupsPanel"
		>
			{groups.length === 0 ? (
				<Caption>No ops yet</Caption>
			) : (
				<DataTable
					rows={groups}
					columns={columns}
					getRowKey={workflowGroupKey}
					selectedKey={selectedKey}
					onRowSelect={onInspect}
				/>
			)}
		</Panel>
	);
}
