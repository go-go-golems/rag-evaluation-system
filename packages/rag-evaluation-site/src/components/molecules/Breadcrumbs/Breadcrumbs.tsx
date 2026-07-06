import type { HTMLAttributes, ReactNode } from "react";
import { Fragment } from "react";
import styles from "./Breadcrumbs.module.css";

export interface BreadcrumbItem {
	id: string;
	label: ReactNode;
}

export interface BreadcrumbsProps extends HTMLAttributes<HTMLElement> {
	items: BreadcrumbItem[];
	onNavigate?: (itemId: string) => void;
	ariaLabel?: string;
}

export function Breadcrumbs({
	items,
	onNavigate,
	ariaLabel = "Breadcrumbs",
	className,
	...rest
}: BreadcrumbsProps) {
	return (
		<nav
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			aria-label={ariaLabel}
			data-rag-molecule="Breadcrumbs"
			{...rest}
		>
			{items.map((item, index) => {
				const isCurrent = index === items.length - 1;
				return (
					<Fragment key={item.id}>
						{index > 0 && (
							<span className={styles.separator} aria-hidden="true">
								/
							</span>
						)}
						{isCurrent || !onNavigate ? (
							<span
								className={isCurrent ? styles.current : styles.item}
								aria-current={isCurrent ? "page" : undefined}
							>
								{item.label}
							</span>
						) : (
							<button type="button" className={styles.link} onClick={() => onNavigate(item.id)}>
								{item.label}
							</button>
						)}
					</Fragment>
				);
			})}
		</nav>
	);
}
