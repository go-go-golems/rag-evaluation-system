import type { SectionBlockWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SectionBlock } from "./SectionBlock";

export const sectionBlockWidget = defineWidget<SectionBlockWidgetProps>({
	type: "SectionBlock",
	module: "widget.dsl",
	render: (props, children, ctx) => (
		<SectionBlock
			className={props.className}
			as={props.as}
			id={props.anchorId}
			label={ctx.renderValue(props.label)}
			caption={ctx.renderValue(props.caption)}
			actions={ctx.renderValue(props.actions)}
			level={props.level}
			rule={props.rule}
			density={props.density}
			divider={props.divider}
		>
			{children}
		</SectionBlock>
	),
});
