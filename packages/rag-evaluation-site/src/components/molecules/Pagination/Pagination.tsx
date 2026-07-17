import type { HTMLAttributes } from "react";
import { Button } from "../../atoms";
import { Caption } from "../../foundation";
import styles from "./Pagination.module.css";

export interface PaginationProps extends HTMLAttributes<HTMLDivElement> {
	page: number;
	pageCount: number;
	onPageChange: (page: number) => void;
	pageSize?: number;
	pageSizes?: number[];
	onPageSizeChange?: (pageSize: number) => void;
	totalItems?: number;
}

export function Pagination({
	page,
	pageCount,
	onPageChange,
	pageSize,
	pageSizes = [],
	onPageSizeChange,
	totalItems,
	className,
	...rest
}: PaginationProps) {
	const clampedCount = Math.max(1, pageCount);
	const clampedPage = Math.min(Math.max(1, page), clampedCount);
	const showTotals = pageSize != null && totalItems != null && totalItems > 0;
	const first = (clampedPage - 1) * (pageSize ?? 0) + 1;
	const last = Math.min(clampedPage * (pageSize ?? 0), totalItems ?? 0);

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="Pagination"
			{...rest}
		>
			<Button
				size="compact"
				disabled={clampedPage <= 1}
				aria-label="Previous page"
				onClick={() => onPageChange(clampedPage - 1)}
			>
				‹ prev
			</Button>
			<Caption className={styles.counter}>
				page {clampedPage} / {clampedCount}
			</Caption>
			<Button
				size="compact"
				disabled={clampedPage >= clampedCount}
				aria-label="Next page"
				onClick={() => onPageChange(clampedPage + 1)}
			>
				next ›
			</Button>
			{showTotals && (
				<Caption className={styles.totals}>
					{first}–{last} of {totalItems}
				</Caption>
			)}
			{pageSizes.length > 0 && onPageSizeChange && (
				<label>
					<Caption>rows</Caption>{" "}
					<select
						aria-label="Rows per page"
						value={pageSize}
						onChange={(event) => onPageSizeChange(Number(event.currentTarget.value))}
					>
						{pageSizes.map((size) => (
							<option key={size} value={size}>
								{size}
							</option>
						))}
					</select>
				</label>
			)}
		</div>
	);
}
