import { CourseStepNav } from "./CourseStepNav";
import { defineWidget } from "../../../widgets/registry";
import type { CourseStepNavWidgetProps } from "../../../widgets/ir";

export const courseStepNavWidget = defineWidget<CourseStepNavWidgetProps>({
	type: "CourseStepNav",
	module: "course.dsl",
	render: (props, _children, ctx) => (
		<CourseStepNav
			className={props.className}
			items={props.items}
			activeItemId={props.activeItemId}
			onItemSelect={
				props.onItemSelectAction
					? (itemId) =>
							ctx.dispatchAction(props.onItemSelectAction!, {
								itemId,
								value: itemId,
								componentType: "CourseStepNav",
							})
					: undefined
			}
		/>
	),
});
