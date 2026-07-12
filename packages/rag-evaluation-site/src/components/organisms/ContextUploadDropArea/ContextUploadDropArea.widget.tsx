import type { ActionSpec, ContextUploadDropAreaWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";
import { defineWidget } from "../../../widgets/registry";
import { serializeUploadFiles } from "../../../widgets/uploadSerialization";
import { ContextUploadDropArea } from "./ContextUploadDropArea";

export type { SerializedUploadFile } from "../../../widgets/uploadSerialization";

export const contextUploadDropAreaWidget = defineWidget<ContextUploadDropAreaWidgetProps>({
	type: "ContextUploadDropArea",
	module: "context_window.dsl",
	render: (props, _children, ctx) => {
		const onFilesSelectedAction = props.onFilesSelectedAction;
		const onDeleteAction = props.onDeleteAction;
		return (
			<ContextUploadDropArea
				className={props.className}
				title={ctx.renderValue(props.title)}
				description={ctx.renderValue(props.description)}
				accept={props.accept}
				disabled={props.disabled}
				active={props.active}
				items={props.items?.map((item) => ({ ...item, label: ctx.renderValue(item.label) }))}
				onDelete={
					onDeleteAction
						? (itemId) =>
								ctx.dispatchAction(onDeleteAction, {
									assetId: itemId,
									asset: { id: itemId },
									value: itemId,
									componentType: "ContextUploadDropArea",
								})
						: undefined
				}
				onFilesSelected={
					onFilesSelectedAction
						? (files) => {
								void runFileSelectionAction(onFilesSelectedAction, files, ctx);
							}
						: undefined
				}
			/>
		);
	},
});

async function runFileSelectionAction(
	action: ActionSpec,
	files: File[],
	ctx: RenderContext,
): Promise<void> {
	const serializedFiles = await serializeUploadFiles(files);
	ctx.dispatchAction(action, {
		componentType: "ContextUploadDropArea",
		files: serializedFiles,
		fileNames: serializedFiles.map((file) => file.name),
		fileCount: serializedFiles.length,
	});
}
