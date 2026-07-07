import type { ChangeEvent, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import type { FieldOption, FieldType, FieldValue } from "../../../crm/types";
import { SelectInput, TextareaInput, TextInput } from "../../atoms";
import styles from "./FieldRenderer.module.css";

export type FieldMode = "read" | "edit";

/** How a relation/user id resolves to something displayable. */
export interface FieldRef {
	label: string;
	avatarUrl?: string;
	href?: string;
}

/**
 * The stable payload every field receives — the CRM analogue of
 * `MatrixCellPayload`. It is the seam that keeps the renderer domain-blind: the
 * renderer owns "how is a typed value shown and edited", and knows nothing about
 * which object the field belongs to. Any control honoring this shape is a valid
 * field control.
 */
export interface FieldRenderPayload {
	/** Which key in record.fields this value came from. */
	fieldKey: string;
	type: FieldType;
	value: FieldValue;
	mode: FieldMode;
	label?: ReactNode;
	options?: FieldOption[];
	relatedObject?: string;
	readOnly?: boolean;
	invalid?: boolean;
	unit?: string;
	/** Palette for select/tag values (reuse the ContextStyleSet contract). */
	styleSet?: ContextStyleSet;
	/** edit mode reports changes back. */
	onChange?: (next: FieldValue) => void;
	/** blur / enter. */
	onCommit?: () => void;
	/** Resolve a relation/user id to a display label + avatar. */
	resolveRef?: (id: string) => FieldRef | undefined;
}

export type FieldRendererProps = FieldRenderPayload;

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function formatDate(iso: string, withTime: boolean): string {
	const date = iso.slice(0, 10).split("-");
	if (date.length !== 3) return iso;
	const [y, m, d] = date;
	const base = `${MONTHS[Number(m) - 1] ?? "?"} ${Number(d)}, ${y}`;
	if (withTime && iso.length >= 16) return `${base} · ${iso.slice(11, 16)}`;
	return base;
}

function formatCurrency(value: number, unit?: string): string {
	const symbol = unit === "EUR" ? "€" : unit === "GBP" ? "£" : "$";
	return `${symbol}${value.toLocaleString("en-US")}`;
}

function asString(value: FieldValue): string {
	if (value == null) return "";
	if (Array.isArray(value)) return value.join(", ");
	if (typeof value === "object") return Object.values(value).filter(Boolean).join(", ");
	return String(value);
}

function toArray(value: FieldValue): string[] {
	if (Array.isArray(value)) return value.map(String);
	if (value == null || value === "") return [];
	return [String(value)];
}

/** A colored pill for a select/tag value, tinted via the ContextStyleSet. */
function OptionPill({
	value,
	options,
	styleSet,
}: {
	value: string;
	options?: FieldOption[];
	styleSet?: ContextStyleSet;
}) {
	const option = options?.find((o) => o.value === value);
	const label = option?.label ?? value;
	const styleKey = option?.colorKey;
	const visual = styleKey ? styleSet?.styles[styleKey] : undefined;
	return (
		<span
			className={styles.pill}
			data-styled={visual ? "true" : undefined}
			style={visual ? contextVisualStyleToCssVars(visual) : undefined}
		>
			{label}
		</span>
	);
}

function RefChip({ ref, empty }: { ref?: FieldRef; empty: string }) {
	if (!ref) return <span className={styles.dim}>{empty}</span>;
	const inner = (
		<>
			<span className={styles.avatar} aria-hidden="true">
				{ref.avatarUrl ? <img src={ref.avatarUrl} alt="" /> : ref.label.charAt(0)}
			</span>
			<span>{ref.label}</span>
		</>
	);
	return ref.href ? (
		<a className={styles.chip} href={ref.href}>
			{inner}
		</a>
	) : (
		<span className={styles.chip}>{inner}</span>
	);
}

function ReadField(p: FieldRenderPayload): ReactNode {
	const { type, value, options, styleSet, unit, resolveRef } = p;
	const str = asString(value);
	switch (type) {
		case "email":
			return str ? (
				<a className={styles.link} href={`mailto:${str}`}>
					{str}
				</a>
			) : (
				<span className={styles.dim}>—</span>
			);
		case "phone":
			return str ? (
				<a className={styles.link} href={`tel:${str.replace(/\s+/g, "")}`}>
					{str}
				</a>
			) : (
				<span className={styles.dim}>—</span>
			);
		case "url":
			return str ? (
				<a className={styles.link} href={str} target="_blank" rel="noreferrer">
					{str.replace(/^https?:\/\//, "")}
				</a>
			) : (
				<span className={styles.dim}>—</span>
			);
		case "currency":
			return typeof value === "number" ? (
				<span className={styles.numeric}>{formatCurrency(value, unit)}</span>
			) : (
				<span className={styles.dim}>—</span>
			);
		case "number":
			return value == null || value === "" ? (
				<span className={styles.dim}>—</span>
			) : (
				<span className={styles.numeric}>{str}</span>
			);
		case "percent":
			return value == null || value === "" ? (
				<span className={styles.dim}>—</span>
			) : (
				<span className={styles.numeric}>{Number(value)}%</span>
			);
		case "date":
			return str ? formatDate(str, false) : <span className={styles.dim}>—</span>;
		case "datetime":
			return str ? formatDate(str, true) : <span className={styles.dim}>—</span>;
		case "boolean":
			return <span className={styles.numeric}>{value ? "✓" : "—"}</span>;
		case "select":
			return str ? (
				<OptionPill value={str} options={options} styleSet={styleSet} />
			) : (
				<span className={styles.dim}>—</span>
			);
		case "multiselect":
		case "tags": {
			const items = toArray(value);
			return items.length ? (
				<span className={styles.pillRow}>
					{items.map((v) => (
						<OptionPill key={v} value={v} options={options} styleSet={styleSet} />
					))}
				</span>
			) : (
				<span className={styles.dim}>—</span>
			);
		}
		case "relation":
		case "user": {
			const ids = toArray(value);
			if (!ids.length) return <span className={styles.dim}>—</span>;
			return (
				<span className={styles.pillRow}>
					{ids.map((id) => (
						<RefChip key={id} ref={resolveRef?.(id)} empty={id} />
					))}
				</span>
			);
		}
		case "address":
			return str ? (
				<span className={styles.address}>{str}</span>
			) : (
				<span className={styles.dim}>—</span>
			);
		default:
			return str ? <span>{str}</span> : <span className={styles.dim}>—</span>;
	}
}

function EditField(p: FieldRenderPayload): ReactNode {
	const { type, value, options, onChange, onCommit, invalid, readOnly, unit, fieldKey, label } = p;
	const ariaLabel = typeof label === "string" ? label : fieldKey;
	const commit = () => onCommit?.();
	const invalidProps = invalid ? { "aria-invalid": true as const, "data-invalid": "true" } : {};

	const onText = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
		onChange?.(e.target.value);
	const onNumber = (e: ChangeEvent<HTMLInputElement>) =>
		onChange?.(e.target.value === "" ? null : Number(e.target.value));

	if (type === "longtext" || type === "address") {
		return (
			<TextareaInput
				aria-label={ariaLabel}
				value={asString(value)}
				disabled={readOnly}
				onChange={onText}
				onBlur={commit}
				{...invalidProps}
			/>
		);
	}
	if (type === "boolean") {
		return (
			<input
				type="checkbox"
				aria-label={ariaLabel}
				checked={Boolean(value)}
				disabled={readOnly}
				onChange={(e) => {
					onChange?.(e.target.checked);
					commit();
				}}
			/>
		);
	}
	if (type === "select") {
		return (
			<SelectInput
				aria-label={ariaLabel}
				value={asString(value)}
				disabled={readOnly}
				onChange={(e) => {
					onChange?.(e.target.value);
					commit();
				}}
				{...invalidProps}
			>
				<option value="">—</option>
				{options?.map((o) => (
					<option key={o.value} value={o.value}>
						{o.label}
					</option>
				))}
			</SelectInput>
		);
	}
	if (type === "multiselect" || type === "tags") {
		// v1: comma-separated entry. A richer pill multi-select (TagListInput) can
		// swap in later without changing the payload contract.
		return (
			<TextInput
				aria-label={ariaLabel}
				value={toArray(value).join(", ")}
				disabled={readOnly}
				placeholder="comma, separated, values"
				onChange={(e) =>
					onChange?.(
						e.target.value
							.split(",")
							.map((s) => s.trim())
							.filter(Boolean),
					)
				}
				onBlur={commit}
				{...invalidProps}
			/>
		);
	}
	if (type === "date" || type === "datetime") {
		return (
			<TextInput
				type={type === "date" ? "date" : "datetime-local"}
				aria-label={ariaLabel}
				value={type === "date" ? asString(value).slice(0, 10) : asString(value).slice(0, 16)}
				disabled={readOnly}
				onChange={onText}
				onBlur={commit}
				{...invalidProps}
			/>
		);
	}
	if (type === "number" || type === "currency" || type === "percent") {
		return (
			<span className={styles.numericInput} data-unit={type}>
				{type === "currency" ? (
					<span className={styles.prefix}>{unit === "EUR" ? "€" : "$"}</span>
				) : null}
				<TextInput
					type="number"
					aria-label={ariaLabel}
					value={value == null ? "" : String(value)}
					disabled={readOnly}
					min={type === "percent" ? 0 : undefined}
					max={type === "percent" ? 100 : undefined}
					onChange={onNumber}
					onBlur={commit}
					{...invalidProps}
				/>
				{type === "percent" ? <span className={styles.suffix}>%</span> : null}
			</span>
		);
	}
	// text / email / phone / url / relation / user (relation edit is a stub input)
	const inputType =
		type === "email" ? "email" : type === "url" ? "url" : type === "phone" ? "tel" : "text";
	return (
		<TextInput
			type={inputType}
			aria-label={ariaLabel}
			value={asString(value)}
			disabled={readOnly}
			onChange={onText}
			onBlur={commit}
			{...invalidProps}
		/>
	);
}

/**
 * The field-system engine: renders one typed CRM value in read or edit mode.
 * Switches on `type × mode` to produce the correct control. Reuses the
 * ContextStyleSet palette for select/tag colors. Domain-blind — it does not know
 * whether the field belongs to a Contact, Company, or Deal.
 */
export function FieldRenderer(props: FieldRendererProps) {
	const { mode, type, invalid } = props;
	return (
		<span
			className={styles.root}
			data-rag-molecule="FieldRenderer"
			data-type={type}
			data-mode={mode}
			data-invalid={invalid || undefined}
		>
			{mode === "edit" ? EditField(props) : ReadField(props)}
		</span>
	);
}
