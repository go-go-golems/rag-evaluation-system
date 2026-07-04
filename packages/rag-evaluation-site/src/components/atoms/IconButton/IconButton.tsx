import type { ButtonHTMLAttributes, ReactNode } from "react";
import styles from "./IconButton.module.css";

export type IconButtonSize = "compact" | "normal" | "large";
export type IconButtonVariant = "bare" | "boxed";

export interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
	size?: IconButtonSize;
	variant?: IconButtonVariant;
	label?: string;
	children?: ReactNode;
}

export function IconButton({
	size = "compact",
	variant = "bare",
	label,
	className,
	children,
	...rest
}: IconButtonProps) {
	return (
		<button
			type="button"
			aria-label={label}
			title={rest.title ?? label}
			className={[
				styles.root,
				styles[size],
				variant === "boxed" ? styles.boxed : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="IconButton"
			data-variant={variant}
			{...rest}
		>
			{children}
		</button>
	);
}
