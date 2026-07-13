import type { DocumentListPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { DocumentListPanel } from "./DocumentListPanel";

export const documentListPanelWidget = defineWidget<DocumentListPanelWidgetProps>({
	type: "DocumentListPanel",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<DocumentListPanel
			className={props.className}
			title={ctx.renderValue(props.title)}
			description={ctx.renderValue(props.description)}
			items={props.items.map((item) => ({
				...item,
				title: ctx.renderValue(item.title),
				format: ctx.renderValue(item.format),
				size: ctx.renderValue(item.size),
				description: ctx.renderValue(item.description),
				icon: ctx.renderValue(item.icon),
			}))}
			selectedItemId={props.selectedItemId}
			onItemSelect={
				props.onItemSelectAction
					? (itemId) =>
							ctx.dispatchAction(props.onItemSelectAction!, {
								itemId,
								value: itemId,
								componentType: "DocumentListPanel",
							})
					: undefined
			}
			onDownloadAll={
				props.onDownloadAllAction
					? () =>
							ctx.dispatchAction(props.onDownloadAllAction!, { componentType: "DocumentListPanel" })
					: undefined
			}
			downloadAllLabel={ctx.renderValue(props.downloadAllLabel)}
			emptyMessage={ctx.renderValue(props.emptyMessage)}
			showItemDescriptions={props.showItemDescriptions}
		/>
	),
});
