import type { CaptionWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Caption } from "./Caption";

export const captionWidget = defineWidget<CaptionWidgetProps>({
	type: "Caption",
	module: "widget.dsl",
	render: (props, children) => (
		<Caption
			className={props.className}
			tone={props.tone}
			transform={props.transform}
			truncate={props.truncate}
		>
			{children}
		</Caption>
	),
});
