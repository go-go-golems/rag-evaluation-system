import type { ContextBudgetBarWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextBudgetBar } from "./ContextBudgetBar";

export const contextBudgetBarWidget = defineWidget<ContextBudgetBarWidgetProps>({
	type: "ContextBudgetBar",
	module: "context_window.dsl",
	render: (props) => (
		<ContextBudgetBar
			className={props.className}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			showLegend={props.showLegend}
			selectedPartId={props.selectedPartId}
		/>
	),
});
