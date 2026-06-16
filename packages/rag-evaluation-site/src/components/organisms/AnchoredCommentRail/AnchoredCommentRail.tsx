import type { HTMLAttributes } from "react";
import type { AnchoredComment } from "../../../context";
import { Caption } from "../../foundation";
import { Panel, Stack } from "../../layout";
import { AnchoredCommentCard } from "../../molecules";

export interface AnchoredCommentRailProps extends Omit<HTMLAttributes<HTMLDivElement>, "title"> {
	title?: string;
	comments: AnchoredComment[];
	selectedCommentId?: string;
	onCommentSelect?: (commentId: string) => void;
}

export function AnchoredCommentRail({
	title = "Comments",
	comments,
	selectedCommentId,
	onCommentSelect,
	...rest
}: AnchoredCommentRailProps) {
	return (
		<Panel title={title} fill data-rag-organism="AnchoredCommentRail" {...rest}>
			<Stack gap="sm">
				<Caption>
					{comments.length} anchored comment{comments.length === 1 ? "" : "s"}
				</Caption>
				{comments.map((comment, index) => (
					<button
						key={comment.id}
						type="button"
						style={{
							appearance: "none",
							border: 0,
							background: "transparent",
							padding: 0,
							textAlign: "left",
							cursor: "pointer",
						}}
						onClick={() => onCommentSelect?.(comment.id)}
					>
						<AnchoredCommentCard
							comment={comment}
							index={index + 1}
							selected={comment.id === selectedCommentId}
						/>
					</button>
				))}
			</Stack>
		</Panel>
	);
}
