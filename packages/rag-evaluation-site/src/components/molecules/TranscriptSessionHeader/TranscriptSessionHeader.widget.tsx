import type { TranscriptSessionHeaderWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TranscriptSessionHeader } from "./TranscriptSessionHeader";

export const transcriptSessionHeaderWidget = defineWidget<TranscriptSessionHeaderWidgetProps>({
	type: "TranscriptSessionHeader",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<TranscriptSessionHeader
			className={props.className}
			title={ctx.renderValue(props.title)}
			subtitle={ctx.renderValue(props.subtitle)}
			messageCount={props.messageCount}
			annotationCount={props.annotationCount}
			tokenTotal={props.tokenTotal}
			rightSlot={props.rightSlot ? ctx.renderNode(props.rightSlot) : undefined}
		/>
	),
});
