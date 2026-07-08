import type { HTMLAttributes, ReactNode } from "react";
import styles from "./MatrixGrid.module.css";

export type MatrixRow = Record<string, unknown>;

export interface MatrixColumnSpec {
	id: string;
	header: ReactNode;
	meta?: Record<string, unknown>;
}

/**
 * The stable payload every cell receives. This is the seam that keeps the grid
 * domain-blind: the grid owns geometry + the `onAction` plumbing; the cell owns
 * everything visual and semantic. Any component honoring this shape is a valid
 * cell (availability toggle, calendar day, rating, avatar, ...).
 */
export interface MatrixCellPayload<Row = MatrixRow> {
	row: Row;
	col: MatrixColumnSpec;
	value: unknown;
	rowKey: string;
	rowIndex: number;
	colIndex: number;
	selected: boolean;
	editable: boolean;
	/** Notify the grid a cell changed. `value` in `extra` overrides the resolved value. */
	onAction: (extra?: { value?: unknown } & Record<string, unknown>) => void;
}

export interface MatrixGridFooterSpec {
	/** Content for the bottom-left corner cell of the footer row. */
	header?: ReactNode;
	render: (col: MatrixColumnSpec, colIndex: number) => ReactNode;
}

export interface MatrixGridProps<Row = MatrixRow>
	extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	rows: Row[];
	columns: MatrixColumnSpec[];
	getRowKey?: (row: Row, index: number) => string;
	renderRowHeader?: (row: Row, index: number) => ReactNode;
	/** Accessor: value at (row, col). Defaults to `row[col.id]`. */
	valueAt?: (row: Row, col: MatrixColumnSpec) => unknown;
	/** Mode A — one renderer applied per (row, col). */
	renderCell?: (payload: MatrixCellPayload<Row>) => ReactNode;
	/** Mode B — an explicit matrix of prebuilt nodes (`cells[rowIndex][colIndex]`). */
	cells?: ReactNode[][];
	footer?: MatrixGridFooterSpec;
	selectedCell?: { rowKey: string; colId: string } | null;
	editableRowKey?: string;
	stickyHeader?: boolean;
	/** Content for the top-left corner cell. */
	cornerCell?: ReactNode;
	onCell?: (payload: { rowKey: string; colId: string; value: unknown }) => void;
	ariaLabel?: string;
}

function defaultRowKey(row: unknown, index: number): string {
	const id = (row as MatrixRow).id;
	return typeof id === "string" || typeof id === "number" ? String(id) : String(index);
}

export function MatrixGrid<Row = MatrixRow>({
	rows,
	columns,
	getRowKey = defaultRowKey,
	renderRowHeader,
	valueAt = (row, col) => (row as MatrixRow)[col.id],
	renderCell,
	cells,
	footer,
	selectedCell,
	editableRowKey,
	stickyHeader = true,
	cornerCell,
	onCell,
	ariaLabel,
	className,
	...rest
}: MatrixGridProps<Row>) {
	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="MatrixGrid"
			{...rest}
		>
			<table className={styles.table} aria-label={ariaLabel}>
				<thead className={stickyHeader ? styles.stickyHead : undefined}>
					<tr>
						<th className={[styles.corner, styles.rowHeadCell].join(" ")} scope="col">
							{cornerCell}
						</th>
						{columns.map((col) => (
							<th key={col.id} className={styles.colHeadCell} scope="col">
								{col.header}
							</th>
						))}
					</tr>
				</thead>
				<tbody>
					{rows.map((row, rowIndex) => {
						const rowKey = getRowKey(row, rowIndex);
						const editable = editableRowKey != null && editableRowKey === rowKey;
						return (
							<tr key={rowKey} data-editable={editable || undefined}>
								<th className={styles.rowHeadCell} scope="row">
									{renderRowHeader ? renderRowHeader(row, rowIndex) : rowKey}
								</th>
								{columns.map((col, colIndex) => {
									const value = valueAt(row, col);
									const selected =
										selectedCell != null &&
										selectedCell.rowKey === rowKey &&
										selectedCell.colId === col.id;
									const explicit = cells?.[rowIndex]?.[colIndex];
									const content =
										explicit !== undefined
											? explicit
											: renderCell?.({
													row,
													col,
													value,
													rowKey,
													rowIndex,
													colIndex,
													selected,
													editable,
													onAction: (extra) =>
														onCell?.({
															rowKey,
															colId: col.id,
															value: extra && "value" in extra ? extra.value : value,
														}),
												});
									return (
										<td key={col.id} className={styles.cell} data-selected={selected || undefined}>
											{content}
										</td>
									);
								})}
							</tr>
						);
					})}
				</tbody>
				{footer ? (
					<tfoot className={styles.foot}>
						<tr>
							<th className={styles.rowHeadCell} scope="row">
								{footer.header}
							</th>
							{columns.map((col, colIndex) => (
								<td key={col.id} className={styles.footCell}>
									{footer.render(col, colIndex)}
								</td>
							))}
						</tr>
					</tfoot>
				) : null}
			</table>
		</div>
	);
}
