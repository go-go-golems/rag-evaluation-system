import type { DividerWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Divider } from "./Divider";

export const dividerWidget = defineWidget<DividerWidgetProps>({
	type: "Divider",
	module: "ui.dsl",
	render: (props) => <Divider className={props.className} orientation={props.orientation} />,
});
