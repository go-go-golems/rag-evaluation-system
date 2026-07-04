import type { ReactNode } from "react";
import { Button } from "../../atoms";
import { Text } from "../../foundation";
import { DialogShell, type DialogShellMode } from "../../layout";
import styles from "./ConfirmDialog.module.css";

export interface ConfirmDialogProps {
	open: boolean;
	title: ReactNode;
	message: ReactNode;
	detail?: ReactNode;
	confirmLabel?: ReactNode;
	cancelLabel?: ReactNode;
	destructive?: boolean;
	onConfirm: () => void;
	onCancel: () => void;
	mode?: DialogShellMode;
	className?: string;
}

export function ConfirmDialog({
	open,
	title,
	message,
	detail,
	confirmLabel = "Confirm",
	cancelLabel = "Cancel",
	destructive = false,
	onConfirm,
	onCancel,
	mode = "modal",
	className,
}: ConfirmDialogProps) {
	return (
		<DialogShell
			open={open}
			onClose={onCancel}
			title={title}
			size="sm"
			mode={mode}
			className={className}
			data-rag-organism="ConfirmDialog"
			data-destructive={destructive || undefined}
			footer={
				<>
					<Button size="compact" onClick={onCancel}>
						{cancelLabel}
					</Button>
					<Button
						size="compact"
						variant={destructive ? "default" : "primary"}
						className={destructive ? styles.destructive : undefined}
						onClick={onConfirm}
					>
						{confirmLabel}
					</Button>
				</>
			}
		>
			<Text>{message}</Text>
			{detail && (
				<Text size="compact" tone={destructive ? "danger" : "muted"} className={styles.detail}>
					{detail}
				</Text>
			)}
		</DialogShell>
	);
}
