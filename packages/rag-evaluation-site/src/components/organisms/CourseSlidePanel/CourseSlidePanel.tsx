import type { HTMLAttributes } from 'react';
import type { ContextSlide, ContextWindowSnapshot } from '../../../context';
import { Button } from '../../atoms';
import { Caption, Text } from '../../foundation';
import { Inline, Panel, Stack } from '../../layout';
import { ContextBudgetBar, ContextStackDiagram, ContextStripDiagram, ContextTreemap } from '../../molecules';
import styles from './CourseSlidePanel.module.css';

export interface CourseSlidePanelProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  slide: ContextSlide;
  snapshot: ContextWindowSnapshot;
  index?: number;
  total?: number;
  onPrevious?: () => void;
  onNext?: () => void;
}

export function CourseSlidePanel({ slide, snapshot, index, total, onPrevious, onNext, className, ...rest }: CourseSlidePanelProps) {
  return (
    <Panel title={slide.title} className={className} data-rag-organism="CourseSlidePanel" {...rest}>
      <div className={styles.root}>
        <header className={styles.header}>
          <div><Caption transform="uppercase">{slide.kicker}</Caption><Text size="metric" weight="bold">{slide.title}</Text></div>
          {index != null && total != null && <Caption>{String(index + 1).padStart(2, '0')} / {String(total).padStart(2, '0')}</Caption>}
        </header>
        <div className={styles.body}>
          <div className={styles.diagram}>
            {slide.view === 'strip' && <ContextStripDiagram snapshot={snapshot} />}
            {slide.view === 'stack' && <ContextStackDiagram snapshot={snapshot} />}
            {slide.view === 'budget' && <ContextBudgetBar snapshot={snapshot} />}
            {slide.view === 'treemap' && <ContextTreemap snapshot={snapshot} />}
          </div>
          <Stack gap="md">
            {slide.notes.map((note, i) => <div key={note} className={styles.note}><span>{String(i + 1).padStart(2, '0')}</span><Text>{note}</Text></div>)}
          </Stack>
        </div>
        {(onPrevious || onNext) && <Inline justify="between"><Button size="compact" onClick={onPrevious} disabled={!onPrevious}>◂ Prev</Button><Button size="compact" onClick={onNext} disabled={!onNext}>Next ▸</Button></Inline>}
      </div>
    </Panel>
  );
}
