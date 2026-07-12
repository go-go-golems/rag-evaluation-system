import { useState } from "react";
import type { MarkdownEditorWidgetProps } from "../../../widgets/ir";
import type { RenderContext } from "../../../widgets/registry";
import { defineWidget } from "../../../widgets/registry";
import { Button } from "../../atoms";
import { ScrollRegion, SplitPane } from "../../layout";
import { MarkdownArticle } from "../MarkdownArticle";
import { MarkdownEditor } from "./MarkdownEditor";

export const markdownEditorWidget = defineWidget<MarkdownEditorWidgetProps>({
	type: "MarkdownEditor",
	module: "cms.dsl",
	render: (props, _children, ctx) => <MarkdownEditorWidgetHost props={props} ctx={ctx} />,
});

// The editor value is browser-local state: the toolbar and the live
// MarkdownArticle preview work without any server round-trip, while the
// named textarea inside MarkdownEditor carries the value in native form
// posts (formPanel({method:"post", formAction})).
function MarkdownEditorWidgetHost({
	props,
	ctx,
}: {
	props: MarkdownEditorWidgetProps;
	ctx: RenderContext;
}) {
	const [value, setValue] = useState(props.defaultValue ?? "");
	const onSubmitAction = props.onSubmitAction;

	const updateValue = (next: string) => {
		setValue(next);
		if (props.onChangeAction) {
			ctx.dispatchAction(props.onChangeAction, {
				value: next,
				componentType: "MarkdownEditor",
			});
		}
	};

	const editor = (
		<MarkdownEditor
			className={props.className}
			name={props.name}
			value={value}
			onValueChange={updateValue}
			minRows={props.minRows}
			maxLength={props.maxLength}
			disabled={props.disabled}
			textareaAriaLabel={props.textareaAriaLabel}
			toolbarSlot={
				onSubmitAction ? (
					<Button
						size="compact"
						onClick={() =>
							ctx.dispatchAction(onSubmitAction, {
								value,
								componentType: "MarkdownEditor",
							})
						}
					>
						Save
					</Button>
				) : undefined
			}
		/>
	);

	if (props.preview === "hidden") return editor;

	return (
		<SplitPane
			divider
			left={editor}
			right={
				<ScrollRegion style={{ maxHeight: 420 }}>
					<MarkdownArticle source={value} style={{ padding: "8px 12px" }} />
				</ScrollRegion>
			}
		/>
	);
}
