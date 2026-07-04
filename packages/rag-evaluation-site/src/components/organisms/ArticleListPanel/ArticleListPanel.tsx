import type { HTMLAttributes, ReactNode } from "react";
import { useState } from "react";
import type { CmsArticleSummary, CmsContentStatus } from "../../../cms/types";
import { Button, ContentStatusBadge, IconButton, SelectInput, Tag } from "../../atoms";
import { Caption, CodeText } from "../../foundation";
import { Panel, Stack } from "../../layout";
import {
	DataTable,
	type DataTableColumn,
	EmptyState,
	Pagination,
	SearchField,
} from "../../molecules";
import { ConfirmDialog } from "../ConfirmDialog";
import styles from "./ArticleListPanel.module.css";

export type ArticleRowAction = "edit" | "publish" | "archive" | "delete";
export type ArticleStatusFilter = CmsContentStatus | "all";

export interface ArticleListPanelProps
	extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect" | "title"> {
	articles: CmsArticleSummary[];
	selectedArticleId?: string;
	onArticleSelect?: (articleId: string) => void;
	onCreate?: () => void;
	onRowAction?: (articleId: string, action: ArticleRowAction) => void;
	statusFilter?: ArticleStatusFilter;
	onStatusFilterChange?: (status: ArticleStatusFilter) => void;
	query?: string;
	onQueryChange?: (query: string) => void;
	onQuerySubmit?: (query: string) => void;
	page?: number;
	pageCount?: number;
	onPageChange?: (page: number) => void;
	emptyMessage?: ReactNode;
	title?: ReactNode;
	maxVisibleTags?: number;
}

const STATUS_OPTIONS: ArticleStatusFilter[] = [
	"all",
	"draft",
	"published",
	"scheduled",
	"archived",
];

interface PendingConfirm {
	article: CmsArticleSummary;
	action: Extract<ArticleRowAction, "archive" | "delete">;
}

export function ArticleListPanel({
	articles,
	selectedArticleId,
	onArticleSelect,
	onCreate,
	onRowAction,
	statusFilter = "all",
	onStatusFilterChange,
	query,
	onQueryChange,
	onQuerySubmit,
	page,
	pageCount,
	onPageChange,
	emptyMessage = "No articles yet",
	title = "Articles",
	maxVisibleTags = 2,
	className,
	...rest
}: ArticleListPanelProps) {
	const [pending, setPending] = useState<PendingConfirm>();

	const columns: DataTableColumn<CmsArticleSummary>[] = [
		{
			id: "title",
			header: "Title",
			cell: (row) => (
				<span className={styles.titleCell}>
					<span className={styles.titleText}>{row.title}</span>
					<CodeText className={styles.slug}>{row.slug}</CodeText>
				</span>
			),
		},
		{
			id: "status",
			header: "Status",
			cell: (row) => <ContentStatusBadge status={row.status} icon={false} />,
		},
		{
			id: "tags",
			header: "Tags",
			cell: (row) => (
				<span className={styles.tagsCell}>
					{row.tags.slice(0, maxVisibleTags).map((tag) => (
						<Tag key={tag} label={tag} />
					))}
					{row.tags.length > maxVisibleTags && (
						<Caption>+{row.tags.length - maxVisibleTags}</Caption>
					)}
				</span>
			),
		},
		{
			id: "author",
			header: "Author",
			cell: (row) => row.author ?? "—",
		},
		{
			id: "updated",
			header: "Updated",
			cell: (row) => <CodeText>{row.updatedAt.slice(0, 10)}</CodeText>,
			align: "end",
		},
	];

	if (onRowAction) {
		columns.push({
			id: "actions",
			header: "",
			align: "end",
			cell: (row) => (
				<span className={styles.actionsCell}>
					<IconButton
						size="large"
						variant="boxed"
						label={`Edit ${row.title}`}
						onClick={() => onRowAction(row.id, "edit")}
					>
						✎
					</IconButton>
					{row.status !== "published" && (
						<IconButton
							size="large"
							variant="boxed"
							label={`Publish ${row.title}`}
							onClick={() => onRowAction(row.id, "publish")}
						>
							●
						</IconButton>
					)}
					{row.status !== "archived" && (
						<IconButton
							size="large"
							variant="boxed"
							label={`Archive ${row.title}`}
							onClick={() => setPending({ article: row, action: "archive" })}
						>
							▣
						</IconButton>
					)}
					<IconButton
						size="large"
						variant="boxed"
						label={`Delete ${row.title}`}
						onClick={() => setPending({ article: row, action: "delete" })}
					>
						×
					</IconButton>
				</span>
			),
		});
	}

	return (
		<Panel
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			title={title}
			actions={
				onCreate && (
					<Button size="compact" variant="primary" onClick={onCreate}>
						＋ New article
					</Button>
				)
			}
			data-rag-organism="ArticleListPanel"
			{...rest}
		>
			<Stack gap="sm">
				{(onQueryChange || onStatusFilterChange || onPageChange) && (
					<div className={styles.toolbar}>
						{onQueryChange && (
							<SearchField
								className={styles.search}
								value={query ?? ""}
								onValueChange={onQueryChange}
								onSubmit={onQuerySubmit}
								placeholder="search articles…"
							/>
						)}
						{onStatusFilterChange && (
							<SelectInput
								aria-label="Filter by status"
								value={statusFilter}
								onChange={(event) =>
									onStatusFilterChange(event.target.value as ArticleStatusFilter)
								}
							>
								{STATUS_OPTIONS.map((option) => (
									<option key={option} value={option}>
										{option}
									</option>
								))}
							</SelectInput>
						)}
						{onPageChange && page != null && pageCount != null && (
							<Pagination
								className={styles.pager}
								page={page}
								pageCount={pageCount}
								onPageChange={onPageChange}
							/>
						)}
					</div>
				)}
				{articles.length === 0 ? (
					<EmptyState glyph="¶" title={emptyMessage} />
				) : (
					<DataTable
						columns={columns}
						rows={articles}
						getRowKey={(row) => row.id}
						selectedKey={selectedArticleId ?? null}
						onRowSelect={onArticleSelect ? (row) => onArticleSelect(row.id) : undefined}
						emptyMessage={emptyMessage}
					/>
				)}
			</Stack>
			<ConfirmDialog
				open={pending != null}
				title={pending?.action === "delete" ? "Delete article" : "Archive article"}
				message={
					pending
						? `${pending.action === "delete" ? "Delete" : "Archive"} “${pending.article.title}”?`
						: ""
				}
				detail={pending?.action === "delete" ? "This cannot be undone." : undefined}
				confirmLabel={pending?.action === "delete" ? "Delete" : "Archive"}
				destructive={pending?.action === "delete"}
				onConfirm={() => {
					if (pending) onRowAction?.(pending.article.id, pending.action);
					setPending(undefined);
				}}
				onCancel={() => setPending(undefined)}
			/>
		</Panel>
	);
}
