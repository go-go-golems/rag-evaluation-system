import type { TranscriptRoleBadgeWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TranscriptRoleBadge } from "./TranscriptRoleBadge";

export const transcriptRoleBadgeWidget = defineWidget<TranscriptRoleBadgeWidgetProps>({
	type: "TranscriptRoleBadge",
	module: "context_window.dsl",
	render: (props) => (
		<TranscriptRoleBadge className={props.className} role={props.role} name={props.name} />
	),
});
