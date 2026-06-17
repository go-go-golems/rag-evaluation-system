import { ContextUploadDropArea } from "./ContextUploadDropArea";
import { defineWidget } from "../../../widgets/registry";
import type { ActionSpec, ContextUploadDropAreaWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";

export interface SerializedUploadFile {
	name: string;
	size: number;
	type: string;
	lastModified: number;
	encoding: "utf8" | "base64";
	text?: string;
	base64?: string;
}

const textMimePrefixes = ["text/"];
const textMimeTypes = new Set([
	"application/json",
	"application/ld+json",
	"application/xml",
	"application/xhtml+xml",
	"application/yaml",
	"application/x-yaml",
	"image/svg+xml",
]);
const textExtensions = new Set([
	".csv",
	".css",
	".html",
	".js",
	".json",
	".jsonl",
	".jsx",
	".log",
	".md",
	".markdown",
	".mjs",
	".svg",
	".text",
	".toml",
	".ts",
	".tsx",
	".txt",
	".xml",
	".yaml",
	".yml",
]);

export const contextUploadDropAreaWidget = defineWidget<ContextUploadDropAreaWidgetProps>({
	type: "ContextUploadDropArea",
	module: "context_window.dsl",
	render: (props, _children, ctx) => {
		const onFilesSelectedAction = props.onFilesSelectedAction;
		return (
			<ContextUploadDropArea
				className={props.className}
				title={ctx.renderValue(props.title)}
				description={ctx.renderValue(props.description)}
				accept={props.accept}
				disabled={props.disabled}
				active={props.active}
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
	const serializedFiles = await Promise.all(files.map(serializeUploadFile));
	ctx.dispatchAction(action, {
		componentType: "ContextUploadDropArea",
		files: serializedFiles,
		fileNames: serializedFiles.map((file) => file.name),
		fileCount: serializedFiles.length,
	});
}

async function serializeUploadFile(file: File): Promise<SerializedUploadFile> {
	const base = {
		name: file.name,
		size: file.size,
		type: file.type,
		lastModified: file.lastModified,
	};

	if (isLikelyTextFile(file)) {
		return {
			...base,
			encoding: "utf8",
			text: await file.text(),
		};
	}

	return {
		...base,
		encoding: "base64",
		base64: arrayBufferToBase64(await file.arrayBuffer()),
	};
}

function isLikelyTextFile(file: File): boolean {
	const mimeType = file.type.toLowerCase();
	if (textMimePrefixes.some((prefix) => mimeType.startsWith(prefix))) return true;
	if (textMimeTypes.has(mimeType)) return true;
	const lowerName = file.name.toLowerCase();
	return Array.from(textExtensions).some((extension) => lowerName.endsWith(extension));
}

function arrayBufferToBase64(buffer: ArrayBuffer): string {
	const bytes = new Uint8Array(buffer);
	const chunkSize = 0x8000;
	let binary = "";
	for (let index = 0; index < bytes.length; index += chunkSize) {
		const chunk = bytes.subarray(index, index + chunkSize);
		binary += String.fromCharCode(...chunk);
	}
	return btoa(binary);
}
