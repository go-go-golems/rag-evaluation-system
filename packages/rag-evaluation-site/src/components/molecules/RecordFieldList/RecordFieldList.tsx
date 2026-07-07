import type { ReactNode } from "react";
import type { ContextStyleSet } from "../../../context";
import type { FieldOption, FieldType, FieldValue } from "../../../crm/types";
import { type FieldMode, type FieldRef, FieldRenderer } from "../FieldRenderer";
import styles from "./RecordFieldList.module.css";

/** One field's view configuration (the React-side analogue of a FieldSpec). */
export interface RecordFieldSpecView {
	key: string;
	type: FieldType;
	label?: ReactNode;
	options?: FieldOption[];
	relatedObject?: string;
	readOnly?: boolean;
	unit?: string;
	styleSet?: ContextStyleSet;
}

export interface RecordFieldListSection {
	label?: ReactNode;
	fields: RecordFieldSpecView[];
}

export interface RecordFieldListProps {
	/** The record's field values, keyed by FieldSpec.key. */
	values: Record<string, FieldValue>;
	sections: RecordFieldListSection[];
	mode?: FieldMode;
	/** Keys that failed validation (edit mode). */
	invalidKeys?: string[];
	resolveRef?: (id: string) => FieldRef | undefined;
	onFieldChange?: (key: string, next: FieldValue) => void;
	onFieldCommit?: (key: string) => void;
	/** Layout of a single label+control row. */
	rowLayout?: "inline" | "stacked";
}

/**
 * Lays out many fields as label + control rows, grouped into sections, and hands
 * each one a `FieldRenderPayload`. It owns arrangement only; the `FieldRenderer`
 * owns rendering. Above it sit the record presets (contactRecord, dealRecord)
 * that supply the field specs and values.
 */
export function RecordFieldList({
	values,
	sections,
	mode = "read",
	invalidKeys,
	resolveRef,
	onFieldChange,
	onFieldCommit,
	rowLayout = "inline",
}: RecordFieldListProps) {
	const invalid = new Set(invalidKeys ?? []);
	return (
		<div
			className={styles.root}
			data-rag-molecule="RecordFieldList"
			data-mode={mode}
			data-layout={rowLayout}
		>
			{sections.map((section, sIndex) => (
				<section
					key={section.label != null ? String(section.label) : `s${sIndex}`}
					className={styles.section}
				>
					{section.label != null ? <h4 className={styles.sectionLabel}>{section.label}</h4> : null}
					<dl className={styles.rows}>
						{section.fields.map((field) => (
							<div key={field.key} className={styles.row}>
								<dt className={styles.label}>{field.label ?? field.key}</dt>
								<dd className={styles.control}>
									<FieldRenderer
										fieldKey={field.key}
										type={field.type}
										value={values[field.key] ?? null}
										mode={mode}
										label={field.label ?? field.key}
										options={field.options}
										relatedObject={field.relatedObject}
										readOnly={field.readOnly}
										unit={field.unit}
										styleSet={field.styleSet}
										invalid={invalid.has(field.key)}
										resolveRef={resolveRef}
										onChange={onFieldChange ? (next) => onFieldChange(field.key, next) : undefined}
										onCommit={onFieldCommit ? () => onFieldCommit(field.key) : undefined}
									/>
								</dd>
							</div>
						))}
					</dl>
				</section>
			))}
		</div>
	);
}
