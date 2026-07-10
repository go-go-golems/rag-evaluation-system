import { type CSSProperties, type DragEvent, type ReactNode, useState } from "react";
import styles from "./BoardEngine.module.css";

export interface BoardColumnSpec {
	id: string;
	header: ReactNode;
	footer?: ReactNode;
	/** Resolved accent (e.g. from a ContextStyleSet) applied to the column head. */
	accentStyle?: CSSProperties;
}

/**
 * The stable payload every card receives — the seam that keeps the board
 * domain-blind. The engine owns columns, drag, drop targets, and selection; the
 * card owns everything visual and semantic. Any component honoring this shape is
 * a valid card (a DealCard, a ContactCard, ...).
 */
export interface BoardCardPayload<Card> {
	card: Card;
	columnId: string;
	selected: boolean;
	dragging: boolean;
	onSelect: () => void;
}

export interface BoardMove {
	cardId: string;
	from: string;
	to: string;
	/** the card the moved card was dropped before, if any. */
	beforeId?: string;
}

export interface BoardEngineProps<Card> {
	columns: BoardColumnSpec[];
	cards: Card[];
	/** which column a card is in. */
	columnOf: (card: Card) => string;
	getCardId: (card: Card) => string;
	renderCard: (payload: BoardCardPayload<Card>) => ReactNode;
	selectedCardId?: string;
	onMove?: (move: BoardMove) => void;
	onCardSelect?: (cardId: string) => void;
	ariaLabel?: string;
}

/**
 * Generic kanban engine: columns of cards where a card can be dragged from one
 * column to another. Blind to deals, stages, or any domain — the CRM equivalent
 * of what MatrixGrid is for scheduling. Configure it with a preset
 * (pipelineBoard) to get a domain board.
 */
export function BoardEngine<Card>({
	columns,
	cards,
	columnOf,
	getCardId,
	renderCard,
	selectedCardId,
	onMove,
	onCardSelect,
	ariaLabel,
}: BoardEngineProps<Card>) {
	const [draggingId, setDraggingId] = useState<string | null>(null);
	const [dropTarget, setDropTarget] = useState<string | null>(null);

	const from = draggingId ? cards.find((c) => getCardId(c) === draggingId) : undefined;
	const fromColumn = from ? columnOf(from) : undefined;

	const endDrag = () => {
		setDraggingId(null);
		setDropTarget(null);
	};

	const drop = (toColumn: string, beforeId?: string) => {
		if (draggingId && fromColumn && onMove) {
			// no-op if dropped exactly where it already sits
			if (
				!(fromColumn === toColumn && beforeId === undefined && from && columnOf(from) === toColumn)
			) {
				onMove({ cardId: draggingId, from: fromColumn, to: toColumn, beforeId });
			}
		}
		endDrag();
	};

	const onColumnDragOver = (columnId: string) => (e: DragEvent) => {
		if (draggingId) {
			e.preventDefault();
			setDropTarget(columnId);
		}
	};

	return (
		<div
			className={styles.root}
			data-rag-molecule="BoardEngine"
			role="group"
			aria-label={ariaLabel}
		>
			{columns.map((column) => {
				const columnCards = cards.filter((c) => columnOf(c) === column.id);
				return (
					<section
						key={column.id}
						className={styles.column}
						data-drop-active={dropTarget === column.id || undefined}
						onDragOver={onColumnDragOver(column.id)}
						onDrop={() => drop(column.id)}
					>
						<header className={styles.columnHead} style={column.accentStyle}>
							{column.header}
						</header>
						<ul className={styles.cards}>
							{columnCards.map((card) => {
								const cardId = getCardId(card);
								const selected = selectedCardId === cardId;
								const dragging = draggingId === cardId;
								return (
									// biome-ignore lint/a11y/useKeyWithClickEvents: selection also reachable via card content
									<li
										key={cardId}
										className={styles.cardSlot}
										draggable
										onDragStart={() => setDraggingId(cardId)}
										onDragEnd={endDrag}
										onDragOver={(e) => {
											if (draggingId && draggingId !== cardId) {
												e.preventDefault();
												e.stopPropagation();
												setDropTarget(column.id);
											}
										}}
										onDrop={(e) => {
											e.stopPropagation();
											drop(column.id, cardId);
										}}
										onClick={() => onCardSelect?.(cardId)}
									>
										{renderCard({
											card,
											columnId: column.id,
											selected,
											dragging,
											onSelect: () => onCardSelect?.(cardId),
										})}
									</li>
								);
							})}
							{columnCards.length === 0 ? <li className={styles.empty}>Drop here</li> : null}
						</ul>
						{column.footer != null ? (
							<footer className={styles.columnFoot}>{column.footer}</footer>
						) : null}
					</section>
				);
			})}
		</div>
	);
}
