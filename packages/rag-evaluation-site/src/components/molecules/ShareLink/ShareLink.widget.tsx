import type { ActionSpec, ShareLinkWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ShareLink } from "./ShareLink";

export const shareLinkWidget = defineWidget<ShareLinkWidgetProps>({
	type: "ShareLink",
	module: "widget.dsl",
	render: (props, _children, ctx) => {
		const copyAction: ActionSpec = props.copyAction ?? { kind: "copy", value: props.href };
		return (
			<ShareLink
				className={props.className}
				label={props.label != null ? ctx.renderValue(props.label) : undefined}
				description={props.description != null ? ctx.renderValue(props.description) : undefined}
				href={props.href}
				displayHref={props.displayHref != null ? ctx.renderValue(props.displayHref) : undefined}
				copyLabel={props.copyLabel != null ? ctx.renderValue(props.copyLabel) : undefined}
				copiedLabel={props.copiedLabel != null ? ctx.renderValue(props.copiedLabel) : undefined}
				copied={props.copied}
				onCopy={ctx.bindAction(copyAction, {
					componentType: "ShareLink",
					value: props.href,
				})}
			/>
		);
	},
});
