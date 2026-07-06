import type { PaginationWidgetProps } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import { Pagination } from "./Pagination";

export const paginationWidget = defineWidget<PaginationWidgetProps>({
	type: "Pagination",
	module: "cms.dsl",
	render: (props, _children, ctx) => {
		const onPageChangeAction = props.onPageChangeAction;
		return (
			<Pagination
				className={props.className}
				page={props.page}
				pageCount={props.pageCount}
				pageSize={props.pageSize}
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
