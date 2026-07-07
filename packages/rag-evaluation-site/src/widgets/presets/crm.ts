import {
	ACTIVITY_GLYPHS,
	type Activity,
	activityStyleSet,
	type Company,
	type Contact,
	type CrmUser,
	type Deal,
	type FieldDef,
	type Pipeline,
	stageStyleSet,
	type StageSummary,
	tagStyleSet,
} from "../../crm";
import {
	type ActivityFeedWidgetProps,
	type BoardColumnWidgetSpec,
	type BoardEngineWidgetProps,
	component,
	type FieldRendererWidgetProps,
	type FieldRefSpec,
	type FieldSpec,
	type JsonObject,
	type RecordFieldListSectionSpec,
	type RecordFieldListWidgetProps,
	type SegmentedBarSegmentSpec,
	type SegmentedBarWidgetProps,
	type StatTileWidgetProps,
	text,
	type WidgetNode,
} from "../ir";
import type { Stage, Task } from "../../crm/types";

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

// ── Record-page presets (IR compositions of registered widgets) ──────────────

/** `crm.dsl` preset: a record's activity timeline. */
export function activityFeed(activities: Activity[]): WidgetNode {
	const props: ActivityFeedWidgetProps = {
		activities: activities.map((a) => ({
			id: a.id,
			kind: a.kind,
			title: text(a.title),
			body: a.body != null ? text(a.body) : undefined,
			atISO: a.atISO,
			actor: { id: a.actor.id, name: a.actor.name, avatarUrl: a.actor.avatarUrl },
		})),
		styleSet: activityStyleSet,
		glyphs: ACTIVITY_GLYPHS,
		onOpenAction: { kind: "event", event: "activity.open" },
	};
	return component("ActivityFeed", props);
}

/** `crm.dsl` preset: a record's field list from its FieldDefs + values. */
export function recordFieldList(
	values: Record<string, unknown>,
	defs: FieldDef[],
	options: { mode?: "read" | "edit"; refs?: Record<string, FieldRefSpec> } = {},
): WidgetNode {
	const props: RecordFieldListWidgetProps = {
		values: values as JsonObject,
		sections: fieldSections(defs),
		mode: options.mode ?? "read",
		refs: options.refs,
		onFieldChangeAction: { kind: "server", name: "field.update" },
	};
	return component("RecordFieldList", props);
}

export interface RecordPageOptions {
	activities?: Activity[];
	users?: CrmUser[];
	companies?: Company[];
	mode?: "read" | "edit";
	related?: WidgetNode;
}

/**
 * `crm.dsl` preset: the record page as an IR composition of already-registered
 * widgets (Panel + SplitPane + RecordFieldList + ActivityFeed) — no bespoke
 * RecordShell node, mirroring how the scheduling presets compose Stack +
 * SegmentedBar. The React `RecordShell` organism is the hand-authored twin.
 */
export function contactRecord(
	contact: Contact,
	defs: FieldDef[],
	options: RecordPageOptions = {},
): WidgetNode {
	const refs = buildRefs(options.users ?? [], options.companies ?? []);
	const subtitle = [contact.title, ...(contact.tags ?? []).map((t) => `🏷 ${t}`)]
		.filter(Boolean)
		.join(" · ");

	const rightChildren: WidgetNode[] = [
		component("Panel", { title: "Activity", density: "condensed" }, [
			activityFeed(options.activities ?? []),
		]),
	];
	if (options.related) rightChildren.push(options.related);

	const splitPane = component("SplitPane", {
		ratio: "leftNarrow",
		gutter: "lg",
		left: component("Panel", { title: "Details", density: "condensed" }, [
			recordFieldList(contact.fields, defs, { mode: options.mode, refs }),
		]),
		right: component("Stack", { gap: "md" }, rightChildren),
	});

	return component("Panel", { title: contact.name, density: "condensed" }, [
		component("Stack", { gap: "sm" }, [
			...(subtitle ? [component("Caption", {}, [text(subtitle)])] : []),
			splitPane,
		]),
	]) as WidgetNode;
}

