import type { HTMLAttributes } from 'react';
import type { ContextSlide, ContextWindowSnapshot } from '../../../context';
import { Button } from '../../atoms';
import { Inline, SlideShell } from '../../layout';
import { ContextBudgetBar, ContextLegend, ContextStackDiagram, ContextStripDiagram, ContextTreemap, FigureBlock, KeyPointList } from '../../molecules';

export interface CourseSlidePanelProps extends Omit<HTMLAttributes<HTMLElement>, 'title'> {
  slide: ContextSlide;
  snapshot: ContextWindowSnapshot;
  index?: number;
  total?: number;
  visualSide?: 'left' | 'right';
  onPrevious?: () => void;
  onNext?: () => void;
}

function renderSlideVisual(slide: ContextSlide, snapshot: ContextWindowSnapshot) {
  if (slide.view === 'strip') {
    return <ContextStripDiagram snapshot={snapshot} />;
  }
  if (slide.view === 'stack') {
    return <ContextStackDiagram snapshot={snapshot} />;
  }
  if (slide.view === 'budget') {
    return <ContextBudgetBar snapshot={snapshot} />;
  }
  return <ContextTreemap snapshot={snapshot} />;
}

export function CourseSlidePanel({
  slide,
  snapshot,
  index,
  total,
  visualSide = 'left',
  onPrevious,
  onNext,
  className,
  ...rest
}: CourseSlidePanelProps) {
  const counter = index != null && total != null ? `${String(index + 1).padStart(2, '0')} / ${String(total).padStart(2, '0')}` : undefined;
  const visual = (
    <FigureBlock caption={`${snapshot.title} — ${snapshot.limit.toLocaleString()} token window`} legend={<ContextLegend compact />}>
      {renderSlideVisual(slide, snapshot)}
    </FigureBlock>
  );
  const notes = <KeyPointList items={slide.notes} />;
  const footer = (onPrevious || onNext) ? (
    <Inline justify="between">
      <Button size="compact" onClick={onPrevious} disabled={!onPrevious}>◂ Prev</Button>
      <Button size="compact" onClick={onNext} disabled={!onNext}>Next ▸</Button>
    </Inline>
  ) : undefined;

  return (
    <SlideShell
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
