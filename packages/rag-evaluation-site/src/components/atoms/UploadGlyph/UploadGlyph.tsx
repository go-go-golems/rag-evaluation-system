import type { SVGAttributes } from "react";
import styles from "./UploadGlyph.module.css";

export interface UploadGlyphProps extends SVGAttributes<SVGSVGElement> {
	title?: string;
}

export function UploadGlyph({ title = "Upload file", className, ...rest }: UploadGlyphProps) {
	const titleId = title ? "rag-upload-glyph-title" : undefined;
	return (
		<svg
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			viewBox="0 0 40 34"
			role={title ? "img" : "presentation"}
			aria-labelledby={titleId}
			data-rag-atom="UploadGlyph"
			{...rest}
		>
			{title && <title id={titleId}>{title}</title>}
			<path className={styles.folder} d="M2 8 L2 32 L38 32 L38 4 L18 4 L14 8 Z" />
			<line className={styles.arrow} x1="20" y1="14" x2="20" y2="26" />
			<path className={styles.arrow} d="M15 19 L20 14 L25 19" />
		</svg>
	);
}
