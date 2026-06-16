import { SelectInput } from "@go-go-golems/rag-evaluation-site";
import { Caption, StatusText } from "@go-go-golems/rag-evaluation-site";
import { Panel, ScrollRegion } from "@go-go-golems/rag-evaluation-site";
import { DataTable, type DataTableColumn } from "@go-go-golems/rag-evaluation-site";
import type { WorkflowListItem } from "../../../services/api";
import { timeAgo, workflowStatusToRagStatus } from "../../workflows/workflowFormat";

export interface WorkflowListPanelProps {
	workflows: WorkflowListItem[];
	total: number;
	isLoading?: boolean;
	statusFilter: string;
	onStatusFilterChange: (status: string) => void;
	onSelect: (workflowId: string) => void;
}

export function WorkflowListPanel({
	workflows,
	total,
	isLoading = false,
	statusFilter,
	onStatusFilterChange,
	onSelect,
}: WorkflowListPanelProps) {
	const columns: DataTableColumn<WorkflowListItem>[] = [
		{
			id: "status",
			header: "Status",
			cell: (item) => (
				<StatusText status={workflowStatusToRagStatus(item.workflow.Status)} icon>
					{item.workflow.Status}
				</StatusText>
			),
		},
		{ id: "id", header: "Workflow ID", maxWidth: 250, cell: (item) => item.workflow.ID },
		{
			id: "strategy",
			header: "Strategy",
			cell: (item) => <Caption>{item.workflow.Metadata?.strategy ?? "—"}</Caption>,
		},
		{ id: "ops", header: "Ops", align: "end", cell: (item) => `${item.opDone}/${item.opTotal}` },
		{
			id: "age",
			header: "Age",
			cell: (item) => <Caption>{timeAgo(item.workflow.CreatedAt)}</Caption>,
		},
	];

	return (
		<Panel
			fill
			title={`Workflows (${total})`}
			actions={
				<SelectInput
					value={statusFilter}
					onChange={(event) => onStatusFilterChange(event.target.value)}
				>
					<option value="">all</option>
					<option value="pending">pending</option>
					<option value="running">running</option>
					<option value="succeeded">succeeded</option>
					<option value="failed">failed</option>
					<option value="canceled">canceled</option>
				</SelectInput>
			}
			density="condensed"
			data-rag-component="WorkflowListPanel"
		>
			<ScrollRegion axis="y" style={{ maxHeight: 400 }}>
				{isLoading ? (
					<Caption>Loading...</Caption>
				) : workflows.length === 0 ? (
					<Caption>
						No workflows{statusFilter ? ` with status "${statusFilter}"` : " yet"}
						{!statusFilter && ". Submit one!"}
					</Caption>
				) : (
					<DataTable
						rows={workflows}
						columns={columns}
						getRowKey={(item) => item.workflow.ID}
						onRowSelect={(item) => onSelect(item.workflow.ID)}
					/>
				)}
			</ScrollRegion>
		</Panel>
	);
}
