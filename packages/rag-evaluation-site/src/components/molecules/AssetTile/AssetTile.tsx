import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cmsAssetMeta } from "../../../cms/fixtures";
import type { CmsAsset } from "../../../cms/types";
import { MediaThumb } from "../../atoms";
import styles from "./AssetTile.module.css";

export interface AssetTileProps
	extends Omit<ButtonHTMLAttributes<HTMLButtonElement>, "onSelect" | "onClick"> {
	asset: CmsAsset;
	selected?: boolean;
	onSelect?: (assetId: string) => void;
	onOpen?: (assetId: string) => void;
	footerSlot?: ReactNode;
}

const FILE_GLYPH: Record<string, string> = {
	pdf: "▤",
	json: "{ }",
	markdown: "¶",
	zip: "▣",
};

function fileGlyph(mime: string): string {
	const subtype = mime.split("/")[1] ?? "";
	return FILE_GLYPH[subtype] ?? "□";
}

export function AssetTile({
	asset,
	selected = false,
	onSelect,
	onOpen,
	footerSlot,
	className,
	...rest
}: AssetTileProps) {
	return (
		<button
			type="button"
			className={[styles.root, selected ? styles.selected : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			onClick={() => onSelect?.(asset.id)}
			onDoubleClick={() => onOpen?.(asset.id)}
			aria-pressed={selected}
			data-rag-molecule="AssetTile"
			data-rag-asset-id={asset.id}
			data-active={selected || undefined}
			{...rest}
		>
			{asset.kind === "image" ? (
				<MediaThumb
					className={styles.thumb}
					src={asset.thumbSrc ?? asset.src}
					alt={asset.alt ?? asset.title}
					frame="none"
				/>
			) : (
				<div className={styles.fileThumb} data-state="file">
					<span className={styles.fileGlyph} aria-hidden="true">
						{fileGlyph(asset.mime)}
					</span>
				</div>
			)}
			<span className={styles.titleRow}>
				<span className={styles.title}>{asset.title}</span>
			</span>
			<span className={styles.metaRow}>
				<span className={styles.meta}>{cmsAssetMeta(asset)}</span>
				{footerSlot}
			</span>
		</button>
	);
}
