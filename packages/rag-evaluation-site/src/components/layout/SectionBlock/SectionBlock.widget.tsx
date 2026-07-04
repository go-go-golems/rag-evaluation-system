import type { SectionBlockWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SectionBlock } from "./SectionBlock";

export const sectionBlockWidget = defineWidget<SectionBlockWidgetProps>({
	type: "SectionBlock",
	module: "ui.dsl",
	render: (props, children, ctx) => (
		<SectionBlock
			className={props.className}
			as={props.as}
			label={ctx.renderValue(props.label)}
			caption={ctx.renderValue(props.caption)}
			density={props.density}
			divider={props.divider}
		>
			{children}
		</SectionBlock>
	),
});
