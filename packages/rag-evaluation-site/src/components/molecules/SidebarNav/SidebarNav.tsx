import type { HTMLAttributes, ReactNode } from "react";
import styles from "./SidebarNav.module.css";

export interface SidebarNavItem {
	id: string;
	label: ReactNode;
	icon?: ReactNode;
	badge?: ReactNode;
	disabled?: boolean;
}

export interface SidebarNavSection {
	id: string;
	label: ReactNode;
	items: SidebarNavItem[];
}

export interface SidebarNavProps extends Omit<HTMLAttributes<HTMLElement>, "onSelect"> {
	sections: SidebarNavSection[];
	activeItemId?: string;
	onItemSelect?: (itemId: string) => void;
	ariaLabel?: string;
}

export function SidebarNav({
	sections,
	activeItemId,
	onItemSelect,
	ariaLabel = "Course",
	className,
	...rest
}: SidebarNavProps) {
	return (
		<nav
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			aria-label={ariaLabel}
			data-rag-molecule="SidebarNav"
			{...rest}
		>
			{sections.map((section) => (
				<section key={section.id} className={styles.section}>
					<div className={styles.sectionLabel}>{section.label}</div>
					<div className={styles.items}>
						{section.items.map((item) => {
							const active = item.id === activeItemId;
							return (
								<button
									key={item.id}
									type="button"
									className={styles.item}
									data-active={active ? "true" : undefined}
									disabled={item.disabled}
									aria-current={active ? "page" : undefined}
									onClick={() => onItemSelect?.(item.id)}
								>
									<span className={styles.icon} aria-hidden="true">
										{item.icon}
									</span>
									<span className={styles.label}>{item.label}</span>
									{item.badge && <span className={styles.badge}>{item.badge}</span>}
								</button>
							);
						})}
					</div>
				</section>
			))}
		</nav>
	);
}
