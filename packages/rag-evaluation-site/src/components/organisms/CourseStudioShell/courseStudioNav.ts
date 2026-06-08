import type { SidebarNavSection } from '../../molecules';

export const courseStudioNavSections: SidebarNavSection[] = [
  { id: 'present', label: 'Present', items: [{ id: 'course', label: 'Course', icon: '◰' }, { id: 'slides', label: 'Slides', icon: '▣' }] },
  { id: 'analyze', label: 'Analyze', items: [{ id: 'visualize', label: 'Visualize', icon: '◫' }, { id: 'upload', label: 'Upload', icon: '↥' }] },
  { id: 'review', label: 'Review', items: [{ id: 'transcript', label: 'Transcript', icon: '☰' }, { id: 'comments', label: 'Comments', icon: '◇' }] },
  { id: 'take-home', label: 'Take-home', items: [{ id: 'handout', label: 'Handout', icon: '□' }] },
];
