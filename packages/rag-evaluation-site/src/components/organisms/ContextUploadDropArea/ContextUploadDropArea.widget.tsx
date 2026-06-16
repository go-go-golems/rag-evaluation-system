import { ContextUploadDropArea } from "./ContextUploadDropArea";
import { defineWidget } from "../../../widgets/registry";
import type { ActionSpec, ContextUploadDropAreaWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";

export const contextUploadDropAreaWidget = defineWidget<ContextUploadDropAreaWidgetProps>({
	type: "ContextUploadDropArea",
	module: "context_window.dsl",
	render: (props, _children, ctx) => (
		<ContextUploadDropArea
			className={props.className}
			title={ctx.renderValue(props.title)}
			description={ctx.renderValue(props.description)}
			accept={props.accept}
			disabled={props.disabled}
			active={props.active}
			onFilesSelected={
				props.onFilesSelectedAction
					? (files) => {
							void runFileSelectionAction(props.onFilesSelectedAction!, files, ctx);
						}
					: undefined
			}
		/>
	),
});

async function runFileSelectionAction(
	action: ActionSpec,
	files: File[],
	ctx: RenderContext,
): Promise<void> {
	const serializedFiles = await Promise.all(
		files.map(async (file) => ({
			name: file.name,
			size: file.size,
			type: file.type,
			lastModified: file.lastModified,
			text: await file.text(),
		})),
	);
	ctx.dispatchAction(action, {
		componentType: "ContextUploadDropArea",
		files: serializedFiles,
		fileNames: serializedFiles.map((file) => file.name),
		fileCount: serializedFiles.length,
	});
}
