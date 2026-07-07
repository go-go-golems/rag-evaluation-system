import {
	type Company,
	type CrmUser,
	type Deal,
	type FieldDef,
	type Pipeline,
	stageStyleSet,
	type StageSummary,
	tagStyleSet,
} from "../../crm";
import {
	type BoardColumnWidgetSpec,
	type BoardEngineWidgetProps,
	component,
	type FieldRefSpec,
	type FieldSpec,
	type JsonObject,
	type RecordFieldListSectionSpec,
	text,
	type WidgetNode,
} from "../ir";

// ── Field-system presets (FieldDef -> FieldSpec / sections) ──────────────────

/** A CRM field's colors come from the tag palette when it is a select/tag type. */
function styleSetFor(def: FieldDef) {
	return def.type === "select" || def.type === "multiselect" || def.type === "tags"
		? tagStyleSet
		: undefined;
}

/** Turn a workspace FieldDef into the serializable FieldSpec the renderer reads. */
export function toFieldSpec(def: FieldDef): FieldSpec {
	return {
		key: def.key,
		type: def.type,
		label: text(def.label),
		options: def.options,
		relatedObject: def.relatedObject,
		readOnly: def.readOnly,
		unit: def.unit,
		styleSet: styleSetFor(def),
	};
}

/** Group FieldDefs into RecordFieldList sections by their `group`, order preserved. */
export function fieldSections(defs: FieldDef[]): RecordFieldListSectionSpec[] {
	const order: string[] = [];
	const byGroup = new Map<string, FieldSpec[]>();
	for (const def of defs) {
		const group = def.group ?? "Details";
		if (!byGroup.has(group)) {
			byGroup.set(group, []);
			order.push(group);
		}
		byGroup.get(group)!.push(toFieldSpec(def));
	}
	return order.map((group) => ({ label: text(group), fields: byGroup.get(group)! }));
}

/** Build the id -> display map for relation/user read-mode chips. */
export function buildRefs(
	users: CrmUser[],
	companies: Company[] = [],
): Record<string, FieldRefSpec> {
	const refs: Record<string, FieldRefSpec> = {};
	for (const u of users) refs[u.id] = { label: u.name, avatarUrl: u.avatarUrl };
	for (const c of companies) refs[c.id] = { label: c.name };
	return refs;
}

function formatAmount(n: number): string {
	if (n >= 1000) return `$${Math.round(n / 1000)}k`;
	return `$${n}`;
}

export interface PipelineBoardOptions {
	summaries?: StageSummary[];
	selectedDealId?: string;
}

/**
 * `crm.dsl` preset: a pipeline kanban. Turns a Pipeline + its Deals into a
 * configured BoardEngine whose columns are stages and whose cards are deals.
 * The words "deal" and "stage" live here and nowhere in the engine — the whole
 * point of the layering. A different preset (contactsByStatusBoard) could reuse
 * the identical engine for a lead-status board.
 */
export function pipelineBoard(
	pipeline: Pipeline,
	deals: Deal[],
	options: PipelineBoardOptions = {},
): WidgetNode {
	const summaryByStage = new Map((options.summaries ?? []).map((s) => [s.stageId, s]));
	const stages = [...pipeline.stages].sort((a, b) => a.order - b.order);

	const columns: BoardColumnWidgetSpec[] = stages.map((stage) => {
		const summary = summaryByStage.get(stage.id);
		const header = summary
			? `${stage.name}  ·  ${formatAmount(summary.amountTotal)} · ${summary.count}`
			: stage.name;
		return { id: stage.id, header: text(header), accent: stage.colorKey };
	});

	const props: BoardEngineWidgetProps = {
		ariaLabel: pipeline.name,
		columns,
		cards: deals as unknown as JsonObject[],
		columnField: "stageId",
		getCardId: { field: "id" },
		styleSet: stageStyleSet,
		selectedCardId: options.selectedDealId,
		card: {
			title: { kind: "field", field: "title" },
			subtitle: { kind: "number", field: "amount", format: "integer", fallback: "—" },
			meta: { kind: "field", field: "ownerId", fallback: "unassigned" },
			accentField: "status",
		},
		onMoveAction: {
			kind: "server",
			name: "deal.move",
			payload: {
				dealId: { kind: "path", path: "cardId" },
				fromStage: { kind: "path", path: "from" },
				toStage: { kind: "path", path: "to" },
				beforeId: { kind: "path", path: "beforeId" },
			} as JsonObject,
		},
		onCardSelectAction: {
			kind: "event",
			event: "deal.open",
			detail: { dealId: { kind: "path", path: "cardId" } } as JsonObject,
		},
	};

	return component("BoardEngine", props);
}

/** `crm.dsl` preset: wrap the pipeline board in a Panel with a header action. */
export function pipelineBoardPanel(
	pipeline: Pipeline,
	deals: Deal[],
	options: PipelineBoardOptions = {},
): WidgetNode {
	return component("Panel", { title: `Pipeline · ${pipeline.name}`, density: "condensed" }, [
		pipelineBoard(pipeline, deals, options),
	]) as WidgetNode;
}
