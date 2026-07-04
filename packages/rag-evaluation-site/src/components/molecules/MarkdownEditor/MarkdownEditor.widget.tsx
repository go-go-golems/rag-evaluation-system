import { useState } from "react";
import type { MarkdownEditorWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { ScrollRegion, SplitPane } from "../../layout";
import { MarkdownArticle } from "../MarkdownArticle";
import { MarkdownEditor } from "./MarkdownEditor";

export const markdownEditorWidget = defineWidget<MarkdownEditorWidgetProps>({
	type: "MarkdownEditor",
	module: "cms.dsl",
	render: (props) => <MarkdownEditorWidgetHost props={props} />,
});

// The editor value is browser-local state: the toolbar and the live
// MarkdownArticle preview work without any server round-trip, while the
// named textarea inside MarkdownEditor carries the value in native form
// posts (formPanel({method:"post", formAction})).
function MarkdownEditorWidgetHost({ props }: { props: MarkdownEditorWidgetProps }) {
	const [value, setValue] = useState(props.defaultValue ?? "");

	const editor = (
		<MarkdownEditor
			className={props.className}
			name={props.name}
			value={value}
			onValueChange={setValue}
			minRows={props.minRows}
			maxLength={props.maxLength}
			disabled={props.disabled}
			textareaAriaLabel={props.textareaAriaLabel}
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
