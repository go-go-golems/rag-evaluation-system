import type { FieldValue } from "../../../crm/types";
import type { FieldRendererWidgetProps, FieldRefSpec, FieldSpec } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import type { RenderContext } from "../../../widgets/registry";
import { type FieldRef, FieldRenderer } from "./FieldRenderer";

/** Coerce a possibly-complex FieldValue to the primitive an ActionContext accepts. */
export function actionValue(value: FieldValue): string | number | boolean | null {
	if (value == null) return null;
	if (typeof value === "object") return JSON.stringify(value);
	return value;
}

export function makeResolveRef(
	refs: Record<string, FieldRefSpec> | undefined,
): ((id: string) => FieldRef | undefined) | undefined {
	if (!refs) return undefined;
	return (id) => refs[id];
}

export function renderFieldSpec(
	spec: FieldSpec,
	value: FieldValue,
	props: FieldRendererWidgetProps,
	ctx: RenderContext,
) {
	return (
		<FieldRenderer
			fieldKey={spec.key}
			type={spec.type}
			value={value}
			mode={props.mode ?? "read"}
			label={spec.label != null ? ctx.renderValue(spec.label) : spec.key}
			options={spec.options}
			relatedObject={spec.relatedObject}
			readOnly={spec.readOnly}
			unit={spec.unit}
			styleSet={spec.styleSet}
			invalid={props.invalid}
			resolveRef={makeResolveRef(props.refs)}
			onChange={
				props.onChangeAction
					? (next) =>
							ctx.dispatchAction(props.onChangeAction!, {
								key: spec.key,
								value: actionValue(next),
								componentType: "FieldRenderer",
							})
					: undefined
			}
		/>
	);
}

export const fieldRendererWidget = defineWidget<FieldRendererWidgetProps>({
	type: "FieldRenderer",
	module: "crm.dsl",
	render: (props, _children, ctx) => renderFieldSpec(props.spec, props.value, props, ctx),
});
