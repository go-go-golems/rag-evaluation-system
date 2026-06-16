import { TextInput } from "./TextInput";
import { defineWidget } from "../../../widgets/registry";
import type { TextInputWidgetProps } from "../../../widgets/ir";

export const textInputWidget = defineWidget<TextInputWidgetProps>({
	type: "TextInput",
	module: "ui.dsl",
	render: (props) => {
		const readOnly = props.readOnly ?? true;
		const sharedProps = {
			className: props.className,
			name: props.name,
			placeholder: props.placeholder,
			type: props.type,
			disabled: props.disabled,
			required: props.required,
			min: props.min,
			max: props.max,
			minLength: props.minLength,
			maxLength: props.maxLength,
			autoComplete: props.autoComplete,
			"aria-invalid": props.ariaInvalid || undefined,
		};

		if (readOnly) {
			return <TextInput {...sharedProps} value={props.value ?? props.defaultValue} readOnly />;
		}

		return (
			<TextInput
				{...sharedProps}
				defaultValue={props.defaultValue ?? props.value}
				readOnly={false}
			/>
		);
	},
});
