import type { DisclosureWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Disclosure } from "./Disclosure";

export const disclosureWidget = defineWidget<DisclosureWidgetProps>({
	type: "Disclosure",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<Disclosure className={props.className} title={ctx.renderValue(props.title)} open={props.open}>
			{children}
		</Disclosure>
	),
});
