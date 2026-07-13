import { useState } from "react";
import type { SearchFieldWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";
import { defineWidget } from "../../../widgets/registry";
import { SearchField } from "./SearchField";

export const searchFieldWidget = defineWidget<SearchFieldWidgetProps>({
	type: "SearchField",
	module: "data.dsl",
	render: (props, _children, ctx) => <SearchFieldWidgetHost props={props} ctx={ctx} />,
});

// IR pages are stateless; the adapter owns the input value locally and only
// dispatches on submit (Enter). `name` still participates in native form posts.
function SearchFieldWidgetHost({
	props,
	ctx,
}: {
	props: SearchFieldWidgetProps;
	ctx: RenderContext;
}) {
	const [value, setValue] = useState(props.defaultValue ?? "");
	const onSubmitAction = props.onSubmitAction;
	const onClearAction = props.onClearAction;
	return (
		<SearchField
			className={props.className}
			name={props.name}
			value={value}
			onValueChange={setValue}
			placeholder={props.placeholder}
			disabled={props.disabled}
			aria-label={props.resultCount == null ? undefined : `Search, ${props.resultCount} results`}
			onClear={
				onClearAction
					? () =>
							ctx.dispatchAction(onClearAction, {
								query: "",
								value: "",
								componentType: "SearchField",
							})
					: undefined
			}
			onSubmit={
				onSubmitAction
					? (query) =>
							ctx.dispatchAction(onSubmitAction, {
								query,
								value: query,
								componentType: "SearchField",
							})
					: undefined
			}
		/>
	);
}
