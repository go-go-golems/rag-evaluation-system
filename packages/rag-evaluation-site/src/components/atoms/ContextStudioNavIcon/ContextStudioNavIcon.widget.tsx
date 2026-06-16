import { ContextStudioNavIcon } from "./ContextStudioNavIcon";
import { defineWidget } from "../../../widgets/registry";
import type { ContextStudioNavIconWidgetProps } from "../../../widgets/ir";

export const contextStudioNavIconWidget = defineWidget<ContextStudioNavIconWidgetProps>({
	type: "ContextStudioNavIcon",
	module: "course.dsl",
	render: (props) => (
		<ContextStudioNavIcon className={props.className} id={props.id} title={props.title} />
	),
});
