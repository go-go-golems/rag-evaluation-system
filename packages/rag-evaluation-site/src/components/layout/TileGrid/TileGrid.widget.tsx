import type { TileGridWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TileGrid } from "./TileGrid";

export const tileGridWidget = defineWidget<TileGridWidgetProps>({
	type: "TileGrid",
	module: "cms.dsl",
	render: (props, children) => (
		<TileGrid className={props.className} minTileWidth={props.minTileWidth} gap={props.gap}>
			{children}
		</TileGrid>
	),
});
