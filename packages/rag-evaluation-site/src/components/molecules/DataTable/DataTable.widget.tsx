import { renderCell, rowKey } from "../../../widgets/cellRenderers";
import type { DataTableWidgetProps, JsonObject } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { DataTable } from "./DataTable";

export const dataTableWidget = defineWidget<DataTableWidgetProps>({
	type: "DataTable",
	module: "data.dsl",
	render: (props, _children, ctx) => {
		const rowSelectAction = props.onRowSelect;
		return (
			<DataTable<JsonObject>
				className={props.className}
				rows={props.rows}
				columns={props.columns.map((column) => ({
					id: column.id,
					header: ctx.renderValue(column.header),
					align: column.align,
					maxWidth: column.maxWidth,
					cell: (row) =>
						renderCell(column.cell, row, ctx.renderNode, (action, context) =>
							ctx.dispatchAction(action, context),
						),
				}))}
				getRowKey={(row) => rowKey(row, props.getRowKey)}
				selectedKey={props.selectedKey == null ? props.selectedKey : String(props.selectedKey)}
				emptyMessage={ctx.renderValue(props.emptyMessage)}
				onRowSelect={
					rowSelectAction
						? (row) =>
								ctx.dispatchAction(rowSelectAction, {
									row,
									rowKey: rowKey(row, props.getRowKey),
									componentType: "DataTable",
								})
						: undefined
				}
			/>
		);
	},
});
