import type { HTMLAttributes, ReactNode } from "react";
import { useEffect, useRef } from "react";
import { IconButton } from "../../atoms";
import styles from "./DialogShell.module.css";

export type DialogShellSize = "sm" | "md" | "lg";
export type DialogShellMode = "modal" | "inline";

export interface DialogShellProps extends Omit<HTMLAttributes<HTMLDialogElement>, "title"> {
	open: boolean;
	onClose: () => void;
	title: ReactNode;
	actions?: ReactNode;
	footer?: ReactNode;
	size?: DialogShellSize;
	/** "inline" renders the dialog in-flow (no showModal) — for stories/visual diffing. */
	mode?: DialogShellMode;
	children?: ReactNode;
}

export function DialogShell({
	open,
	onClose,
	title,
	actions,
	footer,
	size = "md",
	mode = "modal",
	className,
	children,
	...rest
}: DialogShellProps) {
	const ref = useRef<HTMLDialogElement>(null);

	useEffect(() => {
		if (mode !== "modal") return;
		const dialog = ref.current;
		if (!dialog) return;
		if (open && !dialog.open) {
			dialog.showModal();
		} else if (!open && dialog.open) {
			dialog.close();
		}
	}, [open, mode]);

	if (mode === "inline" && !open) return null;

	return (
		<dialog
			ref={ref}
			className={[
				styles.root,
				styles[size],
				mode === "inline" ? styles.inline : "",
				className ?? "",
			]
				.filter(Boolean)
				.join(" ")}
			open={mode === "inline" ? true : undefined}
			onCancel={(event) => {
				event.preventDefault();
				onClose();
			}}
			data-rag-layout="DialogShell"
			data-size={size}
			data-mode={mode}
			{...rest}
		>
			<header className={styles.header}>
				<span className={styles.title}>{title}</span>
				<span className={styles.headerActions}>
					{actions}
					<IconButton className={styles.close} label="Close dialog" onClick={onClose}>
						×
					</IconButton>
				</span>
			</header>
			<div className={styles.body}>{children}</div>
			{footer && <footer className={styles.footer}>{footer}</footer>}
		</dialog>
	);
}
