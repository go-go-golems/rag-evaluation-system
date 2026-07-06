import type { TextareaInputWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { TextareaInput } from "./TextareaInput";

export const textareaInputWidget = defineWidget<TextareaInputWidgetProps>({
	type: "TextareaInput",
	module: "ui.dsl",
	render: (props) => {
		const readOnly = props.readOnly ?? true;
		const sharedProps = {
			className: props.className,
			name: props.name,
			placeholder: props.placeholder,
			disabled: props.disabled,
			required: props.required,
			minLength: props.minLength,
			maxLength: props.maxLength,
			rows: props.rows,
			resize: props.resize,
			"aria-invalid": props.ariaInvalid || undefined,
		};

		if (readOnly) {
			return <TextareaInput {...sharedProps} value={props.value ?? props.defaultValue} readOnly />;
		}

		return (
			<TextareaInput
				{...sharedProps}
				defaultValue={props.defaultValue ?? props.value}
				readOnly={false}
			/>
		);
	},
});
