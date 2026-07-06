import type { DocumentPreviewToolbarWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { DocumentPreviewToolbar } from "./DocumentPreviewToolbar";

export const documentPreviewToolbarWidget = defineWidget<DocumentPreviewToolbarWidgetProps>({
	type: "DocumentPreviewToolbar",
	module: "course.dsl",
	render: (props, _children, ctx) => (
		<DocumentPreviewToolbar
			className={props.className}
			file={ctx.renderValue(props.file)}
			format={ctx.renderValue(props.format)}
			size={ctx.renderValue(props.size)}
			onDownload={
				props.onDownloadAction
					? () =>
							ctx.dispatchAction(props.onDownloadAction!, {
								componentType: "DocumentPreviewToolbar",
							})
					: undefined
			}
			downloadLabel={ctx.renderValue(props.downloadLabel)}
			onPrint={
				props.onPrintAction
					? () =>
							ctx.dispatchAction(props.onPrintAction!, { componentType: "DocumentPreviewToolbar" })
					: undefined
			}
			printLabel={ctx.renderValue(props.printLabel)}
			rightSlot={props.rightSlot ? ctx.renderNode(props.rightSlot) : undefined}
		/>
	),
});
