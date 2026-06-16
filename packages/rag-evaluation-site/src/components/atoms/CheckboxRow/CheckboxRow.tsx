import type { InputHTMLAttributes, LabelHTMLAttributes, ReactNode } from "react";
import styles from "./CheckboxRow.module.css";

export interface CheckboxRowProps extends Omit<LabelHTMLAttributes<HTMLLabelElement>, "onChange"> {
	checked: boolean;
	onChange: InputHTMLAttributes<HTMLInputElement>["onChange"];
	disabled?: boolean;
	children?: ReactNode;
}

export function CheckboxRow({
	checked,
	onChange,
	disabled = false,
	className,
	children,
	...rest
}: CheckboxRowProps) {
	return (
		<label
			className={[styles.root, disabled ? styles.disabled : "", className ?? ""]
				.filter(Boolean)
				.join(" ")}
			data-rag-atom="CheckboxRow"
			data-disabled={disabled ? "true" : undefined}
			{...rest}
		>
			<input type="checkbox" checked={checked} onChange={onChange} disabled={disabled} />
			{children}
		</label>
	);
}
