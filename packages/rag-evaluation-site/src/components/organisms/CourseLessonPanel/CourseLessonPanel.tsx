import type { HTMLAttributes } from 'react';
import type { ContextCourse } from '../../../context';
import { Caption, Text } from '../../foundation';
import { DashboardGrid, Panel, Stack } from '../../layout';
import { MetadataGrid } from '../../molecules';
import { CourseStepNav } from '../../molecules/CourseStepNav';

export interface CourseLessonPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, 'title'> {
  course: ContextCourse;
  activeAgendaItemId?: string;
  onAgendaItemSelect?: (itemId: string) => void;
}

export function CourseLessonPanel({ course, activeAgendaItemId, onAgendaItemSelect, ...rest }: CourseLessonPanelProps) {
  return (
    <Panel title={course.title} data-rag-organism="CourseLessonPanel" {...rest}>
      <DashboardGrid recipe="twoColumn">
        <Stack gap="md">
          {course.kicker && <Caption transform="uppercase">{course.kicker}</Caption>}
          <Text size="metric" weight="bold">{course.tagline}</Text>
          {course.blurb && <Text>{course.blurb}</Text>}
          <MetadataGrid items={[
            { key: 'When', value: course.when ?? 'TBD' },
            { key: 'Where', value: course.where ?? 'TBD' },
            { key: 'Format', value: course.format ?? 'TBD' },
            { key: 'Price', value: course.price ?? 'TBD' },
          ]} />
          {course.instructor && <Panel title="Instructor" density="condensed"><Text weight="bold">{course.instructor.name}</Text><Caption>{course.instructor.role}</Caption><Text>{course.instructor.bio}</Text></Panel>}
        </Stack>
        <Stack gap="md">
          <Panel title="Outcomes" density="condensed"><Stack gap="xs">{course.outcomes.map((outcome) => <Text key={outcome}>▸ {outcome}</Text>)}</Stack></Panel>
          <Panel title="Agenda" density="condensed"><CourseStepNav items={course.agenda} activeItemId={activeAgendaItemId} onItemSelect={onAgendaItemSelect} /></Panel>
        </Stack>
      </DashboardGrid>
    </Panel>
  );
}
