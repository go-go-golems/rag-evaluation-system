import { useEffect, useState } from "react";
import type { ActionSpec, ArticleListPanelWidgetProps } from "../../../widgets/ir";
import { defineWidget, type RenderContext } from "../../../widgets/registry";
import type { ArticleStatusFilter } from "./ArticleListPanel";
import { ArticleListPanel } from "./ArticleListPanel";

export const articleListPanelWidget = defineWidget<ArticleListPanelWidgetProps>({
	type: "ArticleListPanel",
	module: "cms.dsl",
	render: (props, _children, ctx) => <ArticleListPanelWidgetHost props={props} ctx={ctx} />,
});

function ArticleListPanelWidgetHost({
	props,
	ctx,
}: {
	props: ArticleListPanelWidgetProps;
	ctx: RenderContext;
}) {
	const [query, setQuery] = useState(props.query ?? "");
	useEffect(() => setQuery(props.query ?? ""), [props.query]);
	const onArticleSelectAction = props.onArticleSelectAction;
	const onCreateAction = props.onCreateAction;
	const onRowActionAction = props.onRowActionAction;
	const onStatusFilterChangeAction = props.onStatusFilterChangeAction;
	const onQuerySubmitAction = props.onQuerySubmitAction;
	const onPageChangeAction = props.onPageChangeAction;
	const dispatch = (action: ActionSpec, extra: Record<string, unknown>) =>
		ctx.dispatchAction(action, { componentType: "ArticleListPanel", ...extra });

	return (
		<ArticleListPanel
			className={props.className}
			articles={props.articles}
			selectedArticleId={props.selectedArticleId}
			onArticleSelect={
				onArticleSelectAction
					? (articleId) => dispatch(onArticleSelectAction, { articleId, value: articleId })
					: undefined
			}
			onCreate={onCreateAction ? () => dispatch(onCreateAction, {}) : undefined}
			onRowAction={
				onRowActionAction
					? (articleId, rowAction) =>
							dispatch(onRowActionAction, { articleId, rowAction, value: articleId })
					: undefined
			}
			statusFilter={props.statusFilter}
			onStatusFilterChange={
				onStatusFilterChangeAction
					? (status: ArticleStatusFilter) =>
							dispatch(onStatusFilterChangeAction, { status, value: status })
					: undefined
			}
			query={query}
			onQueryChange={onQuerySubmitAction ? setQuery : undefined}
			onQuerySubmit={
				onQuerySubmitAction
					? (submittedQuery) =>
							dispatch(onQuerySubmitAction, { query: submittedQuery, value: submittedQuery })
					: undefined
			}
			page={props.page}
			pageCount={props.pageCount}
			onPageChange={
				onPageChangeAction
					? (page) => dispatch(onPageChangeAction, { page, value: page })
					: undefined
			}
			emptyMessage={props.emptyMessage != null ? ctx.renderValue(props.emptyMessage) : undefined}
			title={props.title != null ? ctx.renderValue(props.title) : undefined}
			maxVisibleTags={props.maxVisibleTags}
		/>
	);
}
