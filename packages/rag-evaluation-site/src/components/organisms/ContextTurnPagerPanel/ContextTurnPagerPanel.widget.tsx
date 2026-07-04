import type { ContextTurnPagerPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextTurnPagerPanel } from "./ContextTurnPagerPanel";

export const contextTurnPagerPanelWidget = defineWidget<ContextTurnPagerPanelWidgetProps>({
	type: "ContextTurnPagerPanel",
	module: "context_window.dsl",
	render: (props) => (
		<ContextTurnPagerPanel
			className={props.className}
			snapshots={props.snapshots}
			styleSet={props.styleSet}
			initialSnapshotId={props.initialSnapshotId}
			selectedPartId={props.selectedPartId}
			diagram={props.diagram}
			groupBy={props.groupBy}
			mode={props.mode}
			includeGlobalParts={props.includeGlobalParts}
			showLegend={props.showLegend}
			title={props.title}
		/>
	),
});
