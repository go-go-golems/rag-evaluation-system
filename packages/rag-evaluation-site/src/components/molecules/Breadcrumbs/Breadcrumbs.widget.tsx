import type { BreadcrumbsWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Breadcrumbs } from "./Breadcrumbs";

export const breadcrumbsWidget = defineWidget<BreadcrumbsWidgetProps>({
	type: "Breadcrumbs",
	module: "widget.dsl",
	render: (props, _children, ctx) => {
		const onNavigateAction = props.onNavigateAction;
		return (
			<Breadcrumbs
				className={props.className}
				ariaLabel={props.ariaLabel}
				items={props.items.map((item) => ({ id: item.id, label: ctx.renderValue(item.label) }))}
				onNavigate={
					onNavigateAction
						? (itemId) =>
								ctx.dispatchAction(onNavigateAction, {
									itemId,
									value: itemId,
									componentType: "Breadcrumbs",
								})
						: undefined
				}
			/>
		);
	},
});
