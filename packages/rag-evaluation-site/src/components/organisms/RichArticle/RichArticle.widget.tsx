import { RichArticle } from "./RichArticle";
import { defineWidget } from "../../../widgets/registry";
import type { RichArticleWidgetProps } from "../../../widgets/ir";

export const richArticleWidget = defineWidget<RichArticleWidgetProps>({
	type: "RichArticle",
	module: "course.dsl",
	render: (props) => (
		<RichArticle className={props.className} blocks={props.blocks} styleSet={props.styleSet} />
	),
});
