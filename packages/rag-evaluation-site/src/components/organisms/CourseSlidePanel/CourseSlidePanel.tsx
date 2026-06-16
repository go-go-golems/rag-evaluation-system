import type { HTMLAttributes } from "react";
import type {
	ContextLegendItemSpec,
	ContextSlide,
	ContextStyleSet,
	ContextWindowSnapshot,
} from "../../../context";
import { contextDefaultStyleSet } from "../../../context";
import { Button } from "../../atoms";
import { Inline, SlideShell } from "../../layout";
import {
	ContextBudgetBar,
	ContextLegend,
	ContextStackDiagram,
	ContextStripDiagram,
	ContextTreemap,
	FigureBlock,
	KeyPointList,
} from "../../molecules";

export interface CourseSlidePanelProps extends Omit<HTMLAttributes<HTMLElement>, "title"> {
	slide: ContextSlide;
	snapshot: ContextWindowSnapshot;
	styleSet?: ContextStyleSet;
	index?: number;
	total?: number;
	visualSide?: "left" | "right";
	mode?: "default" | "present";
	onPrevious?: () => void;
	onNext?: () => void;
	onPresent?: () => void;
	onFullscreen?: () => void;
}

function renderSlideVisual(
	slide: ContextSlide,
	snapshot: ContextWindowSnapshot,
	styleSet: ContextStyleSet,
) {
	if (slide.view === "strip")
		return <ContextStripDiagram snapshot={snapshot} styleSet={styleSet} />;
	if (slide.view === "stack")
		return <ContextStackDiagram snapshot={snapshot} styleSet={styleSet} />;
	if (slide.view === "budget") return <ContextBudgetBar snapshot={snapshot} styleSet={styleSet} />;
	return <ContextTreemap snapshot={snapshot} styleSet={styleSet} />;
}

function legendItemsForSnapshot(snapshot: ContextWindowSnapshot): ContextLegendItemSpec[] {
	return snapshot.parts
		.filter((part) => part.tokens > 0)
		.map((part, order) => ({
			id: part.id,
			label: part.label,
			styleKey: part.styleKey,
			order,
		}));
}

export function CourseSlidePanel({
	slide,
	snapshot,
	styleSet = contextDefaultStyleSet,
	index,
	total,
	visualSide = "left",
	mode = "default",
	onPrevious,
	onNext,
	onPresent,
	onFullscreen,
	className,
	...rest
}: CourseSlidePanelProps) {
	const counter =
		index != null && total != null
			? `${String(index + 1).padStart(2, "0")} / ${String(total).padStart(2, "0")}`
			: undefined;
	const legend =
		slide.view === "budget" ? undefined : (
			<ContextLegend
				items={legendItemsForSnapshot(snapshot)}
				styles={styleSet.styles}
				size="sm"
				selectedId={snapshot.selectedPartId}
			/>
		);
	const visual = (
		<FigureBlock
			caption={`${snapshot.title} — ${snapshot.limit.toLocaleString()} token window`}
			legend={legend}
		>
			{renderSlideVisual(slide, snapshot, styleSet)}
		</FigureBlock>
	);
	const notes = <KeyPointList items={slide.notes} />;
	const footer =
		onPrevious || onNext || onPresent || onFullscreen ? (
			<Inline justify="between">
				<Inline gap="sm">
					<Button size="compact" onClick={onPrevious} disabled={!onPrevious}>
						◂ Prev
					</Button>
					<Button size="compact" onClick={onNext} disabled={!onNext}>
						Next ▸
					</Button>
				</Inline>
				<Inline gap="sm" justify="end">
					{onPresent && (
						<Button size="compact" onClick={onPresent}>
							Present
						</Button>
					)}
					{onFullscreen && (
						<Button size="compact" onClick={onFullscreen}>
							Fullscreen
						</Button>
					)}
				</Inline>
			</Inline>
		) : undefined;

	return (
		<SlideShell
			mode={mode}
			eyebrow={slide.kicker}
			counter={counter}
			title={slide.title}
			primary={visual}
			secondary={notes}
			primarySide={visualSide}
			ratio="primaryWide"
			footer={footer}
			className={className}
			data-rag-organism="CourseSlidePanel"
			{...rest}
		/>
	);
}
