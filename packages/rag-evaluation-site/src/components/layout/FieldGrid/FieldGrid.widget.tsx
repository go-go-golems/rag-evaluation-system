import type { FieldGridWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { FieldGrid } from "./FieldGrid";

export const fieldGridWidget = defineWidget<FieldGridWidgetProps>({
	type: "FieldGrid",
	module: "ui.dsl",
	render: (props, children) => (
		<FieldGrid className={props.className} columns={props.columns} gap={props.gap}>
			{children}
		</FieldGrid>
	),
});
