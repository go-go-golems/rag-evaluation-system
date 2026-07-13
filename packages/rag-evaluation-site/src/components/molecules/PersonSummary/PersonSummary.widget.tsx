import type { PersonSummaryWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { PersonSummary } from "./PersonSummary";

export const personSummaryWidget = defineWidget<PersonSummaryWidgetProps>({
	type: "PersonSummary",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<PersonSummary
			className={props.className}
			name={ctx.renderValue(props.name)}
			subtitle={ctx.renderValue(props.subtitle)}
			bio={ctx.renderValue(props.bio)}
			avatar={ctx.renderValue(props.avatar)}
		/>
	),
});
