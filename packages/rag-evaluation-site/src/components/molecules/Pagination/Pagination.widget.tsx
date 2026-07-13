import type { PaginationWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Pagination } from "./Pagination";

export const paginationWidget = defineWidget<PaginationWidgetProps>({
	type: "Pagination",
	module: "data.dsl",
	render: (props, _children, ctx) => {
		const onPageChangeAction = props.onPageChangeAction;
		const onPageSizeChangeAction = props.onPageSizeChangeAction;
		return (
			<Pagination
				className={props.className}
				page={props.page}
				pageCount={props.pageCount}
				pageSize={props.pageSize}
				pageSizes={props.pageSizes}
				onPageSizeChange={(pageSize) => {
					if (onPageSizeChangeAction)
						ctx.dispatchAction(onPageSizeChangeAction, {
							page: 1,
							pageSize,
							value: pageSize,
							componentType: "Pagination",
						});
				}}
				totalItems={props.totalItems}
				onPageChange={(page) => {
					if (onPageChangeAction) {
						ctx.dispatchAction(onPageChangeAction, {
							page,
							value: page,
							componentType: "Pagination",
						});
					}
				}}
			/>
		);
	},
});
