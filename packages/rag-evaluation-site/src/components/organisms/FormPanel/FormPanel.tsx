import type { FormHTMLAttributes, ReactNode } from "react";
import { Button } from "../../atoms";
import { Caption, StatusText } from "../../foundation";
import { Inline, Panel, Stack } from "../../layout";
import styles from "./FormPanel.module.css";

export type FormPanelStatus = "idle" | "saving" | "success" | "error";

export interface FormPanelProps extends Omit<FormHTMLAttributes<HTMLFormElement>, "title"> {
	title: ReactNode;
	subtitle?: ReactNode;
	actions?: ReactNode;
	status?: FormPanelStatus;
	statusMessage?: ReactNode;
	submitLabel?: ReactNode;
	resetLabel?: ReactNode;
	footer?: ReactNode;
	disabled?: boolean;
	children?: ReactNode;
}

const statusKind: Record<Exclude<FormPanelStatus, "idle">, string> = {
	saving: "running",
	success: "succeeded",
	error: "error",
};

const defaultStatusMessage: Record<Exclude<FormPanelStatus, "idle">, string> = {
	saving: "Saving…",
	success: "Saved.",
	error: "Could not save.",
};

export function FormPanel({
	title,
	subtitle,
	actions,
	status = "idle",
	statusMessage,
	submitLabel = "Save",
	resetLabel,
	footer,
	disabled = false,
	className,
	children,
	...rest
}: FormPanelProps) {
	const resolvedSubmitLabel = submitLabel ?? "Save";
	const showStatus = status !== "idle" || Boolean(statusMessage);
	const normalizedStatus = status === "idle" ? "ready" : statusKind[status];
	const resolvedStatusMessage =
		statusMessage ?? (status === "idle" ? undefined : defaultStatusMessage[status]);
	const isSaving = status === "saving";
	const controlsDisabled = disabled || isSaving;

	return (
		<Panel title={title} actions={actions} density="condensed" className={className}>
			<form
				className={styles.root}
				data-rag-organism="FormPanel"
				data-status={status}
				aria-busy={isSaving || undefined}
				{...rest}
			>
				<Stack gap="sm">
					{subtitle && <Caption tone="muted">{subtitle}</Caption>}
					{children}
					<div className={styles.footer}>
						<div className={styles.status} aria-live={status === "error" ? "assertive" : "polite"}>
							{showStatus && resolvedStatusMessage && (
								<StatusText status={normalizedStatus} icon>
									{resolvedStatusMessage}
								</StatusText>
							)}
						</div>
						<Inline gap="xs" justify="end" wrap>
							{resetLabel && (
								<Button type="reset" size="compact" disabled={controlsDisabled}>
									{resetLabel}
								</Button>
							)}
							<Button type="submit" size="normal" variant="primary" disabled={controlsDisabled}>
								{resolvedSubmitLabel}
							</Button>
						</Inline>
					</div>
					{footer && <div className={styles.extraFooter}>{footer}</div>}
				</Stack>
			</form>
		</Panel>
	);
}
