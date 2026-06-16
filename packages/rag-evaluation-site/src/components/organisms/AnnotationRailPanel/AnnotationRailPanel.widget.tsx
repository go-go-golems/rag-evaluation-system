import { AnnotationRailPanel } from "./AnnotationRailPanel";
import { defineWidget } from "../../../widgets/registry";
import type { AnnotationRailPanelWidgetProps } from "../../../widgets/ir";

export const annotationRailPanelWidget = defineWidget<AnnotationRailPanelWidgetProps>({
	type: "AnnotationRailPanel",
	module: "context_window.dsl",
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
