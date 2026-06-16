import type { SVGAttributes } from "react";
import styles from "./ContextStudioNavIcon.module.css";

export type ContextStudioNavIconId =
	| "course"
	| "slides"
	| "visualize"
	| "upload"
	| "transcript"
	| "comments"
	| "handout"
	| "handouts";

export interface ContextStudioNavIconProps extends Omit<SVGAttributes<SVGSVGElement>, "id"> {
	id: ContextStudioNavIconId;
	title?: string;
}

export function ContextStudioNavIcon({ id, title, className, ...rest }: ContextStudioNavIconProps) {
	const titleId = title ? `rag-context-studio-nav-icon-${id}` : undefined;
	const common = {
		className: [styles.root, className ?? ""].filter(Boolean).join(" "),
		viewBox: "0 0 15 15",
		role: title ? "img" : "presentation",
		"aria-labelledby": titleId,
		"data-rag-atom": "ContextStudioNavIcon",
		"data-icon-id": id,
		...rest,
	} as const;

	return (
		<svg {...common}>
			{title && <title id={titleId}>{title}</title>}
			{renderIcon(id)}
		</svg>
	);
}

function renderIcon(id: ContextStudioNavIconId) {
	switch (id) {
		case "course":
			return (
				<>
					<rect x="2" y="1.5" width="11" height="12" />
					<line x1="4.5" y1="5" x2="10.5" y2="5" />
					<line x1="4.5" y1="8" x2="10.5" y2="8" />
					<line x1="4.5" y1="11" x2="8" y2="11" />
				</>
			);
		case "slides":
			return (
				<>
					<rect x="1.5" y="2.5" width="12" height="8.5" />
					<line x1="6" y1="13" x2="9" y2="13" />
					<line x1="7.5" y1="11" x2="7.5" y2="13" />
				</>
			);
		case "visualize":
			return (
				<>
					<rect className={styles.filled} x="1.5" y="4" width="3" height="7" />
					<rect x="5.5" y="2" width="3" height="9" />
					<rect className={styles.filled} x="9.5" y="6" width="3" height="5" />
				</>
			);
		case "upload":
			return (
				<>
					<path d="M2 9 L2 13 L13 13 L13 9" />
					<line x1="7.5" y1="2" x2="7.5" y2="9.5" />
					<path d="M4.5 5 L7.5 2 L10.5 5" />
				</>
			);
		case "transcript":
			return (
				<>
					<rect x="1.5" y="2" width="12" height="11" />
					<line x1="4" y1="5" x2="11" y2="5" />
					<line x1="4" y1="7.5" x2="11" y2="7.5" />
					<line x1="4" y1="10" x2="8" y2="10" />
				</>
			);
		case "comments":
			return <path d="M2 2 L13 2 L13 9 L7 9 L4 12 L4 9 L2 9 Z" />;
		case "handout":
		case "handouts":
			return (
				<>
					<path d="M3 1.5 L10 1.5 L12.5 4 L12.5 13.5 L3 13.5 Z" />
					<path d="M10 1.5 L10 4 L12.5 4" />
					<line x1="5" y1="7" x2="10.5" y2="7" />
					<line x1="5" y1="9.5" x2="10.5" y2="9.5" />
				</>
			);
	}
}
