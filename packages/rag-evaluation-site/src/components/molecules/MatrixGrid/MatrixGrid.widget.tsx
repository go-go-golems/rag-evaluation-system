import type { CSSProperties, ReactNode } from "react";
import { renderCell, rowKey } from "../../../widgets/cellRenderers";
import type {
	CycleCellSpec,
	JsonObject,
	MatrixCellSpec,
	MatrixColumnWidgetSpec,
	MatrixGridWidgetProps,
	MatrixValueSpec,
	ValueCellSpec,
} from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import type { RenderContext } from "../../../widgets/registry";
import { resolveStyleByVars } from "../../../widgets/styleBy";
import { CycleCell } from "../../atoms";
import { type MatrixCellPayload, MatrixGrid } from "./MatrixGrid";

function isCycleSpec(spec: MatrixCellSpec): spec is CycleCellSpec {
	return (spec as CycleCellSpec).kind === "cycle";
}

function isValueSpec(spec: MatrixCellSpec): spec is ValueCellSpec {
	return (spec as ValueCellSpec).kind === "value";
}

function resolveValue(spec: MatrixValueSpec | undefined, row: JsonObject, colId: string): unknown {
	if (!spec) return row[colId];
	if ("mapField" in spec) {
		const map = row[spec.mapField];
		return map && typeof map === "object" ? (map as Record<string, unknown>)[colId] : undefined;
	}
	// template — interpolate ${field} / ${colId} against the row.
	return spec.template.replace(/\$\{([^}]+)\}/g, (_m, path: string) => {
		if (path === "colId") return colId;
		const value = row[path];
		return value == null ? "" : String(value);
	});
}

function renderMatrixCell(
	spec: MatrixCellSpec,
	props: MatrixGridWidgetProps,
	payload: MatrixCellPayload,
	ctx: RenderContext,
): ReactNode {
	if (isCycleSpec(spec)) {
		const styleSet = spec.styleSet ?? props.styleSet;
		const glyphs = spec.glyphs
			? Object.fromEntries(
					Object.entries(spec.glyphs).map(([key, value]) => [key, ctx.renderValue(value)]),
				)
			: undefined;
		return (
			<CycleCell
				value={String(payload.value ?? spec.states[0] ?? "")}
				states={spec.states}
				glyphs={glyphs}
				styleSet={styleSet}
				readOnly={!payload.editable}
				selected={payload.selected}
				onCycle={(next) => payload.onAction({ value: next })}
			/>
		);
	}

	// Value cell: render the resolved (row,col) value; colorBy tints the cell.
	if (isValueSpec(spec)) {
		const content = payload.value == null ? "" : String(payload.value);
		if (props.colorBy) {
			const vars = resolveStyleByVars(
				props.colorBy,
				payload.value,
				payload.row as unknown as JsonObject,
			);
			const style: CSSProperties = {
				...vars,
				display: "block",
				width: "100%",
				height: "100%",
				padding: "4px 6px",
				background: "var(--ctx-fill, transparent)",
				color: "var(--ctx-label, var(--mac-text))",
			};
			return <span style={style}>{content}</span>;
		}
		return content;
	}

	// Otherwise: a DataTable-style CellSpec evaluated against the row.
	return renderCell(spec, payload.row as unknown as JsonObject, ctx.renderNode, (action, context) =>
		ctx.dispatchAction(action, context),
	);
}

export const matrixGridWidget = defineWidget<MatrixGridWidgetProps>({
	type: "MatrixGrid",
	module: "data.dsl",
	render: (props, _children, ctx) => {
		const columns = props.columns.map((column: MatrixColumnWidgetSpec) => ({
			id: column.id,
			header: ctx.renderValue(column.header),
			meta: column.meta,
		}));
		return (
			<MatrixGrid
				className={props.className}
				ariaLabel={props.ariaLabel}
				rows={props.rows}
				columns={columns}
				stickyHeader={props.stickyHeader}
				cornerCell={ctx.renderValue(props.cornerCell)}
				getRowKey={props.getRowKey ? (row) => rowKey(row, props.getRowKey!) : undefined}
				valueAt={(row, col) => resolveValue(props.valueAt, row, col.id)}
				renderRowHeader={
					props.rowHeader
						? (row) => renderCell(props.rowHeader!, row as unknown as JsonObject, ctx.renderNode)
						: undefined
				}
				cells={props.cells?.map((r) => r.map((n) => ctx.renderNode(n)))}
				renderCell={
					props.cells || !props.cell
						? undefined
						: (payload) => renderMatrixCell(props.cell!, props, payload, ctx)
				}
				editableRowKey={props.editableRowKey}
				selectedCell={props.selectedCell}
				footer={
					props.footer
						? {
								header: ctx.renderValue(props.footer.header),
								render: (col) =>
									renderCell(props.footer!.cell, (col.meta ?? {}) as JsonObject, ctx.renderNode),
							}
						: undefined
				}
				onCell={
					props.onCellAction
						? (payload) =>
								ctx.dispatchAction(props.onCellAction!, {
									rowKey: payload.rowKey,
									colId: payload.colId,
									value: payload.value as string | number | boolean | null,
									componentType: "MatrixGrid",
								})
						: undefined
				}
			/>
		);
	},
});
