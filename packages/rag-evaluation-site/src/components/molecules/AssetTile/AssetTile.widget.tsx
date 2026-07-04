import type { AssetTileWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { AssetTile } from "./AssetTile";

export const assetTileWidget = defineWidget<AssetTileWidgetProps>({
	type: "AssetTile",
	module: "cms.dsl",
	render: (props, _children, ctx) => {
		const onSelectAction = props.onSelectAction;
		const onOpenAction = props.onOpenAction;
		return (
			<AssetTile
				className={props.className}
				asset={props.asset}
				selected={props.selected}
				onSelect={
					onSelectAction
						? (assetId) =>
								ctx.dispatchAction(onSelectAction, {
									assetId,
									value: assetId,
									componentType: "AssetTile",
								})
						: undefined
				}
				onOpen={
					onOpenAction
						? (assetId) =>
								ctx.dispatchAction(onOpenAction, {
									assetId,
									value: assetId,
									componentType: "AssetTile",
								})
						: undefined
				}
				footerSlot={props.footer ? ctx.renderNode(props.footer) : undefined}
			/>
		);
	},
});
