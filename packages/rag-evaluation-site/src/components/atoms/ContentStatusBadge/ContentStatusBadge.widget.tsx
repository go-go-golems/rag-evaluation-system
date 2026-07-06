import type { ContentStatusBadgeWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ContentStatusBadge } from "./ContentStatusBadge";

export const contentStatusBadgeWidget = defineWidget<ContentStatusBadgeWidgetProps>({
	type: "ContentStatusBadge",
	module: "cms.dsl",
	render: (props) => (
		<ContentStatusBadge className={props.className} status={props.status} icon={props.icon} />
	),
});
