import type { SelectInputWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { SelectInput } from "./SelectInput";

export const selectInputWidget = defineWidget<SelectInputWidgetProps>({
	type: "SelectInput",
	module: "ui.dsl",
	render: (props, children, ctx) => {
		const selectedOption = (props.options ?? []).find((option) => option.selected);
		const valueProps =
			props.value !== undefined
				? { value: props.value }
				: { defaultValue: props.defaultValue ?? selectedOption?.value };

		return (
			<SelectInput
				className={props.className}
				name={props.name}
				{...valueProps}
				disabled={props.disabled}
				required={props.required}
				aria-invalid={props.ariaInvalid || undefined}
			>
				{(props.options ?? []).map((option) => (
					<option key={String(option.value)} value={option.value} disabled={option.disabled}>
						{ctx.renderValue(option.label)}
					</option>
				))}
				{children}
			</SelectInput>
		);
	},
});
