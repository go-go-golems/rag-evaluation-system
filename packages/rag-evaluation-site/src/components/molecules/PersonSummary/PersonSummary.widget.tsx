import { PersonSummary } from "./PersonSummary";
import { defineWidget } from "../../../widgets/registry";
import type { PersonSummaryWidgetProps } from "../../../widgets/ir";

export const personSummaryWidget = defineWidget<PersonSummaryWidgetProps>({
	type: "PersonSummary",
	module: "ui.dsl",
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
