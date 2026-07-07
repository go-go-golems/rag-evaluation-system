import type { ReactNode } from "react";
import type { ActivityFeedWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ActivityFeed } from "./ActivityFeed";

export const activityFeedWidget = defineWidget<ActivityFeedWidgetProps>({
	type: "ActivityFeed",
	module: "data.dsl",
	render: (props, _children, ctx) => {
		const glyphs: Record<string, ReactNode> | undefined = props.glyphs
			? Object.fromEntries(
					Object.entries(props.glyphs).map(([key, value]) => [key, ctx.renderValue(value)]),
				)
			: undefined;
		return (
			<ActivityFeed
				activities={props.activities.map((a) => ({
					id: a.id,
					kind: a.kind,
					title: ctx.renderValue(a.title),
					body: a.body != null ? ctx.renderValue(a.body) : undefined,
					atISO: a.atISO,
					actor: a.actor,
				}))}
				glyphs={glyphs}
				styleSet={props.styleSet}
				groupByDay={props.groupByDay}
				onOpen={
					props.onOpenAction
						? (id) =>
								ctx.dispatchAction(props.onOpenAction!, {
									activityId: id,
									componentType: "ActivityFeed",
								} as unknown as Record<string, unknown>)
						: undefined
				}
				onLoadMore={
					props.onLoadMoreAction
						? () => ctx.dispatchAction(props.onLoadMoreAction!, { componentType: "ActivityFeed" })
						: undefined
				}
			/>
		);
	},
});