// ── Dashboard + tasks-inbox presets (compositions of registered widgets) ─────

/** `crm.dsl` preset: a single dashboard metric tile. */
export function statTile(
	label: string,
	value: string,
	options: {
		delta?: number;
		deltaLabel?: string;
		trend?: "up" | "down" | "flat";
		progress?: number;
		tone?: "accent" | "success" | "danger";
	} = {},
): WidgetNode {
	const props: StatTileWidgetProps = {
		label: text(label),
		value: text(value),
		delta: options.delta,
		deltaLabel: options.deltaLabel != null ? text(options.deltaLabel) : undefined,
		trend: options.trend,
		progress: options.progress,
		tone: options.tone,
	};
	return component("StatTile", props);
}

/**
 * `crm.dsl` preset: stage counts as a funnel — a SegmentedBar colored by stage.
 * No new engine, per the design.
 */
export function pipelineFunnel(pipeline: Pipeline, summaries: StageSummary[]): WidgetNode {
	const byStage = new Map(summaries.map((s) => [s.stageId, s]));
	const stages = [...pipeline.stages].sort((a: Stage, b: Stage) => a.order - b.order);
	const segments: SegmentedBarSegmentSpec[] = stages.map((stage) => ({
		value: byStage.get(stage.id)?.count ?? 0,
		styleKey: stage.colorKey,
		label: text(stage.name),
	}));
	const props: SegmentedBarWidgetProps = {
		segments,
		styleSet: stageStyleSet,
		showCounts: true,
		size: "lg",
	};
	return component("SegmentedBar", props);
}

/**
 * `crm.dsl` preset: the CRM dashboard — a TileGrid of StatTiles, a pipeline
 * funnel, and a recent-activity feed. Pure composition of existing engines.
 */
export function crmDashboard(
	pipeline: Pipeline,
	summaries: StageSummary[],
	options: { activities?: Activity[]; tiles?: WidgetNode[] } = {},
): WidgetNode {
	const totalOpen = summaries.reduce((sum, s) => sum + s.amountTotal, 0);
	const totalDeals = summaries.reduce((sum, s) => sum + s.count, 0);
	const tiles = options.tiles ?? [
		statTile("Open pipeline", formatAmount(totalOpen), { delta: 12, progress: 0.62 }),
		statTile("Open deals", String(totalDeals), { delta: 4, tone: "success", progress: 0.5 }),
		statTile("Win rate", "38%", { delta: -3, tone: "danger", progress: 0.38 }),
	];
	return component("Stack", { gap: "md" }, [
		component("TileGrid", { minTileWidth: 160 }, tiles),
		component("Panel", { title: "Pipeline funnel", density: "condensed" }, [
			pipelineFunnel(pipeline, summaries),
		]),
		component("Panel", { title: "Recent activity", density: "condensed" }, [
			activityFeed(options.activities ?? []),
		]),
	]) as WidgetNode;
}

/**
 * `crm.dsl` preset: the tasks inbox — tasks with an inline FieldRenderer-driven
 * "mark done", exercising inline field editing outside a record page. Composed
 * from Stack + FieldRenderer + Text (no ItemList engine yet).
 */
export function tasksInbox(tasks: Task[]): WidgetNode {
	const rows = tasks.map((task) => {
		const doneField: FieldRendererWidgetProps = {
			spec: { key: "done", type: "boolean" },
			value: task.status === "done",
			mode: "edit",
			onChangeAction: {
				kind: "server",
				name: "task.complete",
				payload: { taskId: task.id } as JsonObject,
			},
		};
		const due = task.dueISO ? task.dueISO.slice(5) : "—";
		return component("Inline", { gap: "sm", justify: "between" }, [
			component("Inline", { gap: "sm" }, [
				component("FieldRenderer", doneField),
				component("Text", { size: "compact" }, [text(task.title)]),
			]),
			component("Caption", {}, [text(`${task.priority ?? ""} · ${due}`)]),
		]);
	});
	return component("Panel", { title: "Tasks", density: "condensed" }, [
		component("Stack", { gap: "xs" }, rows),
	]) as WidgetNode;
}
