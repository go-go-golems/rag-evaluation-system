import { type KeyboardEvent, type ReactNode, useEffect, useMemo, useRef, useState } from "react";
import styles from "./DataTable.module.css";

export interface DataTableColumn<T> {
	id: string;
	header: ReactNode;
	cell: (row: T) => ReactNode;
	align?: "start" | "end" | "center";
	maxWidth?: number | string;
	sortable?: boolean;
	sortDirection?: "ascending" | "descending";
	onSort?: () => void;
}

export interface DataTableKeyboard {
	mode?: "rows";
	selection?: "manual" | "followFocus";
	vimAliases?: boolean;
	enterSelect?: boolean;
}

export interface DataTableCommand {
	id: string;
	key: string;
	label: string;
	danger?: boolean;
}

export interface DataTableProps<T> {
	columns: DataTableColumn<T>[];
	rows: T[];
	getRowKey: (row: T) => string;
	selectedKey?: string | null;
	onRowSelect?: (row: T) => void;
	keyboard?: DataTableKeyboard;
	commands?: DataTableCommand[];
	onCommand?: (command: DataTableCommand, row: T) => void;
	rowTone?: (row: T) => "muted" | "success" | "warning" | "danger" | "accent" | undefined;
	emptyMessage?: ReactNode;
	className?: string;
}

function isEditableTarget(target: EventTarget | null): boolean {
	if (!(target instanceof HTMLElement)) return false;
	return (
		target.isContentEditable ||
		Boolean(target.closest("input, textarea, select, [contenteditable='true'], [role='dialog']"))
	);
}

export function DataTable<T>({
	columns,
	rows,
	getRowKey,
	selectedKey,
	onRowSelect,
	keyboard,
	commands = [],
	onCommand,
	rowTone,
	emptyMessage,
	className,
}: DataTableProps<T>) {
	const keys = useMemo(() => rows.map(getRowKey), [rows, getRowKey]);
	const [focusedKey, setFocusedKey] = useState<string | null>(selectedKey ?? keys[0] ?? null);
	const rowRefs = useRef(new Map<string, HTMLTableRowElement>());
	const shouldRestoreFocus = useRef(false);

	useEffect(() => {
		setFocusedKey((current) => {
			if (current && keys.includes(current)) return current;
			if (selectedKey && keys.includes(selectedKey)) return selectedKey;
			return keys[0] ?? null;
		});
	}, [keys, selectedKey]);

	useEffect(() => {
		if (shouldRestoreFocus.current && focusedKey) {
			rowRefs.current.get(focusedKey)?.focus();
			shouldRestoreFocus.current = false;
		}
	}, [focusedKey]);

	const moveFocus = (delta: number) => {
		if (keys.length === 0) return;
		const current = Math.max(0, keys.indexOf(focusedKey ?? ""));
		const next = Math.min(keys.length - 1, Math.max(0, current + delta));
		shouldRestoreFocus.current = true;
		const nextKey = keys[next];
		const nextRow = rows[next];
		if (nextKey == null || nextRow == null) return;
		setFocusedKey(nextKey);
		if (keyboard?.selection === "followFocus") onRowSelect?.(nextRow);
	};

	const handleKeyDown = (event: KeyboardEvent<HTMLTableRowElement>, row: T) => {
		if (!keyboard || isEditableTarget(event.target)) return;
		if (event.key === "ArrowDown" || (keyboard.vimAliases && event.key.toLowerCase() === "j")) {
			event.preventDefault();
			moveFocus(1);
			return;
		}
		if (event.key === "ArrowUp" || (keyboard.vimAliases && event.key.toLowerCase() === "k")) {
			event.preventDefault();
			moveFocus(-1);
			return;
		}
		if (event.key === "Enter" && keyboard.enterSelect !== false) {
			event.preventDefault();
			onRowSelect?.(row);
			return;
		}
		const command = commands.find(
			(candidate) => candidate.key.toLowerCase() === event.key.toLowerCase(),
		);
		if (command) {
			event.preventDefault();
			onCommand?.(command, row);
		}
	};

	return (
		<table
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-component="DataTable"
			data-rag-keyboard-scope={keyboard ? "DataTable" : undefined}
		>
			<thead>
				<tr>
					{columns.map((column) => (
						<th
							key={column.id}
							aria-sort={column.sortDirection ?? "none"}
							className={styles[column.align ?? "start"]}
							style={column.maxWidth ? { maxWidth: column.maxWidth } : undefined}
						>
							{column.sortable ? (
								<button className={styles.sortButton} onClick={column.onSort} type="button">
									{column.header}
									<span aria-hidden="true" className={styles.sortIndicator}>
										{column.sortDirection === "ascending"
											? "↑"
											: column.sortDirection === "descending"
												? "↓"
												: "↕"}
									</span>
								</button>
							) : (
								column.header
							)}
						</th>
					))}
				</tr>
			</thead>
			<tbody>
				{rows.length === 0 && emptyMessage && (
					<tr>
						<td colSpan={columns.length} className={styles.empty}>
							{emptyMessage}
						</td>
					</tr>
				)}
				{rows.map((row) => {
					const key = getRowKey(row);
					const focused = Boolean(keyboard) && focusedKey === key;
					const tone = rowTone?.(row);
					return (
						<tr
							key={key}
							className={[
								onRowSelect ? styles.selectable : "",
								selectedKey === key ? styles.selected : "",
								focused ? styles.focused : "",
								tone ? styles[`tone-${tone}`] : "",
							]
								.filter(Boolean)
								.join(" ")}
							ref={(element) => {
								if (element) rowRefs.current.set(key, element);
								else rowRefs.current.delete(key);
							}}
							aria-selected={selectedKey === key}
							tabIndex={keyboard ? (focused ? 0 : -1) : undefined}
							onFocus={() => setFocusedKey(key)}
							onKeyDown={(event) => handleKeyDown(event, row)}
							onClick={onRowSelect ? () => onRowSelect(row) : undefined}
						>
							{columns.map((column) => (
								<td
									key={column.id}
									className={styles[column.align ?? "start"]}
									style={column.maxWidth ? { maxWidth: column.maxWidth } : undefined}
								>
									{column.cell(row)}
								</td>
							))}
						</tr>
					);
				})}
			</tbody>
			{commands.length > 0 && (
				<caption className={styles.commandHelp}>
					{commands.map((command) => `${command.key}: ${command.label}`).join(" · ")}
				</caption>
			)}
		</table>
	);
}
