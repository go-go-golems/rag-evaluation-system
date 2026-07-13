import type { DashboardGridWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { DashboardGrid } from "./DashboardGrid";

export const dashboardGridWidget = defineWidget<DashboardGridWidgetProps>({
	type: "DashboardGrid",
	module: "widget.dsl",
	render: (props, children) => (
		<DashboardGrid className={props.className} recipe={props.recipe}>
			{children}
		</DashboardGrid>
	),
});
