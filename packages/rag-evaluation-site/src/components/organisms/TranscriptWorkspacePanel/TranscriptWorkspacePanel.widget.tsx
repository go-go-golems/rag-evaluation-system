import type { TranscriptWorkspacePanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TranscriptWorkspacePanel } from "./TranscriptWorkspacePanel";

export const transcriptWorkspacePanelWidget = defineWidget<TranscriptWorkspacePanelWidgetProps>({
	type: "TranscriptWorkspacePanel",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<TranscriptWorkspacePanel
			className={props.className}
			title={props.title}
			subtitle={props.subtitle}
			messages={props.messages}
			annotations={props.annotations}
			selectedAnnotationId={props.selectedAnnotationId}
			showNotes={props.showNotes}
			styleSet={props.styleSet}
			onAnnotationSelect={
				props.onAnnotationSelectAction
					? (annotationId) =>
							ctx.dispatchAction(props.onAnnotationSelectAction!, {
								annotationId,
								value: annotationId,
								componentType: "TranscriptWorkspacePanel",
							})
					: undefined
			}
		/>
	),
});
