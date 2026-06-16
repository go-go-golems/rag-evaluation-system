import { ContextBudgetBar } from "./ContextBudgetBar";
import { defineWidget } from "../../../widgets/registry";
import type { ContextBudgetBarWidgetProps } from "../../../widgets/ir";

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
