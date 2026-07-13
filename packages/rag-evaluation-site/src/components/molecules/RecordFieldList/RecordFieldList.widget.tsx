import type { FieldValue } from "../../../crm/types";
import type {
	FieldSpec,
	RecordFieldListSectionSpec,
	RecordFieldListWidgetProps,
} from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import type { RenderContext } from "../../../widgets/registry";
import { actionValue, makeResolveRef } from "../FieldRenderer/FieldRenderer.widget";
import { RecordFieldList, type RecordFieldSpecView } from "./RecordFieldList";

function toView(spec: FieldSpec, ctx: RenderContext): RecordFieldSpecView {
	return {
		key: spec.key,
		type: spec.type,
		label: spec.label != null ? ctx.renderValue(spec.label) : spec.key,
		options: spec.options,
		relatedObject: spec.relatedObject,
		readOnly: spec.readOnly,
		unit: spec.unit,
		styleSet: spec.styleSet,
	};
}

export const recordFieldListWidget = defineWidget<RecordFieldListWidgetProps>({
	type: "RecordFieldList",
	module: "widget.dsl",
	render: (props, _children, ctx) => {
		const sectionSpecs: RecordFieldListSectionSpec[] =
			props.sections ?? (props.fields ? [{ fields: props.fields }] : []);
		const sections = sectionSpecs.map((section) => ({
			label: section.label != null ? ctx.renderValue(section.label) : undefined,
			fields: section.fields.map((spec) => toView(spec, ctx)),
		}));
		const values = props.values as Record<string, FieldValue>;
		return (
			<RecordFieldList
				values={values}
				sections={sections}
				mode={props.mode ?? "read"}
				rowLayout={props.rowLayout}
				invalidKeys={props.invalidKeys}
				resolveRef={makeResolveRef(props.refs)}
				onFieldChange={
					props.onFieldChangeAction
						? (key, next) =>
								ctx.dispatchAction(props.onFieldChangeAction!, {
									key,
									value: actionValue(next),
									componentType: "RecordFieldList",
								})
						: undefined
				}
			/>
		);
	},
});
