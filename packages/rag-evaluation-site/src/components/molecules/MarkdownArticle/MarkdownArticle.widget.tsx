import type { MarkdownArticleWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { MarkdownArticle } from "./MarkdownArticle";

export const markdownArticleWidget = defineWidget<MarkdownArticleWidgetProps>({
	type: "MarkdownArticle",
	module: "course.dsl",
	render: (props) => <MarkdownArticle className={props.className} source={props.source} />,
});
