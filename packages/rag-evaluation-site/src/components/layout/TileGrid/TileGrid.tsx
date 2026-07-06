import type { CSSProperties, HTMLAttributes, ReactNode } from "react";
import styles from "./TileGrid.module.css";

export type TileGridGap = "sm" | "md";

export interface TileGridProps extends HTMLAttributes<HTMLDivElement> {
	minTileWidth?: number;
	gap?: TileGridGap;
	children?: ReactNode;
}

export function TileGrid({
	minTileWidth = 160,
	gap = "md",
	className,
	style,
	children,
	...rest
}: TileGridProps) {
	return (
		<div
			className={[styles.root, gap === "sm" ? styles.gapSm : styles.gapMd, className ?? ""]
				.filter(Boolean)
				.join(" ")}
			style={{ "--rag-tile-min": `${minTileWidth}px`, ...style } as CSSProperties}
			data-rag-layout="TileGrid"
			{...rest}
		>
			{children}
		</div>
	);
}
