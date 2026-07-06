import type { AnnotationBadgeWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AnnotationBadge } from "./AnnotationBadge";

export const annotationBadgeWidget = defineWidget<AnnotationBadgeWidgetProps>({
	type: "AnnotationBadge",
	module: "context_window.dsl",
	render: (props) => (
		<AnnotationBadge
			className={props.className}
			visualStyle={props.visualStyle}
			label={props.label}
			selected={props.selected}
		/>
	),
});
