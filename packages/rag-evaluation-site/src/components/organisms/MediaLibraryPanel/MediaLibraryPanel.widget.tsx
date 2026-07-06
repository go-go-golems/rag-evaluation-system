import { useEffect, useState } from "react";
import type { ActionSpec, MediaLibraryPanelWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";
import { defineWidget } from "../../../widgets/registry";
import { serializeUploadFiles } from "../../../widgets/uploadSerialization";
import type { MediaLibraryKindFilter } from "./MediaLibraryPanel";
import { MediaLibraryPanel } from "./MediaLibraryPanel";

export const mediaLibraryPanelWidget = defineWidget<MediaLibraryPanelWidgetProps>({
	type: "MediaLibraryPanel",
	module: "cms.dsl",
	render: (props, _children, ctx) => <MediaLibraryPanelWidgetHost props={props} ctx={ctx} />,
});

function MediaLibraryPanelWidgetHost({
	props,
	ctx,
}: {
	props: MediaLibraryPanelWidgetProps;
	ctx: RenderContext;
}) {
	const [query, setQuery] = useState(props.query ?? "");
	useEffect(() => setQuery(props.query ?? ""), [props.query]);
	const onAssetSelectAction = props.onAssetSelectAction;
	const onAssetOpenAction = props.onAssetOpenAction;
	const onQuerySubmitAction = props.onQuerySubmitAction;
	const onKindFilterChangeAction = props.onKindFilterChangeAction;
	const onPageChangeAction = props.onPageChangeAction;
	const onFilesSelectedAction = props.onFilesSelectedAction;
	const dispatchWithValue = (action: ActionSpec, extra: Record<string, unknown>) =>
		ctx.dispatchAction(action, { componentType: "MediaLibraryPanel", ...extra });
	return (
		<MediaLibraryPanel
			className={props.className}
			assets={props.assets}
			selectedAssetIds={props.selectedAssetIds}
			selectionMode={props.selectionMode}
			onAssetSelect={
				onAssetSelectAction
					? (assetId) => dispatchWithValue(onAssetSelectAction, { assetId, value: assetId })
					: undefined
			}
			onAssetOpen={
				onAssetOpenAction
					? (assetId) => dispatchWithValue(onAssetOpenAction, { assetId, value: assetId })
					: undefined
			}
			query={query}
			onQueryChange={onQuerySubmitAction ? setQuery : undefined}
			onQuerySubmit={
				onQuerySubmitAction
					? (submittedQuery) =>
							dispatchWithValue(onQuerySubmitAction, {
								query: submittedQuery,
								value: submittedQuery,
							})
					: undefined
			}
			kindFilter={props.kindFilter}
			onKindFilterChange={
				onKindFilterChangeAction
					? (kind: MediaLibraryKindFilter) =>
							dispatchWithValue(onKindFilterChangeAction, { kind, value: kind })
					: undefined
			}
			page={props.page}
			pageCount={props.pageCount}
			onPageChange={
				onPageChangeAction
					? (page) => dispatchWithValue(onPageChangeAction, { page, value: page })
					: undefined
			}
			onFilesSelected={
				onFilesSelectedAction
					? (files) => {
							void runFileSelectionAction(onFilesSelectedAction, files, ctx);
						}
					: undefined
			}
			uploads={props.uploads}
			showStatusBadges={props.showStatusBadges}
			emptyMessage={props.emptyMessage != null ? ctx.renderValue(props.emptyMessage) : undefined}
			title={props.title != null ? ctx.renderValue(props.title) : undefined}
			minTileWidth={props.minTileWidth}
		/>
	);
}

async function runFileSelectionAction(
	action: ActionSpec,
	files: File[],
	ctx: RenderContext,
): Promise<void> {
	const serializedFiles = await serializeUploadFiles(files);
	ctx.dispatchAction(action, {
		componentType: "MediaLibraryPanel",
		files: serializedFiles,
		fileNames: serializedFiles.map((file) => file.name),
		fileCount: serializedFiles.length,
	});
}
