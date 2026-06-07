import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { contextCourseFixture } from '../../../context';
import { CourseLessonPanel } from './CourseLessonPanel';
const meta = { title: 'Component Library/Organisms/CourseLessonPanel', component: CourseLessonPanel, args: { course: contextCourseFixture } } satisfies Meta<typeof CourseLessonPanel>;
export default meta; type Story = StoryObj<typeof meta>;
export const WorkshopLanding: Story = { render: () => { const [active, setActive] = useState(contextCourseFixture.agenda[0]?.id); return <CourseLessonPanel course={contextCourseFixture} activeAgendaItemId={active} onAgendaItemSelect={setActive} />; } };
