import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { contextVisualStyleToCssVars } from "../../../context";
import { sampleDeals, sampleSalesPipeline, stageStyleSet } from "../../../crm";
import type { Deal } from "../../../crm/types";
import { DealCard } from "../DealCard";
import { BoardEngine, type BoardColumnSpec } from "./BoardEngine";

const stages = [...sampleSalesPipeline.stages].sort((a, b) => a.order - b.order);

function columnsFor(deals: Deal[]): BoardColumnSpec[] {
	return stages.map((stage) => {
		const inStage = deals.filter((d) => d.stageId === stage.id);
		const total = inStage.reduce((sum, d) => sum + (d.amount ?? 0), 0);
		return {
			id: stage.id,
			header: `${stage.name} · $${Math.round(total / 1000)}k · ${inStage.length}`,
			accentStyle: contextVisualStyleToCssVars(stageStyleSet.styles[stage.colorKey]!),
		};
	});
}

function ownerName(id?: string): string {
	return id ? id.replace("u-", "🧑 ") : "unassigned";
}

const meta = {
	title: "Component Library/Molecules/BoardEngine",
	component: BoardEngine<Deal>,
	args: {
		ariaLabel: "Sales pipeline",
		columns: columnsFor(sampleDeals),
		cards: sampleDeals,
		columnOf: (d: Deal) => d.stageId,
		getCardId: (d: Deal) => d.id,
		renderCard: ({ card, selected, dragging }) => (
			<DealCard
				title={card.title}
				subtitle={card.amount != null ? `$${card.amount.toLocaleString("en-US")}` : undefined}
				meta={ownerName(card.ownerId)}
				status={card.status}
				accentStyle={contextVisualStyleToCssVars(stageStyleSet.styles.qualified!)}
				selected={selected}
				dragging={dragging}
			/>
		),
	},
	decorators: [
		(Story) => (
			<div style={{ height: 460 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof BoardEngine<Deal>>;

export default meta;
type Story = StoryObj<typeof meta>;

/** Live board: drag cards between columns, click to select. */
export const Default: Story = {
	render: () => {
		const [deals, setDeals] = useState<Deal[]>(sampleDeals);
		const [selected, setSelected] = useState<string | undefined>();
		return (
			<BoardEngine<Deal>
				ariaLabel="Sales pipeline"
				columns={columnsFor(deals)}
				cards={deals}
				columnOf={(d) => d.stageId}
				getCardId={(d) => d.id}
				selectedCardId={selected}
				onCardSelect={setSelected}
				onMove={({ cardId, to }) =>
					setDeals((ds) => ds.map((d) => (d.id === cardId ? { ...d, stageId: to } : d)))
				}
				renderCard={({ card, selected: sel, dragging }) => (
					<DealCard
						title={card.title}
						subtitle={card.amount != null ? `$${card.amount.toLocaleString("en-US")}` : undefined}
						meta={ownerName(card.ownerId)}
						status={card.status}
						accentStyle={contextVisualStyleToCssVars(
							stageStyleSet.styles[
								card.status === "won" ? "won" : card.status === "lost" ? "lost" : "qualified"
							]!,
						)}
						selected={sel}
						dragging={dragging}
					/>
				)}
			/>
		);
	},
};

/** Empty board — every column shows its drop placeholder. */
export const Empty: Story = {
	render: () => (
		<BoardEngine<Deal>
			ariaLabel="Empty pipeline"
			columns={columnsFor([])}
			cards={[]}
			columnOf={(d) => d.stageId}
			getCardId={(d) => d.id}
			renderCard={({ card }) => <DealCard title={card.title} />}
		/>
	),
};

/** Dense column — many cards, per-column scroll kicks in. */
export const Dense: Story = {
	render: () => {
		const many: Deal[] = Array.from({ length: 12 }, (_, i) => ({
			...sampleDeals[0]!,
			id: `dense-${i}`,
			title: `Deal ${i + 1}`,
			amount: (i + 1) * 1000,
			stageId: "st-qualified",
		}));
		return (
			<BoardEngine<Deal>
				ariaLabel="Dense pipeline"
				columns={columnsFor(many)}
				cards={many}
				columnOf={(d) => d.stageId}
				getCardId={(d) => d.id}
				renderCard={({ card, dragging, selected }) => (
					<DealCard
						title={card.title}
						subtitle={`$${card.amount?.toLocaleString("en-US")}`}
						accentStyle={contextVisualStyleToCssVars(stageStyleSet.styles.qualified!)}
						dragging={dragging}
						selected={selected}
					/>
				)}
			/>
		);
	},
};

/** Selected state — one card highlighted, static. */
export const Selected: Story = {
	render: () => (
		<BoardEngine<Deal>
			ariaLabel="Pipeline with selection"
			columns={columnsFor(sampleDeals)}
			cards={sampleDeals}
			columnOf={(d) => d.stageId}
			getCardId={(d) => d.id}
			selectedCardId="d-globex"
			renderCard={({ card, selected }) => (
				<DealCard
					title={card.title}
					subtitle={card.amount != null ? `$${card.amount.toLocaleString("en-US")}` : undefined}
					meta={ownerName(card.ownerId)}
					status={card.status}
					accentStyle={contextVisualStyleToCssVars(stageStyleSet.styles.proposal!)}
					selected={selected}
				/>
			)}
		/>
	),
};
