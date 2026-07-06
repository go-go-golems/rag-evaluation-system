import type { ContextStudioNavIconWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContextStudioNavIcon } from "./ContextStudioNavIcon";

export const contextStudioNavIconWidget = defineWidget<ContextStudioNavIconWidgetProps>({
	type: "ContextStudioNavIcon",
	module: "course.dsl",
	render: (props) => (
		<ContextStudioNavIcon className={props.className} id={props.id} title={props.title} />
	),
});
