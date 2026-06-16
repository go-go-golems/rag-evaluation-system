import { DashboardGrid } from "./DashboardGrid";
import { defineWidget } from "../../../widgets/registry";
import type { DashboardGridWidgetProps } from "../../../widgets/ir";

export const dashboardGridWidget = defineWidget<DashboardGridWidgetProps>({
	type: "DashboardGrid",
	module: "ui.dsl",
	render: (props, children) => (
		<DashboardGrid className={props.className} recipe={props.recipe}>
			{children}
		</DashboardGrid>
	),
});
