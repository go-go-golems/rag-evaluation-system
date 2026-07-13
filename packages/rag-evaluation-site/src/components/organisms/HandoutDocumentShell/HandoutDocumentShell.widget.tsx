import type { HandoutDocumentShellWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { HandoutDocumentShell } from "./HandoutDocumentShell";

export const handoutDocumentShellWidget = defineWidget<HandoutDocumentShellWidgetProps>({
	type: "HandoutDocumentShell",
	module: "widget.dsl",
	render: (props, _children, ctx) => (
		<HandoutDocumentShell
			className={props.className}
			intro={ctx.renderValue(props.intro)}
			documents={props.documents}
			selectedDocumentId={props.selectedDocumentId}
			onDocumentSelect={
				props.onDocumentSelectAction
					? (documentId) =>
							ctx.dispatchAction(props.onDocumentSelectAction!, {
								documentId,
								document: { id: documentId },
								value: documentId,
								componentType: "HandoutDocumentShell",
							})
					: undefined
			}
			onDownload={
				props.onDownloadAction
					? (documentId) =>
							ctx.dispatchAction(props.onDownloadAction!, {
								documentId,
								document: { id: documentId },
								value: documentId,
								componentType: "HandoutDocumentShell",
							})
					: undefined
			}
			onDownloadAll={
				props.onDownloadAllAction
					? () =>
							ctx.dispatchAction(props.onDownloadAllAction!, {
								componentType: "HandoutDocumentShell",
							})
					: undefined
			}
			onPrint={
				props.onPrintAction
					? (documentId) =>
							ctx.dispatchAction(props.onPrintAction!, {
								documentId,
								document: { id: documentId },
								value: documentId,
								componentType: "HandoutDocumentShell",
							})
					: undefined
			}
			styleSet={props.styleSet}
			title={ctx.renderValue(props.title)}
			emptyMessage={ctx.renderValue(props.emptyMessage)}
		/>
	),
});
