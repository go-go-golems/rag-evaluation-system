import { useEffect, useId, useRef, useState } from "react";
import { Button, CheckboxRow } from "../../atoms";
import { CodeText, Text } from "../../foundation";
import styles from "./KeyboardShortcutHelp.module.css";

export interface KeyboardShortcutHelpItem {
	id: string;
	label: string;
	chord: string;
}

export interface KeyboardShortcutHelpProps {
	items: KeyboardShortcutHelpItem[];
	enabled: boolean;
	onEnabledChange: (enabled: boolean) => void;
}

export function KeyboardShortcutHelp({
	items,
	enabled,
	onEnabledChange,
}: KeyboardShortcutHelpProps) {
	const [open, setOpen] = useState(false);
	const dialogRef = useRef<HTMLDialogElement>(null);
	const titleId = useId();

	useEffect(() => {
		const dialog = dialogRef.current;
		if (!dialog) return;
		if (open && !dialog.open) dialog.showModal();
		if (!open && dialog.open) dialog.close();
	}, [open]);

	if (items.length === 0) return null;

	return (
		<div className={styles.root} data-rag-component="KeyboardShortcutHelp">
			<Button
				size="compact"
				aria-expanded={open}
				aria-haspopup="dialog"
				onClick={() => setOpen(true)}
			>
				Keyboard shortcuts
			</Button>
			<dialog
				ref={dialogRef}
				className={styles.dialog}
				aria-labelledby={titleId}
				onClose={() => setOpen(false)}
			>
				<div className={styles.header}>
					<Text id={titleId} as="strong" size="compact" weight="bold">
						Keyboard shortcuts
					</Text>
					<Button autoFocus size="compact" onClick={() => setOpen(false)}>
						Close
					</Button>
				</div>
				<ul className={styles.list}>
					{items.map((item) => (
						<li key={item.id} className={styles.item}>
							<Text as="span" size="compact">
								{item.label}
							</Text>
							<CodeText>{item.chord}</CodeText>
						</li>
					))}
				</ul>
				<CheckboxRow
					checked={enabled}
					onChange={(event) => onEnabledChange(event.currentTarget.checked)}
				>
					<Text as="span" size="compact">
						Enable page keyboard shortcuts
					</Text>
				</CheckboxRow>
			</dialog>
		</div>
	);
}
