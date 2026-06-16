import type { TextareaHTMLAttributes } from "react";
import styles from "./TextareaInput.module.css";

export type TextareaInputResize = "vertical" | "none";

export interface TextareaInputProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
	resize?: TextareaInputResize;
}

export function TextareaInput({ className, resize = "vertical", ...rest }: TextareaInputProps) {
	return (
		<textarea
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-atom="TextareaInput"
			data-resize={resize}
			{...rest}
		/>
	);
}
