import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { DataTable } from "./DataTable";

const meta = {
	title: "Component Library/Molecules/DataTable",
	component: DataTable,
} satisfies Meta<typeof DataTable>;
export default meta;
type Story = StoryObj;

interface Row {
	id: string;
	rank: number;
	title: string;
	score: number;
}
const rows: Row[] = [
	{ id: "a", rank: 1, title: "Fast Growing Trees", score: 0.0387 },
	{ id: "b", rank: 2, title: "Arborvitae Spacing", score: 0.0321 },
];

const columns = [
	{ id: "rank", header: "#", align: "end" as const, cell: (row: Row) => row.rank },
	{ id: "title", header: "Title", cell: (row: Row) => row.title },
	{ id: "score", header: "Score", align: "end" as const, cell: (row: Row) => row.score.toFixed(4) },
];

export const RetrievalRows: Story = {
	render: () => (
		<DataTable rows={rows} getRowKey={(row) => row.id} selectedKey="a" columns={columns} />
	),
};

function MultiSelectionExample() {
	const [activeKey, setActiveKey] = useState("a");
	const [selectedKeys, setSelectedKeys] = useState<string[]>(["b"]);
	return (
		<DataTable
			rows={rows}
			getRowKey={(row) => row.id}
			selectedKey={activeKey}
			onRowSelect={(row) => setActiveKey(row.id)}
			multiSelection={{
				mode: "multi",
				selectedKeys,
				onSelectionChange: setSelectedKeys,
				bulkActions: [
					{
						id: "archive",
						label: "Archive selected",
						danger: true,
						onInvoke: () => setSelectedKeys([]),
					},
				],
			}}
			keyboard={{ mode: "rows", selection: "followFocus", vimAliases: true }}
			columns={columns}
		/>
	);
}

export const MultiSelectionWithActiveRow: Story = {
	render: () => <MultiSelectionExample />,
};
