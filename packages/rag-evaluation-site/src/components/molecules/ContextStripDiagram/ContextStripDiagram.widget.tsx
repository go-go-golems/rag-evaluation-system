import type { ContextStripDiagramWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextStripDiagram } from "./ContextStripDiagram";

export const contextStripDiagramWidget = defineWidget<ContextStripDiagramWidgetProps>({
	type: "ContextStripDiagram",
	module: "widget.dsl",
	render: (props) => (
		<ContextStripDiagram
			className={props.className}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			selectedPartId={props.selectedPartId}
			showLabels={props.showLabels}
			showSelection={props.showSelection}
		/>
	),
});
