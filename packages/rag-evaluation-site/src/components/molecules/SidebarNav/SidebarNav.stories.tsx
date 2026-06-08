import type { Meta, StoryObj } from '@storybook/react-vite';
import { useState } from 'react';
import { courseStudioNavSections } from '../../organisms/CourseStudioShell/courseStudioNav';
import { SidebarNav } from './SidebarNav';

const meta = { title: 'Component Library/Molecules/SidebarNav', component: SidebarNav, args: { sections: courseStudioNavSections, activeItemId: 'slides' } } satisfies Meta<typeof SidebarNav>;
export default meta;
type Story = StoryObj<typeof meta>;

export const CourseNavigation: Story = {};

export const Interactive: Story = {
  render: () => {
    const [active, setActive] = useState('slides');
    return <SidebarNav sections={courseStudioNavSections} activeItemId={active} onItemSelect={setActive} />;
  },
};
