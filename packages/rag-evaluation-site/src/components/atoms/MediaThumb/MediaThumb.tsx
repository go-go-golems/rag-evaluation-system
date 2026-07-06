import type { HTMLAttributes, ReactNode } from "react";
import { useEffect, useState } from "react";
import { Caption } from "../../foundation";
import styles from "./MediaThumb.module.css";

export type MediaThumbAspect = "square" | "wide" | "natural";
export type MediaThumbFit = "cover" | "contain";
export type MediaThumbFrame = "bordered" | "none";
export type MediaThumbState = "empty" | "loading" | "loaded" | "broken";

export interface MediaThumbProps extends Omit<HTMLAttributes<HTMLDivElement>, "onError"> {
	src?: string;
	alt?: string;
	aspect?: MediaThumbAspect;
	fit?: MediaThumbFit;
	frame?: MediaThumbFrame;
	fallbackGlyph?: ReactNode;
	fallbackLabel?: ReactNode;
	selected?: boolean;
}

export function MediaThumb({
	src,
	alt = "",
	aspect = "square",
	fit = "cover",
	frame = "bordered",
	fallbackGlyph = "▨",
	fallbackLabel,
	selected = false,
	className,
	...rest
}: MediaThumbProps) {
	const [state, setState] = useState<MediaThumbState>(src ? "loading" : "empty");

	useEffect(() => {
		setState(src ? "loading" : "empty");
	}, [src]);

	const showFallback = state === "empty" || state === "broken";

	return (
		<div
			className={[
				styles.root,
				styles[aspect],
				frame === "bordered" ? styles.bordered : "",
				selected ? styles.selected : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="MediaThumb"
			data-state={state}
			data-active={selected || undefined}
			{...rest}
		>
			{src && state !== "broken" && (
				<img
					className={[styles.image, fit === "contain" ? styles.contain : styles.cover].join(" ")}
					src={src}
					alt={alt}
					loading="lazy"
					onLoad={() => setState("loaded")}
					onError={() => setState("broken")}
				/>
			)}
			{showFallback && (
				<div className={styles.fallback}>
					<span className={styles.glyph} aria-hidden="true">
						{fallbackGlyph}
					</span>
					<Caption>{fallbackLabel ?? (state === "broken" ? "missing image" : "no image")}</Caption>
				</div>
			)}
		</div>
	);
}
