import type { ContextTreemapWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextTreemap } from "./ContextTreemap";

export const contextTreemapWidget = defineWidget<ContextTreemapWidgetProps>({
	type: "ContextTreemap",
	module: "context_window.dsl",
	render: (props) => (
		<ContextTreemap
			className={props.className}
			snapshot={props.snapshot}
			styleSet={props.styleSet}
			selectedPartId={props.selectedPartId}
		/>
	),
});
