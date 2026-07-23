import {
	type ChangeEvent,
	type KeyboardEvent,
	type MouseEvent,
	type ReactNode,
	useEffect,
	useMemo,
	useRef,
	useState,
} from "react";
import { Button } from "../../atoms/Button";
import { Caption } from "../../foundation/Caption";
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

export type DataTableSelectionReason =
	| "toggle"
	| "range"
	| "selectAll"
	| "clearAll"
	| "keyboardToggle"
	| "keyboardRange"
	| "clear";

export interface DataTableBulkAction {
	id: string;
	label: ReactNode;
	danger?: boolean;
	disabled?: boolean;
	onInvoke: (selectedKeys: readonly string[]) => void;
}

export interface DataTableMultiSelection {
	mode: "multi";
	selectedKeys: readonly string[];
	onSelectionChange: (nextKeys: string[], reason: DataTableSelectionReason) => void;
	bulkActions?: readonly DataTableBulkAction[];
	ariaLabel?: string;
}

export interface DataTableProps<T> {
	columns: DataTableColumn<T>[];
	rows: T[];
	getRowKey: (row: T) => string;
	selectedKey?: string | null;
	onRowSelect?: (row: T) => void;
	multiSelection?: DataTableMultiSelection;
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

function orderedSelection(keys: readonly string[], selected: Iterable<string>): string[] {
	const selectedSet = new Set(selected);
	return keys.filter((key) => selectedSet.has(key));
}

function rangeKeys(keys: readonly string[], start: string, end: string): string[] {
	const startIndex = keys.indexOf(start);
	const endIndex = keys.indexOf(end);
	if (startIndex < 0 || endIndex < 0) return [end];
	return keys.slice(Math.min(startIndex, endIndex), Math.max(startIndex, endIndex) + 1);
}

function IndeterminateCheckbox({
	checked,
	indeterminate = false,
	label,
	onChange,
	onClick,
}: {
	checked: boolean;
	indeterminate?: boolean;
	label: string;
	onChange: (event: ChangeEvent<HTMLInputElement>) => void;
	onClick?: (event: MouseEvent<HTMLInputElement>) => void;
}) {
	const ref = useRef<HTMLInputElement>(null);
	useEffect(() => {
		if (ref.current) ref.current.indeterminate = indeterminate;
	}, [indeterminate]);
	return (
		<input
			ref={ref}
			aria-label={label}
			checked={checked}
			onChange={onChange}
			onClick={onClick}
			type="checkbox"
		/>
	);
}

export function DataTable<T>({
	columns,
	rows,
	getRowKey,
	selectedKey,
	onRowSelect,
	multiSelection,
	keyboard,
	commands = [],
	onCommand,
	rowTone,
	emptyMessage,
	className,
}: DataTableProps<T>) {
	const keys = useMemo(() => rows.map(getRowKey), [rows, getRowKey]);
	const [focusedKey, setFocusedKey] = useState<string | null>(selectedKey ?? keys[0] ?? null);
	const [anchorKey, setAnchorKey] = useState<string | null>(null);
	const rowRefs = useRef(new Map<string, HTMLTableRowElement>());
	const shouldRestoreFocus = useRef(false);
	const selectedKeys = useMemo(
		() => orderedSelection(keys, multiSelection?.selectedKeys ?? []),
		[keys, multiSelection?.selectedKeys],
	);
	const selectedSet = useMemo(() => new Set(selectedKeys), [selectedKeys]);
	const multi = multiSelection?.mode === "multi";
	const hasSelection = selectedKeys.length > 0;

	useEffect(() => {
		setFocusedKey((current) => {
			if (current && keys.includes(current)) return current;
			if (!multi && selectedKey && keys.includes(selectedKey)) return selectedKey;
			return keys[0] ?? null;
		});
		setAnchorKey((current) => (current && keys.includes(current) ? current : null));
	}, [keys, multi, selectedKey]);

	useEffect(() => {
		if (shouldRestoreFocus.current && focusedKey) {
			rowRefs.current.get(focusedKey)?.focus();
			shouldRestoreFocus.current = false;
		}
	}, [focusedKey]);

	const emitSelection = (next: Iterable<string>, reason: DataTableSelectionReason) => {
		multiSelection?.onSelectionChange(orderedSelection(keys, next), reason);
	};
	const toggleSelection = (key: string, reason: DataTableSelectionReason) => {
		const next = new Set(selectedSet);
		if (next.has(key)) next.delete(key);
		else next.add(key);
		setAnchorKey(key);
		emitSelection(next, reason);
	};
	const extendSelection = (key: string, reason: "range" | "keyboardRange") => {
		const start = anchorKey && keys.includes(anchorKey) ? anchorKey : (focusedKey ?? key);
		setAnchorKey(start);
		emitSelection(new Set([...selectedSet, ...rangeKeys(keys, start, key)]), reason);
	};
	const moveFocus = (delta: number, extend = false) => {
		if (keys.length === 0) return;
		const current = Math.max(0, keys.indexOf(focusedKey ?? ""));
		const next = Math.min(keys.length - 1, Math.max(0, current + delta));
		const nextKey = keys[next];
		const nextRow = rows[next];
		if (nextKey == null || nextRow == null) return;
		shouldRestoreFocus.current = true;
		setFocusedKey(nextKey);
		if (multi && extend) extendSelection(nextKey, "keyboardRange");
		else if (!multi && keyboard?.selection === "followFocus") onRowSelect?.(nextRow);
	};

	const handleKeyDown = (event: KeyboardEvent<HTMLTableRowElement>, row: T) => {
		if (!keyboard || isEditableTarget(event.target)) return;
		if (event.ctrlKey || event.metaKey || event.altKey) return;
		if (event.key === "ArrowDown" || (keyboard.vimAliases && event.key.toLowerCase() === "j")) {
			event.preventDefault();
			moveFocus(1, multi && event.shiftKey);
			return;
		}
		if (event.key === "ArrowUp" || (keyboard.vimAliases && event.key.toLowerCase() === "k")) {
			event.preventDefault();
			moveFocus(-1, multi && event.shiftKey);
			return;
		}
		if (multi && event.key === " ") {
			event.preventDefault();
			toggleSelection(getRowKey(row), "keyboardToggle");
			return;
		}
		if (multi && event.key === "Escape") {
			event.preventDefault();
			setAnchorKey(null);
			emitSelection([], "clear");
			return;
		}
		if (!multi && event.key === "Enter" && keyboard.enterSelect !== false) {
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

	const allSelected = keys.length > 0 && selectedKeys.length === keys.length;
	return (
		<div className={styles.container} data-rag-component="DataTable">
			{multi && hasSelection && (
				<div
					aria-label="Bulk table actions"
					role="toolbar"
					className={styles.bulkActions}
					data-rag-component="DataTableBulkActions"
				>
					<Caption>{selectedKeys.length} selected</Caption>
					{multiSelection.bulkActions?.map((action) => (
						<Button
							disabled={action.disabled}
							key={action.id}
							onClick={() => action.onInvoke(selectedKeys)}
							size="compact"
							variant={action.danger ? "danger" : "default"}
						>
							{action.label}
						</Button>
					))}
					<Button
						onClick={() => {
							setAnchorKey(null);
							emitSelection([], "clear");
						}}
						size="compact"
					>
						Clear selection
					</Button>
				</div>
			)}
			<table
				className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
				data-rag-keyboard-scope={keyboard ? "DataTable" : undefined}
			>
				<thead>
					<tr>
						{multi && (
							<th className={styles.selectionCell}>
								<IndeterminateCheckbox
									checked={allSelected}
									indeterminate={!allSelected && hasSelection}
									label={allSelected ? "Clear visible row selection" : "Select all visible rows"}
									onChange={() => {
										setAnchorKey(keys[0] ?? null);
										emitSelection(allSelected ? [] : keys, allSelected ? "clearAll" : "selectAll");
									}}
								/>
							</th>
						)}
						{columns.map((column) => (
							<th
								aria-sort={column.sortDirection ?? "none"}
								className={styles[column.align ?? "start"]}
								key={column.id}
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
							<td colSpan={columns.length + (multi ? 1 : 0)} className={styles.empty}>
								{emptyMessage}
							</td>
						</tr>
					)}
					{rows.map((row) => {
						const key = getRowKey(row);
						const focused = Boolean(keyboard) && focusedKey === key;
						const tone = rowTone?.(row);
						const selected = multi ? selectedSet.has(key) : selectedKey === key;
						return (
							<tr
								aria-selected={selected}
								className={[
									!multi && onRowSelect ? styles.selectable : "",
									selected ? styles.selected : "",
									focused ? styles.focused : "",
									tone ? styles[`tone-${tone}`] : "",
								]
									.filter(Boolean)
									.join(" ")}
								key={key}
								onClick={!multi && onRowSelect ? () => onRowSelect(row) : undefined}
								onFocus={() => setFocusedKey(key)}
								onKeyDown={(event) => handleKeyDown(event, row)}
								ref={(element) => {
									if (element) rowRefs.current.set(key, element);
									else rowRefs.current.delete(key);
								}}
								tabIndex={keyboard ? (focused ? 0 : -1) : undefined}
							>
								{multi && (
									<td className={styles.selectionCell}>
										<IndeterminateCheckbox
											checked={selected}
											label={`${selected ? "Deselect" : "Select"} row`}
											onChange={() => {}}
											onClick={(event) => {
												event.stopPropagation();
												if (event.shiftKey) extendSelection(key, "range");
												else toggleSelection(key, "toggle");
											}}
										/>
									</td>
								)}
								{columns.map((column) => (
									<td
										className={styles[column.align ?? "start"]}
										key={column.id}
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
		</div>
	);
}
