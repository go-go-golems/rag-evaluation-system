import type { HTMLAttributes } from 'react';
import type { ContextCourse } from '../../../context';
import { Button } from '../../atoms';
import { Caption } from '../../foundation';
import { SectionBlock, SplitPane } from '../../layout';
import { CheckList, KeyValueStrip, PersonSummary, StepList } from '../../molecules';
import styles from './CourseLessonPanel.module.css';

export interface CourseLessonPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  course: ContextCourse;
  activeAgendaItemId?: string;
  onAgendaItemSelect?: (itemId: string) => void;
}

export function CourseLessonPanel({ course, activeAgendaItemId, onAgendaItemSelect, className, ...rest }: CourseLessonPanelProps) {
  const agendaItems = course.agenda.map((item) => ({
    id: item.id,
    index: item.number,
    title: item.title,
    description: item.description,
    meta: item.duration,
  }));

  return (
    <article className={[styles.root, className ?? ''].filter(Boolean).join(' ')} data-rag-organism="CourseLessonPanel" {...rest}>
      <SectionBlock density="spacious" divider="bottom" className={styles.hero}>
        {course.kicker && <div className={styles.eyebrow}>{course.kicker}</div>}
        <h1 className={styles.title}>{course.title}</h1>
        <p className={styles.lead}>{course.tagline}</p>
        <div className={styles.ctaRow}>
          <Button variant="primary">Reserve a seat</Button>
          <Button>Preview the deck</Button>
          {course.price && <Caption>{course.price}</Caption>}
        </div>
        <KeyValueStrip className={styles.metaStrip} items={[
          { key: 'When', value: course.when ?? 'TBD' },
          { key: 'Where', value: course.where ?? 'TBD' },
          { key: 'Format', value: course.format ?? 'TBD' },
        ]} />
      </SectionBlock>

      <SplitPane
        ratio="course"
        divider
        left={(
          <SectionBlock label="What you'll leave with" density="spacious">
            <CheckList items={course.outcomes} />
            {course.instructor && (
              <div className={styles.instructorBlock}>
                <div className={styles.sectionLabel}>Instructor</div>
                <PersonSummary name={course.instructor.name} subtitle={course.instructor.role} bio={course.instructor.bio} />
              </div>
            )}
          </SectionBlock>
        )}
        right={(
          <SectionBlock label="Agenda · 2 hours" density="spacious">
            <StepList items={agendaItems} activeItemId={activeAgendaItemId} onItemSelect={onAgendaItemSelect} />
          </SectionBlock>
        )}
      />

      {course.blurb && (
        <SectionBlock density="spacious" divider="top">
          <p className={styles.blurb}>{course.blurb}</p>
        </SectionBlock>
      )}
    </article>
  );
}
