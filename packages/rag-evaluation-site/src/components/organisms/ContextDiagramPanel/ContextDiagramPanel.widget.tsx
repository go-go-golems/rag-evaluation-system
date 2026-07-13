import type { ContextDiagramPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextDiagramPanel } from "./ContextDiagramPanel";

export const contextDiagramPanelWidget = defineWidget<ContextDiagramPanelWidgetProps>({
	type: "ContextDiagramPanel",
	module: "widget.dsl",
	render: (props, _children, ctx) => {
		const onPartSelectAction = props.onPartSelectAction;
		return (
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
				onPartSelect={
					onPartSelectAction
						? (partId) =>
								ctx.dispatchAction(onPartSelectAction, {
									partId,
									value: partId,
									componentType: "ContextDiagramPanel",
								})
						: undefined
				}
			/>
		);
	},
});
