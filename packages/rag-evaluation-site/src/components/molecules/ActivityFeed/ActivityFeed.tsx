import type { CSSProperties, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import styles from "./ActivityFeed.module.css";

export interface ActivityFeedItem {
	id: string;
	/** e.g. "note" | "email" | "call" | "stage_change" — drives glyph + color. */
	kind: string;
	title: ReactNode;
	body?: ReactNode;
	atISO: string;
	actor?: { name: string; avatarUrl?: string };
}

export interface ActivityFeedProps {
	activities: ActivityFeedItem[];
	/** kind -> glyph for the spine icon. Falls back to the kind's first letter. */
	glyphs?: Record<string, ReactNode>;
	/** kind -> color, via the shared ContextStyleSet contract. */
	styleSet?: ContextStyleSet;
	/** Group items under a per-day marker (default true). */
	groupByDay?: boolean;
	onOpen?: (id: string) => void;
	onLoadMore?: () => void;
	loadMoreLabel?: ReactNode;
	emptyLabel?: ReactNode;
}

const MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

function dayLabel(iso: string): string {
	const [y, m, d] = iso.slice(0, 10).split("-");
	return `${MONTHS[Number(m) - 1] ?? "?"} ${Number(d)}, ${y}`;
}

function timeLabel(iso: string): string {
	return iso.length >= 16 ? iso.slice(11, 16) : "";
}

function accentVars(kind: string, styleSet?: ContextStyleSet): CSSProperties | undefined {
	const visual = styleSet?.styles[kind];
	return visual ? contextVisualStyleToCssVars(visual) : undefined;
}

function ActivityRow({
	item,
	glyphs,
	styleSet,
	isLast,
	onOpen,
}: {
	item: ActivityFeedItem;
	glyphs?: Record<string, ReactNode>;
	styleSet?: ContextStyleSet;
	isLast: boolean;
	onOpen?: (id: string) => void;
}) {
	const glyph = glyphs?.[item.kind] ?? item.kind.charAt(0).toUpperCase();
	return (
		<li className={styles.row} data-kind={item.kind} data-last={isLast || undefined}>
			<span className={styles.spine} aria-hidden="true">
				<span className={styles.dot} style={accentVars(item.kind, styleSet)}>
					{glyph}
				</span>
			</span>
			<div className={styles.content}>
				<div className={styles.head}>
					{onOpen ? (
						<button type="button" className={styles.title} onClick={() => onOpen(item.id)}>
							{item.title}
						</button>
					) : (
						<span className={styles.title}>{item.title}</span>
					)}
					<span className={styles.metaRight}>
						{timeLabel(item.atISO) ? (
							<span className={styles.time}>{timeLabel(item.atISO)}</span>
						) : null}
						{item.actor ? <span className={styles.actor}>🧑 {item.actor.name}</span> : null}
					</span>
				</div>
				{item.body != null ? <div className={styles.body}>{item.body}</div> : null}
			</div>
		</li>
	);
}

/**
 * A record's history as one reverse-chronological stream of activities, with a
 * connective spine and per-day grouping. The engine owns the spine, the grouping,
 * and "load more"; each item's glyph/color comes from its `kind` via the shared
 * palette. Domain-blind — the CRM analogue of the transcript message list.
 */
export function ActivityFeed({
	activities,
	glyphs,
	styleSet,
	groupByDay = true,
	onOpen,
	onLoadMore,
	loadMoreLabel = "Load earlier",
	emptyLabel = "No activity yet",
}: ActivityFeedProps) {
	if (activities.length === 0) {
		return (
			<div className={styles.root} data-rag-molecule="ActivityFeed">
				<p className={styles.empty}>{emptyLabel}</p>
			</div>
		);
	}

	const ordered = [...activities].sort((a, b) => (a.atISO < b.atISO ? 1 : -1));

	// Group consecutively by day, preserving order.
	const groups: Array<{ day: string; items: ActivityFeedItem[] }> = [];
	for (const item of ordered) {
		const day = dayLabel(item.atISO);
		const last = groups[groups.length - 1];
		if (groupByDay && last && last.day === day) last.items.push(item);
		else groups.push({ day, items: [item] });
	}

	return (
		<div className={styles.root} data-rag-molecule="ActivityFeed">
			{groups.map((group, gi) => (
				<section key={`${group.day}-${gi}`} className={styles.group}>
					{groupByDay ? <h5 className={styles.dayLabel}>{group.day}</h5> : null}
					<ul className={styles.list}>
						{group.items.map((item, ii) => (
							<ActivityRow
								key={item.id}
								item={item}
								glyphs={glyphs}
								styleSet={styleSet}
								isLast={gi === groups.length - 1 && ii === group.items.length - 1}
								onOpen={onOpen}
							/>
						))}
					</ul>
				</section>
			))}
			{onLoadMore ? (
				<button type="button" className={styles.loadMore} onClick={onLoadMore}>
					{loadMoreLabel}
				</button>
			) : null}
		</div>
	);
}
