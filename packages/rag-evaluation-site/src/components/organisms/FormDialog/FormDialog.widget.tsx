import { useCallback, useEffect, useRef, useState } from "react";
import { Button } from "../../atoms";
import { Caption } from "../../foundation";
import { Inline, Stack } from "../../layout";
import type { WidgetActionContext } from "../../../widgets/actions";
import type { ActionSpec, FormDialogWidgetProps } from "../../../widgets/ir";
import { defineWidget, type RenderContext } from "../../../widgets/registry";
import styles from "./FormDialog.module.css";

type OverlayEvent = CustomEvent<{
	action: Extract<ActionSpec, { kind: "openOverlay" | "closeOverlay" }>;
	context: WidgetActionContext;
}>;
type ResultEvent = CustomEvent<{
	context: WidgetActionContext;
	responseOk: boolean;
	result?: { ok?: boolean; error?: string; fieldErrors?: Record<string, string> };
}>;

export const formDialogWidget = defineWidget<FormDialogWidgetProps>({
	type: "FormDialog",
	module: "widget.dsl",
	render: (props, _children, ctx) => <FormDialogHost props={props} ctx={ctx} />,
});

function FormDialogHost({ props, ctx }: { props: FormDialogWidgetProps; ctx: RenderContext }) {
	const dialogRef = useRef<HTMLDialogElement>(null);
	const openerRef = useRef<HTMLElement | null>(null);
	const frozenContext = useRef<WidgetActionContext>({});
	const [error, setError] = useState<string>();
	const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

	const close = useCallback(() => {
		dialogRef.current?.close();
		openerRef.current?.focus();
	}, []);

	useEffect(() => {
		const onOverlay = (event: Event) => {
			const { action, context } = (event as OverlayEvent).detail;
			if (action.kind === "openOverlay" && action.target === props.id) {
				openerRef.current =
					document.activeElement instanceof HTMLElement ? document.activeElement : null;
				frozenContext.current = { ...context, overlayId: props.id };
				setError(undefined);
				setFieldErrors({});
				dialogRef.current?.showModal();
				requestAnimationFrame(() =>
					dialogRef.current
						?.querySelector<HTMLElement>(
							props.initialFocus
								? `[name="${CSS.escape(props.initialFocus)}"]`
								: "input, textarea, select, button",
						)
						?.focus(),
				);
			}
			if (action.kind === "closeOverlay" && (!action.target || action.target === props.id)) close();
		};
		const onResult = (event: Event) => {
			const detail = (event as ResultEvent).detail;
			if (detail.context.overlayId !== props.id) return;
			if (detail.responseOk && detail.result?.ok !== false && !detail.result?.fieldErrors) close();
			else {
				setError(detail.result?.error ?? "Could not save. Check the form and try again.");
				setFieldErrors(detail.result?.fieldErrors ?? {});
			}
		};
		window.addEventListener("widget:overlay", onOverlay);
		window.addEventListener("widget:action-result", onResult);
		return () => {
			window.removeEventListener("widget:overlay", onOverlay);
			window.removeEventListener("widget:action-result", onResult);
		};
	}, [props.id, props.initialFocus, close]);

	return (
		<dialog
			ref={dialogRef}
			className={styles.dialog}
			aria-labelledby={`${props.id}-title`}
			onCancel={(event) => {
				event.preventDefault();
				close();
			}}
		>
			<form
				method="dialog"
				onSubmit={(event) => {
					event.preventDefault();
					if (!event.currentTarget.reportValidity()) return;
					const form = Object.fromEntries(new FormData(event.currentTarget).entries());
					ctx.dispatchAction(props.onSubmitAction, {
						...frozenContext.current,
						form,
						overlayId: props.id,
						componentType: "FormDialog",
					});
				}}
			>
				<Stack gap="sm">
					<h2 id={`${props.id}-title`} className={styles.title}>
						{ctx.renderValue(props.title)}
					</h2>
					{props.body ? ctx.renderNode(props.body) : null}
					{Object.keys(fieldErrors).length > 0 && (
						<div className={styles.errors}>
							{Object.entries(fieldErrors).map(([name, message]) => (
								<Caption key={name}>
									{name}: {message}
								</Caption>
							))}
						</div>
					)}
					{error && (
						<div role="alert" className={styles.error}>
							{error}
						</div>
					)}
					<Inline gap="xs" justify="end">
						<Button type="button" onClick={close}>
							{ctx.renderValue(props.cancelLabel ?? "Cancel")}
						</Button>
						<Button type="submit" variant="primary">
							{ctx.renderValue(props.submitLabel ?? "Save")}
						</Button>
					</Inline>
				</Stack>
			</form>
		</dialog>
	);
}
