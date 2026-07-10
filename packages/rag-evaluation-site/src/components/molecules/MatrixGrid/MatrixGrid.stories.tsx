import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import type { ContextStyleSet } from "../../../context";
import { CycleCell, DateTile, RatioBadge, Tag } from "../../atoms";
import { Stack } from "../../layout";
import { type MatrixColumnSpec, MatrixGrid } from "./MatrixGrid";

// ─── shared demo palette + data ─────────────────────────────────────────────

const AVAILABILITY_STATES = ["yes", "ifneedbe", "no", "unknown"];
const AVAILABILITY_GLYPHS: Record<string, string> = {
	yes: "✓",
	ifneedbe: "~",
	no: "✕",
	unknown: "·",
};
const availabilityStyleSet: ContextStyleSet = {
	id: "availability",
	styles: {
		yes: { fill: "var(--mac-green)", labelColor: "var(--mac-text-inv)" },
		ifneedbe: { fill: "var(--mac-amber)", labelColor: "var(--mac-text)" },
		no: { fill: "var(--mac-accent-2)", labelColor: "var(--mac-text-inv)" },
		unknown: { fill: "var(--mac-surface)", labelColor: "var(--mac-text-dim)" },
	},
	legend: [],
};

interface Slot {
	id: string;
	dateISO: string;
	time: string;
}
interface Respondent {
	id: string;
	name: string;
	cells: Record<string, string>;
}

const SLOTS: Slot[] = [
	{ id: "s1", dateISO: "2026-07-09", time: "14:00" },
	{ id: "s2", dateISO: "2026-07-10", time: "10:00" },
	{ id: "s3", dateISO: "2026-07-10", time: "16:00" },
	{ id: "s4", dateISO: "2026-07-11", time: "09:00" },
];

const columns: MatrixColumnSpec[] = SLOTS.map((slot) => ({
	id: slot.id,
	header: (
		<Stack gap="xs" align="center">
			<DateTile dateISO={slot.dateISO} size="sm" />
			<span style={{ font: "var(--rag-font-role-metadata)" }}>{slot.time}</span>
		</Stack>
	),
	meta: { slot },
}));

const RESPONDENTS: Respondent[] = [
	{ id: "alice", name: "Alice", cells: { s1: "yes", s2: "ifneedbe", s3: "no", s4: "yes" } },
	{ id: "bob", name: "Bob", cells: { s1: "yes", s2: "yes", s3: "no", s4: "unknown" } },
	{ id: "chen", name: "Chen", cells: { s1: "yes", s2: "yes", s3: "yes", s4: "yes" } },
];

function tallyYes(rows: Respondent[], colId: string): number {
	return rows.filter((r) => r.cells[colId] === "yes").length;
}

const meta = {
	title: "Component Library/Molecules/MatrixGrid",
	component: MatrixGrid,
	args: { rows: [], columns: [] },
} satisfies Meta<typeof MatrixGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

// ─── Mode A: renderCell + CycleCell (the Doodle poll) ───────────────────────

export const PollWithCycleCells: Story = {
	render: () => {
		const [rows, setRows] = useState<Respondent[]>([
			...RESPONDENTS,
			{ id: "you", name: "You", cells: { s1: "yes", s2: "ifneedbe", s3: "no", s4: "unknown" } },
		]);
		const setCell = (rowKey: string, colId: string, value: string) =>
			setRows((prev) =>
				prev.map((r) => (r.id === rowKey ? { ...r, cells: { ...r.cells, [colId]: value } } : r)),
			);
		return (
			<MatrixGrid<Respondent>
				ariaLabel="Team sync availability"
				rows={rows}
				columns={columns}
				valueAt={(row, col) => row.cells[col.id] ?? "unknown"}
				renderRowHeader={(row) => row.name}
				editableRowKey="you"
				renderCell={(p) => (
					<CycleCell
						value={String(p.value)}
						states={AVAILABILITY_STATES}
						glyphs={AVAILABILITY_GLYPHS}
						styleSet={availabilityStyleSet}
						readOnly={!p.editable}
						onCycle={(next) => p.onAction({ value: next })}
					/>
				)}
				onCell={({ rowKey, colId, value }) => setCell(rowKey, colId, String(value))}
				footer={{
					header: "yes",
					render: (col) => <RatioBadge count={tallyYes(rows, col.id)} total={rows.length} />,
				}}
			/>
		);
	},
};

// ─── Mode B: explicit cells matrix (domain-blind reuse) ─────────────────────

export const ModeBExplicitCells: Story = {
	render: () => {
		const featureColumns: MatrixColumnSpec[] = [
			{ id: "free", header: "Free" },
			{ id: "pro", header: "Pro" },
			{ id: "team", header: "Team" },
		];
		const featureRows = [
			{ id: "sso", name: "SSO" },
			{ id: "api", name: "API access" },
			{ id: "support", name: "Priority support" },
		];
		const yes = <Tag label="✓" />;
		const cells = [
			[null, yes, yes],
			[null, yes, yes],
			[null, null, yes],
		];
		return (
			<MatrixGrid
				ariaLabel="Plan comparison"
				rows={featureRows}
				columns={featureColumns}
				renderRowHeader={(row) => String((row as { name: string }).name)}
				cells={cells}
				cornerCell="Feature"
			/>
		);
	},
};

export const SelectedCell: Story = {
	render: () => (
		<MatrixGrid<Respondent>
			rows={RESPONDENTS}
			columns={columns}
			valueAt={(row, col) => row.cells[col.id] ?? "unknown"}
			renderRowHeader={(row) => row.name}
			selectedCell={{ rowKey: "chen", colId: "s3" }}
			renderCell={(p) => (
				<CycleCell
					value={String(p.value)}
					states={AVAILABILITY_STATES}
					glyphs={AVAILABILITY_GLYPHS}
					styleSet={availabilityStyleSet}
					readOnly
				/>
			)}
		/>
	),
};

export const Empty: Story = {
	render: () => <MatrixGrid rows={[]} columns={columns} cornerCell="—" />,
};
