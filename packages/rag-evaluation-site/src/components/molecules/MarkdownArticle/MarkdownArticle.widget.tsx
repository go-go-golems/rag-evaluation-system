import { MarkdownArticle } from './MarkdownArticle';
import { defineWidget } from '../../../widgets/registry';
import type { MarkdownArticleWidgetProps } from '../../../widgets/ir';

export const markdownArticleWidget = defineWidget<MarkdownArticleWidgetProps>({
  type: 'MarkdownArticle',
  module: 'course.dsl',
  render: (props) => <MarkdownArticle className={props.className} source={props.source} />,
});
