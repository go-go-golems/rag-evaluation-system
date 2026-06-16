import { ContextGroupedStripDiagram } from "./ContextGroupedStripDiagram";
import { defineWidget } from "../../../widgets/registry";
import type { ContextGroupedStripDiagramWidgetProps } from "../../../widgets/ir";

export const contextGroupedStripDiagramWidget = defineWidget<ContextGroupedStripDiagramWidgetProps>(
	{
		type: "ContextGroupedStripDiagram",
		module: "context_window.dsl",
		render: (props) => (
			<ContextGroupedStripDiagram
				className={props.className}
				snapshot={props.snapshot}
				styleSet={props.styleSet}
				selectedPartId={props.selectedPartId}
				groupBy={props.groupBy}
				showGroupLabels={props.showGroupLabels}
				showPartLabels={props.showPartLabels}
				showSelection={props.showSelection}
			/>
		),
	},
);
