import type { DetailsHTMLAttributes, ReactNode } from "react";
import styles from "./Disclosure.module.css";

export interface DisclosureProps extends Omit<DetailsHTMLAttributes<HTMLDetailsElement>, "title"> {
	title: ReactNode;
	children?: ReactNode;
}

export function Disclosure({ title, children, className, ...rest }: DisclosureProps) {
	return (
		<details
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="Disclosure"
			{...rest}
		>
			<summary className={styles.summary}>{title}</summary>
			<div className={styles.content}>{children}</div>
		</details>
	);
}
