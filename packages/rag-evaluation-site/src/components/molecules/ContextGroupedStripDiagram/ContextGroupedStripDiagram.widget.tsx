import type { ContextGroupedStripDiagramWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextGroupedStripDiagram } from "./ContextGroupedStripDiagram";

export const contextGroupedStripDiagramWidget = defineWidget<ContextGroupedStripDiagramWidgetProps>(
	{
		type: "ContextGroupedStripDiagram",
		module: "widget.dsl",
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
