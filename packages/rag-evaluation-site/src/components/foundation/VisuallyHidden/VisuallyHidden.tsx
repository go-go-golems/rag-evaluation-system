import type { HTMLAttributes, ReactNode } from "react";
import styles from "./VisuallyHidden.module.css";

export interface VisuallyHiddenProps extends HTMLAttributes<HTMLSpanElement> {
	children?: ReactNode;
}

export function VisuallyHidden({ className, children, ...rest }: VisuallyHiddenProps) {
	return (
		<span className={[styles.root, className ?? ""].filter(Boolean).join(" ")} {...rest}>
			{children}
		</span>
	);
}
