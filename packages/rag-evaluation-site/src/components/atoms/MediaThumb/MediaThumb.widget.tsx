import type { MediaThumbWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { MediaThumb } from "./MediaThumb";

export const mediaThumbWidget = defineWidget<MediaThumbWidgetProps>({
	type: "MediaThumb",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<MediaThumb
			className={props.className}
			src={props.src}
			alt={props.alt}
			aspect={props.aspect}
			fit={props.fit}
			frame={props.frame}
			selected={props.selected}
			fallbackGlyph={props.fallbackGlyph != null ? ctx.renderValue(props.fallbackGlyph) : undefined}
			fallbackLabel={props.fallbackLabel != null ? ctx.renderValue(props.fallbackLabel) : undefined}
		/>
	),
});
