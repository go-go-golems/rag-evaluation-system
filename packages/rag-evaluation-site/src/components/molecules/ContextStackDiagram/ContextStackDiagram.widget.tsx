import type { ContextStackDiagramWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextStackDiagram } from "./ContextStackDiagram";

export const contextStackDiagramWidget = defineWidget<ContextStackDiagramWidgetProps>({
	type: "ContextStackDiagram",
	module: "context_window.dsl",
	render: (props) => (
		<ContextStackDiagram
			className={props.className}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			selectedPartId={props.selectedPartId}
		/>
	),
});
