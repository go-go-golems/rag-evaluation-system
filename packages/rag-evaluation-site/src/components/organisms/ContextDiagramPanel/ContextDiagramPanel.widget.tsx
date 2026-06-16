import { ContextDiagramPanel } from "./ContextDiagramPanel";
import { defineWidget } from "../../../widgets/registry";
import type { ContextDiagramPanelWidgetProps } from "../../../widgets/ir";

export const contextDiagramPanelWidget = defineWidget<ContextDiagramPanelWidgetProps>({
	type: "ContextDiagramPanel",
	module: "context_window.dsl",
	render: (props) => (
		<ContextDiagramPanel
			className={props.className}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			initialView={props.initialView}
			selectedPartId={props.selectedPartId}
			views={props.views}
			showLegend={props.showLegend}
			showPartDetails={props.showPartDetails}
			chrome={props.chrome}
		/>
	),
});
