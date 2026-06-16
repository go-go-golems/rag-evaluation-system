import { Caption } from "./Caption";
import { defineWidget } from "../../../widgets/registry";
import type { CaptionWidgetProps } from "../../../widgets/ir";

export const captionWidget = defineWidget<CaptionWidgetProps>({
	type: "Caption",
	module: "ui.dsl",
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
