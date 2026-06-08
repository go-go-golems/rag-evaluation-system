import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { contextSlides, contextWindowSnapshots } from '../../../context';
import { Caption } from '../../foundation';
import { SlideShell } from '../../layout';
import { ContextStackDiagram, FigureBlock, KeyPointList } from '../../molecules';
import { CourseStudioShell } from './CourseStudioShell';
import { courseStudioNavSections } from './courseStudioNav';

const slide = contextSlides[1]!;
const snapshot = contextWindowSnapshots.find((item) => item.id === slide.snapshotId)!;

const slideContent = (
  <SlideShell
    eyebrow={slide.kicker}
    counter="02 / 06"
    title={slide.title}
    primary={(
      <FigureBlock caption="context window — 200,000 tokens">
        <ContextStackDiagram snapshot={snapshot} />
      </FigureBlock>
    )}
    secondary={<KeyPointList items={slide.notes} />}
    primarySide="right"
  />
);

const meta = {
  title: 'Component Library/Organisms/CourseStudioShell',
  component: CourseStudioShell,
  args: {
    sections: courseStudioNavSections,
    activeItemId: 'slides',
    children: slideContent,
    sidebarFooter: <Caption>Workshop draft · v0.4</Caption>,
  },
} satisfies Meta<typeof CourseStudioShell>;
export default meta;
type Story = StoryObj<typeof meta>;

export const SlidesActive: Story = {};

export const Interactive: Story = {
  render: () => {
    const [active, setActive] = useState('slides');
    return (
      <CourseStudioShell sections={courseStudioNavSections} activeItemId={active} onNavigate={setActive} sidebarFooter={<Caption>Selected: {active}</Caption>} >
        {slideContent}
      </CourseStudioShell>
    );
  },
};
