import type { TranscriptReaderPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TranscriptReaderPanel } from "./TranscriptReaderPanel";

export const transcriptReaderPanelWidget = defineWidget<TranscriptReaderPanelWidgetProps>({
	type: "TranscriptReaderPanel",
	module: "context_window.dsl",
	render: (props, _children, ctx) => (
		<TranscriptReaderPanel
			className={props.className}
			title={props.title}
			subtitle={props.subtitle}
			messages={props.messages}
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
								componentType: "TranscriptReaderPanel",
							})
					: undefined
			}
		/>
	),
});
