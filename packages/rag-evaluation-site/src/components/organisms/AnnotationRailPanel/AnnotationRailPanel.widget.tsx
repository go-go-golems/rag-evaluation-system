import type { AnnotationRailPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AnnotationRailPanel } from "./AnnotationRailPanel";

export const annotationRailPanelWidget = defineWidget<AnnotationRailPanelWidgetProps>({
	type: "AnnotationRailPanel",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<AnnotationRailPanel
			className={props.className}
			title={props.title}
			description={props.description}
			annotations={props.annotations}
			selectedAnnotationId={props.selectedAnnotationId}
			styleSet={props.styleSet}
			onAnnotationSelect={
				props.onAnnotationSelectAction
					? (annotationId) =>
							ctx.dispatchAction(props.onAnnotationSelectAction!, {
								annotationId,
								value: annotationId,
								componentType: "AnnotationRailPanel",
							})
					: undefined
			}
		/>
	),
});
