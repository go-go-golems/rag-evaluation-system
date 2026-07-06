import type { RichArticleWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { RichArticle } from "./RichArticle";

export const richArticleWidget = defineWidget<RichArticleWidgetProps>({
	type: "RichArticle",
	module: "course.dsl",
	render: (props) => (
		<RichArticle className={props.className} blocks={props.blocks} styleSet={props.styleSet} />
	),
});
