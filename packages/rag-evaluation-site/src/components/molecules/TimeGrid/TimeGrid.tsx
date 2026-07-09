import type { CSSProperties, HTMLAttributes, ReactNode } from "react";
import { type ContextStyleSet, contextVisualStyleToCssVars } from "../../../context";
import { packTimeGridColumn, timeParts } from "./TimeGrid.logic";
import styles from "./TimeGrid.module.css";

export interface TimeGridBlock {
	id: string;
	/** `YYYY-MM-DD` selecting which day column this block lives in. */
	dayISO: string;
	/** ISO datetime; only the wall-clock `HH:MM` is used for vertical position. */
	startISO: string;
	endISO: string;
	styleKey: string;
	label: ReactNode;
	meta?: Record<string, unknown>;
}

export interface TimeGridColumnSpec {
	dayISO: string;
	header?: ReactNode;
}

/**
 * The stable payload every block receives. TimeGrid owns the time geometry
 * (position + lane packing); the block renderer owns how a block looks. Mirrors
 * the MatrixGrid / MonthGrid cell contracts.
 */
export interface TimeGridBlockPayload {
	block: TimeGridBlock;
	selected: boolean;
	onSelect: () => void;
}

export interface TimeGridProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	/** Columns, in order. Accepts `YYYY-MM-DD` strings or `{ dayISO, header }`. */
	days: Array<string | TimeGridColumnSpec>;
	blocks: TimeGridBlock[];
	styleSet: ContextStyleSet;
	hourStart?: number;
	hourEnd?: number;
	/** Pixels per hour. */
	hourHeight?: number;
	/** ISO datetime for the "now" indicator; omit to hide it. */
	nowISO?: string;
	selectedBlockId?: string;
	/** Mode A — custom block renderer. Omit for the default block. */
	renderBlock?: (payload: TimeGridBlockPayload) => ReactNode;
	onBlockSelect?: (blockId: string) => void;
	/** Click on empty space → create. Fires `{ dayISO, hour }`. */
	onSlotCreate?: (slot: { dayISO: string; hour: number }) => void;
}

function pad2(n: number): string {
	return String(n).padStart(2, "0");
}

function normalizeColumn(day: string | TimeGridColumnSpec): TimeGridColumnSpec {
	return typeof day === "string" ? { dayISO: day } : day;
}

export function TimeGrid({
	days,
	blocks,
	styleSet,
	hourStart = 8,
	hourEnd = 20,
	hourHeight = 40,
	nowISO,
	selectedBlockId,
	renderBlock,
	onBlockSelect,
	onSlotCreate,
	className,
	style,
	...rest
}: TimeGridProps) {
	const columns = days.map(normalizeColumn);
	const rangeStart = hourStart * 60;
	const rangeMinutes = Math.max(60, (hourEnd - hourStart) * 60);
	const bodyHeight = (rangeMinutes / 60) * hourHeight;
	const hourLabels = Array.from({ length: hourEnd - hourStart + 1 }, (_, i) => hourStart + i);

	const now = nowISO ? timeParts(nowISO) : null;

	const rootStyle = {
		...style,
		"--rag-hour-height": `${hourHeight}px`,
	} as CSSProperties;

	return (
		<div
			className={[styles.root, className ?? ""].filter(Boolean).join(" ")}
			data-rag-molecule="TimeGrid"
			style={rootStyle}
			{...rest}
		>
			<div className={styles.corner} />
			<div className={styles.headers}>
				{columns.map((col) => (
					<div key={col.dayISO} className={styles.colHeader}>
						{col.header ?? col.dayISO}
					</div>
				))}
			</div>

			<div className={styles.gutter} style={{ height: bodyHeight }}>
				{hourLabels.slice(0, -1).map((hour, i) => (
					<div key={hour} className={styles.hourLabel} style={{ top: i * hourHeight }}>
						<span>{pad2(hour)}:00</span>
					</div>
				))}
			</div>

			<div className={styles.columns}>
				{columns.map((col) => {
					const packed = packTimeGridColumn(
						blocks.filter((b) => b.dayISO === col.dayISO),
						rangeStart,
						rangeMinutes,
					);
					const showNow = now != null && now.date === col.dayISO;
					return (
						<div key={col.dayISO} className={styles.column} style={{ height: bodyHeight }}>
							{hourLabels.slice(0, -1).map((hour) => (
								<button
									key={hour}
									type="button"
									className={styles.hourSlot}
									style={{ height: hourHeight }}
									aria-label={`Create ${col.dayISO} ${pad2(hour)}:00`}
									tabIndex={onSlotCreate ? 0 : -1}
									disabled={!onSlotCreate}
									onClick={() => onSlotCreate?.({ dayISO: col.dayISO, hour })}
								/>
							))}
							{packed.map(({ block, topPct, heightPct, lane, lanes }) => {
								const selected = selectedBlockId === block.id;
								const visualStyle = styleSet.styles[block.styleKey] ?? styleSet.fallbackStyle;
								const blockStyle: CSSProperties = {
									top: `${topPct}%`,
									height: `${heightPct}%`,
									left: `${(lane / lanes) * 100}%`,
									width: `${(1 / lanes) * 100}%`,
									...(visualStyle ? contextVisualStyleToCssVars(visualStyle) : {}),
								};
								if (renderBlock) {
									return (
										<div key={block.id} className={styles.blockSlot} style={blockStyle}>
											{renderBlock({
												block,
												selected,
												onSelect: () => onBlockSelect?.(block.id),
											})}
										</div>
									);
								}
								return (
									<button
										key={block.id}
										type="button"
										className={styles.block}
										style={blockStyle}
										data-selected={selected || undefined}
										onClick={() => onBlockSelect?.(block.id)}
									>
										<span className={styles.blockLabel}>{block.label}</span>
									</button>
								);
							})}
							{showNow && now ? (
								<div
									className={styles.nowLine}
									style={{ top: `${((now.minutes - rangeStart) / rangeMinutes) * 100}%` }}
									aria-hidden="true"
								/>
							) : null}
						</div>
					);
				})}
			</div>
		</div>
	);
}
