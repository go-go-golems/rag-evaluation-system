import type { HTMLAttributes, ReactNode } from "react";
import { Caption } from "../../foundation";
import styles from "./SlideShell.module.css";

export interface SlideShellProps extends Omit<HTMLAttributes<HTMLElement>, "title"> {
	as?: "article" | "section" | "div";
	eyebrow?: ReactNode;
	counter?: ReactNode;
	title: ReactNode;
	subtitle?: ReactNode;
	primary: ReactNode;
	secondary?: ReactNode;
	primarySide?: "left" | "right";
	ratio?: "balanced" | "primaryWide" | "secondaryWide";
	divider?: boolean;
	footer?: ReactNode;
	mode?: "default" | "present";
}

export function SlideShell({
	as: Element = "article",
	eyebrow,
	counter,
	title,
	subtitle,
	primary,
	secondary,
	primarySide = "left",
	ratio = "primaryWide",
	divider = true,
	footer,
	mode = "default",
	className,
	...rest
}: SlideShellProps) {
	return (
		<Element
			className={[styles.root, mode === "present" ? styles.present : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-layout="SlideShell"
			data-primary-side={primarySide}
			data-slide-mode={mode}
			{...rest}
		>
			<header className={styles.header}>
				<div className={styles.metaRow}>
					{eyebrow ? <Caption transform="uppercase">{eyebrow}</Caption> : <span />}
					{counter && <Caption>{counter}</Caption>}
				</div>
				<h2 className={styles.title}>{title}</h2>
				{subtitle && <p className={styles.subtitle}>{subtitle}</p>}
			</header>
			<div
				className={[
					styles.body,
					styles[ratio],
					secondary ? styles.twoPane : styles.singlePane,
					divider && secondary ? styles.divider : "",
					styles[`primary-${primarySide}`],
				]
					.filter(Boolean)
					.join(" ")}
			>
				<div className={styles.primary}>{primary}</div>
				{secondary && <div className={styles.secondary}>{secondary}</div>}
			</div>
			{footer && <footer className={styles.footer}>{footer}</footer>}
		</Element>
	);
}
