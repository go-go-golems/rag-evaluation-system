import type { HTMLAttributes, ReactNode } from "react";
import { Caption } from "../../foundation";
import { SidebarShell } from "../../layout";
import { SidebarNav, type SidebarNavSection } from "../../molecules";
import styles from "./CmsShell.module.css";
import { cmsNavSections } from "./cmsNav";

export interface CmsShellProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	sections?: SidebarNavSection[];
	activeItemId?: string;
	onNavigate?: (itemId: string) => void;
	title?: ReactNode;
	subtitle?: ReactNode;
	headerSlot?: ReactNode;
	sidebarFooter?: ReactNode;
	contentPadding?: "default" | "none";
	children?: ReactNode;
}

export function CmsShell({
	sections = cmsNavSections,
	activeItemId,
	onNavigate,
	title = "Content Studio",
	subtitle = "CMS",
	headerSlot,
	sidebarFooter,
	contentPadding = "default",
	className,
	children,
	...rest
}: CmsShellProps) {
	const header = (
		<div className={styles.header}>
			<div className={styles.title}>{title}</div>
			{subtitle && <Caption>{subtitle}</Caption>}
			{headerSlot}
		</div>
	);

	return (
		<SidebarShell
			className={className}
			sidebarWidth={188}
			sidebarAriaLabel="CMS navigation"
			header={header}
			footer={sidebarFooter}
			sidebar={
				<SidebarNav
					sections={sections}
					activeItemId={activeItemId}
					onItemSelect={onNavigate}
					ariaLabel="CMS navigation"
				/>
			}
			data-rag-organism="CmsShell"
			{...rest}
		>
			<div
				className={[styles.content, contentPadding === "none" ? styles.contentPaddingNone : ""]
					.filter(Boolean)
					.join(" ")}
			>
				{children}
			</div>
		</SidebarShell>
	);
}
