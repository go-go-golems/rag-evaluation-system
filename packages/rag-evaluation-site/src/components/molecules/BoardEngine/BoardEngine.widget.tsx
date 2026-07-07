import type { CSSProperties } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import { renderCell, rowKey } from "../../../widgets/cellRenderers";
import type { BoardEngineWidgetProps, JsonObject } from "../../../widgets/ir";
import { defineWidget } from "../../../widgets/registry";
import type { RenderContext } from "../../../widgets/registry";
import { DealCard, type DealCardStatus } from "../DealCard";
import { BoardEngine } from "./BoardEngine";

function accentFor(
	styleKey: string | undefined,
	styleSet: ContextStyleSet | undefined,
): CSSProperties | undefined {
	if (!styleKey) return undefined;
	const visual = styleSet?.styles[styleKey];
	return visual ? contextVisualStyleToCssVars(visual) : undefined;
}

export const boardEngineWidget = defineWidget<BoardEngineWidgetProps>({
	type: "BoardEngine",
	module: "data.dsl",
	render: (props, _children, ctx: RenderContext) => {
		const columns = props.columns.map((column) => ({
			id: column.id,
			header: ctx.renderValue(column.header),
			footer: column.footer != null ? ctx.renderValue(column.footer) : undefined,
			accentStyle: accentFor(column.accent, props.styleSet),
		}));
		const cardSpec = props.card;
		return (
			<BoardEngine<JsonObject>
				ariaLabel={props.ariaLabel}
				columns={columns}
				cards={props.cards}
				columnOf={(card) => String(card[props.columnField] ?? "")}
				getCardId={(card) =>
					props.getCardId ? rowKey(card, props.getCardId) : String(card.id ?? "")
				}
				selectedCardId={props.selectedCardId}
				renderCard={({ card, selected, dragging }) => {
					const accentKey = cardSpec.accentField
						? String(card[cardSpec.accentField] ?? "")
						: undefined;
					const status = String(card.status ?? "open") as DealCardStatus;
					return (
						<DealCard
							title={renderCell(cardSpec.title, card, ctx.renderNode)}
							subtitle={
								cardSpec.subtitle ? renderCell(cardSpec.subtitle, card, ctx.renderNode) : undefined
							}
							meta={cardSpec.meta ? renderCell(cardSpec.meta, card, ctx.renderNode) : undefined}
							status={status === "won" || status === "lost" ? status : "open"}
							accentStyle={accentFor(accentKey, props.styleSet)}
							selected={selected}
							dragging={dragging}
						/>
					);
				}}
				onMove={
					props.onMoveAction
						? (move) =>
								ctx.dispatchAction(props.onMoveAction!, {
									cardId: move.cardId,
									from: move.from,
									to: move.to,
									beforeId: move.beforeId,
									componentType: "BoardEngine",
								} as unknown as Record<string, unknown>)
						: undefined
				}
				onCardSelect={
					props.onCardSelectAction
						? (cardId) =>
								ctx.dispatchAction(props.onCardSelectAction!, {
									cardId,
									componentType: "BoardEngine",
								} as unknown as Record<string, unknown>)
						: undefined
				}
			/>
		);
	},
});
