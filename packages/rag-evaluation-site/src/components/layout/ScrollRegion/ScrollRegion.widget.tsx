import type { CSSProperties } from "react";
import type { ScrollRegionWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ScrollRegion } from "./ScrollRegion";

export const scrollRegionWidget = defineWidget<ScrollRegionWidgetProps>({
	type: "ScrollRegion",
	module: "ui.dsl",
	render: (props, children) => (
		<ScrollRegion
			className={props.className}
			axis={props.axis}
			style={props.style as CSSProperties | undefined}
		>
			{children}
		</ScrollRegion>
	),
});
