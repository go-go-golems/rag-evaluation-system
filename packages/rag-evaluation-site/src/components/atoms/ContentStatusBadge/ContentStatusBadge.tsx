import type { HTMLAttributes } from "react";
import type { CmsContentStatus } from "../../../cms/types";
import styles from "./ContentStatusBadge.module.css";

export interface ContentStatusBadgeProps extends HTMLAttributes<HTMLSpanElement> {
	status: CmsContentStatus;
	icon?: boolean;
}

const STATUS_GLYPH: Record<CmsContentStatus, string> = {
	draft: "◌",
	published: "●",
	scheduled: "◔",
	archived: "▣",
};

export function ContentStatusBadge({
	status,
	icon = true,
	className,
	...rest
}: ContentStatusBadgeProps) {
	return (
		<span
			className={[styles.root, styles[status], className ?? ""].filter(Boolean).join(" ")}
			data-rag-atom="ContentStatusBadge"
			data-status={status}
			{...rest}
		>
			{icon && (
				<span className={styles.glyph} aria-hidden="true">
					{STATUS_GLYPH[status]}
				</span>
			)}
			{status}
		</span>
	);
}
