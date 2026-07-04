import type { AnnotationNoteCardWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AnnotationNoteCard } from "./AnnotationNoteCard";

export const annotationNoteCardWidget = defineWidget<AnnotationNoteCardWidgetProps>({
	type: "AnnotationNoteCard",
	module: "context_window.dsl",
	render: (props) => (
		<AnnotationNoteCard
			className={props.className}
			annotation={props.annotation}
			styleSet={props.styleSet}
			selected={props.selected}
			index={props.index}
		/>
	),
});
