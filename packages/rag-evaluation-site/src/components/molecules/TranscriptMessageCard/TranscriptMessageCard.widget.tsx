import type { TranscriptMessageCardWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TranscriptMessageCard } from "./TranscriptMessageCard";

export const transcriptMessageCardWidget = defineWidget<TranscriptMessageCardWidgetProps>({
	type: "TranscriptMessageCard",
	module: "context_window.dsl",
	render: (props, _children, ctx) => (
		<TranscriptMessageCard
			className={props.className}
			message={props.message}
			annotations={props.annotations}
			selectedAnnotationId={props.selectedAnnotationId}
			showAnnotationChips={props.showAnnotationChips}
			styleSet={props.styleSet}
			onAnnotationSelect={
				props.onAnnotationSelectAction
					? (annotationId) =>
							ctx.dispatchAction(props.onAnnotationSelectAction!, {
								annotationId,
								value: annotationId,
								componentType: "TranscriptMessageCard",
							})
					: undefined
			}
		/>
	),
});
