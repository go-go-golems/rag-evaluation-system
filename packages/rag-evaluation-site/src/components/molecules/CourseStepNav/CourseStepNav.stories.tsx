import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { contextCourseFixture } from '../../../context';
import { CourseStepNav } from './CourseStepNav';
const meta = { title: 'Component Library/Molecules/CourseStepNav', component: CourseStepNav, args: { items: contextCourseFixture.agenda } } satisfies Meta<typeof CourseStepNav>;
export default meta; type Story = StoryObj<typeof meta>;
export const Agenda: Story = { render: () => { const [active, setActive] = useState(contextCourseFixture.agenda[0]?.id); return <CourseStepNav items={contextCourseFixture.agenda} activeItemId={active} onItemSelect={setActive} />; } };
